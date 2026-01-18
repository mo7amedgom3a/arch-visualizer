package compute

import (
	"errors"
	"strings"
)

// InstanceState represents the state of a compute instance
type InstanceState string

const (
	InstanceStatePending      InstanceState = "pending"
	InstanceStateRunning      InstanceState = "running"
	InstanceStateStopping     InstanceState = "stopping"
	InstanceStateStopped      InstanceState = "stopped"
	InstanceStateShuttingDown InstanceState = "shutting-down"
	InstanceStateTerminated   InstanceState = "terminated"
)

// Instance represents a cloud-agnostic compute instance (EC2, VM, etc.)
// This is the domain model - no cloud-specific details
type Instance struct {
	ID               string
	ARN              *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	Name             string
	Region           string
	AvailabilityZone *string

	// Compute Configuration
	InstanceType string // e.g., "t3.micro", "m5.large"
	AMI          string // AMI ID or image identifier

	// Networking
	SubnetID         string
	SecurityGroupIDs []string
	PrivateIP        *string
	PublicIP         *string

	// Access & Permissions
	KeyName            *string
	IAMInstanceProfile *string

	// Storage
	// RootVolumeID references a storage volume resource (EBS volume)
	// If nil, a default root volume will be created by the cloud provider
	RootVolumeID *string

	// State
	State InstanceState
}

// Validate performs domain-level validation
func (i *Instance) Validate() error {
	if i.Name == "" {
		return errors.New("instance name is required")
	}
	if i.Region == "" {
		return errors.New("instance region is required")
	}
	if i.InstanceType == "" {
		return errors.New("instance type is required")
	}
	if i.AMI == "" {
		return errors.New("ami is required")
	}
	if i.SubnetID == "" {
		return errors.New("subnet id is required")
	}
	if len(i.SecurityGroupIDs) == 0 {
		return errors.New("at least one security group is required")
	}

	// Validate root volume ID format if provided
	if i.RootVolumeID != nil && *i.RootVolumeID != "" {
		if !strings.HasPrefix(*i.RootVolumeID, "vol-") {
			return errors.New("root volume id must start with 'vol-'")
		}
	}

	return nil
}
