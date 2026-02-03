package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"gorm.io/gorm"
)

// ProjectOutputRepository handles database operations for project outputs
type ProjectOutputRepository struct {
	BaseRepository
}

// NewProjectOutputRepository creates a new project output repository
func NewProjectOutputRepository() (*ProjectOutputRepository, error) {
	base, err := NewBaseRepository()
	if err != nil {
		return nil, err
	}
	return &ProjectOutputRepository{BaseRepository: *base}, nil
}

// NewProjectOutputRepositoryWithDB creates a new project output repository with an injected *gorm.DB
func NewProjectOutputRepositoryWithDB(db *gorm.DB) *ProjectOutputRepository {
	return &ProjectOutputRepository{BaseRepository: BaseRepository{db: db}}
}

// Create creates a new project output
func (r *ProjectOutputRepository) Create(ctx context.Context, output *models.ProjectOutput) error {
	return r.GetDB(ctx).Create(output).Error
}

// FindByProjectID retrieves all outputs for a project
func (r *ProjectOutputRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectOutput, error) {
	var outputs []*models.ProjectOutput
	err := r.GetDB(ctx).Where("project_id = ?", projectID).Find(&outputs).Error
	return outputs, err
}

// DeleteByProjectID deletes all outputs for a project
func (r *ProjectOutputRepository) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error {
	return r.GetDB(ctx).Where("project_id = ?", projectID).Delete(&models.ProjectOutput{}).Error
}
