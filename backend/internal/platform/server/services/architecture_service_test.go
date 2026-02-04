package services

import (
	"context"
	"log/slog"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// mockRuleService is a mock implementation of RuleService
type mockRuleService struct {
	validateArchitectureFunc func(ctx context.Context, architecture interface{}) (map[string]interface{}, error)
	loadRulesFunc            func(ctx context.Context, constraints []serverinterfaces.ConstraintRecord) error
}

func (m *mockRuleService) LoadRulesWithDefaults(ctx context.Context, dbConstraints []serverinterfaces.ConstraintRecord) error {
	if m.loadRulesFunc != nil {
		return m.loadRulesFunc(ctx, dbConstraints)
	}
	return nil
}

func (m *mockRuleService) ValidateArchitecture(ctx context.Context, architecture interface{}) (map[string]interface{}, error) {
	if m.validateArchitectureFunc != nil {
		return m.validateArchitectureFunc(ctx, architecture)
	}
	// Return empty valid results by default
	return map[string]interface{}{}, nil
}

func TestArchitectureService_MapFromDiagram(t *testing.T) {
	ruleService := &mockRuleService{}
	service := NewArchitectureService(ruleService, slog.Default())
	ctx := context.Background()

	// Create a valid diagram graph
	diagramGraph := &graph.DiagramGraph{
		Nodes: map[string]*graph.Node{
			"vpc-1": {
				ID:           "vpc-1",
				Type:         "resourceNode",
				ResourceType: "vpc",
				Label:        "Main VPC",
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
		provider  resource.CloudProvider
		wantError bool
	}{
		{
			name:      "valid graph with AWS provider",
			graph:     diagramGraph,
			provider:  resource.AWS,
			wantError: false,
		},
		{
			name:      "nil graph",
			graph:     nil,
			provider:  resource.AWS,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.MapFromDiagram(ctx, tt.graph, tt.provider)
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

func TestArchitectureService_ValidateRules(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		ruleService serverinterfaces.RuleService
		wantError   bool
		wantValid   bool
		setupMock   func(*mockRuleService)
	}{
		{
			name:        "valid architecture with no rule service",
			ruleService: nil,
			wantError:   false,
			wantValid:   true,
		},
		{
			name:        "valid architecture with rule service",
			ruleService: &mockRuleService{},
			wantError:   false,
			wantValid:   true,
			setupMock: func(m *mockRuleService) {
				m.validateArchitectureFunc = func(ctx context.Context, architecture interface{}) (map[string]interface{}, error) {
					return map[string]interface{}{}, nil
				}
			},
		},
		{
			name:        "rule service returns error",
			ruleService: &mockRuleService{},
			wantError:   true,
			setupMock: func(m *mockRuleService) {
				m.validateArchitectureFunc = func(ctx context.Context, architecture interface{}) (map[string]interface{}, error) {
					return nil, &mockError{message: "validation failed"}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil && tt.ruleService != nil {
				tt.setupMock(tt.ruleService.(*mockRuleService))
			}

			service := NewArchitectureService(tt.ruleService, slog.Default())

			// Create a minimal architecture for testing
			// Note: This is a simplified test - in a real scenario, you'd create a full architecture
			// For now, we'll test with nil to check error handling
			_, err := service.ValidateRules(ctx, nil, resource.AWS)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestArchitectureService_GetSortedResources(t *testing.T) {
	service := NewArchitectureService(nil, slog.Default())
	ctx := context.Background()

	tests := []struct {
		name      string
		arch      *architecture.Architecture
		wantError bool
	}{
		{
			name:      "nil architecture",
			arch:      nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetSortedResources(ctx, tt.arch)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// mockError is a simple error implementation
type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}
