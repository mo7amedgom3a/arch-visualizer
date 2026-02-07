package diagram

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"log/slog"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/parser"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"gorm.io/datatypes"
)

// Service provides the main diagram processing service
type Service struct {
	projectRepo        *repository.ProjectRepository
	resourceRepo       *repository.ResourceRepository
	resourceTypeRepo   *repository.ResourceTypeRepository
	dependencyTypeRepo *repository.DependencyTypeRepository
	logger             *slog.Logger
}

// NewService creates a new diagram service
func NewService(logger *slog.Logger) (*Service, error) {
	projectRepo, err := repository.NewProjectRepository(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create project repository: %w", err)
	}

	resourceRepo, err := repository.NewResourceRepository(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource repository: %w", err)
	}

	resourceTypeRepo, err := repository.NewResourceTypeRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create resource type repository: %w", err)
	}

	dependencyTypeRepo, err := repository.NewDependencyTypeRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create dependency type repository: %w", err)
	}

	return &Service{
		projectRepo:        projectRepo,
		resourceRepo:       resourceRepo,
		resourceTypeRepo:   resourceTypeRepo,
		dependencyTypeRepo: dependencyTypeRepo,
		logger:             logger,
	}, nil
}

// ProcessDiagramRequest processes a diagram JSON request and persists it as a project
// Returns the created project ID
func (s *Service) ProcessDiagramRequest(ctx context.Context, jsonData []byte, userID uuid.UUID, projectName string, iacToolID uint) (uuid.UUID, error) {
	// Step 1: Parse IR diagram
	irDiagram, err := parser.ParseIRDiagram(jsonData)
	fmt.Println("IR Digram", irDiagram)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse diagram: %w", err)
	}

	// Step 2: Normalize to graph
	diagramGraph, err := parser.NormalizeToGraph(irDiagram)
	fmt.Println("Dagram graph ", diagramGraph)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to normalize diagram: %w", err)
	}

	// Step 3: Extract region and provider from diagram (needed for validation)
	regionNode, hasRegion := diagramGraph.FindRegionNode()
	region := "us-east-1"    // default
	provider := resource.AWS // default, could be extracted from config

	if hasRegion {
		if regionName, ok := extractRegionFromConfig(regionNode.Config); ok {
			region = regionName
		}
		// Could extract provider from region node config if available
	}

	// Step 4: Load valid resource types from database and validate graph
	validResourceTypes, err := s.buildValidResourceTypesMap(ctx, provider)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to load valid resource types: %w", err)
	}

	validationOpts := &validator.ValidationOptions{
		ValidResourceTypes: validResourceTypes,
		Provider:           string(provider),
	}
	validationResult := validator.Validate(diagramGraph, validationOpts)
	if !validationResult.Valid {
		return uuid.Nil, fmt.Errorf("diagram validation failed: %v", validationResult.Errors)
	}

	// Step 5: Map to domain architecture
	domainArch, err := architecture.MapDiagramToArchitecture(diagramGraph, provider)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to map diagram to architecture: %w", err)
	}

	// Step 6: Persist project and resources
	projectID, err := s.persistArchitecture(ctx, domainArch, userID, projectName, iacToolID, provider, region)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to persist architecture: %w", err)
	}

	return projectID, nil
}

// persistArchitecture persists the architecture as a project with resources
func (s *Service) persistArchitecture(
	ctx context.Context,
	arch *architecture.Architecture,
	userID uuid.UUID,
	projectName string,
	iacToolID uint,
	provider resource.CloudProvider,
	region string,
) (uuid.UUID, error) {
	// Start transaction
	baseRepo := s.projectRepo.BaseRepository
	tx, txCtx := baseRepo.BeginTransaction(ctx)
	defer func() {
		if r := recover(); r != nil {
			baseRepo.RollbackTransaction(tx)
			panic(r)
		}
	}()

	// Create project
	project := &models.Project{
		UserID:        userID,
		InfraToolID:   iacToolID,
		Name:          projectName,
		CloudProvider: string(provider),
		Region:        region,
	}

	if err := s.projectRepo.Create(txCtx, project); err != nil {
		baseRepo.RollbackTransaction(tx)
		return uuid.Nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Create resources
	irNodeIDToDBResourceID := make(map[string]uuid.UUID) // IR node ID -> DB resource UUID

	for _, domainRes := range arch.Resources {
		// Lookup resource type
		resourceType, err := s.resourceTypeRepo.FindByNameAndProvider(txCtx, domainRes.Type.Name, string(provider))
		if err != nil {
			baseRepo.RollbackTransaction(tx)
			return uuid.Nil, fmt.Errorf("failed to find resource type %s for provider %s: %w", domainRes.Type.Name, provider, err)
		}

		// Extract UI State
		var uiState *models.ResourceUIState
		if ui, ok := domainRes.Metadata["ui"].(*graph.UIState); ok && ui != nil {
			styleJSON, _ := json.Marshal(ui.Style)
			measuredJSON, _ := json.Marshal(ui.Measured)

			uiState = &models.ResourceUIState{
				X:          ui.Position.X,
				Y:          ui.Position.Y,
				Width:      ui.Width,
				Height:     ui.Height,
				Style:      datatypes.JSON(styleJSON),
				Measured:   datatypes.JSON(measuredJSON),
				Selected:   ui.Selected,
				Dragging:   ui.Dragging,
				Resizing:   ui.Resizing,
				Focusable:  ui.Focusable,
				Selectable: ui.Selectable,
				ZIndex:     ui.ZIndex,
			}
		}

		// Extract isVisualOnly flag from metadata
		isVisualOnly := false
		if val, ok := domainRes.Metadata["isVisualOnly"].(bool); ok {
			isVisualOnly = val
		}

		// Convert config to JSON (includes isVisualOnly flag in metadata)
		configJSON, err := json.Marshal(domainRes.Metadata)
		if err != nil {
			baseRepo.RollbackTransaction(tx)
			return uuid.Nil, fmt.Errorf("failed to marshal config: %w", err)
		}

		// Create database resource
		dbResource := &models.Resource{
			ProjectID:      project.ID,
			ResourceTypeID: resourceType.ID,
			Name:           domainRes.Name,
			UIState:        uiState,
			IsVisualOnly:   isVisualOnly,
			Config:         datatypes.JSON(configJSON),
		}

		if err := s.resourceRepo.Create(txCtx, dbResource); err != nil {
			baseRepo.RollbackTransaction(tx)
			return uuid.Nil, fmt.Errorf("failed to create resource %s: %w", domainRes.Name, err)
		}

		irNodeIDToDBResourceID[domainRes.ID] = dbResource.ID
	}

	// Create containment relationships
	for parentIRID, childIRIDs := range arch.Containments {
		parentDBID, parentExists := irNodeIDToDBResourceID[parentIRID]
		if !parentExists {
			continue // Skip if parent not found (might be region)
		}

		for _, childIRID := range childIRIDs {
			childDBID, childExists := irNodeIDToDBResourceID[childIRID]
			if !childExists {
				continue
			}

			if err := s.resourceRepo.CreateContainment(txCtx, parentDBID, childDBID); err != nil {
				baseRepo.RollbackTransaction(tx)
				return uuid.Nil, fmt.Errorf("failed to create containment: %w", err)
			}
		}
	}

	// Create dependency relationships
	for resourceIRID, depIRIDs := range arch.Dependencies {
		resourceDBID, resourceExists := irNodeIDToDBResourceID[resourceIRID]
		if !resourceExists {
			continue
		}

		// Get or create "depends_on" dependency type
		depType, err := s.dependencyTypeRepo.FindByName(txCtx, "depends_on")
		if err != nil {
			// Try to create if it doesn't exist (should be seeded, but handle gracefully)
			baseRepo.RollbackTransaction(tx)
			return uuid.Nil, fmt.Errorf("dependency type 'depends_on' not found: %w", err)
		}

		for _, depIRID := range depIRIDs {
			depDBID, depExists := irNodeIDToDBResourceID[depIRID]
			if !depExists {
				continue
			}

			dependency := &models.ResourceDependency{
				FromResourceID:   resourceDBID,
				ToResourceID:     depDBID,
				DependencyTypeID: depType.ID,
			}

			if err := s.resourceRepo.CreateDependency(txCtx, dependency); err != nil {
				baseRepo.RollbackTransaction(tx)
				return uuid.Nil, fmt.Errorf("failed to create dependency: %w", err)
			}
		}
	}

	// Commit transaction
	if err := baseRepo.CommitTransaction(tx); err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return project.ID, nil
}

// buildValidResourceTypesMap loads resource types from database and converts to IR format
// Returns a map of IR resource type names (lowercase) -> true
func (s *Service) buildValidResourceTypesMap(ctx context.Context, provider resource.CloudProvider) (map[string]bool, error) {
	resourceTypes, err := s.resourceTypeRepo.ListByProvider(ctx, string(provider))
	if err != nil {
		return nil, fmt.Errorf("failed to load resource types: %w", err)
	}

	validTypes := make(map[string]bool)
	// Region is always valid (it's not a real resource type in DB)
	validTypes["region"] = true

	// Convert database resource type names to IR format
	for _, rt := range resourceTypes {
		irType := convertDBNameToIRType(rt.Name)
		validTypes[irType] = true
	}

	return validTypes, nil
}

// convertDBNameToIRType converts database resource type names to IR format
// Examples:
//   - "RouteTable" -> "route-table"
//   - "SecurityGroup" -> "security-group"
//   - "VPC" -> "vpc"
//   - "EC2" -> "ec2"
//   - "AutoScalingGroup" -> "auto-scaling-group"
//   - "DynamoDB" -> "dynamodb" (special case: consecutive capitals stay together)
func convertDBNameToIRType(dbName string) string {
	if dbName == "" {
		return ""
	}

	// Special case mappings for known exceptions
	specialCases := map[string]string{
		"DynamoDB":         "dynamodb",
		"EC2":              "ec2",
		"S3":               "s3",
		"EBS":              "ebs",
		"RDS":              "rds",
		"VPC":              "vpc",
		"AutoScalingGroup": "autoscaling-group",
	}
	if mapped, ok := specialCases[dbName]; ok {
		return mapped
	}

	var result strings.Builder
	runes := []rune(dbName)

	for i, r := range runes {
		// Insert hyphen before capital letters (except the first character)
		if i > 0 && unicode.IsUpper(r) {
			// Check if previous character was also uppercase (like "EC2" -> "ec2", not "e-c-2")
			if unicode.IsUpper(runes[i-1]) {
				// If next character is lowercase, insert hyphen (like "RouteTable" -> "route-table")
				if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					result.WriteRune('-')
				}
			} else {
				result.WriteRune('-')
			}
		}
		result.WriteRune(unicode.ToLower(r))
	}

	return result.String()
}

// extractRegionFromConfig extracts region name from config
func extractRegionFromConfig(config map[string]interface{}) (string, bool) {
	if name, ok := config["name"].(string); ok {
		return name, true
	}
	return "", false
}
