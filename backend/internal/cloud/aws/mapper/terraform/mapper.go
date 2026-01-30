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
		"cidr_block": tfString(cidr),
		"tags":       tfTags(res.Name),
	}

	if v, ok := getBool(res.Metadata, "enable_dns_hostnames"); ok {
		attrs["enable_dns_hostnames"] = tfBool(v)
	}
	if v, ok := getBool(res.Metadata, "enable_dns_support"); ok {
		attrs["enable_dns_support"] = tfBool(v)
	}
	if v, ok := getString(res.Metadata, "instance_tenancy"); ok && v != "" {
		attrs["instance_tenancy"] = tfString(v)
	}

	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_vpc", tfName(res.ID)},
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
		"vpc_id":            tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: tfName(*res.ParentID), Attribute: "id"}.Expr()),
		"cidr_block":        tfString(cidr),
		"availability_zone": tfString(az),
		"tags":              tfTags(res.Name),
	}

	if v, ok := getBool(res.Metadata, "map_public_ip_on_launch"); ok {
		attrs["map_public_ip_on_launch"] = tfBool(v)
	}

	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_subnet", tfName(res.ID)},
			Attributes: attrs,
		},
	}, nil
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
		"ami":           tfString(ami),
		"instance_type": tfString(instanceType),
		"subnet_id":     tfExpr(tfmapper.Reference{ResourceType: "aws_subnet", ResourceName: tfName(*res.ParentID), Attribute: "id"}.Expr()),
		"tags":          tfTags(res.Name),
	}

	if v, ok := getString(res.Metadata, "keyName"); ok && v != "" {
		attrs["key_name"] = tfString(v)
	}
	if v, ok := getString(res.Metadata, "iamInstanceProfile"); ok && v != "" {
		attrs["iam_instance_profile"] = tfString(v)
	}
	if v, ok := getString(res.Metadata, "userData"); ok && v != "" {
		attrs["user_data"] = tfString(v)
	}

	// Handle security groups - can be array of strings or array of objects with "id" field
	var sgIDs []string
	if list, ok := getArray(res.Metadata, "securityGroups"); ok {
		for _, item := range list {
			if sgObj, ok := item.(map[string]interface{}); ok {
				if id, ok := getString(sgObj, "id"); ok && id != "" {
					sgIDs = append(sgIDs, id)
				}
			} else if sgStr, ok := item.(string); ok && sgStr != "" {
				sgIDs = append(sgIDs, sgStr)
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

	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_instance", tfName(res.ID)},
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
		"vpc_id":      tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: tfName(*res.ParentID), Attribute: "id"}.Expr()),
		"tags":        tfTags(res.Name),
	}

	blocks := []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_security_group", tfName(res.ID)},
			Attributes: attrs,
		},
	}

	// Parse and add ingress/egress rules
	if rules, ok := getArray(res.Metadata, "rules"); ok {
		for i, ruleRaw := range rules {
			rule, ok := ruleRaw.(map[string]interface{})
			if !ok {
				continue
			}

			ruleType, _ := getString(rule, "type")
			protocol, _ := getString(rule, "protocol")
			portRange, _ := getString(rule, "portRange")
			cidr, _ := getString(rule, "cidr")

			if ruleType == "" || protocol == "" {
				continue
			}

			// Parse port range (e.g., "22" or "80-443")
			fromPort, toPort := parsePortRange(portRange)

			ruleAttrs := map[string]tfmapper.TerraformValue{
				"security_group_id": tfExpr(tfmapper.Reference{ResourceType: "aws_security_group", ResourceName: tfName(res.ID), Attribute: "id"}.Expr()),
				"protocol":          tfString(protocol),
			}

			if fromPort != nil {
				ruleAttrs["from_port"] = tfNumber(float64(*fromPort))
			}
			if toPort != nil {
				ruleAttrs["to_port"] = tfNumber(float64(*toPort))
			}

			if cidr != "" {
				ruleAttrs["cidr_blocks"] = tfList([]tfmapper.TerraformValue{tfString(cidr)})
			}

			var ruleKind string
			if ruleType == "ingress" {
				ruleKind = "aws_security_group_rule"
				ruleAttrs["type"] = tfString("ingress")
			} else if ruleType == "egress" {
				ruleKind = "aws_security_group_rule"
				ruleAttrs["type"] = tfString("egress")
			} else {
				continue
			}

			blocks = append(blocks, tfmapper.TerraformBlock{
				Kind:       "resource",
				Labels:     []string{ruleKind, fmt.Sprintf("%s_rule_%d", tfName(res.ID), i)},
				Attributes: ruleAttrs,
			})
		}
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
			Labels: []string{"aws_route_table", tfName(res.ID)},
			Attributes: map[string]tfmapper.TerraformValue{
				"vpc_id": tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: tfName(*res.ParentID), Attribute: "id"}.Expr()),
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
				"route_table_id":         tfExpr(tfmapper.Reference{ResourceType: "aws_route_table", ResourceName: tfName(res.ID), Attribute: "id"}.Expr()),
				"destination_cidr_block": tfString(destCIDR),
			}

			// Map target type to Terraform resource
			switch targetType {
			case "InternetGateway":
				if targetResourceID != "" {
					routeAttrs["gateway_id"] = tfExpr(tfmapper.Reference{
						ResourceType: "aws_internet_gateway",
						ResourceName: tfName(targetResourceID),
						Attribute:    "id",
					}.Expr())
				}
			case "NATGateway":
				if targetResourceID != "" {
					routeAttrs["nat_gateway_id"] = tfExpr(tfmapper.Reference{
						ResourceType: "aws_nat_gateway",
						ResourceName: tfName(targetResourceID),
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
				Labels:     []string{"aws_route", fmt.Sprintf("%s_route_%d", tfName(res.ID), i)},
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
				"route_table_id": tfExpr(tfmapper.Reference{ResourceType: "aws_route_table", ResourceName: tfName(res.ID), Attribute: "id"}.Expr()),
			}

			// Map association type to Terraform resource
			switch assocType {
			case "Subnet":
				assocAttrs["subnet_id"] = tfExpr(tfmapper.Reference{
					ResourceType: "aws_subnet",
					ResourceName: tfName(associatedResourceID),
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
				Labels:     []string{"aws_route_table_association", fmt.Sprintf("%s_assoc_%d", tfName(res.ID), i)},
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
		"vpc_id": tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: tfName(*res.ParentID), Attribute: "id"}.Expr()),
		"tags":   tfTags(res.Name),
	}
	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_internet_gateway", tfName(res.ID)},
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
			Labels:     []string{"aws_eip", tfName(res.ID)},
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

	attrs := map[string]tfmapper.TerraformValue{
		"subnet_id": tfExpr(tfmapper.Reference{ResourceType: "aws_subnet", ResourceName: tfName(subnetID), Attribute: "id"}.Expr()),
		"tags":      tfTags(res.Name),
	}

	// allocationId optional; user can reference an EIP allocation id directly.
	if alloc, ok := getString(res.Metadata, "allocationId"); ok && alloc != "" {
		attrs["allocation_id"] = tfString(alloc)
	}

	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_nat_gateway", tfName(res.ID)},
			Attributes: attrs,
		},
	}, nil
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
			Labels:     []string{"aws_s3_bucket", tfName(res.ID)},
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
