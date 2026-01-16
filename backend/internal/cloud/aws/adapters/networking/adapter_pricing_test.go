package networking

import (
	"context"
	"testing"
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

func TestAWSNetworkingAdapter_EstimateResourceCost(t *testing.T) {
	// Use the existing mock from adapter_test.go
	mockService := &mockAWSNetworkingService{}
	adapter := NewAWSNetworkingAdapter(mockService)
	ctx := context.Background()

	tests := []struct {
		name         string
		resourceType string
		config       map[string]interface{}
		duration     time.Duration
		expectError  bool
		expectedCost float64
	}{
		{
			name:         "estimate-nat-gateway-cost",
			resourceType: "nat_gateway",
			config: map[string]interface{}{
				"region": "us-east-1",
			},
			duration:     720 * time.Hour,
			expectError:  false,
			expectedCost: 32.40,
		},
		{
			name:         "estimate-elastic-ip-cost",
			resourceType: "elastic_ip",
			config: map[string]interface{}{
				"region": "us-east-1",
			},
			duration:     720 * time.Hour,
			expectError:  false,
			expectedCost: 3.60,
		},
		{
			name:         "estimate-network-interface-cost",
			resourceType: "network_interface",
			config: map[string]interface{}{
				"region": "us-east-1",
			},
			duration:     720 * time.Hour,
			expectError:  false,
			expectedCost: 7.20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimate, err := adapter.EstimateResourceCost(ctx, tt.resourceType, tt.config, tt.duration)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if estimate == nil {
				t.Fatal("Expected cost estimate but got nil")
			}

			if estimate.TotalCost != tt.expectedCost {
				t.Errorf("Expected total cost %.2f, got %.2f", tt.expectedCost, estimate.TotalCost)
			}
		})
	}
}

func TestAWSNetworkingAdapter_GetResourcePricing(t *testing.T) {
	// Use the existing mock from adapter_test.go
	mockService := &mockAWSNetworkingService{}
	adapter := NewAWSNetworkingAdapter(mockService)
	ctx := context.Background()

	tests := []struct {
		name         string
		resourceType string
		region       string
		expectError  bool
		expectedType string
	}{
		{
			name:         "get-nat-gateway-pricing",
			resourceType: "nat_gateway",
			region:       "us-east-1",
			expectError:  false,
			expectedType: "nat_gateway",
		},
		{
			name:         "get-elastic-ip-pricing",
			resourceType: "elastic_ip",
			region:       "us-east-1",
			expectError:  false,
			expectedType: "elastic_ip",
		},
		{
			name:         "get-network-interface-pricing",
			resourceType: "network_interface",
			region:       "us-east-1",
			expectError:  false,
			expectedType: "network_interface",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing, err := adapter.GetResourcePricing(ctx, tt.resourceType, tt.region)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if pricing == nil {
				t.Fatal("Expected pricing but got nil")
			}

			if pricing.ResourceType != tt.expectedType {
				t.Errorf("Expected resource type %s, got %s", tt.expectedType, pricing.ResourceType)
			}

			if pricing.Provider != domainpricing.AWS {
				t.Errorf("Expected provider AWS, got %s", pricing.Provider)
			}
		})
	}
}
