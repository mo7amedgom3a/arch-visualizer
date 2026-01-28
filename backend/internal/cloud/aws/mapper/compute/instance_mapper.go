package compute

import (
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	awsec2 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2"
	awsec2outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomainInstance converts AWS EC2 instance to domain instance (for backward compatibility)
func ToDomainInstance(awsInstance *awsec2.Instance) *domaincompute.Instance {
	if awsInstance == nil {
		return nil
	}

	domainInstance := &domaincompute.Instance{
		Name:               awsInstance.Name,
		InstanceType:       awsInstance.InstanceType,
		AMI:                awsInstance.AMI,
		SubnetID:           awsInstance.SubnetID,
		SecurityGroupIDs:   awsInstance.VpcSecurityGroupIds,
		KeyName:            awsInstance.KeyName,
		IAMInstanceProfile: awsInstance.IAMInstanceProfile,
		RootVolumeID:       awsInstance.RootVolumeID,
	}

	return domainInstance
}

// ToDomainInstanceFromOutput converts AWS EC2 instance output to domain instance with ID and ARN
func ToDomainInstanceFromOutput(output *awsec2outputs.InstanceOutput) *domaincompute.Instance {
	if output == nil {
		return nil
	}

	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}

	az := &output.AvailabilityZone
	if output.AvailabilityZone == "" {
		az = nil
	}

	domainInstance := &domaincompute.Instance{
		ID:                output.ID,
		ARN:               arn,
		Name:              output.Name,
		Region:            output.Region,
		AvailabilityZone:  az,
		InstanceType:      output.InstanceType,
		AMI:               output.AMI,
		SubnetID:          output.SubnetID,
		SecurityGroupIDs:  output.SecurityGroupIDs,
		PrivateIP:         &output.PrivateIP,
		PublicIP:          output.PublicIP,
		KeyName:           output.KeyName,
		IAMInstanceProfile: output.IAMInstanceProfile,
		State:             domaincompute.InstanceState(output.State),
	}

	return domainInstance
}

// FromDomainInstance converts domain instance to AWS EC2 instance
func FromDomainInstance(domainInstance *domaincompute.Instance) *awsec2.Instance {
	if domainInstance == nil {
		return nil
	}

	awsInstance := &awsec2.Instance{
		Name:                domainInstance.Name,
		AMI:                 domainInstance.AMI,
		InstanceType:        domainInstance.InstanceType,
		SubnetID:            domainInstance.SubnetID,
		VpcSecurityGroupIds: domainInstance.SecurityGroupIDs,
		KeyName:             domainInstance.KeyName,
		IAMInstanceProfile:  domainInstance.IAMInstanceProfile,
		RootVolumeID:        domainInstance.RootVolumeID,
		Tags:                []configs.Tag{{Key: "Name", Value: domainInstance.Name}},
	}

	// Set AssociatePublicIPAddress based on PublicIP field
	// If PublicIP is set in domain, we want to associate public IP
	if domainInstance.PublicIP != nil && *domainInstance.PublicIP != "" {
		associatePublicIP := true
		awsInstance.AssociatePublicIPAddress = &associatePublicIP
	}

	return awsInstance
}

// ToDomainInstanceOutputFromOutput converts AWS EC2 instance output directly to domain InstanceOutput
func ToDomainInstanceOutputFromOutput(output *awsec2outputs.InstanceOutput) *domaincompute.InstanceOutput {
	if output == nil {
		return nil
	}

	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}

	az := &output.AvailabilityZone
	if output.AvailabilityZone == "" {
		az = nil
	}

	privateDNS := &output.PrivateDNS
	if output.PrivateDNS == "" {
		privateDNS = nil
	}

	vpcID := &output.VPCID
	if output.VPCID == "" {
		vpcID = nil
	}

	createdAt := &output.CreationTime

	return &domaincompute.InstanceOutput{
		ID:                 output.ID,
		ARN:                arn,
		Name:               output.Name,
		Region:             output.Region,
		AvailabilityZone:   az,
		InstanceType:       output.InstanceType,
		AMI:                output.AMI,
		SubnetID:           output.SubnetID,
		SecurityGroupIDs:   output.SecurityGroupIDs,
		PrivateIP:          &output.PrivateIP,
		PublicIP:           output.PublicIP,
		PrivateDNS:         privateDNS,
		PublicDNS:          output.PublicDNS,
		VPCID:              vpcID,
		KeyName:            output.KeyName,
		IAMInstanceProfile: output.IAMInstanceProfile,
		State:              domaincompute.InstanceState(output.State),
		CreatedAt:          createdAt,
	}
}
