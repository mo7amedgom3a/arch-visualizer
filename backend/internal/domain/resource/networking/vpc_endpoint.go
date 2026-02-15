package networking

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// VPCEndpointType represents the type of VPC Endpoint
type VPCEndpointType string

const (
	VPCEndpointTypeInterface VPCEndpointType = "Interface"
	VPCEndpointTypeGateway   VPCEndpointType = "Gateway"
)

// VPCEndpoint represents a cloud-agnostic VPC Endpoint
type VPCEndpoint struct {
	ID                string
	ARN               *string
	Name              string
	VPCID             string
	ServiceName       string
	Type              VPCEndpointType
	SubnetIDs         []string
	SecurityGroupIDs  []string
	RouteTableIDs     []string
	PrivateDNSEnabled bool
	Policy            string
	Tags              []configs.Tag
}

// Validate performs domain-level validation
func (e *VPCEndpoint) Validate() error {
	if e.Name == "" {
		return errors.New("vpc endpoint name is required")
	}
	if e.VPCID == "" {
		return errors.New("vpc id is required")
	}
	if e.ServiceName == "" {
		return errors.New("service name is required")
	}

	if e.Type == VPCEndpointTypeInterface {
		if len(e.SubnetIDs) == 0 {
			return errors.New("interface endpoints require at least one subnet")
		}
		if len(e.SecurityGroupIDs) == 0 {
			return errors.New("interface endpoints require at least one security group")
		}
	} else if e.Type == VPCEndpointTypeGateway {
		if len(e.RouteTableIDs) == 0 {
			return errors.New("gateway endpoints require at least one route table")
		}
	} else {
		// Default to Gateway if not specified, or validate if strict
		if e.Type != "" {
			return fmt.Errorf("invalid vpc endpoint type: %s", e.Type)
		}
		// If empty, we might infer or error out. Let's error for clarity.
		return errors.New("vpc endpoint type is required")
	}

	// Validate service name format if possible (basic check)
	if !strings.Contains(e.ServiceName, ".") {
		return errors.New("invalid service name format")
	}

	return nil
}
