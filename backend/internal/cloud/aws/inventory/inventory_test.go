package inventory

import (
	"fmt"
	"testing"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

func TestNewInventory(t *testing.T) {
	inv := NewInventory()

	if inv == nil {
		t.Fatal("NewInventory() returned nil")
	}

	if inv.Classifications == nil {
		t.Error("Classifications map is nil")
	}

	if inv.Functions == nil {
		t.Error("Functions map is nil")
	}

	if inv.ByCategory == nil {
		t.Error("ByCategory map is nil")
	}

	if inv.ByIRType == nil {
		t.Error("ByIRType map is nil")
	}

	// Verify all maps are empty initially
	if len(inv.Classifications) != 0 {
		t.Errorf("Expected empty Classifications, got %d entries", len(inv.Classifications))
	}

	if len(inv.Functions) != 0 {
		t.Errorf("Expected empty Functions, got %d entries", len(inv.Functions))
	}

	if len(inv.ByCategory) != 0 {
		t.Errorf("Expected empty ByCategory, got %d entries", len(inv.ByCategory))
	}

	if len(inv.ByIRType) != 0 {
		t.Errorf("Expected empty ByIRType, got %d entries", len(inv.ByIRType))
	}
}

func TestRegisterResource(t *testing.T) {
	inv := NewInventory()

	classification := ResourceClassification{
		Category:     resource.CategoryNetworking,
		ResourceName: "VPC",
		IRType:       "vpc",
		Aliases:      []string{"vpc", "virtual-private-cloud"},
	}

	functions := FunctionRegistry{
		TerraformMapper: func(*resource.Resource) ([]tfmapper.TerraformBlock, error) {
			return nil, nil
		},
	}

	inv.RegisterResource(classification, functions)

	// Verify classification is stored
	retrieved, ok := inv.GetResourceClassification("VPC")
	if !ok {
		t.Fatal("Failed to retrieve registered resource classification")
	}

	if retrieved.Category != resource.CategoryNetworking {
		t.Errorf("Expected category %s, got %s", resource.CategoryNetworking, retrieved.Category)
	}

	if retrieved.ResourceName != "VPC" {
		t.Errorf("Expected resource name VPC, got %s", retrieved.ResourceName)
	}

	if retrieved.IRType != "vpc" {
		t.Errorf("Expected IR type vpc, got %s", retrieved.IRType)
	}

	// Verify functions are stored
	retrievedFunctions, ok := inv.GetFunctions("VPC")
	if !ok {
		t.Fatal("Failed to retrieve registered functions")
	}

	if retrievedFunctions.TerraformMapper == nil {
		t.Error("TerraformMapper function is nil")
	}

	// Verify category indexing
	networkingResources := inv.GetResourcesByCategory(resource.CategoryNetworking)
	if len(networkingResources) != 1 {
		t.Errorf("Expected 1 networking resource, got %d", len(networkingResources))
	}

	if networkingResources[0] != "VPC" {
		t.Errorf("Expected VPC in networking resources, got %s", networkingResources[0])
	}

	// Verify IR type indexing
	resourceName, ok := inv.GetResourceNameByIRType("vpc")
	if !ok {
		t.Fatal("Failed to retrieve resource name by IR type")
	}

	if resourceName != "VPC" {
		t.Errorf("Expected VPC, got %s", resourceName)
	}

	// Verify alias indexing
	for _, alias := range classification.Aliases {
		resourceName, ok := inv.GetResourceNameByIRType(alias)
		fmt.Println("alias", alias, "resourceName", resourceName)
		if !ok {
			t.Errorf("Failed to retrieve resource name by alias %s", alias)
		}

		if resourceName != "VPC" {
			t.Errorf("Expected VPC for alias %s, got %s", alias, resourceName)
		}
	}
}

func TestRegisterResource_MultipleResources(t *testing.T) {
	inv := NewInventory()

	// Register multiple resources in different categories
	resources := []ResourceClassification{
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "VPC",
			IRType:       "vpc",
			Aliases:      []string{"vpc"},
		},
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "Subnet",
			IRType:       "subnet",
			Aliases:      []string{"subnet"},
		},
		{
			Category:     resource.CategoryCompute,
			ResourceName: "EC2",
			IRType:       "ec2",
			Aliases:      []string{"ec2", "instance"},
		},
		{
			Category:     resource.CategoryStorage,
			ResourceName: "S3",
			IRType:       "s3",
			Aliases:      []string{"s3", "bucket"},
		},
	}

	for _, res := range resources {
		inv.RegisterResource(res, FunctionRegistry{})
	}

	// Verify networking resources
	networkingResources := inv.GetResourcesByCategory(resource.CategoryNetworking)
	if len(networkingResources) != 2 {
		t.Errorf("Expected 2 networking resources, got %d", len(networkingResources))
	}

	// Verify compute resources
	computeResources := inv.GetResourcesByCategory(resource.CategoryCompute)
	if len(computeResources) != 1 {
		t.Errorf("Expected 1 compute resource, got %d", len(computeResources))
	}

	if computeResources[0] != "EC2" {
		t.Errorf("Expected EC2, got %s", computeResources[0])
	}

	// Verify storage resources
	storageResources := inv.GetResourcesByCategory(resource.CategoryStorage)
	if len(storageResources) != 1 {
		t.Errorf("Expected 1 storage resource, got %d", len(storageResources))
	}

	if storageResources[0] != "S3" {
		t.Errorf("Expected S3, got %s", storageResources[0])
	}

	// Verify all IR types are mapped
	testCases := []struct {
		irType      string
		expectedRes string
	}{
		{"vpc", "VPC"},
		{"subnet", "Subnet"},
		{"ec2", "EC2"},
		{"instance", "EC2"}, // alias
		{"s3", "S3"},
		{"bucket", "S3"}, // alias
	}

	for _, tc := range testCases {
		resourceName, ok := inv.GetResourceNameByIRType(tc.irType)
		if !ok {
			t.Errorf("Failed to retrieve resource name for IR type %s", tc.irType)
			continue
		}

		if resourceName != tc.expectedRes {
			t.Errorf("Expected %s for IR type %s, got %s", tc.expectedRes, tc.irType, resourceName)
		}
	}
}

func TestGetResourceClassification(t *testing.T) {
	inv := NewInventory()

	classification := ResourceClassification{
		Category:     resource.CategoryNetworking,
		ResourceName: "VPC",
		IRType:       "vpc",
		Aliases:      []string{"vpc"},
	}

	inv.RegisterResource(classification, FunctionRegistry{})

	// Test existing resource
	retrieved, ok := inv.GetResourceClassification("VPC")
	if !ok {
		t.Fatal("Failed to retrieve existing resource classification")
	}

	if retrieved.ResourceName != "VPC" {
		t.Errorf("Expected VPC, got %s", retrieved.ResourceName)
	}

	// Test non-existing resource
	_, ok = inv.GetResourceClassification("NonExistent")
	if ok {
		t.Error("Expected false for non-existing resource, got true")
	}
}

func TestGetFunctions(t *testing.T) {
	inv := NewInventory()

	classification := ResourceClassification{
		Category:     resource.CategoryCompute,
		ResourceName: "EC2",
		IRType:       "ec2",
		Aliases:      []string{"ec2"},
	}

	functions := FunctionRegistry{
		TerraformMapper: func(*resource.Resource) ([]tfmapper.TerraformBlock, error) {
			return nil, nil
		},
	}

	inv.RegisterResource(classification, functions)

	// Test existing resource
	retrieved, ok := inv.GetFunctions("EC2")
	if !ok {
		t.Fatal("Failed to retrieve existing functions")
	}

	if retrieved.TerraformMapper == nil {
		t.Error("TerraformMapper function is nil")
	}

	// Test non-existing resource
	_, ok = inv.GetFunctions("NonExistent")
	if ok {
		t.Error("Expected false for non-existing resource, got true")
	}
}

func TestGetResourcesByCategory(t *testing.T) {
	inv := NewInventory()

	// Register resources in different categories
	inv.RegisterResource(ResourceClassification{
		Category:     resource.CategoryNetworking,
		ResourceName: "VPC",
		IRType:       "vpc",
	}, FunctionRegistry{})

	inv.RegisterResource(ResourceClassification{
		Category:     resource.CategoryNetworking,
		ResourceName: "Subnet",
		IRType:       "subnet",
	}, FunctionRegistry{})

	inv.RegisterResource(ResourceClassification{
		Category:     resource.CategoryCompute,
		ResourceName: "EC2",
		IRType:       "ec2",
	}, FunctionRegistry{})

	// Test existing category
	networkingResources := inv.GetResourcesByCategory(resource.CategoryNetworking)
	if len(networkingResources) != 2 {
		t.Errorf("Expected 2 networking resources, got %d", len(networkingResources))
	}

	// Test non-existing category
	emptyResources := inv.GetResourcesByCategory("NonExistent")
	if len(emptyResources) != 0 {
		t.Errorf("Expected empty slice for non-existing category, got %d items", len(emptyResources))
	}
}

func TestGetResourceNameByIRType(t *testing.T) {
	inv := NewInventory()

	classification := ResourceClassification{
		Category:     resource.CategoryNetworking,
		ResourceName: "InternetGateway",
		IRType:       "internet-gateway",
		Aliases:      []string{"internet-gateway", "internet_gateway", "igw"},
	}

	inv.RegisterResource(classification, FunctionRegistry{})

	// Test IR type
	resourceName, ok := inv.GetResourceNameByIRType("internet-gateway")
	if !ok {
		t.Fatal("Failed to retrieve resource name by IR type")
	}

	if resourceName != "InternetGateway" {
		t.Errorf("Expected InternetGateway, got %s", resourceName)
	}

	// Test aliases
	aliasTests := []string{"internet-gateway", "internet_gateway", "igw"}
	for _, alias := range aliasTests {
		resourceName, ok := inv.GetResourceNameByIRType(alias)
		if !ok {
			t.Errorf("Failed to retrieve resource name by alias %s", alias)
			continue
		}

		if resourceName != "InternetGateway" {
			t.Errorf("Expected InternetGateway for alias %s, got %s", alias, resourceName)
		}
	}

	// Test non-existing IR type
	_, ok = inv.GetResourceNameByIRType("non-existent")
	if ok {
		t.Error("Expected false for non-existing IR type, got true")
	}
}

func TestSupportsResource(t *testing.T) {
	inv := NewInventory()

	inv.RegisterResource(ResourceClassification{
		Category:     resource.CategoryCompute,
		ResourceName: "EC2",
		IRType:       "ec2",
	}, FunctionRegistry{})

	// Test existing resource
	if !inv.SupportsResource("EC2") {
		t.Error("Expected SupportsResource to return true for EC2")
	}

	// Test non-existing resource
	if inv.SupportsResource("NonExistent") {
		t.Error("Expected SupportsResource to return false for non-existing resource")
	}
}

func TestSetTerraformMapper(t *testing.T) {
	inv := NewInventory()

	classification := ResourceClassification{
		Category:     resource.CategoryCompute,
		ResourceName: "EC2",
		IRType:       "ec2",
	}

	inv.RegisterResource(classification, FunctionRegistry{})

	// Set Terraform mapper
	mapperCalled := false
	mapper := func(*resource.Resource) ([]tfmapper.TerraformBlock, error) {
		mapperCalled = true
		return nil, nil
	}

	inv.SetTerraformMapper("EC2", mapper)

	// Verify mapper is set
	functions, ok := inv.GetFunctions("EC2")
	if !ok {
		t.Fatal("Failed to retrieve functions")
	}

	if functions.TerraformMapper == nil {
		t.Fatal("TerraformMapper is nil after setting")
	}

	// Test that mapper can be called
	_, err := functions.TerraformMapper(nil)
	if err != nil {
		t.Errorf("TerraformMapper returned error: %v", err)
	}

	if !mapperCalled {
		t.Error("TerraformMapper was not called")
	}

	// Test setting mapper for non-existing resource (should not panic)
	inv.SetTerraformMapper("NonExistent", mapper)
}

func TestSetPricingCalculator(t *testing.T) {
	inv := NewInventory()

	classification := ResourceClassification{
		Category:     resource.CategoryCompute,
		ResourceName: "EC2",
		IRType:       "ec2",
	}

	inv.RegisterResource(classification, FunctionRegistry{})

	// Set pricing calculator
	calculatorCalled := false
	calculator := func(*resource.Resource, time.Duration) (*pricing.CostEstimate, error) {
		calculatorCalled = true
		return nil, nil
	}

	inv.SetPricingCalculator("EC2", calculator)

	// Verify calculator is set
	functions, ok := inv.GetFunctions("EC2")
	if !ok {
		t.Fatal("Failed to retrieve functions")
	}

	if functions.PricingCalculator == nil {
		t.Fatal("PricingCalculator is nil after setting")
	}

	// Test that calculator can be called
	_, err := functions.PricingCalculator(nil, time.Hour)
	if err != nil {
		t.Errorf("PricingCalculator returned error: %v", err)
	}

	if !calculatorCalled {
		t.Error("PricingCalculator was not called")
	}

	// Test setting calculator for non-existing resource (should not panic)
	inv.SetPricingCalculator("NonExistent", calculator)
}

func TestSetPricingInfoGetter(t *testing.T) {
	inv := NewInventory()

	classification := ResourceClassification{
		Category:     resource.CategoryCompute,
		ResourceName: "EC2",
		IRType:       "ec2",
	}

	inv.RegisterResource(classification, FunctionRegistry{})

	// Set pricing info getter
	getterCalled := false
	getter := func(string) (*pricing.ResourcePricing, error) {
		getterCalled = true
		return nil, nil
	}

	inv.SetPricingInfoGetter("EC2", getter)

	// Verify getter is set
	functions, ok := inv.GetFunctions("EC2")
	if !ok {
		t.Fatal("Failed to retrieve functions")
	}

	if functions.GetPricingInfo == nil {
		t.Fatal("GetPricingInfo is nil after setting")
	}

	// Test that getter can be called
	_, err := functions.GetPricingInfo("us-east-1")
	if err != nil {
		t.Errorf("GetPricingInfo returned error: %v", err)
	}

	if !getterCalled {
		t.Error("GetPricingInfo was not called")
	}

	// Test setting getter for non-existing resource (should not panic)
	inv.SetPricingInfoGetter("NonExistent", getter)
}

func TestRegisterResource_WithoutIRType(t *testing.T) {
	inv := NewInventory()

	// Register resource without IR type
	classification := ResourceClassification{
		Category:     resource.CategoryNetworking,
		ResourceName: "CustomResource",
		IRType:       "", // Empty IR type
		Aliases:      []string{"custom"},
	}

	inv.RegisterResource(classification, FunctionRegistry{})

	// IR type should not be indexed
	_, ok := inv.GetResourceNameByIRType("")
	if ok {
		t.Error("Empty IR type should not be indexed")
	}

	// But aliases should still work
	resourceName, ok := inv.GetResourceNameByIRType("custom")
	if !ok {
		t.Fatal("Alias should still work even without IR type")
	}

	if resourceName != "CustomResource" {
		t.Errorf("Expected CustomResource, got %s", resourceName)
	}
}

func TestRegisterResource_WithoutAliases(t *testing.T) {
	inv := NewInventory()

	// Register resource without aliases
	classification := ResourceClassification{
		Category:     resource.CategoryNetworking,
		ResourceName: "SimpleResource",
		IRType:       "simple-resource",
		Aliases:      nil, // No aliases
	}

	inv.RegisterResource(classification, FunctionRegistry{})

	// IR type should still work
	resourceName, ok := inv.GetResourceNameByIRType("simple-resource")
	if !ok {
		t.Fatal("IR type should work even without aliases")
	}

	if resourceName != "SimpleResource" {
		t.Errorf("Expected SimpleResource, got %s", resourceName)
	}
}
