package compute

import (
	"time"

	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer/outputs"
)

// FromDomainListener converts domain Listener to AWS Listener
func FromDomainListener(domainListener *domaincompute.Listener) *awsloadbalancer.Listener {
	if domainListener == nil {
		return nil
	}

	protocol := string(domainListener.Protocol)

	actionType := string(domainListener.DefaultAction.Type)
	awsAction := awsloadbalancer.ListenerAction{
		Type: awsloadbalancer.ListenerActionType(actionType),
	}

	if domainListener.DefaultAction.TargetGroupARN != nil {
		awsAction.TargetGroupARN = domainListener.DefaultAction.TargetGroupARN
	}

	if domainListener.DefaultAction.RedirectConfig != nil {
		awsAction.RedirectConfig = &awsloadbalancer.RedirectConfig{
			Protocol:   domainListener.DefaultAction.RedirectConfig.Protocol,
			Port:       domainListener.DefaultAction.RedirectConfig.Port,
			StatusCode: domainListener.DefaultAction.RedirectConfig.StatusCode,
			Host:       domainListener.DefaultAction.RedirectConfig.Host,
			Path:       domainListener.DefaultAction.RedirectConfig.Path,
			Query:      domainListener.DefaultAction.RedirectConfig.Query,
		}
	}

	if domainListener.DefaultAction.FixedResponseConfig != nil {
		awsAction.FixedResponseConfig = &awsloadbalancer.FixedResponseConfig{
			ContentType: domainListener.DefaultAction.FixedResponseConfig.ContentType,
			MessageBody: domainListener.DefaultAction.FixedResponseConfig.MessageBody,
			StatusCode:  domainListener.DefaultAction.FixedResponseConfig.StatusCode,
		}
	}

	awsListener := &awsloadbalancer.Listener{
		LoadBalancerARN: domainListener.LoadBalancerARN,
		Port:            domainListener.Port,
		Protocol:        protocol,
		DefaultAction:   awsAction,
	}

	return awsListener
}

// ToDomainListenerFromOutput converts AWS Listener output to domain Listener
func ToDomainListenerFromOutput(output *awsoutputs.ListenerOutput) *domaincompute.Listener {
	if output == nil {
		return nil
	}

	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}

	protocol := domaincompute.ListenerProtocol(output.Protocol)

	actionType := domaincompute.ListenerActionType(output.DefaultAction.Type)
	domainAction := domaincompute.ListenerAction{
		Type: actionType,
	}

	if output.DefaultAction.TargetGroupARN != nil {
		domainAction.TargetGroupARN = output.DefaultAction.TargetGroupARN
	}

	if output.DefaultAction.RedirectConfig != nil {
		domainAction.RedirectConfig = &domaincompute.RedirectConfig{
			Protocol:   output.DefaultAction.RedirectConfig.Protocol,
			Port:       output.DefaultAction.RedirectConfig.Port,
			StatusCode: output.DefaultAction.RedirectConfig.StatusCode,
			Host:       output.DefaultAction.RedirectConfig.Host,
			Path:       output.DefaultAction.RedirectConfig.Path,
			Query:      output.DefaultAction.RedirectConfig.Query,
		}
	}

	if output.DefaultAction.FixedResponseConfig != nil {
		domainAction.FixedResponseConfig = &domaincompute.FixedResponseConfig{
			ContentType: output.DefaultAction.FixedResponseConfig.ContentType,
			MessageBody: output.DefaultAction.FixedResponseConfig.MessageBody,
			StatusCode:  output.DefaultAction.FixedResponseConfig.StatusCode,
		}
	}

	domainListener := &domaincompute.Listener{
		ID:              output.ID,
		ARN:             arn,
		LoadBalancerARN: output.LoadBalancerARN,
		Port:            output.Port,
		Protocol:        protocol,
		DefaultAction:   domainAction,
		Rules:           []domaincompute.ListenerRule{}, // Rules not implemented yet
	}

	return domainListener
}

// ToDomainListenerOutputFromOutput converts AWS Listener output directly to domain ListenerOutput
func ToDomainListenerOutputFromOutput(output *awsoutputs.ListenerOutput) *domaincompute.ListenerOutput {
	if output == nil {
		return nil
	}

	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}

	protocol := domaincompute.ListenerProtocol(output.Protocol)

	actionType := domaincompute.ListenerActionType(output.DefaultAction.Type)
	domainAction := domaincompute.ListenerAction{
		Type: actionType,
	}

	if output.DefaultAction.TargetGroupARN != nil {
		domainAction.TargetGroupARN = output.DefaultAction.TargetGroupARN
	}

	if output.DefaultAction.RedirectConfig != nil {
		domainAction.RedirectConfig = &domaincompute.RedirectConfig{
			Protocol:   output.DefaultAction.RedirectConfig.Protocol,
			Port:       output.DefaultAction.RedirectConfig.Port,
			StatusCode: output.DefaultAction.RedirectConfig.StatusCode,
			Host:       output.DefaultAction.RedirectConfig.Host,
			Path:       output.DefaultAction.RedirectConfig.Path,
			Query:      output.DefaultAction.RedirectConfig.Query,
		}
	}

	if output.DefaultAction.FixedResponseConfig != nil {
		domainAction.FixedResponseConfig = &domaincompute.FixedResponseConfig{
			ContentType: output.DefaultAction.FixedResponseConfig.ContentType,
			MessageBody: output.DefaultAction.FixedResponseConfig.MessageBody,
			StatusCode:  output.DefaultAction.FixedResponseConfig.StatusCode,
		}
	}

	// ListenerOutput doesn't have CreationTime, so we set it to nil
	var createdAt *time.Time

	return &domaincompute.ListenerOutput{
		ID:              output.ID,
		ARN:             arn,
		LoadBalancerARN: output.LoadBalancerARN,
		Port:            output.Port,
		Protocol:        protocol,
		DefaultAction:   domainAction,
		Rules:           []domaincompute.ListenerRule{}, // Rules not implemented yet
		CreatedAt:       createdAt,
	}
}
