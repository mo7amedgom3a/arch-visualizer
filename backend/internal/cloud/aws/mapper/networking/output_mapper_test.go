package networking

import (
	"fmt"
	"testing"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
)

// TestToDomainVPCFromOutput_RealisticData tests VPC output to domain mapping with realistic data
func TestToDomainVPCFromOutput_RealisticData(t *testing.T) {
	tests := []struct {
		name        string
		description string
		output      *awsoutputs.VPCOutput
		expected    *domainnetworking.VPC
	}{
		{
			name:        "realistic-vpc-output-to-domain",
			description: "Map realistic VPC output to domain model with ID and ARN",
			output: &awsoutputs.VPCOutput{
				ID:                 "vpc-0a1b2c3d4e5f6g7h8",
				ARN:                "arn:aws:ec2:us-east-1:123456789012:vpc/vpc-0a1b2c3d4e5f6g7h8",
				Name:               "production-vpc",
				Region:             "us-east-1",
				CIDR:               "10.0.0.0/16",
				State:              "available",
				IsDefault:          false,
				CreationTime:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				OwnerID:            "123456789012",
				EnableDNSHostnames: true,
				EnableDNSSupport:   true,
				InstanceTenancy:    "default",
				Tags: []configs.Tag{
					{Key: "Name", Value: "production-vpc"},
				},
			},
			expected: &domainnetworking.VPC{
				ID:                 "vpc-0a1b2c3d4e5f6g7h8",
				ARN:                stringPtr("arn:aws:ec2:us-east-1:123456789012:vpc/vpc-0a1b2c3d4e5f6g7h8"),
				Name:               "production-vpc",
				Region:             "us-east-1",
				CIDR:               "10.0.0.0/16",
				EnableDNS:          true,
				EnableDNSHostnames: true,
			},
		},
		{
			name:        "vpc-output-with-empty-arn",
			description: "Map VPC output with empty ARN to domain model (ARN should be nil)",
			output: &awsoutputs.VPCOutput{
				ID:                 "vpc-123",
				ARN:                "",
				Name:               "test-vpc",
				Region:             "us-east-1",
				CIDR:               "10.0.0.0/16",
				EnableDNSHostnames: true,
				EnableDNSSupport:   true,
			},
			expected: &domainnetworking.VPC{
				ID:                 "vpc-123",
				ARN:                nil,
				Name:               "test-vpc",
				Region:             "us-east-1",
				CIDR:               "10.0.0.0/16",
				EnableDNS:          true,
				EnableDNSHostnames: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("Input VPC Output ID: %s\n", test.output.ID)
			fmt.Printf("Input VPC Output ARN: %s\n", test.output.ARN)

			result := ToDomainVPCFromOutput(test.output)

			if result == nil {
				t.Fatal("Expected non-nil result, got nil")
			}

			fmt.Printf("Mapped Domain VPC ID: %s\n", result.ID)
			if result.ARN != nil {
				fmt.Printf("Mapped Domain VPC ARN: %s\n", *result.ARN)
			} else {
				fmt.Printf("Mapped Domain VPC ARN: nil\n")
			}

			// Validate ID
			if result.ID != test.expected.ID {
				t.Errorf("Expected ID %s, got %s", test.expected.ID, result.ID)
				fmt.Printf("❌ FAILED: ID mismatch\n")
			} else {
				fmt.Printf("✅ PASSED: ID mapping correct\n")
			}

			// Validate ARN
			if test.expected.ARN == nil {
				if result.ARN != nil {
					t.Errorf("Expected ARN to be nil, got %s", *result.ARN)
					fmt.Printf("❌ FAILED: ARN should be nil\n")
				} else {
					fmt.Printf("✅ PASSED: ARN is nil as expected\n")
				}
			} else {
				if result.ARN == nil {
					t.Errorf("Expected ARN %s, got nil", *test.expected.ARN)
					fmt.Printf("❌ FAILED: ARN is nil but expected %s\n", *test.expected.ARN)
				} else if *result.ARN != *test.expected.ARN {
					t.Errorf("Expected ARN %s, got %s", *test.expected.ARN, *result.ARN)
					fmt.Printf("❌ FAILED: ARN mismatch\n")
				} else {
					fmt.Printf("✅ PASSED: ARN mapping correct\n")
				}
			}

			// Validate other fields
			if result.Name != test.expected.Name {
				t.Errorf("Expected Name %s, got %s", test.expected.Name, result.Name)
			}
			if result.Region != test.expected.Region {
				t.Errorf("Expected Region %s, got %s", test.expected.Region, result.Region)
			}
			if result.CIDR != test.expected.CIDR {
				t.Errorf("Expected CIDR %s, got %s", test.expected.CIDR, result.CIDR)
			}
			if result.EnableDNS != test.expected.EnableDNS {
				t.Errorf("Expected EnableDNS %v, got %v", test.expected.EnableDNS, result.EnableDNS)
			}
			if result.EnableDNSHostnames != test.expected.EnableDNSHostnames {
				t.Errorf("Expected EnableDNSHostnames %v, got %v", test.expected.EnableDNSHostnames, result.EnableDNSHostnames)
			}

			fmt.Printf("✅ PASSED: All field mappings correct\n")
		})
	}
}

// TestToDomainSubnetFromOutput_RealisticData tests Subnet output to domain mapping
func TestToDomainSubnetFromOutput_RealisticData(t *testing.T) {
	tests := []struct {
		name        string
		description string
		output      *awsoutputs.SubnetOutput
		expected    *domainnetworking.Subnet
	}{
		{
			name:        "realistic-subnet-output-to-domain",
			description: "Map realistic Subnet output to domain model with ID and ARN",
			output: &awsoutputs.SubnetOutput{
				ID:                  "subnet-0a1b2c3d4e5f6g7h8",
				ARN:                 "arn:aws:ec2:us-east-1:123456789012:subnet/subnet-0a1b2c3d4e5f6g7h8",
				Name:                "public-subnet-1a",
				VPCID:               "vpc-0a1b2c3d4e5f6g7h8",
				CIDR:                "10.0.1.0/24",
				AvailabilityZone:    "us-east-1a",
				State:               "available",
				AvailableIPCount:    250,
				MapPublicIPOnLaunch: true,
				CreationTime:        time.Now(),
				Tags:                []configs.Tag{},
			},
			expected: &domainnetworking.Subnet{
				ID:               "subnet-0a1b2c3d4e5f6g7h8",
				ARN:              stringPtr("arn:aws:ec2:us-east-1:123456789012:subnet/subnet-0a1b2c3d4e5f6g7h8"),
				Name:             "public-subnet-1a",
				VPCID:            "vpc-0a1b2c3d4e5f6g7h8",
				CIDR:             "10.0.1.0/24",
				AvailabilityZone: stringPtr("us-east-1a"),
				IsPublic:         true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)

			result := ToDomainSubnetFromOutput(test.output)

			if result == nil {
				t.Fatal("Expected non-nil result, got nil")
			}

			// Validate ID and ARN
			if result.ID != test.expected.ID {
				t.Errorf("Expected ID %s, got %s", test.expected.ID, result.ID)
			}
			if test.expected.ARN != nil && (result.ARN == nil || *result.ARN != *test.expected.ARN) {
				t.Errorf("Expected ARN %s, got %v", *test.expected.ARN, result.ARN)
			}
			if result.IsPublic != test.expected.IsPublic {
				t.Errorf("Expected IsPublic %v, got %v", test.expected.IsPublic, result.IsPublic)
			}

			fmt.Printf("✅ PASSED: Subnet output to domain mapping succeeded\n")
		})
	}
}

// TestToDomainInternetGatewayFromOutput_RealisticData tests IGW output to domain mapping
func TestToDomainInternetGatewayFromOutput_RealisticData(t *testing.T) {
	tests := []struct {
		name        string
		description string
		output      *awsoutputs.InternetGatewayOutput
		expected    *domainnetworking.InternetGateway
	}{
		{
			name:        "realistic-igw-output-to-domain",
			description: "Map realistic Internet Gateway output to domain model",
			output: &awsoutputs.InternetGatewayOutput{
				ID:              "igw-0a1b2c3d4e5f6g7h8",
				ARN:             "arn:aws:ec2:us-east-1:123456789012:internet-gateway/igw-0a1b2c3d4e5f6g7h8",
				Name:            "production-igw",
				VPCID:           "vpc-0a1b2c3d4e5f6g7h8",
				State:           "available",
				AttachmentState: "attached",
				CreationTime:    time.Now(),
				Tags:            []configs.Tag{},
			},
			expected: &domainnetworking.InternetGateway{
				ID:    "igw-0a1b2c3d4e5f6g7h8",
				ARN:   stringPtr("arn:aws:ec2:us-east-1:123456789012:internet-gateway/igw-0a1b2c3d4e5f6g7h8"),
				Name:  "production-igw",
				VPCID: "vpc-0a1b2c3d4e5f6g7h8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)

			result := ToDomainInternetGatewayFromOutput(test.output)

			if result == nil {
				t.Fatal("Expected non-nil result, got nil")
			}

			if result.ID != test.expected.ID {
				t.Errorf("Expected ID %s, got %s", test.expected.ID, result.ID)
			}
			if result.Name != test.expected.Name {
				t.Errorf("Expected Name %s, got %s", test.expected.Name, result.Name)
			}
			if result.VPCID != test.expected.VPCID {
				t.Errorf("Expected VPCID %s, got %s", test.expected.VPCID, result.VPCID)
			}

			fmt.Printf("✅ PASSED: Internet Gateway output to domain mapping succeeded\n")
		})
	}
}

// TestToDomainNATGatewayFromOutput_RealisticData tests NAT Gateway output to domain mapping
func TestToDomainNATGatewayFromOutput_RealisticData(t *testing.T) {
	tests := []struct {
		name        string
		description string
		output      *awsoutputs.NATGatewayOutput
		expected    *domainnetworking.NATGateway
	}{
		{
			name:        "realistic-nat-output-to-domain",
			description: "Map realistic NAT Gateway output to domain model",
			output: &awsoutputs.NATGatewayOutput{
				ID:           "nat-0a1b2c3d4e5f6g7h8",
				ARN:          "arn:aws:ec2:us-east-1:123456789012:nat-gateway/nat-0a1b2c3d4e5f6g7h8",
				Name:         "production-nat",
				SubnetID:     "subnet-0a1b2c3d4e5f6g7h8",
				AllocationID: "eipalloc-0a1b2c3d4e5f6g7h8",
				State:        "available",
				PublicIP:     "54.123.45.67",
				PrivateIP:    "10.0.1.100",
				CreationTime: time.Now(),
				Tags:         []configs.Tag{},
			},
			expected: &domainnetworking.NATGateway{
				ID:           "nat-0a1b2c3d4e5f6g7h8",
				ARN:          stringPtr("arn:aws:ec2:us-east-1:123456789012:nat-gateway/nat-0a1b2c3d4e5f6g7h8"),
				Name:         "production-nat",
				SubnetID:     "subnet-0a1b2c3d4e5f6g7h8",
				AllocationID: stringPtr("eipalloc-0a1b2c3d4e5f6g7h8"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)

			result := ToDomainNATGatewayFromOutput(test.output)

			if result == nil {
				t.Fatal("Expected non-nil result, got nil")
			}

			if result.ID != test.expected.ID {
				t.Errorf("Expected ID %s, got %s", test.expected.ID, result.ID)
			}
			if result.SubnetID != test.expected.SubnetID {
				t.Errorf("Expected SubnetID %s, got %s", test.expected.SubnetID, result.SubnetID)
			}

			fmt.Printf("✅ PASSED: NAT Gateway output to domain mapping succeeded\n")
		})
	}
}

// TestToDomainRouteTableFromOutput_RealisticData tests Route Table output to domain mapping
func TestToDomainRouteTableFromOutput_RealisticData(t *testing.T) {
	igwID := "igw-0a1b2c3d4e5f6g7h8"
	natGatewayID := "nat-0a1b2c3d4e5f6g7h8"

	tests := []struct {
		name        string
		description string
		output      *awsoutputs.RouteTableOutput
		expected    *domainnetworking.RouteTable
	}{
		{
			name:        "realistic-route-table-output-to-domain",
			description: "Map realistic Route Table output to domain model with associations",
			output: &awsoutputs.RouteTableOutput{
				ID:    "rtb-0a1b2c3d4e5f6g7h8",
				ARN:   "arn:aws:ec2:us-east-1:123456789012:route-table/rtb-0a1b2c3d4e5f6g7h8",
				Name:  "public-route-table",
				VPCID: "vpc-0a1b2c3d4e5f6g7h8",
				Routes: []awsnetworking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						GatewayID:            &igwID,
					},
				},
				Associations: []awsoutputs.RouteTableAssociation{
					{
						ID:       "rtbassoc-0a1b2c3d4e5f6g7h8",
						SubnetID: "subnet-0a1b2c3d4e5f6g7h8",
						Main:     false,
					},
				},
				CreationTime: time.Now(),
				Tags:         []configs.Tag{},
			},
			expected: &domainnetworking.RouteTable{
				ID:    "rtb-0a1b2c3d4e5f6g7h8",
				ARN:   stringPtr("arn:aws:ec2:us-east-1:123456789012:route-table/rtb-0a1b2c3d4e5f6g7h8"),
				Name:  "public-route-table",
				VPCID: "vpc-0a1b2c3d4e5f6g7h8",
				Routes: []domainnetworking.Route{
					{
						DestinationCIDR: "0.0.0.0/0",
						TargetType:      "internet_gateway",
						TargetID:        igwID,
					},
				},
				Subnets: []string{"subnet-0a1b2c3d4e5f6g7h8"},
			},
		},
		{
			name:        "private-route-table-with-nat",
			description: "Map private route table with NAT Gateway route",
			output: &awsoutputs.RouteTableOutput{
				ID:    "rtb-private-123",
				ARN:   "arn:aws:ec2:us-east-1:123456789012:route-table/rtb-private-123",
				Name:  "private-route-table",
				VPCID: "vpc-123",
				Routes: []awsnetworking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						NatGatewayID:         &natGatewayID,
					},
				},
				Associations: []awsoutputs.RouteTableAssociation{
					{
						ID:       "rtbassoc-1",
						SubnetID: "subnet-private-1",
						Main:     false,
					},
					{
						ID:       "rtbassoc-2",
						SubnetID: "subnet-private-2",
						Main:     false,
					},
				},
				CreationTime: time.Now(),
				Tags:         []configs.Tag{},
			},
			expected: &domainnetworking.RouteTable{
				ID:    "rtb-private-123",
				ARN:   stringPtr("arn:aws:ec2:us-east-1:123456789012:route-table/rtb-private-123"),
				Name:  "private-route-table",
				VPCID: "vpc-123",
				Routes: []domainnetworking.Route{
					{
						DestinationCIDR: "0.0.0.0/0",
						TargetType:      "nat_gateway",
						TargetID:        natGatewayID,
					},
				},
				Subnets: []string{"subnet-private-1", "subnet-private-2"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)

			result := ToDomainRouteTableFromOutput(test.output)

			if result == nil {
				t.Fatal("Expected non-nil result, got nil")
			}

			// Validate basic fields
			if result.ID != test.expected.ID {
				t.Errorf("Expected ID %s, got %s", test.expected.ID, result.ID)
			}
			if len(result.Subnets) != len(test.expected.Subnets) {
				t.Errorf("Expected %d subnets, got %d", len(test.expected.Subnets), len(result.Subnets))
			}
			if len(result.Routes) != len(test.expected.Routes) {
				t.Errorf("Expected %d routes, got %d", len(test.expected.Routes), len(result.Routes))
			}

			fmt.Printf("✅ PASSED: Route Table output to domain mapping succeeded\n")
		})
	}
}

// TestToDomainSecurityGroupFromOutput_RealisticData tests Security Group output to domain mapping
func TestToDomainSecurityGroupFromOutput_RealisticData(t *testing.T) {
	tests := []struct {
		name        string
		description string
		output      *awsoutputs.SecurityGroupOutput
		expected    *domainnetworking.SecurityGroup
	}{
		{
			name:        "realistic-security-group-output-to-domain",
			description: "Map realistic Security Group output to domain model",
			output: &awsoutputs.SecurityGroupOutput{
				ID:          "sg-0a1b2c3d4e5f6g7h8",
				ARN:         "arn:aws:ec2:us-east-1:123456789012:security-group/sg-0a1b2c3d4e5f6g7h8",
				Name:        "web-server-sg",
				Description: "Security group for web servers",
				VPCID:       "vpc-0a1b2c3d4e5f6g7h8",
				Rules: []awsnetworking.SecurityGroupRule{
					{
						Type:        "ingress",
						Protocol:    "tcp",
						FromPort:    intPtr(80),
						ToPort:      intPtr(80),
						CIDRBlocks:  []string{"0.0.0.0/0"},
						Description: "Allow HTTP",
					},
				},
				CreationTime: time.Now(),
				Tags:         []configs.Tag{},
			},
			expected: &domainnetworking.SecurityGroup{
				ID:          "sg-0a1b2c3d4e5f6g7h8",
				ARN:         stringPtr("arn:aws:ec2:us-east-1:123456789012:security-group/sg-0a1b2c3d4e5f6g7h8"),
				Name:        "web-server-sg",
				Description: "Security group for web servers",
				VPCID:       "vpc-0a1b2c3d4e5f6g7h8",
				Rules: []domainnetworking.SecurityGroupRule{
					{
						Type:        "ingress",
						Protocol:    domainnetworking.ProtocolTCP,
						FromPort:    intPtr(80),
						ToPort:      intPtr(80),
						CIDRBlocks:  []string{"0.0.0.0/0"},
						Description: "Allow HTTP",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)

			result := ToDomainSecurityGroupFromOutput(test.output)

			if result == nil {
				t.Fatal("Expected non-nil result, got nil")
			}

			if result.ID != test.expected.ID {
				t.Errorf("Expected ID %s, got %s", test.expected.ID, result.ID)
			}
			if len(result.Rules) != len(test.expected.Rules) {
				t.Errorf("Expected %d rules, got %d", len(test.expected.Rules), len(result.Rules))
			}

			fmt.Printf("✅ PASSED: Security Group output to domain mapping succeeded\n")
		})
	}
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
