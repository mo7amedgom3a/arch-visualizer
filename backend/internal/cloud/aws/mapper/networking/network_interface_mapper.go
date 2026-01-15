package networking

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
)

// ToDomainNetworkInterface converts AWS Network Interface to domain Network Interface (for backward compatibility)
func ToDomainNetworkInterface(awsENI *awsnetworking.NetworkInterface) *domainnetworking.NetworkInterface {
	if awsENI == nil {
		return nil
	}

	sourceDestCheck := true
	if awsENI.SourceDestCheck != nil {
		sourceDestCheck = *awsENI.SourceDestCheck
	}

	return &domainnetworking.NetworkInterface{
		Description:          awsENI.Description,
		SubnetID:             awsENI.SubnetID,
		InterfaceType:        domainnetworking.NetworkInterfaceType(awsENI.InterfaceType),
		PrivateIPv4Address:   awsENI.PrivateIPv4Address,
		AutoAssignPrivateIP:  awsENI.AutoAssignPrivateIP,
		SecurityGroupIDs:     awsENI.SecurityGroupIDs,
		SourceDestCheck:      sourceDestCheck,
		IPv4PrefixDelegation: awsENI.IPv4PrefixDelegation,
	}
}

// ToDomainNetworkInterfaceFromOutput converts AWS Network Interface output to domain Network Interface with ID and ARN
func ToDomainNetworkInterfaceFromOutput(output *awsoutputs.NetworkInterfaceOutput) *domainnetworking.NetworkInterface {
	if output == nil {
		return nil
	}

	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}

	publicIPv4 := output.PublicIPv4Address
	if publicIPv4 != nil && *publicIPv4 == "" {
		publicIPv4 = nil
	}

	macAddress := &output.MACAddress
	if output.MACAddress == "" {
		macAddress = nil
	}

	az := &output.AvailabilityZone
	if output.AvailabilityZone == "" {
		az = nil
	}

	// Convert attachment
	var attachment *domainnetworking.NetworkInterfaceAttachment
	if output.Attachment != nil {
		attachment = &domainnetworking.NetworkInterfaceAttachment{
			AttachmentID:        output.Attachment.AttachmentID,
			InstanceID:          output.Attachment.InstanceID,
			DeviceIndex:         output.Attachment.DeviceIndex,
			Status:              output.Attachment.Status,
			DeleteOnTermination: output.Attachment.DeleteOnTermination,
		}
	}

	return &domainnetworking.NetworkInterface{
		ID:                   output.ID,
		ARN:                  arn,
		Description:          output.Description,
		SubnetID:             output.SubnetID,
		InterfaceType:        domainnetworking.NetworkInterfaceType(output.InterfaceType),
		PrivateIPv4Address:   &output.PrivateIPv4Address,
		AutoAssignPrivateIP:  true, // If we got an IP, it was assigned
		PublicIPv4Address:    publicIPv4,
		IPv6Addresses:        output.IPv6Addresses,
		SecurityGroupIDs:     output.SecurityGroupIDs,
		SourceDestCheck:      output.SourceDestCheck,
		IPv4PrefixDelegation: output.IPv4PrefixDelegation,
		MACAddress:           macAddress,
		Status:               domainnetworking.NetworkInterfaceStatus(output.Status),
		Attachment:           attachment,
		VPCID:                output.VPCID,
		AvailabilityZone:     az,
	}
}

// FromDomainNetworkInterface converts domain Network Interface to AWS Network Interface
func FromDomainNetworkInterface(domainENI *domainnetworking.NetworkInterface) *awsnetworking.NetworkInterface {
	if domainENI == nil {
		return nil
	}

	sourceDestCheck := &domainENI.SourceDestCheck

	return &awsnetworking.NetworkInterface{
		Description:          domainENI.Description,
		SubnetID:             domainENI.SubnetID,
		InterfaceType:        awsnetworking.NetworkInterfaceType(domainENI.InterfaceType),
		PrivateIPv4Address:   domainENI.PrivateIPv4Address,
		AutoAssignPrivateIP:  domainENI.AutoAssignPrivateIP,
		SecurityGroupIDs:     domainENI.SecurityGroupIDs,
		SourceDestCheck:      sourceDestCheck,
		IPv4PrefixDelegation: domainENI.IPv4PrefixDelegation,
		Tags:                 []configs.Tag{}, // Default empty tags
	}
}
