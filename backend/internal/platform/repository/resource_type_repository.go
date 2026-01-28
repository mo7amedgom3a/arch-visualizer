package repository

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
)

// ResourceTypeRepository provides operations for resource type lookups
type ResourceTypeRepository struct {
	*BaseRepository
}

// NewResourceTypeRepository creates a new resource type repository
func NewResourceTypeRepository() (*ResourceTypeRepository, error) {
	base, err := NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &ResourceTypeRepository{BaseRepository: base}, nil
}

// FindByNameAndProvider finds a resource type by name and cloud provider
func (r *ResourceTypeRepository) FindByNameAndProvider(ctx context.Context, name, provider string) (*models.ResourceType, error) {
	var resourceType models.ResourceType
	err := r.GetDB(ctx).
		Where("name = ? AND cloud_provider = ?", name, provider).
		First(&resourceType).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "resource_type", "ResourceTypeRepository.FindByNameAndProvider")
	}
	return &resourceType, nil
}

// FindByID finds a resource type by ID
func (r *ResourceTypeRepository) FindByID(ctx context.Context, id uint) (*models.ResourceType, error) {
	var resourceType models.ResourceType
	err := r.GetDB(ctx).First(&resourceType, id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "resource_type", "ResourceTypeRepository.FindByID")
	}
	return &resourceType, nil
}

// ListByProvider lists all resource types for a given provider
func (r *ResourceTypeRepository) ListByProvider(ctx context.Context, provider string) ([]*models.ResourceType, error) {
	var resourceTypes []*models.ResourceType
	err := r.GetDB(ctx).
		Where("cloud_provider = ?", provider).
		Find(&resourceTypes).Error
	return resourceTypes, err
}
