package repository_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	resourcerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/resource"
)

func TestResourceRepository_CreateAndFindByProject(t *testing.T) {
	db := newTestDB(t, &models.User{}, &models.Project{}, &models.ResourceType{}, &models.Resource{})
	repo := resourcerepo.NewResourceRepositoryWithDB(db, slog.Default())

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

	project := &models.Project{
		ID:            uuid.New(),
		UserID:        user.ID,
		InfraToolID:   1,
		Name:          "Res Project",
		CloudProvider: "aws",
		Region:        "us-east-1",
		CreatedAt:     time.Now(),
	}
	if err := db.Create(project).Error; err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	rt := &models.ResourceType{
		ID:            1,
		Name:          "EC2",
		CloudProvider: "aws",
	}
	if err := db.Create(rt).Error; err != nil {
		t.Fatalf("failed to create resource type: %v", err)
	}

	resource := &models.Resource{
		ID:             uuid.New(),
		ProjectID:      project.ID,
		ResourceTypeID: rt.ID,
		Name:           "instance-1",
		CreatedAt:      time.Now(),
	}
	if err := repo.Create(ctx, resource); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	list, err := repo.FindByProjectID(ctx, project.ID)
	if err != nil {
		t.Fatalf("FindByProjectID returned error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 resource for project, got %d", len(list))
	}
}
