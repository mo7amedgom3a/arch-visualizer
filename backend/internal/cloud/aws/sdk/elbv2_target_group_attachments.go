package sdk

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer/outputs"
)

// AttachTargetToGroup attaches a target to a Target Group
func AttachTargetToGroup(ctx context.Context, client *AWSClient, attachment *awsloadbalancer.TargetGroupAttachment) error {
	if err := attachment.Validate(); err != nil {
		return fmt.Errorf("target group attachment validation failed: %w", err)
	}

	if client == nil || client.ELBv2 == nil {
		return fmt.Errorf("AWS client not available")
	}

	// Build RegisterTargetsInput
	target := elbv2types.TargetDescription{
		Id: aws.String(attachment.TargetID),
	}

	if attachment.Port != nil {
		target.Port = aws.Int32(int32(*attachment.Port))
	}

	if attachment.AvailabilityZone != nil {
		target.AvailabilityZone = attachment.AvailabilityZone
	}

	input := &elasticloadbalancingv2.RegisterTargetsInput{
		TargetGroupArn: aws.String(attachment.TargetGroupARN),
		Targets:        []elbv2types.TargetDescription{target},
	}

	// Register the target
	_, err := client.ELBv2.RegisterTargets(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to attach target to group: %w", err)
	}

	return nil
}

// DetachTargetFromGroup detaches a target from a Target Group
func DetachTargetFromGroup(ctx context.Context, client *AWSClient, targetGroupARN, targetID string) error {
	if client == nil || client.ELBv2 == nil {
		return fmt.Errorf("AWS client not available")
	}

	input := &elasticloadbalancingv2.DeregisterTargetsInput{
		TargetGroupArn: aws.String(targetGroupARN),
		Targets: []elbv2types.TargetDescription{
			{Id: aws.String(targetID)},
		},
	}

	_, err := client.ELBv2.DeregisterTargets(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to detach target from group: %w", err)
	}

	return nil
}

// ListTargetGroupTargets lists targets in a Target Group
func ListTargetGroupTargets(ctx context.Context, client *AWSClient, targetGroupARN string) ([]*awsoutputs.TargetGroupAttachmentOutput, error) {
	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	input := &elasticloadbalancingv2.DescribeTargetHealthInput{
		TargetGroupArn: aws.String(targetGroupARN),
	}

	result, err := client.ELBv2.DescribeTargetHealth(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list target group targets: %w", err)
	}

	var attachments []*awsoutputs.TargetGroupAttachmentOutput
	for _, health := range result.TargetHealthDescriptions {
		var port *int
		if health.Target.Port != nil {
			p := int(*health.Target.Port)
			port = &p
		}

		attachment := &awsoutputs.TargetGroupAttachmentOutput{
			TargetGroupARN:   targetGroupARN,
			TargetID:         aws.ToString(health.Target.Id),
			Port:             port,
			AvailabilityZone: health.Target.AvailabilityZone,
			HealthStatus:     string(health.TargetHealth.State),
			State:            string(health.TargetHealth.State),
		}

		// Map health state to health status
		switch health.TargetHealth.State {
		case elbv2types.TargetHealthStateEnumHealthy:
			attachment.HealthStatus = "healthy"
		case elbv2types.TargetHealthStateEnumUnhealthy:
			attachment.HealthStatus = "unhealthy"
		case elbv2types.TargetHealthStateEnumInitial:
			attachment.HealthStatus = "initial"
		case elbv2types.TargetHealthStateEnumDraining:
			attachment.HealthStatus = "draining"
		case elbv2types.TargetHealthStateEnumUnused:
			attachment.HealthStatus = "unused"
		case elbv2types.TargetHealthStateEnumUnavailable:
			attachment.HealthStatus = "unavailable"
		}

		attachments = append(attachments, attachment)
	}

	return attachments, nil
}
