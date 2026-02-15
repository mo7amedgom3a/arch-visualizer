package networking

import (
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// CalculateVPCEndpointCost calculates the cost for a VPC Endpoint
// duration: time duration for the cost calculation
// endpointType: "Interface" or "Gateway"
// dataProcessedGB: amount of data processed in GB (optional, can be 0)
// numENIs: number of Availability Zones/Subnets the endpoint is deployed in (for Interface type)
// region: AWS region
func CalculateVPCEndpointCost(duration time.Duration, endpointType string, dataProcessedGB float64, numENIs int, region string) float64 {
	if endpointType == "Gateway" {
		return 0.0 // Gateway endpoints are free
	}

	// Interface Endpoint Pricing (defaulting to US-East-1 rates if not dynamic)
	// $0.01 per ENI-hour
	// $0.01 per GB data processed
	baseHourlyRate := 0.01
	dataProcessingRate := 0.01

	// Regional multiplier (placeholder for now)
	multiplier := 1.0

	// Calculate hourly cost
	hours := duration.Hours()
	hourlyCost := baseHourlyRate * multiplier * float64(numENIs) * hours

	// Calculate data processing cost
	// First 1 PB is $0.01/GB, drops after. Simplified to linear for now.
	dataProcessingCost := dataProcessingRate * multiplier * dataProcessedGB

	return hourlyCost + dataProcessingCost
}

// GetVPCEndpointPricing returns the pricing information for VPC Endpoint
func GetVPCEndpointPricing(endpointType string, region string) *domainpricing.ResourcePricing {
	if endpointType == "Gateway" {
		return &domainpricing.ResourcePricing{
			ResourceType: "vpc_endpoint",
			Provider:     domainpricing.AWS,
			Components: []domainpricing.PriceComponent{
				{
					Name:        "Gateway Endpoint",
					Model:       domainpricing.PerHour,
					Unit:        "n/a",
					Rate:        0.0,
					Currency:    domainpricing.USD,
					Region:      &region,
					Description: "Gateway Endpoints (S3, DynamoDB) are free",
				},
			},
		}
	}

	// Interface Endpoint
	baseHourlyRate := 0.01
	dataProcessingRate := 0.01

	return &domainpricing.ResourcePricing{
		ResourceType: "vpc_endpoint",
		Provider:     domainpricing.AWS,
		Components: []domainpricing.PriceComponent{
			{
				Name:        "Interface Endpoint Hourly (per ENI)",
				Model:       domainpricing.PerHour,
				Unit:        "ENI-hour",
				Rate:        baseHourlyRate,
				Currency:    domainpricing.USD,
				Region:      &region,
				Description: "Hourly charge per ENI used by the endpoint",
			},
			{
				Name:        "Interface Endpoint Data Processing",
				Model:       domainpricing.PerGB,
				Unit:        "GB",
				Rate:        dataProcessingRate,
				Currency:    domainpricing.USD,
				Region:      &region,
				Description: "Data processing charge per GB",
			},
		},
	}
}
