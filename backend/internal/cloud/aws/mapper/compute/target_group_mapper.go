package compute

import (
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// FromDomainTargetGroup converts domain TargetGroup to AWS TargetGroup
func FromDomainTargetGroup(domainTG *domaincompute.TargetGroup) *awsloadbalancer.TargetGroup {
	if domainTG == nil {
		return nil
	}

	targetType := "instance"
	if domainTG.TargetType != "" {
		switch domainTG.TargetType {
		case domaincompute.TargetTypeIP:
			targetType = "ip"
		case domaincompute.TargetTypeLambda:
			targetType = "lambda"
		default:
			targetType = "instance"
		}
	}

	protocol := string(domainTG.Protocol)

	awsTG := &awsloadbalancer.TargetGroup{
		Name:       domainTG.Name,
		Port:       domainTG.Port,
		Protocol:   protocol,
		VPCID:      domainTG.VPCID,
		TargetType: &targetType,
		HealthCheck: awsloadbalancer.HealthCheckConfig{
			Path:               domainTG.HealthCheck.Path,
			Matcher:            domainTG.HealthCheck.Matcher,
			Interval:           domainTG.HealthCheck.Interval,
			Timeout:            domainTG.HealthCheck.Timeout,
			HealthyThreshold:   domainTG.HealthCheck.HealthyThreshold,
			UnhealthyThreshold: domainTG.HealthCheck.UnhealthyThreshold,
			Protocol:           domainTG.HealthCheck.Protocol,
			Port:               domainTG.HealthCheck.Port,
		},
		Tags: []configs.Tag{},
	}

	return awsTG
}

// ToDomainTargetGroupFromOutput converts AWS TargetGroup output to domain TargetGroup
func ToDomainTargetGroupFromOutput(output *awsoutputs.TargetGroupOutput) *domaincompute.TargetGroup {
	if output == nil {
		return nil
	}

	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}

	protocol := domaincompute.TargetGroupProtocol(output.Protocol)

	targetType := domaincompute.TargetTypeInstance
	switch output.TargetType {
	case "ip":
		targetType = domaincompute.TargetTypeIP
	case "lambda":
		targetType = domaincompute.TargetTypeLambda
	default:
		targetType = domaincompute.TargetTypeInstance
	}

	state := domaincompute.TargetGroupStateActive
	switch output.State {
	case "draining":
		state = domaincompute.TargetGroupStateDraining
	case "deleting":
		state = domaincompute.TargetGroupStateDeleting
	case "deleted":
		state = domaincompute.TargetGroupStateDeleted
	}

	domainTG := &domaincompute.TargetGroup{
		ID:   output.ID,
		ARN:  arn,
		Name: output.Name,
		VPCID: output.VPCID,
		Port: output.Port,
		Protocol: protocol,
		TargetType: targetType,
		HealthCheck: domaincompute.HealthCheckConfig{
			Path:               output.HealthCheck.Path,
			Matcher:            output.HealthCheck.Matcher,
			Interval:           output.HealthCheck.Interval,
			Timeout:            output.HealthCheck.Timeout,
			HealthyThreshold:   output.HealthCheck.HealthyThreshold,
			UnhealthyThreshold: output.HealthCheck.UnhealthyThreshold,
			Protocol:           output.HealthCheck.Protocol,
			Port:               output.HealthCheck.Port,
		},
		State: state,
	}

	return domainTG
}

// ToDomainTargetGroupOutputFromOutput converts AWS TargetGroup output directly to domain TargetGroupOutput
func ToDomainTargetGroupOutputFromOutput(output *awsoutputs.TargetGroupOutput) *domaincompute.TargetGroupOutput {
	if output == nil {
		return nil
	}

	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}

	protocol := domaincompute.TargetGroupProtocol(output.Protocol)

	targetType := domaincompute.TargetTypeInstance
	switch output.TargetType {
	case "ip":
		targetType = domaincompute.TargetTypeIP
	case "lambda":
		targetType = domaincompute.TargetTypeLambda
	default:
		targetType = domaincompute.TargetTypeInstance
	}

	state := domaincompute.TargetGroupStateActive
	switch output.State {
	case "draining":
		state = domaincompute.TargetGroupStateDraining
	case "deleting":
		state = domaincompute.TargetGroupStateDeleting
	case "deleted":
		state = domaincompute.TargetGroupStateDeleted
	}

	createdAt := &output.CreatedTime

	return &domaincompute.TargetGroupOutput{
		ID:   output.ID,
		ARN:  arn,
		Name: output.Name,
		VPCID: output.VPCID,
		Port: output.Port,
		Protocol: protocol,
		TargetType: targetType,
		HealthCheck: domaincompute.HealthCheckConfig{
			Path:               output.HealthCheck.Path,
			Matcher:            output.HealthCheck.Matcher,
			Interval:           output.HealthCheck.Interval,
			Timeout:            output.HealthCheck.Timeout,
			HealthyThreshold:   output.HealthCheck.HealthyThreshold,
			UnhealthyThreshold: output.HealthCheck.UnhealthyThreshold,
			Protocol:           output.HealthCheck.Protocol,
			Port:               output.HealthCheck.Port,
		},
		State: state,
		CreatedAt: createdAt,
	}
}
