package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// UserServiceImpl implements UserService interface
type UserServiceImpl struct {
	userRepo serverinterfaces.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo serverinterfaces.UserRepository) serverinterfaces.UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// Create creates a new user
func (s *UserServiceImpl) Create(ctx context.Context, req *serverinterfaces.CreateUserRequest) (*models.User, error) {
	if req == nil {
		return nil, fmt.Errorf("create user request is nil")
	}

	var avatar *string
	if req.AvatarURL != "" {
		avatar = &req.AvatarURL
	}

	user := &models.User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		Auth0ID:   req.Auth0ID,
		Avatar:    avatar,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *UserServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}
