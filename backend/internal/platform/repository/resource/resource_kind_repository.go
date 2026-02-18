package resourcerepo

import (
"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"context"

	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ResourceKindRepository provides operations for resource kind lookups.
type ResourceKindRepository struct {
	*repository.BaseRepository
}

// NewResourceKindRepository creates a new resource kind repository.
func NewResourceKindRepository() (*ResourceKindRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &ResourceKindRepository{BaseRepository: base}, nil
}

// Create creates a new resource kind.
func (r *ResourceKindRepository) Create(ctx context.Context, kind *models.ResourceKind) error {
	return r.GetDB(ctx).Create(kind).Error
}

// FindByID finds a resource kind by ID.
func (r *ResourceKindRepository) FindByID(ctx context.Context, id uint) (*models.ResourceKind, error) {
	var kind models.ResourceKind
	err := r.GetDB(ctx).First(&kind, id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "resource_kind", "ResourceKindRepository.FindByID")
	}
	return &kind, nil
}

// FindByName finds a resource kind by its unique name.
func (r *ResourceKindRepository) FindByName(ctx context.Context, name string) (*models.ResourceKind, error) {
	var kind models.ResourceKind
	err := r.GetDB(ctx).
		Where("name = ?", name).
		First(&kind).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "resource_kind", "ResourceKindRepository.FindByName")
	}
	return &kind, nil
}

// List lists resource kinds with optional pagination.
func (r *ResourceKindRepository) List(ctx context.Context, limit, offset int) ([]*models.ResourceKind, error) {
	var kinds []*models.ResourceKind
	db := r.GetDB(ctx)
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err := db.Find(&kinds).Error
	return kinds, err
}

