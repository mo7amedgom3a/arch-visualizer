package architecture

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// IRTypeMapper provides mapping from IR resource types to domain resource names
// This interface allows cloud providers to provide their own IR type mappings
// via their inventory system
type IRTypeMapper interface {
	// GetResourceNameByIRType maps an IR type (kebab-case) to a domain resource name (PascalCase)
	// Returns empty string and false if the IR type is not found
	GetResourceNameByIRType(irType string) (string, bool)
}

// providerMappers is a registry of provider-specific IR type mappers
var providerMappers = make(map[resource.CloudProvider]IRTypeMapper)

// RegisterIRTypeMapper registers an IR type mapper for a cloud provider
func RegisterIRTypeMapper(provider resource.CloudProvider, mapper IRTypeMapper) {
	providerMappers[provider] = mapper
}

// GetIRTypeMapper retrieves the IR type mapper for a cloud provider
func GetIRTypeMapper(provider resource.CloudProvider) (IRTypeMapper, bool) {
	mapper, ok := providerMappers[provider]
	return mapper, ok
}
