package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	awscomputeadapter "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/adapters/compute"
	awsnetworkingadapter "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/adapters/networking"
	awsstorageadapter "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/adapters/storage"
	awscomputeservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	awsnetworkingservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/networking"
	awsstorageservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/storage"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/seeder"
)

// SeedDatabaseWithAdapters seeds the database using adapters and strong types
// This approach uses domain models and adapters instead of interface{} types
func SeedDatabaseWithAdapters() error {
	ctx := context.Background()

	// Connect to database
	if _, err := database.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("🌱 Starting database seeding with adapters...")

	// Ensure reference data exists
	if err := seedReferenceData(ctx); err != nil {
		return fmt.Errorf("failed to seed reference data: %w", err)
	}

	// Seed Resource Constraints
	fmt.Println("🔒 Seeding resource constraints...")
	constraintRepo, err := repository.NewResourceConstraintRepository()
	if err != nil {
		return fmt.Errorf("failed to create constraint repository: %w", err)
	}
	resourceTypeRepo, err := repository.NewResourceTypeRepository()
	if err != nil {
		return fmt.Errorf("failed to create resource type repository: %w", err)
	}

	if err := seeder.SeedResourceConstraints(ctx, constraintRepo, resourceTypeRepo); err != nil {
		return fmt.Errorf("failed to seed resource constraints: %w", err)
	}

	// Seed users
	users, err := seedUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// Seed projects and resources using adapters
	if err := seedScenariosWithAdapters(ctx, users); err != nil {
		return fmt.Errorf("failed to seed scenarios with adapters: %w", err)
	}

	fmt.Println("✅ Database seeding with adapters completed successfully!")
	return nil
}

// seedScenariosWithAdapters creates projects and resources using adapters and strong types
func seedScenariosWithAdapters(ctx context.Context, users []*models.User) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	projectRepo, err := repository.NewProjectRepository(slog.Default())
	if err != nil {
		return err
	}

	resourceRepo, err := repository.NewResourceRepository(slog.Default())
	if err != nil {
		return err
	}

	fmt.Println("\n🏗️  Seeding projects and resources with adapters...")

	// Get resource types
	var vpcType, subnetType, igwType, sgType, ec2Type, lbType, asgType, natType, eipType, lambdaType, s3Type models.ResourceType
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "VPC", "aws").First(&vpcType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "Subnet", "aws").First(&subnetType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "InternetGateway", "aws").First(&igwType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "SecurityGroup", "aws").First(&sgType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "EC2", "aws").First(&ec2Type)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "LoadBalancer", "aws").First(&lbType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "AutoScalingGroup", "aws").First(&asgType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "NATGateway", "aws").First(&natType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "ElasticIP", "aws").First(&eipType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "Lambda", "aws").First(&lambdaType)
	db.WithContext(ctx).Where("name = ? AND cloud_provider = ?", "S3", "aws").First(&s3Type)

	// Get IAC target (Terraform)
	var terraformTarget models.IACTarget
	db.WithContext(ctx).Where("name = ?", "Terraform").First(&terraformTarget)

	// Get dependency type
	var usesDepType models.DependencyType
	db.WithContext(ctx).Where("name = ?", "uses").First(&usesDepType)

	// Create virtual services
	computeService := awscomputeservice.NewComputeService()
	networkingService := awsnetworkingservice.NewNetworkingService()
	storageService := awsstorageservice.NewStorageService()

	// Create adapters
	computeAdapter := awscomputeadapter.NewAWSComputeAdapter(computeService)
	networkingAdapter := awsnetworkingadapter.NewAWSNetworkingAdapter(networkingService)
	storageAdapter := awsstorageadapter.NewAWSStorageAdapter(storageService)

	// Scenario 1: Production Web Application with Strong Types
	if len(users) > 0 {
		project1 := &models.Project{
			UserID:        users[0].ID,
			InfraToolID:   terraformTarget.ID,
			Name:          "Production Web App (Adapter-based)",
			CloudProvider: "aws",
			Region:        "us-east-1",
		}
		if err := projectRepo.Create(ctx, project1); err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
		fmt.Printf("  ✓ Created project: %s\n", project1.Name)

		// Create VPC using adapter
		vpcDomain := &domainnetworking.VPC{
			Name:               "prod-web-vpc",
			Region:             "us-east-1",
			CIDR:               "10.0.0.0/16",
			EnableDNS:          true,
			EnableDNSHostnames: true,
		}
		createdVPC, err := networkingAdapter.CreateVPC(ctx, vpcDomain)
		if err != nil {
			return fmt.Errorf("failed to create VPC via adapter: %w", err)
		}

		// Convert domain VPC to database resource
		vpcConfigJSON, err := domainResourceToJSON(createdVPC)
		if err != nil {
			return fmt.Errorf("failed to marshal VPC config: %w", err)
		}

		vpcResource := &models.Resource{
			ProjectID:      project1.ID,
			ResourceTypeID: vpcType.ID,
			Name:           createdVPC.Name,
			Config:         vpcConfigJSON,
		}
		if err := resourceRepo.Create(ctx, vpcResource); err != nil {
			return fmt.Errorf("failed to save VPC to database: %w", err)
		}
		fmt.Printf("    ✓ Created VPC: %s (ID: %s)\n", createdVPC.Name, createdVPC.ID)

		// Create Public Subnet
		az1 := "us-east-1a"
		publicSubnetDomain := &domainnetworking.Subnet{
			Name:             "prod-public-subnet-1a",
			VPCID:            createdVPC.ID,
			CIDR:             "10.0.1.0/24",
			AvailabilityZone: &az1,
			IsPublic:         true,
		}
		createdPublicSubnet, err := networkingAdapter.CreateSubnet(ctx, publicSubnetDomain)
		if err != nil {
			return fmt.Errorf("failed to create public subnet via adapter: %w", err)
		}

		publicSubnetConfigJSON, err := domainResourceToJSON(createdPublicSubnet)
		if err != nil {
			return fmt.Errorf("failed to marshal subnet config: %w", err)
		}

		publicSubnetResource := &models.Resource{
			ProjectID:      project1.ID,
			ResourceTypeID: subnetType.ID,
			Name:           createdPublicSubnet.Name,
			Config:         publicSubnetConfigJSON,
		}
		if err := resourceRepo.Create(ctx, publicSubnetResource); err != nil {
			return fmt.Errorf("failed to save subnet to database: %w", err)
		}
		if err := resourceRepo.CreateContainment(ctx, vpcResource.ID, publicSubnetResource.ID); err != nil {
			return fmt.Errorf("failed to create containment: %w", err)
		}
		fmt.Printf("    ✓ Created Public Subnet: %s (ID: %s)\n", createdPublicSubnet.Name, createdPublicSubnet.ID)

		// Create Private Subnet
		privateSubnetDomain := &domainnetworking.Subnet{
			Name:             "prod-private-subnet-1a",
			VPCID:            createdVPC.ID,
			CIDR:             "10.0.10.0/24",
			AvailabilityZone: &az1,
			IsPublic:         false,
		}
		createdPrivateSubnet, err := networkingAdapter.CreateSubnet(ctx, privateSubnetDomain)
		if err != nil {
			return fmt.Errorf("failed to create private subnet via adapter: %w", err)
		}

		privateSubnetConfigJSON, err := domainResourceToJSON(createdPrivateSubnet)
		if err != nil {
			return fmt.Errorf("failed to marshal subnet config: %w", err)
		}

		privateSubnetResource := &models.Resource{
			ProjectID:      project1.ID,
			ResourceTypeID: subnetType.ID,
			Name:           createdPrivateSubnet.Name,
			Config:         privateSubnetConfigJSON,
		}
		if err := resourceRepo.Create(ctx, privateSubnetResource); err != nil {
			return fmt.Errorf("failed to save subnet to database: %w", err)
		}
		if err := resourceRepo.CreateContainment(ctx, vpcResource.ID, privateSubnetResource.ID); err != nil {
			return fmt.Errorf("failed to create containment: %w", err)
		}
		fmt.Printf("    ✓ Created Private Subnet: %s (ID: %s)\n", createdPrivateSubnet.Name, createdPrivateSubnet.ID)

		// Create Internet Gateway
		igwDomain := &domainnetworking.InternetGateway{
			Name:  "prod-web-igw",
			VPCID: createdVPC.ID,
		}
		createdIGW, err := networkingAdapter.CreateInternetGateway(ctx, igwDomain)
		if err != nil {
			return fmt.Errorf("failed to create IGW via adapter: %w", err)
		}
		if err := networkingAdapter.AttachInternetGateway(ctx, createdIGW.ID, createdVPC.ID); err != nil {
			return fmt.Errorf("failed to attach IGW: %w", err)
		}

		igwConfigJSON, err := domainResourceToJSON(createdIGW)
		if err != nil {
			return fmt.Errorf("failed to marshal IGW config: %w", err)
		}

		igwResource := &models.Resource{
			ProjectID:      project1.ID,
			ResourceTypeID: igwType.ID,
			Name:           createdIGW.Name,
			Config:         igwConfigJSON,
		}
		if err := resourceRepo.Create(ctx, igwResource); err != nil {
			return fmt.Errorf("failed to save IGW to database: %w", err)
		}
		fmt.Printf("    ✓ Created Internet Gateway: %s (ID: %s)\n", createdIGW.Name, createdIGW.ID)

		// Create Security Group
		sgDomain := &domainnetworking.SecurityGroup{
			Name:        "prod-web-sg",
			VPCID:       createdVPC.ID,
			Description: "Security group for web tier",
		}
		createdSG, err := networkingAdapter.CreateSecurityGroup(ctx, sgDomain)
		if err != nil {
			return fmt.Errorf("failed to create security group via adapter: %w", err)
		}

		sgConfigJSON, err := domainResourceToJSON(createdSG)
		if err != nil {
			return fmt.Errorf("failed to marshal security group config: %w", err)
		}

		sgResource := &models.Resource{
			ProjectID:      project1.ID,
			ResourceTypeID: sgType.ID,
			Name:           createdSG.Name,
			Config:         sgConfigJSON,
		}
		if err := resourceRepo.Create(ctx, sgResource); err != nil {
			return fmt.Errorf("failed to save security group to database: %w", err)
		}
		fmt.Printf("    ✓ Created Security Group: %s (ID: %s)\n", createdSG.Name, createdSG.ID)

		// Create EC2 Instance
		instanceDomain := &domaincompute.Instance{
			Name:             "prod-web-server-1",
			Region:           "us-east-1",
			InstanceType:     "t3.micro",
			AMI:              "ami-0c55b159cbfafe1f0",
			SubnetID:         createdPublicSubnet.ID,
			SecurityGroupIDs: []string{createdSG.ID},
		}
		createdInstance, err := computeAdapter.CreateInstance(ctx, instanceDomain)
		if err != nil {
			return fmt.Errorf("failed to create instance via adapter: %w", err)
		}

		instanceConfigJSON, err := domainResourceToJSON(createdInstance)
		if err != nil {
			return fmt.Errorf("failed to marshal instance config: %w", err)
		}

		instanceResource := &models.Resource{
			ProjectID:      project1.ID,
			ResourceTypeID: ec2Type.ID,
			Name:           createdInstance.Name,
			Config:         instanceConfigJSON,
		}
		if err := resourceRepo.Create(ctx, instanceResource); err != nil {
			return fmt.Errorf("failed to save instance to database: %w", err)
		}
		if err := resourceRepo.CreateContainment(ctx, publicSubnetResource.ID, instanceResource.ID); err != nil {
			return fmt.Errorf("failed to create containment: %w", err)
		}
		fmt.Printf("    ✓ Created EC2 Instance: %s (ID: %s)\n", createdInstance.Name, createdInstance.ID)

		fmt.Printf("    ✓ Created %d resources for Production Web App\n", 6)
	}

	// Scenario 2: High Availability Architecture with Load Balancer
	if len(users) > 1 {
		project2 := &models.Project{
			UserID:        users[1].ID,
			InfraToolID:   terraformTarget.ID,
			Name:          "HA Architecture (Adapter-based)",
			CloudProvider: "aws",
			Region:        "us-east-1",
		}
		if err := projectRepo.Create(ctx, project2); err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
		fmt.Printf("  ✓ Created project: %s\n", project2.Name)

		// Create VPC
		vpc2Domain := &domainnetworking.VPC{
			Name:               "ha-vpc",
			Region:             "us-east-1",
			CIDR:               "10.1.0.0/16",
			EnableDNS:          true,
			EnableDNSHostnames: true,
		}
		createdVPC2, err := networkingAdapter.CreateVPC(ctx, vpc2Domain)
		if err != nil {
			return fmt.Errorf("failed to create VPC via adapter: %w", err)
		}

		vpc2ConfigJSON, err := domainResourceToJSON(createdVPC2)
		if err != nil {
			return fmt.Errorf("failed to marshal VPC config: %w", err)
		}

		vpc2Resource := &models.Resource{
			ProjectID:      project2.ID,
			ResourceTypeID: vpcType.ID,
			Name:           createdVPC2.Name,
			Config:         vpc2ConfigJSON,
		}
		if err := resourceRepo.Create(ctx, vpc2Resource); err != nil {
			return fmt.Errorf("failed to save VPC to database: %w", err)
		}

		// Create Subnets for HA (multiple AZs)
		az1 := "us-east-1a"
		az2 := "us-east-1b"
		publicSubnet1Domain := &domainnetworking.Subnet{
			Name:             "ha-public-subnet-1a",
			VPCID:            createdVPC2.ID,
			CIDR:             "10.1.1.0/24",
			AvailabilityZone: &az1,
			IsPublic:         true,
		}
		createdPublicSubnet1, err := networkingAdapter.CreateSubnet(ctx, publicSubnet1Domain)
		if err != nil {
			return fmt.Errorf("failed to create subnet: %w", err)
		}

		publicSubnet1ConfigJSON, err := domainResourceToJSON(createdPublicSubnet1)
		if err != nil {
			return fmt.Errorf("failed to marshal subnet config: %w", err)
		}

		publicSubnet1Resource := &models.Resource{
			ProjectID:      project2.ID,
			ResourceTypeID: subnetType.ID,
			Name:           createdPublicSubnet1.Name,
			Config:         publicSubnet1ConfigJSON,
		}
		if err := resourceRepo.Create(ctx, publicSubnet1Resource); err != nil {
			return fmt.Errorf("failed to save subnet: %w", err)
		}
		if err := resourceRepo.CreateContainment(ctx, vpc2Resource.ID, publicSubnet1Resource.ID); err != nil {
			return fmt.Errorf("failed to create containment: %w", err)
		}

		publicSubnet2Domain := &domainnetworking.Subnet{
			Name:             "ha-public-subnet-1b",
			VPCID:            createdVPC2.ID,
			CIDR:             "10.1.2.0/24",
			AvailabilityZone: &az2,
			IsPublic:         true,
		}
		createdPublicSubnet2, err := networkingAdapter.CreateSubnet(ctx, publicSubnet2Domain)
		if err != nil {
			return fmt.Errorf("failed to create subnet: %w", err)
		}

		publicSubnet2ConfigJSON, err := domainResourceToJSON(createdPublicSubnet2)
		if err != nil {
			return fmt.Errorf("failed to marshal subnet config: %w", err)
		}

		publicSubnet2Resource := &models.Resource{
			ProjectID:      project2.ID,
			ResourceTypeID: subnetType.ID,
			Name:           createdPublicSubnet2.Name,
			Config:         publicSubnet2ConfigJSON,
		}
		if err := resourceRepo.Create(ctx, publicSubnet2Resource); err != nil {
			return fmt.Errorf("failed to save subnet: %w", err)
		}
		if err := resourceRepo.CreateContainment(ctx, vpc2Resource.ID, publicSubnet2Resource.ID); err != nil {
			return fmt.Errorf("failed to create containment: %w", err)
		}

		// Create Security Group
		sg2Domain := &domainnetworking.SecurityGroup{
			Name:        "ha-web-sg",
			VPCID:       createdVPC2.ID,
			Description: "Security group for HA web tier",
		}
		createdSG2, err := networkingAdapter.CreateSecurityGroup(ctx, sg2Domain)
		if err != nil {
			return fmt.Errorf("failed to create security group: %w", err)
		}

		sg2ConfigJSON, err := domainResourceToJSON(createdSG2)
		if err != nil {
			return fmt.Errorf("failed to marshal security group config: %w", err)
		}

		sg2Resource := &models.Resource{
			ProjectID:      project2.ID,
			ResourceTypeID: sgType.ID,
			Name:           createdSG2.Name,
			Config:         sg2ConfigJSON,
		}
		if err := resourceRepo.Create(ctx, sg2Resource); err != nil {
			return fmt.Errorf("failed to save security group: %w", err)
		}

		// Create Load Balancer
		lbDomain := &domaincompute.LoadBalancer{
			Name:             "ha-alb",
			Region:           "us-east-1",
			Type:             domaincompute.LoadBalancerTypeApplication,
			Internal:         false,
			SecurityGroupIDs: []string{createdSG2.ID},
			SubnetIDs:        []string{createdPublicSubnet1.ID, createdPublicSubnet2.ID},
		}
		createdLB, err := computeAdapter.CreateLoadBalancer(ctx, lbDomain)
		if err != nil {
			return fmt.Errorf("failed to create load balancer via adapter: %w", err)
		}

		lbConfigJSON, err := domainResourceToJSON(createdLB)
		if err != nil {
			return fmt.Errorf("failed to marshal load balancer config: %w", err)
		}

		lbResource := &models.Resource{
			ProjectID:      project2.ID,
			ResourceTypeID: lbType.ID,
			Name:           createdLB.Name,
			Config:         lbConfigJSON,
		}
		if err := resourceRepo.Create(ctx, lbResource); err != nil {
			return fmt.Errorf("failed to save load balancer: %w", err)
		}
		fmt.Printf("    ✓ Created Load Balancer: %s (ID: %s)\n", createdLB.Name, createdLB.ID)

		// Create Launch Template for ASG
		launchTemplateDomain := &domaincompute.LaunchTemplate{
			Name:             "ha-asg-launch-template",
			Region:           "us-east-1",
			ImageID:          "ami-0c55b159cbfafe1f0",
			InstanceType:     "t3.small",
			SecurityGroupIDs: []string{createdSG2.ID},
		}
		createdLaunchTemplate, err := computeAdapter.CreateLaunchTemplate(ctx, launchTemplateDomain)
		if err != nil {
			return fmt.Errorf("failed to create launch template via adapter: %w", err)
		}

		// Create Auto Scaling Group
		latestVersion := "$Latest"
		asgDomain := &domaincompute.AutoScalingGroup{
			Name:              "ha-asg",
			Region:            "us-east-1",
			MinSize:           2,
			MaxSize:           6,
			DesiredCapacity:   intPtr(3),
			VPCZoneIdentifier: []string{createdPublicSubnet1.ID, createdPublicSubnet2.ID},
			LaunchTemplate: &domaincompute.LaunchTemplateSpecification{
				ID:      createdLaunchTemplate.ID,
				Version: &latestVersion,
			},
			HealthCheckType: domaincompute.AutoScalingGroupHealthCheckTypeEC2,
		}
		createdASG, err := computeAdapter.CreateAutoScalingGroup(ctx, asgDomain)
		if err != nil {
			return fmt.Errorf("failed to create ASG via adapter: %w", err)
		}

		asgConfigJSON, err := domainResourceToJSON(createdASG)
		if err != nil {
			return fmt.Errorf("failed to marshal ASG config: %w", err)
		}

		asgResource := &models.Resource{
			ProjectID:      project2.ID,
			ResourceTypeID: asgType.ID,
			Name:           createdASG.Name,
			Config:         asgConfigJSON,
		}
		if err := resourceRepo.Create(ctx, asgResource); err != nil {
			return fmt.Errorf("failed to save ASG: %w", err)
		}
		fmt.Printf("    ✓ Created Auto Scaling Group: %s (ID: %s)\n", createdASG.Name, createdASG.ID)

		// Allocate Elastic IP for NAT Gateway
		eipDomain := &domainnetworking.ElasticIP{
			Region: "us-east-1",
		}
		createdEIP, err := networkingAdapter.AllocateElasticIP(ctx, eipDomain)
		if err != nil {
			return fmt.Errorf("failed to allocate Elastic IP via adapter: %w", err)
		}

		// Create NAT Gateway with Elastic IP allocation ID
		natDomain := &domainnetworking.NATGateway{
			Name:         "ha-nat-gateway",
			SubnetID:     createdPublicSubnet1.ID,
			AllocationID: &createdEIP.ID, // Use the Elastic IP allocation ID
		}
		createdNAT, err := networkingAdapter.CreateNATGateway(ctx, natDomain)
		if err != nil {
			return fmt.Errorf("failed to create NAT gateway via adapter: %w", err)
		}

		natConfigJSON, err := domainResourceToJSON(createdNAT)
		if err != nil {
			return fmt.Errorf("failed to marshal NAT gateway config: %w", err)
		}

		natResource := &models.Resource{
			ProjectID:      project2.ID,
			ResourceTypeID: natType.ID,
			Name:           createdNAT.Name,
			Config:         natConfigJSON,
		}
		if err := resourceRepo.Create(ctx, natResource); err != nil {
			return fmt.Errorf("failed to save NAT gateway: %w", err)
		}
		fmt.Printf("    ✓ Created NAT Gateway: %s (ID: %s)\n", createdNAT.Name, createdNAT.ID)

		fmt.Printf("    ✓ Created %d resources for HA Architecture\n", 6)
	}

	// Scenario 3: Serverless Lambda + S3
	if len(users) > 2 {
		project3 := &models.Project{
			UserID:        users[2].ID,
			InfraToolID:   terraformTarget.ID,
			Name:          "Serverless App (Adapter-based)",
			CloudProvider: "aws",
			Region:        "us-east-1",
		}
		if err := projectRepo.Create(ctx, project3); err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
		fmt.Printf("  ✓ Created project: %s\n", project3.Name)

		// Create Lambda Function
		runtime := "python3.9"
		handler := "index.handler"
		memorySize := 128
		timeout := 30
		s3Bucket := "lambda-code-bucket"
		s3Key := "data-processor.zip"
		lambdaDomain := &domaincompute.LambdaFunction{
			FunctionName: "data-processor",
			Region:       "us-east-1",
			RoleARN:      "arn:aws:iam::123456789012:role/lambda-execution-role",
			S3Bucket:     &s3Bucket,
			S3Key:        &s3Key,
			Runtime:      &runtime,
			Handler:      &handler,
			MemorySize:   &memorySize,
			Timeout:      &timeout,
		}
		createdLambda, err := computeAdapter.CreateLambdaFunction(ctx, lambdaDomain)
		if err != nil {
			return fmt.Errorf("failed to create Lambda via adapter: %w", err)
		}

		lambdaConfigJSON, err := domainResourceToJSON(createdLambda)
		if err != nil {
			return fmt.Errorf("failed to marshal Lambda config: %w", err)
		}

		lambdaResource := &models.Resource{
			ProjectID:      project3.ID,
			ResourceTypeID: lambdaType.ID,
			Name:           createdLambda.FunctionName,
			Config:         lambdaConfigJSON,
		}
		if err := resourceRepo.Create(ctx, lambdaResource); err != nil {
			return fmt.Errorf("failed to save Lambda: %w", err)
		}
		fmt.Printf("    ✓ Created Lambda Function: %s (ARN: %s)\n", createdLambda.FunctionName, getARN(createdLambda.ARN))

		// Create S3 Bucket
		bucketName := "serverless-data-bucket"
		s3Domain := &domainstorage.S3Bucket{
			Name:   bucketName,
			Region: "us-east-1",
		}
		createdS3, err := storageAdapter.CreateS3Bucket(ctx, s3Domain)
		if err != nil {
			return fmt.Errorf("failed to create S3 bucket via adapter: %w", err)
		}

		s3ConfigJSON, err := domainResourceToJSON(createdS3)
		if err != nil {
			return fmt.Errorf("failed to marshal S3 config: %w", err)
		}

		s3Resource := &models.Resource{
			ProjectID:      project3.ID,
			ResourceTypeID: s3Type.ID,
			Name:           createdS3.Name,
			Config:         s3ConfigJSON,
		}
		if err := resourceRepo.Create(ctx, s3Resource); err != nil {
			return fmt.Errorf("failed to save S3 bucket: %w", err)
		}
		fmt.Printf("    ✓ Created S3 Bucket: %s (ID: %s)\n", createdS3.Name, createdS3.ID)

		// Create dependency: Lambda uses S3
		dependency := &models.ResourceDependency{
			FromResourceID:   lambdaResource.ID,
			ToResourceID:     s3Resource.ID,
			DependencyTypeID: usesDepType.ID,
		}
		if err := resourceRepo.CreateDependency(ctx, dependency); err != nil {
			return fmt.Errorf("failed to create dependency: %w", err)
		}
		fmt.Printf("    ✓ Created dependency: Lambda uses S3\n")

		fmt.Printf("    ✓ Created %d resources for Serverless App\n", 2)
	}

	return nil
}

// domainResourceToJSON converts a domain resource to JSON for database storage
func domainResourceToJSON(resource interface{}) ([]byte, error) {
	return json.Marshal(resource)
}

// getARN safely extracts ARN string from pointer
func getARN(arn *string) string {
	if arn == nil {
		return "N/A"
	}
	return *arn
}

// intPtr returns a pointer to an int
func intPtr(i int) *int {
	return &i
}
