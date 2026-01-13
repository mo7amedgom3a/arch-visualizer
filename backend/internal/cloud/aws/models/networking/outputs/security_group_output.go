package outputs

import (
	"time"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// SecurityGroupOutput represents AWS Security Group output/response data after creation
type SecurityGroupOutput struct {
	// AWS-generated identifiers
	ID     string `json:"id"`     // e.g., "sg-12345678"
	ARN    string `json:"arn"`    // e.g., "arn:aws:ec2:us-east-1:123456789012:security-group/sg-12345678"
	
	// Configuration (from input)
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	VPCID       string                    `json:"vpc_id"`
	Rules       []networking.SecurityGroupRule `json:"rules"`
	
	// AWS-specific output fields
	CreationTime time.Time     `json:"creation_time"`
	Tags         []configs.Tag `json:"tags"`
}
