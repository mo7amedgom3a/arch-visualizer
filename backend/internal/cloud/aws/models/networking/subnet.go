package networking

import (
	"errors"
	"fmt"
	"net"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// Subnet represents an AWS-specific subnet
type Subnet struct {
	Name                string        `json:"name"`
	VPCID               string        `json:"vpc_id"`
	CIDR                string        `json:"cidr"`
	AvailabilityZone    string        `json:"availability_zone"`
	MapPublicIPOnLaunch bool          `json:"map_public_ip_on_launch"`
	Tags                []configs.Tag `json:"tags"`
}

// Validate performs AWS-specific validation
func (s *Subnet) Validate() error {
	if s.Name == "" {
		return errors.New("subnet name is required")
	}
	if s.VPCID == "" {
		return errors.New("subnet vpc_id is required")
	}
	if s.CIDR == "" {
		return errors.New("subnet cidr is required")
	}
	if s.AvailabilityZone == "" {
		return errors.New("subnet availability_zone is required")
	}

	// Validate CIDR format
	_, ipNet, err := net.ParseCIDR(s.CIDR)
	if err != nil {
		return fmt.Errorf("invalid cidr format: %w", err)
	}

	// Validate CIDR is IPv4
	if ipNet.IP.To4() == nil {
		return errors.New("subnet cidr must be IPv4")
	}

	return nil
}

// GetCIDRBlock returns the parsed CIDR block
func (s *Subnet) GetCIDRBlock() (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(s.CIDR)
	if err != nil {
		return nil, fmt.Errorf("invalid cidr format: %w", err)
	}
	return ipNet, nil
}

// ValidateCIDRInVPC validates that subnet CIDR is within VPC CIDR
func (s *Subnet) ValidateCIDRInVPC(vpcCIDR string) error {
	contains, err := CIDRContains(vpcCIDR, s.CIDR)
	if err != nil {
		return err
	}
	if !contains {
		return errors.New("subnet cidr must be within vpc cidr")
	}
	return nil
}
