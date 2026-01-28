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
// This is where we map IR resource types to domain resources and extract intent
func MapDiagramToArchitecture(diagramGraph *graph.DiagramGraph, provider resource.CloudProvider) (*Architecture, error) {
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
	resourceCounter := 0

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

		// Map IR resource type to domain resource type
		domainResourceType, err := mapIRResourceTypeToDomain(node.ResourceType)
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
		resourceCounter++
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
func mapIRResourceTypeToDomain(irType string) (*resource.ResourceType, error) {
	// Map IR resource type (kebab-case) to domain resource type (PascalCase)
	typeMapping := map[string]resource.ResourceType{
		"vpc": {
			ID:         "vpc",
			Name:       "VPC",
			Category:   "Networking",
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"subnet": {
			ID:         "subnet",
			Name:       "Subnet",
			Category:   "Networking",
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"ec2": {
			ID:         "ec2",
			Name:       "EC2",
			Category:   "Compute",
			Kind:       "VirtualMachine",
			IsRegional: true,
			IsGlobal:   false,
		},
		"route-table": {
			ID:         "route-table",
			Name:       "RouteTable",
			Category:   "Networking",
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"security-group": {
			ID:         "security-group",
			Name:       "SecurityGroup",
			Category:   "Networking",
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"nat-gateway": {
			ID:         "nat-gateway",
			Name:       "NATGateway",
			Category:   "Networking",
			Kind:       "Gateway",
			IsRegional: true,
			IsGlobal:   false,
		},
		"internet-gateway": {
			ID:         "internet-gateway",
			Name:       "InternetGateway",
			Category:   "Networking",
			Kind:       "Gateway",
			IsRegional: true,
			IsGlobal:   false,
		},
		"elastic-ip": {
			ID:         "elastic-ip",
			Name:       "ElasticIP",
			Category:   "Networking",
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"lambda": {
			ID:         "lambda",
			Name:       "Lambda",
			Category:   "Compute",
			Kind:       "Function",
			IsRegional: true,
			IsGlobal:   false,
		},
		"s3": {
			ID:         "s3",
			Name:       "S3",
			Category:   "Storage",
			Kind:       "Storage",
			IsRegional: false,
			IsGlobal:   true,
		},
		"ebs": {
			ID:         "ebs",
			Name:       "EBS",
			Category:   "Storage",
			Kind:       "Storage",
			IsRegional: true,
			IsGlobal:   false,
		},
		"rds": {
			ID:         "rds",
			Name:       "RDS",
			Category:   "Database",
			Kind:       "Database",
			IsRegional: true,
			IsGlobal:   false,
		},
		"dynamodb": {
			ID:         "dynamodb",
			Name:       "DynamoDB",
			Category:   "Database",
			Kind:       "Database",
			IsRegional: true,
			IsGlobal:   false,
		},
		"load-balancer": {
			ID:         "load-balancer",
			Name:       "LoadBalancer",
			Category:   "Compute",
			Kind:       "LoadBalancer",
			IsRegional: true,
			IsGlobal:   false,
		},
		"auto-scaling-group": {
			ID:         "auto-scaling-group",
			Name:       "AutoScalingGroup",
			Category:   "Compute",
			Kind:       "VirtualMachine",
			IsRegional: true,
			IsGlobal:   false,
		},
	}

	rt, exists := typeMapping[irType]
	if !exists {
		return nil, fmt.Errorf("unknown resource type: %s", irType)
	}

	return &rt, nil
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
