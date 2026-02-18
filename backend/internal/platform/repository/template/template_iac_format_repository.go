package templaterepo

import (
"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// TemplateIACFormatRepository provides operations for template-IaC format associations.
type TemplateIACFormatRepository struct {
	*repository.BaseRepository
}

// NewTemplateIACFormatRepository creates a new template-IaC format repository.
func NewTemplateIACFormatRepository() (*TemplateIACFormatRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &TemplateIACFormatRepository{BaseRepository: base}, nil
}

// Create creates a new template-IaC format association.
func (r *TemplateIACFormatRepository) Create(ctx context.Context, ti *models.TemplateIACFormat) error {
	return r.GetDB(ctx).Create(ti).Error
}

// Delete removes a template-IaC format association.
func (r *TemplateIACFormatRepository) Delete(ctx context.Context, templateID, formatID uuid.UUID) error {
	return r.GetDB(ctx).
		Where("template_id = ? AND iac_format_id = ?", templateID, formatID).
		Delete(&models.TemplateIACFormat{}).Error
}

// FindByTemplate lists IaC format associations for a given template.
func (r *TemplateIACFormatRepository) FindByTemplate(ctx context.Context, templateID uuid.UUID) ([]*models.TemplateIACFormat, error) {
	var items []*models.TemplateIACFormat
	err := r.GetDB(ctx).
		Where("template_id = ?", templateID).
		Preload("IACFormat").
		Find(&items).Error
	return items, err
}

