package outputs

import (
	"time"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// InternetGatewayOutput represents AWS Internet Gateway output/response data after creation
type InternetGatewayOutput struct {
	// AWS-generated identifiers
	ID     string `json:"id"`     // e.g., "igw-12345678"
	ARN    string `json:"arn"`    // e.g., "arn:aws:ec2:us-east-1:123456789012:internet-gateway/igw-12345678"
	
	// Configuration (from input)
	Name   string `json:"name"`
	VPCID  string `json:"vpc_id"`
	
	// AWS-specific output fields
	State          string    `json:"state"`           // available, attached, detaching
	AttachmentState string   `json:"attachment_state"` // attached, detached
	CreationTime   time.Time `json:"creation_time"`
	Tags           []configs.Tag `json:"tags"`
}
