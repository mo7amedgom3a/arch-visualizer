package networking

import (
	"testing"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
)

func TestVPCMapper(t *testing.T) {
	// Test domain -> AWS mapping
	domainVPC := &domainnetworking.VPC{
		Name:               "test-vpc",
		Region:             "us-east-1",
		CIDR:               "10.0.0.0/16",
		EnableDNS:          true,
		EnableDNSHostnames: true,
	}
	
	awsVPC := FromDomainVPC(domainVPC)
	if awsVPC == nil {
		t.Fatal("Expected AWS VPC, got nil")
	}
	
	if awsVPC.Name != domainVPC.Name {
		t.Errorf("Expected name %s, got %s", domainVPC.Name, awsVPC.Name)
	}
	if awsVPC.Region != domainVPC.Region {
		t.Errorf("Expected region %s, got %s", domainVPC.Region, awsVPC.Region)
	}
	if awsVPC.CIDR != domainVPC.CIDR {
		t.Errorf("Expected CIDR %s, got %s", domainVPC.CIDR, awsVPC.CIDR)
	}
	if awsVPC.EnableDNSSupport != domainVPC.EnableDNS {
		t.Errorf("Expected EnableDNSSupport %v, got %v", domainVPC.EnableDNS, awsVPC.EnableDNSSupport)
	}
	
	// Test AWS -> domain mapping
	convertedDomainVPC := ToDomainVPC(awsVPC)
	if convertedDomainVPC == nil {
		t.Fatal("Expected domain VPC, got nil")
	}
	
	if convertedDomainVPC.Name != domainVPC.Name {
		t.Errorf("Expected name %s, got %s", domainVPC.Name, convertedDomainVPC.Name)
	}
	if convertedDomainVPC.Region != domainVPC.Region {
		t.Errorf("Expected region %s, got %s", domainVPC.Region, convertedDomainVPC.Region)
	}
}

func TestSubnetMapper(t *testing.T) {
	az := "us-east-1a"
	domainSubnet := &domainnetworking.Subnet{
		Name:            "test-subnet",
		VPCID:           "vpc-123",
		CIDR:            "10.0.1.0/24",
		AvailabilityZone: &az,
		IsPublic:        true,
	}
	
	awsSubnet := FromDomainSubnet(domainSubnet, az)
	if awsSubnet == nil {
		t.Fatal("Expected AWS Subnet, got nil")
	}
	
	if awsSubnet.Name != domainSubnet.Name {
		t.Errorf("Expected name %s, got %s", domainSubnet.Name, awsSubnet.Name)
	}
	if awsSubnet.VPCID != domainSubnet.VPCID {
		t.Errorf("Expected VPCID %s, got %s", domainSubnet.VPCID, awsSubnet.VPCID)
	}
	if awsSubnet.MapPublicIPOnLaunch != domainSubnet.IsPublic {
		t.Errorf("Expected MapPublicIPOnLaunch %v, got %v", domainSubnet.IsPublic, awsSubnet.MapPublicIPOnLaunch)
	}
	
	// Test AWS -> domain mapping
	convertedDomainSubnet := ToDomainSubnet(awsSubnet)
	if convertedDomainSubnet == nil {
		t.Fatal("Expected domain Subnet, got nil")
	}
	
	if convertedDomainSubnet.Name != domainSubnet.Name {
		t.Errorf("Expected name %s, got %s", domainSubnet.Name, convertedDomainSubnet.Name)
	}
}

func TestInternetGatewayMapper(t *testing.T) {
	domainIGW := &domainnetworking.InternetGateway{
		Name:  "test-igw",
		VPCID: "vpc-123",
	}
	
	awsIGW := FromDomainInternetGateway(domainIGW)
	if awsIGW == nil {
		t.Fatal("Expected AWS IGW, got nil")
	}
	
	if awsIGW.Name != domainIGW.Name {
		t.Errorf("Expected name %s, got %s", domainIGW.Name, awsIGW.Name)
	}
	
	convertedDomainIGW := ToDomainInternetGateway(awsIGW)
	if convertedDomainIGW == nil {
		t.Fatal("Expected domain IGW, got nil")
	}
	
	if convertedDomainIGW.Name != domainIGW.Name {
		t.Errorf("Expected name %s, got %s", domainIGW.Name, convertedDomainIGW.Name)
	}
}
