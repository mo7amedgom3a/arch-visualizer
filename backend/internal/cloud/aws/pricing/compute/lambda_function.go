package compute

import (
	"math"
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// LambdaComputeRate is the per GB-second rate for Lambda compute
// AWS Lambda charges $0.0000166667 per GB-second
const LambdaComputeRate = 0.0000166667 // $0.0000166667 per GB-second

// LambdaRequestRate is the per million requests rate for Lambda invocations
// AWS Lambda charges $0.20 per million requests (first 1M requests/month are free)
const LambdaRequestRate = 0.20 // $0.20 per million requests

// LambdaFreeTierRequests is the number of free requests per month
const LambdaFreeTierRequests = 1000000.0 // 1 million requests per month

// LambdaDataTransferRate is the outbound data transfer rate (per GB)
// First 1GB/month is free
const LambdaDataTransferRate = 0.09 // $0.09 per GB (after free tier)

// LambdaRegionalMultipliers contains regional pricing multipliers for Lambda
var LambdaRegionalMultipliers = map[string]float64{
	"us-east-1":      1.0, // Base rate multiplier
	"us-west-2":      1.0, // Base rate multiplier
	"eu-west-1":      1.0, // Base rate multiplier
	"ap-southeast-1": 1.1, // Slightly higher in some regions
}

// getLambdaComputeRate returns the per GB-second rate for Lambda compute
func getLambdaComputeRate(region string) float64 {
	multiplier := 1.0
	if m, ok := LambdaRegionalMultipliers[region]; ok {
		multiplier = m
	}
	return LambdaComputeRate * multiplier
}

// getLambdaRequestRate returns the per million requests rate for Lambda
func getLambdaRequestRate(region string) float64 {
	multiplier := 1.0
	if m, ok := LambdaRegionalMultipliers[region]; ok {
		multiplier = m
	}
	return LambdaRequestRate * multiplier
}

// getLambdaDataTransferRate returns the per GB rate for Lambda data transfer
func getLambdaDataTransferRate(region string) float64 {
	multiplier := 1.0
	if m, ok := LambdaRegionalMultipliers[region]; ok {
		multiplier = m
	}
	return LambdaDataTransferRate * multiplier
}

// CalculateLambdaFunctionCost calculates the total cost for a Lambda function
// duration: time duration for the cost calculation
// memorySizeMB: allocated memory in MB (e.g., 128, 256, 512)
// averageDurationMs: average execution duration in milliseconds
// requestCount: number of function invocations
// dataTransferGB: amount of data transferred out in GB
// region: AWS region
func CalculateLambdaFunctionCost(
	duration time.Duration,
	memorySizeMB, averageDurationMs, requestCount, dataTransferGB float64,
	region string,
) float64 {
	var totalCost float64

	// Calculate compute cost (GB-seconds)
	// Formula: (memorySizeGB * durationSeconds * requestCount) * computeRate
	memorySizeGB := memorySizeMB / 1024.0 // Convert MB to GB
	durationSeconds := averageDurationMs / 1000.0 // Convert ms to seconds
	
	// Total GB-seconds = memory (GB) * duration (seconds) * requests
	totalGBSeconds := memorySizeGB * durationSeconds * requestCount
	
	computeRate := getLambdaComputeRate(region)
	totalCost += computeRate * totalGBSeconds

	// Calculate request cost (per million requests)
	// First 1M requests per month are free
	requestRate := getLambdaRequestRate(region)
	hoursPerMonth := 720.0
	months := duration.Hours() / hoursPerMonth
	
	// Free tier: 1M requests per month
	freeTierPerMonth := LambdaFreeTierRequests
	freeTierTotal := freeTierPerMonth * months
	chargeableRequests := math.Max(0, requestCount-freeTierTotal)
	
	// Rate is per million requests
	if chargeableRequests > 0 {
		totalCost += (requestRate / 1000000.0) * chargeableRequests
	}

	// Calculate data transfer cost (per GB, first 1GB/month free)
	dataTransferRate := getLambdaDataTransferRate(region)
	if dataTransferGB > 0 {
		// First 1GB per month is free
		freeTierPerMonth := 1.0
		chargeableGB := math.Max(0, dataTransferGB-(freeTierPerMonth*months))
		totalCost += dataTransferRate * chargeableGB
	}

	return totalCost
}

// GetLambdaFunctionPricing returns the pricing information for Lambda functions
func GetLambdaFunctionPricing(memorySizeMB float64, region string) *domainpricing.ResourcePricing {
	computeRate := getLambdaComputeRate(region)
	requestRate := getLambdaRequestRate(region)
	dataTransferRate := getLambdaDataTransferRate(region)

	components := []domainpricing.PriceComponent{
		{
			Name:        "Lambda Compute",
			Model:       domainpricing.PerGB, // Using PerGB model for GB-second pricing
			Unit:        "GB-second",
			Rate:        computeRate,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "Charge per GB-second of compute time",
		},
		{
			Name:        "Lambda Requests",
			Model:       domainpricing.PerRequest,
			Unit:        "per million requests",
			Rate:        requestRate,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "Charge per million requests (first 1M requests/month free)",
		},
		{
			Name:        "Lambda Data Transfer Out",
			Model:       domainpricing.PerGB,
			Unit:        "GB",
			Rate:        dataTransferRate,
			Currency:    domainpricing.USD,
			Region:      &region,
			Description: "Charge per GB for outbound data transfer (first 1GB/month free)",
		},
	}

	metadata := map[string]interface{}{
		"memory_size_mb":        memorySizeMB,
		"compute_rate_per_gbs":  computeRate,
		"request_rate_per_1m":   requestRate,
		"data_transfer_rate":    dataTransferRate,
		"free_tier_requests":    LambdaFreeTierRequests,
	}

	return &domainpricing.ResourcePricing{
		ResourceType: "lambda_function",
		Provider:     domainpricing.AWS,
		Components:   components,
		Metadata:     metadata,
	}
}
