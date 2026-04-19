package rules

import (
	"context"
	"fmt"
)

// AllowedDependenciesRule validates that dependencies are of allowed or forbidden types.
type AllowedDependenciesRule struct {
	ResourceType   string
	AllowedTypes   []string // Allowed dependency types
	ForbiddenTypes []string // Forbidden dependency types
}

func (r *AllowedDependenciesRule) GetType() RuleType {
	// If this rule was constructed as a "forbidden" rule (no allowed types, only forbidden),
	// report its type as ForbiddenDependencies so factories/tests can distinguish it.
	if len(r.ForbiddenTypes) > 0 && len(r.AllowedTypes) == 0 {
		return RuleTypeForbiddenDependencies
	}
	return RuleTypeAllowedDependencies
}

func (r *AllowedDependenciesRule) GetResourceType() string {
	return r.ResourceType
}

func (r *AllowedDependenciesRule) GetValue() string {
	// Return comma-separated list
	result := ""
	for i, t := range r.AllowedTypes {
		if i > 0 {
			result += ","
		}
		result += t
	}
	return result
}

func (r *AllowedDependenciesRule) Evaluate(ctx context.Context, evalCtx *EvaluationContext) error {
	if evalCtx.Resource == nil {
		return fmt.Errorf("resource is required for evaluation")
	}

	// Check forbidden types first
	for _, dep := range evalCtx.Dependencies {
		for _, forbiddenType := range r.ForbiddenTypes {
			if dep.Type.Kind == forbiddenType || dep.Type.Name == forbiddenType {
				return &RuleError{
					RuleType:     RuleTypeForbiddenDependencies,
					ResourceID:   evalCtx.Resource.ID,
					ResourceName: evalCtx.Resource.Name,
					ResourceType: r.ResourceType,
					Message:      fmt.Sprintf("resource has forbidden dependency on type '%s'", forbiddenType),
					Value:        forbiddenType,
				}
			}
		}
	}

	// Check allowed types (if specified)
	if len(r.AllowedTypes) > 0 {
		for _, dep := range evalCtx.Dependencies {
			allowed := false
			for _, allowedType := range r.AllowedTypes {
				if dep.Type.Kind == allowedType || dep.Type.Name == allowedType {
					allowed = true
					break
				}
			}

			if !allowed {
				return &RuleError{
					RuleType:     RuleTypeAllowedDependencies,
					ResourceID:   evalCtx.Resource.ID,
					ResourceName: evalCtx.Resource.Name,
					ResourceType: r.ResourceType,
					Message:      fmt.Sprintf("resource has dependency on type '%s' which is not allowed. Allowed types: %v", dep.Type.Name, r.AllowedTypes),
					Value:        r.GetValue(),
				}
			}
		}
	}

	return nil
}

// NewAllowedDependenciesRule creates a new AllowedDependenciesRule
func NewAllowedDependenciesRule(resourceType string, allowedTypes []string) *AllowedDependenciesRule {
	return &AllowedDependenciesRule{
		ResourceType:   resourceType,
		AllowedTypes:   allowedTypes,
		ForbiddenTypes: []string{},
	}
}

// NewForbiddenDependenciesRule creates a rule that forbids certain dependency types
func NewForbiddenDependenciesRule(resourceType string, forbiddenTypes []string) *AllowedDependenciesRule {
	return &AllowedDependenciesRule{
		ResourceType:   resourceType,
		AllowedTypes:   []string{},
		ForbiddenTypes: forbiddenTypes,
	}
}

// RequiresDependencyRule validates that a resource has a required dependency.
type RequiresDependencyRule struct {
	ResourceType string
	RequiredType string
}

func (r *RequiresDependencyRule) GetType() RuleType {
	return RuleTypeRequiresDependency
}

func (r *RequiresDependencyRule) GetResourceType() string {
	return r.ResourceType
}

func (r *RequiresDependencyRule) GetValue() string {
	return r.RequiredType
}

func (r *RequiresDependencyRule) Evaluate(ctx context.Context, evalCtx *EvaluationContext) error {
	if evalCtx.Resource == nil {
		return fmt.Errorf("resource is required for evaluation")
	}

	found := false
	for _, dep := range evalCtx.Dependencies {
		if dep.Type.Name == r.RequiredType {
			found = true
			break
		}
	}

	if !found {
		return &RuleError{
			RuleType:     RuleTypeRequiresDependency,
			ResourceID:   evalCtx.Resource.ID,
			ResourceName: evalCtx.Resource.Name,
			ResourceType: r.ResourceType,
			Message:      fmt.Sprintf("resource requires dependency of type '%s'", r.RequiredType),
			Value:        r.RequiredType,
		}
	}

	return nil
}

// NewRequiresDependencyRule creates a new RequiresDependencyRule
func NewRequiresDependencyRule(resourceType string, requiredType string) *RequiresDependencyRule {
	return &RequiresDependencyRule{
		ResourceType: resourceType,
		RequiredType: requiredType,
	}
}
