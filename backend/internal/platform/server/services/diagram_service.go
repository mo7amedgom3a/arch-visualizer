package services

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/parser"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// DiagramServiceImpl implements DiagramService interface
type DiagramServiceImpl struct {
	// No dependencies needed - parser and validator are stateless
}

// NewDiagramService creates a new diagram service
func NewDiagramService() serverinterfaces.DiagramService {
	return &DiagramServiceImpl{}
}

// Parse parses diagram JSON into a DiagramGraph
func (s *DiagramServiceImpl) Parse(ctx context.Context, jsonData []byte) (*graph.DiagramGraph, error) {
	// Parse IR diagram
	irDiagram, err := parser.ParseIRDiagram(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IR diagram: %w", err)
	}

	// Normalize to graph
	diagramGraph, err := parser.NormalizeToGraph(irDiagram)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize diagram: %w", err)
	}

	return diagramGraph, nil
}

// Validate validates a diagram graph against structural and schema rules
func (s *DiagramServiceImpl) Validate(ctx context.Context, graph *graph.DiagramGraph, opts *validator.ValidationOptions) (*validator.ValidationResult, error) {
	if graph == nil {
		return nil, fmt.Errorf("diagram graph is nil")
	}

	result := validator.Validate(graph, opts)
	return result, nil
}
