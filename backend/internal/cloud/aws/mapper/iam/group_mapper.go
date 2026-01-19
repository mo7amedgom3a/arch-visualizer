package iam

import (
	domainiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/iam"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// FromDomainGroup converts domain Group to AWS Group
func FromDomainGroup(domainGroup *domainiam.Group) *awsiam.Group {
	if domainGroup == nil {
		return nil
	}

	awsGroup := &awsiam.Group{
		Name: domainGroup.Name,
		Path: domainGroup.Path,
	}

	// Convert tags
	if len(domainGroup.Tags) > 0 {
		awsGroup.Tags = make([]configs.Tag, len(domainGroup.Tags))
		for i, tag := range domainGroup.Tags {
			awsGroup.Tags[i] = configs.Tag{
				Key:   tag.Key,
				Value: tag.Value,
			}
		}
	}

	return awsGroup
}

// ToDomainGroupFromOutput converts AWS Group output to domain Group
func ToDomainGroupFromOutput(output *awsoutputs.GroupOutput) *domainiam.Group {
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

	domainGroup := &domainiam.Group{
		ID:   output.ID,
		ARN:  arn,
		Name: output.Name,
		Path: path,
	}

	// Convert tags
	if len(output.Tags) > 0 {
		domainGroup.Tags = make([]domainiam.PolicyTag, len(output.Tags))
		for i, tag := range output.Tags {
			domainGroup.Tags[i] = domainiam.PolicyTag{
				Key:   tag.Key,
				Value: tag.Value,
			}
		}
	}

	return domainGroup
}
