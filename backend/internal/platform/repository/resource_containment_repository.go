package repository

import (
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ResourceContainmentRepository provides operations for resource containment relationships.
type ResourceContainmentRepository struct {
	*BaseRepository
}

// NewResourceContainmentRepository creates a new resource containment repository.
func NewResourceContainmentRepository() (*ResourceContainmentRepository, error) {
	base, err := NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &ResourceContainmentRepository{BaseRepository: base}, nil
}

// Create creates a new containment relationship.
func (r *ResourceContainmentRepository) Create(ctx context.Context, containment *models.ResourceContainment) error {
	return r.GetDB(ctx).Create(containment).Error
}

// Delete removes a containment relationship.
func (r *ResourceContainmentRepository) Delete(ctx context.Context, parentID, childID uuid.UUID) error {
	return r.GetDB(ctx).
		Where("parent_resource_id = ? AND child_resource_id = ?", parentID, childID).
		Delete(&models.ResourceContainment{}).Error
}

// FindChildren lists all child relationships for a given parent resource.
func (r *ResourceContainmentRepository) FindChildren(ctx context.Context, parentID uuid.UUID) ([]*models.ResourceContainment, error) {
	var items []*models.ResourceContainment
	err := r.GetDB(ctx).
		Where("parent_resource_id = ?", parentID).
		Preload("ChildResource").
		Find(&items).Error
	return items, err
}

// FindParents lists all parent relationships for a given child resource.
func (r *ResourceContainmentRepository) FindParents(ctx context.Context, childID uuid.UUID) ([]*models.ResourceContainment, error) {
	var items []*models.ResourceContainment
	err := r.GetDB(ctx).
		Where("child_resource_id = ?", childID).
		Preload("ParentResource").
		Find(&items).Error
	return items, err
}

