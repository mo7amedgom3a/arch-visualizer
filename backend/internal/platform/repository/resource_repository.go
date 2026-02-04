package repository

import (
	"context"

	"log/slog"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ResourceRepository defines operations for resource management
type ResourceRepository struct {
	*BaseRepository
	logger *slog.Logger
}

// NewResourceRepository creates a new resource repository
func NewResourceRepository(logger *slog.Logger) (*ResourceRepository, error) {
	base, err := NewBaseRepository()
	if err != nil {
		return nil, errors.NewDatabaseConnectionFailed(err)
	}
	return &ResourceRepository{BaseRepository: base, logger: logger}, nil
}

// Create creates a new resource
func (r *ResourceRepository) Create(ctx context.Context, resource *models.Resource) error {
	r.logger.Info("Creating resource", "resource_id", resource.ID, "type", resource.ResourceTypeID)
	err := r.GetDB(ctx).Create(resource).Error
	if err != nil {
		return errors.HandleGormError(err, "resource", "ResourceRepository.Create")
	}
	return nil
}

// FindByID finds a resource by ID with related data
func (r *ResourceRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Resource, error) {
	var resource models.Resource
	err := r.GetDB(ctx).
		Preload("Project").
		Preload("ResourceType").
		Preload("ResourceType.Category").
		Preload("ResourceType.Kind").
		Preload("ParentResources").
		Preload("ChildResources").
		First(&resource, "id = ?", id).Error
	if err != nil {
		return nil, errors.HandleGormError(err, "resource", "ResourceRepository.FindByID")
	}
	return &resource, nil
}

// FindByProjectID finds all resources for a project
func (r *ResourceRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Resource, error) {
	var resources []*models.Resource
	err := r.GetDB(ctx).
		Where("project_id = ?", projectID).
		Preload("ResourceType").
		Find(&resources).Error
	return resources, err
}

// FindByProjectIDAndType finds resources by project and type
func (r *ResourceRepository) FindByProjectIDAndType(ctx context.Context, projectID uuid.UUID, resourceTypeID uint) ([]*models.Resource, error) {
	var resources []*models.Resource
	err := r.GetDB(ctx).
		Where("project_id = ? AND resource_type_id = ?", projectID, resourceTypeID).
		Preload("ResourceType").
		Find(&resources).Error
	return resources, err
}

// Update updates an existing resource
func (r *ResourceRepository) Update(ctx context.Context, resource *models.Resource) error {
	err := r.GetDB(ctx).Save(resource).Error
	if err != nil {
		return errors.HandleGormError(err, "resource", "ResourceRepository.Update")
	}
	return nil
}

// Delete deletes a resource by ID
func (r *ResourceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.GetDB(ctx).Delete(&models.Resource{}, "id = ?", id).Error
	if err != nil {
		return errors.HandleGormError(err, "resource", "ResourceRepository.Delete")
	}
	return nil
}

// CreateContainment creates a parent-child relationship
func (r *ResourceRepository) CreateContainment(ctx context.Context, parentID, childID uuid.UUID) error {
	containment := &models.ResourceContainment{
		ParentResourceID: parentID,
		ChildResourceID:  childID,
	}
	return r.GetDB(ctx).Create(containment).Error
}

// DeleteContainment deletes a parent-child relationship
func (r *ResourceRepository) DeleteContainment(ctx context.Context, parentID, childID uuid.UUID) error {
	return r.GetDB(ctx).
		Where("parent_resource_id = ? AND child_resource_id = ?", parentID, childID).
		Delete(&models.ResourceContainment{}).Error
}

// CreateDependency creates a dependency relationship
func (r *ResourceRepository) CreateDependency(ctx context.Context, dependency *models.ResourceDependency) error {
	return r.GetDB(ctx).Create(dependency).Error
}

// DeleteDependency deletes a dependency relationship
func (r *ResourceRepository) DeleteDependency(ctx context.Context, fromID, toID uuid.UUID) error {
	return r.GetDB(ctx).
		Where("from_resource_id = ? AND to_resource_id = ?", fromID, toID).
		Delete(&models.ResourceDependency{}).Error
}
