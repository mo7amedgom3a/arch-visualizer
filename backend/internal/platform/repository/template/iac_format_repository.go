package templaterepo

import (
"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// IACFormatRepository provides operations for IaC formats in the marketplace.
type IACFormatRepository struct {
	*repository.BaseRepository
}

// NewIACFormatRepository creates a new IaC format repository.
func NewIACFormatRepository() (*IACFormatRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &IACFormatRepository{BaseRepository: base}, nil
}

// Create creates a new IaC format.
func (r *IACFormatRepository) Create(ctx context.Context, format *models.IACFormat) error {
	return r.GetDB(ctx).Create(format).Error
}

// FindByID finds an IaC format by ID.
func (r *IACFormatRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.IACFormat, error) {
	var format models.IACFormat
	err := r.GetDB(ctx).First(&format, "id = ?", id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "iac_format", "IACFormatRepository.FindByID")
	}
	return &format, nil
}

// FindBySlug finds an IaC format by its unique slug.
func (r *IACFormatRepository) FindBySlug(ctx context.Context, slug string) (*models.IACFormat, error) {
	var format models.IACFormat
	err := r.GetDB(ctx).
		Where("slug = ?", slug).
		First(&format).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "iac_format", "IACFormatRepository.FindBySlug")
	}
	return &format, nil
}

// List lists IaC formats with optional pagination.
func (r *IACFormatRepository) List(ctx context.Context, limit, offset int) ([]*models.IACFormat, error) {
	var formats []*models.IACFormat
	db := r.GetDB(ctx)
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err := db.Order("name ASC").Find(&formats).Error
	return formats, err
}

