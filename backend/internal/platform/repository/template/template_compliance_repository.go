package templaterepo

import (
"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// TemplateComplianceRepository provides operations for template-compliance associations.
type TemplateComplianceRepository struct {
	*repository.BaseRepository
}

// NewTemplateComplianceRepository creates a new template-compliance repository.
func NewTemplateComplianceRepository() (*TemplateComplianceRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &TemplateComplianceRepository{BaseRepository: base}, nil
}

// Create creates a new template-compliance association.
func (r *TemplateComplianceRepository) Create(ctx context.Context, tc *models.TemplateCompliance) error {
	return r.GetDB(ctx).Create(tc).Error
}

// Delete removes a template-compliance association.
func (r *TemplateComplianceRepository) Delete(ctx context.Context, templateID, complianceID uuid.UUID) error {
	return r.GetDB(ctx).
		Where("template_id = ? AND compliance_id = ?", templateID, complianceID).
		Delete(&models.TemplateCompliance{}).Error
}

// FindByTemplate lists compliance associations for a given template.
func (r *TemplateComplianceRepository) FindByTemplate(ctx context.Context, templateID uuid.UUID) ([]*models.TemplateCompliance, error) {
	var items []*models.TemplateCompliance
	err := r.GetDB(ctx).
		Where("template_id = ?", templateID).
		Preload("Compliance").
		Find(&items).Error
	return items, err
}

