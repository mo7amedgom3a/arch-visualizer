package compute

import (
	"context"
)

// ComputeService defines the interface for compute resource operations
// This is cloud-agnostic and can be implemented by any cloud provider
type ComputeService interface {
	// Instance operations
	CreateInstance(ctx context.Context, instance *Instance) (*Instance, error)
	GetInstance(ctx context.Context, id string) (*Instance, error)
	UpdateInstance(ctx context.Context, instance *Instance) (*Instance, error)
	DeleteInstance(ctx context.Context, id string) error
	ListInstances(ctx context.Context, filters map[string]string) ([]*Instance, error)

	// Instance lifecycle operations
	StartInstance(ctx context.Context, id string) error
	StopInstance(ctx context.Context, id string) error
	RebootInstance(ctx context.Context, id string) error

	// Launch Template operations
	CreateLaunchTemplate(ctx context.Context, template *LaunchTemplate) (*LaunchTemplate, error)
	GetLaunchTemplate(ctx context.Context, id string) (*LaunchTemplate, error)
	UpdateLaunchTemplate(ctx context.Context, id string, template *LaunchTemplate) (*LaunchTemplate, error)
	DeleteLaunchTemplate(ctx context.Context, id string) error
	ListLaunchTemplates(ctx context.Context, filters map[string]string) ([]*LaunchTemplate, error)
	GetLaunchTemplateVersion(ctx context.Context, id string, version int) (*LaunchTemplate, error)
	ListLaunchTemplateVersions(ctx context.Context, id string) ([]*LaunchTemplateVersion, error)

	// Load Balancer operations
	CreateLoadBalancer(ctx context.Context, lb *LoadBalancer) (*LoadBalancer, error)
	GetLoadBalancer(ctx context.Context, arn string) (*LoadBalancer, error)
	UpdateLoadBalancer(ctx context.Context, lb *LoadBalancer) (*LoadBalancer, error)
	DeleteLoadBalancer(ctx context.Context, arn string) error
	ListLoadBalancers(ctx context.Context, filters map[string]string) ([]*LoadBalancer, error)

	// Target Group operations
	CreateTargetGroup(ctx context.Context, tg *TargetGroup) (*TargetGroup, error)
	GetTargetGroup(ctx context.Context, arn string) (*TargetGroup, error)
	UpdateTargetGroup(ctx context.Context, tg *TargetGroup) (*TargetGroup, error)
	DeleteTargetGroup(ctx context.Context, arn string) error
	ListTargetGroups(ctx context.Context, filters map[string]string) ([]*TargetGroup, error)

	// Listener operations
	CreateListener(ctx context.Context, listener *Listener) (*Listener, error)
	GetListener(ctx context.Context, arn string) (*Listener, error)
	UpdateListener(ctx context.Context, listener *Listener) (*Listener, error)
	DeleteListener(ctx context.Context, arn string) error
	ListListeners(ctx context.Context, loadBalancerARN string) ([]*Listener, error)

	// Target Group Attachment operations
	AttachTargetToGroup(ctx context.Context, attachment *TargetGroupAttachment) error
	DetachTargetFromGroup(ctx context.Context, targetGroupARN, targetID string) error
	ListTargetGroupTargets(ctx context.Context, targetGroupARN string) ([]*TargetGroupAttachment, error)

	// Auto Scaling Group operations
	CreateAutoScalingGroup(ctx context.Context, asg *AutoScalingGroup) (*AutoScalingGroup, error)
	GetAutoScalingGroup(ctx context.Context, name string) (*AutoScalingGroup, error)
	UpdateAutoScalingGroup(ctx context.Context, asg *AutoScalingGroup) (*AutoScalingGroup, error)
	DeleteAutoScalingGroup(ctx context.Context, name string) error
	ListAutoScalingGroups(ctx context.Context, filters map[string]string) ([]*AutoScalingGroup, error)

	// Scaling operations
	SetDesiredCapacity(ctx context.Context, asgName string, capacity int) error
	AttachInstances(ctx context.Context, asgName string, instanceIDs []string) error
	DetachInstances(ctx context.Context, asgName string, instanceIDs []string) error

	// Lambda Function operations
	CreateLambdaFunction(ctx context.Context, function *LambdaFunction) (*LambdaFunction, error)
	GetLambdaFunction(ctx context.Context, name string) (*LambdaFunction, error)
	UpdateLambdaFunction(ctx context.Context, function *LambdaFunction) (*LambdaFunction, error)
	DeleteLambdaFunction(ctx context.Context, name string) error
	ListLambdaFunctions(ctx context.Context, filters map[string]string) ([]*LambdaFunction, error)
}

// ComputeRepository defines the interface for compute resource persistence
// This abstracts data access and can be implemented for different storage backends
type ComputeRepository interface {
	// Instance persistence
	SaveInstance(ctx context.Context, instance *Instance) error
	FindInstanceByID(ctx context.Context, id string) (*Instance, error)
	FindInstancesByFilters(ctx context.Context, filters map[string]string) ([]*Instance, error)
	DeleteInstance(ctx context.Context, id string) error
}
