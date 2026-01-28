package compute

import (
	"regexp"

	domainerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/errors"
)

// LambdaFunction represents a cloud-agnostic Lambda function
// This is the domain model - no cloud-specific details
type LambdaFunction struct {
	// Identity (Required)
	FunctionName string // Required: Unique name for the function
	RoleARN      string // Required: IAM Role ARN that gives the function permission to run
	Region       string // Required: AWS region

	// Code Source (Mutually Exclusive - Choose One)
	// Method A: S3 Bucket
	S3Bucket       *string // S3 bucket name containing the code
	S3Key          *string // Path/filename in the bucket (e.g., builds/v1/app.zip)
	S3ObjectVersion *string // Optional: Specific version ID of the S3 object

	// Method B: Container Image
	PackageType *string // "Image" for container image deployment
	ImageURI    *string // ECR URL (e.g., ...amazonaws.com/my-app:latest)

	// Runtime Configuration (for S3/zip deployments)
	Runtime *string // Runtime identifier (e.g., "python3.9", "nodejs18.x")
	Handler *string // Handler function (e.g., "index.handler")

	// Configuration (Optional)
	MemorySize  *int              // RAM in MB (128-10240), default: 128
	Timeout     *int              // Execution time limit in seconds (1-900), default: 3
	Environment map[string]string // Environment variables
	Layers      []string          // List of Lambda Layer ARNs
	VPCConfig   *LambdaVPCConfig // VPC configuration
	Tags        map[string]string // Tags for the function

	// Output fields (populated after creation)
	ARN         *string // The Amazon Resource Name
	InvokeARN   *string // The Invocation ARN (critical for API Gateway)
	QualifiedARN *string // ARN with version suffix
	Version     *string // The latest published version
	LastModified *string // Last modification timestamp
	CodeSize    *int64  // Size of deployment package in bytes
	CodeSHA256  *string // SHA256 hash of the deployment package
}

// LambdaVPCConfig represents VPC configuration for Lambda function
type LambdaVPCConfig struct {
	SubnetIDs        []string // List of private subnet IDs
	SecurityGroupIDs []string // List of security group IDs
}

// Validate performs domain-level validation
func (f *LambdaFunction) Validate() error {
	// Function name is required
	if f.FunctionName == "" {
		return domainerrors.New(domainerrors.CodeLambdaFunctionNameRequired, domainerrors.KindValidation, "function_name is required")
	}

	// Validate function name format
	if err := validateFunctionName(f.FunctionName); err != nil {
		return domainerrors.Wrap(err, domainerrors.CodeLambdaInvalidFunctionName, domainerrors.KindValidation, "invalid function_name")
	}

	// Role ARN is required
	if f.RoleARN == "" {
		return domainerrors.New(domainerrors.CodeLambdaRoleARNRequired, domainerrors.KindValidation, "role_arn is required")
	}

	// Validate role ARN format
	if err := validateRoleARN(f.RoleARN); err != nil {
		return domainerrors.Wrap(err, domainerrors.CodeLambdaInvalidRoleARN, domainerrors.KindValidation, "invalid role_arn")
	}

	// Region is required
	if f.Region == "" {
		return domainerrors.New(domainerrors.CodeLambdaRegionRequired, domainerrors.KindValidation, "region is required")
	}

	// Code source validation - must have either S3 OR Container image (not both)
	hasS3Code := f.S3Bucket != nil && *f.S3Bucket != "" && f.S3Key != nil && *f.S3Key != ""
	hasContainerImage := f.PackageType != nil && *f.PackageType == "Image" && f.ImageURI != nil && *f.ImageURI != ""

	if !hasS3Code && !hasContainerImage {
		return domainerrors.New(domainerrors.CodeLambdaCodeSourceRequired, domainerrors.KindValidation, "either S3 code (s3_bucket and s3_key) or container image (package_type='Image' and image_uri) is required")
	}

	if hasS3Code && hasContainerImage {
		return domainerrors.New(domainerrors.CodeLambdaCodeSourceConflict, domainerrors.KindValidation, "cannot specify both S3 code and container image - choose one deployment method")
	}

	// S3 code validation
	if hasS3Code {
		if f.S3Bucket == nil || *f.S3Bucket == "" {
			return domainerrors.New(domainerrors.CodeLambdaS3BucketRequired, domainerrors.KindValidation, "s3_bucket is required when using S3 code deployment")
		}
		if f.S3Key == nil || *f.S3Key == "" {
			return domainerrors.New(domainerrors.CodeLambdaS3KeyRequired, domainerrors.KindValidation, "s3_key is required when using S3 code deployment")
		}
		// Runtime and Handler are required for S3/zip deployments
		if f.Runtime == nil || *f.Runtime == "" {
			return domainerrors.New(domainerrors.CodeLambdaRuntimeRequired, domainerrors.KindValidation, "runtime is required when using S3 code deployment")
		}
		if f.Handler == nil || *f.Handler == "" {
			return domainerrors.New(domainerrors.CodeLambdaHandlerRequired, domainerrors.KindValidation, "handler is required when using S3 code deployment")
		}
	}

	// Container image validation
	if hasContainerImage {
		if f.PackageType == nil || *f.PackageType != "Image" {
			return domainerrors.New(domainerrors.CodeLambdaPackageTypeInvalid, domainerrors.KindValidation, "package_type must be 'Image' when using container image deployment")
		}
		if f.ImageURI == nil || *f.ImageURI == "" {
			return domainerrors.New(domainerrors.CodeLambdaImageURIRequired, domainerrors.KindValidation, "image_uri is required when using container image deployment")
		}
		// Runtime and Handler should not be set for container images
		if f.Runtime != nil && *f.Runtime != "" {
			return domainerrors.New(domainerrors.CodeLambdaRuntimeNotAllowed, domainerrors.KindValidation, "runtime should not be set when using container image deployment")
		}
		if f.Handler != nil && *f.Handler != "" {
			return domainerrors.New(domainerrors.CodeLambdaHandlerNotAllowed, domainerrors.KindValidation, "handler should not be set when using container image deployment")
		}
	}

	// Memory size validation
	if f.MemorySize != nil {
		if *f.MemorySize < 128 || *f.MemorySize > 10240 {
			return domainerrors.New(domainerrors.CodeLambdaInvalidMemorySize, domainerrors.KindValidation, "memory_size must be between 128 and 10240 MB").
				WithMeta("memory_size", *f.MemorySize)
		}
		if *f.MemorySize%64 != 0 {
			return domainerrors.New(domainerrors.CodeLambdaInvalidMemorySize, domainerrors.KindValidation, "memory_size must be a multiple of 64 MB").
				WithMeta("memory_size", *f.MemorySize)
		}
	}

	// Timeout validation
	if f.Timeout != nil {
		if *f.Timeout < 1 || *f.Timeout > 900 {
			return domainerrors.New(domainerrors.CodeLambdaInvalidTimeout, domainerrors.KindValidation, "timeout must be between 1 and 900 seconds").
				WithMeta("timeout", *f.Timeout)
		}
	}

	// VPC config validation
	if f.VPCConfig != nil {
		if len(f.VPCConfig.SubnetIDs) == 0 && len(f.VPCConfig.SecurityGroupIDs) == 0 {
			return domainerrors.New(domainerrors.CodeLambdaInvalidVPCConfig, domainerrors.KindValidation, "vpc_config must have at least one subnet_id or security_group_id")
		}
	}

	return nil
}

// validateFunctionName validates Lambda function name
// Rules: 1-64 characters, alphanumeric + hyphens/underscores
func validateFunctionName(name string) error {
	if len(name) < 1 {
		return domainerrors.New(domainerrors.CodeLambdaInvalidFunctionName, domainerrors.KindValidation, "function name must be at least 1 character long")
	}
	if len(name) > 64 {
		return domainerrors.New(domainerrors.CodeLambdaInvalidFunctionName, domainerrors.KindValidation, "function name must be at most 64 characters long")
	}

	// Must start with a letter, number, or underscore
	firstChar := name[0]
	if !isAlphanumericOrUnderscore(firstChar) {
		return domainerrors.New(domainerrors.CodeLambdaInvalidFunctionName, domainerrors.KindValidation, "function name must start with a letter, number, or underscore")
	}

	// Can only contain alphanumeric characters, hyphens, and underscores
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validPattern.MatchString(name) {
		return domainerrors.New(domainerrors.CodeLambdaInvalidFunctionName, domainerrors.KindValidation, "function name can only contain alphanumeric characters, hyphens, and underscores")
	}

	return nil
}

// validateRoleARN validates IAM role ARN format
func validateRoleARN(arn string) error {
	if arn == "" {
		return domainerrors.New(domainerrors.CodeLambdaInvalidRoleARN, domainerrors.KindValidation, "role ARN cannot be empty")
	}

	// Basic ARN format validation: arn:aws:iam::account-id:role/role-name
	arnPattern := regexp.MustCompile(`^arn:aws:iam::\d{12}:role/[\w+=,.@-]+$`)
	if !arnPattern.MatchString(arn) {
		return domainerrors.New(domainerrors.CodeLambdaInvalidRoleARN, domainerrors.KindValidation, "role ARN must be in format: arn:aws:iam::account-id:role/role-name")
	}

	return nil
}

// isAlphanumericOrUnderscore checks if a character is alphanumeric or underscore
func isAlphanumericOrUnderscore(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// Ensure LambdaFunction implements ComputeResource
var _ ComputeResource = (*LambdaFunction)(nil)

// GetID returns the function name (used as ID)
func (f *LambdaFunction) GetID() string {
	return f.FunctionName
}

// GetName returns the function name
func (f *LambdaFunction) GetName() string {
	return f.FunctionName
}

// GetRegion returns the function region
func (f *LambdaFunction) GetRegion() string {
	return f.Region
}

// GetSubnetID returns the first subnet ID from VPC config if available
// Lambda functions don't have a single subnet ID, but may have VPC config with multiple subnets
func (f *LambdaFunction) GetSubnetID() string {
	if f.VPCConfig != nil && len(f.VPCConfig.SubnetIDs) > 0 {
		return f.VPCConfig.SubnetIDs[0]
	}
	return ""
}

// HasS3Code returns true if function uses S3 code deployment
func (f *LambdaFunction) HasS3Code() bool {
	return f.S3Bucket != nil && *f.S3Bucket != "" && f.S3Key != nil && *f.S3Key != ""
}

// HasContainerImage returns true if function uses container image deployment
func (f *LambdaFunction) HasContainerImage() bool {
	return f.PackageType != nil && *f.PackageType == "Image" && f.ImageURI != nil && *f.ImageURI != ""
}
