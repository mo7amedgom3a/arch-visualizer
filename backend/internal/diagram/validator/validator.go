package validator

import (
	"fmt"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
)

// ValidationError represents a validation error
type ValidationError struct {
	Code    string
	Message string
	NodeID  string
}

func (e *ValidationError) Error() string {
	if e.NodeID != "" {
		return fmt.Sprintf("[%s] %s (node: %s)", e.Code, e.Message, e.NodeID)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// ValidationResult contains the result of validation
type ValidationResult struct {
	Valid    bool
	Errors   []*ValidationError
	Warnings []*ValidationError
}

// ValidationOptions contains options for validation
type ValidationOptions struct {
	// ValidResourceTypes maps resource type names (in IR format, lowercase) to validity
	// e.g., map["vpc"] = true, map["route-table"] = true
	ValidResourceTypes map[string]bool
	// Provider is the cloud provider (aws, azure, gcp)
	Provider string
}

// Validate performs comprehensive validation on the diagram graph
// If opts is nil, uses default validation (backward compatible)
func Validate(g *graph.DiagramGraph, opts *ValidationOptions) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   make([]*ValidationError, 0),
		Warnings: make([]*ValidationError, 0),
	}

	// Run all validation checks
	validateMissingParents(g, result)
	validateContainmentCycles(g, result)
	validateEdgeReferences(g, result)
	validateResourceTypes(g, result, opts)
	validateRegionNode(g, result)

	// Set valid to false if there are any errors
	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return result
}

// validateMissingParents checks that all parentId references point to existing nodes
func validateMissingParents(g *graph.DiagramGraph, result *ValidationResult) {
	for _, node := range g.Nodes {
		if node.ParentID != nil {
			if _, exists := g.Nodes[*node.ParentID]; !exists {
				result.Errors = append(result.Errors, &ValidationError{
					Code:    "MISSING_PARENT",
					Message: fmt.Sprintf("Node '%s' references non-existent parent '%s'", node.ID, *node.ParentID),
					NodeID:  node.ID,
				})
			}
		}
	}
}

// validateContainmentCycles checks for cycles in the containment tree
func validateContainmentCycles(g *graph.DiagramGraph, result *ValidationResult) {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var hasCycle func(nodeID string) bool
	hasCycle = func(nodeID string) bool {
		visited[nodeID] = true
		recStack[nodeID] = true

		// Check children (containment goes parent -> child)
		children := g.GetChildren(nodeID)
		for _, child := range children {
			if !visited[child.ID] {
				if hasCycle(child.ID) {
					return true
				}
			} else if recStack[child.ID] {
				// Found a cycle: we're trying to visit a node that's already in the recursion stack
				return true
			}
		}

		recStack[nodeID] = false
		return false
	}

	// Check all nodes for cycles (not just roots, since cycles might not include roots)
	// But we need to check from each unvisited node to catch all cycles
	for nodeID := range g.Nodes {
		if !visited[nodeID] {
			if hasCycle(nodeID) {
				result.Errors = append(result.Errors, &ValidationError{
					Code:    "CONTAINMENT_CYCLE",
					Message: fmt.Sprintf("Cycle detected in containment tree involving node '%s'", nodeID),
					NodeID:  nodeID,
				})
				return // Found a cycle, no need to continue
			}
		}
	}
}

// validateEdgeReferences checks that all edges reference existing nodes
func validateEdgeReferences(g *graph.DiagramGraph, result *ValidationResult) {
	for _, edge := range g.Edges {
		if _, exists := g.Nodes[edge.Source]; !exists {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "INVALID_EDGE_SOURCE",
				Message: fmt.Sprintf("Edge references non-existent source node '%s'", edge.Source),
			})
		}
		if _, exists := g.Nodes[edge.Target]; !exists {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "INVALID_EDGE_TARGET",
				Message: fmt.Sprintf("Edge references non-existent target node '%s'", edge.Target),
			})
		}
	}
}

// validateResourceTypes checks that resource types are valid
// If opts is nil or ValidResourceTypes is nil, skips validation (backward compatible)
func validateResourceTypes(g *graph.DiagramGraph, result *ValidationResult, opts *ValidationOptions) {
	// If no options provided, skip resource type validation (backward compatible)
	if opts == nil || opts.ValidResourceTypes == nil || len(opts.ValidResourceTypes) == 0 {
		return
	}

	for _, node := range g.Nodes {
		if node.ResourceType == "" {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "MISSING_RESOURCE_TYPE",
				Message: fmt.Sprintf("Node '%s' has no resource type", node.ID),
				NodeID:  node.ID,
			})
		} else {
			resourceTypeLower := strings.ToLower(node.ResourceType)
			if !opts.ValidResourceTypes[resourceTypeLower] {
				providerMsg := ""
				if opts.Provider != "" {
					providerMsg = fmt.Sprintf(" for provider '%s'", opts.Provider)
				}
				result.Warnings = append(result.Warnings, &ValidationError{
					Code:    "UNKNOWN_RESOURCE_TYPE",
					Message: fmt.Sprintf("Node '%s' has unknown resource type '%s'%s", node.ID, node.ResourceType, providerMsg),
					NodeID:  node.ID,
				})
			}
		}
	}
}

// validateRegionNode checks that there is exactly one region node
func validateRegionNode(g *graph.DiagramGraph, result *ValidationResult) {
	regionNodes := make([]*graph.Node, 0)
	for _, node := range g.Nodes {
		if node.IsRegion() {
			regionNodes = append(regionNodes, node)
		}
	}

	if len(regionNodes) == 0 {
		result.Warnings = append(result.Warnings, &ValidationError{
			Code:    "NO_REGION_NODE",
			Message: "No region node found in diagram",
		})
	} else if len(regionNodes) > 1 {
		result.Warnings = append(result.Warnings, &ValidationError{
			Code:    "MULTIPLE_REGION_NODES",
			Message: fmt.Sprintf("Multiple region nodes found (%d), only one will be used", len(regionNodes)),
		})
	}
}
