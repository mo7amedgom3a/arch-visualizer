package architecture

import (
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// Architecture represents a complete cloud architecture aggregate
// This is the domain model that represents the user's intent
type Architecture struct {
	// Resources in the architecture
	Resources []*resource.Resource

	// Region configuration (extracted from region node)
	Region   string
	Provider resource.CloudProvider

	// Containment relationships (parent -> child)
	Containments map[string][]string // parentID -> []childIDs

	// Dependency relationships
	Dependencies map[string][]string // resourceID -> []dependencyIDs
}

// NewArchitecture creates a new architecture aggregate
func NewArchitecture() *Architecture {
	return &Architecture{
		Resources:    make([]*resource.Resource, 0),
		Containments: make(map[string][]string),
		Dependencies: make(map[string][]string),
	}
}

// Intent represents the extracted intent from the diagram
// This includes region, provider, and other high-level configuration
type Intent struct {
	Region   string
	Provider resource.CloudProvider
}

// MapDiagramToArchitecture converts a diagram graph into a domain architecture
// This function uses cloud provider-specific generators when available,
// otherwise falls back to a default implementation for backward compatibility
func MapDiagramToArchitecture(diagramGraph *graph.DiagramGraph, provider resource.CloudProvider) (*Architecture, error) {
	// Try to use provider-specific generator (preferred)
	if generator, ok := GetGenerator(provider); ok {
		return generator.Generate(diagramGraph)
	}

	// Fall back to default implementation for backward compatibility
	return mapDiagramToArchitectureDefault(diagramGraph, provider)
}

// mapDiagramToArchitectureDefault provides the default implementation for backward compatibility
// This is used when no provider-specific generator is registered
func mapDiagramToArchitectureDefault(diagramGraph *graph.DiagramGraph, provider resource.CloudProvider) (*Architecture, error) {
	arch := NewArchitecture()
	arch.Provider = provider

	// Extract region and intent from region node
	regionNode, hasRegion := diagramGraph.FindRegionNode()
	if hasRegion {
		// Extract region from config
		if regionName, ok := extractRegionFromConfig(regionNode.Config); ok {
			arch.Region = regionName
		}
	}

	// Map nodes to domain resources (excluding region node and visual-only nodes)
	nodeIDToResourceID := make(map[string]string) // IR node ID -> domain resource ID

	// First pass: build complete nodeIDToResourceID map for all nodes
	// This ensures parent lookups work regardless of processing order
	for _, node := range diagramGraph.Nodes {
		// Skip region node - it's handled as project-level config
		if node.IsRegion() {
			continue
		}

		// Skip visual-only nodes - they are tracked but not persisted as real infrastructure
		if node.IsVisualOnly {
			continue
		}

		// Generate domain resource ID (using node ID as base)
		resourceID := node.ID
		nodeIDToResourceID[node.ID] = resourceID
	}

	// Second pass: create domain resources (now all parent IDs are available in the map)
	for _, node := range diagramGraph.Nodes {
		// Skip region node - it's handled as project-level config
		if node.IsRegion() {
			continue
		}

		// Skip visual-only nodes - they are tracked but not persisted as real infrastructure
		if node.IsVisualOnly {
			continue
		}

		// Get resource ID from map (already set in first pass)
		resourceID := nodeIDToResourceID[node.ID]

		// Map IR resource type to domain resource type using provider-specific mapper
		domainResourceType, err := mapIRResourceTypeToDomain(node.ResourceType, provider)
		if err != nil {
			return nil, fmt.Errorf("failed to map resource type for node %s: %w", node.ID, err)
		}

		// Extract name from config
		name := extractNameFromConfig(node.Config, node.Label)

		// Extract parent ID (if not region)
		var parentID *string
		if node.ParentID != nil {
			// If parent is region, don't set parentID (region is project-level)
			if parentNode, exists := diagramGraph.GetNode(*node.ParentID); exists && !parentNode.IsRegion() {
				if mappedParentID, ok := nodeIDToResourceID[*node.ParentID]; ok {
					parentID = &mappedParentID
				}
			}
		}

		// Build dependencies list
		dependencies := make([]string, 0)
		for _, edge := range diagramGraph.GetDependencyEdges() {
			if edge.Source == node.ID {
				if depID, ok := nodeIDToResourceID[edge.Target]; ok {
					dependencies = append(dependencies, depID)
				}
			}
		}

		// Prepare metadata with position and isVisualOnly flag
		metadata := make(map[string]interface{})
		// Copy existing config
		for k, v := range node.Config {
			metadata[k] = v
		}
		// Add position
		metadata["position"] = map[string]interface{}{
			"x": node.PositionX,
			"y": node.PositionY,
		}
		// Add isVisualOnly flag
		metadata["isVisualOnly"] = node.IsVisualOnly

		// Create domain resource
		domainResource := &resource.Resource{
			ID:        resourceID,
			Name:      name,
			Type:      *domainResourceType,
			Provider:  provider,
			Region:    arch.Region,
			ParentID:  parentID,
			DependsOn: dependencies,
			Metadata:  metadata,
		}

		arch.Resources = append(arch.Resources, domainResource)
	}

	// Build containment relationships
	for _, node := range diagramGraph.Nodes {
		if node.IsRegion() {
			continue
		}

		if node.ParentID != nil {
			parentNode, exists := diagramGraph.GetNode(*node.ParentID)
			if exists && !parentNode.IsRegion() {
				parentResourceID := nodeIDToResourceID[*node.ParentID]
				childResourceID := nodeIDToResourceID[node.ID]

				if _, exists := arch.Containments[parentResourceID]; !exists {
					arch.Containments[parentResourceID] = make([]string, 0)
				}
				arch.Containments[parentResourceID] = append(arch.Containments[parentResourceID], childResourceID)
			}
		}
	}

	// Build dependency relationships
	for _, resource := range arch.Resources {
		if len(resource.DependsOn) > 0 {
			arch.Dependencies[resource.ID] = resource.DependsOn
		}
	}

	return arch, nil
}

// mapIRResourceTypeToDomain maps IR resource type names to domain ResourceType
// Uses provider-specific resource type mapper (no fallback - each provider must define mappings)
func mapIRResourceTypeToDomain(irType string, provider resource.CloudProvider) (*resource.ResourceType, error) {
	mapper, ok := GetResourceTypeMapper(provider)
	if !ok {
		return nil, fmt.Errorf("no resource type mapper registered for provider: %s", provider)
	}

	return mapper.MapIRTypeToResourceType(irType)
}

// extractRegionFromConfig extracts the region name from a region node's config
func extractRegionFromConfig(config map[string]interface{}) (string, bool) {
	if name, ok := config["name"].(string); ok {
		return name, true
	}
	return "", false
}

// extractNameFromConfig extracts the resource name from config, falling back to label
func extractNameFromConfig(config map[string]interface{}, label string) string {
	if name, ok := config["name"].(string); ok && name != "" {
		return name
	}
	if label != "" {
		return label
	}
	return "unnamed-resource"
}
