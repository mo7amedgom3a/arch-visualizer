package resourcerepo

import (
"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"context"

	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// HiddenDependencyRepository provides operations for hidden dependencies
type HiddenDependencyRepository struct {
	*repository.BaseRepository
}

// NewHiddenDependencyRepository creates a new hidden dependency repository
func NewHiddenDependencyRepository() (*HiddenDependencyRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &HiddenDependencyRepository{BaseRepository: base}, nil
}

// Create creates a new hidden dependency
func (r *HiddenDependencyRepository) Create(ctx context.Context, dep *models.HiddenDependency) error {
	return r.GetDB(ctx).Create(dep).Error
}

// FindByID finds a hidden dependency by ID
func (r *HiddenDependencyRepository) FindByID(ctx context.Context, id uint) (*models.HiddenDependency, error) {
	var dep models.HiddenDependency
	err := r.GetDB(ctx).First(&dep, "id = ?", id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "hidden_dependency", "HiddenDependencyRepository.FindByID")
	}
	return &dep, nil
}

// FindByParentResourceType finds all hidden dependencies for a parent resource type
func (r *HiddenDependencyRepository) FindByParentResourceType(ctx context.Context, provider, parentResourceType string) ([]*models.HiddenDependency, error) {
	var deps []*models.HiddenDependency
	err := r.GetDB(ctx).
		Where("provider = ?", provider).
		Where("parent_resource_type = ?", parentResourceType).
		Find(&deps).Error
	return deps, err
}

// FindByProvider finds all hidden dependencies for a provider
func (r *HiddenDependencyRepository) FindByProvider(ctx context.Context, provider string) ([]*models.HiddenDependency, error) {
	var deps []*models.HiddenDependency
	err := r.GetDB(ctx).
		Where("provider = ?", provider).
		Find(&deps).Error
	return deps, err
}

// Update updates a hidden dependency
func (r *HiddenDependencyRepository) Update(ctx context.Context, dep *models.HiddenDependency) error {
	return r.GetDB(ctx).Save(dep).Error
}

// Delete deletes a hidden dependency
func (r *HiddenDependencyRepository) Delete(ctx context.Context, id uint) error {
	return r.GetDB(ctx).Delete(&models.HiddenDependency{}, id).Error
}
