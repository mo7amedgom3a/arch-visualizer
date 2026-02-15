package common

import (
	"context"
	"fmt"

	awscomputemapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/compute"
	awsiammapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/iam"
	awsnetworkingmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/networking"
	awsstoragemapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/storage"
	awscomputeoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling/outputs"
	awslttemplateoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template/outputs"
	awsec2outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/outputs"
	awslambdaoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda/outputs"
	awsloadbalanceroutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer/outputs"
	awsiamoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
	awsnetworkingoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	awss3outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3/outputs"
	awscomputeservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	awsiamservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/iam"
	awsnetworkingservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/networking"
	awsstorageservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/storage"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	domainiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/iam"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
)

// Networking helpers using AWS output models

func CreateVPCWithOutput(ctx context.Context, service awsnetworkingservice.AWSNetworkingService, vpc *domainnetworking.VPC) (*domainnetworking.VPC, *awsnetworkingoutputs.VPCOutput, error) {
	if vpc == nil {
		return nil, nil, fmt.Errorf("vpc is nil")
	}
	if err := vpc.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsVPC := awsnetworkingmapper.FromDomainVPC(vpc)
	if err := awsVPC.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateVPC(ctx, awsVPC)
	if err != nil {
		return nil, nil, err
	}
	return awsnetworkingmapper.ToDomainVPCFromOutput(output), output, nil
}

func CreateSubnetWithOutput(ctx context.Context, service awsnetworkingservice.AWSNetworkingService, subnet *domainnetworking.Subnet) (*domainnetworking.Subnet, *awsnetworkingoutputs.SubnetOutput, error) {
	if subnet == nil {
		return nil, nil, fmt.Errorf("subnet is nil")
	}
	if err := subnet.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	if subnet.AvailabilityZone == nil || *subnet.AvailabilityZone == "" {
		return nil, nil, fmt.Errorf("availability zone is required for subnet")
	}
	awsSubnet := awsnetworkingmapper.FromDomainSubnet(subnet, *subnet.AvailabilityZone)
	if err := awsSubnet.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateSubnet(ctx, awsSubnet)
	if err != nil {
		return nil, nil, err
	}
	return awsnetworkingmapper.ToDomainSubnetFromOutput(output), output, nil
}

func CreateInternetGatewayWithOutput(ctx context.Context, service awsnetworkingservice.AWSNetworkingService, igw *domainnetworking.InternetGateway) (*domainnetworking.InternetGateway, *awsnetworkingoutputs.InternetGatewayOutput, error) {
	if igw == nil {
		return nil, nil, fmt.Errorf("internet gateway is nil")
	}
	if err := igw.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsIGW := awsnetworkingmapper.FromDomainInternetGateway(igw)
	if err := awsIGW.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateInternetGateway(ctx, awsIGW)
	if err != nil {
		return nil, nil, err
	}
	return awsnetworkingmapper.ToDomainInternetGatewayFromOutput(output), output, nil
}

func CreateRouteTableWithOutput(ctx context.Context, service awsnetworkingservice.AWSNetworkingService, rt *domainnetworking.RouteTable) (*domainnetworking.RouteTable, *awsnetworkingoutputs.RouteTableOutput, error) {
	if rt == nil {
		return nil, nil, fmt.Errorf("route table is nil")
	}
	if err := rt.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsRT := awsnetworkingmapper.FromDomainRouteTable(rt)
	if err := awsRT.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateRouteTable(ctx, awsRT)
	if err != nil {
		return nil, nil, err
	}
	return awsnetworkingmapper.ToDomainRouteTableFromOutput(output), output, nil
}

func CreateSecurityGroupWithOutput(ctx context.Context, service awsnetworkingservice.AWSNetworkingService, sg *domainnetworking.SecurityGroup) (*domainnetworking.SecurityGroup, *awsnetworkingoutputs.SecurityGroupOutput, error) {
	if sg == nil {
		return nil, nil, fmt.Errorf("security group is nil")
	}
	if err := sg.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsSG := awsnetworkingmapper.FromDomainSecurityGroup(sg)
	if err := awsSG.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateSecurityGroup(ctx, awsSG)
	if err != nil {
		return nil, nil, err
	}
	return awsnetworkingmapper.ToDomainSecurityGroupFromOutput(output), output, nil
}

func AllocateElasticIPWithOutput(ctx context.Context, service awsnetworkingservice.AWSNetworkingService, eip *domainnetworking.ElasticIP) (*domainnetworking.ElasticIP, *awsnetworkingoutputs.ElasticIPOutput, error) {
	if eip == nil {
		return nil, nil, fmt.Errorf("elastic ip is nil")
	}
	if err := eip.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsEIP := awsnetworkingmapper.FromDomainElasticIP(eip)
	if err := awsEIP.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.AllocateElasticIP(ctx, awsEIP)
	if err != nil {
		return nil, nil, err
	}
	return awsnetworkingmapper.ToDomainElasticIPFromOutput(output), output, nil
}

func CreateNATGatewayWithOutput(ctx context.Context, service awsnetworkingservice.AWSNetworkingService, ngw *domainnetworking.NATGateway) (*domainnetworking.NATGateway, *awsnetworkingoutputs.NATGatewayOutput, error) {
	if ngw == nil {
		return nil, nil, fmt.Errorf("nat gateway is nil")
	}
	if err := ngw.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsNAT := awsnetworkingmapper.FromDomainNATGateway(ngw)
	if err := awsNAT.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateNATGateway(ctx, awsNAT)
	if err != nil {
		return nil, nil, err
	}
	return awsnetworkingmapper.ToDomainNATGatewayFromOutput(output), output, nil
}

func AttachInternetGateway(ctx context.Context, service awsnetworkingservice.AWSNetworkingService, igwID, vpcID string) error {
	return service.AttachInternetGateway(ctx, igwID, vpcID)
}

func AssociateRouteTable(ctx context.Context, service awsnetworkingservice.AWSNetworkingService, rtID, subnetID string) error {
	return service.AssociateRouteTable(ctx, rtID, subnetID)
}

func CreateVPCEndpointWithOutput(ctx context.Context, service awsnetworkingservice.AWSNetworkingService, vpce *domainnetworking.VPCEndpoint) (*domainnetworking.VPCEndpoint, *awsnetworkingoutputs.VPCOutput, error) {
	if vpce == nil {
		return nil, nil, fmt.Errorf("vpc endpoint is nil")
	}
	if err := vpce.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsVPCE := awsnetworkingmapper.FromDomainVPCEndpoint(vpce)
	if err := awsVPCE.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	// Note: Currently returns VPCOutput as placeholder
	output, err := service.CreateVPCEndpoint(ctx, awsVPCE)
	if err != nil {
		return nil, nil, err
	}
	return awsnetworkingmapper.ToDomainVPCEndpointFromOutput(output, vpce), output, nil
}

// Compute helpers using AWS output models

func CreateInstanceWithOutput(ctx context.Context, service awscomputeservice.AWSComputeService, instance *domaincompute.Instance) (*domaincompute.Instance, *awsec2outputs.InstanceOutput, error) {
	if instance == nil {
		return nil, nil, fmt.Errorf("instance is nil")
	}
	if err := instance.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsInstance := awscomputemapper.FromDomainInstance(instance)
	if err := awsInstance.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateInstance(ctx, awsInstance)
	if err != nil {
		return nil, nil, err
	}
	return awscomputemapper.ToDomainInstanceFromOutput(output), output, nil
}

func CreateLoadBalancerWithOutput(ctx context.Context, service awscomputeservice.AWSComputeService, lb *domaincompute.LoadBalancer) (*domaincompute.LoadBalancer, *awsloadbalanceroutputs.LoadBalancerOutput, error) {
	if lb == nil {
		return nil, nil, fmt.Errorf("load balancer is nil")
	}
	if err := lb.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsLB := awscomputemapper.FromDomainLoadBalancer(lb)
	if err := awsLB.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateLoadBalancer(ctx, awsLB)
	if err != nil {
		return nil, nil, err
	}
	return awscomputemapper.ToDomainLoadBalancerFromOutput(output), output, nil
}

func CreateTargetGroupWithOutput(ctx context.Context, service awscomputeservice.AWSComputeService, tg *domaincompute.TargetGroup) (*domaincompute.TargetGroup, *awsloadbalanceroutputs.TargetGroupOutput, error) {
	if tg == nil {
		return nil, nil, fmt.Errorf("target group is nil")
	}
	if err := tg.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsTG := awscomputemapper.FromDomainTargetGroup(tg)
	if err := awsTG.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateTargetGroup(ctx, awsTG)
	if err != nil {
		return nil, nil, err
	}
	return awscomputemapper.ToDomainTargetGroupFromOutput(output), output, nil
}

func AttachTargetToGroup(ctx context.Context, service awscomputeservice.AWSComputeService, attachment *domaincompute.TargetGroupAttachment) error {
	if attachment == nil {
		return fmt.Errorf("target group attachment is nil")
	}
	if err := attachment.Validate(); err != nil {
		return fmt.Errorf("domain validation failed: %w", err)
	}
	awsAttachment := awscomputemapper.FromDomainTargetGroupAttachment(attachment)
	if err := awsAttachment.Validate(); err != nil {
		return fmt.Errorf("aws validation failed: %w", err)
	}
	return service.AttachTargetToGroup(ctx, awsAttachment)
}

func CreateLaunchTemplateWithOutput(ctx context.Context, service awscomputeservice.AWSComputeService, template *domaincompute.LaunchTemplate) (*domaincompute.LaunchTemplate, *awslttemplateoutputs.LaunchTemplateOutput, error) {
	if template == nil {
		return nil, nil, fmt.Errorf("launch template is nil")
	}
	if err := template.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsTemplate := awscomputemapper.FromDomainLaunchTemplate(template)
	if err := awsTemplate.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateLaunchTemplate(ctx, awsTemplate)
	if err != nil {
		return nil, nil, err
	}
	domainTemplate := awscomputemapper.ToDomainLaunchTemplateFromOutput(output)
	domainTemplate.Region = template.Region
	return domainTemplate, output, nil
}

func CreateAutoScalingGroupWithOutput(ctx context.Context, service awscomputeservice.AWSComputeService, asg *domaincompute.AutoScalingGroup) (*domaincompute.AutoScalingGroup, *awscomputeoutputs.AutoScalingGroupOutput, error) {
	if asg == nil {
		return nil, nil, fmt.Errorf("auto scaling group is nil")
	}
	if err := asg.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsASG := awscomputemapper.FromDomainAutoScalingGroup(asg)
	if err := awsASG.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateAutoScalingGroup(ctx, awsASG)
	if err != nil {
		return nil, nil, err
	}
	return awscomputemapper.ToDomainAutoScalingGroupFromOutput(output), output, nil
}

func CreateLambdaFunctionWithOutput(ctx context.Context, service awscomputeservice.AWSComputeService, function *domaincompute.LambdaFunction) (*domaincompute.LambdaFunction, *awslambdaoutputs.FunctionOutput, error) {
	if function == nil {
		return nil, nil, fmt.Errorf("lambda function is nil")
	}
	if err := function.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsFunction := awscomputemapper.FromDomainLambdaFunction(function)
	if err := awsFunction.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateLambdaFunction(ctx, awsFunction)
	if err != nil {
		return nil, nil, err
	}
	return awscomputemapper.ToDomainLambdaFunctionFromOutput(output), output, nil
}

// IAM helpers using AWS output models

func CreateRoleWithOutput(ctx context.Context, service awsiamservice.AWSIAMService, role *domainiam.Role) (*domainiam.Role, *awsiamoutputs.RoleOutput, error) {
	if role == nil {
		return nil, nil, fmt.Errorf("role is nil")
	}
	if err := role.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsRole := awsiammapper.FromDomainRole(role)
	if err := awsRole.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateRole(ctx, awsRole)
	if err != nil {
		return nil, nil, err
	}
	return awsiammapper.ToDomainRoleFromOutput(output), output, nil
}

func CreateInstanceProfileWithOutput(ctx context.Context, service awsiamservice.AWSIAMService, profile *domainiam.InstanceProfile) (*domainiam.InstanceProfile, *awsiamoutputs.InstanceProfileOutput, error) {
	if profile == nil {
		return nil, nil, fmt.Errorf("instance profile is nil")
	}
	if err := profile.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsProfile := awsiammapper.FromDomainInstanceProfile(profile)
	if err := awsProfile.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateInstanceProfile(ctx, awsProfile)
	if err != nil {
		return nil, nil, err
	}
	return awsiammapper.ToDomainInstanceProfileFromOutput(output), output, nil
}

func CreatePolicyWithOutput(ctx context.Context, service awsiamservice.AWSIAMService, policy *domainiam.Policy) (*domainiam.Policy, *awsiamoutputs.PolicyOutput, error) {
	if policy == nil {
		return nil, nil, fmt.Errorf("policy is nil")
	}
	if err := policy.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsPolicy := awsiammapper.FromDomainPolicy(policy)
	if err := awsPolicy.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreatePolicy(ctx, awsPolicy)
	if err != nil {
		return nil, nil, err
	}
	return awsiammapper.ToDomainPolicyFromOutput(output), output, nil
}

func AddRoleToInstanceProfile(ctx context.Context, service awsiamservice.AWSIAMService, profileName, roleName string) error {
	return service.AddRoleToInstanceProfile(ctx, profileName, roleName)
}

func AttachPolicyToRole(ctx context.Context, service awsiamservice.AWSIAMService, policyARN, roleName string) error {
	return service.AttachPolicyToRole(ctx, policyARN, roleName)
}

// Storage helpers using AWS output models

func CreateS3BucketWithOutput(ctx context.Context, service awsstorageservice.AWSStorageService, bucket *domainstorage.S3Bucket) (*domainstorage.S3Bucket, *awss3outputs.BucketOutput, error) {
	if bucket == nil {
		return nil, nil, fmt.Errorf("s3 bucket is nil")
	}
	if err := bucket.Validate(); err != nil {
		return nil, nil, fmt.Errorf("domain validation failed: %w", err)
	}
	awsBucket := awsstoragemapper.FromDomainS3Bucket(bucket)
	if err := awsBucket.Validate(); err != nil {
		return nil, nil, fmt.Errorf("aws validation failed: %w", err)
	}
	output, err := service.CreateS3Bucket(ctx, awsBucket)
	if err != nil {
		return nil, nil, err
	}
	domainBucket := awsstoragemapper.ToDomainS3BucketFromOutput(output)
	domainBucket.Region = bucket.Region
	return domainBucket, output, nil
}
