package pricing

import (
	"context"
	"testing"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

func TestAWSPricingService_GetPricing(t *testing.T) {
	service := NewAWSPricingService()
	ctx := context.Background()

	tests := []struct {
		name          string
		resourceType  string
		provider      string
		region        string
		expectError   bool
		expectedType  string
	}{
		{
			name:         "get-nat-gateway-pricing",
			resourceType: "nat_gateway",
			provider:     "aws",
			region:       "us-east-1",
			expectError:  false,
			expectedType: "nat_gateway",
		},
		{
			name:         "get-elastic-ip-pricing",
			resourceType: "elastic_ip",
			provider:     "aws",
			region:       "us-east-1",
			expectError:  false,
			expectedType: "elastic_ip",
		},
		{
			name:         "get-network-interface-pricing",
			resourceType: "network_interface",
			provider:     "aws",
			region:       "us-east-1",
			expectError:  false,
			expectedType: "network_interface",
		},
		{
			name:         "get-data-transfer-pricing",
			resourceType: "data_transfer",
			provider:     "aws",
			region:       "us-east-1",
			expectError:  false,
			expectedType: "data_transfer",
		},
		{
			name:         "unsupported-provider",
			resourceType: "nat_gateway",
			provider:     "azure",
			region:       "us-east-1",
			expectError:  true,
		},
		{
			name:         "unsupported-resource-type",
			resourceType: "unknown_resource",
			provider:     "aws",
			region:       "us-east-1",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing, err := service.GetPricing(ctx, tt.resourceType, tt.provider, tt.region)

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
		})
	}
}

func TestAWSPricingService_EstimateCost(t *testing.T) {
	service := NewAWSPricingService()
	ctx := context.Background()

	tests := []struct {
		name         string
		resource     *resource.Resource
		duration     time.Duration
		expectError  bool
		expectedCost float64
	}{
		{
			name: "estimate-nat-gateway-cost-720-hours",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "nat_gateway",
				},
				Provider: "aws",
				Region:   "us-east-1",
			},
			duration:     720 * time.Hour,
			expectError:  false,
			expectedCost: 32.40,
		},
		{
			name: "estimate-elastic-ip-cost-720-hours",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "elastic_ip",
				},
				Provider: "aws",
				Region:   "us-east-1",
			},
			duration:     720 * time.Hour,
			expectError:  false,
			expectedCost: 3.60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimate, err := service.EstimateCost(ctx, tt.resource, tt.duration)

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

func TestAWSPricingService_EstimateArchitectureCost(t *testing.T) {
	service := NewAWSPricingService()
	ctx := context.Background()

	resources := []*resource.Resource{
		{
			Type: resource.ResourceType{
				Name: "nat_gateway",
			},
			Provider: "aws",
			Region:   "us-east-1",
		},
		{
			Type: resource.ResourceType{
				Name: "elastic_ip",
			},
			Provider: "aws",
			Region:   "us-east-1",
		},
	}

	duration := 720 * time.Hour // 30 days

	estimate, err := service.EstimateArchitectureCost(ctx, resources, duration)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if estimate == nil {
		t.Fatal("Expected cost estimate but got nil")
	}

	// Expected: 32.40 (NAT) + 3.60 (EIP) = 36.00
	expectedTotal := 36.00
	if estimate.TotalCost != expectedTotal {
		t.Errorf("Expected total cost %.2f, got %.2f", expectedTotal, estimate.TotalCost)
	}
}

func TestAWSPricingService_ListSupportedResources(t *testing.T) {
	service := NewAWSPricingService()
	ctx := context.Background()

	tests := []struct {
		name         string
		provider     string
		expectError  bool
		minResources int
	}{
		{
			name:         "list-aws-resources",
			provider:     "aws",
			expectError:  false,
			minResources: 4, // nat_gateway, elastic_ip, network_interface, data_transfer
		},
		{
			name:        "unsupported-provider",
			provider:    "azure",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resources, err := service.ListSupportedResources(ctx, tt.provider)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(resources) < tt.minResources {
				t.Errorf("Expected at least %d resources, got %d", tt.minResources, len(resources))
			}

			// Check for expected resources
			expectedResources := map[string]bool{
				"nat_gateway":       false,
				"elastic_ip":        false,
				"network_interface": false,
				"data_transfer":     false,
			}

			for _, res := range resources {
				if _, exists := expectedResources[res]; exists {
					expectedResources[res] = true
				}
			}

			for res, found := range expectedResources {
				if !found {
					t.Errorf("Expected resource %s not found in list", res)
				}
			}
		})
	}
}
