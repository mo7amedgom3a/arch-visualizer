package storage

import (
	"math"
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// S3StorageRates contains static pricing rates for AWS S3 storage classes (per GB-month)
// These rates are based on AWS public pricing as of 2024
var S3StorageRates = map[string]float64{
	"standard":    0.023,  // $0.023 per GB-month (Standard storage)
	"standard-ia": 0.0125, // $0.0125 per GB-month (Standard-IA storage)
	"glacier":     0.004,  // $0.004 per GB-month (Glacier storage)
}

// S3RequestRates contains static pricing rates for AWS S3 requests (per 1,000 requests)
var S3RequestRates = map[string]float64{
	"PUT": 0.005,   // $0.005 per 1,000 PUT requests
	"GET": 0.0004,  // $0.0004 per 1,000 GET requests
}

// S3DataTransferRate is the outbound data transfer rate (per GB)
// First 1GB/month is free
const S3DataTransferRate = 0.09 // $0.09 per GB (after free tier)

// S3RegionalMultipliers contains regional pricing multipliers for S3
var S3RegionalMultipliers = map[string]float64{
	"us-east-1":      1.0,  // Base rate multiplier
	"us-west-2":      1.0,  // Base rate multiplier
	"eu-west-1":      1.0,  // Base rate multiplier
	"ap-southeast-1": 1.1,  // Slightly higher in some regions
}

// getS3StorageRate returns the per GB-month rate for an S3 storage class
func getS3StorageRate(storageClass, region string) (float64, bool) {
	baseRate, exists := S3StorageRates[storageClass]
	if !exists {
		// Default to standard if not specified
		baseRate = S3StorageRates["standard"]
		exists = true
	}

	multiplier := 1.0
	if m, ok := S3RegionalMultipliers[region]; ok {
		multiplier = m
	}

	return baseRate * multiplier, true
}

// getS3RequestRate returns the per 1,000 requests rate for an S3 request type
func getS3RequestRate(requestType, region string) (float64, bool) {
	baseRate, exists := S3RequestRates[requestType]
	if !exists {
		return 0.0, false
	}

	multiplier := 1.0
	if m, ok := S3RegionalMultipliers[region]; ok {
		multiplier = m
	}

	return baseRate * multiplier, true
}

// getS3DataTransferRate returns the per GB rate for S3 data transfer
func getS3DataTransferRate(region string) (float64, bool) {
	multiplier := 1.0
	if m, ok := S3RegionalMultipliers[region]; ok {
		multiplier = m
	}

	return S3DataTransferRate * multiplier, true
}

// CalculateS3BucketCost calculates the total cost for an S3 bucket
// duration: time duration for the cost calculation
// sizeGB: bucket size in GB
// putRequests: number of PUT requests
// getRequests: number of GET requests
// dataTransferGB: amount of data transferred out in GB
// storageClass: S3 storage class (e.g., "standard", "standard-ia", "glacier")
// region: AWS region
func CalculateS3BucketCost(
	duration time.Duration,
	sizeGB, putRequests, getRequests, dataTransferGB float64,
	storageClass, region string,
) float64 {
	var totalCost float64

	// Calculate storage cost (per GB-month)
	storageRate, exists := getS3StorageRate(storageClass, region)
	if exists {
		// S3 pricing is per GB-month, so we need to prorate based on duration
		// 720 hours = 30 days = 1 month
		hoursPerMonth := 720.0
		months := duration.Hours() / hoursPerMonth
		totalCost += storageRate * sizeGB * months
	}

	// Calculate PUT request cost (per 1,000 requests)
	putRate, exists := getS3RequestRate("PUT", region)
	if exists && putRequests > 0 {
		// Rate is per 1,000 requests
		totalCost += (putRate / 1000.0) * putRequests
	}

	// Calculate GET request cost (per 1,000 requests)
	getRate, exists := getS3RequestRate("GET", region)
	if exists && getRequests > 0 {
		// Rate is per 1,000 requests
		totalCost += (getRate / 1000.0) * getRequests
	}

	// Calculate data transfer cost (per GB, first 1GB/month free)
	dataTransferRate, exists := getS3DataTransferRate(region)
	if exists && dataTransferGB > 0 {
		hoursPerMonth := 720.0
		months := duration.Hours() / hoursPerMonth
		// First 1GB per month is free
		freeTierPerMonth := 1.0
		chargeableGB := math.Max(0, dataTransferGB-(freeTierPerMonth*months))
		totalCost += dataTransferRate * chargeableGB
	}

	return totalCost
}

// GetS3BucketPricing returns the pricing information for S3 buckets
func GetS3BucketPricing(storageClass, region string) *domainpricing.ResourcePricing {
	storageRate, exists := getS3StorageRate(storageClass, region)
	if !exists {
		// Default to standard if not found
		storageRate, _ = getS3StorageRate("standard", region)
	}

	putRate, _ := getS3RequestRate("PUT", region)
	getRate, _ := getS3RequestRate("GET", region)
	dataTransferRate, _ := getS3DataTransferRate(region)

	components := []domainpricing.PriceComponent{
		{
			Name:        "S3 Storage",
			Model:       domainpricing.PerGB,
			Unit:        "GB-month",
			Rate:        storageRate,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "Monthly charge per GB for S3 storage",
		},
		{
			Name:        "S3 PUT Requests",
			Model:       domainpricing.PerRequest,
			Unit:        "per 1,000 requests",
			Rate:        putRate,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "Charge per 1,000 PUT requests",
		},
		{
			Name:        "S3 GET Requests",
			Model:       domainpricing.PerRequest,
			Unit:        "per 1,000 requests",
			Rate:        getRate,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "Charge per 1,000 GET requests",
		},
		{
			Name:        "S3 Data Transfer Out",
			Model:       domainpricing.PerGB,
			Unit:        "GB",
			Rate:        dataTransferRate,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "Charge per GB for outbound data transfer (first 1GB/month free)",
		},
	}

	metadata := map[string]interface{}{
		"storage_class":      storageClass,
		"rate_per_gb_month":  storageRate,
		"put_rate_per_1k":    putRate,
		"get_rate_per_1k":    getRate,
		"data_transfer_rate": dataTransferRate,
	}

	return &domainpricing.ResourcePricing{
		ResourceType: "s3_bucket",
		Provider:     domainpricing.AWS,
		Components:   components,
		Metadata:     metadata,
	}
}
