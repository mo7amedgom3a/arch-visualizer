package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
)

func TestPricingRepository_CreateAndQuery(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(
		t,
		&models.Project{},
		&models.Resource{},
		&models.ProjectPricing{},
		&models.ServicePricing{},
		&models.ServiceTypePricing{},
		&models.ResourcePricing{},
		&models.PricingComponent{},
		&models.User{},
		&models.IACTarget{},
		&models.ResourceType{},
		&models.ResourceCategory{},
	)
	base := repository.NewBaseRepositoryWithDB(db)
	repo := &repository.PricingRepository{BaseRepository: base}

	user := &models.User{ID: uuid.New(), Name: "Owner", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	target := &models.IACTarget{ID: 1, Name: "terraform"}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if err := db.Create(target).Error; err != nil {
		t.Fatalf("failed to create iac target: %v", err)
	}

	project := &models.Project{
		ID:            uuid.New(),
		UserID:        user.ID,
		InfraToolID:   target.ID,
		Name:          "PricingProj",
		CloudProvider: "aws",
		Region:        "us-east-1",
		CreatedAt:     time.Now(),
	}
	if err := db.Create(project).Error; err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	resource := &models.Resource{
		ID:             uuid.New(),
		ProjectID:      project.ID,
		ResourceTypeID: 1,
		Name:           "r1",
		CreatedAt:      time.Now(),
	}
	if err := db.Create(resource).Error; err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	category := &models.ResourceCategory{Name: "Compute"}
	if err := db.Create(category).Error; err != nil {
		t.Fatalf("failed to create category: %v", err)
	}

	rt := &models.ResourceType{Name: "EC2", CloudProvider: "aws"}
	if err := db.Create(rt).Error; err != nil {
		t.Fatalf("failed to create resource type: %v", err)
	}

	pp := &models.ProjectPricing{
		ProjectID:      project.ID,
		TotalCost:      10,
		Currency:       "USD",
		Period:         "monthly",
		DurationSeconds: 30 * 24 * 3600,
		Provider:       "aws",
		CalculatedAt:   time.Now(),
	}
	if err := repo.CreateProjectPricing(ctx, pp); err != nil {
		t.Fatalf("CreateProjectPricing error: %v", err)
	}

	sp := &models.ServicePricing{
		ProjectID:      project.ID,
		CategoryID:     category.ID,
		TotalCost:      5,
		Currency:       "USD",
		Period:         "monthly",
		DurationSeconds: 30 * 24 * 3600,
		Provider:       "aws",
		CalculatedAt:   time.Now(),
	}
	if err := repo.CreateServicePricing(ctx, sp); err != nil {
		t.Fatalf("CreateServicePricing error: %v", err)
	}

	stp := &models.ServiceTypePricing{
		ProjectID:      project.ID,
		ResourceTypeID: rt.ID,
		TotalCost:      3,
		Currency:       "USD",
		Period:         "monthly",
		DurationSeconds: 30 * 24 * 3600,
		Provider:       "aws",
		CalculatedAt:   time.Now(),
	}
	if err := repo.CreateServiceTypePricing(ctx, stp); err != nil {
		t.Fatalf("CreateServiceTypePricing error: %v", err)
	}

	rp := &models.ResourcePricing{
		ProjectID:      project.ID,
		ResourceID:     resource.ID,
		TotalCost:      2,
		Currency:       "USD",
		Period:         "monthly",
		DurationSeconds: 30 * 24 * 3600,
		Provider:       "aws",
		CalculatedAt:   time.Now(),
	}
	if err := repo.CreateResourcePricing(ctx, rp); err != nil {
		t.Fatalf("CreateResourcePricing error: %v", err)
	}

	pc := &models.PricingComponent{
		ResourcePricingID: rp.ID,
		ComponentName:     "hours",
		Model:             "per_hour",
		Unit:              "h",
		Quantity:          10,
		UnitRate:          0.2,
		Subtotal:          2,
		Currency:          "USD",
	}
	if err := repo.CreatePricingComponent(ctx, pc); err != nil {
		t.Fatalf("CreatePricingComponent error: %v", err)
	}

	if _, err := repo.FindProjectPricingByProjectID(ctx, project.ID); err != nil {
		t.Fatalf("FindProjectPricingByProjectID error: %v", err)
	}
	if _, err := repo.FindServicePricingByProjectID(ctx, project.ID); err != nil {
		t.Fatalf("FindServicePricingByProjectID error: %v", err)
	}
	if _, err := repo.FindServiceTypePricingByProjectID(ctx, project.ID); err != nil {
		t.Fatalf("FindServiceTypePricingByProjectID error: %v", err)
	}
	if _, err := repo.FindResourcePricingByProjectID(ctx, project.ID); err != nil {
		t.Fatalf("FindResourcePricingByProjectID error: %v", err)
	}
	if _, err := repo.FindResourcePricingByResourceID(ctx, resource.ID); err != nil {
		t.Fatalf("FindResourcePricingByResourceID error: %v", err)
	}
}

