package repository_test

import (
	"context"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	infrastructurerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/infrastructure"
	resourcerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/resource"
)

func TestLookupRepositories_BasicOperations(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(
		t,
		&models.IACTarget{},
		&models.ResourceCategory{},
		&models.ResourceKind{},
		&models.DependencyType{},
		&models.ResourceType{},
		&models.ResourceConstraint{},
	)
	base := repository.NewBaseRepositoryWithDB(db)

	// IACTargetRepository
	iaRepo := &infrastructurerepo.IACTargetRepository{BaseRepository: base}
	target := &models.IACTarget{Name: "terraform"}
	if err := iaRepo.Create(ctx, target); err != nil {
		t.Fatalf("IACTargetRepository.Create error: %v", err)
	}
	if _, err := iaRepo.FindByName(ctx, "terraform"); err != nil {
		t.Fatalf("IACTargetRepository.FindByName error: %v", err)
	}

	// ResourceCategoryRepository
	rcRepo := &resourcerepo.ResourceCategoryRepository{BaseRepository: base}
	cat := &models.ResourceCategory{Name: "Compute"}
	if err := rcRepo.Create(ctx, cat); err != nil {
		t.Fatalf("ResourceCategoryRepository.Create error: %v", err)
	}
	if _, err := rcRepo.FindByName(ctx, "Compute"); err != nil {
		t.Fatalf("ResourceCategoryRepository.FindByName error: %v", err)
	}

	// ResourceKindRepository
	rkRepo := &resourcerepo.ResourceKindRepository{BaseRepository: base}
	kind := &models.ResourceKind{Name: "VirtualMachine"}
	if err := rkRepo.Create(ctx, kind); err != nil {
		t.Fatalf("ResourceKindRepository.Create error: %v", err)
	}
	if _, err := rkRepo.FindByName(ctx, "VirtualMachine"); err != nil {
		t.Fatalf("ResourceKindRepository.FindByName error: %v", err)
	}

	// DependencyTypeRepository
	dtRepo := &resourcerepo.DependencyTypeRepository{BaseRepository: base}
	depType := &models.DependencyType{Name: "uses"}
	if err := db.Create(depType).Error; err != nil {
		t.Fatalf("failed to create dependency type directly: %v", err)
	}
	if _, err := dtRepo.FindByName(ctx, "uses"); err != nil {
		t.Fatalf("DependencyTypeRepository.FindByName error: %v", err)
	}

	// ResourceTypeRepository
	rtRepo := &resourcerepo.ResourceTypeRepository{BaseRepository: base}
	rt := &models.ResourceType{Name: "S3", CloudProvider: "aws"}
	if err := db.Create(rt).Error; err != nil {
		t.Fatalf("failed to create resource type directly: %v", err)
	}
	if _, err := rtRepo.FindByNameAndProvider(ctx, "S3", "aws"); err != nil {
		t.Fatalf("ResourceTypeRepository.FindByNameAndProvider error: %v", err)
	}

	// ResourceConstraintRepository
	rcnRepo := &resourcerepo.ResourceConstraintRepository{BaseRepository: base}
	constraint := &models.ResourceConstraint{
		ResourceTypeID:  rt.ID,
		ConstraintType:  "region",
		ConstraintValue: "us-east-1",
	}
	if err := rcnRepo.Create(ctx, constraint); err != nil {
		t.Fatalf("ResourceConstraintRepository.Create error: %v", err)
	}
	if _, err := rcnRepo.FindByResourceType(ctx, rt.ID); err != nil {
		t.Fatalf("ResourceConstraintRepository.FindByResourceType error: %v", err)
	}
}
