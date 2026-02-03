package architecture

import (
	"fmt"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// AWSResourceTypeMapper implements ResourceTypeMapper for AWS
type AWSResourceTypeMapper struct{}

// NewAWSResourceTypeMapper creates a new AWS resource type mapper
func NewAWSResourceTypeMapper() *AWSResourceTypeMapper {
	return &AWSResourceTypeMapper{}
}

// MapIRTypeToResourceType maps an IR type (kebab-case) to ResourceType for AWS
func (m *AWSResourceTypeMapper) MapIRTypeToResourceType(irType string) (*resource.ResourceType, error) {
	// First, try to get resource name from inventory
	if mapper, ok := architecture.GetIRTypeMapper(resource.AWS); ok {
		if resourceName, found := mapper.GetResourceNameByIRType(irType); found {
			// Map resource name to ResourceType
			return m.MapResourceNameToResourceType(resourceName)
		}
		// Try lowercase fallback
		if resourceName, found := mapper.GetResourceNameByIRType(strings.ToLower(irType)); found {
			return m.MapResourceNameToResourceType(resourceName)
		}
	}

	// Final fallback: try to map based on the input string being a valid ResourceName already
	// This handles cases where the FE sends "IAMPolicy" or "Lambda" directly as the type
	if rt, err := m.MapResourceNameToResourceType(irType); err == nil {
		return rt, nil
	}

	return nil, fmt.Errorf("unknown IR type for AWS: %s", irType)
}

// MapResourceNameToResourceType maps a resource name (PascalCase) to ResourceType for AWS
func (m *AWSResourceTypeMapper) MapResourceNameToResourceType(resourceName string) (*resource.ResourceType, error) {
	resourceTypeMap := map[string]resource.ResourceType{
		"VPC": {
			ID:         "vpc",
			Name:       "VPC",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"Subnet": {
			ID:         "subnet",
			Name:       "Subnet",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"EC2": {
			ID:         "ec2",
			Name:       "EC2",
			Category:   string(resource.CategoryCompute),
			Kind:       "VirtualMachine",
			IsRegional: true,
			IsGlobal:   false,
		},
		"RouteTable": {
			ID:         "route-table",
			Name:       "RouteTable",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"SecurityGroup": {
			ID:         "security-group",
			Name:       "SecurityGroup",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"NATGateway": {
			ID:         "nat-gateway",
			Name:       "NATGateway",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Gateway",
			IsRegional: true,
			IsGlobal:   false,
		},
		"InternetGateway": {
			ID:         "internet-gateway",
			Name:       "InternetGateway",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Gateway",
			IsRegional: true,
			IsGlobal:   false,
		},
		"ElasticIP": {
			ID:         "elastic-ip",
			Name:       "ElasticIP",
			Category:   string(resource.CategoryNetworking),
			Kind:       "Network",
			IsRegional: true,
			IsGlobal:   false,
		},
		"Lambda": {
			ID:         "lambda",
			Name:       "Lambda",
			Category:   string(resource.CategoryCompute),
			Kind:       "Function",
			IsRegional: true,
			IsGlobal:   false,
		},
		"S3": {
			ID:         "s3",
			Name:       "S3",
			Category:   string(resource.CategoryStorage),
			Kind:       "Storage",
			IsRegional: false,
			IsGlobal:   true,
		},
		"EBS": {
			ID:         "ebs",
			Name:       "EBS",
			Category:   string(resource.CategoryStorage),
			Kind:       "Storage",
			IsRegional: true,
			IsGlobal:   false,
		},
		"RDS": {
			ID:         "rds",
			Name:       "RDS",
			Category:   string(resource.CategoryDatabase),
			Kind:       "Database",
			IsRegional: true,
			IsGlobal:   false,
		},
		"DynamoDB": {
			ID:         "dynamodb",
			Name:       "DynamoDB",
			Category:   string(resource.CategoryDatabase),
			Kind:       "Database",
			IsRegional: true,
			IsGlobal:   false,
		},
		"LoadBalancer": {
			ID:         "load-balancer",
			Name:       "LoadBalancer",
			Category:   string(resource.CategoryCompute),
			Kind:       "LoadBalancer",
			IsRegional: true,
			IsGlobal:   false,
		},
		"AutoScalingGroup": {
			ID:         "auto-scaling-group",
			Name:       "AutoScalingGroup",
			Category:   string(resource.CategoryCompute),
			Kind:       "VirtualMachine",
			IsRegional: true,
			IsGlobal:   false,
		},
		"IAMPolicy": {
			ID:         "iam-policy",
			Name:       "IAMPolicy",
			Category:   string(resource.CategoryIAM),
			Kind:       "Policy",
			IsRegional: false,
			IsGlobal:   true,
		},
		"IAMUser": {
			ID:         "iam-user",
			Name:       "IAMUser",
			Category:   string(resource.CategoryIAM),
			Kind:       "User",
			IsRegional: false,
			IsGlobal:   true,
		},
		"IAMRole": {
			ID:         "iam-role",
			Name:       "IAMRole",
			Category:   string(resource.CategoryIAM),
			Kind:       "Role",
			IsRegional: false,
			IsGlobal:   true,
		},
		"IAMRolePolicyAttachment": {
			ID:         "iam-role-policy-attachment",
			Name:       "IAMRolePolicyAttachment",
			Category:   string(resource.CategoryIAM),
			Kind:       "Attachment",
			IsRegional: false,
			IsGlobal:   true,
		},
	}

	rt, exists := resourceTypeMap[resourceName]
	if !exists {
		return nil, fmt.Errorf("unknown AWS resource name: %s", resourceName)
	}

	return &rt, nil
}
