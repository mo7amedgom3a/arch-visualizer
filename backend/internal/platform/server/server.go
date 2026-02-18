package server

import (
	"context"
	"fmt"

	"log/slog"

	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/architecture" // Register AWS architecture generator
	awscompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	awsdatabase "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/iam"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/networking"
	awsstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/storage"
	infrastructurerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/infrastructure"
	pricingrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/pricing"
	projectrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/project"
	resourcerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/resource"
	userrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/user"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/orchestrator"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/services"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/seeder"
)

// Server represents the service layer server with all dependencies wired
type Server struct {
	// Services
	DiagramService          serverinterfaces.DiagramService
	ArchitectureService     serverinterfaces.ArchitectureService
	CodegenService          serverinterfaces.CodegenService
	ProjectService          serverinterfaces.ProjectService
	PricingService          serverinterfaces.PricingService
	OptimizationService     serverinterfaces.OptimizationService
	UserService             serverinterfaces.UserService
	StaticDataService       serverinterfaces.StaticDataService
	ResourceMetadataService serverinterfaces.ResourceMetadataService
	IAMService              iam.AWSIAMService

	// Orchestrator
	PipelineOrchestrator serverinterfaces.PipelineOrchestrator

	// Repositories (kept for reference, but services use interfaces)
	projectRepo        *projectrepo.ProjectRepository
	resourceRepo       *resourcerepo.ResourceRepository
	resourceTypeRepo   *resourcerepo.ResourceTypeRepository
	containmentRepo    *resourcerepo.ResourceContainmentRepository
	dependencyRepo     *resourcerepo.ResourceDependencyRepository
	dependencyTypeRepo *resourcerepo.DependencyTypeRepository
	userRepo           *userrepo.UserRepository
	iacTargetRepo      *infrastructurerepo.IACTargetRepository
	variableRepo       *projectrepo.ProjectVariableRepository
	outputRepo         *projectrepo.ProjectOutputRepository
	pricingRepo        *pricingrepo.PricingRepository
	pricingRateRepo    *pricingrepo.PricingRateRepository
	hiddenDepRepo      *resourcerepo.HiddenDependencyRepository
	constraintRepo     *resourcerepo.ResourceConstraintRepository
}

// NewServer creates a new server with all dependencies wired
func NewServer(logger *slog.Logger) (*Server, error) {
	// Initialize repositories
	projectRepo, err := projectrepo.NewProjectRepository(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create project repository: %w", err)
	}

	resourceRepo, err := resourcerepo.NewResourceRepository(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource repository: %w", err)
	}

	resourceTypeRepo, err := resourcerepo.NewResourceTypeRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create resource type repository: %w", err)
	}

	containmentRepo, err := resourcerepo.NewResourceContainmentRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create containment repository: %w", err)
	}

	dependencyRepo, err := resourcerepo.NewResourceDependencyRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create dependency repository: %w", err)
	}

	dependencyTypeRepo, err := resourcerepo.NewDependencyTypeRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create dependency type repository: %w", err)
	}

	userRepo, err := userrepo.NewUserRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create user repository: %w", err)
	}

	iacTargetRepo, err := infrastructurerepo.NewIACTargetRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create iac target repository: %w", err)
	}

	pricingRepo, err := pricingrepo.NewPricingRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create pricing repository: %w", err)
	}

	pricingRateRepo, err := pricingrepo.NewPricingRateRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create pricing rate repository: %w", err)
	}

	hiddenDepRepo, err := resourcerepo.NewHiddenDependencyRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create hidden dependency repository: %w", err)
	}

	constraintRepo, err := resourcerepo.NewResourceConstraintRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create constraint repository: %w", err)
	}
	variableRepo, err := projectrepo.NewProjectVariableRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create variable repository: %w", err)
	}
	outputRepo, err := projectrepo.NewProjectOutputRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create output repository: %w", err)
	}

	// Create repository adapters
	projectRepoAdapter := &services.ProjectRepositoryAdapter{Repo: projectRepo}

	// Create project version repository
	versionRepo, err := projectrepo.NewProjectVersionRepository()
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
	diagramService := services.NewDiagramService(logger)

	// Create rule service adapter
	ruleService := services.NewAWSRuleServiceAdapter()

	architectureService := services.NewArchitectureService(ruleService, logger)
	codegenService := services.NewCodegenService(logger)
	optimizationService := services.NewOptimizationService()

	// Create pricing service with DB-driven rates and hidden dependencies
	pricingService := services.NewPricingServiceWithRepos(
		pricingRepoAdapter,
		pricingRateRepo,
		hiddenDepRepo,
	)

	// Create constraint service
	constraintService := services.NewConstraintService(constraintRepo)

	// Run seeder for resource constraints
	// We need a background context for seeding
	seedCtx := context.Background()
	if err := seeder.SeedResourceConstraints(seedCtx, constraintRepo, resourceTypeRepo); err != nil {
		fmt.Printf("Warning: Failed to seed resource constraints: %v\n", err)
		// Continue anyway, don't crash
	}

	// Load constraints into rule service
	// We do this after seeding to ensure we have the latest constraints
	if constraints, err := constraintService.GetAllConstraints(seedCtx); err != nil {
		fmt.Printf("Warning: Failed to load resource constraints: %v\n", err)
	} else {
		if err := ruleService.LoadRulesWithDefaults(seedCtx, constraints); err != nil {
			fmt.Printf("Warning: Failed to apply resource constraints: %v\n", err)
		} else {
			fmt.Printf("✓ Loaded %d resource constraints from database\n", len(constraints))
		}
	}

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

	// Create resource metadata services for all service categories
	networkingMetaSvc := awsnetworking.NewNetworkingMetadataService()
	computeMetaSvc := awscompute.NewComputeMetadataService()
	storageMetaSvc := awsstorage.NewStorageMetadataService()
	databaseMetaSvc := awsdatabase.NewDatabaseMetadataService()
	iamMetaSvc := iam.NewIAMMetadataService()
	resourceMetadataService := services.NewResourceMetadataService(
		networkingMetaSvc,
		computeMetaSvc,
		storageMetaSvc,
		databaseMetaSvc,
		iamMetaSvc,
	)

	iamService := iam.NewIAMService()

	return &Server{
		DiagramService:      diagramService,
		ArchitectureService: architectureService,
		CodegenService:      codegenService,
		ProjectService:      projectService,

		PricingService:          pricingService,
		OptimizationService:     optimizationService,
		UserService:             userService,
		StaticDataService:       staticDataService,
		ResourceMetadataService: resourceMetadataService,
		IAMService:              iamService,
		PipelineOrchestrator:    pipelineOrchestrator,
		projectRepo:             projectRepo,
		resourceRepo:            resourceRepo,
		resourceTypeRepo:        resourceTypeRepo,
		containmentRepo:         containmentRepo,
		dependencyRepo:          dependencyRepo,
		dependencyTypeRepo:      dependencyTypeRepo,
		userRepo:                userRepo,
		iacTargetRepo:           iacTargetRepo,
		variableRepo:            variableRepo,
		outputRepo:              outputRepo,
		pricingRepo:             pricingRepo,
		pricingRateRepo:         pricingRateRepo,
		hiddenDepRepo:           hiddenDepRepo,
		constraintRepo:          constraintRepo,
	}, nil
}
