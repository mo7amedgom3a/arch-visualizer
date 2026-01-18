package compute

import (
	"context"
)

// ComputeService defines the interface for compute resource operations
// This is cloud-agnostic and can be implemented by any cloud provider
type ComputeService interface {
	// Instance operations
	CreateInstance(ctx context.Context, instance *Instance) (*Instance, error)
	GetInstance(ctx context.Context, id string) (*Instance, error)
	UpdateInstance(ctx context.Context, instance *Instance) (*Instance, error)
	DeleteInstance(ctx context.Context, id string) error
	ListInstances(ctx context.Context, filters map[string]string) ([]*Instance, error)

	// Instance lifecycle operations
	StartInstance(ctx context.Context, id string) error
	StopInstance(ctx context.Context, id string) error
	RebootInstance(ctx context.Context, id string) error

	// Launch Template operations
	CreateLaunchTemplate(ctx context.Context, template *LaunchTemplate) (*LaunchTemplate, error)
	GetLaunchTemplate(ctx context.Context, id string) (*LaunchTemplate, error)
	UpdateLaunchTemplate(ctx context.Context, id string, template *LaunchTemplate) (*LaunchTemplate, error)
	DeleteLaunchTemplate(ctx context.Context, id string) error
	ListLaunchTemplates(ctx context.Context, filters map[string]string) ([]*LaunchTemplate, error)
	GetLaunchTemplateVersion(ctx context.Context, id string, version int) (*LaunchTemplate, error)
	ListLaunchTemplateVersions(ctx context.Context, id string) ([]*LaunchTemplateVersion, error)
}

// ComputeRepository defines the interface for compute resource persistence
// This abstracts data access and can be implemented for different storage backends
type ComputeRepository interface {
	// Instance persistence
	SaveInstance(ctx context.Context, instance *Instance) error
	FindInstanceByID(ctx context.Context, id string) (*Instance, error)
	FindInstancesByFilters(ctx context.Context, filters map[string]string) ([]*Instance, error)
	DeleteInstance(ctx context.Context, id string) error
}
