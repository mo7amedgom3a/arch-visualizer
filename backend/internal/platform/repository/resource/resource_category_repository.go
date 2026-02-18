package resourcerepo

import (
"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"context"

	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ResourceCategoryRepository provides operations for resource category lookups.
type ResourceCategoryRepository struct {
	*repository.BaseRepository
}

// NewResourceCategoryRepository creates a new resource category repository.
func NewResourceCategoryRepository() (*ResourceCategoryRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &ResourceCategoryRepository{BaseRepository: base}, nil
}

// Create creates a new resource category.
func (r *ResourceCategoryRepository) Create(ctx context.Context, category *models.ResourceCategory) error {
	return r.GetDB(ctx).Create(category).Error
}

// FindByID finds a resource category by ID.
func (r *ResourceCategoryRepository) FindByID(ctx context.Context, id uint) (*models.ResourceCategory, error) {
	var category models.ResourceCategory
	err := r.GetDB(ctx).First(&category, id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "resource_category", "ResourceCategoryRepository.FindByID")
	}
	return &category, nil
}

// FindByName finds a resource category by its unique name.
func (r *ResourceCategoryRepository) FindByName(ctx context.Context, name string) (*models.ResourceCategory, error) {
	var category models.ResourceCategory
	err := r.GetDB(ctx).
		Where("name = ?", name).
		First(&category).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "resource_category", "ResourceCategoryRepository.FindByName")
	}
	return &category, nil
}

// List lists resource categories with optional pagination.
func (r *ResourceCategoryRepository) List(ctx context.Context, limit, offset int) ([]*models.ResourceCategory, error) {
	var categories []*models.ResourceCategory
	db := r.GetDB(ctx)
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err := db.Find(&categories).Error
	return categories, err
}

