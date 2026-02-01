package seed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"gorm.io/gorm"
)

// SeedDatabase seeds the database with realistic use case data
func SeedDatabase() error {
	ctx := context.Background()

	// Connect to database
	if _, err := database.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("🌱 Starting database seeding...")

	// Seed reference data first
	if err := seedReferenceData(ctx); err != nil {
		return fmt.Errorf("failed to seed reference data: %w", err)
	}

	// Seed users
	users, err := seedUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// Seed marketplace data
	if err := seedMarketplaceData(ctx, users); err != nil {
		return fmt.Errorf("failed to seed marketplace data: %w", err)
	}

	// Seed pricing data (rates and hidden dependencies)
	if err := SeedPricingData(ctx); err != nil {
		return fmt.Errorf("failed to seed pricing data: %w", err)
	}

	// Seed projects and resources based on scenarios
	if err := seedScenarios(ctx, users); err != nil {
		return fmt.Errorf("failed to seed scenarios: %w", err)
	}

	fmt.Println("✅ Database seeding completed successfully!")
	return nil
}

// seedReferenceData seeds categories, kinds, types, dependency types, and IAC targets
func seedReferenceData(ctx context.Context) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	fmt.Println("\n📋 Seeding reference data...")

	// Seed Resource Categories
	categories := []models.ResourceCategory{
		{Name: "Compute"},
		{Name: "Networking"},
		{Name: "Storage"},
		{Name: "Database"},
		{Name: "Security"},
		{Name: "Analytics"},
		{Name: "Application Integration"},
	}
	for _, cat := range categories {
		if err := db.WithContext(ctx).FirstOrCreate(&cat, models.ResourceCategory{Name: cat.Name}).Error; err != nil {
			return fmt.Errorf("failed to seed category %s: %w", cat.Name, err)
		}
	}
	fmt.Printf("  ✓ Seeded %d resource categories\n", len(categories))

	// Seed Resource Kinds
	kinds := []models.ResourceKind{
		{Name: "VirtualMachine"},
		{Name: "Container"},
		{Name: "Function"},
		{Name: "Network"},
		{Name: "LoadBalancer"},
		{Name: "Database"},
		{Name: "Storage"},
		{Name: "Gateway"},
	}
	for _, kind := range kinds {
		if err := db.WithContext(ctx).FirstOrCreate(&kind, models.ResourceKind{Name: kind.Name}).Error; err != nil {
			return fmt.Errorf("failed to seed kind %s: %w", kind.Name, err)
		}
	}
	fmt.Printf("  ✓ Seeded %d resource kinds\n", len(kinds))

	// Get category and kind IDs
	var computeCat, networkCat, storageCat, dbCat, securityCat models.ResourceCategory
	db.WithContext(ctx).Where("name = ?", "Compute").First(&computeCat)
	db.WithContext(ctx).Where("name = ?", "Networking").First(&networkCat)
	db.WithContext(ctx).Where("name = ?", "Storage").First(&storageCat)
	db.WithContext(ctx).Where("name = ?", "Database").First(&dbCat)
	db.WithContext(ctx).Where("name = ?", "Security").First(&securityCat)

	var vmKind, containerKind, functionKind, networkKind, lbKind, dbKind, storageKind, gatewayKind models.ResourceKind
	db.WithContext(ctx).Where("name = ?", "VirtualMachine").First(&vmKind)
	db.WithContext(ctx).Where("name = ?", "Container").First(&containerKind)
	db.WithContext(ctx).Where("name = ?", "Function").First(&functionKind)
	db.WithContext(ctx).Where("name = ?", "Network").First(&networkKind)
	db.WithContext(ctx).Where("name = ?", "LoadBalancer").First(&lbKind)
	db.WithContext(ctx).Where("name = ?", "Database").First(&dbKind)
	db.WithContext(ctx).Where("name = ?", "Storage").First(&storageKind)
	db.WithContext(ctx).Where("name = ?", "Gateway").First(&gatewayKind)

	// Seed Resource Types (AWS)
	resourceTypes := []models.ResourceType{
		// Compute
		{Name: "EC2", CloudProvider: "aws", CategoryID: &computeCat.ID, KindID: &vmKind.ID, IsRegional: true, IsGlobal: false},
		{Name: "Lambda", CloudProvider: "aws", CategoryID: &computeCat.ID, KindID: &functionKind.ID, IsRegional: true, IsGlobal: false},
		{Name: "ECS", CloudProvider: "aws", CategoryID: &computeCat.ID, KindID: &containerKind.ID, IsRegional: true, IsGlobal: false},
		{Name: "EKS", CloudProvider: "aws", CategoryID: &computeCat.ID, KindID: &containerKind.ID, IsRegional: true, IsGlobal: false},
		{Name: "AutoScalingGroup", CloudProvider: "aws", CategoryID: &computeCat.ID, KindID: &vmKind.ID, IsRegional: true, IsGlobal: false},
		{Name: "LoadBalancer", CloudProvider: "aws", CategoryID: &computeCat.ID, KindID: &lbKind.ID, IsRegional: true, IsGlobal: false},
		// Networking
		{Name: "VPC", CloudProvider: "aws", CategoryID: &networkCat.ID, KindID: &networkKind.ID, IsRegional: true, IsGlobal: false},
		{Name: "Subnet", CloudProvider: "aws", CategoryID: &networkCat.ID, KindID: &networkKind.ID, IsRegional: true, IsGlobal: false},
		{Name: "InternetGateway", CloudProvider: "aws", CategoryID: &networkCat.ID, KindID: &gatewayKind.ID, IsRegional: true, IsGlobal: false},
		{Name: "NATGateway", CloudProvider: "aws", CategoryID: &networkCat.ID, KindID: &gatewayKind.ID, IsRegional: true, IsGlobal: false},
		{Name: "RouteTable", CloudProvider: "aws", CategoryID: &networkCat.ID, KindID: &networkKind.ID, IsRegional: true, IsGlobal: false},
		{Name: "SecurityGroup", CloudProvider: "aws", CategoryID: &networkCat.ID, KindID: &networkKind.ID, IsRegional: true, IsGlobal: false},
		{Name: "ElasticIP", CloudProvider: "aws", CategoryID: &networkCat.ID, KindID: &networkKind.ID, IsRegional: true, IsGlobal: false},
		// Storage
		{Name: "S3", CloudProvider: "aws", CategoryID: &storageCat.ID, KindID: &storageKind.ID, IsRegional: false, IsGlobal: true},
		{Name: "EBS", CloudProvider: "aws", CategoryID: &storageCat.ID, KindID: &storageKind.ID, IsRegional: true, IsGlobal: false},
		// Database
		{Name: "RDS", CloudProvider: "aws", CategoryID: &dbCat.ID, KindID: &dbKind.ID, IsRegional: true, IsGlobal: false},
		{Name: "DynamoDB", CloudProvider: "aws", CategoryID: &dbCat.ID, KindID: &dbKind.ID, IsRegional: true, IsGlobal: false},
	}
	for _, rt := range resourceTypes {
		if err := db.WithContext(ctx).FirstOrCreate(&rt, models.ResourceType{Name: rt.Name, CloudProvider: rt.CloudProvider}).Error; err != nil {
			return fmt.Errorf("failed to seed resource type %s: %w", rt.Name, err)
		}
	}
	fmt.Printf("  ✓ Seeded %d resource types\n", len(resourceTypes))

	// Seed Dependency Types
	dependencyTypes := []models.DependencyType{
		{Name: "uses"},
		{Name: "depends_on"},
		{Name: "connects_to"},
		{Name: "references"},
		{Name: "contains"},
	}
	for _, dt := range dependencyTypes {
		if err := db.WithContext(ctx).FirstOrCreate(&dt, models.DependencyType{Name: dt.Name}).Error; err != nil {
			return fmt.Errorf("failed to seed dependency type %s: %w", dt.Name, err)
		}
	}
	fmt.Printf("  ✓ Seeded %d dependency types\n", len(dependencyTypes))

	// Seed IAC Targets
	iacTargets := []models.IACTarget{
		{Name: "Terraform"},
		{Name: "Pulumi"},
		{Name: "CDK"},
		{Name: "CloudFormation"},
	}
	for _, iac := range iacTargets {
		if err := db.WithContext(ctx).FirstOrCreate(&iac, models.IACTarget{Name: iac.Name}).Error; err != nil {
			return fmt.Errorf("failed to seed IAC target %s: %w", iac.Name, err)
		}
	}
	fmt.Printf("  ✓ Seeded %d IAC targets\n", len(iacTargets))

	return nil
}

// seedUsers creates sample users
func seedUsers(ctx context.Context) ([]*models.User, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	fmt.Println("\n👥 Seeding users...")

	users := []*models.User{
		{Name: "Alice Johnson", Avatar: stringPtr("https://example.com/avatars/alice.png"), IsVerified: true},
		{Name: "Bob Smith", Avatar: stringPtr("https://example.com/avatars/bob.png"), IsVerified: false},
		{Name: "Charlie Brown", Avatar: stringPtr("https://example.com/avatars/charlie.png"), IsVerified: true},
	}

	var createdUsers []*models.User
	for _, user := range users {
		var existing models.User
		err := db.WithContext(ctx).Where("name = ?", user.Name).First(&existing).Error
		if err == nil {
			fmt.Printf("  ⚠ User %s already exists, skipping\n", user.Name)
			existingUser := existing
			createdUsers = append(createdUsers, &existingUser)
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check user %s: %w", user.Name, err)
		}

		if err := db.WithContext(ctx).Create(user).Error; err != nil {
			return nil, fmt.Errorf("failed to create user %s: %w", user.Name, err)
		}
		createdUsers = append(createdUsers, user)
		fmt.Printf("  ✓ Created user: %s\n", user.Name)
	}

	return createdUsers, nil
}

// seedMarketplaceData seeds marketplace tables (categories, templates, and related data)
func seedMarketplaceData(ctx context.Context, users []*models.User) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	fmt.Println("\n🛍️  Seeding marketplace data...")

	categories := []models.Category{
		{Name: "Web Applications", Slug: "web-applications"},
		{Name: "Data Analytics", Slug: "data-analytics"},
		{Name: "Security", Slug: "security"},
	}
	for _, cat := range categories {
		if err := db.WithContext(ctx).FirstOrCreate(&cat, models.Category{Slug: cat.Slug}).Error; err != nil {
			return fmt.Errorf("failed to seed marketplace category %s: %w", cat.Name, err)
		}
	}
	fmt.Printf("  ✓ Seeded %d marketplace categories\n", len(categories))

	technologies := []models.Technology{
		{Name: "Kubernetes", Slug: "kubernetes"},
		{Name: "PostgreSQL", Slug: "postgresql"},
		{Name: "Redis", Slug: "redis"},
		{Name: "Amazon S3", Slug: "amazon-s3"},
		{Name: "AWS Glue", Slug: "aws-glue"},
	}
	for _, tech := range technologies {
		if err := db.WithContext(ctx).FirstOrCreate(&tech, models.Technology{Slug: tech.Slug}).Error; err != nil {
			return fmt.Errorf("failed to seed marketplace technology %s: %w", tech.Name, err)
		}
	}
	fmt.Printf("  ✓ Seeded %d marketplace technologies\n", len(technologies))

	iacFormats := []models.IACFormat{
		{Name: "Terraform", Slug: "terraform"},
		{Name: "Pulumi", Slug: "pulumi"},
		{Name: "CloudFormation", Slug: "cloudformation"},
	}
	for _, format := range iacFormats {
		if err := db.WithContext(ctx).FirstOrCreate(&format, models.IACFormat{Slug: format.Slug}).Error; err != nil {
			return fmt.Errorf("failed to seed marketplace IaC format %s: %w", format.Name, err)
		}
	}
	fmt.Printf("  ✓ Seeded %d marketplace IaC formats\n", len(iacFormats))

	complianceStandards := []models.ComplianceStandard{
		{Name: "SOC 2", Slug: "soc-2"},
		{Name: "HIPAA", Slug: "hipaa"},
	}
	for _, compliance := range complianceStandards {
		if err := db.WithContext(ctx).FirstOrCreate(&compliance, models.ComplianceStandard{Slug: compliance.Slug}).Error; err != nil {
			return fmt.Errorf("failed to seed compliance standard %s: %w", compliance.Name, err)
		}
	}
	fmt.Printf("  ✓ Seeded %d compliance standards\n", len(complianceStandards))

	if len(users) == 0 {
		return nil
	}

	var webCat, dataCat models.Category
	db.WithContext(ctx).Where("slug = ?", "web-applications").First(&webCat)
	db.WithContext(ctx).Where("slug = ?", "data-analytics").First(&dataCat)

	var terraformFormat, pulumiFormat models.IACFormat
	db.WithContext(ctx).Where("slug = ?", "terraform").First(&terraformFormat)
	db.WithContext(ctx).Where("slug = ?", "pulumi").First(&pulumiFormat)

	var soc2Standard models.ComplianceStandard
	db.WithContext(ctx).Where("slug = ?", "soc-2").First(&soc2Standard)

	var kubernetesTech, postgresTech, redisTech, s3Tech, glueTech models.Technology
	db.WithContext(ctx).Where("slug = ?", "kubernetes").First(&kubernetesTech)
	db.WithContext(ctx).Where("slug = ?", "postgresql").First(&postgresTech)
	db.WithContext(ctx).Where("slug = ?", "redis").First(&redisTech)
	db.WithContext(ctx).Where("slug = ?", "amazon-s3").First(&s3Tech)
	db.WithContext(ctx).Where("slug = ?", "aws-glue").First(&glueTech)

	authorOne := users[0]
	authorTwo := users[0]
	if len(users) > 1 {
		authorTwo = users[1]
	}

	webAppTemplate := &models.Template{
		Title:            "Scalable Web App",
		Description:      "Highly available web application stack with autoscaling and caching.",
		CategoryID:       webCat.ID,
		CloudProvider:    "AWS",
		EstimatedCostMin: 50,
		EstimatedCostMax: 250,
		AuthorID:         authorOne.ID,
		ImageURL:         stringPtr("https://example.com/images/scalable-web-app.png"),
		IsPopular:        true,
		IsNew:            true,
		LastUpdated:      time.Now(),
		Resources:        8,
		DeploymentTime:   stringPtr("30-45 minutes"),
		Regions:          stringPtr("us-east-1, us-west-2"),
	}
	if err := db.WithContext(ctx).Where("title = ? AND author_id = ?", webAppTemplate.Title, webAppTemplate.AuthorID).FirstOrCreate(webAppTemplate).Error; err != nil {
		return fmt.Errorf("failed to seed template %s: %w", webAppTemplate.Title, err)
	}

	dataLakeTemplate := &models.Template{
		Title:            "Data Lake Starter",
		Description:      "Starter data lake with ingestion, catalog, and analytics components.",
		CategoryID:       dataCat.ID,
		CloudProvider:    "AWS",
		EstimatedCostMin: 75,
		EstimatedCostMax: 300,
		AuthorID:         authorTwo.ID,
		ImageURL:         stringPtr("https://example.com/images/data-lake-starter.png"),
		IsPopular:        false,
		IsNew:            true,
		LastUpdated:      time.Now(),
		Resources:        6,
		DeploymentTime:   stringPtr("45-60 minutes"),
		Regions:          stringPtr("us-east-1"),
	}
	if err := db.WithContext(ctx).Where("title = ? AND author_id = ?", dataLakeTemplate.Title, dataLakeTemplate.AuthorID).FirstOrCreate(dataLakeTemplate).Error; err != nil {
		return fmt.Errorf("failed to seed template %s: %w", dataLakeTemplate.Title, err)
	}

	if err := db.WithContext(ctx).Model(webAppTemplate).Association("Technologies").Replace(
		&kubernetesTech, &postgresTech, &redisTech,
	); err != nil {
		return fmt.Errorf("failed to seed template technologies for %s: %w", webAppTemplate.Title, err)
	}
	if err := db.WithContext(ctx).Model(webAppTemplate).Association("IACFormats").Replace(&terraformFormat); err != nil {
		return fmt.Errorf("failed to seed template IaC formats for %s: %w", webAppTemplate.Title, err)
	}
	if err := db.WithContext(ctx).Model(webAppTemplate).Association("ComplianceStandards").Replace(&soc2Standard); err != nil {
		return fmt.Errorf("failed to seed template compliance for %s: %w", webAppTemplate.Title, err)
	}

	if err := db.WithContext(ctx).Model(dataLakeTemplate).Association("Technologies").Replace(&s3Tech, &glueTech); err != nil {
		return fmt.Errorf("failed to seed template technologies for %s: %w", dataLakeTemplate.Title, err)
	}
	if err := db.WithContext(ctx).Model(dataLakeTemplate).Association("IACFormats").Replace(&pulumiFormat); err != nil {
		return fmt.Errorf("failed to seed template IaC formats for %s: %w", dataLakeTemplate.Title, err)
	}

	webUseCases := []models.TemplateUseCase{
		{TemplateID: webAppTemplate.ID, Icon: stringPtr("rocket"), Title: "Launch MVPs", Description: stringPtr("Deploy quickly with production-ready defaults."), DisplayOrder: 1},
		{TemplateID: webAppTemplate.ID, Icon: stringPtr("scale"), Title: "Auto Scaling", Description: stringPtr("Scale seamlessly with traffic spikes."), DisplayOrder: 2},
	}
	for _, useCase := range webUseCases {
		if err := db.WithContext(ctx).FirstOrCreate(&useCase, models.TemplateUseCase{TemplateID: useCase.TemplateID, Title: useCase.Title}).Error; err != nil {
			return fmt.Errorf("failed to seed template use case %s: %w", useCase.Title, err)
		}
	}

	webFeatures := []models.TemplateFeature{
		{TemplateID: webAppTemplate.ID, Feature: "Multi-AZ load balancer with autoscaling", DisplayOrder: 1},
		{TemplateID: webAppTemplate.ID, Feature: "Managed PostgreSQL with read replicas", DisplayOrder: 2},
		{TemplateID: webAppTemplate.ID, Feature: "Redis cache for session storage", DisplayOrder: 3},
	}
	for _, feature := range webFeatures {
		if err := db.WithContext(ctx).FirstOrCreate(&feature, models.TemplateFeature{TemplateID: feature.TemplateID, Feature: feature.Feature}).Error; err != nil {
			return fmt.Errorf("failed to seed template feature %s: %w", feature.Feature, err)
		}
	}

	webComponents := []models.TemplateComponent{
		{TemplateID: webAppTemplate.ID, Name: "Application Load Balancer", Service: "ELB", MonthlyCost: 18.5, Purpose: stringPtr("Distribute traffic across instances"), DisplayOrder: 1},
		{TemplateID: webAppTemplate.ID, Name: "Auto Scaling Group", Service: "EC2", MonthlyCost: 120.0, Purpose: stringPtr("Scale compute capacity"), DisplayOrder: 2},
		{TemplateID: webAppTemplate.ID, Name: "PostgreSQL Cluster", Service: "RDS", MonthlyCost: 65.0, Purpose: stringPtr("Primary data store"), DisplayOrder: 3},
	}
	for _, component := range webComponents {
		if err := db.WithContext(ctx).FirstOrCreate(&component, models.TemplateComponent{TemplateID: component.TemplateID, Name: component.Name}).Error; err != nil {
			return fmt.Errorf("failed to seed template component %s: %w", component.Name, err)
		}
	}

	dataFeatures := []models.TemplateFeature{
		{TemplateID: dataLakeTemplate.ID, Feature: "S3 data lake with lifecycle policies", DisplayOrder: 1},
		{TemplateID: dataLakeTemplate.ID, Feature: "Glue catalog for metadata management", DisplayOrder: 2},
		{TemplateID: dataLakeTemplate.ID, Feature: "Athena-ready query layer", DisplayOrder: 3},
	}
	for _, feature := range dataFeatures {
		if err := db.WithContext(ctx).FirstOrCreate(&feature, models.TemplateFeature{TemplateID: feature.TemplateID, Feature: feature.Feature}).Error; err != nil {
			return fmt.Errorf("failed to seed template feature %s: %w", feature.Feature, err)
		}
	}

	if len(users) > 1 {
		review := models.Review{
			TemplateID:     webAppTemplate.ID,
			UserID:         users[1].ID,
			Rating:         5,
			Title:          "Great starting point",
			Content:        "The template saved us a ton of setup time and scaled well.",
			UseCase:        stringPtr("SaaS web app"),
			TeamSize:       stringPtr("5-10"),
			DeploymentTime: stringPtr("40 minutes"),
		}
		if err := db.WithContext(ctx).FirstOrCreate(&review, models.Review{TemplateID: review.TemplateID, UserID: review.UserID}).Error; err != nil {
			return fmt.Errorf("failed to seed review: %w", err)
		}
	}

	fmt.Println("  ✓ Seeded marketplace templates and related data")
	return nil
}

// seedScenarios creates projects and resources based on the use case scenarios
func seedScenarios(ctx context.Context, users []*models.User) error {
	projectRepo, err := repository.NewProjectRepository()
	if err != nil {
		return err
	}

	resourceRepo, err := repository.NewResourceRepository()
	if err != nil {
		return err
	}

	db, err := database.Connect()
	if err != nil {
		return err
	}

	fmt.Println("\n🏗️  Seeding projects and resources...")

	// Get resource types
	var vpcType, subnetType, igwType, sgType, ec2Type, lbType, asgType, natType, eipType models.ResourceType
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "VPC", "aws").First(&vpcType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "Subnet", "aws").First(&subnetType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "InternetGateway", "aws").First(&igwType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "SecurityGroup", "aws").First(&sgType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "EC2", "aws").First(&ec2Type)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "LoadBalancer", "aws").First(&lbType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "AutoScalingGroup", "aws").First(&asgType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "NATGateway", "aws").First(&natType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "ElasticIP", "aws").First(&eipType)

	// Get IAC target (Terraform)
	var terraformTarget models.IACTarget
	db.WithContext(ctx).Where("name = ?", "Terraform").First(&terraformTarget)

	// Scenario 1: Basic Web Application
	if len(users) > 0 {
		project1 := &models.Project{
			UserID:        users[0].ID,
			InfraToolID:   terraformTarget.ID,
			Name:          "Basic Web Application",
			CloudProvider: "aws",
			Region:        "us-east-1",
		}
		if err := projectRepo.Create(ctx, project1); err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
		fmt.Printf("  ✓ Created project: %s\n", project1.Name)

		// Create VPC
		vpcConfigJSON, _ := json.Marshal(map[string]interface{}{"cidr": "10.0.0.0/16"})
		vpc := &models.Resource{
			ProjectID:      project1.ID,
			ResourceTypeID: vpcType.ID,
			Name:           "web-app-vpc",
			PosX:           100,
			PosY:           100,
			Config:         vpcConfigJSON,
		}
		if err := resourceRepo.Create(ctx, vpc); err != nil {
			return fmt.Errorf("failed to create VPC: %w", err)
		}

		// Create Subnets
		publicSubnet1ConfigJSON, _ := json.Marshal(map[string]interface{}{"cidr": "10.0.1.0/24", "availability_zone": "us-east-1a", "public": true})
		publicSubnet1 := &models.Resource{
			ProjectID:      project1.ID,
			ResourceTypeID: subnetType.ID,
			Name:           "public-subnet-1",
			PosX:           200,
			PosY:           150,
			Config:         publicSubnet1ConfigJSON,
		}
		if err := resourceRepo.Create(ctx, publicSubnet1); err != nil {
			return fmt.Errorf("failed to create public subnet: %w", err)
		}
		if err := resourceRepo.CreateContainment(ctx, vpc.ID, publicSubnet1.ID); err != nil {
			return fmt.Errorf("failed to create containment: %w", err)
		}

		privateSubnet1ConfigJSON, _ := json.Marshal(map[string]interface{}{"cidr": "10.0.10.0/24", "availability_zone": "us-east-1a", "public": false})
		privateSubnet1 := &models.Resource{
			ProjectID:      project1.ID,
			ResourceTypeID: subnetType.ID,
			Name:           "private-subnet-1",
			PosX:           200,
			PosY:           250,
			Config:         privateSubnet1ConfigJSON,
		}
		if err := resourceRepo.Create(ctx, privateSubnet1); err != nil {
			return fmt.Errorf("failed to create private subnet: %w", err)
		}
		if err := resourceRepo.CreateContainment(ctx, vpc.ID, privateSubnet1.ID); err != nil {
			return fmt.Errorf("failed to create containment: %w", err)
		}

		// Create Internet Gateway
		igwConfigJSON, _ := json.Marshal(map[string]interface{}{"attached": true})
		igw := &models.Resource{
			ProjectID:      project1.ID,
			ResourceTypeID: igwType.ID,
			Name:           "web-app-igw",
			PosX:           50,
			PosY:           100,
			Config:         igwConfigJSON,
		}
		if err := resourceRepo.Create(ctx, igw); err != nil {
			return fmt.Errorf("failed to create IGW: %w", err)
		}

		// Create Security Groups
		webSGConfigJSON, _ := json.Marshal(map[string]interface{}{"description": "Security group for web tier"})
		webSG := &models.Resource{
			ProjectID:      project1.ID,
			ResourceTypeID: sgType.ID,
			Name:           "web-sg",
			PosX:           300,
			PosY:           150,
			Config:         webSGConfigJSON,
		}
		if err := resourceRepo.Create(ctx, webSG); err != nil {
			return fmt.Errorf("failed to create security group: %w", err)
		}

		// Create EC2 Instances
		ec2ConfigJSON, _ := json.Marshal(map[string]interface{}{"instance_type": "t3.micro", "ami": "ami-0c55b159cbfafe1f0"})
		ec2 := &models.Resource{
			ProjectID:      project1.ID,
			ResourceTypeID: ec2Type.ID,
			Name:           "web-server-1",
			PosX:           350,
			PosY:           150,
			Config:         ec2ConfigJSON,
		}
		if err := resourceRepo.Create(ctx, ec2); err != nil {
			return fmt.Errorf("failed to create EC2: %w", err)
		}
		if err := resourceRepo.CreateContainment(ctx, publicSubnet1.ID, ec2.ID); err != nil {
			return fmt.Errorf("failed to create containment: %w", err)
		}

		fmt.Printf("    ✓ Created %d resources for Basic Web Application\n", 6)
	}

	// Scenario 2: High Availability Architecture
	if len(users) > 1 {
		project2 := &models.Project{
			UserID:        users[1].ID,
			InfraToolID:   terraformTarget.ID,
			Name:          "High Availability Architecture",
			CloudProvider: "aws",
			Region:        "us-east-1",
		}
		if err := projectRepo.Create(ctx, project2); err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
		fmt.Printf("  ✓ Created project: %s\n", project2.Name)

		// Create VPC
		vpc2ConfigJSON, _ := json.Marshal(map[string]interface{}{"cidr": "10.0.0.0/16"})
		vpc2 := &models.Resource{
			ProjectID:      project2.ID,
			ResourceTypeID: vpcType.ID,
			Name:           "ha-vpc",
			PosX:           100,
			PosY:           100,
			Config:         vpc2ConfigJSON,
		}
		if err := resourceRepo.Create(ctx, vpc2); err != nil {
			return fmt.Errorf("failed to create VPC: %w", err)
		}

		// Create Load Balancer
		lbConfigJSON, _ := json.Marshal(map[string]interface{}{"type": "application", "scheme": "internet-facing"})
		lb := &models.Resource{
			ProjectID:      project2.ID,
			ResourceTypeID: lbType.ID,
			Name:           "ha-alb",
			PosX:           50,
			PosY:           100,
			Config:         lbConfigJSON,
		}
		if err := resourceRepo.Create(ctx, lb); err != nil {
			return fmt.Errorf("failed to create load balancer: %w", err)
		}

		// Create Auto Scaling Group
		asgConfigJSON, _ := json.Marshal(map[string]interface{}{"min_size": 2, "max_size": 6, "desired_capacity": 3, "instance_type": "t3.small"})
		asg := &models.Resource{
			ProjectID:      project2.ID,
			ResourceTypeID: asgType.ID,
			Name:           "ha-asg",
			PosX:           200,
			PosY:           200,
			Config:         asgConfigJSON,
		}
		if err := resourceRepo.Create(ctx, asg); err != nil {
			return fmt.Errorf("failed to create ASG: %w", err)
		}

		// Create NAT Gateway
		natConfigJSON, _ := json.Marshal(map[string]interface{}{})
		nat := &models.Resource{
			ProjectID:      project2.ID,
			ResourceTypeID: natType.ID,
			Name:           "ha-nat-gateway",
			PosX:           150,
			PosY:           150,
			Config:         natConfigJSON,
		}
		if err := resourceRepo.Create(ctx, nat); err != nil {
			return fmt.Errorf("failed to create NAT gateway: %w", err)
		}

		fmt.Printf("    ✓ Created %d resources for High Availability Architecture\n", 4)
	}

	// Scenario 3: Lambda + S3 (Serverless)
	if len(users) > 2 {
		var lambdaType, s3Type models.ResourceType
		db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "Lambda", "aws").First(&lambdaType)
		db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "S3", "aws").First(&s3Type)

		project3 := &models.Project{
			UserID:        users[2].ID,
			InfraToolID:   terraformTarget.ID,
			Name:          "Serverless Lambda + S3",
			CloudProvider: "aws",
			Region:        "us-east-1",
		}
		if err := projectRepo.Create(ctx, project3); err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
		fmt.Printf("  ✓ Created project: %s\n", project3.Name)

		// Create Lambda Function
		lambdaConfigJSON, _ := json.Marshal(map[string]interface{}{"runtime": "python3.9", "handler": "index.handler", "memory_size_mb": 128})
		lambda := &models.Resource{
			ProjectID:      project3.ID,
			ResourceTypeID: lambdaType.ID,
			Name:           "data-processor",
			PosX:           100,
			PosY:           100,
			Config:         lambdaConfigJSON,
		}
		if err := resourceRepo.Create(ctx, lambda); err != nil {
			return fmt.Errorf("failed to create Lambda: %w", err)
		}

		// Create S3 Bucket
		s3ConfigJSON, _ := json.Marshal(map[string]interface{}{"bucket_name": "my-data-bucket", "versioning": true})
		s3 := &models.Resource{
			ProjectID:      project3.ID,
			ResourceTypeID: s3Type.ID,
			Name:           "data-bucket",
			PosX:           200,
			PosY:           100,
			Config:         s3ConfigJSON,
		}
		if err := resourceRepo.Create(ctx, s3); err != nil {
			return fmt.Errorf("failed to create S3: %w", err)
		}

		// Create dependency: Lambda uses S3
		var usesDepType models.DependencyType
		db.WithContext(ctx).Where("name = ?", "uses").First(&usesDepType)
		dependency := &models.ResourceDependency{
			FromResourceID:   lambda.ID,
			ToResourceID:     s3.ID,
			DependencyTypeID: usesDepType.ID,
		}
		if err := resourceRepo.CreateDependency(ctx, dependency); err != nil {
			return fmt.Errorf("failed to create dependency: %w", err)
		}

		fmt.Printf("    ✓ Created %d resources for Serverless Lambda + S3\n", 2)
	}

	return nil
}

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}
