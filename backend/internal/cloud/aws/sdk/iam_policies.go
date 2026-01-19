package sdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
)

// CreatePolicy creates a new IAM policy using AWS SDK
func CreatePolicy(ctx context.Context, client *AWSClient, policy *awsiam.Policy) (*awsoutputs.PolicyOutput, error) {
	if err := policy.Validate(); err != nil {
		return nil, fmt.Errorf("policy validation failed: %w", err)
	}

	// Convert policy document from JSON string to bytes
	policyDocument := []byte(policy.PolicyDocument)

	// Build CreatePolicyInput
	input := &iam.CreatePolicyInput{
		PolicyName:     aws.String(policy.Name),
		PolicyDocument: aws.String(string(policyDocument)),
	}

	// Set description if provided
	if policy.Description != nil {
		input.Description = policy.Description
	}

	// Set path (default to "/" if not provided)
	path := "/"
	if policy.Path != nil && *policy.Path != "" {
		path = *policy.Path
	}
	input.Path = aws.String(path)

	// Add tags if provided
	if len(policy.Tags) > 0 {
		var tagList []types.Tag
		for _, tag := range policy.Tags {
			tagList = append(tagList, types.Tag{
				Key:   aws.String(tag.Key),
				Value: aws.String(tag.Value),
			})
		}
		input.Tags = tagList
	}

	// Create the policy
	result, err := client.IAM.CreatePolicy(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create policy: %w", err)
	}

	if result.Policy == nil {
		return nil, fmt.Errorf("policy creation returned nil")
	}

	// Get the default policy version to retrieve the policy document
	policyOutput := convertPolicyToOutput(result.Policy, false)

	// Fetch the policy document from the default version
	if result.Policy.DefaultVersionId != nil {
		versionOutput, err := GetPolicyVersion(ctx, client, aws.ToString(result.Policy.Arn), aws.ToString(result.Policy.DefaultVersionId))
		if err == nil && versionOutput != nil {
			policyOutput.PolicyDocument = versionOutput.PolicyDocument
			policyOutput.DefaultVersionID = result.Policy.DefaultVersionId
		}
	}

	return policyOutput, nil
}

// GetPolicy retrieves an IAM policy by ARN, including its policy document
// Checks PolicyService static data first, then falls back to AWS SDK
func GetPolicy(ctx context.Context, client *AWSClient, arn string) (*awsoutputs.PolicyOutput, error) {
	// Check if we have a policy service with static data
	// Try to get from static data first (fast, no API call)
	if client != nil {
		// Create a temporary policy service to check static data
		// In production, this would be initialized once and reused
		policyService, err := NewPolicyService(client)
		if err == nil {
			if staticPolicy, err := policyService.GetPolicy(ctx, arn); err == nil {
				// Found in static data, but try to enhance with SDK data if available
				if client.IAM != nil {
					// Try to get from SDK to get latest metadata (attachment count, etc.)
					sdkPolicy, err := getPolicyFromSDK(ctx, client, arn)
					if err == nil && sdkPolicy != nil {
						// Merge: use SDK data but keep static policy document if SDK fails
						if sdkPolicy.PolicyDocument != "" {
							staticPolicy.PolicyDocument = sdkPolicy.PolicyDocument
						}
						staticPolicy.AttachmentCount = sdkPolicy.AttachmentCount
						staticPolicy.CreateDate = sdkPolicy.CreateDate
						staticPolicy.UpdateDate = sdkPolicy.UpdateDate
						if sdkPolicy.DefaultVersionID != nil {
							staticPolicy.DefaultVersionID = sdkPolicy.DefaultVersionID
						}
						return staticPolicy, nil
					}
				}
				// Return static data if SDK is not available
				return staticPolicy, nil
			}
		}
	}

	// Fall back to SDK
	return getPolicyFromSDK(ctx, client, arn)
}

// getPolicyFromSDK retrieves policy directly from AWS SDK
func getPolicyFromSDK(ctx context.Context, client *AWSClient, arn string) (*awsoutputs.PolicyOutput, error) {
	if client == nil || client.IAM == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	input := &iam.GetPolicyInput{
		PolicyArn: aws.String(arn),
	}

	result, err := client.IAM.GetPolicy(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy %s: %w", arn, err)
	}

	if result.Policy == nil {
		return nil, fmt.Errorf("get policy %s returned nil", arn)
	}

	isAWSManaged := strings.HasPrefix(arn, "arn:aws:iam::aws:policy/")
	output := convertPolicyToOutput(result.Policy, isAWSManaged)

	// Get the policy document for the default version
	policyVersionOutput, err := GetPolicyVersion(ctx, client, arn, aws.ToString(result.Policy.DefaultVersionId))
	if err != nil {
		return nil, fmt.Errorf("failed to get policy version for %s: %w", arn, err)
	}
	output.PolicyDocument = policyVersionOutput.PolicyDocument

	return output, nil
}

// UpdatePolicy updates a policy by creating a new version
func UpdatePolicy(ctx context.Context, client *AWSClient, arn string, policy *awsiam.Policy) (*awsoutputs.PolicyOutput, error) {
	if err := policy.Validate(); err != nil {
		return nil, fmt.Errorf("policy validation failed: %w", err)
	}

	// Check if this is an AWS managed policy (cannot be updated)
	if strings.HasPrefix(arn, "arn:aws:iam::aws:policy/") {
		return nil, fmt.Errorf("cannot update AWS managed policy")
	}

	// Convert policy document from JSON string to bytes
	policyDocument := []byte(policy.PolicyDocument)

	// Create a new policy version
	createVersionInput := &iam.CreatePolicyVersionInput{
		PolicyArn:      aws.String(arn),
		PolicyDocument: aws.String(string(policyDocument)),
		SetAsDefault:   true, // Set as default version
	}

	_, err := client.IAM.CreatePolicyVersion(ctx, createVersionInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create policy version: %w", err)
	}

	// Get the updated policy
	return GetPolicy(ctx, client, arn)
}

// DeletePolicy deletes an IAM policy
func DeletePolicy(ctx context.Context, client *AWSClient, arn string) error {
	// Check if this is an AWS managed policy (cannot be deleted)
	if strings.HasPrefix(arn, "arn:aws:iam::aws:policy/") {
		return fmt.Errorf("cannot delete AWS managed policy")
	}

	input := &iam.DeletePolicyInput{
		PolicyArn: aws.String(arn),
	}

	_, err := client.IAM.DeletePolicy(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}

	return nil
}

// ListPolicies lists IAM policies with optional path prefix filter
// If includePolicyDocument is false, policy documents are not fetched (faster for listing)
func ListPolicies(ctx context.Context, client *AWSClient, pathPrefix *string, scope types.PolicyScopeType, includePolicyDocument bool) ([]*awsoutputs.PolicyOutput, error) {
	var allPolicies []*awsoutputs.PolicyOutput
	var nextToken *string

	for {
		input := &iam.ListPoliciesInput{
			Scope: scope,
		}

		if pathPrefix != nil && *pathPrefix != "" {
			input.PathPrefix = pathPrefix
		}

		if nextToken != nil {
			input.Marker = nextToken
		}

		result, err := client.IAM.ListPolicies(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list policies: %w", err)
		}

		// Convert each policy to output model
		for _, policy := range result.Policies {
			isAWSManaged := scope == types.PolicyScopeTypeAws || strings.HasPrefix(aws.ToString(policy.Arn), "arn:aws:iam::aws:policy/")
			policyOutput := convertPolicyToOutput(&policy, isAWSManaged)

			// Only fetch policy document if explicitly requested (skip for listing to improve performance)
			if includePolicyDocument && policy.DefaultVersionId != nil {
				versionOutput, err := GetPolicyVersion(ctx, client, aws.ToString(policy.Arn), aws.ToString(policy.DefaultVersionId))
				if err == nil && versionOutput != nil {
					policyOutput.PolicyDocument = versionOutput.PolicyDocument
					policyOutput.DefaultVersionID = policy.DefaultVersionId
				}
			} else if policy.DefaultVersionId != nil {
				// Store version ID even if we don't fetch the document
				policyOutput.DefaultVersionID = policy.DefaultVersionId
			}

			allPolicies = append(allPolicies, policyOutput)
		}

		// Check if there are more pages
		if !result.IsTruncated || result.Marker == nil {
			break
		}
		nextToken = result.Marker
	}

	return allPolicies, nil
}

// ListAWSManagedPolicies lists AWS managed policies
// By default, policy documents are not fetched for performance (use GetPolicy for full details)
// Merges static data with SDK results
func ListAWSManagedPolicies(ctx context.Context, client *AWSClient, pathPrefix *string) ([]*awsoutputs.PolicyOutput, error) {
	var results []*awsoutputs.PolicyOutput
	seenARNs := make(map[string]bool)

	// First, get policies from static data (fast, no API call)
	if client != nil {
		policyService, err := NewPolicyService(client)
		if err == nil {
			staticPolicies, err := policyService.ListAllStaticPolicies(ctx)
			if err == nil {
				for _, policy := range staticPolicies {
					// Filter by path prefix if provided
					if pathPrefix == nil || *pathPrefix == "" || strings.HasPrefix(policy.Path, *pathPrefix) {
						// Only include AWS managed policies
						if policy.IsAWSManaged {
							results = append(results, policy)
							seenARNs[strings.ToLower(policy.ARN)] = true
						}
					}
				}
			}
		}
	}

	// Then, get from SDK and merge (avoid duplicates)
	if client != nil && client.IAM != nil {
		sdkPolicies, err := ListPolicies(ctx, client, pathPrefix, types.PolicyScopeTypeAws, false)
		if err == nil {
			for _, policy := range sdkPolicies {
				arnLower := strings.ToLower(policy.ARN)
				if !seenARNs[arnLower] {
					results = append(results, policy)
					seenARNs[arnLower] = true
				} else {
					// Update existing entry with SDK metadata (attachment count, dates, etc.)
					for i, existing := range results {
						if strings.EqualFold(existing.ARN, policy.ARN) {
							results[i].AttachmentCount = policy.AttachmentCount
							results[i].CreateDate = policy.CreateDate
							results[i].UpdateDate = policy.UpdateDate
							if policy.DefaultVersionID != nil {
								results[i].DefaultVersionID = policy.DefaultVersionID
							}
							break
						}
					}
				}
			}
		}
	}

	return results, nil
}

// GetAWSManagedPolicy retrieves an AWS managed policy by ARN
// Checks static data first, then falls back to SDK
func GetAWSManagedPolicy(ctx context.Context, client *AWSClient, arn string) (*awsoutputs.PolicyOutput, error) {
	// Verify it's an AWS managed policy
	if !strings.HasPrefix(arn, "arn:aws:iam::aws:policy/") {
		return nil, fmt.Errorf("not an AWS managed policy ARN")
	}

	// Try static data first (fast)
	if client != nil {
		policyService, err := NewPolicyService(client)
		if err == nil {
			if staticPolicy, err := policyService.GetPolicy(ctx, arn); err == nil {
				// Found in static data, enhance with SDK if available
				if client.IAM != nil {
					sdkPolicy, err := getPolicyFromSDK(ctx, client, arn)
					if err == nil && sdkPolicy != nil {
						// Merge SDK metadata
						if sdkPolicy.PolicyDocument != "" {
							staticPolicy.PolicyDocument = sdkPolicy.PolicyDocument
						}
						staticPolicy.AttachmentCount = sdkPolicy.AttachmentCount
						staticPolicy.CreateDate = sdkPolicy.CreateDate
						staticPolicy.UpdateDate = sdkPolicy.UpdateDate
						if sdkPolicy.DefaultVersionID != nil {
							staticPolicy.DefaultVersionID = sdkPolicy.DefaultVersionID
						}
						return staticPolicy, nil
					}
				}
				return staticPolicy, nil
			}
		}
	}

	// Fall back to SDK
	return getPolicyFromSDK(ctx, client, arn)
}

// GetPolicyVersion retrieves a specific policy version document
func GetPolicyVersion(ctx context.Context, client *AWSClient, policyARN, versionID string) (*awsoutputs.PolicyOutput, error) {
	input := &iam.GetPolicyVersionInput{
		PolicyArn: aws.String(policyARN),
		VersionId: aws.String(versionID),
	}

	result, err := client.IAM.GetPolicyVersion(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy version: %w", err)
	}

	if result.PolicyVersion == nil {
		return nil, fmt.Errorf("policy version not found")
	}

	// Create a minimal output with just the policy document
	output := &awsoutputs.PolicyOutput{
		PolicyDocument: aws.ToString(result.PolicyVersion.Document),
	}

	return output, nil
}

// convertPolicyToOutput converts AWS SDK Policy to output model
func convertPolicyToOutput(policy *types.Policy, isAWSManaged bool) *awsoutputs.PolicyOutput {
	output := &awsoutputs.PolicyOutput{
		ARN:             aws.ToString(policy.Arn),
		ID:              aws.ToString(policy.Arn),
		Name:            aws.ToString(policy.PolicyName),
		Description:     policy.Description,
		Path:            aws.ToString(policy.Path),
		CreateDate:      aws.ToTime(policy.CreateDate),
		UpdateDate:      aws.ToTime(policy.UpdateDate),
		AttachmentCount: int(aws.ToInt32(policy.AttachmentCount)),
		IsAttachable:    policy.IsAttachable,
		IsAWSManaged:    isAWSManaged,
	}

	// Convert tags
	if len(policy.Tags) > 0 {
		output.Tags = make([]configs.Tag, len(policy.Tags))
		for i, tag := range policy.Tags {
			output.Tags[i] = configs.Tag{
				Key:   aws.ToString(tag.Key),
				Value: aws.ToString(tag.Value),
			}
		}
	}

	return output
}
