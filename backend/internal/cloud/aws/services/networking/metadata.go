package networking

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services"
)

// ---------------------------------------------------------------------------
// Internal registry
// ---------------------------------------------------------------------------

var networkingSchemaRegistry = map[string]*services.ResourceSchema{}

func registerNetworkingSchema(label string, schema *services.ResourceSchema) {
	networkingSchemaRegistry[label] = schema
}

// ---------------------------------------------------------------------------
// Service interface & implementation
// ---------------------------------------------------------------------------

// NetworkingMetadataService exposes structured schemas for AWS networking resources.
type NetworkingMetadataService interface {
	GetResourceSchema(ctx context.Context, resource string) (*services.ResourceSchema, error)
	ListResourceSchemas(ctx context.Context) ([]*services.ResourceSchema, error)
}

type networkingMetadataServiceImpl struct{}

// NewNetworkingMetadataService returns a ready-to-use metadata service.
func NewNetworkingMetadataService() NetworkingMetadataService {
	return &networkingMetadataServiceImpl{}
}

func (s *networkingMetadataServiceImpl) GetResourceSchema(_ context.Context, resource string) (*services.ResourceSchema, error) {
	schema, ok := networkingSchemaRegistry[resource]
	if !ok {
		return nil, fmt.Errorf("unknown networking resource: %s", resource)
	}
	return schema, nil
}

func (s *networkingMetadataServiceImpl) ListResourceSchemas(_ context.Context) ([]*services.ResourceSchema, error) {
	out := make([]*services.ResourceSchema, 0, len(networkingSchemaRegistry))
	for _, schema := range networkingSchemaRegistry {
		out = append(out, schema)
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Schema registrations (run once at import time)
// ---------------------------------------------------------------------------

func init() {
	// ---- VPC ----
	registerNetworkingSchema("vpc", &services.ResourceSchema{
		Label: "vpc",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the VPC"},
			{Name: "region", Type: "string", Required: true, Enum: []string{}, Description: "AWS region"},
			{Name: "cidr", Type: "string", Required: true, Enum: []string{}, Description: "IPv4 CIDR block (e.g. 10.0.0.0/16)"},
			{Name: "enable_dns_hostnames", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Enable DNS hostnames"},
			{Name: "enable_dns_support", Type: "bool", Required: false, Enum: []string{}, Default: true, Description: "Enable DNS support"},
			{Name: "instance_tenancy", Type: "string", Required: false, Enum: []string{"default", "dedicated", "host"}, Default: "default", Description: "Tenancy of instances"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":         "string",
			"arn":        "string",
			"state":      "string",
			"owner_id":   "string",
			"is_default": "bool",
		},
	})

	// ---- Subnet ----
	registerNetworkingSchema("subnet", &services.ResourceSchema{
		Label: "subnet",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the subnet"},
			{Name: "vpc_id", Type: "string", Required: true, Enum: []string{}, Description: "VPC ID"},
			{Name: "cidr", Type: "string", Required: true, Enum: []string{}, Description: "IPv4 CIDR block"},
			{Name: "availability_zone", Type: "string", Required: true, Enum: []string{}, Description: "Availability zone"},
			{Name: "map_public_ip_on_launch", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Auto-assign public IP"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":                 "string",
			"arn":                "string",
			"state":              "string",
			"available_ip_count": "int",
		},
	})

	// ---- Internet Gateway ----
	registerNetworkingSchema("internet_gateway", &services.ResourceSchema{
		Label: "internet_gateway",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the internet gateway"},
			{Name: "vpc_id", Type: "string", Required: true, Enum: []string{}, Description: "VPC to attach to"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":               "string",
			"arn":              "string",
			"state":            "string",
			"attachment_state": "string",
		},
	})

	// ---- Route Table ----
	registerNetworkingSchema("route_table", &services.ResourceSchema{
		Label: "route_table",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the route table"},
			{Name: "vpc_id", Type: "string", Required: true, Enum: []string{}, Description: "VPC ID"},
			{Name: "routes", Type: "object", Required: false, Enum: []string{}, Description: "Route entries"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":  "string",
			"arn": "string",
		},
	})

	// ---- Security Group ----
	registerNetworkingSchema("security_group", &services.ResourceSchema{
		Label: "security_group",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the security group"},
			{Name: "description", Type: "string", Required: true, Enum: []string{}, Description: "Description (max 255 chars)"},
			{Name: "vpc_id", Type: "string", Required: true, Enum: []string{}, Description: "VPC ID"},
			{Name: "rules", Type: "object", Required: false, Enum: []string{}, Description: "Ingress and egress rules"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":  "string",
			"arn": "string",
		},
	})

	// ---- NAT Gateway ----
	registerNetworkingSchema("nat_gateway", &services.ResourceSchema{
		Label: "nat_gateway",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the NAT gateway"},
			{Name: "subnet_id", Type: "string", Required: true, Enum: []string{}, Description: "Subnet to place the NAT gateway in"},
			{Name: "allocation_id", Type: "string", Required: true, Enum: []string{}, Description: "Elastic IP allocation ID"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":         "string",
			"arn":        "string",
			"state":      "string",
			"public_ip":  "string",
			"private_ip": "string",
		},
	})

	// ---- Elastic IP ----
	registerNetworkingSchema("elastic_ip", &services.ResourceSchema{
		Label: "elastic_ip",
		Fields: []services.FieldDescriptor{
			{Name: "region", Type: "string", Required: true, Enum: []string{}, Description: "AWS region"},
			{Name: "allocation_id", Type: "string", Required: false, Enum: []string{}, Description: "Use existing EIP allocation ID"},
			{Name: "address_pool_type", Type: "string", Required: false, Enum: []string{"amazon", "byoip", "customer_owned", "ipam"}, Default: "amazon", Description: "IP address pool type"},
			{Name: "address_pool_id", Type: "string", Required: false, Enum: []string{}, Description: "Pool ID for non-amazon pool types"},
			{Name: "network_border_group", Type: "string", Required: false, Enum: []string{}, Description: "Network border group"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":            "string",
			"arn":           "string",
			"public_ip":     "string",
			"allocation_id": "string",
			"state":         "string",
			"domain":        "string",
		},
	})

	// ---- Network ACL ----
	registerNetworkingSchema("network_acl", &services.ResourceSchema{
		Label: "network_acl",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the network ACL"},
			{Name: "vpc_id", Type: "string", Required: true, Enum: []string{}, Description: "VPC ID"},
			{Name: "inbound_rules", Type: "object", Required: false, Enum: []string{}, Description: "Inbound ACL rules"},
			{Name: "outbound_rules", Type: "object", Required: false, Enum: []string{}, Description: "Outbound ACL rules"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":         "string",
			"arn":        "string",
			"is_default": "bool",
		},
	})

	// ---- Network Interface ----
	registerNetworkingSchema("network_interface", &services.ResourceSchema{
		Label: "network_interface",
		Fields: []services.FieldDescriptor{
			{Name: "description", Type: "string", Required: false, Enum: []string{}, Description: "Description"},
			{Name: "subnet_id", Type: "string", Required: true, Enum: []string{}, Description: "Subnet to create the interface in"},
			{Name: "interface_type", Type: "string", Required: true, Enum: []string{"elastic", "attached"}, Description: "Type of network interface"},
			{Name: "private_ipv4_address", Type: "string", Required: false, Enum: []string{}, Description: "Custom private IPv4 address"},
			{Name: "auto_assign_private_ip", Type: "bool", Required: false, Enum: []string{}, Default: true, Description: "Auto-assign private IP"},
			{Name: "security_group_ids", Type: "[]string", Required: true, Enum: []string{}, Description: "Security groups to attach (1-5)"},
			{Name: "source_dest_check", Type: "bool", Required: false, Enum: []string{}, Default: true, Description: "Source/destination check"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":                   "string",
			"arn":                  "string",
			"status":               "string",
			"private_ipv4_address": "string",
			"mac_address":          "string",
			"vpc_id":               "string",
			"availability_zone":    "string",
		},
	})

	// ---- VPC Endpoint ----
	registerNetworkingSchema("vpc_endpoint", &services.ResourceSchema{
		Label: "vpc_endpoint",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the VPC endpoint"},
			{Name: "vpc_id", Type: "string", Required: true, Enum: []string{}, Description: "VPC ID"},
			{Name: "service_name", Type: "string", Required: true, Enum: []string{}, Description: "AWS service name (e.g. com.amazonaws.us-east-1.s3)"},
			{Name: "vpc_endpoint_type", Type: "string", Required: false, Enum: []string{"Gateway", "Interface"}, Default: "Gateway", Description: "Endpoint type"},
			{Name: "subnet_ids", Type: "[]string", Required: false, Enum: []string{}, Description: "Subnet IDs (required for Interface type)"},
			{Name: "security_group_ids", Type: "[]string", Required: false, Enum: []string{}, Description: "Security group IDs"},
			{Name: "private_dns_enabled", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Enable private DNS"},
			{Name: "route_table_ids", Type: "[]string", Required: false, Enum: []string{}, Description: "Route table IDs (for Gateway type)"},
			{Name: "policy", Type: "string", Required: false, Enum: []string{}, Description: "VPC endpoint policy JSON"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":  "string",
			"arn": "string",
		},
	})
}
