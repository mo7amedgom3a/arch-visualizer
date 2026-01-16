package pricing

import (
	"testing"
	"time"
)

func TestPriceComponent(t *testing.T) {
	tests := []struct {
		name      string
		component PriceComponent
		wantValid bool
	}{
		{
			name: "valid-per-hour-component",
			component: PriceComponent{
				Name:      "NAT Gateway Hourly",
				Model:     PerHour,
				Unit:      "hour",
				Rate:      0.045,
				Currency:  USD,
				Description: "Base hourly charge",
			},
			wantValid: true,
		},
		{
			name: "valid-per-gb-component",
			component: PriceComponent{
				Name:      "Data Transfer Outbound",
				Model:     PerGB,
				Unit:      "GB",
				Rate:      0.09,
				Currency:  USD,
				Description: "Outbound data transfer",
			},
			wantValid: true,
		},
		{
			name: "valid-per-request-component",
			component: PriceComponent{
				Name:      "API Gateway Request",
				Model:     PerRequest,
				Unit:      "request",
				Rate:      0.000001,
				Currency:  USD,
				Description: "Per request charge",
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.component.Name == "" {
				t.Error("PriceComponent should have a name")
			}
			if tt.component.Rate < 0 {
				t.Error("PriceComponent rate should be non-negative")
			}
			if tt.component.Currency == "" {
				t.Error("PriceComponent should have a currency")
			}
		})
	}
}

func TestResourcePricing(t *testing.T) {
	tests := []struct {
		name     string
		pricing  ResourcePricing
		wantValid bool
	}{
		{
			name: "valid-nat-gateway-pricing",
			pricing: ResourcePricing{
				ResourceType: "nat_gateway",
				Provider:     AWS,
				Components: []PriceComponent{
					{
						Name:     "NAT Gateway Hourly",
						Model:    PerHour,
						Unit:     "hour",
						Rate:     0.045,
						Currency: USD,
					},
					{
						Name:     "NAT Gateway Data Processing",
						Model:    PerGB,
						Unit:     "GB",
						Rate:     0.045,
						Currency: USD,
					},
				},
			},
			wantValid: true,
		},
		{
			name: "valid-elastic-ip-pricing",
			pricing: ResourcePricing{
				ResourceType: "elastic_ip",
				Provider:     AWS,
				Components: []PriceComponent{
					{
						Name:     "Elastic IP Hourly (Unattached)",
						Model:    PerHour,
						Unit:     "hour",
						Rate:     0.005,
						Currency: USD,
					},
				},
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.pricing.ResourceType == "" {
				t.Error("ResourcePricing should have a resource type")
			}
			if tt.pricing.Provider == "" {
				t.Error("ResourcePricing should have a provider")
			}
			if len(tt.pricing.Components) == 0 {
				t.Error("ResourcePricing should have at least one component")
			}
		})
	}
}

func TestCostEstimate(t *testing.T) {
	now := time.Now()
	estimate := CostEstimate{
		TotalCost:    32.40,
		Currency:     USD,
		Breakdown: []CostComponent{
			{
				ComponentName: "NAT Gateway Hourly",
				Model:         PerHour,
				Quantity:      720.0,
				UnitRate:      0.045,
				Subtotal:      32.40,
				Currency:      USD,
			},
		},
		Period:       Monthly,
		Duration:     720 * time.Hour,
		CalculatedAt: now,
		ResourceType: stringPtr("nat_gateway"),
		Provider:     AWS,
		Region:       stringPtr("us-east-1"),
	}

	if estimate.TotalCost != 32.40 {
		t.Errorf("Expected TotalCost 32.40, got %f", estimate.TotalCost)
	}
	if len(estimate.Breakdown) != 1 {
		t.Errorf("Expected 1 breakdown component, got %d", len(estimate.Breakdown))
	}
	if estimate.Currency != USD {
		t.Errorf("Expected currency USD, got %s", estimate.Currency)
	}
}

func stringPtr(s string) *string {
	return &s
}
