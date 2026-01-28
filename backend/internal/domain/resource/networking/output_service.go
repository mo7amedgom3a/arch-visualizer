package networking

import (
	"context"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// NetworkingOutputService defines the interface for networking resource operations that return output DTOs
// This is a parallel interface to NetworkingService, providing output-specific models
type NetworkingOutputService interface {
	// VPC operations
	CreateVPCOutput(ctx context.Context, vpc *VPC) (*VPCOutput, error)
	GetVPCOutput(ctx context.Context, id string) (*VPCOutput, error)
	UpdateVPCOutput(ctx context.Context, vpc *VPC) (*VPCOutput, error)
	ListVPCsOutput(ctx context.Context, region string) ([]*VPCOutput, error)

	// Subnet operations
	CreateSubnetOutput(ctx context.Context, subnet *Subnet) (*SubnetOutput, error)
	GetSubnetOutput(ctx context.Context, id string) (*SubnetOutput, error)
	UpdateSubnetOutput(ctx context.Context, subnet *Subnet) (*SubnetOutput, error)
	ListSubnetsOutput(ctx context.Context, vpcID string) ([]*SubnetOutput, error)

	// Internet Gateway operations
	CreateInternetGatewayOutput(ctx context.Context, igw *InternetGateway) (*InternetGatewayOutput, error)

	// Route Table operations
	CreateRouteTableOutput(ctx context.Context, rt *RouteTable) (*RouteTableOutput, error)
	GetRouteTableOutput(ctx context.Context, id string) (*RouteTableOutput, error)

	// Security Group operations
	CreateSecurityGroupOutput(ctx context.Context, sg *SecurityGroup) (*SecurityGroupOutput, error)
	GetSecurityGroupOutput(ctx context.Context, id string) (*SecurityGroupOutput, error)
	UpdateSecurityGroupOutput(ctx context.Context, sg *SecurityGroup) (*SecurityGroupOutput, error)

	// NAT Gateway operations
	CreateNATGatewayOutput(ctx context.Context, ngw *NATGateway) (*NATGatewayOutput, error)
	GetNATGatewayOutput(ctx context.Context, id string) (*NATGatewayOutput, error)

	// Elastic IP operations
	AllocateElasticIPOutput(ctx context.Context, eip *ElasticIP) (*ElasticIPOutput, error)
	GetElasticIPOutput(ctx context.Context, id string) (*ElasticIPOutput, error)
	ListElasticIPsOutput(ctx context.Context, region string) ([]*ElasticIPOutput, error)

	// Network ACL operations
	CreateNetworkACLOutput(ctx context.Context, acl *NetworkACL) (*NetworkACLOutput, error)
	GetNetworkACLOutput(ctx context.Context, id string) (*NetworkACLOutput, error)
	ListNetworkACLsOutput(ctx context.Context, vpcID string) ([]*NetworkACLOutput, error)

	// Network Interface operations
	CreateNetworkInterfaceOutput(ctx context.Context, eni *NetworkInterface) (*NetworkInterfaceOutput, error)
	GetNetworkInterfaceOutput(ctx context.Context, id string) (*NetworkInterfaceOutput, error)
	ListNetworkInterfacesOutput(ctx context.Context, subnetID string) ([]*NetworkInterfaceOutput, error)

	// Pricing operations
	EstimateResourceCost(ctx context.Context, resourceType string, config map[string]interface{}, duration time.Duration) (*pricing.CostEstimate, error)
	GetResourcePricing(ctx context.Context, resourceType string, region string) (*pricing.ResourcePricing, error)
}
