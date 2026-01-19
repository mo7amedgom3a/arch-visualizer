package iam

import (
	domainiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/iam"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// FromDomainUser converts domain User to AWS User
func FromDomainUser(domainUser *domainiam.User) *awsiam.User {
	if domainUser == nil {
		return nil
	}

	awsUser := &awsiam.User{
		Name:                domainUser.Name,
		Path:                domainUser.Path,
		PermissionsBoundary: domainUser.PermissionsBoundary,
		ForceDestroy:        domainUser.ForceDestroy,
	}

	// Convert tags
	if len(domainUser.Tags) > 0 {
		awsUser.Tags = make([]configs.Tag, len(domainUser.Tags))
		for i, tag := range domainUser.Tags {
			awsUser.Tags[i] = configs.Tag{
				Key:   tag.Key,
				Value: tag.Value,
			}
		}
	}

	return awsUser
}

// ToDomainUserFromOutput converts AWS User output to domain User
func ToDomainUserFromOutput(output *awsoutputs.UserOutput) *domainiam.User {
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

	domainUser := &domainiam.User{
		ID:                output.ID,
		ARN:               arn,
		Name:              output.Name,
		Path:              path,
		PermissionsBoundary: output.PermissionsBoundary,
	}

	// Convert tags
	if len(output.Tags) > 0 {
		domainUser.Tags = make([]domainiam.PolicyTag, len(output.Tags))
		for i, tag := range output.Tags {
			domainUser.Tags[i] = domainiam.PolicyTag{
				Key:   tag.Key,
				Value: tag.Value,
			}
		}
	}

	return domainUser
}
