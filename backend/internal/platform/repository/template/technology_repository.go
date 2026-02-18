package templaterepo

import (
"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// TechnologyRepository provides operations for marketplace technologies.
type TechnologyRepository struct {
	*repository.BaseRepository
}

// NewTechnologyRepository creates a new technology repository.
func NewTechnologyRepository() (*TechnologyRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &TechnologyRepository{BaseRepository: base}, nil
}

// Create creates a new technology.
func (r *TechnologyRepository) Create(ctx context.Context, tech *models.Technology) error {
	return r.GetDB(ctx).Create(tech).Error
}

// FindByID finds a technology by ID.
func (r *TechnologyRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Technology, error) {
	var tech models.Technology
	err := r.GetDB(ctx).First(&tech, "id = ?", id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "technology", "TechnologyRepository.FindByID")
	}
	return &tech, nil
}

// FindBySlug finds a technology by its unique slug.
func (r *TechnologyRepository) FindBySlug(ctx context.Context, slug string) (*models.Technology, error) {
	var tech models.Technology
	err := r.GetDB(ctx).
		Where("slug = ?", slug).
		First(&tech).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "technology", "TechnologyRepository.FindBySlug")
	}
	return &tech, nil
}

// List lists technologies with optional pagination.
func (r *TechnologyRepository) List(ctx context.Context, limit, offset int) ([]*models.Technology, error) {
	var techs []*models.Technology
	db := r.GetDB(ctx)
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err := db.Order("name ASC").Find(&techs).Error
	return techs, err
}

