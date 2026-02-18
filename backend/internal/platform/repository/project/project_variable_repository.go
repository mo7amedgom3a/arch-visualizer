package projectrepo

import (
"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"context"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"gorm.io/gorm"
)

// ProjectVariableRepository handles database operations for project variables
type ProjectVariableRepository struct {
	*repository.BaseRepository
}

// NewProjectVariableRepository creates a new project variable repository
func NewProjectVariableRepository() (*ProjectVariableRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, err
	}
	return &ProjectVariableRepository{BaseRepository: base}, nil
}

// NewProjectVariableRepositoryWithDB creates a new project variable repository with an injected *gorm.DB
func NewProjectVariableRepositoryWithDB(db *gorm.DB) *ProjectVariableRepository {
	return &ProjectVariableRepository{BaseRepository: repository.NewBaseRepositoryWithDB(db)}
}

// Create creates a new project variable
func (r *ProjectVariableRepository) Create(ctx context.Context, variable *models.ProjectVariable) error {
	return r.GetDB(ctx).Create(variable).Error
}

// FindByProjectID retrieves all variables for a project
func (r *ProjectVariableRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectVariable, error) {
	var variables []*models.ProjectVariable
	err := r.GetDB(ctx).Where("project_id = ?", projectID).Find(&variables).Error
	return variables, err
}

// DeleteByProjectID deletes all variables for a project
func (r *ProjectVariableRepository) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error {
	return r.GetDB(ctx).Where("project_id = ?", projectID).Delete(&models.ProjectVariable{}).Error
}
