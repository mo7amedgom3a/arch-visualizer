package pricing

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// AWSPricingRate represents AWS-specific pricing rate information
type AWSPricingRate struct {
	// BaseHourlyRate is the base hourly rate (if applicable)
	BaseHourlyRate float64 `json:"base_hourly_rate,omitempty"`
	// DataProcessingRate is the rate for data processing (per GB)
	DataProcessingRate float64 `json:"data_processing_rate,omitempty"`
	// FreeTierAllowance is the free tier allowance (e.g., first 1GB free)
	FreeTierAllowance *FreeTierAllowance `json:"free_tier_allowance,omitempty"`
	// RegionalVariations contains region-specific pricing adjustments
	RegionalVariations map[string]float64 `json:"regional_variations,omitempty"`
	// AdditionalComponents contains additional pricing components
	AdditionalComponents []pricing.PriceComponent `json:"additional_components,omitempty"`
}

// FreeTierAllowance represents free tier allowances
type FreeTierAllowance struct {
	// Amount is the free amount (e.g., 1.0 for 1GB)
	Amount float64 `json:"amount"`
	// Unit is the unit of the free tier (GB, hours, requests, etc.)
	Unit string `json:"unit"`
	// Period is the period for the free tier (monthly, yearly)
	Period pricing.Period `json:"period"`
}

// AWSResourcePricing extends domain ResourcePricing with AWS-specific fields
type AWSResourcePricing struct {
	// ResourcePricing is the base pricing information
	*pricing.ResourcePricing
	// AWSSpecific contains AWS-specific pricing details
	AWSSpecific *AWSPricingRate `json:"aws_specific,omitempty"`
}
