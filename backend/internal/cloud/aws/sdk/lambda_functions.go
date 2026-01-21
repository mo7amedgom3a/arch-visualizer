package sdk

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	awslambda "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda"
	awslambdaoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda/outputs"
)

// CreateLambdaFunction creates a new Lambda function using AWS SDK
func CreateLambdaFunction(ctx context.Context, client *AWSClient, function *awslambda.Function) (*awslambdaoutputs.FunctionOutput, error) {
	if err := function.Validate(); err != nil {
		return nil, fmt.Errorf("lambda function validation failed: %w", err)
	}

	if client == nil || client.Lambda == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	// Build CreateFunctionInput
	input := &lambda.CreateFunctionInput{
		FunctionName: aws.String(function.FunctionName),
		Role:         aws.String(function.RoleARN),
	}

	// Code Source - S3
	if function.S3Bucket != nil && *function.S3Bucket != "" && function.S3Key != nil && *function.S3Key != "" {
		input.Code = &types.FunctionCode{
			S3Bucket:        function.S3Bucket,
			S3Key:           function.S3Key,
			S3ObjectVersion: function.S3ObjectVersion,
		}
		if function.Runtime != nil {
			input.Runtime = types.Runtime(*function.Runtime)
		}
		if function.Handler != nil {
			input.Handler = function.Handler
		}
	}

	// Code Source - Container Image
	if function.PackageType != nil && *function.PackageType == "Image" && function.ImageURI != nil && *function.ImageURI != "" {
		input.PackageType = types.PackageTypeImage
		input.Code = &types.FunctionCode{
			ImageUri: function.ImageURI,
		}
	}

	// Configuration
	if function.MemorySize != nil {
		input.MemorySize = function.MemorySize
	}
	if function.Timeout != nil {
		input.Timeout = function.Timeout
	}

	// Environment variables
	if function.Environment != nil && len(function.Environment) > 0 {
		envVars := make(map[string]string)
		for k, v := range function.Environment {
			envVars[k] = v
		}
		input.Environment = &types.Environment{
			Variables: envVars,
		}
	}

	// Layers
	if function.Layers != nil && len(function.Layers) > 0 {
		input.Layers = function.Layers
	}

	// VPC Config
	if function.VPCConfig != nil {
		input.VpcConfig = &types.VpcConfig{
			SubnetIds:        function.VPCConfig.SubnetIDs,
			SecurityGroupIds: function.VPCConfig.SecurityGroupIDs,
		}
	}

	// Tags
	if len(function.Tags) > 0 {
		tags := make(map[string]string)
		for _, tag := range function.Tags {
			tags[tag.Key] = tag.Value
		}
		input.Tags = tags
	}

	// Create the function
	result, err := client.Lambda.CreateFunction(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create lambda function: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("lambda function creation returned nil")
	}

	// CreateFunctionOutput embeds FunctionConfiguration fields
	// Convert to FunctionConfiguration for consistent handling
	config := &types.FunctionConfiguration{
		FunctionName: result.FunctionName,
		FunctionArn:  result.FunctionArn,
		Runtime:      result.Runtime,
		Role:         result.Role,
		Handler:      result.Handler,
		CodeSize:     result.CodeSize,
		CodeSha256:   result.CodeSha256,
		Description:  result.Description,
		Timeout:      result.Timeout,
		MemorySize:   result.MemorySize,
		LastModified: result.LastModified,
		Version:      result.Version,
		VpcConfig:    result.VpcConfig,
		Environment:  result.Environment,
		Layers:       result.Layers,
		PackageType:  result.PackageType,
	}
	return convertLambdaFunctionToOutput(config), nil
}

// GetLambdaFunction retrieves a Lambda function by name
func GetLambdaFunction(ctx context.Context, client *AWSClient, name string) (*awslambdaoutputs.FunctionOutput, error) {
	if client == nil || client.Lambda == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	input := &lambda.GetFunctionInput{
		FunctionName: aws.String(name),
	}

	result, err := client.Lambda.GetFunction(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get lambda function %s: %w", name, err)
	}

	if result == nil || result.Configuration == nil {
		return nil, fmt.Errorf("lambda function %s not found", name)
	}

	return convertLambdaFunctionToOutput(result.Configuration), nil
}

// UpdateLambdaFunction updates a Lambda function configuration
func UpdateLambdaFunction(ctx context.Context, client *AWSClient, name string, function *awslambda.Function) (*awslambdaoutputs.FunctionOutput, error) {
	if err := function.Validate(); err != nil {
		return nil, fmt.Errorf("lambda function validation failed: %w", err)
	}

	if client == nil || client.Lambda == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	// Build UpdateFunctionConfigurationInput
	configInput := &lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(name),
	}

	// Update role if changed
	if function.RoleARN != "" {
		configInput.Role = aws.String(function.RoleARN)
	}

	// Update memory size
	if function.MemorySize != nil {
		configInput.MemorySize = function.MemorySize
	}

	// Update timeout
	if function.Timeout != nil {
		configInput.Timeout = function.Timeout
	}

	// Update environment variables
	if function.Environment != nil {
		envVars := make(map[string]string)
		for k, v := range function.Environment {
			envVars[k] = v
		}
		configInput.Environment = &types.Environment{
			Variables: envVars,
		}
	}

	// Update layers
	if function.Layers != nil {
		configInput.Layers = function.Layers
	}

	// Update VPC config
	if function.VPCConfig != nil {
		configInput.VpcConfig = &types.VpcConfig{
			SubnetIds:        function.VPCConfig.SubnetIDs,
			SecurityGroupIds: function.VPCConfig.SecurityGroupIDs,
		}
	}

	// Update function configuration
	_, err := client.Lambda.UpdateFunctionConfiguration(ctx, configInput)
	if err != nil {
		return nil, fmt.Errorf("failed to update lambda function configuration: %w", err)
	}

	// If code is being updated, update function code separately
	if (function.S3Bucket != nil && *function.S3Bucket != "" && function.S3Key != nil && *function.S3Key != "") ||
		(function.PackageType != nil && *function.PackageType == "Image" && function.ImageURI != nil && *function.ImageURI != "") {
		codeInput := &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(name),
		}

		// S3 code update
		if function.S3Bucket != nil && *function.S3Bucket != "" && function.S3Key != nil && *function.S3Key != "" {
			codeInput.S3Bucket = function.S3Bucket
			codeInput.S3Key = function.S3Key
			codeInput.S3ObjectVersion = function.S3ObjectVersion
		}

		// Container image update
		if function.PackageType != nil && *function.PackageType == "Image" && function.ImageURI != nil && *function.ImageURI != "" {
			codeInput.ImageUri = function.ImageURI
		}

		_, err = client.Lambda.UpdateFunctionCode(ctx, codeInput)
		if err != nil {
			return nil, fmt.Errorf("failed to update lambda function code: %w", err)
		}
	}

	// Get updated function
	return GetLambdaFunction(ctx, client, name)
}

// DeleteLambdaFunction deletes a Lambda function
func DeleteLambdaFunction(ctx context.Context, client *AWSClient, name string) error {
	if client == nil || client.Lambda == nil {
		return fmt.Errorf("AWS client not available")
	}

	input := &lambda.DeleteFunctionInput{
		FunctionName: aws.String(name),
	}

	_, err := client.Lambda.DeleteFunction(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete lambda function %s: %w", name, err)
	}

	return nil
}

// ListLambdaFunctions lists Lambda functions with optional filters
func ListLambdaFunctions(ctx context.Context, client *AWSClient, filters map[string][]string) ([]*awslambdaoutputs.FunctionOutput, error) {
	if client == nil || client.Lambda == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	input := &lambda.ListFunctionsInput{}

	// Apply filters if provided
	if filters != nil {
		if functionNames, ok := filters["function_name"]; ok && len(functionNames) > 0 {
			// Note: AWS Lambda ListFunctions doesn't support filtering by name directly
			// We'll filter after fetching
		}
	}

	var allFunctions []*awslambdaoutputs.FunctionOutput
	var marker *string

	for {
		if marker != nil {
			input.Marker = marker
		}

		result, err := client.Lambda.ListFunctions(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list lambda functions: %w", err)
		}

		// Convert each function to output
		for _, fn := range result.Functions {
			allFunctions = append(allFunctions, convertLambdaFunctionToOutput(&fn))
		}

		// Check if there are more pages
		if result.NextMarker == nil {
			break
		}
		marker = result.NextMarker
	}

	// Apply filters after fetching
	if filters != nil {
		if functionNames, ok := filters["function_name"]; ok && len(functionNames) > 0 {
			filtered := make([]*awslambdaoutputs.FunctionOutput, 0)
			for _, fn := range allFunctions {
				for _, name := range functionNames {
					if fn.FunctionName == name {
						filtered = append(filtered, fn)
						break
					}
				}
			}
			allFunctions = filtered
		}
	}

	return allFunctions, nil
}

// convertLambdaFunctionToOutput converts AWS Lambda FunctionConfiguration to output model
func convertLambdaFunctionToOutput(config *types.FunctionConfiguration) *awslambdaoutputs.FunctionOutput {
	if config == nil {
		return nil
	}

	output := &awslambdaoutputs.FunctionOutput{
		ARN:          aws.ToString(config.FunctionArn),
		FunctionName: aws.ToString(config.FunctionName),
		RoleARN:      aws.ToString(config.Role),
		Version:      aws.ToString(config.Version),
	}

	// Build invoke ARN: arn:aws:apigateway:region:lambda:path/2015-03-31/functions/function-arn/invocations
	if config.FunctionArn != nil {
		invokeARN := *config.FunctionArn + ":$LATEST"
		output.InvokeARN = invokeARN
	}

	// Qualified ARN (with version)
	if config.FunctionArn != nil && config.Version != nil {
		qualifiedARN := *config.FunctionArn + ":" + *config.Version
		output.QualifiedARN = &qualifiedARN
	}

	// Runtime configuration
	if config.Runtime != "" {
		runtime := string(config.Runtime)
		output.Runtime = &runtime
	}
	if config.Handler != nil {
		output.Handler = config.Handler
	}

	// Configuration
	if config.MemorySize != nil {
		output.MemorySize = config.MemorySize
	}
	if config.Timeout != nil {
		output.Timeout = config.Timeout
	}

	// Environment variables
	if config.Environment != nil && config.Environment.Variables != nil {
		output.Environment = config.Environment.Variables
	}

	// Layers
	if config.Layers != nil && len(config.Layers) > 0 {
		output.Layers = make([]string, len(config.Layers))
		for i, layer := range config.Layers {
			output.Layers[i] = aws.ToString(layer.Arn)
		}
	}

	// VPC Config
	if config.VpcConfig != nil {
		output.VPCConfig = &awslambdaoutputs.FunctionVPCConfigOutput{
			SubnetIDs:        config.VpcConfig.SubnetIds,
			SecurityGroupIDs: config.VpcConfig.SecurityGroupIds,
		}
		if config.VpcConfig.VpcId != nil {
			output.VPCConfig.VPCID = config.VpcConfig.VpcId
		}
	}

	// Metadata
	if config.LastModified != nil {
		output.LastModified = config.LastModified
	}
	if config.CodeSize != 0 {
		codeSize := config.CodeSize
		output.CodeSize = &codeSize
	}
	if config.CodeSha256 != nil {
		output.CodeSHA256 = config.CodeSha256
	}
	if config.Description != nil {
		output.Description = config.Description
	}

	// Package type (for container images)
	if config.PackageType != "" {
		packageType := string(config.PackageType)
		output.PackageType = &packageType
	}
	// Note: Image URI is typically retrieved from GetFunctionCode call
	// For now, we'll leave it empty in the output model

	return output
}
