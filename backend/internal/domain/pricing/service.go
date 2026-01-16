package pricing

import (
	"context"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// PricingService defines the interface for pricing operations
type PricingService interface {
	// GetPricing retrieves the pricing information for a specific resource type
	GetPricing(ctx context.Context, resourceType string, provider string, region string) (*ResourcePricing, error)
	
	// EstimateCost estimates the cost for a resource over a given duration
	EstimateCost(ctx context.Context, res *resource.Resource, duration time.Duration) (*CostEstimate, error)
	
	// EstimateArchitectureCost estimates the total cost for multiple resources over a given duration
	EstimateArchitectureCost(ctx context.Context, resources []*resource.Resource, duration time.Duration) (*CostEstimate, error)
	
	// ListSupportedResources returns a list of resource types that have pricing information
	ListSupportedResources(ctx context.Context, provider string) ([]string, error)
}
