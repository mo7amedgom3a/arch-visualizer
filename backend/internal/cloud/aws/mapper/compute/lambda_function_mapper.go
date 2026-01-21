package compute

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awslambda "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda"
	awslambdaoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda/outputs"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
)

// FromDomainLambdaFunction converts domain LambdaFunction to AWS Function input model
func FromDomainLambdaFunction(domain *domaincompute.LambdaFunction) *awslambda.Function {
	if domain == nil {
		return nil
	}

	awsFunction := &awslambda.Function{
		FunctionName: domain.FunctionName,
		RoleARN:      domain.RoleARN,
	}

	// Code Source - S3
	if domain.HasS3Code() {
		awsFunction.S3Bucket = domain.S3Bucket
		awsFunction.S3Key = domain.S3Key
		awsFunction.S3ObjectVersion = domain.S3ObjectVersion
		awsFunction.Runtime = domain.Runtime
		awsFunction.Handler = domain.Handler
	}

	// Code Source - Container Image
	if domain.HasContainerImage() {
		awsFunction.PackageType = domain.PackageType
		awsFunction.ImageURI = domain.ImageURI
	}

	// Configuration
	if domain.MemorySize != nil {
		memorySize := int32(*domain.MemorySize)
		awsFunction.MemorySize = &memorySize
	}
	if domain.Timeout != nil {
		timeout := int32(*domain.Timeout)
		awsFunction.Timeout = &timeout
	}

	// Environment variables
	if domain.Environment != nil && len(domain.Environment) > 0 {
		awsFunction.Environment = make(map[string]string)
		for k, v := range domain.Environment {
			awsFunction.Environment[k] = v
		}
	}

	// Layers
	if domain.Layers != nil && len(domain.Layers) > 0 {
		awsFunction.Layers = make([]string, len(domain.Layers))
		copy(awsFunction.Layers, domain.Layers)
	}

	// VPC Config
	if domain.VPCConfig != nil {
		awsFunction.VPCConfig = &awslambda.FunctionVPCConfig{
			SubnetIDs:        make([]string, len(domain.VPCConfig.SubnetIDs)),
			SecurityGroupIDs: make([]string, len(domain.VPCConfig.SecurityGroupIDs)),
		}
		copy(awsFunction.VPCConfig.SubnetIDs, domain.VPCConfig.SubnetIDs)
		copy(awsFunction.VPCConfig.SecurityGroupIDs, domain.VPCConfig.SecurityGroupIDs)
	}

	// Tags - convert from map[string]string to []configs.Tag
	if domain.Tags != nil && len(domain.Tags) > 0 {
		awsFunction.Tags = make([]configs.Tag, 0, len(domain.Tags))
		for key, value := range domain.Tags {
			awsFunction.Tags = append(awsFunction.Tags, configs.Tag{
				Key:   key,
				Value: value,
			})
		}
	}

	// Add Name tag if function name is set
	if domain.FunctionName != "" {
		hasNameTag := false
		for _, tag := range awsFunction.Tags {
			if tag.Key == "Name" {
				hasNameTag = true
				break
			}
		}
		if !hasNameTag {
			awsFunction.Tags = append(awsFunction.Tags, configs.Tag{
				Key:   "Name",
				Value: domain.FunctionName,
			})
		}
	}

	return awsFunction
}

// ToDomainLambdaFunction converts AWS Function input model to domain LambdaFunction
func ToDomainLambdaFunction(aws *awslambda.Function) *domaincompute.LambdaFunction {
	if aws == nil {
		return nil
	}

	domain := &domaincompute.LambdaFunction{
		FunctionName: aws.FunctionName,
		RoleARN:      aws.RoleARN,
	}

	// Code Source - S3
	if aws.S3Bucket != nil && *aws.S3Bucket != "" && aws.S3Key != nil && *aws.S3Key != "" {
		domain.S3Bucket = aws.S3Bucket
		domain.S3Key = aws.S3Key
		domain.S3ObjectVersion = aws.S3ObjectVersion
		domain.Runtime = aws.Runtime
		domain.Handler = aws.Handler
	}

	// Code Source - Container Image
	if aws.PackageType != nil && *aws.PackageType == "Image" && aws.ImageURI != nil && *aws.ImageURI != "" {
		domain.PackageType = aws.PackageType
		domain.ImageURI = aws.ImageURI
	}

	// Configuration
	if aws.MemorySize != nil {
		memorySize := int(*aws.MemorySize)
		domain.MemorySize = &memorySize
	}
	if aws.Timeout != nil {
		timeout := int(*aws.Timeout)
		domain.Timeout = &timeout
	}

	// Environment variables
	if aws.Environment != nil && len(aws.Environment) > 0 {
		domain.Environment = make(map[string]string)
		for k, v := range aws.Environment {
			domain.Environment[k] = v
		}
	}

	// Layers
	if aws.Layers != nil && len(aws.Layers) > 0 {
		domain.Layers = make([]string, len(aws.Layers))
		copy(domain.Layers, aws.Layers)
	}

	// VPC Config
	if aws.VPCConfig != nil {
		domain.VPCConfig = &domaincompute.LambdaVPCConfig{
			SubnetIDs:        make([]string, len(aws.VPCConfig.SubnetIDs)),
			SecurityGroupIDs: make([]string, len(aws.VPCConfig.SecurityGroupIDs)),
		}
		copy(domain.VPCConfig.SubnetIDs, aws.VPCConfig.SubnetIDs)
		copy(domain.VPCConfig.SecurityGroupIDs, aws.VPCConfig.SecurityGroupIDs)
	}

	// Tags - convert from []configs.Tag to map[string]string
	if aws.Tags != nil && len(aws.Tags) > 0 {
		domain.Tags = make(map[string]string)
		for _, tag := range aws.Tags {
			domain.Tags[tag.Key] = tag.Value
		}
	}

	return domain
}

// ToDomainLambdaFunctionFromOutput converts AWS FunctionOutput to domain LambdaFunction
func ToDomainLambdaFunctionFromOutput(output *awslambdaoutputs.FunctionOutput) *domaincompute.LambdaFunction {
	if output == nil {
		return nil
	}

	domain := &domaincompute.LambdaFunction{
		FunctionName: output.FunctionName,
		RoleARN:      output.RoleARN,
		Region:       output.Region,
	}

	// Output fields
	if output.ARN != "" {
		domain.ARN = &output.ARN
	}
	if output.InvokeARN != "" {
		domain.InvokeARN = &output.InvokeARN
	}
	if output.QualifiedARN != nil {
		domain.QualifiedARN = output.QualifiedARN
	}
	if output.Version != "" {
		domain.Version = &output.Version
	}
	if output.LastModified != nil {
		domain.LastModified = output.LastModified
	}
	if output.CodeSize != nil {
		domain.CodeSize = output.CodeSize
	}
	if output.CodeSHA256 != nil {
		domain.CodeSHA256 = output.CodeSHA256
	}

	// Code Source - S3
	if output.S3Bucket != nil && *output.S3Bucket != "" && output.S3Key != nil && *output.S3Key != "" {
		domain.S3Bucket = output.S3Bucket
		domain.S3Key = output.S3Key
		domain.S3ObjectVersion = output.S3ObjectVersion
		domain.Runtime = output.Runtime
		domain.Handler = output.Handler
	}

	// Code Source - Container Image
	if output.PackageType != nil && *output.PackageType == "Image" && output.ImageURI != nil && *output.ImageURI != "" {
		domain.PackageType = output.PackageType
		domain.ImageURI = output.ImageURI
	}

	// Configuration
	if output.MemorySize != nil {
		memorySize := int(*output.MemorySize)
		domain.MemorySize = &memorySize
	}
	if output.Timeout != nil {
		timeout := int(*output.Timeout)
		domain.Timeout = &timeout
	}

	// Environment variables
	if output.Environment != nil && len(output.Environment) > 0 {
		domain.Environment = make(map[string]string)
		for k, v := range output.Environment {
			domain.Environment[k] = v
		}
	}

	// Layers
	if output.Layers != nil && len(output.Layers) > 0 {
		domain.Layers = make([]string, len(output.Layers))
		copy(domain.Layers, output.Layers)
	}

	// VPC Config
	if output.VPCConfig != nil {
		domain.VPCConfig = &domaincompute.LambdaVPCConfig{
			SubnetIDs:        make([]string, len(output.VPCConfig.SubnetIDs)),
			SecurityGroupIDs: make([]string, len(output.VPCConfig.SecurityGroupIDs)),
		}
		copy(domain.VPCConfig.SubnetIDs, output.VPCConfig.SubnetIDs)
		copy(domain.VPCConfig.SecurityGroupIDs, output.VPCConfig.SecurityGroupIDs)
	}

	return domain
}
