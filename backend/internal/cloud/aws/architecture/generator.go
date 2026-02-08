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

	// 6. Enrich architecture with default resources and fallbacks
	g.enrichArchitecture(arch)

	return arch, nil
}

// enrichArchitecture checks for missing required resources and creates default fallbacks
func (g *AWSArchitectureGenerator) enrichArchitecture(arch *architecture.Architecture) {
	// 1. Build map for quick lookups
	resMap := make(map[string]*resource.Resource)
	for _, res := range arch.Resources {
		resMap[res.ID] = res
	}

	// Helper to find VPC from resource (traversing up)
	findVPC := func(r *resource.Resource) *resource.Resource {
		curr := r
		for curr.ParentID != nil && *curr.ParentID != "" {
			parent, ok := resMap[*curr.ParentID]
			if !ok {
				break
			}
			if parent.Type.Name == "VPC" {
				return parent
			}
			curr = parent
		}
		return nil
	}

	// 2. Iterate resources
	newResources := make([]*resource.Resource, 0)

	for _, res := range arch.Resources {
		// --- EC2 Security Group Fallback ---
		if res.Type.Name == "EC2" {
			hasSG := false
			// Check metadata
			if _, ok := res.Metadata["securityGroups"]; ok {
				hasSG = true
			} else if _, ok := res.Metadata["securityGroupIds"]; ok {
				hasSG = true
			} else if v, ok := res.Metadata["vpc_security_group_ids"]; ok {
				if l, ok := v.([]interface{}); ok && len(l) > 0 {
					hasSG = true
				}
			}

			// Check depends_on
			if !hasSG {
				for _, depID := range res.DependsOn {
					if dep, ok := resMap[depID]; ok && dep.Type.Name == "SecurityGroup" {
						hasSG = true
						break
					}
				}
			}

			if !hasSG {
				// Fallback: Create/Use default SG in the parent VPC
				vpc := findVPC(res)
				if vpc == nil {
					// Cannot create SG without VPC
					continue
				}

				defaultSGID := fmt.Sprintf("default-sg-%s", vpc.ID)

				// Create if not exists in map (check both existing map and newResources if we were using a map for them)
				// Since we add to resMap immediately, we check resMap
				var defaultSG *resource.Resource
				if existing, ok := resMap[defaultSGID]; ok {
					defaultSG = existing
				} else {
					defaultSG = &resource.Resource{
						ID:   defaultSGID,
						Name: "default-security-group",
						Type: resource.ResourceType{
							ID:       "SecurityGroup",
							Name:     "SecurityGroup",
							Category: "Networking",
							Kind:     "Security",
						},
						Provider: resource.AWS,
						Region:   arch.Region,
						ParentID: &vpc.ID,
						Metadata: map[string]interface{}{
							"name":        "default-generated-sg",
							"description": "Default security group generated by fallback",
							"vpcId":       vpc.ID,
						},
					}
					newResources = append(newResources, defaultSG)
					resMap[defaultSGID] = defaultSG
				}

				// Link EC2 to this SG
				if res.Metadata == nil {
					res.Metadata = make(map[string]interface{})
				}

				// Update securityGroupIds list
				if list, ok := res.Metadata["securityGroupIds"].([]string); ok {
					res.Metadata["securityGroupIds"] = append(list, defaultSGID)
				} else {
					res.Metadata["securityGroupIds"] = []string{defaultSGID}
				}

				// Add explicit dependency
				res.DependsOn = append(res.DependsOn, defaultSGID)
				if arch.Dependencies == nil {
					arch.Dependencies = make(map[string][]string)
				}
				arch.Dependencies[res.ID] = append(arch.Dependencies[res.ID], defaultSGID)

				// Create Warning
				arch.Warnings = append(arch.Warnings, architecture.Warning{
					Message:    fmt.Sprintf("EC2 Instance '%s' is missing a security group. Using default security group '%s'.", res.Name, defaultSG.Name),
					ResourceID: res.ID,
				})
			}
		}

		// --- Auto Scaling Group Launch Template Fallback ---
		if res.Type.Name == "AutoScalingGroup" {
			hasLT := false
			if _, ok := res.Metadata["launchTemplate"]; ok {
				hasLT = true
			}

			if !hasLT {
				for _, depID := range res.DependsOn {
					if dep, ok := resMap[depID]; ok && dep.Type.Name == "LaunchTemplate" {
						hasLT = true
						break
					}
				}
			}

			if !hasLT {
				ltID := fmt.Sprintf("default-lt-%s", res.ID)
				var lt *resource.Resource
				if existing, ok := resMap[ltID]; ok {
					lt = existing
				} else {
					lt = &resource.Resource{
						ID:       ltID,
						Name:     fmt.Sprintf("lt-%s", res.Name),
						Type:     resource.ResourceType{ID: "LaunchTemplate", Name: "LaunchTemplate", Category: "Compute", Kind: "Template"},
						Provider: resource.AWS,
						Region:   arch.Region,
						Metadata: map[string]interface{}{
							"name":         fmt.Sprintf("lt-%s", res.Name),
							"imageId":      "ami-0c55b159cbfafe1f0", // Placeholder
							"instanceType": "t2.micro",
						},
					}
					newResources = append(newResources, lt)
					resMap[ltID] = lt
				}

				res.Metadata["launchTemplate"] = map[string]interface{}{
					"id": ltID,
				}
				res.DependsOn = append(res.DependsOn, ltID)

				arch.Warnings = append(arch.Warnings, architecture.Warning{
					Message:    fmt.Sprintf("Auto Scaling Group '%s' is missing a Launch Template. Using default launch template '%s'.", res.Name, lt.Name),
					ResourceID: res.ID,
				})
			}
		}

		// --- Elastic IP for NAT Gateway ---
		if res.Type.Name == "NATGateway" {
			hasAlloc := false
			if _, ok := res.Metadata["allocationId"]; ok {
				hasAlloc = true
			}

			if !hasAlloc {
				eipID := fmt.Sprintf("default-eip-%s", res.ID)
				var eip *resource.Resource
				if existing, ok := resMap[eipID]; ok {
					eip = existing
				} else {
					eip = &resource.Resource{
						ID:       eipID,
						Name:     fmt.Sprintf("eip-%s", res.Name),
						Type:     resource.ResourceType{ID: "ElasticIP", Name: "ElasticIP", Category: "Networking", Kind: "IP"},
						Provider: resource.AWS,
						Region:   arch.Region,
						Metadata: map[string]interface{}{
							"domain": "vpc",
						},
					}
					newResources = append(newResources, eip)
					resMap[eipID] = eip
				}

				res.Metadata["allocationId"] = eipID
				res.DependsOn = append(res.DependsOn, eipID)

				arch.Warnings = append(arch.Warnings, architecture.Warning{
					Message:    fmt.Sprintf("NAT Gateway '%s' is missing an Elastic IP. Using default Elastic IP '%s'.", res.Name, eip.Name),
					ResourceID: res.ID,
				})
			}
		}
	}

	// Append new resources
	arch.Resources = append(arch.Resources, newResources...)
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
