package templaterepo

import (
"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// TemplateTechnologyRepository provides operations for template-technology associations.
type TemplateTechnologyRepository struct {
	*repository.BaseRepository
}

// NewTemplateTechnologyRepository creates a new template-technology repository.
func NewTemplateTechnologyRepository() (*TemplateTechnologyRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &TemplateTechnologyRepository{BaseRepository: base}, nil
}

// Create creates a new template-technology association.
func (r *TemplateTechnologyRepository) Create(ctx context.Context, tt *models.TemplateTechnology) error {
	return r.GetDB(ctx).Create(tt).Error
}

// Delete removes a template-technology association.
func (r *TemplateTechnologyRepository) Delete(ctx context.Context, templateID, technologyID uuid.UUID) error {
	return r.GetDB(ctx).
		Where("template_id = ? AND technology_id = ?", templateID, technologyID).
		Delete(&models.TemplateTechnology{}).Error
}

// FindByTemplate lists technology associations for a given template.
func (r *TemplateTechnologyRepository) FindByTemplate(ctx context.Context, templateID uuid.UUID) ([]*models.TemplateTechnology, error) {
	var items []*models.TemplateTechnology
	err := r.GetDB(ctx).
		Where("template_id = ?", templateID).
		Preload("Technology").
		Find(&items).Error
	return items, err
}

