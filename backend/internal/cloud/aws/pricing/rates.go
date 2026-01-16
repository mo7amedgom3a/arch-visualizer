package pricing

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// NetworkingPricingRates contains static pricing rates for AWS networking resources
// These rates are based on AWS public pricing as of 2024
// Note: Rates may vary by region and can be updated via AWS Pricing API integration in the future
var NetworkingPricingRates = map[string]AWSPricingRate{
	"nat_gateway": {
		BaseHourlyRate:     0.045, // $0.045 per hour
		DataProcessingRate: 0.045, // $0.045 per GB of data processed
		RegionalVariations: map[string]float64{
			// Some regions may have different rates
			"us-east-1": 1.0,  // Base rate multiplier
			"us-west-2": 1.0,  // Base rate multiplier
			"eu-west-1": 1.0,  // Base rate multiplier
		},
	},
	"elastic_ip": {
		BaseHourlyRate: 0.005, // $0.005 per hour when not attached to a running instance
		// Note: Elastic IPs are free when attached to a running instance
		RegionalVariations: map[string]float64{
			"us-east-1": 1.0,
			"us-west-2": 1.0,
			"eu-west-1": 1.0,
		},
	},
	"network_interface": {
		BaseHourlyRate: 0.01, // $0.01 per hour when not attached
		// Note: Network interfaces are free when attached to an instance
		RegionalVariations: map[string]float64{
			"us-east-1": 1.0,
			"us-west-2": 1.0,
			"eu-west-1": 1.0,
		},
	},
	"data_transfer": {
		// Data transfer pricing is complex and depends on direction and destination
		// These are base rates, actual calculation is in data_transfer.go
		FreeTierAllowance: &FreeTierAllowance{
			Amount: 1.0,              // First 1GB free per month
			Unit:   "GB",
			Period: pricing.Monthly,
		},
		RegionalVariations: map[string]float64{
			"us-east-1": 1.0,
			"us-west-2": 1.0,
			"eu-west-1": 1.0,
		},
	},
}

// GetNetworkingPricingRates returns the pricing rates for networking resources
func GetNetworkingPricingRates() map[string]AWSPricingRate {
	return NetworkingPricingRates
}

// GetResourcePricingRate returns the pricing rate for a specific resource type
func GetResourcePricingRate(resourceType string) (*AWSPricingRate, bool) {
	rate, exists := NetworkingPricingRates[resourceType]
	return &rate, exists
}

// GetRegionalMultiplier returns the regional pricing multiplier for a resource
func GetRegionalMultiplier(resourceType, region string) float64 {
	rate, exists := NetworkingPricingRates[resourceType]
	if !exists {
		return 1.0
	}
	if multiplier, ok := rate.RegionalVariations[region]; ok {
		return multiplier
	}
	return 1.0
}
