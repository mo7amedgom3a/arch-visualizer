package networking

import (
	"context"
	"errors"
	"testing"

	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
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

func (m *mockAWSNetworkingService) CreateVPC(ctx context.Context, vpc *awsnetworking.VPC) (*awsnetworking.VPC, error) {
	if m.createVPCError != nil {
		return nil, m.createVPCError
	}
	m.vpc = vpc
	return vpc, nil
}

func (m *mockAWSNetworkingService) GetVPC(ctx context.Context, id string) (*awsnetworking.VPC, error) {
	if m.getVPCError != nil {
		return nil, m.getVPCError
	}
	return m.vpc, nil
}

func (m *mockAWSNetworkingService) UpdateVPC(ctx context.Context, vpc *awsnetworking.VPC) (*awsnetworking.VPC, error) {
	m.vpc = vpc
	return vpc, nil
}

func (m *mockAWSNetworkingService) DeleteVPC(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSNetworkingService) ListVPCs(ctx context.Context, region string) ([]*awsnetworking.VPC, error) {
	if m.vpc != nil {
		return []*awsnetworking.VPC{m.vpc}, nil
	}
	return []*awsnetworking.VPC{}, nil
}

func (m *mockAWSNetworkingService) CreateSubnet(ctx context.Context, subnet *awsnetworking.Subnet) (*awsnetworking.Subnet, error) {
	m.subnet = subnet
	return subnet, nil
}

func (m *mockAWSNetworkingService) GetSubnet(ctx context.Context, id string) (*awsnetworking.Subnet, error) {
	return m.subnet, nil
}

func (m *mockAWSNetworkingService) UpdateSubnet(ctx context.Context, subnet *awsnetworking.Subnet) (*awsnetworking.Subnet, error) {
	m.subnet = subnet
	return subnet, nil
}

func (m *mockAWSNetworkingService) DeleteSubnet(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSNetworkingService) ListSubnets(ctx context.Context, vpcID string) ([]*awsnetworking.Subnet, error) {
	if m.subnet != nil {
		return []*awsnetworking.Subnet{m.subnet}, nil
	}
	return []*awsnetworking.Subnet{}, nil
}

func (m *mockAWSNetworkingService) CreateInternetGateway(ctx context.Context, igw *awsnetworking.InternetGateway) (*awsnetworking.InternetGateway, error) {
	m.igw = igw
	return igw, nil
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

func (m *mockAWSNetworkingService) CreateRouteTable(ctx context.Context, rt *awsnetworking.RouteTable) (*awsnetworking.RouteTable, error) {
	m.rt = rt
	return rt, nil
}

func (m *mockAWSNetworkingService) GetRouteTable(ctx context.Context, id string) (*awsnetworking.RouteTable, error) {
	return m.rt, nil
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

func (m *mockAWSNetworkingService) CreateSecurityGroup(ctx context.Context, sg *awsnetworking.SecurityGroup) (*awsnetworking.SecurityGroup, error) {
	m.sg = sg
	return sg, nil
}

func (m *mockAWSNetworkingService) GetSecurityGroup(ctx context.Context, id string) (*awsnetworking.SecurityGroup, error) {
	return m.sg, nil
}

func (m *mockAWSNetworkingService) UpdateSecurityGroup(ctx context.Context, sg *awsnetworking.SecurityGroup) (*awsnetworking.SecurityGroup, error) {
	m.sg = sg
	return sg, nil
}

func (m *mockAWSNetworkingService) DeleteSecurityGroup(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSNetworkingService) CreateNATGateway(ctx context.Context, ngw *awsnetworking.NATGateway) (*awsnetworking.NATGateway, error) {
	m.nat = ngw
	return ngw, nil
}

func (m *mockAWSNetworkingService) GetNATGateway(ctx context.Context, id string) (*awsnetworking.NATGateway, error) {
	return m.nat, nil
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
