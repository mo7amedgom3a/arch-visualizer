package storage

// StorageResource is a marker interface for all storage resources
// This allows polymorphic handling of different storage resource types
type StorageResource interface {
	// GetID returns the unique identifier of the resource
	GetID() string

	// GetName returns the name of the resource
	GetName() string

	// GetRegion returns the region where the resource is located
	GetRegion() string

	// Validate performs domain-level validation
	Validate() error
}
