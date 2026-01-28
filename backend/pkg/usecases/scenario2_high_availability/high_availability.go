package scenario2_high_availability

import (
	"context"
	"fmt"
	"time"

	awspricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing"
	awscomputeservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	awsnetworkingservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	usecasescommon "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/common"
)

// HighAvailabilityRunner demonstrates a multi-AZ high availability architecture with load balancer
func HighAvailabilityRunner() {
	ctx := context.Background()
	region := usecasescommon.SelectRegion("us-east-1")

	fmt.Println("============================================")
	fmt.Println("SCENARIO 2: HIGH AVAILABILITY ARCHITECTURE")
	fmt.Println("============================================")
	fmt.Printf("Region: %s\n", usecasescommon.FormatRegionName(region))
	fmt.Println("\n[OUTPUT MODE] Domain models + AWS output models")

	// Initialize virtual services
	networkingService := awsnetworkingservice.NewNetworkingService()
	computeService := awscomputeservice.NewComputeService()

	// Step 1: Get availability zones
	fmt.Println("\n--- Step 1: Selecting Region with Multiple AZs ---")
	azs := usecasescommon.GetDefaultAvailabilityZones(region)
	fmt.Printf("✓ Selected region with %d availability zones: %v\n", len(azs), azs)

	// Step 2: Create VPC
	fmt.Println("\n--- Step 2: Creating VPC ---")
	vpc, vpcOutput, err := usecasescommon.CreateVPCWithOutput(ctx, networkingService, &domainnetworking.VPC{
		Name:               "ha-vpc",
		Region:             region,
		CIDR:               "10.0.0.0/16",
		EnableDNS:          true,
		EnableDNSHostnames: true,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create VPC: %v\n", err)
		return
	}
	fmt.Printf("✓ VPC created: %s (%s)\n", vpcOutput.Name, vpcOutput.ID)
	fmt.Printf("  CIDR: %s\n", vpcOutput.CIDR)

	// Step 3: Create public subnets across 3 AZs
	fmt.Println("\n--- Step 3: Creating Public Subnets (Multi-AZ) ---")
	publicSubnets := []*domainnetworking.Subnet{}
	for i, az := range azs {
		subnet, subnetOutput, err := usecasescommon.CreateSubnetWithOutput(ctx, networkingService, &domainnetworking.Subnet{
			Name:             fmt.Sprintf("public-subnet-%s", az[len(az)-1:]),
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
		fmt.Printf("✓ Created: %s (%s) in %s\n", subnetOutput.Name, subnetOutput.ID, az)
	}

	// Step 4: Create private subnets across 3 AZs
	fmt.Println("\n--- Step 4: Creating Private Subnets (Multi-AZ) ---")
	privateSubnets := []*domainnetworking.Subnet{}
	for i, az := range azs {
		subnet, subnetOutput, err := usecasescommon.CreateSubnetWithOutput(ctx, networkingService, &domainnetworking.Subnet{
			Name:             fmt.Sprintf("private-subnet-%s", az[len(az)-1:]),
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
		fmt.Printf("✓ Created: %s (%s) in %s\n", subnetOutput.Name, subnetOutput.ID, az)
	}

	// Step 5: Create Internet Gateway
	fmt.Println("\n--- Step 5: Creating Internet Gateway ---")
	igw, igwOutput, err := usecasescommon.CreateInternetGatewayWithOutput(ctx, networkingService, &domainnetworking.InternetGateway{
		Name:  "ha-igw",
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
	fmt.Printf("✓ Internet Gateway created: %s (%s)\n", igwOutput.Name, igwOutput.ID)
	fmt.Printf("  Attached to VPC: %s\n", igw.VPCID)

	// Step 6: Create NAT Gateway in public subnet
	fmt.Println("\n--- Step 6: Creating NAT Gateway ---")
	eip, eipOutput, err := usecasescommon.AllocateElasticIPWithOutput(ctx, networkingService, &domainnetworking.ElasticIP{
		Region: region,
	})
	if err != nil {
		fmt.Printf("✗ Failed to allocate Elastic IP: %v\n", err)
		return
	}
	natGateway, natOutput, err := usecasescommon.CreateNATGatewayWithOutput(ctx, networkingService, &domainnetworking.NATGateway{
		Name:         "ha-nat-gateway",
		SubnetID:     publicSubnets[0].ID,
		AllocationID: &eip.ID,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create NAT Gateway: %v\n", err)
		return
	}
	fmt.Printf("✓ NAT Gateway created: %s (%s)\n", natOutput.Name, natOutput.ID)
	fmt.Printf("  Located in: %s\n", natGateway.SubnetID)
	if natGateway.ARN != nil {
		fmt.Printf("  ARN: %s\n", *natGateway.ARN)
	}
	_ = eipOutput

	// Step 7: Configure route tables
	fmt.Println("\n--- Step 7: Configuring Route Tables ---")
	_, publicRTOutput, err := usecasescommon.CreateRouteTableWithOutput(ctx, networkingService, &domainnetworking.RouteTable{
		Name:  "public-route-table",
		VPCID: vpc.ID,
		Routes: []domainnetworking.Route{
			{
				DestinationCIDR: "0.0.0.0/0",
				TargetID:        igw.ID,
				TargetType:      "internet_gateway",
			},
		},
	})
	if err != nil {
		fmt.Printf("✗ Failed to create public route table: %v\n", err)
		return
	}
	_, privateRTOutput, err := usecasescommon.CreateRouteTableWithOutput(ctx, networkingService, &domainnetworking.RouteTable{
		Name:  "private-route-table",
		VPCID: vpc.ID,
		Routes: []domainnetworking.Route{
			{
				DestinationCIDR: "0.0.0.0/0",
				TargetID:        natGateway.ID,
				TargetType:      "nat_gateway",
			},
		},
	})
	if err != nil {
		fmt.Printf("✗ Failed to create private route table: %v\n", err)
		return
	}
	fmt.Printf("✓ Public Route Table: %s\n", publicRTOutput.Name)
	fmt.Printf("  Route: 0.0.0.0/0 -> %s (Internet Gateway)\n", igw.ID)
	fmt.Printf("✓ Private Route Table: %s\n", privateRTOutput.Name)
	fmt.Printf("  Route: 0.0.0.0/0 -> %s (NAT Gateway)\n", natGateway.ID)

	// Step 8: Create security groups
	fmt.Println("\n--- Step 8: Creating Security Groups ---")
	albSG, albSGOutput, err := usecasescommon.CreateSecurityGroupWithOutput(ctx, networkingService, &domainnetworking.SecurityGroup{
		Name:        "alb-sg",
		VPCID:       vpc.ID,
		Description: "Security group for Application Load Balancer",
	})
	if err != nil {
		fmt.Printf("✗ Failed to create ALB security group: %v\n", err)
		return
	}
	appSG, appSGOutput, err := usecasescommon.CreateSecurityGroupWithOutput(ctx, networkingService, &domainnetworking.SecurityGroup{
		Name:        "app-sg",
		VPCID:       vpc.ID,
		Description: "Security group for application servers",
	})
	if err != nil {
		fmt.Printf("✗ Failed to create app security group: %v\n", err)
		return
	}
	fmt.Printf("✓ Created security groups:\n")
	fmt.Printf("  - ALB: %s (%s)\n", albSGOutput.Name, albSGOutput.ID)
	fmt.Printf("  - App: %s (%s)\n", appSGOutput.Name, appSGOutput.ID)

	// Step 9: Create Application Load Balancer
	fmt.Println("\n--- Step 9: Creating Application Load Balancer ---")
	publicSubnetIDs := []string{}
	for _, subnet := range publicSubnets {
		publicSubnetIDs = append(publicSubnetIDs, subnet.ID)
	}
	alb, albOutput, err := usecasescommon.CreateLoadBalancerWithOutput(ctx, computeService, &domaincompute.LoadBalancer{
		Name:             "ha-alb",
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

	// Step 10: Create Target Group
	fmt.Println("\n--- Step 10: Creating Target Group ---")
	healthPath := "/health"
	healthPort := "80"
	healthProtocol := "HTTP"
	targetGroup, targetOutput, err := usecasescommon.CreateTargetGroupWithOutput(ctx, computeService, &domaincompute.TargetGroup{
		Name:       "ha-target-group",
		VPCID:      vpc.ID,
		Protocol:   domaincompute.TargetGroupProtocolHTTP,
		Port:       80,
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
	fmt.Printf("✓ Target Group created: %s\n", targetOutput.Name)
	if targetOutput.ARN != "" {
		fmt.Printf("  ARN: %s\n", targetOutput.ARN)
	}
	fmt.Printf("  Protocol: %s, Port: %d\n", targetGroup.Protocol, targetGroup.Port)
	if targetGroup.HealthCheck.Path != nil {
		protocol := "HTTP"
		if targetGroup.HealthCheck.Protocol != nil {
			protocol = *targetGroup.HealthCheck.Protocol
		}
		fmt.Printf("  Health Check: %s%s\n", protocol, *targetGroup.HealthCheck.Path)
	}

	// Step 11: Create EC2 instances in private subnets
	fmt.Println("\n--- Step 11: Creating EC2 Instances (Behind ALB) ---")
	appInstances := []*domaincompute.Instance{}
	for _, subnet := range privateSubnets {
		instance, instanceOutput, err := usecasescommon.CreateInstanceWithOutput(ctx, computeService, &domaincompute.Instance{
			Name:             fmt.Sprintf("app-server-%s", *subnet.AvailabilityZone),
			Region:           region,
			AvailabilityZone: subnet.AvailabilityZone,
			InstanceType:     "t3.small",
			AMI:              "ami-0c55b159cbfafe1f0",
			SubnetID:         subnet.ID,
			SecurityGroupIDs: []string{appSG.ID},
		})
		if err != nil {
			fmt.Printf("✗ Failed to create app instance: %v\n", err)
			return
		}
		appInstances = append(appInstances, instance)
		fmt.Printf("✓ Created: %s (%s) in %s\n", instanceOutput.Name, instanceOutput.ID, *subnet.AvailabilityZone)
	}

	// Step 12: Attach instances to target group
	fmt.Println("\n--- Step 12: Attaching Instances to Target Group ---")
	targetGroupARN := targetGroup.ID
	if targetGroup.ARN != nil {
		targetGroupARN = *targetGroup.ARN
	}
	for _, instance := range appInstances {
		if err := usecasescommon.AttachTargetToGroup(ctx, computeService, &domaincompute.TargetGroupAttachment{
			TargetGroupARN: targetGroupARN,
			TargetID:       instance.ID,
		}); err != nil {
			fmt.Printf("✗ Failed to attach instance to target group: %v\n", err)
			return
		}
	}
	fmt.Printf("✓ Attached %d instances to target group:\n", len(appInstances))
	for _, instance := range appInstances {
		fmt.Printf("  - %s (%s)\n", instance.Name, instance.ID)
	}

	// Step 13: Create Launch Template
	fmt.Println("\n--- Step 13: Creating Launch Template ---")
	launchTemplate, launchTemplateOutput, err := usecasescommon.CreateLaunchTemplateWithOutput(ctx, computeService, &domaincompute.LaunchTemplate{
		Name:             "ha-launch-template",
		Region:           region,
		ImageID:          "ami-0c55b159cbfafe1f0",
		InstanceType:     "t3.small",
		SecurityGroupIDs: []string{appSG.ID},
	})
	if err != nil {
		fmt.Printf("✗ Failed to create launch template: %v\n", err)
		return
	}
	fmt.Printf("✓ Launch Template created: %s (%s)\n", launchTemplateOutput.Name, launchTemplateOutput.ID)
	if launchTemplateOutput.ARN != "" {
		fmt.Printf("  ARN: %s\n", launchTemplateOutput.ARN)
	}
	fmt.Printf("  Instance Type: %s\n", launchTemplate.InstanceType)

	// Step 14: Create Auto Scaling Group
	fmt.Println("\n--- Step 14: Creating Auto Scaling Group ---")
	privateSubnetIDs := []string{}
	for _, subnet := range privateSubnets {
		privateSubnetIDs = append(privateSubnetIDs, subnet.ID)
	}
	desiredCapacity := 3
	version := "$Latest"
	asg, asgOutput, err := usecasescommon.CreateAutoScalingGroupWithOutput(ctx, computeService, &domaincompute.AutoScalingGroup{
		Name:              "ha-asg",
		Region:            region,
		MinSize:           2,
		MaxSize:           6,
		DesiredCapacity:   &desiredCapacity,
		VPCZoneIdentifier: privateSubnetIDs,
		LaunchTemplate: &domaincompute.LaunchTemplateSpecification{
			ID:      launchTemplate.ID,
			Version: &version,
		},
		HealthCheckType: domaincompute.AutoScalingGroupHealthCheckTypeEC2,
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
	fmt.Printf("  Health Check Type: %s\n", asg.HealthCheckType)
	fmt.Printf("  Subnets: %d across %d AZs\n", len(privateSubnets), len(azs))

	// Step 15: Display HA architecture summary
	fmt.Println("\n============================================")
	fmt.Println("HIGH AVAILABILITY ARCHITECTURE SUMMARY")
	fmt.Println("============================================")
	fmt.Printf("VPC: %s (%s)\n", vpc.Name, vpc.CIDR)
	fmt.Printf("Internet Gateway: %s\n", igw.Name)
	fmt.Printf("NAT Gateway: %s\n", natGateway.Name)
	fmt.Printf("\nNetworking:\n")
	fmt.Printf("  Public Subnets: %d (across %d AZs)\n", len(publicSubnets), len(azs))
	fmt.Printf("  Private Subnets: %d (across %d AZs)\n", len(privateSubnets), len(azs))
	fmt.Printf("  Route Tables: 2 (public with IGW, private with NAT)\n")
	fmt.Printf("  Security Groups: 2\n")
	fmt.Printf("\nLoad Balancing:\n")
	fmt.Printf("  Application Load Balancer: %s\n", alb.Name)
	fmt.Printf("  Target Group: %s\n", targetGroup.Name)
	fmt.Printf("  Instances Attached: %d\n", len(appInstances))
	fmt.Printf("\nAuto Scaling:\n")
	fmt.Printf("  Auto Scaling Group: %s\n", asg.Name)
	fmt.Printf("  Capacity: %d-%d instances (desired: %d)\n", asg.MinSize, asg.MaxSize, *asg.DesiredCapacity)
	fmt.Printf("  Launch Template: %s\n", launchTemplate.Name)
	fmt.Printf("\nHigh Availability Features:\n")
	fmt.Printf("  ✓ Multi-AZ deployment (%d availability zones)\n", len(azs))
	fmt.Printf("  ✓ Load balancer distributes traffic\n")
	fmt.Printf("  ✓ Auto scaling for dynamic capacity\n")
	fmt.Printf("  ✓ Health checks ensure instance availability\n")

	// Step 16: Calculate costs
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
	fmt.Printf("Note: VPC, subnets, route tables, security groups, and IGW are free\n")

	fmt.Println("\n============================================")
	fmt.Println("SCENARIO COMPLETE")
	fmt.Println("============================================")
}
