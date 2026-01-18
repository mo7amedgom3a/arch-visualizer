package compute

import (
	"strings"
)

// CategorizeInstanceType maps an instance type name to its category based on prefix and special cases
// This is the public function that should be used to categorize instance types
func CategorizeInstanceType(name string) InstanceCategory {
	return categorizeInstanceType(name)
}

// categorizeInstanceType maps an instance type name to its category based on prefix and special cases
func categorizeInstanceType(name string) InstanceCategory {
	// Convert to lowercase for comparison
	name = strings.ToLower(name)

	// Special cases: Free tier eligible types
	freeTierTypes := map[string]bool{
		"t2.micro":  true,
		"t3.micro":  true,
		"t4g.micro": true,
	}
	if freeTierTypes[name] {
		return CategoryFreeTier
	}

	// Extract prefix (family) from instance type (e.g., "t3" from "t3.micro")
	parts := strings.Split(name, ".")
	if len(parts) == 0 {
		return CategoryPreviousGeneration
	}
	family := parts[0]

	// Remove generation suffix if present (e.g., "m5" from "m5a", "m5n")
	// Keep only the base family letter(s)
	baseFamily := ""
	for i, r := range family {
		if i == 0 {
			baseFamily += string(r)
		} else if r >= '0' && r <= '9' {
			// Stop at first digit (generation number)
			break
		} else if (r >= 'a' && r <= 'z') && i == 1 {
			// Second character might be a modifier (e.g., 'a', 'n', 'd')
			// We'll check the full prefix for known patterns
			baseFamily += string(r)
		}
	}

	// Category mapping based on prefix
	switch {
	// General Purpose: t* (burstable), m*
	case strings.HasPrefix(family, "t"):
		// t2 is previous generation, but t2.micro is free tier (handled above)
		if strings.HasPrefix(family, "t2") {
			return CategoryPreviousGeneration
		}
		return CategoryGeneralPurpose
	case strings.HasPrefix(family, "m"):
		return CategoryGeneralPurpose

	// Compute Optimized: c*
	case strings.HasPrefix(family, "c"):
		return CategoryComputeOptimized

	// Memory Optimized: r*, x*, u*
	case strings.HasPrefix(family, "r"):
		return CategoryMemoryOptimized
	case strings.HasPrefix(family, "x"):
		return CategoryMemoryOptimized
	case strings.HasPrefix(family, "u"):
		return CategoryMemoryOptimized

	// Accelerated Computing: p*, g*, f*, inf*
	// Check "inf" before "i" to avoid matching with storage optimized
	case strings.HasPrefix(family, "inf"):
		return CategoryAcceleratedComputing
	case strings.HasPrefix(family, "p"):
		return CategoryAcceleratedComputing
	case strings.HasPrefix(family, "g"):
		return CategoryAcceleratedComputing
	case strings.HasPrefix(family, "f"):
		return CategoryAcceleratedComputing

	// Storage Optimized: i*, d*, h*
	case strings.HasPrefix(family, "i"):
		return CategoryStorageOptimized
	case strings.HasPrefix(family, "d"):
		return CategoryStorageOptimized
	case strings.HasPrefix(family, "h"):
		return CategoryStorageOptimized

	// High Performance Computing: z*
	case strings.HasPrefix(family, "z"):
		return CategoryHighPerformanceComputing

	// Previous Generation: anything else or unknown
	default:
		return CategoryPreviousGeneration
	}
}

// IsFreeTierEligible checks if an instance type is eligible for AWS Free Tier
func IsFreeTierEligible(instanceType string, region string) bool {
	// Free tier eligible types (region-specific, but we'll use a general list)
	freeTierTypes := map[string]bool{
		"t2.micro":  true,
		"t3.micro":  true,
		"t4g.micro": true,
	}
	return freeTierTypes[strings.ToLower(instanceType)]
}
