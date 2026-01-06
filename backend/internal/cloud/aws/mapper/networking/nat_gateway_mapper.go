package networking

import (
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomainNATGateway converts AWS NAT Gateway to domain NAT Gateway
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
