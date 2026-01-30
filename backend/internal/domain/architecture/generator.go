package architecture

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// ArchitectureGenerator is an interface for cloud provider-specific architecture generation
// Each cloud provider can implement its own generator to handle provider-specific logic
type ArchitectureGenerator interface {
	// Provider returns the cloud provider this generator supports
	Provider() resource.CloudProvider

	// Generate converts a diagram graph into a domain architecture
	// This allows cloud providers to customize the mapping process
	Generate(diagramGraph *graph.DiagramGraph) (*Architecture, error)
}

// generatorRegistry stores registered architecture generators by provider
var generatorRegistry = make(map[resource.CloudProvider]ArchitectureGenerator)

// RegisterGenerator registers an architecture generator for a cloud provider
func RegisterGenerator(generator ArchitectureGenerator) {
	if generator == nil {
		panic("architecture generator cannot be nil")
	}
	generatorRegistry[generator.Provider()] = generator
}

// GetGenerator retrieves the architecture generator for a cloud provider
func GetGenerator(provider resource.CloudProvider) (ArchitectureGenerator, bool) {
	generator, ok := generatorRegistry[provider]
	return generator, ok
}

// MustGetGenerator retrieves the architecture generator for a cloud provider, panicking if not found
func MustGetGenerator(provider resource.CloudProvider) ArchitectureGenerator {
	generator, ok := generatorRegistry[provider]
	if !ok {
		panic("no architecture generator registered for provider: " + string(provider))
	}
	return generator
}
