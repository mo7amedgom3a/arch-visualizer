package outputs

import (
	"time"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ElasticIPOutput represents AWS Elastic IP output/response data after allocation
type ElasticIPOutput struct {
	// AWS-generated identifiers
	ID           string `json:"id"`            // Allocation ID (e.g., "eipalloc-12345678")
	ARN          string `json:"arn"`           // e.g., "arn:aws:ec2:us-east-1:123456789012:elastic-ip/eipalloc-12345678"
	PublicIP     string `json:"public_ip"`     // Public IPv4 address (e.g., "54.123.45.67")
	
	// Configuration (from input)
	Region             string `json:"region"`
	NetworkBorderGroup *string `json:"network_border_group,omitempty"`
	
	// AWS-specific output fields
	Domain             string    `json:"domain"`              // "vpc" or "standard"
	AssociationID      *string   `json:"association_id,omitempty"` // If associated with an instance/interface
	InstanceID         *string   `json:"instance_id,omitempty"`    // Associated EC2 instance ID
	NetworkInterfaceID *string   `json:"network_interface_id,omitempty"` // Associated network interface ID
	PrivateIPAddress   *string   `json:"private_ip_address,omitempty"`   // Private IP if associated
	AllocationID       string    `json:"allocation_id"`        // Same as ID
	State              string    `json:"state"`               // "available", "in-use", "released"
	CreationTime       time.Time `json:"creation_time"`
	Tags               []configs.Tag `json:"tags"`
}
