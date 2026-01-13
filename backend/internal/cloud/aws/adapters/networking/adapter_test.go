package networking

import (
	"context"
	"errors"
	"testing"
	"time"

	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/networking"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
)

// mockAWSNetworkingService is a mock implementation of AWSNetworkingService for testing
type mockAWSNetworkingService struct {
	vpc            *awsnetworking.VPC
	subnet         *awsnetworking.Subnet
	igw            *awsnetworking.InternetGateway
	rt             *awsnetworking.RouteTable
	sg             *awsnetworking.SecurityGroup
	nat            *awsnetworking.NATGateway
	createVPCError error
	getVPCError    error
}

// Ensure mockAWSNetworkingService implements AWSNetworkingService
var _ awsservice.AWSNetworkingService = (*mockAWSNetworkingService)(nil)

// Helper function to convert VPC input to output
func vpcToOutput(vpc *awsnetworking.VPC) *awsoutputs.VPCOutput {
	if vpc == nil {
		return nil
	}
	return &awsoutputs.VPCOutput{
		ID:                 "vpc-mock-123",
		ARN:                "arn:aws:ec2:us-east-1:123456789012:vpc/vpc-mock-123",
		Name:               vpc.Name,
		Region:             vpc.Region,
		CIDR:               vpc.CIDR,
		State:              "available",
		IsDefault:          false,
		CreationTime:       time.Now(),
		OwnerID:            "123456789012",
		EnableDNSHostnames: vpc.EnableDNSHostnames,
		EnableDNSSupport:   vpc.EnableDNSSupport,
		InstanceTenancy:    vpc.InstanceTenancy,
		Tags:               vpc.Tags,
	}
}

// Helper function to convert Subnet input to output
func subnetToOutput(subnet *awsnetworking.Subnet) *awsoutputs.SubnetOutput {
	if subnet == nil {
		return nil
	}
	return &awsoutputs.SubnetOutput{
		ID:                  "subnet-mock-123",
		ARN:                 "arn:aws:ec2:us-east-1:123456789012:subnet/subnet-mock-123",
		Name:                subnet.Name,
		VPCID:               subnet.VPCID,
		CIDR:                subnet.CIDR,
		AvailabilityZone:    subnet.AvailabilityZone,
		State:               "available",
		AvailableIPCount:    250,
		MapPublicIPOnLaunch: subnet.MapPublicIPOnLaunch,
		CreationTime:        time.Now(),
		Tags:                subnet.Tags,
	}
}

// Helper function to convert InternetGateway input to output
func igwToOutput(igw *awsnetworking.InternetGateway) *awsoutputs.InternetGatewayOutput {
	if igw == nil {
		return nil
	}
	return &awsoutputs.InternetGatewayOutput{
		ID:              "igw-mock-123",
		ARN:             "arn:aws:ec2:us-east-1:123456789012:internet-gateway/igw-mock-123",
		Name:            igw.Name,
		VPCID:           igw.VPCID,
		State:           "available",
		AttachmentState: "attached",
		CreationTime:    time.Now(),
		Tags:            igw.Tags,
	}
}

// Helper function to convert RouteTable input to output
func routeTableToOutput(rt *awsnetworking.RouteTable) *awsoutputs.RouteTableOutput {
	if rt == nil {
		return nil
	}
	return &awsoutputs.RouteTableOutput{
		ID:           "rtb-mock-123",
		ARN:          "arn:aws:ec2:us-east-1:123456789012:route-table/rtb-mock-123",
		Name:         rt.Name,
		VPCID:        rt.VPCID,
		Routes:       rt.Routes,
		Associations: []awsoutputs.RouteTableAssociation{},
		CreationTime: time.Now(),
		Tags:         rt.Tags,
	}
}

// Helper function to convert SecurityGroup input to output
func securityGroupToOutput(sg *awsnetworking.SecurityGroup) *awsoutputs.SecurityGroupOutput {
	if sg == nil {
		return nil
	}
	return &awsoutputs.SecurityGroupOutput{
		ID:           "sg-mock-123",
		ARN:          "arn:aws:ec2:us-east-1:123456789012:security-group/sg-mock-123",
		Name:         sg.Name,
		Description:  sg.Description,
		VPCID:        sg.VPCID,
		Rules:        sg.Rules,
		CreationTime: time.Now(),
		Tags:         sg.Tags,
	}
}

// Helper function to convert NATGateway input to output
func natGatewayToOutput(ngw *awsnetworking.NATGateway) *awsoutputs.NATGatewayOutput {
	if ngw == nil {
		return nil
	}
	return &awsoutputs.NATGatewayOutput{
		ID:           "nat-mock-123",
		ARN:          "arn:aws:ec2:us-east-1:123456789012:nat-gateway/nat-mock-123",
		Name:         ngw.Name,
		SubnetID:     ngw.SubnetID,
		AllocationID: ngw.AllocationID,
		State:        "available",
		PublicIP:     "1.2.3.4",
		PrivateIP:    "10.0.1.100",
		CreationTime: time.Now(),
		Tags:         ngw.Tags,
	}
}

func (m *mockAWSNetworkingService) CreateVPC(ctx context.Context, vpc *awsnetworking.VPC) (*awsoutputs.VPCOutput, error) {
	if m.createVPCError != nil {
		return nil, m.createVPCError
	}
	m.vpc = vpc
	return vpcToOutput(vpc), nil
}

func (m *mockAWSNetworkingService) GetVPC(ctx context.Context, id string) (*awsoutputs.VPCOutput, error) {
	if m.getVPCError != nil {
		return nil, m.getVPCError
	}
	return vpcToOutput(m.vpc), nil
}

func (m *mockAWSNetworkingService) UpdateVPC(ctx context.Context, vpc *awsnetworking.VPC) (*awsoutputs.VPCOutput, error) {
	m.vpc = vpc
	return vpcToOutput(vpc), nil
}

func (m *mockAWSNetworkingService) DeleteVPC(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSNetworkingService) ListVPCs(ctx context.Context, region string) ([]*awsoutputs.VPCOutput, error) {
	if m.vpc != nil {
		return []*awsoutputs.VPCOutput{vpcToOutput(m.vpc)}, nil
	}
	return []*awsoutputs.VPCOutput{}, nil
}

func (m *mockAWSNetworkingService) CreateSubnet(ctx context.Context, subnet *awsnetworking.Subnet) (*awsoutputs.SubnetOutput, error) {
	m.subnet = subnet
	return subnetToOutput(subnet), nil
}

func (m *mockAWSNetworkingService) GetSubnet(ctx context.Context, id string) (*awsoutputs.SubnetOutput, error) {
	return subnetToOutput(m.subnet), nil
}

func (m *mockAWSNetworkingService) UpdateSubnet(ctx context.Context, subnet *awsnetworking.Subnet) (*awsoutputs.SubnetOutput, error) {
	m.subnet = subnet
	return subnetToOutput(subnet), nil
}

func (m *mockAWSNetworkingService) DeleteSubnet(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSNetworkingService) ListSubnets(ctx context.Context, vpcID string) ([]*awsoutputs.SubnetOutput, error) {
	if m.subnet != nil {
		return []*awsoutputs.SubnetOutput{subnetToOutput(m.subnet)}, nil
	}
	return []*awsoutputs.SubnetOutput{}, nil
}

func (m *mockAWSNetworkingService) CreateInternetGateway(ctx context.Context, igw *awsnetworking.InternetGateway) (*awsoutputs.InternetGatewayOutput, error) {
	m.igw = igw
	return igwToOutput(igw), nil
}

func (m *mockAWSNetworkingService) AttachInternetGateway(ctx context.Context, igwID, vpcID string) error {
	return nil
}

func (m *mockAWSNetworkingService) DetachInternetGateway(ctx context.Context, igwID, vpcID string) error {
	return nil
}

func (m *mockAWSNetworkingService) DeleteInternetGateway(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSNetworkingService) CreateRouteTable(ctx context.Context, rt *awsnetworking.RouteTable) (*awsoutputs.RouteTableOutput, error) {
	m.rt = rt
	return routeTableToOutput(rt), nil
}

func (m *mockAWSNetworkingService) GetRouteTable(ctx context.Context, id string) (*awsoutputs.RouteTableOutput, error) {
	return routeTableToOutput(m.rt), nil
}

func (m *mockAWSNetworkingService) AssociateRouteTable(ctx context.Context, rtID, subnetID string) error {
	return nil
}

func (m *mockAWSNetworkingService) DisassociateRouteTable(ctx context.Context, associationID string) error {
	return nil
}

func (m *mockAWSNetworkingService) DeleteRouteTable(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSNetworkingService) CreateSecurityGroup(ctx context.Context, sg *awsnetworking.SecurityGroup) (*awsoutputs.SecurityGroupOutput, error) {
	m.sg = sg
	return securityGroupToOutput(sg), nil
}

func (m *mockAWSNetworkingService) GetSecurityGroup(ctx context.Context, id string) (*awsoutputs.SecurityGroupOutput, error) {
	return securityGroupToOutput(m.sg), nil
}

func (m *mockAWSNetworkingService) UpdateSecurityGroup(ctx context.Context, sg *awsnetworking.SecurityGroup) (*awsoutputs.SecurityGroupOutput, error) {
	m.sg = sg
	return securityGroupToOutput(sg), nil
}

func (m *mockAWSNetworkingService) DeleteSecurityGroup(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSNetworkingService) CreateNATGateway(ctx context.Context, ngw *awsnetworking.NATGateway) (*awsoutputs.NATGatewayOutput, error) {
	m.nat = ngw
	return natGatewayToOutput(ngw), nil
}

func (m *mockAWSNetworkingService) GetNATGateway(ctx context.Context, id string) (*awsoutputs.NATGatewayOutput, error) {
	return natGatewayToOutput(m.nat), nil
}

func (m *mockAWSNetworkingService) DeleteNATGateway(ctx context.Context, id string) error {
	return nil
}

func TestAWSNetworkingAdapter_CreateVPC(t *testing.T) {
	mockService := &mockAWSNetworkingService{}
	adapter := NewAWSNetworkingAdapter(mockService)

	domainVPC := &domainnetworking.VPC{
		Name:               "test-vpc",
		Region:             "us-east-1",
		CIDR:               "10.0.0.0/16",
		EnableDNS:          true,
		EnableDNSHostnames: true,
	}

	ctx := context.Background()
	createdVPC, err := adapter.CreateVPC(ctx, domainVPC)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdVPC == nil {
		t.Fatal("Expected created VPC, got nil")
	}

	if createdVPC.Name != domainVPC.Name {
		t.Errorf("Expected name %s, got %s", domainVPC.Name, createdVPC.Name)
	}

	if createdVPC.Region != domainVPC.Region {
		t.Errorf("Expected region %s, got %s", domainVPC.Region, createdVPC.Region)
	}

	if createdVPC.CIDR != domainVPC.CIDR {
		t.Errorf("Expected CIDR %s, got %s", domainVPC.CIDR, createdVPC.CIDR)
	}
}

func TestAWSNetworkingAdapter_CreateVPC_ValidationError(t *testing.T) {
	mockService := &mockAWSNetworkingService{}
	adapter := NewAWSNetworkingAdapter(mockService)

	invalidVPC := &domainnetworking.VPC{
		Name:   "", // Invalid: empty name
		Region: "us-east-1",
		CIDR:   "10.0.0.0/16",
	}

	ctx := context.Background()
	_, err := adapter.CreateVPC(ctx, invalidVPC)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}

	if err.Error() == "" {
		t.Error("Expected error message, got empty string")
	}
}

func TestAWSNetworkingAdapter_GetVPC(t *testing.T) {
	mockService := &mockAWSNetworkingService{
		vpc: &awsnetworking.VPC{
			Name:               "test-vpc",
			Region:             "us-east-1",
			CIDR:               "10.0.0.0/16",
			EnableDNSSupport:   true,
			EnableDNSHostnames: true,
		},
	}
	adapter := NewAWSNetworkingAdapter(mockService)

	ctx := context.Background()
	vpc, err := adapter.GetVPC(ctx, "vpc-123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if vpc == nil {
		t.Fatal("Expected VPC, got nil")
	}

	if vpc.Name != "test-vpc" {
		t.Errorf("Expected name test-vpc, got %s", vpc.Name)
	}
}

func TestAWSNetworkingAdapter_CreateSubnet(t *testing.T) {
	mockService := &mockAWSNetworkingService{}
	adapter := NewAWSNetworkingAdapter(mockService)

	az := "us-east-1a"
	domainSubnet := &domainnetworking.Subnet{
		Name:             "test-subnet",
		VPCID:            "vpc-123",
		CIDR:             "10.0.1.0/24",
		AvailabilityZone: &az,
		IsPublic:         true,
	}

	ctx := context.Background()
	createdSubnet, err := adapter.CreateSubnet(ctx, domainSubnet)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdSubnet == nil {
		t.Fatal("Expected created subnet, got nil")
	}

	if createdSubnet.Name != domainSubnet.Name {
		t.Errorf("Expected name %s, got %s", domainSubnet.Name, createdSubnet.Name)
	}
}

func TestAWSNetworkingAdapter_ListVPCs(t *testing.T) {
	mockService := &mockAWSNetworkingService{
		vpc: &awsnetworking.VPC{
			Name:   "test-vpc",
			Region: "us-east-1",
			CIDR:   "10.0.0.0/16",
		},
	}
	adapter := NewAWSNetworkingAdapter(mockService)

	ctx := context.Background()
	vpcs, err := adapter.ListVPCs(ctx, "us-east-1")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(vpcs) != 1 {
		t.Errorf("Expected 1 VPC, got %d", len(vpcs))
	}

	if vpcs[0].Name != "test-vpc" {
		t.Errorf("Expected name test-vpc, got %s", vpcs[0].Name)
	}
}

func TestAWSNetworkingAdapter_ErrorHandling(t *testing.T) {
	mockService := &mockAWSNetworkingService{
		getVPCError: errors.New("aws service error"),
	}
	adapter := NewAWSNetworkingAdapter(mockService)

	ctx := context.Background()
	_, err := adapter.GetVPC(ctx, "vpc-123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Verify error is wrapped
	if err.Error() == "" {
		t.Error("Expected error message, got empty string")
	}
}
