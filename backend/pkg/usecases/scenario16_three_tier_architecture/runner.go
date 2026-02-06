package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/services"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	// Register mappers
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/terraform"
)

// -- Adapters --

type ProjectRepoAdapter struct {
	*repository.ProjectRepository
}

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
	return fmt.Errorf("Create not implemented in adapter")
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	ctx := context.Background()
	logger := slog.Default()
	fmt.Println("Starting Scenario 16: Three-Tier Architecture (Database Persistence & Generation)...")

	// 1. Connect to Database
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	fmt.Println("✓ Database connected")

	// Clean up previous run data for this specific project name to avoid clutter/dupes if re-run without ID cleaning
	// For now, simpler to just create valid data.

	// 2. Initialize Infrastructure (Repos & Services)
	projectRepoRaw, _ := repository.NewProjectRepository(logger)
	projectRepo := &ProjectRepoAdapter{projectRepoRaw}
	versionRepoRaw, _ := repository.NewProjectVersionRepository()
	versionRepo := &services.ProjectVersionRepositoryAdapter{Repo: versionRepoRaw}
	resourceRepo, _ := repository.NewResourceRepository(logger)
	resourceTypeRepo, _ := repository.NewResourceTypeRepository()
	containmentRepo, _ := repository.NewResourceContainmentRepository()
	dependencyRepo, _ := repository.NewResourceDependencyRepository()
	dependencyTypeRepoRaw, _ := repository.NewDependencyTypeRepository()
	dependencyTypeRepo := &DependencyTypeRepoAdapter{dependencyTypeRepoRaw}
	userRepo, _ := repository.NewUserRepository()
	iacTargetRepo, _ := repository.NewIACTargetRepository()
	variableRepo, _ := repository.NewProjectVariableRepository()
	outputRepo, _ := repository.NewProjectOutputRepository()

	codegenService := services.NewCodegenService(logger)
	projectService := services.NewProjectService(
		projectRepo, versionRepo, resourceRepo, resourceTypeRepo,
		containmentRepo, dependencyRepo, dependencyTypeRepo,
		userRepo, iacTargetRepo, variableRepo, outputRepo,
	)

	// 3. Setup Data
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	// Ensure User
	var user models.User
	if err := db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		user = models.User{ID: userID, Name: "Demo User", IsVerified: true}
		db.WithContext(ctx).Create(&user)
	}

	var iacTarget models.IACTarget
	db.WithContext(ctx).FirstOrCreate(&iacTarget, models.IACTarget{Name: "Terraform"})

	project := &models.Project{
		ID:            uuid.New(),
		UserID:        userID,
		InfraToolID:   iacTarget.ID,
		Name:          "Three-Tier Architecture Project",
		CloudProvider: "aws",
		Region:        "us-east-1",
	}
	if err := db.WithContext(ctx).Create(project).Error; err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	fmt.Printf("✓ Created Project: %s (%s)\n", project.Name, project.ID)

	// Helper to get or create Category/Kind
	getOrCreateCategory := func(name string) (uint, error) {
		var cat models.ResourceCategory
		if err := db.WithContext(ctx).FirstOrCreate(&cat, models.ResourceCategory{Name: name}).Error; err != nil {
			return 0, err
		}
		return cat.ID, nil
	}
	getOrCreateKind := func(name string) (uint, error) {
		var kind models.ResourceKind
		if err := db.WithContext(ctx).FirstOrCreate(&kind, models.ResourceKind{Name: name}).Error; err != nil {
			return 0, err
		}
		return kind.ID, nil
	}

	// Helper to create resource
	mapper := architecture.NewAWSResourceTypeMapper()
	createResource := func(name, typeName string, metadata map[string]interface{}) (*models.Resource, error) {
		rtDomain, err := mapper.MapResourceNameToResourceType(typeName)
		if err != nil {
			return nil, err
		}

		// Ensure ResourceType exists
		var rt models.ResourceType
		if err := db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", rtDomain.Name, "aws").First(&rt).Error; err != nil {
			// Create it
			catID, _ := getOrCreateCategory(rtDomain.Category)
			kindID, _ := getOrCreateKind(rtDomain.Kind)

			rt = models.ResourceType{
				Name:          rtDomain.Name,
				CloudProvider: "aws",
				CategoryID:    &catID,
				KindID:        &kindID,
				IsRegional:    rtDomain.IsRegional,
				IsGlobal:      rtDomain.IsGlobal,
			}
			if err := db.WithContext(ctx).Create(&rt).Error; err != nil {
				return nil, fmt.Errorf("create resource type %s: %w", rtDomain.Name, err)
			}
		}

		metaJSON, _ := json.Marshal(metadata)
		res := &models.Resource{
			ID:             uuid.New(),
			ProjectID:      project.ID,
			ResourceTypeID: rt.ID,
			Name:           name,
			Config:         datatypes.JSON(metaJSON),
		}
		if err := db.WithContext(ctx).Create(res).Error; err != nil {
			return nil, err
		}
		return res, nil
	}

	// 4. Create Resources
	fmt.Println("Creating Resources...")

	// VPC
	vpc, _ := createResource("main-vpc", "VPC", map[string]interface{}{
		"cidr":               "10.0.0.0/16",
		"enableDnsHostnames": true,
		"enableDnsSupport":   true,
	})

	// Subnets
	pubSub1, _ := createResource("public-subnet-1", "Subnet", map[string]interface{}{
		"cidr": "10.0.1.0/24", "availabilityZoneId": "us-east-1a", "map_public_ip_on_launch": true,
	})
	pubSub2, _ := createResource("public-subnet-2", "Subnet", map[string]interface{}{
		"cidr": "10.0.2.0/24", "availabilityZoneId": "us-east-1b", "map_public_ip_on_launch": true,
	})
	privSub1, _ := createResource("private-subnet-1", "Subnet", map[string]interface{}{
		"cidr": "10.0.3.0/24", "availabilityZoneId": "us-east-1a",
	})
	privSub2, _ := createResource("private-subnet-2", "Subnet", map[string]interface{}{
		"cidr": "10.0.4.0/24", "availabilityZoneId": "us-east-1b",
	})

	// Gateways
	igw, _ := createResource("main-igw", "InternetGateway", nil)
	// NAT Gateway needs a subnet (usually public) to reside in
	natGw, _ := createResource("nat-gw", "NATGateway", map[string]interface{}{
		"subnetId": pubSub1.ID.String(),
	})

	// Route Tables
	pubRT, _ := createResource("public-rtb", "RouteTable", nil)
	privRT, _ := createResource("private-rtb", "RouteTable", nil)

	// Security Groups
	ec2SG, _ := createResource("ec2-sg", "SecurityGroup", map[string]interface{}{
		"description": "Allow HTTP",
		"rules": []map[string]interface{}{
			{"type": "ingress", "protocol": "tcp", "portRange": "80", "cidr": "0.0.0.0/0", "description": "HTTP"},
		},
	})
	rdsSG, _ := createResource("rds-sg", "SecurityGroup", map[string]interface{}{"description": "Allow EC2"})

	// Compute & DB
	ec2, _ := createResource("web-server", "EC2", map[string]interface{}{
		"instanceType": "t3.micro", "ami": "ami-12345678",
		"keyName": "my-key",
	})

	rdsPrimary, _ := createResource("primary-db", "RDS", map[string]interface{}{
		"engine": "postgres", "instance_class": "db.t3.micro", "allocated_storage": 20,
		"engine_version": "13.7",
		"db_name":        "mydb", "username": "admin", "password": "password", "multi_az": true,
		"backup_retention_period": 7,
	})

	// 5. Containments
	fmt.Println("Establishing Hierarchy...")
	createContainment := func(parent, child *models.Resource) {
		db.WithContext(ctx).Create(&models.ResourceContainment{
			ParentResourceID: parent.ID, ChildResourceID: child.ID,
		})
	}

	// VPC Hierarchy
	for _, res := range []*models.Resource{pubSub1, pubSub2, privSub1, privSub2, igw, pubRT, privRT, ec2SG, rdsSG} {
		createContainment(vpc, res)
	}

	// Subnet Placements
	// NAT GW is technically "in" a subnet, but containment vs depends_on/property is tricky.
	// AWS adapter expects containment for subnet -> resource mapping for some things, but NAT GW has explicit subnetId property too.
	// For visualizer, containment usually implies "deployed in".
	createContainment(pubSub1, natGw)
	createContainment(privSub1, ec2)
	createContainment(privSub1, rdsPrimary)

	// 6. Dependencies (Explicit)
	fmt.Println("Linking Resources...")
	var depType models.DependencyType
	if err := db.WithContext(ctx).FirstOrCreate(&depType, models.DependencyType{Name: "depends_on"}).Error; err != nil {
		return fmt.Errorf("create dependency type: %w", err)
	}

	createDependency := func(from, to *models.Resource) {
		db.WithContext(ctx).Create(&models.ResourceDependency{
			FromResourceID: from.ID, ToResourceID: to.ID, DependencyTypeID: depType.ID,
		})
	}

	// EC2 depends on RDS (simulating app dependency)
	createDependency(ec2, rdsPrimary)

	// Security Group Rule: RDS SG allows traffic from EC2 SG
	// This is typically done via ingress rule referencing SG ID.
	// Updating RDS SG metadata to include rule
	rdsSGRules := []map[string]interface{}{
		{"type": "ingress", "protocol": "tcp", "portRange": "5432", "sourceSecurityGroupId": ec2SG.ID.String(), "description": "PostgreSQL"}, // Reference EC2 SG
	}
	rdsSGRulesJSON, _ := json.Marshal(map[string]interface{}{
		"description": "Allow EC2",
		"rules":       rdsSGRules,
	})
	db.Model(rdsSG).Update("config", datatypes.JSON(rdsSGRulesJSON))

	// 7. Load Architecture via Service (Tests loading layer)
	fmt.Println("Loading Architecture from DB...")
	archDomain, err := projectService.LoadArchitecture(ctx, project.ID)
	if err != nil {
		return fmt.Errorf("LoadArchitecture failed: %w", err)
	}
	fmt.Printf("✓ Loaded Architecture. Resources: %d\n", len(archDomain.Resources))

	// 8. Generate Terraform (Tests Generation layer)
	fmt.Println("Generating Terraform...")
	output, err := codegenService.Generate(ctx, archDomain, "terraform")
	if err != nil {
		return fmt.Errorf("Codegen failed: %w", err)
	}

	// 9. Output results
	outDir := "terraform_output_scenario16"
	os.MkdirAll(outDir, 0755)

	for _, f := range output.Files {
		path := fmt.Sprintf("%s/%s", outDir, f.Path)
		if err := os.WriteFile(path, []byte(f.Content), 0644); err != nil {
			return err
		}
		fmt.Printf("✓ Wrote %s\n", path)
		if f.Path == "main.tf" {
			// Print snippet
			fmt.Println("\n--- main.tf snippet ---")
			if len(f.Content) > 500 {
				fmt.Println(f.Content[:500] + "...\n(truncated)")
			} else {
				fmt.Println(f.Content)
			}
		}
	}

	fmt.Println("SUCCESS: Scenario 16 Completed.")
	return nil
}
