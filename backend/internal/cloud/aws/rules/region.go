package rules

import (
	"context"
	"fmt"
)

// RequiresRegionRule validates that a resource has a region specified
type RequiresRegionRule struct {
	ResourceType string
	Required     bool // Whether region is required (true) or forbidden (false)
}

func (r *RequiresRegionRule) GetType() RuleType {
	return RuleTypeRequiresRegion
}

func (r *RequiresRegionRule) GetResourceType() string {
	return r.ResourceType
}

func (r *RequiresRegionRule) GetValue() string {
	if r.Required {
		return "true"
	}
	return "false"
}

func (r *RequiresRegionRule) Evaluate(ctx context.Context, evalCtx *EvaluationContext) error {
	if evalCtx.Resource == nil {
		return fmt.Errorf("resource is required for evaluation")
	}

	hasRegion := evalCtx.Resource.Region != ""

	if r.Required && !hasRegion {
		return &RuleError{
			RuleType:     RuleTypeRequiresRegion,
			ResourceID:   evalCtx.Resource.ID,
			ResourceName: evalCtx.Resource.Name,
			ResourceType: r.ResourceType,
			Message:      "resource requires a region but none is specified",
			Value:        "true",
		}
	}

	if !r.Required && hasRegion {
		return &RuleError{
			RuleType:     RuleTypeRequiresRegion,
			ResourceID:   evalCtx.Resource.ID,
			ResourceName: evalCtx.Resource.Name,
			ResourceType: r.ResourceType,
			Message:      "resource is global and should not have a region specified",
			Value:        "false",
		}
	}

	return nil
}

// NewRequiresRegionRule creates a new RequiresRegionRule
func NewRequiresRegionRule(resourceType string, required bool) *RequiresRegionRule {
	return &RequiresRegionRule{
		ResourceType: resourceType,
		Required:     required,
	}
}
