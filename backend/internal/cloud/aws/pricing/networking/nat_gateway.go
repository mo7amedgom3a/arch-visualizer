package networking

import (
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// CalculateNATGatewayCost calculates the cost for a NAT Gateway
// duration: time duration for the cost calculation
// dataProcessedGB: amount of data processed in GB (optional, can be 0)
// region: AWS region
func CalculateNATGatewayCost(duration time.Duration, dataProcessedGB float64, region string) float64 {
	// Base rates (from rates.go, but accessed directly to avoid import cycle)
	baseHourlyRate := 0.045
	dataProcessingRate := 0.045
	
	// Get regional multiplier (default to 1.0)
	multiplier := 1.0
	
	// Calculate hourly cost
	hours := duration.Hours()
	hourlyCost := baseHourlyRate * multiplier * hours
	
	// Calculate data processing cost
	dataProcessingCost := dataProcessingRate * multiplier * dataProcessedGB
	
	return hourlyCost + dataProcessingCost
}

// GetNATGatewayPricing returns the pricing information for NAT Gateway
func GetNATGatewayPricing(region string) *domainpricing.ResourcePricing {
	// Base rates
	baseHourlyRate := 0.045
	dataProcessingRate := 0.045
	
	// Get regional multiplier (default to 1.0)
	multiplier := 1.0
	hourlyRate := baseHourlyRate * multiplier
	dataProcRate := dataProcessingRate * multiplier

	components := []domainpricing.PriceComponent{
		{
			Name:      "NAT Gateway Hourly",
			Model:     domainpricing.PerHour,
			Unit:      "hour",
			Rate:      hourlyRate,
			Currency:  domainpricing.USD,
			Region:    &region,
			Description: "Base hourly charge for NAT Gateway",
		},
		{
			Name:      "NAT Gateway Data Processing",
			Model:     domainpricing.PerGB,
			Unit:      "GB",
			Rate:      dataProcRate,
			Currency:  domainpricing.USD,
			Region:    &region,
			Description: "Data processing charge per GB",
		},
	}

	return &domainpricing.ResourcePricing{
		ResourceType: "nat_gateway",
		Provider:     domainpricing.AWS,
		Components:  components,
		Metadata: map[string]interface{}{
			"base_hourly_rate": hourlyRate,
			"data_processing_rate": dataProcRate,
		},
	}
}
