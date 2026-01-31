package server

import (
	"fmt"

	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/architecture" // Register AWS architecture generator
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/orchestrator"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/services"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// Server represents the service layer server with all dependencies wired
type Server struct {
	// Services
	DiagramService      serverinterfaces.DiagramService
	ArchitectureService serverinterfaces.ArchitectureService
	CodegenService      serverinterfaces.CodegenService
	ProjectService      serverinterfaces.ProjectService

	// Orchestrator
	PipelineOrchestrator serverinterfaces.PipelineOrchestrator

	// Repositories (kept for reference, but services use interfaces)
	projectRepo        *repository.ProjectRepository
	resourceRepo       *repository.ResourceRepository
	resourceTypeRepo   *repository.ResourceTypeRepository
	containmentRepo    *repository.ResourceContainmentRepository
	dependencyRepo     *repository.ResourceDependencyRepository
	dependencyTypeRepo *repository.DependencyTypeRepository
	userRepo           *repository.UserRepository
	iacTargetRepo      *repository.IACTargetRepository
}

// NewServer creates a new server with all dependencies wired
func NewServer() (*Server, error) {
	// Initialize repositories
	projectRepo, err := repository.NewProjectRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create project repository: %w", err)
	}

	resourceRepo, err := repository.NewResourceRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create resource repository: %w", err)
	}

	resourceTypeRepo, err := repository.NewResourceTypeRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create resource type repository: %w", err)
	}

	containmentRepo, err := repository.NewResourceContainmentRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create containment repository: %w", err)
	}

	dependencyRepo, err := repository.NewResourceDependencyRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create dependency repository: %w", err)
	}

	dependencyTypeRepo, err := repository.NewDependencyTypeRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create dependency type repository: %w", err)
	}

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create user repository: %w", err)
	}

	iacTargetRepo, err := repository.NewIACTargetRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create iac target repository: %w", err)
	}

	// Create repository adapters
	projectRepoAdapter := &services.ProjectRepositoryAdapter{Repo: projectRepo}
	resourceRepoAdapter := &services.ResourceRepositoryAdapter{Repo: resourceRepo}
	resourceTypeRepoAdapter := &services.ResourceTypeRepositoryAdapter{Repo: resourceTypeRepo}
	containmentRepoAdapter := &services.ResourceContainmentRepositoryAdapter{Repo: containmentRepo}
	dependencyRepoAdapter := &services.ResourceDependencyRepositoryAdapter{Repo: dependencyRepo}
	dependencyTypeRepoAdapter := &services.DependencyTypeRepositoryAdapter{Repo: dependencyTypeRepo}
	userRepoAdapter := &services.UserRepositoryAdapter{Repo: userRepo}
	iacTargetRepoAdapter := &services.IACTargetRepositoryAdapter{Repo: iacTargetRepo}

	// Initialize services
	diagramService := services.NewDiagramService()

	// Create rule service adapter
	ruleService := services.NewAWSRuleServiceAdapter()

	architectureService := services.NewArchitectureService(ruleService)
	codegenService := services.NewCodegenService()

	projectService := services.NewProjectService(
		projectRepoAdapter,
		resourceRepoAdapter,
		resourceTypeRepoAdapter,
		containmentRepoAdapter,
		dependencyRepoAdapter,
		dependencyTypeRepoAdapter,
		userRepoAdapter,
		iacTargetRepoAdapter,
	)

	// Create pipeline orchestrator
	pipelineOrchestrator := orchestrator.NewPipelineOrchestrator(
		diagramService,
		architectureService,
		codegenService,
		projectService,
	)

	return &Server{
		DiagramService:      diagramService,
		ArchitectureService: architectureService,
		CodegenService:      codegenService,
		ProjectService:      projectService,
		PipelineOrchestrator: pipelineOrchestrator,
		projectRepo:          projectRepo,
		resourceRepo:         resourceRepo,
		resourceTypeRepo:     resourceTypeRepo,
		containmentRepo:      containmentRepo,
		dependencyRepo:       dependencyRepo,
		dependencyTypeRepo:   dependencyTypeRepo,
		userRepo:             userRepo,
		iacTargetRepo:        iacTargetRepo,
	}, nil
}
