package architecture

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/resource"
)

func init() {
	// Register AWS architecture generator
	generator := NewAWSArchitectureGenerator()
	architecture.RegisterGenerator(generator)

	// Register AWS resource type mapper
	mapper := NewAWSResourceTypeMapper()
	architecture.RegisterResourceTypeMapper(resource.AWS, mapper)
}
