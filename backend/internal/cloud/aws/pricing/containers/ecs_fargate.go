package containers

import (
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// Fargate pricing rates (as of 2024)
// Pricing is per vCPU-hour and per GB-hour for memory

// FargatevCPURates contains per-vCPU-hour pricing by region
var FargatevCPURates = map[string]float64{
	"us-east-1":      0.04048, // $0.04048 per vCPU per hour
	"us-west-2":      0.04048,
	"eu-west-1":      0.04656,
	"ap-southeast-1": 0.05056,
}

// FargateMemoryRates contains per-GB-hour pricing by region
var FargateMemoryRates = map[string]float64{
	"us-east-1":      0.004445, // $0.004445 per GB per hour
	"us-west-2":      0.004445,
	"eu-west-1":      0.005111,
	"ap-southeast-1": 0.005556,
}

// FargateSpotDiscount is the approximate discount for Fargate Spot (up to 70%)
const FargateSpotDiscount = 0.70

// getFargateRates returns the vCPU and memory rates for a region
func getFargateRates(region string) (vcpuRate, memoryRate float64) {
	vcpuRate, exists := FargatevCPURates[region]
	if !exists {
		vcpuRate = FargatevCPURates["us-east-1"] // Default to us-east-1
	}

	memoryRate, exists = FargateMemoryRates[region]
	if !exists {
		memoryRate = FargateMemoryRates["us-east-1"]
	}

	return vcpuRate, memoryRate
}

// CalculateFargateCost calculates the cost for Fargate tasks
// vcpu: number of vCPUs (e.g., 0.25, 0.5, 1, 2, 4)
// memoryGB: memory in GB (e.g., 0.5, 1, 2, 4, 8)
// duration: time duration for the cost calculation
// region: AWS region
// spot: whether using Fargate Spot
func CalculateFargateCost(vcpu, memoryGB float64, duration time.Duration, region string, spot bool) float64 {
	vcpuRate, memoryRate := getFargateRates(region)
	hours := duration.Hours()

	vcpuCost := vcpu * vcpuRate * hours
	memoryCost := memoryGB * memoryRate * hours
	totalCost := vcpuCost + memoryCost

	if spot {
		totalCost *= (1 - FargateSpotDiscount)
	}

	return totalCost
}

// CalculateFargateMonthlyCost calculates the monthly cost for Fargate tasks
func CalculateFargateMonthlyCost(vcpu, memoryGB float64, region string, spot bool) float64 {
	// 730 hours per month (average)
	return CalculateFargateCost(vcpu, memoryGB, 730*time.Hour, region, spot)
}

// GetFargatePricing returns the pricing information for Fargate
func GetFargatePricing(vcpu, memoryGB float64, region string, spot bool) *domainpricing.ResourcePricing {
	vcpuRate, memoryRate := getFargateRates(region)

	if spot {
		vcpuRate *= (1 - FargateSpotDiscount)
		memoryRate *= (1 - FargateSpotDiscount)
	}

	pricingModel := "on_demand"
	if spot {
		pricingModel = "spot"
	}

	components := []domainpricing.PriceComponent{
		{
			Name:        "Fargate vCPU",
			Model:       domainpricing.PerHour,
			Unit:        "vCPU-hour",
			Rate:        vcpuRate * vcpu,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "Fargate vCPU per hour charge",
		},
		{
			Name:        "Fargate Memory",
			Model:       domainpricing.PerHour,
			Unit:        "GB-hour",
			Rate:        memoryRate * memoryGB,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "Fargate Memory per GB per hour charge",
		},
	}

	hourlyTotal := (vcpu * vcpuRate) + (memoryGB * memoryRate)
	monthlyEstimate := hourlyTotal * 730

	metadata := map[string]interface{}{
		"vcpu":             vcpu,
		"memory_gb":        memoryGB,
		"pricing_model":    pricingModel,
		"hourly_total":     hourlyTotal,
		"monthly_estimate": monthlyEstimate,
	}

	return &domainpricing.ResourcePricing{
		ResourceType: "ecs_fargate",
		Provider:     domainpricing.AWS,
		Components:   components,
		Metadata:     metadata,
	}
}

// GetECSTaskDefinitionPricing returns pricing based on task definition configuration
func GetECSTaskDefinitionPricing(cpu, memory string, region string, launchType string) *domainpricing.ResourcePricing {
	// Parse CPU and memory from task definition format
	// CPU is in "cpu units" (1024 = 1 vCPU)
	// Memory is in MB for EC2, or specific values for Fargate
	vcpu := parseCPUToVCPU(cpu)
	memoryGB := parseMemoryToGB(memory)

	if launchType == "FARGATE" {
		return GetFargatePricing(vcpu, memoryGB, region, false)
	}

	// EC2 launch type - task definition doesn't incur ECS charges
	// Cost is based on underlying EC2 instances
	return &domainpricing.ResourcePricing{
		ResourceType: "ecs_task_definition",
		Provider:     domainpricing.AWS,
		Components:   []domainpricing.PriceComponent{},
		Metadata: map[string]interface{}{
			"launch_type": launchType,
			"note":        "EC2 launch type: no additional ECS charges, uses EC2 instance pricing",
		},
	}
}

// parseCPUToVCPU converts ECS CPU units to vCPU
func parseCPUToVCPU(cpu string) float64 {
	cpuMap := map[string]float64{
		"256":  0.25,
		"512":  0.5,
		"1024": 1.0,
		"2048": 2.0,
		"4096": 4.0,
		"8192": 8.0,
	}
	if vcpu, ok := cpuMap[cpu]; ok {
		return vcpu
	}
	return 0.25 // Default to minimum
}

// parseMemoryToGB converts ECS memory string to GB
func parseMemoryToGB(memory string) float64 {
	memMap := map[string]float64{
		"512":   0.5,
		"1024":  1.0,
		"2048":  2.0,
		"3072":  3.0,
		"4096":  4.0,
		"5120":  5.0,
		"6144":  6.0,
		"7168":  7.0,
		"8192":  8.0,
		"16384": 16.0,
		"30720": 30.0,
	}
	if gb, ok := memMap[memory]; ok {
		return gb
	}
	return 0.5 // Default to minimum
}
