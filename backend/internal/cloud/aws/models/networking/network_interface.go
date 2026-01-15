package networking

import (
	"errors"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// NetworkInterfaceType represents the type of network interface
type NetworkInterfaceType string

const (
	// NetworkInterfaceTypeElastic Elastic network interface (manually created)
	NetworkInterfaceTypeElastic NetworkInterfaceType = "elastic"
	// NetworkInterfaceTypeAttached Attached to instance (automatically created)
	NetworkInterfaceTypeAttached NetworkInterfaceType = "attached"
)

// NetworkInterface represents an AWS-specific Network Interface (ENI) creation request
type NetworkInterface struct {
	Description          *string              `json:"description,omitempty"`            // Optional descriptive name
	SubnetID             string               `json:"subnet_id"`                        // Required: Subnet in which to create the interface
	InterfaceType        NetworkInterfaceType `json:"interface_type"`                   // Elastic or Attached
	PrivateIPv4Address   *string              `json:"private_ipv4_address,omitempty"`   // Custom IP or nil for auto-assign
	AutoAssignPrivateIP  bool                 `json:"auto_assign_private_ip"`           // Auto-assign private IP (default: true)
	SecurityGroupIDs     []string             `json:"security_group_ids"`               // Security groups to attach
	SourceDestCheck      *bool                `json:"source_dest_check,omitempty"`      // Source/destination check (default: true)
	IPv4PrefixDelegation *string              `json:"ipv4_prefix_delegation,omitempty"` // IPv4 prefix delegation (auto-assign, custom, or nil)
	Tags                 []configs.Tag        `json:"tags,omitempty"`
}

// Validate performs AWS-specific validation
func (eni *NetworkInterface) Validate() error {
	if eni.SubnetID == "" {
		return errors.New("network interface subnet_id is required")
	}

	if len(eni.SecurityGroupIDs) == 0 {
		return errors.New("network interface must have at least one security group")
	}

	// Validate security group ID format
	for i, sgID := range eni.SecurityGroupIDs {
		if len(sgID) < 3 || sgID[:3] != "sg-" {
			return errors.New("invalid security group id format")
		}
		if i > 0 && len(eni.SecurityGroupIDs) > 5 {
			return errors.New("network interface can have at most 5 security groups")
		}
	}

	// If private IP is provided, auto-assign should be false
	if eni.PrivateIPv4Address != nil && *eni.PrivateIPv4Address != "" && eni.AutoAssignPrivateIP {
		return errors.New("cannot specify both private_ipv4_address and auto_assign_private_ip")
	}

	// Validate interface type
	if eni.InterfaceType != NetworkInterfaceTypeElastic && eni.InterfaceType != NetworkInterfaceTypeAttached {
		return errors.New("invalid interface type, must be 'elastic' or 'attached'")
	}

	return nil
}
