package compute

import (
	"context"
	"fmt"

	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/compute"
	awsInstanceTypes "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/instance_types"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
)

// AWSComputeAdapter adapts AWS-specific compute service to domain compute service
// This implements the Adapter pattern, allowing the domain layer to work with cloud-specific implementations
type AWSComputeAdapter struct {
	awsService awsservice.AWSComputeService
}

// NewAWSComputeAdapter creates a new AWS compute adapter
func NewAWSComputeAdapter(awsService awsservice.AWSComputeService) domaincompute.ComputeService {
	return &AWSComputeAdapter{
		awsService: awsService,
	}
}

// Ensure AWSComputeAdapter implements ComputeService
var _ domaincompute.ComputeService = (*AWSComputeAdapter)(nil)

// Instance Operations

func (a *AWSComputeAdapter) CreateInstance(ctx context.Context, instance *domaincompute.Instance) (*domaincompute.Instance, error) {
	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsInstance := awsmapper.FromDomainInstance(instance)
	if err := awsInstance.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsInstanceOutput, err := a.awsService.CreateInstance(ctx, awsInstance)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainInstanceFromOutput(awsInstanceOutput), nil
}

func (a *AWSComputeAdapter) GetInstance(ctx context.Context, id string) (*domaincompute.Instance, error) {
	awsInstanceOutput, err := a.awsService.GetInstance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainInstanceFromOutput(awsInstanceOutput), nil
}

func (a *AWSComputeAdapter) UpdateInstance(ctx context.Context, instance *domaincompute.Instance) (*domaincompute.Instance, error) {
	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsInstance := awsmapper.FromDomainInstance(instance)
	if err := awsInstance.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsInstanceOutput, err := a.awsService.UpdateInstance(ctx, awsInstance)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainInstanceFromOutput(awsInstanceOutput), nil
}

func (a *AWSComputeAdapter) DeleteInstance(ctx context.Context, id string) error {
	if err := a.awsService.DeleteInstance(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) ListInstances(ctx context.Context, filters map[string]string) ([]*domaincompute.Instance, error) {
	// Convert domain filters to AWS filters format
	awsFilters := make(map[string][]string)
	for key, value := range filters {
		awsFilters[key] = []string{value}
	}

	awsInstanceOutputs, err := a.awsService.ListInstances(ctx, awsFilters)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainInstances := make([]*domaincompute.Instance, len(awsInstanceOutputs))
	for i, awsInstanceOutput := range awsInstanceOutputs {
		domainInstances[i] = awsmapper.ToDomainInstanceFromOutput(awsInstanceOutput)
	}

	return domainInstances, nil
}

// Instance Lifecycle Operations

func (a *AWSComputeAdapter) StartInstance(ctx context.Context, id string) error {
	if err := a.awsService.StartInstance(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) StopInstance(ctx context.Context, id string) error {
	if err := a.awsService.StopInstance(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) RebootInstance(ctx context.Context, id string) error {
	if err := a.awsService.RebootInstance(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

// Instance Type Operations

func (a *AWSComputeAdapter) GetInstanceTypeInfo(ctx context.Context, instanceType string, region string) (*awsInstanceTypes.InstanceTypeInfo, error) {
	info, err := a.awsService.GetInstanceTypeInfo(ctx, instanceType, region)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}
	return info, nil
}

func (a *AWSComputeAdapter) ListInstanceTypesByCategory(ctx context.Context, category awsInstanceTypes.InstanceCategory, region string) ([]*awsInstanceTypes.InstanceTypeInfo, error) {
	types, err := a.awsService.ListInstanceTypesByCategory(ctx, category, region)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}
	return types, nil
}

// Launch Template Operations

func (a *AWSComputeAdapter) CreateLaunchTemplate(ctx context.Context, template *domaincompute.LaunchTemplate) (*domaincompute.LaunchTemplate, error) {
	if err := template.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsTemplate := awsmapper.FromDomainLaunchTemplate(template)
	if err := awsTemplate.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsTemplateOutput, err := a.awsService.CreateLaunchTemplate(ctx, awsTemplate)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainTemplate := awsmapper.ToDomainLaunchTemplateFromOutput(awsTemplateOutput)
	// Preserve region from input
	domainTemplate.Region = template.Region
	return domainTemplate, nil
}

func (a *AWSComputeAdapter) GetLaunchTemplate(ctx context.Context, id string) (*domaincompute.LaunchTemplate, error) {
	awsTemplateOutput, err := a.awsService.GetLaunchTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainLaunchTemplateFromOutput(awsTemplateOutput), nil
}

func (a *AWSComputeAdapter) UpdateLaunchTemplate(ctx context.Context, id string, template *domaincompute.LaunchTemplate) (*domaincompute.LaunchTemplate, error) {
	if err := template.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsTemplate := awsmapper.FromDomainLaunchTemplate(template)
	if err := awsTemplate.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsTemplateOutput, err := a.awsService.UpdateLaunchTemplate(ctx, id, awsTemplate)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainTemplate := awsmapper.ToDomainLaunchTemplateFromOutput(awsTemplateOutput)
	// Preserve region from input
	domainTemplate.Region = template.Region
	return domainTemplate, nil
}

func (a *AWSComputeAdapter) DeleteLaunchTemplate(ctx context.Context, id string) error {
	if err := a.awsService.DeleteLaunchTemplate(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) ListLaunchTemplates(ctx context.Context, filters map[string]string) ([]*domaincompute.LaunchTemplate, error) {
	// Convert domain filters to AWS filters format
	awsFilters := make(map[string][]string)
	for key, value := range filters {
		awsFilters[key] = []string{value}
	}

	awsTemplateOutputs, err := a.awsService.ListLaunchTemplates(ctx, awsFilters)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainTemplates := make([]*domaincompute.LaunchTemplate, len(awsTemplateOutputs))
	for i, awsTemplateOutput := range awsTemplateOutputs {
		domainTemplates[i] = awsmapper.ToDomainLaunchTemplateFromOutput(awsTemplateOutput)
	}

	return domainTemplates, nil
}

func (a *AWSComputeAdapter) GetLaunchTemplateVersion(ctx context.Context, id string, version int) (*domaincompute.LaunchTemplate, error) {
	awsVersion, err := a.awsService.GetLaunchTemplateVersion(ctx, id, version)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainTemplate := awsmapper.ToDomainLaunchTemplate(awsVersion.TemplateData)
	if domainTemplate != nil {
		domainTemplate.ID = awsVersion.TemplateID
		domainTemplate.Version = &awsVersion.VersionNumber
	}

	return domainTemplate, nil
}

func (a *AWSComputeAdapter) ListLaunchTemplateVersions(ctx context.Context, id string) ([]*domaincompute.LaunchTemplateVersion, error) {
	awsVersions, err := a.awsService.ListLaunchTemplateVersions(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainVersions := make([]*domaincompute.LaunchTemplateVersion, len(awsVersions))
	for i := range awsVersions {
		domainVersions[i] = awsmapper.ToDomainLaunchTemplateVersion(awsVersions[i])
	}

	return domainVersions, nil
}

// Load Balancer Operations

func (a *AWSComputeAdapter) CreateLoadBalancer(ctx context.Context, lb *domaincompute.LoadBalancer) (*domaincompute.LoadBalancer, error) {
	if err := lb.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsLB := awsmapper.FromDomainLoadBalancer(lb)
	if err := awsLB.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsLBOutput, err := a.awsService.CreateLoadBalancer(ctx, awsLB)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainLB := awsmapper.ToDomainLoadBalancerFromOutput(awsLBOutput)
	if domainLB != nil {
		domainLB.Region = lb.Region // Preserve region from input
	}
	return domainLB, nil
}

func (a *AWSComputeAdapter) GetLoadBalancer(ctx context.Context, arn string) (*domaincompute.LoadBalancer, error) {
	awsLBOutput, err := a.awsService.GetLoadBalancer(ctx, arn)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainLoadBalancerFromOutput(awsLBOutput), nil
}

func (a *AWSComputeAdapter) UpdateLoadBalancer(ctx context.Context, lb *domaincompute.LoadBalancer) (*domaincompute.LoadBalancer, error) {
	if err := lb.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}
	if lb.ARN == nil {
		return nil, fmt.Errorf("load balancer ARN is required for update")
	}

	awsLB := awsmapper.FromDomainLoadBalancer(lb)
	if err := awsLB.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsLBOutput, err := a.awsService.UpdateLoadBalancer(ctx, *lb.ARN, awsLB)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainLB := awsmapper.ToDomainLoadBalancerFromOutput(awsLBOutput)
	if domainLB != nil {
		domainLB.Region = lb.Region // Preserve region from input
	}
	return domainLB, nil
}

func (a *AWSComputeAdapter) DeleteLoadBalancer(ctx context.Context, arn string) error {
	if err := a.awsService.DeleteLoadBalancer(ctx, arn); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) ListLoadBalancers(ctx context.Context, filters map[string]string) ([]*domaincompute.LoadBalancer, error) {
	// Convert domain filters to AWS filters format
	awsFilters := make(map[string][]string)
	for key, value := range filters {
		awsFilters[key] = []string{value}
	}

	awsLBOutputs, err := a.awsService.ListLoadBalancers(ctx, awsFilters)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainLBs := make([]*domaincompute.LoadBalancer, len(awsLBOutputs))
	for i, awsLBOutput := range awsLBOutputs {
		domainLBs[i] = awsmapper.ToDomainLoadBalancerFromOutput(awsLBOutput)
	}

	return domainLBs, nil
}

// Target Group Operations

func (a *AWSComputeAdapter) CreateTargetGroup(ctx context.Context, tg *domaincompute.TargetGroup) (*domaincompute.TargetGroup, error) {
	if err := tg.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsTG := awsmapper.FromDomainTargetGroup(tg)
	if err := awsTG.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsTGOutput, err := a.awsService.CreateTargetGroup(ctx, awsTG)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainTargetGroupFromOutput(awsTGOutput), nil
}

func (a *AWSComputeAdapter) GetTargetGroup(ctx context.Context, arn string) (*domaincompute.TargetGroup, error) {
	awsTGOutput, err := a.awsService.GetTargetGroup(ctx, arn)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainTargetGroupFromOutput(awsTGOutput), nil
}

func (a *AWSComputeAdapter) UpdateTargetGroup(ctx context.Context, tg *domaincompute.TargetGroup) (*domaincompute.TargetGroup, error) {
	if err := tg.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}
	if tg.ARN == nil {
		return nil, fmt.Errorf("target group ARN is required for update")
	}

	awsTG := awsmapper.FromDomainTargetGroup(tg)
	if err := awsTG.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsTGOutput, err := a.awsService.UpdateTargetGroup(ctx, *tg.ARN, awsTG)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainTargetGroupFromOutput(awsTGOutput), nil
}

func (a *AWSComputeAdapter) DeleteTargetGroup(ctx context.Context, arn string) error {
	if err := a.awsService.DeleteTargetGroup(ctx, arn); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) ListTargetGroups(ctx context.Context, filters map[string]string) ([]*domaincompute.TargetGroup, error) {
	// Convert domain filters to AWS filters format
	awsFilters := make(map[string][]string)
	for key, value := range filters {
		awsFilters[key] = []string{value}
	}

	awsTGOutputs, err := a.awsService.ListTargetGroups(ctx, awsFilters)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainTGs := make([]*domaincompute.TargetGroup, len(awsTGOutputs))
	for i, awsTGOutput := range awsTGOutputs {
		domainTGs[i] = awsmapper.ToDomainTargetGroupFromOutput(awsTGOutput)
	}

	return domainTGs, nil
}

// Listener Operations

func (a *AWSComputeAdapter) CreateListener(ctx context.Context, listener *domaincompute.Listener) (*domaincompute.Listener, error) {
	if err := listener.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsListener := awsmapper.FromDomainListener(listener)
	if err := awsListener.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsListenerOutput, err := a.awsService.CreateListener(ctx, awsListener)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainListenerFromOutput(awsListenerOutput), nil
}

func (a *AWSComputeAdapter) GetListener(ctx context.Context, arn string) (*domaincompute.Listener, error) {
	awsListenerOutput, err := a.awsService.GetListener(ctx, arn)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainListenerFromOutput(awsListenerOutput), nil
}

func (a *AWSComputeAdapter) UpdateListener(ctx context.Context, listener *domaincompute.Listener) (*domaincompute.Listener, error) {
	if err := listener.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}
	if listener.ARN == nil {
		return nil, fmt.Errorf("listener ARN is required for update")
	}

	awsListener := awsmapper.FromDomainListener(listener)
	if err := awsListener.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsListenerOutput, err := a.awsService.UpdateListener(ctx, *listener.ARN, awsListener)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainListenerFromOutput(awsListenerOutput), nil
}

func (a *AWSComputeAdapter) DeleteListener(ctx context.Context, arn string) error {
	if err := a.awsService.DeleteListener(ctx, arn); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) ListListeners(ctx context.Context, loadBalancerARN string) ([]*domaincompute.Listener, error) {
	awsListenerOutputs, err := a.awsService.ListListeners(ctx, loadBalancerARN)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainListeners := make([]*domaincompute.Listener, len(awsListenerOutputs))
	for i, awsListenerOutput := range awsListenerOutputs {
		domainListeners[i] = awsmapper.ToDomainListenerFromOutput(awsListenerOutput)
	}

	return domainListeners, nil
}

// Target Group Attachment Operations

func (a *AWSComputeAdapter) AttachTargetToGroup(ctx context.Context, attachment *domaincompute.TargetGroupAttachment) error {
	if err := attachment.Validate(); err != nil {
		return fmt.Errorf("domain validation failed: %w", err)
	}

	awsAttachment := awsmapper.FromDomainTargetGroupAttachment(attachment)
	if err := awsAttachment.Validate(); err != nil {
		return fmt.Errorf("aws validation failed: %w", err)
	}

	if err := a.awsService.AttachTargetToGroup(ctx, awsAttachment); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) DetachTargetFromGroup(ctx context.Context, targetGroupARN, targetID string) error {
	if err := a.awsService.DetachTargetFromGroup(ctx, targetGroupARN, targetID); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) ListTargetGroupTargets(ctx context.Context, targetGroupARN string) ([]*domaincompute.TargetGroupAttachment, error) {
	awsAttachmentOutputs, err := a.awsService.ListTargetGroupTargets(ctx, targetGroupARN)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainAttachments := make([]*domaincompute.TargetGroupAttachment, len(awsAttachmentOutputs))
	for i, awsAttachmentOutput := range awsAttachmentOutputs {
		domainAttachments[i] = awsmapper.ToDomainTargetGroupAttachmentFromOutput(awsAttachmentOutput)
	}

	return domainAttachments, nil
}

// Auto Scaling Group operations

func (a *AWSComputeAdapter) CreateAutoScalingGroup(ctx context.Context, asg *domaincompute.AutoScalingGroup) (*domaincompute.AutoScalingGroup, error) {
	if err := asg.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsASG := awsmapper.FromDomainAutoScalingGroup(asg)
	if err := awsASG.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsASGOutput, err := a.awsService.CreateAutoScalingGroup(ctx, awsASG)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainAutoScalingGroupFromOutput(awsASGOutput), nil
}

func (a *AWSComputeAdapter) GetAutoScalingGroup(ctx context.Context, name string) (*domaincompute.AutoScalingGroup, error) {
	awsASGOutput, err := a.awsService.GetAutoScalingGroup(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainAutoScalingGroupFromOutput(awsASGOutput), nil
}

func (a *AWSComputeAdapter) UpdateAutoScalingGroup(ctx context.Context, asg *domaincompute.AutoScalingGroup) (*domaincompute.AutoScalingGroup, error) {
	if err := asg.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsASG := awsmapper.FromDomainAutoScalingGroup(asg)
	if err := awsASG.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	// Use name from domain ASG (ID field contains the name)
	asgName := asg.ID
	if asgName == "" {
		asgName = asg.Name
	}

	awsASGOutput, err := a.awsService.UpdateAutoScalingGroup(ctx, asgName, awsASG)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainAutoScalingGroupFromOutput(awsASGOutput), nil
}

func (a *AWSComputeAdapter) DeleteAutoScalingGroup(ctx context.Context, name string) error {
	if err := a.awsService.DeleteAutoScalingGroup(ctx, name); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) ListAutoScalingGroups(ctx context.Context, filters map[string]string) ([]*domaincompute.AutoScalingGroup, error) {
	// Convert domain filters to AWS filters format
	awsFilters := make(map[string][]string)
	for k, v := range filters {
		awsFilters[k] = []string{v}
	}

	awsASGs, err := a.awsService.ListAutoScalingGroups(ctx, awsFilters)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainASGs := make([]*domaincompute.AutoScalingGroup, len(awsASGs))
	for i, awsASG := range awsASGs {
		domainASGs[i] = awsmapper.ToDomainAutoScalingGroupFromOutput(awsASG)
	}

	return domainASGs, nil
}

// Scaling operations

func (a *AWSComputeAdapter) SetDesiredCapacity(ctx context.Context, asgName string, capacity int) error {
	if err := a.awsService.SetDesiredCapacity(ctx, asgName, capacity); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) AttachInstances(ctx context.Context, asgName string, instanceIDs []string) error {
	if err := a.awsService.AttachInstances(ctx, asgName, instanceIDs); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) DetachInstances(ctx context.Context, asgName string, instanceIDs []string) error {
	if err := a.awsService.DetachInstances(ctx, asgName, instanceIDs); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

// Lambda Function Operations

func (a *AWSComputeAdapter) CreateLambdaFunction(ctx context.Context, function *domaincompute.LambdaFunction) (*domaincompute.LambdaFunction, error) {
	if err := function.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsFunction := awsmapper.FromDomainLambdaFunction(function)
	if err := awsFunction.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsFunctionOutput, err := a.awsService.CreateLambdaFunction(ctx, awsFunction)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainFunction := awsmapper.ToDomainLambdaFunctionFromOutput(awsFunctionOutput)
	// Preserve region from input
	domainFunction.Region = function.Region
	return domainFunction, nil
}

func (a *AWSComputeAdapter) GetLambdaFunction(ctx context.Context, name string) (*domaincompute.LambdaFunction, error) {
	awsFunctionOutput, err := a.awsService.GetLambdaFunction(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainLambdaFunctionFromOutput(awsFunctionOutput), nil
}

func (a *AWSComputeAdapter) UpdateLambdaFunction(ctx context.Context, function *domaincompute.LambdaFunction) (*domaincompute.LambdaFunction, error) {
	if err := function.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsFunction := awsmapper.FromDomainLambdaFunction(function)
	if err := awsFunction.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsFunctionOutput, err := a.awsService.UpdateLambdaFunction(ctx, function.FunctionName, awsFunction)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainFunction := awsmapper.ToDomainLambdaFunctionFromOutput(awsFunctionOutput)
	// Preserve region from input
	domainFunction.Region = function.Region
	return domainFunction, nil
}

func (a *AWSComputeAdapter) DeleteLambdaFunction(ctx context.Context, name string) error {
	if err := a.awsService.DeleteLambdaFunction(ctx, name); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSComputeAdapter) ListLambdaFunctions(ctx context.Context, filters map[string]string) ([]*domaincompute.LambdaFunction, error) {
	// Convert domain filters to AWS filters format
	awsFilters := make(map[string][]string)
	for k, v := range filters {
		awsFilters[k] = []string{v}
	}

	awsFunctionOutputs, err := a.awsService.ListLambdaFunctions(ctx, awsFilters)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainFunctions := make([]*domaincompute.LambdaFunction, len(awsFunctionOutputs))
	for i, output := range awsFunctionOutputs {
		domainFunctions[i] = awsmapper.ToDomainLambdaFunctionFromOutput(output)
	}

	return domainFunctions, nil
}
