package sdk_test

import (
	"context"
	"os"
	"testing"

	awsInstanceTypes "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/instance_types"
	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
)

func TestInstanceTypeService_GetInstanceType_Static(t *testing.T) {
	// Create a service with a nil client (will use static data)
	client, err := awssdk.NewAWSClient(context.Background())
	if err != nil {
		// If we can't create a client, create a service with nil client for static data testing
		// This test focuses on static data fallback
		t.Skip("Skipping test: AWS client creation failed, but this is OK for static data tests")
	}

	service, err := awssdk.NewInstanceTypeService(client)
	if err != nil {
		t.Fatalf("Failed to create InstanceTypeService: %v", err)
	}

	ctx := context.Background()
	region := "us-east-1"

	// Test getting a known instance type from static data
	info, err := service.GetInstanceType(ctx, "t3.micro", region)
	if err != nil {
		t.Fatalf("Failed to get t3.micro: %v", err)
	}

	if info.Name != "t3.micro" {
		t.Errorf("Expected name t3.micro, got %s", info.Name)
	}

	if info.Category != awsInstanceTypes.CategoryFreeTier {
		t.Errorf("Expected category %v, got %v", awsInstanceTypes.CategoryFreeTier, info.Category)
	}

	if !info.FreeTierEligible {
		t.Error("Expected t3.micro to be free tier eligible")
	}

	if info.VCPU != 2 {
		t.Errorf("Expected 2 vCPU, got %d", info.VCPU)
	}

	if info.MemoryGiB != 1.0 {
		t.Errorf("Expected 1.0 GiB memory, got %.2f", info.MemoryGiB)
	}
}

func TestInstanceTypeService_ListByCategory(t *testing.T) {
	client, err := awssdk.NewAWSClient(context.Background())
	if err != nil {
		t.Skip("Skipping test: AWS client creation failed")
	}

	service, err := awssdk.NewInstanceTypeService(client)
	if err != nil {
		t.Fatalf("Failed to create InstanceTypeService: %v", err)
	}

	ctx := context.Background()
	region := "us-east-1"

	// Test listing free tier instances
	freeTier, err := service.ListFreeTier(ctx, region)
	if err != nil {
		t.Fatalf("Failed to list free tier instances: %v", err)
	}

	if len(freeTier) == 0 {
		t.Error("Expected at least one free tier instance type")
	}

	// Verify all returned instances are free tier eligible
	for _, info := range freeTier {
		if !info.FreeTierEligible {
			t.Errorf("Instance %s should be free tier eligible", info.Name)
		}
		if info.Category != awsInstanceTypes.CategoryFreeTier {
			t.Errorf("Instance %s should have category %v, got %v", info.Name, awsInstanceTypes.CategoryFreeTier, info.Category)
		}
	}
}

func TestInstanceTypeService_ListByCategory_GeneralPurpose(t *testing.T) {
	client, err := awssdk.NewAWSClient(context.Background())
	if err != nil {
		t.Skip("Skipping test: AWS client creation failed")
	}

	service, err := awssdk.NewInstanceTypeService(client)
	if err != nil {
		t.Fatalf("Failed to create InstanceTypeService: %v", err)
	}

	ctx := context.Background()
	region := "us-east-1"

	// Test listing general purpose instances
	generalPurpose, err := service.ListByCategory(ctx, awsInstanceTypes.CategoryGeneralPurpose, region)
	if err != nil {
		t.Fatalf("Failed to list general purpose instances: %v", err)
	}

	if len(generalPurpose) == 0 {
		t.Error("Expected at least one general purpose instance type")
	}

	// Verify all returned instances are general purpose
	for _, info := range generalPurpose {
		if info.Category != awsInstanceTypes.CategoryGeneralPurpose {
			t.Errorf("Instance %s should have category %v, got %v", info.Name, awsInstanceTypes.CategoryGeneralPurpose, info.Category)
		}
	}
}

func TestInstanceTypeService_ListInstanceTypes_WithFilters(t *testing.T) {
	// Skip if credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping test: AWS_ACCESS_KEY_ID not set")
	}

	client, err := awssdk.NewAWSClient(context.Background())
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	service, err := awssdk.NewInstanceTypeService(client)
	if err != nil {
		t.Fatalf("Failed to create InstanceTypeService: %v", err)
	}

	ctx := context.Background()
	region := "us-east-1"

	// Test filtering by minimum VCPU
	minVCPU := 4
	filters := &awsInstanceTypes.InstanceTypeFilters{
		MinVCPU: &minVCPU,
	}

	types, err := service.ListInstanceTypes(ctx, region, filters)
	if err != nil {
		t.Fatalf("Failed to list instance types with filters: %v", err)
	}

	// Verify all returned instances have at least minVCPU
	for _, info := range types {
		if info.VCPU < minVCPU {
			t.Errorf("Instance %s has %d vCPU, expected at least %d", info.Name, info.VCPU, minVCPU)
		}
	}

	t.Logf("Found %d instance types with at least %d vCPU", len(types), minVCPU)
}

func TestInstanceTypeService_GetInstanceType_AWS(t *testing.T) {
	// Skip if credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping test: AWS_ACCESS_KEY_ID not set")
	}

	client, err := awssdk.NewAWSClient(context.Background())
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	service, err := awssdk.NewInstanceTypeService(client)
	if err != nil {
		t.Fatalf("Failed to create InstanceTypeService: %v", err)
	}

	ctx := context.Background()
	region := client.GetRegion()

	// Test getting a specific instance type from AWS API
	info, err := service.GetInstanceType(ctx, "t3.micro", region)
	if err != nil {
		t.Fatalf("Failed to get t3.micro from AWS: %v", err)
	}

	if info.Name != "t3.micro" {
		t.Errorf("Expected name t3.micro, got %s", info.Name)
	}

	if info.VCPU <= 0 {
		t.Errorf("Expected VCPU > 0, got %d", info.VCPU)
	}

	if info.MemoryGiB <= 0 {
		t.Errorf("Expected MemoryGiB > 0, got %.2f", info.MemoryGiB)
	}

	t.Logf("Retrieved t3.micro: %d vCPU, %.2f GiB memory, category: %v", info.VCPU, info.MemoryGiB, info.Category)
}
