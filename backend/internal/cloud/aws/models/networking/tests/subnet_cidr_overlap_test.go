package tests

import (
	"fmt"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
)

type SubnetCIDROverlapTest struct {
	name          string
	subnets       []*networking.Subnet
	expectedError error
	description   string
}

func TestSubnetCIDROverlap(t *testing.T) {
	tests := []SubnetCIDROverlapTest{
		{
			name: "valid-four-non-overlapping-subnets",
			subnets: []*networking.Subnet{
				{
					Name:             "public-subnet",
					VPCID:            "vpc-123",
					CIDR:             "10.0.1.0/24",
					AvailabilityZone: "us-east-1a",
					Tags:             []configs.Tag{{Key: "Name", Value: "public-subnet"}},
				},
				{
					Name:             "private-subnet-1",
					VPCID:            "vpc-123",
					CIDR:             "10.0.2.0/24",
					AvailabilityZone: "us-east-1a",
					Tags:             []configs.Tag{{Key: "Name", Value: "private-subnet-1"}},
				},
				{
					Name:             "private-subnet-2",
					VPCID:            "vpc-123",
					CIDR:             "10.0.3.0/24",
					AvailabilityZone: "us-east-1b",
					Tags:             []configs.Tag{{Key: "Name", Value: "private-subnet-2"}},
				},
				{
					Name:             "private-subnet-3",
					VPCID:            "vpc-123",
					CIDR:             "10.0.4.0/24",
					AvailabilityZone: "us-east-1a",
					Tags:             []configs.Tag{{Key: "Name", Value: "private-subnet-3"}},
				},
			},
			expectedError: nil,
			description:   "Four non-overlapping subnets: 10.0.1.0/24, 10.0.2.0/24, 10.0.3.0/24, 10.0.4.0/24",
		},
		{
			name: "invalid-overlapping-subnets-same-cidr",
			subnets: []*networking.Subnet{
				{
					Name:             "subnet-1",
					VPCID:            "vpc-123",
					CIDR:             "10.0.1.0/24",
					AvailabilityZone: "us-east-1a",
				},
				{
					Name:             "subnet-2",
					VPCID:            "vpc-123",
					CIDR:             "10.0.1.0/24",
					AvailabilityZone: "us-east-1b",
				},
			},
			expectedError: fmt.Errorf("subnets have overlapping CIDR blocks"),
			description:   "Two subnets with identical CIDR blocks should overlap",
		},
		{
			name: "invalid-overlapping-subnets-partial-overlap",
			subnets: []*networking.Subnet{
				{
					Name:             "subnet-1",
					VPCID:            "vpc-123",
					CIDR:             "10.0.1.0/24",
					AvailabilityZone: "us-east-1a",
				},
				{
					Name:             "subnet-2",
					VPCID:            "vpc-123",
					CIDR:             "10.0.1.128/25",
					AvailabilityZone: "us-east-1b",
				},
			},
			expectedError: fmt.Errorf("subnets have overlapping CIDR blocks"),
			description:   "Subnet 10.0.1.128/25 overlaps with 10.0.1.0/24",
		},
		{
			name: "invalid-overlapping-subnets-one-contains-other",
			subnets: []*networking.Subnet{
				{
					Name:             "subnet-1",
					VPCID:            "vpc-123",
					CIDR:             "10.0.1.0/24",
					AvailabilityZone: "us-east-1a",
				},
				{
					Name:             "subnet-2",
					VPCID:            "vpc-123",
					CIDR:             "10.0.1.0/26",
					AvailabilityZone: "us-east-1b",
				},
			},
			expectedError: fmt.Errorf("subnets have overlapping CIDR blocks"),
			description:   "Subnet 10.0.1.0/26 is contained within 10.0.1.0/24",
		},
		{
			name: "valid-adjacent-non-overlapping-subnets",
			subnets: []*networking.Subnet{
				{
					Name:             "subnet-1",
					VPCID:            "vpc-123",
					CIDR:             "10.0.1.0/24",
					AvailabilityZone: "us-east-1a",
				},
				{
					Name:             "subnet-2",
					VPCID:            "vpc-123",
					CIDR:             "10.0.2.0/24",
					AvailabilityZone: "us-east-1b",
				},
			},
			expectedError: nil,
			description:   "Adjacent CIDR blocks (10.0.1.0/24 and 10.0.2.0/24) do not overlap",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("Number of subnets: %d\n", len(test.subnets))

			// Print subnet details
			for i, subnet := range test.subnets {
				fmt.Printf("  Subnet %d: %s - CIDR: %s, AZ: %s\n", i+1, subnet.Name, subnet.CIDR, subnet.AvailabilityZone)
			}

			// Validate each subnet individually first
			fmt.Printf("\nValidating individual subnets...\n")
			for i, subnet := range test.subnets {
				err := subnet.Validate()
				if err != nil {
					t.Errorf("Subnet %d (%s) validation failed: %v", i+1, subnet.Name, err)
					fmt.Printf("  ❌ FAILED: Subnet %d validation error: %v\n", i+1, err)
				} else {
					fmt.Printf("  ✅ PASSED: Subnet %d validation succeeded\n", i+1)
				}
			}

			// Check for CIDR overlaps
			fmt.Printf("\nChecking for CIDR overlaps...\n")
			hasOverlap := false
			var overlapError error
			for i := 0; i < len(test.subnets); i++ {
				for j := i + 1; j < len(test.subnets); j++ {
					overlaps, err := networking.CIDROverlaps(test.subnets[i].CIDR, test.subnets[j].CIDR)
					if err != nil {
						t.Errorf("Error checking CIDR overlap: %v", err)
						fmt.Printf("  ❌ ERROR: Failed to check overlap between %s and %s: %v\n", test.subnets[i].CIDR, test.subnets[j].CIDR, err)
						continue
					}
					if overlaps {
						hasOverlap = true
						overlapError = fmt.Errorf("subnets have overlapping CIDR blocks: %s (%s) overlaps with %s (%s)", test.subnets[i].Name, test.subnets[i].CIDR, test.subnets[j].Name, test.subnets[j].CIDR)
						fmt.Printf("  ❌ OVERLAP DETECTED: %s (%s) overlaps with %s (%s)\n", test.subnets[i].Name, test.subnets[i].CIDR, test.subnets[j].Name, test.subnets[j].CIDR)
					} else {
						fmt.Printf("  ✅ NO OVERLAP: %s (%s) and %s (%s)\n", test.subnets[i].Name, test.subnets[i].CIDR, test.subnets[j].Name, test.subnets[j].CIDR)
					}
				}
			}

			// Check expected result
			if test.expectedError != nil {
				if !hasOverlap {
					t.Errorf("Expected overlap error, but no overlap was detected")
					fmt.Printf("❌ FAILED: Expected overlap but none detected\n")
				} else {
					fmt.Printf("✅ PASSED: Overlap detected as expected\n")
				}
			} else {
				if hasOverlap {
					t.Errorf("Expected no overlap, but overlap was detected: %v", overlapError)
					fmt.Printf("❌ FAILED: Unexpected overlap detected: %v\n", overlapError)
				} else {
					fmt.Printf("✅ PASSED: No overlaps detected as expected\n")
				}
			}

			fmt.Printf("=== Test completed: %s ===\n\n", test.name)
		})
	}
}
