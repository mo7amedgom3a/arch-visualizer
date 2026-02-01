package pricing

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// HiddenDependency represents an implicit/hidden dependency that a resource requires
type HiddenDependency struct {
	// ParentResourceType is the resource type that requires this dependency
	ParentResourceType string `json:"parent_resource_type"`
	// ChildResourceType is the type of the hidden dependency resource
	ChildResourceType string `json:"child_resource_type"`
	// QuantityExpression is an expression to calculate quantity (e.g., "1", "metadata.size_gb")
	QuantityExpression string `json:"quantity_expression"`
	// ConditionExpression is an optional condition to check if dependency applies
	ConditionExpression string `json:"condition_expression,omitempty"`
	// IsAttached indicates if the dependency is attached to the parent (affects pricing)
	IsAttached bool `json:"is_attached"`
	// Description explains why this dependency exists
	Description string `json:"description,omitempty"`
}

// HiddenDependencyResolver defines the interface for resolving hidden dependencies
type HiddenDependencyResolver interface {
	// ResolveHiddenDependencies resolves all hidden dependencies for a resource
	ResolveHiddenDependencies(ctx context.Context, res *resource.Resource, architecture interface{}) ([]*HiddenDependencyResource, error)

	// GetHiddenDependenciesForResourceType returns all hidden dependencies for a resource type
	GetHiddenDependenciesForResourceType(ctx context.Context, provider CloudProvider, resourceType string) ([]*HiddenDependency, error)
}

// HiddenDependencyResource represents a resolved hidden dependency with its resource instance
type HiddenDependencyResource struct {
	// Dependency is the dependency definition
	Dependency *HiddenDependency `json:"dependency"`
	// Resource is the virtual resource instance created for this dependency
	Resource *resource.Resource `json:"resource"`
	// Quantity is the calculated quantity for this dependency
	Quantity float64 `json:"quantity"`
}
