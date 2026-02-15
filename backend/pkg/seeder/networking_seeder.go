package seeder

import (
	"context"
	"log"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// NetworkingResourceType defines a Networking resource type for seeding
type NetworkingResourceType struct {
	Name       string
	Category   string
	Kind       string
	IsRegional bool
	IsGlobal   bool
}

// NetworkingPricingRate defines pricing rates for Networking resources
type NetworkingPricingRate struct {
	ResourceType  string
	ComponentName string
	PricingModel  string
	Unit          string
	Rate          float64
	Region        string
}

// SeedNetworkingData seeds all Networking-related data to the database
func SeedNetworkingData(ctx context.Context) error {
	log.Println("Starting Networking data seeding...")

	// Seed categories
	if err := seedNetworkingCategories(ctx); err != nil {
		return err
	}

	// Seed kinds
	if err := seedNetworkingKinds(ctx); err != nil {
		return err
	}

	// Seed resource types
	if err := seedNetworkingResourceTypes(ctx); err != nil {
		return err
	}

	// Seed pricing rates
	if err := seedNetworkingPricingRates(ctx); err != nil {
		return err
	}

	log.Println("Networking data seeding completed successfully!")
	return nil
}

func seedNetworkingCategories(ctx context.Context) error {
	log.Println("Seeding Networking categories...")

	db := database.DB
	categories := []string{"Networking"}

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

func seedNetworkingKinds(ctx context.Context) error {
	log.Println("Seeding Networking kinds...")

	db := database.DB
	// Add other networking kinds if needed, for now just VPCEndpoint
	kinds := []string{"VPCEndpoint"}

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

func seedNetworkingResourceTypes(ctx context.Context) error {
	log.Println("Seeding Networking resource types...")

	db := database.DB

	resourceTypes := []NetworkingResourceType{
		{Name: "VPCEndpoint", Category: "Networking", Kind: "VPCEndpoint", IsRegional: true, IsGlobal: false},
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

func seedNetworkingPricingRates(ctx context.Context) error {
	log.Println("Seeding Networking pricing rates...")

	db := database.DB

	// US East 1 pricing
	region := "us-east-1"
	pricingRates := []NetworkingPricingRate{
		// VPC Endpoint pricing (Header: $0.01/hr per ENI per AZ, Data processing: $0.01/GB)
		// Simplified for now
		{ResourceType: "VPCEndpoint", ComponentName: "Interface Endpoint per ENI", PricingModel: "per_hour", Unit: "ENI-hour", Rate: 0.01, Region: region},
		{ResourceType: "VPCEndpoint", ComponentName: "Interface Endpoint Data Processing", PricingModel: "per_gb", Unit: "GB", Rate: 0.01, Region: region}, // $0.01 per GB
		{ResourceType: "VPCEndpoint", ComponentName: "Gateway Endpoint", PricingModel: "free", Unit: "n/a", Rate: 0.00, Region: region},                    // Gateways are free
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
