package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// Repository interfaces abstract the platform repositories for dependency injection

// ProjectRepository defines project repository operations
type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Project, error)
	FindAll(ctx context.Context, userID uuid.UUID, page, limit int, sort, order, search string) ([]*models.Project, int64, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	BeginTransaction(ctx context.Context) (interface{}, context.Context)
	CommitTransaction(tx interface{}) error
	RollbackTransaction(tx interface{}) error
}

// ProjectVersionRepository defines project version repository operations
type ProjectVersionRepository interface {
	Create(ctx context.Context, version *models.ProjectVersion) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.ProjectVersion, error)
	ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectVersion, error)
	Delete(ctx context.Context, id uuid.UUID) error
} // added ProjectVersionRepository definitions

// ResourceRepository defines resource repository operations
type ResourceRepository interface {
	Create(ctx context.Context, resource *models.Resource) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Resource, error)
	FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Resource, error)
	CreateContainment(ctx context.Context, parentID, childID uuid.UUID) error
	CreateDependency(ctx context.Context, dependency *models.ResourceDependency) error
}

// ResourceTypeRepository defines resource type repository operations
type ResourceTypeRepository interface {
	FindByNameAndProvider(ctx context.Context, name, provider string) (*models.ResourceType, error)
	ListByProvider(ctx context.Context, provider string) ([]*models.ResourceType, error)
}

// ResourceConstraintRepository defines constraint repository operations
type ResourceConstraintRepository interface {
	FindByResourceType(ctx context.Context, resourceTypeID uint) ([]*models.ResourceConstraint, error)
}

// DependencyTypeRepository defines dependency type repository operations
type DependencyTypeRepository interface {
	FindByName(ctx context.Context, name string) (*models.DependencyType, error)
	Create(ctx context.Context, depType *models.DependencyType) error
}

// UserRepository defines user repository operations
type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

// IACTargetRepository defines IaC target repository operations
type IACTargetRepository interface {
	FindByName(ctx context.Context, name string) (*models.IACTarget, error)
	Create(ctx context.Context, target *models.IACTarget) error
}

// ResourceContainmentRepository defines containment repository operations
type ResourceContainmentRepository interface {
	Create(ctx context.Context, containment *models.ResourceContainment) error
	FindChildren(ctx context.Context, parentID uuid.UUID) ([]*models.ResourceContainment, error)
	FindParents(ctx context.Context, childID uuid.UUID) ([]*models.ResourceContainment, error)
}

// ResourceDependencyRepository defines dependency repository operations
type ResourceDependencyRepository interface {
	Create(ctx context.Context, dependency *models.ResourceDependency) error
	FindByFromResource(ctx context.Context, fromID uuid.UUID) ([]*models.ResourceDependency, error)
	FindByToResource(ctx context.Context, toID uuid.UUID) ([]*models.ResourceDependency, error)
}

// ProjectVariableRepository defines project variable repository operations
type ProjectVariableRepository interface {
	Create(ctx context.Context, variable *models.ProjectVariable) error
	FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectVariable, error)
	DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error
}

// ProjectOutputRepository defines project output repository operations
type ProjectOutputRepository interface {
	Create(ctx context.Context, output *models.ProjectOutput) error
	FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectOutput, error)
	DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error
}
