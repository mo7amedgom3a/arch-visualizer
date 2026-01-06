package constraints

import (
	"context"
	"fmt"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
)

// RequiresParentRule validates that a resource has a required parent
type RequiresParentRule struct {
	ResourceType string
	ParentType   string // The type of parent required
	MinCount     int    // Minimum number of parents required (default: 1)
}

func (r *RequiresParentRule) GetType() rules.RuleType {
	return rules.RuleTypeRequiresParent
}

func (r *RequiresParentRule) GetResourceType() string {
	return r.ResourceType
}

func (r *RequiresParentRule) GetValue() string {
	return r.ParentType
}

func (r *RequiresParentRule) Evaluate(ctx context.Context, evalCtx *rules.EvaluationContext) error {
	if evalCtx.Resource == nil {
		return fmt.Errorf("resource is required for evaluation")
	}
	
	// Count parents of the required type
	count := 0
	for _, parent := range evalCtx.Parents {
		if parent.Type.Kind == r.ParentType || parent.Type.Name == r.ParentType {
			count++
		}
	}
	
	minCount := r.MinCount
	if minCount == 0 {
		minCount = 1 // Default to 1
	}
	
	if count < minCount {
		return &rules.RuleError{
			RuleType:     rules.RuleTypeRequiresParent,
			ResourceID:   evalCtx.Resource.ID,
			ResourceName: evalCtx.Resource.Name,
			ResourceType: r.ResourceType,
			Message:      fmt.Sprintf("resource requires at least %d parent(s) of type '%s', but has %d", minCount, r.ParentType, count),
			Value:        r.ParentType,
		}
	}
	
	return nil
}

// AllowedParentRule validates that a resource only has allowed parent types
type AllowedParentRule struct {
	ResourceType  string
	AllowedTypes  []string // List of allowed parent types
	AllowMultiple bool    // Whether multiple parents are allowed
}

func (r *AllowedParentRule) GetType() rules.RuleType {
	return rules.RuleTypeAllowedParent
}

func (r *AllowedParentRule) GetResourceType() string {
	return r.ResourceType
}

func (r *AllowedParentRule) GetValue() string {
	// Return comma-separated list of allowed types
	result := ""
	for i, t := range r.AllowedTypes {
		if i > 0 {
			result += ","
		}
		result += t
	}
	return result
}

func (r *AllowedParentRule) Evaluate(ctx context.Context, evalCtx *rules.EvaluationContext) error {
	if evalCtx.Resource == nil {
		return fmt.Errorf("resource is required for evaluation")
	}
	
	// Check if all parents are in the allowed list
	for _, parent := range evalCtx.Parents {
		allowed := false
		for _, allowedType := range r.AllowedTypes {
			if parent.Type.Kind == allowedType || parent.Type.Name == allowedType {
				allowed = true
				break
			}
		}
		
		if !allowed {
			return &rules.RuleError{
				RuleType:     rules.RuleTypeAllowedParent,
				ResourceID:   evalCtx.Resource.ID,
				ResourceName: evalCtx.Resource.Name,
				ResourceType: r.ResourceType,
				Message:      fmt.Sprintf("resource has parent of type '%s' which is not allowed. Allowed types: %v", parent.Type.Name, r.AllowedTypes),
				Value:        r.GetValue(),
			}
		}
	}
	
	// Check multiple parents constraint
	if !r.AllowMultiple && len(evalCtx.Parents) > 1 {
		return &rules.RuleError{
			RuleType:     rules.RuleTypeAllowedParent,
			ResourceID:   evalCtx.Resource.ID,
			ResourceName: evalCtx.Resource.Name,
			ResourceType: r.ResourceType,
			Message:      fmt.Sprintf("resource has %d parents but only 1 is allowed", len(evalCtx.Parents)),
			Value:        r.GetValue(),
		}
	}
	
	return nil
}

// NewRequiresParentRule creates a new RequiresParentRule
func NewRequiresParentRule(resourceType, parentType string) *RequiresParentRule {
	return &RequiresParentRule{
		ResourceType: resourceType,
		ParentType:   parentType,
		MinCount:     1,
	}
}

// NewAllowedParentRule creates a new AllowedParentRule
func NewAllowedParentRule(resourceType string, allowedTypes []string) *AllowedParentRule {
	return &AllowedParentRule{
		ResourceType:  resourceType,
		AllowedTypes:  allowedTypes,
		AllowMultiple: false,
	}
}
