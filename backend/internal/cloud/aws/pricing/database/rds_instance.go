package database

import (
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// RDSInstanceRates contains static pricing rates for AWS RDS instance types
// These are indicative rates (On-Demand, Single-AZ)
var RDSInstanceRates = map[string]float64{
	// Burstable Performance (T3)
	"db.t3.micro":   0.017, // ~$0.017/hr
	"db.t3.small":   0.034, // ~$0.034/hr
	"db.t3.medium":  0.068, // ~$0.068/hr
	"db.t3.large":   0.136, // ~$0.136/hr
	"db.t3.xlarge":  0.272, // ~$0.272/hr
	"db.t3.2xlarge": 0.544, // ~$0.544/hr

	// General Purpose (M5)
	"db.m5.large":   0.176, // ~$0.176/hr
	"db.m5.xlarge":  0.352, // ~$0.352/hr
	"db.m5.2xlarge": 0.704, // ~$0.704/hr
	"db.m5.4xlarge": 1.408, // ~$1.408/hr

	// Memory Optimized (R5)
	"db.r5.large":   0.24, // ~$0.240/hr
	"db.r5.xlarge":  0.48, // ~$0.480/hr
	"db.r5.2xlarge": 0.96, // ~$0.960/hr
}

// RDSStorageRates per GB-month
var RDSStorageRates = map[string]float64{
	"gp2":      0.115, // General Purpose SSD
	"gp3":      0.08,  // General Purpose SSD (gp3)
	"io1":      0.125, // Provisioned IOPS SSD
	"standard": 0.23,  // Magnetic (old) - using approximate
}

// RDSMultiAZMultiplier
const RDSMultiAZMultiplier = 2.0 // Multi-AZ is roughly 2x Single-AZ cost for instance

// GetRDSInstanceRate returns the hourly rate for an RDS instance
func GetRDSInstanceRate(instanceClass, engine string, multiAZ bool, region string) (float64, bool) {
	baseRate, exists := RDSInstanceRates[instanceClass]
	if !exists {
		return 0, false
	}

	// Multiplier logic
	rate := baseRate
	if multiAZ {
		rate *= RDSMultiAZMultiplier
	}

	// Region adjustment (simplified)
	if region == "ap-southeast-1" {
		rate *= 1.1
	}

	return rate, true
}

// CalculateRDSInstanceCost calculates the total cost (Instance + Storage)
func CalculateRDSInstanceCost(duration time.Duration, instanceClass, engine string, multiAZ bool, allocatedStorage float64, storageType string, region string) float64 {
	// Instance Cost
	rate, _ := GetRDSInstanceRate(instanceClass, engine, multiAZ, region)
	instanceCost := rate * duration.Hours()

	return instanceCost
}

// GetRDSInstancePricing returns the detailed pricing structure
func GetRDSInstancePricing(instanceClass, engine string, multiAZ bool, allocatedStorage float64, storageType, region string) *domainpricing.ResourcePricing {
	rate, _ := GetRDSInstanceRate(instanceClass, engine, multiAZ, region)

	storageRate, ok := RDSStorageRates[storageType]
	if !ok {
		storageRate = RDSStorageRates["gp2"]
	}
	if multiAZ {
		storageRate *= 2.0
	}

	components := []domainpricing.PriceComponent{
		{
			Name:        "RDS Instance Hourly",
			Model:       domainpricing.PerHour,
			Unit:        "hour",
			Rate:        rate,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "Hourly charge for RDS DB instance",
		},
	}

	return &domainpricing.ResourcePricing{
		ResourceType: "rds_instance",
		Provider:     domainpricing.AWS,
		Components:   components,
		Metadata: map[string]interface{}{
			"instance_class": instanceClass,
			"engine":         engine,
			"multi_az":       multiAZ,
			"storage_size":   allocatedStorage,
			"storage_type":   storageType,
		},
	}
}
