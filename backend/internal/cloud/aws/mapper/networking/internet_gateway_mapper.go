package networking

import (
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomainInternetGateway converts AWS Internet Gateway to domain Internet Gateway (for backward compatibility)
func ToDomainInternetGateway(awsIGW *awsnetworking.InternetGateway) *domainnetworking.InternetGateway {
	if awsIGW == nil {
		return nil
	}
	
	return &domainnetworking.InternetGateway{
		Name:  awsIGW.Name,
		VPCID: awsIGW.VPCID,
	}
}

// ToDomainInternetGatewayFromOutput converts AWS Internet Gateway output to domain Internet Gateway with ID and ARN
func ToDomainInternetGatewayFromOutput(output *awsoutputs.InternetGatewayOutput) *domainnetworking.InternetGateway {
	if output == nil {
		return nil
	}
	
	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}
	
	return &domainnetworking.InternetGateway{
		ID:    output.ID,
		ARN:   arn,
		Name:  output.Name,
		VPCID: output.VPCID,
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
