package networking

import (
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomainRouteTable converts AWS Route Table to domain Route Table (for backward compatibility)
func ToDomainRouteTable(awsRT *awsnetworking.RouteTable) *domainnetworking.RouteTable {
	if awsRT == nil {
		return nil
	}
	
	domainRoutes := make([]domainnetworking.Route, len(awsRT.Routes))
	for i, awsRoute := range awsRT.Routes {
		route := domainnetworking.Route{
			DestinationCIDR: awsRoute.DestinationCIDRBlock,
		}
		
		// Determine target type and ID
		if awsRoute.GatewayID != nil && *awsRoute.GatewayID != "" {
			route.TargetType = "internet_gateway"
			route.TargetID = *awsRoute.GatewayID
		} else if awsRoute.NatGatewayID != nil && *awsRoute.NatGatewayID != "" {
			route.TargetType = "nat_gateway"
			route.TargetID = *awsRoute.NatGatewayID
		} else if awsRoute.TransitGatewayID != nil && *awsRoute.TransitGatewayID != "" {
			route.TargetType = "transit_gateway"
			route.TargetID = *awsRoute.TransitGatewayID
		} else if awsRoute.VpcPeeringConnectionID != nil && *awsRoute.VpcPeeringConnectionID != "" {
			route.TargetType = "vpc_peering"
			route.TargetID = *awsRoute.VpcPeeringConnectionID
		}
		
		domainRoutes[i] = route
	}
	
	return &domainnetworking.RouteTable{
		Name:   awsRT.Name,
		VPCID:  awsRT.VPCID,
		Routes: domainRoutes,
	}
}

// ToDomainRouteTableFromOutput converts AWS Route Table output to domain Route Table with ID and ARN
func ToDomainRouteTableFromOutput(output *awsoutputs.RouteTableOutput) *domainnetworking.RouteTable {
	if output == nil {
		return nil
	}
	
	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}
	
	domainRoutes := make([]domainnetworking.Route, len(output.Routes))
	for i, awsRoute := range output.Routes {
		route := domainnetworking.Route{
			DestinationCIDR: awsRoute.DestinationCIDRBlock,
		}
		
		// Determine target type and ID
		if awsRoute.GatewayID != nil && *awsRoute.GatewayID != "" {
			route.TargetType = "internet_gateway"
			route.TargetID = *awsRoute.GatewayID
		} else if awsRoute.NatGatewayID != nil && *awsRoute.NatGatewayID != "" {
			route.TargetType = "nat_gateway"
			route.TargetID = *awsRoute.NatGatewayID
		} else if awsRoute.TransitGatewayID != nil && *awsRoute.TransitGatewayID != "" {
			route.TargetType = "transit_gateway"
			route.TargetID = *awsRoute.TransitGatewayID
		} else if awsRoute.VpcPeeringConnectionID != nil && *awsRoute.VpcPeeringConnectionID != "" {
			route.TargetType = "vpc_peering"
			route.TargetID = *awsRoute.VpcPeeringConnectionID
		}
		
		domainRoutes[i] = route
	}
	
	// Extract subnet IDs from associations
	subnetIDs := make([]string, 0, len(output.Associations))
	for _, assoc := range output.Associations {
		subnetIDs = append(subnetIDs, assoc.SubnetID)
	}
	
	return &domainnetworking.RouteTable{
		ID:      output.ID,
		ARN:     arn,
		Name:    output.Name,
		VPCID:   output.VPCID,
		Routes:  domainRoutes,
		Subnets: subnetIDs,
	}
}

// FromDomainRouteTable converts domain Route Table to AWS Route Table
func FromDomainRouteTable(domainRT *domainnetworking.RouteTable) *awsnetworking.RouteTable {
	if domainRT == nil {
		return nil
	}
	
	awsRoutes := make([]awsnetworking.Route, len(domainRT.Routes))
	for i, domainRoute := range domainRT.Routes {
		awsRoute := awsnetworking.Route{
			DestinationCIDRBlock: domainRoute.DestinationCIDR,
		}
		
		// Set target based on type
		switch domainRoute.TargetType {
		case "internet_gateway":
			awsRoute.GatewayID = &domainRoute.TargetID
		case "nat_gateway":
			awsRoute.NatGatewayID = &domainRoute.TargetID
		case "transit_gateway":
			awsRoute.TransitGatewayID = &domainRoute.TargetID
		case "vpc_peering":
			awsRoute.VpcPeeringConnectionID = &domainRoute.TargetID
		}
		
		awsRoutes[i] = awsRoute
	}
	
	return &awsnetworking.RouteTable{
		Name:   domainRT.Name,
		VPCID:  domainRT.VPCID,
		Routes: awsRoutes,
		Tags:   []configs.Tag{{Key: "Name", Value: domainRT.Name}},
	}
}
