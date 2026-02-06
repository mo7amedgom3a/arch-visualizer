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
	inv.SetTerraformMapper("Lambda", mapper.mapLambda)

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

	// Resolve Parent VPC Name
	vpcName := tfName(*res.ParentID)
	if pName, ok := res.Metadata["_parentName"].(string); ok {
		vpcName = tfName(pName)
	}

	attrs := map[string]tfmapper.TerraformValue{
		"vpc_id":            tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: vpcName, Attribute: "id"}.Expr()),
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
		"subnet_id":     tfExpr(tfmapper.Reference{ResourceType: "aws_subnet", ResourceName: tfName(*res.ParentID), Attribute: "id"}.Expr()),
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
		subnetName := tfName(*res.ParentID) // Default to ID
		if n, ok := res.Metadata["_parentName"].(string); ok {
			subnetName = tfName(n)
		}
		attrs["subnet_id"] = tfExpr(tfmapper.Reference{ResourceType: "aws_subnet", ResourceName: subnetName, Attribute: "id"}.Expr())
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
			// Reference the security group resource by its ID
			// The ID should match the resource ID in the architecture
			sgVals = append(sgVals, tfExpr(tfmapper.Reference{
				ResourceType: "aws_security_group",
				ResourceName: tfName(id),
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

	// Resolve Parent VPC Name
	vpcName := tfName(*res.ParentID)
	if pName, ok := res.Metadata["_parentName"].(string); ok {
		vpcName = tfName(pName)
	}

	attrs := map[string]tfmapper.TerraformValue{
		"name":        tfString(res.Name),
		"description": tfString(desc),
		"vpc_id":      tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: vpcName, Attribute: "id"}.Expr()),
		"tags":        tfTags(res.Name),
	}

	// Initialize nested blocks for inline ingress/egress rules
	nestedBlocks := make(map[string][]tfmapper.NestedBlock)
	var ingressBlocks []tfmapper.NestedBlock
	var egressBlocks []tfmapper.NestedBlock
	hasEgressRule := false

	// Parse and add inline ingress/egress rules
	if rules, ok := getArray(res.Metadata, "rules"); ok {
		for _, ruleRaw := range rules {
			rule, ok := ruleRaw.(map[string]interface{})
			if !ok {
				continue
			}

			ruleType, _ := getString(rule, "type")
			protocol, _ := getString(rule, "protocol")
			portRange, _ := getString(rule, "portRange")
			cidr, _ := getString(rule, "cidr")
			description, _ := getString(rule, "description")

			if ruleType == "" || protocol == "" {
				continue
			}

			// Parse port range (e.g., "22" or "80-443")
			fromPort, toPort := parsePortRange(portRange)

			ruleAttrs := map[string]tfmapper.TerraformValue{
				"protocol": tfString(protocol),
			}

			// Add description if provided
			if description != "" {
				ruleAttrs["description"] = tfString(description)
			}

			// For protocol "-1" (all), from_port and to_port must be 0
			if protocol == "-1" {
				ruleAttrs["from_port"] = tfNumber(0)
				ruleAttrs["to_port"] = tfNumber(0)
			} else {
				if fromPort != nil {
					ruleAttrs["from_port"] = tfNumber(float64(*fromPort))
				}
				if toPort != nil {
					ruleAttrs["to_port"] = tfNumber(float64(*toPort))
				}
			}

			if cidr != "" {
				ruleAttrs["cidr_blocks"] = tfList([]tfmapper.TerraformValue{tfString(cidr)})
			}

			nestedBlock := tfmapper.NestedBlock{Attributes: ruleAttrs}

			if ruleType == "ingress" {
				ingressBlocks = append(ingressBlocks, nestedBlock)
			} else if ruleType == "egress" {
				egressBlocks = append(egressBlocks, nestedBlock)
				hasEgressRule = true
			}
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

	vpcName := tfName(*res.ParentID)
	if pName, ok := res.Metadata["_parentName"].(string); ok {
		vpcName = tfName(pName)
	}

	blocks := []tfmapper.TerraformBlock{
		{
			Kind:   "resource",
			Labels: []string{"aws_route_table", tfBlockName(res)},
			Attributes: map[string]tfmapper.TerraformValue{
				"vpc_id": tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: vpcName, Attribute: "id"}.Expr()),
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
						ResourceName: tfName(targetResourceID), // Use tfName for targetResourceID as it's an ID
						Attribute:    "id",
					}.Expr())
				}
			case "NATGateway":
				if targetResourceID != "" {
					routeAttrs["nat_gateway_id"] = tfExpr(tfmapper.Reference{
						ResourceType: "aws_nat_gateway",
						ResourceName: tfName(targetResourceID), // Use tfName for targetResourceID as it's an ID
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
			associatedResourceName, _ := getString(assoc, "_associatedResourceName") // Get name if available

			if associatedResourceID == "" {
				continue
			}

			assocAttrs := map[string]tfmapper.TerraformValue{
				"route_table_id": tfExpr(tfmapper.Reference{ResourceType: "aws_route_table", ResourceName: tfBlockName(res), Attribute: "id"}.Expr()),
			}

			// Map association type to Terraform resource
			switch assocType {
			case "Subnet":
				// Prefer name if available, otherwise fall back to ID
				resourceName := tfName(associatedResourceID)
				if associatedResourceName != "" {
					resourceName = tfName(associatedResourceName)
				}
				assocAttrs["subnet_id"] = tfExpr(tfmapper.Reference{
					ResourceType: "aws_subnet",
					ResourceName: resourceName,
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

	// Resolve Parent VPC Name
	vpcName := tfName(*res.ParentID)
	if pName, ok := res.Metadata["_parentName"].(string); ok {
		vpcName = tfName(pName)
	}

	attrs := map[string]tfmapper.TerraformValue{
		"vpc_id": tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: vpcName, Attribute: "id"}.Expr()),
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
		"subnet_id": tfExpr(tfmapper.Reference{ResourceType: "aws_subnet", ResourceName: tfName(subnetID), Attribute: "id"}.Expr()),
		"tags":      tfTags(res.Name),
	}

	// Use _subnetName if available
	if n, ok := res.Metadata["_subnetName"].(string); ok {
		natGatewayAttrs["subnet_id"] = tfExpr(tfmapper.Reference{ResourceType: "aws_subnet", ResourceName: tfName(n), Attribute: "id"}.Expr())
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

	attrs := map[string]tfmapper.TerraformValue{
		"identifier":          tfString(tfName(res.ID)),
		"instance_class":      tfString(instanceClass),
		"engine":              tfString(engine),
		"tags":                tfTags(res.Name),
		"skip_final_snapshot": tfBool(true), // Avoid hanging on destroy in dev environments
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

	// Handle db_subnet_group_name
	if v, ok := getString(res.Metadata, "db_subnet_group_name"); ok && v != "" {
		attrs["db_subnet_group_name"] = tfString(v)
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
	if len(sgIDs) > 0 {
		// Use _securityGroupNames if available
		sgNamesMap, _ := res.Metadata["_securityGroupNames"].(map[string]string)

		sgVals := make([]tfmapper.TerraformValue, 0, len(sgIDs))
		for _, id := range sgIDs {
			refName := tfName(id)
			if name, ok := sgNamesMap[id]; ok {
				refName = tfName(name)
			}

			sgVals = append(sgVals, tfExpr(tfmapper.Reference{
				ResourceType: "aws_security_group",
				ResourceName: refName,
				Attribute:    "id",
			}.Expr()))
		}
		attrs["vpc_security_group_ids"] = tfList(sgVals)
	}

	addDependsOn(attrs, res)

	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_db_instance", tfBlockName(res)},
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
	case "Lambda":
		return "aws_lambda_function"
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
