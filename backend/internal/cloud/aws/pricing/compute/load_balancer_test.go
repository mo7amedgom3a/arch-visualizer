package compute

import (
	"math"
	"testing"
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

const epsilon = 0.0001

func TestCalculateLoadBalancerCost(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		lbType   string
		region   string
		expected float64
	}{
		{
			name:     "ALB-1-hour",
			duration: 1 * time.Hour,
			lbType:   "application",
			region:   "us-east-1",
			expected: 0.0225,
		},
		{
			name:     "ALB-1-month",
			duration: 30 * 24 * time.Hour,
			lbType:   "application",
			region:   "us-east-1",
			expected: 16.2, // 0.0225 * 720
		},
		{
			name:     "NLB-1-hour",
			duration: 1 * time.Hour,
			lbType:   "network",
			region:   "us-east-1",
			expected: 0.0225,
		},
		{
			name:     "CLB-1-hour",
			duration: 1 * time.Hour,
			lbType:   "classic",
			region:   "us-east-1",
			expected: 0.025,
		},
		{
			name:     "ALB-invalid-type",
			duration: 1 * time.Hour,
			lbType:   "invalid",
			region:   "us-east-1",
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := CalculateLoadBalancerCost(tt.duration, tt.lbType, tt.region)
			if math.Abs(cost-tt.expected) > epsilon {
				t.Errorf("Expected cost %.4f, got %.4f", tt.expected, cost)
			}
		})
	}
}

func TestGetLoadBalancerPricing(t *testing.T) {
	tests := []struct {
		name         string
		lbType       string
		region       string
		expectedType string
		expectedRate float64
	}{
		{
			name:         "ALB-pricing",
			lbType:       "application",
			region:       "us-east-1",
			expectedType: "load_balancer",
			expectedRate: 0.0225,
		},
		{
			name:         "NLB-pricing",
			lbType:       "network",
			region:       "us-east-1",
			expectedType: "load_balancer",
			expectedRate: 0.0225,
		},
		{
			name:         "CLB-pricing",
			lbType:       "classic",
			region:       "us-east-1",
			expectedType: "load_balancer",
			expectedRate: 0.025,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetLoadBalancerPricing(tt.lbType, tt.region)

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

			// Verify metadata
			if pricing.Metadata == nil {
				t.Fatal("Expected metadata but got nil")
			}

			if lbType, ok := pricing.Metadata["load_balancer_type"].(string); !ok || lbType != tt.lbType {
				t.Errorf("Expected load_balancer_type %s in metadata, got %v", tt.lbType, pricing.Metadata["load_balancer_type"])
			}
		})
	}
}

func TestGetLoadBalancerRate(t *testing.T) {
	tests := []struct {
		name     string
		lbType   LoadBalancerType
		region   string
		expected float64
		exists   bool
	}{
		{
			name:     "ALB-us-east-1",
			lbType:   LoadBalancerTypeALB,
			region:   "us-east-1",
			expected: 0.0225,
			exists:   true,
		},
		{
			name:     "NLB-us-east-1",
			lbType:   LoadBalancerTypeNLB,
			region:   "us-east-1",
			expected: 0.0225,
			exists:   true,
		},
		{
			name:     "CLB-us-east-1",
			lbType:   LoadBalancerTypeCLB,
			region:   "us-east-1",
			expected: 0.025,
			exists:   true,
		},
		{
			name:     "ALB-ap-southeast-1",
			lbType:   LoadBalancerTypeALB,
			region:   "ap-southeast-1",
			expected: 0.02475, // 0.0225 * 1.1
			exists:   true,
		},
		{
			name:     "Invalid-type",
			lbType:   LoadBalancerType("invalid"),
			region:   "us-east-1",
			expected: 0.0,
			exists:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate, exists := getLoadBalancerRate(tt.lbType, tt.region)
			if exists != tt.exists {
				t.Errorf("Expected exists %v, got %v", tt.exists, exists)
			}
			if math.Abs(rate-tt.expected) > epsilon {
				t.Errorf("Expected rate %.4f, got %.4f", tt.expected, rate)
			}
		})
	}
}
