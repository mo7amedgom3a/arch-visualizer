package architecture

import (
	"errors"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// TopologicalSortResult contains the sorted resources and metadata
type TopologicalSortResult struct {
	Resources []*resource.Resource // Sorted in provisioning order (dependencies first)
	Levels    [][]string           // Resources grouped by dependency level
	HasCycle  bool                 // True if cycle detected
	CycleInfo []string              // IDs involved in cycle (if any)
}

// TopologicalSort performs Kahn's algorithm on the architecture
// Returns resources in dependency order (dependencies first)
// Combines both explicit dependencies and containment relationships
// (children depend on their parents)
func (g *Graph) TopologicalSort() (*TopologicalSortResult, error) {
	if g.architecture == nil {
		return nil, errors.New("architecture is nil")
	}

	// Build resource ID to resource map for quick lookup
	resourceMap := make(map[string]*resource.Resource)
	for _, res := range g.architecture.Resources {
		resourceMap[res.ID] = res
	}

	if len(resourceMap) == 0 {
		return &TopologicalSortResult{
			Resources: []*resource.Resource{},
			Levels:    [][]string{},
			HasCycle:  false,
			CycleInfo: nil,
		}, nil
	}

	// Build adjacency list: for each resource, what depends on it
	// This is the reverse of dependencies (if A depends on B, then B -> A in adjacency)
	adjacency := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize in-degree for all resources
	for id := range resourceMap {
		inDegree[id] = 0
		adjacency[id] = make([]string, 0)
	}

	// Process explicit dependencies
	// Dependencies map: resourceID -> []dependencyIDs
	// This means resourceID depends on dependencyIDs
	// So dependencyID -> resourceID is an edge in our graph
	for resourceID, depIDs := range g.architecture.Dependencies {
		if _, exists := resourceMap[resourceID]; !exists {
			continue // Skip if resource doesn't exist
		}

		for _, depID := range depIDs {
			if _, exists := resourceMap[depID]; !exists {
				continue // Skip if dependency doesn't exist
			}
			// Add edge: depID -> resourceID (resource depends on dep)
			adjacency[depID] = append(adjacency[depID], resourceID)
			inDegree[resourceID]++
		}
	}

	// Process containment relationships
	// Containments map: parentID -> []childIDs
	// Children depend on parents, so parentID -> childID is an edge
	for parentID, childIDs := range g.architecture.Containments {
		if _, exists := resourceMap[parentID]; !exists {
			continue // Skip if parent doesn't exist
		}

		for _, childID := range childIDs {
			if _, exists := resourceMap[childID]; !exists {
				continue // Skip if child doesn't exist
			}
			// Add edge: parentID -> childID (child depends on parent)
			adjacency[parentID] = append(adjacency[parentID], childID)
			inDegree[childID]++
		}
	}

	// Also process ParentID relationships from resources themselves
	// This handles cases where ParentID is set but not in Containments map
	for _, res := range g.architecture.Resources {
		if res.ParentID != nil {
			parentID := *res.ParentID
			if _, parentExists := resourceMap[parentID]; parentExists {
				// Check if this edge already exists in adjacency
				// (might already be in Containments)
				edgeExists := false
				for _, childID := range adjacency[parentID] {
					if childID == res.ID {
						edgeExists = true
						break
					}
				}
				if !edgeExists {
					adjacency[parentID] = append(adjacency[parentID], res.ID)
					inDegree[res.ID]++
				}
			}
		}
	}

	// Initialize queue with resources having in-degree 0
	queue := make([]string, 0)
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}

	// Process queue using Kahn's algorithm
	result := make([]*resource.Resource, 0)
	levels := make([][]string, 0)
	currentLevel := make([]string, 0)

	processedCount := 0
	for len(queue) > 0 {
		// Process all nodes at current level
		currentLevelSize := len(queue)
		currentLevel = make([]string, 0)

		for i := 0; i < currentLevelSize; i++ {
			// Dequeue
			currentID := queue[0]
			queue = queue[1:]

			// Add to result
			if res, exists := resourceMap[currentID]; exists {
				result = append(result, res)
				currentLevel = append(currentLevel, currentID)
				processedCount++
			}

			// Decrement in-degree of dependents
			for _, dependentID := range adjacency[currentID] {
				inDegree[dependentID]--
				if inDegree[dependentID] == 0 {
					queue = append(queue, dependentID)
				}
			}
		}

		if len(currentLevel) > 0 {
			levels = append(levels, currentLevel)
		}
	}

	// Detect cycles: if not all resources were processed, there's a cycle
	hasCycle := processedCount != len(resourceMap)
	var cycleInfo []string

	if hasCycle {
		// Find resources that weren't processed (part of cycle)
		processedSet := make(map[string]bool)
		for _, res := range result {
			processedSet[res.ID] = true
		}

		for id := range resourceMap {
			if !processedSet[id] {
				cycleInfo = append(cycleInfo, id)
			}
		}
	}

	return &TopologicalSortResult{
		Resources: result,
		Levels:    levels,
		HasCycle:  hasCycle,
		CycleInfo: cycleInfo,
	}, nil
}

// GetSortedResources returns resources sorted in topological order
// This is a convenience method that returns just the sorted resources
// Returns an error if a cycle is detected
func (g *Graph) GetSortedResources() ([]*resource.Resource, error) {
	result, err := g.TopologicalSort()
	if err != nil {
		return nil, err
	}

	if result.HasCycle {
		return nil, fmt.Errorf("circular dependency detected involving resources: %v", result.CycleInfo)
	}

	return result.Resources, nil
}
