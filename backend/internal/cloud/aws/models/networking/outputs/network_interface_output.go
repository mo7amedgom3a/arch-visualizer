package outputs

import (
	"time"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// NetworkInterfaceAttachment represents attachment information for a network interface
type NetworkInterfaceAttachment struct {
	AttachmentID        string    `json:"attachment_id,omitempty"`
	InstanceID          *string    `json:"instance_id,omitempty"`
	DeviceIndex         *int       `json:"device_index,omitempty"`
	Status              string     `json:"status,omitempty"` // attaching, attached, detaching, detached
	DeleteOnTermination *bool      `json:"delete_on_termination,omitempty"`
	AttachmentTime      *time.Time `json:"attachment_time,omitempty"`
}

// NetworkInterfaceOutput represents AWS Network Interface output/response data after creation
type NetworkInterfaceOutput struct {
	// AWS-generated identifiers
	ID     string `json:"id"`     // e.g., "eni-0805e5fd003cda740"
	ARN    string `json:"arn"`    // e.g., "arn:aws:ec2:us-east-1:123456789012:network-interface/eni-0805e5fd003cda740"
	
	// Configuration (from input)
	Description          *string   `json:"description,omitempty"`
	SubnetID             string    `json:"subnet_id"`
	InterfaceType        string    `json:"interface_type"` // "elastic" or "attached"
	SecurityGroupIDs     []string  `json:"security_group_ids"`
	
	// AWS-specific output fields
	Status               string    `json:"status"`                // available, attaching, in-use, detaching
	VPCID                string    `json:"vpc_id"`
	AvailabilityZone     string    `json:"availability_zone"`
	OwnerID              string    `json:"owner_id"`
	RequesterID          *string   `json:"requester_id,omitempty"`
	RequesterManaged     bool      `json:"requester_managed"`
	SourceDestCheck      bool      `json:"source_dest_check"`
	PrivateIPv4Address   string    `json:"private_ipv4_address"`  // Primary private IP
	PublicIPv4Address    *string   `json:"public_ipv4_address,omitempty"`
	SecondaryPrivateIPv4Addresses []string `json:"secondary_private_ipv4_addresses,omitempty"`
	SecondaryPublicIPv4Addresses  []string `json:"secondary_public_ipv4_addresses,omitempty"`
	IPv6Addresses        []string  `json:"ipv6_addresses,omitempty"`
	MACAddress           string    `json:"mac_address"`
	IPv4PrefixDelegation *string   `json:"ipv4_prefix_delegation,omitempty"`
	Attachment           *NetworkInterfaceAttachment `json:"attachment,omitempty"`
	CreationTime         time.Time `json:"creation_time"`
	Tags                 []configs.Tag `json:"tags"`
}
