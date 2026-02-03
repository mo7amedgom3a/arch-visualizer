package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/orchestrator"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/services"

	// Register types
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/architecture"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/iam"

	"gorm.io/gorm"
)

// -- Adapters to fix interface mismatches --

type ProjectRepoAdapter struct {
	*repository.ProjectRepository
}

// BeginTransaction matches interface signature (returns interface{})
func (a *ProjectRepoAdapter) BeginTransaction(ctx context.Context) (interface{}, context.Context) {
	tx, txCtx := a.ProjectRepository.BeginTransaction(ctx)
	return tx, txCtx
}

func (a *ProjectRepoAdapter) CommitTransaction(tx interface{}) error {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("transaction is not *gorm.DB")
	}
	return a.ProjectRepository.CommitTransaction(gormTx)
}

func (a *ProjectRepoAdapter) RollbackTransaction(tx interface{}) error {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("transaction is not *gorm.DB")
	}
	return a.ProjectRepository.RollbackTransaction(gormTx)
}

type DependencyTypeRepoAdapter struct {
	*repository.DependencyTypeRepository
}

func (a *DependencyTypeRepoAdapter) Create(ctx context.Context, depType *models.DependencyType) error {
	return fmt.Errorf("Create not implemented in adapter (DependencyTypeRepository)")
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	ctx := context.Background()
	fmt.Println("Starting Full Flow Simulation (JSON -> DB -> Terraform)...")

	// 1. Connect to Database
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	fmt.Println("✓ Database connected")

	// 1.5 Clean up data to avoid migration issues (duplicates)
	// WARNING: This deletes data, suitable for this simulation script only
	if err := db.Exec("DELETE FROM resource_dependencies").Error; err != nil {
		fmt.Printf("Warning cleaning dependencies: %v\n", err)
	}
	if err := db.Exec("DELETE FROM resource_containments").Error; err != nil {
		fmt.Printf("Warning cleaning containments: %v\n", err)
	}
	if err := db.Exec("DELETE FROM resources").Error; err != nil {
		fmt.Printf("Warning cleaning resources: %v\n", err)
	}
	if err := db.Exec("DELETE FROM project_versions").Error; err != nil {
		fmt.Printf("Warning cleaning project_versions: %v\n", err)
	}
	if err := db.Exec("DELETE FROM projects").Error; err != nil {
		fmt.Printf("Warning cleaning projects: %v\n", err)
	}
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		fmt.Printf("Warning cleaning users: %v\n", err)
	}

	// 1.6 Manual Schema Update (Bypassing AutoMigrate issues)
	// Add missing columns to projects
	db.Exec("ALTER TABLE projects ADD COLUMN IF NOT EXISTS description text")
	db.Exec("ALTER TABLE projects ADD COLUMN IF NOT EXISTS thumbnail text")
	db.Exec("ALTER TABLE projects ADD COLUMN IF NOT EXISTS tags text[]")

	// Create project_versions table
	db.Exec(`CREATE TABLE IF NOT EXISTS project_versions (
		id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
		project_id uuid NOT NULL,
		created_at timestamptz DEFAULT now(),
		created_by uuid,
		changes text,
		snapshot jsonb
	)`)
	db.Exec("CREATE INDEX IF NOT EXISTS idx_project_versions_project_id ON project_versions (project_id)")

	fmt.Println("✓ Database schema manually updated")

	// 2. Seed Missing Resource Types (IAMPolicy, IAMRolePolicyAttachment)
	resourceTypes := []models.ResourceType{
		{Name: "IAMPolicy", CloudProvider: "aws", IsRegional: true},
		{Name: "IAMRolePolicyAttachment", CloudProvider: "aws", IsRegional: true},
		{Name: "Lambda", CloudProvider: "aws", IsRegional: true},      // Ensure exists
		{Name: "S3", CloudProvider: "aws", IsRegional: true},          // Ensure exists
		{Name: "GenericEdge", CloudProvider: "aws", IsRegional: true}, // New Type for Edges
	}
	for _, rt := range resourceTypes {
		var count int64
		if err := db.Model(&models.ResourceType{}).Where("name = ? AND cloud_provider = ?", rt.Name, rt.CloudProvider).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := db.Create(&rt).Error; err != nil {
				return fmt.Errorf("failed to seed resource type %s: %w", rt.Name, err)
			}
			fmt.Printf("✓ Seeded ResourceType: %s\n", rt.Name)
		}
	}

	// also seed "depends_on" dependency type
	var depTypeCount int64
	if err := db.Model(&models.DependencyType{}).Where("name = ?", "depends_on").Count(&depTypeCount).Error; err == nil && depTypeCount == 0 {
		if err := db.Create(&models.DependencyType{Name: "depends_on"}).Error; err != nil {
			fmt.Printf("Warning: failed to seed dependency type: %v\n", err)
		} else {
			fmt.Println("✓ Seeded DependencyType: depends_on")
		}
	}

	// 3. Initialize Repositories
	projectRepoRaw, err := repository.NewProjectRepository()
	if err != nil {
		return err
	}
	projectRepo := &ProjectRepoAdapter{projectRepoRaw}

	versionRepoRaw, err := repository.NewProjectVersionRepository()
	if err != nil {
		return err
	}
	// Use the adapter from services package if possible, or define locally if needed.
	// Since we are in main package and importing services, let's try to use services.ProjectVersionRepositoryAdapter?
	// But struct fields are not exported? No, "Repo" field is exported in services.ProjectVersionRepositoryAdapter
	versionRepo := &services.ProjectVersionRepositoryAdapter{Repo: versionRepoRaw}

	resourceRepo, err := repository.NewResourceRepository()
	if err != nil {
		return err
	}
	resourceTypeRepo, err := repository.NewResourceTypeRepository()
	if err != nil {
		return err
	}
	containmentRepo, err := repository.NewResourceContainmentRepository()
	if err != nil {
		return err
	}
	dependencyRepo, err := repository.NewResourceDependencyRepository()
	if err != nil {
		return err
	}

	dependencyTypeRepoRaw, err := repository.NewDependencyTypeRepository()
	if err != nil {
		return err
	}
	dependencyTypeRepo := &DependencyTypeRepoAdapter{dependencyTypeRepoRaw}

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		return err
	}
	iacTargetRepo, err := repository.NewIACTargetRepository()
	if err != nil {
		return err
	}

	// 4. Initialize Services
	diagramService := services.NewDiagramService()
	// Pass nil for ruleService as we don't need strict rule validation for this simulation
	architectureService := services.NewArchitectureService(nil)
	codegenService := services.NewCodegenService()
	projectService := services.NewProjectService(
		projectRepo, // Wrapped
		versionRepo, // Wrapper
		resourceRepo,
		resourceTypeRepo,
		containmentRepo,
		dependencyRepo,
		dependencyTypeRepo, // Wrapped
		userRepo,
		iacTargetRepo,
	)

	// 5. Initialize Orchestrator
	orchestrator := orchestrator.NewPipelineOrchestrator(
		diagramService,
		architectureService,
		codegenService,
		projectService,
	)

	// 6. Setup Test User
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		fmt.Printf("User %s not found, creating...\n", userID)
		newUser := &models.User{
			ID:         userID,
			Name:       "Simulation User",
			Email:      "simulation@example.com",
			IsVerified: true,
		}
		if err := userRepo.Create(ctx, newUser); err != nil {
			// If already exists (race condition or unique constraint), ignore
			fmt.Printf("Warning creating user: %v\n", err)
		}
		user = newUser
	}
	fmt.Printf("✓ Using User: %s (%s)\n", user.Name, user.ID)

	// 7. Read Diagram JSON
	jsonPath := "backend/json-request-edges-s3-lambda.json"
	// Check if file exists relative to where we run, might need adjustment
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		// Try absolute path or assume running from root
		jsonPath = "json-request-edges-s3-lambda.json"
	}

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		// Fallback for running from backend dir
		jsonData, err = os.ReadFile("json-request-edges-s3-lambda.json")
		if err != nil {
			return fmt.Errorf("failed to read diagram json: %w", err)
		}
	}
	fmt.Printf("✓ Read diagram JSON (%d bytes)\n", len(jsonData))

	// 8. Process Diagram (Save to DB)
	projectName := "Simulation Project Policies"
	processReq := &serverinterfaces.ProcessDiagramRequest{
		JSONData:      jsonData,
		UserID:        userID,
		ProjectName:   projectName,
		IACToolID:     1, // Terraform
		CloudProvider: "aws",
		Region:        "us-east-1",
	}

	result, err := orchestrator.ProcessDiagram(ctx, processReq)
	if err != nil {
		return fmt.Errorf("ProcessDiagram failed: %w", err)
	}
	fmt.Printf("✓ Diagram processed. Project ID: %s\n", result.ProjectID)

	// 8.5 Verify Edge Persistence in DB
	var edgeResource models.Resource
	// The edge ID from JSON is "edge-lambda-1-s3-3", used as Name if label is missing
	if err := db.Where("name = ?", "edge-lambda-1-s3-3").First(&edgeResource).Error; err != nil {
		return fmt.Errorf("failed to find edge resource in DB: %w", err)
	}
	fmt.Printf("✓ Found Edge Resource in DB: %s (Type: %s)\n", edgeResource.Name, edgeResource.ResourceTypeID)

	// 9. Generate Code
	genReq := &serverinterfaces.GenerateCodeRequest{
		ProjectID:     result.ProjectID,
		Engine:        "terraform",
		CloudProvider: "aws",
	}

	output, err := orchestrator.GenerateCode(ctx, genReq)
	if err != nil {
		return fmt.Errorf("GenerateCode failed: %w", err)
	}
	fmt.Printf("✓ Code generated. Files: %d\n", len(output.Files))

	// 10. Validate Output
	var mainTfContent string
	for _, f := range output.Files {
		if f.Path == "main.tf" {
			mainTfContent = f.Content
			break
		}
	}

	if mainTfContent == "" {
		return fmt.Errorf("main.tf not found in output")
	}

	fmt.Println("\n--- Validating Generated Terraform ---")

	checks := []struct {
		Term string
		Desc string
	}{
		{"aws_iam_policy", "IAM Policy Resource"},
		{"aws_iam_role_policy_attachment", "IAM Attachment Resource"},
		{"s3:GetObject", "Policy Action s3:GetObject"},
		{"\\\"Effect\\\":\\\"Allow\\\"", "Policy Effect"},
	}

	allPassed := true
	for _, check := range checks {
		// We might need to unescape json inside hcl string to match strictly, but substring check is usually enough
		if strings.Contains(mainTfContent, check.Term) {
			fmt.Printf("✓ Found %s\n", check.Desc)
		} else {
			fmt.Printf("❌ Missing %s\n", check.Desc)
			allPassed = false
		}
	}

	// Save to file for inspection
	outDir := "terraform_output_scenario15"
	os.MkdirAll(outDir, 0755)
	os.WriteFile(outDir+"/main.tf", []byte(mainTfContent), 0644)
	fmt.Printf("Saved main.tf to %s/main.tf\n", outDir)

	if !allPassed {
		return fmt.Errorf("simulation verification failed")
	}

	fmt.Println("SUCCESS! Full flow simulation passed.")
	return nil
}
