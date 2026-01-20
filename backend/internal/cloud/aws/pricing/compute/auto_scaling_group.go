package compute

import (
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// CalculateAutoScalingGroupCost calculates the cost for an Auto Scaling Group
// duration: time duration for the cost calculation
// instanceType: EC2 instance type (e.g., "t3.micro", "m5.large")
// minSize: minimum number of instances in the ASG
// maxSize: maximum number of instances in the ASG
// region: AWS region
// Note: ASG itself has no cost, only the instances it manages
// We calculate based on average capacity: (minSize + maxSize) / 2
func CalculateAutoScalingGroupCost(duration time.Duration, instanceType string, minSize, maxSize int, region string) float64 {
	// Calculate average capacity
	avgCapacity := float64(minSize+maxSize) / 2.0

	// Get instance hourly rate
	instanceRate, exists := getEC2InstanceRate(instanceType, region)
	if !exists {
		return 0.0
	}

	// Calculate hourly cost: average capacity * instance hourly rate * hours
	hours := duration.Hours()
	return avgCapacity * instanceRate * hours
}

// GetAutoScalingGroupPricing returns the pricing information for Auto Scaling Groups
// instanceType: EC2 instance type used by the ASG
// minSize: minimum number of instances
// maxSize: maximum number of instances
// region: AWS region
func GetAutoScalingGroupPricing(instanceType string, minSize, maxSize int, region string) *domainpricing.ResourcePricing {
	// Get instance hourly rate
	instanceRate, exists := getEC2InstanceRate(instanceType, region)
	if !exists {
		// Return default pricing structure even if instance type not found
		instanceRate = 0.0
	}

	// Calculate average capacity
	avgCapacity := float64(minSize+maxSize) / 2.0

	// Calculate effective hourly rate (average capacity * instance rate)
	effectiveHourlyRate := avgCapacity * instanceRate

	components := []domainpricing.PriceComponent{
		{
			Name:        "Auto Scaling Group Hourly",
			Model:       domainpricing.PerHour,
			Unit:        "hour",
			Rate:        effectiveHourlyRate,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "Hourly charge based on average capacity (min+max)/2 and instance type",
		},
	}

	metadata := map[string]interface{}{
		"instance_type":      instanceType,
		"min_size":           minSize,
		"max_size":           maxSize,
		"average_capacity":   avgCapacity,
		"instance_hourly_rate": instanceRate,
		"effective_hourly_rate": effectiveHourlyRate,
		"pricing_model":      "on_demand",
		"note":               "ASG itself has no cost; pricing is for managed EC2 instances",
	}

	return &domainpricing.ResourcePricing{
		ResourceType: "auto_scaling_group",
		Provider:     domainpricing.AWS,
		Components:   components,
		Metadata:     metadata,
	}
}
