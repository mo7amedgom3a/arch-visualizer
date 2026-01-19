package compute

import (
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// EC2InstanceRates contains static pricing rates for AWS EC2 instance types
// These rates are based on AWS public pricing as of 2024 (On-Demand pricing)
var EC2InstanceRates = map[string]float64{
	// General Purpose instances
	"t3.micro":   0.0104, // $0.0104 per hour
	"t3.small":   0.0208, // $0.0208 per hour
	"t3.medium":  0.0416, // $0.0416 per hour
	"t3.large":   0.0832, // $0.0832 per hour
	"t3.xlarge":  0.1664, // $0.1664 per hour
	"t3.2xlarge": 0.3328, // $0.3328 per hour
	"m5.large":   0.096,  // $0.096 per hour
	"m5.xlarge":  0.192,  // $0.192 per hour
	"m5.2xlarge": 0.384,  // $0.384 per hour
	"m5.4xlarge": 0.768,  // $0.768 per hour
	// Compute Optimized instances
	"c5.large":   0.085,  // $0.085 per hour
	"c5.xlarge":  0.17,   // $0.17 per hour
	"c5.2xlarge": 0.34,   // $0.34 per hour
	"c5.4xlarge": 0.68,   // $0.68 per hour
	// Memory Optimized instances
	"r5.large":   0.126,  // $0.126 per hour
	"r5.xlarge":  0.252,  // $0.252 per hour
	"r5.2xlarge": 0.504,  // $0.504 per hour
	"r5.4xlarge": 1.008,  // $1.008 per hour
}

// EC2RegionalMultipliers contains regional pricing multipliers for EC2 instances
var EC2RegionalMultipliers = map[string]float64{
	"us-east-1":      1.0,  // Base rate multiplier
	"us-west-2":      1.0,  // Base rate multiplier
	"eu-west-1":      1.0,  // Base rate multiplier
	"ap-southeast-1": 1.1,  // Slightly higher in some regions
}

// getEC2InstanceRate returns the hourly rate for an EC2 instance type
func getEC2InstanceRate(instanceType, region string) (float64, bool) {
	baseRate, exists := EC2InstanceRates[instanceType]
	if !exists {
		return 0, false
	}
	
	multiplier := 1.0
	if m, ok := EC2RegionalMultipliers[region]; ok {
		multiplier = m
	}
	
	return baseRate * multiplier, true
}

// CalculateEC2InstanceCost calculates the cost for an EC2 instance
// duration: time duration for the cost calculation
// instanceType: EC2 instance type (e.g., "t3.micro", "m5.large")
// region: AWS region
func CalculateEC2InstanceCost(duration time.Duration, instanceType, region string) float64 {
	rate, exists := getEC2InstanceRate(instanceType, region)
	if !exists {
		return 0.0
	}

	// Calculate hourly cost
	hours := duration.Hours()
	return rate * hours
}

// GetEC2InstancePricing returns the pricing information for EC2 instances
func GetEC2InstancePricing(instanceType, region string) *domainpricing.ResourcePricing {
	rate, exists := getEC2InstanceRate(instanceType, region)
	if !exists {
		// Return default pricing structure even if instance type not found
		rate = 0.0
	}

	components := []domainpricing.PriceComponent{
		{
			Name:        "EC2 Instance Hourly",
			Model:       domainpricing.PerHour,
			Unit:        "hour",
			Rate:        rate,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "On-Demand hourly charge for EC2 instance",
		},
	}

	metadata := map[string]interface{}{
		"instance_type": instanceType,
		"hourly_rate":    rate,
		"pricing_model":  "on_demand",
	}

	return &domainpricing.ResourcePricing{
		ResourceType: "ec2_instance",
		Provider:     domainpricing.AWS,
		Components:   components,
		Metadata:     metadata,
	}
}
