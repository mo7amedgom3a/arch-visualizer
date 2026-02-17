package interfaces

import (
	"context"
)

// ResourceSchemaDTO is the service-layer representation of a resource schema.
type ResourceSchemaDTO struct {
	Label   string               `json:"label"`
	Fields  []FieldDescriptorDTO `json:"fields"`
	Outputs map[string]string    `json:"outputs"`
}

// FieldDescriptorDTO is the service-layer representation of a field descriptor.
type FieldDescriptorDTO struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Enum        []string `json:"enum"`
	Default     any      `json:"default,omitempty"`
	Description string   `json:"description,omitempty"`
}

// ResourceMetadataService exposes structured resource schemas.
type ResourceMetadataService interface {
	GetResourceSchema(ctx context.Context, provider, service, resource string) (*ResourceSchemaDTO, error)
	ListResourceSchemas(ctx context.Context, provider, service string) ([]*ResourceSchemaDTO, error)
}
