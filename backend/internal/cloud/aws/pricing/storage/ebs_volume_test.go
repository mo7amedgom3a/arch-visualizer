package storage

import (
	"testing"
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

func TestCalculateEBSVolumeCost(t *testing.T) {
	tests := []struct {
		name       string
		duration   time.Duration
		sizeGB     float64
		volumeType string
		region     string
		expectedCost float64
	}{
		{
			name:        "gp3-100gb-1-month",
			duration:    720 * time.Hour, // 30 days = 1 month
			sizeGB:      100.0,
			volumeType:  "gp3",
			region:      "us-east-1",
			expectedCost: 8.0, // 0.08 * 100 * 1
		},
		{
			name:        "gp3-100gb-2-months",
			duration:    1440 * time.Hour, // 60 days = 2 months
			sizeGB:      100.0,
			volumeType:  "gp3",
			region:      "us-east-1",
			expectedCost: 16.0, // 0.08 * 100 * 2
		},
		{
			name:        "gp3-100gb-15-days",
			duration:    360 * time.Hour, // 15 days = 0.5 months
			sizeGB:      100.0,
			volumeType:  "gp3",
			region:      "us-east-1",
			expectedCost: 4.0, // 0.08 * 100 * 0.5
		},
		{
			name:        "gp2-50gb-1-month",
			duration:    720 * time.Hour,
			sizeGB:      50.0,
			volumeType:  "gp2",
			region:      "us-east-1",
			expectedCost: 5.0, // 0.10 * 50 * 1
		},
		{
			name:        "io1-200gb-1-month",
			duration:    720 * time.Hour,
			sizeGB:      200.0,
			volumeType:  "io1",
			region:      "us-east-1",
			expectedCost: 25.0, // 0.125 * 200 * 1
		},
		{
			name:        "unknown-volume-type-defaults-to-gp3",
			duration:    720 * time.Hour,
			sizeGB:      100.0,
			volumeType:  "unknown",
			region:      "us-east-1",
			expectedCost: 8.0, // Defaults to gp3 rate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := CalculateEBSVolumeCost(tt.duration, tt.sizeGB, tt.volumeType, tt.region)
			if cost != tt.expectedCost {
				t.Errorf("Expected cost %.2f, got %.2f", tt.expectedCost, cost)
			}
		})
	}
}

func TestGetEBSVolumePricing(t *testing.T) {
	tests := []struct {
		name          string
		volumeType    string
		region        string
		expectedType  string
		expectedRate  float64
		expectError   bool
	}{
		{
			name:         "gp3-pricing",
			volumeType:   "gp3",
			region:       "us-east-1",
			expectedType: "ebs_volume",
			expectedRate: 0.08,
			expectError:  false,
		},
		{
			name:         "gp2-pricing",
			volumeType:   "gp2",
			region:       "us-east-1",
			expectedType: "ebs_volume",
			expectedRate: 0.10,
			expectError:  false,
		},
		{
			name:         "io1-pricing",
			volumeType:   "io1",
			region:       "us-east-1",
			expectedType: "ebs_volume",
			expectedRate: 0.125,
			expectError:  false,
		},
		{
			name:         "unknown-volume-type-defaults-to-gp3",
			volumeType:   "unknown",
			region:       "us-east-1",
			expectedType: "ebs_volume",
			expectedRate: 0.08, // Defaults to gp3
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetEBSVolumePricing(tt.volumeType, tt.region)

			if pricing == nil {
				t.Fatal("Expected pricing but got nil")
			}

			if pricing.ResourceType != tt.expectedType {
				t.Errorf("Expected resource type %s, got %s", tt.expectedType, pricing.ResourceType)
			}

			if pricing.Provider != domainpricing.AWS {
				t.Errorf("Expected provider AWS, got %s", pricing.Provider)
			}

			if len(pricing.Components) == 0 {
				t.Fatal("Expected at least one pricing component")
			}

			if pricing.Components[0].Rate != tt.expectedRate {
				t.Errorf("Expected rate %.3f, got %.3f", tt.expectedRate, pricing.Components[0].Rate)
			}

			if pricing.Components[0].Model != domainpricing.PerGB {
				t.Errorf("Expected model PerGB, got %s", pricing.Components[0].Model)
			}

			if pricing.Components[0].Currency != domainpricing.USD {
				t.Errorf("Expected currency USD, got %s", pricing.Components[0].Currency)
			}
		})
	}
}
