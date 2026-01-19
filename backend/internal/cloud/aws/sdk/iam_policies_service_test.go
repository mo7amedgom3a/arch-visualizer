package sdk

import (
	"context"
	"testing"
)

func TestPolicyService_LoadFromStatic(t *testing.T) {
	// Create a service with a nil client (will use static data only)
	service, err := NewPolicyService(nil)
	if err != nil {
		t.Fatalf("Failed to create PolicyService: %v", err)
	}

	ctx := context.Background()

	// Test getting a known policy from static data
	policy, err := service.GetPolicy(ctx, "arn:aws:iam::aws:policy/ReadOnlyAccess")
	if err != nil {
		t.Fatalf("Failed to get ReadOnlyAccess policy: %v", err)
	}

	if policy == nil {
		t.Fatal("Policy is nil")
	}

	if policy.Name != "ReadOnlyAccess" {
		t.Errorf("Expected name ReadOnlyAccess, got %s", policy.Name)
	}

	if !policy.IsAWSManaged {
		t.Error("Expected IsAWSManaged to be true")
	}

	if policy.ARN != "arn:aws:iam::aws:policy/ReadOnlyAccess" {
		t.Errorf("Expected ARN arn:aws:iam::aws:policy/ReadOnlyAccess, got %s", policy.ARN)
	}
}

func TestPolicyService_ListAllStaticPolicies(t *testing.T) {
	service, err := NewPolicyService(nil)
	if err != nil {
		t.Fatalf("Failed to create PolicyService: %v", err)
	}

	ctx := context.Background()

	policies, err := service.ListAllStaticPolicies(ctx)
	if err != nil {
		t.Fatalf("Failed to list policies: %v", err)
	}

	if len(policies) == 0 {
		t.Error("Expected at least one policy from static data")
	}

	// Verify all policies are AWS managed
	for _, policy := range policies {
		if !policy.IsAWSManaged {
			t.Errorf("Policy %s should be marked as AWS managed", policy.ARN)
		}
	}
}

func TestPolicyService_ListPoliciesByResource(t *testing.T) {
	service, err := NewPolicyService(nil)
	if err != nil {
		t.Fatalf("Failed to create PolicyService: %v", err)
	}

	ctx := context.Background()

	// Test listing EC2 policies
	ec2Policies, err := service.ListPoliciesByResource(ctx, "ec2")
	if err != nil {
		t.Fatalf("Failed to list EC2 policies: %v", err)
	}

	// Should have at least AmazonEC2FullAccess and AmazonEC2ReadOnlyAccess
	if len(ec2Policies) < 2 {
		t.Errorf("Expected at least 2 EC2 policies, got %d", len(ec2Policies))
	}

	// Verify policies are EC2-related
	foundEC2Full := false
	foundEC2ReadOnly := false
	for _, policy := range ec2Policies {
		if policy.Name == "AmazonEC2FullAccess" {
			foundEC2Full = true
		}
		if policy.Name == "AmazonEC2ReadOnlyAccess" {
			foundEC2ReadOnly = true
		}
	}

	if !foundEC2Full {
		t.Error("Expected to find AmazonEC2FullAccess")
	}
	if !foundEC2ReadOnly {
		t.Error("Expected to find AmazonEC2ReadOnlyAccess")
	}
}

func TestPolicyService_HasPolicy(t *testing.T) {
	service, err := NewPolicyService(nil)
	if err != nil {
		t.Fatalf("Failed to create PolicyService: %v", err)
	}

	// Test with known policy
	if !service.HasPolicy("arn:aws:iam::aws:policy/ReadOnlyAccess") {
		t.Error("Expected HasPolicy to return true for ReadOnlyAccess")
	}

	// Test with unknown policy
	if service.HasPolicy("arn:aws:iam::aws:policy/NonExistentPolicy") {
		t.Error("Expected HasPolicy to return false for non-existent policy")
	}
}
