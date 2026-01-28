package repository

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
)

// DependencyTypeRepository provides operations for dependency type lookups
type DependencyTypeRepository struct {
	*BaseRepository
}

// NewDependencyTypeRepository creates a new dependency type repository
func NewDependencyTypeRepository() (*DependencyTypeRepository, error) {
	base, err := NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &DependencyTypeRepository{BaseRepository: base}, nil
}

// FindByName finds a dependency type by name
func (r *DependencyTypeRepository) FindByName(ctx context.Context, name string) (*models.DependencyType, error) {
	var dependencyType models.DependencyType
	err := r.GetDB(ctx).
		Where("name = ?", name).
		First(&dependencyType).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "dependency_type", "DependencyTypeRepository.FindByName")
	}
	return &dependencyType, nil
}

// FindByID finds a dependency type by ID
func (r *DependencyTypeRepository) FindByID(ctx context.Context, id uint) (*models.DependencyType, error) {
	var dependencyType models.DependencyType
	err := r.GetDB(ctx).First(&dependencyType, id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "dependency_type", "DependencyTypeRepository.FindByID")
	}
	return &dependencyType, nil
}

// ListAll lists all dependency types
func (r *DependencyTypeRepository) ListAll(ctx context.Context) ([]*models.DependencyType, error) {
	var dependencyTypes []*models.DependencyType
	err := r.GetDB(ctx).Find(&dependencyTypes).Error
	return dependencyTypes, err
}
