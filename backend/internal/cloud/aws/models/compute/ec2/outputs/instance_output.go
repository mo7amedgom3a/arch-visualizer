package outputs

import (
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// InstanceOutput represents AWS EC2 instance output/response data after creation
type InstanceOutput struct {
	// AWS-generated identifiers
	ID     string `json:"id"`     // e.g., "i-1234567890abcdef0"
	ARN    string `json:"arn"`    // e.g., "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0"

	// Configuration (from input)
	Name         string `json:"name"`
	Region       string `json:"region"`
	InstanceType string `json:"instance_type"`
	AMI          string `json:"ami"`

	// AWS-specific output fields
	State            string    `json:"state"`              // pending, running, stopping, stopped, shutting-down, terminated
	AvailabilityZone string    `json:"availability_zone"`  // e.g., "us-east-1a"
	CreationTime     time.Time `json:"creation_time"`

	// Networking
	PublicIP   *string  `json:"public_ip,omitempty"`   // Public IPv4 address
	PrivateIP  string   `json:"private_ip"`            // Private IPv4 address
	PublicDNS  *string  `json:"public_dns,omitempty"`  // Public DNS hostname
	PrivateDNS string   `json:"private_dns"`           // Private DNS hostname
	SubnetID   string   `json:"subnet_id"`
	VPCID      string   `json:"vpc_id"`
	SecurityGroupIDs []string `json:"security_group_ids"`

	// Access & Permissions
	KeyName           *string `json:"key_name,omitempty"`
	IAMInstanceProfile *string `json:"iam_instance_profile,omitempty"`

	// Tags
	Tags []configs.Tag `json:"tags"`
}
