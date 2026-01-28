package errors

// Domain-specific error codes
const (
	// Validation errors
	CodeValidationFailed        = "VALIDATION_FAILED"
	CodeRequiredFieldMissing    = "REQUIRED_FIELD_MISSING"
	CodeInvalidFormat           = "INVALID_FORMAT"
	CodeInvalidValue            = "INVALID_VALUE"
	CodeValueOutOfRange         = "VALUE_OUT_OF_RANGE"

	// Lambda function errors
	CodeLambdaFunctionNameRequired     = "LAMBDA_FUNCTION_NAME_REQUIRED"
	CodeLambdaInvalidFunctionName      = "LAMBDA_INVALID_FUNCTION_NAME"
	CodeLambdaRoleARNRequired          = "LAMBDA_ROLE_ARN_REQUIRED"
	CodeLambdaInvalidRoleARN           = "LAMBDA_INVALID_ROLE_ARN"
	CodeLambdaRegionRequired           = "LAMBDA_REGION_REQUIRED"
	CodeLambdaCodeSourceRequired       = "LAMBDA_CODE_SOURCE_REQUIRED"
	CodeLambdaCodeSourceConflict       = "LAMBDA_CODE_SOURCE_CONFLICT"
	CodeLambdaS3BucketRequired         = "LAMBDA_S3_BUCKET_REQUIRED"
	CodeLambdaS3KeyRequired            = "LAMBDA_S3_KEY_REQUIRED"
	CodeLambdaRuntimeRequired          = "LAMBDA_RUNTIME_REQUIRED"
	CodeLambdaHandlerRequired          = "LAMBDA_HANDLER_REQUIRED"
	CodeLambdaPackageTypeInvalid       = "LAMBDA_PACKAGE_TYPE_INVALID"
	CodeLambdaImageURIRequired         = "LAMBDA_IMAGE_URI_REQUIRED"
	CodeLambdaRuntimeNotAllowed        = "LAMBDA_RUNTIME_NOT_ALLOWED"
	CodeLambdaHandlerNotAllowed        = "LAMBDA_HANDLER_NOT_ALLOWED"
	CodeLambdaInvalidMemorySize        = "LAMBDA_INVALID_MEMORY_SIZE"
	CodeLambdaInvalidTimeout           = "LAMBDA_INVALID_TIMEOUT"
	CodeLambdaInvalidVPCConfig         = "LAMBDA_INVALID_VPC_CONFIG"

	// Storage errors
	CodeS3BucketRequired              = "S3_BUCKET_REQUIRED"
	CodeS3InvalidSSEAlgorithm         = "S3_INVALID_SSE_ALGORITHM"
)
