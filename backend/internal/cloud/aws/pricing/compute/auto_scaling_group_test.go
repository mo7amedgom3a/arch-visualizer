package compute

import (
	"math"
	"testing"
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

func TestCalculateAutoScalingGroupCost(t *testing.T) {
	tests := []struct {
		name         string
		duration     time.Duration
		instanceType string
		minSize      int
		maxSize      int
		region       string
		expected     float64
	}{
		{
			name:         "t3.micro-1-hour-min1-max3",
			duration:     1 * time.Hour,
			instanceType: "t3.micro",
			minSize:      1,
			maxSize:      3,
			region:       "us-east-1",
			expected:     0.0208, // (1+3)/2 * 0.0104
		},
		{
			name:         "t3.micro-1-month-min1-max3",
			duration:     30 * 24 * time.Hour,
			instanceType: "t3.micro",
			minSize:      1,
			maxSize:      3,
			region:       "us-east-1",
			expected:     14.976, // 0.0208 * 720
		},
		{
			name:         "m5.large-1-hour-min2-max5",
			duration:     1 * time.Hour,
			instanceType: "m5.large",
			minSize:      2,
			maxSize:      5,
			region:       "us-east-1",
			expected:     0.336, // (2+5)/2 * 0.096
		},
		{
			name:         "t3.micro-zero-capacity",
			duration:     1 * time.Hour,
			instanceType: "t3.micro",
			minSize:      0,
			maxSize:      0,
			region:       "us-east-1",
			expected:     0.0,
		},
		{
			name:         "t3.micro-same-min-max",
			duration:     1 * time.Hour,
			instanceType: "t3.micro",
			minSize:      2,
			maxSize:      2,
			region:       "us-east-1",
			expected:     0.0208, // 2 * 0.0104
		},
		{
			name:         "invalid-instance-type",
			duration:     1 * time.Hour,
			instanceType: "invalid.type",
			minSize:      1,
			maxSize:      3,
			region:       "us-east-1",
			expected:     0.0,
		},
	}

	epsilon := 0.0001
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := CalculateAutoScalingGroupCost(tt.duration, tt.instanceType, tt.minSize, tt.maxSize, tt.region)
			if math.Abs(cost-tt.expected) > epsilon {
				t.Errorf("Expected cost %.4f, got %.4f", tt.expected, cost)
			}
		})
	}
}

func TestGetAutoScalingGroupPricing(t *testing.T) {
	tests := []struct {
		name         string
		instanceType string
		minSize      int
		maxSize      int
		region       string
		expectedType string
		expectedRate float64
	}{
		{
			name:         "t3.micro-min1-max3",
			instanceType: "t3.micro",
			minSize:      1,
			maxSize:      3,
			region:       "us-east-1",
			expectedType: "auto_scaling_group",
			expectedRate: 0.0208, // (1+3)/2 * 0.0104
		},
		{
			name:         "m5.large-min2-max5",
			instanceType: "m5.large",
			minSize:      2,
			maxSize:      5,
			region:       "us-east-1",
			expectedType: "auto_scaling_group",
			expectedRate: 0.336, // (2+5)/2 * 0.096
		},
		{
			name:         "c5.xlarge-min1-max10",
			instanceType: "c5.xlarge",
			minSize:      1,
			maxSize:      10,
			region:       "us-east-1",
			expectedType: "auto_scaling_group",
			expectedRate: 0.935, // (1+10)/2 * 0.17
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetAutoScalingGroupPricing(tt.instanceType, tt.minSize, tt.maxSize, tt.region)

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

			epsilon := 0.0001
			if math.Abs(pricing.Components[0].Rate-tt.expectedRate) > epsilon {
				t.Errorf("Expected rate %.4f, got %.4f", tt.expectedRate, pricing.Components[0].Rate)
			}

			if pricing.Components[0].Model != domainpricing.PerHour {
				t.Errorf("Expected model PerHour, got %s", pricing.Components[0].Model)
			}

			if pricing.Components[0].Currency != domainpricing.USD {
				t.Errorf("Expected currency USD, got %s", pricing.Components[0].Currency)
			}

			// Verify metadata
			if pricing.Metadata == nil {
				t.Fatal("Expected metadata but got nil")
			}

			if instanceType, ok := pricing.Metadata["instance_type"].(string); !ok || instanceType != tt.instanceType {
				t.Errorf("Expected instance_type %s in metadata, got %v", tt.instanceType, pricing.Metadata["instance_type"])
			}

			if minSize, ok := pricing.Metadata["min_size"].(int); !ok || minSize != tt.minSize {
				t.Errorf("Expected min_size %d in metadata, got %v", tt.minSize, pricing.Metadata["min_size"])
			}

			if maxSize, ok := pricing.Metadata["max_size"].(int); !ok || maxSize != tt.maxSize {
				t.Errorf("Expected max_size %d in metadata, got %v", tt.maxSize, pricing.Metadata["max_size"])
			}

			expectedAvgCapacity := float64(tt.minSize+tt.maxSize) / 2.0
			if avgCapacity, ok := pricing.Metadata["average_capacity"].(float64); !ok || math.Abs(avgCapacity-expectedAvgCapacity) > epsilon {
				t.Errorf("Expected average_capacity %.2f in metadata, got %v", expectedAvgCapacity, pricing.Metadata["average_capacity"])
			}
		})
	}
}
