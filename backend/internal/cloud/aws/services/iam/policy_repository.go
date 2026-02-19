package iam

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services"
)

// PolicyRepository manages access to static AWS managed policies
type PolicyRepository struct {
	policies []*PolicyDefinition
	mu       sync.RWMutex
	basePath string
}

// PolicyDefinition matches the JSON structure in data/policies.json
type PolicyDefinition struct {
	ARN                string   `json:"arn"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	Path               string   `json:"path"`
	PolicyDocument     string   `json:"policy_document"`
	IsAWSManaged       bool     `json:"is_aws_managed"`
	ResourceCategories []string `json:"resource_categories"`
	RelatedResources   []string `json:"related_resources"`
}

// NewPolicyRepository creates a new repository instance
func NewPolicyRepository() *PolicyRepository {
	return &PolicyRepository{
		policies: make([]*PolicyDefinition, 0),
	}
}

// LoadPolicies reads policies from the JSON file system
func (r *PolicyRepository) LoadPolicies() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Determine the path to the data directory
	// logic to find backend/internal/cloud/aws/models/iam/data
	// This assumes running from backend root or similar, but we try to be robust

	// Default path relative to where this code might run (dev environment)
	// We can try multiple paths or use an absolute path provided by config
	// For now, let's try to locate it relative to this file

	baseDir, err := resolveDataDir()
	if err != nil {
		return fmt.Errorf("resolve data dir: %w", err)
	}

	r.basePath = baseDir
	r.policies = make([]*PolicyDefinition, 0)

	// Walk through the directory to find all policies.json files
	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "policies.json" {
			if err := r.loadPoliciesFromFile(path); err != nil {
				return fmt.Errorf("load policies from %s: %w", path, err)
			}
		}
		return nil
	})

	return err
}

func (r *PolicyRepository) loadPoliciesFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var filePolicies []*PolicyDefinition
	if err := json.Unmarshal(data, &filePolicies); err != nil {
		return fmt.Errorf("unmarshal policies: %w", err)
	}

	r.policies = append(r.policies, filePolicies...)

	// Infer service from directory name if applicable
	dir := filepath.Base(filepath.Dir(path))
	// If the parent dir is not the data root (we can't easily check against r.basePath here without absolute paths,
	// but generally if it's not "data", it's likely a service dir)
	if dir != "data" && dir != "." {
		for _, p := range filePolicies {
			// Check if dir name is already in RelatedResources
			found := false
			for _, rr := range p.RelatedResources {
				if strings.EqualFold(rr, dir) {
					found = true
					break
				}
			}
			if !found {
				p.RelatedResources = append(p.RelatedResources, dir)
			}
		}
	}

	return nil
}

// ListPolicies returns policies, optionally filtering by service (related_resource)
func (r *PolicyRepository) ListPolicies(service string) []*awsoutputs.PolicyOutput {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*awsoutputs.PolicyOutput, 0)
	service = strings.ToLower(service)

	for _, def := range r.policies {
		if service != "" {
			// Check if this policy relates to the requested service
			match := false
			for _, rr := range def.RelatedResources {
				if strings.ToLower(rr) == service {
					match = true
					break
				}
			}
			// Also check categories if needed, but 'related_resources' is more specific
			if !match {
				continue
			}
		}

		result = append(result, r.toOutput(def))
	}

	return result
}

func (r *PolicyRepository) toOutput(def *PolicyDefinition) *awsoutputs.PolicyOutput {
	return &awsoutputs.PolicyOutput{
		ARN:              def.ARN,
		ID:               def.ARN,
		Name:             def.Name,
		Description:      services.StringPtr(def.Description),
		Path:             def.Path,
		PolicyDocument:   def.PolicyDocument,
		CreateDate:       services.GetFixedTimestamp(),
		UpdateDate:       services.GetFixedTimestamp(),
		DefaultVersionID: services.StringPtr("v1"),
		AttachmentCount:  0,
		IsAttachable:     true,
		Tags:             []configs.Tag{},
		IsAWSManaged:     def.IsAWSManaged,
	}
}

// resolveDataDir attempts to find the models/iam/data directory
func resolveDataDir() (string, error) {
	// Try finding it relative to the caller location during development
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get caller info")
	}

	// filename is .../backend/internal/cloud/aws/services/iam/policy_repository.go
	// goal is .../backend/internal/cloud/aws/models/iam/data

	currDir := filepath.Dir(filename)
	// Up 3 levels: internal/cloud/aws
	awsDir := filepath.Clean(filepath.Join(currDir, "..", ".."))

	dataDir := filepath.Join(awsDir, "models", "iam", "data")

	if info, err := os.Stat(dataDir); err == nil && info.IsDir() {
		return dataDir, nil
	}

	return "", fmt.Errorf("data directory not found at expected location: %s", dataDir)
}

// GetPolicy returns a single policy by ARN
func (r *PolicyRepository) GetPolicy(arn string) *awsoutputs.PolicyOutput {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, def := range r.policies {
		if def.ARN == arn {
			return r.toOutput(def)
		}
	}
	return nil
}

// Count returns the total number of loaded policies
func (r *PolicyRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.policies)
}

// ListPoliciesByService returns policies filtered by source and destination service
func (r *PolicyRepository) ListPoliciesByService(sourceService, destinationService string) []*awsoutputs.PolicyOutput {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*awsoutputs.PolicyOutput, 0)
	sourceService = strings.ToLower(sourceService)
	destinationService = strings.ToLower(destinationService)

	for _, def := range r.policies {
		// Heuristic 1: Policy name should contain source service (e.g. "Lambda" in "AWSLambdaExecute")
		if sourceService != "" && !strings.Contains(strings.ToLower(def.Name), sourceService) {
			continue
		}

		// Heuristic 2: Policy document should contain destination service (e.g. "s3:" in Action or Resource)
		if destinationService != "" && !strings.Contains(strings.ToLower(def.PolicyDocument), destinationService) {
			continue
		}

		result = append(result, r.toOutput(def))
	}

	return result
}
