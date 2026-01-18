package compute

// InstanceCategory represents the category/class of an EC2 instance type
type InstanceCategory string

const (
	// CategoryGeneralPurpose represents general-purpose instances (balanced compute, memory, networking)
	CategoryGeneralPurpose InstanceCategory = "general_purpose"

	// CategoryComputeOptimized represents compute-optimized instances (high CPU performance)
	CategoryComputeOptimized InstanceCategory = "compute_optimized"

	// CategoryMemoryOptimized represents memory-optimized instances (high memory-to-CPU ratio)
	CategoryMemoryOptimized InstanceCategory = "memory_optimized"

	// CategoryStorageOptimized represents storage-optimized instances (high storage I/O)
	CategoryStorageOptimized InstanceCategory = "storage_optimized"

	// CategoryAcceleratedComputing represents instances with GPUs or other accelerators
	CategoryAcceleratedComputing InstanceCategory = "accelerated_computing"

	// CategoryHighPerformanceComputing represents HPC-optimized instances
	CategoryHighPerformanceComputing InstanceCategory = "high_performance_computing"

	// CategoryPreviousGeneration represents older generation instance types
	CategoryPreviousGeneration InstanceCategory = "previous_generation"

	// CategoryFreeTier represents free tier eligible instance types
	CategoryFreeTier InstanceCategory = "free_tier"
)

// String returns the string representation of the category
func (c InstanceCategory) String() string {
	return string(c)
}

// IsValid checks if the category is a valid instance category
func (c InstanceCategory) IsValid() bool {
	switch c {
	case CategoryGeneralPurpose,
		CategoryComputeOptimized,
		CategoryMemoryOptimized,
		CategoryStorageOptimized,
		CategoryAcceleratedComputing,
		CategoryHighPerformanceComputing,
		CategoryPreviousGeneration,
		CategoryFreeTier:
		return true
	default:
		return false
	}
}

// GetDescription returns a human-readable description of the category
func (c InstanceCategory) GetDescription() string {
	switch c {
	case CategoryGeneralPurpose:
		return "General Purpose - Balanced compute, memory, and networking resources"
	case CategoryComputeOptimized:
		return "Compute Optimized - High CPU performance for compute-intensive workloads"
	case CategoryMemoryOptimized:
		return "Memory Optimized - High memory-to-CPU ratio for memory-intensive workloads"
	case CategoryStorageOptimized:
		return "Storage Optimized - High storage I/O performance for storage-intensive workloads"
	case CategoryAcceleratedComputing:
		return "Accelerated Computing - Instances with GPUs or other accelerators"
	case CategoryHighPerformanceComputing:
		return "High Performance Computing - Optimized for HPC workloads"
	case CategoryPreviousGeneration:
		return "Previous Generation - Older generation instance types"
	case CategoryFreeTier:
		return "Free Tier - Eligible for AWS Free Tier"
	default:
		return "Unknown category"
	}
}

// AllCategories returns all valid instance categories
func AllCategories() []InstanceCategory {
	return []InstanceCategory{
		CategoryGeneralPurpose,
		CategoryComputeOptimized,
		CategoryMemoryOptimized,
		CategoryStorageOptimized,
		CategoryAcceleratedComputing,
		CategoryHighPerformanceComputing,
		CategoryPreviousGeneration,
		CategoryFreeTier,
	}
}
