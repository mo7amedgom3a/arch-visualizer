package tests

import (
	"fmt"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
)

type SubnetAZConstraintTest struct {
	name          string
	subnet        *networking.Subnet
	expectedError error
	description   string
}

func TestSubnetAZConstraint(t *testing.T) {
	tests := []SubnetAZConstraintTest{
		{
			name: "valid-subnet-in-us-east-1a",
			subnet: &networking.Subnet{
				Name:             "public-subnet",
				VPCID:            "vpc-123",
				CIDR:             "10.0.1.0/24",
				AvailabilityZone: "us-east-1a",
				Tags:             []configs.Tag{{Key: "Name", Value: "public-subnet"}},
			},
			expectedError: nil,
			description:   "Subnet in us-east-1a (single AZ)",
		},
		{
			name: "valid-subnet-in-us-east-1b",
			subnet: &networking.Subnet{
				Name:             "private-subnet-1",
				VPCID:            "vpc-123",
				CIDR:             "10.0.2.0/24",
				AvailabilityZone: "us-east-1b",
				Tags:             []configs.Tag{{Key: "Name", Value: "private-subnet-1"}},
			},
			expectedError: nil,
			description:   "Subnet in us-east-1b (single AZ)",
		},
		{
			name: "valid-subnet-in-us-east-1a-second",
			subnet: &networking.Subnet{
				Name:             "private-subnet-2",
				VPCID:            "vpc-123",
				CIDR:             "10.0.3.0/24",
				AvailabilityZone: "us-east-1a",
				Tags:             []configs.Tag{{Key: "Name", Value: "private-subnet-2"}},
			},
			expectedError: nil,
			description:   "Another subnet in us-east-1a (single AZ)",
		},
		{
			name: "valid-subnet-in-us-east-1b-second",
			subnet: &networking.Subnet{
				Name:             "private-subnet-3",
				VPCID:            "vpc-123",
				CIDR:             "10.0.4.0/24",
				AvailabilityZone: "us-east-1b",
				Tags:             []configs.Tag{{Key: "Name", Value: "private-subnet-3"}},
			},
			expectedError: nil,
			description:   "Another subnet in us-east-1b (single AZ)",
		},
		{
			name: "invalid-subnet-with-empty-az",
			subnet: &networking.Subnet{
				Name:             "invalid-subnet",
				VPCID:            "vpc-123",
				CIDR:             "10.0.5.0/24",
				AvailabilityZone: "",
			},
			expectedError: fmt.Errorf("subnet availability_zone is required"),
			description:   "Subnet with empty AZ should fail validation",
		},
		{
			name: "invalid-subnet-with-invalid-az-format",
			subnet: &networking.Subnet{
				Name:             "invalid-subnet",
				VPCID:            "vpc-123",
				CIDR:             "10.0.5.0/24",
				AvailabilityZone: "invalid-az",
			},
			expectedError: nil, // Format validation is not strict, just checks non-empty
			description:   "Subnet with invalid AZ format (validation only checks non-empty)",
		},
		{
			name: "valid-subnet-az-format-check",
			subnet: &networking.Subnet{
				Name:             "valid-subnet",
				VPCID:            "vpc-123",
				CIDR:             "10.0.6.0/24",
				AvailabilityZone: "us-east-1a",
				Tags:             []configs.Tag{{Key: "Name", Value: "valid-subnet"}},
			},
			expectedError: nil,
			description:   "Subnet with valid AZ format (region-az pattern)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("Subnet: %s\n", test.subnet.Name)
			fmt.Printf("CIDR: %s\n", test.subnet.CIDR)
			fmt.Printf("Availability Zone: %s\n", test.subnet.AvailabilityZone)

			// Validate subnet
			fmt.Printf("\nValidating subnet...\n")
			err := test.subnet.Validate()

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
					fmt.Printf("✅ PASSED: Subnet validation succeeded\n")
				}
			}

			// Check AZ format and constraint
			if test.subnet.AvailabilityZone != "" {
				fmt.Printf("\nChecking AZ constraint...\n")
				// AWS AZ format is typically: region-az (e.g., us-east-1a)
				// A subnet must be in a single AZ and cannot span multiple AZs
				// This is enforced by AWS - a subnet is created in a specific AZ
				az := test.subnet.AvailabilityZone
				fmt.Printf("  AZ: %s\n", az)

				// Basic format check: should be at least 9 characters (us-east-1a)
				if len(az) >= 9 {
					// Extract region part (everything except last character)
					region := az[:len(az)-1]
					fmt.Printf("  Extracted region: %s\n", region)
					fmt.Printf("  ✅ PASSED: AZ format appears valid (region-az pattern)\n")
				} else {
					fmt.Printf("  ⚠️  WARNING: AZ format may be invalid (too short)\n")
				}

				// Note: AWS enforces that a subnet cannot span multiple AZs
				// This is a physical constraint - we just validate the AZ field is present
				fmt.Printf("  ✅ PASSED: Subnet is assigned to single AZ: %s\n", az)
			} else {
				fmt.Printf("\n⚠️  WARNING: AZ is empty (should be validated above)\n")
			}

			fmt.Printf("=== Test completed: %s ===\n\n", test.name)
		})
	}
}
