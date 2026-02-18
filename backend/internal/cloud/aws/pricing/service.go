package pricing

import (
	"context"
	"fmt"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/inventory"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing/compute"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing/storage"
	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	pricingrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/pricing"
	resourcerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/resource"
)

// AWSPricingService implements the PricingService interface for AWS
type AWSPricingService struct {
	calculator *AWSPricingCalculator
}

// NewAWSPricingService creates a new AWS pricing service
func NewAWSPricingService() *AWSPricingService {
	service := &AWSPricingService{}
	service.calculator = NewAWSPricingCalculator(service)

	// Note: We intentionally do NOT set up inventory pricing functions here
	// to avoid circular dependencies. The calculator's fallback mechanism
	// (calculateResourceCostFallback) handles pricing calculations directly.
	// If inventory-based pricing is needed, it should be set up separately
	// with proper guards against recursion.

	return service
}

// NewAWSPricingServiceWithRepos creates a new AWS pricing service with repositories for DB-driven pricing
func NewAWSPricingServiceWithRepos(
	pricingRateRepo *pricingrepo.PricingRateRepository,
	hiddenDepRepo *resourcerepo.HiddenDependencyRepository,
) *AWSPricingService {
	service := &AWSPricingService{}
	service.calculator = NewAWSPricingCalculatorWithRepos(service, pricingRateRepo, hiddenDepRepo)
	return service
}

// GetCalculator returns the pricing calculator
func (s *AWSPricingService) GetCalculator() *AWSPricingCalculator {
	return s.calculator
}

// mapPricingTypeToResourceNameReverse maps domain resource name to pricing service resource type
func mapPricingTypeToResourceNameReverse(resourceName string) string {
	mapping := map[string]string{
		"VPC":              "vpc",
		"Subnet":           "subnet",
		"RouteTable":       "route_table",
		"SecurityGroup":    "security_group",
		"InternetGateway":  "internet_gateway",
		"NATGateway":       "nat_gateway",
		"ElasticIP":        "elastic_ip",
		"EC2":              "ec2_instance",
		"Lambda":           "lambda_function",
		"LoadBalancer":     "load_balancer",
		"AutoScalingGroup": "auto_scaling_group",
		"S3":               "s3_bucket",
		"EBS":              "ebs_volume",
		"RDS":              "rds_instance",
		"DynamoDB":         "dynamodb_table",
	}

	if mapped, ok := mapping[resourceName]; ok {
		return mapped
	}

	// Default: convert PascalCase to snake_case
	return toSnakeCase(resourceName)
}

// toSnakeCase converts PascalCase to snake_case (simple implementation)
func toSnakeCase(s string) string {
	result := ""
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		result += string(r)
	}
	return result
}

// GetPricing retrieves the pricing information for a specific resource type
func (s *AWSPricingService) GetPricing(ctx context.Context, resourceType string, provider string, region string) (*domainpricing.ResourcePricing, error) {
	if provider != "aws" {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Try to use inventory first
	inv := inventory.GetDefaultInventory()

	// Map pricing resource type to domain resource name
	resourceName := mapPricingTypeToResourceName(resourceType)
	if resourceName != "" {
		if functions, ok := inv.GetFunctions(resourceName); ok && functions.GetPricingInfo != nil {
			return functions.GetPricingInfo(region)
		}
	}

	// Fallback to switch-based pricing (for backward compatibility)
	return s.getPricingFallback(resourceType, region)
}

// getPricingFallback provides backward compatibility with switch-based pricing
func (s *AWSPricingService) getPricingFallback(resourceType string, region string) (*domainpricing.ResourcePricing, error) {
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

// mapPricingTypeToResourceName maps pricing service resource type to domain resource name
func mapPricingTypeToResourceName(pricingType string) string {
	mapping := map[string]string{
		"nat_gateway":        "NATGateway",
		"elastic_ip":         "ElasticIP",
		"network_interface":  "NetworkInterface",
		"ec2_instance":       "EC2",
		"ebs_volume":         "EBS",
		"s3_bucket":          "S3",
		"load_balancer":      "LoadBalancer",
		"auto_scaling_group": "AutoScalingGroup",
		"lambda_function":    "Lambda",
	}

	if mapped, ok := mapping[pricingType]; ok {
		return mapped
	}
	return ""
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
