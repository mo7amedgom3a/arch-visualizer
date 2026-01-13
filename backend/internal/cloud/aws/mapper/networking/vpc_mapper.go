package networking

import (
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomain converts AWS VPC to domain VPC (for backward compatibility)
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

// ToDomainVPCFromOutput converts AWS VPC output to domain VPC with ID and ARN
func ToDomainVPCFromOutput(output *awsoutputs.VPCOutput) *domainnetworking.VPC {
	if output == nil {
		return nil
	}
	
	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}
	
	domainVPC := &domainnetworking.VPC{
		ID:                 output.ID,
		ARN:                arn,
		Name:               output.Name,
		Region:             output.Region,
		CIDR:               output.CIDR,
		EnableDNS:          output.EnableDNSSupport,
		EnableDNSHostnames: output.EnableDNSHostnames,
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
