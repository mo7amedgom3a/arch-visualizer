package sdk

import (
	"context"
	"os"
	"testing"

	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
)

// TestCreateInstanceProfile tests instance profile creation (integration test - requires AWS credentials)
func TestCreateInstanceProfile(t *testing.T) {
	// Skip if AWS credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping integration test: AWS credentials not available")
	}

	ctx := context.Background()
	client, err := NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	// Create a test instance profile
	profile := &awsiam.InstanceProfile{
		Name: "test-instance-profile-sdk",
		Path: stringPtr("/test/"),
	}

	output, err := CreateInstanceProfile(ctx, client, profile)
	if err != nil {
		t.Fatalf("Failed to create instance profile: %v", err)
	}

	if output == nil {
		t.Fatal("Instance profile output is nil")
	}

	if output.Name != profile.Name {
		t.Errorf("Expected instance profile name %s, got %s", profile.Name, output.Name)
	}

	if output.ARN == "" {
		t.Error("Instance profile ARN is empty")
	}

	// Cleanup: delete the instance profile
	if err := DeleteInstanceProfile(ctx, client, output.Name); err != nil {
		t.Logf("Warning: Failed to cleanup test instance profile: %v", err)
	}
}

// TestGetInstanceProfile tests instance profile retrieval (integration test)
func TestGetInstanceProfile(t *testing.T) {
	// Skip if AWS credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping integration test: AWS credentials not available")
	}

	ctx := context.Background()
	client, err := NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	// Create a test instance profile first
	profile := &awsiam.InstanceProfile{
		Name: "test-get-profile-sdk",
		Path: stringPtr("/test/"),
	}

	created, err := CreateInstanceProfile(ctx, client, profile)
	if err != nil {
		t.Fatalf("Failed to create instance profile: %v", err)
	}

	// Get the instance profile
	retrieved, err := GetInstanceProfile(ctx, client, created.Name)
	if err != nil {
		t.Fatalf("Failed to get instance profile: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Retrieved instance profile is nil")
	}

	if retrieved.Name != created.Name {
		t.Errorf("Expected name %s, got %s", created.Name, retrieved.Name)
	}

	if retrieved.ARN != created.ARN {
		t.Errorf("Expected ARN %s, got %s", created.ARN, retrieved.ARN)
	}

	// Cleanup
	if err := DeleteInstanceProfile(ctx, client, created.Name); err != nil {
		t.Logf("Warning: Failed to cleanup test instance profile: %v", err)
	}
}

// TestListInstanceProfiles tests listing instance profiles (integration test)
func TestListInstanceProfiles(t *testing.T) {
	// Skip if AWS credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping integration test: AWS credentials not available")
	}

	ctx := context.Background()
	client, err := NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	profiles, err := ListInstanceProfiles(ctx, client, nil)
	if err != nil {
		t.Fatalf("Failed to list instance profiles: %v", err)
	}

	// Should have at least 0 profiles (may be empty)
	if profiles == nil {
		t.Error("Profiles list is nil")
	}
}

