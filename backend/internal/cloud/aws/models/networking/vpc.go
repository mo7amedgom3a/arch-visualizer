package networking

import (
	"errors"
	"fmt"
	"net"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

type VPC struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	CIDR   string `json:"cidr"`
	// +optional
	EnableDNSHostnames bool `json:"enable_dns_hostnames"`
	// +optional
	EnableDNSSupport bool `json:"enable_dns_support"`
	// +optional
	InstanceTenancy string `json:"instance_tenancy"`
	// +optional
	Tags []configs.Tag `json:"tags"`
}

func (vpc *VPC) Validate() error {
	if vpc.Name == "" {
		return errors.New("name is required")
	}
	if vpc.Region == "" {
		return errors.New("region is required")
	}
	if vpc.CIDR == "" {
		return errors.New("cidr is required")
	}

	// Validate CIDR format
	_, ipNet, err := net.ParseCIDR(vpc.CIDR)
	if err != nil {
		return fmt.Errorf("invalid cidr format: %w", err)
	}

	// Validate CIDR is IPv4 (AWS VPCs currently only support IPv4)
	if ipNet.IP.To4() == nil {
		return errors.New("cidr must be IPv4")
	}

	return nil
}
