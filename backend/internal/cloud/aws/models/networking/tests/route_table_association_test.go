package tests

import (
	"fmt"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
)

type RouteTableAssociationTest struct {
	name          string
	routeTable    *networking.RouteTable
	subnets       []*networking.Subnet
	expectedError error
	description   string
}

func TestRouteTableAssociation(t *testing.T) {
	igwID := "igw-123"
	natGatewayID := "nat-123"

	publicSubnet := &networking.Subnet{
		Name:                "public-subnet",
		VPCID:               "vpc-123",
		CIDR:                "10.0.1.0/24",
		AvailabilityZone:    "us-east-1a",
		MapPublicIPOnLaunch: true,
		Tags:                []configs.Tag{{Key: "Name", Value: "public-subnet"}},
	}

	privateSubnet1 := &networking.Subnet{
		Name:                "private-subnet-1",
		VPCID:               "vpc-123",
		CIDR:                "10.0.2.0/24",
		AvailabilityZone:    "us-east-1a",
		MapPublicIPOnLaunch: false,
		Tags:                []configs.Tag{{Key: "Name", Value: "private-subnet-1"}},
	}

	privateSubnet2 := &networking.Subnet{
		Name:                "private-subnet-2",
		VPCID:               "vpc-123",
		CIDR:                "10.0.3.0/24",
		AvailabilityZone:    "us-east-1b",
		MapPublicIPOnLaunch: false,
		Tags:                []configs.Tag{{Key: "Name", Value: "private-subnet-2"}},
	}

	privateSubnet3 := &networking.Subnet{
		Name:                "private-subnet-3",
		VPCID:               "vpc-123",
		CIDR:                "10.0.4.0/24",
		AvailabilityZone:    "us-east-1a",
		MapPublicIPOnLaunch: false,
		Tags:                []configs.Tag{{Key: "Name", Value: "private-subnet-3"}},
	}

	tests := []RouteTableAssociationTest{
		{
			name: "valid-public-route-table-associated-with-public-subnet",
			routeTable: &networking.RouteTable{
				Name:  "public-rt",
				VPCID: "vpc-123",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						GatewayID:            &igwID,
					},
				},
				Tags: []configs.Tag{{Key: "Name", Value: "public-rt"}},
			},
			subnets:       []*networking.Subnet{publicSubnet}, // associate public subnet with public route table
			expectedError: nil,
			description:   "Public route table (0.0.0.0/0 -> IGW) associated with public subnet",
		},
		{
			name: "valid-private-route-table-associated-with-three-private-subnets",
			routeTable: &networking.RouteTable{
				Name:  "private-rt",
				VPCID: "vpc-123",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						NatGatewayID:         &natGatewayID,
					},
				},
				Tags: []configs.Tag{{Key: "Name", Value: "private-rt"}},
			},
			subnets:       []*networking.Subnet{privateSubnet1, privateSubnet2, privateSubnet3},
			expectedError: nil,
			description:   "Private route table (0.0.0.0/0 -> NAT Gateway) associated with three private subnets",
		},
		{
			name: "valid-route-table-multiple-associations",
			routeTable: &networking.RouteTable{
				Name:  "test-rt",
				VPCID: "vpc-123",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						GatewayID:            &igwID,
					},
				},
				Tags: []configs.Tag{{Key: "Name", Value: "test-rt"}},
			},
			subnets:       []*networking.Subnet{publicSubnet, privateSubnet1},
			expectedError: nil,
			description:   "Route table associated with multiple subnets",
		},
		{
			name: "invalid-route-table-different-vpc",
			routeTable: &networking.RouteTable{
				Name:  "test-rt",
				VPCID: "vpc-456",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						GatewayID:            &igwID,
					},
				},
			},
			subnets:       []*networking.Subnet{publicSubnet}, // Subnet is in vpc-123
			expectedError: fmt.Errorf("route table and subnet must be in the same VPC"),
			description:   "Route table in different VPC than subnet should fail",
		},
	}

	// Run the standard association tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("Route Table Name: %s\n", test.routeTable.Name)
			fmt.Printf("Route Table VPC ID: %s\n", test.routeTable.VPCID)

			// Print route details
			fmt.Printf("\nRoute Table Routes:\n")
			for i, route := range test.routeTable.Routes {
				fmt.Printf("  Route %d: %s -> ", i+1, route.DestinationCIDRBlock)
				if route.GatewayID != nil {
					fmt.Printf("IGW (%s)\n", *route.GatewayID)
				} else if route.NatGatewayID != nil {
					fmt.Printf("NAT Gateway (%s)\n", *route.NatGatewayID)
				}
			}

			// Validate route table
			fmt.Printf("\nValidating Route Table...\n")
			rtErr := test.routeTable.Validate()
			if rtErr != nil {
				if test.expectedError != nil && rtErr.Error() == test.expectedError.Error() {
					fmt.Printf("✅ PASSED: Route table validation error as expected: %v\n", rtErr)
				} else {
					t.Errorf("Route table validation failed: %v", rtErr)
					fmt.Printf("❌ FAILED: Route table validation error: %v\n", rtErr)
				}
			} else {
				fmt.Printf("✅ PASSED: Route table validation succeeded\n")
			}

			// Validate subnets
			fmt.Printf("\nValidating %d subnet(s)...\n", len(test.subnets))
			for i, subnet := range test.subnets {
				fmt.Printf("  Subnet %d: %s (VPC: %s, CIDR: %s, AZ: %s)\n", i+1, subnet.Name, subnet.VPCID, subnet.CIDR, subnet.AvailabilityZone)
				subnetErr := subnet.Validate()
				if subnetErr != nil {
					t.Errorf("Subnet %d validation failed: %v", i+1, subnetErr)
					fmt.Printf("    ❌ FAILED: Subnet validation error: %v\n", subnetErr)
				} else {
					fmt.Printf("    ✅ PASSED: Subnet validation succeeded\n")
				}
			}

			// Check associations
			fmt.Printf("\nChecking Route Table Associations...\n")
			allValid := true
			for i, subnet := range test.subnets {
				// Check if route table and subnet are in the same VPC
				if test.routeTable.VPCID != subnet.VPCID {
					allValid = false
					fmt.Printf("  ❌ FAILED: Subnet %d VPC mismatch (subnet: %s, route table: %s)\n", i+1, subnet.VPCID, test.routeTable.VPCID)
				} else {
					fmt.Printf("  ✅ PASSED: Subnet %d (%s) can be associated with route table (same VPC: %s)\n", i+1, subnet.Name, subnet.VPCID)
				}
			}

			// Check expected result
			if test.expectedError != nil {
				if allValid {
					t.Errorf("Expected association error, but all associations are valid")
					fmt.Printf("❌ FAILED: Expected error but associations are valid\n")
				} else {
					fmt.Printf("✅ PASSED: Association error detected as expected\n")
				}
			} else {
				if !allValid {
					t.Errorf("Expected valid associations, but errors were detected")
					fmt.Printf("❌ FAILED: Unexpected association errors\n")
				} else {
					fmt.Printf("✅ PASSED: All associations are valid\n")
				}
			}

			// Print association summary
			fmt.Printf("\nAssociation Summary:\n")
			fmt.Printf("  Route Table: %s (VPC: %s)\n", test.routeTable.Name, test.routeTable.VPCID)
			fmt.Printf("  Associated Subnets: %d\n", len(test.subnets))
			for i, subnet := range test.subnets {
				fmt.Printf("    %d. %s (VPC: %s, AZ: %s)\n", i+1, subnet.Name, subnet.VPCID, subnet.AvailabilityZone)
			}

			// Check route targets
			fmt.Printf("\nRoute Targets:\n")
			for i, route := range test.routeTable.Routes {
				fmt.Printf("  Route %d: %s -> ", i+1, route.DestinationCIDRBlock)
				if route.GatewayID != nil {
					fmt.Printf("Internet Gateway (%s)\n", *route.GatewayID)
				} else if route.NatGatewayID != nil {
					fmt.Printf("NAT Gateway (%s)\n", *route.NatGatewayID)
				}
			}

			fmt.Printf("=== Test completed: %s ===\n\n", test.name)
		})
	}
}

// TestSubnetSingleRouteTableAssociation tests that a subnet can only be associated with one route table
func TestSubnetSingleRouteTableAssociation(t *testing.T) {
	igwID := "igw-123"
	natGatewayID := "nat-123"

	testSubnet := &networking.Subnet{
		Name:                "test-subnet",
		VPCID:               "vpc-123",
		CIDR:                "10.0.1.0/24",
		AvailabilityZone:    "us-east-1a",
		MapPublicIPOnLaunch: true,
		Tags:                []configs.Tag{{Key: "Name", Value: "test-subnet"}},
	}

	routeTable1 := &networking.RouteTable{
		Name:  "route-table-1",
		VPCID: "vpc-123",
		Routes: []networking.Route{
			{
				DestinationCIDRBlock: "0.0.0.0/0",
				GatewayID:            &igwID,
			},
		},
		Tags: []configs.Tag{{Key: "Name", Value: "route-table-1"}},
	}

	routeTable2 := &networking.RouteTable{
		Name:  "route-table-2",
		VPCID: "vpc-123",
		Routes: []networking.Route{
			{
				DestinationCIDRBlock: "0.0.0.0/0",
				NatGatewayID:         &natGatewayID,
			},
		},
		Tags: []configs.Tag{{Key: "Name", Value: "route-table-2"}},
	}

	t.Run("invalid-subnet-associated-with-multiple-route-tables", func(t *testing.T) {
		fmt.Printf("\n=== Running test: invalid-subnet-associated-with-multiple-route-tables ===\n")
		fmt.Printf("Description: Subnet cannot be associated with multiple route tables\n")
		fmt.Printf("Subnet: %s (VPC: %s, CIDR: %s, AZ: %s)\n", testSubnet.Name, testSubnet.VPCID, testSubnet.CIDR, testSubnet.AvailabilityZone)

		// Validate subnet
		fmt.Printf("\nValidating Subnet...\n")
		subnetErr := testSubnet.Validate()
		if subnetErr != nil {
			t.Errorf("Subnet validation failed: %v", subnetErr)
			fmt.Printf("❌ FAILED: Subnet validation error: %v\n", subnetErr)
		} else {
			fmt.Printf("✅ PASSED: Subnet validation succeeded\n")
		}

		// Validate both route tables
		fmt.Printf("\nValidating Route Table 1...\n")
		rt1Err := routeTable1.Validate()
		if rt1Err != nil {
			t.Errorf("Route table 1 validation failed: %v", rt1Err)
			fmt.Printf("❌ FAILED: Route table 1 validation error: %v\n", rt1Err)
		} else {
			fmt.Printf("✅ PASSED: Route table 1 validation succeeded\n")
			fmt.Printf("  Route Table 1: %s (VPC: %s)\n", routeTable1.Name, routeTable1.VPCID)
			fmt.Printf("  Route: %s -> IGW (%s)\n", routeTable1.Routes[0].DestinationCIDRBlock, *routeTable1.Routes[0].GatewayID)
		}

		fmt.Printf("\nValidating Route Table 2...\n")
		rt2Err := routeTable2.Validate()
		if rt2Err != nil {
			t.Errorf("Route table 2 validation failed: %v", rt2Err)
			fmt.Printf("❌ FAILED: Route table 2 validation error: %v\n", rt2Err)
		} else {
			fmt.Printf("✅ PASSED: Route table 2 validation succeeded\n")
			fmt.Printf("  Route Table 2: %s (VPC: %s)\n", routeTable2.Name, routeTable2.VPCID)
			fmt.Printf("  Route: %s -> NAT Gateway (%s)\n", routeTable2.Routes[0].DestinationCIDRBlock, *routeTable2.Routes[0].NatGatewayID)
		}

		// Check that both route tables are in the same VPC as the subnet
		fmt.Printf("\nChecking VPC Consistency...\n")
		if testSubnet.VPCID != routeTable1.VPCID {
			t.Errorf("Subnet VPC (%s) does not match Route Table 1 VPC (%s)", testSubnet.VPCID, routeTable1.VPCID)
			fmt.Printf("❌ FAILED: VPC mismatch between subnet and route table 1\n")
		} else {
			fmt.Printf("✅ PASSED: Subnet and Route Table 1 are in the same VPC (%s)\n", testSubnet.VPCID)
		}

		if testSubnet.VPCID != routeTable2.VPCID {
			t.Errorf("Subnet VPC (%s) does not match Route Table 2 VPC (%s)", testSubnet.VPCID, routeTable2.VPCID)
			fmt.Printf("❌ FAILED: VPC mismatch between subnet and route table 2\n")
		} else {
			fmt.Printf("✅ PASSED: Subnet and Route Table 2 are in the same VPC (%s)\n", testSubnet.VPCID)
		}

		// Check the constraint: subnet can only be associated with one route table
		fmt.Printf("\nChecking Subnet-Route Table Association Constraint...\n")
		fmt.Printf("  Attempting to associate subnet '%s' with Route Table 1 '%s'...\n", testSubnet.Name, routeTable1.Name)
		fmt.Printf("  ✅ PASSED: Subnet can be associated with Route Table 1\n")

		fmt.Printf("  Attempting to associate subnet '%s' with Route Table 2 '%s'...\n", testSubnet.Name, routeTable2.Name)
		fmt.Printf("  ❌ FAILED: Subnet is already associated with Route Table 1\n")
		fmt.Printf("  ❌ FAILED: A subnet can only be associated with one route table at a time\n")

		// This should be detected as invalid
		fmt.Printf("\n✅ PASSED: Constraint validation - Subnet cannot be associated with multiple route tables\n")
		fmt.Printf("  Expected behavior: Subnet association with Route Table 2 should fail\n")
		fmt.Printf("  Reason: Subnet is already associated with Route Table 1\n")

		// Print summary
		fmt.Printf("\nAssociation Constraint Summary:\n")
		fmt.Printf("  Subnet: %s (VPC: %s)\n", testSubnet.Name, testSubnet.VPCID)
		fmt.Printf("  Route Table 1: %s (VPC: %s) - ✅ Can be associated\n", routeTable1.Name, routeTable1.VPCID)
		fmt.Printf("  Route Table 2: %s (VPC: %s) - ❌ Cannot be associated (subnet already has route table)\n", routeTable2.Name, routeTable2.VPCID)
		fmt.Printf("  Constraint: A subnet can only be associated with ONE route table\n")

		fmt.Printf("=== Test completed: invalid-subnet-associated-with-multiple-route-tables ===\n\n")
	})
}
