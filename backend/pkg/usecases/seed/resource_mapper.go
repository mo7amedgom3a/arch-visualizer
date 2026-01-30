package seed

import (
	"encoding/json"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/inventory"
	awsmappercompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/compute"
	awsmappernetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/networking"
	awsmapperstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/storage"
	awsec2 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// ResourceMapper handles mapping between database resources and domain/AWS models
type ResourceMapper struct{}

// NewResourceMapper creates a new resource mapper
func NewResourceMapper() *ResourceMapper {
	return &ResourceMapper{}
}

// DatabaseResourceToDomain converts a database resource to a domain model
// based on the resource type name
func (rm *ResourceMapper) DatabaseResourceToDomain(dbResource *models.Resource) (interface{}, error) {
	if dbResource == nil {
		return nil, fmt.Errorf("database resource is nil")
	}

	// Get resource type name
	resourceTypeName := dbResource.ResourceType.Name
	if resourceTypeName == "" {
		return nil, fmt.Errorf("resource type name is empty")
	}

	// Try to use inventory first
	inv := inventory.GetDefaultInventory()
	if functions, ok := inv.GetFunctions(resourceTypeName); ok && functions.DomainMapper != nil {
		return functions.DomainMapper(dbResource.Config)
	}

	// Fallback to switch-based unmarshaling
	return rm.databaseResourceToDomainFallback(dbResource)
}

// databaseResourceToDomainFallback provides backward compatibility with switch-based unmarshaling
func (rm *ResourceMapper) databaseResourceToDomainFallback(dbResource *models.Resource) (interface{}, error) {
	resourceTypeName := dbResource.ResourceType.Name
	
	switch resourceTypeName {
	case "VPC":
		var vpc domainnetworking.VPC
		if err := json.Unmarshal(dbResource.Config, &vpc); err != nil {
			return nil, fmt.Errorf("failed to unmarshal VPC: %w", err)
		}
		return &vpc, nil

	case "Subnet":
		var subnet domainnetworking.Subnet
		if err := json.Unmarshal(dbResource.Config, &subnet); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Subnet: %w", err)
		}
		return &subnet, nil

	case "InternetGateway":
		var igw domainnetworking.InternetGateway
		if err := json.Unmarshal(dbResource.Config, &igw); err != nil {
			return nil, fmt.Errorf("failed to unmarshal InternetGateway: %w", err)
		}
		return &igw, nil

	case "NATGateway":
		var nat domainnetworking.NATGateway
		if err := json.Unmarshal(dbResource.Config, &nat); err != nil {
			return nil, fmt.Errorf("failed to unmarshal NATGateway: %w", err)
		}
		return &nat, nil

	case "SecurityGroup":
		var sg domainnetworking.SecurityGroup
		if err := json.Unmarshal(dbResource.Config, &sg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal SecurityGroup: %w", err)
		}
		return &sg, nil

	case "ElasticIP":
		var eip domainnetworking.ElasticIP
		if err := json.Unmarshal(dbResource.Config, &eip); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ElasticIP: %w", err)
		}
		return &eip, nil

	case "RouteTable":
		var rt domainnetworking.RouteTable
		if err := json.Unmarshal(dbResource.Config, &rt); err != nil {
			return nil, fmt.Errorf("failed to unmarshal RouteTable: %w", err)
		}
		return &rt, nil

	case "EC2":
		var instance domaincompute.Instance
		if err := json.Unmarshal(dbResource.Config, &instance); err != nil {
			return nil, fmt.Errorf("failed to unmarshal EC2 Instance: %w", err)
		}
		return &instance, nil

	case "Lambda":
		var lambda domaincompute.LambdaFunction
		if err := json.Unmarshal(dbResource.Config, &lambda); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Lambda Function: %w", err)
		}
		return &lambda, nil

	case "LoadBalancer":
		var lb domaincompute.LoadBalancer
		if err := json.Unmarshal(dbResource.Config, &lb); err != nil {
			return nil, fmt.Errorf("failed to unmarshal LoadBalancer: %w", err)
		}
		return &lb, nil

	case "AutoScalingGroup":
		var asg domaincompute.AutoScalingGroup
		if err := json.Unmarshal(dbResource.Config, &asg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal AutoScalingGroup: %w", err)
		}
		return &asg, nil

	case "S3":
		var s3 domainstorage.S3Bucket
		if err := json.Unmarshal(dbResource.Config, &s3); err != nil {
			return nil, fmt.Errorf("failed to unmarshal S3 Bucket: %w", err)
		}
		return &s3, nil

	case "EBS":
		var ebs domainstorage.EBSVolume
		if err := json.Unmarshal(dbResource.Config, &ebs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal EBS Volume: %w", err)
		}
		return &ebs, nil

	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceTypeName)
	}
}

// DomainToAWSModel converts a domain model to an AWS model
// This simulates the mapping that would happen when calling AWS services
func (rm *ResourceMapper) DomainToAWSModel(domainResource interface{}) (interface{}, error) {
	if domainResource == nil {
		return nil, fmt.Errorf("domain resource is nil")
	}

	// Try to use inventory first - need to determine resource type from the domain resource
	// For now, fallback to switch-based type assertion
	return rm.domainToAWSModelFallback(domainResource)
}

// domainToAWSModelFallback provides backward compatibility with switch-based type assertion
func (rm *ResourceMapper) domainToAWSModelFallback(domainResource interface{}) (interface{}, error) {
	switch resource := domainResource.(type) {
	case *domainnetworking.VPC:
		return awsmappernetworking.FromDomainVPC(resource), nil

	case *domainnetworking.Subnet:
		// FromDomainSubnet requires availability zone, use default if not set
		az := "us-east-1a"
		if resource.AvailabilityZone != nil && *resource.AvailabilityZone != "" {
			az = *resource.AvailabilityZone
		}
		return awsmappernetworking.FromDomainSubnet(resource, az), nil

	case *domainnetworking.InternetGateway:
		return awsmappernetworking.FromDomainInternetGateway(resource), nil

	case *domainnetworking.NATGateway:
		return awsmappernetworking.FromDomainNATGateway(resource), nil

	case *domainnetworking.SecurityGroup:
		return awsmappernetworking.FromDomainSecurityGroup(resource), nil

	case *domainnetworking.ElasticIP:
		return awsmappernetworking.FromDomainElasticIP(resource), nil

	case *domainnetworking.RouteTable:
		return awsmappernetworking.FromDomainRouteTable(resource), nil

	case *domaincompute.Instance:
		return awsmappercompute.FromDomainInstance(resource), nil

	case *domaincompute.LambdaFunction:
		return awsmappercompute.FromDomainLambdaFunction(resource), nil

	case *domaincompute.LoadBalancer:
		return awsmappercompute.FromDomainLoadBalancer(resource), nil

	case *domaincompute.AutoScalingGroup:
		return awsmappercompute.FromDomainAutoScalingGroup(resource), nil

	case *domainstorage.S3Bucket:
		return awsmapperstorage.FromDomainS3Bucket(resource), nil

	case *domainstorage.EBSVolume:
		return awsmapperstorage.FromDomainEBSVolume(resource), nil

	default:
		return nil, fmt.Errorf("unsupported domain resource type: %T", domainResource)
	}
}

// DatabaseResourceToAWSModel converts a database resource directly to an AWS model
// This is a convenience method that combines DatabaseResourceToDomain and DomainToAWSModel
func (rm *ResourceMapper) DatabaseResourceToAWSModel(dbResource *models.Resource) (interface{}, error) {
	domainResource, err := rm.DatabaseResourceToDomain(dbResource)
	if err != nil {
		return nil, fmt.Errorf("failed to convert database resource to domain: %w", err)
	}

	awsModel, err := rm.DomainToAWSModel(domainResource)
	if err != nil {
		return nil, fmt.Errorf("failed to convert domain resource to AWS model: %w", err)
	}

	return awsModel, nil
}

// SimulateAWSResponse simulates an AWS service response by converting a database resource
// through domain model to AWS output model
// This demonstrates the full flow: Database -> Domain -> AWS Input -> AWS Output (simulated)
func (rm *ResourceMapper) SimulateAWSResponse(dbResource *models.Resource) (interface{}, error) {
	// Step 1: Database resource to domain model
	domainResource, err := rm.DatabaseResourceToDomain(dbResource)
	if err != nil {
		return nil, fmt.Errorf("failed to convert database resource to domain: %w", err)
	}

	// Step 2: Domain model to AWS input model
	awsInput, err := rm.DomainToAWSModel(domainResource)
	if err != nil {
		return nil, fmt.Errorf("failed to convert domain to AWS input: %w", err)
	}

	// Step 3: Simulate AWS service response (in a real scenario, this would call AWS SDK)
	// For now, we return the AWS input model as if it were the output
	// In production, you would call the virtual service here
	return awsInput, nil
}

// GetResourceTypeName returns the resource type name from a database resource
func (rm *ResourceMapper) GetResourceTypeName(dbResource *models.Resource) string {
	if dbResource == nil || dbResource.ResourceType.Name == "" {
		return ""
	}
	return dbResource.ResourceType.Name
}

// Example usage functions for demonstration

// ExampleReadAndMapVPC demonstrates reading a VPC from database and mapping to domain/AWS
func ExampleReadAndMapVPC(dbResource *models.Resource) error {
	mapper := NewResourceMapper()

	// Read from database and convert to domain model
	domainVPC, err := mapper.DatabaseResourceToDomain(dbResource)
	if err != nil {
		return err
	}

	vpc := domainVPC.(*domainnetworking.VPC)
	fmt.Printf("Domain VPC: %s (ID: %s, CIDR: %s)\n", vpc.Name, vpc.ID, vpc.CIDR)

	// Convert domain model to AWS model
	awsVPC := awsmappernetworking.FromDomainVPC(vpc)
	fmt.Printf("AWS VPC: %s (CIDR: %s)\n", awsVPC.Name, awsVPC.CIDR)

	return nil
}

// ExampleReadAndMapInstance demonstrates reading an EC2 instance from database
func ExampleReadAndMapInstance(dbResource *models.Resource) error {
	mapper := NewResourceMapper()

	// Direct conversion to AWS model
	awsInstance, err := mapper.DatabaseResourceToAWSModel(dbResource)
	if err != nil {
		return err
	}

	instance := awsInstance.(*awsec2.Instance)
	fmt.Printf("AWS Instance: %s (Type: %s, AMI: %s)\n", instance.Name, instance.InstanceType, instance.AMI)

	return nil
}

// ExampleReadAndMapS3Bucket demonstrates reading an S3 bucket from database
func ExampleReadAndMapS3Bucket(dbResource *models.Resource) error {
	mapper := NewResourceMapper()

	// Read to domain model
	domainS3, err := mapper.DatabaseResourceToDomain(dbResource)
	if err != nil {
		return err
	}

	s3 := domainS3.(*domainstorage.S3Bucket)
	fmt.Printf("Domain S3: %s (Region: %s)\n", s3.Name, s3.Region)

	// Convert to AWS model
	awsS3 := awsmapperstorage.FromDomainS3Bucket(s3)
	bucketName := "N/A"
	if awsS3.Bucket != nil {
		bucketName = *awsS3.Bucket
	}
	fmt.Printf("AWS S3: %s\n", bucketName)

	return nil
}
