package outputs

import (
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// PolicyOutput represents AWS IAM policy output/response data after creation
type PolicyOutput struct {
	// AWS-generated identifiers
	ARN string `json:"arn"` // e.g., "arn:aws:iam::123456789012:policy/my-policy"
	ID  string `json:"id"`  // Same as ARN for policies

	// Configuration (from input)
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	Path          string  `json:"path"`
	PolicyDocument string `json:"policy_document"` // JSON string

	// AWS-specific output fields
	CreateDate    time.Time     `json:"create_date"`
	UpdateDate    time.Time     `json:"update_date"`
	DefaultVersionID *string     `json:"default_version_id,omitempty"`
	AttachmentCount int          `json:"attachment_count"`
	IsAttachable   bool         `json:"is_attachable"`
	Tags           []configs.Tag `json:"tags,omitempty"`
	
	// Policy type information
	IsAWSManaged  bool          `json:"is_aws_managed"` // True if this is an AWS managed policy
}
