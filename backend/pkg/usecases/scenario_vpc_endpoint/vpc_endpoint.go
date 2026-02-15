package scenario_vpc_endpoint

import (
	"context"
	"fmt"
	"strings"

	awscomputeservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	awsnetworkingservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/networking"
	awsstorageservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/storage"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
	usecasescommon "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/common"
)

// VPCEndpointRunner demonstrates a VPC architecture with S3 Gateway Endpoint
func VPCEndpointRunner() {
	ctx := context.Background()
	region := usecasescommon.SelectRegion("us-east-1")

	fmt.Println("============================================")
	fmt.Println("SCENARIO: VPC ENDPOINT (GATEWAY) FOR S3")
	fmt.Println("============================================")
	fmt.Printf("Region: %s\n", usecasescommon.FormatRegionName(region))
	fmt.Println("\n[OUTPUT MODE] Domain models + AWS output models")

	// Initialize services
	networkingService := awsnetworkingservice.NewNetworkingService()
	computeService := awscomputeservice.NewComputeService()
	storageService := awsstorageservice.NewStorageService()

	// Step 1: Create VPC
	fmt.Println("\n--- Step 1: Creating VPC ---")
	vpc, _, err := usecasescommon.CreateVPCWithOutput(ctx, networkingService, &domainnetworking.VPC{
		Name:               "endpoint-demo-vpc",
		Region:             region,
		CIDR:               "10.0.0.0/16",
		EnableDNS:          true,
		EnableDNSHostnames: true,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create VPC: %v\n", err)
		return
	}
	fmt.Printf("✓ VPC created: %s (ID: %s)\n", vpc.Name, vpc.ID)

	// Step 2: Create Subnets (Public & Private)
	fmt.Println("\n--- Step 2: Creating Subnets ---")
	az1 := region + "a"

	// Public Subnet
	publicSubnet, _, err := usecasescommon.CreateSubnetWithOutput(ctx, networkingService, &domainnetworking.Subnet{
		Name:             "public-subnet",
		VPCID:            vpc.ID,
		CIDR:             "10.0.1.0/24",
		AvailabilityZone: &az1,
		IsPublic:         true,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create Public Subnet: %v\n", err)
		return
	}
	fmt.Printf("✓ Public Subnet created: %s (ID: %s)\n", publicSubnet.Name, publicSubnet.ID)

	// Private Subnet
	privateSubnet, _, err := usecasescommon.CreateSubnetWithOutput(ctx, networkingService, &domainnetworking.Subnet{
		Name:             "private-subnet",
		VPCID:            vpc.ID,
		CIDR:             "10.0.2.0/24",
		AvailabilityZone: &az1,
		IsPublic:         false,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create Private Subnet: %v\n", err)
		return
	}
	fmt.Printf("✓ Private Subnet created: %s (ID: %s)\n", privateSubnet.Name, privateSubnet.ID)

	// Step 3: Create Internet Gateway & NAT Gateway
	fmt.Println("\n--- Step 3: Creating Gateways ---")

	// Internet Gateway
	igw, _, err := usecasescommon.CreateInternetGatewayWithOutput(ctx, networkingService, &domainnetworking.InternetGateway{
		Name:  "demo-igw",
		VPCID: vpc.ID,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create Internet Gateway: %v\n", err)
		return
	}
	fmt.Printf("✓ Internet Gateway created: %s (ID: %s)\n", igw.Name, igw.ID)

	// Attach IGW
	if err := usecasescommon.AttachInternetGateway(ctx, networkingService, igw.ID, vpc.ID); err != nil {
		fmt.Printf("✗ Failed to attach IGW: %v\n", err)
		return
	}
	fmt.Printf("✓ Internet Gateway attached to VPC\n")

	// Elastic IP for NAT
	eip, _, err := usecasescommon.AllocateElasticIPWithOutput(ctx, networkingService, &domainnetworking.ElasticIP{
		Region: region,
	})
	if err != nil {
		fmt.Printf("✗ Failed to allocate Elastic IP: %v\n", err)
		return
	}

	// NAT Gateway
	nat, _, err := usecasescommon.CreateNATGatewayWithOutput(ctx, networkingService, &domainnetworking.NATGateway{
		Name:         "demo-nat",
		SubnetID:     publicSubnet.ID,
		AllocationID: eip.AllocationID,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create NAT Gateway: %v\n", err)
		return
	}
	fmt.Printf("✓ NAT Gateway created: %s (ID: %s)\n", nat.Name, nat.ID)

	// Step 4: Configure Route Tables
	fmt.Println("\n--- Step 4: Configuring Route Tables ---")

	// Public Route Table
	publicRT, _, err := usecasescommon.CreateRouteTableWithOutput(ctx, networkingService, &domainnetworking.RouteTable{
		Name:  "public-rtb",
		VPCID: vpc.ID,
		Routes: []domainnetworking.Route{
			{DestinationCIDR: "0.0.0.0/0", TargetID: igw.ID, TargetType: "internet_gateway"},
		},
	})
	if err != nil {
		fmt.Printf("✗ Failed to create Public Route Table: %v\n", err)
		return
	}
	usecasescommon.AssociateRouteTable(ctx, networkingService, publicRT.ID, publicSubnet.ID)
	fmt.Printf("✓ Public Route Table configured and associated with Public Subnet\n")

	// Private Route Table
	privateRT, _, err := usecasescommon.CreateRouteTableWithOutput(ctx, networkingService, &domainnetworking.RouteTable{
		Name:  "private-rtb",
		VPCID: vpc.ID,
		Routes: []domainnetworking.Route{
			{DestinationCIDR: "0.0.0.0/0", TargetID: nat.ID, TargetType: "nat_gateway"},
		},
	})
	if err != nil {
		fmt.Printf("✗ Failed to create Private Route Table: %v\n", err)
		return
	}
	usecasescommon.AssociateRouteTable(ctx, networkingService, privateRT.ID, privateSubnet.ID)
	fmt.Printf("✓ Private Route Table configured and associated with Private Subnet\n")

	// Step 5: Create Resources (S3 & EC2)
	fmt.Println("\n--- Step 5: Creating Resources ---")

	// S3 Bucket
	s3Bucket, _, err := usecasescommon.CreateS3BucketWithOutput(ctx, storageService, &domainstorage.S3Bucket{
		Name:   "private-data-bucket-" + generateRandomSuffix(6),
		Region: region,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create S3 bucket: %v\n", err)
		return
	}
	fmt.Printf("✓ S3 Bucket created: %s\n", s3Bucket.Name)

	// Security Group for EC2
	sg, _, err := usecasescommon.CreateSecurityGroupWithOutput(ctx, networkingService, &domainnetworking.SecurityGroup{
		Name:        "private-sg",
		Description: "Allow internal traffic",
		VPCID:       vpc.ID,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create Security Group: %v\n", err)
		return
	}

	// EC2 Instance
	ami := "ami-0c55b159cbfafe1f0" // Amazon Linux 2
	instanceType := "t3.micro"
	instance, _, err := usecasescommon.CreateInstanceWithOutput(ctx, computeService, &domaincompute.Instance{
		Name:             "private-instance",
		Region:           region,
		AMI:              ami,
		InstanceType:     instanceType,
		SubnetID:         privateSubnet.ID,
		AvailabilityZone: &az1,
		SecurityGroupIDs: []string{sg.ID},
	})
	if err != nil {
		fmt.Printf("✗ Failed to create EC2 Instance: %v\n", err)
		return
	}
	fmt.Printf("✓ EC2 Instance created: %s (ID: %s) in Private Subnet\n", instance.Name, instance.ID)

	// Step 6: Create VPC Endpoint (Gateway)
	fmt.Println("\n--- Step 6: Creating VPC Endpoint ---")

	serviceName := fmt.Sprintf("com.amazonaws.%s.s3", region)

	// Create Custom Policy
	policyDoc := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": ["s3:GetObject", "s3:PutObject"],
				"Resource": ["arn:aws:s3:::%s", "arn:aws:s3:::%s/*"]
			}
		]
	}`, s3Bucket.Name, s3Bucket.Name)

	vpce, _, err := usecasescommon.CreateVPCEndpointWithOutput(ctx, networkingService, &domainnetworking.VPCEndpoint{
		Name:          "s3-gateway-endpoint",
		VPCID:         vpc.ID,
		ServiceName:   serviceName,
		Type:          domainnetworking.VPCEndpointTypeGateway,
		RouteTableIDs: []string{privateRT.ID},
		Policy:        policyDoc,
	})
	if err != nil {
		fmt.Printf("✗ Failed to create VPC Endpoint: %v\n", err)
		return
	}

	fmt.Printf("✓ VPC Endpoint created: %s (ID: %s)\n", vpce.Name, vpce.ID)
	fmt.Printf("  Type: %s\n", vpce.Type)
	fmt.Printf("  Service: %s\n", vpce.ServiceName)
	fmt.Printf("  Associated Route Tables: %s\n", strings.Join(vpce.RouteTableIDs, ", "))

	// Summary
	fmt.Println("\n============================================")
	fmt.Println("ARCHITECTURE SUMMARY")
	fmt.Println("============================================")
	fmt.Printf("VPC: %s (%s)\n", vpc.Name, vpc.CIDR)
	fmt.Printf("Subnets:\n")
	fmt.Printf("  - Public: %s (%s) -> IGW\n", publicSubnet.Name, publicSubnet.CIDR)
	fmt.Printf("  - Private: %s (%s) -> NATGW\n", privateSubnet.Name, privateSubnet.CIDR)
	fmt.Printf("Connectivity:\n")
	fmt.Printf("  - Internet Gateway: %s\n", igw.Name)
	fmt.Printf("  - NAT Gateway: %s\n", nat.Name)
	fmt.Printf("  - VPC Endpoint: %s (S3 Gateway) attached to Private Route Table\n", vpce.Name)
	fmt.Printf("Resources:\n")
	fmt.Printf("  - EC2: %s (Private Subnet)\n", instance.Name)
	fmt.Printf("  - S3: %s (Region %s)\n", s3Bucket.Name, s3Bucket.Region)
	fmt.Println("============================================")
	fmt.Println("Traffic Flow:")
	fmt.Println("1. EC2 Instance -> Internet: Routes via NAT Gateway -> Internet Gateway")
	fmt.Println("2. EC2 Instance -> S3 Bucket: Routes via VPC Endpoint (Gateway) -> S3 Service (Traffic stays within AWS network)")
}

func generateRandomSuffix(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[i%len(letters)]) // Simple deterministic suffix for demo
	}
	return sb.String()
}
