package networking

import (
	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// DataTransferDirection represents the direction of data transfer
type DataTransferDirection string

const (
	// Inbound data transfer (from internet to AWS)
	Inbound DataTransferDirection = "inbound"
	// Outbound data transfer (from AWS to internet)
	Outbound DataTransferDirection = "outbound"
	// InterAZ data transfer (between availability zones)
	InterAZ DataTransferDirection = "inter_az"
	// IntraRegion data transfer (within same region, different AZs)
	IntraRegion DataTransferDirection = "intra_region"
)

// DataTransferPricing contains data transfer pricing rates
var DataTransferPricing = map[DataTransferDirection]float64{
	Inbound:    0.0,  // Free
	Outbound:   0.09, // $0.09 per GB (after free tier)
	InterAZ:    0.01, // $0.01 per GB
	IntraRegion: 0.01, // $0.01 per GB
}

// FreeTierDataTransfer is the free tier allowance for data transfer (per month)
const FreeTierDataTransfer = 1.0 // 1 GB per month

// CalculateDataTransferCost calculates the cost for data transfer
// amountGB: amount of data in GB
// direction: direction of data transfer
// region: AWS region
func CalculateDataTransferCost(amountGB float64, direction DataTransferDirection, region string) float64 {
	// Inbound is always free
	if direction == Inbound {
		return 0.0
	}

	// Get base rate for direction
	rate, exists := DataTransferPricing[direction]
	if !exists {
		rate = 0.09 // Default to outbound rate
	}

	// Apply regional multiplier (default to 1.0 for now)
	// In a full implementation, this would fetch from a rates map
	multiplier := 1.0
	rate = rate * multiplier

	// Apply free tier (only for outbound)
	if direction == Outbound && amountGB > FreeTierDataTransfer {
		// First 1GB is free
		chargeableAmount := amountGB - FreeTierDataTransfer
		return chargeableAmount * rate
	} else if direction == Outbound {
		// Within free tier
		return 0.0
	}

	// For inter-AZ and intra-region, no free tier
	return amountGB * rate
}

// GetDataTransferPricing returns the pricing information for data transfer
func GetDataTransferPricing(region string) *domainpricing.ResourcePricing {
	components := []domainpricing.PriceComponent{
		{
			Name:     "Data Transfer Inbound",
			Model:    domainpricing.PerGB,
			Unit:     "GB",
			Rate:     0.0,
			Currency: domainpricing.USD,
			Region:   &region,
			Description: "Inbound data transfer is free",
		},
		{
			Name:     "Data Transfer Outbound",
			Model:    domainpricing.PerGB,
			Unit:     "GB",
			Rate:     0.09,
			Currency: domainpricing.USD,
			Region:   &region,
			Description: "First 1GB per month free, then $0.09/GB",
		},
		{
			Name:     "Data Transfer Inter-AZ",
			Model:    domainpricing.PerGB,
			Unit:     "GB",
			Rate:     0.01,
			Currency: domainpricing.USD,
			Region:   &region,
			Description: "Data transfer between availability zones",
		},
	}

	return &domainpricing.ResourcePricing{
		ResourceType: "data_transfer",
		Provider:     domainpricing.AWS,
		Components:  components,
		Metadata: map[string]interface{}{
			"free_tier_gb": FreeTierDataTransfer,
			"free_tier_period": "monthly",
		},
	}
}
