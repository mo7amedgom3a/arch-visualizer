package networking

import (
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomainInternetGateway converts AWS Internet Gateway to domain Internet Gateway
func ToDomainInternetGateway(awsIGW *awsnetworking.InternetGateway) *domainnetworking.InternetGateway {
	if awsIGW == nil {
		return nil
	}
	
	return &domainnetworking.InternetGateway{
		Name:  awsIGW.Name,
		VPCID: awsIGW.VPCID,
	}
}

// FromDomainInternetGateway converts domain Internet Gateway to AWS Internet Gateway
func FromDomainInternetGateway(domainIGW *domainnetworking.InternetGateway) *awsnetworking.InternetGateway {
	if domainIGW == nil {
		return nil
	}
	
	return &awsnetworking.InternetGateway{
		Name:  domainIGW.Name,
		VPCID: domainIGW.VPCID,
		Tags:  []configs.Tag{{Key: "Name", Value: domainIGW.Name}},
	}
}
