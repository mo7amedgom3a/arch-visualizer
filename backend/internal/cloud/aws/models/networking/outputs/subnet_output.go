package outputs

import (
	"time"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// SubnetOutput represents AWS Subnet output/response data after creation
type SubnetOutput struct {
	// AWS-generated identifiers
	ID     string `json:"id"`     // e.g., "subnet-12345678"
	ARN    string `json:"arn"`    // e.g., "arn:aws:ec2:us-east-1:123456789012:subnet/subnet-12345678"
	
	// Configuration (from input)
	Name             string `json:"name"`
	VPCID            string `json:"vpc_id"`
	CIDR             string `json:"cidr"`
	AvailabilityZone string `json:"availability_zone"`
	
	// AWS-specific output fields
	State            string    `json:"state"`             // available, pending, etc.
	AvailableIPCount int       `json:"available_ip_count"`
	MapPublicIPOnLaunch bool   `json:"map_public_ip_on_launch"`
	CreationTime     time.Time `json:"creation_time"`
	Tags             []configs.Tag `json:"tags"`
}
