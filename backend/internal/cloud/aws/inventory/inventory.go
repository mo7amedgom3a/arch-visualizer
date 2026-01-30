package inventory

import (
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

// ResourceClassification represents how a resource is classified
type ResourceClassification struct {
	Category     string   // "Networking", "Compute", etc.
	ResourceName string   // "VPC", "Subnet", "EC2", etc.
	Aliases      []string // ["vpc", "internet-gateway"] for IR mapping
	IRType       string   // IR resource type (kebab-case) for mapping from diagram
}

// FunctionRegistry holds function references for dynamic dispatch
type FunctionRegistry struct {
	// TerraformMapper maps a domain resource to Terraform blocks
	TerraformMapper func(*resource.Resource) ([]tfmapper.TerraformBlock, error)
	
	// PricingCalculator calculates cost for a resource
	PricingCalculator func(*resource.Resource, time.Duration) (*pricing.CostEstimate, error)
	
	// GetPricingInfo retrieves pricing information for a resource type
	GetPricingInfo func(string) (*pricing.ResourcePricing, error)
	
	// DomainMapper converts database resource to domain model (optional)
	DomainMapper func(interface{}) (interface{}, error)
	
	// AWSModelMapper converts domain model to AWS model (optional)
	AWSModelMapper func(interface{}) (interface{}, error)
}

// Inventory holds all resource classifications and function mappings for a cloud provider
type Inventory struct {
	// Classifications maps resource name to its classification
	Classifications map[string]ResourceClassification
	
	// Functions maps resource name to its function registry
	Functions map[string]FunctionRegistry
	
	// ByCategory maps category to list of resource names
	ByCategory map[string][]string
	
	// ByIRType maps IR type (kebab-case) to resource name for diagram parsing
	ByIRType map[string]string
}

// NewInventory creates a new empty inventory
func NewInventory() *Inventory {
	return &Inventory{
		Classifications: make(map[string]ResourceClassification),
		Functions:        make(map[string]FunctionRegistry),
		ByCategory:       make(map[string][]string),
		ByIRType:         make(map[string]string),
	}
}

// RegisterResource registers a resource in the inventory
func (inv *Inventory) RegisterResource(classification ResourceClassification, functions FunctionRegistry) {
	resourceName := classification.ResourceName
	
	// Store classification
	inv.Classifications[resourceName] = classification
	
	// Store functions
	inv.Functions[resourceName] = functions
	
	// Index by category
	category := classification.Category
	if inv.ByCategory[category] == nil {
		inv.ByCategory[category] = make([]string, 0)
	}
	inv.ByCategory[category] = append(inv.ByCategory[category], resourceName)
	
	// Index by IR type
	if classification.IRType != "" {
		inv.ByIRType[classification.IRType] = resourceName
	}
	
	// Index by aliases
	for _, alias := range classification.Aliases {
		inv.ByIRType[alias] = resourceName
	}
}

// GetResourceClassification retrieves classification for a resource
func (inv *Inventory) GetResourceClassification(resourceName string) (ResourceClassification, bool) {
	classification, ok := inv.Classifications[resourceName]
	return classification, ok
}

// GetFunctions retrieves function registry for a resource
func (inv *Inventory) GetFunctions(resourceName string) (FunctionRegistry, bool) {
	functions, ok := inv.Functions[resourceName]
	return functions, ok
}

// GetResourcesByCategory returns all resource names in a category
func (inv *Inventory) GetResourcesByCategory(category string) []string {
	return inv.ByCategory[category]
}

// GetResourceNameByIRType maps IR type to resource name
func (inv *Inventory) GetResourceNameByIRType(irType string) (string, bool) {
	resourceName, ok := inv.ByIRType[irType]
	return resourceName, ok
}

// SupportsResource checks if a resource type is supported
func (inv *Inventory) SupportsResource(resourceName string) bool {
	_, ok := inv.Classifications[resourceName]
	return ok
}
