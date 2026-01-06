package networking

import (
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomain converts AWS VPC to domain VPC
func ToDomainVPC(awsVPC *awsnetworking.VPC) *domainnetworking.VPC {
	if awsVPC == nil {
		return nil
	}
	
	domainVPC := &domainnetworking.VPC{
		Name:               awsVPC.Name,
		Region:             awsVPC.Region,
		CIDR:               awsVPC.CIDR,
		EnableDNS:          awsVPC.EnableDNSSupport,
		EnableDNSHostnames: awsVPC.EnableDNSHostnames,
	}
	
	return domainVPC
}

// FromDomain converts domain VPC to AWS VPC
func FromDomainVPC(domainVPC *domainnetworking.VPC) *awsnetworking.VPC {
	if domainVPC == nil {
		return nil
	}
	
	awsVPC := &awsnetworking.VPC{
		Name:               domainVPC.Name,
		Region:             domainVPC.Region,
		CIDR:               domainVPC.CIDR,
		EnableDNSSupport:   domainVPC.EnableDNS,
		EnableDNSHostnames: domainVPC.EnableDNSHostnames,
		InstanceTenancy:    "default", // AWS default
		Tags:               []configs.Tag{{Key: "Name", Value: domainVPC.Name}},
	}
	
	return awsVPC
}
