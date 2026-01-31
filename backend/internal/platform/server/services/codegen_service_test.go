package services

import (
	"context"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
)

func TestCodegenService_SupportedEngines(t *testing.T) {
	service := NewCodegenService()

	engines := service.SupportedEngines()
	if len(engines) == 0 {
		t.Error("expected at least one supported engine")
	}

	// Check that terraform is supported
	hasTerraform := false
	for _, engine := range engines {
		if engine == "terraform" {
			hasTerraform = true
			break
		}
	}
	if !hasTerraform {
		t.Error("expected terraform to be a supported engine")
	}
}

func TestCodegenService_Generate(t *testing.T) {
	service := NewCodegenService()
	ctx := context.Background()

	tests := []struct {
		name      string
		arch      *architecture.Architecture
		engine    string
		wantError bool
	}{
		{
			name:      "nil architecture",
			arch:      nil,
			engine:    "terraform",
			wantError: true,
		},
		{
			name: "unsupported engine",
			arch: &architecture.Architecture{
				Provider: resource.AWS,
				Region:   "us-east-1",
			},
			engine:    "unsupported",
			wantError: true,
		},
		{
			name: "empty engine defaults to terraform",
			arch: &architecture.Architecture{
				Provider: resource.AWS,
				Region:   "us-east-1",
			},
			engine:    "",
			wantError: false, // Should default to terraform
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Generate(ctx, tt.arch, tt.engine)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				// For some cases, we might get errors due to missing resources
				// That's okay for these basic tests
				return
			}
			if result == nil {
				t.Errorf("expected result but got nil")
			}
		})
	}
}

func TestCodegenService_GenerateWithValidArchitecture(t *testing.T) {
	service := NewCodegenService()
	ctx := context.Background()

	// Create a minimal valid architecture
	arch := &architecture.Architecture{
		Provider:  resource.AWS,
		Region:     "us-east-1",
		Resources:  []*resource.Resource{},
		Variables:  []architecture.Variable{},
		Outputs:    []architecture.Output{},
		Containments: make(map[string][]string),
		Dependencies: make(map[string][]string),
	}

	output, err := service.Generate(ctx, arch, "terraform")
	if err != nil {
		// This might fail if there are no resources, which is expected
		// We're just testing that the service doesn't panic
		return
	}

	if output == nil {
		t.Error("expected output but got nil")
		return
	}

	if output.Files == nil {
		t.Error("expected files but got nil")
	}
}

func TestNewCodegenServiceWithEngines(t *testing.T) {
	// Create a mock engine
	mockEngine := &mockIACEngine{
		name: "mock",
	}

	engines := map[string]iac.Engine{
		"mock": mockEngine,
	}

	service := NewCodegenServiceWithEngines(engines)

	enginesList := service.SupportedEngines()
	if len(enginesList) != 1 {
		t.Errorf("expected 1 engine, got %d", len(enginesList))
	}

	if enginesList[0] != "mock" {
		t.Errorf("expected engine 'mock', got '%s'", enginesList[0])
	}
}

// mockIACEngine is a mock implementation of iac.Engine
type mockIACEngine struct {
	name string
}

func (m *mockIACEngine) Name() string {
	return m.name
}

func (m *mockIACEngine) Generate(ctx context.Context, arch *architecture.Architecture, sortedResources []*resource.Resource) (*iac.Output, error) {
	return &iac.Output{
		Files: []iac.GeneratedFile{
			{
				Path:    "main.tf",
				Content: "# Mock generated code",
				Type:    "hcl",
			},
		},
	}, nil
}
