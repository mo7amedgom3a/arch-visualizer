package scenario6_terraform_with_persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/terraform"
	awsrules "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/rules"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/parser"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	rulesengine "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/engine"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
	tfgen "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/generator"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"gorm.io/datatypes"
)

// TerraformWithPersistenceRunner demonstrates the end-to-end pipeline from a diagram IR JSON
// to generated Terraform files, while also persisting the project and architecture to the database.
//
// Steps:
//  1. Parse IR JSON into diagram graph
//  2. Validate diagram (structure, schemas, relationships)
//  3. Map to domain Architecture aggregate
//  4. Validate domain rules/constraints (AWS networking defaults)
//  5. Build domain graph + topologically sort resources
//  6. Run Terraform engine to produce IaC files
//  7. Persist project, resources, containments, and dependencies to database
//  8. Write Terraform files to ./terraform_output/
func TerraformWithPersistenceRunner(ctx context.Context) error {
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("SCENARIO 6: Terraform Code Generation with Database Persistence")
	fmt.Println(strings.Repeat("=", 100))

	// 1) Read IR JSON from file (using complete diagram for full feature testing)
	jsonPath, err := resolveDiagramJSONPath("json-request-fiagram-complete.json")
	if err != nil {
		return err
	}
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("read IR json %s: %w", jsonPath, err)
	}
	fmt.Printf("✓ Read diagram JSON from: %s\n", jsonPath)

	// Extract diagram from project-wrapped JSON structure if needed
	diagramData, err := extractDiagramFromProjectJSON(data)
	if err != nil {
		return fmt.Errorf("extract diagram from project JSON: %w", err)
	}

	// 2) Parse & normalize to diagram graph
	irDiagram, err := parser.ParseIRDiagram(diagramData)
	if err != nil {
		return fmt.Errorf("parse IR diagram: %w", err)
	}
	fmt.Printf("✓ Parsed IR diagram: %d nodes, %d variables, %d outputs\n",
		len(irDiagram.Nodes), len(irDiagram.Variables), len(irDiagram.Outputs))

	diagramGraph, err := parser.NormalizeToGraph(irDiagram)
	if err != nil {
		return fmt.Errorf("normalize diagram: %w", err)
	}
	fmt.Printf("✓ Normalized to graph: %d nodes, %d edges\n", len(diagramGraph.Nodes), len(diagramGraph.Edges))

	// 3) Validate diagram (structure + schemas + relationships)
	validationResult := validator.Validate(diagramGraph, nil)
	if !validationResult.Valid {
		return fmt.Errorf("diagram validation failed:\n%s", formatValidationErrors(validationResult))
	}
	fmt.Println("✓ Diagram validation passed")

	// 4) Map to domain architecture (cloud-agnostic core)
	arch, err := architecture.MapDiagramToArchitecture(diagramGraph, resource.AWS)
	if err != nil {
		return fmt.Errorf("map diagram to architecture: %w", err)
	}
	fmt.Printf("✓ Mapped to architecture: %d resources, region=%s, provider=%s\n",
		len(arch.Resources), arch.Region, arch.Provider)

	// Basic domain validation hook
	if err := arch.Validate(); err != nil {
		return fmt.Errorf("architecture validation failed: %w", err)
	}

	// 5) Validate domain rules/constraints using AWS default networking rules
	if err := validateAWSRules(ctx, arch); err != nil {
		return fmt.Errorf("AWS rules validation: %w", err)
	}
	fmt.Println("✓ AWS rules validation passed")

	// 6) Build domain graph + topologically sort resources
	graph := architecture.NewGraph(arch)
	sorted, err := graph.GetSortedResources()
	if err != nil {
		return fmt.Errorf("topological sort failed: %w", err)
	}
	fmt.Printf("✓ Topologically sorted: %d resources\n", len(sorted))

	// 7) Generate Terraform using Terraform engine (HCL generation)
	output, err := generateTerraform(ctx, arch, sorted)
	if err != nil {
		return fmt.Errorf("terraform generation: %w", err)
	}
	fmt.Printf("✓ Generated Terraform: %d files\n", len(output.Files))

	// 8) Persist to database
	projectID, err := persistToDatabase(ctx, arch, diagramGraph)
	if err != nil {
		return fmt.Errorf("database persistence: %w", err)
	}
	fmt.Printf("✓ Persisted to database: project_id=%s\n", projectID.String())

	// 9) Write generated files to disk
	outDir := "terraform_output"
	if err := writeTerraformOutput(outDir, output); err != nil {
		return err
	}

	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("SUCCESS: Terraform code generated and project persisted!")
	fmt.Printf("  Project ID: %s\n", projectID.String())
	fmt.Printf("  Output directory: ./%s/\n", outDir)
	for _, f := range output.Files {
		fmt.Printf("    - %s\n", f.Path)
	}
	fmt.Println(strings.Repeat("=", 100))

	return nil
}

// persistToDatabase saves the project, resources, containments, and dependencies to the database
func persistToDatabase(ctx context.Context, arch *architecture.Architecture, diagramGraph *graph.DiagramGraph) (uuid.UUID, error) {
	// Initialize repositories
	projectRepo, err := repository.NewProjectRepository(slog.Default())
	if err != nil {
		return uuid.Nil, fmt.Errorf("create project repository: %w", err)
	}

	resourceRepo, err := repository.NewResourceRepository(slog.Default())
	if err != nil {
		return uuid.Nil, fmt.Errorf("create resource repository: %w", err)
	}

	resourceTypeRepo, err := repository.NewResourceTypeRepository()
	if err != nil {
		return uuid.Nil, fmt.Errorf("create resource type repository: %w", err)
	}

	iacTargetRepo, err := repository.NewIACTargetRepository()
	if err != nil {
		return uuid.Nil, fmt.Errorf("create iac target repository: %w", err)
	}

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		return uuid.Nil, fmt.Errorf("create user repository: %w", err)
	}

	containmentRepo, err := repository.NewResourceContainmentRepository()
	if err != nil {
		return uuid.Nil, fmt.Errorf("create containment repository: %w", err)
	}

	dependencyRepo, err := repository.NewResourceDependencyRepository()
	if err != nil {
		return uuid.Nil, fmt.Errorf("create dependency repository: %w", err)
	}

	dependencyTypeRepo, err := repository.NewDependencyTypeRepository()
	if err != nil {
		return uuid.Nil, fmt.Errorf("create dependency type repository: %w", err)
	}

	// Start transaction
	tx, txCtx := projectRepo.BeginTransaction(ctx)
	defer func() {
		if r := recover(); r != nil {
			projectRepo.RollbackTransaction(tx)
		}
	}()

	// 1. Ensure user exists (create demo user if not exists)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	user, err := userRepo.FindByID(txCtx, userID)
	if err != nil {
		// Create demo user
		user = &models.User{
			ID:         userID,
			Name:       "Demo User",
			IsVerified: true,
		}
		if err := userRepo.Create(txCtx, user); err != nil {
			projectRepo.RollbackTransaction(tx)
			return uuid.Nil, fmt.Errorf("create demo user: %w", err)
		}
		fmt.Println("  → Created demo user")
	}

	// 2. Get or create IaC target (Terraform)
	iacTarget, err := iacTargetRepo.FindByName(txCtx, "Terraform")
	if err != nil {
		// Create Terraform target
		iacTarget = &models.IACTarget{
			Name: "Terraform",
		}
		if err := iacTargetRepo.Create(txCtx, iacTarget); err != nil {
			projectRepo.RollbackTransaction(tx)
			return uuid.Nil, fmt.Errorf("create iac target: %w", err)
		}
		fmt.Println("  → Created IaC target: Terraform")
	}

	// 3. Create project
	project := &models.Project{
		ID:            uuid.New(),
		UserID:        userID,
		InfraToolID:   iacTarget.ID,
		Name:          "Generated Architecture Project",
		CloudProvider: string(arch.Provider),
		Region:        arch.Region,
	}
	if err := projectRepo.Create(txCtx, project); err != nil {
		projectRepo.RollbackTransaction(tx)
		return uuid.Nil, fmt.Errorf("create project: %w", err)
	}
	fmt.Printf("  → Created project: %s\n", project.Name)

	// 4. Create resources and build ID mapping
	domainIDToDBID := make(map[string]uuid.UUID)
	resourceTypeCache := make(map[string]uint) // cache resource type IDs

	for _, res := range arch.Resources {
		// Get or cache resource type ID
		resourceTypeID, ok := resourceTypeCache[res.Type.Name]
		if !ok {
			resourceType, err := resourceTypeRepo.FindByNameAndProvider(txCtx, res.Type.Name, string(arch.Provider))
			if err != nil {
				// Create resource type if not exists
				newResourceType := &models.ResourceType{
					Name:          res.Type.Name,
					CloudProvider: string(arch.Provider),
					IsRegional:    true,
					IsGlobal:      false,
				}
				if err := tx.Create(newResourceType).Error; err != nil {
					projectRepo.RollbackTransaction(tx)
					return uuid.Nil, fmt.Errorf("create resource type %s: %w", res.Type.Name, err)
				}
				resourceTypeID = newResourceType.ID
				fmt.Printf("  → Created resource type: %s\n", res.Type.Name)
			} else {
				resourceTypeID = resourceType.ID
			}
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
			projectRepo.RollbackTransaction(tx)
			return uuid.Nil, fmt.Errorf("marshal resource config: %w", err)
		}

		// Create resource
		dbResource := &models.Resource{
			ID:             uuid.New(),
			ProjectID:      project.ID,
			ResourceTypeID: resourceTypeID,
			Name:           res.Name,
			PosX:           posX,
			PosY:           posY,
			IsVisualOnly:   isVisualOnly,
			Config:         datatypes.JSON(configJSON),
		}
		if err := resourceRepo.Create(txCtx, dbResource); err != nil {
			projectRepo.RollbackTransaction(tx)
			return uuid.Nil, fmt.Errorf("create resource %s: %w", res.Name, err)
		}

		domainIDToDBID[res.ID] = dbResource.ID
	}
	fmt.Printf("  → Created %d resources\n", len(arch.Resources))

	// 5. Create containment relationships
	containmentCount := 0
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
			if err := containmentRepo.Create(txCtx, containment); err != nil {
				projectRepo.RollbackTransaction(tx)
				return uuid.Nil, fmt.Errorf("create containment: %w", err)
			}
			containmentCount++
		}
	}
	fmt.Printf("  → Created %d containment relationships\n", containmentCount)

	// 6. Create dependency relationships
	// Get or create "depends_on" dependency type
	dependencyType, err := dependencyTypeRepo.FindByName(txCtx, "depends_on")
	if err != nil {
		// Create dependency type
		newDepType := &models.DependencyType{
			Name: "depends_on",
		}
		if err := tx.Create(newDepType).Error; err != nil {
			projectRepo.RollbackTransaction(tx)
			return uuid.Nil, fmt.Errorf("create dependency type: %w", err)
		}
		dependencyType = newDepType
		fmt.Println("  → Created dependency type: depends_on")
	}

	dependencyCount := 0
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
			if err := dependencyRepo.Create(txCtx, dependency); err != nil {
				projectRepo.RollbackTransaction(tx)
				return uuid.Nil, fmt.Errorf("create dependency: %w", err)
			}
			dependencyCount++
		}
	}
	fmt.Printf("  → Created %d dependency relationships\n", dependencyCount)

	// Commit transaction
	if err := projectRepo.CommitTransaction(tx); err != nil {
		return uuid.Nil, fmt.Errorf("commit transaction: %w", err)
	}

	return project.ID, nil
}

func validateAWSRules(ctx context.Context, arch *architecture.Architecture) error {
	ruleService := awsrules.NewAWSRuleService()

	// Start with code-defined defaults; no DB overrides for this use case.
	if err := ruleService.LoadRulesWithDefaults(ctx, nil); err != nil {
		return fmt.Errorf("load AWS default rules: %w", err)
	}

	// Adapt domain architecture to rules engine architecture view.
	engineArch := &rulesengine.Architecture{
		Resources: arch.Resources,
	}

	results, err := ruleService.ValidateArchitecture(ctx, engineArch)
	if err != nil {
		return fmt.Errorf("validate architecture rules: %w", err)
	}

	var messages []string
	for resID, resResult := range results {
		if !resResult.Valid {
			for _, re := range resResult.Errors {
				messages = append(messages, fmt.Sprintf("resource %s (%s): %s", resID, re.ResourceType, re.Message))
			}
		}
	}

	if len(messages) > 0 {
		return fmt.Errorf("rule/constraint validation failed:\n%s", strings.Join(messages, "\n"))
	}

	return nil
}

func generateTerraform(ctx context.Context, arch *architecture.Architecture, sorted []*resource.Resource) (*iac.Output, error) {
	// Wire Terraform mapper registry with AWS Terraform mapper.
	mapperRegistry := tfmapper.NewRegistry()
	if err := mapperRegistry.Register(terraform.New()); err != nil {
		return nil, fmt.Errorf("register aws terraform mapper: %w", err)
	}

	engine := tfgen.NewEngine(mapperRegistry)
	output, err := engine.Generate(ctx, arch, sorted)
	if err != nil {
		return nil, fmt.Errorf("terraform engine generate: %w", err)
	}

	return output, nil
}

func writeTerraformOutput(dir string, out *iac.Output) error {
	if out == nil {
		return fmt.Errorf("nil terraform output")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create output dir %s: %w", dir, err)
	}

	for _, f := range out.Files {
		target := filepath.Join(dir, f.Path)
		if err := os.WriteFile(target, []byte(f.Content), 0o644); err != nil {
			return fmt.Errorf("write file %s: %w", target, err)
		}
	}

	return nil
}

func formatValidationErrors(result *validator.ValidationResult) string {
	if result == nil || len(result.Errors) == 0 {
		return ""
	}
	var b strings.Builder
	for _, e := range result.Errors {
		if e == nil {
			continue
		}
		b.WriteString("- ")
		b.WriteString(e.Code)
		if e.NodeID != "" {
			b.WriteString(" (node ")
			b.WriteString(e.NodeID)
			b.WriteString(")")
		}
		if e.Message != "" {
			b.WriteString(": ")
			b.WriteString(e.Message)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// resolveDiagramJSONPath resolves the JSON file path
// relative to the backend module root, regardless of the current working dir.
func resolveDiagramJSONPath(filename string) (string, error) {
	// Use runtime.Caller to get this file's directory, then walk up to the backend root.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine caller for resolving diagram JSON path")
	}

	// thisFile = .../backend/pkg/usecases/scenario6_terraform_with_persistence/terraform_with_persistence.go
	// backend root = thisFile/../../../..
	dir := filepath.Dir(thisFile)
	root := filepath.Clean(filepath.Join(dir, "..", "..", ".."))
	jsonPath := filepath.Join(root, filename)

	return jsonPath, nil
}

// extractDiagramFromProjectJSON extracts the diagram structure from project-wrapped JSON.
// Handles both formats:
//   - Direct format: {"nodes": [...], "edges": [...]}
//   - Project-wrapped: {"cloud-canvas-project-...": {"nodes": [...], "edges": [...]}}
func extractDiagramFromProjectJSON(data []byte) ([]byte, error) {
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Check if this is a direct diagram format (has "nodes" at root)
	if _, hasNodes := rawData["nodes"]; hasNodes {
		// Already in the correct format, return as-is
		return data, nil
	}

	// Otherwise, look for project-wrapped structure
	// Find the first key that contains a nested object with "nodes"
	for _, value := range rawData {
		if projectData, ok := value.(map[string]interface{}); ok {
			if _, hasNodes := projectData["nodes"]; hasNodes {
				// Found the diagram structure, extract it
				diagramBytes, err := json.Marshal(projectData)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal extracted diagram: %w", err)
				}
				return diagramBytes, nil
			}
		}
	}

	// If we get here, couldn't find the diagram structure
	return nil, fmt.Errorf("could not find diagram structure in JSON (expected 'nodes' field at root or nested under project key)")
}
