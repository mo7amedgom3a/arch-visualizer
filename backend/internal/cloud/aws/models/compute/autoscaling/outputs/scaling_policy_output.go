package outputs

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling"
)

// Alarm represents a CloudWatch alarm associated with a scaling policy
type Alarm struct {
	AlarmARN  string `json:"alarm_arn"`
	AlarmName string `json:"alarm_name"`
}

// ScalingPolicyOutput represents AWS scaling policy output/response data after creation
type ScalingPolicyOutput struct {
	// AWS-generated identifiers
	PolicyARN string `json:"policy_arn"`
	PolicyName string `json:"policy_name"`

	// Configuration (from input)
	AutoScalingGroupName string                      `json:"auto_scaling_group_name"`
	PolicyType           autoscaling.ScalingPolicyType `json:"policy_type"`
	TargetTrackingConfiguration *autoscaling.TargetTrackingConfiguration `json:"target_tracking_configuration,omitempty"`
	StepScalingConfiguration    *autoscaling.StepScalingConfiguration   `json:"step_scaling_configuration,omitempty"`
	SimpleScalingConfiguration  *autoscaling.SimpleScalingConfiguration  `json:"simple_scaling_configuration,omitempty"`
	AdjustmentType              *autoscaling.AdjustmentType              `json:"adjustment_type,omitempty"`
	Cooldown                    *int                                    `json:"cooldown,omitempty"`
	MinAdjustmentMagnitude      *int                                    `json:"min_adjustment_magnitude,omitempty"`

	// AWS-specific output fields
	Alarms []Alarm `json:"alarms"` // CloudWatch alarms created for this policy
}
