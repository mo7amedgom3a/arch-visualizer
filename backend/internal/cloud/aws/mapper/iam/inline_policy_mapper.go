package iam

import (
	domainiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/iam"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
)

// FromDomainInlinePolicy converts domain InlinePolicy to AWS InlinePolicy
func FromDomainInlinePolicy(domainPolicy *domainiam.InlinePolicy) *awsiam.InlinePolicy {
	if domainPolicy == nil {
		return nil
	}

	return &awsiam.InlinePolicy{
		Name:   domainPolicy.Name,
		Policy: domainPolicy.Policy,
	}
}

// ToDomainInlinePolicy converts AWS InlinePolicy to domain InlinePolicy
func ToDomainInlinePolicy(awsPolicy *awsiam.InlinePolicy) *domainiam.InlinePolicy {
	if awsPolicy == nil {
		return nil
	}

	return &domainiam.InlinePolicy{
		Name:   awsPolicy.Name,
		Policy: awsPolicy.Policy,
	}
}
