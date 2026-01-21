package lambda

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// Function represents an AWS Lambda function configuration
type Function struct {
	FunctionName string // Required: Unique name for the function
	RoleARN      string // Required: IAM Role ARN

	// Code Source (Mutually Exclusive)
	// Method A: S3 Bucket
	S3Bucket        *string
	S3Key           *string
	S3ObjectVersion *string

	// Method B: Container Image
	PackageType *string // "Image" for container image
	ImageURI    *string // ECR URL

	// Runtime Configuration (for S3/zip)
	Runtime *string
	Handler *string

	// Configuration
	MemorySize  *int32             // MB (128-10240)
	Timeout     *int32             // Seconds (1-900)
	Environment map[string]string  // Environment variables
	Layers      []string           // Layer ARNs
	VPCConfig   *FunctionVPCConfig // VPC configuration

	// Tags
	Tags []configs.Tag
}

// FunctionVPCConfig represents VPC configuration for Lambda function
type FunctionVPCConfig struct {
	SubnetIDs        []string
	SecurityGroupIDs []string
}

// Validate performs AWS-specific validation
func (f *Function) Validate() error {
	// Function name is required
	if f.FunctionName == "" {
		return errors.New("function_name is required")
	}

	// Validate function name format
	if err := validateFunctionName(f.FunctionName); err != nil {
		return fmt.Errorf("invalid function_name: %w", err)
	}

	// Role ARN is required
	if f.RoleARN == "" {
		return errors.New("role_arn is required")
	}

	// Validate role ARN format
	if err := validateRoleARN(f.RoleARN); err != nil {
		return fmt.Errorf("invalid role_arn: %w", err)
	}

	// Code source validation - must have either S3 OR Container image (not both)
	hasS3Code := f.S3Bucket != nil && *f.S3Bucket != "" && f.S3Key != nil && *f.S3Key != ""
	hasContainerImage := f.PackageType != nil && *f.PackageType == "Image" && f.ImageURI != nil && *f.ImageURI != ""

	if !hasS3Code && !hasContainerImage {
		return errors.New("either S3 code (s3_bucket and s3_key) or container image (package_type='Image' and image_uri) is required")
	}

	if hasS3Code && hasContainerImage {
		return errors.New("cannot specify both S3 code and container image - choose one deployment method")
	}

	// S3 code validation
	if hasS3Code {
		if f.S3Bucket == nil || *f.S3Bucket == "" {
			return errors.New("s3_bucket is required when using S3 code deployment")
		}
		if f.S3Key == nil || *f.S3Key == "" {
			return errors.New("s3_key is required when using S3 code deployment")
		}
		// Runtime and Handler are required for S3/zip deployments
		if f.Runtime == nil || *f.Runtime == "" {
			return errors.New("runtime is required when using S3 code deployment")
		}
		if f.Handler == nil || *f.Handler == "" {
			return errors.New("handler is required when using S3 code deployment")
		}
	}

	// Container image validation
	if hasContainerImage {
		if f.PackageType == nil || *f.PackageType != "Image" {
			return errors.New("package_type must be 'Image' when using container image deployment")
		}
		if f.ImageURI == nil || *f.ImageURI == "" {
			return errors.New("image_uri is required when using container image deployment")
		}
		// Runtime and Handler should not be set for container images
		if f.Runtime != nil && *f.Runtime != "" {
			return errors.New("runtime should not be set when using container image deployment")
		}
		if f.Handler != nil && *f.Handler != "" {
			return errors.New("handler should not be set when using container image deployment")
		}
	}

	// Memory size validation
	if f.MemorySize != nil {
		if *f.MemorySize < 128 || *f.MemorySize > 10240 {
			return errors.New("memory_size must be between 128 and 10240 MB")
		}
		if *f.MemorySize%64 != 0 {
			return errors.New("memory_size must be a multiple of 64 MB")
		}
	}

	// Timeout validation
	if f.Timeout != nil {
		if *f.Timeout < 1 || *f.Timeout > 900 {
			return errors.New("timeout must be between 1 and 900 seconds")
		}
	}

	// VPC config validation
	if f.VPCConfig != nil {
		if len(f.VPCConfig.SubnetIDs) == 0 && len(f.VPCConfig.SecurityGroupIDs) == 0 {
			return errors.New("vpc_config must have at least one subnet_id or security_group_id")
		}
	}

	return nil
}

// validateFunctionName validates Lambda function name
func validateFunctionName(name string) error {
	if len(name) < 1 {
		return errors.New("function name must be at least 1 character long")
	}
	if len(name) > 64 {
		return errors.New("function name must be at most 64 characters long")
	}

	// Must start with a letter, number, or underscore
	firstChar := name[0]
	if !isAlphanumericOrUnderscore(firstChar) {
		return errors.New("function name must start with a letter, number, or underscore")
	}

	// Can only contain alphanumeric characters, hyphens, and underscores
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validPattern.MatchString(name) {
		return errors.New("function name can only contain alphanumeric characters, hyphens, and underscores")
	}

	return nil
}

// validateRoleARN validates IAM role ARN format
func validateRoleARN(arn string) error {
	if arn == "" {
		return errors.New("role ARN cannot be empty")
	}

	// Basic ARN format validation: arn:aws:iam::account-id:role/role-name
	arnPattern := regexp.MustCompile(`^arn:aws:iam::\d{12}:role/[\w+=,.@-]+$`)
	if !arnPattern.MatchString(arn) {
		return errors.New("role ARN must be in format: arn:aws:iam::account-id:role/role-name")
	}

	return nil
}

// isAlphanumericOrUnderscore checks if a character is alphanumeric or underscore
func isAlphanumericOrUnderscore(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}
