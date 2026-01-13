package outputs

import (
	"time"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// VPCOutput represents AWS VPC output/response data after creation
type VPCOutput struct {
	// AWS-generated identifiers
	ID     string `json:"id"`     // e.g., "vpc-12345678"
	ARN    string `json:"arn"`    // e.g., "arn:aws:ec2:us-east-1:123456789012:vpc/vpc-12345678"
	
	// Configuration (from input)
	Name   string `json:"name"`
	Region string `json:"region"`
	CIDR   string `json:"cidr"`
	
	// AWS-specific output fields
	State         string    `json:"state"`          // available, pending, etc.
	IsDefault     bool      `json:"is_default"`
	CreationTime  time.Time `json:"creation_time"`
	OwnerID       string    `json:"owner_id"`
	
	// Configuration fields
	EnableDNSHostnames bool       `json:"enable_dns_hostnames"`
	EnableDNSSupport   bool       `json:"enable_dns_support"`
	InstanceTenancy    string     `json:"instance_tenancy"`
	Tags               []configs.Tag `json:"tags"`
}
