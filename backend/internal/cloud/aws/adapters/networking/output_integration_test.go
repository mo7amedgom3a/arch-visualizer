package networking

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/networking"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
)

// realisticAWSNetworkingService is a realistic implementation that returns proper output models
type realisticAWSNetworkingService struct{}

var _ awsservice.AWSNetworkingService = (*realisticAWSNetworkingService)(nil)

func (s *realisticAWSNetworkingService) CreateVPC(ctx context.Context, vpc *awsnetworking.VPC) (*awsoutputs.VPCOutput, error) {
	// Simulate realistic AWS VPC creation
	return &awsoutputs.VPCOutput{
		ID:                 "vpc-0a1b2c3d4e5f6g7h8",
		ARN:                "arn:aws:ec2:us-east-1:123456789012:vpc/vpc-0a1b2c3d4e5f6g7h8",
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
	}, nil
}

func (s *realisticAWSNetworkingService) GetVPC(ctx context.Context, id string) (*awsoutputs.VPCOutput, error) {
	return &awsoutputs.VPCOutput{
		ID:                 id,
		ARN:                fmt.Sprintf("arn:aws:ec2:us-east-1:123456789012:vpc/%s", id),
		Name:               "test-vpc",
		Region:             "us-east-1",
		CIDR:               "10.0.0.0/16",
		State:              "available",
		IsDefault:          false,
		CreationTime:       time.Now(),
		OwnerID:            "123456789012",
		EnableDNSHostnames: true,
		EnableDNSSupport:   true,
		InstanceTenancy:    "default",
		Tags:               []configs.Tag{},
	}, nil
}

func (s *realisticAWSNetworkingService) UpdateVPC(ctx context.Context, vpc *awsnetworking.VPC) (*awsoutputs.VPCOutput, error) {
	return s.CreateVPC(ctx, vpc)
}

func (s *realisticAWSNetworkingService) DeleteVPC(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSNetworkingService) ListVPCs(ctx context.Context, region string) ([]*awsoutputs.VPCOutput, error) {
	return []*awsoutputs.VPCOutput{
		{
			ID:                 "vpc-0a1b2c3d4e5f6g7h8",
			ARN:                "arn:aws:ec2:us-east-1:123456789012:vpc/vpc-0a1b2c3d4e5f6g7h8",
			Name:               "test-vpc",
			Region:             region,
			CIDR:               "10.0.0.0/16",
			State:              "available",
			IsDefault:          false,
			CreationTime:       time.Now(),
			OwnerID:            "123456789012",
			EnableDNSHostnames: true,
			EnableDNSSupport:   true,
			InstanceTenancy:    "default",
			Tags:               []configs.Tag{},
		},
	}, nil
}

func (s *realisticAWSNetworkingService) CreateSubnet(ctx context.Context, subnet *awsnetworking.Subnet) (*awsoutputs.SubnetOutput, error) {
	return &awsoutputs.SubnetOutput{
		ID:                  "subnet-0a1b2c3d4e5f6g7h8",
		ARN:                 "arn:aws:ec2:us-east-1:123456789012:subnet/subnet-0a1b2c3d4e5f6g7h8",
		Name:                subnet.Name,
		VPCID:               subnet.VPCID,
		CIDR:                subnet.CIDR,
		AvailabilityZone:    subnet.AvailabilityZone,
		State:               "available",
		AvailableIPCount:    250,
		MapPublicIPOnLaunch: subnet.MapPublicIPOnLaunch,
		CreationTime:        time.Now(),
		Tags:                subnet.Tags,
	}, nil
}

func (s *realisticAWSNetworkingService) GetSubnet(ctx context.Context, id string) (*awsoutputs.SubnetOutput, error) {
	return &awsoutputs.SubnetOutput{
		ID:                  id,
		ARN:                 fmt.Sprintf("arn:aws:ec2:us-east-1:123456789012:subnet/%s", id),
		Name:                "test-subnet",
		VPCID:               "vpc-123",
		CIDR:                "10.0.1.0/24",
		AvailabilityZone:    "us-east-1a",
		State:               "available",
		AvailableIPCount:    250,
		MapPublicIPOnLaunch: true,
		CreationTime:        time.Now(),
		Tags:                []configs.Tag{},
	}, nil
}

func (s *realisticAWSNetworkingService) UpdateSubnet(ctx context.Context, subnet *awsnetworking.Subnet) (*awsoutputs.SubnetOutput, error) {
	return s.CreateSubnet(ctx, subnet)
}

func (s *realisticAWSNetworkingService) DeleteSubnet(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSNetworkingService) ListSubnets(ctx context.Context, vpcID string) ([]*awsoutputs.SubnetOutput, error) {
	return []*awsoutputs.SubnetOutput{
		{
			ID:                  "subnet-0a1b2c3d4e5f6g7h8",
			ARN:                 "arn:aws:ec2:us-east-1:123456789012:subnet/subnet-0a1b2c3d4e5f6g7h8",
			Name:                "test-subnet",
			VPCID:               vpcID,
			CIDR:                "10.0.1.0/24",
			AvailabilityZone:    "us-east-1a",
			State:               "available",
			AvailableIPCount:    250,
			MapPublicIPOnLaunch: true,
			CreationTime:        time.Now(),
			Tags:                []configs.Tag{},
		},
	}, nil
}

func (s *realisticAWSNetworkingService) CreateInternetGateway(ctx context.Context, igw *awsnetworking.InternetGateway) (*awsoutputs.InternetGatewayOutput, error) {
	return &awsoutputs.InternetGatewayOutput{
		ID:              "igw-0a1b2c3d4e5f6g7h8",
		ARN:             "arn:aws:ec2:us-east-1:123456789012:internet-gateway/igw-0a1b2c3d4e5f6g7h8",
		Name:            igw.Name,
		VPCID:           igw.VPCID,
		State:           "available",
		AttachmentState: "attached",
		CreationTime:    time.Now(),
		Tags:            igw.Tags,
	}, nil
}

func (s *realisticAWSNetworkingService) AttachInternetGateway(ctx context.Context, igwID, vpcID string) error {
	return nil
}

func (s *realisticAWSNetworkingService) DetachInternetGateway(ctx context.Context, igwID, vpcID string) error {
	return nil
}

func (s *realisticAWSNetworkingService) DeleteInternetGateway(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSNetworkingService) CreateRouteTable(ctx context.Context, rt *awsnetworking.RouteTable) (*awsoutputs.RouteTableOutput, error) {
	return &awsoutputs.RouteTableOutput{
		ID:           "rtb-0a1b2c3d4e5f6g7h8",
		ARN:          "arn:aws:ec2:us-east-1:123456789012:route-table/rtb-0a1b2c3d4e5f6g7h8",
		Name:         rt.Name,
		VPCID:        rt.VPCID,
		Routes:       rt.Routes,
		Associations: []awsoutputs.RouteTableAssociation{},
		CreationTime: time.Now(),
		Tags:         rt.Tags,
	}, nil
}

func (s *realisticAWSNetworkingService) GetRouteTable(ctx context.Context, id string) (*awsoutputs.RouteTableOutput, error) {
	igwID := "igw-123"
	return &awsoutputs.RouteTableOutput{
		ID:    id,
		ARN:   fmt.Sprintf("arn:aws:ec2:us-east-1:123456789012:route-table/%s", id),
		Name:  "test-route-table",
		VPCID: "vpc-123",
		Routes: []awsnetworking.Route{
			{
				DestinationCIDRBlock: "0.0.0.0/0",
				GatewayID:            &igwID,
			},
		},
		Associations: []awsoutputs.RouteTableAssociation{},
		CreationTime: time.Now(),
		Tags:         []configs.Tag{},
	}, nil
}

func (s *realisticAWSNetworkingService) AssociateRouteTable(ctx context.Context, rtID, subnetID string) error {
	return nil
}

func (s *realisticAWSNetworkingService) DisassociateRouteTable(ctx context.Context, associationID string) error {
	return nil
}

func (s *realisticAWSNetworkingService) DeleteRouteTable(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSNetworkingService) CreateSecurityGroup(ctx context.Context, sg *awsnetworking.SecurityGroup) (*awsoutputs.SecurityGroupOutput, error) {
	return &awsoutputs.SecurityGroupOutput{
		ID:           "sg-0a1b2c3d4e5f6g7h8",
		ARN:          "arn:aws:ec2:us-east-1:123456789012:security-group/sg-0a1b2c3d4e5f6g7h8",
		Name:         sg.Name,
		Description:  sg.Description,
		VPCID:        sg.VPCID,
		Rules:        sg.Rules,
		CreationTime: time.Now(),
		Tags:         sg.Tags,
	}, nil
}

func (s *realisticAWSNetworkingService) GetSecurityGroup(ctx context.Context, id string) (*awsoutputs.SecurityGroupOutput, error) {
	return &awsoutputs.SecurityGroupOutput{
		ID:           id,
		ARN:          fmt.Sprintf("arn:aws:ec2:us-east-1:123456789012:security-group/%s", id),
		Name:         "test-sg",
		Description:  "Test security group",
		VPCID:        "vpc-123",
		Rules:        []awsnetworking.SecurityGroupRule{},
		CreationTime: time.Now(),
		Tags:         []configs.Tag{},
	}, nil
}

func (s *realisticAWSNetworkingService) UpdateSecurityGroup(ctx context.Context, sg *awsnetworking.SecurityGroup) (*awsoutputs.SecurityGroupOutput, error) {
	return s.CreateSecurityGroup(ctx, sg)
}

func (s *realisticAWSNetworkingService) DeleteSecurityGroup(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSNetworkingService) CreateNATGateway(ctx context.Context, ngw *awsnetworking.NATGateway) (*awsoutputs.NATGatewayOutput, error) {
	return &awsoutputs.NATGatewayOutput{
		ID:           "nat-0a1b2c3d4e5f6g7h8",
		ARN:          "arn:aws:ec2:us-east-1:123456789012:nat-gateway/nat-0a1b2c3d4e5f6g7h8",
		Name:         ngw.Name,
		SubnetID:     ngw.SubnetID,
		AllocationID: ngw.AllocationID,
		State:        "available",
		PublicIP:     "54.123.45.67",
		PrivateIP:    "10.0.1.100",
		CreationTime: time.Now(),
		Tags:         ngw.Tags,
	}, nil
}

func (s *realisticAWSNetworkingService) GetNATGateway(ctx context.Context, id string) (*awsoutputs.NATGatewayOutput, error) {
	return &awsoutputs.NATGatewayOutput{
		ID:           id,
		ARN:          fmt.Sprintf("arn:aws:ec2:us-east-1:123456789012:nat-gateway/%s", id),
		Name:         "test-nat",
		SubnetID:     "subnet-123",
		AllocationID: "eipalloc-123",
		State:        "available",
		PublicIP:     "54.123.45.67",
		PrivateIP:    "10.0.1.100",
		CreationTime: time.Now(),
		Tags:         []configs.Tag{},
	}, nil
}

func (s *realisticAWSNetworkingService) DeleteNATGateway(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSNetworkingService) AllocateElasticIP(ctx context.Context, eip *awsnetworking.ElasticIP) (*awsoutputs.ElasticIPOutput, error) {
	allocationID := "eipalloc-0a1b2c3d4e5f6g7h8"
	if eip.AllocationID != nil && *eip.AllocationID != "" {
		allocationID = *eip.AllocationID
	}

	return &awsoutputs.ElasticIPOutput{
		ID:                 allocationID,
		ARN:                fmt.Sprintf("arn:aws:ec2:%s:123456789012:elastic-ip/%s", eip.Region, allocationID),
		PublicIP:           "54.123.45.67",
		Region:             eip.Region,
		NetworkBorderGroup: eip.NetworkBorderGroup,
		AllocationID:       allocationID,
		State:              "available",
		Domain:             "vpc",
		CreationTime:       time.Now(),
		Tags:               eip.Tags,
	}, nil
}

func (s *realisticAWSNetworkingService) GetElasticIP(ctx context.Context, id string) (*awsoutputs.ElasticIPOutput, error) {
	return &awsoutputs.ElasticIPOutput{
		ID:           id,
		ARN:          fmt.Sprintf("arn:aws:ec2:us-east-1:123456789012:elastic-ip/%s", id),
		PublicIP:     "54.123.45.67",
		Region:       "us-east-1",
		AllocationID: id,
		State:        "available",
		Domain:       "vpc",
		CreationTime: time.Now(),
		Tags:         []configs.Tag{},
	}, nil
}

func (s *realisticAWSNetworkingService) ReleaseElasticIP(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSNetworkingService) AssociateElasticIP(ctx context.Context, allocationID, instanceID string) error {
	return nil
}

func (s *realisticAWSNetworkingService) DisassociateElasticIP(ctx context.Context, associationID string) error {
	return nil
}

func (s *realisticAWSNetworkingService) ListElasticIPs(ctx context.Context, region string) ([]*awsoutputs.ElasticIPOutput, error) {
	return []*awsoutputs.ElasticIPOutput{
		{
			ID:           "eipalloc-0a1b2c3d4e5f6g7h8",
			ARN:          fmt.Sprintf("arn:aws:ec2:%s:123456789012:elastic-ip/eipalloc-0a1b2c3d4e5f6g7h8", region),
			PublicIP:     "54.123.45.67",
			Region:       region,
			AllocationID: "eipalloc-0a1b2c3d4e5f6g7h8",
			State:        "available",
			Domain:       "vpc",
			CreationTime: time.Now(),
			Tags:         []configs.Tag{},
		},
	}, nil
}

func (s *realisticAWSNetworkingService) CreateNetworkACL(ctx context.Context, acl *awsnetworking.NetworkACL) (*awsoutputs.NetworkACLOutput, error) {
	return &awsoutputs.NetworkACLOutput{
		ID:            "acl-0a1b2c3d4e5f6g7h8",
		ARN:           "arn:aws:ec2:us-east-1:123456789012:network-acl/acl-0a1b2c3d4e5f6g7h8",
		Name:          acl.Name,
		VPCID:         acl.VPCID,
		InboundRules:  acl.InboundRules,
		OutboundRules: acl.OutboundRules,
		IsDefault:     false,
		Associations:  []awsoutputs.NetworkACLAssociation{},
		CreationTime:  time.Now(),
		Tags:          acl.Tags,
	}, nil
}

func (s *realisticAWSNetworkingService) GetNetworkACL(ctx context.Context, id string) (*awsoutputs.NetworkACLOutput, error) {
	return &awsoutputs.NetworkACLOutput{
		ID:            id,
		ARN:           fmt.Sprintf("arn:aws:ec2:us-east-1:123456789012:network-acl/%s", id),
		Name:          "test-acl",
		VPCID:         "vpc-123",
		InboundRules:  []awsnetworking.ACLRule{},
		OutboundRules: []awsnetworking.ACLRule{},
		IsDefault:     false,
		Associations:  []awsoutputs.NetworkACLAssociation{},
		CreationTime:  time.Now(),
		Tags:          []configs.Tag{},
	}, nil
}

func (s *realisticAWSNetworkingService) DeleteNetworkACL(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSNetworkingService) AddNetworkACLRule(ctx context.Context, aclID string, rule awsnetworking.ACLRule) error {
	return nil
}

func (s *realisticAWSNetworkingService) RemoveNetworkACLRule(ctx context.Context, aclID string, ruleNumber int, ruleType awsnetworking.ACLRuleType) error {
	return nil
}

func (s *realisticAWSNetworkingService) AssociateNetworkACLWithSubnet(ctx context.Context, aclID, subnetID string) error {
	return nil
}

func (s *realisticAWSNetworkingService) DisassociateNetworkACLFromSubnet(ctx context.Context, associationID string) error {
	return nil
}

func (s *realisticAWSNetworkingService) ListNetworkACLs(ctx context.Context, vpcID string) ([]*awsoutputs.NetworkACLOutput, error) {
	return []*awsoutputs.NetworkACLOutput{
		{
			ID:            "acl-0a1b2c3d4e5f6g7h8",
			ARN:           "arn:aws:ec2:us-east-1:123456789012:network-acl/acl-0a1b2c3d4e5f6g7h8",
			Name:          "test-acl",
			VPCID:         vpcID,
			InboundRules:  []awsnetworking.ACLRule{},
			OutboundRules: []awsnetworking.ACLRule{},
			IsDefault:     false,
			Associations:  []awsoutputs.NetworkACLAssociation{},
			CreationTime:  time.Now(),
			Tags:          []configs.Tag{},
		},
	}, nil
}

func (s *realisticAWSNetworkingService) CreateNetworkInterface(ctx context.Context, eni *awsnetworking.NetworkInterface) (*awsoutputs.NetworkInterfaceOutput, error) {
	privateIP := "10.0.1.100"
	if eni.PrivateIPv4Address != nil && *eni.PrivateIPv4Address != "" {
		privateIP = *eni.PrivateIPv4Address
	}

	return &awsoutputs.NetworkInterfaceOutput{
		ID:                   "eni-0a1b2c3d4e5f6g7h8",
		ARN:                  "arn:aws:ec2:us-east-1:123456789012:network-interface/eni-0a1b2c3d4e5f6g7h8",
		Description:          eni.Description,
		SubnetID:             eni.SubnetID,
		InterfaceType:        string(eni.InterfaceType),
		SecurityGroupIDs:     eni.SecurityGroupIDs,
		Status:               "available",
		VPCID:                "vpc-123",
		AvailabilityZone:     "us-east-1a",
		OwnerID:              "123456789012",
		RequesterManaged:     false,
		SourceDestCheck:      true,
		PrivateIPv4Address:   privateIP,
		MACAddress:           "0a:ff:fe:97:e7:61",
		IPv4PrefixDelegation: eni.IPv4PrefixDelegation,
		CreationTime:         time.Now(),
		Tags:                 eni.Tags,
	}, nil
}

func (s *realisticAWSNetworkingService) GetNetworkInterface(ctx context.Context, id string) (*awsoutputs.NetworkInterfaceOutput, error) {
	return &awsoutputs.NetworkInterfaceOutput{
		ID:                 id,
		ARN:                fmt.Sprintf("arn:aws:ec2:us-east-1:123456789012:network-interface/%s", id),
		SubnetID:           "subnet-123",
		InterfaceType:      "elastic",
		SecurityGroupIDs:   []string{"sg-123"},
		Status:             "available",
		VPCID:              "vpc-123",
		AvailabilityZone:   "us-east-1a",
		OwnerID:            "123456789012",
		RequesterManaged:   false,
		SourceDestCheck:    true,
		PrivateIPv4Address: "10.0.1.100",
		MACAddress:         "0a:ff:fe:97:e7:61",
		CreationTime:       time.Now(),
		Tags:               []configs.Tag{},
	}, nil
}

func (s *realisticAWSNetworkingService) DeleteNetworkInterface(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSNetworkingService) AttachNetworkInterface(ctx context.Context, eniID, instanceID string, deviceIndex int) error {
	return nil
}

func (s *realisticAWSNetworkingService) DetachNetworkInterface(ctx context.Context, attachmentID string) error {
	return nil
}

func (s *realisticAWSNetworkingService) AssignPrivateIPAddress(ctx context.Context, eniID, privateIP string) error {
	return nil
}

func (s *realisticAWSNetworkingService) UnassignPrivateIPAddress(ctx context.Context, eniID, privateIP string) error {
	return nil
}

func (s *realisticAWSNetworkingService) ListNetworkInterfaces(ctx context.Context, subnetID string) ([]*awsoutputs.NetworkInterfaceOutput, error) {
	return []*awsoutputs.NetworkInterfaceOutput{
		{
			ID:                 "eni-0a1b2c3d4e5f6g7h8",
			ARN:                "arn:aws:ec2:us-east-1:123456789012:network-interface/eni-0a1b2c3d4e5f6g7h8",
			SubnetID:           subnetID,
			InterfaceType:      "elastic",
			SecurityGroupIDs:   []string{"sg-123"},
			Status:             "available",
			VPCID:              "vpc-123",
			AvailabilityZone:   "us-east-1a",
			OwnerID:            "123456789012",
			RequesterManaged:   false,
			SourceDestCheck:    true,
			PrivateIPv4Address: "10.0.1.100",
			MACAddress:         "0a:ff:fe:97:e7:61",
			CreationTime:       time.Now(),
			Tags:               []configs.Tag{},
		},
	}, nil
}

// TestFullFlow_VPC_CreateToDomain tests the complete flow: domain input → AWS service → output → domain with ID/ARN
func TestFullFlow_VPC_CreateToDomain(t *testing.T) {
	fmt.Printf("\n=== Running Integration Test: Full Flow VPC Creation ===\n")

	realisticService := &realisticAWSNetworkingService{}
	adapter := NewAWSNetworkingAdapter(realisticService)

	// Step 1: Create domain VPC (input)
	domainVPC := &domainnetworking.VPC{
		Name:               "production-vpc",
		Region:             "us-east-1",
		CIDR:               "10.0.0.0/16",
		EnableDNS:          true,
		EnableDNSHostnames: true,
	}

	fmt.Printf("\nStep 1: Domain Input VPC\n")
	fmt.Printf("  Name: %s\n", domainVPC.Name)
	fmt.Printf("  Region: %s\n", domainVPC.Region)
	fmt.Printf("  CIDR: %s\n", domainVPC.CIDR)
	fmt.Printf("  ID: %s (empty before creation)\n", domainVPC.ID)
	fmt.Printf("  ARN: %v (nil before creation)\n", domainVPC.ARN)

	// Step 2: Create through adapter (simulates full flow)
	ctx := context.Background()
	createdVPC, err := adapter.CreateVPC(ctx, domainVPC)
	if err != nil {
		t.Fatalf("Failed to create VPC: %v", err)
	}

	fmt.Printf("\nStep 2: Domain Output VPC (after creation)\n")
	fmt.Printf("  Name: %s\n", createdVPC.Name)
	fmt.Printf("  Region: %s\n", createdVPC.Region)
	fmt.Printf("  CIDR: %s\n", createdVPC.CIDR)
	fmt.Printf("  ID: %s (populated from AWS output)\n", createdVPC.ID)
	if createdVPC.ARN != nil {
		fmt.Printf("  ARN: %s (populated from AWS output)\n", *createdVPC.ARN)
	} else {
		fmt.Printf("  ARN: nil\n")
	}

	// Step 3: Validate that ID and ARN are populated
	if createdVPC.ID == "" {
		t.Error("Expected VPC ID to be populated, got empty string")
		fmt.Printf("❌ FAILED: VPC ID is empty\n")
	} else {
		fmt.Printf("✅ PASSED: VPC ID is populated: %s\n", createdVPC.ID)
	}

	if createdVPC.ARN == nil {
		t.Error("Expected VPC ARN to be populated, got nil")
		fmt.Printf("❌ FAILED: VPC ARN is nil\n")
	} else {
		fmt.Printf("✅ PASSED: VPC ARN is populated: %s\n", *createdVPC.ARN)
	}

	// Step 4: Validate realistic AWS ID format
	if len(createdVPC.ID) < 4 || createdVPC.ID[:4] != "vpc-" {
		t.Errorf("Expected VPC ID to start with 'vpc-', got: %s", createdVPC.ID)
		fmt.Printf("❌ FAILED: Invalid VPC ID format\n")
	} else {
		fmt.Printf("✅ PASSED: VPC ID has correct format\n")
	}

	// Step 5: Validate realistic ARN format
	if createdVPC.ARN != nil {
		if len(*createdVPC.ARN) < 20 || (*createdVPC.ARN)[:4] != "arn:" {
			t.Errorf("Expected valid ARN format, got: %s", *createdVPC.ARN)
			fmt.Printf("❌ FAILED: Invalid ARN format\n")
		} else {
			fmt.Printf("✅ PASSED: ARN has correct format\n")
		}
	}

	fmt.Printf("\n✅ PASSED: Full flow test completed successfully\n")
}

// TestFullFlow_Subnet_CreateToDomain tests the complete flow for subnet creation
func TestFullFlow_Subnet_CreateToDomain(t *testing.T) {
	fmt.Printf("\n=== Running Integration Test: Full Flow Subnet Creation ===\n")

	realisticService := &realisticAWSNetworkingService{}
	adapter := NewAWSNetworkingAdapter(realisticService)

	az := "us-east-1a"
	domainSubnet := &domainnetworking.Subnet{
		Name:             "public-subnet-1a",
		VPCID:            "vpc-0a1b2c3d4e5f6g7h8",
		CIDR:             "10.0.1.0/24",
		AvailabilityZone: &az,
		IsPublic:         true,
	}

	fmt.Printf("\nStep 1: Domain Input Subnet\n")
	fmt.Printf("  Name: %s\n", domainSubnet.Name)
	fmt.Printf("  VPC ID: %s\n", domainSubnet.VPCID)
	fmt.Printf("  CIDR: %s\n", domainSubnet.CIDR)
	fmt.Printf("  ID: %s (empty before creation)\n", domainSubnet.ID)

	ctx := context.Background()
	createdSubnet, err := adapter.CreateSubnet(ctx, domainSubnet)
	if err != nil {
		t.Fatalf("Failed to create subnet: %v", err)
	}

	fmt.Printf("\nStep 2: Domain Output Subnet (after creation)\n")
	fmt.Printf("  Name: %s\n", createdSubnet.Name)
	fmt.Printf("  ID: %s (populated from AWS output)\n", createdSubnet.ID)
	if createdSubnet.ARN != nil {
		fmt.Printf("  ARN: %s (populated from AWS output)\n", *createdSubnet.ARN)
	}

	// Validate
	if createdSubnet.ID == "" {
		t.Error("Expected Subnet ID to be populated, got empty string")
	} else {
		fmt.Printf("✅ PASSED: Subnet ID is populated\n")
	}

	if createdSubnet.ARN == nil {
		t.Error("Expected Subnet ARN to be populated, got nil")
	} else {
		fmt.Printf("✅ PASSED: Subnet ARN is populated\n")
	}

	fmt.Printf("\n✅ PASSED: Full flow test completed successfully\n")
}

// TestFullFlow_InternetGateway_CreateToDomain tests the complete flow for IGW creation
func TestFullFlow_InternetGateway_CreateToDomain(t *testing.T) {
	fmt.Printf("\n=== Running Integration Test: Full Flow Internet Gateway Creation ===\n")

	realisticService := &realisticAWSNetworkingService{}
	adapter := NewAWSNetworkingAdapter(realisticService)

	domainIGW := &domainnetworking.InternetGateway{
		Name:  "production-igw",
		VPCID: "vpc-0a1b2c3d4e5f6g7h8",
	}

	fmt.Printf("\nStep 1: Domain Input Internet Gateway\n")
	fmt.Printf("  Name: %s\n", domainIGW.Name)
	fmt.Printf("  VPC ID: %s\n", domainIGW.VPCID)
	fmt.Printf("  ID: %s (empty before creation)\n", domainIGW.ID)

	ctx := context.Background()
	createdIGW, err := adapter.CreateInternetGateway(ctx, domainIGW)
	if err != nil {
		t.Fatalf("Failed to create Internet Gateway: %v", err)
	}

	fmt.Printf("\nStep 2: Domain Output Internet Gateway (after creation)\n")
	fmt.Printf("  Name: %s\n", createdIGW.Name)
	fmt.Printf("  ID: %s (populated from AWS output)\n", createdIGW.ID)
	if createdIGW.ARN != nil {
		fmt.Printf("  ARN: %s (populated from AWS output)\n", *createdIGW.ARN)
	}

	// Validate
	if createdIGW.ID == "" {
		t.Error("Expected IGW ID to be populated, got empty string")
	} else {
		fmt.Printf("✅ PASSED: IGW ID is populated\n")
	}

	if createdIGW.ARN == nil {
		t.Error("Expected IGW ARN to be populated, got nil")
	} else {
		fmt.Printf("✅ PASSED: IGW ARN is populated\n")
	}

	fmt.Printf("\n✅ PASSED: Full flow test completed successfully\n")
}

// TestFullFlow_GetVPC_WithOutput tests getting a VPC and verifying output model mapping
func TestFullFlow_GetVPC_WithOutput(t *testing.T) {
	fmt.Printf("\n=== Running Integration Test: Get VPC with Output Model ===\n")

	realisticService := &realisticAWSNetworkingService{}
	adapter := NewAWSNetworkingAdapter(realisticService)

	ctx := context.Background()
	vpcID := "vpc-0a1b2c3d4e5f6g7h8"

	fmt.Printf("\nStep 1: Get VPC by ID\n")
	fmt.Printf("  VPC ID: %s\n", vpcID)

	vpc, err := adapter.GetVPC(ctx, vpcID)
	if err != nil {
		t.Fatalf("Failed to get VPC: %v", err)
	}

	fmt.Printf("\nStep 2: Domain VPC (from output model)\n")
	fmt.Printf("  ID: %s\n", vpc.ID)
	if vpc.ARN != nil {
		fmt.Printf("  ARN: %s\n", *vpc.ARN)
	}
	fmt.Printf("  Name: %s\n", vpc.Name)
	fmt.Printf("  Region: %s\n", vpc.Region)
	fmt.Printf("  CIDR: %s\n", vpc.CIDR)

	// Validate
	if vpc.ID != vpcID {
		t.Errorf("Expected VPC ID %s, got %s", vpcID, vpc.ID)
	} else {
		fmt.Printf("✅ PASSED: VPC ID matches\n")
	}

	if vpc.ARN == nil {
		t.Error("Expected VPC ARN to be populated, got nil")
	} else {
		fmt.Printf("✅ PASSED: VPC ARN is populated\n")
	}

	fmt.Printf("\n✅ PASSED: Get VPC test completed successfully\n")
}
