package seeder

import (
	"context"
	"log"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/rules"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ECSResourceType defines an ECS resource type for seeding
type ECSResourceType struct {
	Name       string
	Category   string
	Kind       string
	IsRegional bool
	IsGlobal   bool
}

// ECSPricingRate defines pricing rates for ECS resources
type ECSPricingRate struct {
	ResourceType  string
	ComponentName string
	PricingModel  string
	Unit          string
	Rate          float64
	Region        string
}

// SeedECSData seeds all ECS-related data to the database
func SeedECSData(ctx context.Context) error {
	log.Println("Starting ECS data seeding...")

	// Seed categories first
	if err := seedECSCategories(ctx); err != nil {
		return err
	}

	// Seed kinds
	if err := seedECSKinds(ctx); err != nil {
		return err
	}

	// Seed resource types
	if err := seedECSResourceTypes(ctx); err != nil {
		return err
	}

	// Seed pricing rates
	if err := seedECSPricingRates(ctx); err != nil {
		return err
	}

	// Seed constraints (reuse existing constraint seeder logic)
	if err := seedECSConstraints(ctx); err != nil {
		return err
	}

	log.Println("ECS data seeding completed successfully!")
	return nil
}

func seedECSCategories(ctx context.Context) error {
	log.Println("Seeding ECS categories...")

	db := database.DB
	categories := []string{"Containers"}

	for _, name := range categories {
		var existing models.ResourceCategory
		if err := db.Where("name = ?", name).First(&existing).Error; err == nil {
			log.Printf("Category '%s' already exists, skipping", name)
			continue
		}

		category := &models.ResourceCategory{Name: name}
		if err := db.Create(category).Error; err != nil {
			log.Printf("Error creating category '%s': %v", name, err)
			return err
		}
		log.Printf("Created category: %s", name)
	}

	return nil
}

func seedECSKinds(ctx context.Context) error {
	log.Println("Seeding ECS kinds...")

	db := database.DB
	kinds := []string{"Container", "ContainerCluster", "ContainerService", "TaskDefinition", "CapacityProvider"}

	for _, name := range kinds {
		var existing models.ResourceKind
		if err := db.Where("name = ?", name).First(&existing).Error; err == nil {
			log.Printf("Kind '%s' already exists, skipping", name)
			continue
		}

		kind := &models.ResourceKind{Name: name}
		if err := db.Create(kind).Error; err != nil {
			log.Printf("Error creating kind '%s': %v", name, err)
			return err
		}
		log.Printf("Created kind: %s", name)
	}

	return nil
}

func seedECSResourceTypes(ctx context.Context) error {
	log.Println("Seeding ECS resource types...")

	db := database.DB

	resourceTypes := []ECSResourceType{
		{Name: "ECSCluster", Category: "Containers", Kind: "ContainerCluster", IsRegional: true, IsGlobal: false},
		{Name: "ECSTaskDefinition", Category: "Containers", Kind: "TaskDefinition", IsRegional: true, IsGlobal: false},
		{Name: "ECSService", Category: "Containers", Kind: "ContainerService", IsRegional: true, IsGlobal: false},
		{Name: "ECSCapacityProvider", Category: "Containers", Kind: "CapacityProvider", IsRegional: true, IsGlobal: false},
		{Name: "ECSClusterCapacityProviders", Category: "Containers", Kind: "Container", IsRegional: true, IsGlobal: false},
	}

	for _, rt := range resourceTypes {
		// Check if exists
		var existing models.ResourceType
		if err := db.Where("name = ? AND cloud_provider = ?", rt.Name, "aws").First(&existing).Error; err == nil {
			log.Printf("Resource type '%s' already exists, skipping", rt.Name)
			continue
		}

		// Get category ID
		var category models.ResourceCategory
		if err := db.Where("name = ?", rt.Category).First(&category).Error; err != nil {
			log.Printf("Warning: Category '%s' not found for resource type '%s'", rt.Category, rt.Name)
			continue
		}

		// Get kind ID
		var kind models.ResourceKind
		if err := db.Where("name = ?", rt.Kind).First(&kind).Error; err != nil {
			log.Printf("Warning: Kind '%s' not found for resource type '%s'", rt.Kind, rt.Name)
			continue
		}

		newResourceType := &models.ResourceType{
			Name:          rt.Name,
			CloudProvider: "aws",
			CategoryID:    &category.ID,
			KindID:        &kind.ID,
			IsRegional:    rt.IsRegional,
			IsGlobal:      rt.IsGlobal,
		}

		if err := db.Create(newResourceType).Error; err != nil {
			log.Printf("Error creating resource type '%s': %v", rt.Name, err)
			return err
		}
		log.Printf("Created resource type: %s", rt.Name)
	}

	return nil
}

func seedECSPricingRates(ctx context.Context) error {
	log.Println("Seeding ECS pricing rates...")

	db := database.DB

	// Fargate pricing for us-east-1
	region := "us-east-1"
	pricingRates := []ECSPricingRate{
		// Fargate vCPU pricing
		{ResourceType: "ECSService", ComponentName: "Fargate vCPU", PricingModel: "per_hour", Unit: "vCPU-hour", Rate: 0.04048, Region: region},
		{ResourceType: "ECSService", ComponentName: "Fargate Memory", PricingModel: "per_hour", Unit: "GB-hour", Rate: 0.004445, Region: region},

		// Fargate Spot pricing (approx 70% discount)
		{ResourceType: "ECSService", ComponentName: "Fargate Spot vCPU", PricingModel: "per_hour", Unit: "vCPU-hour", Rate: 0.01214, Region: region},
		{ResourceType: "ECSService", ComponentName: "Fargate Spot Memory", PricingModel: "per_hour", Unit: "GB-hour", Rate: 0.00133, Region: region},

		// ECS Cluster (Container Insights pricing)
		{ResourceType: "ECSCluster", ComponentName: "Container Insights", PricingModel: "per_month", Unit: "container-month", Rate: 0.50, Region: region},
	}

	for _, pr := range pricingRates {
		// Check if pricing rate exists
		var existing models.PricingRate
		if err := db.Where("resource_type = ? AND component_name = ? AND region = ?",
			pr.ResourceType, pr.ComponentName, pr.Region).First(&existing).Error; err == nil {
			log.Printf("Pricing rate '%s / %s' already exists, skipping", pr.ResourceType, pr.ComponentName)
			continue
		}

		newRate := &models.PricingRate{
			Provider:      "aws",
			ResourceType:  pr.ResourceType,
			ComponentName: pr.ComponentName,
			PricingModel:  pr.PricingModel,
			Unit:          pr.Unit,
			Rate:          pr.Rate,
			Currency:      "USD",
			Region:        &pr.Region,
			EffectiveFrom: time.Now(),
		}

		if err := db.Create(newRate).Error; err != nil {
			log.Printf("Warning: pricing rate '%s / %s' creation skipped: %v", pr.ResourceType, pr.ComponentName, err)
			continue
		}
		log.Printf("Created pricing rate: %s / %s @ $%.5f %s", pr.ResourceType, pr.ComponentName, pr.Rate, pr.Unit)
	}

	return nil
}

func seedECSConstraints(ctx context.Context) error {
	log.Println("Seeding ECS constraints...")

	db := database.DB

	// Get ECS rules from the defaults
	ecsRules := rules.DefaultContainerRules()

	created := 0
	skipped := 0
	failed := 0

	for _, rule := range ecsRules {
		// Find resource type
		var resourceType models.ResourceType
		if err := db.Where("name = ? AND cloud_provider = ?", rule.ResourceType, "aws").First(&resourceType).Error; err != nil {
			log.Printf("Warning: Resource type '%s' not found for constraint seeding", rule.ResourceType)
			failed++
			continue
		}

		// Check if constraint exists
		var existing models.ResourceConstraint
		if err := db.Where("resource_type_id = ? AND constraint_type = ? AND constraint_value = ?",
			resourceType.ID, rule.ConstraintType, rule.ConstraintValue).First(&existing).Error; err == nil {
			skipped++
			continue
		}

		// Create new constraint
		constraint := &models.ResourceConstraint{
			ResourceTypeID:  resourceType.ID,
			ConstraintType:  rule.ConstraintType,
			ConstraintValue: rule.ConstraintValue,
		}

		if err := db.Create(constraint).Error; err != nil {
			log.Printf("Error creating constraint for %s: %v", rule.ResourceType, err)
			failed++
			continue
		}
		created++
	}

	log.Printf("ECS constraints: %d created, %d skipped, %d failed", created, skipped, failed)
	return nil
}
