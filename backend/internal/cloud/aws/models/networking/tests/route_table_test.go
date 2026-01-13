package tests

import (
	"fmt"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
)

type RouteTableTest struct {
	name          string
	routeTable    *networking.RouteTable
	expectedError error
	description   string
}

func TestRouteTable(t *testing.T) {
	igwID := "igw-123"
	natGatewayID := "nat-123"

	tests := []RouteTableTest{
		{
			name: "valid-public-route-table-with-igw",
			routeTable: &networking.RouteTable{
				Name:  "public-rt",
				VPCID: "vpc-123",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						GatewayID:           &igwID,
					},
				},
				Tags: []configs.Tag{{Key: "Name", Value: "public-rt"}},
			},
			expectedError: nil,
			description:   "Public route table with 0.0.0.0/0 -> IGW",
		},
		{
			name: "valid-private-route-table-with-nat-gateway",
			routeTable: &networking.RouteTable{
				Name:  "private-rt",
				VPCID: "vpc-123",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						NatGatewayID:        &natGatewayID,
					},
				},
				Tags: []configs.Tag{{Key: "Name", Value: "private-rt"}},
			},
			expectedError: nil,
			description:   "Private route table with 0.0.0.0/0 -> NAT Gateway",
		},
		{
			name: "invalid-route-table-missing-name",
			routeTable: &networking.RouteTable{
				Name:  "",
				VPCID: "vpc-123",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						GatewayID:           &igwID,
					},
				},
			},
			expectedError: fmt.Errorf("route table name is required"),
			description:   "Route table with empty name should fail validation",
		},
		{
			name: "invalid-route-table-missing-vpc-id",
			routeTable: &networking.RouteTable{
				Name:  "test-rt",
				VPCID: "",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						GatewayID:           &igwID,
					},
				},
			},
			expectedError: fmt.Errorf("route table vpc_id is required"),
			description:   "Route table with empty VPC ID should fail validation",
		},
		{
			name: "invalid-route-missing-destination-cidr",
			routeTable: &networking.RouteTable{
				Name:  "test-rt",
				VPCID: "vpc-123",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "",
						GatewayID:           &igwID,
					},
				},
			},
			expectedError: fmt.Errorf("route destination_cidr_block is required"),
			description:   "Route with empty destination CIDR should fail validation",
		},
		{
			name: "invalid-route-missing-target",
			routeTable: &networking.RouteTable{
				Name:  "test-rt",
				VPCID: "vpc-123",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
					},
				},
			},
			expectedError: fmt.Errorf("route must have at least one target (gateway_id, nat_gateway_id, etc.)"),
			description:   "Route with no target should fail validation",
		},
		{
			name: "invalid-route-multiple-targets",
			routeTable: &networking.RouteTable{
				Name:  "test-rt",
				VPCID: "vpc-123",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						GatewayID:           &igwID,
						NatGatewayID:        &natGatewayID,
					},
				},
			},
			expectedError: fmt.Errorf("route can only have one target"),
			description:   "Route with multiple targets should fail validation",
		},
		{
			name: "valid-route-table-with-local-route",
			routeTable: &networking.RouteTable{
				Name:  "test-rt",
				VPCID: "vpc-123",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "10.0.0.0/16",
						// Local route has no target (handled by AWS)
					},
				},
			},
			expectedError: fmt.Errorf("route must have at least one target (gateway_id, nat_gateway_id, etc.)"),
			description:   "Route table with local route (no target) - AWS handles this differently",
		},
		{
			name: "valid-route-table-with-multiple-routes",
			routeTable: &networking.RouteTable{
				Name:  "test-rt",
				VPCID: "vpc-123",
				Routes: []networking.Route{
					{
						DestinationCIDRBlock: "10.0.0.0/16",
						GatewayID:           nil, // Local route
					},
					{
						DestinationCIDRBlock: "0.0.0.0/0",
						GatewayID:           &igwID,
					},
				},
			},
			expectedError: fmt.Errorf("route must have at least one target (gateway_id, nat_gateway_id, etc.)"),
			description:   "Route table with multiple routes including local route",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Printf("\n=== Running test: %s ===\n", test.name)
			fmt.Printf("Description: %s\n", test.description)
			fmt.Printf("Route Table Name: %s\n", test.routeTable.Name)
			fmt.Printf("Route Table VPC ID: %s\n", test.routeTable.VPCID)
			fmt.Printf("Number of routes: %d\n", len(test.routeTable.Routes))

			// Print route details
			for i, route := range test.routeTable.Routes {
				fmt.Printf("\n  Route %d:\n", i+1)
				fmt.Printf("    Destination: %s\n", route.DestinationCIDRBlock)
				if route.GatewayID != nil {
					fmt.Printf("    Target: IGW (%s)\n", *route.GatewayID)
				}
				if route.NatGatewayID != nil {
					fmt.Printf("    Target: NAT Gateway (%s)\n", *route.NatGatewayID)
				}
				if route.TransitGatewayID != nil {
					fmt.Printf("    Target: Transit Gateway (%s)\n", *route.TransitGatewayID)
				}
				if route.VpcPeeringConnectionID != nil {
					fmt.Printf("    Target: VPC Peering (%s)\n", *route.VpcPeeringConnectionID)
				}
				if route.GatewayID == nil && route.NatGatewayID == nil && route.TransitGatewayID == nil && route.VpcPeeringConnectionID == nil {
					fmt.Printf("    Target: None\n")
				}
			}

			// Validate route table
			fmt.Printf("\nValidating Route Table...\n")
			err := test.routeTable.Validate()

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
					fmt.Printf("✅ PASSED: Route Table validation succeeded\n")
				}
			}

			// Additional validation checks
			if err == nil {
				fmt.Printf("\nAdditional validation checks...\n")
				targetCount := 0
				for _, route := range test.routeTable.Routes {
					if route.GatewayID != nil && *route.GatewayID != "" {
						targetCount++
					}
					if route.NatGatewayID != nil && *route.NatGatewayID != "" {
						targetCount++
					}
				}
				if targetCount > 0 {
					fmt.Printf("  ✅ PASSED: Route table has valid targets\n")
				}
			}

			fmt.Printf("=== Test completed: %s ===\n\n", test.name)
		})
	}
}
