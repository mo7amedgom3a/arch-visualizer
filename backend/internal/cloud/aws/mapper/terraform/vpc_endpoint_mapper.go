package terraform

import (
	"fmt"

	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

// MapVPCEndpoint maps a VPC Endpoint resource to Terraform blocks
func MapVPCEndpoint(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
	if res == nil {
		return nil, fmt.Errorf("resource is nil")
	}

	vpcID, ok := getString(res.Metadata, "vpc_id")
	if !ok || vpcID == "" {
		// Try to fallback to ParentID if vpc_id is not in metadata
		if res.ParentID != nil && *res.ParentID != "" {
			vpcID = *res.ParentID
		} else {
			return nil, fmt.Errorf("vpc_endpoint requires vpc_id or parent vpc")
		}
	}

	serviceName, ok := getString(res.Metadata, "service_name")
	if !ok || serviceName == "" {
		return nil, fmt.Errorf("missing required config %q", "service_name")
	}

	endpointType, _ := getString(res.Metadata, "vpc_endpoint_type")
	if endpointType == "" {
		endpointType = "Gateway" // Default
	}

	attrs := map[string]tfmapper.TerraformValue{
		"vpc_id":            tfExpr(tfmapper.Reference{ResourceType: "aws_vpc", ResourceName: resolveRef(vpcID, res.Metadata), Attribute: "id"}.Expr()),
		"service_name":      tfStringOrVar(res.Metadata, "service_name", serviceName),
		"vpc_endpoint_type": tfString(endpointType),
		"tags":              tfTags(res.Name),
	}

	// Handle Policy
	if policy, ok := getString(res.Metadata, "policy"); ok && policy != "" {
		attrs["policy"] = tfString(policy)
	}

	// Handle Private DNS (Interface only)
	if endpointType == string(awsnetworking.VPCEndpointTypeInterface) {
		if v, ok := getBool(res.Metadata, "private_dns_enabled"); ok {
			attrs["private_dns_enabled"] = tfBool(v)
		}
	}

	// Handle Subnets (Interface only)
	if endpointType == string(awsnetworking.VPCEndpointTypeInterface) {
		var subnetIDs []string
		if list, ok := getArray(res.Metadata, "subnet_ids"); ok {
			for _, item := range list {
				if id, ok := item.(string); ok {
					subnetIDs = append(subnetIDs, id)
				}
			}
		}
		// Fallback to subnets (array of objects) if simple ID list is empty?
		// For now simple list.

		if len(subnetIDs) > 0 {
			subnetVals := make([]tfmapper.TerraformValue, 0, len(subnetIDs))
			for _, id := range subnetIDs {
				subnetVals = append(subnetVals, tfExpr(tfmapper.Reference{
					ResourceType: "aws_subnet",
					ResourceName: resolveRef(id, res.Metadata),
					Attribute:    "id",
				}.Expr()))
			}
			attrs["subnet_ids"] = tfList(subnetVals)
		}
	}

	// Handle Security Groups (Interface only)
	if endpointType == string(awsnetworking.VPCEndpointTypeInterface) {
		var sgIDs []string
		if list, ok := getArray(res.Metadata, "security_group_ids"); ok {
			for _, item := range list {
				if id, ok := item.(string); ok {
					sgIDs = append(sgIDs, id)
				}
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
			attrs["security_group_ids"] = tfList(sgVals)
		}
	}

	// Handle Route Tables (Gateway only)
	if endpointType == string(awsnetworking.VPCEndpointTypeGateway) {
		var rtIDs []string
		if list, ok := getArray(res.Metadata, "route_table_ids"); ok {
			for _, item := range list {
				if id, ok := item.(string); ok {
					rtIDs = append(rtIDs, id)
				}
			}
		}

		if len(rtIDs) > 0 {
			rtVals := make([]tfmapper.TerraformValue, 0, len(rtIDs))
			for _, id := range rtIDs {
				rtVals = append(rtVals, tfExpr(tfmapper.Reference{
					ResourceType: "aws_route_table",
					ResourceName: resolveRef(id, res.Metadata),
					Attribute:    "id",
				}.Expr()))
			}
			attrs["route_table_ids"] = tfList(rtVals)
		}
	}

	return []tfmapper.TerraformBlock{
		{
			Kind:       "resource",
			Labels:     []string{"aws_vpc_endpoint", tfBlockName(res)},
			Attributes: attrs,
		},
	}, nil
}
