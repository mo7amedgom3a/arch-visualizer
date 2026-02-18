package templaterepo

import (
"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// TemplateRepository provides operations for marketplace templates.
type TemplateRepository struct {
	*repository.BaseRepository
}

// NewTemplateRepository creates a new template repository.
func NewTemplateRepository() (*TemplateRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &TemplateRepository{BaseRepository: base}, nil
}

// Create creates a new template.
func (r *TemplateRepository) Create(ctx context.Context, template *models.Template) error {
	return r.GetDB(ctx).Create(template).Error
}

// FindByID finds a template by ID with related data preloaded.
func (r *TemplateRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	var template models.Template
	err := r.GetDB(ctx).
		Preload("Category").
		Preload("Author").
		Preload("Technologies").
		Preload("IACFormats").
		Preload("ComplianceStandards").
		Preload("UseCases").
		Preload("Features").
		Preload("Components").
		Preload("Reviews").
		First(&template, "id = ?", id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "template", "TemplateRepository.FindByID")
	}
	return &template, nil
}

// FindByCategory lists templates for a given category.
func (r *TemplateRepository) FindByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*models.Template, error) {
	var templates []*models.Template
	db := r.GetDB(ctx).
		Where("category_id = ?", categoryID)
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err := db.Order("created_at DESC").Find(&templates).Error
	return templates, err
}

// FindByAuthor lists templates created by a specific author.
func (r *TemplateRepository) FindByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int) ([]*models.Template, error) {
	var templates []*models.Template
	db := r.GetDB(ctx).
		Where("author_id = ?", authorID)
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err := db.Order("created_at DESC").Find(&templates).Error
	return templates, err
}

// ListPopular lists popular templates ordered by popularity flags and metrics.
func (r *TemplateRepository) ListPopular(ctx context.Context, limit, offset int) ([]*models.Template, error) {
	var templates []*models.Template
	db := r.GetDB(ctx).
		Order("is_popular DESC").
		Order("downloads DESC").
		Order("rating DESC")
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err := db.Find(&templates).Error
	return templates, err
}

