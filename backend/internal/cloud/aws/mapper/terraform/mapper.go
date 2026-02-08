package terraform

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/inventory"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

// AWSMapper maps domain resources (AWS provider) into Terraform blocks.
//
// Naming strategy: Terraform local names use the domain resource ID (sanitized).
// This makes inter-resource references possible without needing global lookup.
type AWSMapper struct{}

func New() *AWSMapper {
	mapper := &AWSMapper{}
	// Initialize Terraform mapper functions in inventory
	inv := inventory.GetDefaultInventory()

	// Register specific mapper functions for each resource type
	inv.SetTerraformMapper("VPC", mapper.mapVPC)
	inv.SetTerraformMapper("Subnet", mapper.mapSubnet)
	inv.SetTerraformMapper("EC2", mapper.mapEC2)
	inv.SetTerraformMapper("SecurityGroup", mapper.mapSecurityGroup)
	inv.SetTerraformMapper("RouteTable", mapper.mapRouteTable)
	inv.SetTerraformMapper("InternetGateway", mapper.mapInternetGateway)
	inv.SetTerraformMapper("NATGateway", mapper.mapNATGateway)
	inv.SetTerraformMapper("ElasticIP", mapper.mapElasticIP)
	inv.SetTerraformMapper("S3", mapper.mapS3)
	inv.SetTerraformMapper("RDS", mapper.mapRDS)
	inv.SetTerraformMapper("AutoScalingGroup", mapper.mapAutoScalingGroup)
	inv.SetTerraformMapper("Lambda", mapper.mapLambda)
	inv.SetTerraformMapper("AvailabilityZone", mapper.mapAvailabilityZone)
	inv.SetTerraformMapper("LoadBalancer", mapper.mapLoadBalancer)
	inv.SetTerraformMapper("LaunchTemplate", mapper.mapLaunchTemplate)
	inv.SetTerraformMapper("Listener", mapper.mapListener)
	inv.SetTerraformMapper("TargetGroup", mapper.mapTargetGroup)

	return mapper
}

func (m *AWSMapper) Provider() string { return "aws" }

func (m *AWSMapper) SupportsResource(resourceType string) bool {
	inv := inventory.GetDefaultInventory()
	return inv.SupportsResource(resourceType)
}

func (m *AWSMapper) MapResource(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	if res == nil {
		return nil, fmt.Errorf("resource is nil")
	}
	if res.ID == "" {
		return nil, fmt.Errorf("resource id is empty")
	}

	// Use inventory for dynamic dispatch
	inv := inventory.GetDefaultInventory()
	functions, ok := inv.GetFunctions(res.Type.Name)
	if !ok || functions.TerraformMapper == nil {
		// Fallback to switch-based mapping for backward compatibility
		return m.mapResourceFallback(res)
	}

	// Use inventory function registry
	return functions.TerraformMapper(res)
}

// mapResourceFallback provides backward compatibility with switch-based mapping
func (m *AWSMapper) mapResourceFallback(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	switch res.Type.Name {
	case "VPC":
		return m.mapVPC(res)
	case "Subnet":
		return m.mapSubnet(res)
	case "EC2":
		return m.mapEC2(res)
	case "SecurityGroup":
		return m.mapSecurityGroup(res)
	case "RouteTable":
		return m.mapRouteTable(res)
	case "InternetGateway":
		return m.mapInternetGateway(res)
	case "NATGateway":
		return m.mapNATGateway(res)
	case "ElasticIP":
		return m.mapElasticIP(res)
	case "S3":
		return m.mapS3(res)
	case "RDS":
		return m.mapRDS(res)
	case "AutoScalingGroup":
		return m.mapAutoScalingGroup(res)
	case "AvailabilityZone":
		return m.mapAvailabilityZone(res)
	case "LoadBalancer":
		return m.mapLoadBalancer(res)
	case "LaunchTemplate":
		return m.mapLaunchTemplate(res)
	case "Listener":
		return m.mapListener(res)
	case "TargetGroup":
		return m.mapTargetGroup(res)
	default:
		return nil, fmt.Errorf("unsupported resource type %q", res.Type.Name)
	}
}

func (m *AWSMapper) mapVPC(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	cidr, ok := getString(res.Metadata, "cidr")
	if !ok || cidr == "" {
		return nil, fmt.Errorf("missing required config %q", "cidr")
	}

	attrs := map[string]tfmapper.TerraformValue{
		// Use variable reference if available, otherwise use resolved value
		"cidr_block": tfStringOrVar(res.Metadata, "cidr", cidr),
		"tags":       tfTags(res.Name),
	}

	// Enable DNS support and hostnames by default (CloudCanvas rule)
	// Only disable if user explicitly sets them to false
	if v, ok := getBool(res.Metadata, "enable_dns_hostnames"); ok {
		attrs["enable_dns_hostnames"] = tfBoolOrVar(res.Metadata, "enable_dns_hostnames", v)
	} else if v, ok := getBool(res.Metadata, "enableDnsHostnames"); ok {
		attrs["enable_dns_hostnames"] = tfBoolOrVar(res.Metadata, "enableDnsHostnames", v)
	} else {
		// Default: enable DNS hostnames
		attrs["enable_dns_hostnames"] = tfBool(true)
	}
	if v, ok := getBool(res.Metadata, "enable_dns_support"); ok {
		attrs["enable_dns_support"] = tfBoolOrVar(res.Metadata, "enable_dns_support", v)
	} else if v, ok := getBool(res.Metadata, "enableDnsSupport"); ok {
		attrs["enable_dns_support"] = tfBoolOrVar(res.Metadata, "enableDnsSupport", v)
	} else {
		// Default: enable DNS support
		attrs["enable_dns_support"] = tfBool(true)
	}
	if v, ok := getString(res.Metadata, "instance_tenancy"); ok && v != "" {
		attrs["instance_tenancy"] = tfStringOrVar(res.Metadata, "instance_tenancy", v)
	} else if v, ok := getString(res.Metadata, "tenancy"); ok && v != "" {
		attrs["instance_tenancy"] = tfStringOrVar(res.Metadata, "tenancy", v)
	}

	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_vpc", tfBlockName(res)},
			Attributes: attrs,
		},
	}, nil
}

func (m *AWSMapper) mapSubnet(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	cidr, ok := getString(res.Metadata, "cidr")
	if !ok || cidr == "" {
		return nil, fmt.Errorf("missing required config %q", "cidr")
	}
	az, ok := getString(res.Metadata, "availabilityZoneId")
	if !ok || az == "" {
		return nil, fmt.Errorf("missing required config %q", "availabilityZoneId")
	}
	if res.ParentID == nil || *res.ParentID == "" {
		return nil, fmt.Errorf("subnet requires parent vpc (parentID missing)")
	}

	attrs := map[string]tfmapper.TerraformValue{
		"vpc_id":            tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: resolveRef(*res.ParentID, res.Metadata), Attribute: "id"}.Expr()),
		"cidr_block":        tfStringOrVar(res.Metadata, "cidr", cidr),
		"availability_zone": tfStringOrVar(res.Metadata, "availabilityZoneId", az),
		"tags":              tfTags(res.Name),
	}

	// Determine if subnet is public and set map_public_ip_on_launch accordingly
	// Priority:
	// 1. Explicit map_public_ip_on_launch in metadata (highest priority - user override)
	// 2. Route table association with Internet Gateway (_isPublicByRouteTable set by generator)
	//    This is the actual source of truth - if subnet routes to IGW, it's public
	// 3. Explicit isPublic metadata (fallback)
	// 4. Name-based detection (lowest priority, fallback)
	if v, ok := getBool(res.Metadata, "map_public_ip_on_launch"); ok {
		attrs["map_public_ip_on_launch"] = tfBoolOrVar(res.Metadata, "map_public_ip_on_launch", v)
	} else if v, ok := getBool(res.Metadata, "_isPublicByRouteTable"); ok {
		// Route table association is the source of truth for public/private
		// If subnet has route to IGW, it's public; otherwise it's private
		attrs["map_public_ip_on_launch"] = tfBool(v)
	} else if v, ok := getBool(res.Metadata, "isPublic"); ok {
		attrs["map_public_ip_on_launch"] = tfBoolOrVar(res.Metadata, "isPublic", v)
	} else if isPublicSubnet(res.Name) {
		// Fallback: auto-detect public subnet based on name containing "public"
		attrs["map_public_ip_on_launch"] = tfBool(true)
	}
	// Note: if none of the above conditions are met, map_public_ip_on_launch is not set
	// which means it defaults to false (private subnet behavior)

	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_subnet", tfBlockName(res)},
			Attributes: attrs,
		},
	}, nil
}

// isPublicSubnet checks if a subnet name indicates it's a public subnet
func isPublicSubnet(name string) bool {
	nameLower := strings.ToLower(name)
	return strings.Contains(nameLower, "public")
}

func (m *AWSMapper) mapEC2(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	ami, ok := getString(res.Metadata, "ami")
	if !ok || ami == "" {
		return nil, fmt.Errorf("missing required config %q", "ami")
	}
	instanceType, ok := getString(res.Metadata, "instanceType")
	if !ok || instanceType == "" {
		return nil, fmt.Errorf("missing required config %q", "instanceType")
	}
	if res.ParentID == nil || *res.ParentID == "" {
		return nil, fmt.Errorf("ec2 requires parent subnet (parentID missing)")
	}

	attrs := map[string]tfmapper.TerraformValue{
		"ami":           tfStringOrVar(res.Metadata, "ami", ami),
		"instance_type": tfStringOrVar(res.Metadata, "instanceType", instanceType),
		"subnet_id":     tfExpr(tfmapper.Reference{ResourceType: "aws_subnet", ResourceName: resolveRef(*res.ParentID, res.Metadata), Attribute: "id"}.Expr()),
		"tags":          tfTags(res.Name),
	}

	// Handle associate_public_ip_address based on subnet type
	// Check if explicitly set in metadata first
	if v, ok := getBool(res.Metadata, "associate_public_ip_address"); ok {
		attrs["associate_public_ip_address"] = tfBoolOrVar(res.Metadata, "associate_public_ip_address", v)
	} else {
		// Determine based on parent subnet info (set by generator) or subnet name
		isInPublicSubnet := false
		if v, ok := getBool(res.Metadata, "_parentSubnetIsPublic"); ok {
			isInPublicSubnet = v
		} else if parentSubnetName, ok := getString(res.Metadata, "_parentSubnetName"); ok {
			isInPublicSubnet = isPublicSubnet(parentSubnetName)
		}

		// Always explicitly set associate_public_ip_address
		// true for public subnets, false for private subnets
		attrs["associate_public_ip_address"] = tfBool(isInPublicSubnet)
	}

	if v, ok := getString(res.Metadata, "keyName"); ok && v != "" {
		attrs["key_name"] = tfStringOrVar(res.Metadata, "keyName", v)
	}
	if v, ok := getString(res.Metadata, "iamInstanceProfile"); ok && v != "" {
		attrs["iam_instance_profile"] = tfString(v)
	}
	if v, ok := getString(res.Metadata, "userData"); ok && v != "" {
		attrs["user_data"] = tfString(v)
	}

	// Handle subnet
	if res.ParentID != nil && *res.ParentID != "" {
		attrs["subnet_id"] = tfExpr(tfmapper.Reference{ResourceType: "aws_subnet", ResourceName: resolveRef(*res.ParentID, res.Metadata), Attribute: "id"}.Expr())
	}

	// Handle security groups
	var sgIDs []string
	if list, ok := getArray(res.Metadata, "securityGroups"); ok {
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				if id, ok := m["id"].(string); ok {
					sgIDs = append(sgIDs, id)
				}
			}
		}
	}
	// Fallback to securityGroupIds for backward compatibility
	if len(sgIDs) == 0 {
		if list, ok := getStringSlice(res.Metadata, "securityGroupIds"); ok {
			sgIDs = list
		}
	}

	// Convert security group IDs to Terraform references
	// We need to find the security group resource by ID and reference it
	if len(sgIDs) > 0 {
		sgVals := make([]tfmapper.TerraformValue, 0, len(sgIDs))
		for _, id := range sgIDs {
			if id == "" {
				continue
			}
			// Reference the security group resource
			sgVals = append(sgVals, tfExpr(tfmapper.Reference{
				ResourceType: "aws_security_group",
				ResourceName: resolveRef(id, res.Metadata),
				Attribute:    "id",
			}.Expr()))
		}
		if len(sgVals) > 0 {
			attrs["vpc_security_group_ids"] = tfList(sgVals)
		}
	}

	addDependsOn(attrs, res)

	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_instance", tfBlockName(res)},
			Attributes: attrs,
		},
	}, nil
}

func (m *AWSMapper) mapSecurityGroup(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	// SGs are modeled under VPC; expect ParentID -> VPC.
	if res.ParentID == nil || *res.ParentID == "" {
		return nil, fmt.Errorf("security-group requires parent vpc (parentID missing)")
	}

	desc, _ := getString(res.Metadata, "description")
	if desc == "" {
		desc = "managed by arch-visualizer"
	}

	attrs := map[string]tfmapper.TerraformValue{
		"name":        tfString(res.Name),
		"description": tfString(desc),
		"vpc_id":      tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: resolveRef(*res.ParentID, res.Metadata), Attribute: "id"}.Expr()),
		"tags":        tfTags(res.Name),
	}

	// Initialize nested blocks for inline ingress/egress rules
	nestedBlocks := make(map[string][]tfmapper.NestedBlock)
	var ingressBlocks []tfmapper.NestedBlock
	var egressBlocks []tfmapper.NestedBlock
	hasEgressRule := false

	// Try parsing "ingressRules" and "egressRules" arrays as well
	if ingressRules, ok := getArray(res.Metadata, "ingressRules"); ok {
		for _, ruleRaw := range ingressRules {
			if rule, ok := ruleRaw.(map[string]interface{}); ok {
				rule["type"] = "ingress" // Ensure type is set
				// Append to "rules" to reuse parsing logic below?
				// Better to just call a helper or duplicate loop.
				// For valid JSON, we should just process them.
				// Let's adapt the loop below to handle a combined list.
			}
		}
	}

	// Collect all rule sources
	var allRules []interface{}
	if rules, ok := getArray(res.Metadata, "rules"); ok {
		allRules = append(allRules, rules...)
	}
	if rules, ok := getArray(res.Metadata, "ingressRules"); ok {
		for _, r := range rules {
			if m, ok := r.(map[string]interface{}); ok {
				m["type"] = "ingress"
				allRules = append(allRules, m)
			}
		}
	}
	if rules, ok := getArray(res.Metadata, "egressRules"); ok {
		for _, r := range rules {
			if m, ok := r.(map[string]interface{}); ok {
				m["type"] = "egress"
				allRules = append(allRules, m)
			}
		}
	}

	// Parse and add inline ingress/egress rules
	for _, ruleRaw := range allRules {
		rule, ok := ruleRaw.(map[string]interface{})
		if !ok {
			continue
		}

		ruleType, _ := getString(rule, "type")
		protocol, _ := getString(rule, "protocol")
		portRange, _ := getString(rule, "portRange")
		// Check for specific ports if range not present
		fromPortVal, hasFrom := getInt(rule, "fromPort")
		toPortVal, hasTo := getInt(rule, "toPort")

		cidr, _ := getString(rule, "cidr")
		description, _ := getString(rule, "description")

		if ruleType == "" || protocol == "" {
			continue
		}

		ruleAttrs := map[string]tfmapper.TerraformValue{
			"protocol": tfString(protocol),
		}

		// Handle Ports
		if portRange != "" {
			fromPort, toPort := parsePortRange(portRange)
			if fromPort != nil {
				ruleAttrs["from_port"] = tfNumber(float64(*fromPort))
			}
			if toPort != nil {
				ruleAttrs["to_port"] = tfNumber(float64(*toPort))
			}
		} else if hasFrom && hasTo {
			ruleAttrs["from_port"] = tfNumber(float64(fromPortVal))
			ruleAttrs["to_port"] = tfNumber(float64(toPortVal))
		}

		// Add description if provided
		if description != "" {
			ruleAttrs["description"] = tfString(description)
		}

		// For protocol "-1" (all), from_port and to_port must be 0
		if protocol == "-1" {
			ruleAttrs["from_port"] = tfNumber(0)
			ruleAttrs["to_port"] = tfNumber(0)
		}

		if cidr != "" {
			ruleAttrs["cidr_blocks"] = tfList([]tfmapper.TerraformValue{tfString(cidr)})
		}

		// Handle source security group reference
		if sourceSGID, ok := getString(rule, "sourceSecurityGroupId"); ok && sourceSGID != "" {
			ruleAttrs["security_groups"] = tfList([]tfmapper.TerraformValue{
				tfExpr(tfmapper.Reference{
					ResourceType: "aws_security_group",
					ResourceName: resolveRef(sourceSGID, res.Metadata),
					Attribute:    "id",
				}.Expr()),
			})
		}

		nestedBlock := tfmapper.NestedBlock{Attributes: ruleAttrs}

		if ruleType == "ingress" {
			ingressBlocks = append(ingressBlocks, nestedBlock)
		} else if ruleType == "egress" {
			egressBlocks = append(egressBlocks, nestedBlock)
			hasEgressRule = true
		}
	}

	// Add default egress rule (allow all outbound) if no egress rules are defined
	// This is AWS best practice - be explicit about egress rules
	if !hasEgressRule {
		defaultEgress := tfmapper.NestedBlock{
			Attributes: map[string]tfmapper.TerraformValue{
				"from_port":   tfNumber(0),
				"to_port":     tfNumber(0),
				"protocol":    tfString("-1"),
				"cidr_blocks": tfList([]tfmapper.TerraformValue{tfString("0.0.0.0/0")}),
				"description": tfString("Allow all outbound traffic"),
			},
		}
		egressBlocks = append(egressBlocks, defaultEgress)
	}

	// Add nested blocks to the map
	if len(ingressBlocks) > 0 {
		nestedBlocks["ingress"] = ingressBlocks
	}
	if len(egressBlocks) > 0 {
		nestedBlocks["egress"] = egressBlocks
	}

	blocks := []tfmapper.TerraformBlock{
		{
			Kind:         "resource",
			Labels:       []string{"aws_security_group", tfBlockName(res)},
			Attributes:   attrs,
			NestedBlocks: nestedBlocks,
		},
	}

	return blocks, nil
}

func (m *AWSMapper) mapRouteTable(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	if res.ParentID == nil || *res.ParentID == "" {
		return nil, fmt.Errorf("route-table requires parent vpc (parentID missing)")
	}

	// Skip main route table - it's automatically created by AWS for each VPC
	// Only generate custom route tables
	if isMain, ok := getBool(res.Metadata, "isMain"); ok && isMain {
		return []tfmapper.TerraformBlock{}, nil
	}

	blocks := []tfmapper.TerraformBlock{
		{
			Kind:   "resource",
			Labels: []string{"aws_route_table", tfBlockName(res)},
			Attributes: map[string]tfmapper.TerraformValue{
				"vpc_id": tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: resolveRef(*res.ParentID, res.Metadata), Attribute: "id"}.Expr()),
				"tags":   tfTags(res.Name),
			},
		},
	}

	// Parse and add routes
	if routes, ok := getArray(res.Metadata, "routes"); ok {
		for i, routeRaw := range routes {
			route, ok := routeRaw.(map[string]interface{})
			if !ok {
				continue
			}

			destObj, _ := route["destination"].(map[string]interface{})
			targetObj, _ := route["target"].(map[string]interface{})

			if destObj == nil || targetObj == nil {
				continue
			}

			destCIDR, _ := getString(destObj, "cidr")
			isLocal, _ := getBool(destObj, "isLocal")
			targetType, _ := getString(targetObj, "type")
			targetResourceID, _ := getString(targetObj, "resourceId")

			if destCIDR == "" {
				continue
			}

			// Skip local routes (they're automatically created by AWS)
			if isLocal {
				continue
			}

			routeAttrs := map[string]tfmapper.TerraformValue{
				"route_table_id":         tfExpr(tfmapper.Reference{ResourceType: "aws_route_table", ResourceName: tfBlockName(res), Attribute: "id"}.Expr()),
				"destination_cidr_block": tfString(destCIDR),
			}

			// Map target type to Terraform resource
			switch targetType {
			case "InternetGateway":
				if targetResourceID != "" {
					routeAttrs["gateway_id"] = tfExpr(tfmapper.Reference{
						ResourceType: "aws_internet_gateway",
						ResourceName: resolveRef(targetResourceID, res.Metadata),
						Attribute:    "id",
					}.Expr())
				}
			case "NATGateway":
				if targetResourceID != "" {
					routeAttrs["nat_gateway_id"] = tfExpr(tfmapper.Reference{
						ResourceType: "aws_nat_gateway",
						ResourceName: resolveRef(targetResourceID, res.Metadata),
						Attribute:    "id",
					}.Expr())
				}
			case "Local":
				// Local routes are handled automatically, skip
				continue
			default:
				// Unknown target type, skip
				continue
			}

			blocks = append(blocks, tfmapper.TerraformBlock{
				Kind:       "resource",
				Labels:     []string{"aws_route", fmt.Sprintf("%s_route_%d", tfBlockName(res), i)},
				Attributes: routeAttrs,
			})
		}
	}

	// Parse and add route table associations
	if associations, ok := getArray(res.Metadata, "associations"); ok {
		for i, assocRaw := range associations {
			assoc, ok := assocRaw.(map[string]interface{})
			if !ok {
				continue
			}

			assocType, _ := getString(assoc, "associationType")
			associatedResourceID, _ := getString(assoc, "associatedResourceId")

			if associatedResourceID == "" {
				continue
			}

			assocAttrs := map[string]tfmapper.TerraformValue{
				"route_table_id": tfExpr(tfmapper.Reference{ResourceType: "aws_route_table", ResourceName: tfBlockName(res), Attribute: "id"}.Expr()),
			}

			// Map association type to Terraform resource
			switch assocType {
			case "Subnet":
				assocAttrs["subnet_id"] = tfExpr(tfmapper.Reference{
					ResourceType: "aws_subnet",
					ResourceName: resolveRef(associatedResourceID, res.Metadata),
					Attribute:    "id",
				}.Expr())
			case "Gateway":
				// Gateway associations are typically handled via routes, skip for now
				continue
			default:
				// Unknown association type, skip
				continue
			}

			blocks = append(blocks, tfmapper.TerraformBlock{
				Kind:       "resource",
				Labels:     []string{"aws_route_table_association", fmt.Sprintf("%s_assoc_%d", tfBlockName(res), i)},
				Attributes: assocAttrs,
			})
		}
	}

	return blocks, nil
}

func (m *AWSMapper) mapInternetGateway(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	if res.ParentID == nil || *res.ParentID == "" {
		return nil, fmt.Errorf("internet-gateway requires parent vpc (parentID missing)")
	}

	attrs := map[string]tfmapper.TerraformValue{
		"vpc_id": tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: resolveRef(*res.ParentID, res.Metadata), Attribute: "id"}.Expr()),
		"tags":   tfTags(res.Name),
	}
	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_internet_gateway", tfBlockName(res)},
			Attributes: attrs,
		},
	}, nil
}

func (m *AWSMapper) mapElasticIP(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	domain, _ := getString(res.Metadata, "domain")
	attrs := map[string]tfmapper.TerraformValue{
		"tags": tfTags(res.Name),
	}
	if domain != "" {
		attrs["domain"] = tfString(domain)
	}
	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_eip", tfBlockName(res)},
			Attributes: attrs,
		},
	}, nil
}

func (m *AWSMapper) mapAvailabilityZone(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	// Implicitly used, no explicit Terraform code needed
	return []tfmapper.TerraformBlock{}, nil
}

func (m *AWSMapper) mapLoadBalancer(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	// Determine LB type (application or network)
	lbType, _ := getString(res.Metadata, "load_balancer_type")
	if lbType == "" {
		lbType = "application" // Default
	}

	internal, _ := getBool(res.Metadata, "internal")

	attrs := map[string]tfmapper.TerraformValue{
		"name":               tfStringOrVar(res.Metadata, "name", res.Name),
		"internal":           tfBool(internal),
		"load_balancer_type": tfString(lbType),
		"tags":               tfTags(res.Name),
	}

	// Security Groups (Application LBs only)
	if lbType == "application" {
		var sgIDs []string
		if list, ok := getArray(res.Metadata, "securityGroups"); ok {
			for _, item := range list {
				if m, ok := item.(map[string]interface{}); ok {
					if id, ok := m["id"].(string); ok {
						sgIDs = append(sgIDs, id)
					}
				}
			}
		}
		// Fallback to securityGroupIds
		if len(sgIDs) == 0 {
			if list, ok := getStringSlice(res.Metadata, "securityGroupIds"); ok {
				sgIDs = list
			}
		}

		if len(sgIDs) > 0 {
			sgVals := make([]tfmapper.TerraformValue, 0, len(sgIDs))
			for _, id := range sgIDs {
				if id == "" {
					continue
				}
				sgVals = append(sgVals, tfExpr(tfmapper.Reference{
					ResourceType: "aws_security_group",
					ResourceName: resolveRef(id, res.Metadata),
					Attribute:    "id",
				}.Expr()))
			}
			attrs["security_groups"] = tfList(sgVals)
		}
	}

	// Subnets
	var subnetIDs []string
	if subnetsRaw, ok := getArray(res.Metadata, "subnets"); ok {
		for _, sRaw := range subnetsRaw {
			if sMap, ok := sRaw.(map[string]interface{}); ok {
				if sid, ok := sMap["subnetId"].(string); ok {
					subnetIDs = append(subnetIDs, sid)
				}
			}
		}
	}

	if len(subnetIDs) > 0 {
		subnetRefs := make([]tfmapper.TerraformValue, 0, len(subnetIDs))
		for _, sid := range subnetIDs {
			subnetRefs = append(subnetRefs, tfExpr(tfmapper.Reference{
				ResourceType: "aws_subnet",
				ResourceName: resolveRef(sid, res.Metadata),
				Attribute:    "id",
			}.Expr()))
		}
		attrs["subnets"] = tfList(subnetRefs)
	}

	// Enable deletion protection
	if v, ok := getBool(res.Metadata, "enable_deletion_protection"); ok {
		attrs["enable_deletion_protection"] = tfBool(v)
	}

	addDependsOn(attrs, res)

	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_lb", tfBlockName(res)},
			Attributes: attrs,
		},
	}, nil
}

func (m *AWSMapper) mapNATGateway(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	// NAT GW must be in subnet; allow both parent subnet or explicit subnetId.
	subnetID, _ := getString(res.Metadata, "subnetId")
	if subnetID == "" && res.ParentID != nil {
		subnetID = *res.ParentID
	}
	if subnetID == "" {
		return nil, fmt.Errorf("nat-gateway requires subnetId (or parent subnet)")
	}

	blocks := []tfmapper.TerraformBlock{}

	natGatewayAttrs := map[string]tfmapper.TerraformValue{
		"subnet_id": tfExpr(tfmapper.Reference{ResourceType: "aws_subnet", ResourceName: resolveRef(subnetID, res.Metadata), Attribute: "id"}.Expr()),
		"tags":      tfTags(res.Name),
	}

	// Check if user provided an allocation_id
	if alloc, ok := getString(res.Metadata, "allocationId"); ok && alloc != "" {
		// User provided allocation_id, use it directly
		natGatewayAttrs["allocation_id"] = tfString(alloc)
	} else {
		// No allocation_id provided, create an EIP automatically
		eipName := fmt.Sprintf("%s_eip", tfName(res.ID))
		eipBlock := tfmapper.TerraformBlock{
			Kind:   "resource",
			Labels: []string{"aws_eip", eipName},
			Attributes: map[string]tfmapper.TerraformValue{
				"domain": tfString("vpc"),
				"tags":   tfTags(res.Name + "-eip"),
			},
		}
		blocks = append(blocks, eipBlock)

		// Reference the auto-created EIP
		natGatewayAttrs["allocation_id"] = tfExpr(tfmapper.Reference{
			ResourceType: "aws_eip",
			ResourceName: eipName,
			Attribute:    "id",
		}.Expr())
	}

	// Add NAT Gateway block
	blocks = append(blocks, tfmapper.TerraformBlock{
		Kind:       "resource",
		Labels:     []string{"aws_nat_gateway", tfBlockName(res)},
		Attributes: natGatewayAttrs,
	})

	return blocks, nil
}

func (m *AWSMapper) mapS3(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	// S3 bucket name must be globally unique; we expect the user-provided name.
	bucket := res.Name
	if n, ok := getString(res.Metadata, "name"); ok && n != "" {
		bucket = n
	}
	bucket = sanitizeS3BucketName(bucket)
	if bucket == "" {
		return nil, fmt.Errorf("invalid s3 bucket name")
	}

	attrs := map[string]tfmapper.TerraformValue{
		"bucket": tfString(bucket),
		"tags":   tfTags(res.Name),
	}
	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_s3_bucket", tfName(res.Name)},
			Attributes: attrs,
		},
	}, nil
}

var _ tfmapper.ResourceMapper = (*AWSMapper)(nil)

var tfNameSanitizer = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

func tfName(id string) string {
	if id == "" {
		return "resource"
	}
	s := tfNameSanitizer.ReplaceAllString(id, "_")
	s = strings.Trim(s, "_")
	s = strings.ToLower(s)
	if s == "" {
		return "resource"
	}
	// Terraform identifiers must not start with a digit.
	if s[0] >= '0' && s[0] <= '9' {
		s = "r_" + s
	}
	return s
}

// tfBlockName returns the terraform local name for the resource choice: Name > ID
// DEPRECATED: This implementation often causes collisions.
// Replaced by implementation at end of file.
// OLD IMPLEMENTATION DELETED

func sanitizeS3BucketName(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	// remove invalid chars
	allowed := make([]rune, 0, len(s))
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '.' {
			allowed = append(allowed, r)
		}
	}
	s = string(allowed)
	s = strings.Trim(s, "-.")
	if len(s) < 3 || len(s) > 63 {
		return ""
	}
	return s
}

func tfString(s string) tfmapper.TerraformValue {
	return tfmapper.TerraformValue{String: &s}
}

// tfStringOrVar returns a TerraformValue that uses a variable reference if one exists,
// otherwise uses the resolved string value.
// fieldKey is the key in metadata to check (e.g., "cidr", "name")
// metadata is the resource metadata containing _varRefs
// resolvedValue is the already-resolved value to use if no var ref exists
func tfStringOrVar(metadata map[string]interface{}, fieldKey string, resolvedValue string) tfmapper.TerraformValue {
	if varRef := getVarRef(metadata, fieldKey); varRef != "" {
		// Use the variable reference as an expression
		return tfExpr(tfmapper.TerraformExpr(varRef))
	}
	return tfString(resolvedValue)
}

// tfBoolOrVar returns a TerraformValue that uses a variable reference if one exists,
// otherwise uses the resolved bool value.
func tfBoolOrVar(metadata map[string]interface{}, fieldKey string, resolvedValue bool) tfmapper.TerraformValue {
	if varRef := getVarRef(metadata, fieldKey); varRef != "" {
		// Use the variable reference as an expression
		return tfExpr(tfmapper.TerraformExpr(varRef))
	}
	return tfBool(resolvedValue)
}

// getVarRef retrieves the original variable reference for a field from _varRefs metadata.
// Returns empty string if no variable reference exists for that field.
func getVarRef(metadata map[string]interface{}, fieldKey string) string {
	if metadata == nil {
		return ""
	}
	varRefs, ok := metadata["_varRefs"].(map[string]interface{})
	if !ok {
		// Try as map[string]string
		if varRefsStr, ok := metadata["_varRefs"].(map[string]string); ok {
			return varRefsStr[fieldKey]
		}
		return ""
	}
	if ref, ok := varRefs[fieldKey].(string); ok {
		return ref
	}
	return ""
}

// hasVarRef checks if a field has a variable reference in metadata
func hasVarRefForField(metadata map[string]interface{}, fieldKey string) bool {
	return getVarRef(metadata, fieldKey) != ""
}

func tfBool(b bool) tfmapper.TerraformValue {
	return tfmapper.TerraformValue{Bool: &b}
}

func tfExpr(e tfmapper.TerraformExpr) tfmapper.TerraformValue {
	return tfmapper.TerraformValue{Expr: &e}
}

func tfList(items []tfmapper.TerraformValue) tfmapper.TerraformValue {
	return tfmapper.TerraformValue{List: items}
}

func tfTags(name string) tfmapper.TerraformValue {
	if name == "" {
		name = "arch-resource"
	}
	return tfmapper.TerraformValue{
		Map: map[string]tfmapper.TerraformValue{
			"Name": tfString(name),
		},
	}
}

func getString(m map[string]interface{}, key string) (string, bool) {
	if m == nil {
		return "", false
	}
	v, ok := m[key]
	if !ok || v == nil {
		return "", false
	}
	switch t := v.(type) {
	case string:
		return t, true
	default:
		return "", false
	}
}

func getBool(m map[string]interface{}, key string) (bool, bool) {
	if m == nil {
		return false, false
	}
	v, ok := m[key]
	if !ok || v == nil {
		return false, false
	}
	switch t := v.(type) {
	case bool:
		return t, true
	default:
		return false, false
	}
}

func getStringSlice(m map[string]interface{}, key string) ([]string, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m[key]
	if !ok || v == nil {
		return nil, false
	}
	switch t := v.(type) {
	case []string:
		return t, true
	case []interface{}:
		out := make([]string, 0, len(t))
		for _, item := range t {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out, true
	default:
		return nil, false
	}
}

func getArray(m map[string]interface{}, key string) ([]interface{}, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m[key]
	if !ok || v == nil {
		return nil, false
	}
	switch t := v.(type) {
	case []interface{}:
		return t, true
	default:
		return nil, false
	}
}

func tfNumber(n float64) tfmapper.TerraformValue {
	return tfmapper.TerraformValue{Number: &n}
}

// parsePortRange parses a port range string like "22" or "80-443" into from/to ports
func parsePortRange(portRange string) (fromPort *int, toPort *int) {
	if portRange == "" {
		return nil, nil
	}

	// Handle single port
	if !strings.Contains(portRange, "-") {
		if port, err := parseInt(portRange); err == nil {
			fromPort = &port
			toPort = &port
			return
		}
		return nil, nil
	}

	// Handle port range
	parts := strings.Split(portRange, "-")
	if len(parts) == 2 {
		if from, err := parseInt(parts[0]); err == nil {
			fromPort = &from
		}
		if to, err := parseInt(parts[1]); err == nil {
			toPort = &to
		}
	}
	return
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &result)
	return result, err
}

func (m *AWSMapper) mapLambda(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	runtime, ok := getString(res.Metadata, "runtime")
	if !ok || runtime == "" {
		return nil, fmt.Errorf("lambda requires runtime")
	}
	handler, ok := getString(res.Metadata, "handler")
	if !ok || handler == "" {
		return nil, fmt.Errorf("lambda requires handler")
	}
	role, ok := getString(res.Metadata, "role")
	if !ok || role == "" {
		return nil, fmt.Errorf("lambda requires role")
	}

	attrs := map[string]tfmapper.TerraformValue{
		"function_name": tfString(res.Name),
		"role":          tfString(role),
		"handler":       tfString(handler),
		"runtime":       tfString(runtime),
		"filename":      tfString("function.zip"), // Placeholder
		"tags":          tfTags(res.Name),
	}

	if mem, ok := res.Metadata["memory"].(float64); ok {
		attrs["memory_size"] = tfNumber(mem)
	}

	addDependsOn(attrs, res)

	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_lambda_function", tfName(res.Name)},
			Attributes: attrs,
		},
	}, nil
}

func (m *AWSMapper) mapRDS(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	// Check required fields
	instanceClass, _ := getString(res.Metadata, "instance_class")
	if instanceClass == "" {
		instanceClass = "db.t3.micro" // Default
	}
	engine, _ := getString(res.Metadata, "engine")
	if engine == "" {
		return nil, fmt.Errorf("rds requires engine")
	}

	blocks := []tfmapper.TerraformBlock{}

	// DB Subnet Group
	// Auto-create a subnet group using all subnets in the parent VPC (or user-defined logic)
	// For now, let's create a subnet group for this DB if not provided
	subnetGroupName := fmt.Sprintf("%s_subnet_group", tfBlockName(res))
	if v, ok := getString(res.Metadata, "db_subnet_group_name"); ok && v != "" {
		// If user provided a name, use it for the resource name attribute
		// But we still need a sanitized Terraform identifier
		subnetGroupName = tfName(v)
	}

	// Calculate subnets from metadata
	var subnetIDs []string
	if list, ok := getArray(res.Metadata, "subnets"); ok {
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				if id, ok := m["subnetId"].(string); ok {
					subnetIDs = append(subnetIDs, id)
				}
			}
		}
	} else if list, ok := getArray(res.Metadata, "subnetIds"); ok {
		// Fallback to string array
		for _, item := range list {
			if s, ok := item.(string); ok {
				subnetIDs = append(subnetIDs, s)
			}
		}
	}

	// If no subnets specified, use private subnets injected by generator
	// RDS instances should always use private subnets
	if len(subnetIDs) == 0 {
		if privateSubnets, ok := res.Metadata["_privateSubnetIDs"].([]string); ok && len(privateSubnets) >= 2 {
			// AWS requires at least 2 subnets for DB subnet groups
			subnetIDs = privateSubnets
		}
	}

	// Always try to create the subnet group if we have subnets
	// This ensures "Required generator output" is present
	if len(subnetIDs) > 0 {
		subnetRefs := make([]tfmapper.TerraformValue, 0, len(subnetIDs))
		for _, sid := range subnetIDs {
			subnetRefs = append(subnetRefs, tfExpr(tfmapper.Reference{
				ResourceType: "aws_subnet",
				ResourceName: resolveRef(sid, res.Metadata), // Use resolveRef to properly resolve frontend node IDs
				Attribute:    "id",
			}.Expr()))
		}

		blocks = append(blocks, tfmapper.TerraformBlock{
			Kind:   "resource",
			Labels: []string{"aws_db_subnet_group", subnetGroupName},
			Attributes: map[string]tfmapper.TerraformValue{
				"name":       tfString(subnetGroupName),
				"subnet_ids": tfList(subnetRefs),
				"tags":       tfTags(res.Name + " subnet group"),
			},
		})
	}

	attrs := map[string]tfmapper.TerraformValue{
		"identifier":          tfString(tfName(res.ID)),
		"instance_class":      tfString(instanceClass),
		"engine":              tfString(engine),
		"tags":                tfTags(res.Name),
		"skip_final_snapshot": tfBool(true), // Avoid hanging on destroy in dev environments
		"db_subnet_group_name": tfExpr(tfmapper.Reference{
			ResourceType: "aws_db_subnet_group",
			ResourceName: subnetGroupName,
			Attribute:    "name",
		}.Expr()),
	}

	if v, ok := getString(res.Metadata, "engine_version"); ok && v != "" {
		attrs["engine_version"] = tfString(v)
	}
	if v, ok := getInt(res.Metadata, "allocated_storage"); ok {
		attrs["allocated_storage"] = tfNumber(float64(v))
	} else if v, ok := getFloat(res.Metadata, "allocated_storage"); ok {
		attrs["allocated_storage"] = tfNumber(v)
	} else {
		// Default storage if not provided
		attrs["allocated_storage"] = tfNumber(20)
	}

	if v, ok := getString(res.Metadata, "db_name"); ok && v != "" {
		attrs["db_name"] = tfString(v)
	}
	if v, ok := getString(res.Metadata, "username"); ok && v != "" {
		attrs["username"] = tfString(v)
	}
	if v, ok := getString(res.Metadata, "password"); ok && v != "" {
		attrs["password"] = tfString(v)
	}

	if v, ok := getBool(res.Metadata, "multi_az"); ok {
		attrs["multi_az"] = tfBool(v)
	}
	if v, ok := getBool(res.Metadata, "publicly_accessible"); ok {
		attrs["publicly_accessible"] = tfBool(v)
	}
	if v, ok := getInt(res.Metadata, "backup_retention_period"); ok {
		attrs["backup_retention_period"] = tfNumber(float64(v))
	}

	// Handle security groups
	var sgIDs []string
	if list, ok := getArray(res.Metadata, "vpc_security_group_ids"); ok {
		for _, item := range list {
			if s, ok := item.(string); ok {
				sgIDs = append(sgIDs, s)
			}
		}
	}
	// Also check securityGroups (object array) or securityGroupIds (string array)
	if len(sgIDs) == 0 {
		if list, ok := getArray(res.Metadata, "securityGroups"); ok {
			for _, item := range list {
				if m, ok := item.(map[string]interface{}); ok {
					if id, ok := m["id"].(string); ok {
						sgIDs = append(sgIDs, id)
					}
				}
			}
		}
	}
	if len(sgIDs) == 0 {
		if list, ok := getStringSlice(res.Metadata, "securityGroupIds"); ok {
			sgIDs = list
		}
	}

	if len(sgIDs) > 0 {
		sgVals := make([]tfmapper.TerraformValue, 0, len(sgIDs))
		for _, id := range sgIDs {
			sgVals = append(sgVals, tfExpr(tfmapper.Reference{
				ResourceType: "aws_security_group",
				ResourceName: resolveRef(id, res.Metadata),
				Attribute:    "id",
			}.Expr()))
		}
		attrs["vpc_security_group_ids"] = tfList(sgVals)
	}

	addDependsOn(attrs, res)

	blocks = append(blocks, tfmapper.TerraformBlock{
		Kind:       "resource",
		Labels:     []string{"aws_db_instance", tfBlockName(res)},
		Attributes: attrs,
	})

	return blocks, nil
}

func (m *AWSMapper) mapAutoScalingGroup(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	minSize, _ := getInt(res.Metadata, "minSize")
	maxSize, _ := getInt(res.Metadata, "maxSize")
	desiredCapacity, _ := getInt(res.Metadata, "desiredCapacity")

	attrs := map[string]tfmapper.TerraformValue{
		"name":             tfStringOrVar(res.Metadata, "name", res.Name),
		"min_size":         tfNumber(float64(minSize)),
		"max_size":         tfNumber(float64(maxSize)),
		"desired_capacity": tfNumber(float64(desiredCapacity)),
	}

	var extraBlocks []tfmapper.TerraformBlock

	// Check if Launch Template is explicitly referenced
	ltID, _ := getString(res.Metadata, "launchTemplateId")

	nestedBlocks := make(map[string][]tfmapper.NestedBlock)

	// If launchTemplateId is provided, use it
	if ltID != "" && ltID != "implicit" {
		ltBlock := tfmapper.NestedBlock{
			Attributes: map[string]tfmapper.TerraformValue{
				"id":      tfExpr(tfmapper.Reference{ResourceType: "aws_launch_template", ResourceName: resolveRef(ltID, res.Metadata), Attribute: "id"}.Expr()),
				"version": tfString("$Latest"),
			},
		}
		nestedBlocks["launch_template"] = []tfmapper.NestedBlock{ltBlock}

		// Dependencies will be handled by addDependsOn via _dependsOn
	} else {
		// FALLBACK: Create implicit Launch Template (legacy behavior)
		ltName := res.Name + "-lt"
		if n, ok := getString(res.Metadata, "launchTemplateName"); ok && n != "" {
			ltName = n
		}

		// Create Launch Template resource
		ltAttrs := map[string]tfmapper.TerraformValue{
			"name_prefix":   tfString(ltName + "-"),
			"image_id":      tfStringOrVar(res.Metadata, "ami", "ami-0123456789"), // Default AMI if missing
			"instance_type": tfStringOrVar(res.Metadata, "instanceType", "t3.micro"),
		}

		// Add Launch Template block
		ltResourceType := "aws_launch_template"
		ltResourceName := tfName(res.ID) + "_lt"

		ltBlock := tfmapper.TerraformBlock{
			Kind:       "resource",
			Labels:     []string{ltResourceType, ltResourceName},
			Attributes: ltAttrs,
		}
		extraBlocks = append(extraBlocks, ltBlock)

		// Add reference to ASG nested block
		asgLtBlock := tfmapper.NestedBlock{
			Attributes: map[string]tfmapper.TerraformValue{
				"id":      tfExpr(tfmapper.Reference{ResourceType: ltResourceType, ResourceName: ltResourceName, Attribute: "id"}.Expr()),
				"version": tfString("$Latest"),
			},
		}
		nestedBlocks["launch_template"] = []tfmapper.NestedBlock{asgLtBlock}

		// Store ltBlock for return at the end
		// We'll handle this by checking if we made one
		// Wait, we need to return it.
		// Let's defer implicit LT return logic to later, but we need to keep ltBlock accessible.
		// Actually, the previous implementation returned directly inside the ELSE block.
		// But now I'm constructing `nestedBlocks` which is used at the end.

		// To fix lint and logic:
		// I will re-implement the return logic for implicit LT inside this block OR store it.
		// Let's store the fact that we have an implicit LT block to return.

		// Refactor: Just return here for implicit LT case?
		// But I need to process tags, subnets, etc.
		// So I should continue processing and return both blocks at the end.

		// Let's assign ltBlock to a variable outside?
		// The scope is tricky.

		// Better approach:
		// Revert to original structure slightly or use a variable
	}

	// Handle tags as `tag` blocks for ASG
	// Add Name tag
	tagBlocks := []tfmapper.NestedBlock{
		{
			Attributes: map[string]tfmapper.TerraformValue{
				"key":                 tfString("Name"),
				"value":               tfString(res.Name),
				"propagate_at_launch": tfBool(true),
			},
		},
	}

	// Add other tags if present (omitted for brevity effectively, but should be added here ideally if we parsed them)
	// For now just standard Name tag compliance as per user request
	nestedBlocks["tag"] = tagBlocks

	// Subnets
	var subnetIDs []string
	if subnetsRaw, ok := getArray(res.Metadata, "subnets"); ok {
		for _, sRaw := range subnetsRaw {
			if sMap, ok := sRaw.(map[string]interface{}); ok {
				if sid, ok := sMap["subnetId"].(string); ok {
					subnetIDs = append(subnetIDs, sid)
				}
			}
		}
	}

	if len(subnetIDs) > 0 {
		subnetRefs := make([]tfmapper.TerraformValue, 0, len(subnetIDs))
		for _, sid := range subnetIDs {
			subnetRefs = append(subnetRefs, tfExpr(tfmapper.Reference{
				ResourceType: "aws_subnet",
				ResourceName: resolveRef(sid, res.Metadata),
				Attribute:    "id",
			}.Expr()))
		}
		attrs["vpc_zone_identifier"] = tfList(subnetRefs)
	}

	// Target Groups
	var targetGroupARNs []string
	if list, ok := getArray(res.Metadata, "targetGroupArns"); ok {
		// If explicit ARNs provided (unlikely in our internal model, but possible)
		for _, item := range list {
			if s, ok := item.(string); ok {
				targetGroupARNs = append(targetGroupARNs, s)
			}
		}
	} else if list, ok := getStringSlice(res.Metadata, "targetGroupArns"); ok {
		targetGroupARNs = list
	}

	// Logic: If we have targetGroupIds, verify/convert them to references
	// This covers the case where usecase provides IDs
	var tgIDs []string
	if list, ok := getArray(res.Metadata, "targetGroupIds"); ok {
		for _, item := range list {
			if s, ok := item.(string); ok {
				tgIDs = append(tgIDs, s)
			}
		}
	} else if list, ok := getStringSlice(res.Metadata, "targetGroupIds"); ok {
		tgIDs = list
	}

	// If we have IDs, create references
	if len(tgIDs) > 0 {
		tgRefs := make([]tfmapper.TerraformValue, 0, len(tgIDs))
		for _, id := range tgIDs {
			tgRefs = append(tgRefs, tfExpr(tfmapper.Reference{
				ResourceType: "aws_lb_target_group",
				ResourceName: resolveRef(id, res.Metadata),
				Attribute:    "arn",
			}.Expr()))
		}
		attrs["target_group_arns"] = tfList(tgRefs)
	} else if len(targetGroupARNs) > 0 {
		// Use direct ARNs if provided
		validARNs := make([]tfmapper.TerraformValue, 0, len(targetGroupARNs))
		for _, arn := range targetGroupARNs {
			validARNs = append(validARNs, tfString(arn))
		}
		attrs["target_group_arns"] = tfList(validARNs)
	}

	addDependsOn(attrs, res)

	// If implicit LT was created, it's already in extraBlocks

	asgBlock := tfmapper.TerraformBlock{
		Kind:         "resource",
		Labels:       []string{"aws_autoscaling_group", tfBlockName(res)},
		Attributes:   attrs,
		NestedBlocks: nestedBlocks,
	}

	return append(extraBlocks, asgBlock), nil
}

func (m *AWSMapper) mapLaunchTemplate(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	attrs := map[string]tfmapper.TerraformValue{
		"name_prefix":   tfString(tfName(res.Name) + "-"),
		"image_id":      tfStringOrVar(res.Metadata, "image_id", "ami-0123456789"),
		"instance_type": tfStringOrVar(res.Metadata, "instance_type", "t3.micro"),
		"tags":          tfTags(res.Name),
	}

	// Allow overriding name_prefix with explicit Name if provided, but name_prefix is safer for conflicts
	if v, ok := getString(res.Metadata, "name"); ok && v != "" {
		// If explicit name is requested, use it, but LaunchTemplates are immutable so prefixes are better
		// We'll stick to name_prefix based on res.Name for now as per best practice
	}

	if v, ok := getString(res.Metadata, "update_default_version"); ok && v != "" {
		// defaults to true usually
	} else {
		attrs["update_default_version"] = tfBool(true)
	}

	if v, ok := getString(res.Metadata, "key_name"); ok && v != "" {
		attrs["key_name"] = tfStringOrVar(res.Metadata, "key_name", v)
	}

	if v, ok := getString(res.Metadata, "user_data"); ok && v != "" {
		attrs["user_data"] = tfString(v)
	}

	// IAM Instance Profile
	if v, ok := getString(res.Metadata, "iam_instance_profile"); ok && v != "" {
		attrs["iam_instance_profile"] = tfmapper.TerraformValue{
			Map: map[string]tfmapper.TerraformValue{
				"name": tfString(v),
			},
		}
	}

	// Security Groups
	var sgIDs []string
	if list, ok := getArray(res.Metadata, "vpc_security_group_ids"); ok {
		for _, item := range list {
			if s, ok := item.(string); ok {
				sgIDs = append(sgIDs, s)
			}
		}
	}
	// Also check securityGroups (object array)
	if len(sgIDs) == 0 {
		if list, ok := getArray(res.Metadata, "securityGroups"); ok {
			for _, item := range list {
				if m, ok := item.(map[string]interface{}); ok {
					if id, ok := m["id"].(string); ok {
						sgIDs = append(sgIDs, id)
					}
				}
			}
		}
	}
	// Fallback to securityGroupIds
	if len(sgIDs) == 0 {
		if list, ok := getStringSlice(res.Metadata, "securityGroupIds"); ok {
			sgIDs = list
		}
	}

	if len(sgIDs) > 0 {
		sgVals := make([]tfmapper.TerraformValue, 0, len(sgIDs))
		for _, id := range sgIDs {
			// Check if we have a name for this SG (provided by generator)
			sgName := resolveRef(id, res.Metadata)

			sgVals = append(sgVals, tfExpr(tfmapper.Reference{
				ResourceType: "aws_security_group",
				ResourceName: sgName,
				Attribute:    "id",
			}.Expr()))
		}
		attrs["vpc_security_group_ids"] = tfList(sgVals)
	}

	addDependsOn(attrs, res)

	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_launch_template", tfBlockName(res)},
			Attributes: attrs,
		},
	}, nil
}

// addDependsOn checks for _dependsOn metadata (injected by generator) and adds depends_on attribute
func addDependsOn(attrs map[string]tfmapper.TerraformValue, res *resource.Resource) {
	if res.Metadata == nil {
		return
	}
	deps, ok := res.Metadata["_dependsOn"].([]map[string]string)
	if !ok {
		// Try interface slice if unmarshaling weirdness
		if rawDeps, ok := res.Metadata["_dependsOn"].([]interface{}); ok {
			for _, d := range rawDeps {
				if dm, ok := d.(map[string]string); ok {
					deps = append(deps, dm)
				} else if dm, ok := d.(map[string]interface{}); ok {
					// Convert map[string]interface{} to map[string]string
					converted := make(map[string]string)
					if id, ok := dm["id"].(string); ok {
						converted["id"] = id
					}
					if typ, ok := dm["type"].(string); ok {
						converted["type"] = typ
					}
					deps = append(deps, converted)
				}
			}
		}
	}

	if len(deps) == 0 {
		return
	}

	var tfDeps []tfmapper.TerraformValue
	for _, dep := range deps {
		id := dep["id"]
		typ := dep["type"]

		tfType := getTerraformType(typ)
		if tfType != "" && id != "" {
			// depends_on = [aws_s3_bucket.bucket, ...] (whole resource reference)
			// Expression: "aws_type.name"
			refName := tfName(id)
			if name, ok := dep["name"]; ok && name != "" {
				s := tfName(name)
				if s != "resource" {
					refName = s
				}
			}

			tfDeps = append(tfDeps, tfExpr(tfmapper.TerraformExpr(fmt.Sprintf("%s.%s", tfType, refName))))
		}
	}

	if len(tfDeps) > 0 {
		attrs["depends_on"] = tfList(tfDeps)
	}
}

// getTerraformType maps domain resource type to Terraform resource type
func getTerraformType(domainType string) string {
	switch domainType {
	case "VPC":
		return "aws_vpc"
	case "Subnet":
		return "aws_subnet"
	case "EC2":
		return "aws_instance"
	case "SecurityGroup":
		return "aws_security_group"
	case "RouteTable":
		return "aws_route_table"
	case "InternetGateway":
		return "aws_internet_gateway"
	case "NATGateway":
		return "aws_nat_gateway"
	case "ElasticIP":
		return "aws_eip"
	case "S3":
		return "aws_s3_bucket"
	case "RDS":
		return "aws_db_instance"
	case "AutoScalingGroup":
		return "aws_autoscaling_group"
	case "LoadBalancer":
		return "aws_lb"
	case "Lambda":
		return "aws_lambda_function"
	case "Listener":
		return "aws_lb_listener"
	case "TargetGroup":
		return "aws_lb_target_group"
	default:
		return ""
	}
}

func getInt(m map[string]interface{}, key string) (int, bool) {
	if m == nil {
		return 0, false
	}
	v, ok := m[key]
	if !ok || v == nil {
		return 0, false
	}
	switch t := v.(type) {
	case int:
		return t, true
	case float64:
		return int(t), true
	default:
		return 0, false
	}
}

func getFloat(m map[string]interface{}, key string) (float64, bool) {
	if m == nil {
		return 0, false
	}
	v, ok := m[key]
	if !ok || v == nil {
		return 0, false
	}
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	default:
		return 0, false
	}
}

// tfBlockName returns the terraform local name for the resource.
// We prioritize Name to match user request, falling back to ID.
func tfBlockName(res *resource.Resource) string {
	// Prefer Name
	if res.Name != "" {
		s := tfName(res.Name)
		if s != "resource" {
			return s
		}
	}
	return tfName(res.ID)
}

// resolveRef resolves an ID to a Terraform reference name.
// It checks multiple mappings in order:
// 1. _resourceNames (domain ID -> name, injected by generator)
// 2. _originalIDToName (frontend node ID -> name, injected by LoadArchitecture)
// Falls back to tfName(id).
func resolveRef(id string, metadata map[string]interface{}) string {
	if metadata != nil {
		// First, check _resourceNames (domain/DB ID -> name)
		if mapRaw, ok := metadata["_resourceNames"]; ok {
			if mapping, ok := mapRaw.(map[string]string); ok {
				if name, found := mapping[id]; found {
					s := tfName(name)
					if s != "resource" {
						return s
					}
				}
			}
		}
		// Second, check _originalIDToName (frontend node ID -> name)
		// This handles references that use the original frontend IDs like "subnet-5"
		if mapRaw, ok := metadata["_originalIDToName"]; ok {
			if mapping, ok := mapRaw.(map[string]string); ok {
				if name, found := mapping[id]; found {
					s := tfName(name)
					if s != "resource" {
						return s
					}
				}
			}
		}
	}
	return tfName(id)
}

func (m *AWSMapper) mapListener(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	// Listener requires a Load Balancer (ParentID)
	if res.ParentID == nil || *res.ParentID == "" {
		return nil, fmt.Errorf("listener requires parent load balancer (parentID missing)")
	}

	port, _ := getInt(res.Metadata, "port")
	if port == 0 {
		port = 80 // Default
	}
	protocol, _ := getString(res.Metadata, "protocol")
	if protocol == "" {
		protocol = "HTTP" // Default
	}

	attrs := map[string]tfmapper.TerraformValue{
		"load_balancer_arn": tfExpr(tfmapper.Reference{ResourceType: "aws_lb", ResourceName: resolveRef(*res.ParentID, res.Metadata), Attribute: "arn"}.Expr()),
		"port":              tfNumber(float64(port)),
		"protocol":          tfString(protocol),
		"tags":              tfTags(res.Name),
	}

	// Default action
	// For now, we support "forward" to a Target Group
	defaultActionType, _ := getString(res.Metadata, "defaultActionType")
	if defaultActionType == "" {
		defaultActionType = "forward"
	}

	var defaultActionBlock tfmapper.NestedBlock

	if defaultActionType == "forward" {
		targetGroupARN := ""
		// Check explicit target group ID
		if tgID, ok := getString(res.Metadata, "targetGroupId"); ok && tgID != "" {
			// Resolve reference
			targetGroupARN = fmt.Sprintf("aws_lb_target_group.%s.arn", resolveRef(tgID, res.Metadata))
		}

		if targetGroupARN == "" {
			// Try to find a target group in dependencies?
			// Or maybe the user provided "targetGroupArn" directly?
			if arn, ok := getString(res.Metadata, "targetGroupArn"); ok && arn != "" {
				targetGroupARN = arn
			}
		}

		if targetGroupARN != "" {
			defaultActionBlock = tfmapper.NestedBlock{
				Attributes: map[string]tfmapper.TerraformValue{
					"type":             tfString("forward"),
					"target_group_arn": tfExpr(tfmapper.TerraformExpr(targetGroupARN)),
				},
			}
		} else {
			// Fallback: fixed-response if no TG found (to avoid invalid TF)
			defaultActionBlock = tfmapper.NestedBlock{
				Attributes: map[string]tfmapper.TerraformValue{
					"type": tfString("fixed-response"),
					"fixed_response": tfmapper.TerraformValue{
						Map: map[string]tfmapper.TerraformValue{
							"content_type": tfString("text/plain"),
							"message_body": tfString("No target group configured"),
							"status_code":  tfString("200"),
						},
					},
				},
			}
		}
	} else {
		// Implement other action types as needed
		defaultActionBlock = tfmapper.NestedBlock{
			Attributes: map[string]tfmapper.TerraformValue{
				"type": tfString("fixed-response"),
				"fixed_response": tfmapper.TerraformValue{
					Map: map[string]tfmapper.TerraformValue{
						"content_type": tfString("text/plain"),
						"message_body": tfString("Not Implemented"),
						"status_code":  tfString("501"),
					},
				},
			},
		}
	}

	return []tfmapper.TerraformBlock{
		{
			Kind:         "resource",
			Labels:       []string{"aws_lb_listener", tfBlockName(res)},
			Attributes:   attrs,
			NestedBlocks: map[string][]tfmapper.NestedBlock{"default_action": {defaultActionBlock}},
		},
	}, nil
}

func (m *AWSMapper) mapTargetGroup(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	// Target Group needs VPC ID (usually)
	// We can get VPC ID from parent if parent is VPC, or explicit "vpcId"
	vpcID := ""
	if res.ParentID != nil && *res.ParentID != "" {
		// If attached to VPC directly (not typical, but possible in graph)
		vpcID = *res.ParentID
	}
	if v, ok := getString(res.Metadata, "vpcId"); ok && v != "" {
		vpcID = v
	}

	port, _ := getInt(res.Metadata, "port")
	if port == 0 {
		port = 80
	}
	protocol, _ := getString(res.Metadata, "protocol")
	if protocol == "" {
		protocol = "HTTP"
	}
	targetType, _ := getString(res.Metadata, "targetType")
	if targetType == "" {
		targetType = "instance" // Default
	}

	attrs := map[string]tfmapper.TerraformValue{
		"name":        tfStringOrVar(res.Metadata, "name", res.Name),
		"port":        tfNumber(float64(port)),
		"protocol":    tfString(protocol),
		"target_type": tfString(targetType),
		"tags":        tfTags(res.Name),
	}

	if vpcID != "" {
		attrs["vpc_id"] = tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: resolveRef(vpcID, res.Metadata), Attribute: "id"}.Expr())
	}

	// Health Check
	nestedBlocks := make(map[string][]tfmapper.NestedBlock)
	if hc, ok := res.Metadata["healthCheck"].(map[string]interface{}); ok {
		hcAttrs := make(map[string]tfmapper.TerraformValue)
		if v, ok := getString(hc, "path"); ok {
			hcAttrs["path"] = tfString(v)
		}
		if v, ok := getString(hc, "protocol"); ok {
			hcAttrs["protocol"] = tfString(v)
		}
		if v, ok := getInt(hc, "port"); ok {
			hcAttrs["port"] = tfString(fmt.Sprintf("%d", v))
		} else if v, ok := getString(hc, "port"); ok {
			hcAttrs["port"] = tfString(v)
		}

		if len(hcAttrs) > 0 {
			// Add defaults if missing
			if _, ok := hcAttrs["path"]; !ok {
				hcAttrs["path"] = tfString("/")
			}
			block := tfmapper.NestedBlock{Attributes: hcAttrs}
			nestedBlocks["health_check"] = []tfmapper.NestedBlock{block}
		}
	}

	return []tfmapper.TerraformBlock{
		{
			Kind:         "resource",
			Labels:       []string{"aws_lb_target_group", tfBlockName(res)},
			Attributes:   attrs,
			NestedBlocks: nestedBlocks,
		},
	}, nil
}
