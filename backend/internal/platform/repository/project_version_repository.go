package repository

import (
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ProjectVersionRepository defines operations for project version management
type ProjectVersionRepository struct {
	*BaseRepository
}

// NewProjectVersionRepository creates a new project version repository
func NewProjectVersionRepository() (*ProjectVersionRepository, error) {
	base, err := NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &ProjectVersionRepository{BaseRepository: base}, nil
}

// Create creates a new project version
func (r *ProjectVersionRepository) Create(ctx context.Context, version *models.ProjectVersion) error {
	return r.GetDB(ctx).Create(version).Error
}

// FindByID finds a version by ID
func (r *ProjectVersionRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.ProjectVersion, error) {
	var version models.ProjectVersion
	err := r.GetDB(ctx).First(&version, "id = ?", id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "project_version", "ProjectVersionRepository.FindByID")
	}
	return &version, nil
}

// ListByProjectID lists all versions for a project
func (r *ProjectVersionRepository) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectVersion, error) {
	var versions []*models.ProjectVersion
	err := r.GetDB(ctx).
		Where("project_id = ?", projectID).
		Order("created_at desc").
		Find(&versions).Error
	return versions, err
}

// Delete deletes a project version
func (r *ProjectVersionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.GetDB(ctx).Delete(&models.ProjectVersion{}, "id = ?", id).Error
}
