package sdk

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	awsautoscaling "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling/outputs"
)

// PutScalingPolicy creates or updates a scaling policy
func PutScalingPolicy(ctx context.Context, client *AWSClient, policy *awsautoscaling.ScalingPolicy) (*awsoutputs.ScalingPolicyOutput, error) {
	if err := policy.Validate(); err != nil {
		return nil, fmt.Errorf("scaling policy validation failed: %w", err)
	}

	if client == nil || client.AutoScaling == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	input := &autoscaling.PutScalingPolicyInput{
		PolicyName:           aws.String(policy.PolicyName),
		AutoScalingGroupName: aws.String(policy.AutoScalingGroupName),
		PolicyType:           aws.String(string(policy.PolicyType)),
	}

	// Set policy type specific configurations
	switch policy.PolicyType {
	case awsautoscaling.ScalingPolicyTypeTargetTrackingScaling:
		if policy.TargetTrackingConfiguration != nil {
			ttc := policy.TargetTrackingConfiguration
			ttcInput := &autoscalingtypes.TargetTrackingConfiguration{
				TargetValue: aws.Float64(ttc.TargetValue),
			}

			if ttc.DisableScaleIn != nil {
				ttcInput.DisableScaleIn = ttc.DisableScaleIn
			}

			// Set predefined metric specification
			if ttc.PredefinedMetricSpecification != nil {
				pmSpec := &autoscalingtypes.PredefinedMetricSpecification{
					PredefinedMetricType: autoscalingtypes.MetricType(ttc.PredefinedMetricSpecification.PredefinedMetricType),
				}
				if ttc.PredefinedMetricSpecification.ResourceLabel != nil {
					pmSpec.ResourceLabel = ttc.PredefinedMetricSpecification.ResourceLabel
				}
				ttcInput.PredefinedMetricSpecification = pmSpec
			}

			// Set customized metric specification
			if ttc.CustomizedMetricSpecification != nil {
				cmSpec := &autoscalingtypes.CustomizedMetricSpecification{
					MetricName: aws.String(ttc.CustomizedMetricSpecification.MetricName),
					Namespace:  aws.String(ttc.CustomizedMetricSpecification.Namespace),
					Statistic:  autoscalingtypes.MetricStatistic(ttc.CustomizedMetricSpecification.Statistic),
				}
				if len(ttc.CustomizedMetricSpecification.Dimensions) > 0 {
					cmSpec.Dimensions = make([]autoscalingtypes.MetricDimension, len(ttc.CustomizedMetricSpecification.Dimensions))
					for i, dim := range ttc.CustomizedMetricSpecification.Dimensions {
						cmSpec.Dimensions[i] = autoscalingtypes.MetricDimension{
							Name:  aws.String(dim.Name),
							Value: aws.String(dim.Value),
						}
					}
				}
				if ttc.CustomizedMetricSpecification.Unit != nil {
					cmSpec.Unit = ttc.CustomizedMetricSpecification.Unit
				}
				ttcInput.CustomizedMetricSpecification = cmSpec
			}

			input.TargetTrackingConfiguration = ttcInput
		}

	case awsautoscaling.ScalingPolicyTypeStepScaling:
		if policy.StepScalingConfiguration != nil {
			ssc := policy.StepScalingConfiguration
			// Convert AdjustmentType enum to string pointer
			adjTypeStr := string(ssc.AdjustmentType)
			input.AdjustmentType = aws.String(adjTypeStr)

			if len(ssc.StepAdjustments) > 0 {
				input.StepAdjustments = make([]autoscalingtypes.StepAdjustment, len(ssc.StepAdjustments))
				for i, stepAdj := range ssc.StepAdjustments {
					input.StepAdjustments[i] = autoscalingtypes.StepAdjustment{
						ScalingAdjustment: aws.Int32(int32(stepAdj.ScalingAdjustment)),
					}
					if stepAdj.MetricIntervalLowerBound != nil {
						input.StepAdjustments[i].MetricIntervalLowerBound = stepAdj.MetricIntervalLowerBound
					}
					if stepAdj.MetricIntervalUpperBound != nil {
						input.StepAdjustments[i].MetricIntervalUpperBound = stepAdj.MetricIntervalUpperBound
					}
				}
			}

			if ssc.MinAdjustmentMagnitude != nil {
				input.MinAdjustmentMagnitude = aws.Int32(int32(*ssc.MinAdjustmentMagnitude))
			}
			if ssc.Cooldown != nil {
				input.Cooldown = aws.Int32(int32(*ssc.Cooldown))
			}
			// Note: MetricAggregationType is not directly available in PutScalingPolicyInput
			// It's typically configured via CloudWatch alarms associated with the policy
		}

	case awsautoscaling.ScalingPolicyTypeSimpleScaling:
		if policy.SimpleScalingConfiguration != nil {
			ssc := policy.SimpleScalingConfiguration
			// Convert AdjustmentType enum to string pointer
			adjTypeStr := string(ssc.AdjustmentType)
			input.AdjustmentType = aws.String(adjTypeStr)
			input.ScalingAdjustment = aws.Int32(int32(ssc.ScalingAdjustment))
			if ssc.Cooldown != nil {
				input.Cooldown = aws.Int32(int32(*ssc.Cooldown))
			}
		}
	}

	// Set common fields
	if policy.Cooldown != nil {
		input.Cooldown = aws.Int32(int32(*policy.Cooldown))
	}
	if policy.MinAdjustmentMagnitude != nil {
		input.MinAdjustmentMagnitude = aws.Int32(int32(*policy.MinAdjustmentMagnitude))
	}

	// Create/update the scaling policy
	_, err := client.AutoScaling.PutScalingPolicy(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to put scaling policy: %w", err)
	}

	// Retrieve the policy to get full details including alarms
	return DescribeScalingPolicy(ctx, client, policy.PolicyName, policy.AutoScalingGroupName)
}

// DescribeScalingPolicy retrieves a scaling policy by name
func DescribeScalingPolicy(ctx context.Context, client *AWSClient, policyName, asgName string) (*awsoutputs.ScalingPolicyOutput, error) {
	if client == nil || client.AutoScaling == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	input := &autoscaling.DescribePoliciesInput{
		AutoScalingGroupName: aws.String(asgName),
		PolicyNames:          []string{policyName},
	}

	result, err := client.AutoScaling.DescribePolicies(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe scaling policy %s: %w", policyName, err)
	}

	if len(result.ScalingPolicies) == 0 {
		return nil, fmt.Errorf("scaling policy %s not found", policyName)
	}

	return convertScalingPolicyToOutput(&result.ScalingPolicies[0]), nil
}

// DescribeScalingPolicies lists scaling policies for an Auto Scaling Group
func DescribeScalingPolicies(ctx context.Context, client *AWSClient, asgName string) ([]*awsoutputs.ScalingPolicyOutput, error) {
	if client == nil || client.AutoScaling == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	input := &autoscaling.DescribePoliciesInput{
		AutoScalingGroupName: aws.String(asgName),
	}

	result, err := client.AutoScaling.DescribePolicies(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe scaling policies for ASG %s: %w", asgName, err)
	}

	policies := make([]*awsoutputs.ScalingPolicyOutput, len(result.ScalingPolicies))
	for i, policy := range result.ScalingPolicies {
		policies[i] = convertScalingPolicyToOutput(&policy)
	}

	return policies, nil
}

// DeleteScalingPolicy deletes a scaling policy
func DeleteScalingPolicy(ctx context.Context, client *AWSClient, policyName, asgName string) error {
	if client == nil || client.AutoScaling == nil {
		return fmt.Errorf("AWS client not available")
	}

	input := &autoscaling.DeletePolicyInput{
		PolicyName:           aws.String(policyName),
		AutoScalingGroupName: aws.String(asgName),
	}

	_, err := client.AutoScaling.DeletePolicy(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete scaling policy %s: %w", policyName, err)
	}

	return nil
}

// ExecuteScalingPolicy manually executes a scaling policy
func ExecuteScalingPolicy(ctx context.Context, client *AWSClient, policyName, asgName string, honorCooldown bool) error {
	if client == nil || client.AutoScaling == nil {
		return fmt.Errorf("AWS client not available")
	}

	input := &autoscaling.ExecutePolicyInput{
		PolicyName:           aws.String(policyName),
		AutoScalingGroupName: aws.String(asgName),
		HonorCooldown:        aws.Bool(honorCooldown),
	}

	_, err := client.AutoScaling.ExecutePolicy(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to execute scaling policy %s: %w", policyName, err)
	}

	return nil
}

// convertScalingPolicyToOutput converts AWS SDK ScalingPolicy to output model
func convertScalingPolicyToOutput(policy *autoscalingtypes.ScalingPolicy) *awsoutputs.ScalingPolicyOutput {
	output := &awsoutputs.ScalingPolicyOutput{
		PolicyARN:            aws.ToString(policy.PolicyARN),
		PolicyName:            aws.ToString(policy.PolicyName),
		AutoScalingGroupName:  aws.ToString(policy.AutoScalingGroupName),
		PolicyType:            awsautoscaling.ScalingPolicyType(aws.ToString(policy.PolicyType)),
	}

	// Convert alarms
	if len(policy.Alarms) > 0 {
		output.Alarms = make([]awsoutputs.Alarm, len(policy.Alarms))
		for i, alarm := range policy.Alarms {
			output.Alarms[i] = awsoutputs.Alarm{
				AlarmARN:  aws.ToString(alarm.AlarmARN),
				AlarmName: aws.ToString(alarm.AlarmName),
			}
		}
	}

	// Note: TargetTrackingConfiguration, StepScalingConfiguration, and SimpleScalingConfiguration
	// are not directly available in the ScalingPolicy output type from DescribePolicies.
	// They would need to be reconstructed from the policy ARN or stored separately.

	return output
}
