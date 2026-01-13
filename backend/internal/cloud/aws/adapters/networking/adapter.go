package networking

import (
	"context"
	"fmt"

	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/networking"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/networking"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
)

// AWSNetworkingAdapter adapts AWS-specific networking service to domain networking service
// This implements the Adapter pattern, allowing the domain layer to work with cloud-specific implementations
type AWSNetworkingAdapter struct {
	awsService awsservice.AWSNetworkingService
}

// NewAWSNetworkingAdapter creates a new AWS networking adapter
func NewAWSNetworkingAdapter(awsService awsservice.AWSNetworkingService) domainnetworking.NetworkingService {
	return &AWSNetworkingAdapter{
		awsService: awsService,
	}
}

// Ensure AWSNetworkingAdapter implements NetworkingService
var _ domainnetworking.NetworkingService = (*AWSNetworkingAdapter)(nil)

// VPC Operations

func (a *AWSNetworkingAdapter) CreateVPC(ctx context.Context, vpc *domainnetworking.VPC) (*domainnetworking.VPC, error) {
	if err := vpc.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsVPC := awsmapper.FromDomainVPC(vpc)
	if err := awsVPC.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsVPCOutput, err := a.awsService.CreateVPC(ctx, awsVPC)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainVPCFromOutput(awsVPCOutput), nil
}

func (a *AWSNetworkingAdapter) GetVPC(ctx context.Context, id string) (*domainnetworking.VPC, error) {
	awsVPCOutput, err := a.awsService.GetVPC(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainVPCFromOutput(awsVPCOutput), nil
}

func (a *AWSNetworkingAdapter) UpdateVPC(ctx context.Context, vpc *domainnetworking.VPC) (*domainnetworking.VPC, error) {
	if err := vpc.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsVPC := awsmapper.FromDomainVPC(vpc)
	if err := awsVPC.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsVPCOutput, err := a.awsService.UpdateVPC(ctx, awsVPC)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainVPCFromOutput(awsVPCOutput), nil
}

func (a *AWSNetworkingAdapter) DeleteVPC(ctx context.Context, id string) error {
	if err := a.awsService.DeleteVPC(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSNetworkingAdapter) ListVPCs(ctx context.Context, region string) ([]*domainnetworking.VPC, error) {
	awsVPCOutputs, err := a.awsService.ListVPCs(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainVPCs := make([]*domainnetworking.VPC, len(awsVPCOutputs))
	for i, awsVPCOutput := range awsVPCOutputs {
		domainVPCs[i] = awsmapper.ToDomainVPCFromOutput(awsVPCOutput)
	}

	return domainVPCs, nil
}

// Subnet Operations

func (a *AWSNetworkingAdapter) CreateSubnet(ctx context.Context, subnet *domainnetworking.Subnet) (*domainnetworking.Subnet, error) {
	if err := subnet.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	// Get VPC to determine availability zone if not provided
	// For now, we'll require AZ to be set in the domain model
	// In a real implementation, you might fetch the VPC to get available AZs
	az := ""
	if subnet.AvailabilityZone != nil {
		az = *subnet.AvailabilityZone
	} else {
		return nil, fmt.Errorf("availability zone is required for subnet")
	}

	awsSubnet := awsmapper.FromDomainSubnet(subnet, az)
	if err := awsSubnet.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsSubnetOutput, err := a.awsService.CreateSubnet(ctx, awsSubnet)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainSubnetFromOutput(awsSubnetOutput), nil
}

func (a *AWSNetworkingAdapter) GetSubnet(ctx context.Context, id string) (*domainnetworking.Subnet, error) {
	awsSubnetOutput, err := a.awsService.GetSubnet(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainSubnetFromOutput(awsSubnetOutput), nil
}

func (a *AWSNetworkingAdapter) UpdateSubnet(ctx context.Context, subnet *domainnetworking.Subnet) (*domainnetworking.Subnet, error) {
	if err := subnet.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	az := ""
	if subnet.AvailabilityZone != nil {
		az = *subnet.AvailabilityZone
	}

	awsSubnet := awsmapper.FromDomainSubnet(subnet, az)
	if err := awsSubnet.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsSubnetOutput, err := a.awsService.UpdateSubnet(ctx, awsSubnet)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainSubnetFromOutput(awsSubnetOutput), nil
}

func (a *AWSNetworkingAdapter) DeleteSubnet(ctx context.Context, id string) error {
	if err := a.awsService.DeleteSubnet(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSNetworkingAdapter) ListSubnets(ctx context.Context, vpcID string) ([]*domainnetworking.Subnet, error) {
	awsSubnetOutputs, err := a.awsService.ListSubnets(ctx, vpcID)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainSubnets := make([]*domainnetworking.Subnet, len(awsSubnetOutputs))
	for i, awsSubnetOutput := range awsSubnetOutputs {
		domainSubnets[i] = awsmapper.ToDomainSubnetFromOutput(awsSubnetOutput)
	}

	return domainSubnets, nil
}

// Internet Gateway Operations

func (a *AWSNetworkingAdapter) CreateInternetGateway(ctx context.Context, igw *domainnetworking.InternetGateway) (*domainnetworking.InternetGateway, error) {
	if err := igw.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsIGW := awsmapper.FromDomainInternetGateway(igw)
	if err := awsIGW.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsIGWOutput, err := a.awsService.CreateInternetGateway(ctx, awsIGW)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainInternetGatewayFromOutput(awsIGWOutput), nil
}

func (a *AWSNetworkingAdapter) AttachInternetGateway(ctx context.Context, igwID, vpcID string) error {
	if err := a.awsService.AttachInternetGateway(ctx, igwID, vpcID); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSNetworkingAdapter) DetachInternetGateway(ctx context.Context, igwID, vpcID string) error {
	if err := a.awsService.DetachInternetGateway(ctx, igwID, vpcID); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSNetworkingAdapter) DeleteInternetGateway(ctx context.Context, id string) error {
	if err := a.awsService.DeleteInternetGateway(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

// Route Table Operations

func (a *AWSNetworkingAdapter) CreateRouteTable(ctx context.Context, rt *domainnetworking.RouteTable) (*domainnetworking.RouteTable, error) {
	if err := rt.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsRT := awsmapper.FromDomainRouteTable(rt)
	if err := awsRT.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsRTOutput, err := a.awsService.CreateRouteTable(ctx, awsRT)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainRouteTableFromOutput(awsRTOutput), nil
}

func (a *AWSNetworkingAdapter) GetRouteTable(ctx context.Context, id string) (*domainnetworking.RouteTable, error) {
	awsRTOutput, err := a.awsService.GetRouteTable(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainRouteTableFromOutput(awsRTOutput), nil
}

func (a *AWSNetworkingAdapter) AssociateRouteTable(ctx context.Context, rtID, subnetID string) error {
	if err := a.awsService.AssociateRouteTable(ctx, rtID, subnetID); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSNetworkingAdapter) DisassociateRouteTable(ctx context.Context, associationID string) error {
	if err := a.awsService.DisassociateRouteTable(ctx, associationID); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSNetworkingAdapter) DeleteRouteTable(ctx context.Context, id string) error {
	if err := a.awsService.DeleteRouteTable(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

// Security Group Operations

func (a *AWSNetworkingAdapter) CreateSecurityGroup(ctx context.Context, sg *domainnetworking.SecurityGroup) (*domainnetworking.SecurityGroup, error) {
	if err := sg.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsSG := awsmapper.FromDomainSecurityGroup(sg)
	if err := awsSG.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsSGOutput, err := a.awsService.CreateSecurityGroup(ctx, awsSG)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainSecurityGroupFromOutput(awsSGOutput), nil
}

func (a *AWSNetworkingAdapter) GetSecurityGroup(ctx context.Context, id string) (*domainnetworking.SecurityGroup, error) {
	awsSGOutput, err := a.awsService.GetSecurityGroup(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainSecurityGroupFromOutput(awsSGOutput), nil
}

func (a *AWSNetworkingAdapter) UpdateSecurityGroup(ctx context.Context, sg *domainnetworking.SecurityGroup) (*domainnetworking.SecurityGroup, error) {
	if err := sg.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsSG := awsmapper.FromDomainSecurityGroup(sg)
	if err := awsSG.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsSGOutput, err := a.awsService.UpdateSecurityGroup(ctx, awsSG)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainSecurityGroupFromOutput(awsSGOutput), nil
}

func (a *AWSNetworkingAdapter) DeleteSecurityGroup(ctx context.Context, id string) error {
	if err := a.awsService.DeleteSecurityGroup(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

// NAT Gateway Operations

func (a *AWSNetworkingAdapter) CreateNATGateway(ctx context.Context, ngw *domainnetworking.NATGateway) (*domainnetworking.NATGateway, error) {
	if err := ngw.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsNAT := awsmapper.FromDomainNATGateway(ngw)
	if err := awsNAT.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsNATOutput, err := a.awsService.CreateNATGateway(ctx, awsNAT)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainNATGatewayFromOutput(awsNATOutput), nil
}

func (a *AWSNetworkingAdapter) GetNATGateway(ctx context.Context, id string) (*domainnetworking.NATGateway, error) {
	awsNATOutput, err := a.awsService.GetNATGateway(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainNATGatewayFromOutput(awsNATOutput), nil
}

func (a *AWSNetworkingAdapter) DeleteNATGateway(ctx context.Context, id string) error {
	if err := a.awsService.DeleteNATGateway(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}
