package networking

import (
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomainNATGateway converts AWS NAT Gateway to domain NAT Gateway (for backward compatibility)
func ToDomainNATGateway(awsNAT *awsnetworking.NATGateway) *domainnetworking.NATGateway {
	if awsNAT == nil {
		return nil
	}
	
	return &domainnetworking.NATGateway{
		Name:          awsNAT.Name,
		SubnetID:      awsNAT.SubnetID,
		AllocationID:  &awsNAT.AllocationID,
	}
}

// ToDomainNATGatewayFromOutput converts AWS NAT Gateway output to domain NAT Gateway with ID and ARN
func ToDomainNATGatewayFromOutput(output *awsoutputs.NATGatewayOutput) *domainnetworking.NATGateway {
	if output == nil {
		return nil
	}
	
	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}
	
	allocationID := &output.AllocationID
	if output.AllocationID == "" {
		allocationID = nil
	}
	
	return &domainnetworking.NATGateway{
		ID:           output.ID,
		ARN:          arn,
		Name:         output.Name,
		SubnetID:     output.SubnetID,
		AllocationID: allocationID,
	}
}

// FromDomainNATGateway converts domain NAT Gateway to AWS NAT Gateway
func FromDomainNATGateway(domainNAT *domainnetworking.NATGateway) *awsnetworking.NATGateway {
	if domainNAT == nil {
		return nil
	}
	
	allocationID := ""
	if domainNAT.AllocationID != nil {
		allocationID = *domainNAT.AllocationID
	}
	
	return &awsnetworking.NATGateway{
		Name:         domainNAT.Name,
		SubnetID:     domainNAT.SubnetID,
		AllocationID: allocationID,
		Tags:         []configs.Tag{{Key: "Name", Value: domainNAT.Name}},
	}
}
