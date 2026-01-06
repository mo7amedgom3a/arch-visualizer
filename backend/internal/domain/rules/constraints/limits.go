package constraints

import (
	"context"
	"fmt"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
)

// MaxChildrenRule validates that a resource doesn't exceed maximum children
type MaxChildrenRule struct {
	ResourceType string
	MaxCount     int
}

func (r *MaxChildrenRule) GetType() rules.RuleType {
	return rules.RuleTypeMaxChildren
}

func (r *MaxChildrenRule) GetResourceType() string {
	return r.ResourceType
}

func (r *MaxChildrenRule) GetValue() string {
	return fmt.Sprintf("%d", r.MaxCount)
}

func (r *MaxChildrenRule) Evaluate(ctx context.Context, evalCtx *rules.EvaluationContext) error {
	if evalCtx.Resource == nil {
		return fmt.Errorf("resource is required for evaluation")
	}
	
	count := len(evalCtx.Children)
	if count > r.MaxCount {
		return &rules.RuleError{
			RuleType:     rules.RuleTypeMaxChildren,
			ResourceID:   evalCtx.Resource.ID,
			ResourceName: evalCtx.Resource.Name,
			ResourceType: r.ResourceType,
			Message:      fmt.Sprintf("resource has %d children but maximum allowed is %d", count, r.MaxCount),
			Value:        r.GetValue(),
		}
	}
	
	return nil
}

// MinChildrenRule validates that a resource has minimum required children
type MinChildrenRule struct {
	ResourceType string
	MinCount     int
}

func (r *MinChildrenRule) GetType() rules.RuleType {
	return rules.RuleTypeMinChildren
}

func (r *MinChildrenRule) GetResourceType() string {
	return r.ResourceType
}

func (r *MinChildrenRule) GetValue() string {
	return fmt.Sprintf("%d", r.MinCount)
}

func (r *MinChildrenRule) Evaluate(ctx context.Context, evalCtx *rules.EvaluationContext) error {
	if evalCtx.Resource == nil {
		return fmt.Errorf("resource is required for evaluation")
	}
	
	count := len(evalCtx.Children)
	if count < r.MinCount {
		return &rules.RuleError{
			RuleType:     rules.RuleTypeMinChildren,
			ResourceID:   evalCtx.Resource.ID,
			ResourceName: evalCtx.Resource.Name,
			ResourceType: r.ResourceType,
			Message:      fmt.Sprintf("resource has %d children but minimum required is %d", count, r.MinCount),
			Value:        r.GetValue(),
		}
	}
	
	return nil
}

// NewMaxChildrenRule creates a new MaxChildrenRule
func NewMaxChildrenRule(resourceType string, maxCount int) *MaxChildrenRule {
	return &MaxChildrenRule{
		ResourceType: resourceType,
		MaxCount:     maxCount,
	}
}

// NewMinChildrenRule creates a new MinChildrenRule
func NewMinChildrenRule(resourceType string, minCount int) *MinChildrenRule {
	return &MinChildrenRule{
		ResourceType: resourceType,
		MinCount:     minCount,
	}
}
