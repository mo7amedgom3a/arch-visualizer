package compute

import (
	"math"
	"testing"
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

func TestCalculateEC2InstanceCost(t *testing.T) {
	tests := []struct {
		name         string
		duration     time.Duration
		instanceType string
		region       string
		expectedCost float64
	}{
		{
			name:         "t3.micro-1-hour",
			duration:     1 * time.Hour,
			instanceType: "t3.micro",
			region:       "us-east-1",
			expectedCost: 0.0104,
		},
		{
			name:         "t3.micro-720-hours",
			duration:     720 * time.Hour, // 30 days
			instanceType: "t3.micro",
			region:       "us-east-1",
			expectedCost: 7.488, // 0.0104 * 720
		},
		{
			name:         "m5.large-1-hour",
			duration:     1 * time.Hour,
			instanceType: "m5.large",
			region:       "us-east-1",
			expectedCost: 0.096,
		},
		{
			name:         "m5.large-720-hours",
			duration:     720 * time.Hour,
			instanceType: "m5.large",
			region:       "us-east-1",
			expectedCost: 69.12, // 0.096 * 720
		},
		{
			name:         "unknown-instance-type",
			duration:     1 * time.Hour,
			instanceType: "unknown.type",
			region:       "us-east-1",
			expectedCost: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := CalculateEC2InstanceCost(tt.duration, tt.instanceType, tt.region)
			// Use epsilon for floating point comparison
			epsilon := 0.0001
			if math.Abs(cost-tt.expectedCost) > epsilon {
				t.Errorf("Expected cost %.4f, got %.4f", tt.expectedCost, cost)
			}
		})
	}
}

func TestGetEC2InstancePricing(t *testing.T) {
	tests := []struct {
		name         string
		instanceType string
		region       string
		expectedType string
		expectedRate float64
		expectError  bool
	}{
		{
			name:         "t3.micro-pricing",
			instanceType: "t3.micro",
			region:       "us-east-1",
			expectedType: "ec2_instance",
			expectedRate: 0.0104,
			expectError:  false,
		},
		{
			name:         "m5.large-pricing",
			instanceType: "m5.large",
			region:       "us-east-1",
			expectedType: "ec2_instance",
			expectedRate: 0.096,
			expectError:  false,
		},
		{
			name:         "c5.xlarge-pricing",
			instanceType: "c5.xlarge",
			region:       "us-east-1",
			expectedType: "ec2_instance",
			expectedRate: 0.17,
			expectError:  false,
		},
		{
			name:         "unknown-instance-type",
			instanceType: "unknown.type",
			region:       "us-east-1",
			expectedType: "ec2_instance",
			expectedRate: 0.0,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetEC2InstancePricing(tt.instanceType, tt.region)

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
				t.Errorf("Expected rate %.4f, got %.4f", tt.expectedRate, pricing.Components[0].Rate)
			}

			if pricing.Components[0].Model != domainpricing.PerHour {
				t.Errorf("Expected model PerHour, got %s", pricing.Components[0].Model)
			}

			if pricing.Components[0].Currency != domainpricing.USD {
				t.Errorf("Expected currency USD, got %s", pricing.Components[0].Currency)
			}
		})
	}
}
