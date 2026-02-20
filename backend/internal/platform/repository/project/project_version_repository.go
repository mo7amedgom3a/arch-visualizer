package projectrepo

import (
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
)

// ProjectVersionRepository defines operations for project version management
type ProjectVersionRepository struct {
	*repository.BaseRepository
}

// NewProjectVersionRepository creates a new project version repository
func NewProjectVersionRepository() (*ProjectVersionRepository, error) {
	base, err := repository.NewBaseRepository()
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

// ListByProjectID lists all version entries whose project_id matches the given project UUID.
// This is a direct lookup — for the full history chain use ListByRootProjectID.
func (r *ProjectVersionRepository) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectVersion, error) {
	var versions []*models.ProjectVersion
	err := r.GetDB(ctx).
		Where("project_id = ?", projectID).
		Order("version_number asc").
		Find(&versions).Error
	return versions, err
}

// ListByRootProjectID returns all version entries across the entire version chain for a root project.
// It finds every project that shares the same root_project_id (or whose id IS the root), then
// returns all version entries for those projects ordered by version_number ASC.
func (r *ProjectVersionRepository) ListByRootProjectID(ctx context.Context, rootProjectID uuid.UUID) ([]*models.ProjectVersion, error) {
	var versions []*models.ProjectVersion
	err := r.GetDB(ctx).
		Joins("JOIN projects ON projects.id = project_versions.project_id").
		Where("projects.root_project_id = ? OR projects.id = ?", rootProjectID, rootProjectID).
		Order("project_versions.version_number asc").
		Find(&versions).Error
	return versions, err
}

// GetLatestVersionForProject returns the most recent version entry for a given project snapshot UUID.
func (r *ProjectVersionRepository) GetLatestVersionForProject(ctx context.Context, projectID uuid.UUID) (*models.ProjectVersion, error) {
	var version models.ProjectVersion
	err := r.GetDB(ctx).
		Where("project_id = ?", projectID).
		Order("version_number desc").
		First(&version).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "project_version", "ProjectVersionRepository.GetLatestVersionForProject")
	}
	return &version, nil
}

// Delete deletes a project version
func (r *ProjectVersionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.GetDB(ctx).Delete(&models.ProjectVersion{}, "id = ?", id).Error
}
