package networking

import (
	"context"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
)

// AWSNetworkingService defines AWS-specific networking operations
// This implements cloud provider-specific logic while maintaining domain compatibility
type AWSNetworkingService interface {
	// VPC operations
	CreateVPC(ctx context.Context, vpc *awsnetworking.VPC) (*awsnetworking.VPC, error)
	GetVPC(ctx context.Context, id string) (*awsnetworking.VPC, error)
	UpdateVPC(ctx context.Context, vpc *awsnetworking.VPC) (*awsnetworking.VPC, error)
	DeleteVPC(ctx context.Context, id string) error
	ListVPCs(ctx context.Context, region string) ([]*awsnetworking.VPC, error)
	
	// Subnet operations
	CreateSubnet(ctx context.Context, subnet *awsnetworking.Subnet) (*awsnetworking.Subnet, error)
	GetSubnet(ctx context.Context, id string) (*awsnetworking.Subnet, error)
	UpdateSubnet(ctx context.Context, subnet *awsnetworking.Subnet) (*awsnetworking.Subnet, error)
	DeleteSubnet(ctx context.Context, id string) error
	ListSubnets(ctx context.Context, vpcID string) ([]*awsnetworking.Subnet, error)
	
	// Internet Gateway operations
	CreateInternetGateway(ctx context.Context, igw *awsnetworking.InternetGateway) (*awsnetworking.InternetGateway, error)
	AttachInternetGateway(ctx context.Context, igwID, vpcID string) error
	DetachInternetGateway(ctx context.Context, igwID, vpcID string) error
	DeleteInternetGateway(ctx context.Context, id string) error
	
	// Route Table operations
	CreateRouteTable(ctx context.Context, rt *awsnetworking.RouteTable) (*awsnetworking.RouteTable, error)
	GetRouteTable(ctx context.Context, id string) (*awsnetworking.RouteTable, error)
	AssociateRouteTable(ctx context.Context, rtID, subnetID string) error
	DisassociateRouteTable(ctx context.Context, associationID string) error
	DeleteRouteTable(ctx context.Context, id string) error
	
	// Security Group operations
	CreateSecurityGroup(ctx context.Context, sg *awsnetworking.SecurityGroup) (*awsnetworking.SecurityGroup, error)
	GetSecurityGroup(ctx context.Context, id string) (*awsnetworking.SecurityGroup, error)
	UpdateSecurityGroup(ctx context.Context, sg *awsnetworking.SecurityGroup) (*awsnetworking.SecurityGroup, error)
	DeleteSecurityGroup(ctx context.Context, id string) error
	
	// NAT Gateway operations
	CreateNATGateway(ctx context.Context, ngw *awsnetworking.NATGateway) (*awsnetworking.NATGateway, error)
	GetNATGateway(ctx context.Context, id string) (*awsnetworking.NATGateway, error)
	DeleteNATGateway(ctx context.Context, id string) error
}
