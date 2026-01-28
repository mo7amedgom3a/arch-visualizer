package compute

import (
	"strings"

	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	awsautoscaling "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling/outputs"
)

// FromDomainAutoScalingGroup converts domain AutoScalingGroup to AWS AutoScalingGroup
func FromDomainAutoScalingGroup(domainASG *domaincompute.AutoScalingGroup) *awsautoscaling.AutoScalingGroup {
	if domainASG == nil {
		return nil
	}

	awsASG := &awsautoscaling.AutoScalingGroup{
		MinSize:         domainASG.MinSize,
		MaxSize:         domainASG.MaxSize,
		DesiredCapacity: domainASG.DesiredCapacity,
		VPCZoneIdentifier: domainASG.VPCZoneIdentifier,
	}

	// Set name or name prefix
	if domainASG.Name != "" {
		awsASG.AutoScalingGroupName = &domainASG.Name
	} else if domainASG.NamePrefix != nil {
		awsASG.AutoScalingGroupNamePrefix = domainASG.NamePrefix
	}

	// Convert Launch Template
	if domainASG.LaunchTemplate != nil {
		awsASG.LaunchTemplate = &awsautoscaling.LaunchTemplateSpecification{
			LaunchTemplateId: domainASG.LaunchTemplate.ID,
			Version:          domainASG.LaunchTemplate.Version,
		}
	}

	// Convert Health Check Type
	if domainASG.HealthCheckType != "" {
		healthCheckType := string(domainASG.HealthCheckType)
		awsASG.HealthCheckType = &healthCheckType
	} else {
		// Default to EC2
		defaultHealthCheck := "EC2"
		awsASG.HealthCheckType = &defaultHealthCheck
	}

	// Convert Health Check Grace Period
	awsASG.HealthCheckGracePeriod = domainASG.HealthCheckGracePeriod

	// Convert Target Group ARNs
	awsASG.TargetGroupARNs = domainASG.TargetGroupARNs

	// Convert Tags
	if len(domainASG.Tags) > 0 {
		awsASG.Tags = make([]awsautoscaling.Tag, len(domainASG.Tags))
		for i, tag := range domainASG.Tags {
			awsASG.Tags[i] = awsautoscaling.Tag{
				Key:               tag.Key,
				Value:             tag.Value,
				PropagateAtLaunch: tag.PropagateAtLaunch,
			}
		}
	}

	return awsASG
}

// ToDomainAutoScalingGroupFromOutput converts AWS AutoScalingGroup output to domain AutoScalingGroup
func ToDomainAutoScalingGroupFromOutput(output *awsoutputs.AutoScalingGroupOutput) *domaincompute.AutoScalingGroup {
	if output == nil {
		return nil
	}

	arn := &output.AutoScalingGroupARN
	if output.AutoScalingGroupARN == "" {
		arn = nil
	}

	// Extract region from ARN if available
	// ARN format: arn:aws:autoscaling:REGION:ACCOUNT:autoScalingGroup:UUID:autoScalingGroupName/ASG_NAME
	region := ""
	if output.AutoScalingGroupARN != "" {
		parts := strings.Split(output.AutoScalingGroupARN, ":")
		if len(parts) >= 4 {
			region = parts[3]
		}
	}

	// Convert Launch Template
	var launchTemplate *domaincompute.LaunchTemplateSpecification
	if output.LaunchTemplate != nil {
		launchTemplate = &domaincompute.LaunchTemplateSpecification{
			ID:      output.LaunchTemplate.LaunchTemplateId,
			Version: output.LaunchTemplate.Version,
		}
	}

	// Convert Health Check Type
	healthCheckType := domaincompute.AutoScalingGroupHealthCheckTypeEC2
	if output.HealthCheckType != "" {
		healthCheckType = domaincompute.AutoScalingGroupHealthCheckType(strings.ToUpper(output.HealthCheckType))
	}

	// Convert State
	state := domaincompute.AutoScalingGroupStateActive
	switch strings.ToLower(output.Status) {
	case "deleting":
		state = domaincompute.AutoScalingGroupStateDeleting
	case "updating":
		state = domaincompute.AutoScalingGroupStateUpdating
	}

	// Convert Tags
	tags := make([]domaincompute.Tag, len(output.Tags))
	for i, tag := range output.Tags {
		tags[i] = domaincompute.Tag{
			Key:               tag.Key,
			Value:             tag.Value,
			PropagateAtLaunch: tag.PropagateAtLaunch,
		}
	}

	// Convert CreatedTime to string
	var createdTime *string
	if !output.CreatedTime.IsZero() {
		timeStr := output.CreatedTime.Format("2006-01-02T15:04:05Z07:00")
		createdTime = &timeStr
	}

	domainASG := &domaincompute.AutoScalingGroup{
		ID:                  output.AutoScalingGroupName,
		ARN:                 arn,
		Name:                output.AutoScalingGroupName,
		Region:              region,
		MinSize:             output.MinSize,
		MaxSize:             output.MaxSize,
		DesiredCapacity:     &output.DesiredCapacity,
		VPCZoneIdentifier:   output.VPCZoneIdentifier,
		LaunchTemplate:      launchTemplate,
		HealthCheckType:     healthCheckType,
		HealthCheckGracePeriod: output.HealthCheckGracePeriod,
		TargetGroupARNs:     output.TargetGroupARNs,
		Tags:                tags,
		State:               state,
		CreatedTime:         createdTime,
	}

	return domainASG
}

// ToDomainAutoScalingGroupOutputFromOutput converts AWS AutoScalingGroup output directly to domain AutoScalingGroupOutput
func ToDomainAutoScalingGroupOutputFromOutput(output *awsoutputs.AutoScalingGroupOutput) *domaincompute.AutoScalingGroupOutput {
	if output == nil {
		return nil
	}

	arn := &output.AutoScalingGroupARN
	if output.AutoScalingGroupARN == "" {
		arn = nil
	}

	// Extract region from ARN if available
	region := ""
	if output.AutoScalingGroupARN != "" {
		parts := strings.Split(output.AutoScalingGroupARN, ":")
		if len(parts) >= 4 {
			region = parts[3]
		}
	}

	// Convert Launch Template
	var launchTemplate *domaincompute.LaunchTemplateSpecification
	if output.LaunchTemplate != nil {
		launchTemplate = &domaincompute.LaunchTemplateSpecification{
			ID:      output.LaunchTemplate.LaunchTemplateId,
			Version: output.LaunchTemplate.Version,
		}
	}

	// Convert Health Check Type
	healthCheckType := domaincompute.AutoScalingGroupHealthCheckTypeEC2
	if output.HealthCheckType != "" {
		healthCheckType = domaincompute.AutoScalingGroupHealthCheckType(strings.ToUpper(output.HealthCheckType))
	}

	// Convert State
	state := domaincompute.AutoScalingGroupStateActive
	switch strings.ToLower(output.Status) {
	case "deleting":
		state = domaincompute.AutoScalingGroupStateDeleting
	case "updating":
		state = domaincompute.AutoScalingGroupStateUpdating
	}

	// Convert CreatedTime to string
	var createdTime *string
	if !output.CreatedTime.IsZero() {
		timeStr := output.CreatedTime.Format("2006-01-02T15:04:05Z07:00")
		createdTime = &timeStr
	}

	// AutoScalingGroupOutput doesn't have NamePrefix field
	var namePrefix *string

	return &domaincompute.AutoScalingGroupOutput{
		ID:                    output.AutoScalingGroupName,
		ARN:                   arn,
		Name:                  output.AutoScalingGroupName,
		Region:                region,
		NamePrefix:            namePrefix,
		MinSize:               output.MinSize,
		MaxSize:               output.MaxSize,
		DesiredCapacity:       &output.DesiredCapacity,
		VPCZoneIdentifier:     output.VPCZoneIdentifier,
		LaunchTemplate:        launchTemplate,
		HealthCheckType:       healthCheckType,
		HealthCheckGracePeriod: output.HealthCheckGracePeriod,
		TargetGroupARNs:       output.TargetGroupARNs,
		State:                 state,
		CreatedTime:           createdTime,
	}
}
