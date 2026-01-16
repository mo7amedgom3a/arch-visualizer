package networking

import (
	"testing"
	"time"
)

func TestCalculateElasticIPCost(t *testing.T) {
	tests := []struct {
		name         string
		duration     time.Duration
		isAttached   bool
		region       string
		expectedCost float64
	}{
		{
			name:         "elastic-ip-attached-1-hour",
			duration:     1 * time.Hour,
			isAttached:   true,
			region:       "us-east-1",
			expectedCost: 0.0, // Free when attached
		},
		{
			name:         "elastic-ip-unattached-1-hour",
			duration:     1 * time.Hour,
			isAttached:   false,
			region:       "us-east-1",
			expectedCost: 0.005, // $0.005 per hour
		},
		{
			name:         "elastic-ip-attached-720-hours",
			duration:     720 * time.Hour, // 30 days
			isAttached:   true,
			region:       "us-east-1",
			expectedCost: 0.0, // Free when attached
		},
		{
			name:         "elastic-ip-unattached-720-hours",
			duration:     720 * time.Hour, // 30 days
			isAttached:   false,
			region:       "us-east-1",
			expectedCost: 3.60, // $0.005 * 720
		},
		{
			name:         "elastic-ip-unattached-24-hours",
			duration:     24 * time.Hour,
			isAttached:   false,
			region:       "us-east-1",
			expectedCost: 0.12, // $0.005 * 24
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := CalculateElasticIPCost(tt.duration, tt.isAttached, tt.region)
			if cost != tt.expectedCost {
				t.Errorf("Expected cost %.3f, got %.3f", tt.expectedCost, cost)
			}
		})
	}
}

func TestGetElasticIPPricing(t *testing.T) {
	tests := []struct {
		name    string
		region  string
		wantNil bool
	}{
		{
			name:    "get-elastic-ip-pricing-us-east-1",
			region:  "us-east-1",
			wantNil: false,
		},
		{
			name:    "get-elastic-ip-pricing-eu-west-1",
			region:  "eu-west-1",
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetElasticIPPricing(tt.region)
			if tt.wantNil && pricing != nil {
				t.Error("Expected nil pricing")
			}
			if !tt.wantNil && pricing == nil {
				t.Error("Expected non-nil pricing")
			}
			if pricing != nil {
				if pricing.ResourceType != "elastic_ip" {
					t.Errorf("Expected resource type 'elastic_ip', got '%s'", pricing.ResourceType)
				}
				if len(pricing.Components) < 1 {
					t.Error("Expected at least 1 component")
				}
				// Check that it mentions free when attached
				if metadata, ok := pricing.Metadata["free_when_attached"].(bool); !ok || !metadata {
					t.Error("Expected metadata to indicate free when attached")
				}
			}
		})
	}
}
