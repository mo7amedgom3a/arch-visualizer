package interfaces

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// ArchitectureService handles architecture mapping and validation
type ArchitectureService interface {
	// MapFromDiagram converts a diagram graph to a domain architecture
	MapFromDiagram(ctx context.Context, graph *graph.DiagramGraph, provider resource.CloudProvider) (*architecture.Architecture, error)

	// ValidateRules validates an architecture against domain rules and constraints
	ValidateRules(ctx context.Context, arch *architecture.Architecture, provider resource.CloudProvider) (*RuleValidationResult, error)

	// GetSortedResources returns resources sorted by dependencies (topological sort)
	GetSortedResources(ctx context.Context, arch *architecture.Architecture) ([]*resource.Resource, error)
}

// RuleValidationResult contains the result of rule validation
type RuleValidationResult struct {
	Valid   bool
	Results map[string]*ResourceValidationResult
	Errors  []string
}

// ResourceValidationResult contains validation result for a single resource
type ResourceValidationResult struct {
	ResourceID   string
	ResourceType string
	Valid        bool
	Errors       []ValidationError
}

// ValidationError represents a validation error
type ValidationError struct {
	ResourceID   string
	ResourceType string
	Message      string
	Code         string
}
