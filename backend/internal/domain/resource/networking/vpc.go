package networking

import (
	"errors"
	"fmt"
	"net"
)

// VPC represents a cloud-agnostic Virtual Private Cloud
// This is the domain model - no cloud-specific details
type VPC struct {
	ID               string
	ARN              *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	Name             string
	Region           string
	CIDR             string // IPv4 CIDR block (e.g., "10.0.0.0/16")
	IPv6CIDR         *string
	EnableDNS        bool   // Enable DNS resolution
	EnableDNSHostnames bool // Enable DNS hostnames
}

// Validate performs domain-level validation
func (v *VPC) Validate() error {
	if v.Name == "" {
		return errors.New("vpc name is required")
	}
	if v.Region == "" {
		return errors.New("vpc region is required")
	}
	if v.CIDR == "" {
		return errors.New("vpc cidr is required")
	}
	
	// Validate CIDR format
	_, ipNet, err := net.ParseCIDR(v.CIDR)
	if err != nil {
		return fmt.Errorf("invalid cidr format: %w", err)
	}
	
	// Validate CIDR is IPv4 (for now)
	if ipNet.IP.To4() == nil {
		return errors.New("vpc cidr must be IPv4")
	}
	
	return nil
}

// GetCIDRBlock returns the parsed CIDR block
func (v *VPC) GetCIDRBlock() (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(v.CIDR)
	if err != nil {
		return nil, err
	}
	return ipNet, nil
}
