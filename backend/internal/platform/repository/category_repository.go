package repository

import (
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// CategoryRepository provides operations for template marketplace categories.
type CategoryRepository struct {
	*BaseRepository
}

// NewCategoryRepository creates a new category repository.
func NewCategoryRepository() (*CategoryRepository, error) {
	base, err := NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &CategoryRepository{BaseRepository: base}, nil
}

// Create creates a new category.
func (r *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	return r.GetDB(ctx).Create(category).Error
}

// FindByID finds a category by ID.
func (r *CategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.GetDB(ctx).First(&category, "id = ?", id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "category", "CategoryRepository.FindByID")
	}
	return &category, nil
}

// FindBySlug finds a category by its unique slug.
func (r *CategoryRepository) FindBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var category models.Category
	err := r.GetDB(ctx).
		Where("slug = ?", slug).
		First(&category).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "category", "CategoryRepository.FindBySlug")
	}
	return &category, nil
}

// List lists categories with optional pagination.
func (r *CategoryRepository) List(ctx context.Context, limit, offset int) ([]*models.Category, error) {
	var categories []*models.Category
	db := r.GetDB(ctx)
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err := db.Order("name ASC").Find(&categories).Error
	return categories, err
}

