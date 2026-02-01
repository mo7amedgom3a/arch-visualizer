package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// PricingService handles pricing calculations and persistence
type PricingService interface {
	// CalculateResourceCost calculates the cost for a single resource over a given duration
	CalculateResourceCost(ctx context.Context, res *resource.Resource, duration time.Duration) (*domainpricing.CostEstimate, error)

	// CalculateArchitectureCost calculates the total cost for an architecture over a given duration
	CalculateArchitectureCost(ctx context.Context, arch *architecture.Architecture, duration time.Duration) (*ArchitectureCostEstimate, error)

	// PersistResourcePricing saves resource pricing to the database
	PersistResourcePricing(ctx context.Context, projectID, resourceID uuid.UUID, estimate *domainpricing.CostEstimate, provider, region string) error

	// PersistProjectPricing saves project-level pricing to the database
	PersistProjectPricing(ctx context.Context, projectID uuid.UUID, estimate *domainpricing.CostEstimate, provider, region string) error

	// GetProjectPricing retrieves pricing for a project
	GetProjectPricing(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectPricing, error)

	// GetResourcePricing retrieves pricing for a resource
	GetResourcePricing(ctx context.Context, resourceID uuid.UUID) ([]*models.ResourcePricing, error)
}

// ArchitectureCostEstimate contains the cost estimate for an entire architecture
type ArchitectureCostEstimate struct {
	// TotalCost is the total estimated cost for the architecture
	TotalCost float64 `json:"total_cost"`
	// Currency is the currency used
	Currency string `json:"currency"`
	// Period is the time period for the estimate
	Period string `json:"period"`
	// Duration is the duration used for the calculation
	Duration time.Duration `json:"duration"`
	// ResourceEstimates contains individual resource cost estimates
	ResourceEstimates map[string]*ResourceCostEstimate `json:"resource_estimates"`
	// Provider is the cloud provider
	Provider string `json:"provider"`
	// Region is the region (if applicable)
	Region string `json:"region,omitempty"`
}

// ResourceCostEstimate contains the cost estimate for a single resource
type ResourceCostEstimate struct {
	// ResourceID is the domain resource ID
	ResourceID string `json:"resource_id"`
	// ResourceName is the name of the resource
	ResourceName string `json:"resource_name"`
	// ResourceType is the type of resource
	ResourceType string `json:"resource_type"`
	// TotalCost is the total estimated cost for this resource
	TotalCost float64 `json:"total_cost"`
	// Currency is the currency used
	Currency string `json:"currency"`
	// Breakdown contains the cost breakdown components
	Breakdown []CostBreakdownComponent `json:"breakdown"`
}

// CostBreakdownComponent represents a single cost component in the breakdown
type CostBreakdownComponent struct {
	// ComponentName is the name of the cost component
	ComponentName string `json:"component_name"`
	// Model is the pricing model used
	Model string `json:"model"`
	// Quantity is the amount consumed
	Quantity float64 `json:"quantity"`
	// UnitRate is the rate per unit
	UnitRate float64 `json:"unit_rate"`
	// Subtotal is the calculated cost for this component
	Subtotal float64 `json:"subtotal"`
	// Currency is the currency used
	Currency string `json:"currency"`
}

// PricingRepository defines pricing repository operations
type PricingRepository interface {
	// CreateProjectPricing creates project-level pricing
	CreateProjectPricing(ctx context.Context, pricing *models.ProjectPricing) error
	// FindProjectPricingByProjectID finds pricing for a project
	FindProjectPricingByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectPricing, error)
	// CreateResourcePricing creates resource-level pricing
	CreateResourcePricing(ctx context.Context, pricing *models.ResourcePricing) error
	// FindResourcePricingByResourceID finds pricing for a resource
	FindResourcePricingByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*models.ResourcePricing, error)
	// FindResourcePricingByProjectID finds all resource pricing for a project
	FindResourcePricingByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ResourcePricing, error)
	// CreatePricingComponent creates a pricing component
	CreatePricingComponent(ctx context.Context, component *models.PricingComponent) error
}
