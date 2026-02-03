package iam

import (
	"testing"
)

func TestPolicyRepository_LoadPolicies(t *testing.T) {
	repo := NewPolicyRepository()
	if err := repo.LoadPolicies(); err != nil {
		t.Fatalf("Failed to load policies: %v", err)
	}

	count := repo.Count()
	if count == 0 {
		t.Error("Expected policies to be loaded, got 0")
	}
	t.Logf("Loaded %d policies", count)

	// Test filtering
	s3Policies := repo.ListPolicies("s3")
	if len(s3Policies) == 0 {
		// Note: This relies on policies.json actually having s3 policies.
		// If the static file doesn't have them, this might fail or need adjustment.
		// Based on previous file view, we saw "general", "compute" (EC2), "storage" (EBS), etc.
		// Let's check for "ec2" which we know exists from the file view.
		t.Log("No s3 policies found (might be expected if not in sample json), checking ec2")
	}

	ec2Policies := repo.ListPolicies("ec2")
	if len(ec2Policies) == 0 {
		t.Error("Expected to find EC2 related policies")
	} else {
		t.Logf("Found %d EC2 policies", len(ec2Policies))
	}

	// Test case insensitivity
	ec2PoliciesUpper := repo.ListPolicies("EC2")
	if len(ec2PoliciesUpper) != len(ec2Policies) {
		t.Errorf("Expected case insensitive filtering: got %d for EC2, %d for ec2", len(ec2PoliciesUpper), len(ec2Policies))
	}
}
