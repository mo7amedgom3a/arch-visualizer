package storage

import (
	"testing"
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

func TestCalculateS3BucketCost(t *testing.T) {
	tests := []struct {
		name           string
		duration       time.Duration
		sizeGB         float64
		putRequests    float64
		getRequests    float64
		dataTransferGB float64
		storageClass   string
		region         string
		expectedCost   float64
	}{
		{
			name:           "storage-only-100gb-1-month",
			duration:       720 * time.Hour, // 30 days = 1 month
			sizeGB:         100.0,
			putRequests:    0.0,
			getRequests:    0.0,
			dataTransferGB: 0.0,
			storageClass:   "standard",
			region:         "us-east-1",
			expectedCost:   2.3, // 0.023 * 100 * 1
		},
		{
			name:           "storage-only-100gb-2-months",
			duration:       1440 * time.Hour, // 60 days = 2 months
			sizeGB:         100.0,
			putRequests:    0.0,
			getRequests:    0.0,
			dataTransferGB: 0.0,
			storageClass:   "standard",
			region:         "us-east-1",
			expectedCost:   4.6, // 0.023 * 100 * 2
		},
		{
			name:           "storage-only-50gb-15-days",
			duration:       360 * time.Hour, // 15 days = 0.5 months
			sizeGB:         50.0,
			putRequests:    0.0,
			getRequests:    0.0,
			dataTransferGB: 0.0,
			storageClass:   "standard",
			region:         "us-east-1",
			expectedCost:   0.575, // 0.023 * 50 * 0.5
		},
		{
			name:           "storage-with-requests",
			duration:       720 * time.Hour,
			sizeGB:         100.0,
			putRequests:    5000.0,  // 5,000 PUT requests
			getRequests:    10000.0, // 10,000 GET requests
			dataTransferGB: 0.0,
			storageClass:   "standard",
			region:         "us-east-1",
			expectedCost:   2.3 + 0.025 + 0.004, // storage + PUT + GET = 2.3 + (0.005/1000)*5000 + (0.0004/1000)*10000
		},
		{
			name:           "storage-with-data-transfer",
			duration:       720 * time.Hour,
			sizeGB:         100.0,
			putRequests:    0.0,
			getRequests:    0.0,
			dataTransferGB: 50.0, // 50 GB transfer, first 1GB free
			storageClass:   "standard",
			region:         "us-east-1",
			expectedCost:   2.3 + 4.41, // storage + data transfer = 2.3 + (50-1)*0.09
		},
		{
			name:           "full-cost-calculation",
			duration:       720 * time.Hour,
			sizeGB:         100.0,
			putRequests:    5000.0,
			getRequests:    10000.0,
			dataTransferGB: 50.0,
			storageClass:   "standard",
			region:         "us-east-1",
			expectedCost:   2.3 + 0.025 + 0.004 + 4.41, // all components
		},
		{
			name:           "standard-ia-storage",
			duration:       720 * time.Hour,
			sizeGB:         100.0,
			putRequests:    0.0,
			getRequests:    0.0,
			dataTransferGB: 0.0,
			storageClass:   "standard-ia",
			region:         "us-east-1",
			expectedCost:   1.25, // 0.0125 * 100 * 1
		},
		{
			name:           "glacier-storage",
			duration:       720 * time.Hour,
			sizeGB:         100.0,
			putRequests:    0.0,
			getRequests:    0.0,
			dataTransferGB: 0.0,
			storageClass:   "glacier",
			region:         "us-east-1",
			expectedCost:   0.4, // 0.004 * 100 * 1
		},
		{
			name:           "unknown-storage-class-defaults-to-standard",
			duration:       720 * time.Hour,
			sizeGB:         100.0,
			putRequests:    0.0,
			getRequests:    0.0,
			dataTransferGB: 0.0,
			storageClass:   "unknown",
			region:         "us-east-1",
			expectedCost:   2.3, // Defaults to standard rate
		},
		{
			name:           "data-transfer-within-free-tier",
			duration:       720 * time.Hour,
			sizeGB:         100.0,
			putRequests:    0.0,
			getRequests:    0.0,
			dataTransferGB: 0.5, // Less than 1GB, should be free
			storageClass:   "standard",
			region:         "us-east-1",
			expectedCost:   2.3, // Only storage cost
		},
		{
			name:           "data-transfer-exactly-1gb-free",
			duration:       720 * time.Hour,
			sizeGB:         100.0,
			putRequests:    0.0,
			getRequests:    0.0,
			dataTransferGB: 1.0, // Exactly 1GB, should be free
			storageClass:   "standard",
			region:         "us-east-1",
			expectedCost:   2.3, // Only storage cost
		},
		{
			name:           "regional-variation",
			duration:       720 * time.Hour,
			sizeGB:         100.0,
			putRequests:    0.0,
			getRequests:    0.0,
			dataTransferGB: 0.0,
			storageClass:   "standard",
			region:         "ap-southeast-1",
			expectedCost:   2.53, // 0.023 * 1.1 * 100 * 1
		},
		{
			name:           "zero-values",
			duration:       720 * time.Hour,
			sizeGB:         0.0,
			putRequests:    0.0,
			getRequests:    0.0,
			dataTransferGB: 0.0,
			storageClass:   "standard",
			region:         "us-east-1",
			expectedCost:   0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := CalculateS3BucketCost(
				tt.duration,
				tt.sizeGB,
				tt.putRequests,
				tt.getRequests,
				tt.dataTransferGB,
				tt.storageClass,
				tt.region,
			)
			// Allow small floating point differences
			if cost < tt.expectedCost-0.01 || cost > tt.expectedCost+0.01 {
				t.Errorf("Expected cost %.4f, got %.4f", tt.expectedCost, cost)
			}
		})
	}
}

func TestGetS3BucketPricing(t *testing.T) {
	tests := []struct {
		name         string
		storageClass string
		region       string
		checkFunc    func(*testing.T, *domainpricing.ResourcePricing)
	}{
		{
			name:         "standard-storage-class",
			storageClass: "standard",
			region:       "us-east-1",
			checkFunc: func(t *testing.T, pricing *domainpricing.ResourcePricing) {
				if pricing == nil {
					t.Fatal("Expected pricing, got nil")
				}
				if pricing.ResourceType != "s3_bucket" {
					t.Errorf("Expected resource type 's3_bucket', got '%s'", pricing.ResourceType)
				}
				if pricing.Provider != domainpricing.AWS {
					t.Errorf("Expected provider AWS, got '%s'", pricing.Provider)
				}
				if len(pricing.Components) != 4 {
					t.Errorf("Expected 4 components, got %d", len(pricing.Components))
				}
				// Check storage component
				if pricing.Components[0].Name != "S3 Storage" {
					t.Errorf("Expected first component 'S3 Storage', got '%s'", pricing.Components[0].Name)
				}
				if pricing.Components[0].Rate != 0.023 {
					t.Errorf("Expected storage rate 0.023, got %.4f", pricing.Components[0].Rate)
				}
				// Check PUT requests component
				if pricing.Components[1].Name != "S3 PUT Requests" {
					t.Errorf("Expected second component 'S3 PUT Requests', got '%s'", pricing.Components[1].Name)
				}
				if pricing.Components[1].Rate != 0.005 {
					t.Errorf("Expected PUT rate 0.005, got %.4f", pricing.Components[1].Rate)
				}
				// Check GET requests component
				if pricing.Components[2].Name != "S3 GET Requests" {
					t.Errorf("Expected third component 'S3 GET Requests', got '%s'", pricing.Components[2].Name)
				}
				if pricing.Components[2].Rate != 0.0004 {
					t.Errorf("Expected GET rate 0.0004, got %.4f", pricing.Components[2].Rate)
				}
				// Check data transfer component
				if pricing.Components[3].Name != "S3 Data Transfer Out" {
					t.Errorf("Expected fourth component 'S3 Data Transfer Out', got '%s'", pricing.Components[3].Name)
				}
				if pricing.Components[3].Rate != 0.09 {
					t.Errorf("Expected data transfer rate 0.09, got %.4f", pricing.Components[3].Rate)
				}
			},
		},
		{
			name:         "standard-ia-storage-class",
			storageClass: "standard-ia",
			region:       "us-east-1",
			checkFunc: func(t *testing.T, pricing *domainpricing.ResourcePricing) {
				if pricing == nil {
					t.Fatal("Expected pricing, got nil")
				}
				if pricing.Components[0].Rate != 0.0125 {
					t.Errorf("Expected storage rate 0.0125, got %.4f", pricing.Components[0].Rate)
				}
			},
		},
		{
			name:         "unknown-storage-class-defaults-to-standard",
			storageClass: "unknown",
			region:       "us-east-1",
			checkFunc: func(t *testing.T, pricing *domainpricing.ResourcePricing) {
				if pricing == nil {
					t.Fatal("Expected pricing, got nil")
				}
				// Should default to standard rate
				if pricing.Components[0].Rate != 0.023 {
					t.Errorf("Expected default storage rate 0.023, got %.4f", pricing.Components[0].Rate)
				}
			},
		},
		{
			name:         "regional-variation",
			storageClass: "standard",
			region:       "ap-southeast-1",
			checkFunc: func(t *testing.T, pricing *domainpricing.ResourcePricing) {
				if pricing == nil {
					t.Fatal("Expected pricing, got nil")
				}
				// Should apply regional multiplier
				expectedRate := 0.023 * 1.1
				// Use epsilon for floating point comparison
				epsilon := 0.0001
				diff := pricing.Components[0].Rate - expectedRate
				if diff < -epsilon || diff > epsilon {
					t.Errorf("Expected storage rate %.4f, got %.4f", expectedRate, pricing.Components[0].Rate)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetS3BucketPricing(tt.storageClass, tt.region)
			if tt.checkFunc != nil {
				tt.checkFunc(t, pricing)
			}
		})
	}
}
