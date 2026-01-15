package outputs

import (
	"fmt"
	"testing"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
)

// TestVPCOutput_RealisticData tests VPC output with realistic AWS data
func TestVPCOutput_RealisticData(t *testing.T) {
	tests := []struct {
		name        string
		description string
		vpcOutput   *VPCOutput
		validate    func(*VPCOutput) error
	}{
		{
			name:        "realistic-vpc-output",
			description: "VPC output with realistic AWS identifiers and metadata",
			vpcOutput: &VPCOutput{
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
					{Key: "Environment", Value: "production"},
					{Key: "ManagedBy", Value: "terraform"},
				},
			},
			validate: func(v *VPCOutput) error {
				if v.ID == "" {
					return fmt.Errorf("VPC ID is required")
				}
				if v.ARN == "" {
					return fmt.Errorf("VPC ARN is required")
				}
				if !isValidVPCID(v.ID) {
					return fmt.Errorf("invalid VPC ID format: %s", v.ID)
				}
				if !isValidARN(v.ARN, "vpc") {
					return fmt.Errorf("invalid VPC ARN format: %s", v.ARN)
				}
				if v.State != "available" && v.State != "pending" {
					return fmt.Errorf("invalid VPC state: %s", v.State)
				}
				return nil
			},
		},
		{
			name:        "default-vpc-output",
			description: "Default VPC output with IsDefault flag set",
			vpcOutput: &VPCOutput{
				ID:                 "vpc-default-123",
				ARN:                "arn:aws:ec2:us-east-1:123456789012:vpc/vpc-default-123",
				Name:               "default",
				Region:             "us-east-1",
				CIDR:               "172.31.0.0/16",
				State:              "available",
				IsDefault:          true,
				CreationTime:       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				OwnerID:            "123456789012",
				EnableDNSHostnames: true,
				EnableDNSSupport:   true,
				InstanceTenancy:    "default",
				Tags:               []configs.Tag{},
			},
			validate: func(v *VPCOutput) error {
				if !v.IsDefault {
					return fmt.Errorf("expected IsDefault to be true for default VPC")
				}
				return nil
			},
		},
		{
			name:        "pending-vpc-output",
			description: "VPC output in pending state",
			vpcOutput: &VPCOutput{
				ID:                 "vpc-pending-456",
				ARN:                "arn:aws:ec2:us-east-1:123456789012:vpc/vpc-pending-456",
				Name:               "new-vpc",
				Region:             "us-west-2",
				CIDR:               "10.1.0.0/16",
				State:              "pending",
				IsDefault:          false,
				CreationTime:       time.Now(),
				OwnerID:            "123456789012",
				EnableDNSHostnames: true,
				EnableDNSSupport:   true,
				InstanceTenancy:    "default",
				Tags:               []configs.Tag{{Key: "Name", Value: "new-vpc"}},
			},
			validate: func(v *VPCOutput) error {
				if v.State != "pending" {
					return fmt.Errorf("expected state to be pending, got: %s", v.State)
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("VPC ID: %s\n", test.vpcOutput.ID)
			fmt.Printf("VPC ARN: %s\n", test.vpcOutput.ARN)
			fmt.Printf("VPC Name: %s\n", test.vpcOutput.Name)
			fmt.Printf("VPC Region: %s\n", test.vpcOutput.Region)
			fmt.Printf("VPC CIDR: %s\n", test.vpcOutput.CIDR)
			fmt.Printf("VPC State: %s\n", test.vpcOutput.State)
			fmt.Printf("Is Default: %v\n", test.vpcOutput.IsDefault)
			fmt.Printf("Owner ID: %s\n", test.vpcOutput.OwnerID)
			fmt.Printf("Creation Time: %s\n", test.vpcOutput.CreationTime.Format(time.RFC3339))

			if err := test.validate(test.vpcOutput); err != nil {
				t.Errorf("Validation failed: %v", err)
				fmt.Printf("❌ FAILED: %v\n", err)
			} else {
				fmt.Printf("✅ PASSED: VPC output validation succeeded\n")
			}
		})
	}
}

// TestSubnetOutput_RealisticData tests Subnet output with realistic AWS data
func TestSubnetOutput_RealisticData(t *testing.T) {
	tests := []struct {
		name         string
		description  string
		subnetOutput *SubnetOutput
		validate     func(*SubnetOutput) error
	}{
		{
			name:        "realistic-subnet-output",
			description: "Subnet output with realistic AWS identifiers and metadata",
			subnetOutput: &SubnetOutput{
				ID:                  "subnet-0a1b2c3d4e5f6g7h8",
				ARN:                 "arn:aws:ec2:us-east-1:123456789012:subnet/subnet-0a1b2c3d4e5f6g7h8",
				Name:                "public-subnet-1a",
				VPCID:               "vpc-0a1b2c3d4e5f6g7h8",
				CIDR:                "10.0.1.0/24",
				AvailabilityZone:    "us-east-1a",
				State:               "available",
				AvailableIPCount:    250,
				MapPublicIPOnLaunch: true,
				CreationTime:        time.Date(2024, 1, 15, 10, 35, 0, 0, time.UTC),
				Tags: []configs.Tag{
					{Key: "Name", Value: "public-subnet-1a"},
					{Key: "Type", Value: "public"},
				},
			},
			validate: func(s *SubnetOutput) error {
				if s.ID == "" {
					return fmt.Errorf("Subnet ID is required")
				}
				if s.ARN == "" {
					return fmt.Errorf("Subnet ARN is required")
				}
				if !isValidSubnetID(s.ID) {
					return fmt.Errorf("invalid Subnet ID format: %s", s.ID)
				}
				if !isValidARN(s.ARN, "subnet") {
					return fmt.Errorf("invalid Subnet ARN format: %s", s.ARN)
				}
				if s.State != "available" && s.State != "pending" {
					return fmt.Errorf("invalid Subnet state: %s", s.State)
				}
				if s.AvailableIPCount < 0 {
					return fmt.Errorf("AvailableIPCount cannot be negative: %d", s.AvailableIPCount)
				}
				return nil
			},
		},
		{
			name:        "private-subnet-output",
			description: "Private subnet output with MapPublicIPOnLaunch false",
			subnetOutput: &SubnetOutput{
				ID:                  "subnet-private-1b",
				ARN:                 "arn:aws:ec2:us-east-1:123456789012:subnet/subnet-private-1b",
				Name:                "private-subnet-1b",
				VPCID:               "vpc-0a1b2c3d4e5f6g7h8",
				CIDR:                "10.0.2.0/24",
				AvailabilityZone:    "us-east-1b",
				State:               "available",
				AvailableIPCount:    250,
				MapPublicIPOnLaunch: false,
				CreationTime:        time.Date(2024, 1, 15, 10, 40, 0, 0, time.UTC),
				Tags: []configs.Tag{
					{Key: "Name", Value: "private-subnet-1b"},
					{Key: "Type", Value: "private"},
				},
			},
			validate: func(s *SubnetOutput) error {
				if s.MapPublicIPOnLaunch {
					return fmt.Errorf("expected MapPublicIPOnLaunch to be false for private subnet")
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("Subnet ID: %s\n", test.subnetOutput.ID)
			fmt.Printf("Subnet ARN: %s\n", test.subnetOutput.ARN)
			fmt.Printf("Subnet Name: %s\n", test.subnetOutput.Name)
			fmt.Printf("VPC ID: %s\n", test.subnetOutput.VPCID)
			fmt.Printf("CIDR: %s\n", test.subnetOutput.CIDR)
			fmt.Printf("Availability Zone: %s\n", test.subnetOutput.AvailabilityZone)
			fmt.Printf("State: %s\n", test.subnetOutput.State)
			fmt.Printf("Available IP Count: %d\n", test.subnetOutput.AvailableIPCount)
			fmt.Printf("Map Public IP: %v\n", test.subnetOutput.MapPublicIPOnLaunch)

			if err := test.validate(test.subnetOutput); err != nil {
				t.Errorf("Validation failed: %v", err)
				fmt.Printf("❌ FAILED: %v\n", err)
			} else {
				fmt.Printf("✅ PASSED: Subnet output validation succeeded\n")
			}
		})
	}
}

// TestInternetGatewayOutput_RealisticData tests Internet Gateway output with realistic AWS data
func TestInternetGatewayOutput_RealisticData(t *testing.T) {
	tests := []struct {
		name        string
		description string
		igwOutput   *InternetGatewayOutput
		validate    func(*InternetGatewayOutput) error
	}{
		{
			name:        "realistic-igw-output",
			description: "Internet Gateway output with realistic AWS identifiers and metadata",
			igwOutput: &InternetGatewayOutput{
				ID:              "igw-0a1b2c3d4e5f6g7h8",
				ARN:             "arn:aws:ec2:us-east-1:123456789012:internet-gateway/igw-0a1b2c3d4e5f6g7h8",
				Name:            "production-igw",
				VPCID:           "vpc-0a1b2c3d4e5f6g7h8",
				State:           "available",
				AttachmentState: "attached",
				CreationTime:    time.Date(2024, 1, 15, 10, 45, 0, 0, time.UTC),
				Tags: []configs.Tag{
					{Key: "Name", Value: "production-igw"},
				},
			},
			validate: func(i *InternetGatewayOutput) error {
				if i.ID == "" {
					return fmt.Errorf("IGW ID is required")
				}
				if i.ARN == "" {
					return fmt.Errorf("IGW ARN is required")
				}
				if !isValidIGWID(i.ID) {
					return fmt.Errorf("invalid IGW ID format: %s", i.ID)
				}
				if !isValidARN(i.ARN, "internet-gateway") {
					return fmt.Errorf("invalid IGW ARN format: %s", i.ARN)
				}
				if i.State != "available" && i.State != "attached" && i.State != "detaching" {
					return fmt.Errorf("invalid IGW state: %s", i.State)
				}
				return nil
			},
		},
		{
			name:        "detaching-igw-output",
			description: "Internet Gateway in detaching state",
			igwOutput: &InternetGatewayOutput{
				ID:              "igw-detaching-123",
				ARN:             "arn:aws:ec2:us-east-1:123456789012:internet-gateway/igw-detaching-123",
				Name:            "detaching-igw",
				VPCID:           "vpc-123",
				State:           "detaching",
				AttachmentState: "detaching",
				CreationTime:    time.Now(),
				Tags:            []configs.Tag{},
			},
			validate: func(i *InternetGatewayOutput) error {
				if i.State != "detaching" {
					return fmt.Errorf("expected state to be detaching, got: %s", i.State)
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("IGW ID: %s\n", test.igwOutput.ID)
			fmt.Printf("IGW ARN: %s\n", test.igwOutput.ARN)
			fmt.Printf("IGW Name: %s\n", test.igwOutput.Name)
			fmt.Printf("VPC ID: %s\n", test.igwOutput.VPCID)
			fmt.Printf("State: %s\n", test.igwOutput.State)
			fmt.Printf("Attachment State: %s\n", test.igwOutput.AttachmentState)

			if err := test.validate(test.igwOutput); err != nil {
				t.Errorf("Validation failed: %v", err)
				fmt.Printf("❌ FAILED: %v\n", err)
			} else {
				fmt.Printf("✅ PASSED: IGW output validation succeeded\n")
			}
		})
	}
}

// TestNATGatewayOutput_RealisticData tests NAT Gateway output with realistic AWS data
func TestNATGatewayOutput_RealisticData(t *testing.T) {
	tests := []struct {
		name        string
		description string
		natOutput   *NATGatewayOutput
		validate    func(*NATGatewayOutput) error
	}{
		{
			name:        "realistic-nat-output",
			description: "NAT Gateway output with realistic AWS identifiers and metadata",
			natOutput: &NATGatewayOutput{
				ID:           "nat-0a1b2c3d4e5f6g7h8",
				ARN:          "arn:aws:ec2:us-east-1:123456789012:nat-gateway/nat-0a1b2c3d4e5f6g7h8",
				Name:         "production-nat",
				SubnetID:     "subnet-0a1b2c3d4e5f6g7h8",
				AllocationID: "eipalloc-0a1b2c3d4e5f6g7h8",
				State:        "available",
				PublicIP:     "54.123.45.67",
				PrivateIP:    "10.0.1.100",
				CreationTime: time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
				Tags: []configs.Tag{
					{Key: "Name", Value: "production-nat"},
				},
			},
			validate: func(n *NATGatewayOutput) error {
				if n.ID == "" {
					return fmt.Errorf("NAT Gateway ID is required")
				}
				if n.ARN == "" {
					return fmt.Errorf("NAT Gateway ARN is required")
				}
				if !isValidNATGatewayID(n.ID) {
					return fmt.Errorf("invalid NAT Gateway ID format: %s", n.ID)
				}
				if !isValidARN(n.ARN, "nat-gateway") {
					return fmt.Errorf("invalid NAT Gateway ARN format: %s", n.ARN)
				}
				if n.State != "available" && n.State != "pending" && n.State != "failed" && n.State != "deleting" && n.State != "deleted" {
					return fmt.Errorf("invalid NAT Gateway state: %s", n.State)
				}
				if n.PublicIP == "" {
					return fmt.Errorf("PublicIP is required for NAT Gateway")
				}
				return nil
			},
		},
		{
			name:        "pending-nat-output",
			description: "NAT Gateway in pending state",
			natOutput: &NATGatewayOutput{
				ID:           "nat-pending-123",
				ARN:          "arn:aws:ec2:us-east-1:123456789012:nat-gateway/nat-pending-123",
				Name:         "pending-nat",
				SubnetID:     "subnet-123",
				AllocationID: "eipalloc-123",
				State:        "pending",
				PublicIP:     "",
				PrivateIP:    "",
				CreationTime: time.Now(),
				Tags:         []configs.Tag{},
			},
			validate: func(n *NATGatewayOutput) error {
				if n.State != "pending" {
					return fmt.Errorf("expected state to be pending, got: %s", n.State)
				}
				// PublicIP may be empty when pending
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("NAT Gateway ID: %s\n", test.natOutput.ID)
			fmt.Printf("NAT Gateway ARN: %s\n", test.natOutput.ARN)
			fmt.Printf("NAT Gateway Name: %s\n", test.natOutput.Name)
			fmt.Printf("Subnet ID: %s\n", test.natOutput.SubnetID)
			fmt.Printf("Allocation ID: %s\n", test.natOutput.AllocationID)
			fmt.Printf("State: %s\n", test.natOutput.State)
			fmt.Printf("Public IP: %s\n", test.natOutput.PublicIP)
			fmt.Printf("Private IP: %s\n", test.natOutput.PrivateIP)

			if err := test.validate(test.natOutput); err != nil {
				t.Errorf("Validation failed: %v", err)
				fmt.Printf("❌ FAILED: %v\n", err)
			} else {
				fmt.Printf("✅ PASSED: NAT Gateway output validation succeeded\n")
			}
		})
	}
}

// TestRouteTableOutput_RealisticData tests Route Table output with realistic AWS data
func TestRouteTableOutput_RealisticData(t *testing.T) {
	igwID := "igw-0a1b2c3d4e5f6g7h8"
	natGatewayID := "nat-0a1b2c3d4e5f6g7h8"

	tests := []struct {
		name        string
		description string
		rtOutput    *RouteTableOutput
		validate    func(*RouteTableOutput) error
	}{
		{
			name:        "realistic-route-table-output",
			description: "Route Table output with realistic AWS identifiers and associations",
			rtOutput: &RouteTableOutput{
				ID:    "rtb-0a1b2c3d4e5f6g7h8",
				ARN:   "arn:aws:ec2:us-east-1:123456789012:route-table/rtb-0a1b2c3d4e5f6g7h8",
				Name:  "public-route-table",
				VPCID: "vpc-0a1b2c3d4e5f6g7h8",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						GatewayID:            &igwID,
					},
					{
						DestinationCIDRBlock: "10.0.0.0/16",
						GatewayID:            nil,
					},
				},
				Associations: []RouteTableAssociation{
					{
						ID:       "rtbassoc-0a1b2c3d4e5f6g7h8",
						SubnetID: "subnet-0a1b2c3d4e5f6g7h8",
						Main:     false,
					},
				},
				CreationTime: time.Date(2024, 1, 15, 11, 15, 0, 0, time.UTC),
				Tags: []configs.Tag{
					{Key: "Name", Value: "public-route-table"},
				},
			},
			validate: func(r *RouteTableOutput) error {
				if r.ID == "" {
					return fmt.Errorf("Route Table ID is required")
				}
				if r.ARN == "" {
					return fmt.Errorf("Route Table ARN is required")
				}
				if !isValidRouteTableID(r.ID) {
					return fmt.Errorf("invalid Route Table ID format: %s", r.ID)
				}
				if !isValidARN(r.ARN, "route-table") {
					return fmt.Errorf("invalid Route Table ARN format: %s", r.ARN)
				}
				return nil
			},
		},
		{
			name:        "private-route-table-output",
			description: "Private route table with NAT Gateway route",
			rtOutput: &RouteTableOutput{
				ID:    "rtb-private-123",
				ARN:   "arn:aws:ec2:us-east-1:123456789012:route-table/rtb-private-123",
				Name:  "private-route-table",
				VPCID: "vpc-123",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						NatGatewayID:         &natGatewayID,
					},
				},
				Associations: []RouteTableAssociation{
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
			validate: func(r *RouteTableOutput) error {
				if len(r.Associations) != 2 {
					return fmt.Errorf("expected 2 associations, got %d", len(r.Associations))
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("Route Table ID: %s\n", test.rtOutput.ID)
			fmt.Printf("Route Table ARN: %s\n", test.rtOutput.ARN)
			fmt.Printf("Route Table Name: %s\n", test.rtOutput.Name)
			fmt.Printf("VPC ID: %s\n", test.rtOutput.VPCID)
			fmt.Printf("Number of Routes: %d\n", len(test.rtOutput.Routes))
			fmt.Printf("Number of Associations: %d\n", len(test.rtOutput.Associations))

			if err := test.validate(test.rtOutput); err != nil {
				t.Errorf("Validation failed: %v", err)
				fmt.Printf("❌ FAILED: %v\n", err)
			} else {
				fmt.Printf("✅ PASSED: Route Table output validation succeeded\n")
			}
		})
	}
}

// TestSecurityGroupOutput_RealisticData tests Security Group output with realistic AWS data
func TestSecurityGroupOutput_RealisticData(t *testing.T) {
	tests := []struct {
		name        string
		description string
		sgOutput    *SecurityGroupOutput
		validate    func(*SecurityGroupOutput) error
	}{
		{
			name:        "realistic-security-group-output",
			description: "Security Group output with realistic AWS identifiers and rules",
			sgOutput: &SecurityGroupOutput{
				ID:          "sg-0a1b2c3d4e5f6g7h8",
				ARN:         "arn:aws:ec2:us-east-1:123456789012:security-group/sg-0a1b2c3d4e5f6g7h8",
				Name:        "web-server-sg",
				Description: "Security group for web servers",
				VPCID:       "vpc-0a1b2c3d4e5f6g7h8",
				Rules: []networking.SecurityGroupRule{
					{
						Type:        "ingress",
						Protocol:    "tcp",
						FromPort:    intPtr(80),
						ToPort:      intPtr(80),
						CIDRBlocks:  []string{"0.0.0.0/0"},
						Description: "Allow HTTP from anywhere",
					},
					{
						Type:        "ingress",
						Protocol:    "tcp",
						FromPort:    intPtr(443),
						ToPort:      intPtr(443),
						CIDRBlocks:  []string{"0.0.0.0/0"},
						Description: "Allow HTTPS from anywhere",
					},
				},
				CreationTime: time.Date(2024, 1, 15, 11, 30, 0, 0, time.UTC),
				Tags: []configs.Tag{
					{Key: "Name", Value: "web-server-sg"},
				},
			},
			validate: func(s *SecurityGroupOutput) error {
				if s.ID == "" {
					return fmt.Errorf("Security Group ID is required")
				}
				if s.ARN == "" {
					return fmt.Errorf("Security Group ARN is required")
				}
				if !isValidSecurityGroupID(s.ID) {
					return fmt.Errorf("invalid Security Group ID format: %s", s.ID)
				}
				if !isValidARN(s.ARN, "security-group") {
					return fmt.Errorf("invalid Security Group ARN format: %s", s.ARN)
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("Security Group ID: %s\n", test.sgOutput.ID)
			fmt.Printf("Security Group ARN: %s\n", test.sgOutput.ARN)
			fmt.Printf("Security Group Name: %s\n", test.sgOutput.Name)
			fmt.Printf("Description: %s\n", test.sgOutput.Description)
			fmt.Printf("VPC ID: %s\n", test.sgOutput.VPCID)
			fmt.Printf("Number of Rules: %d\n", len(test.sgOutput.Rules))

			if err := test.validate(test.sgOutput); err != nil {
				t.Errorf("Validation failed: %v", err)
				fmt.Printf("❌ FAILED: %v\n", err)
			} else {
				fmt.Printf("✅ PASSED: Security Group output validation succeeded\n")
			}
		})
	}
}

// Helper functions for validation

func isValidVPCID(id string) bool {
	return len(id) > 4 && id[:4] == "vpc-"
}

func isValidSubnetID(id string) bool {
	return len(id) > 7 && id[:7] == "subnet-"
}

func isValidIGWID(id string) bool {
	return len(id) > 4 && id[:4] == "igw-"
}

func isValidNATGatewayID(id string) bool {
	return len(id) > 4 && id[:4] == "nat-"
}

func isValidRouteTableID(id string) bool {
	return len(id) > 4 && id[:4] == "rtb-"
}

func isValidSecurityGroupID(id string) bool {
	return len(id) > 3 && id[:3] == "sg-"
}

func isValidARN(arn, resourceType string) bool {
	if len(arn) < 20 {
		return false
	}
	// Basic ARN format: arn:aws:service:region:account-id:resource-type/resource-id
	return arn[:4] == "arn:" && contains(arn, resourceType)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func intPtr(i int) *int {
	return &i
}
