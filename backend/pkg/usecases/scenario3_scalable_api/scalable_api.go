package scenario3_scalable_api

import (
	"context"
	"fmt"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	awspricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing"
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
	fmt.Println("\n[MOCK MODE] All resources are simulated - no AWS SDK calls")

	// Initialize mock ID generator
	gen := usecasescommon.NewMockIDGenerator()

	// Step 1: Select region (already done above)
	fmt.Println("\n--- Step 1: Region Selection ---")
	fmt.Printf("✓ Selected region: %s\n", usecasescommon.FormatRegionName(region))

	// Step 2: Create VPC and networking (similar to scenario 2)
	fmt.Println("\n--- Step 2: Creating VPC and Networking ---")
	vpc := usecasescommon.CreateMockVPC(region, "api-vpc", "10.0.0.0/16", gen)
	azs := usecasescommon.GetDefaultAvailabilityZones(region)

	// Create public subnets
	publicSubnets := []*domainnetworking.Subnet{}
	for i, az := range azs {
		subnet := usecasescommon.CreateMockSubnet(
			vpc.ID,
			fmt.Sprintf("api-public-%s", az[len(az)-1:]),
			fmt.Sprintf("10.0.%d.0/24", i+1),
			az,
			gen,
		)
		publicSubnets = append(publicSubnets, subnet)
	}

	// Create private subnets
	privateSubnets := []*domainnetworking.Subnet{}
	for i, az := range azs {
		subnet := usecasescommon.CreateMockSubnet(
			vpc.ID,
			fmt.Sprintf("api-private-%s", az[len(az)-1:]),
			fmt.Sprintf("10.0.%d.0/24", i+10),
			az,
			gen,
		)
		privateSubnets = append(privateSubnets, subnet)
	}

	igw := usecasescommon.CreateMockInternetGateway(vpc.ID, "api-igw", gen)
	natGateway := usecasescommon.CreateMockNATGateway(publicSubnets[0].ID, "api-nat-gateway", gen)

	fmt.Printf("✓ VPC: %s (%s)\n", vpc.Name, vpc.CIDR)
	fmt.Printf("✓ Public Subnets: %d across %d AZs\n", len(publicSubnets), len(azs))
	fmt.Printf("✓ Private Subnets: %d across %d AZs\n", len(privateSubnets), len(azs))
	fmt.Printf("✓ Internet Gateway: %s\n", igw.Name)
	fmt.Printf("✓ NAT Gateway: %s\n", natGateway.Name)

	// Step 3: Create IAM Role for EC2 instances
	fmt.Println("\n--- Step 3: Creating IAM Role for EC2 Instances ---")
	iamRole := usecasescommon.CreateMockIAMRole(
		"api-ec2-role",
		"Role for API EC2 instances to access AWS services",
		gen,
	)
	fmt.Printf("✓ IAM Role created: %s\n", iamRole.Name)
	if iamRole.ARN != nil {
		fmt.Printf("  ARN: %s\n", *iamRole.ARN)
	}
	if iamRole.Description != nil {
		fmt.Printf("  Description: %s\n", *iamRole.Description)
	}
	fmt.Printf("  Path: %s\n", iamRole.Path)

	// Step 4: Create Instance Profile
	fmt.Println("\n--- Step 4: Creating IAM Instance Profile ---")
	instanceProfile := usecasescommon.CreateMockIAMInstanceProfile(
		"api-instance-profile",
		iamRole.Name,
		gen,
	)
	fmt.Printf("✓ Instance Profile created: %s\n", instanceProfile.Name)
	if instanceProfile.ARN != nil {
		fmt.Printf("  ARN: %s\n", *instanceProfile.ARN)
	}
	fmt.Printf("  Associated Role: %s\n", iamRole.Name)
	fmt.Printf("  Path: %s\n", instanceProfile.Path)

	// Step 5: Create Launch Template with IAM role
	fmt.Println("\n--- Step 5: Creating Launch Template with IAM Role ---")
	appSG := usecasescommon.CreateMockSecurityGroup(vpc.ID, "api-sg", "Security group for API servers", gen)
	launchTemplate := usecasescommon.CreateMockLaunchTemplate(
		"api-launch-template",
		"t3.medium",
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
	fmt.Printf("  Security Groups: %v\n", launchTemplate.SecurityGroupIDs)
	fmt.Printf("  Note: IAM Instance Profile (%s) will be attached to instances\n", instanceProfile.Name)

	// Step 6: Create Application Load Balancer
	fmt.Println("\n--- Step 6: Creating Application Load Balancer ---")
	albSG := usecasescommon.CreateMockSecurityGroup(vpc.ID, "api-alb-sg", "Security group for API Load Balancer", gen)
	publicSubnetIDs := []string{}
	for _, subnet := range publicSubnets {
		publicSubnetIDs = append(publicSubnetIDs, subnet.ID)
	}
	alb := usecasescommon.CreateMockLoadBalancer("api-alb", "application", publicSubnetIDs, []string{albSG.ID}, region, gen)
	fmt.Printf("✓ Application Load Balancer created: %s\n", alb.Name)
	if alb.ARN != nil {
		fmt.Printf("  ARN: %s\n", *alb.ARN)
	}
	if alb.DNSName != nil {
		fmt.Printf("  DNS Name: %s\n", *alb.DNSName)
	}
	fmt.Printf("  Subnets: %d across %d AZs\n", len(publicSubnets), len(azs))

	// Step 7: Create Target Group
	fmt.Println("\n--- Step 7: Creating Target Group ---")
	targetGroup := usecasescommon.CreateMockTargetGroup("api-target-group", vpc.ID, "HTTP", 8080, region, gen)
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

	// Step 8: Create Auto Scaling Group with scaling policies
	fmt.Println("\n--- Step 8: Creating Auto Scaling Group ---")
	privateSubnetIDs := []string{}
	for _, subnet := range privateSubnets {
		privateSubnetIDs = append(privateSubnetIDs, subnet.ID)
	}
	asg := usecasescommon.CreateMockAutoScalingGroup(
		"api-asg",
		2,  // min size
		10, // max size
		4,  // desired capacity
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
	fmt.Printf("  Launch Template: %s\n", launchTemplate.Name)
	fmt.Printf("  Subnets: %d across %d AZs\n", len(privateSubnets), len(azs))

	// Step 9: Configure health checks (ELB-based)
	fmt.Println("\n--- Step 9: Configuring Health Checks ---")
	fmt.Printf("✓ Health Check Type: ELB (Load Balancer)\n")
	fmt.Printf("✓ Health Check Grace Period: %d seconds\n", *asg.HealthCheckGracePeriod)
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
	fmt.Printf("  - Target Group ARN: %s\n", *targetGroup.ARN)

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
	fmt.Printf("Note: VPC, subnets, route tables, security groups, IGW, and IAM are free\n")
	fmt.Printf("      Cost varies based on actual instance count (scales dynamically)\n")

	fmt.Println("\n============================================")
	fmt.Println("SCENARIO COMPLETE")
	fmt.Println("============================================")
}
