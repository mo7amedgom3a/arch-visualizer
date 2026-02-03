package services

import (
	"context"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	ec2_outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	net_outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/common"
)

// StaticDataServiceImpl implements StaticDataService interface
type StaticDataServiceImpl struct {
	resourceTypeRepo serverinterfaces.ResourceTypeRepository
}

// NewStaticDataService creates a new static data service
func NewStaticDataService(resourceTypeRepo serverinterfaces.ResourceTypeRepository) serverinterfaces.StaticDataService {
	return &StaticDataServiceImpl{
		resourceTypeRepo: resourceTypeRepo,
	}
}

// ListResourceTypes retrieves all resource types grouped by category
func (s *StaticDataServiceImpl) ListResourceTypes(ctx context.Context) ([]serverinterfaces.ResourceTypeGroup, error) {
	// Quick hack: getting AWS, Azure, GCP types
	// Since the repo returns a flat list, we need to group them.
	// For this mock implementation, we'll try to fetch all providers we know about.

	var allTypes []*models.ResourceType

	// Fetch for AWS
	awsTypes, err := s.resourceTypeRepo.ListByProvider(ctx, "aws")
	if err == nil {
		allTypes = append(allTypes, awsTypes...)
	}

	// Fetch for GCP
	gcpTypes, err := s.resourceTypeRepo.ListByProvider(ctx, "gcp")
	if err == nil {
		allTypes = append(allTypes, gcpTypes...)
	}

	// Fetch for Azure
	azureTypes, err := s.resourceTypeRepo.ListByProvider(ctx, "azure")
	if err == nil {
		allTypes = append(allTypes, azureTypes...)
	}

	return s.groupResourceTypes(allTypes), nil
}

// ListResourceTypesByProvider retrieves resource types for a provider grouped by category
func (s *StaticDataServiceImpl) ListResourceTypesByProvider(ctx context.Context, provider string) ([]serverinterfaces.ResourceTypeGroup, error) {
	types, err := s.resourceTypeRepo.ListByProvider(ctx, provider)
	if err != nil {
		return nil, err
	}
	return s.groupResourceTypes(types), nil
}

// ListResourceOutputModels retrieves output models for resources with default values grouped by category
func (s *StaticDataServiceImpl) ListResourceOutputModels(ctx context.Context, provider string) ([]serverinterfaces.ResourceModelGroup, error) {
	// Only AWS supported for now
	if provider != "aws" && provider != "" {
		return []serverinterfaces.ResourceModelGroup{}, nil
	}

	now := time.Now()
	var groups []serverinterfaces.ResourceModelGroup

	// --- Networking Group ---
	var netResources []serverinterfaces.ResourceModel

	netResources = append(netResources, serverinterfaces.ResourceModel{
		Name: "VPC",
		Model: net_outputs.VPCOutput{
			ID:           "vpc-12345678",
			ARN:          "arn:aws:ec2:region:account:vpc/vpc-12345678",
			Name:         "MyVPC",
			Region:       "us-east-1",
			CIDR:         "10.0.0.0/16",
			State:        "available",
			CreationTime: now,
			Tags:         []configs.Tag{{Key: "Environment", Value: "Dev"}},
		},
	})

	netResources = append(netResources, serverinterfaces.ResourceModel{
		Name: "Subnet",
		Model: net_outputs.SubnetOutput{
			ID:               "subnet-12345678",
			VPCID:            "vpc-12345678",
			CIDR:             "10.0.1.0/24",
			AvailableIPCount: 251,
			State:            "available",
			AvailabilityZone: "us-east-1a",
			Tags:             []configs.Tag{},
		},
	})

	netResources = append(netResources, serverinterfaces.ResourceModel{
		Name: "InternetGateway",
		Model: net_outputs.InternetGatewayOutput{
			ID:    "igw-12345678",
			VPCID: "vpc-12345678",
			State: "attached",
		},
	})

	netResources = append(netResources, serverinterfaces.ResourceModel{
		Name: "NATGateway",
		Model: net_outputs.NATGatewayOutput{
			ID:           "nat-12345678",
			SubnetID:     "subnet-12345678",
			State:        "available",
			PublicIP:     "1.2.3.4",
			CreationTime: now,
		},
	})

	gatewayID := "igw-12345678"
	netResources = append(netResources, serverinterfaces.ResourceModel{
		Name: "RouteTable",
		Model: net_outputs.RouteTableOutput{
			ID:    "rtb-12345678",
			VPCID: "vpc-12345678",
			Routes: []networking.Route{
				{DestinationCIDRBlock: "10.0.0.0/16", GatewayID: nil},
				{DestinationCIDRBlock: "0.0.0.0/0", GatewayID: &gatewayID},
			},
		},
	})

	fromPortHTTP := 80
	toPortHTTP := 80
	netResources = append(netResources, serverinterfaces.ResourceModel{
		Name: "SecurityGroup",
		Model: net_outputs.SecurityGroupOutput{
			ID:          "sg-12345678",
			VPCID:       "vpc-12345678",
			Name:        "MySecurityGroup",
			Description: "Allow HTTP",
			Rules: []networking.SecurityGroupRule{
				{
					Type:        "ingress",
					Protocol:    "tcp",
					FromPort:    &fromPortHTTP,
					ToPort:      &toPortHTTP,
					CIDRBlocks:  []string{"0.0.0.0/0"},
					Description: "HTTP",
				},
				{
					Type:        "egress",
					Protocol:    "-1",
					CIDRBlocks:  []string{"0.0.0.0/0"},
					Description: "All Traffic",
				},
			},
		},
	})

	groups = append(groups, serverinterfaces.ResourceModelGroup{
		ServiceType: "Networking",
		Resources:   netResources,
	})

	// --- Compute Group ---
	var computeResources []serverinterfaces.ResourceModel

	publicIP := "1.2.3.4"
	publicDNS := "ec2-1-2-3-4.compute-1.amazonaws.com"
	keyName := "my-key"

	computeResources = append(computeResources, serverinterfaces.ResourceModel{
		Name: "EC2",
		Model: ec2_outputs.InstanceOutput{
			ID:               "i-1234567890abcdef0",
			ARN:              "arn:aws:ec2:us-east-1:123456789012:instance/i-0",
			Name:             "MyInstance",
			Region:           "us-east-1",
			InstanceType:     "t3.micro",
			AMI:              "ami-12345678",
			State:            "running",
			AvailabilityZone: "us-east-1a",
			CreationTime:     now,
			PublicIP:         &publicIP,
			PrivateIP:        "10.0.1.10",
			PublicDNS:        &publicDNS,
			PrivateDNS:       "ip-10-0-1-10.ec2.internal",
			SubnetID:         "subnet-12345678",
			VPCID:            "vpc-12345678",
			SecurityGroupIDs: []string{"sg-12345678"},
			KeyName:          &keyName,
			Tags:             []configs.Tag{{Key: "Name", Value: "MyInstance"}},
		},
	})

	groups = append(groups, serverinterfaces.ResourceModelGroup{
		ServiceType: "Compute",
		Resources:   computeResources,
	})

	return groups, nil
}

// groupResourceTypes helper to group resources by category
func (s *StaticDataServiceImpl) groupResourceTypes(types []*models.ResourceType) []serverinterfaces.ResourceTypeGroup {
	groupsMap := make(map[string][]*models.ResourceType)

	for _, t := range types {
		categoryName := "Uncategorized"
		if t.Category != nil {
			categoryName = t.Category.Name
		} else if t.CategoryID != nil {
			// If Category struct is not loaded but ID is present, we might want to fetch or just use "Unknown".
			// Given this is a simple list, maybe we rely on Preloading in Repo.
			// If not preloaded, well... "Uncategorized" for now.
			// Ideally repo should Preload("Category").
		}

		groupsMap[categoryName] = append(groupsMap[categoryName], t)
	}

	// Convert map to slice
	var groups []serverinterfaces.ResourceTypeGroup
	for serviceType, resources := range groupsMap {
		groups = append(groups, serverinterfaces.ResourceTypeGroup{
			ServiceType: serviceType,
			Resources:   resources,
		})
	}

	// Optional: Sort groups by ServiceType for consistent output?
	// Leaving order random for now unless required.

	return groups
}

// ListProviders retrieves supported cloud providers
func (s *StaticDataServiceImpl) ListProviders(ctx context.Context) ([]string, error) {
	return []string{"aws", "azure", "gcp"}, nil
}

// ListCloudConfiguration retrieves global cloud configurations
func (s *StaticDataServiceImpl) ListCloudConfiguration(ctx context.Context, provider string) (*serverinterfaces.CloudConfig, error) {
	config := &serverinterfaces.CloudConfig{
		Provider: provider,
	}

	if provider == "aws" {
		// Regions
		for _, regionCode := range common.SupportedRegions {
			config.Regions = append(config.Regions, serverinterfaces.RegionConfig{
				Code: regionCode,
				Name: common.FormatRegionName(regionCode),
				AZs:  common.GetAZsForRegion(regionCode),
			})
		}

		// Instance Types (Mock data for commonly used types)
		config.InstanceTypes = []serverinterfaces.InstanceTypeConfig{
			{Name: "t3.nano", VCPU: 2, MemoryGiB: 0.5},
			{Name: "t3.micro", VCPU: 2, MemoryGiB: 1.0},
			{Name: "t3.small", VCPU: 2, MemoryGiB: 2.0},
			{Name: "t3.medium", VCPU: 2, MemoryGiB: 4.0},
			{Name: "t3.large", VCPU: 2, MemoryGiB: 8.0},
			{Name: "m5.large", VCPU: 2, MemoryGiB: 8.0},
			{Name: "m5.xlarge", VCPU: 4, MemoryGiB: 16.0},
			{Name: "c5.large", VCPU: 2, MemoryGiB: 4.0},
			{Name: "r5.large", VCPU: 2, MemoryGiB: 16.0},
		}

		// Storage Types
		config.StorageTypes = []string{"gp2", "gp3", "io1", "io2", "st1", "sc1", "standard"}
	} else {
		// Return empty for other providers for now
		config.Regions = []serverinterfaces.RegionConfig{}
		config.InstanceTypes = []serverinterfaces.InstanceTypeConfig{}
		config.StorageTypes = []string{}
	}

	return config, nil
}
