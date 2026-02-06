package generator

import (
	"context"
	"fmt"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/writer"
)

type Engine struct {
	mappers *tfmapper.MapperRegistry
}

func NewEngine(mappers *tfmapper.MapperRegistry) *Engine {
	return &Engine{mappers: mappers}
}

func (e *Engine) Name() string {
	return "terraform"
}

func (e *Engine) Generate(ctx context.Context, arch *architecture.Architecture, sortedResources []*resource.Resource) (*iac.Output, error) {
	_ = ctx

	if arch == nil {
		return nil, fmt.Errorf("architecture is nil")
	}
	if e.mappers == nil {
		return nil, fmt.Errorf("terraform mapper registry is nil")
	}

	provider := string(arch.Provider)
	mapper, ok := e.mappers.Get(provider)
	if !ok {
		return nil, fmt.Errorf("no terraform mapper registered for provider %q", provider)
	}

	blocks := make([]tfmapper.TerraformBlock, 0, len(sortedResources)+1)
	// if not provider block, add it explicitly ex) provider "aws" {
	//   region = var.aws_region  (if variable exists)
	// }
	if pb, ok := providerBlockWithVars(provider, arch.Region, arch.Variables); ok {
		blocks = append(blocks, pb)
	}

	// Build a lookup map for security groups: metadata "id" -> resource ID
	// This helps resolve security group references in EC2 instances
	sgIDToResourceID := make(map[string]string)
	for _, res := range arch.Resources {
		if res.Type.Name == "SecurityGroup" {
			if sgID, ok := res.Metadata["id"].(string); ok && sgID != "" {
				sgIDToResourceID[sgID] = res.ID
			}
			// Also map resource ID to itself for direct references
			sgIDToResourceID[res.ID] = res.ID
		}
	}

	// Build a map of subnet ID -> whether it has a route to Internet Gateway
	// This is determined by analyzing route table associations
	subnetHasIGWRoute := buildSubnetIGWRouteMap(arch.Resources)

	// Enrich subnet resources with route table-based public/private status
	// This allows the subnet mapper to use this information for map_public_ip_on_launch
	for _, res := range arch.Resources {
		if res.Type.Name == "Subnet" {
			if res.Metadata == nil {
				res.Metadata = make(map[string]interface{})
			}
			if hasIGW, found := subnetHasIGWRoute[res.ID]; found {
				res.Metadata["_isPublicByRouteTable"] = hasIGW
			}
		}
	}

	// Build a lookup map for subnets: resource ID -> subnet info
	// This helps determine if EC2 instances are in public/private subnets
	subnetInfo := make(map[string]subnetMetadata)
	for _, res := range arch.Resources {
		if res.Type.Name == "Subnet" {
			info := subnetMetadata{
				name: res.Name,
			}
			// Determine if subnet is public using priority:
			// 1. Route table association with Internet Gateway (source of truth)
			// 2. Explicit isPublic metadata
			// 3. map_public_ip_on_launch metadata
			// 4. Name-based detection (lowest priority, fallback)
			if hasIGW, found := subnetHasIGWRoute[res.ID]; found {
				// Route table association is the source of truth
				// Subnet is public if it has a route to Internet Gateway
				info.isPublic = hasIGW
			} else if isPublic, ok := res.Metadata["isPublic"].(bool); ok {
				info.isPublic = isPublic
			} else if mapPublicIP, ok := res.Metadata["map_public_ip_on_launch"].(bool); ok {
				info.isPublic = mapPublicIP
			} else {
				// Fallback: auto-detect based on name containing "public"
				info.isPublic = isPublicSubnetName(res.Name)
			}
			subnetInfo[res.ID] = info
		}
	}

	// Build a lookup map for all resources: ID -> Resource
	// This helps resolve dependencies
	resourceMap := make(map[string]*resource.Resource)
	resourceIDToName := make(map[string]string)
	for _, res := range arch.Resources {
		resourceMap[res.ID] = res
		// Pre-calculate terraform name for each resource
		// We use a local helper tfName (duplicated/matched from mapper) or just the Name if suitable
		// But here we just store the raw name and let mapper decide, OR we decide here.
		// User wants "provided name or default".
		// Let's store the raw Name.
		resourceIDToName[res.ID] = res.Name
	}

	// Enrich resources with explicit dependencies (DependsOn field)
	// We pass the dependency's ID and Type to the mapper via metadata
	for _, res := range arch.Resources {
		if len(res.DependsOn) > 0 {
			if res.Metadata == nil {
				res.Metadata = make(map[string]interface{})
			}
			var deps []map[string]string
			for _, depID := range res.DependsOn {
				if depRes, ok := resourceMap[depID]; ok {
					deps = append(deps, map[string]string{
						"id":   depRes.ID,
						"type": depRes.Type.Name,
						"name": depRes.Name,
					})
				}
			}
			if len(deps) > 0 {
				res.Metadata["_dependsOn"] = deps
			}
		}
	}

	for _, res := range sortedResources {
		if res == nil {
			continue
		}

		// Skip visual-only resources - they are for diagram display only, not Terraform generation
		if isVisualOnly, ok := res.Metadata["isVisualOnly"].(bool); ok && isVisualOnly {
			continue
		}

		if !mapper.SupportsResource(res.Type.Name) {
			return nil, fmt.Errorf("terraform mapper for %q does not support resource type %q (resource id %q)", provider, res.Type.Name, res.ID)
		}

		// Resolve security group IDs in EC2 resources before mapping
		if res.Type.Name == "EC2" {
			resolveSecurityGroupIDs(res, sgIDToResourceID)
			// Enrich EC2 with parent subnet info for associate_public_ip_address
			enrichEC2WithSubnetInfo(res, subnetInfo)
		}

		// Enrich with parent name for references
		if res.ParentID != nil && *res.ParentID != "" {
			if pName, ok := resourceIDToName[*res.ParentID]; ok {
				if res.Metadata == nil {
					res.Metadata = make(map[string]interface{})
				}
				res.Metadata["_parentName"] = pName
			}
		}

		// Enrich Security Groups with names
		if validSGs := getSecurityGroupRefs(res); len(validSGs) > 0 {
			sgNames := make(map[string]string)
			for _, sgID := range validSGs {
				if name, ok := resourceIDToName[sgID]; ok {
					sgNames[sgID] = name
				}
			}
			if len(sgNames) > 0 {
				if res.Metadata == nil {
					res.Metadata = make(map[string]interface{})
				}
				res.Metadata["_securityGroupNames"] = sgNames
			}
		}

		// Enrich Subnet ID references (e.g. for NAT Gateway)
		if subnetID, ok := res.Metadata["subnetId"].(string); ok && subnetID != "" {
			if name, ok := resourceIDToName[subnetID]; ok {
				res.Metadata["_subnetName"] = name
			}
		}

		bs, err := mapper.MapResource(res)
		if err != nil {
			return nil, fmt.Errorf("map resource %q (%s): %w", res.ID, res.Type.Name, err)
		}
		blocks = append(blocks, bs...)
	}

	mainTF, err := writer.RenderMainTF(blocks)
	if err != nil {
		return nil, err
	}

	out := &iac.Output{
		Files: []iac.GeneratedFile{
			{Path: "main.tf", Content: mainTF, Type: "hcl"},
		},
	}

	// Generate variables.tf if there are variables
	if len(arch.Variables) > 0 {
		vars := convertArchVariablesToTF(arch.Variables)
		varsTF, err := writer.RenderVariablesTF(vars)
		if err != nil {
			return nil, fmt.Errorf("render variables.tf: %w", err)
		}
		if varsTF != "" {
			out.Files = append(out.Files, iac.GeneratedFile{
				Path:    "variables.tf",
				Content: varsTF,
				Type:    "hcl",
			})
		}
	}

	// Generate outputs.tf if there are outputs
	if len(arch.Outputs) > 0 {
		outputs := convertArchOutputsToTF(arch.Outputs, arch.Resources)
		outputsTF, err := writer.RenderOutputsTF(outputs)
		if err != nil {
			return nil, fmt.Errorf("render outputs.tf: %w", err)
		}
		if outputsTF != "" {
			out.Files = append(out.Files, iac.GeneratedFile{
				Path:    "outputs.tf",
				Content: outputsTF,
				Type:    "hcl",
			})
		}
	}

	return out, nil
}

// providerBlockWithVars creates a provider block, using variable references if available
func providerBlockWithVars(provider, region string, variables []architecture.Variable) (tfmapper.TerraformBlock, bool) {
	if provider == "" {
		return tfmapper.TerraformBlock{}, false
	}
	attrs := map[string]tfmapper.TerraformValue{}
	if region != "" {
		// Check if there's a variable that matches the region value
		// Common variable names for region: aws_region, region
		regionVarRef := findRegionVariable(region, variables)
		if regionVarRef != "" {
			// Use variable reference
			expr := tfmapper.TerraformExpr(regionVarRef)
			attrs["region"] = tfmapper.TerraformValue{Expr: &expr}
		} else {
			r := region
			attrs["region"] = tfmapper.TerraformValue{String: &r}
		}
	}

	return tfmapper.TerraformBlock{
		Kind:       "provider",
		Labels:     []string{provider},
		Attributes: attrs,
	}, true
}

// resolveSecurityGroupIDs resolves security group IDs in EC2 metadata
// by mapping metadata "id" values to actual resource IDs
func resolveSecurityGroupIDs(res *resource.Resource, sgIDToResourceID map[string]string) {
	if res.Metadata == nil {
		return
	}

	// Handle securityGroups array
	if securityGroups, ok := res.Metadata["securityGroups"].([]interface{}); ok {
		for _, sgRaw := range securityGroups {
			if sgObj, ok := sgRaw.(map[string]interface{}); ok {
				if sgID, ok := sgObj["id"].(string); ok && sgID != "" {
					// Try to resolve the security group ID
					if resolvedID, found := sgIDToResourceID[sgID]; found {
						sgObj["id"] = resolvedID
					}
				}
			}
		}
	}
}

// subnetMetadata holds information about a subnet for EC2 enrichment
type subnetMetadata struct {
	name     string
	isPublic bool
}

// enrichEC2WithSubnetInfo adds parent subnet information to EC2 metadata
// This allows the mapper to determine associate_public_ip_address
func enrichEC2WithSubnetInfo(res *resource.Resource, subnetInfo map[string]subnetMetadata) {
	if res.Metadata == nil {
		res.Metadata = make(map[string]interface{})
	}

	if res.ParentID == nil || *res.ParentID == "" {
		return
	}

	if info, found := subnetInfo[*res.ParentID]; found {
		res.Metadata["_parentSubnetName"] = info.name
		res.Metadata["_parentSubnetIsPublic"] = info.isPublic
	}
}

// isPublicSubnetName checks if a subnet name indicates it's a public subnet
func isPublicSubnetName(name string) bool {
	nameLower := strings.ToLower(name)
	return strings.Contains(nameLower, "public")
}

// buildSubnetIGWRouteMap analyzes route tables to determine which subnets
// have routes to an Internet Gateway (making them public subnets)
func buildSubnetIGWRouteMap(resources []*resource.Resource) map[string]bool {
	result := make(map[string]bool)

	// First, find all route tables and check if they have IGW routes
	routeTableHasIGW := make(map[string]bool)
	routeTableToSubnets := make(map[string][]string)

	for _, res := range resources {
		if res.Type.Name == "RouteTable" {
			routeTableID := res.ID
			hasIGWRoute := false

			// Check routes for Internet Gateway target
			if routes, ok := res.Metadata["routes"].([]interface{}); ok {
				for _, routeRaw := range routes {
					if route, ok := routeRaw.(map[string]interface{}); ok {
						if target, ok := route["target"].(map[string]interface{}); ok {
							if targetType, ok := target["type"].(string); ok {
								if targetType == "InternetGateway" {
									hasIGWRoute = true
									break
								}
							}
						}
					}
				}
			}
			routeTableHasIGW[routeTableID] = hasIGWRoute

			// Extract subnet associations
			if associations, ok := res.Metadata["associations"].([]interface{}); ok {
				for _, assocRaw := range associations {
					if assoc, ok := assocRaw.(map[string]interface{}); ok {
						assocType, _ := assoc["associationType"].(string)
						subnetID, _ := assoc["associatedResourceId"].(string)
						if assocType == "Subnet" && subnetID != "" {
							routeTableToSubnets[routeTableID] = append(routeTableToSubnets[routeTableID], subnetID)
						}
					}
				}
			}
		}
	}

	// Map subnets to their public/private status based on route table associations
	for routeTableID, subnets := range routeTableToSubnets {
		hasIGW := routeTableHasIGW[routeTableID]
		for _, subnetID := range subnets {
			// If subnet is associated with a route table that has IGW route, it's public
			// If already marked as public by another route table, keep it public
			if existing, found := result[subnetID]; found {
				result[subnetID] = existing || hasIGW
			} else {
				result[subnetID] = hasIGW
			}
		}
	}

	return result
}

// findRegionVariable checks if there's a variable whose default matches the region value
// and returns the variable reference (var.name) if found
func findRegionVariable(region string, variables []architecture.Variable) string {
	// Common variable names for AWS region
	regionVarNames := []string{"aws_region", "region", "aws-region"}

	for _, v := range variables {
		// Check if variable name matches common region variable names
		nameLower := strings.ToLower(v.Name)
		for _, regionName := range regionVarNames {
			if nameLower == regionName {
				// Check if default value matches the region
				if defaultStr, ok := v.Default.(string); ok && defaultStr == region {
					return "var." + v.Name
				}
			}
		}
	}
	return ""
}

// convertArchVariablesToTF converts architecture variables to terraform mapper variables
func convertArchVariablesToTF(archVars []architecture.Variable) []tfmapper.Variable {
	vars := make([]tfmapper.Variable, 0, len(archVars))
	for _, v := range archVars {
		tfVar := tfmapper.Variable{
			Name:        v.Name,
			Type:        v.Type,
			Description: v.Description,
			Sensitive:   v.Sensitive,
		}
		// Convert default value if present
		if v.Default != nil {
			tfVal := convertToTerraformValue(v.Default)
			tfVar.Default = &tfVal
		}
		vars = append(vars, tfVar)
	}
	return vars
}

// convertArchOutputsToTF converts architecture outputs to terraform mapper outputs
func convertArchOutputsToTF(archOutputs []architecture.Output, resources []*resource.Resource) []tfmapper.Output {
	outputs := make([]tfmapper.Output, 0, len(archOutputs))

	// Build a map of resource ID to terraform resource name for reference resolution
	resourceIDToTFName := make(map[string]string)
	for _, res := range resources {
		resourceIDToTFName[res.ID] = tfName(res.ID)
	}

	for _, o := range archOutputs {
		// The value should be a terraform expression (e.g., "aws_vpc.vpc_2.id")
		// We need to resolve resource references if they use resource IDs
		valueExpr := resolveOutputValue(o.Value, resourceIDToTFName)

		expr := tfmapper.TerraformExpr(valueExpr)
		tfOutput := tfmapper.Output{
			Name:        o.Name,
			Value:       tfmapper.TerraformValue{Expr: &expr},
			Description: o.Description,
			Sensitive:   o.Sensitive,
		}
		outputs = append(outputs, tfOutput)
	}
	return outputs
}

// resolveOutputValue resolves resource IDs in output values to terraform resource names
func resolveOutputValue(value string, resourceIDToTFName map[string]string) string {
	// The value might be in format "resource_type.resource_id.attribute"
	// We need to convert resource_id to the terraform name (sanitized)
	parts := strings.Split(value, ".")
	if len(parts) >= 2 {
		// Check if the second part is a resource ID we know about
		if tfName, ok := resourceIDToTFName[parts[1]]; ok {
			parts[1] = tfName
			return strings.Join(parts, ".")
		}
	}
	return value
}

// convertToTerraformValue converts a Go interface{} to a TerraformValue
func convertToTerraformValue(val interface{}) tfmapper.TerraformValue {
	switch v := val.(type) {
	case string:
		return tfmapper.TerraformValue{String: &v}
	case float64:
		return tfmapper.TerraformValue{Number: &v}
	case int:
		f := float64(v)
		return tfmapper.TerraformValue{Number: &f}
	case bool:
		return tfmapper.TerraformValue{Bool: &v}
	case []interface{}:
		list := make([]tfmapper.TerraformValue, 0, len(v))
		for _, item := range v {
			list = append(list, convertToTerraformValue(item))
		}
		return tfmapper.TerraformValue{List: list}
	case map[string]interface{}:
		m := make(map[string]tfmapper.TerraformValue)
		for k, item := range v {
			m[k] = convertToTerraformValue(item)
		}
		return tfmapper.TerraformValue{Map: m}
	default:
		// Fallback to empty string
		empty := ""
		return tfmapper.TerraformValue{String: &empty}
	}
}

// tfName sanitizes an ID for use as a Terraform resource name
func tfName(id string) string {
	if id == "" {
		return "resource"
	}
	s := strings.ReplaceAll(id, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ToLower(s)
	if s == "" {
		return "resource"
	}
	// Terraform identifiers must not start with a digit
	if s[0] >= '0' && s[0] <= '9' {
		s = "r_" + s
	}
	return s
}

// getSecurityGroupRefs extracts security group IDs from metadata (generic)
func getSecurityGroupRefs(res *resource.Resource) []string {
	var ids []string

	// Check "vpc_security_group_ids" (RDS, etc)
	if list, ok := res.Metadata["vpc_security_group_ids"].([]interface{}); ok {
		for _, item := range list {
			if s, ok := item.(string); ok {
				ids = append(ids, s)
			}
		}
	}

	// Check "securityGroups" (EC2)
	// Note: resolveSecurityGroupIDs might have already run or not.
	// We handle explicit IDs.
	if sgs, ok := res.Metadata["securityGroups"].([]interface{}); ok {
		for _, item := range sgs {
			if m, ok := item.(map[string]interface{}); ok {
				if id, ok := m["id"].(string); ok {
					ids = append(ids, id)
				}
			}
		}
	}

	// Check "rules" (SecurityGroup) for "sourceSecurityGroupId"
	if rules, ok := res.Metadata["rules"].([]interface{}); ok {
		for _, ruleRaw := range rules {
			if rule, ok := ruleRaw.(map[string]interface{}); ok {
				if sourceSGID, ok := rule["sourceSecurityGroupId"].(string); ok && sourceSGID != "" {
					ids = append(ids, sourceSGID)
				}
			}
		}
	}

	return ids
}

var _ iac.Engine = (*Engine)(nil)
