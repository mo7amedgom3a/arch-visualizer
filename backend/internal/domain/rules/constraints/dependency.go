package constraints

import (
	"context"
	"fmt"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
)

// AllowedDependenciesRule validates that dependencies are of allowed types
type AllowedDependenciesRule struct {
	ResourceType      string
	AllowedTypes      []string // Allowed dependency types
	ForbiddenTypes    []string // Forbidden dependency types
}

func (r *AllowedDependenciesRule) GetType() rules.RuleType {
	return rules.RuleTypeAllowedDependencies
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

func (r *AllowedDependenciesRule) Evaluate(ctx context.Context, evalCtx *rules.EvaluationContext) error {
	if evalCtx.Resource == nil {
		return fmt.Errorf("resource is required for evaluation")
	}
	
	// Check forbidden types first
	for _, dep := range evalCtx.Dependencies {
		for _, forbiddenType := range r.ForbiddenTypes {
			if dep.Type.Kind == forbiddenType || dep.Type.Name == forbiddenType {
				return &rules.RuleError{
					RuleType:     rules.RuleTypeForbiddenDependencies,
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
				return &rules.RuleError{
					RuleType:     rules.RuleTypeAllowedDependencies,
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
