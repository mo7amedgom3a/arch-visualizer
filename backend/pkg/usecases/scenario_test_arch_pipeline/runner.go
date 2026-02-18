package scenario_test_arch_pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/parser"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/services"
	infrastructurerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/infrastructure"
	projectrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/project"
	resourcerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/resource"
	userrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/user"
)

// Run executes the architecture pipeline test scenario
func Run(ctx context.Context) error {
	fmt.Println("==================================================================================")
	fmt.Println("SCENARIO: Architecture Pipeline Test with Graph Export")
	fmt.Println("==================================================================================")

	// 1. Read input JSON
	jsonPath, err := resolveDiagramJSONPath("json_request_with_ui.json")
	if err != nil {
		return err
	}
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("read IR json %s: %w", jsonPath, err)
	}
	fmt.Printf("✓ Read diagram JSON from: %s\n", jsonPath)

	// Extract diagram if wrapped in project structure
	diagramData, err := extractDiagramFromProjectJSON(data)
	if err != nil {
		return fmt.Errorf("extract diagram from project JSON: %w", err)
	}

	// 2. Parse IR
	irDiagram, err := parser.ParseIRDiagram(diagramData)
	if err != nil {
		return fmt.Errorf("parse IR diagram: %w", err)
	}
	fmt.Printf("✓ Parsed IR diagram: %d nodes\n", len(irDiagram.Nodes))

	// 3. Normalize to Graph
	diagramGraph, err := parser.NormalizeToGraph(irDiagram)
	if err != nil {
		return fmt.Errorf("normalize diagram: %w", err)
	}
	fmt.Printf("✓ Normalized to graph: %d nodes, %d edges\n", len(diagramGraph.Nodes), len(diagramGraph.Edges))

	// 4. Validate
	validationResult := validator.Validate(diagramGraph, nil)
	if !validationResult.Valid {
		return fmt.Errorf("diagram validation failed: %v", validationResult.Errors)
	}
	fmt.Println("✓ Diagram validation passed")

	// 5. Map to Architecture
	arch, err := architecture.MapDiagramToArchitecture(diagramGraph, resource.AWS)
	if err != nil {
		return fmt.Errorf("map diagram to architecture: %w", err)
	}
	fmt.Printf("✓ Mapped to architecture: %d resources\n", len(arch.Resources))

	// 6. Persist to Database
	projectID, err := persistProject(ctx, arch, diagramGraph)
	if err != nil {
		return fmt.Errorf("persist project: %w", err)
	}
	fmt.Printf("✓ Persisted project to DB with ID: %s\n", projectID)

	// 7. Export Graph to JSON
	if err := exportGraphJSON(diagramGraph, "pipeline_output.json"); err != nil {
		return fmt.Errorf("export graph json: %w", err)
	}
	fmt.Println("✓ Exported graph to pipeline_output.json")

	return nil
}

func persistProject(ctx context.Context, arch *architecture.Architecture, diagramGraph interface{}) (uuid.UUID, error) {
	// Initialize repositories with adapters
	logger := slog.Default()

	// Base repositories
	projectBase, _ := projectrepo.NewProjectRepository(logger)
	resourceBase, _ := resourcerepo.NewResourceRepository(logger)
	verBase, _ := projectrepo.NewProjectVersionRepository()
	resTypeBase, _ := resourcerepo.NewResourceTypeRepository()
	contBase, _ := resourcerepo.NewResourceContainmentRepository()
	depBase, _ := resourcerepo.NewResourceDependencyRepository()
	depTypeBase, _ := resourcerepo.NewDependencyTypeRepository()
	userBase, _ := userrepo.NewUserRepository()
	iacBase, _ := infrastructurerepo.NewIACTargetRepository()
	varBase, _ := projectrepo.NewProjectVariableRepository()
	outBase, _ := projectrepo.NewProjectOutputRepository()

	// Adapters
	projectRepo := &services.ProjectRepositoryAdapter{Repo: projectBase}
	resourceRepo := &services.ResourceRepositoryAdapter{Repo: resourceBase}
	verRepo := &services.ProjectVersionRepositoryAdapter{Repo: verBase}
	resTypeRepo := &services.ResourceTypeRepositoryAdapter{Repo: resTypeBase}
	contRepo := &services.ResourceContainmentRepositoryAdapter{Repo: contBase}
	depRepo := &services.ResourceDependencyRepositoryAdapter{Repo: depBase}
	depTypeRepo := &services.DependencyTypeRepositoryAdapter{Repo: depTypeBase}
	userRepo := &services.UserRepositoryAdapter{Repo: userBase}
	iacRepo := &services.IACTargetRepositoryAdapter{Repo: iacBase}
	varRepo := &services.ProjectVariableRepositoryAdapter{Repo: varBase}
	outRepo := &services.ProjectOutputRepositoryAdapter{Repo: outBase}

	// Initialize Service
	projectService := services.NewProjectService(
		projectRepo, verRepo, resourceRepo, resTypeRepo,
		contRepo, depRepo, depTypeRepo, userRepo, iacRepo, varRepo, outRepo,
	)

	// Ensure user exists
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	_, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		user := &models.User{
			ID:         userID,
			Name:       "Test User",
			Email:      "test@example.com",
			IsVerified: true,
		}
		if err := userRepo.Create(ctx, user); err != nil {
			return uuid.Nil, fmt.Errorf("create user: %w", err)
		}
	}

	// Ensure IACTarget (InfraTool) exists
	iacTargetName := "Terraform"
	iacTarget, err := iacRepo.FindByName(ctx, iacTargetName)
	if err != nil {
		// Create if not exists
		iacTarget = &models.IACTarget{
			Name: iacTargetName,
		}
		if err := iacRepo.Create(ctx, iacTarget); err != nil {
			return uuid.Nil, fmt.Errorf("create iac target: %w", err)
		}
		// Refetch to get ID
		iacTarget, err = iacRepo.FindByName(ctx, iacTargetName)
		if err != nil {
			return uuid.Nil, fmt.Errorf("find created iac target: %w", err)
		}
	}

	// Create Project Request
	req := &serverinterfaces.CreateProjectRequest{
		UserID:        userID,
		Name:          "Pipeline Test Project",
		Description:   "Test project for architecture pipeline with UI State",
		CloudProvider: string(arch.Provider),
		Region:        arch.Region,
		IACTargetID:   iacTarget.ID,
	}

	// Create Project
	project, err := projectService.Create(ctx, req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create project service: %w", err)
	}

	// Persist Architecture
	if err := projectService.PersistArchitecture(ctx, project.ID, arch, diagramGraph); err != nil {
		return uuid.Nil, fmt.Errorf("persist architecture service: %w", err)
	}

	return project.ID, nil
}

func exportGraphJSON(graph interface{}, filename string) error {
	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Helpers (copied/adapted from scenario6)

func resolveDiagramJSONPath(filename string) (string, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine caller")
	}
	dir := filepath.Dir(thisFile)
	// Up 3 levels to backend root
	root := filepath.Clean(filepath.Join(dir, "..", "..", ".."))
	return filepath.Join(root, filename), nil
}

func extractDiagramFromProjectJSON(data []byte) ([]byte, error) {
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	if _, hasNodes := rawData["nodes"]; hasNodes {
		return data, nil
	}
	for _, value := range rawData {
		if projectData, ok := value.(map[string]interface{}); ok {
			if _, hasNodes := projectData["nodes"]; hasNodes {
				diagramBytes, err := json.Marshal(projectData)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal extracted diagram: %w", err)
				}
				return diagramBytes, nil
			}
		}
	}
	// Also check if root has uiState and nodes, if so return data as is
	if _, hasUI := rawData["uiState"]; hasUI {
		if _, hasNodes := rawData["nodes"]; hasNodes {
			return data, nil
		}
	}

	return nil, fmt.Errorf("could not find diagram structure")
}
