package sdk

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// CreateInstanceProfile creates a new IAM instance profile using AWS SDK
func CreateInstanceProfile(ctx context.Context, client *AWSClient, profile *awsiam.InstanceProfile) (*awsoutputs.InstanceProfileOutput, error) {
	if err := profile.Validate(); err != nil {
		return nil, fmt.Errorf("instance profile validation failed: %w", err)
	}

	// Build CreateInstanceProfileInput
	input := &iam.CreateInstanceProfileInput{}

	// Set name or name prefix
	if profile.Name != "" {
		input.InstanceProfileName = aws.String(profile.Name)
	} else if profile.NamePrefix != nil && *profile.NamePrefix != "" {
		input.InstanceProfileName = aws.String(*profile.NamePrefix)
	}

	// Set path (default to "/" if not provided)
	path := "/"
	if profile.Path != nil && *profile.Path != "" {
		path = *profile.Path
	}
	input.Path = aws.String(path)

	// Add tags if provided
	if len(profile.Tags) > 0 {
		var tagList []types.Tag
		for _, tag := range profile.Tags {
			tagList = append(tagList, types.Tag{
				Key:   aws.String(tag.Key),
				Value: aws.String(tag.Value),
			})
		}
		input.Tags = tagList
	}

	// Create the instance profile
	result, err := client.IAM.CreateInstanceProfile(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance profile: %w", err)
	}

	if result.InstanceProfile == nil {
		return nil, fmt.Errorf("instance profile creation returned nil")
	}

	output := convertInstanceProfileToOutput(result.InstanceProfile)

	// If role was provided, add it to the profile
	if profile.Role != nil && *profile.Role != "" {
		if err := AddRoleToInstanceProfile(ctx, client, output.Name, *profile.Role); err != nil {
			// Log warning but don't fail - profile was created successfully
			fmt.Printf("Warning: Failed to add role %s to instance profile %s: %v\n", *profile.Role, output.Name, err)
		} else {
			// Refresh the profile to get the role information
			updatedProfile, err := GetInstanceProfile(ctx, client, output.Name)
			if err == nil {
				output = updatedProfile
			}
		}
	}

	return output, nil
}

// GetInstanceProfile retrieves an IAM instance profile by name
func GetInstanceProfile(ctx context.Context, client *AWSClient, name string) (*awsoutputs.InstanceProfileOutput, error) {
	if client == nil || client.IAM == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	input := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	}

	result, err := client.IAM.GetInstanceProfile(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance profile %s: %w", name, err)
	}

	if result.InstanceProfile == nil {
		return nil, fmt.Errorf("instance profile %s not found", name)
	}

	return convertInstanceProfileToOutput(result.InstanceProfile), nil
}

// UpdateInstanceProfile updates an IAM instance profile
// Note: AWS only allows updating path and tags, not name or roles
func UpdateInstanceProfile(ctx context.Context, client *AWSClient, name string, profile *awsiam.InstanceProfile) (*awsoutputs.InstanceProfileOutput, error) {
	if err := profile.Validate(); err != nil {
		return nil, fmt.Errorf("instance profile validation failed: %w", err)
	}

	// Update path if provided
	if profile.Path != nil && *profile.Path != "" {
		// Note: AWS doesn't have a direct UpdateInstanceProfile API
		// Path cannot be changed after creation, so we'll just update tags
	}

	// Update tags if provided
	if len(profile.Tags) > 0 {
		// Tag instance profile
		var tagList []types.Tag
		for _, tag := range profile.Tags {
			tagList = append(tagList, types.Tag{
				Key:   aws.String(tag.Key),
				Value: aws.String(tag.Value),
			})
		}

		tagInput := &iam.TagInstanceProfileInput{
			InstanceProfileName: aws.String(name),
			Tags:                tagList,
		}

		_, err := client.IAM.TagInstanceProfile(ctx, tagInput)
		if err != nil {
			return nil, fmt.Errorf("failed to tag instance profile: %w", err)
		}
	}

	// Return updated profile
	return GetInstanceProfile(ctx, client, name)
}

// DeleteInstanceProfile deletes an IAM instance profile
// Note: All roles must be removed from the profile before deletion
func DeleteInstanceProfile(ctx context.Context, client *AWSClient, name string) error {
	if client == nil || client.IAM == nil {
		return fmt.Errorf("AWS client not available")
	}

	// First, get the profile to check for attached roles
	profile, err := GetInstanceProfile(ctx, client, name)
	if err != nil {
		return fmt.Errorf("failed to get instance profile: %w", err)
	}

	// Remove all roles from the profile
	for _, role := range profile.Roles {
		if err := RemoveRoleFromInstanceProfile(ctx, client, name, role.Name); err != nil {
			return fmt.Errorf("failed to remove role %s from instance profile: %w", role.Name, err)
		}
	}

	// Now delete the profile
	deleteInput := &iam.DeleteInstanceProfileInput{
		InstanceProfileName: aws.String(name),
	}

	_, err = client.IAM.DeleteInstanceProfile(ctx, deleteInput)
	if err != nil {
		return fmt.Errorf("failed to delete instance profile %s: %w", name, err)
	}

	return nil
}

// ListInstanceProfiles lists IAM instance profiles with optional path prefix filter
func ListInstanceProfiles(ctx context.Context, client *AWSClient, pathPrefix *string) ([]*awsoutputs.InstanceProfileOutput, error) {
	if client == nil || client.IAM == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	var allProfiles []*awsoutputs.InstanceProfileOutput
	var nextToken *string

	for {
		input := &iam.ListInstanceProfilesInput{
			PathPrefix: pathPrefix,
		}

		if nextToken != nil {
			input.Marker = nextToken
		}

		result, err := client.IAM.ListInstanceProfiles(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list instance profiles: %w", err)
		}

		// Convert each profile to output model
		for _, profile := range result.InstanceProfiles {
			allProfiles = append(allProfiles, convertInstanceProfileToOutput(&profile))
		}

		// Check if there are more pages
		if !result.IsTruncated || result.Marker == nil {
			break
		}
		nextToken = result.Marker
	}

	return allProfiles, nil
}

// AddRoleToInstanceProfile adds a role to an instance profile
func AddRoleToInstanceProfile(ctx context.Context, client *AWSClient, profileName, roleName string) error {
	if client == nil || client.IAM == nil {
		return fmt.Errorf("AWS client not available")
	}

	input := &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
		RoleName:            aws.String(roleName),
	}

	_, err := client.IAM.AddRoleToInstanceProfile(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to add role %s to instance profile %s: %w", roleName, profileName, err)
	}

	return nil
}

// RemoveRoleFromInstanceProfile removes a role from an instance profile
func RemoveRoleFromInstanceProfile(ctx context.Context, client *AWSClient, profileName, roleName string) error {
	if client == nil || client.IAM == nil {
		return fmt.Errorf("AWS client not available")
	}

	input := &iam.RemoveRoleFromInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
		RoleName:            aws.String(roleName),
	}

	_, err := client.IAM.RemoveRoleFromInstanceProfile(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to remove role %s from instance profile %s: %w", roleName, profileName, err)
	}

	return nil
}

// GetInstanceProfileRoles retrieves all roles attached to an instance profile
func GetInstanceProfileRoles(ctx context.Context, client *AWSClient, profileName string) ([]*awsoutputs.RoleOutput, error) {
	profile, err := GetInstanceProfile(ctx, client, profileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance profile: %w", err)
	}

	return profile.Roles, nil
}

// convertInstanceProfileToOutput converts AWS SDK InstanceProfile to output model
func convertInstanceProfileToOutput(profile *types.InstanceProfile) *awsoutputs.InstanceProfileOutput {
	output := &awsoutputs.InstanceProfileOutput{
		ARN:        aws.ToString(profile.Arn),
		ID:         aws.ToString(profile.InstanceProfileName),
		Name:       aws.ToString(profile.InstanceProfileName),
		Path:       aws.ToString(profile.Path),
		CreateDate: aws.ToTime(profile.CreateDate),
	}

	// Convert roles
	if len(profile.Roles) > 0 {
		output.Roles = make([]*awsoutputs.RoleOutput, len(profile.Roles))
		for i, role := range profile.Roles {
			roleOutput := &awsoutputs.RoleOutput{
				ARN:        aws.ToString(role.Arn),
				ID:         aws.ToString(role.RoleName),
				Name:       aws.ToString(role.RoleName),
				UniqueID:   aws.ToString(role.RoleId),
				Path:       aws.ToString(role.Path),
				CreateDate: aws.ToTime(role.CreateDate),
			}
			if role.Description != nil {
				roleOutput.Description = role.Description
			}
			if role.PermissionsBoundary != nil {
				roleOutput.PermissionsBoundary = role.PermissionsBoundary.PermissionsBoundaryArn
			}
			output.Roles[i] = roleOutput
		}
	}

	// Convert tags
	if len(profile.Tags) > 0 {
		output.Tags = make([]configs.Tag, len(profile.Tags))
		for i, tag := range profile.Tags {
			output.Tags[i] = configs.Tag{
				Key:   aws.ToString(tag.Key),
				Value: aws.ToString(tag.Value),
			}
		}
	}

	return output
}
