package compute

import (
	"context"
	"errors"
	"testing"
	"time"

	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsec2 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2"
	awslttemplate "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template"
	awslttemplateoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template/outputs"
	awsec2outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/outputs"
	awsInstanceTypes "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/instance_types"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awsloadbalanceroutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer/outputs"
	awsautoscaling "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling"
	awsautoscalingoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling/outputs"
	awslambda "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda"
	awslambdaoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda/outputs"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
)

// mockAWSComputeService is a mock implementation of AWSComputeService for testing
type mockAWSComputeService struct {
	instance              *awsec2.Instance
	launchTemplate        *awslttemplate.LaunchTemplate
	loadBalancer          *awsloadbalancer.LoadBalancer
	targetGroup           *awsloadbalancer.TargetGroup
	listener              *awsloadbalancer.Listener
	targetGroupAttachment *awsloadbalancer.TargetGroupAttachment
	autoScalingGroup      *awsautoscaling.AutoScalingGroup
	lambdaFunction        *awslambda.Function
	createError           error
	getError              error
	instanceTypeInfo      *awsInstanceTypes.InstanceTypeInfo
}

// Ensure mockAWSComputeService implements AWSComputeService
var _ awsservice.AWSComputeService = (*mockAWSComputeService)(nil)

// Helper function to convert Instance input to output
func instanceToOutput(instance *awsec2.Instance) *awsec2outputs.InstanceOutput {
	if instance == nil {
		return nil
	}
	return &awsec2outputs.InstanceOutput{
		ID:                 "i-mock-1234567890abcdef0",
		ARN:                "arn:aws:ec2:us-east-1:123456789012:instance/i-mock-1234567890abcdef0",
		Name:               instance.Name,
		Region:             "us-east-1",
		AvailabilityZone:   "us-east-1a",
		InstanceType:       instance.InstanceType,
		AMI:                instance.AMI,
		SubnetID:           instance.SubnetID,
		SecurityGroupIDs:   instance.VpcSecurityGroupIds,
		PrivateIP:          "10.0.1.100",
		PublicIP:           nil,
		KeyName:            instance.KeyName,
		IAMInstanceProfile: instance.IAMInstanceProfile,
		State:              "running",
		CreationTime:       time.Now(),
		Tags:               instance.Tags,
	}
}

// Helper function to convert LaunchTemplate input to output
func launchTemplateToOutput(template *awslttemplate.LaunchTemplate) *awslttemplateoutputs.LaunchTemplateOutput {
	if template == nil {
		return nil
	}
	name := "test-template"
	if template.Name != nil {
		name = *template.Name
	} else if template.NamePrefix != nil {
		name = *template.NamePrefix + "-12345"
	}
	return &awslttemplateoutputs.LaunchTemplateOutput{
		ID:             "lt-mock-1234567890abcdef0",
		ARN:            "arn:aws:ec2:us-east-1:123456789012:launch-template/lt-mock-1234567890abcdef0",
		Name:           name,
		DefaultVersion: 1,
		LatestVersion:  1,
		CreateTime:     time.Now(),
		CreatedBy:      "arn:aws:iam::123456789012:user/test",
		Tags: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			{Key: "Name", Value: name},
		},
	}
}

// Instance Operations

func (m *mockAWSComputeService) CreateInstance(ctx context.Context, instance *awsec2.Instance) (*awsec2outputs.InstanceOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.instance = instance
	return instanceToOutput(instance), nil
}

func (m *mockAWSComputeService) GetInstance(ctx context.Context, id string) (*awsec2outputs.InstanceOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return instanceToOutput(m.instance), nil
}

func (m *mockAWSComputeService) UpdateInstance(ctx context.Context, instance *awsec2.Instance) (*awsec2outputs.InstanceOutput, error) {
	m.instance = instance
	return instanceToOutput(instance), nil
}

func (m *mockAWSComputeService) DeleteInstance(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSComputeService) ListInstances(ctx context.Context, filters map[string][]string) ([]*awsec2outputs.InstanceOutput, error) {
	if m.instance != nil {
		return []*awsec2outputs.InstanceOutput{instanceToOutput(m.instance)}, nil
	}
	return []*awsec2outputs.InstanceOutput{}, nil
}

// Instance Lifecycle Operations

func (m *mockAWSComputeService) StartInstance(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSComputeService) StopInstance(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSComputeService) RebootInstance(ctx context.Context, id string) error {
	return nil
}

// Instance Type Operations

func (m *mockAWSComputeService) GetInstanceTypeInfo(ctx context.Context, instanceType string, region string) (*awsInstanceTypes.InstanceTypeInfo, error) {
	if m.instanceTypeInfo != nil {
		return m.instanceTypeInfo, nil
	}
	return &awsInstanceTypes.InstanceTypeInfo{
		Name:           instanceType,
		Category:       awsInstanceTypes.CategoryGeneralPurpose,
		VCPU:           2,
		MemoryGiB:      4.0,
		MaxNetworkGbps: 5,
		StorageType:    "EBS only",
	}, nil
}

func (m *mockAWSComputeService) ListInstanceTypesByCategory(ctx context.Context, category awsInstanceTypes.InstanceCategory, region string) ([]*awsInstanceTypes.InstanceTypeInfo, error) {
	return []*awsInstanceTypes.InstanceTypeInfo{
		{
			Name:           "t3.micro",
			Category:       category,
			VCPU:           2,
			MemoryGiB:      1.0,
			MaxNetworkGbps: 5,
			StorageType:    "EBS",
		},
	}, nil
}

// Launch Template Operations

func (m *mockAWSComputeService) CreateLaunchTemplate(ctx context.Context, template *awslttemplate.LaunchTemplate) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.launchTemplate = template
	return launchTemplateToOutput(template), nil
}

func (m *mockAWSComputeService) GetLaunchTemplate(ctx context.Context, id string) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return launchTemplateToOutput(m.launchTemplate), nil
}

func (m *mockAWSComputeService) UpdateLaunchTemplate(ctx context.Context, id string, template *awslttemplate.LaunchTemplate) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	m.launchTemplate = template
	return launchTemplateToOutput(template), nil
}

func (m *mockAWSComputeService) DeleteLaunchTemplate(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSComputeService) ListLaunchTemplates(ctx context.Context, filters map[string][]string) ([]*awslttemplateoutputs.LaunchTemplateOutput, error) {
	if m.launchTemplate != nil {
		return []*awslttemplateoutputs.LaunchTemplateOutput{launchTemplateToOutput(m.launchTemplate)}, nil
	}
	return []*awslttemplateoutputs.LaunchTemplateOutput{}, nil
}

func (m *mockAWSComputeService) GetLaunchTemplateVersion(ctx context.Context, id string, version int) (*awslttemplate.LaunchTemplateVersion, error) {
	return &awslttemplate.LaunchTemplateVersion{
		TemplateID:    id,
		VersionNumber: version,
		IsDefault:     version == 1,
		CreateTime:    time.Now(),
		CreatedBy:     stringPtr("arn:aws:iam::123456789012:user/test"),
		TemplateData:  m.launchTemplate,
	}, nil
}

func (m *mockAWSComputeService) ListLaunchTemplateVersions(ctx context.Context, id string) ([]*awslttemplate.LaunchTemplateVersion, error) {
	return []*awslttemplate.LaunchTemplateVersion{
		{
			TemplateID:    id,
			VersionNumber: 1,
			IsDefault:     true,
			CreateTime:    time.Now(),
			CreatedBy:     stringPtr("arn:aws:iam::123456789012:user/test"),
			TemplateData:  m.launchTemplate,
		},
	}, nil
}

// Helper function to convert LoadBalancer input to output
func loadBalancerToOutput(lb *awsloadbalancer.LoadBalancer) *awsloadbalanceroutputs.LoadBalancerOutput {
	if lb == nil {
		return nil
	}
	internal := false
	if lb.Internal != nil {
		internal = *lb.Internal
	}
	return &awsloadbalanceroutputs.LoadBalancerOutput{
		ARN:              "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef",
		ID:               "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef",
		Name:             lb.Name,
		DNSName:          "test-lb-1234567890.us-east-1.elb.amazonaws.com",
		ZoneID:           "Z35SXDOTRQ7X7K",
		Type:             lb.LoadBalancerType,
		Internal:         internal,
		SecurityGroupIDs: lb.SecurityGroupIDs,
		SubnetIDs:        lb.SubnetIDs,
		State:            "active",
		CreatedTime:      time.Now(),
	}
}

// Helper function to convert TargetGroup input to output
func targetGroupToOutput(tg *awsloadbalancer.TargetGroup) *awsloadbalanceroutputs.TargetGroupOutput {
	if tg == nil {
		return nil
	}
	targetType := "instance"
	if tg.TargetType != nil {
		targetType = *tg.TargetType
	}
	return &awsloadbalanceroutputs.TargetGroupOutput{
		ARN:         "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef",
		ID:          "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef",
		Name:        tg.Name,
		Port:        tg.Port,
		Protocol:    tg.Protocol,
		VPCID:       tg.VPCID,
		TargetType:  targetType,
		HealthCheck: tg.HealthCheck,
		State:       "active",
		CreatedTime: time.Now(),
	}
}

// Helper function to convert Listener input to output
func listenerToOutput(listener *awsloadbalancer.Listener) *awsloadbalanceroutputs.ListenerOutput {
	if listener == nil {
		return nil
	}
	return &awsloadbalanceroutputs.ListenerOutput{
		ARN:             "arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/test-lb/1234567890abcdef/1234567890abcdef",
		ID:              "arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/test-lb/1234567890abcdef/1234567890abcdef",
		LoadBalancerARN: listener.LoadBalancerARN,
		Port:            listener.Port,
		Protocol:        listener.Protocol,
		DefaultAction:   listener.DefaultAction,
	}
}

// Load Balancer operations
func (m *mockAWSComputeService) CreateLoadBalancer(ctx context.Context, lb *awsloadbalancer.LoadBalancer) (*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.loadBalancer = lb
	return loadBalancerToOutput(lb), nil
}

func (m *mockAWSComputeService) GetLoadBalancer(ctx context.Context, arn string) (*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return loadBalancerToOutput(m.loadBalancer), nil
}

func (m *mockAWSComputeService) UpdateLoadBalancer(ctx context.Context, arn string, lb *awsloadbalancer.LoadBalancer) (*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	m.loadBalancer = lb
	return loadBalancerToOutput(lb), nil
}

func (m *mockAWSComputeService) DeleteLoadBalancer(ctx context.Context, arn string) error {
	return nil
}

func (m *mockAWSComputeService) ListLoadBalancers(ctx context.Context, filters map[string][]string) ([]*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	if m.loadBalancer != nil {
		return []*awsloadbalanceroutputs.LoadBalancerOutput{loadBalancerToOutput(m.loadBalancer)}, nil
	}
	return []*awsloadbalanceroutputs.LoadBalancerOutput{}, nil
}

// Target Group operations
func (m *mockAWSComputeService) CreateTargetGroup(ctx context.Context, tg *awsloadbalancer.TargetGroup) (*awsloadbalanceroutputs.TargetGroupOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.targetGroup = tg
	return targetGroupToOutput(tg), nil
}

func (m *mockAWSComputeService) GetTargetGroup(ctx context.Context, arn string) (*awsloadbalanceroutputs.TargetGroupOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return targetGroupToOutput(m.targetGroup), nil
}

func (m *mockAWSComputeService) UpdateTargetGroup(ctx context.Context, arn string, tg *awsloadbalancer.TargetGroup) (*awsloadbalanceroutputs.TargetGroupOutput, error) {
	m.targetGroup = tg
	return targetGroupToOutput(tg), nil
}

func (m *mockAWSComputeService) DeleteTargetGroup(ctx context.Context, arn string) error {
	return nil
}

func (m *mockAWSComputeService) ListTargetGroups(ctx context.Context, filters map[string][]string) ([]*awsloadbalanceroutputs.TargetGroupOutput, error) {
	if m.targetGroup != nil {
		return []*awsloadbalanceroutputs.TargetGroupOutput{targetGroupToOutput(m.targetGroup)}, nil
	}
	return []*awsloadbalanceroutputs.TargetGroupOutput{}, nil
}

// Listener operations
func (m *mockAWSComputeService) CreateListener(ctx context.Context, listener *awsloadbalancer.Listener) (*awsloadbalanceroutputs.ListenerOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.listener = listener
	return listenerToOutput(listener), nil
}

func (m *mockAWSComputeService) GetListener(ctx context.Context, arn string) (*awsloadbalanceroutputs.ListenerOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return listenerToOutput(m.listener), nil
}

func (m *mockAWSComputeService) UpdateListener(ctx context.Context, arn string, listener *awsloadbalancer.Listener) (*awsloadbalanceroutputs.ListenerOutput, error) {
	m.listener = listener
	return listenerToOutput(listener), nil
}

func (m *mockAWSComputeService) DeleteListener(ctx context.Context, arn string) error {
	return nil
}

func (m *mockAWSComputeService) ListListeners(ctx context.Context, loadBalancerARN string) ([]*awsloadbalanceroutputs.ListenerOutput, error) {
	if m.listener != nil {
		return []*awsloadbalanceroutputs.ListenerOutput{listenerToOutput(m.listener)}, nil
	}
	return []*awsloadbalanceroutputs.ListenerOutput{}, nil
}

// Target Group Attachment operations
func (m *mockAWSComputeService) AttachTargetToGroup(ctx context.Context, attachment *awsloadbalancer.TargetGroupAttachment) error {
	m.targetGroupAttachment = attachment
	return nil
}

func (m *mockAWSComputeService) DetachTargetFromGroup(ctx context.Context, targetGroupARN, targetID string) error {
	return nil
}

func (m *mockAWSComputeService) ListTargetGroupTargets(ctx context.Context, targetGroupARN string) ([]*awsloadbalanceroutputs.TargetGroupAttachmentOutput, error) {
	if m.targetGroupAttachment != nil {
		return []*awsloadbalanceroutputs.TargetGroupAttachmentOutput{
			{
				TargetGroupARN:   m.targetGroupAttachment.TargetGroupARN,
				TargetID:         m.targetGroupAttachment.TargetID,
				Port:             m.targetGroupAttachment.Port,
				AvailabilityZone: m.targetGroupAttachment.AvailabilityZone,
				HealthStatus:     "healthy",
				State:            "healthy",
			},
		}, nil
	}
	return []*awsloadbalanceroutputs.TargetGroupAttachmentOutput{}, nil
}

// Tests

func TestAWSComputeAdapter_CreateInstance(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	domainInstance := &domaincompute.Instance{
		Name:             "test-instance",
		Region:           "us-east-1",
		InstanceType:     "t3.micro",
		AMI:              "ami-0c55b159cbfafe1f0",
		SubnetID:         "subnet-123",
		SecurityGroupIDs: []string{"sg-123"},
	}

	ctx := context.Background()
	createdInstance, err := adapter.CreateInstance(ctx, domainInstance)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdInstance == nil {
		t.Fatal("Expected created instance, got nil")
	}

	if createdInstance.Name != domainInstance.Name {
		t.Errorf("Expected name %s, got %s", domainInstance.Name, createdInstance.Name)
	}

	if createdInstance.ID == "" {
		t.Error("Expected instance ID to be populated")
	}
}

func TestAWSComputeAdapter_CreateInstance_ValidationError(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	invalidInstance := &domaincompute.Instance{
		Name:   "", // Invalid: empty name
		Region: "us-east-1",
		AMI:    "ami-0c55b159cbfafe1f0",
	}

	ctx := context.Background()
	_, err := adapter.CreateInstance(ctx, invalidInstance)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
}

func TestAWSComputeAdapter_GetInstance(t *testing.T) {
	mockService := &mockAWSComputeService{
		instance: &awsec2.Instance{
			Name:                "test-instance",
			InstanceType:        "t3.micro",
			AMI:                 "ami-0c55b159cbfafe1f0",
			SubnetID:            "subnet-123",
			VpcSecurityGroupIds: []string{"sg-123"},
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	instance, err := adapter.GetInstance(ctx, "i-123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if instance == nil {
		t.Fatal("Expected instance, got nil")
	}

	if instance.Name != "test-instance" {
		t.Errorf("Expected name test-instance, got %s", instance.Name)
	}
}

func TestAWSComputeAdapter_ListInstances(t *testing.T) {
	mockService := &mockAWSComputeService{
		instance: &awsec2.Instance{
			Name:                "test-instance",
			InstanceType:        "t3.micro",
			AMI:                 "ami-0c55b159cbfafe1f0",
			SubnetID:            "subnet-123",
			VpcSecurityGroupIds: []string{"sg-123"},
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	instances, err := adapter.ListInstances(ctx, map[string]string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(instances) != 1 {
		t.Errorf("Expected 1 instance, got %d", len(instances))
	}

	if instances[0].Name != "test-instance" {
		t.Errorf("Expected name test-instance, got %s", instances[0].Name)
	}
}

func TestAWSComputeAdapter_CreateLaunchTemplate(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	namePrefix := "test-template"
	domainTemplate := &domaincompute.LaunchTemplate{
		Name:             "test-template",
		Region:           "us-east-1",
		NamePrefix:       &namePrefix,
		ImageID:          "ami-0c55b159cbfafe1f0",
		InstanceType:     "t3.micro",
		SecurityGroupIDs: []string{"sg-123"},
	}

	ctx := context.Background()
	createdTemplate, err := adapter.CreateLaunchTemplate(ctx, domainTemplate)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdTemplate == nil {
		t.Fatal("Expected created template, got nil")
	}

	if createdTemplate.ID == "" {
		t.Error("Expected template ID to be populated")
	}
}

func TestAWSComputeAdapter_ErrorHandling(t *testing.T) {
	mockService := &mockAWSComputeService{
		getError: errors.New("aws service error"),
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	_, err := adapter.GetInstance(ctx, "i-123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Verify error is wrapped
	if err.Error() == "" {
		t.Error("Expected error message, got empty string")
	}
}

// Load Balancer Tests

func TestAWSComputeAdapter_CreateLoadBalancer(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	domainLB := &domaincompute.LoadBalancer{
		Name:             "test-alb",
		Region:           "us-east-1",
		Type:             domaincompute.LoadBalancerTypeApplication,
		Internal:         false,
		SecurityGroupIDs: []string{"sg-123", "sg-456"},
		SubnetIDs:        []string{"subnet-123", "subnet-456"},
	}

	ctx := context.Background()
	createdLB, err := adapter.CreateLoadBalancer(ctx, domainLB)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdLB == nil {
		t.Fatal("Expected created load balancer, got nil")
	}

	if createdLB.Name != domainLB.Name {
		t.Errorf("Expected name %s, got %s", domainLB.Name, createdLB.Name)
	}

	if createdLB.ID == "" {
		t.Error("Expected load balancer ID to be populated")
	}

	if createdLB.ARN == nil {
		t.Error("Expected load balancer ARN to be populated")
	}

	if createdLB.DNSName == nil {
		t.Error("Expected DNS name to be populated")
	}
}

func TestAWSComputeAdapter_CreateLoadBalancer_ValidationError(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	invalidLB := &domaincompute.LoadBalancer{
		Name:      "", // Invalid: empty name
		Region:    "us-east-1",
		Type:      domaincompute.LoadBalancerTypeApplication,
		SubnetIDs: []string{"subnet-123"}, // Invalid: only one subnet
	}

	ctx := context.Background()
	_, err := adapter.CreateLoadBalancer(ctx, invalidLB)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
}

func TestAWSComputeAdapter_GetLoadBalancer(t *testing.T) {
	mockService := &mockAWSComputeService{
		loadBalancer: &awsloadbalancer.LoadBalancer{
			Name:             "test-alb",
			LoadBalancerType: "application",
			Internal:         boolPtr(false),
			SecurityGroupIDs: []string{"sg-123"},
			SubnetIDs:        []string{"subnet-123", "subnet-456"},
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	lb, err := adapter.GetLoadBalancer(ctx, "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if lb == nil {
		t.Fatal("Expected load balancer, got nil")
	}

	if lb.Name != "test-alb" {
		t.Errorf("Expected name test-alb, got %s", lb.Name)
	}

	if lb.ID == "" {
		t.Error("Expected load balancer ID to be populated")
	}
}

func TestAWSComputeAdapter_ListLoadBalancers(t *testing.T) {
	mockService := &mockAWSComputeService{
		loadBalancer: &awsloadbalancer.LoadBalancer{
			Name:             "test-alb",
			LoadBalancerType: "application",
			Internal:         boolPtr(false),
			SecurityGroupIDs: []string{"sg-123"},
			SubnetIDs:        []string{"subnet-123", "subnet-456"},
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	lbs, err := adapter.ListLoadBalancers(ctx, map[string]string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(lbs) != 1 {
		t.Errorf("Expected 1 load balancer, got %d", len(lbs))
	}

	if lbs[0].Name != "test-alb" {
		t.Errorf("Expected name test-alb, got %s", lbs[0].Name)
	}
}

func TestAWSComputeAdapter_CreateTargetGroup(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	domainTG := &domaincompute.TargetGroup{
		Name:       "test-tg",
		VPCID:      "vpc-123",
		Port:       80,
		Protocol:   domaincompute.TargetGroupProtocolHTTP,
		TargetType: domaincompute.TargetTypeInstance,
		HealthCheck: domaincompute.HealthCheckConfig{
			Path:    stringPtr("/health"),
			Matcher: stringPtr("200"),
		},
	}

	ctx := context.Background()
	createdTG, err := adapter.CreateTargetGroup(ctx, domainTG)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdTG == nil {
		t.Fatal("Expected created target group, got nil")
	}

	if createdTG.Name != domainTG.Name {
		t.Errorf("Expected name %s, got %s", domainTG.Name, createdTG.Name)
	}

	if createdTG.ID == "" {
		t.Error("Expected target group ID to be populated")
	}

	if createdTG.ARN == nil {
		t.Error("Expected target group ARN to be populated")
	}
}

func TestAWSComputeAdapter_GetTargetGroup(t *testing.T) {
	mockService := &mockAWSComputeService{
		targetGroup: &awsloadbalancer.TargetGroup{
			Name:       "test-tg",
			VPCID:      "vpc-123",
			Port:       80,
			Protocol:   "HTTP",
			TargetType: stringPtr("instance"),
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	tg, err := adapter.GetTargetGroup(ctx, "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if tg == nil {
		t.Fatal("Expected target group, got nil")
	}

	if tg.Name != "test-tg" {
		t.Errorf("Expected name test-tg, got %s", tg.Name)
	}
}

func TestAWSComputeAdapter_CreateListener(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	targetGroupARN := "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef"
	domainListener := &domaincompute.Listener{
		LoadBalancerARN: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef",
		Port:            80,
		Protocol:        domaincompute.ListenerProtocolHTTP,
		DefaultAction: domaincompute.ListenerAction{
			Type:           domaincompute.ListenerActionTypeForward,
			TargetGroupARN: &targetGroupARN,
		},
	}

	ctx := context.Background()
	createdListener, err := adapter.CreateListener(ctx, domainListener)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdListener == nil {
		t.Fatal("Expected created listener, got nil")
	}

	if createdListener.Port != domainListener.Port {
		t.Errorf("Expected port %d, got %d", domainListener.Port, createdListener.Port)
	}

	if createdListener.ID == "" {
		t.Error("Expected listener ID to be populated")
	}

	if createdListener.ARN == nil {
		t.Error("Expected listener ARN to be populated")
	}
}

func TestAWSComputeAdapter_AttachTargetToGroup(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	domainAttachment := &domaincompute.TargetGroupAttachment{
		TargetGroupARN: "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef",
		TargetID:       "i-1234567890abcdef0",
		Port:           intPtr(8080),
	}

	ctx := context.Background()
	err := adapter.AttachTargetToGroup(ctx, domainAttachment)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAWSComputeAdapter_ListTargetGroupTargets(t *testing.T) {
	mockService := &mockAWSComputeService{
		targetGroupAttachment: &awsloadbalancer.TargetGroupAttachment{
			TargetGroupARN: "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef",
			TargetID:       "i-1234567890abcdef0",
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	targets, err := adapter.ListTargetGroupTargets(ctx, "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(targets) != 1 {
		t.Errorf("Expected 1 target, got %d", len(targets))
	}

	if targets[0].TargetID != "i-1234567890abcdef0" {
		t.Errorf("Expected target ID i-1234567890abcdef0, got %s", targets[0].TargetID)
	}
}

func TestAWSComputeAdapter_UpdateLoadBalancer(t *testing.T) {
	mockService := &mockAWSComputeService{
		loadBalancer: &awsloadbalancer.LoadBalancer{
			Name:             "test-alb",
			LoadBalancerType: "application",
			Internal:         boolPtr(false),
			SecurityGroupIDs: []string{"sg-123"},
			SubnetIDs:        []string{"subnet-123", "subnet-456"},
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	arn := "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef"
	domainLB := &domaincompute.LoadBalancer{
		ID:               arn,
		ARN:              &arn,
		Name:             "updated-alb",
		Region:           "us-east-1",
		Type:             domaincompute.LoadBalancerTypeApplication,
		Internal:         false,
		SecurityGroupIDs: []string{"sg-789"},
		SubnetIDs:        []string{"subnet-123", "subnet-456"},
	}

	ctx := context.Background()
	updatedLB, err := adapter.UpdateLoadBalancer(ctx, domainLB)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if updatedLB == nil {
		t.Fatal("Expected updated load balancer, got nil")
	}

	if updatedLB.Name != "updated-alb" {
		t.Errorf("Expected name updated-alb, got %s", updatedLB.Name)
	}
}

func TestAWSComputeAdapter_DeleteLoadBalancer(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	err := adapter.DeleteLoadBalancer(ctx, "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAWSComputeAdapter_UpdateTargetGroup(t *testing.T) {
	mockService := &mockAWSComputeService{
		targetGroup: &awsloadbalancer.TargetGroup{
			Name:     "test-tg",
			VPCID:    "vpc-123",
			Port:     80,
			Protocol: "HTTP",
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	arn := "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef"
	domainTG := &domaincompute.TargetGroup{
		ID:         arn,
		ARN:        &arn,
		Name:       "updated-tg",
		VPCID:      "vpc-123",
		Port:       8080,
		Protocol:   domaincompute.TargetGroupProtocolHTTP,
		TargetType: domaincompute.TargetTypeInstance,
	}

	ctx := context.Background()
	updatedTG, err := adapter.UpdateTargetGroup(ctx, domainTG)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if updatedTG == nil {
		t.Fatal("Expected updated target group, got nil")
	}

	if updatedTG.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", updatedTG.Port)
	}
}

func TestAWSComputeAdapter_DeleteTargetGroup(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	err := adapter.DeleteTargetGroup(ctx, "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAWSComputeAdapter_ListTargetGroups(t *testing.T) {
	mockService := &mockAWSComputeService{
		targetGroup: &awsloadbalancer.TargetGroup{
			Name:     "test-tg",
			VPCID:    "vpc-123",
			Port:     80,
			Protocol: "HTTP",
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	tgs, err := adapter.ListTargetGroups(ctx, map[string]string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(tgs) != 1 {
		t.Errorf("Expected 1 target group, got %d", len(tgs))
	}

	if tgs[0].Name != "test-tg" {
		t.Errorf("Expected name test-tg, got %s", tgs[0].Name)
	}
}

func TestAWSComputeAdapter_GetListener(t *testing.T) {
	targetGroupARN := "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef"
	mockService := &mockAWSComputeService{
		listener: &awsloadbalancer.Listener{
			LoadBalancerARN: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef",
			Port:            80,
			Protocol:        "HTTP",
			DefaultAction: awsloadbalancer.ListenerAction{
				Type:           awsloadbalancer.ListenerActionTypeForward,
				TargetGroupARN: &targetGroupARN,
			},
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	listener, err := adapter.GetListener(ctx, "arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/test-lb/1234567890abcdef/1234567890abcdef")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if listener == nil {
		t.Fatal("Expected listener, got nil")
	}

	if listener.Port != 80 {
		t.Errorf("Expected port 80, got %d", listener.Port)
	}
}

func TestAWSComputeAdapter_UpdateListener(t *testing.T) {
	targetGroupARN := "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef"
	certARN := "arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012"
	mockService := &mockAWSComputeService{
		listener: &awsloadbalancer.Listener{
			LoadBalancerARN: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef",
			Port:            80,
			Protocol:        "HTTP",
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	listenerARN := "arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/test-lb/1234567890abcdef/1234567890abcdef"
	domainListener := &domaincompute.Listener{
		ID:              listenerARN,
		ARN:             &listenerARN,
		LoadBalancerARN: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef",
		Port:            8080, // Keep HTTP for simplicity
		Protocol:        domaincompute.ListenerProtocolHTTP,
		DefaultAction: domaincompute.ListenerAction{
			Type:           domaincompute.ListenerActionTypeForward,
			TargetGroupARN: &targetGroupARN,
		},
	}

	ctx := context.Background()
	updatedListener, err := adapter.UpdateListener(ctx, domainListener)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if updatedListener == nil {
		t.Fatal("Expected updated listener, got nil")
	}

	if updatedListener.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", updatedListener.Port)
	}

	_ = certARN // Suppress unused variable warning
}

func TestAWSComputeAdapter_DeleteListener(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	err := adapter.DeleteListener(ctx, "arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/test-lb/1234567890abcdef/1234567890abcdef")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAWSComputeAdapter_ListListeners(t *testing.T) {
	targetGroupARN := "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef"
	mockService := &mockAWSComputeService{
		listener: &awsloadbalancer.Listener{
			LoadBalancerARN: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef",
			Port:            80,
			Protocol:        "HTTP",
			DefaultAction: awsloadbalancer.ListenerAction{
				Type:           awsloadbalancer.ListenerActionTypeForward,
				TargetGroupARN: &targetGroupARN,
			},
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	listeners, err := adapter.ListListeners(ctx, "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(listeners) != 1 {
		t.Errorf("Expected 1 listener, got %d", len(listeners))
	}

	if listeners[0].Port != 80 {
		t.Errorf("Expected port 80, got %d", listeners[0].Port)
	}
}

func TestAWSComputeAdapter_DetachTargetFromGroup(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	err := adapter.DetachTargetFromGroup(ctx, "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef", "i-1234567890abcdef0")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// Helper function to convert AutoScalingGroup input to output
func autoScalingGroupToOutput(asg *awsautoscaling.AutoScalingGroup) *awsautoscalingoutputs.AutoScalingGroupOutput {
	if asg == nil {
		return nil
	}
	name := "test-asg"
	if asg.AutoScalingGroupName != nil {
		name = *asg.AutoScalingGroupName
	}
	desiredCapacity := 2
	if asg.DesiredCapacity != nil {
		desiredCapacity = *asg.DesiredCapacity
	}
	healthCheckType := "EC2"
	if asg.HealthCheckType != nil {
		healthCheckType = *asg.HealthCheckType
	}
	return &awsautoscalingoutputs.AutoScalingGroupOutput{
		AutoScalingGroupARN:  "arn:aws:autoscaling:us-east-1:123456789012:autoScalingGroup:uuid:autoScalingGroupName/" + name,
		AutoScalingGroupName:  name,
		MinSize:              asg.MinSize,
		MaxSize:              asg.MaxSize,
		DesiredCapacity:     desiredCapacity,
		VPCZoneIdentifier:   asg.VPCZoneIdentifier,
		LaunchTemplate:      asg.LaunchTemplate,
		HealthCheckType:     healthCheckType,
		HealthCheckGracePeriod: asg.HealthCheckGracePeriod,
		TargetGroupARNs:     asg.TargetGroupARNs,
		Status:              "active",
		CreatedTime:         time.Now(),
		Instances:           []awsautoscalingoutputs.Instance{},
		Tags:                asg.Tags,
	}
}

// Auto Scaling Group operations for mock service
func (m *mockAWSComputeService) CreateAutoScalingGroup(ctx context.Context, asg *awsautoscaling.AutoScalingGroup) (*awsautoscalingoutputs.AutoScalingGroupOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.autoScalingGroup = asg
	return autoScalingGroupToOutput(asg), nil
}

func (m *mockAWSComputeService) GetAutoScalingGroup(ctx context.Context, name string) (*awsautoscalingoutputs.AutoScalingGroupOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return autoScalingGroupToOutput(m.autoScalingGroup), nil
}

func (m *mockAWSComputeService) UpdateAutoScalingGroup(ctx context.Context, name string, asg *awsautoscaling.AutoScalingGroup) (*awsautoscalingoutputs.AutoScalingGroupOutput, error) {
	m.autoScalingGroup = asg
	return autoScalingGroupToOutput(asg), nil
}

func (m *mockAWSComputeService) DeleteAutoScalingGroup(ctx context.Context, name string) error {
	return nil
}

func (m *mockAWSComputeService) ListAutoScalingGroups(ctx context.Context, filters map[string][]string) ([]*awsautoscalingoutputs.AutoScalingGroupOutput, error) {
	if m.autoScalingGroup != nil {
		return []*awsautoscalingoutputs.AutoScalingGroupOutput{autoScalingGroupToOutput(m.autoScalingGroup)}, nil
	}
	return []*awsautoscalingoutputs.AutoScalingGroupOutput{}, nil
}

func (m *mockAWSComputeService) SetDesiredCapacity(ctx context.Context, asgName string, capacity int) error {
	return nil
}

func (m *mockAWSComputeService) AttachInstances(ctx context.Context, asgName string, instanceIDs []string) error {
	return nil
}

func (m *mockAWSComputeService) DetachInstances(ctx context.Context, asgName string, instanceIDs []string) error {
	return nil
}

func (m *mockAWSComputeService) PutScalingPolicy(ctx context.Context, policy *awsautoscaling.ScalingPolicy) (*awsautoscalingoutputs.ScalingPolicyOutput, error) {
	return &awsautoscalingoutputs.ScalingPolicyOutput{
		PolicyARN:            "arn:aws:autoscaling:us-east-1:123456789012:scalingPolicy:uuid:autoScalingGroupName/" + policy.AutoScalingGroupName + ":policyName/" + policy.PolicyName,
		PolicyName:           policy.PolicyName,
		AutoScalingGroupName: policy.AutoScalingGroupName,
		PolicyType:           policy.PolicyType,
		Alarms:               []awsautoscalingoutputs.Alarm{},
	}, nil
}

func (m *mockAWSComputeService) DescribeScalingPolicies(ctx context.Context, asgName string) ([]*awsautoscalingoutputs.ScalingPolicyOutput, error) {
	return []*awsautoscalingoutputs.ScalingPolicyOutput{}, nil
}

func (m *mockAWSComputeService) DeleteScalingPolicy(ctx context.Context, policyName, asgName string) error {
	return nil
}

// Lambda Function mock methods

func (m *mockAWSComputeService) CreateLambdaFunction(ctx context.Context, function *awslambda.Function) (*awslambdaoutputs.FunctionOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.lambdaFunction = function
	return lambdaFunctionToOutput(function), nil
}

func (m *mockAWSComputeService) GetLambdaFunction(ctx context.Context, name string) (*awslambdaoutputs.FunctionOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	if m.lambdaFunction == nil || m.lambdaFunction.FunctionName != name {
		return nil, errors.New("lambda function not found")
	}
	return lambdaFunctionToOutput(m.lambdaFunction), nil
}

func (m *mockAWSComputeService) UpdateLambdaFunction(ctx context.Context, name string, function *awslambda.Function) (*awslambdaoutputs.FunctionOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.lambdaFunction = function
	return lambdaFunctionToOutput(function), nil
}

func (m *mockAWSComputeService) DeleteLambdaFunction(ctx context.Context, name string) error {
	if m.lambdaFunction == nil || m.lambdaFunction.FunctionName != name {
		return errors.New("lambda function not found")
	}
	m.lambdaFunction = nil
	return nil
}

func (m *mockAWSComputeService) ListLambdaFunctions(ctx context.Context, filters map[string][]string) ([]*awslambdaoutputs.FunctionOutput, error) {
	if m.lambdaFunction == nil {
		return []*awslambdaoutputs.FunctionOutput{}, nil
	}
	return []*awslambdaoutputs.FunctionOutput{lambdaFunctionToOutput(m.lambdaFunction)}, nil
}

// Helper function to convert Lambda Function input to output
func lambdaFunctionToOutput(function *awslambda.Function) *awslambdaoutputs.FunctionOutput {
	if function == nil {
		return nil
	}
	arn := "arn:aws:lambda:us-east-1:123456789012:function:" + function.FunctionName
	invokeARN := "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/" + arn + "/invocations"
	version := "1"
	qualifiedARN := arn + ":" + version
	
	output := &awslambdaoutputs.FunctionOutput{
		ARN:          arn,
		InvokeARN:    invokeARN,
		QualifiedARN: &qualifiedARN,
		Version:      version,
		FunctionName: function.FunctionName,
		Region:       "us-east-1",
		RoleARN:      function.RoleARN,
	}

	if function.S3Bucket != nil {
		output.S3Bucket = function.S3Bucket
		output.S3Key = function.S3Key
		output.S3ObjectVersion = function.S3ObjectVersion
		output.Runtime = function.Runtime
		output.Handler = function.Handler
	}

	if function.PackageType != nil {
		output.PackageType = function.PackageType
		output.ImageURI = function.ImageURI
	}

	if function.MemorySize != nil {
		output.MemorySize = function.MemorySize
	}
	if function.Timeout != nil {
		output.Timeout = function.Timeout
	}
	if function.Environment != nil {
		output.Environment = function.Environment
	}
	if function.Layers != nil {
		output.Layers = function.Layers
	}
	if function.VPCConfig != nil {
		output.VPCConfig = &awslambdaoutputs.FunctionVPCConfigOutput{
			SubnetIDs:        function.VPCConfig.SubnetIDs,
			SecurityGroupIDs: function.VPCConfig.SecurityGroupIDs,
		}
	}

	return output
}

// Auto Scaling Group Tests

func TestAWSComputeAdapter_CreateAutoScalingGroup(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	version := "$Latest"
	domainASG := &domaincompute.AutoScalingGroup{
		Name:             "test-asg",
		Region:           "us-east-1",
		MinSize:          1,
		MaxSize:          5,
		DesiredCapacity:  intPtr(2),
		VPCZoneIdentifier: []string{"subnet-123", "subnet-456"},
		LaunchTemplate: &domaincompute.LaunchTemplateSpecification{
			ID:      "lt-1234567890abcdef0",
			Version: &version,
		},
		HealthCheckType:        domaincompute.AutoScalingGroupHealthCheckTypeEC2,
		HealthCheckGracePeriod: intPtr(300),
	}

	ctx := context.Background()
	createdASG, err := adapter.CreateAutoScalingGroup(ctx, domainASG)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdASG == nil {
		t.Fatal("Expected created auto scaling group, got nil")
	}

	if createdASG.Name != domainASG.Name {
		t.Errorf("Expected name %s, got %s", domainASG.Name, createdASG.Name)
	}

	if createdASG.ID == "" {
		t.Error("Expected auto scaling group ID to be populated")
	}
}

func TestAWSComputeAdapter_GetAutoScalingGroup(t *testing.T) {
	version := "$Latest"
	mockService := &mockAWSComputeService{
		autoScalingGroup: &awsautoscaling.AutoScalingGroup{
			AutoScalingGroupName: stringPtr("test-asg"),
			MinSize:              1,
			MaxSize:              5,
			DesiredCapacity:      intPtr(2),
			VPCZoneIdentifier:    []string{"subnet-123", "subnet-456"},
			LaunchTemplate: &awsautoscaling.LaunchTemplateSpecification{
				LaunchTemplateId: "lt-1234567890abcdef0",
				Version:          &version,
			},
			HealthCheckType: stringPtr("EC2"),
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	asg, err := adapter.GetAutoScalingGroup(ctx, "test-asg")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if asg == nil {
		t.Fatal("Expected auto scaling group, got nil")
	}

	if asg.Name != "test-asg" {
		t.Errorf("Expected name test-asg, got %s", asg.Name)
	}
}

func TestAWSComputeAdapter_SetDesiredCapacity(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	err := adapter.SetDesiredCapacity(ctx, "test-asg", 3)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAWSComputeAdapter_AttachInstances(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	err := adapter.AttachInstances(ctx, "test-asg", []string{"i-123", "i-456"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAWSComputeAdapter_DetachInstances(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	err := adapter.DetachInstances(ctx, "test-asg", []string{"i-123"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// Lambda Function Tests

func TestAWSComputeAdapter_CreateLambdaFunction(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	domainFunction := &domaincompute.LambdaFunction{
		FunctionName: "test-function",
		RoleARN:      "arn:aws:iam::123456789012:role/test-role",
		Region:       "us-east-1",
		S3Bucket:     stringPtr("my-bucket"),
		S3Key:        stringPtr("code.zip"),
		Runtime:      stringPtr("python3.9"),
		Handler:      stringPtr("index.handler"),
		MemorySize:   intPtr(256),
		Timeout:      intPtr(30),
	}

	ctx := context.Background()
	createdFunction, err := adapter.CreateLambdaFunction(ctx, domainFunction)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdFunction == nil {
		t.Fatal("Expected created lambda function, got nil")
	}

	if createdFunction.FunctionName != domainFunction.FunctionName {
		t.Errorf("Expected function name %s, got %s", domainFunction.FunctionName, createdFunction.FunctionName)
	}

	if createdFunction.ARN == nil {
		t.Error("Expected ARN to be populated")
	}

	if createdFunction.InvokeARN == nil {
		t.Error("Expected InvokeARN to be populated")
	}
}

func TestAWSComputeAdapter_CreateLambdaFunction_WithContainerImage(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	domainFunction := &domaincompute.LambdaFunction{
		FunctionName: "test-function",
		RoleARN:      "arn:aws:iam::123456789012:role/test-role",
		Region:       "us-east-1",
		PackageType:  stringPtr("Image"),
		ImageURI:     stringPtr("123456789012.dkr.ecr.us-east-1.amazonaws.com/my-app:latest"),
		MemorySize:   intPtr(512),
	}

	ctx := context.Background()
	createdFunction, err := adapter.CreateLambdaFunction(ctx, domainFunction)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdFunction == nil {
		t.Fatal("Expected created lambda function, got nil")
	}

	if createdFunction.PackageType == nil || *createdFunction.PackageType != "Image" {
		t.Error("Expected package type to be 'Image'")
	}
}

func TestAWSComputeAdapter_CreateLambdaFunction_ValidationError(t *testing.T) {
	mockService := &mockAWSComputeService{}
	adapter := NewAWSComputeAdapter(mockService)

	domainFunction := &domaincompute.LambdaFunction{
		FunctionName: "", // Invalid: empty name
		RoleARN:      "arn:aws:iam::123456789012:role/test-role",
		Region:       "us-east-1",
	}

	ctx := context.Background()
	_, err := adapter.CreateLambdaFunction(ctx, domainFunction)

	if err == nil {
		t.Error("Expected validation error, got nil")
	}
}

func TestAWSComputeAdapter_GetLambdaFunction(t *testing.T) {
	mockService := &mockAWSComputeService{
		lambdaFunction: &awslambda.Function{
			FunctionName: "test-function",
			RoleARN:      "arn:aws:iam::123456789012:role/test-role",
			S3Bucket:     stringPtr("my-bucket"),
			S3Key:        stringPtr("code.zip"),
			Runtime:      stringPtr("python3.9"),
			Handler:      stringPtr("index.handler"),
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	function, err := adapter.GetLambdaFunction(ctx, "test-function")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if function == nil {
		t.Fatal("Expected lambda function, got nil")
	}

	if function.FunctionName != "test-function" {
		t.Errorf("Expected function name test-function, got %s", function.FunctionName)
	}
}

func TestAWSComputeAdapter_UpdateLambdaFunction(t *testing.T) {
	mockService := &mockAWSComputeService{
		lambdaFunction: &awslambda.Function{
			FunctionName: "test-function",
			RoleARN:      "arn:aws:iam::123456789012:role/test-role",
			S3Bucket:     stringPtr("my-bucket"),
			S3Key:        stringPtr("code.zip"),
			Runtime:      stringPtr("python3.9"),
			Handler:      stringPtr("index.handler"),
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	domainFunction := &domaincompute.LambdaFunction{
		FunctionName: "test-function",
		RoleARN:      "arn:aws:iam::123456789012:role/test-role",
		Region:       "us-east-1",
		S3Bucket:     stringPtr("my-bucket"),
		S3Key:        stringPtr("code-v2.zip"),
		Runtime:      stringPtr("python3.9"),
		Handler:      stringPtr("index.handler"),
		MemorySize:   intPtr(512), // Updated memory
	}

	ctx := context.Background()
	updatedFunction, err := adapter.UpdateLambdaFunction(ctx, domainFunction)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if updatedFunction == nil {
		t.Fatal("Expected updated lambda function, got nil")
	}

	if updatedFunction.MemorySize == nil || *updatedFunction.MemorySize != 512 {
		t.Errorf("Expected memory size 512, got %v", updatedFunction.MemorySize)
	}
}

func TestAWSComputeAdapter_DeleteLambdaFunction(t *testing.T) {
	mockService := &mockAWSComputeService{
		lambdaFunction: &awslambda.Function{
			FunctionName: "test-function",
			RoleARN:      "arn:aws:iam::123456789012:role/test-role",
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	err := adapter.DeleteLambdaFunction(ctx, "test-function")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAWSComputeAdapter_ListLambdaFunctions(t *testing.T) {
	mockService := &mockAWSComputeService{
		lambdaFunction: &awslambda.Function{
			FunctionName: "test-function",
			RoleARN:      "arn:aws:iam::123456789012:role/test-role",
		},
	}
	adapter := NewAWSComputeAdapter(mockService)

	ctx := context.Background()
	functions, err := adapter.ListLambdaFunctions(ctx, map[string]string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(functions))
	}
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
