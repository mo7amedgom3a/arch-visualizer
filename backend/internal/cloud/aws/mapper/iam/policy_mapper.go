package iam

import (
	domainiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/iam"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// FromDomainPolicy converts domain Policy to AWS Policy
func FromDomainPolicy(domainPolicy *domainiam.Policy) *awsiam.Policy {
	if domainPolicy == nil {
		return nil
	}

	awsPolicy := &awsiam.Policy{
		Name:          domainPolicy.Name,
		Description:   domainPolicy.Description,
		Path:          domainPolicy.Path,
		PolicyDocument: domainPolicy.PolicyDocument,
	}

	// Convert tags
	if len(domainPolicy.Tags) > 0 {
		awsPolicy.Tags = make([]configs.Tag, len(domainPolicy.Tags))
		for i, tag := range domainPolicy.Tags {
			awsPolicy.Tags[i] = configs.Tag{
				Key:   tag.Key,
				Value: tag.Value,
			}
		}
	}

	return awsPolicy
}

// ToDomainPolicyFromOutput converts AWS Policy output to domain Policy
func ToDomainPolicyFromOutput(output *awsoutputs.PolicyOutput) *domainiam.Policy {
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

	// Determine policy type based on ARN
	var policyType domainiam.PolicyType
	if output.IsAWSManaged {
		policyType = domainiam.PolicyTypeAWSManaged
	} else {
		policyType = domainiam.PolicyTypeCustomerManaged
	}

	isAttachable := output.IsAttachable

	domainPolicy := &domainiam.Policy{
		ID:            output.ID,
		ARN:           arn,
		Name:          output.Name,
		Description:   output.Description,
		Path:          path,
		PolicyDocument: output.PolicyDocument,
		Type:          policyType,
		IsAttachable:  &isAttachable,
	}

	// Convert tags
	if len(output.Tags) > 0 {
		domainPolicy.Tags = make([]domainiam.PolicyTag, len(output.Tags))
		for i, tag := range output.Tags {
			domainPolicy.Tags[i] = domainiam.PolicyTag{
				Key:   tag.Key,
				Value: tag.Value,
			}
		}
	}

	return domainPolicy
}
