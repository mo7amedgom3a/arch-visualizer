package scenario1_basic_web_app

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

// BasicWebAppRunner demonstrates a simple 3-tier web application architecture
func BasicWebAppRunner() {
	ctx := context.Background()
	region := usecasescommon.SelectRegion("us-east-1")

	fmt.Println("============================================")
	fmt.Println("SCENARIO 1: BASIC WEB APPLICATION")
	fmt.Println("============================================")
	fmt.Printf("Region: %s\n", usecasescommon.FormatRegionName(region))
	fmt.Println("\n[OUTPUT MODE] Domain models + AWS output models")

	// Initialize virtual services
	networkingService := awsnetworkingservice.NewNetworkingService()
	computeService := awscomputeservice.NewComputeService()

	// Step 1: Create VPC
	fmt.Println("\n--- Step 1: Creating VPC ---")
	vpcDomain := &domainnetworking.VPC{
		Name:               "web-app-vpc",
		Region:             region,
		CIDR:               "10.0.0.0/16",
		EnableDNS:          true,
		EnableDNSHostnames: true,
	}
	vpc, vpcOutput, err := usecasescommon.CreateVPCWithOutput(ctx, networkingService, vpcDomain)
	if err != nil {
		fmt.Printf("✗ Failed to create VPC: %v\n", err)
		return
	}
	fmt.Printf("✓ VPC created: %s (%s)\n", vpcOutput.Name, vpcOutput.ID)
	fmt.Printf("  CIDR: %s\n", vpcOutput.CIDR)
	if vpcOutput.ARN != "" {
		fmt.Printf("  ARN: %s\n", vpcOutput.ARN)
	}

	// Step 2: Get availability zones
	fmt.Println("\n--- Step 2: Getting Availability Zones ---")
	azs := usecasescommon.GetDefaultAvailabilityZones(region)
	fmt.Printf("✓ Available AZs: %v\n", azs)

	// Step 3: Create public subnets
	fmt.Println("\n--- Step 3: Creating Public Subnets ---")
	publicSubnets := []*domainnetworking.Subnet{}
	publicSubnetConfigs := []struct {
		Name string
		CIDR string
		AZ   string
	}{
		{Name: "public-subnet-1", CIDR: "10.0.1.0/24", AZ: azs[0]},
		{Name: "public-subnet-2", CIDR: "10.0.2.0/24", AZ: azs[1]},
	}
	for _, cfg := range publicSubnetConfigs {
		subnet, subnetOutput, err := usecasescommon.CreateSubnetWithOutput(ctx, networkingService, &domainnetworking.Subnet{
			Name:             cfg.Name,
			VPCID:            vpc.ID,
			CIDR:             cfg.CIDR,
			AvailabilityZone: &cfg.AZ,
			IsPublic:         true,
		})
		if err != nil {
			fmt.Printf("✗ Failed to create public subnet %s: %v\n", cfg.Name, err)
			return
		}
		publicSubnets = append(publicSubnets, subnet)
		_ = subnetOutput
	}
	fmt.Printf("✓ Created %d public subnets:\n", len(publicSubnets))
	for i, subnet := range publicSubnets {
		fmt.Printf("  %d. %s (%s) in %s\n", i+1, subnet.Name, subnet.ID, *subnet.AvailabilityZone)
	}

	// Step 4: Create private subnets
	fmt.Println("\n--- Step 4: Creating Private Subnets ---")
	privateSubnets := []*domainnetworking.Subnet{}
	privateSubnetConfigs := []struct {
		Name string
		CIDR string
		AZ   string
	}{
		{Name: "private-subnet-1", CIDR: "10.0.10.0/24", AZ: azs[0]},
		{Name: "private-subnet-2", CIDR: "10.0.11.0/24", AZ: azs[1]},
	}
	for _, cfg := range privateSubnetConfigs {
		subnet, subnetOutput, err := usecasescommon.CreateSubnetWithOutput(ctx, networkingService, &domainnetworking.Subnet{
			Name:             cfg.Name,
			VPCID:            vpc.ID,
			CIDR:             cfg.CIDR,
			AvailabilityZone: &cfg.AZ,
			IsPublic:         false,
		})
		if err != nil {
			fmt.Printf("✗ Failed to create private subnet %s: %v\n", cfg.Name, err)
			return
		}
		privateSubnets = append(privateSubnets, subnet)
		_ = subnetOutput
	}
	fmt.Printf("✓ Created %d private subnets:\n", len(privateSubnets))
	for i, subnet := range privateSubnets {
		fmt.Printf("  %d. %s (%s) in %s\n", i+1, subnet.Name, subnet.ID, *subnet.AvailabilityZone)
	}

	// Step 5: Create Internet Gateway
	fmt.Println("\n--- Step 5: Creating Internet Gateway ---")
	igw, igwOutput, err := usecasescommon.CreateInternetGatewayWithOutput(ctx, networkingService, &domainnetworking.InternetGateway{
		Name:  "web-app-igw",
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
	if igwOutput.ARN != "" {
		fmt.Printf("  ARN: %s\n", igwOutput.ARN)
	}
	fmt.Printf("  Attached to VPC: %s\n", igw.VPCID)

	// Step 6: Create public route table
	fmt.Println("\n--- Step 6: Creating Public Route Table ---")
	publicRT, publicRTOutput, err := usecasescommon.CreateRouteTableWithOutput(ctx, networkingService, &domainnetworking.RouteTable{
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
	fmt.Printf("✓ Public Route Table created: %s (%s)\n", publicRTOutput.Name, publicRTOutput.ID)
	fmt.Printf("  Route: 0.0.0.0/0 -> %s (Internet Gateway)\n", igw.ID)

	// Step 7: Create private route table
	fmt.Println("\n--- Step 7: Creating Private Route Table ---")
	privateRT, privateRTOutput, err := usecasescommon.CreateRouteTableWithOutput(ctx, networkingService, &domainnetworking.RouteTable{
		Name:  "private-route-table",
		VPCID: vpc.ID,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create private route table: %v\n", err)
		return
	}
	fmt.Printf("✓ Private Route Table created: %s (%s)\n", privateRTOutput.Name, privateRTOutput.ID)
	fmt.Printf("  Note: No internet gateway route (private subnets)\n")

	// Step 8: Associate route tables with subnets
	fmt.Println("\n--- Step 8: Associating Route Tables with Subnets ---")
	for _, subnet := range publicSubnets {
		if err := usecasescommon.AssociateRouteTable(ctx, networkingService, publicRT.ID, subnet.ID); err != nil {
			fmt.Printf("✗ Failed to associate public route table: %v\n", err)
			return
		}
		fmt.Printf("  - %s -> %s\n", subnet.Name, publicRT.Name)
	}
	for _, subnet := range privateSubnets {
		if err := usecasescommon.AssociateRouteTable(ctx, networkingService, privateRT.ID, subnet.ID); err != nil {
			fmt.Printf("✗ Failed to associate private route table: %v\n", err)
			return
		}
		fmt.Printf("  - %s -> %s\n", subnet.Name, privateRT.Name)
	}

	// Step 9: Create security groups
	fmt.Println("\n--- Step 9: Creating Security Groups ---")
	webSG, webSGOutput, err := usecasescommon.CreateSecurityGroupWithOutput(ctx, networkingService, &domainnetworking.SecurityGroup{
		Name:        "web-sg",
		VPCID:       vpc.ID,
		Description: "Security group for web tier",
	})
	if err != nil {
		fmt.Printf("✗ Failed to create web security group: %v\n", err)
		return
	}
	appSG, appSGOutput, err := usecasescommon.CreateSecurityGroupWithOutput(ctx, networkingService, &domainnetworking.SecurityGroup{
		Name:        "app-sg",
		VPCID:       vpc.ID,
		Description: "Security group for application tier",
	})
	if err != nil {
		fmt.Printf("✗ Failed to create app security group: %v\n", err)
		return
	}
	dbSG, dbSGOutput, err := usecasescommon.CreateSecurityGroupWithOutput(ctx, networkingService, &domainnetworking.SecurityGroup{
		Name:        "db-sg",
		VPCID:       vpc.ID,
		Description: "Security group for database tier",
	})
	if err != nil {
		fmt.Printf("✗ Failed to create db security group: %v\n", err)
		return
	}
	securityGroups := map[string]*domainnetworking.SecurityGroup{
		"web": webSG,
		"app": appSG,
		"db":  dbSG,
	}
	fmt.Printf("✓ Created %d security groups:\n", len(securityGroups))
	fmt.Printf("  - web: %s (%s) - %s\n", webSGOutput.Name, webSGOutput.ID, webSGOutput.Description)
	fmt.Printf("  - app: %s (%s) - %s\n", appSGOutput.Name, appSGOutput.ID, appSGOutput.Description)
	fmt.Printf("  - db: %s (%s) - %s\n", dbSGOutput.Name, dbSGOutput.ID, dbSGOutput.Description)

	// Step 10: Create EC2 instances in public subnets (web tier)
	fmt.Println("\n--- Step 10: Creating EC2 Instances (Web Tier) ---")
	webInstances := []*domaincompute.Instance{}
	for i, subnet := range publicSubnets {
		instance, instanceOutput, err := usecasescommon.CreateInstanceWithOutput(ctx, computeService, &domaincompute.Instance{
			Name:             fmt.Sprintf("web-server-%d", i+1),
			Region:           region,
			AvailabilityZone: subnet.AvailabilityZone,
			InstanceType:     "t3.micro",
			AMI:              "ami-0c55b159cbfafe1f0",
			SubnetID:         subnet.ID,
			SecurityGroupIDs: []string{webSG.ID},
		})
		if err != nil {
			fmt.Printf("✗ Failed to create web instance: %v\n", err)
			return
		}
		webInstances = append(webInstances, instance)
		fmt.Printf("✓ Created: %s (%s) in %s\n", instanceOutput.Name, instanceOutput.ID, *subnet.AvailabilityZone)
	}

	// Step 11: Create EC2 instances in private subnets (app tier)
	fmt.Println("\n--- Step 11: Creating EC2 Instances (Application Tier) ---")
	appInstances := []*domaincompute.Instance{}
	for i, subnet := range privateSubnets {
		instance, instanceOutput, err := usecasescommon.CreateInstanceWithOutput(ctx, computeService, &domaincompute.Instance{
			Name:             fmt.Sprintf("app-server-%d", i+1),
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
