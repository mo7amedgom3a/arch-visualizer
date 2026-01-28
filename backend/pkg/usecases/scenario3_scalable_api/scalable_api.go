package scenario3_scalable_api

import (
	"context"
	"fmt"
	"time"

	awspricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing"
	awscomputeservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	awsiamservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/iam"
	awsnetworkingservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	domainiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/iam"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	usecasescommon "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/common"
)

// ScalableAPIRunner demonstrates a scalable API backend with auto-scaling and IAM roles
func ScalableAPIRunner() {
	ctx := context.Background()
	region := usecasescommon.SelectRegion("us-east-1")

	fmt.Println("============================================")
	fmt.Println("SCENARIO 3: SCALABLE API ARCHITECTURE")
	fmt.Println("============================================")
	fmt.Printf("Region: %s\n", usecasescommon.FormatRegionName(region))
	fmt.Println("\n[OUTPUT MODE] Domain models + AWS output models")

	// Initialize virtual services
	networkingService := awsnetworkingservice.NewNetworkingService()
	computeService := awscomputeservice.NewComputeService()
	iamService := awsiamservice.NewIAMService()

	// Step 1: Select region (already done above)
	fmt.Println("\n--- Step 1: Region Selection ---")
	fmt.Printf("✓ Selected region: %s\n", usecasescommon.FormatRegionName(region))

	// Step 2: Create VPC and networking (similar to scenario 2)
	fmt.Println("\n--- Step 2: Creating VPC and Networking ---")
	vpc, vpcOutput, err := usecasescommon.CreateVPCWithOutput(ctx, networkingService, &domainnetworking.VPC{
		Name:               "api-vpc",
		Region:             region,
		CIDR:               "10.0.0.0/16",
		EnableDNS:          true,
		EnableDNSHostnames: true,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create VPC: %v\n", err)
		return
	}
	azs := usecasescommon.GetDefaultAvailabilityZones(region)

	// Create public subnets
	publicSubnets := []*domainnetworking.Subnet{}
	for i, az := range azs {
		subnet, subnetOutput, err := usecasescommon.CreateSubnetWithOutput(ctx, networkingService, &domainnetworking.Subnet{
			Name:             fmt.Sprintf("api-public-%s", az[len(az)-1:]),
			VPCID:            vpc.ID,
			CIDR:             fmt.Sprintf("10.0.%d.0/24", i+1),
			AvailabilityZone: &az,
			IsPublic:         true,
		})
		if err != nil {
			fmt.Printf("✗ Failed to create public subnet: %v\n", err)
			return
		}
		publicSubnets = append(publicSubnets, subnet)
		_ = subnetOutput
	}

	// Create private subnets
	privateSubnets := []*domainnetworking.Subnet{}
	for i, az := range azs {
		subnet, subnetOutput, err := usecasescommon.CreateSubnetWithOutput(ctx, networkingService, &domainnetworking.Subnet{
			Name:             fmt.Sprintf("api-private-%s", az[len(az)-1:]),
			VPCID:            vpc.ID,
			CIDR:             fmt.Sprintf("10.0.%d.0/24", i+10),
			AvailabilityZone: &az,
			IsPublic:         false,
		})
		if err != nil {
			fmt.Printf("✗ Failed to create private subnet: %v\n", err)
			return
		}
		privateSubnets = append(privateSubnets, subnet)
		_ = subnetOutput
	}

	igw, igwOutput, err := usecasescommon.CreateInternetGatewayWithOutput(ctx, networkingService, &domainnetworking.InternetGateway{
		Name:  "api-igw",
		VPCID: vpc.ID,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create Internet Gateway: %v\n", err)
		return
	}
	if err := usecasescommon.AttachInternetGateway(ctx, networkingService, igw.ID, vpc.ID); err != nil {
		fmt.Printf("✗ Failed to attach Internet Gateway: %v\n", err)
		return
	}
	eip, eipOutput, err := usecasescommon.AllocateElasticIPWithOutput(ctx, networkingService, &domainnetworking.ElasticIP{
		Region: region,
	})
	if err != nil {
		fmt.Printf("✗ Failed to allocate Elastic IP: %v\n", err)
		return
	}
	natGateway, natOutput, err := usecasescommon.CreateNATGatewayWithOutput(ctx, networkingService, &domainnetworking.NATGateway{
		Name:         "api-nat-gateway",
		SubnetID:     publicSubnets[0].ID,
		AllocationID: &eip.ID,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create NAT Gateway: %v\n", err)
		return
	}

	fmt.Printf("✓ VPC: %s (%s)\n", vpcOutput.Name, vpcOutput.CIDR)
	fmt.Printf("✓ Public Subnets: %d across %d AZs\n", len(publicSubnets), len(azs))
	fmt.Printf("✓ Private Subnets: %d across %d AZs\n", len(privateSubnets), len(azs))
	fmt.Printf("✓ Internet Gateway: %s\n", igwOutput.Name)
	fmt.Printf("✓ NAT Gateway: %s\n", natOutput.Name)
	_ = eipOutput

	// Step 3: Create IAM Role for EC2 instances
	fmt.Println("\n--- Step 3: Creating IAM Role for EC2 Instances ---")
	assumeRolePolicy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"Service": "ec2.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}
		]
	}`
	roleDescription := "Role for API EC2 instances to access AWS services"
	iamRole, roleOutput, err := usecasescommon.CreateRoleWithOutput(ctx, iamService, &domainiam.Role{
		Name:             "api-ec2-role",
		Description:      &roleDescription,
		Path:             stringPtr("/"),
		AssumeRolePolicy: assumeRolePolicy,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create IAM Role: %v\n", err)
		return
	}
	fmt.Printf("✓ IAM Role created: %s\n", roleOutput.Name)
	if roleOutput.ARN != "" {
		fmt.Printf("  ARN: %s\n", roleOutput.ARN)
	}
	if iamRole.Description != nil {
		fmt.Printf("  Description: %s\n", *iamRole.Description)
	}
	if iamRole.Path != nil {
		fmt.Printf("  Path: %s\n", *iamRole.Path)
	}

	// Step 4: Create Instance Profile
	fmt.Println("\n--- Step 4: Creating IAM Instance Profile ---")
	instanceProfile, profileOutput, err := usecasescommon.CreateInstanceProfileWithOutput(ctx, iamService, &domainiam.InstanceProfile{
		Name: "api-instance-profile",
		Path: stringPtr("/"),
	})
	if err != nil {
		fmt.Printf("✗ Failed to create IAM Instance Profile: %v\n", err)
		return
	}
	if err := usecasescommon.AddRoleToInstanceProfile(ctx, iamService, instanceProfile.Name, iamRole.Name); err != nil {
		fmt.Printf("✗ Failed to attach role to instance profile: %v\n", err)
		return
	}
	fmt.Printf("✓ Instance Profile created: %s\n", profileOutput.Name)
	if profileOutput.ARN != "" {
		fmt.Printf("  ARN: %s\n", profileOutput.ARN)
	}
	fmt.Printf("  Associated Role: %s\n", iamRole.Name)
	if instanceProfile.Path != nil {
		fmt.Printf("  Path: %s\n", *instanceProfile.Path)
	}

	// Step 5: Create Launch Template with IAM role
	fmt.Println("\n--- Step 5: Creating Launch Template with IAM Role ---")
	appSG, appSGOutput, err := usecasescommon.CreateSecurityGroupWithOutput(ctx, networkingService, &domainnetworking.SecurityGroup{
		Name:        "api-sg",
		VPCID:       vpc.ID,
		Description: "Security group for API servers",
	})
	if err != nil {
		fmt.Printf("✗ Failed to create API security group: %v\n", err)
		return
	}
	launchTemplate, ltOutput, err := usecasescommon.CreateLaunchTemplateWithOutput(ctx, computeService, &domaincompute.LaunchTemplate{
		Name:               "api-launch-template",
		Region:             region,
		ImageID:            "ami-0c55b159cbfafe1f0",
		InstanceType:       "t3.medium",
		SecurityGroupIDs:   []string{appSG.ID},
		IAMInstanceProfile: &instanceProfile.Name,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create launch template: %v\n", err)
		return
	}
	fmt.Printf("✓ Launch Template created: %s (%s)\n", ltOutput.Name, ltOutput.ID)
	if ltOutput.ARN != "" {
		fmt.Printf("  ARN: %s\n", ltOutput.ARN)
	}
	fmt.Printf("  Instance Type: %s\n", launchTemplate.InstanceType)
	fmt.Printf("  Security Groups: %v\n", launchTemplate.SecurityGroupIDs)
	fmt.Printf("  Note: IAM Instance Profile (%s) will be attached to instances\n", instanceProfile.Name)
	_ = appSGOutput

	// Step 6: Create Application Load Balancer
	fmt.Println("\n--- Step 6: Creating Application Load Balancer ---")
	albSG, albSGOutput, err := usecasescommon.CreateSecurityGroupWithOutput(ctx, networkingService, &domainnetworking.SecurityGroup{
		Name:        "api-alb-sg",
		VPCID:       vpc.ID,
		Description: "Security group for API Load Balancer",
	})
	if err != nil {
		fmt.Printf("✗ Failed to create ALB security group: %v\n", err)
		return
	}
	publicSubnetIDs := []string{}
	for _, subnet := range publicSubnets {
		publicSubnetIDs = append(publicSubnetIDs, subnet.ID)
	}
	alb, albOutput, err := usecasescommon.CreateLoadBalancerWithOutput(ctx, computeService, &domaincompute.LoadBalancer{
		Name:             "api-alb",
		Region:           region,
		Type:             domaincompute.LoadBalancerTypeApplication,
		Internal:         false,
		SecurityGroupIDs: []string{albSG.ID},
		SubnetIDs:        publicSubnetIDs,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create load balancer: %v\n", err)
		return
	}
	fmt.Printf("✓ Application Load Balancer created: %s\n", albOutput.Name)
	if albOutput.ARN != "" {
		fmt.Printf("  ARN: %s\n", albOutput.ARN)
	}
	if albOutput.DNSName != "" {
		fmt.Printf("  DNS Name: %s\n", albOutput.DNSName)
	}
	fmt.Printf("  Subnets: %d across %d AZs\n", len(publicSubnets), len(azs))
	_ = albSGOutput

	// Step 7: Create Target Group
	fmt.Println("\n--- Step 7: Creating Target Group ---")
	healthPath := "/health"
	healthPort := "8080"
	healthProtocol := "HTTP"
	targetGroup, tgOutput, err := usecasescommon.CreateTargetGroupWithOutput(ctx, computeService, &domaincompute.TargetGroup{
		Name:       "api-target-group",
		VPCID:      vpc.ID,
		Protocol:   domaincompute.TargetGroupProtocolHTTP,
		Port:       8080,
		TargetType: domaincompute.TargetTypeInstance,
		HealthCheck: domaincompute.HealthCheckConfig{
			Path:     &healthPath,
			Protocol: &healthProtocol,
			Port:     &healthPort,
		},
	})
	if err != nil {
		fmt.Printf("✗ Failed to create target group: %v\n", err)
		return
	}
	fmt.Printf("✓ Target Group created: %s\n", tgOutput.Name)
	if tgOutput.ARN != "" {
		fmt.Printf("  ARN: %s\n", tgOutput.ARN)
	}
	fmt.Printf("  Protocol: %s, Port: %d\n", targetGroup.Protocol, targetGroup.Port)
	if targetGroup.HealthCheck.Path != nil {
		protocol := "HTTP"
		if targetGroup.HealthCheck.Protocol != nil {
			protocol = *targetGroup.HealthCheck.Protocol
		}
		fmt.Printf("  Health Check: %s%s\n", protocol, *targetGroup.HealthCheck.Path)
	}

	// Step 8: Create Auto Scaling Group with scaling policies
	fmt.Println("\n--- Step 8: Creating Auto Scaling Group ---")
	privateSubnetIDs := []string{}
	for _, subnet := range privateSubnets {
		privateSubnetIDs = append(privateSubnetIDs, subnet.ID)
	}
	desiredCapacity := 4
	version := "$Latest"
	targetGroupARN := targetGroup.ID
	if targetGroup.ARN != nil {
		targetGroupARN = *targetGroup.ARN
	}
	gracePeriod := 300
	asg, asgOutput, err := usecasescommon.CreateAutoScalingGroupWithOutput(ctx, computeService, &domaincompute.AutoScalingGroup{
		Name:              "api-asg",
		Region:            region,
		MinSize:           2,
		MaxSize:           10,
		DesiredCapacity:   &desiredCapacity,
		VPCZoneIdentifier: privateSubnetIDs,
		LaunchTemplate: &domaincompute.LaunchTemplateSpecification{
			ID:      launchTemplate.ID,
			Version: &version,
		},
		HealthCheckType:        domaincompute.AutoScalingGroupHealthCheckTypeELB,
		HealthCheckGracePeriod: &gracePeriod,
		TargetGroupARNs:        []string{targetGroupARN},
	})
	if err != nil {
		fmt.Printf("✗ Failed to create Auto Scaling Group: %v\n", err)
		return
	}
	fmt.Printf("✓ Auto Scaling Group created: %s\n", asgOutput.AutoScalingGroupName)
	if asgOutput.AutoScalingGroupARN != "" {
		fmt.Printf("  ARN: %s\n", asgOutput.AutoScalingGroupARN)
	}
	fmt.Printf("  Min Size: %d\n", asg.MinSize)
	fmt.Printf("  Max Size: %d\n", asg.MaxSize)
	if asg.DesiredCapacity != nil {
		fmt.Printf("  Desired Capacity: %d\n", *asg.DesiredCapacity)
	}
	fmt.Printf("  Launch Template: %s\n", launchTemplate.Name)
	fmt.Printf("  Subnets: %d across %d AZs\n", len(privateSubnets), len(azs))

	// Step 9: Configure health checks (ELB-based)
	fmt.Println("\n--- Step 9: Configuring Health Checks ---")
	fmt.Printf("✓ Health Check Type: ELB (Load Balancer)\n")
	if asg.HealthCheckGracePeriod != nil {
		fmt.Printf("✓ Health Check Grace Period: %d seconds\n", *asg.HealthCheckGracePeriod)
	}
	fmt.Printf("✓ Target Group Health Check:\n")
	if targetGroup.HealthCheck.Protocol != nil {
		fmt.Printf("  - Protocol: %s\n", *targetGroup.HealthCheck.Protocol)
	}
	if targetGroup.HealthCheck.Port != nil {
		fmt.Printf("  - Port: %s\n", *targetGroup.HealthCheck.Port)
	}
	if targetGroup.HealthCheck.Path != nil {
		fmt.Printf("  - Path: %s\n", *targetGroup.HealthCheck.Path)
	}
	fmt.Printf("  - Target Group ARN: %s\n", targetGroupARN)

	// Step 10: Display scaling configuration
	fmt.Println("\n--- Step 10: Scaling Configuration ---")
	fmt.Printf("✓ Auto Scaling Configuration:\n")
	fmt.Printf("  - Capacity Range: %d - %d instances\n", asg.MinSize, asg.MaxSize)
	fmt.Printf("  - Current Desired: %d instances\n", *asg.DesiredCapacity)
	fmt.Printf("  - Scaling Policies: (simulated)\n")
	fmt.Printf("    * Target Tracking: CPU utilization at 70%%\n")
	fmt.Printf("    * Scale Out: Add 2 instances when CPU > 80%%\n")
	fmt.Printf("    * Scale In: Remove 1 instance when CPU < 30%%\n")
	fmt.Printf("  - Cooldown Period: 300 seconds\n")

	// Step 11: Simulate scaling events
	fmt.Println("\n--- Step 11: Simulating Scaling Events ---")
	fmt.Println("  [Simulation] High traffic detected...")
	fmt.Printf("  → CPU utilization: 85%%\n")
	fmt.Printf("  → Scaling out: Adding 2 instances\n")
	fmt.Printf("  → New desired capacity: %d instances\n", *asg.DesiredCapacity+2)
	fmt.Println("\n  [Simulation] Traffic normalized...")
	fmt.Printf("  → CPU utilization: 25%%\n")
	fmt.Printf("  → Scaling in: Removing 1 instance\n")
	fmt.Printf("  → New desired capacity: %d instances\n", *asg.DesiredCapacity+1)

	// Step 12: Display architecture summary
	fmt.Println("\n============================================")
	fmt.Println("SCALABLE API ARCHITECTURE SUMMARY")
	fmt.Println("============================================")
	fmt.Printf("VPC: %s (%s)\n", vpc.Name, vpc.CIDR)
	fmt.Printf("\nNetworking:\n")
	fmt.Printf("  Public Subnets: %d (across %d AZs)\n", len(publicSubnets), len(azs))
	fmt.Printf("  Private Subnets: %d (across %d AZs)\n", len(privateSubnets), len(azs))
	fmt.Printf("  Internet Gateway: %s\n", igw.Name)
	fmt.Printf("  NAT Gateway: %s\n", natGateway.Name)
	fmt.Printf("\nIAM:\n")
	fmt.Printf("  IAM Role: %s\n", iamRole.Name)
	fmt.Printf("  Instance Profile: %s\n", instanceProfile.Name)
	fmt.Printf("\nLoad Balancing:\n")
	fmt.Printf("  Application Load Balancer: %s\n", alb.Name)
	fmt.Printf("  Target Group: %s\n", targetGroup.Name)
	fmt.Printf("\nAuto Scaling:\n")
	fmt.Printf("  Auto Scaling Group: %s\n", asg.Name)
	fmt.Printf("  Capacity: %d-%d instances (desired: %d)\n", asg.MinSize, asg.MaxSize, *asg.DesiredCapacity)
	fmt.Printf("  Launch Template: %s\n", launchTemplate.Name)
	fmt.Printf("  Health Check: ELB-based\n")
	fmt.Printf("\nScalability Features:\n")
	fmt.Printf("  ✓ Multi-AZ deployment for high availability\n")
	fmt.Printf("  ✓ Auto scaling based on CPU utilization\n")
	fmt.Printf("  ✓ Load balancer distributes traffic\n")
	fmt.Printf("  ✓ IAM roles for secure AWS service access\n")
	fmt.Printf("  ✓ Health checks ensure instance availability\n")

	// Step 13: Calculate costs
	fmt.Println("\n============================================")
	fmt.Println("COST ESTIMATION (30 days)")
	fmt.Println("============================================")
	pricingService := awspricing.NewAWSPricingService()
	duration := 30 * 24 * time.Hour // 30 days

	fmt.Println("\nResource Costs:")
	totalCost := 0.0

	// Load Balancer
	lbRes := &resource.Resource{
		Type:     resource.ResourceType{Name: "load_balancer"},
		Provider: "aws",
		Region:   region,
		Metadata: map[string]interface{}{
			"load_balancer_type": "application",
		},
	}
	lbEstimate, err := pricingService.EstimateCost(ctx, lbRes, duration)
	if err == nil {
		fmt.Printf("  Application Load Balancer: $%.2f\n", lbEstimate.TotalCost)
		totalCost += lbEstimate.TotalCost
	}

	// NAT Gateway
	natRes := &resource.Resource{
		Type:     resource.ResourceType{Name: "nat_gateway"},
		Provider: "aws",
		Region:   region,
	}
	natEstimate, err := pricingService.EstimateCost(ctx, natRes, duration)
	if err == nil {
		fmt.Printf("  NAT Gateway: $%.2f\n", natEstimate.TotalCost)
		totalCost += natEstimate.TotalCost
	}

	// Auto Scaling Group (based on average capacity)
	asgRes := &resource.Resource{
		Type:     resource.ResourceType{Name: "auto_scaling_group"},
		Provider: "aws",
		Region:   region,
		Metadata: map[string]interface{}{
			"instance_type": launchTemplate.InstanceType,
			"min_size":      asg.MinSize,
			"max_size":      asg.MaxSize,
		},
	}
	asgEstimate, err := pricingService.EstimateCost(ctx, asgRes, duration)
	if err == nil {
		fmt.Printf("  Auto Scaling Group (avg capacity): $%.2f\n", asgEstimate.TotalCost)
		totalCost += asgEstimate.TotalCost
	}

	fmt.Printf("\nTotal Estimated Cost (30 days): $%.2f\n", totalCost)
	fmt.Printf("Note: VPC, subnets, route tables, security groups, IGW, and IAM are free\n")
	fmt.Printf("      Cost varies based on actual instance count (scales dynamically)\n")

	fmt.Println("\n============================================")
	fmt.Println("SCENARIO COMPLETE")
	fmt.Println("============================================")
}

func stringPtr(s string) *string {
	return &s
}
