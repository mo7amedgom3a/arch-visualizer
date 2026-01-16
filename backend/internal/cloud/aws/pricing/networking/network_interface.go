package networking

import (
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// CalculateNetworkInterfaceCost calculates the cost for a Network Interface
// duration: time duration for the cost calculation
// isAttached: whether the ENI is attached to an instance (free if attached)
// region: AWS region
func CalculateNetworkInterfaceCost(duration time.Duration, isAttached bool, region string) float64 {
	// Network interfaces are free when attached to an instance
	if isAttached {
		return 0.0
	}

	// Base rate
	baseHourlyRate := 0.01
	
	// Get regional multiplier (default to 1.0)
	multiplier := 1.0
	
	// Calculate hourly cost
	hours := duration.Hours()
	return baseHourlyRate * multiplier * hours
}

// GetNetworkInterfacePricing returns the pricing information for Network Interface
func GetNetworkInterfacePricing(region string) *domainpricing.ResourcePricing {
	// Base rate
	baseHourlyRate := 0.01
	
	// Get regional multiplier (default to 1.0)
	multiplier := 1.0
	hourlyRate := baseHourlyRate * multiplier

	components := []domainpricing.PriceComponent{
		{
			Name:      "Network Interface Hourly (Unattached)",
			Model:     domainpricing.PerHour,
			Unit:      "hour",
			Rate:      hourlyRate,
			Currency:  domainpricing.USD,
			Region:    &region,
			Description: "Hourly charge when not attached to an instance. Free when attached.",
		},
	}

	return &domainpricing.ResourcePricing{
		ResourceType: "network_interface",
		Provider:     domainpricing.AWS,
		Components:  components,
		Metadata: map[string]interface{}{
			"base_hourly_rate": hourlyRate,
			"free_when_attached": true,
		},
	}
}
