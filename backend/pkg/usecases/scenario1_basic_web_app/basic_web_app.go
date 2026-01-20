package scenario1_basic_web_app

import (
	"context"
	"fmt"
	"time"

	awspricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	usecasescommon "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/common"
)

// BasicWebAppRunner demonstrates a simple 3-tier web application architecture
func BasicWebAppRunner() {
	ctx := context.Background()
	region := usecasescommon.SelectRegion("us-east-1")

	fmt.Println("============================================")
	fmt.Println("SCENARIO 1: BASIC WEB APPLICATION")
	fmt.Println("============================================")
	fmt.Printf("Region: %s\n", usecasescommon.FormatRegionName(region))
	fmt.Println("\n[MOCK MODE] All resources are simulated - no AWS SDK calls")

	// Initialize mock ID generator
	gen := usecasescommon.NewMockIDGenerator()

	// Step 1: Create VPC
	fmt.Println("\n--- Step 1: Creating VPC ---")
	vpc := usecasescommon.CreateMockVPC(region, "web-app-vpc", "10.0.0.0/16", gen)
	fmt.Printf("✓ VPC created: %s (%s)\n", vpc.Name, vpc.ID)
	fmt.Printf("  CIDR: %s\n", vpc.CIDR)
	if vpc.ARN != nil {
		fmt.Printf("  ARN: %s\n", *vpc.ARN)
	}

	// Step 2: Get availability zones
	fmt.Println("\n--- Step 2: Getting Availability Zones ---")
	azs := usecasescommon.GetDefaultAvailabilityZones(region)
	fmt.Printf("✓ Available AZs: %v\n", azs)

	// Step 3: Create public subnets
	fmt.Println("\n--- Step 3: Creating Public Subnets ---")
	publicSubnet1 := usecasescommon.CreateMockSubnet(vpc.ID, "public-subnet-1", "10.0.1.0/24", azs[0], gen)
	publicSubnet2 := usecasescommon.CreateMockSubnet(vpc.ID, "public-subnet-2", "10.0.2.0/24", azs[1], gen)
	publicSubnets := []*domainnetworking.Subnet{publicSubnet1, publicSubnet2}
	fmt.Printf("✓ Created %d public subnets:\n", len(publicSubnets))
	for i, subnet := range publicSubnets {
		fmt.Printf("  %d. %s (%s) in %s\n", i+1, subnet.Name, subnet.ID, *subnet.AvailabilityZone)
	}

	// Step 4: Create private subnets
	fmt.Println("\n--- Step 4: Creating Private Subnets ---")
	privateSubnet1 := usecasescommon.CreateMockSubnet(vpc.ID, "private-subnet-1", "10.0.10.0/24", azs[0], gen)
	privateSubnet2 := usecasescommon.CreateMockSubnet(vpc.ID, "private-subnet-2", "10.0.11.0/24", azs[1], gen)
	privateSubnets := []*domainnetworking.Subnet{privateSubnet1, privateSubnet2}
	fmt.Printf("✓ Created %d private subnets:\n", len(privateSubnets))
	for i, subnet := range privateSubnets {
		fmt.Printf("  %d. %s (%s) in %s\n", i+1, subnet.Name, subnet.ID, *subnet.AvailabilityZone)
	}

	// Step 5: Create Internet Gateway
	fmt.Println("\n--- Step 5: Creating Internet Gateway ---")
	igw := usecasescommon.CreateMockInternetGateway(vpc.ID, "web-app-igw", gen)
	fmt.Printf("✓ Internet Gateway created: %s (%s)\n", igw.Name, igw.ID)
	if igw.ARN != nil {
		fmt.Printf("  ARN: %s\n", *igw.ARN)
	}
	fmt.Printf("  Attached to VPC: %s\n", igw.VPCID)

	// Step 6: Create public route table
	fmt.Println("\n--- Step 6: Creating Public Route Table ---")
	publicRT := usecasescommon.CreateMockRouteTable(vpc.ID, "public-route-table", gen)
	fmt.Printf("✓ Public Route Table created: %s (%s)\n", publicRT.Name, publicRT.ID)
	fmt.Printf("  Route: 0.0.0.0/0 -> %s (Internet Gateway)\n", igw.ID)

	// Step 7: Create private route table
	fmt.Println("\n--- Step 7: Creating Private Route Table ---")
	privateRT := usecasescommon.CreateMockRouteTable(vpc.ID, "private-route-table", gen)
	fmt.Printf("✓ Private Route Table created: %s (%s)\n", privateRT.Name, privateRT.ID)
	fmt.Printf("  Note: No internet gateway route (private subnets)\n")

	// Step 8: Associate route tables with subnets
	fmt.Println("\n--- Step 8: Associating Route Tables with Subnets ---")
	fmt.Printf("✓ Associated public route table with public subnets\n")
	for _, subnet := range publicSubnets {
		fmt.Printf("  - %s -> %s\n", subnet.Name, publicRT.Name)
	}
	fmt.Printf("✓ Associated private route table with private subnets\n")
	for _, subnet := range privateSubnets {
		fmt.Printf("  - %s -> %s\n", subnet.Name, privateRT.Name)
	}

	// Step 9: Create security groups
	fmt.Println("\n--- Step 9: Creating Security Groups ---")
	webSG := usecasescommon.CreateMockSecurityGroup(vpc.ID, "web-sg", "Security group for web tier", gen)
	appSG := usecasescommon.CreateMockSecurityGroup(vpc.ID, "app-sg", "Security group for application tier", gen)
	dbSG := usecasescommon.CreateMockSecurityGroup(vpc.ID, "db-sg", "Security group for database tier", gen)
	securityGroups := map[string]*domainnetworking.SecurityGroup{
		"web": webSG,
		"app": appSG,
		"db":  dbSG,
	}
	fmt.Printf("✓ Created %d security groups:\n", len(securityGroups))
	for tier, sg := range securityGroups {
		fmt.Printf("  - %s: %s (%s) - %s\n", tier, sg.Name, sg.ID, sg.Description)
	}

	// Step 10: Create EC2 instances in public subnets (web tier)
	fmt.Println("\n--- Step 10: Creating EC2 Instances (Web Tier) ---")
	webInstances := []*domaincompute.Instance{}
	for i, subnet := range publicSubnets {
		instance := usecasescommon.CreateMockEC2Instance(
			fmt.Sprintf("web-server-%d", i+1),
			"t3.micro",
			subnet.ID,
			webSG.ID,
			region,
			*subnet.AvailabilityZone,
			gen,
		)
		webInstances = append(webInstances, instance)
		fmt.Printf("✓ Created: %s (%s) in %s\n", instance.Name, instance.ID, *subnet.AvailabilityZone)
	}

	// Step 11: Create EC2 instances in private subnets (app tier)
	fmt.Println("\n--- Step 11: Creating EC2 Instances (Application Tier) ---")
	appInstances := []*domaincompute.Instance{}
	for i, subnet := range privateSubnets {
		instance := usecasescommon.CreateMockEC2Instance(
			fmt.Sprintf("app-server-%d", i+1),
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

	// Step 12: Display architecture summary
	fmt.Println("\n============================================")
	fmt.Println("ARCHITECTURE SUMMARY")
	fmt.Println("============================================")
	fmt.Printf("VPC: %s (%s)\n", vpc.Name, vpc.CIDR)
	fmt.Printf("Internet Gateway: %s\n", igw.Name)
	fmt.Printf("\nNetworking:\n")
	fmt.Printf("  Public Subnets: %d\n", len(publicSubnets))
	fmt.Printf("  Private Subnets: %d\n", len(privateSubnets))
	fmt.Printf("  Route Tables: 2 (public, private)\n")
	fmt.Printf("  Security Groups: %d\n", len(securityGroups))
	fmt.Printf("\nCompute:\n")
	fmt.Printf("  Web Tier Instances: %d (t3.micro)\n", len(webInstances))
	fmt.Printf("  App Tier Instances: %d (t3.small)\n", len(appInstances))
	fmt.Printf("  Total Instances: %d\n", len(webInstances)+len(appInstances))

	// Step 13: Calculate estimated costs
	fmt.Println("\n============================================")
	fmt.Println("COST ESTIMATION (30 days)")
	fmt.Println("============================================")
	pricingService := awspricing.NewAWSPricingService()
	duration := 30 * 24 * time.Hour // 30 days

	// Calculate VPC cost (VPC itself is free, but we'll show other resources)
	fmt.Println("\nResource Costs:")
	totalCost := 0.0

	// EC2 instances
	for _, instance := range webInstances {
		res := &resource.Resource{
			Type:     resource.ResourceType{Name: "ec2_instance"},
			Provider: "aws",
			Region:   region,
			Metadata: map[string]interface{}{
				"instance_type": instance.InstanceType,
			},
		}
		estimate, err := pricingService.EstimateCost(ctx, res, duration)
		if err == nil {
			fmt.Printf("  %s (%s): $%.2f\n", instance.Name, instance.InstanceType, estimate.TotalCost)
			totalCost += estimate.TotalCost
		}
	}
	for _, instance := range appInstances {
		res := &resource.Resource{
			Type:     resource.ResourceType{Name: "ec2_instance"},
			Provider: "aws",
			Region:   region,
			Metadata: map[string]interface{}{
				"instance_type": instance.InstanceType,
			},
		}
		estimate, err := pricingService.EstimateCost(ctx, res, duration)
		if err == nil {
			fmt.Printf("  %s (%s): $%.2f\n", instance.Name, instance.InstanceType, estimate.TotalCost)
			totalCost += estimate.TotalCost
		}
	}

	fmt.Printf("\nTotal Estimated Cost (30 days): $%.2f\n", totalCost)
	fmt.Printf("Note: VPC, subnets, route tables, and security groups are free\n")
	fmt.Printf("      Internet Gateway is free when attached to VPC\n")

	fmt.Println("\n============================================")
	fmt.Println("SCENARIO COMPLETE")
	fmt.Println("============================================")
}
