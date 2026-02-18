package infrastructurerepo

import (
"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"context"

	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// IACTargetRepository defines operations for IaC target lookups (Terraform, Pulumi, CDK, etc.).
type IACTargetRepository struct {
	*repository.BaseRepository
}

// NewIACTargetRepository creates a new IaC target repository.
func NewIACTargetRepository() (*IACTargetRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &IACTargetRepository{BaseRepository: base}, nil
}

// Create creates a new IaC target.
func (r *IACTargetRepository) Create(ctx context.Context, target *models.IACTarget) error {
	return r.GetDB(ctx).Create(target).Error
}

// FindByID finds an IaC target by its ID.
func (r *IACTargetRepository) FindByID(ctx context.Context, id uint) (*models.IACTarget, error) {
	var target models.IACTarget
	err := r.GetDB(ctx).First(&target, id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "iac_target", "IACTargetRepository.FindByID")
	}
	return &target, nil
}

// FindByName finds an IaC target by its unique name.
func (r *IACTargetRepository) FindByName(ctx context.Context, name string) (*models.IACTarget, error) {
	var target models.IACTarget
	err := r.GetDB(ctx).
		Where("name = ?", name).
		First(&target).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "iac_target", "IACTargetRepository.FindByName")
	}
	return &target, nil
}

// List lists IaC targets with optional pagination.
func (r *IACTargetRepository) List(ctx context.Context, limit, offset int) ([]*models.IACTarget, error) {
	var targets []*models.IACTarget
	db := r.GetDB(ctx)
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err := db.Find(&targets).Error
	return targets, err
}

