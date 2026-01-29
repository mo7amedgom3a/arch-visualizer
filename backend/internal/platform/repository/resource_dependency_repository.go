package repository

import (
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ResourceDependencyRepository provides operations for resource dependency relationships.
type ResourceDependencyRepository struct {
	*BaseRepository
}

// NewResourceDependencyRepository creates a new resource dependency repository.
func NewResourceDependencyRepository() (*ResourceDependencyRepository, error) {
	base, err := NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &ResourceDependencyRepository{BaseRepository: base}, nil
}

// Create creates a new dependency relationship.
func (r *ResourceDependencyRepository) Create(ctx context.Context, dep *models.ResourceDependency) error {
	return r.GetDB(ctx).Create(dep).Error
}

// Delete removes a dependency relationship.
func (r *ResourceDependencyRepository) Delete(ctx context.Context, fromID, toID uuid.UUID) error {
	return r.GetDB(ctx).
		Where("from_resource_id = ? AND to_resource_id = ?", fromID, toID).
		Delete(&models.ResourceDependency{}).Error
}

// FindByFromResource lists dependencies originating from a given resource.
func (r *ResourceDependencyRepository) FindByFromResource(ctx context.Context, fromID uuid.UUID) ([]*models.ResourceDependency, error) {
	var deps []*models.ResourceDependency
	err := r.GetDB(ctx).
		Where("from_resource_id = ?", fromID).
		Preload("ToResource").
		Preload("DependencyType").
		Find(&deps).Error
	return deps, err
}

// FindByToResource lists dependencies targeting a given resource.
func (r *ResourceDependencyRepository) FindByToResource(ctx context.Context, toID uuid.UUID) ([]*models.ResourceDependency, error) {
	var deps []*models.ResourceDependency
	err := r.GetDB(ctx).
		Where("to_resource_id = ?", toID).
		Preload("FromResource").
		Preload("DependencyType").
		Find(&deps).Error
	return deps, err
}

