package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// UserService handles user management
type UserService interface {
	// Create creates a new user
	Create(ctx context.Context, req *CreateUserRequest) (*models.User, error)
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

// CreateUserRequest contains data needed to create a user
type CreateUserRequest struct {
	Name      string
	Email     string
	Auth0ID   string
	AvatarURL string
}

// User represents a user in the system (alias to model for now)
type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Auth0ID   string
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}
