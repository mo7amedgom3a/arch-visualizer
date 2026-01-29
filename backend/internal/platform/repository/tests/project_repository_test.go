package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
)

func TestProjectRepository_CreateAndFind(t *testing.T) {
	db := newTestDB(t, &models.User{}, &models.IACTarget{}, &models.Project{})
	base := repository.NewBaseRepositoryWithDB(db)
	repo := &repository.ProjectRepository{BaseRepository: base}

	ctx := context.Background()

	user := &models.User{
		ID:        uuid.New(),
		Name:      "Owner",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	target := &models.IACTarget{
		ID:   1,
		Name: "terraform",
	}
	if err := db.Create(target).Error; err != nil {
		t.Fatalf("failed to create iac target: %v", err)
	}

	project := &models.Project{
		ID:            uuid.New(),
		UserID:        user.ID,
		InfraToolID:   target.ID,
		Name:          "Test Project",
		CloudProvider: "aws",
		Region:        "us-east-1",
		CreatedAt:     time.Now(),
	}
	if err := repo.Create(ctx, project); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	found, err := repo.FindByID(ctx, project.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if found.Name != project.Name {
		t.Fatalf("expected name %q, got %q", project.Name, found.Name)
	}

	byUser, err := repo.FindByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("FindByUserID returned error: %v", err)
	}
	if len(byUser) != 1 {
		t.Fatalf("expected 1 project for user, got %d", len(byUser))
	}
}

