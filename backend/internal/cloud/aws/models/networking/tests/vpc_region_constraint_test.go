package tests

import (
	"fmt"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
)

type VPCRegionTest struct {
	name          string
	vpc           *networking.VPC
	subnets       []*networking.Subnet
	expectedError error
	description   string
}

func TestVPCRegionConstraint(t *testing.T) {
	validVPC := &networking.VPC{
		Name:               "test-vpc",
		Region:             "us-east-1",
		CIDR:               "10.0.0.0/16",
		EnableDNSHostnames: true,
		EnableDNSSupport:   true,
		InstanceTenancy:    "default",
		Tags:               []configs.Tag{{Key: "Name", Value: "test-vpc"}},
	}

	validSubnet1 := &networking.Subnet{
		Name:             "public-subnet",
		VPCID:            "vpc-123",
		CIDR:             "10.0.1.0/24",
		AvailabilityZone: "us-east-1a",
		Tags:             []configs.Tag{{Key: "Name", Value: "public-subnet"}},
	}

	validSubnet2 := &networking.Subnet{
		Name:             "private-subnet-1",
		VPCID:            "vpc-123",
		CIDR:             "10.0.2.0/24",
		AvailabilityZone: "us-east-1b",
		Tags:             []configs.Tag{{Key: "Name", Value: "private-subnet-1"}},
	}

	tests := []VPCRegionTest{
		{
			name:          "valid-vpc-in-us-east-1-with-subnets-in-same-region",
			vpc:           validVPC,
			subnets:       []*networking.Subnet{validSubnet1, validSubnet2},
			expectedError: nil,
			description:   "VPC in us-east-1 with subnets in us-east-1a and us-east-1b (same region)",
		},
		{
			name: "vpc-with-empty-region",
			vpc: &networking.VPC{
				Name:               "test-vpc",
				Region:             "",
				CIDR:               "10.0.0.0/16",
				EnableDNSHostnames: true,
				EnableDNSSupport:   true,
				InstanceTenancy:    "default",
			},
			subnets:       []*networking.Subnet{},
			expectedError: fmt.Errorf("region is required"),
			description:   "VPC with empty region should fail validation",
		},
		{
			name: "vpc-with-invalid-region-format",
			vpc: &networking.VPC{
				Name:               "test-vpc",
				Region:             "invalid-region",
				CIDR:               "10.0.0.0/16",
				EnableDNSHostnames: true,
				EnableDNSSupport:   true,
				InstanceTenancy:    "default",
			},
			subnets:       []*networking.Subnet{},
			expectedError: nil, // Region format validation is not strict, just checks non-empty
			description:   "VPC with invalid region format (validation only checks non-empty)",
		},
		{
			name:          "vpc-in-us-east-1-with-subnets-in-us-east-1a-and-us-east-1b",
			vpc:           validVPC,
			subnets:       []*networking.Subnet{validSubnet1, validSubnet2},
			expectedError: nil,
			description:   "Valid scenario: VPC in us-east-1 with subnets in us-east-1a and us-east-1b",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("VPC Region: %s\n", test.vpc.Region)

			// Validate VPC
			fmt.Printf("Validating VPC...\n")
			err := test.vpc.Validate()
			if test.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error: %v, but got nil", test.expectedError)
					fmt.Printf("❌ FAILED: Expected error but got nil\n")
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
					fmt.Printf("✅ PASSED: VPC validation succeeded\n")
				}
			}

			// Validate subnets if provided
			if len(test.subnets) > 0 {
				fmt.Printf("\nValidating %d subnet(s)...\n", len(test.subnets))
				for i, subnet := range test.subnets {
					fmt.Printf("  Subnet %d: %s (AZ: %s)\n", i+1, subnet.Name, subnet.AvailabilityZone)
					subnetErr := subnet.Validate()
					if subnetErr != nil {
						t.Errorf("Subnet %d validation failed: %v", i+1, subnetErr)
						fmt.Printf("  ❌ FAILED: Subnet validation error: %v\n", subnetErr)
					} else {
						fmt.Printf("  ✅ PASSED: Subnet validation succeeded\n")
					}

					// Check AZ format (should be region-az format like us-east-1a)
					if subnet.AvailabilityZone != "" {
						// Extract region from AZ (us-east-1a -> us-east-1)
						azRegion := ""
						if len(subnet.AvailabilityZone) >= 9 {
							azRegion = subnet.AvailabilityZone[:len(subnet.AvailabilityZone)-1]
						}
						if azRegion != "" && test.vpc.Region != "" && azRegion != test.vpc.Region {
							t.Errorf("Subnet %d AZ %s is not in VPC region %s", i+1, subnet.AvailabilityZone, test.vpc.Region)
							fmt.Printf("  ❌ FAILED: Subnet AZ %s is not in VPC region %s\n", subnet.AvailabilityZone, test.vpc.Region)
						} else if azRegion != "" {
							fmt.Printf("  ✅ PASSED: Subnet AZ %s is in VPC region %s\n", subnet.AvailabilityZone, test.vpc.Region)
						}
					}
				}
			}

			fmt.Printf("=== Test completed: %s ===\n\n", test.name)
		})
	}
}
