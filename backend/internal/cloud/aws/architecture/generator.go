package architecture

import (
	"encoding/json"
	"fmt"
	"strings"

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

	// Convert diagram variables to architecture variables
	for _, v := range diagramGraph.Variables {
		arch.Variables = append(arch.Variables, architecture.Variable{
			Name:        v.Name,
			Type:        v.Type,
			Description: v.Description,
			Default:     v.Default,
			Sensitive:   v.Sensitive,
		})
	}

	// Convert diagram outputs to architecture outputs
	for _, o := range diagramGraph.Outputs {
		arch.Outputs = append(arch.Outputs, architecture.Output{
			Name:        o.Name,
			Value:       o.Value,
			Description: o.Description,
			Sensitive:   o.Sensitive,
		})
	}

	// Build node ID to resource ID mapping (first pass)
	// Include ALL nodes (including visual-only) for database persistence
	nodeIDToResourceID := make(map[string]string)
	for _, node := range diagramGraph.Nodes {
		if node.IsRegion() {
			continue
		}
		nodeIDToResourceID[node.ID] = node.ID
	}

	// Create domain resources (second pass)
	// Include ALL nodes (including visual-only) for database persistence
	for _, node := range diagramGraph.Nodes {
		if node.IsRegion() {
			continue
		}

		resourceID := nodeIDToResourceID[node.ID]

		// Map IR resource type to domain resource type using AWS resource type mapper
		// For visual-only nodes, use a generic "VisualIcon" type if mapping fails
		domainResourceType, err := g.mapIRResourceTypeToDomain(node.ResourceType)
		if err != nil {
			if node.IsVisualOnly {
				// Create a generic visual icon type for visual-only nodes
				domainResourceType = &resource.ResourceType{
					ID:         node.ResourceType,
					Name:       node.ResourceType,
					Category:   "Visual",
					Kind:       "Icon",
					IsRegional: false,
					IsGlobal:   false,
				}
			} else {
				return nil, fmt.Errorf("failed to map resource type for node %s: %w", node.ID, err)
			}
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
		// Add UI State
		if node.UI != nil {
			metadata["ui"] = node.UI
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

	// Build containment relationships (include visual-only nodes)
	for _, node := range diagramGraph.Nodes {
		if node.IsRegion() {
			continue
		}

		if node.ParentID != nil {
			parentNode, exists := diagramGraph.GetNode(*node.ParentID)
			if exists && !parentNode.IsRegion() {
				parentResourceID, parentOk := nodeIDToResourceID[*node.ParentID]
				childResourceID, childOk := nodeIDToResourceID[node.ID]

				if parentOk && childOk {
					if _, exists := arch.Containments[parentResourceID]; !exists {
						arch.Containments[parentResourceID] = make([]string, 0)
					}
					arch.Containments[parentResourceID] = append(arch.Containments[parentResourceID], childResourceID)
				}
			}
		}
	}

	// Build dependency relationships
	for _, res := range arch.Resources {
		if len(res.DependsOn) > 0 {
			arch.Dependencies[res.ID] = res.DependsOn
		}
	}

	// Process policies
	for _, policy := range diagramGraph.Policies {
		// 1. Create IAM Policy Resource
		policyDocBytes, err := json.Marshal(policy.PolicyDocument)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal policy document for edge %s: %w", policy.EdgeID, err)
		}
		policyDocJSON := string(policyDocBytes)

		// Sanitize ID
		safeEdgeID := strings.ReplaceAll(policy.EdgeID, "edge-", "")
		policyName := fmt.Sprintf("policy-%s", safeEdgeID)
		policyResourceID := fmt.Sprintf("policy-%s", safeEdgeID)

		policyMetadata := map[string]interface{}{
			"name":        policyName,
			"policy":      policyDocJSON,
			"description": fmt.Sprintf("Policy for edge %s", policy.EdgeID),
		}

		policyResource := &resource.Resource{
			ID:   policyResourceID,
			Name: policyName,
			Type: resource.ResourceType{
				ID:       "IAMPolicy", // Must match ResourceName in inventory
				Name:     "IAMPolicy",
				Category: "IAM",
				Kind:     "IAM",
			},
			Provider: resource.AWS,
			Region:   arch.Region,
			Metadata: policyMetadata,
		}
		arch.Resources = append(arch.Resources, policyResource)

		// 2. Create Attachment Resource
		if policy.Role != "" {
			attachmentName := fmt.Sprintf("attach-%s", safeEdgeID)
			attachmentResourceID := fmt.Sprintf("attach-%s", safeEdgeID)

			attachmentMetadata := map[string]interface{}{
				"name":       attachmentName,
				"role":       policy.Role,
				"policy_arn": fmt.Sprintf("aws_iam_policy.%s.arn", policyName),
			}

			attachmentResource := &resource.Resource{
				ID:   attachmentResourceID,
				Name: attachmentName,
				Type: resource.ResourceType{
					ID:       "IAMRolePolicyAttachment", // Must match ResourceName in inventory
					Name:     "IAMRolePolicyAttachment",
					Category: "IAM",
					Kind:     "IAM",
				},
				Provider:  resource.AWS,
				Region:    arch.Region,
				DependsOn: []string{policyResourceID},
				Metadata:  attachmentMetadata,
			}
			arch.Resources = append(arch.Resources, attachmentResource)

			// Update dependency map
			if arch.Dependencies == nil {
				arch.Dependencies = make(map[string][]string)
			}
			arch.Dependencies[attachmentResourceID] = []string{policyResourceID}
		}
	}

	// 5. Create Explicit Edge Resources
	// Store edges as generic resources to preserve metadata (style, markers, etc.)
	for _, edge := range diagramGraph.Edges {
		// Only persist dependency edges as resources for now
		if edge.IsDependency() {
			edgeResourceID := edge.ID
			if edgeResourceID == "" {
				edgeResourceID = fmt.Sprintf("edge-%s-%s", edge.Source, edge.Target)
			}

			edgeName := edgeResourceID
			if label, ok := edge.Config["label"].(string); ok && label != "" {
				edgeName = label
			}

			edgeMetadata := make(map[string]interface{})
			for k, v := range edge.Config {
				edgeMetadata[k] = v
			}
			// Set mandatory internal fields
			edgeMetadata["source"] = edge.Source
			edgeMetadata["target"] = edge.Target
			edgeMetadata["isVisualOnly"] = true // Edges are visual/structural, not cloud infrastructure

			edgeResource := &resource.Resource{
				ID:   edgeResourceID,
				Name: edgeName,
				Type: resource.ResourceType{
					ID:         "GenericEdge",
					Name:       "GenericEdge",
					Category:   "Visual",
					Kind:       "Connection",
					IsRegional: false,
					IsGlobal:   false,
				},
				Provider: resource.AWS,
				Region:   arch.Region,
				Metadata: edgeMetadata,
				// Dependencies are implicit in graph, but we can list source/target if needed
				// For now, leave DependsOn empty to avoid circular logic or double-counting
			}
			arch.Resources = append(arch.Resources, edgeResource)
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
