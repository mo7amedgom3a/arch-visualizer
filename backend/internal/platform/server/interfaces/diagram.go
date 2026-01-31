package interfaces

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator"
)

// DiagramService handles diagram parsing and validation
type DiagramService interface {
	// Parse parses diagram JSON into a DiagramGraph
	Parse(ctx context.Context, jsonData []byte) (*graph.DiagramGraph, error)

	// Validate validates a diagram graph against structural and schema rules
	Validate(ctx context.Context, graph *graph.DiagramGraph, opts *validator.ValidationOptions) (*validator.ValidationResult, error)
}
