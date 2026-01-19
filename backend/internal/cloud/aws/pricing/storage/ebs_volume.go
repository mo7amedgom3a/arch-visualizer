package storage

import (
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// EBSVolumeRates contains static pricing rates for AWS EBS volume types
// These rates are based on AWS public pricing as of 2024 (per GB-month)
var EBSVolumeRates = map[string]float64{
	"gp2":      0.10,  // $0.10 per GB-month
	"gp3":      0.08,  // $0.08 per GB-month
	"io1":      0.125, // $0.125 per GB-month
	"io2":      0.125, // $0.125 per GB-month
	"sc1":      0.015, // $0.015 per GB-month
	"st1":      0.045, // $0.045 per GB-month
	"standard": 0.05,  // $0.05 per GB-month
}

// EBSRegionalMultipliers contains regional pricing multipliers for EBS volumes
var EBSRegionalMultipliers = map[string]float64{
	"us-east-1":      1.0,  // Base rate multiplier
	"us-west-2":      1.0,  // Base rate multiplier
	"eu-west-1":      1.0,  // Base rate multiplier
	"ap-southeast-1": 1.1,  // Slightly higher in some regions
}

// getEBSVolumeRate returns the per GB-month rate for an EBS volume type
func getEBSVolumeRate(volumeType, region string) (float64, bool) {
	baseRate, exists := EBSVolumeRates[volumeType]
	if !exists {
		// Default to gp3 if not specified
		baseRate = EBSVolumeRates["gp3"]
		exists = true
	}
	
	multiplier := 1.0
	if m, ok := EBSRegionalMultipliers[region]; ok {
		multiplier = m
	}
	
	return baseRate * multiplier, true
}

// CalculateEBSVolumeCost calculates the cost for an EBS volume
// duration: time duration for the cost calculation
// sizeGB: volume size in GB
// volumeType: EBS volume type (e.g., "gp3", "gp2", "io1")
// region: AWS region
func CalculateEBSVolumeCost(duration time.Duration, sizeGB float64, volumeType, region string) float64 {
	rate, exists := getEBSVolumeRate(volumeType, region)
	if !exists {
		return 0.0
	}

	// EBS pricing is per GB-month, so we need to prorate based on duration
	// 720 hours = 30 days = 1 month
	hoursPerMonth := 720.0
	months := duration.Hours() / hoursPerMonth

	return rate * sizeGB * months
}

// GetEBSVolumePricing returns the pricing information for EBS volumes
func GetEBSVolumePricing(volumeType, region string) *domainpricing.ResourcePricing {
	rate, exists := getEBSVolumeRate(volumeType, region)
	if !exists {
		// Default to gp3 if not found
		rate, _ = getEBSVolumeRate("gp3", region)
	}

	components := []domainpricing.PriceComponent{
		{
			Name:        "EBS Volume Storage",
			Model:       domainpricing.PerGB,
			Unit:        "GB-month",
			Rate:        rate,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "Monthly charge per GB for EBS volume storage",
		},
	}

	metadata := map[string]interface{}{
		"volume_type":   volumeType,
		"rate_per_gb_month": rate,
	}

	return &domainpricing.ResourcePricing{
		ResourceType: "ebs_volume",
		Provider:     domainpricing.AWS,
		Components:   components,
		Metadata:     metadata,
	}
}
