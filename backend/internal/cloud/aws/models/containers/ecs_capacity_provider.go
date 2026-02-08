package containers

import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"

// ECSCapacityProvider represents an AWS ECS Capacity Provider
type ECSCapacityProvider struct {
	Name                     string                    `json:"name"`
	AutoScalingGroupProvider *AutoScalingGroupProvider `json:"auto_scaling_group_provider,omitempty"`
	Tags                     []configs.Tag             `json:"tags,omitempty"`
}

// AutoScalingGroupProvider defines the ASG connection
type AutoScalingGroupProvider struct {
	AutoScalingGroupARN          string          `json:"auto_scaling_group_arn"`
	ManagedScaling               *ManagedScaling `json:"managed_scaling,omitempty"`
	ManagedTerminationProtection string          `json:"managed_termination_protection,omitempty"` // ENABLED, DISABLED
	ManagedDraining              string          `json:"managed_draining,omitempty"`               // ENABLED, DISABLED
}

// ManagedScaling defines ECS-managed scaling behavior
type ManagedScaling struct {
	Status                 string `json:"status,omitempty"`          // ENABLED, DISABLED
	TargetCapacity         int    `json:"target_capacity,omitempty"` // 1-100
	MinimumScalingStepSize int    `json:"minimum_scaling_step_size,omitempty"`
	MaximumScalingStepSize int    `json:"maximum_scaling_step_size,omitempty"`
	InstanceWarmupPeriod   int    `json:"instance_warmup_period,omitempty"` // Seconds
}
