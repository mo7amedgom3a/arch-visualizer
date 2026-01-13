package networking

import (
	"errors"
)

// Route represents a single route entry
type Route struct {
	DestinationCIDR string // Destination CIDR block (e.g., "0.0.0.0/0")
	TargetID        string // Target resource ID (IGW, NAT Gateway, etc.)
	TargetType      string // Target type (internet_gateway, nat_gateway, etc.)
}

// RouteTable represents a cloud-agnostic route table
type RouteTable struct {
	ID      string
	ARN     *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	Name    string
	VPCID   string
	Routes  []Route
	Subnets []string // Associated subnet IDs
}

// Validate performs domain-level validation
func (rt *RouteTable) Validate() error {
	if rt.Name == "" {
		return errors.New("route table name is required")
	}
	if rt.VPCID == "" {
		return errors.New("route table vpc_id is required")
	}
	
	// Validate routes
	for i, route := range rt.Routes {
		if route.DestinationCIDR == "" {
			return errors.New("route destination_cidr is required")
		}
		if route.TargetID == "" {
			return errors.New("route target_id is required")
		}
		if route.TargetType == "" {
			return errors.New("route target_type is required")
		}
		
		// Validate target types
		validTargetTypes := map[string]bool{
			"internet_gateway": true,
			"nat_gateway":     true,
			"vpc_peering":     true,
			"transit_gateway": true,
			"local":           true,
		}
		if !validTargetTypes[route.TargetType] {
			return errors.New("invalid route target_type")
		}
		
		// Prevent duplicate routes
		for j := i + 1; j < len(rt.Routes); j++ {
			if rt.Routes[j].DestinationCIDR == route.DestinationCIDR {
				return errors.New("duplicate route destination")
			}
		}
	}
	
	return nil
}
