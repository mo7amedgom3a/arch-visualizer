package outputs

import (
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// InstanceProfileOutput represents AWS IAM instance profile output/response data after creation
type InstanceProfileOutput struct {
	// AWS-generated identifiers
	ARN string `json:"arn"` // e.g., "arn:aws:iam::123456789012:instance-profile/my-profile"
	ID  string `json:"id"`  // Same as name for instance profiles
	Name string `json:"name"`

	// Configuration (from input)
	Path string `json:"path"`

	// AWS-specific output fields
	CreateDate time.Time     `json:"create_date"`
	Tags       []configs.Tag `json:"tags,omitempty"`

	// Attached roles
	Roles []*RoleOutput `json:"roles,omitempty"` // Roles attached to this instance profile
}
