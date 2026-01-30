package inventory

import (
	"testing"
	"time"

	architecture "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

func TestGetDefaultInventory(t *testing.T) {
	inv := GetDefaultInventory()

	if inv == nil {
		t.Fatal("GetDefaultInventory() returned nil")
	}

	// Verify default inventory has resources registered
	if len(inv.Classifications) == 0 {
		t.Error("Default inventory should have at least one resource classification")
	}

	// Verify some expected resources are present
	expectedResources := []string{"VPC", "Subnet", "EC2", "S3"}
	for _, resName := range expectedResources {
		if !inv.SupportsResource(resName) {
			t.Errorf("Expected resource %s not found in default inventory", resName)
		}
	}
}

func TestDefaultInventory_IRTypeMapping(t *testing.T) {
	inv := GetDefaultInventory()

	// Test that default inventory can map IR types
	testCases := []struct {
		irType      string
		expectedRes string
	}{
		{"vpc", "VPC"},
		{"subnet", "Subnet"},
		{"ec2", "EC2"},
		{"s3", "S3"},
		{"igw", "InternetGateway"},
		{"sg", "SecurityGroup"},
	}

	for _, tc := range testCases {
		resourceName, ok := inv.GetResourceNameByIRType(tc.irType)
		if !ok {
			t.Errorf("Failed to map IR type %s in default inventory", tc.irType)
			continue
		}

		if resourceName != tc.expectedRes {
			t.Errorf("Expected %s for IR type %s, got %s", tc.expectedRes, tc.irType, resourceName)
		}
	}
}

func TestDefaultInventory_CategoryIndexing(t *testing.T) {
	inv := GetDefaultInventory()

	// Verify categories are indexed
	categories := []string{
		resource.CategoryNetworking,
		resource.CategoryCompute,
		resource.CategoryStorage,
		resource.CategoryDatabase,
	}

	for _, category := range categories {
		resources := inv.GetResourcesByCategory(category)
		if len(resources) == 0 {
			t.Errorf("Expected at least one resource in category %s", category)
		}
	}
}

func TestAWSInventoryMapperAdapter(t *testing.T) {
	inv := NewInventory()

	// Register a test resource
	classification := ResourceClassification{
		Category:     resource.CategoryNetworking,
		ResourceName: "VPC",
		IRType:       "vpc",
		Aliases:      []string{"vpc", "virtual-private-cloud"},
	}

	inv.RegisterResource(classification, FunctionRegistry{})

	// Create adapter
	adapter := &awsInventoryMapperAdapter{inventory: inv}

	// Test GetResourceNameByIRType
	resourceName, ok := adapter.GetResourceNameByIRType("vpc")
	if !ok {
		t.Fatal("Failed to retrieve resource name by IR type")
	}

	if resourceName != "VPC" {
		t.Errorf("Expected VPC, got %s", resourceName)
	}

	// Test with alias
	resourceName, ok = adapter.GetResourceNameByIRType("virtual-private-cloud")
	if !ok {
		t.Fatal("Failed to retrieve resource name by alias")
	}

	if resourceName != "VPC" {
		t.Errorf("Expected VPC for alias, got %s", resourceName)
	}

	// Test non-existing IR type
	_, ok = adapter.GetResourceNameByIRType("non-existent")
	if ok {
		t.Error("Expected false for non-existing IR type")
	}
}

func TestDefaultInventory_IRTypeMapperRegistration(t *testing.T) {
	// Verify that the default inventory is registered as an IR type mapper
	mapper, ok := architecture.GetIRTypeMapper(resource.AWS)
	if !ok {
		t.Fatal("AWS IR type mapper should be registered")
	}

	if mapper == nil {
		t.Fatal("IR type mapper should not be nil")
	}

	// Test that the mapper works
	resourceName, ok := mapper.GetResourceNameByIRType("vpc")
	if !ok {
		t.Fatal("Failed to retrieve resource name via registered mapper")
	}

	if resourceName != "VPC" {
		t.Errorf("Expected VPC, got %s", resourceName)
	}
}

func TestSetTerraformMapper_Integration(t *testing.T) {
	inv := GetDefaultInventory()

	// Set a test mapper
	mapperCalled := false
	testMapper := func(*resource.Resource) ([]tfmapper.TerraformBlock, error) {
		mapperCalled = true
		return []tfmapper.TerraformBlock{
			{
				Kind:   "resource",
				Labels: []string{"aws_vpc", "test"},
			},
		}, nil
	}

	inv.SetTerraformMapper("VPC", testMapper)

	// Retrieve and test
	functions, ok := inv.GetFunctions("VPC")
	if !ok {
		t.Fatal("Failed to retrieve functions")
	}

	if functions.TerraformMapper == nil {
		t.Fatal("TerraformMapper should be set")
	}

	blocks, err := functions.TerraformMapper(nil)
	if err != nil {
		t.Errorf("TerraformMapper returned error: %v", err)
	}

	if len(blocks) != 1 {
		t.Errorf("Expected 1 block, got %d", len(blocks))
	}

	if !mapperCalled {
		t.Error("Test mapper was not called")
	}
}

func TestSetPricingCalculator_Integration(t *testing.T) {
	inv := GetDefaultInventory()

	// Set a test calculator
	calculatorCalled := false
	testCalculator := func(*resource.Resource, time.Duration) (*domainpricing.CostEstimate, error) {
		calculatorCalled = true
		return &domainpricing.CostEstimate{
			TotalCost: 10.0,
			Currency:  domainpricing.USD,
			Period:    domainpricing.Monthly,
			Duration:  time.Hour * 24 * 30,
		}, nil
	}

	inv.SetPricingCalculator("EC2", testCalculator)

	// Retrieve and test
	functions, ok := inv.GetFunctions("EC2")
	if !ok {
		t.Fatal("Failed to retrieve functions")
	}

	if functions.PricingCalculator == nil {
		t.Fatal("PricingCalculator should be set")
	}

	estimate, err := functions.PricingCalculator(nil, time.Hour)
	if err != nil {
		t.Errorf("PricingCalculator returned error: %v", err)
	}

	if estimate == nil {
		t.Fatal("CostEstimate should not be nil")
	}

	if estimate.TotalCost != 10.0 {
		t.Errorf("Expected total cost 10.0, got %f", estimate.TotalCost)
	}

	if !calculatorCalled {
		t.Error("Test calculator was not called")
	}
}

func TestSetPricingInfoGetter_Integration(t *testing.T) {
	inv := GetDefaultInventory()

	// Set a test getter
	getterCalled := false
	testGetter := func(string) (*domainpricing.ResourcePricing, error) {
		getterCalled = true
		return &domainpricing.ResourcePricing{
			ResourceType: "ec2",
			Provider:     domainpricing.AWS,
			Components: []domainpricing.PriceComponent{
				{
					Name:     "instance",
					Rate:     0.01,
					Unit:     "hour",
					Model:    domainpricing.PerHour,
					Currency: domainpricing.USD,
				},
			},
		}, nil
	}

	inv.SetPricingInfoGetter("EC2", testGetter)

	// Retrieve and test
	functions, ok := inv.GetFunctions("EC2")
	if !ok {
		t.Fatal("Failed to retrieve functions")
	}

	if functions.GetPricingInfo == nil {
		t.Fatal("GetPricingInfo should be set")
	}

	pricing, err := functions.GetPricingInfo("us-east-1")
	if err != nil {
		t.Errorf("GetPricingInfo returned error: %v", err)
	}

	if pricing == nil {
		t.Fatal("ResourcePricing should not be nil")
	}

	if pricing.ResourceType != "ec2" {
		t.Errorf("Expected resource type ec2, got %s", pricing.ResourceType)
	}

	if len(pricing.Components) == 0 {
		t.Error("Expected at least one price component")
	}

	if !getterCalled {
		t.Error("Test getter was not called")
	}
}

func TestAddingNewResource(t *testing.T) {
	// Create a new inventory for testing (don't modify default)
	inv := NewInventory()

	// Add a new custom resource
	newResource := ResourceClassification{
		Category:     resource.CategoryNetworking,
		ResourceName: "CustomGateway",
		IRType:       "custom-gateway",
		Aliases:      []string{"custom-gateway", "cgw"},
	}

	functions := FunctionRegistry{
		TerraformMapper: func(*resource.Resource) ([]tfmapper.TerraformBlock, error) {
			return []tfmapper.TerraformBlock{
				{
					Kind:   "resource",
					Labels: []string{"aws_custom_gateway", "test"},
				},
			}, nil
		},
	}

	inv.RegisterResource(newResource, functions)

	// Verify resource is registered
	if !inv.SupportsResource("CustomGateway") {
		t.Error("CustomGateway should be supported")
	}

	// Verify IR type mapping
	resourceName, ok := inv.GetResourceNameByIRType("custom-gateway")
	if !ok {
		t.Fatal("Failed to map custom-gateway IR type")
	}

	if resourceName != "CustomGateway" {
		t.Errorf("Expected CustomGateway, got %s", resourceName)
	}

	// Verify alias mapping
	resourceName, ok = inv.GetResourceNameByIRType("cgw")
	if !ok {
		t.Fatal("Failed to map cgw alias")
	}

	if resourceName != "CustomGateway" {
		t.Errorf("Expected CustomGateway for alias, got %s", resourceName)
	}

	// Verify category indexing
	customResources := inv.GetResourcesByCategory(resource.CategoryNetworking)
	found := false
	for _, res := range customResources {
		if res == "CustomGateway" {
			found = true
			break
		}
	}

	if !found {
		t.Error("CustomGateway should be in networking category")
	}

	// Verify functions work
	retrievedFunctions, ok := inv.GetFunctions("CustomGateway")
	if !ok {
		t.Fatal("Failed to retrieve functions for CustomGateway")
	}

	if retrievedFunctions.TerraformMapper == nil {
		t.Fatal("TerraformMapper should be set")
	}

	blocks, err := retrievedFunctions.TerraformMapper(nil)
	if err != nil {
		t.Errorf("TerraformMapper returned error: %v", err)
	}

	if len(blocks) != 1 {
		t.Errorf("Expected 1 block, got %d", len(blocks))
	}
}

func TestMultipleResourcesSameCategory(t *testing.T) {
	inv := NewInventory()

	// Register multiple resources in the same category
	resources := []ResourceClassification{
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "VPC1",
			IRType:       "vpc1",
		},
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "VPC2",
			IRType:       "vpc2",
		},
		{
			Category:     resource.CategoryNetworking,
			ResourceName: "VPC3",
			IRType:       "vpc3",
		},
	}

	for _, res := range resources {
		inv.RegisterResource(res, FunctionRegistry{})
	}

	// Verify all are in the category
	networkingResources := inv.GetResourcesByCategory(resource.CategoryNetworking)
	if len(networkingResources) != 3 {
		t.Errorf("Expected 3 networking resources, got %d", len(networkingResources))
	}

	// Verify all can be retrieved
	for _, res := range resources {
		if !inv.SupportsResource(res.ResourceName) {
			t.Errorf("Resource %s should be supported", res.ResourceName)
		}
	}
}
