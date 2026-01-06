package networking

import (
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomainSubnet converts AWS Subnet to domain Subnet
func ToDomainSubnet(awsSubnet *awsnetworking.Subnet) *domainnetworking.Subnet {
	if awsSubnet == nil {
		return nil
	}
	
	domainSubnet := &domainnetworking.Subnet{
		Name:            awsSubnet.Name,
		VPCID:           awsSubnet.VPCID,
		CIDR:            awsSubnet.CIDR,
		AvailabilityZone: &awsSubnet.AvailabilityZone,
		IsPublic:        awsSubnet.MapPublicIPOnLaunch,
	}
	
	return domainSubnet
}

// FromDomainSubnet converts domain Subnet to AWS Subnet
func FromDomainSubnet(domainSubnet *domainnetworking.Subnet, availabilityZone string) *awsnetworking.Subnet {
	if domainSubnet == nil {
		return nil
	}
	
	az := availabilityZone
	if domainSubnet.AvailabilityZone != nil {
		az = *domainSubnet.AvailabilityZone
	}
	
	awsSubnet := &awsnetworking.Subnet{
		Name:              domainSubnet.Name,
		VPCID:             domainSubnet.VPCID,
		CIDR:              domainSubnet.CIDR,
		AvailabilityZone:  az,
		MapPublicIPOnLaunch: domainSubnet.IsPublic,
		Tags:              []configs.Tag{{Key: "Name", Value: domainSubnet.Name}},
	}
	
	return awsSubnet
}
