package networking

import (
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomainElasticIP converts AWS Elastic IP to domain Elastic IP (for backward compatibility)
func ToDomainElasticIP(awsEIP *awsnetworking.ElasticIP) *domainnetworking.ElasticIP {
	if awsEIP == nil {
		return nil
	}
	
	var poolType *domainnetworking.ElasticIPAddressPoolType
	if awsEIP.AddressPoolType != nil {
		dt := domainnetworking.ElasticIPAddressPoolType(*awsEIP.AddressPoolType)
		poolType = &dt
	}
	
	return &domainnetworking.ElasticIP{
		AllocationID:       awsEIP.AllocationID,
		AddressPoolType:    poolType,
		AddressPoolID:      awsEIP.AddressPoolID,
		NetworkBorderGroup: awsEIP.NetworkBorderGroup,
		Region:             awsEIP.Region,
	}
}

// ToDomainElasticIPFromOutput converts AWS Elastic IP output to domain Elastic IP with ID and ARN
func ToDomainElasticIPFromOutput(output *awsoutputs.ElasticIPOutput) *domainnetworking.ElasticIP {
	if output == nil {
		return nil
	}
	
	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}
	
	publicIP := &output.PublicIP
	if output.PublicIP == "" {
		publicIP = nil
	}
	
	allocationID := &output.AllocationID
	if output.AllocationID == "" {
		allocationID = nil
	}
	
	return &domainnetworking.ElasticIP{
		ID:                 output.ID,
		ARN:                arn,
		PublicIP:           publicIP,
		AllocationID:       allocationID,
		NetworkBorderGroup: output.NetworkBorderGroup,
		Region:             output.Region,
	}
}

// FromDomainElasticIP converts domain Elastic IP to AWS Elastic IP
func FromDomainElasticIP(domainEIP *domainnetworking.ElasticIP) *awsnetworking.ElasticIP {
	if domainEIP == nil {
		return nil
	}
	
	var poolType *awsnetworking.ElasticIPAddressPoolType
	if domainEIP.AddressPoolType != nil {
		at := awsnetworking.ElasticIPAddressPoolType(*domainEIP.AddressPoolType)
		poolType = &at
	}
	
	return &awsnetworking.ElasticIP{
		AllocationID:       domainEIP.AllocationID,
		AddressPoolType:    poolType,
		AddressPoolID:      domainEIP.AddressPoolID,
		NetworkBorderGroup: domainEIP.NetworkBorderGroup,
		Region:             domainEIP.Region,
		Tags:               []configs.Tag{}, // Default empty tags
	}
}
