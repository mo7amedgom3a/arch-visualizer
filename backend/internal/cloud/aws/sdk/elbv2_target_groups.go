package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer/outputs"
)

// CreateTargetGroup creates a new Target Group using AWS SDK
func CreateTargetGroup(ctx context.Context, client *AWSClient, tg *awsloadbalancer.TargetGroup) (*awsoutputs.TargetGroupOutput, error) {
	if err := tg.Validate(); err != nil {
		return nil, fmt.Errorf("target group validation failed: %w", err)
	}

	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	// Build CreateTargetGroupInput
	input := &elasticloadbalancingv2.CreateTargetGroupInput{
		Name:     aws.String(tg.Name),
		Port:     aws.Int32(int32(tg.Port)),
		Protocol: elbv2types.ProtocolEnum(tg.Protocol),
		VpcId:    aws.String(tg.VPCID),
	}

	// Set target type
	if tg.TargetType != nil {
		input.TargetType = elbv2types.TargetTypeEnum(*tg.TargetType)
	} else {
		input.TargetType = elbv2types.TargetTypeEnumInstance
	}

	// Configure health check
	healthCheck := tg.HealthCheck
	input.HealthCheckEnabled = aws.Bool(true)
	input.HealthCheckPath = healthCheck.Path
	if healthCheck.Matcher != nil {
		input.Matcher = &elbv2types.Matcher{
			HttpCode: healthCheck.Matcher,
		}
	}
	if healthCheck.Interval != nil {
		input.HealthCheckIntervalSeconds = aws.Int32(int32(*healthCheck.Interval))
	}
	if healthCheck.Timeout != nil {
		input.HealthCheckTimeoutSeconds = aws.Int32(int32(*healthCheck.Timeout))
	}
	if healthCheck.HealthyThreshold != nil {
		input.HealthyThresholdCount = aws.Int32(int32(*healthCheck.HealthyThreshold))
	}
	if healthCheck.UnhealthyThreshold != nil {
		input.UnhealthyThresholdCount = aws.Int32(int32(*healthCheck.UnhealthyThreshold))
	}
	if healthCheck.Protocol != nil {
		input.HealthCheckProtocol = elbv2types.ProtocolEnum(*healthCheck.Protocol)
	}
	if healthCheck.Port != nil {
		input.HealthCheckPort = healthCheck.Port
	}

	// Add tags if provided
	if len(tg.Tags) > 0 {
		var tagList []elbv2types.Tag
		for _, tag := range tg.Tags {
			tagList = append(tagList, elbv2types.Tag{
				Key:   aws.String(tag.Key),
				Value: aws.String(tag.Value),
			})
		}
		input.Tags = tagList
	}

	// Create the target group
	result, err := client.ELBv2.CreateTargetGroup(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create target group: %w", err)
	}

	if len(result.TargetGroups) == 0 {
		return nil, fmt.Errorf("target group creation returned no target groups")
	}

	return convertTargetGroupToOutput(&result.TargetGroups[0]), nil
}

// GetTargetGroup retrieves a Target Group by ARN
func GetTargetGroup(ctx context.Context, client *AWSClient, arn string) (*awsoutputs.TargetGroupOutput, error) {
	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	input := &elasticloadbalancingv2.DescribeTargetGroupsInput{
		TargetGroupArns: []string{arn},
	}

	result, err := client.ELBv2.DescribeTargetGroups(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get target group %s: %w", arn, err)
	}

	if len(result.TargetGroups) == 0 {
		return nil, fmt.Errorf("target group %s not found", arn)
	}

	return convertTargetGroupToOutput(&result.TargetGroups[0]), nil
}

// UpdateTargetGroup updates a Target Group
func UpdateTargetGroup(ctx context.Context, client *AWSClient, arn string, tg *awsloadbalancer.TargetGroup) (*awsoutputs.TargetGroupOutput, error) {
	if err := tg.Validate(); err != nil {
		return nil, fmt.Errorf("target group validation failed: %w", err)
	}

	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	// Build ModifyTargetGroupInput
	input := &elasticloadbalancingv2.ModifyTargetGroupInput{
		TargetGroupArn: aws.String(arn),
	}

	// Update health check
	healthCheck := tg.HealthCheck
	if healthCheck.Path != nil {
		input.HealthCheckPath = healthCheck.Path
	}
	if healthCheck.Matcher != nil {
		input.Matcher = &elbv2types.Matcher{
			HttpCode: healthCheck.Matcher,
		}
	}
	if healthCheck.Interval != nil {
		input.HealthCheckIntervalSeconds = aws.Int32(int32(*healthCheck.Interval))
	}
	if healthCheck.Timeout != nil {
		input.HealthCheckTimeoutSeconds = aws.Int32(int32(*healthCheck.Timeout))
	}
	if healthCheck.HealthyThreshold != nil {
		input.HealthyThresholdCount = aws.Int32(int32(*healthCheck.HealthyThreshold))
	}
	if healthCheck.UnhealthyThreshold != nil {
		input.UnhealthyThresholdCount = aws.Int32(int32(*healthCheck.UnhealthyThreshold))
	}
	if healthCheck.Protocol != nil {
		input.HealthCheckProtocol = elbv2types.ProtocolEnum(*healthCheck.Protocol)
	}
	if healthCheck.Port != nil {
		input.HealthCheckPort = healthCheck.Port
	}

	// Update the target group
	_, err := client.ELBv2.ModifyTargetGroup(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update target group: %w", err)
	}

	// Return updated target group
	return GetTargetGroup(ctx, client, arn)
}

// DeleteTargetGroup deletes a Target Group
func DeleteTargetGroup(ctx context.Context, client *AWSClient, arn string) error {
	if client == nil || client.ELBv2 == nil {
		return fmt.Errorf("AWS client not available")
	}

	input := &elasticloadbalancingv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(arn),
	}

	_, err := client.ELBv2.DeleteTargetGroup(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete target group %s: %w", arn, err)
	}

	return nil
}

// ListTargetGroups lists Target Groups with optional filters
func ListTargetGroups(ctx context.Context, client *AWSClient, filters map[string][]string) ([]*awsoutputs.TargetGroupOutput, error) {
	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	var allTargetGroups []*awsoutputs.TargetGroupOutput
	var nextToken *string

	for {
		input := &elasticloadbalancingv2.DescribeTargetGroupsInput{}

		// Apply filters if provided
		if names, ok := filters["name"]; ok && len(names) > 0 {
			input.Names = names
		}
		if loadBalancerARNs, ok := filters["load_balancer_arn"]; ok && len(loadBalancerARNs) > 0 {
			input.LoadBalancerArn = aws.String(loadBalancerARNs[0])
		}

		if nextToken != nil {
			input.Marker = nextToken
		}

		result, err := client.ELBv2.DescribeTargetGroups(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list target groups: %w", err)
		}

		// Convert each target group to output model
		for _, tg := range result.TargetGroups {
			allTargetGroups = append(allTargetGroups, convertTargetGroupToOutput(&tg))
		}

		// Check if there are more pages
		if result.NextMarker == nil {
			break
		}
		nextToken = result.NextMarker
	}

	return allTargetGroups, nil
}

// convertTargetGroupToOutput converts AWS SDK TargetGroup to output model
func convertTargetGroupToOutput(tg *elbv2types.TargetGroup) *awsoutputs.TargetGroupOutput {
	output := &awsoutputs.TargetGroupOutput{
		ARN:        aws.ToString(tg.TargetGroupArn),
		ID:         aws.ToString(tg.TargetGroupArn),
		Name:       aws.ToString(tg.TargetGroupName),
		Port:       int(aws.ToInt32(tg.Port)),
		Protocol:   string(tg.Protocol),
		VPCID:      aws.ToString(tg.VpcId),
		TargetType: string(tg.TargetType),
		State:      "active", // Default state, actual state would require separate health check call
		CreatedTime: time.Now(), // CreatedTime not available in TargetGroup type, use current time as placeholder
	}

	// Convert health check
	if tg.HealthCheckPath != nil {
		output.HealthCheck.Path = tg.HealthCheckPath
	}
	if tg.Matcher != nil && tg.Matcher.HttpCode != nil {
		output.HealthCheck.Matcher = tg.Matcher.HttpCode
	}
	if tg.HealthCheckIntervalSeconds != nil {
		interval := int(*tg.HealthCheckIntervalSeconds)
		output.HealthCheck.Interval = &interval
	}
	if tg.HealthCheckTimeoutSeconds != nil {
		timeout := int(*tg.HealthCheckTimeoutSeconds)
		output.HealthCheck.Timeout = &timeout
	}
	if tg.HealthyThresholdCount != nil {
		healthyThreshold := int(*tg.HealthyThresholdCount)
		output.HealthCheck.HealthyThreshold = &healthyThreshold
	}
	if tg.UnhealthyThresholdCount != nil {
		unhealthyThreshold := int(*tg.UnhealthyThresholdCount)
		output.HealthCheck.UnhealthyThreshold = &unhealthyThreshold
	}
	if tg.HealthCheckProtocol != "" {
		protocol := string(tg.HealthCheckProtocol)
		output.HealthCheck.Protocol = &protocol
	}
	if tg.HealthCheckPort != nil {
		port := *tg.HealthCheckPort
		output.HealthCheck.Port = &port
	}

	return output
}
