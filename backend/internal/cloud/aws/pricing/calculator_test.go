package pricing

import (
	"context"
	"math"
	"testing"
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

func TestAWSPricingCalculator_CalculateResourceCost(t *testing.T) {
	service := NewAWSPricingService()
	calculator := NewAWSPricingCalculator(service)
	ctx := context.Background()

	tests := []struct {
		name         string
		resource     *resource.Resource
		duration     time.Duration
		expectError  bool
		expectedCost float64
	}{
		{
			name: "nat-gateway-1-hour",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "nat_gateway",
				},
				Provider: "aws",
				Region:   "us-east-1",
			},
			duration:     1 * time.Hour,
			expectError:  false,
			expectedCost: 0.045, // $0.045 per hour
		},
		{
			name: "nat-gateway-720-hours",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "nat_gateway",
				},
				Provider: "aws",
				Region:   "us-east-1",
			},
			duration:     720 * time.Hour, // 30 days
			expectError:  false,
			expectedCost: 32.40, // $0.045 * 720
		},
		{
			name: "elastic-ip-unattached-720-hours",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "elastic_ip",
				},
				Provider: "aws",
				Region:   "us-east-1",
			},
			duration:     720 * time.Hour,
			expectError:  false,
			expectedCost: 3.60, // $0.005 * 720
		},
		{
			name: "network-interface-unattached-720-hours",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "network_interface",
				},
				Provider: "aws",
				Region:   "us-east-1",
			},
			duration:     720 * time.Hour,
			expectError:  false,
			expectedCost: 7.20, // $0.01 * 720
		},
		{
			name: "unsupported-provider",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "nat_gateway",
				},
				Provider: "azure",
				Region:   "us-east-1",
			},
			duration:    1 * time.Hour,
			expectError: true,
		},
		{
			name: "ec2-instance-t3.micro-1-hour",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "ec2_instance",
				},
				Provider: "aws",
				Region:   "us-east-1",
				Metadata: map[string]interface{}{
					"instance_type": "t3.micro",
				},
			},
			duration:     1 * time.Hour,
			expectError:  false,
			expectedCost: 0.0104,
		},
		{
			name: "ec2-instance-t3.micro-720-hours",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "ec2_instance",
				},
				Provider: "aws",
				Region:   "us-east-1",
				Metadata: map[string]interface{}{
					"instance_type": "t3.micro",
				},
			},
			duration:     720 * time.Hour,
			expectError:  false,
			expectedCost: 7.488, // 0.0104 * 720
		},
		{
			name: "ec2-instance-m5.large-720-hours",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "ec2_instance",
				},
				Provider: "aws",
				Region:   "us-east-1",
				Metadata: map[string]interface{}{
					"instance_type": "m5.large",
				},
			},
			duration:     720 * time.Hour,
			expectError:  false,
			expectedCost: 69.12, // 0.096 * 720
		},
		{
			name: "ebs-volume-gp3-100gb-720-hours",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "ebs_volume",
				},
				Provider: "aws",
				Region:   "us-east-1",
				Metadata: map[string]interface{}{
					"size_gb":     100.0,
					"volume_type": "gp3",
				},
			},
			duration:     720 * time.Hour, // 1 month
			expectError:  false,
			expectedCost: 8.0, // 0.08 * 100 * 1
		},
		{
			name: "ebs-volume-gp2-50gb-720-hours",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "ebs_volume",
				},
				Provider: "aws",
				Region:   "us-east-1",
				Metadata: map[string]interface{}{
					"size_gb":     50.0,
					"volume_type": "gp2",
				},
			},
			duration:     720 * time.Hour,
			expectError:  false,
			expectedCost: 5.0, // 0.10 * 50 * 1
		},
		{
			name: "ebs-volume-missing-size-gb",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "ebs_volume",
				},
				Provider: "aws",
				Region:   "us-east-1",
				Metadata: map[string]interface{}{
					"volume_type": "gp3",
				},
			},
			duration:    720 * time.Hour,
			expectError: true,
		},
		{
			name: "s3-bucket-storage-only-100gb-720-hours",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "s3_bucket",
				},
				Provider: "aws",
				Region:   "us-east-1",
				Metadata: map[string]interface{}{
					"size_gb":       100.0,
					"storage_class": "standard",
				},
			},
			duration:     720 * time.Hour, // 1 month
			expectError:  false,
			expectedCost: 2.3, // 0.023 * 100 * 1
		},
		{
			name: "s3-bucket-full-cost-calculation",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "s3_bucket",
				},
				Provider: "aws",
				Region:   "us-east-1",
				Metadata: map[string]interface{}{
					"size_gb":          100.0,
					"storage_class":    "standard",
					"put_requests":     5000.0,
					"get_requests":     10000.0,
					"data_transfer_gb": 50.0,
				},
			},
			duration:     720 * time.Hour,
			expectError:  false,
			expectedCost: 2.3 + 0.025 + 0.004 + 4.41, // storage + PUT + GET + data transfer
		},
		{
			name: "s3-bucket-with-data-transfer-free-tier",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "s3_bucket",
				},
				Provider: "aws",
				Region:   "us-east-1",
				Metadata: map[string]interface{}{
					"size_gb":          100.0,
					"storage_class":    "standard",
					"data_transfer_gb": 0.5, // Less than 1GB, should be free
				},
			},
			duration:     720 * time.Hour,
			expectError:  false,
			expectedCost: 2.3, // Only storage cost
		},
		{
			name: "unsupported-resource-type",
			resource: &resource.Resource{
				Type: resource.ResourceType{
					Name: "unknown_resource",
				},
				Provider: "aws",
				Region:   "us-east-1",
			},
			duration:    1 * time.Hour,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimate, err := calculator.CalculateResourceCost(ctx, tt.resource, tt.duration)

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

			// Use epsilon for floating point comparison
			epsilon := 0.01
			if math.Abs(estimate.TotalCost-tt.expectedCost) > epsilon {
				t.Errorf("Expected total cost %.2f, got %.2f", tt.expectedCost, estimate.TotalCost)
			}

			if estimate.Currency != domainpricing.USD {
				t.Errorf("Expected currency USD, got %s", estimate.Currency)
			}

			if estimate.Provider != domainpricing.AWS {
				t.Errorf("Expected provider AWS, got %s", estimate.Provider)
			}

			if len(estimate.Breakdown) == 0 {
				t.Error("Expected breakdown to have at least one component")
			}
		})
	}
}

func TestAWSPricingCalculator_CalculateArchitectureCost(t *testing.T) {
	service := NewAWSPricingService()
	calculator := NewAWSPricingCalculator(service)
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
		{
			Type: resource.ResourceType{
				Name: "network_interface",
			},
			Provider: "aws",
			Region:   "us-east-1",
		},
	}

	duration := 720 * time.Hour // 30 days

	estimate, err := calculator.CalculateArchitectureCost(ctx, resources, duration)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if estimate == nil {
		t.Fatal("Expected cost estimate but got nil")
	}

	// Expected: 32.40 (NAT) + 3.60 (EIP) + 7.20 (ENI) = 43.20
	// Note: This test only includes NAT, EIP, and ENI
	expectedTotal := 43.20
	if estimate.TotalCost != expectedTotal {
		t.Errorf("Expected total cost %.2f, got %.2f", expectedTotal, estimate.TotalCost)
	}

	if len(estimate.Breakdown) < 3 {
		t.Errorf("Expected at least 3 breakdown components, got %d", len(estimate.Breakdown))
	}
}

func TestAWSPricingCalculator_GetResourcePricing(t *testing.T) {
	service := NewAWSPricingService()
	calculator := NewAWSPricingCalculator(service)
	ctx := context.Background()

	tests := []struct {
		name         string
		resourceType string
		provider     string
		region       string
		expectError  bool
		expectedType string
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
			name:         "get-ec2-instance-pricing",
			resourceType: "ec2_instance",
			provider:     "aws",
			region:       "us-east-1",
			expectError:  false,
			expectedType: "ec2_instance",
		},
		{
			name:         "get-ebs-volume-pricing",
			resourceType: "ebs_volume",
			provider:     "aws",
			region:       "us-east-1",
			expectError:  false,
			expectedType: "ebs_volume",
		},
		{
			name:         "get-s3-bucket-pricing",
			resourceType: "s3_bucket",
			provider:     "aws",
			region:       "us-east-1",
			expectError:  false,
			expectedType: "s3_bucket",
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
			pricing, err := calculator.GetResourcePricing(ctx, tt.resourceType, tt.provider, tt.region)

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

			if len(pricing.Components) == 0 {
				t.Error("Expected at least one pricing component")
			}
		})
	}
}
