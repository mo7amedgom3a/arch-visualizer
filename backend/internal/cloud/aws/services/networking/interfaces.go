package networking

import (
	"context"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
)

// AWSNetworkingService defines AWS-specific networking operations
// This implements cloud provider-specific logic while maintaining domain compatibility
type AWSNetworkingService interface {
	// VPC operations
	CreateVPC(ctx context.Context, vpc *awsnetworking.VPC) (*awsoutputs.VPCOutput, error)
	GetVPC(ctx context.Context, id string) (*awsoutputs.VPCOutput, error)
	UpdateVPC(ctx context.Context, vpc *awsnetworking.VPC) (*awsoutputs.VPCOutput, error)
	DeleteVPC(ctx context.Context, id string) error
	ListVPCs(ctx context.Context, region string) ([]*awsoutputs.VPCOutput, error)
	
	// Subnet operations
	CreateSubnet(ctx context.Context, subnet *awsnetworking.Subnet) (*awsoutputs.SubnetOutput, error)
	GetSubnet(ctx context.Context, id string) (*awsoutputs.SubnetOutput, error)
	UpdateSubnet(ctx context.Context, subnet *awsnetworking.Subnet) (*awsoutputs.SubnetOutput, error)
	DeleteSubnet(ctx context.Context, id string) error
	ListSubnets(ctx context.Context, vpcID string) ([]*awsoutputs.SubnetOutput, error)
	
	// Internet Gateway operations
	CreateInternetGateway(ctx context.Context, igw *awsnetworking.InternetGateway) (*awsoutputs.InternetGatewayOutput, error)
	AttachInternetGateway(ctx context.Context, igwID, vpcID string) error
	DetachInternetGateway(ctx context.Context, igwID, vpcID string) error
	DeleteInternetGateway(ctx context.Context, id string) error
	
	// Route Table operations
	CreateRouteTable(ctx context.Context, rt *awsnetworking.RouteTable) (*awsoutputs.RouteTableOutput, error)
	GetRouteTable(ctx context.Context, id string) (*awsoutputs.RouteTableOutput, error)
	AssociateRouteTable(ctx context.Context, rtID, subnetID string) error
	DisassociateRouteTable(ctx context.Context, associationID string) error
	DeleteRouteTable(ctx context.Context, id string) error
	
	// Security Group operations
	CreateSecurityGroup(ctx context.Context, sg *awsnetworking.SecurityGroup) (*awsoutputs.SecurityGroupOutput, error)
	GetSecurityGroup(ctx context.Context, id string) (*awsoutputs.SecurityGroupOutput, error)
	UpdateSecurityGroup(ctx context.Context, sg *awsnetworking.SecurityGroup) (*awsoutputs.SecurityGroupOutput, error)
	DeleteSecurityGroup(ctx context.Context, id string) error
	
	// NAT Gateway operations
	CreateNATGateway(ctx context.Context, ngw *awsnetworking.NATGateway) (*awsoutputs.NATGatewayOutput, error)
	GetNATGateway(ctx context.Context, id string) (*awsoutputs.NATGatewayOutput, error)
	DeleteNATGateway(ctx context.Context, id string) error
}
