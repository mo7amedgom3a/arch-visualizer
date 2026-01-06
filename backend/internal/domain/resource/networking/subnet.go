package networking

import (
	"errors"
	"fmt"
	"net"
)

// Subnet represents a cloud-agnostic subnet
type Subnet struct {
	ID              string
	Name            string
	VPCID           string // Parent VPC
	CIDR            string // IPv4 CIDR block (must be within VPC CIDR)
	AvailabilityZone *string // Optional AZ for regional subnets
	IsPublic        bool   // Whether subnet is public-facing
}

// Validate performs domain-level validation
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
		return nil, err
	}
	return ipNet, nil
}

// ValidateCIDRInVPC validates that subnet CIDR is within VPC CIDR
func (s *Subnet) ValidateCIDRInVPC(vpcCIDR string) error {
	vpcIPNet, err := parseCIDR(vpcCIDR)
	if err != nil {
		return fmt.Errorf("invalid vpc cidr: %w", err)
	}
	
	subnetIPNet, err := s.GetCIDRBlock()
	if err != nil {
		return err
	}
	
	if !vpcIPNet.Contains(subnetIPNet.IP) {
		return errors.New("subnet cidr must be within vpc cidr")
	}
	
	// Check if subnet mask is smaller than VPC mask (more specific)
	vpcMaskSize, _ := vpcIPNet.Mask.Size()
	subnetMaskSize, _ := subnetIPNet.Mask.Size()
	if subnetMaskSize <= vpcMaskSize {
		return errors.New("subnet cidr mask must be more specific than vpc cidr mask")
	}
	
	return nil
}

func parseCIDR(cidr string) (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	return ipNet, err
}
