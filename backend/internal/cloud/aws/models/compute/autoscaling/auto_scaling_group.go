package autoscaling

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// LaunchTemplateSpecification represents AWS Launch Template specification
type LaunchTemplateSpecification struct {
	LaunchTemplateId string  `json:"launch_template_id"` // Required
	Version          *string `json:"version,omitempty"`  // Optional: "$Latest", "$Default", or version number
}

// AutoScalingGroup represents an AWS Auto Scaling Group configuration
type AutoScalingGroup struct {
	AutoScalingGroupName      *string  `json:"auto_scaling_group_name,omitempty"`      // Optional: exact name
	AutoScalingGroupNamePrefix *string `json:"auto_scaling_group_name_prefix,omitempty"` // Optional: name prefix
	MinSize                   int      `json:"min_size"`                               // Required
	MaxSize                   int      `json:"max_size"`                               // Required
	DesiredCapacity           *int     `json:"desired_capacity,omitempty"`             // Optional
	VPCZoneIdentifier         []string `json:"vpc_zone_identifier"`                    // Required: subnet IDs
	LaunchTemplate            *LaunchTemplateSpecification `json:"launch_template"`  // Required
	HealthCheckType           *string  `json:"health_check_type,omitempty"`           // Optional: "EC2" or "ELB" (default: "EC2")
	HealthCheckGracePeriod    *int     `json:"health_check_grace_period,omitempty"`  // Optional: seconds (default: 300)
	TargetGroupARNs          []string `json:"target_group_arns,omitempty"`          // Optional: required for ELB health checks
	Tags                     []Tag    `json:"tags,omitempty"`                        // Optional: tags with propagation
}

// Tag represents an AWS tag with propagation flag
type Tag struct {
	Key               string `json:"key"`
	Value             string `json:"value"`
	PropagateAtLaunch bool   `json:"propagate_at_launch"`
}

// Validate performs AWS-specific validation
func (asg *AutoScalingGroup) Validate() error {
	// Name or NamePrefix must be provided
	if (asg.AutoScalingGroupName == nil || *asg.AutoScalingGroupName == "") &&
		(asg.AutoScalingGroupNamePrefix == nil || *asg.AutoScalingGroupNamePrefix == "") {
		return errors.New("auto_scaling_group_name or auto_scaling_group_name_prefix is required")
	}

	// Validate name format if provided
	if asg.AutoScalingGroupName != nil && *asg.AutoScalingGroupName != "" {
		name := *asg.AutoScalingGroupName
		if len(name) < 1 || len(name) > 255 {
			return errors.New("auto_scaling_group_name must be between 1 and 255 characters")
		}
		// AWS allows alphanumeric, hyphens, underscores, and periods
		namePattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
		if !namePattern.MatchString(name) {
			return errors.New("auto_scaling_group_name can only contain alphanumeric characters, hyphens, underscores, and periods")
		}
	}

	// Validate name prefix format if provided
	if asg.AutoScalingGroupNamePrefix != nil && *asg.AutoScalingGroupNamePrefix != "" {
		prefix := *asg.AutoScalingGroupNamePrefix
		if len(prefix) < 1 || len(prefix) > 255 {
			return errors.New("auto_scaling_group_name_prefix must be between 1 and 255 characters")
		}
		prefixPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
		if !prefixPattern.MatchString(prefix) {
			return errors.New("auto_scaling_group_name_prefix can only contain alphanumeric characters, hyphens, underscores, and periods")
		}
	}

	// Validate capacity constraints
	if asg.MinSize < 0 {
		return errors.New("min_size must be non-negative")
	}
	if asg.MaxSize < 0 {
		return errors.New("max_size must be non-negative")
	}
	if asg.MaxSize > 10000 {
		return errors.New("max_size cannot exceed 10000")
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
	for _, subnetID := range asg.VPCZoneIdentifier {
		if !strings.HasPrefix(subnetID, "subnet-") {
			return fmt.Errorf("invalid subnet ID format: %s (must start with 'subnet-')", subnetID)
		}
	}

	// Validate Launch Template
	if asg.LaunchTemplate == nil {
		return errors.New("launch_template is required")
	}
	if asg.LaunchTemplate.LaunchTemplateId == "" {
		return errors.New("launch_template.launch_template_id is required")
	}
	if !strings.HasPrefix(asg.LaunchTemplate.LaunchTemplateId, "lt-") {
		return errors.New("launch_template_id must start with 'lt-'")
	}

	// Validate Health Check Type
	if asg.HealthCheckType != nil {
		healthCheckType := strings.ToUpper(*asg.HealthCheckType)
		if healthCheckType != "EC2" && healthCheckType != "ELB" {
			return errors.New("health_check_type must be EC2 or ELB")
		}

		// Validate ELB health check requires Target Group ARNs
		if healthCheckType == "ELB" && len(asg.TargetGroupARNs) == 0 {
			return errors.New("target_group_arns is required when health_check_type is ELB")
		}
	}

	// Validate Health Check Grace Period
	if asg.HealthCheckGracePeriod != nil {
		if *asg.HealthCheckGracePeriod < 0 {
			return errors.New("health_check_grace_period must be non-negative")
		}
		if *asg.HealthCheckGracePeriod > 7200 {
			return errors.New("health_check_grace_period cannot exceed 7200 seconds")
		}
	}

	// Validate Target Group ARNs format
	for _, arn := range asg.TargetGroupARNs {
		if !strings.HasPrefix(arn, "arn:aws:elasticloadbalancing:") {
			return fmt.Errorf("invalid target group ARN format: %s", arn)
		}
		if !strings.Contains(arn, ":targetgroup/") {
			return fmt.Errorf("invalid target group ARN format: %s (must contain ':targetgroup/')", arn)
		}
	}

	return nil
}
