package compute

import (
	"testing"
)

func TestCategorizeInstanceType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected InstanceCategory
	}{
		// Free Tier
		{"t2.micro free tier", "t2.micro", CategoryFreeTier},
		{"t3.micro free tier", "t3.micro", CategoryFreeTier},
		{"t4g.micro free tier", "t4g.micro", CategoryFreeTier},

		// General Purpose
		{"t3.small general purpose", "t3.small", CategoryGeneralPurpose},
		{"t3.medium general purpose", "t3.medium", CategoryGeneralPurpose},
		{"t3.large general purpose", "t3.large", CategoryGeneralPurpose},
		{"m5.large general purpose", "m5.large", CategoryGeneralPurpose},
		{"m5.xlarge general purpose", "m5.xlarge", CategoryGeneralPurpose},
		{"m6i.large general purpose", "m6i.large", CategoryGeneralPurpose},
		{"m6a.xlarge general purpose", "m6a.xlarge", CategoryGeneralPurpose},

		// Previous Generation (t2 except micro)
		{"t2.small previous generation", "t2.small", CategoryPreviousGeneration},
		{"t2.medium previous generation", "t2.medium", CategoryPreviousGeneration},

		// Compute Optimized
		{"c5.large compute optimized", "c5.large", CategoryComputeOptimized},
		{"c5.xlarge compute optimized", "c5.xlarge", CategoryComputeOptimized},
		{"c6i.large compute optimized", "c6i.large", CategoryComputeOptimized},
		{"c6a.xlarge compute optimized", "c6a.xlarge", CategoryComputeOptimized},

		// Memory Optimized
		{"r5.large memory optimized", "r5.large", CategoryMemoryOptimized},
		{"r5.xlarge memory optimized", "r5.xlarge", CategoryMemoryOptimized},
		{"r6i.large memory optimized", "r6i.large", CategoryMemoryOptimized},
		{"x1e.xlarge memory optimized", "x1e.xlarge", CategoryMemoryOptimized},
		{"u-6tb1.metal memory optimized", "u-6tb1.metal", CategoryMemoryOptimized},

		// Storage Optimized
		{"i3.large storage optimized", "i3.large", CategoryStorageOptimized},
		{"i3.xlarge storage optimized", "i3.xlarge", CategoryStorageOptimized},
		{"d2.xlarge storage optimized", "d2.xlarge", CategoryStorageOptimized},
		{"h1.2xlarge storage optimized", "h1.2xlarge", CategoryStorageOptimized},

		// Accelerated Computing
		{"p3.2xlarge accelerated computing", "p3.2xlarge", CategoryAcceleratedComputing},
		{"g4dn.xlarge accelerated computing", "g4dn.xlarge", CategoryAcceleratedComputing},
		{"f1.2xlarge accelerated computing", "f1.2xlarge", CategoryAcceleratedComputing},
		{"inf1.xlarge accelerated computing", "inf1.xlarge", CategoryAcceleratedComputing},

		// High Performance Computing
		{"z1d.large high performance", "z1d.large", CategoryHighPerformanceComputing},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CategorizeInstanceType(tt.input)
			if result != tt.expected {
				t.Errorf("CategorizeInstanceType(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsFreeTierEligible(t *testing.T) {
	tests := []struct {
		name         string
		instanceType string
		region       string
		expected     bool
	}{
		{"t2.micro eligible", "t2.micro", "us-east-1", true},
		{"t3.micro eligible", "t3.micro", "us-east-1", true},
		{"t4g.micro eligible", "t4g.micro", "us-east-1", true},
		{"t3.small not eligible", "t3.small", "us-east-1", false},
		{"m5.large not eligible", "m5.large", "us-east-1", false},
		{"case insensitive", "T3.MICRO", "us-east-1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsFreeTierEligible(tt.instanceType, tt.region)
			if result != tt.expected {
				t.Errorf("IsFreeTierEligible(%q, %q) = %v, want %v", tt.instanceType, tt.region, result, tt.expected)
			}
		})
	}
}
