package outputs

import (
	"time"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// NATGatewayOutput represents AWS NAT Gateway output/response data after creation
type NATGatewayOutput struct {
	// AWS-generated identifiers
	ID     string `json:"id"`     // e.g., "nat-12345678"
	ARN    string `json:"arn"`    // e.g., "arn:aws:ec2:us-east-1:123456789012:nat-gateway/nat-12345678"
	
	// Configuration (from input)
	Name         string `json:"name"`
	SubnetID     string `json:"subnet_id"`
	AllocationID string `json:"allocation_id"` // Elastic IP allocation ID
	
	// AWS-specific output fields
	State        string    `json:"state"`         // available, pending, failed, deleting, deleted
	PublicIP     string    `json:"public_ip"`
	PrivateIP    string    `json:"private_ip"`
	CreationTime time.Time `json:"creation_time"`
	Tags         []configs.Tag `json:"tags"`
}
