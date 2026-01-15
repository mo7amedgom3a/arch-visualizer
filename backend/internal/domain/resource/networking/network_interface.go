package networking

import "errors"

// NetworkInterfaceType represents the type of network interface
type NetworkInterfaceType string

const (
	// NetworkInterfaceTypeElastic Elastic network interface (manually created)
	NetworkInterfaceTypeElastic NetworkInterfaceType = "elastic"
	// NetworkInterfaceTypeAttached Attached to instance (automatically created)
	NetworkInterfaceTypeAttached NetworkInterfaceType = "attached"
)

// NetworkInterfaceStatus represents the status of a network interface
type NetworkInterfaceStatus string

const (
	NetworkInterfaceStatusAvailable NetworkInterfaceStatus = "available"
	NetworkInterfaceStatusAttaching NetworkInterfaceStatus = "attaching"
	NetworkInterfaceStatusInUse     NetworkInterfaceStatus = "in-use"
	NetworkInterfaceStatusDetaching NetworkInterfaceStatus = "detaching"
)

// NetworkInterfaceAttachment represents attachment information for a network interface
type NetworkInterfaceAttachment struct {
	AttachmentID       string  `json:"attachment_id,omitempty"`
	InstanceID         *string `json:"instance_id,omitempty"`
	DeviceIndex        *int    `json:"device_index,omitempty"`
	Status             string  `json:"status,omitempty"`
	DeleteOnTermination *bool  `json:"delete_on_termination,omitempty"`
}

// NetworkInterface represents a cloud-agnostic Network Interface (ENI)
type NetworkInterface struct {
	ID                    string                      `json:"id"`
	ARN                   *string                     `json:"arn,omitempty"` // Cloud-specific ARN
	Description           *string                     `json:"description,omitempty"`
	SubnetID              string                      `json:"subnet_id"`      // Required: Subnet in which to create the interface
	InterfaceType         NetworkInterfaceType         `json:"interface_type"` // Elastic or Attached
	PrivateIPv4Address    *string                     `json:"private_ipv4_address,omitempty"` // Custom IP or auto-assign
	AutoAssignPrivateIP   bool                        `json:"auto_assign_private_ip"`         // Auto-assign private IP
	PublicIPv4Address     *string                     `json:"public_ipv4_address,omitempty"`  // Populated after creation
	IPv6Addresses        []string                    `json:"ipv6_addresses,omitempty"`
	SecurityGroupIDs      []string                    `json:"security_group_ids"`            // Security groups to attach
	SourceDestCheck       bool                         `json:"source_dest_check"`            // Source/destination check enabled
	IPv4PrefixDelegation  *string                     `json:"ipv4_prefix_delegation,omitempty"` // IPv4 prefix delegation
	MACAddress            *string                     `json:"mac_address,omitempty"`          // Populated after creation
	Status                NetworkInterfaceStatus      `json:"status"`                        // available, attaching, in-use, detaching
	Attachment            *NetworkInterfaceAttachment `json:"attachment,omitempty"`          // Attachment info if attached to instance
	VPCID                 string                       `json:"vpc_id"`                        // Populated from subnet
	AvailabilityZone      *string                     `json:"availability_zone,omitempty"`  // Populated from subnet
}

// Validate performs domain-level validation
func (eni *NetworkInterface) Validate() error {
	if eni.SubnetID == "" {
		return errors.New("network interface subnet_id is required")
	}

	if len(eni.SecurityGroupIDs) == 0 {
		return errors.New("network interface must have at least one security group")
	}

	// If private IP is provided, auto-assign should be false
	if eni.PrivateIPv4Address != nil && *eni.PrivateIPv4Address != "" && eni.AutoAssignPrivateIP {
		return errors.New("cannot specify both private_ipv4_address and auto_assign_private_ip")
	}

	return nil
}

// IsAttached returns true if the network interface is attached to an instance
func (eni *NetworkInterface) IsAttached() bool {
	return eni.Attachment != nil && eni.Attachment.InstanceID != nil && *eni.Attachment.InstanceID != ""
}
