package iam

import (
	domainiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/iam"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// FromDomainRole converts domain Role to AWS Role
func FromDomainRole(domainRole *domainiam.Role) *awsiam.Role {
	if domainRole == nil {
		return nil
	}

	awsRole := &awsiam.Role{
		Name:              domainRole.Name,
		Description:       domainRole.Description,
		Path:              domainRole.Path,
		AssumeRolePolicy:  domainRole.AssumeRolePolicy,
		PermissionsBoundary: domainRole.PermissionsBoundary,
	}

	// Convert tags
	if len(domainRole.Tags) > 0 {
		awsRole.Tags = make([]configs.Tag, len(domainRole.Tags))
		for i, tag := range domainRole.Tags {
			awsRole.Tags[i] = configs.Tag{
				Key:   tag.Key,
				Value: tag.Value,
			}
		}
	}

	return awsRole
}

// ToDomainRoleFromOutput converts AWS Role output to domain Role
func ToDomainRoleFromOutput(output *awsoutputs.RoleOutput) *domainiam.Role {
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

	domainRole := &domainiam.Role{
		ID:                output.ID,
		ARN:               arn,
		Name:              output.Name,
		Description:       output.Description,
		Path:              path,
		AssumeRolePolicy:  output.AssumeRolePolicy,
		PermissionsBoundary: output.PermissionsBoundary,
	}

	// Convert tags
	if len(output.Tags) > 0 {
		domainRole.Tags = make([]domainiam.PolicyTag, len(output.Tags))
		for i, tag := range output.Tags {
			domainRole.Tags[i] = domainiam.PolicyTag{
				Key:   tag.Key,
				Value: tag.Value,
			}
		}
	}

	return domainRole
}
