package sdk

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
)

// TestCreatePolicy tests policy creation (integration test - requires AWS credentials)
func TestCreatePolicy(t *testing.T) {
	// Skip if AWS credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping integration test: AWS credentials not available")
	}

	ctx := context.Background()
	client, err := NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	// Create a test policy
	policy := &awsiam.Policy{
		Name:           "test-policy-sdk",
		PolicyDocument: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`,
		Path:           stringPtr("/test/"),
	}

	output, err := CreatePolicy(ctx, client, policy)
	if err != nil {
		t.Fatalf("Failed to create policy: %v", err)
	}

	if output == nil {
		t.Fatal("Policy output is nil")
	}

	if output.Name != policy.Name {
		t.Errorf("Expected policy name %s, got %s", policy.Name, output.Name)
	}

	if output.ARN == "" {
		t.Error("Policy ARN is empty")
	}

	// Cleanup: delete the policy
	if err := DeletePolicy(ctx, client, output.ARN); err != nil {
		t.Logf("Warning: Failed to cleanup test policy: %v", err)
	}
}

// TestGetPolicy tests policy retrieval (integration test - requires AWS credentials)
func TestGetPolicy(t *testing.T) {
	// Skip if AWS credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping integration test: AWS credentials not available")
	}

	ctx := context.Background()
	client, err := NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	// Test getting an AWS managed policy (ReadOnlyAccess should exist)
	arn := "arn:aws:iam::aws:policy/ReadOnlyAccess"
	policy, err := GetPolicy(ctx, client, arn)
	if err != nil {
		t.Fatalf("Failed to get policy: %v", err)
	}

	if policy == nil {
		t.Fatal("Policy is nil")
	}

	if policy.ARN != arn {
		t.Errorf("Expected ARN %s, got %s", arn, policy.ARN)
	}

	if !policy.IsAWSManaged {
		t.Error("Expected IsAWSManaged to be true for AWS managed policy")
	}
}

// TestListAWSManagedPolicies tests listing AWS managed policies (integration test)
func TestListAWSManagedPolicies(t *testing.T) {
	// Skip if AWS credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping integration test: AWS credentials not available")
	}

	ctx := context.Background()
	client, err := NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	policies, err := ListAWSManagedPolicies(ctx, client, nil)
	if err != nil {
		t.Fatalf("Failed to list AWS managed policies: %v", err)
	}

	if len(policies) == 0 {
		t.Error("Expected at least one AWS managed policy")
	}

	// Verify all policies are marked as AWS managed
	for _, policy := range policies {
		if !policy.IsAWSManaged {
			t.Errorf("Policy %s should be marked as AWS managed", policy.ARN)
		}
		if !stringHasPrefix(policy.ARN, "arn:aws:iam::aws:policy/") {
			t.Errorf("Policy ARN %s should start with arn:aws:iam::aws:policy/", policy.ARN)
		}
	}
}

// TestListPolicies tests listing customer managed policies (integration test)
func TestListPolicies(t *testing.T) {
	// Skip if AWS credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping integration test: AWS credentials not available")
	}

	ctx := context.Background()
	client, err := NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	// This will list customer managed policies
	// Note: This test may pass even if there are no customer policies
	// Don't fetch policy documents for listing (faster performance)
	policies, err := ListPolicies(ctx, client, nil, types.PolicyScopeTypeLocal, false)
	if err != nil {
		t.Fatalf("Failed to list policies: %v", err)
	}

	// Verify policies are not marked as AWS managed
	for _, policy := range policies {
		if policy.IsAWSManaged {
			t.Errorf("Customer policy %s should not be marked as AWS managed", policy.ARN)
		}
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func stringHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
