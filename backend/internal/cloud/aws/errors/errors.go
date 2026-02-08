package errors

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/errors"
)

// AWS-specific error codes
const (
	// EC2 errors
	CodeEC2InstanceNotFound       = "EC2_INSTANCE_NOT_FOUND"
	CodeEC2InvalidInstanceType    = "EC2_INVALID_INSTANCE_TYPE"
	CodeEC2InvalidAMI             = "EC2_INVALID_AMI"
	CodeEC2InvalidSubnet          = "EC2_INVALID_SUBNET"
	CodeEC2InvalidSecurityGroup   = "EC2_INVALID_SECURITY_GROUP"
	CodeEC2InstanceCreationFailed = "EC2_INSTANCE_CREATION_FAILED"

	// Lambda errors
	CodeLambdaFunctionNotFound       = "LAMBDA_FUNCTION_NOT_FOUND"
	CodeLambdaInvalidFunctionName    = "LAMBDA_INVALID_FUNCTION_NAME"
	CodeLambdaInvalidRoleARN         = "LAMBDA_INVALID_ROLE_ARN"
	CodeLambdaInvalidRuntime         = "LAMBDA_INVALID_RUNTIME"
	CodeLambdaFunctionCreationFailed = "LAMBDA_FUNCTION_CREATION_FAILED"

	// Load Balancer errors
	CodeLoadBalancerNotFound       = "LOAD_BALANCER_NOT_FOUND"
	CodeLoadBalancerInvalidARN     = "LOAD_BALANCER_INVALID_ARN"
	CodeLoadBalancerCreationFailed = "LOAD_BALANCER_CREATION_FAILED"

	// Auto Scaling errors
	CodeAutoScalingGroupNotFound       = "AUTO_SCALING_GROUP_NOT_FOUND"
	CodeAutoScalingGroupCreationFailed = "AUTO_SCALING_GROUP_CREATION_FAILED"

	// Networking errors
	CodeVPCNotFound             = "VPC_NOT_FOUND"
	CodeSubnetNotFound          = "SUBNET_NOT_FOUND"
	CodeSecurityGroupNotFound   = "SECURITY_GROUP_NOT_FOUND"
	CodeNATGatewayNotFound      = "NAT_GATEWAY_NOT_FOUND"
	CodeInternetGatewayNotFound = "INTERNET_GATEWAY_NOT_FOUND"
	CodeRouteTableNotFound      = "ROUTE_TABLE_NOT_FOUND"
	CodeElasticIPNotFound       = "ELASTIC_IP_NOT_FOUND"
	CodeInvalidCIDR             = "INVALID_CIDR"
	CodeCIDROverlap             = "CIDR_OVERLAP"
	CodeInvalidAvailabilityZone = "INVALID_AVAILABILITY_ZONE"

	// Storage errors
	CodeS3BucketNotFound        = "S3_BUCKET_NOT_FOUND"
	CodeS3BucketCreationFailed  = "S3_BUCKET_CREATION_FAILED"
	CodeEBSVolumeNotFound       = "EBS_VOLUME_NOT_FOUND"
	CodeEBSVolumeCreationFailed = "EBS_VOLUME_CREATION_FAILED"

	// IAM errors
	CodeIAMRoleNotFound            = "IAM_ROLE_NOT_FOUND"
	CodeIAMPolicyNotFound          = "IAM_POLICY_NOT_FOUND"
	CodeIAMUserNotFound            = "IAM_USER_NOT_FOUND"
	CodeIAMGroupNotFound           = "IAM_GROUP_NOT_FOUND"
	CodeIAMInstanceProfileNotFound = "IAM_INSTANCE_PROFILE_NOT_FOUND"
	CodeIAMInvalidARN              = "IAM_INVALID_ARN"

	// Instance Type errors
	CodeInstanceTypeNotFound     = "INSTANCE_TYPE_NOT_FOUND"
	CodeInstanceTypeNameRequired = "INSTANCE_TYPE_NAME_REQUIRED"
	CodeInvalidCategory          = "INVALID_CATEGORY"
	CodeInvalidVCPU              = "INVALID_VCPU"
	CodeInvalidMemory            = "INVALID_MEMORY"

	// RDS errors
	CodeRDSInstanceNotFound       = "RDS_INSTANCE_NOT_FOUND"
	CodeRDSInstanceCreationFailed = "RDS_INSTANCE_CREATION_FAILED"

	// ECS errors
	CodeECSClusterNotFound                = "ECS_CLUSTER_NOT_FOUND"
	CodeECSServiceNotFound                = "ECS_SERVICE_NOT_FOUND"
	CodeECSTaskDefinitionNotFound         = "ECS_TASK_DEFINITION_NOT_FOUND"
	CodeECSCapacityProviderNotFound       = "ECS_CAPACITY_PROVIDER_NOT_FOUND"
	CodeECSInvalidLaunchType              = "ECS_INVALID_LAUNCH_TYPE"
	CodeECSClusterCreationFailed          = "ECS_CLUSTER_CREATION_FAILED"
	CodeECSServiceCreationFailed          = "ECS_SERVICE_CREATION_FAILED"
	CodeECSTaskDefinitionCreationFailed   = "ECS_TASK_DEFINITION_CREATION_FAILED"
	CodeECSCapacityProviderCreationFailed = "ECS_CAPACITY_PROVIDER_CREATION_FAILED"

	// SDK/API errors
	CodeAWSSDKError    = "AWS_SDK_ERROR"
	CodeAWSAPIError    = "AWS_API_ERROR"
	CodeAWSConfigError = "AWS_CONFIG_ERROR"
)

// NewEC2InstanceNotFound creates an error for when an EC2 instance is not found
func NewEC2InstanceNotFound(instanceID string) *errors.AppError {
	return errors.New(CodeEC2InstanceNotFound, errors.KindNotFound, "EC2 instance not found").
		WithMeta("instance_id", instanceID)
}

// NewEC2InvalidInstanceType creates an error for invalid instance type
func NewEC2InvalidInstanceType(instanceType string) *errors.AppError {
	return errors.New(CodeEC2InvalidInstanceType, errors.KindValidation, "Invalid EC2 instance type").
		WithMeta("instance_type", instanceType)
}

// NewLambdaFunctionNotFound creates an error for when a Lambda function is not found
func NewLambdaFunctionNotFound(functionName string) *errors.AppError {
	return errors.New(CodeLambdaFunctionNotFound, errors.KindNotFound, "Lambda function not found").
		WithMeta("function_name", functionName)
}

// NewLambdaInvalidFunctionName creates an error for invalid Lambda function name
func NewLambdaInvalidFunctionName(functionName string, reason string) *errors.AppError {
	return errors.New(CodeLambdaInvalidFunctionName, errors.KindValidation, "Invalid Lambda function name").
		WithMeta("function_name", functionName).
		WithMeta("reason", reason)
}

// NewLambdaInvalidRoleARN creates an error for invalid IAM role ARN
func NewLambdaInvalidRoleARN(roleARN string) *errors.AppError {
	return errors.New(CodeLambdaInvalidRoleARN, errors.KindValidation, "Invalid IAM role ARN").
		WithMeta("role_arn", roleARN)
}

// NewVPCNotFound creates an error for when a VPC is not found
func NewVPCNotFound(vpcID string) *errors.AppError {
	return errors.New(CodeVPCNotFound, errors.KindNotFound, "VPC not found").
		WithMeta("vpc_id", vpcID)
}

// NewSubnetNotFound creates an error for when a subnet is not found
func NewSubnetNotFound(subnetID string) *errors.AppError {
	return errors.New(CodeSubnetNotFound, errors.KindNotFound, "Subnet not found").
		WithMeta("subnet_id", subnetID)
}

// NewInvalidCIDR creates an error for invalid CIDR block
func NewInvalidCIDR(cidr string, reason string) *errors.AppError {
	return errors.New(CodeInvalidCIDR, errors.KindValidation, "Invalid CIDR block").
		WithMeta("cidr", cidr).
		WithMeta("reason", reason)
}

// NewCIDROverlap creates an error for CIDR overlap
func NewCIDROverlap(cidr1, cidr2 string) *errors.AppError {
	return errors.New(CodeCIDROverlap, errors.KindConflict, "CIDR blocks overlap").
		WithMeta("cidr1", cidr1).
		WithMeta("cidr2", cidr2)
}

// NewInstanceTypeNotFound creates an error for when an instance type is not found
func NewInstanceTypeNotFound(instanceType string) *errors.AppError {
	return errors.New(CodeInstanceTypeNotFound, errors.KindNotFound, "Instance type not found").
		WithMeta("instance_type", instanceType)
}

// NewInstanceTypeNameRequired creates an error for missing instance type name
func NewInstanceTypeNameRequired() *errors.AppError {
	return errors.New(CodeInstanceTypeNameRequired, errors.KindValidation, "Instance type name is required")
}

// NewInvalidCategory creates an error for invalid instance category
func NewInvalidCategory(category string) *errors.AppError {
	return errors.New(CodeInvalidCategory, errors.KindValidation, "Invalid instance category").
		WithMeta("category", category)
}

// NewRDSInstanceNotFound creates an error for when an RDS instance is not found
func NewRDSInstanceNotFound(instanceID string) *errors.AppError {
	return errors.New(CodeRDSInstanceNotFound, errors.KindNotFound, "RDS instance not found").
		WithMeta("instance_id", instanceID)
}

// NewInvalidVCPU creates an error for invalid VCPU count
func NewInvalidVCPU(vcpu int) *errors.AppError {
	return errors.New(CodeInvalidVCPU, errors.KindValidation, "VCPU count must be greater than 0").
		WithMeta("vcpu", vcpu)
}

// NewInvalidMemory creates an error for invalid memory
func NewInvalidMemory(memory int) *errors.AppError {
	return errors.New(CodeInvalidMemory, errors.KindValidation, "Memory must be greater than 0").
		WithMeta("memory", memory)
}

// NewAWSSDKError wraps an AWS SDK error
func NewAWSSDKError(operation string, cause error) *errors.AppError {
	return errors.Wrap(cause, CodeAWSSDKError, errors.KindInternal, "AWS SDK error").
		WithOp(operation)
}

// NewAWSAPIError creates an error for AWS API failures
func NewAWSAPIError(operation string, cause error) *errors.AppError {
	return errors.Wrap(cause, CodeAWSAPIError, errors.KindInternal, "AWS API error").
		WithOp(operation)
}

// NewAWSConfigError creates an error for AWS configuration issues
func NewAWSConfigError(reason string) *errors.AppError {
	return errors.New(CodeAWSConfigError, errors.KindInternal, "AWS configuration error").
		WithMeta("reason", reason)
}

// ECS Error Helpers

// NewECSClusterNotFound creates an error for when an ECS cluster is not found
func NewECSClusterNotFound(clusterName string) *errors.AppError {
	return errors.New(CodeECSClusterNotFound, errors.KindNotFound, "ECS cluster not found").
		WithMeta("cluster_name", clusterName)
}

// NewECSServiceNotFound creates an error for when an ECS service is not found
func NewECSServiceNotFound(serviceName, clusterName string) *errors.AppError {
	return errors.New(CodeECSServiceNotFound, errors.KindNotFound, "ECS service not found").
		WithMeta("service_name", serviceName).
		WithMeta("cluster_name", clusterName)
}

// NewECSTaskDefinitionNotFound creates an error for when an ECS task definition is not found
func NewECSTaskDefinitionNotFound(family string) *errors.AppError {
	return errors.New(CodeECSTaskDefinitionNotFound, errors.KindNotFound, "ECS task definition not found").
		WithMeta("family", family)
}

// NewECSCapacityProviderNotFound creates an error for when an ECS capacity provider is not found
func NewECSCapacityProviderNotFound(providerName string) *errors.AppError {
	return errors.New(CodeECSCapacityProviderNotFound, errors.KindNotFound, "ECS capacity provider not found").
		WithMeta("provider_name", providerName)
}

// NewECSInvalidLaunchType creates an error for invalid ECS launch type
func NewECSInvalidLaunchType(launchType string) *errors.AppError {
	return errors.New(CodeECSInvalidLaunchType, errors.KindValidation, "Invalid ECS launch type").
		WithMeta("launch_type", launchType).
		WithMeta("valid_types", "FARGATE, EC2")
}

// NewECSClusterCreationFailed creates an error for ECS cluster creation failure
func NewECSClusterCreationFailed(clusterName string, cause error) *errors.AppError {
	return errors.Wrap(cause, CodeECSClusterCreationFailed, errors.KindInternal, "ECS cluster creation failed").
		WithMeta("cluster_name", clusterName)
}

// NewECSServiceCreationFailed creates an error for ECS service creation failure
func NewECSServiceCreationFailed(serviceName, clusterName string, cause error) *errors.AppError {
	return errors.Wrap(cause, CodeECSServiceCreationFailed, errors.KindInternal, "ECS service creation failed").
		WithMeta("service_name", serviceName).
		WithMeta("cluster_name", clusterName)
}
