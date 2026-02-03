package interfaces

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ResourceTypeGroup represents a group of resource types by service category
type ResourceTypeGroup struct {
	ServiceType string                 `json:"service_type"` // e.g., "Compute"
	Resources   []*models.ResourceType `json:"resources"`
}

// ResourceModel represents a resource output model with its name
type ResourceModel struct {
	Name  string      `json:"name"`
	Model interface{} `json:"model"`
}

// ResourceModelGroup represents a group of resource models by service category
type ResourceModelGroup struct {
	ServiceType string          `json:"service_type"` // e.g., "Networking", "Compute"
	Resources   []ResourceModel `json:"resources"`
}

// CloudConfig represents global cloud configurations
type CloudConfig struct {
	Provider      string               `json:"provider"`
	Regions       []RegionConfig       `json:"regions"`
	InstanceTypes []InstanceTypeConfig `json:"instance_types"`
	StorageTypes  []string             `json:"storage_types"`
}

// RegionConfig represents a cloud region
type RegionConfig struct {
	Code string   `json:"code"` // e.g., "us-east-1"
	Name string   `json:"name"` // e.g., "US East (N. Virginia)"
	AZs  []string `json:"azs"`  // e.g., ["us-east-1a", "us-east-1b"]
}

// InstanceTypeConfig represents an EC2 instance type
type InstanceTypeConfig struct {
	Name      string  `json:"name"` // e.g., "t3.micro"
	VCPU      int     `json:"vcpu"`
	MemoryGiB float64 `json:"memory_gib"`
}

// StaticDataService handles static reference data
type StaticDataService interface {
	// ListResourceTypes retrieves all resource types grouped by category
	ListResourceTypes(ctx context.Context) ([]ResourceTypeGroup, error)
	// ListResourceTypesByProvider retrieves resource types for a provider grouped by category
	ListResourceTypesByProvider(ctx context.Context, provider string) ([]ResourceTypeGroup, error)
	// ListResourceOutputModels retrieves output models for resources with default values grouped by category
	ListResourceOutputModels(ctx context.Context, provider string) ([]ResourceModelGroup, error)
	// ListProviders retrieves supported cloud providers
	ListProviders(ctx context.Context) ([]string, error)
	// ListCloudConfiguration retrieves global cloud configurations
	ListCloudConfiguration(ctx context.Context, provider string) (*CloudConfig, error)
}
