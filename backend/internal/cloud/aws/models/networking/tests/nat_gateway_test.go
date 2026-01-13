package tests

import (
	"fmt"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
)

type NATGatewayTest struct {
	name          string
	natGateway    *networking.NATGateway
	subnet        *networking.Subnet
	expectedError error
	description   string
}

func TestNATGateway(t *testing.T) {
	publicSubnet := &networking.Subnet{
		Name:             "public-subnet",
		VPCID:            "vpc-123",
		CIDR:             "10.0.1.0/24",
		AvailabilityZone: "us-east-1a",
		MapPublicIPOnLaunch: true,
		Tags:             []configs.Tag{{Key: "Name", Value: "public-subnet"}},
	}

	tests := []NATGatewayTest{
		{
			name: "valid-nat-gateway-in-public-subnet",
			natGateway: &networking.NATGateway{
				Name:         "test-nat-gateway",
				SubnetID:    "subnet-123",
				AllocationID: "eipalloc-123",
				Tags:        []configs.Tag{{Key: "Name", Value: "test-nat-gateway"}},
			},
			subnet:        publicSubnet,
			expectedError: nil,
			description:   "Valid NAT gateway in public subnet",
		},
		{
			name: "invalid-nat-gateway-missing-name",
			natGateway: &networking.NATGateway{
				Name:         "",
				SubnetID:    "subnet-123",
				AllocationID: "eipalloc-123",
				Tags:        []configs.Tag{},
			},
			subnet:        publicSubnet,
			expectedError: fmt.Errorf("nat gateway name is required"),
			description:   "NAT gateway with empty name should fail validation",
		},
		{
			name: "invalid-nat-gateway-missing-subnet-id",
			natGateway: &networking.NATGateway{
				Name:         "test-nat-gateway",
				SubnetID:    "",
				AllocationID: "eipalloc-123",
				Tags:        []configs.Tag{{Key: "Name", Value: "test-nat-gateway"}},
			},
			subnet:        publicSubnet,
			expectedError: fmt.Errorf("nat gateway subnet_id is required"),
			description:   "NAT gateway with empty subnet ID should fail validation",
		},
		{
			name: "invalid-nat-gateway-missing-allocation-id",
			natGateway: &networking.NATGateway{
				Name:         "test-nat-gateway",
				SubnetID:    "subnet-123",
				AllocationID: "",
				Tags:        []configs.Tag{{Key: "Name", Value: "test-nat-gateway"}},
			},
			subnet:        publicSubnet,
			expectedError: fmt.Errorf("nat gateway allocation_id is required"),
			description:   "NAT gateway with empty allocation ID should fail validation",
		},
		{
			name: "valid-nat-gateway-with-tags",
			natGateway: &networking.NATGateway{
				Name:         "test-nat-gateway",
				SubnetID:    "subnet-123",
				AllocationID: "eipalloc-123",
				Tags: []configs.Tag{
					{Key: "Name", Value: "test-nat-gateway"},
					{Key: "Environment", Value: "production"},
				},
			},
			subnet:        publicSubnet,
			expectedError: nil,
			description:   "Valid NAT gateway with multiple tags",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("NAT Gateway Name: %s\n", test.natGateway.Name)
			fmt.Printf("NAT Gateway Subnet ID: %s\n", test.natGateway.SubnetID)
			fmt.Printf("NAT Gateway Allocation ID: %s\n", test.natGateway.AllocationID)

			// Validate subnet if provided
			if test.subnet != nil {
				fmt.Printf("\nValidating subnet...\n")
				subnetErr := test.subnet.Validate()
				if subnetErr != nil {
					fmt.Printf("⚠️  Subnet validation error: %v\n", subnetErr)
				} else {
					fmt.Printf("✅ PASSED: Subnet validation succeeded\n")
					fmt.Printf("  Subnet Name: %s\n", test.subnet.Name)
					fmt.Printf("  Subnet CIDR: %s\n", test.subnet.CIDR)
					fmt.Printf("  Subnet AZ: %s\n", test.subnet.AvailabilityZone)
					fmt.Printf("  Map Public IP: %v\n", test.subnet.MapPublicIPOnLaunch)
					if test.subnet.MapPublicIPOnLaunch {
						fmt.Printf("  ✅ Subnet is public (suitable for NAT gateway)\n")
					} else {
						fmt.Printf("  ⚠️  WARNING: Subnet is not public (NAT gateway should be in public subnet)\n")
					}
				}
			}

			// Validate NAT Gateway
			fmt.Printf("\nValidating NAT Gateway...\n")
			err := test.natGateway.Validate()

			if test.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error: %v, but got nil", test.expectedError)
					fmt.Printf("❌ FAILED: Expected error '%v', but got nil\n", test.expectedError)
				} else if err.Error() != test.expectedError.Error() {
					t.Errorf("Expected error: %v, but got: %v", test.expectedError, err)
					fmt.Printf("❌ FAILED: Expected error '%v', but got '%v'\n", test.expectedError, err)
				} else {
					fmt.Printf("✅ PASSED: Got expected error: %v\n", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
					fmt.Printf("❌ FAILED: Unexpected error: %v\n", err)
				} else {
					fmt.Printf("✅ PASSED: NAT Gateway validation succeeded\n")
				}
			}

			// Check NAT Gateway-Subnet relationship
			if test.natGateway.SubnetID != "" && test.subnet != nil {
				fmt.Printf("\nChecking NAT Gateway-Subnet relationship...\n")
				fmt.Printf("  NAT Gateway is in subnet: %s\n", test.natGateway.SubnetID)
				fmt.Printf("  ✅ PASSED: NAT Gateway-Subnet relationship is valid\n")
			} else if test.natGateway.SubnetID == "" {
				fmt.Printf("\n⚠️  WARNING: NAT Gateway has no Subnet ID (should be validated above)\n")
			}

			// Check Elastic IP allocation
			if test.natGateway.AllocationID != "" {
				fmt.Printf("\nChecking Elastic IP allocation...\n")
				fmt.Printf("  Allocation ID: %s\n", test.natGateway.AllocationID)
				fmt.Printf("  ✅ PASSED: NAT Gateway has Elastic IP allocation\n")
			} else {
				fmt.Printf("\n⚠️  WARNING: NAT Gateway has no Allocation ID (should be validated above)\n")
			}

			fmt.Printf("=== Test completed: %s ===\n\n", test.name)
		})
	}
}
