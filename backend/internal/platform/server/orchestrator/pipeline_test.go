package orchestrator

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// Mock services for testing
type mockDiagramService struct {
	parseFunc    func(ctx context.Context, jsonData []byte) (*graph.DiagramGraph, error)
	validateFunc func(ctx context.Context, graph *graph.DiagramGraph, opts *validator.ValidationOptions) (*validator.ValidationResult, error)
}

func (m *mockDiagramService) Parse(ctx context.Context, jsonData []byte) (*graph.DiagramGraph, error) {
	if m.parseFunc != nil {
		return m.parseFunc(ctx, jsonData)
	}
	return &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: []*graph.Edge{},
	}, nil
}

func (m *mockDiagramService) Validate(ctx context.Context, graph *graph.DiagramGraph, opts *validator.ValidationOptions) (*validator.ValidationResult, error) {
	if m.validateFunc != nil {
		return m.validateFunc(ctx, graph, opts)
	}
	return &validator.ValidationResult{
		Valid:    true,
		Errors:   []*validator.ValidationError{},
		Warnings: []*validator.ValidationError{},
	}, nil
}

type mockArchitectureService struct {
	mapFromDiagramFunc func(ctx context.Context, graph *graph.DiagramGraph, provider resource.CloudProvider) (*architecture.Architecture, error)
	validateRulesFunc  func(ctx context.Context, arch *architecture.Architecture, provider resource.CloudProvider) (*serverinterfaces.RuleValidationResult, error)
	getSortedFunc      func(ctx context.Context, arch *architecture.Architecture) ([]*resource.Resource, error)
}

func (m *mockArchitectureService) MapFromDiagram(ctx context.Context, graph *graph.DiagramGraph, provider resource.CloudProvider) (*architecture.Architecture, error) {
	if m.mapFromDiagramFunc != nil {
		return m.mapFromDiagramFunc(ctx, graph, provider)
	}
	return &architecture.Architecture{
		Provider:     provider,
		Region:       "us-east-1",
		Resources:    []*resource.Resource{},
		Containments: make(map[string][]string),
		Dependencies: make(map[string][]string),
	}, nil
}

func (m *mockArchitectureService) ValidateRules(ctx context.Context, arch *architecture.Architecture, provider resource.CloudProvider) (*serverinterfaces.RuleValidationResult, error) {
	if m.validateRulesFunc != nil {
		return m.validateRulesFunc(ctx, arch, provider)
	}
	return &serverinterfaces.RuleValidationResult{
		Valid:   true,
		Results: make(map[string]*serverinterfaces.ResourceValidationResult),
	}, nil
}

func (m *mockArchitectureService) GetSortedResources(ctx context.Context, arch *architecture.Architecture) ([]*resource.Resource, error) {
	if m.getSortedFunc != nil {
		return m.getSortedFunc(ctx, arch)
	}
	return []*resource.Resource{}, nil
}

type mockCodegenService struct {
	generateFunc func(ctx context.Context, arch *architecture.Architecture, engine string) (*iac.Output, error)
}

func (m *mockCodegenService) Generate(ctx context.Context, arch *architecture.Architecture, engine string) (*iac.Output, error) {
	if m.generateFunc != nil {
		return m.generateFunc(ctx, arch, engine)
	}
	return &iac.Output{
		Files: []iac.GeneratedFile{},
	}, nil
}

func (m *mockCodegenService) SupportedEngines() []string {
	return []string{"terraform"}
}

type mockProjectService struct {
	createFunc                 func(ctx context.Context, req *serverinterfaces.CreateProjectRequest) (*models.Project, error)
	getByIDFunc                func(ctx context.Context, id uuid.UUID) (*models.Project, error)
	persistArchFunc            func(ctx context.Context, projectID uuid.UUID, arch *architecture.Architecture, diagramGraph interface{}) error
	persistArchWithPricingFunc func(ctx context.Context, projectID uuid.UUID, arch *architecture.Architecture, diagramGraph interface{}, pricingDuration time.Duration) (*serverinterfaces.ArchitecturePersistResult, error)
	loadArchFunc               func(ctx context.Context, projectID uuid.UUID) (*architecture.Architecture, error)
	getProjectPricingFunc      func(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectPricing, error)
	listFunc                   func(ctx context.Context, userID uuid.UUID, page, limit int, sort, order, search string) ([]*models.Project, int64, error)
	duplicateFunc              func(ctx context.Context, projectID uuid.UUID, name string) (*models.Project, *models.ProjectVersion, error)
	updateMetadataFunc         func(ctx context.Context, project *models.Project) (*models.Project, error)
	deleteFunc                 func(ctx context.Context, id uuid.UUID) error
	getArchFunc                func(ctx context.Context, projectID uuid.UUID) (*dto.ArchitectureResponse, error)
	// Version CRUD
	createVersionFunc       func(ctx context.Context, projectID uuid.UUID, req *serverinterfaces.CreateVersionRequest) (*serverinterfaces.ProjectVersionDetail, error)
	getVersionsFunc         func(ctx context.Context, projectID uuid.UUID) ([]*serverinterfaces.ProjectVersionSummary, error)
	getLatestVersionFunc    func(ctx context.Context, projectID uuid.UUID) (*serverinterfaces.ProjectVersionDetail, error)
	getVersionByIDFunc      func(ctx context.Context, projectID, versionID uuid.UUID) (*serverinterfaces.ProjectVersionDetail, error)
	deleteVersionFunc       func(ctx context.Context, projectID, versionID uuid.UUID) error
	validateVersionArchFunc func(ctx context.Context, versionID uuid.UUID) (*dto.ValidationResponse, error)
}

func (m *mockProjectService) Create(ctx context.Context, req *serverinterfaces.CreateProjectRequest) (*models.Project, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, req)
	}
	return &models.Project{
		ID:            uuid.New(),
		Name:          req.Name,
		CloudProvider: req.CloudProvider,
		Region:        req.Region,
	}, nil
}

func (m *mockProjectService) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return &models.Project{
		ID: id,
	}, nil
}

func (m *mockProjectService) PersistArchitecture(ctx context.Context, projectID uuid.UUID, arch *architecture.Architecture, diagramGraph interface{}) error {
	if m.persistArchFunc != nil {
		return m.persistArchFunc(ctx, projectID, arch, diagramGraph)
	}
	return nil
}

func (m *mockProjectService) PersistArchitectureWithPricing(ctx context.Context, projectID uuid.UUID, arch *architecture.Architecture, diagramGraph interface{}, pricingDuration time.Duration) (*serverinterfaces.ArchitecturePersistResult, error) {
	if m.persistArchWithPricingFunc != nil {
		return m.persistArchWithPricingFunc(ctx, projectID, arch, diagramGraph, pricingDuration)
	}
	return &serverinterfaces.ArchitecturePersistResult{
		ResourceIDMapping: make(map[string]uuid.UUID),
		PricingEstimate:   nil,
	}, nil
}

func (m *mockProjectService) LoadArchitecture(ctx context.Context, projectID uuid.UUID) (*architecture.Architecture, error) {
	if m.loadArchFunc != nil {
		return m.loadArchFunc(ctx, projectID)
	}
	return &architecture.Architecture{
		Resources:    []*resource.Resource{},
		Containments: make(map[string][]string),
		Dependencies: make(map[string][]string),
	}, nil
}

func (m *mockProjectService) GetProjectPricing(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectPricing, error) {
	if m.getProjectPricingFunc != nil {
		return m.getProjectPricingFunc(ctx, projectID)
	}
	return []*models.ProjectPricing{}, nil
}

func (m *mockProjectService) List(ctx context.Context, userID uuid.UUID, page, limit int, sort, order, search string) ([]*models.Project, int64, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, userID, page, limit, sort, order, search)
	}
	return []*models.Project{}, 0, nil
}

func (m *mockProjectService) Duplicate(ctx context.Context, projectID uuid.UUID, name string) (*models.Project, *models.ProjectVersion, error) {
	if m.duplicateFunc != nil {
		return m.duplicateFunc(ctx, projectID, name)
	}
	project := &models.Project{ID: uuid.New(), Name: name}
	version := &models.ProjectVersion{ID: uuid.New(), ProjectID: project.ID, VersionNumber: 1}
	return project, version, nil
}

func (m *mockProjectService) GetVersions(ctx context.Context, projectID uuid.UUID) ([]*serverinterfaces.ProjectVersionSummary, error) {
	if m.getVersionsFunc != nil {
		return m.getVersionsFunc(ctx, projectID)
	}
	return []*serverinterfaces.ProjectVersionSummary{}, nil
}

func (m *mockProjectService) CreateVersion(ctx context.Context, projectID uuid.UUID, req *serverinterfaces.CreateVersionRequest) (*serverinterfaces.ProjectVersionDetail, error) {
	if m.createVersionFunc != nil {
		return m.createVersionFunc(ctx, projectID, req)
	}
	return &serverinterfaces.ProjectVersionDetail{
		ProjectVersionSummary: serverinterfaces.ProjectVersionSummary{
			ID:            uuid.New(),
			ProjectID:     projectID,
			VersionNumber: 1,
		},
	}, nil
}

func (m *mockProjectService) GetLatestVersion(ctx context.Context, projectID uuid.UUID) (*serverinterfaces.ProjectVersionDetail, error) {
	if m.getLatestVersionFunc != nil {
		return m.getLatestVersionFunc(ctx, projectID)
	}
	return &serverinterfaces.ProjectVersionDetail{
		ProjectVersionSummary: serverinterfaces.ProjectVersionSummary{ID: uuid.New(), ProjectID: projectID, VersionNumber: 1},
	}, nil
}

func (m *mockProjectService) GetVersionByID(ctx context.Context, projectID, versionID uuid.UUID) (*serverinterfaces.ProjectVersionDetail, error) {
	if m.getVersionByIDFunc != nil {
		return m.getVersionByIDFunc(ctx, projectID, versionID)
	}
	return &serverinterfaces.ProjectVersionDetail{
		ProjectVersionSummary: serverinterfaces.ProjectVersionSummary{ID: versionID, ProjectID: projectID, VersionNumber: 1},
	}, nil
}

func (m *mockProjectService) DeleteVersion(ctx context.Context, projectID, versionID uuid.UUID) error {
	if m.deleteVersionFunc != nil {
		return m.deleteVersionFunc(ctx, projectID, versionID)
	}
	return nil
}

func (m *mockProjectService) UpdateMetadata(ctx context.Context, project *models.Project) (*models.Project, error) {
	if m.updateMetadataFunc != nil {
		return m.updateMetadataFunc(ctx, project)
	}
	return project, nil
}

func (m *mockProjectService) ValidateVersionArchitecture(ctx context.Context, versionID uuid.UUID) (*dto.ValidationResponse, error) {
	if m.validateVersionArchFunc != nil {
		return m.validateVersionArchFunc(ctx, versionID)
	}
	return &dto.ValidationResponse{Valid: true, Errors: []dto.ValidationIssue{}, Warnings: []dto.ValidationIssue{}}, nil
}

func (m *mockProjectService) GetArchitecture(ctx context.Context, projectID uuid.UUID) (*dto.ArchitectureResponse, error) {
	if m.getArchFunc != nil {
		return m.getArchFunc(ctx, projectID)
	}
	return &dto.ArchitectureResponse{}, nil
}

func (m *mockProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func TestPipelineOrchestrator_ProcessDiagram(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		req        *serverinterfaces.ProcessDiagramRequest
		wantError  bool
		setupMocks func(*mockDiagramService, *mockArchitectureService, *mockProjectService)
	}{
		{
			name: "successful processing",
			req: &serverinterfaces.ProcessDiagramRequest{
				JSONData:      []byte(`{"nodes": [], "edges": []}`),
				UserID:        uuid.New(),
				ProjectName:   "Test Project",
				IACToolID:     1,
				CloudProvider: "aws",
				Region:        "us-east-1",
			},
			wantError: false,
		},
		{
			name:      "nil request",
			req:       nil,
			wantError: true,
		},
		{
			name: "parse error",
			req: &serverinterfaces.ProcessDiagramRequest{
				JSONData:      []byte(`invalid json`),
				UserID:        uuid.New(),
				ProjectName:   "Test Project",
				IACToolID:     1,
				CloudProvider: "aws",
				Region:        "us-east-1",
			},
			wantError: true,
			setupMocks: func(ds *mockDiagramService, as *mockArchitectureService, ps *mockProjectService) {
				ds.parseFunc = func(ctx context.Context, jsonData []byte) (*graph.DiagramGraph, error) {
					return nil, errors.New("parse error")
				}
			},
		},
		{
			name: "validation error",
			req: &serverinterfaces.ProcessDiagramRequest{
				JSONData:      []byte(`{"nodes": [], "edges": []}`),
				UserID:        uuid.New(),
				ProjectName:   "Test Project",
				IACToolID:     1,
				CloudProvider: "aws",
				Region:        "us-east-1",
			},
			wantError: true,
			setupMocks: func(ds *mockDiagramService, as *mockArchitectureService, ps *mockProjectService) {
				ds.validateFunc = func(ctx context.Context, graph *graph.DiagramGraph, opts *validator.ValidationOptions) (*validator.ValidationResult, error) {
					return &validator.ValidationResult{
						Valid: false,
						Errors: []*validator.ValidationError{
							{Code: "ERROR", Message: "Validation failed"},
						},
					}, nil
				}
			},
		},
		{
			name: "architecture mapping error",
			req: &serverinterfaces.ProcessDiagramRequest{
				JSONData:      []byte(`{"nodes": [], "edges": []}`),
				UserID:        uuid.New(),
				ProjectName:   "Test Project",
				IACToolID:     1,
				CloudProvider: "aws",
				Region:        "us-east-1",
			},
			wantError: true,
			setupMocks: func(ds *mockDiagramService, as *mockArchitectureService, ps *mockProjectService) {
				as.mapFromDiagramFunc = func(ctx context.Context, graph *graph.DiagramGraph, provider resource.CloudProvider) (*architecture.Architecture, error) {
					return nil, errors.New("mapping error")
				}
			},
		},
		{
			name: "rule validation error",
			req: &serverinterfaces.ProcessDiagramRequest{
				JSONData:      []byte(`{"nodes": [], "edges": []}`),
				UserID:        uuid.New(),
				ProjectName:   "Test Project",
				IACToolID:     1,
				CloudProvider: "aws",
				Region:        "us-east-1",
			},
			wantError: true,
			setupMocks: func(ds *mockDiagramService, as *mockArchitectureService, ps *mockProjectService) {
				as.validateRulesFunc = func(ctx context.Context, arch *architecture.Architecture, provider resource.CloudProvider) (*serverinterfaces.RuleValidationResult, error) {
					return &serverinterfaces.RuleValidationResult{
						Valid:  false,
						Errors: []string{"Rule validation failed"},
					}, nil
				}
			},
		},
		{
			name: "project creation error",
			req: &serverinterfaces.ProcessDiagramRequest{
				JSONData:      []byte(`{"nodes": [], "edges": []}`),
				UserID:        uuid.New(),
				ProjectName:   "Test Project",
				IACToolID:     1,
				CloudProvider: "aws",
				Region:        "us-east-1",
			},
			wantError: true,
			setupMocks: func(ds *mockDiagramService, as *mockArchitectureService, ps *mockProjectService) {
				ps.createFunc = func(ctx context.Context, req *serverinterfaces.CreateProjectRequest) (*models.Project, error) {
					return nil, errors.New("project creation error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagramService := &mockDiagramService{}
			archService := &mockArchitectureService{}
			codegenService := &mockCodegenService{}
			projectService := &mockProjectService{}

			if tt.setupMocks != nil {
				tt.setupMocks(diagramService, archService, projectService)
			}

			orchestrator := NewPipelineOrchestrator(
				diagramService,
				archService,
				codegenService,
				projectService,
			)

			result, err := orchestrator.ProcessDiagram(ctx, tt.req)
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

func TestPipelineOrchestrator_GenerateCode(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		req        *serverinterfaces.GenerateCodeRequest
		wantError  bool
		setupMocks func(*mockProjectService, *mockArchitectureService, *mockCodegenService)
	}{
		{
			name:      "nil request",
			req:       nil,
			wantError: true,
		},
		{
			name: "project not found",
			req: &serverinterfaces.GenerateCodeRequest{
				ProjectID:     uuid.New(),
				Engine:        "terraform",
				CloudProvider: "aws",
			},
			wantError: true,
			setupMocks: func(ps *mockProjectService, as *mockArchitectureService, cs *mockCodegenService) {
				ps.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*models.Project, error) {
					return nil, errors.New("project not found")
				}
			},
		},
		{
			name: "load architecture error",
			req: &serverinterfaces.GenerateCodeRequest{
				ProjectID:     uuid.New(),
				Engine:        "terraform",
				CloudProvider: "aws",
			},
			wantError: true,
			setupMocks: func(ps *mockProjectService, as *mockArchitectureService, cs *mockCodegenService) {
				ps.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*models.Project, error) {
					return &models.Project{
						ID:            id,
						CloudProvider: "aws",
						Region:        "us-east-1",
					}, nil
				}
				ps.loadArchFunc = func(ctx context.Context, projectID uuid.UUID) (*architecture.Architecture, error) {
					return nil, errors.New("failed to load architecture")
				}
			},
		},
		{
			name: "successful code generation",
			req: &serverinterfaces.GenerateCodeRequest{
				ProjectID:     uuid.New(),
				Engine:        "terraform",
				CloudProvider: "aws",
			},
			wantError: false,
			setupMocks: func(ps *mockProjectService, as *mockArchitectureService, cs *mockCodegenService) {
				ps.getByIDFunc = func(ctx context.Context, id uuid.UUID) (*models.Project, error) {
					return &models.Project{
						ID:            id,
						CloudProvider: "aws",
						Region:        "us-east-1",
					}, nil
				}
				ps.loadArchFunc = func(ctx context.Context, projID uuid.UUID) (*architecture.Architecture, error) {
					return &architecture.Architecture{
						Resources:    []*resource.Resource{},
						Containments: make(map[string][]string),
						Dependencies: make(map[string][]string),
						Provider:     resource.AWS,
						Region:       "us-east-1",
					}, nil
				}
				as.validateRulesFunc = func(ctx context.Context, arch *architecture.Architecture, provider resource.CloudProvider) (*serverinterfaces.RuleValidationResult, error) {
					return &serverinterfaces.RuleValidationResult{
						Valid:   true,
						Results: make(map[string]*serverinterfaces.ResourceValidationResult),
					}, nil
				}
				cs.generateFunc = func(ctx context.Context, arch *architecture.Architecture, engine string) (*iac.Output, error) {
					return &iac.Output{
						Files: []iac.GeneratedFile{
							{Path: "main.tf", Content: "# Terraform code", Type: "hcl"},
						},
					}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagramService := &mockDiagramService{}
			archService := &mockArchitectureService{}
			codegenService := &mockCodegenService{}
			projectService := &mockProjectService{}

			if tt.setupMocks != nil {
				tt.setupMocks(projectService, archService, codegenService)
			}

			orchestrator := NewPipelineOrchestrator(
				diagramService,
				archService,
				codegenService,
				projectService,
			)

			output, err := orchestrator.GenerateCode(ctx, tt.req)
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
			if output == nil {
				t.Errorf("expected output but got nil")
			}
		})
	}
}
