package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
)

// ServicePolicyMapping maps service names to policy name patterns
var ServicePolicyMapping = map[string][]string{
	"ec2": {
		"AmazonEC2",
		"EC2",
	},
	"vpc": {
		"AmazonVPC",
		"VPC",
	},
	"ebs": {
		"AmazonElasticBlockStore",
		"EBS",
	},
	"s3": {
		"AmazonS3",
		"S3",
		"SimpleStorageService",
	},
	"iam": {
		"IAM",
		"Identity",
	},
	"general": {
		"ReadOnlyAccess",
		"PowerUserAccess",
		"AdministratorAccess",
		"ViewOnlyAccess",
	},
}

// FetchAndSaveServicePolicies fetches AWS managed policies from SDK and saves them to service-specific JSON files
func FetchAndSaveServicePolicies(ctx context.Context, client *AWSClient) error {
	if client == nil || client.IAM == nil {
		return fmt.Errorf("AWS client not available")
	}

	// List all AWS managed policies (without documents first for speed)
	policies, err := ListPolicies(ctx, client, nil, types.PolicyScopeTypeAws, false)
	if err != nil {
		return fmt.Errorf("failed to list AWS managed policies: %w", err)
	}

	// Fetch policy documents for each policy (this may take a while)
	fmt.Printf("Fetching policy documents for %d policies...\n", len(policies))
	for i, policy := range policies {
		if policy.PolicyDocument == "" {
			fullPolicy, err := GetPolicy(ctx, client, policy.ARN)
			if err == nil && fullPolicy != nil {
				policies[i].PolicyDocument = fullPolicy.PolicyDocument
			}
		}
		if (i+1)%10 == 0 {
			fmt.Printf("  Processed %d/%d policies...\n", i+1, len(policies))
		}
	}

	// Organize policies by service
	servicePolicies := make(map[string][]*awsoutputs.PolicyOutput)
	for service := range ServicePolicyMapping {
		servicePolicies[service] = make([]*awsoutputs.PolicyOutput, 0)
	}

	// Categorize policies
	for _, policy := range policies {
		service := categorizePolicy(policy.Name)
		if service != "" {
			servicePolicies[service] = append(servicePolicies[service], policy)
		}
	}

	// Save each service's policies to its own JSON file
	for service, policies := range servicePolicies {
		if len(policies) == 0 {
			continue // Skip empty services
		}

		if err := saveServicePolicies(service, policies); err != nil {
			fmt.Printf("Warning: Failed to save %s policies: %v\n", service, err)
			continue
		}

		fmt.Printf("Saved %d policies for service: %s\n", len(policies), service)
	}

	return nil
}

// categorizePolicy determines which service a policy belongs to based on its name
func categorizePolicy(policyName string) string {
	policyNameLower := strings.ToLower(policyName)

	// Check each service's patterns
	for service, patterns := range ServicePolicyMapping {
		for _, pattern := range patterns {
			if strings.Contains(policyNameLower, strings.ToLower(pattern)) {
				return service
			}
		}
	}

	// Default to general if no specific service matches
	return "general"
}

// saveServicePolicies saves policies to a service-specific JSON file
func saveServicePolicies(service string, policies []*awsoutputs.PolicyOutput) error {
	// Convert to static policy entries
	entries := make([]StaticPolicyEntry, 0, len(policies))
	for _, policy := range policies {
		// Decode policy document if it's URL-encoded
		policyDoc := policy.PolicyDocument
		if decoded, err := DecodePolicyDocument(policyDoc); err == nil {
			policyDoc = decoded
		}

		entry := StaticPolicyEntry{
			ARN:                policy.ARN,
			Name:               policy.Name,
			Description:        policy.Description,
			Path:               policy.Path,
			PolicyDocument:     policyDoc, // Use decoded document
			IsAWSManaged:       policy.IsAWSManaged,
			ResourceCategories: []string{service},
			RelatedResources:   []string{service},
		}
		entries = append(entries, entry)
	}

	// Get the data directory path
	dataDir := getDataDirectory()
	serviceDir := filepath.Join(dataDir, service)

	// Create service directory if it doesn't exist
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// Write to JSON file
	jsonPath := filepath.Join(serviceDir, "policies.json")
	jsonData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal policies: %w", err)
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// getDataDirectory returns the path to the IAM data directory
func getDataDirectory() string {
	wd, _ := os.Getwd()
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)

	possiblePaths := []string{
		filepath.Join(currentDir, "..", "models", "iam", "data"),
		filepath.Join(wd, "internal", "cloud", "aws", "models", "iam", "data"),
		filepath.Join(wd, "backend", "internal", "cloud", "aws", "models", "iam", "data"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Default to the first path and create it if needed
	defaultPath := possiblePaths[0]
	os.MkdirAll(defaultPath, 0755)
	return defaultPath
}

// LoadServicePolicies loads policies from a service-specific JSON file
func LoadServicePolicies(service string) ([]StaticPolicyEntry, error) {
	dataDir := getDataDirectory()
	jsonPath := filepath.Join(dataDir, service, "policies.json")

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read service policies file: %w", err)
	}

	var entries []StaticPolicyEntry
	if err := json.Unmarshal(jsonData, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse service policies JSON: %w", err)
	}

	return entries, nil
}

// ListAvailableServices returns a list of services that have policy files
func ListAvailableServices() ([]string, error) {
	dataDir := getDataDirectory()

	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	var services []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Check if this directory has a policies.json file
			policiesPath := filepath.Join(dataDir, entry.Name(), "policies.json")
			if _, err := os.Stat(policiesPath); err == nil {
				services = append(services, entry.Name())
			}
		}
	}

	return services, nil
}

// FetchAndSaveS3Policies fetches S3-related AWS managed policies from SDK and saves them to s3/polices.json
func FetchAndSaveS3Policies(ctx context.Context, client *AWSClient) error {
	if client == nil || client.IAM == nil {
		return fmt.Errorf("AWS client not available")
	}

	// List all AWS managed policies (without documents first for speed)
	policies, err := ListPolicies(ctx, client, nil, types.PolicyScopeTypeAws, false)
	if err != nil {
		return fmt.Errorf("failed to list AWS managed policies: %w", err)
	}

	// Filter S3-related policies
	var s3Policies []*awsoutputs.PolicyOutput
	for _, policy := range policies {
		service := categorizePolicy(policy.Name)
		if service == "s3" {
			s3Policies = append(s3Policies, policy)
		}
	}

	if len(s3Policies) == 0 {
		return fmt.Errorf("no S3 policies found")
	}

	// Fetch policy documents for each S3 policy
	fmt.Printf("Fetching policy documents for %d S3 policies...\n", len(s3Policies))
	for i, policy := range s3Policies {
		if policy.PolicyDocument == "" {
			fullPolicy, err := GetPolicy(ctx, client, policy.ARN)
			if err == nil && fullPolicy != nil {
				s3Policies[i].PolicyDocument = fullPolicy.PolicyDocument
			}
		}
		if (i+1)%10 == 0 {
			fmt.Printf("  Processed %d/%d policies...\n", i+1, len(s3Policies))
		}
	}

	// Save S3 policies to polices.json (note: filename has typo as per existing file)
	return saveS3Policies(s3Policies)
}

// saveS3Policies saves S3 policies to s3/polices.json (note: filename has typo)
func saveS3Policies(policies []*awsoutputs.PolicyOutput) error {
	// Convert to static policy entries
	entries := make([]StaticPolicyEntry, 0, len(policies))
	for _, policy := range policies {
		// Decode policy document if it's URL-encoded
		policyDoc := policy.PolicyDocument
		if decoded, err := DecodePolicyDocument(policyDoc); err == nil {
			policyDoc = decoded
		}

		entry := StaticPolicyEntry{
			ARN:                policy.ARN,
			Name:               policy.Name,
			Description:        policy.Description,
			Path:               policy.Path,
			PolicyDocument:     policyDoc, // Use decoded document
			IsAWSManaged:       policy.IsAWSManaged,
			ResourceCategories: []string{"s3"},
			RelatedResources:   []string{"s3"},
		}
		entries = append(entries, entry)
	}

	// Get the data directory path
	dataDir := getDataDirectory()
	serviceDir := filepath.Join(dataDir, "s3")

	// Create service directory if it doesn't exist
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// Write to JSON file (note: using "polices.json" as per existing file name)
	jsonPath := filepath.Join(serviceDir, "polices.json")
	jsonData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal policies: %w", err)
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	fmt.Printf("Saved %d S3 policies to %s\n", len(entries), jsonPath)
	return nil
}
