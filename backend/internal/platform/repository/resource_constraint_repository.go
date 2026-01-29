package repository

import (
	"context"

	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ResourceConstraintRepository provides operations for resource constraints.
type ResourceConstraintRepository struct {
	*BaseRepository
}

// NewResourceConstraintRepository creates a new resource constraint repository.
func NewResourceConstraintRepository() (*ResourceConstraintRepository, error) {
	base, err := NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &ResourceConstraintRepository{BaseRepository: base}, nil
}

// Create creates a new resource constraint.
func (r *ResourceConstraintRepository) Create(ctx context.Context, constraint *models.ResourceConstraint) error {
	return r.GetDB(ctx).Create(constraint).Error
}

// FindByID finds a resource constraint by ID.
func (r *ResourceConstraintRepository) FindByID(ctx context.Context, id uint) (*models.ResourceConstraint, error) {
	var constraint models.ResourceConstraint
	err := r.GetDB(ctx).First(&constraint, id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "resource_constraint", "ResourceConstraintRepository.FindByID")
	}
	return &constraint, nil
}

// FindByResourceType lists all constraints for a given resource type.
func (r *ResourceConstraintRepository) FindByResourceType(ctx context.Context, resourceTypeID uint) ([]*models.ResourceConstraint, error) {
	var constraints []*models.ResourceConstraint
	err := r.GetDB(ctx).
		Where("resource_type_id = ?", resourceTypeID).
		Find(&constraints).Error
	return constraints, err
}

// Delete deletes a constraint by ID.
func (r *ResourceConstraintRepository) Delete(ctx context.Context, id uint) error {
	return r.GetDB(ctx).Delete(&models.ResourceConstraint{}, id).Error
}

