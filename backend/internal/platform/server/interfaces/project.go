package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ProjectService handles project management and persistence
type ProjectService interface {
	// Create creates a new project
	Create(ctx context.Context, req *CreateProjectRequest) (*models.Project, error)

	// GetByID retrieves a project by ID with related data
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)

	// PersistArchitecture persists an architecture to the database as part of a project
	PersistArchitecture(ctx context.Context, projectID uuid.UUID, arch *architecture.Architecture, diagramGraph interface{}) error

	// LoadArchitecture loads an architecture from the database for a project
	LoadArchitecture(ctx context.Context, projectID uuid.UUID) (*architecture.Architecture, error)
}

// CreateProjectRequest contains data needed to create a project
type CreateProjectRequest struct {
	UserID        uuid.UUID
	Name          string
	IACTargetID   uint
	CloudProvider string
	Region        string
}
