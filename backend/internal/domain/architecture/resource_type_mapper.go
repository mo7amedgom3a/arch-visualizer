package architecture

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// ResourceTypeMapper provides mapping from IR resource types and resource names to domain ResourceType
// Each cloud provider implements this interface with provider-specific mappings
// since resource names differ between cloud providers (e.g., AWS VPC vs Azure VirtualNetwork)
type ResourceTypeMapper interface {
	// MapIRTypeToResourceType maps an IR type (kebab-case) directly to a ResourceType
	// This is the primary method for mapping IR types to domain resource types
	MapIRTypeToResourceType(irType string) (*resource.ResourceType, error)

	// MapResourceNameToResourceType maps a resource name (PascalCase) to ResourceType
	// This is used when we already have the resource name from inventory
	MapResourceNameToResourceType(resourceName string) (*resource.ResourceType, error)
}

// resourceTypeMapperRegistry stores registered resource type mappers by provider
var resourceTypeMapperRegistry = make(map[resource.CloudProvider]ResourceTypeMapper)

// RegisterResourceTypeMapper registers a resource type mapper for a cloud provider
func RegisterResourceTypeMapper(provider resource.CloudProvider, mapper ResourceTypeMapper) {
	if mapper == nil {
		panic("resource type mapper cannot be nil")
	}
	resourceTypeMapperRegistry[provider] = mapper
}

// GetResourceTypeMapper retrieves the resource type mapper for a cloud provider
func GetResourceTypeMapper(provider resource.CloudProvider) (ResourceTypeMapper, bool) {
	mapper, ok := resourceTypeMapperRegistry[provider]
	return mapper, ok
}

// MustGetResourceTypeMapper retrieves the resource type mapper for a cloud provider, panicking if not found
func MustGetResourceTypeMapper(provider resource.CloudProvider) ResourceTypeMapper {
	mapper, ok := resourceTypeMapperRegistry[provider]
	if !ok {
		panic("no resource type mapper registered for provider: " + string(provider))
	}
	return mapper
}
