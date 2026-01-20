package scenario2_high_availability

import (
	"context"
	"fmt"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	awspricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing"
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
	fmt.Println("\n[MOCK MODE] All resources are simulated - no AWS SDK calls")

	// Initialize mock ID generator
	gen := usecasescommon.NewMockIDGenerator()

	// Step 1: Get availability zones
	fmt.Println("\n--- Step 1: Selecting Region with Multiple AZs ---")
	azs := usecasescommon.GetDefaultAvailabilityZones(region)
	fmt.Printf("✓ Selected region with %d availability zones: %v\n", len(azs), azs)

	// Step 2: Create VPC
	fmt.Println("\n--- Step 2: Creating VPC ---")
	vpc := usecasescommon.CreateMockVPC(region, "ha-vpc", "10.0.0.0/16", gen)
	fmt.Printf("✓ VPC created: %s (%s)\n", vpc.Name, vpc.ID)
	fmt.Printf("  CIDR: %s\n", vpc.CIDR)

	// Step 3: Create public subnets across 3 AZs
	fmt.Println("\n--- Step 3: Creating Public Subnets (Multi-AZ) ---")
	publicSubnets := []*domainnetworking.Subnet{}
	for i, az := range azs {
		subnet := usecasescommon.CreateMockSubnet(
			vpc.ID,
			fmt.Sprintf("public-subnet-%s", az[len(az)-1:]),
			fmt.Sprintf("10.0.%d.0/24", i+1),
			az,
			gen,
		)
		publicSubnets = append(publicSubnets, subnet)
		fmt.Printf("✓ Created: %s (%s) in %s\n", subnet.Name, subnet.ID, az)
	}

	// Step 4: Create private subnets across 3 AZs
	fmt.Println("\n--- Step 4: Creating Private Subnets (Multi-AZ) ---")
	privateSubnets := []*domainnetworking.Subnet{}
	for i, az := range azs {
		subnet := usecasescommon.CreateMockSubnet(
			vpc.ID,
			fmt.Sprintf("private-subnet-%s", az[len(az)-1:]),
			fmt.Sprintf("10.0.%d.0/24", i+10),
			az,
			gen,
		)
		privateSubnets = append(privateSubnets, subnet)
		fmt.Printf("✓ Created: %s (%s) in %s\n", subnet.Name, subnet.ID, az)
	}

	// Step 5: Create Internet Gateway
	fmt.Println("\n--- Step 5: Creating Internet Gateway ---")
	igw := usecasescommon.CreateMockInternetGateway(vpc.ID, "ha-igw", gen)
	fmt.Printf("✓ Internet Gateway created: %s (%s)\n", igw.Name, igw.ID)
	fmt.Printf("  Attached to VPC: %s\n", igw.VPCID)

	// Step 6: Create NAT Gateway in public subnet
	fmt.Println("\n--- Step 6: Creating NAT Gateway ---")
	natGateway := usecasescommon.CreateMockNATGateway(publicSubnets[0].ID, "ha-nat-gateway", gen)
	fmt.Printf("✓ NAT Gateway created: %s (%s)\n", natGateway.Name, natGateway.ID)
	fmt.Printf("  Located in: %s\n", natGateway.SubnetID)
	if natGateway.ARN != nil {
		fmt.Printf("  ARN: %s\n", *natGateway.ARN)
	}

	// Step 7: Configure route tables
	fmt.Println("\n--- Step 7: Configuring Route Tables ---")
	publicRT := usecasescommon.CreateMockRouteTable(vpc.ID, "public-route-table", gen)
	privateRT := usecasescommon.CreateMockRouteTable(vpc.ID, "private-route-table", gen)
	fmt.Printf("✓ Public Route Table: %s\n", publicRT.Name)
	fmt.Printf("  Route: 0.0.0.0/0 -> %s (Internet Gateway)\n", igw.ID)
	fmt.Printf("✓ Private Route Table: %s\n", privateRT.Name)
	fmt.Printf("  Route: 0.0.0.0/0 -> %s (NAT Gateway)\n", natGateway.ID)

	// Step 8: Create security groups
	fmt.Println("\n--- Step 8: Creating Security Groups ---")
	albSG := usecasescommon.CreateMockSecurityGroup(vpc.ID, "alb-sg", "Security group for Application Load Balancer", gen)
	appSG := usecasescommon.CreateMockSecurityGroup(vpc.ID, "app-sg", "Security group for application servers", gen)
	fmt.Printf("✓ Created security groups:\n")
	fmt.Printf("  - ALB: %s (%s)\n", albSG.Name, albSG.ID)
	fmt.Printf("  - App: %s (%s)\n", appSG.Name, appSG.ID)

	// Step 9: Create Application Load Balancer
	fmt.Println("\n--- Step 9: Creating Application Load Balancer ---")
	publicSubnetIDs := []string{}
	for _, subnet := range publicSubnets {
		publicSubnetIDs = append(publicSubnetIDs, subnet.ID)
	}
	alb := usecasescommon.CreateMockLoadBalancer("ha-alb", "application", publicSubnetIDs, []string{albSG.ID}, region, gen)
	fmt.Printf("✓ Application Load Balancer created: %s\n", alb.Name)
	if alb.ARN != nil {
		fmt.Printf("  ARN: %s\n", *alb.ARN)
	}
	if alb.DNSName != nil {
		fmt.Printf("  DNS Name: %s\n", *alb.DNSName)
	}
	fmt.Printf("  Subnets: %d across %d AZs\n", len(publicSubnets), len(azs))

	// Step 10: Create Target Group
	fmt.Println("\n--- Step 10: Creating Target Group ---")
	targetGroup := usecasescommon.CreateMockTargetGroup("ha-target-group", vpc.ID, "HTTP", 80, region, gen)
	fmt.Printf("✓ Target Group created: %s\n", targetGroup.Name)
	if targetGroup.ARN != nil {
		fmt.Printf("  ARN: %s\n", *targetGroup.ARN)
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
		instance := usecasescommon.CreateMockEC2Instance(
			fmt.Sprintf("app-server-%s", *subnet.AvailabilityZone),
			"t3.small",
			subnet.ID,
			appSG.ID,
			region,
			*subnet.AvailabilityZone,
			gen,
		)
		appInstances = append(appInstances, instance)
		fmt.Printf("✓ Created: %s (%s) in %s\n", instance.Name, instance.ID, *subnet.AvailabilityZone)
	}

	// Step 12: Attach instances to target group
	fmt.Println("\n--- Step 12: Attaching Instances to Target Group ---")
	fmt.Printf("✓ Attached %d instances to target group:\n", len(appInstances))
	for _, instance := range appInstances {
		fmt.Printf("  - %s (%s)\n", instance.Name, instance.ID)
	}

	// Step 13: Create Launch Template
	fmt.Println("\n--- Step 13: Creating Launch Template ---")
	launchTemplate := usecasescommon.CreateMockLaunchTemplate(
		"ha-launch-template",
		"t3.small",
		"ami-0c55b159cbfafe1f0",
		[]string{appSG.ID},
		region,
		gen,
	)
	fmt.Printf("✓ Launch Template created: %s (%s)\n", launchTemplate.Name, launchTemplate.ID)
	if launchTemplate.ARN != nil {
		fmt.Printf("  ARN: %s\n", *launchTemplate.ARN)
	}
	fmt.Printf("  Instance Type: %s\n", launchTemplate.InstanceType)

	// Step 14: Create Auto Scaling Group
	fmt.Println("\n--- Step 14: Creating Auto Scaling Group ---")
	privateSubnetIDs := []string{}
	for _, subnet := range privateSubnets {
		privateSubnetIDs = append(privateSubnetIDs, subnet.ID)
	}
	asg := usecasescommon.CreateMockAutoScalingGroup(
		"ha-asg",
		2, // min size
		6, // max size
		3, // desired capacity
		privateSubnetIDs,
		launchTemplate.ID,
		region,
		gen,
	)
	fmt.Printf("✓ Auto Scaling Group created: %s\n", asg.Name)
	if asg.ARN != nil {
		fmt.Printf("  ARN: %s\n", *asg.ARN)
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
		Type: resource.ResourceType{Name: "load_balancer"},
		Provider: "aws",
		Region: region,
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
		Type: resource.ResourceType{Name: "nat_gateway"},
		Provider: "aws",
		Region: region,
	}
	natEstimate, err := pricingService.EstimateCost(ctx, natRes, duration)
	if err == nil {
		fmt.Printf("  NAT Gateway: $%.2f\n", natEstimate.TotalCost)
		totalCost += natEstimate.TotalCost
	}

	// Auto Scaling Group (based on average capacity)
	asgRes := &resource.Resource{
		Type: resource.ResourceType{Name: "auto_scaling_group"},
		Provider: "aws",
		Region: region,
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
