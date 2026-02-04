package services

import (
	"context"
	"log/slog"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator"
)

func TestDiagramService_Parse(t *testing.T) {
	service := NewDiagramService(slog.Default())
	ctx := context.Background()

	tests := []struct {
		name      string
		jsonData  []byte
		wantError bool
	}{
		{
			name: "valid diagram JSON",
			jsonData: []byte(`{
				"nodes": [
					{
						"id": "node-1",
						"type": "resourceNode",
						"data": {
							"label": "VPC",
							"resourceType": "vpc",
							"config": {
								"cidr": "10.0.0.0/16"
							}
						}
					}
				],
				"edges": [],
				"variables": [],
				"outputs": []
			}`),
			wantError: false,
		},
		{
			name:      "invalid JSON",
			jsonData:  []byte(`invalid json`),
			wantError: true,
		},
		{
			name:      "empty JSON",
			jsonData:  []byte(`{}`),
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Parse(ctx, tt.jsonData)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result == nil {
				t.Errorf("expected result but got nil")
			}
		})
	}
}

func TestDiagramService_Validate(t *testing.T) {
	service := NewDiagramService(slog.Default())
	ctx := context.Background()

	// Create a valid diagram graph
	validGraph := &graph.DiagramGraph{
		Nodes: map[string]*graph.Node{
			"node-1": {
				ID:           "node-1",
				Type:         "resourceNode",
				ResourceType: "vpc",
				Label:        "VPC",
				Config: map[string]interface{}{
					"cidr": "10.0.0.0/16",
				},
			},
		},
		Edges:     []*graph.Edge{},
		Variables: []graph.Variable{},
		Outputs:   []graph.Output{},
	}

	tests := []struct {
		name      string
		graph     *graph.DiagramGraph
		opts      *validator.ValidationOptions
		wantError bool
	}{
		{
			name:      "valid graph",
			graph:     validGraph,
			opts:      nil,
			wantError: false,
		},
		{
			name:      "nil graph",
			graph:     nil,
			opts:      nil,
			wantError: true,
		},
		{
			name:  "graph with validation options",
			graph: validGraph,
			opts: &validator.ValidationOptions{
				Provider: "aws",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Validate(ctx, tt.graph, tt.opts)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result == nil {
				t.Errorf("expected result but got nil")
			}
		})
	}
}

func TestDiagramService_Integration(t *testing.T) {
	service := NewDiagramService(slog.Default())
	ctx := context.Background()

	// Test full flow: Parse -> Validate
	jsonData := []byte(`{
		"nodes": [
			{
				"id": "vpc-1",
				"type": "resourceNode",
				"data": {
					"label": "Main VPC",
					"resourceType": "vpc",
					"config": {
						"cidr": "10.0.0.0/16"
					}
				}
			},
			{
				"id": "subnet-1",
				"type": "resourceNode",
				"parentId": "vpc-1",
				"data": {
					"label": "Public Subnet",
					"resourceType": "subnet",
					"config": {
						"cidr": "10.0.1.0/24"
					}
				}
			}
		],
		"edges": [],
		"variables": [],
		"outputs": []
	}`)

	// Parse
	graph, err := service.Parse(ctx, jsonData)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if graph == nil {
		t.Fatal("expected graph but got nil")
	}

	if len(graph.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(graph.Nodes))
	}

	// Validate
	validationResult, err := service.Validate(ctx, graph, nil)
	if err != nil {
		t.Fatalf("failed to validate: %v", err)
	}

	if validationResult == nil {
		t.Fatal("expected validation result but got nil")
	}
}
