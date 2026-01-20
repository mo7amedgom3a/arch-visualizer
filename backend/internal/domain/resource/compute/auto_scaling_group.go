package compute

import (
	"errors"
	"fmt"
)

// AutoScalingGroupHealthCheckType represents the type of health check for Auto Scaling Group
type AutoScalingGroupHealthCheckType string

const (
	AutoScalingGroupHealthCheckTypeEC2 AutoScalingGroupHealthCheckType = "EC2"
	AutoScalingGroupHealthCheckTypeELB AutoScalingGroupHealthCheckType = "ELB"
)

// AutoScalingGroupState represents the state of an Auto Scaling Group
type AutoScalingGroupState string

const (
	AutoScalingGroupStateActive   AutoScalingGroupState = "active"
	AutoScalingGroupStateDeleting AutoScalingGroupState = "deleting"
	AutoScalingGroupStateUpdating AutoScalingGroupState = "updating"
)

// LaunchTemplateSpecification represents a reference to a Launch Template
type LaunchTemplateSpecification struct {
	ID      string  // Launch Template ID
	Version *string // Version (e.g., "$Latest", "$Default", or specific version number)
}

// Tag represents a tag with propagation flag
type Tag struct {
	Key               string
	Value             string
	PropagateAtLaunch bool
}

// AutoScalingGroup represents a cloud-agnostic Auto Scaling Group
type AutoScalingGroup struct {
	ID       string
	ARN      *string // Cloud-specific ARN
	Name     string
	Region   string
	NamePrefix *string // Optional prefix for unique naming

	// Capacity Configuration
	MinSize         int // Required: minimum number of instances
	MaxSize         int // Required: maximum number of instances
	DesiredCapacity *int // Optional: desired number of instances

	// Location Configuration
	VPCZoneIdentifier []string // Subnet IDs where instances can be created

	// Launch Configuration
	LaunchTemplate *LaunchTemplateSpecification // Reference to Launch Template

	// Health Check Configuration
	HealthCheckType        AutoScalingGroupHealthCheckType // EC2 or ELB
	HealthCheckGracePeriod *int                            // Seconds to ignore health checks after launch (default: 300)

	// Load Balancer Integration
	TargetGroupARNs []string // Target Group ARNs for ELB health checks

	// Tags
	Tags []Tag // Tags with propagation flag

	// Output Fields
	State       AutoScalingGroupState
	CreatedTime *string // ISO 8601 timestamp
}

// Validate performs domain-level validation
func (asg *AutoScalingGroup) Validate() error {
	// Name or NamePrefix must be provided
	if asg.Name == "" && (asg.NamePrefix == nil || *asg.NamePrefix == "") {
		return errors.New("auto scaling group name or name_prefix is required")
	}

	// Validate capacity constraints
	if asg.MinSize < 0 {
		return errors.New("min_size must be non-negative")
	}
	if asg.MaxSize < 0 {
		return errors.New("max_size must be non-negative")
	}
	if asg.MinSize > asg.MaxSize {
		return errors.New("min_size cannot be greater than max_size")
	}

	// Validate desired capacity if provided
	if asg.DesiredCapacity != nil {
		if *asg.DesiredCapacity < asg.MinSize {
			return fmt.Errorf("desired_capacity (%d) cannot be less than min_size (%d)", *asg.DesiredCapacity, asg.MinSize)
		}
		if *asg.DesiredCapacity > asg.MaxSize {
			return fmt.Errorf("desired_capacity (%d) cannot be greater than max_size (%d)", *asg.DesiredCapacity, asg.MaxSize)
		}
	}

	// Validate VPC Zone Identifier
	if len(asg.VPCZoneIdentifier) == 0 {
		return errors.New("vpc_zone_identifier (subnet IDs) is required")
	}

	// Validate Launch Template
	if asg.LaunchTemplate == nil {
		return errors.New("launch_template is required")
	}
	if asg.LaunchTemplate.ID == "" {
		return errors.New("launch_template.id is required")
	}

	// Validate Health Check Type
	if asg.HealthCheckType == "" {
		return errors.New("health_check_type is required")
	}
	if asg.HealthCheckType != AutoScalingGroupHealthCheckTypeEC2 && asg.HealthCheckType != AutoScalingGroupHealthCheckTypeELB {
		return errors.New("health_check_type must be EC2 or ELB")
	}

	// Validate ELB health check requires Target Group ARNs
	if asg.HealthCheckType == AutoScalingGroupHealthCheckTypeELB && len(asg.TargetGroupARNs) == 0 {
		return errors.New("target_group_arns is required when health_check_type is ELB")
	}

	// Validate Health Check Grace Period
	if asg.HealthCheckGracePeriod != nil && *asg.HealthCheckGracePeriod < 0 {
		return errors.New("health_check_grace_period must be non-negative")
	}

	// Validate Region
	if asg.Region == "" {
		return errors.New("region is required")
	}

	return nil
}

// GetID returns the Auto Scaling Group ID (implements ComputeResource interface)
func (asg *AutoScalingGroup) GetID() string {
	return asg.ID
}

// GetName returns the Auto Scaling Group name (implements ComputeResource interface)
func (asg *AutoScalingGroup) GetName() string {
	return asg.Name
}

// GetSubnetID returns the first subnet ID (implements ComputeResource interface)
func (asg *AutoScalingGroup) GetSubnetID() string {
	if len(asg.VPCZoneIdentifier) > 0 {
		return asg.VPCZoneIdentifier[0]
	}
	return ""
}
