package sdk

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awslttemplate "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template"
	awslttemplateoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template/outputs"
)

// CreateLaunchTemplate creates a new Launch Template using AWS SDK
func CreateLaunchTemplate(ctx context.Context, client *AWSClient, template *awslttemplate.LaunchTemplate) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	if err := template.Validate(); err != nil {
		return nil, fmt.Errorf("template validation failed: %w", err)
	}

	// Build LaunchTemplateData
	launchTemplateData := buildLaunchTemplateData(template)

	// Build CreateLaunchTemplateInput
	input := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateData: launchTemplateData,
	}

	// Set name or name prefix
	if template.NamePrefix != nil && *template.NamePrefix != "" {
		input.LaunchTemplateName = template.NamePrefix
	} else if template.Name != nil && *template.Name != "" {
		input.LaunchTemplateName = template.Name
	} else {
		return nil, fmt.Errorf("launch template name or name_prefix is required")
	}

	// Set update default version flag
	if template.UpdateDefaultVersion != nil {
		input.VersionDescription = aws.String("Default version")
	}

	// Add tags
	if len(template.Tags) > 0 {
		var tagSpecs []types.TagSpecification
		var tags []types.Tag
		for _, tag := range template.Tags {
			tags = append(tags, types.Tag{
				Key:   aws.String(tag.Key),
				Value: aws.String(tag.Value),
			})
		}
		tagSpecs = append(tagSpecs, types.TagSpecification{
			ResourceType: types.ResourceTypeLaunchTemplate,
			Tags:         tags,
		})
		input.TagSpecifications = tagSpecs
	}

	// Create the launch template
	result, err := client.EC2.CreateLaunchTemplate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create launch template: %w", err)
	}

	if result.LaunchTemplate == nil {
		return nil, fmt.Errorf("launch template creation returned nil")
	}

	return convertLaunchTemplateToOutput(result.LaunchTemplate), nil
}

// GetLaunchTemplate retrieves a Launch Template by ID
func GetLaunchTemplate(ctx context.Context, client *AWSClient, id string) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	input := &ec2.DescribeLaunchTemplatesInput{
		LaunchTemplateIds: []string{id},
	}

	output, err := client.EC2.DescribeLaunchTemplates(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe launch template: %w", err)
	}

	if len(output.LaunchTemplates) == 0 {
		return nil, fmt.Errorf("launch template not found: %s", id)
	}

	return convertLaunchTemplateToOutput(&output.LaunchTemplates[0]), nil
}

// UpdateLaunchTemplate creates a new version of an existing Launch Template
func UpdateLaunchTemplate(ctx context.Context, client *AWSClient, id string, template *awslttemplate.LaunchTemplate) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	if err := template.Validate(); err != nil {
		return nil, fmt.Errorf("template validation failed: %w", err)
	}

	// Build LaunchTemplateData
	launchTemplateData := buildLaunchTemplateData(template)

	// Create new version
	input := &ec2.CreateLaunchTemplateVersionInput{
		LaunchTemplateId:   aws.String(id),
		LaunchTemplateData: launchTemplateData,
	}

	// Create the new version
	_, err := client.EC2.CreateLaunchTemplateVersion(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create launch template version: %w", err)
	}

	// Note: Setting default version requires a separate ModifyLaunchTemplate call
	// For now, we'll just create the version
	// If UpdateDefaultVersion is true, the caller can set it as default separately

	// Get the updated template
	return GetLaunchTemplate(ctx, client, id)
}

// DeleteLaunchTemplate deletes a Launch Template
func DeleteLaunchTemplate(ctx context.Context, client *AWSClient, id string) error {
	input := &ec2.DeleteLaunchTemplateInput{
		LaunchTemplateId: aws.String(id),
	}

	_, err := client.EC2.DeleteLaunchTemplate(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete launch template: %w", err)
	}

	return nil
}

// ListLaunchTemplates lists Launch Templates with optional filters
func ListLaunchTemplates(ctx context.Context, client *AWSClient, filters map[string][]string) ([]*awslttemplateoutputs.LaunchTemplateOutput, error) {
	input := &ec2.DescribeLaunchTemplatesInput{}

	// Convert filters to AWS filter format
	if len(filters) > 0 {
		var awsFilters []types.Filter
		for key, values := range filters {
			awsFilters = append(awsFilters, types.Filter{
				Name:   aws.String(key),
				Values: values,
			})
		}
		input.Filters = awsFilters
	}

	output, err := client.EC2.DescribeLaunchTemplates(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe launch templates: %w", err)
	}

	templates := make([]*awslttemplateoutputs.LaunchTemplateOutput, len(output.LaunchTemplates))
	for i, lt := range output.LaunchTemplates {
		templates[i] = convertLaunchTemplateToOutput(&lt)
	}

	return templates, nil
}

// GetLaunchTemplateVersion retrieves a specific version of a Launch Template
func GetLaunchTemplateVersion(ctx context.Context, client *AWSClient, id string, version int) (*awslttemplate.LaunchTemplateVersion, error) {
	input := &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateId: aws.String(id),
		Versions:          []string{fmt.Sprintf("%d", version)},
	}

	output, err := client.EC2.DescribeLaunchTemplateVersions(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe launch template version: %w", err)
	}

	if len(output.LaunchTemplateVersions) == 0 {
		return nil, fmt.Errorf("launch template version not found: %s version %d", id, version)
	}

	return convertLaunchTemplateVersion(&output.LaunchTemplateVersions[0], id), nil
}

// ListLaunchTemplateVersions lists all versions of a Launch Template
func ListLaunchTemplateVersions(ctx context.Context, client *AWSClient, id string) ([]*awslttemplate.LaunchTemplateVersion, error) {
	input := &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateId: aws.String(id),
	}

	output, err := client.EC2.DescribeLaunchTemplateVersions(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe launch template versions: %w", err)
	}

	versions := make([]*awslttemplate.LaunchTemplateVersion, len(output.LaunchTemplateVersions))
	for i, version := range output.LaunchTemplateVersions {
		versions[i] = convertLaunchTemplateVersion(&version, id)
	}

	return versions, nil
}

// buildLaunchTemplateData builds RequestLaunchTemplateData from AWS model
func buildLaunchTemplateData(template *awslttemplate.LaunchTemplate) *types.RequestLaunchTemplateData {
	data := &types.RequestLaunchTemplateData{
		ImageId:      aws.String(template.ImageID),
		InstanceType: types.InstanceType(template.InstanceType),
	}

	// Security groups
	if len(template.VpcSecurityGroupIds) > 0 {
		data.SecurityGroupIds = template.VpcSecurityGroupIds
	}

	// Key name
	if template.KeyName != nil {
		data.KeyName = template.KeyName
	}

	// IAM instance profile
	if template.IAMInstanceProfile != nil {
		data.IamInstanceProfile = &types.LaunchTemplateIamInstanceProfileSpecificationRequest{}
		if template.IAMInstanceProfile.Name != nil {
			data.IamInstanceProfile.Name = template.IAMInstanceProfile.Name
		}
		if template.IAMInstanceProfile.ARN != nil {
			data.IamInstanceProfile.Arn = template.IAMInstanceProfile.ARN
		}
	}

	// User data
	if template.UserData != nil {
		data.UserData = template.UserData
	}

	// Note: Root volume and additional volumes are now referenced by ID
	// They should be created separately using the storage service and attached
	// Block device mappings are not set here as volumes are managed separately
	// If RootVolumeID or AdditionalVolumeIDs are provided, they will be attached
	// during instance launch, not during template creation

	// Metadata options
	if template.MetadataOptions != nil {
		data.MetadataOptions = &types.LaunchTemplateInstanceMetadataOptionsRequest{}

		if template.MetadataOptions.HTTPPutResponseHopLimit != nil {
			hopLimit := int32(*template.MetadataOptions.HTTPPutResponseHopLimit)
			data.MetadataOptions.HttpPutResponseHopLimit = &hopLimit
		}

		if template.MetadataOptions.HTTPEndpoint != nil {
			data.MetadataOptions.HttpEndpoint = types.LaunchTemplateInstanceMetadataEndpointState(*template.MetadataOptions.HTTPEndpoint)
		}

		if template.MetadataOptions.HTTPTokens != nil {
			// HTTPTokens is of type LaunchTemplateHttpTokensState
			data.MetadataOptions.HttpTokens = types.LaunchTemplateHttpTokensState(*template.MetadataOptions.HTTPTokens)
		}
	}

	return data
}

// convertLaunchTemplateToOutput converts AWS SDK LaunchTemplate to output model
func convertLaunchTemplateToOutput(lt *types.LaunchTemplate) *awslttemplateoutputs.LaunchTemplateOutput {
	// Construct ARN manually since LaunchTemplate doesn't have ARN field
	// Format: arn:aws:ec2:region:account-id:launch-template/lt-id
	arn := ""
	if lt.LaunchTemplateId != nil {
		// ARN will be constructed if needed, or can be retrieved from DescribeLaunchTemplates
		// For now, we'll leave it empty or construct it if we have region/account info
		arn = fmt.Sprintf("arn:aws:ec2:region:account:launch-template/%s", *lt.LaunchTemplateId)
	}

	output := &awslttemplateoutputs.LaunchTemplateOutput{
		ID:    aws.ToString(lt.LaunchTemplateId),
		ARN:   arn,
		Name:  aws.ToString(lt.LaunchTemplateName),
		CreateTime: aws.ToTime(lt.CreateTime),
	}

	// Version information
	if lt.DefaultVersionNumber != nil {
		output.DefaultVersion = int(*lt.DefaultVersionNumber)
	}
	if lt.LatestVersionNumber != nil {
		output.LatestVersion = int(*lt.LatestVersionNumber)
	}

	// Created by
	if lt.CreatedBy != nil {
		output.CreatedBy = *lt.CreatedBy
	}

	// Tags - convert to configs.Tag format
	if lt.Tags != nil {
		output.Tags = make([]struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}, len(lt.Tags))
		for i, tag := range lt.Tags {
			output.Tags[i] = struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}{
				Key:   aws.ToString(tag.Key),
				Value: aws.ToString(tag.Value),
			}
		}
	}

	return output
}

// convertLaunchTemplateVersion converts AWS SDK LaunchTemplateVersion to model
func convertLaunchTemplateVersion(ltv *types.LaunchTemplateVersion, templateID string) *awslttemplate.LaunchTemplateVersion {
	version := &awslttemplate.LaunchTemplateVersion{
		TemplateID:    templateID,
		IsDefault:     ltv.DefaultVersion != nil && *ltv.DefaultVersion,
		CreateTime:    aws.ToTime(ltv.CreateTime),
	}

	if ltv.VersionNumber != nil {
		version.VersionNumber = int(*ltv.VersionNumber)
	}

	if ltv.CreatedBy != nil {
		version.CreatedBy = ltv.CreatedBy
	}

	// Note: TemplateData would need to be converted from LaunchTemplateData
	// This is complex and may be done lazily when needed
	// For now, we'll leave it nil and can populate it on demand

	return version
}
