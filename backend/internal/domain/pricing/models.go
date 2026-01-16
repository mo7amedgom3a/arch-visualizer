package pricing

import (
	"time"
)

// PriceComponent represents a single pricing component for a resource
// Each resource can have multiple price components (e.g., base price + data transfer)
type PriceComponent struct {
	// Name of the component (e.g., "NAT Gateway Hourly", "Data Transfer Out")
	Name string `json:"name"`
	// Model is the pricing model type (PerHour, PerGB, PerRequest, etc.)
	Model PricingModel `json:"model"`
	// Unit is the unit of measurement (hour, GB, request, etc.)
	Unit string `json:"unit"`
	// Rate is the price per unit (e.g., 0.045 for $0.045/hour)
	Rate float64 `json:"rate"`
	// Currency is the currency code (USD, EUR, etc.)
	Currency Currency `json:"currency"`
	// Region is optional, for regional pricing variations
	Region *string `json:"region,omitempty"`
	// Description provides additional context about the component
	Description string `json:"description,omitempty"`
}

// ResourcePricing represents the pricing information for a specific resource type
type ResourcePricing struct {
	// ResourceType is the type of resource (vpc, nat_gateway, elastic_ip, etc.)
	ResourceType string `json:"resource_type"`
	// Provider is the cloud provider (aws, azure, gcp)
	Provider CloudProvider `json:"provider"`
	// Components is a list of price components for this resource
	Components []PriceComponent `json:"components"`
	// Metadata contains additional provider-specific pricing information
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CostComponent represents a calculated cost component in a cost estimate
type CostComponent struct {
	// ComponentName is the name of the cost component
	ComponentName string `json:"component_name"`
	// Model is the pricing model used
	Model PricingModel `json:"model"`
	// Quantity is the amount consumed (e.g., 720 hours, 100 GB)
	Quantity float64 `json:"quantity"`
	// UnitRate is the rate per unit
	UnitRate float64 `json:"unit_rate"`
	// Subtotal is the calculated cost for this component
	Subtotal float64 `json:"subtotal"`
	// Currency is the currency used
	Currency Currency `json:"currency"`
}

// CostEstimate represents a calculated cost estimate for a resource or architecture
type CostEstimate struct {
	// TotalCost is the total estimated cost
	TotalCost float64 `json:"total_cost"`
	// Currency is the currency used
	Currency Currency `json:"currency"`
	// Breakdown is a detailed breakdown of cost components
	Breakdown []CostComponent `json:"breakdown"`
	// Period is the time period for the estimate (hourly, monthly, yearly)
	Period Period `json:"period"`
	// Duration is the duration used for the calculation
	Duration time.Duration `json:"duration"`
	// CalculatedAt is the timestamp when the estimate was calculated
	CalculatedAt time.Time `json:"calculated_at"`
	// ResourceType is the type of resource (if applicable)
	ResourceType *string `json:"resource_type,omitempty"`
	// Provider is the cloud provider
	Provider CloudProvider `json:"provider"`
	// Region is the region (if applicable)
	Region *string `json:"region,omitempty"`
}
