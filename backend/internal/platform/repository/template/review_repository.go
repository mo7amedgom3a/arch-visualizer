package templaterepo

import (
"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"gorm.io/gorm"
)

// ReviewRepository provides operations for template reviews.
type ReviewRepository struct {
	*repository.BaseRepository
}

// NewReviewRepository creates a new review repository.
func NewReviewRepository() (*ReviewRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &ReviewRepository{BaseRepository: base}, nil
}

// Create creates a new review.
func (r *ReviewRepository) Create(ctx context.Context, review *models.Review) error {
	return r.GetDB(ctx).Create(review).Error
}

// FindByID finds a review by ID with related data preloaded.
func (r *ReviewRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := r.GetDB(ctx).
		Preload("Template").
		Preload("User").
		First(&review, "id = ?", id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "review", "ReviewRepository.FindByID")
	}
	return &review, nil
}

// FindByTemplate lists reviews for a given template.
func (r *ReviewRepository) FindByTemplate(ctx context.Context, templateID uuid.UUID, limit, offset int) ([]*models.Review, error) {
	var reviews []*models.Review
	db := r.GetDB(ctx).
		Where("template_id = ?", templateID).
		Order("created_at DESC")
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err := db.Find(&reviews).Error
	return reviews, err
}

// FindByUser lists reviews created by a given user.
func (r *ReviewRepository) FindByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Review, error) {
	var reviews []*models.Review
	db := r.GetDB(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")
	if limit > 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err := db.Find(&reviews).Error
	return reviews, err
}

// Update updates an existing review.
func (r *ReviewRepository) Update(ctx context.Context, review *models.Review) error {
	return r.GetDB(ctx).Save(review).Error
}

// IncrementHelpfulCount increments the helpful_count for a review.
func (r *ReviewRepository) IncrementHelpfulCount(ctx context.Context, id uuid.UUID) error {
	return r.GetDB(ctx).
		Model(&models.Review{}).
		Where("id = ?", id).
		UpdateColumn("helpful_count", gorm.Expr("helpful_count + 1")).Error
}

