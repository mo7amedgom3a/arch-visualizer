package outputs

// FunctionOutput represents the output from AWS Lambda function operations
type FunctionOutput struct {
	ARN         string  // The Amazon Resource Name
	InvokeARN   string  // The Invocation ARN (critical for API Gateway)
	QualifiedARN *string // ARN with version suffix
	Version     string  // The latest published version
	FunctionName string
	Region      string
	RoleARN     string

	// Code Source
	S3Bucket       *string
	S3Key          *string
	S3ObjectVersion *string
	PackageType    *string
	ImageURI       *string

	// Runtime Configuration
	Runtime *string
	Handler *string

	// Configuration
	MemorySize  *int32
	Timeout     *int32
	Environment map[string]string
	Layers      []string
	VPCConfig   *FunctionVPCConfigOutput

	// Metadata
	LastModified *string
	CodeSize     *int64
	CodeSHA256   *string
	Description  *string
}

// FunctionVPCConfigOutput represents VPC configuration output
type FunctionVPCConfigOutput struct {
	SubnetIDs        []string
	SecurityGroupIDs []string
	VPCID            *string
}
