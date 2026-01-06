package engine

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
)

// BuildEvaluationContext builds an evaluation context for a resource
// This extracts parents, children, and dependencies from the architecture graph
func BuildEvaluationContext(
	res *resource.Resource,
	architecture *Architecture,
	provider string,
) *rules.EvaluationContext {
	ctx := &rules.EvaluationContext{
		Resource:    res,
		Parents:     []*resource.Resource{},
		Children:    []*resource.Resource{},
		Dependencies: []*resource.Resource{},
		Provider:    provider,
		Metadata:    make(map[string]interface{}),
	}
	
	if architecture == nil {
		return ctx
	}
	
	// Find parents (resources that contain this resource)
	if res.ParentID != nil {
		for _, r := range architecture.Resources {
			if r.ID == *res.ParentID {
				ctx.Parents = append(ctx.Parents, r)
				break
			}
		}
	}
	
	// Find children (resources contained by this resource)
	for _, r := range architecture.Resources {
		if r.ParentID != nil && *r.ParentID == res.ID {
			ctx.Children = append(ctx.Children, r)
		}
	}
	
	// Find dependencies
	for _, depID := range res.DependsOn {
		for _, r := range architecture.Resources {
			if r.ID == depID {
				ctx.Dependencies = append(ctx.Dependencies, r)
				break
			}
		}
	}
	
	return ctx
}

// Architecture represents a collection of resources and their relationships
// This is a simplified version - you may have a more complex architecture type
type Architecture struct {
	Resources []*resource.Resource
}
