package server

import (
	"fmt"

	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/architecture" // Register AWS architecture generator
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/iam"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/orchestrator"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/services"
)

// Server represents the service layer server with all dependencies wired
type Server struct {
	// Services
	DiagramService      serverinterfaces.DiagramService
	ArchitectureService serverinterfaces.ArchitectureService
	CodegenService      serverinterfaces.CodegenService
	ProjectService      serverinterfaces.ProjectService
	PricingService      serverinterfaces.PricingService
	UserService         serverinterfaces.UserService
	StaticDataService   serverinterfaces.StaticDataService
	IAMService          iam.AWSIAMService

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
	variableRepo       *repository.ProjectVariableRepository
	outputRepo         *repository.ProjectOutputRepository
	pricingRepo        *repository.PricingRepository
	pricingRateRepo    *repository.PricingRateRepository
	hiddenDepRepo      *repository.HiddenDependencyRepository
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

	pricingRepo, err := repository.NewPricingRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create pricing repository: %w", err)
	}

	pricingRateRepo, err := repository.NewPricingRateRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create pricing rate repository: %w", err)
	}

	hiddenDepRepo, err := repository.NewHiddenDependencyRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create hidden dependency repository: %w", err)
	}
	variableRepo, err := repository.NewProjectVariableRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create variable repository: %w", err)
	}
	outputRepo, err := repository.NewProjectOutputRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create output repository: %w", err)
	}

	// Create repository adapters
	projectRepoAdapter := &services.ProjectRepositoryAdapter{Repo: projectRepo}

	// Create project version repository
	versionRepo, err := repository.NewProjectVersionRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create project version repository: %w", err)
	}
	versionRepoAdapter := &services.ProjectVersionRepositoryAdapter{Repo: versionRepo}

	resourceRepoAdapter := &services.ResourceRepositoryAdapter{Repo: resourceRepo}
	resourceTypeRepoAdapter := &services.ResourceTypeRepositoryAdapter{Repo: resourceTypeRepo}
	containmentRepoAdapter := &services.ResourceContainmentRepositoryAdapter{Repo: containmentRepo}
	dependencyRepoAdapter := &services.ResourceDependencyRepositoryAdapter{Repo: dependencyRepo}
	dependencyTypeRepoAdapter := &services.DependencyTypeRepositoryAdapter{Repo: dependencyTypeRepo}
	userRepoAdapter := &services.UserRepositoryAdapter{Repo: userRepo}
	iacTargetRepoAdapter := &services.IACTargetRepositoryAdapter{Repo: iacTargetRepo}
	variableRepoAdapter := &services.ProjectVariableRepositoryAdapter{Repo: variableRepo}
	outputRepoAdapter := &services.ProjectOutputRepositoryAdapter{Repo: outputRepo}
	pricingRepoAdapter := &services.PricingRepositoryAdapter{Repo: pricingRepo}

	// Initialize services
	diagramService := services.NewDiagramService()

	// Create rule service adapter
	ruleService := services.NewAWSRuleServiceAdapter()

	architectureService := services.NewArchitectureService(ruleService)
	codegenService := services.NewCodegenService()

	// Create pricing service with DB-driven rates and hidden dependencies
	pricingService := services.NewPricingServiceWithRepos(
		pricingRepoAdapter,
		pricingRateRepo,
		hiddenDepRepo,
	)

	// Create project service with pricing support
	projectService := services.NewProjectServiceWithPricing(
		projectRepoAdapter,
		versionRepoAdapter,
		resourceRepoAdapter,
		resourceTypeRepoAdapter,
		containmentRepoAdapter,
		dependencyRepoAdapter,
		dependencyTypeRepoAdapter,
		userRepoAdapter,
		iacTargetRepoAdapter,
		variableRepoAdapter,
		outputRepoAdapter,
		pricingService,
	)

	// Create pipeline orchestrator
	pipelineOrchestrator := orchestrator.NewPipelineOrchestrator(
		diagramService,
		architectureService,
		codegenService,
		projectService,
	)

	userService := services.NewUserService(userRepoAdapter)
	staticDataService := services.NewStaticDataService(resourceTypeRepoAdapter)
	iamService := iam.NewIAMService()

	return &Server{
		DiagramService:       diagramService,
		ArchitectureService:  architectureService,
		CodegenService:       codegenService,
		ProjectService:       projectService,
		PricingService:       pricingService,
		UserService:          userService,
		StaticDataService:    staticDataService,
		IAMService:           iamService,
		PipelineOrchestrator: pipelineOrchestrator,
		projectRepo:          projectRepo,
		resourceRepo:         resourceRepo,
		resourceTypeRepo:     resourceTypeRepo,
		containmentRepo:      containmentRepo,
		dependencyRepo:       dependencyRepo,
		dependencyTypeRepo:   dependencyTypeRepo,
		userRepo:             userRepo,
		iacTargetRepo:        iacTargetRepo,
		variableRepo:         variableRepo,
		outputRepo:           outputRepo,
		pricingRepo:          pricingRepo,
		pricingRateRepo:      pricingRateRepo,
		hiddenDepRepo:        hiddenDepRepo,
	}, nil
}
