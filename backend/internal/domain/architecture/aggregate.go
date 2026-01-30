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

		// Map IR resource type to domain resource type (using dynamic inventory mapping)
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
// This function first tries to use provider-specific inventory mappers (dynamic),
// then falls back to static mapping for backward compatibility
func mapIRResourceTypeToDomain(irType string, provider resource.CloudProvider) (*resource.ResourceType, error) {
	// Step 1: Try to use provider-specific inventory mapper (dynamic)
	if mapper, ok := GetIRTypeMapper(provider); ok {
		if resourceName, found := mapper.GetResourceNameByIRType(irType); found {
			// Map resource name to ResourceType using domain-level mapping
			if rt, err := mapResourceNameToResourceType(resourceName); err == nil {
				return rt, nil
			}
		}
	}

	// Step 2: Fall back to static mapping for backward compatibility
	return mapIRResourceTypeToDomainStatic(irType)
}

// mapResourceNameToResourceType maps a domain resource name (PascalCase) to ResourceType
// This is a domain-level mapping that defines the structure of each resource type
func mapResourceNameToResourceType(resourceName string) (*resource.ResourceType, error) {
	resourceTypeMap := map[string]resource.ResourceType{
		"VPC": {
			ID:         "vpc",
			Name:       "VPC",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"Subnet": {
			ID:         "subnet",
			Name:       "Subnet",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"EC2": {
			ID:         "ec2",
			Name:       "EC2",
			Category:   string(resource.CategoryCompute),
			Kind:       "VirtualMachine",
			IsRegional: true,
			IsGlobal:   false,
		},
		"RouteTable": {
			ID:         "route-table",
			Name:       "RouteTable",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"SecurityGroup": {
			ID:         "security-group",
			Name:       "SecurityGroup",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"NATGateway": {
			ID:         "nat-gateway",
			Name:       "NATGateway",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Gateway",
			IsRegional: true,
			IsGlobal:   false,
		},
		"InternetGateway": {
			ID:         "internet-gateway",
			Name:       "InternetGateway",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Gateway",
			IsRegional: true,
			IsGlobal:   false,
		},
		"ElasticIP": {
			ID:         "elastic-ip",
			Name:       "ElasticIP",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"Lambda": {
			ID:         "lambda",
			Name:       "Lambda",
			Category:   string(resource.CategoryCompute),
			Kind:       "Function",
			IsRegional: true,
			IsGlobal:   false,
		},
		"S3": {
			ID:         "s3",
			Name:       "S3",
			Category:   string(resource.CategoryStorage),
			Kind:       "Storage",
			IsRegional: false,
			IsGlobal:   true,
		},
		"EBS": {
			ID:         "ebs",
			Name:       "EBS",
			Category:   string(resource.CategoryStorage),
			Kind:       "Storage",
			IsRegional: true,
			IsGlobal:   false,
		},
		"RDS": {
			ID:         "rds",
			Name:       "RDS",
			Category:   string(resource.CategoryDatabase),
			Kind:       "Database",
			IsRegional: true,
			IsGlobal:   false,
		},
		"DynamoDB": {
			ID:         "dynamodb",
			Name:       "DynamoDB",
			Category:   string(resource.CategoryDatabase),
			Kind:       "Database",
			IsRegional: true,
			IsGlobal:   false,
		},
		"LoadBalancer": {
			ID:         "load-balancer",
			Name:       "LoadBalancer",
			Category:   string(resource.CategoryCompute),
			Kind:       "LoadBalancer",
			IsRegional: true,
			IsGlobal:   false,
		},
		"AutoScalingGroup": {
			ID:         "auto-scaling-group",
			Name:       "AutoScalingGroup",
			Category:   string(resource.CategoryCompute),
			Kind:       "VirtualMachine",
			IsRegional: true,
			IsGlobal:   false,
		},
	}

	rt, exists := resourceTypeMap[resourceName]
	if !exists {
		return nil, fmt.Errorf("unknown resource name: %s", resourceName)
	}

	return &rt, nil
}

// mapIRResourceTypeToDomainStatic provides static fallback mapping for backward compatibility
// This is used when provider-specific inventory mappers are not available
func mapIRResourceTypeToDomainStatic(irType string) (*resource.ResourceType, error) {
	typeMapping := map[string]resource.ResourceType{
		"vpc": {
			ID:         "vpc",
			Name:       "VPC",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"subnet": {
			ID:         "subnet",
			Name:       "Subnet",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"ec2": {
			ID:         "ec2",
			Name:       "EC2",
			Category:   string(resource.CategoryCompute),
			Kind:       "VirtualMachine",
			IsRegional: true,
			IsGlobal:   false,
		},
		"route-table": {
			ID:         "route-table",
			Name:       "RouteTable",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"security-group": {
			ID:         "security-group",
			Name:       "SecurityGroup",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"nat-gateway": {
			ID:         "nat-gateway",
			Name:       "NATGateway",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Gateway",
			IsRegional: true,
			IsGlobal:   false,
		},
		"internet-gateway": {
			ID:         "internet-gateway",
			Name:       "InternetGateway",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Gateway",
			IsRegional: true,
			IsGlobal:   false,
		},
		"igw": {
			ID:         "internet-gateway",
			Name:       "InternetGateway",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Gateway",
			IsRegional: true,
			IsGlobal:   false,
		},
		"elastic-ip": {
			ID:         "elastic-ip",
			Name:       "ElasticIP",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"lambda": {
			ID:         "lambda",
			Name:       "Lambda",
			Category:   string(resource.CategoryCompute),
			Kind:       "Function",
			IsRegional: true,
			IsGlobal:   false,
		},
		"s3": {
			ID:         "s3",
			Name:       "S3",
			Category:   string(resource.CategoryStorage),
			Kind:       "Storage",
			IsRegional: false,
			IsGlobal:   true,
		},
		"ebs": {
			ID:         "ebs",
			Name:       "EBS",
			Category:   string(resource.CategoryStorage),
			Kind:       "Storage",
			IsRegional: true,
			IsGlobal:   false,
		},
		"rds": {
			ID:         "rds",
			Name:       "RDS",
			Category:   string(resource.CategoryDatabase),
			Kind:       "Database",
			IsRegional: true,
			IsGlobal:   false,
		},
		"dynamodb": {
			ID:         "dynamodb",
			Name:       "DynamoDB",
			Category:   string(resource.CategoryDatabase),
			Kind:       "Database",
			IsRegional: true,
			IsGlobal:   false,
		},
		"load-balancer": {
			ID:         "load-balancer",
			Name:       "LoadBalancer",
			Category:   string(resource.CategoryCompute),
			Kind:       "LoadBalancer",
			IsRegional: true,
			IsGlobal:   false,
		},
		"auto-scaling-group": {
			ID:         "auto-scaling-group",
			Name:       "AutoScalingGroup",
			Category:   string(resource.CategoryCompute),
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
