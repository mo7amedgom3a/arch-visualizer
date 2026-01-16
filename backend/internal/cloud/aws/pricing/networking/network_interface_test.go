package networking

import (
	"testing"
	"time"
)

func TestCalculateNetworkInterfaceCost(t *testing.T) {
	tests := []struct {
		name         string
		duration     time.Duration
		isAttached   bool
		region       string
		expectedCost float64
	}{
		{
			name:         "network-interface-attached-1-hour",
			duration:     1 * time.Hour,
			isAttached:   true,
			region:       "us-east-1",
			expectedCost: 0.0, // Free when attached
		},
		{
			name:         "network-interface-unattached-1-hour",
			duration:     1 * time.Hour,
			isAttached:   false,
			region:       "us-east-1",
			expectedCost: 0.01, // $0.01 per hour
		},
		{
			name:         "network-interface-attached-720-hours",
			duration:     720 * time.Hour, // 30 days
			isAttached:   true,
			region:       "us-east-1",
			expectedCost: 0.0, // Free when attached
		},
		{
			name:         "network-interface-unattached-720-hours",
			duration:     720 * time.Hour, // 30 days
			isAttached:   false,
			region:       "us-east-1",
			expectedCost: 7.20, // $0.01 * 720
		},
		{
			name:         "network-interface-unattached-24-hours",
			duration:     24 * time.Hour,
			isAttached:   false,
			region:       "us-east-1",
			expectedCost: 0.24, // $0.01 * 24
		},
		{
			name:         "network-interface-unattached-168-hours",
			duration:     168 * time.Hour, // 1 week
			isAttached:   false,
			region:       "us-east-1",
			expectedCost: 1.68, // $0.01 * 168
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := CalculateNetworkInterfaceCost(tt.duration, tt.isAttached, tt.region)
			if cost != tt.expectedCost {
				t.Errorf("Expected cost %.2f, got %.2f", tt.expectedCost, cost)
			}
		})
	}
}

func TestGetNetworkInterfacePricing(t *testing.T) {
	tests := []struct {
		name    string
		region  string
		wantNil bool
	}{
		{
			name:    "get-network-interface-pricing-us-east-1",
			region:  "us-east-1",
			wantNil: false,
		},
		{
			name:    "get-network-interface-pricing-us-west-2",
			region:  "us-west-2",
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetNetworkInterfacePricing(tt.region)
			if tt.wantNil && pricing != nil {
				t.Error("Expected nil pricing")
			}
			if !tt.wantNil && pricing == nil {
				t.Error("Expected non-nil pricing")
			}
			if pricing != nil {
				if pricing.ResourceType != "network_interface" {
					t.Errorf("Expected resource type 'network_interface', got '%s'", pricing.ResourceType)
				}
				if len(pricing.Components) < 1 {
					t.Error("Expected at least 1 component")
				}
				// Check that it mentions free when attached
				if metadata, ok := pricing.Metadata["free_when_attached"].(bool); !ok || !metadata {
					t.Error("Expected metadata to indicate free when attached")
				}
				// Check hourly component
				hourlyFound := false
				for _, comp := range pricing.Components {
					if comp.Name == "Network Interface Hourly (Unattached)" {
						hourlyFound = true
						if comp.Rate != 0.01 {
							t.Errorf("Expected hourly rate 0.01, got %f", comp.Rate)
						}
					}
				}
				if !hourlyFound {
					t.Error("Expected to find 'Network Interface Hourly (Unattached)' component")
				}
			}
		})
	}
}
