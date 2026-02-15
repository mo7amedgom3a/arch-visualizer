package networking

import (
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
)

// ToDomainVPCEndpoint converts AWS VPC Endpoint to domain VPC Endpoint
func ToDomainVPCEndpoint(awsVPCE *awsnetworking.VPCEndpoint) *domainnetworking.VPCEndpoint {
	if awsVPCE == nil {
		return nil
	}

	return &domainnetworking.VPCEndpoint{
		Name:              awsVPCE.Name,
		VPCID:             awsVPCE.VPCID,
		ServiceName:       awsVPCE.ServiceName,
		Type:              domainnetworking.VPCEndpointType(awsVPCE.VPCEndpointType),
		SubnetIDs:         awsVPCE.SubnetIDs,
		SecurityGroupIDs:  awsVPCE.SecurityGroupIDs,
		RouteTableIDs:     awsVPCE.RouteTableIDs,
		PrivateDNSEnabled: awsVPCE.PrivateDNSEnabled,
		Policy:            awsVPCE.Policy,
		Tags:              awsVPCE.Tags,
	}
}

// ToDomainVPCEndpointFromOutput converts AWS VPCOutput to domain VPCEndpoint
// Note: VPCOutput is used as a placeholder in the current service implementation
func ToDomainVPCEndpointFromOutput(output *awsoutputs.VPCOutput, original *domainnetworking.VPCEndpoint) *domainnetworking.VPCEndpoint {
	if output == nil || original == nil {
		return nil
	}

	// Copy original fields
	domainVPCE := &domainnetworking.VPCEndpoint{
		Name:              original.Name,
		VPCID:             original.VPCID,
		ServiceName:       original.ServiceName,
		Type:              original.Type,
		SubnetIDs:         original.SubnetIDs,
		SecurityGroupIDs:  original.SecurityGroupIDs,
		RouteTableIDs:     original.RouteTableIDs,
		PrivateDNSEnabled: original.PrivateDNSEnabled,
		Policy:            original.Policy,
		Tags:              original.Tags,
	}

	// Update with output ID and ARN (if available in future output types)
	domainVPCE.ID = output.ID
	if output.ARN != "" {
		domainVPCE.ARN = &output.ARN
	}

	return domainVPCE
}

// FromDomainVPCEndpoint converts domain VPC Endpoint to AWS VPC Endpoint
func FromDomainVPCEndpoint(domainVPCE *domainnetworking.VPCEndpoint) *awsnetworking.VPCEndpoint {
	if domainVPCE == nil {
		return nil
	}

	return &awsnetworking.VPCEndpoint{
		Name:              domainVPCE.Name,
		VPCID:             domainVPCE.VPCID,
		ServiceName:       domainVPCE.ServiceName,
		VPCEndpointType:   awsnetworking.VPCEndpointType(domainVPCE.Type),
		SubnetIDs:         domainVPCE.SubnetIDs,
		SecurityGroupIDs:  domainVPCE.SecurityGroupIDs,
		RouteTableIDs:     domainVPCE.RouteTableIDs,
		PrivateDNSEnabled: domainVPCE.PrivateDNSEnabled,
		Policy:            domainVPCE.Policy,
		Tags:              domainVPCE.Tags,
	}
}
