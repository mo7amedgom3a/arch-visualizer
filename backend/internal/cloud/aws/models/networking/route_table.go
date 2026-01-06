package networking

import (
	"errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// Route represents an AWS route entry
type Route struct {
	DestinationCIDRBlock string `json:"destination_cidr_block"`
	GatewayID           *string `json:"gateway_id,omitempty"` // IGW or NAT Gateway ID
	NatGatewayID        *string `json:"nat_gateway_id,omitempty"`
	TransitGatewayID    *string `json:"transit_gateway_id,omitempty"`
	VpcPeeringConnectionID *string `json:"vpc_peering_connection_id,omitempty"`
}

// RouteTable represents an AWS-specific Route Table
type RouteTable struct {
	Name    string `json:"name"`
	VPCID   string `json:"vpc_id"`
	Routes  []Route `json:"routes"`
	Tags    []configs.Tag `json:"tags"`
}

// Validate performs AWS-specific validation
func (rt *RouteTable) Validate() error {
	if rt.Name == "" {
		return errors.New("route table name is required")
	}
	if rt.VPCID == "" {
		return errors.New("route table vpc_id is required")
	}
	
	// Validate routes
	for _, route := range rt.Routes {
		if route.DestinationCIDRBlock == "" {
			return errors.New("route destination_cidr_block is required")
		}
		
		// At least one target must be specified
		targetCount := 0
		if route.GatewayID != nil && *route.GatewayID != "" {
			targetCount++
		}
		if route.NatGatewayID != nil && *route.NatGatewayID != "" {
			targetCount++
		}
		if route.TransitGatewayID != nil && *route.TransitGatewayID != "" {
			targetCount++
		}
		if route.VpcPeeringConnectionID != nil && *route.VpcPeeringConnectionID != "" {
			targetCount++
		}
		
		if targetCount == 0 {
			return errors.New("route must have at least one target (gateway_id, nat_gateway_id, etc.)")
		}
		
		if targetCount > 1 {
			return errors.New("route can only have one target")
		}
	}
	
	return nil
}
