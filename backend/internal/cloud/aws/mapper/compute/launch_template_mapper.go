package compute

import (
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	awslttemplate "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template"
	awslttemplateoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// FromDomainLaunchTemplate converts domain LaunchTemplate to AWS LaunchTemplate input model
func FromDomainLaunchTemplate(domain *domaincompute.LaunchTemplate) *awslttemplate.LaunchTemplate {
	if domain == nil {
		return nil
	}

	awsTemplate := &awslttemplate.LaunchTemplate{
		ImageID:      domain.ImageID,
		InstanceType: domain.InstanceType,
		VpcSecurityGroupIds: domain.SecurityGroupIDs,
	}

	// Handle naming (prefer name_prefix if available)
	if domain.NamePrefix != nil && *domain.NamePrefix != "" {
		awsTemplate.NamePrefix = domain.NamePrefix
	} else if domain.Name != "" {
		awsTemplate.Name = &domain.Name
	}

	// Key name
	if domain.KeyName != nil {
		awsTemplate.KeyName = domain.KeyName
	}

	// IAM instance profile (convert from string to struct)
	if domain.IAMInstanceProfile != nil {
		awsTemplate.IAMInstanceProfile = &awslttemplate.IAMInstanceProfile{
			Name: domain.IAMInstanceProfile,
		}
	}

	// User data
	if domain.UserData != nil {
		awsTemplate.UserData = domain.UserData
	}

	// Root volume ID
	if domain.RootVolumeID != nil {
		awsTemplate.RootVolumeID = domain.RootVolumeID
	}

	// Additional volume IDs
	if len(domain.AdditionalVolumeIDs) > 0 {
		awsTemplate.AdditionalVolumeIDs = domain.AdditionalVolumeIDs
	}

	// Metadata options
	if domain.MetadataOptions != nil {
		awsTemplate.MetadataOptions = &awslttemplate.MetadataOptions{
			HTTPEndpoint:            domain.MetadataOptions.HTTPEndpoint,
			HTTPTokens:              domain.MetadataOptions.HTTPTokens,
			HTTPPutResponseHopLimit: domain.MetadataOptions.HTTPPutResponseHopLimit,
		}
	}

	// Update default version (default to true)
	updateDefault := true
	awsTemplate.UpdateDefaultVersion = &updateDefault

	// Tags (add Name tag if name is provided)
	if domain.Name != "" {
		awsTemplate.Tags = []configs.Tag{
			{Key: "Name", Value: domain.Name},
		}
	} else if domain.NamePrefix != nil {
		awsTemplate.Tags = []configs.Tag{
			{Key: "Name", Value: *domain.NamePrefix},
		}
	}

	return awsTemplate
}

// ToDomainLaunchTemplate converts AWS LaunchTemplate input model to domain LaunchTemplate
// This is useful for backward compatibility or when reading existing templates
func ToDomainLaunchTemplate(aws *awslttemplate.LaunchTemplate) *domaincompute.LaunchTemplate {
	if aws == nil {
		return nil
	}

	domain := &domaincompute.LaunchTemplate{
		ImageID:      aws.ImageID,
		InstanceType: aws.InstanceType,
		SecurityGroupIDs: aws.VpcSecurityGroupIds,
	}

	// Handle naming
	if aws.NamePrefix != nil {
		domain.NamePrefix = aws.NamePrefix
	}
	if aws.Name != nil {
		domain.Name = *aws.Name
	}

	// Key name
	if aws.KeyName != nil {
		domain.KeyName = aws.KeyName
	}

	// IAM instance profile (convert from struct to string)
	if aws.IAMInstanceProfile != nil && aws.IAMInstanceProfile.Name != nil {
		domain.IAMInstanceProfile = aws.IAMInstanceProfile.Name
	}

	// User data
	if aws.UserData != nil {
		domain.UserData = aws.UserData
	}

	// Root volume ID
	if aws.RootVolumeID != nil {
		domain.RootVolumeID = aws.RootVolumeID
	}

	// Additional volume IDs
	if len(aws.AdditionalVolumeIDs) > 0 {
		domain.AdditionalVolumeIDs = aws.AdditionalVolumeIDs
	}

	// Metadata options
	if aws.MetadataOptions != nil {
		domain.MetadataOptions = &domaincompute.MetadataOptions{
			HTTPEndpoint:            aws.MetadataOptions.HTTPEndpoint,
			HTTPTokens:              aws.MetadataOptions.HTTPTokens,
			HTTPPutResponseHopLimit: aws.MetadataOptions.HTTPPutResponseHopLimit,
		}
	}

	return domain
}

// ToDomainLaunchTemplateFromOutput converts AWS LaunchTemplateOutput to domain LaunchTemplate
// This populates the domain model with AWS-generated identifiers (ID, ARN, versions)
func ToDomainLaunchTemplateFromOutput(output *awslttemplateoutputs.LaunchTemplateOutput) *domaincompute.LaunchTemplate {
	if output == nil {
		return nil
	}

	domain := &domaincompute.LaunchTemplate{
		ID:     output.ID,
		ARN:    &output.ARN,
		Name:   output.Name,
		Region: "", // Region should be set from context
	}

	// Version information
	if output.DefaultVersion > 0 {
		domain.Version = &output.DefaultVersion
	}
	if output.LatestVersion > 0 {
		domain.LatestVersion = &output.LatestVersion
	}

	return domain
}

// ToDomainLaunchTemplateOutputFromOutput converts AWS LaunchTemplateOutput directly to domain LaunchTemplateOutput
func ToDomainLaunchTemplateOutputFromOutput(output *awslttemplateoutputs.LaunchTemplateOutput) *domaincompute.LaunchTemplateOutput {
	if output == nil {
		return nil
	}

	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}

	// LaunchTemplateOutput doesn't have NamePrefix
	var namePrefix *string

	defaultVersion := &output.DefaultVersion
	if output.DefaultVersion == 0 {
		defaultVersion = nil
	}

	latestVersion := &output.LatestVersion
	if output.LatestVersion == 0 {
		latestVersion = nil
	}

	createdAt := &output.CreateTime
	createdBy := &output.CreatedBy

	return &domaincompute.LaunchTemplateOutput{
		ID:            output.ID,
		ARN:           arn,
		Name:          output.Name,
		Region:        "", // Region should be set from context
		NamePrefix:    namePrefix,
		DefaultVersion: defaultVersion,
		LatestVersion:  latestVersion,
		CreatedAt:     createdAt,
		CreatedBy:     createdBy,
	}
}

// ToDomainLaunchTemplateVersion converts AWS LaunchTemplateVersion to domain LaunchTemplateVersion
func ToDomainLaunchTemplateVersion(aws *awslttemplate.LaunchTemplateVersion) *domaincompute.LaunchTemplateVersion {
	if aws == nil {
		return nil
	}

	domain := &domaincompute.LaunchTemplateVersion{
		TemplateID:    aws.TemplateID,
		VersionNumber: aws.VersionNumber,
		IsDefault:     aws.IsDefault,
		CreateTime:    aws.CreateTime,
		CreatedBy:     aws.CreatedBy,
	}

	// Convert template data if available
	if aws.TemplateData != nil {
		domain.TemplateData = ToDomainLaunchTemplate(aws.TemplateData)
	}

	return domain
}
