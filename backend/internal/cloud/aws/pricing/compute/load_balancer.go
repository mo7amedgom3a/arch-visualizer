package compute

import (
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// LoadBalancerType represents the type of load balancer
type LoadBalancerType string

const (
	LoadBalancerTypeALB LoadBalancerType = "application" // Application Load Balancer
	LoadBalancerTypeNLB LoadBalancerType = "network"     // Network Load Balancer
	LoadBalancerTypeCLB LoadBalancerType = "classic"    // Classic Load Balancer
)

// LoadBalancerRates contains static pricing rates for AWS Load Balancers
// These rates are based on AWS public pricing as of 2024 (On-Demand pricing)
var LoadBalancerRates = map[LoadBalancerType]float64{
	LoadBalancerTypeALB: 0.0225, // $0.0225 per hour
	LoadBalancerTypeNLB: 0.0225, // $0.0225 per hour
	LoadBalancerTypeCLB: 0.025,   // $0.025 per hour
}

// LoadBalancerRegionalMultipliers contains regional pricing multipliers for Load Balancers
var LoadBalancerRegionalMultipliers = map[string]float64{
	"us-east-1":      1.0, // Base rate multiplier
	"us-west-2":      1.0, // Base rate multiplier
	"eu-west-1":      1.0, // Base rate multiplier
	"ap-southeast-1": 1.1, // Slightly higher in some regions
}

// getLoadBalancerRate returns the hourly rate for a load balancer type
func getLoadBalancerRate(lbType LoadBalancerType, region string) (float64, bool) {
	baseRate, exists := LoadBalancerRates[lbType]
	if !exists {
		return 0, false
	}

	multiplier := 1.0
	if m, ok := LoadBalancerRegionalMultipliers[region]; ok {
		multiplier = m
	}

	return baseRate * multiplier, true
}

// CalculateLoadBalancerCost calculates the cost for a load balancer
// duration: time duration for the cost calculation
// lbType: Load balancer type (application, network, classic)
// region: AWS region
func CalculateLoadBalancerCost(duration time.Duration, lbType, region string) float64 {
	lbTypeEnum := LoadBalancerType(lbType)
	rate, exists := getLoadBalancerRate(lbTypeEnum, region)
	if !exists {
		return 0.0
	}

	// Calculate hourly cost
	hours := duration.Hours()
	return rate * hours
}

// GetLoadBalancerPricing returns the pricing information for Load Balancers
func GetLoadBalancerPricing(lbType, region string) *domainpricing.ResourcePricing {
	lbTypeEnum := LoadBalancerType(lbType)
	rate, exists := getLoadBalancerRate(lbTypeEnum, region)
	if !exists {
		// Return default pricing structure even if LB type not found
		rate = 0.0
	}

	components := []domainpricing.PriceComponent{
		{
			Name:        "Load Balancer Hourly",
			Model:       domainpricing.PerHour,
			Unit:        "hour",
			Rate:        rate,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "On-Demand hourly charge for Load Balancer",
		},
	}

	metadata := map[string]interface{}{
		"load_balancer_type": lbType,
		"hourly_rate":        rate,
		"pricing_model":      "on_demand",
	}

	return &domainpricing.ResourcePricing{
		ResourceType: "load_balancer",
		Provider:     domainpricing.AWS,
		Components:   components,
		Metadata:     metadata,
	}
}
