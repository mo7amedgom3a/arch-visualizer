package projectrepo

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"gorm.io/gorm"

	"log/slog"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ProjectRepository defines operations for project management
type ProjectRepository struct {
	*repository.BaseRepository
	logger *slog.Logger
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(logger *slog.Logger) (*ProjectRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &ProjectRepository{BaseRepository: base, logger: logger}, nil
}

// NewProjectRepositoryWithDB creates a new project repository with a custom DB (detailed loggers not supported)
func NewProjectRepositoryWithDB(db *gorm.DB, logger *slog.Logger) *ProjectRepository {
	return &ProjectRepository{
		BaseRepository: repository.NewBaseRepositoryWithDB(db),
		logger:         logger,
	}
}

// Create creates a new project
func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	r.logger.Info("Creating project", "project_id", project.ID)
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
		return nil, platformerrors.HandleGormError(err, "project", "ProjectRepository.FindByID")
	}
	return &project, nil
}

// FindByUserID finds all projects for a user
func (r *ProjectRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.GetDB(ctx).Where("user_id = ?", userID).Find(&projects).Error
	return projects, err
}

// FindByRootProjectID returns all project snapshot rows that share the given root project ID.
// This includes the root itself (root_project_id IS NULL but id = rootProjectID) and all
// subsequent versions (root_project_id = rootProjectID).
func (r *ProjectRepository) FindByRootProjectID(ctx context.Context, rootProjectID uuid.UUID) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.GetDB(ctx).
		Where("root_project_id = ? OR id = ?", rootProjectID, rootProjectID).
		Order("created_at asc").
		Find(&projects).Error
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

// FindAll finds all projects with pagination, sorting, and searching
func (r *ProjectRepository) FindAll(ctx context.Context, userID uuid.UUID, page, limit int, sort, order, search string) ([]*models.Project, int64, error) {
	var projects []*models.Project
	var total int64

	db := r.GetDB(ctx).Model(&models.Project{})

	// Filter by user
	if userID != uuid.Nil {
		db = db.Where("user_id = ?", userID)
	}

	// Search
	if search != "" {
		searchParam := "%" + search + "%"
		db = db.Where("name ILIKE ? OR description ILIKE ?", searchParam, searchParam)
	}

	// Count total before pagination
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	if sort != "" {
		if order == "" {
			order = "asc"
		}
		db = db.Order(sort + " " + order)
	} else {
		db = db.Order("created_at desc")
	}

	// Pagination
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	err := db.
		Preload("User").
		Preload("IACTarget").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error

	return projects, total, err
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
