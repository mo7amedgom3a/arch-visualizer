package networking

import (
	"errors"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

type VPCEndpointType string

const (
	VPCEndpointTypeGateway   VPCEndpointType = "Gateway"
	VPCEndpointTypeInterface VPCEndpointType = "Interface"
)

type VPCEndpoint struct {
	Name string `json:"name"`
	// +required
	VPCID string `json:"vpc_id"`
	// +required
	ServiceName string `json:"service_name"`
	// +optional
	VPCEndpointType VPCEndpointType `json:"vpc_endpoint_type"`
	// +optional
	SubnetIDs []string `json:"subnet_ids"`
	// +optional
	SecurityGroupIDs []string `json:"security_group_ids"`
	// +optional
	PrivateDNSEnabled bool `json:"private_dns_enabled"`
	// +optional
	RouteTableIDs []string `json:"route_table_ids"`
	// +optional
	Policy string `json:"policy"`
	// +optional
	Tags []configs.Tag `json:"tags"`
}

func (e *VPCEndpoint) Validate() error {
	if e.Name == "" {
		return errors.New("name is required")
	}
	if e.VPCID == "" {
		return errors.New("vpc_id is required")
	}
	if e.ServiceName == "" {
		return errors.New("service_name is required")
	}

	if e.VPCEndpointType == "" {
		e.VPCEndpointType = VPCEndpointTypeGateway
	}

	if e.VPCEndpointType == VPCEndpointTypeInterface && len(e.SubnetIDs) == 0 {
		return errors.New("subnet_ids are required for Interface type VPC Endpoint")
	}

	return nil
}
