package compute

import (
	"context"

	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awsloadbalanceroutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer/outputs"
)

// ComputeService implements AWSComputeService using AWS SDK
// Note: This is a partial implementation - instance and launch template methods
// are implemented elsewhere or need to be added
type ComputeService struct {
	client *awssdk.AWSClient
}

// NewComputeService creates a new compute service implementation
func NewComputeService(client *awssdk.AWSClient) *ComputeService {
	return &ComputeService{
		client: client,
	}
}

// Note: This service only implements Load Balancer methods.
// Instance and Launch Template methods are implemented elsewhere or need to be added.

// Load Balancer operations

func (s *ComputeService) CreateLoadBalancer(ctx context.Context, lb *awsloadbalancer.LoadBalancer) (*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	return awssdk.CreateLoadBalancer(ctx, s.client, lb)
}

func (s *ComputeService) GetLoadBalancer(ctx context.Context, arn string) (*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	return awssdk.GetLoadBalancer(ctx, s.client, arn)
}

func (s *ComputeService) UpdateLoadBalancer(ctx context.Context, arn string, lb *awsloadbalancer.LoadBalancer) (*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	return awssdk.UpdateLoadBalancer(ctx, s.client, arn, lb)
}

func (s *ComputeService) DeleteLoadBalancer(ctx context.Context, arn string) error {
	return awssdk.DeleteLoadBalancer(ctx, s.client, arn)
}

func (s *ComputeService) ListLoadBalancers(ctx context.Context, filters map[string][]string) ([]*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	return awssdk.ListLoadBalancers(ctx, s.client, filters)
}

// Target Group operations

func (s *ComputeService) CreateTargetGroup(ctx context.Context, tg *awsloadbalancer.TargetGroup) (*awsloadbalanceroutputs.TargetGroupOutput, error) {
	return awssdk.CreateTargetGroup(ctx, s.client, tg)
}

func (s *ComputeService) GetTargetGroup(ctx context.Context, arn string) (*awsloadbalanceroutputs.TargetGroupOutput, error) {
	return awssdk.GetTargetGroup(ctx, s.client, arn)
}

func (s *ComputeService) UpdateTargetGroup(ctx context.Context, arn string, tg *awsloadbalancer.TargetGroup) (*awsloadbalanceroutputs.TargetGroupOutput, error) {
	return awssdk.UpdateTargetGroup(ctx, s.client, arn, tg)
}

func (s *ComputeService) DeleteTargetGroup(ctx context.Context, arn string) error {
	return awssdk.DeleteTargetGroup(ctx, s.client, arn)
}

func (s *ComputeService) ListTargetGroups(ctx context.Context, filters map[string][]string) ([]*awsloadbalanceroutputs.TargetGroupOutput, error) {
	return awssdk.ListTargetGroups(ctx, s.client, filters)
}

// Listener operations

func (s *ComputeService) CreateListener(ctx context.Context, listener *awsloadbalancer.Listener) (*awsloadbalanceroutputs.ListenerOutput, error) {
	return awssdk.CreateListener(ctx, s.client, listener)
}

func (s *ComputeService) GetListener(ctx context.Context, arn string) (*awsloadbalanceroutputs.ListenerOutput, error) {
	return awssdk.GetListener(ctx, s.client, arn)
}

func (s *ComputeService) UpdateListener(ctx context.Context, arn string, listener *awsloadbalancer.Listener) (*awsloadbalanceroutputs.ListenerOutput, error) {
	return awssdk.UpdateListener(ctx, s.client, arn, listener)
}

func (s *ComputeService) DeleteListener(ctx context.Context, arn string) error {
	return awssdk.DeleteListener(ctx, s.client, arn)
}

func (s *ComputeService) ListListeners(ctx context.Context, loadBalancerARN string) ([]*awsloadbalanceroutputs.ListenerOutput, error) {
	return awssdk.ListListeners(ctx, s.client, loadBalancerARN)
}

// Target Group Attachment operations

func (s *ComputeService) AttachTargetToGroup(ctx context.Context, attachment *awsloadbalancer.TargetGroupAttachment) error {
	return awssdk.AttachTargetToGroup(ctx, s.client, attachment)
}

func (s *ComputeService) DetachTargetFromGroup(ctx context.Context, targetGroupARN, targetID string) error {
	return awssdk.DetachTargetFromGroup(ctx, s.client, targetGroupARN, targetID)
}

func (s *ComputeService) ListTargetGroupTargets(ctx context.Context, targetGroupARN string) ([]*awsloadbalanceroutputs.TargetGroupAttachmentOutput, error) {
	return awssdk.ListTargetGroupTargets(ctx, s.client, targetGroupARN)
}
