package sdk

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	awsautoscaling "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling/outputs"
)

// CreateAutoScalingGroup creates a new Auto Scaling Group using AWS SDK
func CreateAutoScalingGroup(ctx context.Context, client *AWSClient, asg *awsautoscaling.AutoScalingGroup) (*awsoutputs.AutoScalingGroupOutput, error) {
	if err := asg.Validate(); err != nil {
		return nil, fmt.Errorf("auto scaling group validation failed: %w", err)
	}

	if client == nil || client.AutoScaling == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	// Build CreateAutoScalingGroupInput
	input := &autoscaling.CreateAutoScalingGroupInput{
		MinSize:         aws.Int32(int32(asg.MinSize)),
		MaxSize:         aws.Int32(int32(asg.MaxSize)),
		VPCZoneIdentifier: aws.String(stringSliceToCommaSeparated(asg.VPCZoneIdentifier)),
	}

	// Set name (AWS SDK v2 doesn't support name prefix in CreateAutoScalingGroupInput)
	// If only prefix is provided, we'll generate a name or use the prefix as-is
	if asg.AutoScalingGroupName != nil {
		input.AutoScalingGroupName = asg.AutoScalingGroupName
	} else if asg.AutoScalingGroupNamePrefix != nil {
		// AWS SDK v2 requires a full name, so we'll use the prefix as the name
		// In production, you might want to append a timestamp or UUID
		input.AutoScalingGroupName = asg.AutoScalingGroupNamePrefix
	}

	// Set desired capacity if provided
	if asg.DesiredCapacity != nil {
		input.DesiredCapacity = aws.Int32(int32(*asg.DesiredCapacity))
	}

	// Set Launch Template
	if asg.LaunchTemplate != nil {
		ltSpec := &autoscalingtypes.LaunchTemplateSpecification{
			LaunchTemplateId: aws.String(asg.LaunchTemplate.LaunchTemplateId),
		}
		if asg.LaunchTemplate.Version != nil {
			ltSpec.Version = asg.LaunchTemplate.Version
		} else {
			// Default to $Latest if not specified
			ltSpec.Version = aws.String("$Latest")
		}
		input.LaunchTemplate = ltSpec
	}

	// Set Health Check Type (AWS SDK expects string pointer)
	if asg.HealthCheckType != nil {
		healthCheckTypeStr := strings.ToUpper(*asg.HealthCheckType)
		if healthCheckTypeStr == "EC2" || healthCheckTypeStr == "ELB" {
			input.HealthCheckType = aws.String(healthCheckTypeStr)
		}
	}

	// Set Health Check Grace Period
	if asg.HealthCheckGracePeriod != nil {
		input.HealthCheckGracePeriod = aws.Int32(int32(*asg.HealthCheckGracePeriod))
	}

	// Set Target Group ARNs
	if len(asg.TargetGroupARNs) > 0 {
		input.TargetGroupARNs = asg.TargetGroupARNs
	}

	// Add tags
	if len(asg.Tags) > 0 {
		var tagList []autoscalingtypes.Tag
		for _, tag := range asg.Tags {
			tagList = append(tagList, autoscalingtypes.Tag{
				Key:               aws.String(tag.Key),
				Value:             aws.String(tag.Value),
				PropagateAtLaunch: aws.Bool(tag.PropagateAtLaunch),
				ResourceType:      aws.String("auto-scaling-group"),
				ResourceId:        input.AutoScalingGroupName, // Will be set after creation
			})
		}
		input.Tags = tagList
	}

	// Create the auto scaling group
	_, err := client.AutoScaling.CreateAutoScalingGroup(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create auto scaling group: %w", err)
	}

	// Retrieve the created ASG to get output fields
	asgName := ""
	if input.AutoScalingGroupName != nil {
		asgName = *input.AutoScalingGroupName
	} else if asg.AutoScalingGroupNamePrefix != nil {
		// For name prefix, use the prefix as the name (AWS SDK v2 requires full name)
		asgName = *asg.AutoScalingGroupNamePrefix
	}

	// Wait a bit for ASG to be created, then fetch it
	time.Sleep(1 * time.Second)
	return GetAutoScalingGroup(ctx, client, asgName)
}

// GetAutoScalingGroup retrieves an Auto Scaling Group by name
func GetAutoScalingGroup(ctx context.Context, client *AWSClient, name string) (*awsoutputs.AutoScalingGroupOutput, error) {
	if client == nil || client.AutoScaling == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{name},
	}

	result, err := client.AutoScaling.DescribeAutoScalingGroups(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get auto scaling group %s: %w", name, err)
	}

	if len(result.AutoScalingGroups) == 0 {
		return nil, fmt.Errorf("auto scaling group %s not found", name)
	}

	return convertAutoScalingGroupToOutput(&result.AutoScalingGroups[0]), nil
}

// UpdateAutoScalingGroup updates an Auto Scaling Group
func UpdateAutoScalingGroup(ctx context.Context, client *AWSClient, name string, asg *awsautoscaling.AutoScalingGroup) (*awsoutputs.AutoScalingGroupOutput, error) {
	if err := asg.Validate(); err != nil {
		return nil, fmt.Errorf("auto scaling group validation failed: %w", err)
	}

	if client == nil || client.AutoScaling == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	// Build UpdateAutoScalingGroupInput
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		MinSize:              aws.Int32(int32(asg.MinSize)),
		MaxSize:              aws.Int32(int32(asg.MaxSize)),
		VPCZoneIdentifier:    aws.String(stringSliceToCommaSeparated(asg.VPCZoneIdentifier)),
	}

	// Set desired capacity if provided
	if asg.DesiredCapacity != nil {
		input.DesiredCapacity = aws.Int32(int32(*asg.DesiredCapacity))
	}

	// Set Launch Template
	if asg.LaunchTemplate != nil {
		ltSpec := &autoscalingtypes.LaunchTemplateSpecification{
			LaunchTemplateId: aws.String(asg.LaunchTemplate.LaunchTemplateId),
		}
		if asg.LaunchTemplate.Version != nil {
			ltSpec.Version = asg.LaunchTemplate.Version
		} else {
			ltSpec.Version = aws.String("$Latest")
		}
		input.LaunchTemplate = ltSpec
	}

	// Set Health Check Type (AWS SDK expects string pointer)
	if asg.HealthCheckType != nil {
		healthCheckTypeStr := strings.ToUpper(*asg.HealthCheckType)
		if healthCheckTypeStr == "EC2" || healthCheckTypeStr == "ELB" {
			input.HealthCheckType = aws.String(healthCheckTypeStr)
		}
	}

	// Set Health Check Grace Period
	if asg.HealthCheckGracePeriod != nil {
		input.HealthCheckGracePeriod = aws.Int32(int32(*asg.HealthCheckGracePeriod))
	}

	// Note: Target Group ARNs cannot be updated via UpdateAutoScalingGroup
	// Use AttachLoadBalancerTargetGroups/DetachLoadBalancerTargetGroups instead

	// Update the auto scaling group
	_, err := client.AutoScaling.UpdateAutoScalingGroup(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update auto scaling group: %w", err)
	}

	// Retrieve and return the updated ASG
	return GetAutoScalingGroup(ctx, client, name)
}

// DeleteAutoScalingGroup deletes an Auto Scaling Group
func DeleteAutoScalingGroup(ctx context.Context, client *AWSClient, name string, forceDelete bool) error {
	if client == nil || client.AutoScaling == nil {
		return fmt.Errorf("AWS client not available")
	}

	input := &autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		ForceDelete:          aws.Bool(forceDelete),
	}

	_, err := client.AutoScaling.DeleteAutoScalingGroup(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete auto scaling group %s: %w", name, err)
	}

	return nil
}

// ListAutoScalingGroups lists Auto Scaling Groups with optional filters
func ListAutoScalingGroups(ctx context.Context, client *AWSClient, filters map[string][]string) ([]*awsoutputs.AutoScalingGroupOutput, error) {
	if client == nil || client.AutoScaling == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	var allASGs []*awsoutputs.AutoScalingGroupOutput
	var nextToken *string

	for {
		input := &autoscaling.DescribeAutoScalingGroupsInput{}

		// Apply filters if provided
		if names, ok := filters["name"]; ok && len(names) > 0 {
			input.AutoScalingGroupNames = names
		}

		if nextToken != nil {
			input.NextToken = nextToken
		}

		result, err := client.AutoScaling.DescribeAutoScalingGroups(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list auto scaling groups: %w", err)
		}

		// Convert each ASG to output model
		for _, asg := range result.AutoScalingGroups {
			allASGs = append(allASGs, convertAutoScalingGroupToOutput(&asg))
		}

		// Check if there are more pages
		if result.NextToken == nil {
			break
		}
		nextToken = result.NextToken
	}

	return allASGs, nil
}

// SetDesiredCapacity updates the desired capacity of an Auto Scaling Group
func SetDesiredCapacity(ctx context.Context, client *AWSClient, asgName string, capacity int) error {
	if client == nil || client.AutoScaling == nil {
		return fmt.Errorf("AWS client not available")
	}

	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(asgName),
		DesiredCapacity:      aws.Int32(int32(capacity)),
	}

	_, err := client.AutoScaling.SetDesiredCapacity(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to set desired capacity for auto scaling group %s: %w", asgName, err)
	}

	return nil
}

// AttachInstances attaches instances to an Auto Scaling Group
func AttachInstances(ctx context.Context, client *AWSClient, asgName string, instanceIDs []string) error {
	if client == nil || client.AutoScaling == nil {
		return fmt.Errorf("AWS client not available")
	}

	if len(instanceIDs) == 0 {
		return fmt.Errorf("instance IDs list cannot be empty")
	}

	input := &autoscaling.AttachInstancesInput{
		AutoScalingGroupName: aws.String(asgName),
		InstanceIds:          instanceIDs,
	}

	_, err := client.AutoScaling.AttachInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to attach instances to auto scaling group %s: %w", asgName, err)
	}

	return nil
}

// DetachInstances detaches instances from an Auto Scaling Group
func DetachInstances(ctx context.Context, client *AWSClient, asgName string, instanceIDs []string, shouldDecrementDesiredCapacity bool) error {
	if client == nil || client.AutoScaling == nil {
		return fmt.Errorf("AWS client not available")
	}

	if len(instanceIDs) == 0 {
		return fmt.Errorf("instance IDs list cannot be empty")
	}

	input := &autoscaling.DetachInstancesInput{
		AutoScalingGroupName:       aws.String(asgName),
		InstanceIds:                instanceIDs,
		ShouldDecrementDesiredCapacity: aws.Bool(shouldDecrementDesiredCapacity),
	}

	_, err := client.AutoScaling.DetachInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to detach instances from auto scaling group %s: %w", asgName, err)
	}

	return nil
}

// convertAutoScalingGroupToOutput converts AWS SDK AutoScalingGroup to output model
func convertAutoScalingGroupToOutput(asg *autoscalingtypes.AutoScalingGroup) *awsoutputs.AutoScalingGroupOutput {
	// Convert VPC Zone Identifier from comma-separated string to slice
	vpcZoneIdentifier := []string{}
	if asg.VPCZoneIdentifier != nil {
		vpcZoneIdentifier = commaSeparatedToStringSlice(*asg.VPCZoneIdentifier)
	}

	// Convert Launch Template
	var launchTemplate *awsautoscaling.LaunchTemplateSpecification
	if asg.LaunchTemplate != nil {
		launchTemplate = &awsautoscaling.LaunchTemplateSpecification{
			LaunchTemplateId: aws.ToString(asg.LaunchTemplate.LaunchTemplateId),
			Version:          asg.LaunchTemplate.Version,
		}
	}

	// Convert Health Check Type
	var healthCheckType string
	if asg.HealthCheckType != nil {
		healthCheckType = string(*asg.HealthCheckType)
	} else {
		healthCheckType = "EC2" // Default
	}

	// Convert Health Check Grace Period
	var healthCheckGracePeriod *int
	if asg.HealthCheckGracePeriod != nil {
		gracePeriod := int(*asg.HealthCheckGracePeriod)
		healthCheckGracePeriod = &gracePeriod
	}

	// Convert Instances
	instances := make([]awsoutputs.Instance, len(asg.Instances))
	for i, instance := range asg.Instances {
		healthStatus := ""
		if instance.HealthStatus != nil {
			healthStatus = *instance.HealthStatus
		}
		instances[i] = awsoutputs.Instance{
			InstanceID:       aws.ToString(instance.InstanceId),
			AvailabilityZone: aws.ToString(instance.AvailabilityZone),
			LifecycleState:   string(instance.LifecycleState),
			HealthStatus:     healthStatus,
		}
	}

	// Convert Tags
	tags := make([]awsautoscaling.Tag, len(asg.Tags))
	for i, tag := range asg.Tags {
		tags[i] = awsautoscaling.Tag{
			Key:               aws.ToString(tag.Key),
			Value:             aws.ToString(tag.Value),
			PropagateAtLaunch: aws.ToBool(tag.PropagateAtLaunch),
		}
	}

	output := &awsoutputs.AutoScalingGroupOutput{
		AutoScalingGroupARN:  aws.ToString(asg.AutoScalingGroupARN),
		AutoScalingGroupName:  aws.ToString(asg.AutoScalingGroupName),
		MinSize:              int(aws.ToInt32(asg.MinSize)),
		MaxSize:              int(aws.ToInt32(asg.MaxSize)),
		DesiredCapacity:      int(aws.ToInt32(asg.DesiredCapacity)),
		VPCZoneIdentifier:    vpcZoneIdentifier,
		LaunchTemplate:       launchTemplate,
		HealthCheckType:      healthCheckType,
		HealthCheckGracePeriod: healthCheckGracePeriod,
		TargetGroupARNs:      asg.TargetGroupARNs,
		Status:               aws.ToString(asg.Status),
		CreatedTime:          aws.ToTime(asg.CreatedTime),
		Instances:            instances,
		Tags:                 tags,
	}

	return output
}

// Helper functions

func stringSliceToCommaSeparated(slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += "," + slice[i]
	}
	return result
}

func commaSeparatedToStringSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := []string{}
	for _, part := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}
