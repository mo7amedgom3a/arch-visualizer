package outputs

import (
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// GroupOutput represents AWS IAM group output/response data after creation
type GroupOutput struct {
	// AWS-generated identifiers
	ARN      string `json:"arn"`       // e.g., "arn:aws:iam::123456789012:group/my-group"
	ID       string `json:"id"`         // Same as name for groups
	Name     string `json:"name"`
	UniqueID string `json:"unique_id"` // Stable unique identifier

	// Configuration (from input)
	Path string `json:"path"`

	// AWS-specific output fields
	CreateDate time.Time     `json:"create_date"`
	Tags       []configs.Tag `json:"tags,omitempty"`
}
