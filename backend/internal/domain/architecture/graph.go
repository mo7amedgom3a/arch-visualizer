package architecture

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// Graph represents the graph structure of an architecture
// It provides methods to query relationships
type Graph struct {
	architecture *Architecture
}

// NewGraph creates a new graph wrapper around an architecture
func NewGraph(arch *Architecture) *Graph {
	return &Graph{architecture: arch}
}

// GetResource returns a resource by ID
func (g *Graph) GetResource(id string) (*resource.Resource, bool) {
	for _, res := range g.architecture.Resources {
		if res.ID == id {
			return res, true
		}
	}
	return nil, false
}

// GetChildren returns all child resources of a given parent
func (g *Graph) GetChildren(parentID string) []*resource.Resource {
	children := make([]*resource.Resource, 0)
	if childIDs, exists := g.architecture.Containments[parentID]; exists {
		for _, childID := range childIDs {
			if child, found := g.GetResource(childID); found {
				children = append(children, child)
			}
		}
	}
	return children
}

// GetParent returns the parent resource of a given child
func (g *Graph) GetParent(childID string) (*resource.Resource, bool) {
	child, exists := g.GetResource(childID)
	if !exists || child.ParentID == nil {
		return nil, false
	}
	return g.GetResource(*child.ParentID)
}

// GetDependencies returns all resources that a resource depends on
func (g *Graph) GetDependencies(resourceID string) []*resource.Resource {
	dependencies := make([]*resource.Resource, 0)
	if depIDs, exists := g.architecture.Dependencies[resourceID]; exists {
		for _, depID := range depIDs {
			if dep, found := g.GetResource(depID); found {
				dependencies = append(dependencies, dep)
			}
		}
	}
	return dependencies
}

// GetRootResources returns all resources that have no parent
func (g *Graph) GetRootResources() []*resource.Resource {
	roots := make([]*resource.Resource, 0)
	for _, res := range g.architecture.Resources {
		if res.ParentID == nil {
			roots = append(roots, res)
		}
	}
	return roots
}

// BuildContainmentTree returns a map of parent ID -> list of child resources
func (g *Graph) BuildContainmentTree() map[string][]*resource.Resource {
	tree := make(map[string][]*resource.Resource)
	for parentID, childIDs := range g.architecture.Containments {
		children := make([]*resource.Resource, 0)
		for _, childID := range childIDs {
			if child, found := g.GetResource(childID); found {
				children = append(children, child)
			}
		}
		tree[parentID] = children
	}
	return tree
}
