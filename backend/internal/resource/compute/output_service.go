package compute

import (
	"context"
)

// ComputeOutputService defines the interface for compute resource operations that return output DTOs
// This is a parallel interface to ComputeService, providing output-specific models
type ComputeOutputService interface {
	// Instance operations
	CreateInstanceOutput(ctx context.Context, instance *Instance) (*InstanceOutput, error)
	GetInstanceOutput(ctx context.Context, id string) (*InstanceOutput, error)
	UpdateInstanceOutput(ctx context.Context, instance *Instance) (*InstanceOutput, error)
	ListInstancesOutput(ctx context.Context, filters map[string]string) ([]*InstanceOutput, error)

	// Launch Template operations
	CreateLaunchTemplateOutput(ctx context.Context, template *LaunchTemplate) (*LaunchTemplateOutput, error)
	GetLaunchTemplateOutput(ctx context.Context, id string) (*LaunchTemplateOutput, error)
	UpdateLaunchTemplateOutput(ctx context.Context, id string, template *LaunchTemplate) (*LaunchTemplateOutput, error)
	ListLaunchTemplatesOutput(ctx context.Context, filters map[string]string) ([]*LaunchTemplateOutput, error)

	// Load Balancer operations
	CreateLoadBalancerOutput(ctx context.Context, lb *LoadBalancer) (*LoadBalancerOutput, error)
	GetLoadBalancerOutput(ctx context.Context, arn string) (*LoadBalancerOutput, error)
	UpdateLoadBalancerOutput(ctx context.Context, lb *LoadBalancer) (*LoadBalancerOutput, error)
	ListLoadBalancersOutput(ctx context.Context, filters map[string]string) ([]*LoadBalancerOutput, error)

	// Target Group operations
	CreateTargetGroupOutput(ctx context.Context, tg *TargetGroup) (*TargetGroupOutput, error)
	GetTargetGroupOutput(ctx context.Context, arn string) (*TargetGroupOutput, error)
	UpdateTargetGroupOutput(ctx context.Context, tg *TargetGroup) (*TargetGroupOutput, error)
	ListTargetGroupsOutput(ctx context.Context, filters map[string]string) ([]*TargetGroupOutput, error)

	// Listener operations
	CreateListenerOutput(ctx context.Context, listener *Listener) (*ListenerOutput, error)
	GetListenerOutput(ctx context.Context, arn string) (*ListenerOutput, error)
	UpdateListenerOutput(ctx context.Context, listener *Listener) (*ListenerOutput, error)
	ListListenersOutput(ctx context.Context, loadBalancerARN string) ([]*ListenerOutput, error)

	// Target Group Attachment operations
	ListTargetGroupTargetsOutput(ctx context.Context, targetGroupARN string) ([]*TargetGroupAttachmentOutput, error)

	// Auto Scaling Group operations
	CreateAutoScalingGroupOutput(ctx context.Context, asg *AutoScalingGroup) (*AutoScalingGroupOutput, error)
	GetAutoScalingGroupOutput(ctx context.Context, name string) (*AutoScalingGroupOutput, error)
	UpdateAutoScalingGroupOutput(ctx context.Context, asg *AutoScalingGroup) (*AutoScalingGroupOutput, error)
	ListAutoScalingGroupsOutput(ctx context.Context, filters map[string]string) ([]*AutoScalingGroupOutput, error)

	// Lambda Function operations
	CreateLambdaFunctionOutput(ctx context.Context, function *LambdaFunction) (*LambdaFunctionOutput, error)
	GetLambdaFunctionOutput(ctx context.Context, name string) (*LambdaFunctionOutput, error)
	UpdateLambdaFunctionOutput(ctx context.Context, function *LambdaFunction) (*LambdaFunctionOutput, error)
	ListLambdaFunctionsOutput(ctx context.Context, filters map[string]string) ([]*LambdaFunctionOutput, error)
}
