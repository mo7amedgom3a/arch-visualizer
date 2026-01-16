package networking

import (
	"testing"
	"time"
)

func TestCalculateNATGatewayCost(t *testing.T) {
	tests := []struct {
		name            string
		duration        time.Duration
		dataProcessedGB float64
		region          string
		expectedCost    float64
	}{
		{
			name:            "nat-gateway-1-hour-no-data",
			duration:        1 * time.Hour,
			dataProcessedGB: 0.0,
			region:          "us-east-1",
			expectedCost:    0.045, // $0.045 per hour
		},
		{
			name:            "nat-gateway-24-hours-no-data",
			duration:        24 * time.Hour,
			dataProcessedGB: 0.0,
			region:          "us-east-1",
			expectedCost:    1.08, // $0.045 * 24
		},
		{
			name:            "nat-gateway-720-hours-no-data",
			duration:        720 * time.Hour, // 30 days
			dataProcessedGB: 0.0,
			region:          "us-east-1",
			expectedCost:    32.40, // $0.045 * 720
		},
		{
			name:            "nat-gateway-1-hour-with-data",
			duration:        1 * time.Hour,
			dataProcessedGB: 100.0,
			region:          "us-east-1",
			expectedCost:    4.545, // $0.045 (hourly) + $0.045 * 100 (data)
		},
		{
			name:            "nat-gateway-720-hours-with-data",
			duration:        720 * time.Hour,
			dataProcessedGB: 500.0,
			region:          "us-east-1",
			expectedCost:    54.90, // $0.045 * 720 + $0.045 * 500
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := CalculateNATGatewayCost(tt.duration, tt.dataProcessedGB, tt.region)
			if cost != tt.expectedCost {
				t.Errorf("Expected cost %.2f, got %.2f", tt.expectedCost, cost)
			}
		})
	}
}

func TestGetNATGatewayPricing(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		wantNil  bool
	}{
		{
			name:    "get-nat-gateway-pricing-us-east-1",
			region:  "us-east-1",
			wantNil: false,
		},
		{
			name:    "get-nat-gateway-pricing-us-west-2",
			region:  "us-west-2",
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetNATGatewayPricing(tt.region)
			if tt.wantNil && pricing != nil {
				t.Error("Expected nil pricing")
			}
			if !tt.wantNil && pricing == nil {
				t.Error("Expected non-nil pricing")
			}
			if pricing != nil {
				if pricing.ResourceType != "nat_gateway" {
					t.Errorf("Expected resource type 'nat_gateway', got '%s'", pricing.ResourceType)
				}
				if len(pricing.Components) < 2 {
					t.Errorf("Expected at least 2 components, got %d", len(pricing.Components))
				}
				// Check hourly component
				hourlyFound := false
				for _, comp := range pricing.Components {
					if comp.Name == "NAT Gateway Hourly" {
						hourlyFound = true
						if comp.Rate != 0.045 {
							t.Errorf("Expected hourly rate 0.045, got %f", comp.Rate)
						}
					}
				}
				if !hourlyFound {
					t.Error("Expected to find 'NAT Gateway Hourly' component")
				}
			}
		})
	}
}
