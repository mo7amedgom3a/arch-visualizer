package compute

import (
	"context"
	"fmt"

	awserrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/errors"
	domainerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsautoscaling "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling"
	awsautoscalingoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling/outputs"
	awsec2 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2"
	awslttemplate "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template"
	awslttemplateoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template/outputs"
	awsec2outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/outputs"
	awsmodel "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/instance_types"
	awslambda "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda"
	awslambdaoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda/outputs"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awsloadbalanceroutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services"
)

// ComputeService implements AWSComputeService with deterministic virtual operations
type ComputeService struct{}

// NewComputeService creates a new compute service implementation
func NewComputeService() *ComputeService {
	return &ComputeService{}
}

// Instance operations

func (s *ComputeService) CreateInstance(ctx context.Context, instance *awsec2.Instance) (*awsec2outputs.InstanceOutput, error) {
	if instance == nil {
		return nil, domainerrors.New(awserrors.CodeEC2InstanceCreationFailed, domainerrors.KindValidation, "instance is nil").
			WithOp("ComputeService.CreateInstance")
	}

	instanceID := services.GenerateInstanceID(instance.Name)
	region := "us-east-1"
	arn := services.GenerateARN("ec2", "instance", instanceID, region)

	availabilityZone := "us-east-1a"
	privateIP := "10.0.1.100"

	return &awsec2outputs.InstanceOutput{
		ID:                 instanceID,
		ARN:                arn,
		Name:               instance.Name,
		Region:             region,
		AvailabilityZone:   availabilityZone,
		InstanceType:       instance.InstanceType,
		AMI:                instance.AMI,
		SubnetID:           instance.SubnetID,
		SecurityGroupIDs:   instance.VpcSecurityGroupIds,
		PrivateIP:          privateIP,
		PublicIP:           services.StringPtr("54.123.45.67"),
		PrivateDNS:         fmt.Sprintf("ip-%s.ec2.internal", privateIP),
		PublicDNS:          services.StringPtr(fmt.Sprintf("ec2-%s.compute-1.amazonaws.com", privateIP)),
		VPCID:              "vpc-12345678", // Default VPC ID
		KeyName:            instance.KeyName,
		IAMInstanceProfile: instance.IAMInstanceProfile,
		State:              "running",
		CreationTime:       services.GetFixedTimestamp(),
		Tags:               instance.Tags,
	}, nil
}

func (s *ComputeService) GetInstance(ctx context.Context, id string) (*awsec2outputs.InstanceOutput, error) {
	region := "us-east-1"
	arn := services.GenerateARN("ec2", "instance", id, region)

	return &awsec2outputs.InstanceOutput{
		ID:                 id,
		ARN:                arn,
		Name:               "test-instance",
		Region:             region,
		AvailabilityZone:   "us-east-1a",
		InstanceType:       "t3.micro",
		AMI:                "ami-0c55b159cbfafe1f0",
		SubnetID:           "subnet-123",
		SecurityGroupIDs:   []string{"sg-123"},
		PrivateIP:          "10.0.1.100",
		PublicIP:           services.StringPtr("54.123.45.67"),
		PrivateDNS:         "ip-10-0-1-100.ec2.internal",
		PublicDNS:          services.StringPtr("ec2-54-123-45-67.compute-1.amazonaws.com"),
		VPCID:              "vpc-123",
		KeyName:            services.StringPtr("my-key"),
		IAMInstanceProfile: nil,
		State:              "running",
		CreationTime:       services.GetFixedTimestamp(),
		Tags:               []configs.Tag{},
	}, nil
}

func (s *ComputeService) UpdateInstance(ctx context.Context, instance *awsec2.Instance) (*awsec2outputs.InstanceOutput, error) {
	return s.CreateInstance(ctx, instance)
}

func (s *ComputeService) DeleteInstance(ctx context.Context, id string) error {
	return nil
}

func (s *ComputeService) ListInstances(ctx context.Context, filters map[string][]string) ([]*awsec2outputs.InstanceOutput, error) {
	return []*awsec2outputs.InstanceOutput{
		{
			ID:                 "i-0a1b2c3d4e5f6g7h8",
			ARN:                "arn:aws:ec2:us-east-1:123456789012:instance/i-0a1b2c3d4e5f6g7h8",
			Name:               "test-instance",
			Region:             "us-east-1",
			AvailabilityZone:   "us-east-1a",
			InstanceType:       "t3.micro",
			AMI:                "ami-0c55b159cbfafe1f0",
			SubnetID:           "subnet-123",
			SecurityGroupIDs:   []string{"sg-123"},
			PrivateIP:          "10.0.1.100",
			PublicIP:           services.StringPtr("54.123.45.67"),
			PrivateDNS:         "ip-10-0-1-100.ec2.internal",
			PublicDNS:          services.StringPtr("ec2-54-123-45-67.compute-1.amazonaws.com"),
			VPCID:              "vpc-123",
			KeyName:            services.StringPtr("my-key"),
			IAMInstanceProfile: nil,
			State:              "running",
			CreationTime:       services.GetFixedTimestamp(),
			Tags:               []configs.Tag{},
		},
	}, nil
}

// Instance lifecycle operations

func (s *ComputeService) StartInstance(ctx context.Context, id string) error {
	return nil
}

func (s *ComputeService) StopInstance(ctx context.Context, id string) error {
	return nil
}

func (s *ComputeService) RebootInstance(ctx context.Context, id string) error {
	return nil
}

// Instance type operations

func (s *ComputeService) GetInstanceTypeInfo(ctx context.Context, instanceType string, region string) (*awsmodel.InstanceTypeInfo, error) {
	return &awsmodel.InstanceTypeInfo{
		Name:           instanceType,
		Category:       awsmodel.CategoryGeneralPurpose,
		VCPU:           2,
		MemoryGiB:      4.0,
		MaxNetworkGbps: 5,
		StorageType:    "EBS only",
	}, nil
}

func (s *ComputeService) ListInstanceTypesByCategory(ctx context.Context, category awsmodel.InstanceCategory, region string) ([]*awsmodel.InstanceTypeInfo, error) {
	return []*awsmodel.InstanceTypeInfo{
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

// Launch Template operations

func (s *ComputeService) CreateLaunchTemplate(ctx context.Context, template *awslttemplate.LaunchTemplate) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	if template == nil {
		return nil, domainerrors.New(awserrors.CodeEC2InstanceCreationFailed, domainerrors.KindValidation, "template is nil").
			WithOp("ComputeService.CreateLaunchTemplate")
	}

	name := "test-template"
	if template.Name != nil {
		name = *template.Name
	} else if template.NamePrefix != nil {
		name = *template.NamePrefix + "-12345"
	}

	templateID := fmt.Sprintf("lt-%s", services.GenerateDeterministicID(name)[:15])
	region := "us-east-1"
	arn := services.GenerateARN("ec2", "launch-template", templateID, region)

	return &awslttemplateoutputs.LaunchTemplateOutput{
		ID:             templateID,
		ARN:            arn,
		Name:           name,
		DefaultVersion: 1,
		LatestVersion:  1,
		CreateTime:     services.GetFixedTimestamp(),
		CreatedBy:      "arn:aws:iam::123456789012:user/test",
		Tags: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			{Key: "Name", Value: name},
		},
	}, nil
}

func (s *ComputeService) GetLaunchTemplate(ctx context.Context, id string) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	region := "us-east-1"
	arn := services.GenerateARN("ec2", "launch-template", id, region)

	return &awslttemplateoutputs.LaunchTemplateOutput{
		ID:             id,
		ARN:            arn,
		Name:           "test-template",
		DefaultVersion: 1,
		LatestVersion:  1,
		CreateTime:     services.GetFixedTimestamp(),
		CreatedBy:      "arn:aws:iam::123456789012:user/test",
		Tags: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			{Key: "Name", Value: "test-template"},
		},
	}, nil
}

func (s *ComputeService) UpdateLaunchTemplate(ctx context.Context, id string, template *awslttemplate.LaunchTemplate) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	return s.CreateLaunchTemplate(ctx, template)
}

func (s *ComputeService) DeleteLaunchTemplate(ctx context.Context, id string) error {
	return nil
}

func (s *ComputeService) ListLaunchTemplates(ctx context.Context, filters map[string][]string) ([]*awslttemplateoutputs.LaunchTemplateOutput, error) {
	return []*awslttemplateoutputs.LaunchTemplateOutput{
		{
			ID:             "lt-0a1b2c3d4e5f6g7h8",
			ARN:            "arn:aws:ec2:us-east-1:123456789012:launch-template/lt-0a1b2c3d4e5f6g7h8",
			Name:           "test-template",
			DefaultVersion: 1,
			LatestVersion:  1,
			CreateTime:     services.GetFixedTimestamp(),
			CreatedBy:      "arn:aws:iam::123456789012:user/test",
			Tags: []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}{
				{Key: "Name", Value: "test-template"},
			},
		},
	}, nil
}

func (s *ComputeService) GetLaunchTemplateVersion(ctx context.Context, id string, version int) (*awslttemplate.LaunchTemplateVersion, error) {
	templateData := &awslttemplate.LaunchTemplate{
		Name:                services.StringPtr("test-template"),
		ImageID:             "ami-0c55b159cbfafe1f0",
		InstanceType:        "t3.micro",
		VpcSecurityGroupIds: []string{"sg-123"},
	}
	return &awslttemplate.LaunchTemplateVersion{
		TemplateID:    id,
		VersionNumber: version,
		IsDefault:     version == 1,
		CreateTime:    services.GetFixedTimestamp(),
		CreatedBy:     services.StringPtr("arn:aws:iam::123456789012:user/test"),
		TemplateData:  templateData,
	}, nil
}

func (s *ComputeService) ListLaunchTemplateVersions(ctx context.Context, id string) ([]*awslttemplate.LaunchTemplateVersion, error) {
	templateData := &awslttemplate.LaunchTemplate{
		Name:                services.StringPtr("test-template"),
		ImageID:             "ami-0c55b159cbfafe1f0",
		InstanceType:        "t3.micro",
		VpcSecurityGroupIds: []string{"sg-123"},
	}
	return []*awslttemplate.LaunchTemplateVersion{
		{
			TemplateID:    id,
			VersionNumber: 1,
			IsDefault:     true,
			CreateTime:    services.GetFixedTimestamp(),
			CreatedBy:     services.StringPtr("arn:aws:iam::123456789012:user/test"),
			TemplateData:  templateData,
		},
	}, nil
}

// Load Balancer operations

func (s *ComputeService) CreateLoadBalancer(ctx context.Context, lb *awsloadbalancer.LoadBalancer) (*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	if lb == nil {
		return nil, domainerrors.New(awserrors.CodeLoadBalancerCreationFailed, domainerrors.KindValidation, "load balancer is nil").
			WithOp("ComputeService.CreateLoadBalancer")
	}

	internal := false
	if lb.Internal != nil {
		internal = *lb.Internal
	}

	lbID := fmt.Sprintf("app/%s/%s", lb.Name, services.GenerateDeterministicID(lb.Name)[:16])
	region := "us-east-1"
	arn := fmt.Sprintf("arn:aws:elasticloadbalancing:%s:123456789012:loadbalancer/%s", region, lbID)

	return &awsloadbalanceroutputs.LoadBalancerOutput{
		ARN:              arn,
		ID:               arn,
		Name:             lb.Name,
		DNSName:          fmt.Sprintf("%s-%s.%s.elb.amazonaws.com", lb.Name, services.GenerateDeterministicID(lb.Name)[:8], region),
		ZoneID:           "Z35SXDOTRQ7X7K",
		Type:             lb.LoadBalancerType,
		Internal:         internal,
		SecurityGroupIDs: lb.SecurityGroupIDs,
		SubnetIDs:        lb.SubnetIDs,
		State:            "active",
		CreatedTime:      services.GetFixedTimestamp(),
	}, nil
}

func (s *ComputeService) GetLoadBalancer(ctx context.Context, arn string) (*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	return &awsloadbalanceroutputs.LoadBalancerOutput{
		ARN:              arn,
		ID:               arn,
		Name:             "test-lb",
		DNSName:          "test-lb-1234567890.us-east-1.elb.amazonaws.com",
		ZoneID:           "Z35SXDOTRQ7X7K",
		Type:             "application",
		Internal:         false,
		SecurityGroupIDs: []string{"sg-123"},
		SubnetIDs:        []string{"subnet-123", "subnet-456"},
		State:            "active",
		CreatedTime:      services.GetFixedTimestamp(),
	}, nil
}

func (s *ComputeService) UpdateLoadBalancer(ctx context.Context, arn string, lb *awsloadbalancer.LoadBalancer) (*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	return s.CreateLoadBalancer(ctx, lb)
}

func (s *ComputeService) DeleteLoadBalancer(ctx context.Context, arn string) error {
	return nil
}

func (s *ComputeService) ListLoadBalancers(ctx context.Context, filters map[string][]string) ([]*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	return []*awsloadbalanceroutputs.LoadBalancerOutput{
		{
			ARN:              "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef",
			ID:               "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef",
			Name:             "test-lb",
			DNSName:          "test-lb-1234567890.us-east-1.elb.amazonaws.com",
			ZoneID:           "Z35SXDOTRQ7X7K",
			Type:             "application",
			Internal:         false,
			SecurityGroupIDs: []string{"sg-123"},
			SubnetIDs:        []string{"subnet-123", "subnet-456"},
			State:            "active",
			CreatedTime:      services.GetFixedTimestamp(),
		},
	}, nil
}

// Target Group operations

func (s *ComputeService) CreateTargetGroup(ctx context.Context, tg *awsloadbalancer.TargetGroup) (*awsloadbalanceroutputs.TargetGroupOutput, error) {
	if tg == nil {
		return nil, domainerrors.New(awserrors.CodeLoadBalancerCreationFailed, domainerrors.KindValidation, "target group is nil").
			WithOp("ComputeService.CreateTargetGroup")
	}

	targetType := "instance"
	if tg.TargetType != nil {
		targetType = *tg.TargetType
	}

	tgID := fmt.Sprintf("targetgroup/%s/%s", tg.Name, services.GenerateDeterministicID(tg.Name)[:16])
	region := "us-east-1"
	arn := fmt.Sprintf("arn:aws:elasticloadbalancing:%s:123456789012:%s", region, tgID)

	return &awsloadbalanceroutputs.TargetGroupOutput{
		ARN:         arn,
		ID:          arn,
		Name:        tg.Name,
		Port:        tg.Port,
		Protocol:    tg.Protocol,
		VPCID:       tg.VPCID,
		TargetType:  targetType,
		HealthCheck: tg.HealthCheck,
		State:       "active",
		CreatedTime: services.GetFixedTimestamp(),
	}, nil
}

func (s *ComputeService) GetTargetGroup(ctx context.Context, arn string) (*awsloadbalanceroutputs.TargetGroupOutput, error) {
	return &awsloadbalanceroutputs.TargetGroupOutput{
		ARN:        arn,
		ID:         arn,
		Name:       "test-tg",
		Port:       80,
		Protocol:   "HTTP",
		VPCID:      "vpc-123",
		TargetType: "instance",
		HealthCheck: awsloadbalancer.HealthCheckConfig{
			Path:    services.StringPtr("/health"),
			Matcher: services.StringPtr("200"),
		},
		State:       "active",
		CreatedTime: services.GetFixedTimestamp(),
	}, nil
}

func (s *ComputeService) UpdateTargetGroup(ctx context.Context, arn string, tg *awsloadbalancer.TargetGroup) (*awsloadbalanceroutputs.TargetGroupOutput, error) {
	return s.CreateTargetGroup(ctx, tg)
}

func (s *ComputeService) DeleteTargetGroup(ctx context.Context, arn string) error {
	return nil
}

func (s *ComputeService) ListTargetGroups(ctx context.Context, filters map[string][]string) ([]*awsloadbalanceroutputs.TargetGroupOutput, error) {
	return []*awsloadbalanceroutputs.TargetGroupOutput{
		{
			ARN:        "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef",
			ID:         "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef",
			Name:       "test-tg",
			Port:       80,
			Protocol:   "HTTP",
			VPCID:      "vpc-123",
			TargetType: "instance",
			HealthCheck: awsloadbalancer.HealthCheckConfig{
				Path:    services.StringPtr("/health"),
				Matcher: services.StringPtr("200"),
			},
			State:       "active",
			CreatedTime: services.GetFixedTimestamp(),
		},
	}, nil
}

// Listener operations

func (s *ComputeService) CreateListener(ctx context.Context, listener *awsloadbalancer.Listener) (*awsloadbalanceroutputs.ListenerOutput, error) {
	if listener == nil {
		return nil, domainerrors.New(awserrors.CodeLoadBalancerCreationFailed, domainerrors.KindValidation, "listener is nil").
			WithOp("ComputeService.CreateListener")
	}

	listenerID := fmt.Sprintf("%s/%s", listener.LoadBalancerARN, services.GenerateDeterministicID(listener.LoadBalancerARN)[:16])
	region := "us-east-1"
	arn := fmt.Sprintf("arn:aws:elasticloadbalancing:%s:123456789012:listener/%s", region, listenerID)

	return &awsloadbalanceroutputs.ListenerOutput{
		ARN:             arn,
		ID:              arn,
		LoadBalancerARN: listener.LoadBalancerARN,
		Port:            listener.Port,
		Protocol:        listener.Protocol,
		DefaultAction:   listener.DefaultAction,
	}, nil
}

func (s *ComputeService) GetListener(ctx context.Context, arn string) (*awsloadbalanceroutputs.ListenerOutput, error) {
	targetGroupARN := "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef"
	return &awsloadbalanceroutputs.ListenerOutput{
		ARN:             arn,
		ID:              arn,
		LoadBalancerARN: "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/1234567890abcdef",
		Port:            80,
		Protocol:        "HTTP",
		DefaultAction: awsloadbalancer.ListenerAction{
			Type:           awsloadbalancer.ListenerActionTypeForward,
			TargetGroupARN: &targetGroupARN,
		},
	}, nil
}

func (s *ComputeService) UpdateListener(ctx context.Context, arn string, listener *awsloadbalancer.Listener) (*awsloadbalanceroutputs.ListenerOutput, error) {
	return s.CreateListener(ctx, listener)
}

func (s *ComputeService) DeleteListener(ctx context.Context, arn string) error {
	return nil
}

func (s *ComputeService) ListListeners(ctx context.Context, loadBalancerARN string) ([]*awsloadbalanceroutputs.ListenerOutput, error) {
	targetGroupARN := "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef"
	return []*awsloadbalanceroutputs.ListenerOutput{
		{
			ARN:             "arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/test-lb/1234567890abcdef/1234567890abcdef",
			ID:              "arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/test-lb/1234567890abcdef/1234567890abcdef",
			LoadBalancerARN: loadBalancerARN,
			Port:            80,
			Protocol:        "HTTP",
			DefaultAction: awsloadbalancer.ListenerAction{
				Type:           awsloadbalancer.ListenerActionTypeForward,
				TargetGroupARN: &targetGroupARN,
			},
		},
	}, nil
}

// Target Group Attachment operations

func (s *ComputeService) AttachTargetToGroup(ctx context.Context, attachment *awsloadbalancer.TargetGroupAttachment) error {
	return nil
}

func (s *ComputeService) DetachTargetFromGroup(ctx context.Context, targetGroupARN, targetID string) error {
	return nil
}

func (s *ComputeService) ListTargetGroupTargets(ctx context.Context, targetGroupARN string) ([]*awsloadbalanceroutputs.TargetGroupAttachmentOutput, error) {
	return []*awsloadbalanceroutputs.TargetGroupAttachmentOutput{
		{
			TargetGroupARN:   targetGroupARN,
			TargetID:         "i-1234567890abcdef0",
			Port:             services.IntPtr(8080),
			AvailabilityZone: services.StringPtr("us-east-1a"),
			HealthStatus:     "healthy",
			State:            "healthy",
		},
	}, nil
}

// Auto Scaling Group operations

func (s *ComputeService) CreateAutoScalingGroup(ctx context.Context, asg *awsautoscaling.AutoScalingGroup) (*awsautoscalingoutputs.AutoScalingGroupOutput, error) {
	if asg == nil {
		return nil, domainerrors.New(awserrors.CodeAutoScalingGroupCreationFailed, domainerrors.KindValidation, "auto scaling group is nil").
			WithOp("ComputeService.CreateAutoScalingGroup")
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

	asgARN := fmt.Sprintf("arn:aws:autoscaling:us-east-1:123456789012:autoScalingGroup:uuid:autoScalingGroupName/%s", name)

	return &awsautoscalingoutputs.AutoScalingGroupOutput{
		AutoScalingGroupARN:    asgARN,
		AutoScalingGroupName:   name,
		MinSize:                asg.MinSize,
		MaxSize:                asg.MaxSize,
		DesiredCapacity:        desiredCapacity,
		VPCZoneIdentifier:      asg.VPCZoneIdentifier,
		LaunchTemplate:         asg.LaunchTemplate,
		HealthCheckType:        healthCheckType,
		HealthCheckGracePeriod: asg.HealthCheckGracePeriod,
		TargetGroupARNs:        asg.TargetGroupARNs,
		Status:                 "active",
		CreatedTime:            services.GetFixedTimestamp(),
		Instances:              []awsautoscalingoutputs.Instance{},
		Tags:                   asg.Tags,
	}, nil
}

func (s *ComputeService) GetAutoScalingGroup(ctx context.Context, name string) (*awsautoscalingoutputs.AutoScalingGroupOutput, error) {
	version := "$Latest"
	asgARN := fmt.Sprintf("arn:aws:autoscaling:us-east-1:123456789012:autoScalingGroup:uuid:autoScalingGroupName/%s", name)

	return &awsautoscalingoutputs.AutoScalingGroupOutput{
		AutoScalingGroupARN:  asgARN,
		AutoScalingGroupName: name,
		MinSize:              1,
		MaxSize:              5,
		DesiredCapacity:      2,
		VPCZoneIdentifier:    []string{"subnet-123", "subnet-456"},
		LaunchTemplate: &awsautoscaling.LaunchTemplateSpecification{
			LaunchTemplateId: "lt-1234567890abcdef0",
			Version:          &version,
		},
		HealthCheckType:        "EC2",
		HealthCheckGracePeriod: services.IntPtr(300),
		TargetGroupARNs:        []string{},
		Status:                 "active",
		CreatedTime:            services.GetFixedTimestamp(),
		Instances:              []awsautoscalingoutputs.Instance{},
		Tags:                   []awsautoscaling.Tag{},
	}, nil
}

func (s *ComputeService) UpdateAutoScalingGroup(ctx context.Context, name string, asg *awsautoscaling.AutoScalingGroup) (*awsautoscalingoutputs.AutoScalingGroupOutput, error) {
	return s.CreateAutoScalingGroup(ctx, asg)
}

func (s *ComputeService) DeleteAutoScalingGroup(ctx context.Context, name string) error {
	return nil
}

func (s *ComputeService) ListAutoScalingGroups(ctx context.Context, filters map[string][]string) ([]*awsautoscalingoutputs.AutoScalingGroupOutput, error) {
	version := "$Latest"
	return []*awsautoscalingoutputs.AutoScalingGroupOutput{
		{
			AutoScalingGroupARN:  "arn:aws:autoscaling:us-east-1:123456789012:autoScalingGroup:uuid:autoScalingGroupName/test-asg",
			AutoScalingGroupName: "test-asg",
			MinSize:              1,
			MaxSize:              5,
			DesiredCapacity:      2,
			VPCZoneIdentifier:    []string{"subnet-123", "subnet-456"},
			LaunchTemplate: &awsautoscaling.LaunchTemplateSpecification{
				LaunchTemplateId: "lt-1234567890abcdef0",
				Version:          &version,
			},
			HealthCheckType:        "EC2",
			HealthCheckGracePeriod: services.IntPtr(300),
			TargetGroupARNs:        []string{},
			Status:                 "active",
			CreatedTime:            services.GetFixedTimestamp(),
			Instances:              []awsautoscalingoutputs.Instance{},
			Tags:                   []awsautoscaling.Tag{},
		},
	}, nil
}

// Scaling operations

func (s *ComputeService) SetDesiredCapacity(ctx context.Context, asgName string, capacity int) error {
	return nil
}

func (s *ComputeService) AttachInstances(ctx context.Context, asgName string, instanceIDs []string) error {
	return nil
}

func (s *ComputeService) DetachInstances(ctx context.Context, asgName string, instanceIDs []string) error {
	return nil
}

// Scaling Policy operations

func (s *ComputeService) PutScalingPolicy(ctx context.Context, policy *awsautoscaling.ScalingPolicy) (*awsautoscalingoutputs.ScalingPolicyOutput, error) {
	if policy == nil {
		return nil, domainerrors.New(awserrors.CodeAutoScalingGroupCreationFailed, domainerrors.KindValidation, "policy is nil").
			WithOp("ComputeService.CreateScalingPolicy")
	}

	policyARN := fmt.Sprintf("arn:aws:autoscaling:us-east-1:123456789012:scalingPolicy:uuid:autoScalingGroupName/%s:policyName/%s", policy.AutoScalingGroupName, policy.PolicyName)

	return &awsautoscalingoutputs.ScalingPolicyOutput{
		PolicyARN:            policyARN,
		PolicyName:           policy.PolicyName,
		AutoScalingGroupName: policy.AutoScalingGroupName,
		PolicyType:           policy.PolicyType,
		Alarms:               []awsautoscalingoutputs.Alarm{},
	}, nil
}

func (s *ComputeService) DescribeScalingPolicies(ctx context.Context, asgName string) ([]*awsautoscalingoutputs.ScalingPolicyOutput, error) {
	return []*awsautoscalingoutputs.ScalingPolicyOutput{}, nil
}

func (s *ComputeService) DeleteScalingPolicy(ctx context.Context, policyName, asgName string) error {
	return nil
}

// Lambda Function operations

func (s *ComputeService) CreateLambdaFunction(ctx context.Context, function *awslambda.Function) (*awslambdaoutputs.FunctionOutput, error) {
	if function == nil {
		return nil, domainerrors.New(awserrors.CodeLambdaFunctionCreationFailed, domainerrors.KindValidation, "function is nil").
			WithOp("ComputeService.CreateLambdaFunction")
	}

	arn := fmt.Sprintf("arn:aws:lambda:us-east-1:123456789012:function:%s", function.FunctionName)
	invokeARN := fmt.Sprintf("arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/%s/invocations", arn)
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
		MemorySize:   function.MemorySize,
		Timeout:      function.Timeout,
		Runtime:      function.Runtime,
		Handler:      function.Handler,
		Environment:  function.Environment,
		Layers:       function.Layers,
	}

	if function.S3Bucket != nil {
		output.S3Bucket = function.S3Bucket
		output.S3Key = function.S3Key
		output.S3ObjectVersion = function.S3ObjectVersion
	}

	if function.PackageType != nil {
		output.PackageType = function.PackageType
		output.ImageURI = function.ImageURI
	}

	if function.VPCConfig != nil {
		output.VPCConfig = &awslambdaoutputs.FunctionVPCConfigOutput{
			SubnetIDs:        function.VPCConfig.SubnetIDs,
			SecurityGroupIDs: function.VPCConfig.SecurityGroupIDs,
			VPCID:            services.StringPtr("vpc-12345678"),
		}
	}

	lastModified := services.GetFixedTimestamp().Format("2006-01-02T15:04:05.000Z")
	output.LastModified = &lastModified
	codeSize := int64(1024)
	output.CodeSize = &codeSize
	codeSHA256 := "abc123def456"
	output.CodeSHA256 = &codeSHA256

	return output, nil
}

func (s *ComputeService) GetLambdaFunction(ctx context.Context, name string) (*awslambdaoutputs.FunctionOutput, error) {
	arn := fmt.Sprintf("arn:aws:lambda:us-east-1:123456789012:function:%s", name)
	invokeARN := fmt.Sprintf("arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/%s/invocations", arn)
	version := "1"
	qualifiedARN := arn + ":" + version

	return &awslambdaoutputs.FunctionOutput{
		ARN:          arn,
		InvokeARN:    invokeARN,
		QualifiedARN: &qualifiedARN,
		Version:      version,
		FunctionName: name,
		Region:       "us-east-1",
		RoleARN:      "arn:aws:iam::123456789012:role/test-role",
		Runtime:      services.StringPtr("python3.9"),
		Handler:      services.StringPtr("index.handler"),
		MemorySize:   services.Int32Ptr(256),
		Timeout:      services.Int32Ptr(30),
	}, nil
}

func (s *ComputeService) UpdateLambdaFunction(ctx context.Context, name string, function *awslambda.Function) (*awslambdaoutputs.FunctionOutput, error) {
	return s.CreateLambdaFunction(ctx, function)
}

func (s *ComputeService) DeleteLambdaFunction(ctx context.Context, name string) error {
	return nil
}

func (s *ComputeService) ListLambdaFunctions(ctx context.Context, filters map[string][]string) ([]*awslambdaoutputs.FunctionOutput, error) {
	return []*awslambdaoutputs.FunctionOutput{
		{
			ARN:          "arn:aws:lambda:us-east-1:123456789012:function:test-function-1",
			InvokeARN:    "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789012:function:test-function-1/invocations",
			FunctionName: "test-function-1",
			Region:       "us-east-1",
			RoleARN:      "arn:aws:iam::123456789012:role/test-role",
			Version:      "1",
		},
		{
			ARN:          "arn:aws:lambda:us-east-1:123456789012:function:test-function-2",
			InvokeARN:    "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789012:function:test-function-2/invocations",
			FunctionName: "test-function-2",
			Region:       "us-east-1",
			RoleARN:      "arn:aws:iam::123456789012:role/test-role",
			Version:      "1",
		},
	}, nil
}
