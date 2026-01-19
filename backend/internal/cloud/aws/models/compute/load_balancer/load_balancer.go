package load_balancer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// LoadBalancer represents an AWS Load Balancer configuration
type LoadBalancer struct {
	Name             string        `json:"name"`              // Required
	LoadBalancerType string        `json:"load_balancer_type"` // Required: "application" or "network"
	Internal         *bool         `json:"internal,omitempty"` // Optional, default false
	SecurityGroupIDs []string      `json:"security_group_ids"` // Required for ALB
	SubnetIDs        []string      `json:"subnet_ids"`        // Required, at least 2 in different AZs
	IPAddressType    *string       `json:"ip_address_type,omitempty"` // Optional: "ipv4" or "dualstack"
	Tags             []configs.Tag `json:"tags,omitempty"`
}

// Validate performs AWS-specific validation
func (lb *LoadBalancer) Validate() error {
	if lb.Name == "" {
		return errors.New("load balancer name is required")
	}

	// Validate name format based on type
	lbType := strings.ToLower(lb.LoadBalancerType)
	maxLength := 32 // ALB default
	if lbType == "network" {
		maxLength = 80 // NLB
	}

	if len(lb.Name) < 1 || len(lb.Name) > maxLength {
		return fmt.Errorf("load balancer name must be between 1 and %d characters", maxLength)
	}

	// Name pattern: alphanumeric and hyphens
	namePattern := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	if !namePattern.MatchString(lb.Name) {
		return errors.New("load balancer name can only contain alphanumeric characters and hyphens")
	}

	// Validate type
	if lb.LoadBalancerType == "" {
		return errors.New("load balancer type is required")
	}
	if lbType != "application" && lbType != "network" {
		return errors.New("load balancer type must be 'application' or 'network'")
	}

	// Validate subnets
	if len(lb.SubnetIDs) < 2 {
		return errors.New("at least 2 subnets are required for load balancer")
	}

	// Validate subnet IDs format
	for i, subnetID := range lb.SubnetIDs {
		if !strings.HasPrefix(subnetID, "subnet-") {
			return fmt.Errorf("subnet ID at index %d must start with 'subnet-'", i)
		}
	}

	// Validate security groups for ALB
	if lbType == "application" && len(lb.SecurityGroupIDs) == 0 {
		return errors.New("at least one security group is required for application load balancer")
	}

	// Validate security group IDs format
	for i, sgID := range lb.SecurityGroupIDs {
		if !strings.HasPrefix(sgID, "sg-") {
			return fmt.Errorf("security group ID at index %d must start with 'sg-'", i)
		}
	}

	// Validate IP address type if provided
	if lb.IPAddressType != nil {
		ipType := strings.ToLower(*lb.IPAddressType)
		if ipType != "ipv4" && ipType != "dualstack" {
			return errors.New("ip address type must be 'ipv4' or 'dualstack'")
		}
	}

	return nil
}
