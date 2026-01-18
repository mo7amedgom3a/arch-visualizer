package ec2

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsec2model "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2"
	awsmodel "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/instance_types"
	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
)

// EC2Runner demonstrates EC2 instance operations and instance type management
func EC2Runner() {
	ctx := context.Background()

	// Initialize AWS client
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		fmt.Printf("Error creating AWS client: %v\n", err)
		return
	}

	fmt.Println("============================================")
	fmt.Println("EC2 INSTANCE TYPES AND CATEGORIES")
	fmt.Println("============================================")

	// Initialize instance type service
	instanceTypeService, err := awssdk.NewInstanceTypeService(client)
	if err != nil {
		fmt.Printf("Error creating instance type service: %v\n", err)
		return
	}

	region := client.GetRegion()
	fmt.Printf("\nRegion: %s\n", region)

	// Display instance categories
	displayInstanceCategories()

	// Display instance types by category
	displayInstanceTypesByCategory(ctx, instanceTypeService, region)

	// Display free tier instances
	displayFreeTierInstances(ctx, instanceTypeService, region)

	// Display specific instance type details
	displayInstanceTypeDetails(ctx, instanceTypeService, region, "t3.micro")
	displayInstanceTypeDetails(ctx, instanceTypeService, region, "m5.large")
	displayInstanceTypeDetails(ctx, instanceTypeService, region, "c5.xlarge")

	// Example: Create EC2 instance configuration
	fmt.Println("\n============================================")
	fmt.Println("EC2 INSTANCE CONFIGURATION EXAMPLE")
	fmt.Println("============================================")
	createInstanceExample(ctx, client, instanceTypeService, region)
}

// displayInstanceCategories displays all available instance categories
func displayInstanceCategories() {
	fmt.Println("\n--- Instance Categories ---")
	categories := awsmodel.AllCategories()
	for i, category := range categories {
		fmt.Printf("%d. %s - %s\n", i+1, category.String(), category.GetDescription())
	}
}

// displayInstanceTypesByCategory displays instance types grouped by category
func displayInstanceTypesByCategory(ctx context.Context, service *awssdk.InstanceTypeService, region string) {
	fmt.Println("\n--- Instance Types by Category ---")

	categories := awsmodel.AllCategories()
	for _, category := range categories {
		types, err := service.ListByCategory(ctx, category, region)
		if err != nil {
			fmt.Printf("Error listing %s instances: %v\n", category.String(), err)
			continue
		}

		if len(types) > 0 {
			fmt.Printf("\n%s (%d types):\n", category.GetDescription(), len(types))
			// Limit display to first 10 types per category
			maxDisplay := 10
			if len(types) < maxDisplay {
				maxDisplay = len(types)
			}
			for i := 0; i < maxDisplay; i++ {
				info := types[i]
				fmt.Printf("  - %s: %d vCPU, %.2f GiB RAM", info.Name, info.VCPU, info.MemoryGiB)
				if info.HasLocalStorage && info.LocalStorageSizeGiB != nil {
					fmt.Printf(", %.2f GiB local storage", *info.LocalStorageSizeGiB)
				}
				fmt.Printf("\n")
			}
			if len(types) > maxDisplay {
				fmt.Printf("  ... and %d more\n", len(types)-maxDisplay)
			}
		}
	}
}

// displayFreeTierInstances displays free tier eligible instance types
func displayFreeTierInstances(ctx context.Context, service *awssdk.InstanceTypeService, region string) {
	fmt.Println("\n--- Free Tier Eligible Instances ---")
	freeTier, err := service.ListFreeTier(ctx, region)
	if err != nil {
		fmt.Printf("Error listing free tier instances: %v\n", err)
		return
	}

	if len(freeTier) == 0 {
		fmt.Println("No free tier instances found")
		return
	}

	for _, info := range freeTier {
		fmt.Printf("  - %s: %d vCPU, %.2f GiB RAM, Category: %s\n",
			info.Name, info.VCPU, info.MemoryGiB, info.Category.String())
	}
}

// displayInstanceTypeDetails displays detailed information about a specific instance type
func displayInstanceTypeDetails(ctx context.Context, service *awssdk.InstanceTypeService, region, instanceType string) {
	fmt.Printf("\n--- Instance Type Details: %s ---\n", instanceType)
	info, err := service.GetInstanceType(ctx, instanceType, region)
	if err != nil {
		fmt.Printf("Error getting instance type %s: %v\n", instanceType, err)
		return
	}

	fmt.Printf("Name: %s\n", info.Name)
	fmt.Printf("Category: %s (%s)\n", info.Category.String(), info.Category.GetDescription())
	fmt.Printf("vCPU: %d\n", info.VCPU)
	fmt.Printf("Memory: %.2f GiB (%.2f MB)\n", info.MemoryGiB, info.GetMemoryMB())
	fmt.Printf("Storage Type: %s\n", info.StorageType)
	fmt.Printf("Has Local Storage: %v\n", info.HasLocalStorage)
	if info.HasLocalStorage && info.LocalStorageSizeGiB != nil {
		fmt.Printf("Local Storage Size: %.2f GiB\n", *info.LocalStorageSizeGiB)
	}
	fmt.Printf("Max Network: %.2f Gbps\n", info.MaxNetworkGbps)
	if info.EBSBandwidthGbps != nil {
		fmt.Printf("EBS Bandwidth: %.2f Gbps\n", *info.EBSBandwidthGbps)
	}
	fmt.Printf("Free Tier Eligible: %v\n", info.FreeTierEligible)
	if len(info.SupportedArchitectures) > 0 {
		fmt.Printf("Supported Architectures: %s\n", strings.Join(info.SupportedArchitectures, ", "))
	}
	if len(info.SupportedVirtualizationTypes) > 0 {
		fmt.Printf("Supported Virtualization Types: %s\n", strings.Join(info.SupportedVirtualizationTypes, ", "))
	}
}

// createInstanceExample demonstrates how to create and configure an EC2 instance
func createInstanceExample(ctx context.Context, client *awssdk.AWSClient, instanceTypeService *awssdk.InstanceTypeService, region string) {
	// Example: Create a t3.micro instance configuration
	instanceType := "t3.micro"

	// Get instance type details to validate
	instanceTypeInfo, err := instanceTypeService.GetInstanceType(ctx, instanceType, region)
	if err != nil {
		fmt.Printf("Error getting instance type info: %v\n", err)
		return
	}

	fmt.Printf("\nCreating instance configuration for: %s\n", instanceType)
	fmt.Printf("Category: %s\n", instanceTypeInfo.Category.String())
	fmt.Printf("Specs: %d vCPU, %.2f GiB RAM\n", instanceTypeInfo.VCPU, instanceTypeInfo.MemoryGiB)

	// Create instance configuration
	// Note: This is a configuration example - actual instance creation would require
	// valid AMI, subnet, and security group IDs
	instance := &awsec2model.Instance{
		Name:                     "example-instance",
		AMI:                      "ami-0c55b159cbfafe1f0", // Example AMI (Amazon Linux 2)
		InstanceType:             instanceType,
		SubnetID:                 "subnet-xxxxxxxxx",       // Replace with actual subnet ID
		VpcSecurityGroupIds:      []string{"sg-xxxxxxxxx"}, // Replace with actual security group ID
		AssociatePublicIPAddress: aws.Bool(true),
		RootVolumeID:             aws.String("vol-123"),
		Tags: []configs.Tag{
			{Key: "Name", Value: "example-instance"},
			{Key: "Environment", Value: "development"},
		},
	}

	// Validate instance configuration
	if err := instance.Validate(); err != nil {
		fmt.Printf("Instance validation error: %v\n", err)
		fmt.Println("\nNote: This is a configuration example. To actually create an instance,")
		fmt.Println("you need valid AMI, subnet, and security group IDs.")
		return
	}

	fmt.Println("\nInstance Configuration:")
	fmt.Printf("  Name: %s\n", instance.Name)
	fmt.Printf("  AMI: %s\n", instance.AMI)
	fmt.Printf("  Instance Type: %s\n", instance.InstanceType)
	fmt.Printf("  Subnet ID: %s\n", instance.SubnetID)
	fmt.Printf("  Security Groups: %v\n", instance.VpcSecurityGroupIds)
	if instance.AssociatePublicIPAddress != nil {
		fmt.Printf("  Associate Public IP: %v\n", *instance.AssociatePublicIPAddress)
	}
	if instance.RootVolumeID != nil {
		fmt.Printf("  Root Volume: %s, %d GB, Encrypted: %v\n",
			*instance.RootVolumeID,
			10,
			false)
	}

	// Example: Create instance using AWS SDK (commented out - requires valid IDs)
	/*
		createInput := &ec2.RunInstancesInput{
			ImageId:      aws.String(instance.AMI),
			InstanceType: types.InstanceType(instance.InstanceType),
			MinCount:     aws.Int32(1),
			MaxCount:     aws.Int32(1),
			SubnetId:     aws.String(instance.SubnetID),
			SecurityGroupIds: instance.VpcSecurityGroupIds,
		}

		if instance.AssociatePublicIPAddress != nil && *instance.AssociatePublicIPAddress {
			createInput.NetworkInterfaces = []types.InstanceNetworkInterfaceSpecification{
				{
					AssociatePublicIpAddress: instance.AssociatePublicIPAddress,
					SubnetId:                 aws.String(instance.SubnetID),
					Groups:                   instance.VpcSecurityGroupIds,
				},
			}
		}

		result, err := client.EC2.RunInstances(ctx, createInput)
		if err != nil {
			fmt.Printf("Error creating instance: %v\n", err)
			return
		}

		if len(result.Instances) > 0 {
			instanceID := result.Instances[0].InstanceId
			fmt.Printf("\nInstance created successfully: %s\n", aws.ToString(instanceID))
		}
	*/
}

// ListInstanceTypesByCategory lists all instance types in a specific category
func ListInstanceTypesByCategory(ctx context.Context, category awsmodel.InstanceCategory, region string) ([]*awsmodel.InstanceTypeInfo, error) {
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	service, err := awssdk.NewInstanceTypeService(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance type service: %w", err)
	}

	return service.ListByCategory(ctx, category, region)
}

// GetInstanceTypeInfo retrieves detailed information about a specific instance type
func GetInstanceTypeInfo(ctx context.Context, instanceType, region string) (*awsmodel.InstanceTypeInfo, error) {
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	service, err := awssdk.NewInstanceTypeService(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance type service: %w", err)
	}

	return service.GetInstanceType(ctx, instanceType, region)
}

// CreateInstance creates an EC2 instance with the given configuration
func CreateInstance(ctx context.Context, instance *awsec2model.Instance) (*ec2.RunInstancesOutput, error) {
	// Validate instance configuration
	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("instance validation failed: %w", err)
	}

	// Create AWS client
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	// Build RunInstances input
	createInput := &ec2.RunInstancesInput{
		ImageId:          aws.String(instance.AMI),
		InstanceType:     types.InstanceType(instance.InstanceType),
		MinCount:         aws.Int32(1),
		MaxCount:         aws.Int32(1),
		SubnetId:         aws.String(instance.SubnetID),
		SecurityGroupIds: instance.VpcSecurityGroupIds,
	}

	// Configure network interface if public IP is requested
	if instance.AssociatePublicIPAddress != nil && *instance.AssociatePublicIPAddress {
		createInput.NetworkInterfaces = []types.InstanceNetworkInterfaceSpecification{
			{
				AssociatePublicIpAddress: instance.AssociatePublicIPAddress,
				SubnetId:                 aws.String(instance.SubnetID),
				Groups:                   instance.VpcSecurityGroupIds,
				DeviceIndex:              aws.Int32(0),
			},
		}
		// Remove SubnetId and SecurityGroupIds from top level when using NetworkInterfaces
		createInput.SubnetId = nil
		createInput.SecurityGroupIds = nil
	}

	// Configure IAM instance profile if provided
	if instance.IAMInstanceProfile != nil {
		createInput.IamInstanceProfile = &types.IamInstanceProfileSpecification{
			Name: instance.IAMInstanceProfile,
		}
	}

	// Configure key pair if provided
	if instance.KeyName != nil {
		createInput.KeyName = instance.KeyName
	}

	// Configure user data if provided
	if instance.UserData != nil {
		createInput.UserData = instance.UserData
	}

	if instance.RootVolumeID != nil {
		blockDeviceMapping := types.BlockDeviceMapping{
			DeviceName: aws.String("/dev/xvda"), // Default root device name
		}
		createInput.BlockDeviceMappings = []types.BlockDeviceMapping{blockDeviceMapping}
	}

	// Add tags
	if len(instance.Tags) > 0 {
		var tagSpecs []types.TagSpecification
		var tags []types.Tag
		for _, tag := range instance.Tags {
			tags = append(tags, types.Tag{
				Key:   aws.String(tag.Key),
				Value: aws.String(tag.Value),
			})
		}
		tagSpecs = append(tagSpecs, types.TagSpecification{
			ResourceType: types.ResourceTypeInstance,
			Tags:         tags,
		})
		createInput.TagSpecifications = tagSpecs
	}

	// Create the instance
	result, err := client.EC2.RunInstances(ctx, createInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	return result, nil
}

// ListInstances lists EC2 instances with optional filters
func ListInstances(ctx context.Context, filters map[string][]string) ([]types.Instance, error) {
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	input := &ec2.DescribeInstancesInput{}

	// Convert filters to AWS filter format
	if len(filters) > 0 {
		var awsFilters []types.Filter
		for key, values := range filters {
			awsFilters = append(awsFilters, types.Filter{
				Name:   aws.String(key),
				Values: values,
			})
		}
		input.Filters = awsFilters
	}

	output, err := client.EC2.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}

	var instances []types.Instance
	for _, reservation := range output.Reservations {
		instances = append(instances, reservation.Instances...)
	}

	return instances, nil
}

// SearchInstanceTypes searches for instance types matching the given criteria
func SearchInstanceTypes(ctx context.Context, filters *awsmodel.InstanceTypeFilters, region string) ([]*awsmodel.InstanceTypeInfo, error) {
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	service, err := awssdk.NewInstanceTypeService(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance type service: %w", err)
	}

	return service.ListInstanceTypes(ctx, region, filters)
}

// DisplayInstanceTypesTable displays instance types in a formatted table
func DisplayInstanceTypesTable(ctx context.Context, region string, category *awsmodel.InstanceCategory) error {
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	service, err := awssdk.NewInstanceTypeService(client)
	if err != nil {
		return fmt.Errorf("failed to create instance type service: %w", err)
	}

	var types []*awsmodel.InstanceTypeInfo
	var err2 error

	if category != nil {
		types, err2 = service.ListByCategory(ctx, *category, region)
	} else {
		types, err2 = service.ListInstanceTypes(ctx, region, nil)
	}

	if err2 != nil {
		return err2
	}

	// Sort by name
	sort.Slice(types, func(i, j int) bool {
		return types[i].Name < types[j].Name
	})

	// Display table header
	fmt.Printf("\n%-20s %-25s %6s %10s %12s %15s\n", "Instance Type", "Category", "vCPU", "Memory(GiB)", "Network(Gbps)", "Free Tier")
	fmt.Println(strings.Repeat("-", 100))

	// Display table rows
	for _, info := range types {
		freeTier := "No"
		if info.FreeTierEligible {
			freeTier = "Yes"
		}
		fmt.Printf("%-20s %-25s %6d %10.2f %12.2f %15s\n",
			info.Name,
			info.Category.String(),
			info.VCPU,
			info.MemoryGiB,
			info.MaxNetworkGbps,
			freeTier,
		)
	}

	fmt.Printf("\nTotal: %d instance types\n", len(types))
	return nil
}

// Launch Template Functions

// CreateLaunchTemplateExample demonstrates how to create a Launch Template
func CreateLaunchTemplateExample(ctx context.Context) error {
	fmt.Println("\n--- Creating an Example Launch Template ---")
	fmt.Println("Launch Template creation would be implemented here")
	fmt.Println("Example: Create template with image_id, instance_type, security groups, etc.")
	return nil
}

// ListLaunchTemplates lists all Launch Templates
func ListLaunchTemplates(ctx context.Context) error {
	fmt.Println("\n--- Listing Launch Templates ---")
	fmt.Println("Launch Template listing would be implemented here")
	return nil
}

// DisplayLaunchTemplateDetails shows detailed information for a specific Launch Template
func DisplayLaunchTemplateDetails(ctx context.Context, templateID string) error {
	fmt.Printf("\n--- Details for Launch Template: %s ---\n", templateID)
	fmt.Println("Launch Template details display would be implemented here")
	return nil
}

// DisplayLaunchTemplateVersions shows version history for a Launch Template
func DisplayLaunchTemplateVersions(ctx context.Context, templateID string) error {
	fmt.Printf("\n--- Version History for Launch Template: %s ---\n", templateID)
	fmt.Println("Launch Template version history display would be implemented here")
	return nil
}
