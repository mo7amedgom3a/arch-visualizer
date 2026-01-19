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

// CreateListener creates a new Listener using AWS SDK
func CreateListener(ctx context.Context, client *AWSClient, listener *awsloadbalancer.Listener) (*awsoutputs.ListenerOutput, error) {
	if err := listener.Validate(); err != nil {
		return nil, fmt.Errorf("listener validation failed: %w", err)
	}

	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	// Build CreateListenerInput
	input := &elasticloadbalancingv2.CreateListenerInput{
		LoadBalancerArn: aws.String(listener.LoadBalancerARN),
		Port:            aws.Int32(int32(listener.Port)),
		Protocol:        elbv2types.ProtocolEnum(listener.Protocol),
		DefaultActions:  convertListenerActions([]awsloadbalancer.ListenerAction{listener.DefaultAction}),
	}

	// Set certificate ARN for HTTPS/TLS
	if listener.CertificateARN != nil && *listener.CertificateARN != "" {
		input.Certificates = []elbv2types.Certificate{
			{CertificateArn: listener.CertificateARN},
		}
	}

	// Set SSL policy if provided
	if listener.SSLPolicy != nil {
		input.SslPolicy = listener.SSLPolicy
	}

	// Create the listener
	result, err := client.ELBv2.CreateListener(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	if len(result.Listeners) == 0 {
		return nil, fmt.Errorf("listener creation returned no listeners")
	}

	return convertListenerToOutput(&result.Listeners[0]), nil
}

// GetListener retrieves a Listener by ARN
func GetListener(ctx context.Context, client *AWSClient, arn string) (*awsoutputs.ListenerOutput, error) {
	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	input := &elasticloadbalancingv2.DescribeListenersInput{
		ListenerArns: []string{arn},
	}

	result, err := client.ELBv2.DescribeListeners(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get listener %s: %w", arn, err)
	}

	if len(result.Listeners) == 0 {
		return nil, fmt.Errorf("listener %s not found", arn)
	}

	return convertListenerToOutput(&result.Listeners[0]), nil
}

// UpdateListener updates a Listener
func UpdateListener(ctx context.Context, client *AWSClient, arn string, listener *awsloadbalancer.Listener) (*awsoutputs.ListenerOutput, error) {
	if err := listener.Validate(); err != nil {
		return nil, fmt.Errorf("listener validation failed: %w", err)
	}

	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	// Build ModifyListenerInput
	input := &elasticloadbalancingv2.ModifyListenerInput{
		ListenerArn:    aws.String(arn),
		Port:           aws.Int32(int32(listener.Port)),
		Protocol:       elbv2types.ProtocolEnum(listener.Protocol),
		DefaultActions: convertListenerActions([]awsloadbalancer.ListenerAction{listener.DefaultAction}),
	}

	// Set certificate ARN for HTTPS/TLS
	if listener.CertificateARN != nil && *listener.CertificateARN != "" {
		input.Certificates = []elbv2types.Certificate{
			{CertificateArn: listener.CertificateARN},
		}
	}

	// Set SSL policy if provided
	if listener.SSLPolicy != nil {
		input.SslPolicy = listener.SSLPolicy
	}

	// Update the listener
	result, err := client.ELBv2.ModifyListener(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update listener: %w", err)
	}

	if len(result.Listeners) == 0 {
		return nil, fmt.Errorf("listener update returned no listeners")
	}

	return convertListenerToOutput(&result.Listeners[0]), nil
}

// DeleteListener deletes a Listener
func DeleteListener(ctx context.Context, client *AWSClient, arn string) error {
	if client == nil || client.ELBv2 == nil {
		return fmt.Errorf("AWS client not available")
	}

	input := &elasticloadbalancingv2.DeleteListenerInput{
		ListenerArn: aws.String(arn),
	}

	_, err := client.ELBv2.DeleteListener(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete listener %s: %w", arn, err)
	}

	return nil
}

// ListListeners lists Listeners for a Load Balancer
func ListListeners(ctx context.Context, client *AWSClient, loadBalancerARN string) ([]*awsoutputs.ListenerOutput, error) {
	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	var allListeners []*awsoutputs.ListenerOutput
	var nextToken *string

	for {
		input := &elasticloadbalancingv2.DescribeListenersInput{
			LoadBalancerArn: aws.String(loadBalancerARN),
		}

		if nextToken != nil {
			input.Marker = nextToken
		}

		result, err := client.ELBv2.DescribeListeners(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list listeners: %w", err)
		}

		// Convert each listener to output model
		for _, listener := range result.Listeners {
			allListeners = append(allListeners, convertListenerToOutput(&listener))
		}

		// Check if there are more pages
		if result.NextMarker == nil {
			break
		}
		nextToken = result.NextMarker
	}

	return allListeners, nil
}

// convertListenerActions converts domain listener actions to AWS SDK actions
func convertListenerActions(actions []awsloadbalancer.ListenerAction) []elbv2types.Action {
	awsActions := make([]elbv2types.Action, len(actions))
	for i, action := range actions {
		awsAction := elbv2types.Action{
			Type: elbv2types.ActionTypeEnum(action.Type),
		}

		if action.TargetGroupARN != nil {
			awsAction.TargetGroupArn = action.TargetGroupARN
		}

		if action.RedirectConfig != nil {
			awsAction.RedirectConfig = &elbv2types.RedirectActionConfig{
				Protocol:   action.RedirectConfig.Protocol,
				Port:       action.RedirectConfig.Port,
				StatusCode: elbv2types.RedirectActionStatusCodeEnum(action.RedirectConfig.StatusCode),
				Host:       action.RedirectConfig.Host,
				Path:       action.RedirectConfig.Path,
				Query:      action.RedirectConfig.Query,
			}
		}

		if action.FixedResponseConfig != nil {
			awsAction.FixedResponseConfig = &elbv2types.FixedResponseActionConfig{
				ContentType: aws.String(action.FixedResponseConfig.ContentType),
				MessageBody: action.FixedResponseConfig.MessageBody,
				StatusCode:  aws.String(action.FixedResponseConfig.StatusCode),
			}
		}

		awsActions[i] = awsAction
	}
	return awsActions
}

// convertListenerToOutput converts AWS SDK Listener to output model
func convertListenerToOutput(listener *elbv2types.Listener) *awsoutputs.ListenerOutput {
	output := &awsoutputs.ListenerOutput{
		ARN:             aws.ToString(listener.ListenerArn),
		ID:              aws.ToString(listener.ListenerArn),
		LoadBalancerARN: aws.ToString(listener.LoadBalancerArn),
		Port:            int(aws.ToInt32(listener.Port)),
		Protocol:        string(listener.Protocol),
	}

	// Convert default action (take first action)
	if len(listener.DefaultActions) > 0 {
		output.DefaultAction = convertSDKActionToDomain(listener.DefaultActions[0])
	}

	return output
}

// convertSDKActionToDomain converts AWS SDK Action to domain action
func convertSDKActionToDomain(action elbv2types.Action) awsloadbalancer.ListenerAction {
	domainAction := awsloadbalancer.ListenerAction{
		Type: awsloadbalancer.ListenerActionType(action.Type),
	}

	if action.TargetGroupArn != nil {
		domainAction.TargetGroupARN = action.TargetGroupArn
	}

	if action.RedirectConfig != nil {
		domainAction.RedirectConfig = &awsloadbalancer.RedirectConfig{
			Protocol:   action.RedirectConfig.Protocol,
			Port:       action.RedirectConfig.Port,
			StatusCode: string(action.RedirectConfig.StatusCode),
			Host:       action.RedirectConfig.Host,
			Path:       action.RedirectConfig.Path,
			Query:      action.RedirectConfig.Query,
		}
	}

	if action.FixedResponseConfig != nil {
		domainAction.FixedResponseConfig = &awsloadbalancer.FixedResponseConfig{
			ContentType: aws.ToString(action.FixedResponseConfig.ContentType),
			MessageBody: action.FixedResponseConfig.MessageBody,
			StatusCode:  aws.ToString(action.FixedResponseConfig.StatusCode),
		}
	}

	return domainAction
}
