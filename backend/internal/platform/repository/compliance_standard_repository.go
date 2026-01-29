package repository

import (
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ComplianceStandardRepository provides operations for compliance standards.
type ComplianceStandardRepository struct {
	*BaseRepository
}

// NewComplianceStandardRepository creates a new compliance standard repository.
func NewComplianceStandardRepository() (*ComplianceStandardRepository, error) {
	base, err := NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &ComplianceStandardRepository{BaseRepository: base}, nil
}

// Create creates a new compliance standard.
func (r *ComplianceStandardRepository) Create(ctx context.Context, standard *models.ComplianceStandard) error {
	return r.GetDB(ctx).Create(standard).Error
}

// FindByID finds a compliance standard by ID.
func (r *ComplianceStandardRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.ComplianceStandard, error) {
	var standard models.ComplianceStandard
	err := r.GetDB(ctx).First(&standard, "id = ?", id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "compliance_standard", "ComplianceStandardRepository.FindByID")
	}
	return &standard, nil
}

// FindBySlug finds a compliance standard by its unique slug.
func (r *ComplianceStandardRepository) FindBySlug(ctx context.Context, slug string) (*models.ComplianceStandard, error) {
	var standard models.ComplianceStandard
	err := r.GetDB(ctx).
		Where("slug = ?", slug).
		First(&standard).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "compliance_standard", "ComplianceStandardRepository.FindBySlug")
	}
	return &standard, nil
}

// List lists compliance standards with optional pagination.
func (r *ComplianceStandardRepository) List(ctx context.Context, limit, offset int) ([]*models.ComplianceStandard, error) {
	var standards []*models.ComplianceStandard
	db := r.GetDB(ctx)
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err := db.Order("name ASC").Find(&standards).Error
	return standards, err
}

