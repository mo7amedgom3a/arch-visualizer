package pricing

import (
	"context"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// PricingCalculator defines the interface for calculating resource costs
type PricingCalculator interface {
	// CalculateResourceCost calculates the cost for a single resource over a given duration
	CalculateResourceCost(ctx context.Context, res *resource.Resource, duration time.Duration) (*CostEstimate, error)

	// CalculateArchitectureCost calculates the total cost for multiple resources over a given duration
	CalculateArchitectureCost(ctx context.Context, resources []*resource.Resource, duration time.Duration) (*CostEstimate, error)

	// GetResourcePricing retrieves the pricing information for a specific resource type
	GetResourcePricing(ctx context.Context, resourceType string, provider string, region string) (*ResourcePricing, error)
}
