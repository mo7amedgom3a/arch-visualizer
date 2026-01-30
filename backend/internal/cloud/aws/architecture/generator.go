package architecture

import (
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// AWSArchitectureGenerator implements ArchitectureGenerator for AWS
type AWSArchitectureGenerator struct{}

// NewAWSArchitectureGenerator creates a new AWS architecture generator
func NewAWSArchitectureGenerator() *AWSArchitectureGenerator {
	return &AWSArchitectureGenerator{}
}

// Provider returns AWS as the cloud provider
func (g *AWSArchitectureGenerator) Provider() resource.CloudProvider {
	return resource.AWS
}

// Generate converts a diagram graph into a domain architecture for AWS
func (g *AWSArchitectureGenerator) Generate(diagramGraph *graph.DiagramGraph) (*architecture.Architecture, error) {
	arch := architecture.NewArchitecture()
	arch.Provider = resource.AWS

	// Extract region from region node
	regionNode, hasRegion := diagramGraph.FindRegionNode()
	if hasRegion {
		if regionName, ok := extractRegionFromConfig(regionNode.Config); ok {
			arch.Region = regionName
		}
	}

	// Build node ID to resource ID mapping (first pass)
	nodeIDToResourceID := make(map[string]string)
	for _, node := range diagramGraph.Nodes {
		if node.IsRegion() || node.IsVisualOnly {
			continue
		}
		nodeIDToResourceID[node.ID] = node.ID
	}

	// Create domain resources (second pass)
	for _, node := range diagramGraph.Nodes {
		if node.IsRegion() || node.IsVisualOnly {
			continue
		}

		resourceID := nodeIDToResourceID[node.ID]

		// Map IR resource type to domain resource type using AWS resource type mapper
		domainResourceType, err := g.mapIRResourceTypeToDomain(node.ResourceType)
		if err != nil {
			return nil, fmt.Errorf("failed to map resource type for node %s: %w", node.ID, err)
		}

		// Extract name from config
		name := extractNameFromConfig(node.Config, node.Label)

		// Extract parent ID (if not region)
		var parentID *string
		if node.ParentID != nil {
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

		// Prepare metadata
		metadata := make(map[string]interface{})
		for k, v := range node.Config {
			metadata[k] = v
		}
		metadata["position"] = map[string]interface{}{
			"x": node.PositionX,
			"y": node.PositionY,
		}
		metadata["isVisualOnly"] = node.IsVisualOnly

		// Create domain resource
		domainResource := &resource.Resource{
			ID:        resourceID,
			Name:      name,
			Type:      *domainResourceType,
			Provider:  resource.AWS,
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
	for _, res := range arch.Resources {
		if len(res.DependsOn) > 0 {
			arch.Dependencies[res.ID] = res.DependsOn
		}
	}

	return arch, nil
}

// mapIRResourceTypeToDomain maps IR resource type to domain ResourceType using AWS resource type mapper
func (g *AWSArchitectureGenerator) mapIRResourceTypeToDomain(irType string) (*resource.ResourceType, error) {
	// Use AWS resource type mapper (no fallback - provider must define all mappings)
	mapper, ok := architecture.GetResourceTypeMapper(resource.AWS)
	if !ok {
		return nil, fmt.Errorf("AWS resource type mapper not registered")
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
