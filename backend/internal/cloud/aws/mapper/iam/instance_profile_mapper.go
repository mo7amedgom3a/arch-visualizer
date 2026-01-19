package iam

import (
	domainiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/iam"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// FromDomainInstanceProfile converts domain InstanceProfile to AWS InstanceProfile
func FromDomainInstanceProfile(domainProfile *domainiam.InstanceProfile) *awsiam.InstanceProfile {
	if domainProfile == nil {
		return nil
	}

	awsProfile := &awsiam.InstanceProfile{
		Name: domainProfile.Name,
		Path: domainProfile.Path,
		Role: domainProfile.RoleName,
	}

	// Convert tags
	if len(domainProfile.Tags) > 0 {
		awsProfile.Tags = make([]configs.Tag, len(domainProfile.Tags))
		for i, tag := range domainProfile.Tags {
			awsProfile.Tags[i] = configs.Tag{
				Key:   tag.Key,
				Value: tag.Value,
			}
		}
	}

	return awsProfile
}

// ToDomainInstanceProfileFromOutput converts AWS InstanceProfile output to domain InstanceProfile
func ToDomainInstanceProfileFromOutput(output *awsoutputs.InstanceProfileOutput) *domainiam.InstanceProfile {
	if output == nil {
		return nil
	}

	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}

	path := &output.Path
	if output.Path == "" {
		defaultPath := "/"
		path = &defaultPath
	}

	// Get role name from attached roles (if any)
	var roleName *string
	if len(output.Roles) > 0 && output.Roles[0] != nil {
		roleName = &output.Roles[0].Name
	}

	domainProfile := &domainiam.InstanceProfile{
		ID:       output.ID,
		ARN:      arn,
		Name:     output.Name,
		Path:     path,
		RoleName: roleName,
	}

	// Convert tags
	if len(output.Tags) > 0 {
		domainProfile.Tags = make([]domainiam.PolicyTag, len(output.Tags))
		for i, tag := range output.Tags {
			domainProfile.Tags[i] = domainiam.PolicyTag{
				Key:   tag.Key,
				Value: tag.Value,
			}
		}
	}

	return domainProfile
}
