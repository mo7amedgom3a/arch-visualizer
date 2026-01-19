package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
)

// StaticPolicyEntry represents a policy entry in the static JSON file
type StaticPolicyEntry struct {
	ARN              string   `json:"arn"`
	Name             string   `json:"name"`
	Description      *string  `json:"description,omitempty"`
	Path             string   `json:"path"`
	PolicyDocument   string   `json:"policy_document"`
	IsAWSManaged     bool     `json:"is_aws_managed"`
	ResourceCategories []string `json:"resource_categories"`
	RelatedResources  []string `json:"related_resources"`
}

// PolicyService provides IAM policy information with caching and fallback to static data
type PolicyService struct {
	client     *AWSClient
	cache      map[string]*awsoutputs.PolicyOutput
	cacheMutex sync.RWMutex
	staticData map[string]*awsoutputs.PolicyOutput
}

// NewPolicyService creates a new policy service
func NewPolicyService(client *AWSClient) (*PolicyService, error) {
	service := &PolicyService{
		client:     client,
		cache:      make(map[string]*awsoutputs.PolicyOutput),
		staticData: make(map[string]*awsoutputs.PolicyOutput),
	}

	// Load static data
	if err := service.loadFromStatic(); err != nil {
		// Log error but don't fail - we can still use AWS API
		fmt.Printf("Warning: Failed to load static IAM policy data: %v\n", err)
	}

	return service, nil
}

// GetPolicy retrieves a policy by ARN
// Checks cache first, then AWS API, then static data
func (s *PolicyService) GetPolicy(ctx context.Context, arn string) (*awsoutputs.PolicyOutput, error) {
	// Normalize ARN (lowercase for consistency)
	arn = strings.ToLower(arn)

	// Check cache first
	s.cacheMutex.RLock()
	if cached, ok := s.cache[arn]; ok {
		s.cacheMutex.RUnlock()
		return cached, nil
	}
	s.cacheMutex.RUnlock()

	// Try AWS API if client is available
	if s.client != nil && s.client.IAM != nil {
		// Note: We don't call SDK here directly to avoid circular dependency
		// The SDK functions will call this service if needed
		// For now, we'll fall back to static data
	}

	// Fall back to static data
	s.cacheMutex.RLock()
	if static, ok := s.staticData[arn]; ok {
		s.cacheMutex.RUnlock()
		// Cache it for future use
		s.cacheMutex.Lock()
		s.cache[arn] = static
		s.cacheMutex.Unlock()
		return static, nil
	}
	s.cacheMutex.RUnlock()

	return nil, fmt.Errorf("policy not found: %s", arn)
}

// ListPoliciesByResource lists policies filtered by resource type
// resourceType can be: "ec2", "ebs", "vpc", "iam", "general", or "" for all
func (s *PolicyService) ListPoliciesByResource(ctx context.Context, resourceType string) ([]*awsoutputs.PolicyOutput, error) {
	var results []*awsoutputs.PolicyOutput

	// If resourceType is empty or "all", return all policies
	if resourceType == "" || strings.EqualFold(resourceType, "all") {
		return s.ListAllStaticPolicies(ctx)
	}

	// Normalize resource type
	resourceType = strings.ToLower(resourceType)

	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	// We need to reload the static entries to check related_resources
	// For now, we'll use a simple approach: reload from JSON to check resource associations
	// In a production system, we'd store this metadata in PolicyOutput or a separate index
	wd, _ := os.Getwd()
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	possiblePaths := []string{
		filepath.Join(currentDir, "..", "models", "iam", "data", "policies.json"),
		filepath.Join(wd, "internal", "cloud", "aws", "models", "iam", "data", "policies.json"),
		filepath.Join(wd, "backend", "internal", "cloud", "aws", "models", "iam", "data", "policies.json"),
	}

	var jsonData []byte
	var err error
	for _, path := range possiblePaths {
		jsonData, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		// If we can't read the file, return all policies
		return s.ListAllStaticPolicies(ctx)
	}

	var staticPolicies []StaticPolicyEntry
	if err := json.Unmarshal(jsonData, &staticPolicies); err != nil {
		// If we can't parse, return all policies
		return s.ListAllStaticPolicies(ctx)
	}

	// Filter by resource type
	for _, entry := range staticPolicies {
		// Check if this policy is related to the requested resource type
		matches := false
		for _, relatedResource := range entry.RelatedResources {
			if strings.EqualFold(relatedResource, resourceType) {
				matches = true
				break
			}
		}

		if matches {
			arn := strings.ToLower(entry.ARN)
			if policy, ok := s.staticData[arn]; ok {
				results = append(results, policy)
			}
		}
	}

	return results, nil
}

// ListAllStaticPolicies returns all policies from static data
func (s *PolicyService) ListAllStaticPolicies(ctx context.Context) ([]*awsoutputs.PolicyOutput, error) {
	var results []*awsoutputs.PolicyOutput

	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	for _, policy := range s.staticData {
		results = append(results, policy)
	}

	return results, nil
}

// HasPolicy checks if a policy exists in static data or cache
func (s *PolicyService) HasPolicy(arn string) bool {
	arn = strings.ToLower(arn)

	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	if _, ok := s.cache[arn]; ok {
		return true
	}

	if _, ok := s.staticData[arn]; ok {
		return true
	}

	return false
}

// loadFromStatic loads policies from service-specific JSON files
func (s *PolicyService) loadFromStatic() error {
	// List available services
	services, err := ListAvailableServices()
	if err != nil {
		// If we can't list services, try to load from the old single file format
		return s.loadFromLegacyFile()
	}

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// Load policies from each service directory
	for _, service := range services {
		entries, err := LoadServicePolicies(service)
		if err != nil {
			fmt.Printf("Warning: Failed to load policies for service %s: %v\n", service, err)
			continue
		}

		for _, entry := range entries {
			arn := strings.ToLower(entry.ARN)
			policyOutput := &awsoutputs.PolicyOutput{
				ARN:            entry.ARN,
				ID:             entry.ARN,
				Name:           entry.Name,
				Description:    entry.Description,
				Path:           entry.Path,
				PolicyDocument: entry.PolicyDocument,
				IsAWSManaged:   entry.IsAWSManaged,
				CreateDate:     time.Time{}, // Static data doesn't have creation date
				UpdateDate:     time.Time{}, // Static data doesn't have update date
				AttachmentCount: 0,          // Unknown for static data
				IsAttachable:   true,        // AWS managed policies are attachable
			}
			s.staticData[arn] = policyOutput
		}
	}

	return nil
}

// loadFromLegacyFile loads from the old single policies.json file (backward compatibility)
func (s *PolicyService) loadFromLegacyFile() error {
	wd, _ := os.Getwd()
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	possiblePaths := []string{
		filepath.Join(currentDir, "..", "models", "iam", "data", "policies.json"),
		filepath.Join(wd, "internal", "cloud", "aws", "models", "iam", "data", "policies.json"),
		filepath.Join(wd, "backend", "internal", "cloud", "aws", "models", "iam", "data", "policies.json"),
	}

	var jsonData []byte
	var err error
	for _, path := range possiblePaths {
		jsonData, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("failed to find policies.json: %w", err)
	}

	var staticPolicies []StaticPolicyEntry
	if err := json.Unmarshal(jsonData, &staticPolicies); err != nil {
		return fmt.Errorf("failed to parse policies.json: %w", err)
	}

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	for _, entry := range staticPolicies {
		arn := strings.ToLower(entry.ARN)
		policyOutput := &awsoutputs.PolicyOutput{
			ARN:            entry.ARN,
			ID:             entry.ARN,
			Name:           entry.Name,
			Description:    entry.Description,
			Path:           entry.Path,
			PolicyDocument: entry.PolicyDocument,
			IsAWSManaged:   entry.IsAWSManaged,
			CreateDate:     time.Time{},
			UpdateDate:     time.Time{},
			AttachmentCount: 0,
			IsAttachable:   true,
		}
		s.staticData[arn] = policyOutput
	}

	return nil
}

// RefreshCache forces a refresh of the cache from static data
func (s *PolicyService) RefreshCache(ctx context.Context) error {
	return s.loadFromStatic()
}
