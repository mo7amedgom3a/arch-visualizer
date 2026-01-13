package networking

import (
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomainSubnet converts AWS Subnet to domain Subnet (for backward compatibility)
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

// ToDomainSubnetFromOutput converts AWS Subnet output to domain Subnet with ID and ARN
func ToDomainSubnetFromOutput(output *awsoutputs.SubnetOutput) *domainnetworking.Subnet {
	if output == nil {
		return nil
	}
	
	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}
	
	domainSubnet := &domainnetworking.Subnet{
		ID:              output.ID,
		ARN:             arn,
		Name:            output.Name,
		VPCID:           output.VPCID,
		CIDR:            output.CIDR,
		AvailabilityZone: &output.AvailabilityZone,
		IsPublic:        output.MapPublicIPOnLaunch,
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
