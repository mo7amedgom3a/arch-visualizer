package outputs

import (
	"time"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// NetworkACLAssociation represents a Network ACL association with a subnet
type NetworkACLAssociation struct {
	ID       string `json:"id"`        // Association ID
	SubnetID string `json:"subnet_id"`
}

// NetworkACLOutput represents AWS Network ACL output/response data after creation
type NetworkACLOutput struct {
	// AWS-generated identifiers
	ID     string `json:"id"`     // e.g., "acl-12345678"
	ARN    string `json:"arn"`    // e.g., "arn:aws:ec2:us-east-1:123456789012:network-acl/acl-12345678"
	
	// Configuration (from input)
	Name          string              `json:"name"`
	VPCID         string              `json:"vpc_id"`
	InboundRules  []networking.ACLRule `json:"inbound_rules"`
	OutboundRules []networking.ACLRule `json:"outbound_rules"`
	
	// AWS-specific output fields
	IsDefault     bool                  `json:"is_default"`     // Whether this is the default ACL for the VPC
	Associations  []NetworkACLAssociation `json:"associations"`  // Associated subnet IDs
	CreationTime  time.Time             `json:"creation_time"`
	Tags          []configs.Tag         `json:"tags"`
}
