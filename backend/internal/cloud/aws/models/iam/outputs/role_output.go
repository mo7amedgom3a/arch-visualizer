package outputs

import (
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// RoleOutput represents AWS IAM role output/response data after creation
type RoleOutput struct {
	// AWS-generated identifiers
	ARN      string `json:"arn"`       // e.g., "arn:aws:iam::123456789012:role/my-role"
	ID       string `json:"id"`         // Same as name for roles
	Name     string `json:"name"`
	UniqueID string `json:"unique_id"` // Stable unique identifier

	// Configuration (from input)
	Description        *string `json:"description,omitempty"`
	Path               string  `json:"path"`
	AssumeRolePolicy   string  `json:"assume_role_policy"` // JSON string
	PermissionsBoundary *string `json:"permissions_boundary,omitempty"`

	// AWS-specific output fields
	CreateDate         time.Time     `json:"create_date"`
	MaxSessionDuration *int          `json:"max_session_duration,omitempty"`
	Tags               []configs.Tag `json:"tags,omitempty"`
}
