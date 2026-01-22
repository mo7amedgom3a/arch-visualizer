package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ProjectRepository defines operations for project management
type ProjectRepository struct {
	*BaseRepository
}

// NewProjectRepository creates a new project repository
func NewProjectRepository() (*ProjectRepository, error) {
	base, err := NewBaseRepository()
	if err != nil {
		return nil, err
	}
	return &ProjectRepository{BaseRepository: base}, nil
}

// Create creates a new project
func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	return r.GetDB(ctx).Create(project).Error
}

// FindByID finds a project by ID with related data
func (r *ProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.GetDB(ctx).
		Preload("User").
		Preload("IACTarget").
		Preload("Resources").
		Preload("Resources.ResourceType").
		First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// FindByUserID finds all projects for a user
func (r *ProjectRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.GetDB(ctx).Where("user_id = ?", userID).Find(&projects).Error
	return projects, err
}

// Update updates an existing project
func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	return r.GetDB(ctx).Save(project).Error
}

// Delete deletes a project by ID
func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.GetDB(ctx).Delete(&models.Project{}, "id = ?", id).Error
}

// List lists all projects with pagination
func (r *ProjectRepository) List(ctx context.Context, limit, offset int) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.GetDB(ctx).
		Preload("User").
		Preload("IACTarget").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error
	return projects, err
}
