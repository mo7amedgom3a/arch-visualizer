package outputs

import (
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
)

// TargetGroupOutput represents AWS Target Group output/response data after creation
type TargetGroupOutput struct {
	// AWS-generated identifiers
	ARN  string `json:"arn"` // e.g., "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/my-tg/..."
	ID   string `json:"id"`  // Same as ARN
	Name string `json:"name"`

	// Configuration (from input)
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	VPCID      string `json:"vpc_id"`
	TargetType string `json:"target_type"`

	// Health check configuration
	HealthCheck load_balancer.HealthCheckConfig `json:"health_check"`

	// AWS-specific output fields
	State       string    `json:"state"` // active, draining, deleting, deleted
	CreatedTime time.Time `json:"created_time"`
}
