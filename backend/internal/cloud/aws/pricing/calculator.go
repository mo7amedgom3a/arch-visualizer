package pricing

import (
	"context"
	"fmt"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing/networking"
	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// AWSPricingCalculator implements the PricingCalculator interface for AWS
type AWSPricingCalculator struct {
	service *AWSPricingService
}

// NewAWSPricingCalculator creates a new AWS pricing calculator
func NewAWSPricingCalculator(service *AWSPricingService) *AWSPricingCalculator {
	return &AWSPricingCalculator{
		service: service,
	}
}

// CalculateResourceCost calculates the cost for a single resource over a given duration
func (c *AWSPricingCalculator) CalculateResourceCost(ctx context.Context, res *resource.Resource, duration time.Duration) (*domainpricing.CostEstimate, error) {
	if res.Provider != "aws" {
		return nil, fmt.Errorf("unsupported provider: %s", res.Provider)
	}

	// Get pricing information for the resource type
	pricingInfo, err := c.GetResourcePricing(ctx, res.Type.Name, "aws", res.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to get pricing for resource type %s: %w", res.Type.Name, err)
	}

	// Calculate cost based on resource type
	var totalCost float64
	var breakdown []domainpricing.CostComponent

	switch res.Type.Name {
	case "nat_gateway":
		// For NAT Gateway, we need to estimate data processing
		// In a real scenario, this would come from usage metrics
		// For now, we default to 0 - this can be extended to accept usage metrics
		estimatedDataGB := 0.0
		cost := networking.CalculateNATGatewayCost(duration, estimatedDataGB, res.Region)
		totalCost = cost
		breakdown = []domainpricing.CostComponent{
			{
				ComponentName: "NAT Gateway Hourly",
				Model:         domainpricing.PerHour,
				Quantity:      duration.Hours(),
				UnitRate:      pricingInfo.Components[0].Rate,
				Subtotal:      pricingInfo.Components[0].Rate * duration.Hours(),
				Currency:      domainpricing.USD,
			},
		}
		if estimatedDataGB > 0 {
			breakdown = append(breakdown, domainpricing.CostComponent{
				ComponentName: "NAT Gateway Data Processing",
				Model:         domainpricing.PerGB,
				Quantity:      estimatedDataGB,
				UnitRate:      pricingInfo.Components[1].Rate,
				Subtotal:      pricingInfo.Components[1].Rate * estimatedDataGB,
				Currency:      domainpricing.USD,
			})
		}

	case "elastic_ip":
		// Check if EIP is attached
		// For now, we default to false (unattached) which incurs charges
		// In a real scenario, this would come from resource state
		isAttached := false
		cost := networking.CalculateElasticIPCost(duration, isAttached, res.Region)
		totalCost = cost
		if !isAttached {
			breakdown = []domainpricing.CostComponent{
				{
					ComponentName: "Elastic IP Hourly (Unattached)",
					Model:         domainpricing.PerHour,
					Quantity:      duration.Hours(),
					UnitRate:      pricingInfo.Components[0].Rate,
					Subtotal:      cost,
					Currency:      domainpricing.USD,
				},
			}
		}

	case "network_interface":
		// Check if ENI is attached
		// For now, we default to false (unattached) which incurs charges
		// In a real scenario, this would come from resource state
		isAttached := false
		cost := networking.CalculateNetworkInterfaceCost(duration, isAttached, res.Region)
		totalCost = cost
		if !isAttached {
			breakdown = []domainpricing.CostComponent{
				{
					ComponentName: "Network Interface Hourly (Unattached)",
					Model:         domainpricing.PerHour,
					Quantity:      duration.Hours(),
					UnitRate:      pricingInfo.Components[0].Rate,
					Subtotal:      cost,
					Currency:      domainpricing.USD,
				},
			}
		}

	default:
		// For other resource types, use generic calculation
		// This can be extended for other resource types
		return nil, fmt.Errorf("pricing calculation not yet implemented for resource type: %s", res.Type.Name)
	}

	// Determine period based on duration
	var period domainpricing.Period
	if duration.Hours() <= 24 {
		period = domainpricing.Hourly
	} else if duration.Hours() <= 720 {
		period = domainpricing.Monthly
	} else {
		period = domainpricing.Yearly
	}

	return &domainpricing.CostEstimate{
		TotalCost:    totalCost,
		Currency:     domainpricing.USD,
		Breakdown:    breakdown,
		Period:       period,
		Duration:     duration,
		CalculatedAt: time.Now(),
		ResourceType: &res.Type.Name,
		Provider:     domainpricing.AWS,
		Region:       &res.Region,
	}, nil
}

// CalculateArchitectureCost calculates the total cost for multiple resources over a given duration
func (c *AWSPricingCalculator) CalculateArchitectureCost(ctx context.Context, resources []*resource.Resource, duration time.Duration) (*domainpricing.CostEstimate, error) {
	var totalCost float64
	var allBreakdown []domainpricing.CostComponent

	for _, res := range resources {
		estimate, err := c.CalculateResourceCost(ctx, res, duration)
		if err != nil {
			// Log error but continue with other resources
			continue
		}
		totalCost += estimate.TotalCost
		allBreakdown = append(allBreakdown, estimate.Breakdown...)
	}

	// Determine period based on duration
	var period domainpricing.Period
	if duration.Hours() <= 24 {
		period = domainpricing.Hourly
	} else if duration.Hours() <= 720 {
		period = domainpricing.Monthly
	} else {
		period = domainpricing.Yearly
	}

	return &domainpricing.CostEstimate{
		TotalCost:    totalCost,
		Currency:     domainpricing.USD,
		Breakdown:    allBreakdown,
		Period:       period,
		Duration:     duration,
		CalculatedAt: time.Now(),
		Provider:     domainpricing.AWS,
	}, nil
}

// GetResourcePricing retrieves the pricing information for a specific resource type
func (c *AWSPricingCalculator) GetResourcePricing(ctx context.Context, resourceType string, provider string, region string) (*domainpricing.ResourcePricing, error) {
	return c.service.GetPricing(ctx, resourceType, provider, region)
}
