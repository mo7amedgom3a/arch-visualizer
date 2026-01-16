package networking

import (
	"testing"
)

func TestCalculateDataTransferCost(t *testing.T) {
	tests := []struct {
		name         string
		amountGB     float64
		direction    DataTransferDirection
		region       string
		expectedCost float64
	}{
		{
			name:         "inbound-data-transfer-free",
			amountGB:     100.0,
			direction:    Inbound,
			region:       "us-east-1",
			expectedCost: 0.0, // Inbound is always free
		},
		{
			name:         "outbound-data-transfer-within-free-tier",
			amountGB:     0.5,
			direction:    Outbound,
			region:       "us-east-1",
			expectedCost: 0.0, // First 1GB is free
		},
		{
			name:         "outbound-data-transfer-exceeds-free-tier",
			amountGB:     10.0,
			direction:    Outbound,
			region:       "us-east-1",
			expectedCost: 0.81, // (10 - 1) * 0.09
		},
		{
			name:         "outbound-data-transfer-exactly-free-tier",
			amountGB:     1.0,
			direction:    Outbound,
			region:       "us-east-1",
			expectedCost: 0.0, // Exactly 1GB is free
		},
		{
			name:         "inter-az-data-transfer",
			amountGB:     50.0,
			direction:    InterAZ,
			region:       "us-east-1",
			expectedCost: 0.50, // 50 * 0.01
		},
		{
			name:         "intra-region-data-transfer",
			amountGB:     25.0,
			direction:    IntraRegion,
			region:       "us-east-1",
			expectedCost: 0.25, // 25 * 0.01
		},
		{
			name:         "outbound-large-amount",
			amountGB:     1000.0,
			direction:    Outbound,
			region:       "us-east-1",
			expectedCost: 89.91, // (1000 - 1) * 0.09
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := CalculateDataTransferCost(tt.amountGB, tt.direction, tt.region)
			// Use a small epsilon for floating point comparison
			epsilon := 0.01
			diff := cost - tt.expectedCost
			if diff < 0 {
				diff = -diff
			}
			if diff > epsilon {
				t.Errorf("Expected cost %.2f, got %.2f (diff: %.4f)", tt.expectedCost, cost, diff)
			}
		})
	}
}

func TestGetDataTransferPricing(t *testing.T) {
	tests := []struct {
		name    string
		region  string
		wantNil bool
	}{
		{
			name:    "get-data-transfer-pricing-us-east-1",
			region:  "us-east-1",
			wantNil: false,
		},
		{
			name:    "get-data-transfer-pricing-us-west-2",
			region:  "us-west-2",
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetDataTransferPricing(tt.region)
			if tt.wantNil && pricing != nil {
				t.Error("Expected nil pricing")
			}
			if !tt.wantNil && pricing == nil {
				t.Error("Expected non-nil pricing")
			}
			if pricing != nil {
				if pricing.ResourceType != "data_transfer" {
					t.Errorf("Expected resource type 'data_transfer', got '%s'", pricing.ResourceType)
				}
				if len(pricing.Components) < 3 {
					t.Errorf("Expected at least 3 components, got %d", len(pricing.Components))
				}
				// Check that inbound is free
				inboundFound := false
				for _, comp := range pricing.Components {
					if comp.Name == "Data Transfer Inbound" {
						inboundFound = true
						if comp.Rate != 0.0 {
							t.Errorf("Expected inbound rate 0.0, got %f", comp.Rate)
						}
					}
				}
				if !inboundFound {
					t.Error("Expected to find 'Data Transfer Inbound' component")
				}
			}
		})
	}
}
