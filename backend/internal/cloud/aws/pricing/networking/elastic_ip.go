package networking

import (
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// CalculateElasticIPCost calculates the cost for an Elastic IP
// duration: time duration for the cost calculation
// isAttached: whether the EIP is attached to a running instance (free if attached)
// region: AWS region
func CalculateElasticIPCost(duration time.Duration, isAttached bool, region string) float64 {
	// Elastic IPs are free when attached to a running instance
	if isAttached {
		return 0.0
	}

	// Base rate
	baseHourlyRate := 0.005
	
	// Get regional multiplier (default to 1.0)
	multiplier := 1.0
	
	// Calculate hourly cost
	hours := duration.Hours()
	return baseHourlyRate * multiplier * hours
}

// GetElasticIPPricing returns the pricing information for Elastic IP
func GetElasticIPPricing(region string) *domainpricing.ResourcePricing {
	// Base rate
	baseHourlyRate := 0.005
	
	// Get regional multiplier (default to 1.0)
	multiplier := 1.0
	hourlyRate := baseHourlyRate * multiplier

	components := []domainpricing.PriceComponent{
		{
			Name:      "Elastic IP Hourly (Unattached)",
			Model:     domainpricing.PerHour,
			Unit:      "hour",
			Rate:      hourlyRate,
			Currency:  domainpricing.USD,
			Region:    &region,
			Description: "Hourly charge when not attached to a running instance. Free when attached.",
		},
	}

	return &domainpricing.ResourcePricing{
		ResourceType: "elastic_ip",
		Provider:     domainpricing.AWS,
		Components:  components,
		Metadata: map[string]interface{}{
			"base_hourly_rate": hourlyRate,
			"free_when_attached": true,
		},
	}
}
