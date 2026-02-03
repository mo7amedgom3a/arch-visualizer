package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
	"gorm.io/datatypes"
)

// ProjectServiceImpl implements ProjectService interface
type ProjectServiceImpl struct {
	projectRepo        serverinterfaces.ProjectRepository
	versionRepo        serverinterfaces.ProjectVersionRepository
	resourceRepo       serverinterfaces.ResourceRepository
	resourceTypeRepo   serverinterfaces.ResourceTypeRepository
	containmentRepo    serverinterfaces.ResourceContainmentRepository
	dependencyRepo     serverinterfaces.ResourceDependencyRepository
	dependencyTypeRepo serverinterfaces.DependencyTypeRepository
	userRepo           serverinterfaces.UserRepository
	iacTargetRepo      serverinterfaces.IACTargetRepository
	pricingService     serverinterfaces.PricingService
}

// NewProjectService creates a new project service
func NewProjectService(
	projectRepo serverinterfaces.ProjectRepository,
	versionRepo serverinterfaces.ProjectVersionRepository,
	resourceRepo serverinterfaces.ResourceRepository,
	resourceTypeRepo serverinterfaces.ResourceTypeRepository,
	containmentRepo serverinterfaces.ResourceContainmentRepository,
	dependencyRepo serverinterfaces.ResourceDependencyRepository,
	dependencyTypeRepo serverinterfaces.DependencyTypeRepository,
	userRepo serverinterfaces.UserRepository,
	iacTargetRepo serverinterfaces.IACTargetRepository,
) serverinterfaces.ProjectService {
	return &ProjectServiceImpl{
		projectRepo:        projectRepo,
		versionRepo:        versionRepo,
		resourceRepo:       resourceRepo,
		resourceTypeRepo:   resourceTypeRepo,
		containmentRepo:    containmentRepo,
		dependencyRepo:     dependencyRepo,
		dependencyTypeRepo: dependencyTypeRepo,
		userRepo:           userRepo,
		iacTargetRepo:      iacTargetRepo,
		pricingService:     nil, // Will be set via SetPricingService
	}
}

// NewProjectServiceWithPricing creates a new project service with pricing support
func NewProjectServiceWithPricing(
	projectRepo serverinterfaces.ProjectRepository,
	versionRepo serverinterfaces.ProjectVersionRepository,
	resourceRepo serverinterfaces.ResourceRepository,
	resourceTypeRepo serverinterfaces.ResourceTypeRepository,
	containmentRepo serverinterfaces.ResourceContainmentRepository,
	dependencyRepo serverinterfaces.ResourceDependencyRepository,
	dependencyTypeRepo serverinterfaces.DependencyTypeRepository,
	userRepo serverinterfaces.UserRepository,
	iacTargetRepo serverinterfaces.IACTargetRepository,
	pricingService serverinterfaces.PricingService,
) serverinterfaces.ProjectService {
	return &ProjectServiceImpl{
		projectRepo:        projectRepo,
		versionRepo:        versionRepo,
		resourceRepo:       resourceRepo,
		resourceTypeRepo:   resourceTypeRepo,
		containmentRepo:    containmentRepo,
		dependencyRepo:     dependencyRepo,
		dependencyTypeRepo: dependencyTypeRepo,
		userRepo:           userRepo,
		iacTargetRepo:      iacTargetRepo,
		pricingService:     pricingService,
	}
}

// SetPricingService sets the pricing service (for backward compatibility)
func (s *ProjectServiceImpl) SetPricingService(pricingService serverinterfaces.PricingService) {
	s.pricingService = pricingService
}

// Create creates a new project
func (s *ProjectServiceImpl) Create(ctx context.Context, req *serverinterfaces.CreateProjectRequest) (*models.Project, error) {
	if req == nil {
		return nil, fmt.Errorf("create project request is nil")
	}

	project := &models.Project{
		ID:            uuid.New(),
		UserID:        req.UserID,
		InfraToolID:   req.IACTargetID,
		Name:          req.Name,
		Description:   req.Description,
		Tags:          req.Tags,
		CloudProvider: req.CloudProvider,
		Region:        req.Region,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

// GetByID retrieves a project by ID with related data
func (s *ProjectServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return s.projectRepo.FindByID(ctx, id)
}

// List retrieves projects with pagination and filtering
func (s *ProjectServiceImpl) List(ctx context.Context, userID uuid.UUID, page, limit int, sort, order, search string) ([]*models.Project, int64, error) {
	return s.projectRepo.FindAll(ctx, userID, page, limit, sort, order, search)
}

// Update updates an existing project
func (s *ProjectServiceImpl) Update(ctx context.Context, project *models.Project) error {
	if project == nil {
		return fmt.Errorf("project is nil")
	}
	// Add validation logic here if needed
	return s.projectRepo.Update(ctx, project)
}

// Delete deletes a project by ID
func (s *ProjectServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return s.projectRepo.Delete(ctx, id)
}

// PersistArchitecture persists an architecture to the database as part of a project
func (s *ProjectServiceImpl) PersistArchitecture(ctx context.Context, projectID uuid.UUID, arch *architecture.Architecture, diagramGraph interface{}) error {
	if arch == nil {
		return fmt.Errorf("architecture is nil")
	}

	// Start transaction
	tx, txCtx := s.projectRepo.BeginTransaction(ctx)
	defer func() {
		if r := recover(); r != nil {
			s.projectRepo.RollbackTransaction(tx)
			panic(r)
		}
	}()

	// Build ID mapping: domain resource ID -> database resource UUID
	domainIDToDBID := make(map[string]uuid.UUID)
	resourceTypeCache := make(map[string]uint) // cache resource type IDs

	// Create resources
	for _, res := range arch.Resources {
		// Get or cache resource type ID
		resourceTypeID, ok := resourceTypeCache[res.Type.Name]
		if !ok {
			resourceType, err := s.resourceTypeRepo.FindByNameAndProvider(txCtx, res.Type.Name, string(arch.Provider))
			if err != nil {
				// Create resource type if not exists
				// Note: This requires direct DB access, which we'll need to handle
				// For now, return error if resource type doesn't exist
				s.projectRepo.RollbackTransaction(tx)
				return fmt.Errorf("resource type %s not found for provider %s: %w", res.Type.Name, arch.Provider, err)
			}
			resourceTypeID = resourceType.ID
			resourceTypeCache[res.Type.Name] = resourceTypeID
		}

		// Get position from metadata
		posX, posY := 0, 0
		if pos, ok := res.Metadata["position"].(map[string]interface{}); ok {
			if x, ok := pos["x"].(float64); ok {
				posX = int(x)
			} else if x, ok := pos["x"].(int); ok {
				posX = x
			}
			if y, ok := pos["y"].(float64); ok {
				posY = int(y)
			} else if y, ok := pos["y"].(int); ok {
				posY = y
			}
		}

		// Get isVisualOnly from metadata
		isVisualOnly := false
		if v, ok := res.Metadata["isVisualOnly"].(bool); ok {
			isVisualOnly = v
		}

		// Convert metadata to JSON
		configJSON, err := json.Marshal(res.Metadata)
		if err != nil {
			s.projectRepo.RollbackTransaction(tx)
			return fmt.Errorf("marshal resource config: %w", err)
		}

		// Create resource
		dbResource := &models.Resource{
			ID:             uuid.New(),
			ProjectID:      projectID,
			ResourceTypeID: resourceTypeID,
			Name:           res.Name,
			PosX:           posX,
			PosY:           posY,
			IsVisualOnly:   isVisualOnly,
			Config:         datatypes.JSON(configJSON),
		}
		if err := s.resourceRepo.Create(txCtx, dbResource); err != nil {
			s.projectRepo.RollbackTransaction(tx)
			return fmt.Errorf("create resource %s: %w", res.Name, err)
		}

		domainIDToDBID[res.ID] = dbResource.ID
	}

	// Create containment relationships
	for parentID, childIDs := range arch.Containments {
		parentDBID, ok := domainIDToDBID[parentID]
		if !ok {
			continue
		}
		for _, childID := range childIDs {
			childDBID, ok := domainIDToDBID[childID]
			if !ok {
				continue
			}
			containment := &models.ResourceContainment{
				ParentResourceID: parentDBID,
				ChildResourceID:  childDBID,
			}
			if err := s.containmentRepo.Create(txCtx, containment); err != nil {
				s.projectRepo.RollbackTransaction(tx)
				return fmt.Errorf("create containment: %w", err)
			}
		}
	}

	// Create dependency relationships
	// Get or create "depends_on" dependency type
	dependencyType, err := s.dependencyTypeRepo.FindByName(txCtx, "depends_on")
	if err != nil {
		// Create dependency type if it doesn't exist
		// Note: This requires direct DB access, which we'll need to handle
		// For now, return error if dependency type doesn't exist
		s.projectRepo.RollbackTransaction(tx)
		return fmt.Errorf("dependency type 'depends_on' not found: %w", err)
	}

	for fromID, toIDs := range arch.Dependencies {
		fromDBID, ok := domainIDToDBID[fromID]
		if !ok {
			continue
		}
		for _, toID := range toIDs {
			toDBID, ok := domainIDToDBID[toID]
			if !ok {
				continue
			}
			dependency := &models.ResourceDependency{
				FromResourceID:   fromDBID,
				ToResourceID:     toDBID,
				DependencyTypeID: dependencyType.ID,
			}
			if err := s.dependencyRepo.Create(txCtx, dependency); err != nil {
				s.projectRepo.RollbackTransaction(tx)
				return fmt.Errorf("create dependency: %w", err)
			}
		}
	}

	// Commit transaction
	if err := s.projectRepo.CommitTransaction(tx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// PersistArchitectureWithPricing persists an architecture with pricing calculation
func (s *ProjectServiceImpl) PersistArchitectureWithPricing(ctx context.Context, projectID uuid.UUID, arch *architecture.Architecture, diagramGraph interface{}, pricingDuration time.Duration) (*serverinterfaces.ArchitecturePersistResult, error) {
	if arch == nil {
		return nil, fmt.Errorf("architecture is nil")
	}

	// Start transaction
	tx, txCtx := s.projectRepo.BeginTransaction(ctx)
	defer func() {
		if r := recover(); r != nil {
			s.projectRepo.RollbackTransaction(tx)
			panic(r)
		}
	}()

	// Build ID mapping: domain resource ID -> database resource UUID
	domainIDToDBID := make(map[string]uuid.UUID)
	resourceTypeCache := make(map[string]uint) // cache resource type IDs

	// Create resources
	for _, res := range arch.Resources {
		// Get or cache resource type ID
		resourceTypeID, ok := resourceTypeCache[res.Type.Name]
		if !ok {
			resourceType, err := s.resourceTypeRepo.FindByNameAndProvider(txCtx, res.Type.Name, string(arch.Provider))
			if err != nil {
				s.projectRepo.RollbackTransaction(tx)
				return nil, fmt.Errorf("resource type %s not found for provider %s: %w", res.Type.Name, arch.Provider, err)
			}
			resourceTypeID = resourceType.ID
			resourceTypeCache[res.Type.Name] = resourceTypeID
		}

		// Get position from metadata
		posX, posY := 0, 0
		if pos, ok := res.Metadata["position"].(map[string]interface{}); ok {
			if x, ok := pos["x"].(float64); ok {
				posX = int(x)
			} else if x, ok := pos["x"].(int); ok {
				posX = x
			}
			if y, ok := pos["y"].(float64); ok {
				posY = int(y)
			} else if y, ok := pos["y"].(int); ok {
				posY = y
			}
		}

		// Get isVisualOnly from metadata
		isVisualOnly := false
		if v, ok := res.Metadata["isVisualOnly"].(bool); ok {
			isVisualOnly = v
		}

		// Convert metadata to JSON
		configJSON, err := json.Marshal(res.Metadata)
		if err != nil {
			s.projectRepo.RollbackTransaction(tx)
			return nil, fmt.Errorf("marshal resource config: %w", err)
		}

		// Create resource
		dbResource := &models.Resource{
			ID:             uuid.New(),
			ProjectID:      projectID,
			ResourceTypeID: resourceTypeID,
			Name:           res.Name,
			PosX:           posX,
			PosY:           posY,
			IsVisualOnly:   isVisualOnly,
			Config:         datatypes.JSON(configJSON),
		}
		if err := s.resourceRepo.Create(txCtx, dbResource); err != nil {
			s.projectRepo.RollbackTransaction(tx)
			return nil, fmt.Errorf("create resource %s: %w", res.Name, err)
		}

		domainIDToDBID[res.ID] = dbResource.ID
	}

	// Create containment relationships
	for parentID, childIDs := range arch.Containments {
		parentDBID, ok := domainIDToDBID[parentID]
		if !ok {
			continue
		}
		for _, childID := range childIDs {
			childDBID, ok := domainIDToDBID[childID]
			if !ok {
				continue
			}
			containment := &models.ResourceContainment{
				ParentResourceID: parentDBID,
				ChildResourceID:  childDBID,
			}
			if err := s.containmentRepo.Create(txCtx, containment); err != nil {
				s.projectRepo.RollbackTransaction(tx)
				return nil, fmt.Errorf("create containment: %w", err)
			}
		}
	}

	// Create dependency relationships
	dependencyType, err := s.dependencyTypeRepo.FindByName(txCtx, "depends_on")
	if err != nil {
		s.projectRepo.RollbackTransaction(tx)
		return nil, fmt.Errorf("dependency type 'depends_on' not found: %w", err)
	}

	for fromID, toIDs := range arch.Dependencies {
		fromDBID, ok := domainIDToDBID[fromID]
		if !ok {
			continue
		}
		for _, toID := range toIDs {
			toDBID, ok := domainIDToDBID[toID]
			if !ok {
				continue
			}
			dependency := &models.ResourceDependency{
				FromResourceID:   fromDBID,
				ToResourceID:     toDBID,
				DependencyTypeID: dependencyType.ID,
			}
			if err := s.dependencyRepo.Create(txCtx, dependency); err != nil {
				s.projectRepo.RollbackTransaction(tx)
				return nil, fmt.Errorf("create dependency: %w", err)
			}
		}
	}

	// Commit transaction before pricing calculation
	if err := s.projectRepo.CommitTransaction(tx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	result := &serverinterfaces.ArchitecturePersistResult{
		ResourceIDMapping: domainIDToDBID,
	}

	// Calculate and persist pricing if pricing service is available
	if s.pricingService != nil && pricingDuration > 0 {
		fmt.Println("\n💵 Calculating pricing for resources...")
		fmt.Println(strings.Repeat("-", 100))

		// Calculate architecture cost
		archEstimate, err := s.pricingService.CalculateArchitectureCost(ctx, arch, pricingDuration)
		if err != nil {
			// Log error but don't fail the operation
			// Pricing calculation failure shouldn't block architecture persistence
			fmt.Printf("  ⚠️  Failed to calculate architecture cost: %v\n", err)
		} else {
			result.PricingEstimate = archEstimate

			// Log resource costs summary
			fmt.Printf("  ✓ Calculated pricing for %d resources\n", len(archEstimate.ResourceEstimates))
			for domainID, resEstimate := range archEstimate.ResourceEstimates {
				// Find the domain resource to get its details
				var domainRes *resource.Resource
				for _, res := range arch.Resources {
					if res.ID == domainID {
						domainRes = res
						break
					}
				}

				resourceName := resEstimate.ResourceName
				if domainRes != nil {
					resourceName = domainRes.Name
				}

				// Calculate base and hidden costs for display
				baseCost := 0.0
				hiddenCost := 0.0
				for _, comp := range resEstimate.Breakdown {
					if strings.Contains(comp.ComponentName, "(") && strings.Contains(comp.ComponentName, ")") {
						hiddenCost += comp.Subtotal
					} else {
						baseCost += comp.Subtotal
					}
				}

				fmt.Printf("    • %s (%s): $%.2f", resourceName, resEstimate.ResourceType, resEstimate.TotalCost)
				if hiddenCost > 0 {
					fmt.Printf(" [Base: $%.2f + Hidden: $%.2f]", baseCost, hiddenCost)
				}
				fmt.Println()
			}
			fmt.Println(strings.Repeat("-", 100))
			fmt.Printf("  💰 Total Architecture Cost: $%.2f %s\n", archEstimate.TotalCost, archEstimate.Currency)
			fmt.Println()

			// Persist individual resource pricing
			for domainID, dbID := range domainIDToDBID {
				if resEstimate, ok := archEstimate.ResourceEstimates[domainID]; ok {
					// Find the domain resource to get its details
					var domainRes *resource.Resource
					for _, res := range arch.Resources {
						if res.ID == domainID {
							domainRes = res
							break
						}
					}

					if domainRes != nil {
						// Calculate individual resource cost
						costEstimate, err := s.pricingService.CalculateResourceCost(ctx, domainRes, pricingDuration)
						if err == nil {
							// Persist resource pricing
							_ = s.pricingService.PersistResourcePricing(ctx, projectID, dbID, costEstimate, string(arch.Provider), arch.Region)
						}
					}
					_ = resEstimate // suppress unused warning
				}
			}

			// Persist project-level pricing
			projectEstimate := &domainpricing.CostEstimate{
				TotalCost:    archEstimate.TotalCost,
				Currency:     domainpricing.Currency(archEstimate.Currency),
				Period:       domainpricing.Period(archEstimate.Period),
				Duration:     archEstimate.Duration,
				CalculatedAt: time.Now(),
				Provider:     domainpricing.CloudProvider(archEstimate.Provider),
				Region:       &archEstimate.Region,
			}
			_ = s.pricingService.PersistProjectPricing(ctx, projectID, projectEstimate, archEstimate.Provider, archEstimate.Region)
		}
	}

	return result, nil
}

// GetProjectPricing retrieves pricing for a project
func (s *ProjectServiceImpl) GetProjectPricing(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectPricing, error) {
	if s.pricingService == nil {
		return nil, fmt.Errorf("pricing service not configured")
	}
	return s.pricingService.GetProjectPricing(ctx, projectID)
}

// LoadArchitecture loads an architecture from the database for a project
func (s *ProjectServiceImpl) LoadArchitecture(ctx context.Context, projectID uuid.UUID) (*architecture.Architecture, error) {
	// Step 1: Load project
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load project: %w", err)
	}

	// Step 2: Load all resources for the project
	dbResources, err := s.resourceRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load resources: %w", err)
	}

	// Step 3: Convert database resources to domain resources
	domainResources := make([]*resource.Resource, 0, len(dbResources))
	dbIDToDomainID := make(map[uuid.UUID]string) // DB UUID -> Domain ID

	provider := resource.CloudProvider(project.CloudProvider)
	if provider == "" {
		provider = resource.AWS // Default
	}

	for _, dbRes := range dbResources {
		// Use database UUID as domain resource ID
		domainID := dbRes.ID.String()
		dbIDToDomainID[dbRes.ID] = domainID

		// Parse config/metadata from JSON
		var metadata map[string]interface{}
		if err := json.Unmarshal(dbRes.Config, &metadata); err != nil {
			// If unmarshal fails, create empty metadata
			metadata = make(map[string]interface{})
		}

		// Add position to metadata
		metadata["position"] = map[string]interface{}{
			"x": dbRes.PosX,
			"y": dbRes.PosY,
		}
		metadata["isVisualOnly"] = dbRes.IsVisualOnly

		// Map resource type
		resourceType := resource.ResourceType{
			ID:         dbRes.ResourceType.Name,
			Name:       dbRes.ResourceType.Name,
			Category:   "",
			Kind:       "",
			IsRegional: dbRes.ResourceType.IsRegional,
			IsGlobal:   dbRes.ResourceType.IsGlobal,
		}

		// Get category and kind if available
		if dbRes.ResourceType.Category != nil {
			resourceType.Category = dbRes.ResourceType.Category.Name
		}
		if dbRes.ResourceType.Kind != nil {
			resourceType.Kind = dbRes.ResourceType.Kind.Name
		}

		domainRes := &resource.Resource{
			ID:        domainID,
			Name:      dbRes.Name,
			Type:      resourceType,
			Provider:  provider,
			Region:    project.Region,
			ParentID:  nil,        // Will be set from containments
			DependsOn: []string{}, // Will be set from dependencies
			Metadata:  metadata,
		}

		domainResources = append(domainResources, domainRes)
	}

	// Step 4: Load containments and build parent-child relationships
	containments := make(map[string][]string) // parentID -> []childIDs
	childToParent := make(map[string]string)  // childID -> parentID

	for _, dbRes := range dbResources {
		childID := dbRes.ID.String()
		parentContainments, err := s.containmentRepo.FindParents(ctx, dbRes.ID)
		if err != nil {
			// Log error but continue
			continue
		}

		for _, containment := range parentContainments {
			// Use ParentResourceID directly (it's always set)
			parentID := containment.ParentResourceID.String()
			childToParent[childID] = parentID

			// Add to containments map
			if _, exists := containments[parentID]; !exists {
				containments[parentID] = make([]string, 0)
			}
			containments[parentID] = append(containments[parentID], childID)

			// Set parent ID on domain resource
			for _, domainRes := range domainResources {
				if domainRes.ID == childID {
					domainRes.ParentID = &parentID
					break
				}
			}
		}
	}

	// Step 5: Load dependencies
	dependencies := make(map[string][]string) // resourceID -> []dependencyIDs

	for _, dbRes := range dbResources {
		fromID := dbRes.ID.String()
		dbDependencies, err := s.dependencyRepo.FindByFromResource(ctx, dbRes.ID)
		if err != nil {
			// Log error but continue
			continue
		}

		depIDs := make([]string, 0)
		for _, dep := range dbDependencies {
			// Use ToResourceID directly (it's always set)
			toID := dep.ToResourceID.String()
			depIDs = append(depIDs, toID)
		}

		if len(depIDs) > 0 {
			dependencies[fromID] = depIDs

			// Set DependsOn on domain resource
			for _, domainRes := range domainResources {
				if domainRes.ID == fromID {
					domainRes.DependsOn = depIDs
					break
				}
			}
		}
	}

	// Step 6: Build architecture aggregate
	arch := &architecture.Architecture{
		Resources:    domainResources,
		Region:       project.Region,
		Provider:     provider,
		Containments: containments,
		Dependencies: dependencies,
		Variables:    []architecture.Variable{}, // Variables not stored in DB yet
		Outputs:      []architecture.Output{},   // Outputs not stored in DB yet
	}

	return arch, nil
}
