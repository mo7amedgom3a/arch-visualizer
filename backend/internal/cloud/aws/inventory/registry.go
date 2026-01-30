package inventory

import (
	"time"

	architecture "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

var (
	// DefaultAWSInventory is the singleton AWS inventory instance
	DefaultAWSInventory *Inventory
)

func init() {
	DefaultAWSInventory = NewInventory()
	
	// Register all resources (without function mappings - those are set lazily)
	classifications := GetAWSResourceClassifications()
	for _, classification := range classifications {
		// Register with empty function registry - functions will be set by the mapper/pricing services
		DefaultAWSInventory.RegisterResource(classification, FunctionRegistry{})
	}
	
	// Register AWS inventory as IR type mapper for the domain architecture layer
	registerAWSInventoryAsMapper()
}

// registerAWSInventoryAsMapper registers the AWS inventory as an IR type mapper
// This allows the domain architecture layer to use AWS inventory for dynamic IR type mapping
func registerAWSInventoryAsMapper() {
	// Create an adapter that implements IRTypeMapper interface
	adapter := &awsInventoryMapperAdapter{inventory: DefaultAWSInventory}
	architecture.RegisterIRTypeMapper(resource.AWS, adapter)
}

// awsInventoryMapperAdapter adapts the Inventory to the IRTypeMapper interface
type awsInventoryMapperAdapter struct {
	inventory *Inventory
}

// GetResourceNameByIRType implements IRTypeMapper interface
func (a *awsInventoryMapperAdapter) GetResourceNameByIRType(irType string) (string, bool) {
	return a.inventory.GetResourceNameByIRType(irType)
}

// SetTerraformMapper sets the Terraform mapper function for a resource
func (inv *Inventory) SetTerraformMapper(resourceName string, mapper func(*resource.Resource) ([]tfmapper.TerraformBlock, error)) {
	if functions, ok := inv.Functions[resourceName]; ok {
		functions.TerraformMapper = mapper
		inv.Functions[resourceName] = functions
	}
}

// SetPricingCalculator sets the pricing calculator function for a resource
func (inv *Inventory) SetPricingCalculator(resourceName string, calculator func(*resource.Resource, time.Duration) (*domainpricing.CostEstimate, error)) {
	if functions, ok := inv.Functions[resourceName]; ok {
		functions.PricingCalculator = calculator
		inv.Functions[resourceName] = functions
	}
}

// SetPricingInfoGetter sets the pricing info getter function for a resource
func (inv *Inventory) SetPricingInfoGetter(resourceName string, getter func(string) (*domainpricing.ResourcePricing, error)) {
	if functions, ok := inv.Functions[resourceName]; ok {
		functions.GetPricingInfo = getter
		inv.Functions[resourceName] = functions
	}
}



// GetDefaultInventory returns the default AWS inventory
func GetDefaultInventory() *Inventory {
	return DefaultAWSInventory
}
