package pricing

import (
	"context"
	"fmt"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing/compute"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing/storage"
	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// AWSPricingService implements the PricingService interface for AWS
type AWSPricingService struct {
	calculator *AWSPricingCalculator
}

// NewAWSPricingService creates a new AWS pricing service
func NewAWSPricingService() *AWSPricingService {
	service := &AWSPricingService{}
	service.calculator = NewAWSPricingCalculator(service)
	return service
}

// GetPricing retrieves the pricing information for a specific resource type
func (s *AWSPricingService) GetPricing(ctx context.Context, resourceType string, provider string, region string) (*domainpricing.ResourcePricing, error) {
	if provider != "aws" {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	switch resourceType {
	case "nat_gateway":
		return networking.GetNATGatewayPricing(region), nil
	case "elastic_ip":
		return networking.GetElasticIPPricing(region), nil
	case "network_interface":
		return networking.GetNetworkInterfacePricing(region), nil
	case "data_transfer":
		return networking.GetDataTransferPricing(region), nil
	case "ec2_instance":
		// Default to t3.micro if instance type not provided
		instanceType := "t3.micro"
		return compute.GetEC2InstancePricing(instanceType, region), nil
	case "ebs_volume":
		// Default to gp3 if volume type not provided
		volumeType := "gp3"
		return storage.GetEBSVolumePricing(volumeType, region), nil
	case "s3_bucket":
		// Default to standard storage class if not provided
		storageClass := "standard"
		return storage.GetS3BucketPricing(storageClass, region), nil
	case "load_balancer":
		// Default to application if LB type not provided
		lbType := "application"
		return compute.GetLoadBalancerPricing(lbType, region), nil
	case "auto_scaling_group":
		// Default values if not provided
		instanceType := "t3.micro"
		minSize := 1
		maxSize := 3
		return compute.GetAutoScalingGroupPricing(instanceType, minSize, maxSize, region), nil
	case "lambda_function":
		// Default to 128 MB memory if not provided
		memorySizeMB := 128.0
		return compute.GetLambdaFunctionPricing(memorySizeMB, region), nil
	default:
		return nil, fmt.Errorf("pricing not available for resource type: %s", resourceType)
	}
}

// EstimateCost estimates the cost for a resource over a given duration
func (s *AWSPricingService) EstimateCost(ctx context.Context, res *resource.Resource, duration time.Duration) (*domainpricing.CostEstimate, error) {
	return s.calculator.CalculateResourceCost(ctx, res, duration)
}

// EstimateArchitectureCost estimates the total cost for multiple resources over a given duration
func (s *AWSPricingService) EstimateArchitectureCost(ctx context.Context, resources []*resource.Resource, duration time.Duration) (*domainpricing.CostEstimate, error) {
	return s.calculator.CalculateArchitectureCost(ctx, resources, duration)
}

// ListSupportedResources returns a list of resource types that have pricing information
func (s *AWSPricingService) ListSupportedResources(ctx context.Context, provider string) ([]string, error) {
	if provider != "aws" {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	return []string{
		"nat_gateway",
		"elastic_ip",
		"network_interface",
		"data_transfer",
		"ec2_instance",
		"ebs_volume",
		"s3_bucket",
		"load_balancer",
		"auto_scaling_group",
		"lambda_function",
	}, nil
}
