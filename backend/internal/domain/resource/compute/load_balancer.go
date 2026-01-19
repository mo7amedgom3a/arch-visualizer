package compute

import (
	"errors"
)

// LoadBalancerType represents the type of load balancer
type LoadBalancerType string

const (
	LoadBalancerTypeApplication LoadBalancerType = "application"
	LoadBalancerTypeNetwork     LoadBalancerType = "network"
)

// LoadBalancerState represents the state of a load balancer
type LoadBalancerState string

const (
	LoadBalancerStateActive       LoadBalancerState = "active"
	LoadBalancerStateProvisioning LoadBalancerState = "provisioning"
	LoadBalancerStateActiveImpaired LoadBalancerState = "active_impaired"
	LoadBalancerStateFailed       LoadBalancerState = "failed"
)

// LoadBalancer represents a cloud-agnostic load balancer
type LoadBalancer struct {
	ID               string
	ARN              *string // Cloud-specific ARN
	Name             string
	Region           string
	Type             LoadBalancerType
	Internal         bool
	SecurityGroupIDs []string
	SubnetIDs        []string
	DNSName          *string // Output: auto-generated DNS name
	ZoneID           *string // Output: Route53 hosted zone ID
	State            LoadBalancerState
}

// Validate performs domain-level validation
func (lb *LoadBalancer) Validate() error {
	if lb.Name == "" {
		return errors.New("load balancer name is required")
	}
	if lb.Region == "" {
		return errors.New("load balancer region is required")
	}
	if lb.Type == "" {
		return errors.New("load balancer type is required")
	}
	if lb.Type != LoadBalancerTypeApplication && lb.Type != LoadBalancerTypeNetwork {
		return errors.New("load balancer type must be 'application' or 'network'")
	}
	if len(lb.SubnetIDs) < 2 {
		return errors.New("at least 2 subnets are required for load balancer")
	}
	if lb.Type == LoadBalancerTypeApplication && len(lb.SecurityGroupIDs) == 0 {
		return errors.New("at least one security group is required for application load balancer")
	}
	return nil
}
