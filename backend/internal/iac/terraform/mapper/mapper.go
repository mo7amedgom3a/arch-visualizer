package mapper

import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"

// ResourceMapper maps a domain resource into one or more Terraform blocks.
// Implementations must be provider-specific (AWS/GCP/Azure...), but the interface
// is cloud-agnostic and lives in the Terraform engine package.
type ResourceMapper interface {
	Provider() string

	// SupportsResource reports whether this mapper can map the given domain resource type.
	SupportsResource(resourceType string) bool

	// MapResource converts a domain resource to Terraform blocks.
	// The returned blocks are expected to be in a stable order.
	MapResource(res *resource.Resource) ([]TerraformBlock, error)
}

