package compute

// InstanceTypeInfo contains detailed information about an EC2 instance type
type InstanceTypeInfo struct {
	// Name is the instance type name (e.g., "t3.micro", "m5.large")
	Name string `json:"name"`

	// Category is the instance category classification
	Category InstanceCategory `json:"category"`

	// VCPU is the number of virtual CPUs
	VCPU int `json:"vcpu"`

	// MemoryGiB is the amount of memory in GiB
	MemoryGiB float64 `json:"memory_gib"`

	// StorageType describes the storage type (e.g., "EBS", "NVMe SSD", "HDD")
	StorageType string `json:"storage_type,omitempty"`

	// HasLocalStorage indicates if the instance has local instance storage
	HasLocalStorage bool `json:"has_local_storage"`

	// LocalStorageSizeGiB is the size of local storage in GiB (if applicable)
	LocalStorageSizeGiB *float64 `json:"local_storage_size_gib,omitempty"`

	// MaxNetworkGbps is the maximum network performance in Gbps
	MaxNetworkGbps float64 `json:"max_network_gbps"`

	// EBSBandwidthGbps is the EBS bandwidth in Gbps (if applicable)
	EBSBandwidthGbps *float64 `json:"ebs_bandwidth_gbps,omitempty"`

	// FreeTierEligible indicates if this instance type is eligible for AWS Free Tier
	FreeTierEligible bool `json:"free_tier_eligible"`

	// SupportedArchitectures lists supported CPU architectures (e.g., "x86_64", "arm64")
	SupportedArchitectures []string `json:"supported_architectures,omitempty"`

	// SupportedVirtualizationTypes lists supported virtualization types (e.g., "hvm", "paravirtual")
	SupportedVirtualizationTypes []string `json:"supported_virtualization_types,omitempty"`

	// Region is the AWS region where this instance type is available (optional, empty means all regions)
	Region string `json:"region,omitempty"`
}

// Validate performs basic validation on the instance type info
func (i *InstanceTypeInfo) Validate() error {
	if i.Name == "" {
		return ErrInstanceTypeNameRequired
	}
	if !i.Category.IsValid() {
		return ErrInvalidCategory
	}
	if i.VCPU <= 0 {
		return ErrInvalidVCPU
	}
	if i.MemoryGiB <= 0 {
		return ErrInvalidMemory
	}
	return nil
}

// IsAvailableInRegion checks if the instance type is available in the specified region
func (i *InstanceTypeInfo) IsAvailableInRegion(region string) bool {
	if i.Region == "" {
		return true // Available in all regions
	}
	return i.Region == region
}

// GetMemoryMB returns memory in MB
func (i *InstanceTypeInfo) GetMemoryMB() float64 {
	return i.MemoryGiB * 1024
}

// GetLocalStorageSizeMB returns local storage size in MB (if applicable)
func (i *InstanceTypeInfo) GetLocalStorageSizeMB() *float64 {
	if i.LocalStorageSizeGiB == nil {
		return nil
	}
	mb := *i.LocalStorageSizeGiB * 1024
	return &mb
}
