package tests

import (
	"fmt"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
)

type SubnetCIDRInVPCTest struct {
	name          string
	vpcCIDR       string
	subnet        *networking.Subnet
	expectedError error
	description   string
}

func TestSubnetCIDRInVPC(t *testing.T) {
	vpcCIDR := "10.0.0.0/16"

	tests := []SubnetCIDRInVPCTest{
		{
			name:    "valid-subnet-cidr-within-vpc",
			vpcCIDR: vpcCIDR,
			subnet: &networking.Subnet{
				Name:             "public-subnet",
				VPCID:            "vpc-123",
				CIDR:             "10.0.1.0/24",
				AvailabilityZone: "us-east-1a",
				Tags:             []configs.Tag{{Key: "Name", Value: "public-subnet"}},
			},
			expectedError: nil,
			description:   "Subnet 10.0.1.0/24 is within VPC 10.0.0.0/16",
		},
		{
			name:    "valid-subnet-cidr-within-vpc-different-range",
			vpcCIDR: vpcCIDR,
			subnet: &networking.Subnet{
				Name:             "private-subnet-1",
				VPCID:            "vpc-123",
				CIDR:             "10.0.2.0/24",
				AvailabilityZone: "us-east-1a",
				Tags:             []configs.Tag{{Key: "Name", Value: "private-subnet-1"}},
			},
			expectedError: nil,
			description:   "Subnet 10.0.2.0/24 is within VPC 10.0.0.0/16",
		},
		{
			name:    "valid-subnet-cidr-within-vpc-another-range",
			vpcCIDR: vpcCIDR,
			subnet: &networking.Subnet{
				Name:             "private-subnet-2",
				VPCID:            "vpc-123",
				CIDR:             "10.0.3.0/24",
				AvailabilityZone: "us-east-1b",
				Tags:             []configs.Tag{{Key: "Name", Value: "private-subnet-2"}},
			},
			expectedError: nil,
			description:   "Subnet 10.0.3.0/24 is within VPC 10.0.0.0/16",
		},
		{
			name:    "valid-subnet-cidr-within-vpc-fourth-subnet",
			vpcCIDR: vpcCIDR,
			subnet: &networking.Subnet{
				Name:             "private-subnet-3",
				VPCID:            "vpc-123",
				CIDR:             "10.0.4.0/24",
				AvailabilityZone: "us-east-1a",
				Tags:             []configs.Tag{{Key: "Name", Value: "private-subnet-3"}},
			},
			expectedError: nil,
			description:   "Subnet 10.0.4.0/24 is within VPC 10.0.0.0/16",
		},
		{
			name:    "invalid-subnet-cidr-outside-vpc",
			vpcCIDR: vpcCIDR,
			subnet: &networking.Subnet{
				Name:             "invalid-subnet",
				VPCID:            "vpc-123",
				CIDR:             "172.16.0.0/24",
				AvailabilityZone: "us-east-1a",
			},
			expectedError: fmt.Errorf("subnet cidr must be within vpc cidr"),
			description:   "Subnet 172.16.0.0/24 is outside VPC 10.0.0.0/16",
		},
		{
			name:    "invalid-subnet-cidr-different-network",
			vpcCIDR: vpcCIDR,
			subnet: &networking.Subnet{
				Name:             "invalid-subnet",
				VPCID:            "vpc-123",
				CIDR:             "192.168.1.0/24",
				AvailabilityZone: "us-east-1a",
			},
			expectedError: fmt.Errorf("subnet cidr must be within vpc cidr"),
			description:   "Subnet 192.168.1.0/24 is outside VPC 10.0.0.0/16",
		},
		{
			name:    "invalid-subnet-cidr-less-specific-mask",
			vpcCIDR: vpcCIDR,
			subnet: &networking.Subnet{
				Name:             "invalid-subnet",
				VPCID:            "vpc-123",
				CIDR:             "10.0.0.0/8",
				AvailabilityZone: "us-east-1a",
			},
			expectedError: fmt.Errorf("subnet cidr must be within vpc cidr"),
			description:   "Subnet 10.0.0.0/8 has less specific mask than VPC 10.0.0.0/16",
		},
		{
			name:    "invalid-subnet-cidr-same-mask",
			vpcCIDR: vpcCIDR,
			subnet: &networking.Subnet{
				Name:             "invalid-subnet",
				VPCID:            "vpc-123",
				CIDR:             "10.0.0.0/16",
				AvailabilityZone: "us-east-1a",
			},
			expectedError: fmt.Errorf("subnet cidr must be within vpc cidr"),
			description:   "Subnet 10.0.0.0/16 has same mask as VPC 10.0.0.0/16 (must be more specific)",
		},
		{
			name:    "invalid-subnet-cidr-partially-outside-vpc",
			vpcCIDR: vpcCIDR,
			subnet: &networking.Subnet{
				Name:             "invalid-subnet",
				VPCID:            "vpc-123",
				CIDR:             "10.0.255.0/24",
				AvailabilityZone: "us-east-1a",
			},
			expectedError: nil, // This should actually be valid as 10.0.255.0/24 is within 10.0.0.0/16
			description:   "Subnet 10.0.255.0/24 is within VPC 10.0.0.0/16 (edge case)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("VPC CIDR: %s\n", test.vpcCIDR)
			fmt.Printf("Subnet: %s (CIDR: %s, AZ: %s)\n", test.subnet.Name, test.subnet.CIDR, test.subnet.AvailabilityZone)

			// Validate subnet first
			fmt.Printf("\nValidating subnet...\n")
			err := test.subnet.Validate()
			if err != nil {
				if test.expectedError == nil {
					t.Errorf("Subnet validation failed unexpectedly: %v", err)
					fmt.Printf("❌ FAILED: Subnet validation error: %v\n", err)
				} else {
					fmt.Printf("⚠️  Subnet validation error (may be expected): %v\n", err)
				}
			} else {
				fmt.Printf("✅ PASSED: Subnet validation succeeded\n")
			}

			// Validate CIDR in VPC
			fmt.Printf("\nValidating subnet CIDR within VPC CIDR...\n")
			cidrErr := test.subnet.ValidateCIDRInVPC(test.vpcCIDR)

			if test.expectedError != nil {
				if cidrErr == nil {
					t.Errorf("Expected error: %v, but got nil", test.expectedError)
					fmt.Printf("❌ FAILED: Expected error '%v', but got nil\n", test.expectedError)
				} else if cidrErr.Error() != test.expectedError.Error() {
					// Check if error message contains expected text
					if cidrErr.Error() == "subnet cidr must be within vpc cidr" && test.expectedError.Error() == "subnet cidr must be within vpc cidr" {
						fmt.Printf("✅ PASSED: Got expected error: %v\n", cidrErr)
					} else {
						t.Errorf("Expected error: %v, but got: %v", test.expectedError, cidrErr)
						fmt.Printf("❌ FAILED: Expected error '%v', but got '%v'\n", test.expectedError, cidrErr)
					}
				} else {
					fmt.Printf("✅ PASSED: Got expected error: %v\n", cidrErr)
				}
			} else {
				if cidrErr != nil {
					t.Errorf("Expected no error, but got: %v", cidrErr)
					fmt.Printf("❌ FAILED: Unexpected error: %v\n", cidrErr)
				} else {
					fmt.Printf("✅ PASSED: Subnet CIDR is within VPC CIDR\n")
				}
			}

			// Additional check using CIDRContains utility
			fmt.Printf("\nVerifying with CIDRContains utility...\n")
			contains, err := networking.CIDRContains(test.vpcCIDR, test.subnet.CIDR)
			if err != nil {
				fmt.Printf("⚠️  Error checking CIDR containment: %v\n", err)
			} else if contains {
				fmt.Printf("✅ VERIFIED: VPC CIDR contains subnet CIDR\n")
			} else {
				fmt.Printf("⚠️  VERIFIED: VPC CIDR does not contain subnet CIDR\n")
			}

			fmt.Printf("=== Test completed: %s ===\n\n", test.name)
		})
	}
}
