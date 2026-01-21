package compute

import (
	"math"
	"testing"
	"time"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

func TestCalculateLambdaFunctionCost(t *testing.T) {
	tests := []struct {
		name             string
		duration         time.Duration
		memorySizeMB     float64
		averageDurationMs float64
		requestCount     float64
		dataTransferGB   float64
		region           string
		expectedCost     float64
		epsilon          float64 // For floating point comparison
	}{
		{
			name:             "compute-only-128mb-100ms-1m-requests-1-month",
			duration:         720 * time.Hour, // 1 month
			memorySizeMB:     128.0,
			averageDurationMs: 100.0,
			requestCount:     1000000.0, // 1M requests
			dataTransferGB:   0.0,
			region:           "us-east-1",
			expectedCost:     0.0000166667 * (128.0/1024.0) * (100.0/1000.0) * 1000000.0, // GB-seconds * rate
			epsilon:          0.0001,
		},
		{
			name:             "compute-only-256mb-200ms-500k-requests-1-month",
			duration:         720 * time.Hour,
			memorySizeMB:     256.0,
			averageDurationMs: 200.0,
			requestCount:     500000.0,
			dataTransferGB:   0.0,
			region:           "us-east-1",
			expectedCost:     0.0000166667 * (256.0/1024.0) * (200.0/1000.0) * 500000.0,
			epsilon:          0.0001,
		},
		{
			name:             "requests-with-free-tier",
			duration:         720 * time.Hour,
			memorySizeMB:     128.0,
			averageDurationMs: 100.0,
			requestCount:     2000000.0, // 2M requests, first 1M free
			dataTransferGB:   0.0,
			region:           "us-east-1",
			expectedCost:     0.0000166667*(128.0/1024.0)*(100.0/1000.0)*2000000.0 + (0.20/1000000.0)*1000000.0,
			epsilon:          0.0001,
		},
		{
			name:             "with-data-transfer",
			duration:         720 * time.Hour,
			memorySizeMB:     128.0,
			averageDurationMs: 100.0,
			requestCount:     1000000.0,
			dataTransferGB:   10.0, // 10 GB, first 1GB free
			region:           "us-east-1",
			expectedCost:     0.0000166667*(128.0/1024.0)*(100.0/1000.0)*1000000.0 + 0.09*9.0,
			epsilon:          0.0001,
		},
		{
			name:             "full-cost-calculation",
			duration:         720 * time.Hour,
			memorySizeMB:     512.0,
			averageDurationMs: 300.0,
			requestCount:     5000000.0, // 5M requests, first 1M free
			dataTransferGB:   20.0,       // 20 GB, first 1GB free
			region:           "us-east-1",
			expectedCost:     0.0000166667*(512.0/1024.0)*(300.0/1000.0)*5000000.0 + (0.20/1000000.0)*4000000.0 + 0.09*19.0,
			epsilon:          0.01,
		},
		{
			name:             "zero-requests",
			duration:         720 * time.Hour,
			memorySizeMB:     128.0,
			averageDurationMs: 100.0,
			requestCount:     0.0,
			dataTransferGB:   0.0,
			region:           "us-east-1",
			expectedCost:     0.0,
			epsilon:          0.0,
		},
		{
			name:             "regional-variation",
			duration:         720 * time.Hour,
			memorySizeMB:     128.0,
			averageDurationMs: 100.0,
			requestCount:     1000000.0,
			dataTransferGB:   0.0,
			region:           "ap-southeast-1", // 10% higher
			expectedCost:     0.0000166667 * 1.1 * (128.0/1024.0) * (100.0/1000.0) * 1000000.0,
			epsilon:          0.0001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualCost := CalculateLambdaFunctionCost(
				tt.duration,
				tt.memorySizeMB,
				tt.averageDurationMs,
				tt.requestCount,
				tt.dataTransferGB,
				tt.region,
			)

			if tt.epsilon > 0 {
				if math.Abs(actualCost-tt.expectedCost) > tt.epsilon {
					t.Errorf("expected cost %.6f, got %.6f (difference: %.6f)", tt.expectedCost, actualCost, math.Abs(actualCost-tt.expectedCost))
				}
			} else {
				if actualCost != tt.expectedCost {
					t.Errorf("expected cost %.6f, got %.6f", tt.expectedCost, actualCost)
				}
			}
		})
	}
}

func TestGetLambdaFunctionPricing(t *testing.T) {
	tests := []struct {
		name         string
		memorySizeMB float64
		region       string
		checkFunc    func(*testing.T, *domainpricing.ResourcePricing)
	}{
		{
			name:         "basic-pricing",
			memorySizeMB: 128.0,
			region:       "us-east-1",
			checkFunc: func(t *testing.T, pricing *domainpricing.ResourcePricing) {
				if pricing == nil {
					t.Fatal("expected pricing, got nil")
				}
				if pricing.ResourceType != "lambda_function" {
					t.Errorf("expected resource type 'lambda_function', got '%s'", pricing.ResourceType)
				}
				if len(pricing.Components) != 3 {
					t.Errorf("expected 3 components, got %d", len(pricing.Components))
				}
				// Check compute component
				if pricing.Components[0].Name != "Lambda Compute" {
					t.Errorf("expected component name 'Lambda Compute', got '%s'", pricing.Components[0].Name)
				}
				if pricing.Components[0].Rate != LambdaComputeRate {
					t.Errorf("expected compute rate %.10f, got %.10f", LambdaComputeRate, pricing.Components[0].Rate)
				}
			},
		},
		{
			name:         "request-component",
			memorySizeMB: 256.0,
			region:       "us-east-1",
			checkFunc: func(t *testing.T, pricing *domainpricing.ResourcePricing) {
				if pricing.Components[1].Name != "Lambda Requests" {
					t.Errorf("expected component name 'Lambda Requests', got '%s'", pricing.Components[1].Name)
				}
				if pricing.Components[1].Rate != LambdaRequestRate {
					t.Errorf("expected request rate %.2f, got %.2f", LambdaRequestRate, pricing.Components[1].Rate)
				}
			},
		},
		{
			name:         "data-transfer-component",
			memorySizeMB: 512.0,
			region:       "us-east-1",
			checkFunc: func(t *testing.T, pricing *domainpricing.ResourcePricing) {
				if pricing.Components[2].Name != "Lambda Data Transfer Out" {
					t.Errorf("expected component name 'Lambda Data Transfer Out', got '%s'", pricing.Components[2].Name)
				}
				if pricing.Components[2].Rate != LambdaDataTransferRate {
					t.Errorf("expected data transfer rate %.2f, got %.2f", LambdaDataTransferRate, pricing.Components[2].Rate)
				}
			},
		},
		{
			name:         "metadata-check",
			memorySizeMB: 1024.0,
			region:       "us-east-1",
			checkFunc: func(t *testing.T, pricing *domainpricing.ResourcePricing) {
				if pricing.Metadata == nil {
					t.Fatal("expected metadata, got nil")
				}
				if mem, ok := pricing.Metadata["memory_size_mb"].(float64); !ok || mem != 1024.0 {
					t.Errorf("expected memory_size_mb 1024.0, got %v", pricing.Metadata["memory_size_mb"])
				}
			},
		},
		{
			name:         "regional-variation",
			memorySizeMB: 128.0,
			region:       "ap-southeast-1",
			checkFunc: func(t *testing.T, pricing *domainpricing.ResourcePricing) {
				// Check that regional multiplier is applied
				expectedRate := LambdaComputeRate * 1.1
				if math.Abs(pricing.Components[0].Rate-expectedRate) > 0.0000001 {
					t.Errorf("expected compute rate %.10f, got %.10f", expectedRate, pricing.Components[0].Rate)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pricing := GetLambdaFunctionPricing(tt.memorySizeMB, tt.region)
			if tt.checkFunc != nil {
				tt.checkFunc(t, pricing)
			}
		})
	}
}
