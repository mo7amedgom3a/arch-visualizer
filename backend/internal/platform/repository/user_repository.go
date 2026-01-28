package repository

import (
	"context"

	"github.com/google/uuid"
	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// UserRepository defines operations for user management
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new user repository
func NewUserRepository() (*UserRepository, error) {
	base, err := NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &UserRepository{BaseRepository: base}, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.GetDB(ctx).Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.GetDB(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "user", "UserRepository.FindByID")
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.GetDB(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "user", "UserRepository.FindByEmail")
	}
	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.GetDB(ctx).Save(user).Error
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.GetDB(ctx).Delete(&models.User{}, "id = ?", id).Error
}

// List lists all users with pagination
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	err := r.GetDB(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}
