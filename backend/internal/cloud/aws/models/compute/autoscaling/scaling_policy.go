package autoscaling

import (
	"errors"
	"fmt"
)

// ScalingPolicyType represents the type of scaling policy
type ScalingPolicyType string

const (
	ScalingPolicyTypeTargetTrackingScaling ScalingPolicyType = "TargetTrackingScaling"
	ScalingPolicyTypeStepScaling            ScalingPolicyType = "StepScaling"
	ScalingPolicyTypeSimpleScaling          ScalingPolicyType = "SimpleScaling"
)

// AdjustmentType represents how to adjust capacity
type AdjustmentType string

const (
	AdjustmentTypeChangeInCapacity        AdjustmentType = "ChangeInCapacity"
	AdjustmentTypeExactCapacity           AdjustmentType = "ExactCapacity"
	AdjustmentTypePercentChangeInCapacity AdjustmentType = "PercentChangeInCapacity"
)

// TargetTrackingConfiguration represents target tracking scaling configuration
type TargetTrackingConfiguration struct {
	TargetValue     float64  `json:"target_value"`     // Required: target metric value
	PredefinedMetricSpecification *PredefinedMetricSpecification `json:"predefined_metric_specification,omitempty"`
	CustomizedMetricSpecification *CustomizedMetricSpecification `json:"customized_metric_specification,omitempty"`
	DisableScaleIn  *bool    `json:"disable_scale_in,omitempty"` // Optional: prevent scale-in
}

// PredefinedMetricSpecification represents a predefined metric for target tracking
type PredefinedMetricSpecification struct {
	PredefinedMetricType string  `json:"predefined_metric_type"` // e.g., "ASGAverageCPUUtilization"
	ResourceLabel        *string `json:"resource_label,omitempty"`
}

// CustomizedMetricSpecification represents a custom metric for target tracking
type CustomizedMetricSpecification struct {
	MetricName string   `json:"metric_name"`
	Namespace  string   `json:"namespace"`
	Statistic  string   `json:"statistic"` // Average, Sum, Minimum, Maximum, SampleCount
	Dimensions []MetricDimension `json:"dimensions,omitempty"`
	Unit       *string  `json:"unit,omitempty"`
}

// MetricDimension represents a metric dimension
type MetricDimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// StepAdjustment represents a step adjustment for step scaling
type StepAdjustment struct {
	MetricIntervalLowerBound *float64 `json:"metric_interval_lower_bound,omitempty"`
	MetricIntervalUpperBound *float64 `json:"metric_interval_upper_bound,omitempty"`
	ScalingAdjustment        int      `json:"scaling_adjustment"` // Required
}

// StepScalingConfiguration represents step scaling configuration
type StepScalingConfiguration struct {
	AdjustmentType       AdjustmentType   `json:"adjustment_type"`       // Required
	StepAdjustments      []StepAdjustment `json:"step_adjustments"`      // Required
	MinAdjustmentMagnitude *int           `json:"min_adjustment_magnitude,omitempty"`
	Cooldown             *int            `json:"cooldown,omitempty"`    // Seconds
	MetricAggregationType *string         `json:"metric_aggregation_type,omitempty"` // Average, Minimum, Maximum
}

// SimpleScalingConfiguration represents simple scaling configuration
type SimpleScalingConfiguration struct {
	AdjustmentType       AdjustmentType `json:"adjustment_type"`       // Required
	ScalingAdjustment    int            `json:"scaling_adjustment"`    // Required
	Cooldown             *int           `json:"cooldown,omitempty"`    // Seconds
}

// ScalingPolicy represents an AWS Auto Scaling scaling policy configuration
type ScalingPolicy struct {
	PolicyName           string                      `json:"policy_name"`            // Required
	AutoScalingGroupName string                      `json:"auto_scaling_group_name"` // Required
	PolicyType           ScalingPolicyType           `json:"policy_type"`            // Required
	TargetTrackingConfiguration *TargetTrackingConfiguration `json:"target_tracking_configuration,omitempty"`
	StepScalingConfiguration    *StepScalingConfiguration   `json:"step_scaling_configuration,omitempty"`
	SimpleScalingConfiguration  *SimpleScalingConfiguration  `json:"simple_scaling_configuration,omitempty"`
	AdjustmentType              *AdjustmentType              `json:"adjustment_type,omitempty"` // For backward compatibility
	Cooldown                    *int                        `json:"cooldown,omitempty"`
	MinAdjustmentMagnitude      *int                        `json:"min_adjustment_magnitude,omitempty"`
}

// Validate performs validation for scaling policy
func (sp *ScalingPolicy) Validate() error {
	if sp.PolicyName == "" {
		return errors.New("policy_name is required")
	}
	if len(sp.PolicyName) > 255 {
		return errors.New("policy_name cannot exceed 255 characters")
	}

	if sp.AutoScalingGroupName == "" {
		return errors.New("auto_scaling_group_name is required")
	}

	if sp.PolicyType == "" {
		return errors.New("policy_type is required")
	}

	// Validate policy type specific configurations
	switch sp.PolicyType {
	case ScalingPolicyTypeTargetTrackingScaling:
		if sp.TargetTrackingConfiguration == nil {
			return errors.New("target_tracking_configuration is required for TargetTrackingScaling policy")
		}
		if sp.TargetTrackingConfiguration.TargetValue <= 0 {
			return errors.New("target_tracking_configuration.target_value must be greater than 0")
		}
		if sp.TargetTrackingConfiguration.PredefinedMetricSpecification == nil &&
			sp.TargetTrackingConfiguration.CustomizedMetricSpecification == nil {
			return errors.New("either predefined_metric_specification or customized_metric_specification is required")
		}

	case ScalingPolicyTypeStepScaling:
		if sp.StepScalingConfiguration == nil {
			return errors.New("step_scaling_configuration is required for StepScaling policy")
		}
		if len(sp.StepScalingConfiguration.StepAdjustments) == 0 {
			return errors.New("step_scaling_configuration.step_adjustments cannot be empty")
		}
		if sp.StepScalingConfiguration.AdjustmentType == "" {
			return errors.New("step_scaling_configuration.adjustment_type is required")
		}

	case ScalingPolicyTypeSimpleScaling:
		if sp.SimpleScalingConfiguration == nil {
			return errors.New("simple_scaling_configuration is required for SimpleScaling policy")
		}
		if sp.SimpleScalingConfiguration.AdjustmentType == "" {
			return errors.New("simple_scaling_configuration.adjustment_type is required")
		}

	default:
		return fmt.Errorf("invalid policy_type: %s (must be TargetTrackingScaling, StepScaling, or SimpleScaling)", sp.PolicyType)
	}

	return nil
}
