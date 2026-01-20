package outputs

import (
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling"
)

// Instance represents an instance in an Auto Scaling Group
type Instance struct {
	InstanceID       string `json:"instance_id"`
	AvailabilityZone string `json:"availability_zone"`
	LifecycleState   string `json:"lifecycle_state"` // Pending, InService, Terminating, etc.
	HealthStatus     string `json:"health_status"`   // Healthy, Unhealthy
}

// AutoScalingGroupOutput represents AWS Auto Scaling Group output/response data after creation
type AutoScalingGroupOutput struct {
	// AWS-generated identifiers
	AutoScalingGroupARN  string `json:"auto_scaling_group_arn"`
	AutoScalingGroupName string `json:"auto_scaling_group_name"`

	// Configuration (from input)
	MinSize         int      `json:"min_size"`
	MaxSize         int      `json:"max_size"`
	DesiredCapacity int      `json:"desired_capacity"`
	VPCZoneIdentifier []string `json:"vpc_zone_identifier"`
	LaunchTemplate  *autoscaling.LaunchTemplateSpecification `json:"launch_template"`
	HealthCheckType string   `json:"health_check_type"`
	HealthCheckGracePeriod *int `json:"health_check_grace_period"`
	TargetGroupARNs []string `json:"target_group_arns"`

	// AWS-specific output fields
	Status      string    `json:"status"`       // Active, Deleting, etc.
	CreatedTime time.Time `json:"created_time"`
	Instances   []Instance `json:"instances"`   // Current instances in ASG
	Tags        []autoscaling.Tag `json:"tags"`
}
