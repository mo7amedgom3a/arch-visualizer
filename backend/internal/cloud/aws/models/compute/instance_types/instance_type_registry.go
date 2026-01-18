package compute

// InstanceTypeRegistry defines the interface for managing instance types
type InstanceTypeRegistry interface {
	// GetInstanceType retrieves information about a specific instance type
	GetInstanceType(name string) (*InstanceTypeInfo, error)

	// ListByCategory returns all instance types in a specific category
	ListByCategory(category InstanceCategory) ([]*InstanceTypeInfo, error)

	// ListFreeTier returns all free tier eligible instance types
	ListFreeTier() ([]*InstanceTypeInfo, error)

	// Search allows filtering instance types by various criteria
	Search(filters *InstanceTypeFilters) ([]*InstanceTypeInfo, error)

	// ListAll returns all registered instance types
	ListAll() ([]*InstanceTypeInfo, error)
}

// InstanceTypeFilters defines filters for searching instance types
type InstanceTypeFilters struct {
	Category              *InstanceCategory
	MinVCPU               *int
	MaxVCPU               *int
	MinMemoryGiB          *float64
	MaxMemoryGiB          *float64
	FreeTierOnly          bool
	HasLocalStorage       *bool
	SupportedArchitecture *string
	Region                string
}
