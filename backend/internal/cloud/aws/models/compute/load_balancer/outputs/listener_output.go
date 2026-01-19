package outputs

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
)

// ListenerOutput represents AWS Load Balancer Listener output/response data after creation
type ListenerOutput struct {
	// AWS-generated identifiers
	ARN string `json:"arn"` // e.g., "arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/my-alb/.../..."
	ID  string `json:"id"`  // Same as ARN

	// Configuration (from input)
	LoadBalancerARN string                `json:"load_balancer_arn"`
	Port            int                   `json:"port"`
	Protocol        string                `json:"protocol"`
	DefaultAction   load_balancer.ListenerAction `json:"default_action"`
}
