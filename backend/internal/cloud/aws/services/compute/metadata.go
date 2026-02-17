package compute

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services"
)

// ---------------------------------------------------------------------------
// Internal registry
// ---------------------------------------------------------------------------

var computeSchemaRegistry = map[string]*services.ResourceSchema{}

func registerComputeSchema(label string, schema *services.ResourceSchema) {
	computeSchemaRegistry[label] = schema
}

// ---------------------------------------------------------------------------
// Service interface & implementation
// ---------------------------------------------------------------------------

// ComputeMetadataService exposes structured schemas for AWS compute resources.
type ComputeMetadataService interface {
	GetResourceSchema(ctx context.Context, resource string) (*services.ResourceSchema, error)
	ListResourceSchemas(ctx context.Context) ([]*services.ResourceSchema, error)
}

type computeMetadataServiceImpl struct{}

// NewComputeMetadataService returns a ready-to-use metadata service.
func NewComputeMetadataService() ComputeMetadataService {
	return &computeMetadataServiceImpl{}
}

func (s *computeMetadataServiceImpl) GetResourceSchema(_ context.Context, resource string) (*services.ResourceSchema, error) {
	schema, ok := computeSchemaRegistry[resource]
	if !ok {
		return nil, fmt.Errorf("unknown compute resource: %s", resource)
	}
	return schema, nil
}

func (s *computeMetadataServiceImpl) ListResourceSchemas(_ context.Context) ([]*services.ResourceSchema, error) {
	out := make([]*services.ResourceSchema, 0, len(computeSchemaRegistry))
	for _, schema := range computeSchemaRegistry {
		out = append(out, schema)
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Schema registrations (run once at import time)
// ---------------------------------------------------------------------------

func init() {
	// ---- EC2 Instance ----
	registerComputeSchema("ec2_instance", &services.ResourceSchema{
		Label: "ec2_instance",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the instance"},
			{Name: "ami", Type: "string", Required: true, Enum: []string{}, Description: "AMI ID (e.g. ami-0abcdef1234567890)"},
			{Name: "instance_type", Type: "string", Required: true, Enum: []string{
				"t2.micro", "t2.small", "t2.medium",
				"t3.micro", "t3.small", "t3.medium", "t3.large", "t3.xlarge", "t3.2xlarge",
				"t4g.micro", "t4g.small",
				"m5.large", "m5.xlarge", "m5.2xlarge", "m5.4xlarge",
				"m6i.large", "m6i.xlarge",
				"c5.large", "c5.xlarge", "c5.2xlarge",
				"c6i.large", "c6i.xlarge",
				"r5.large", "r5.xlarge", "r5.2xlarge",
				"r6i.large", "r6i.xlarge",
				"i3.large", "i3.xlarge", "i3.2xlarge",
				"g4dn.xlarge", "g4dn.2xlarge",
			}, Description: "Instance type (e.g. t3.micro, m5.large)"},
			{Name: "subnet_id", Type: "string", Required: true, Enum: []string{}, Description: "Subnet ID to launch in"},
			{Name: "vpc_security_group_ids", Type: "[]string", Required: true, Enum: []string{}, Description: "Security group IDs (at least one)"},
			{Name: "associate_public_ip_address", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Auto-assign public IP"},
			{Name: "key_name", Type: "string", Required: false, Enum: []string{}, Description: "SSH key pair name"},
			{Name: "iam_instance_profile", Type: "string", Required: false, Enum: []string{}, Description: "IAM instance profile name or ARN"},
			{Name: "user_data", Type: "string", Required: false, Enum: []string{}, Description: "User data script (max 12KB raw / 16KB base64)"},
			{Name: "root_volume_id", Type: "string", Required: false, Enum: []string{}, Description: "Reference to root EBS volume"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":                "string",
			"arn":               "string",
			"public_ip":         "string",
			"private_ip":        "string",
			"state":             "string",
			"availability_zone": "string",
		},
	})

	// ---- Launch Template ----
	registerComputeSchema("launch_template", &services.ResourceSchema{
		Label: "launch_template",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: false, Enum: []string{}, Description: "Name of the launch template (alternative to name_prefix)"},
			{Name: "name_prefix", Type: "string", Required: false, Enum: []string{}, Description: "Name prefix for unique naming (recommended)"},
			{Name: "image_id", Type: "string", Required: true, Enum: []string{}, Description: "AMI ID"},
			{Name: "instance_type", Type: "string", Required: true, Enum: []string{}, Description: "Instance type (e.g. t3.micro)"},
			{Name: "vpc_security_group_ids", Type: "[]string", Required: true, Enum: []string{}, Description: "Security group IDs"},
			{Name: "key_name", Type: "string", Required: false, Enum: []string{}, Description: "SSH key pair name"},
			{Name: "iam_instance_profile", Type: "object", Required: false, Enum: []string{}, Description: "IAM instance profile configuration"},
			{Name: "user_data", Type: "string", Required: false, Enum: []string{}, Description: "Base64-encoded user data (max 16KB)"},
			{Name: "root_volume_id", Type: "string", Required: false, Enum: []string{}, Description: "Reference to root storage volume"},
			{Name: "additional_volume_ids", Type: "[]string", Required: false, Enum: []string{}, Description: "References to additional storage volumes"},
			{Name: "metadata_options", Type: "object", Required: false, Enum: []string{}, Description: "IMDSv2 metadata settings"},
			{Name: "update_default_version", Type: "bool", Required: false, Enum: []string{}, Default: true, Description: "Update default version on changes"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":              "string",
			"arn":             "string",
			"latest_version":  "int",
			"default_version": "int",
		},
	})

	// ---- Load Balancer ----
	registerComputeSchema("load_balancer", &services.ResourceSchema{
		Label: "load_balancer",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the load balancer"},
			{Name: "load_balancer_type", Type: "string", Required: true, Enum: []string{"application", "network"}, Description: "Load balancer type"},
			{Name: "internal", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Whether the LB is internal"},
			{Name: "security_group_ids", Type: "[]string", Required: false, Enum: []string{}, Description: "Security group IDs (required for ALB)"},
			{Name: "subnet_ids", Type: "[]string", Required: true, Enum: []string{}, Description: "Subnet IDs (at least 2 in different AZs)"},
			{Name: "ip_address_type", Type: "string", Required: false, Enum: []string{"ipv4", "dualstack"}, Default: "ipv4", Description: "IP address type"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":       "string",
			"arn":      "string",
			"dns_name": "string",
			"zone_id":  "string",
		},
	})

	// ---- Target Group ----
	registerComputeSchema("target_group", &services.ResourceSchema{
		Label: "target_group",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the target group"},
			{Name: "port", Type: "int", Required: true, Enum: []string{}, Description: "Port (1-65535)"},
			{Name: "protocol", Type: "string", Required: true, Enum: []string{"HTTP", "HTTPS", "TCP", "TLS"}, Description: "Protocol"},
			{Name: "vpc_id", Type: "string", Required: true, Enum: []string{}, Description: "VPC ID"},
			{Name: "target_type", Type: "string", Required: false, Enum: []string{"instance", "ip", "lambda"}, Default: "instance", Description: "Target type"},
			{Name: "health_check", Type: "object", Required: false, Enum: []string{}, Description: "Health check configuration"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":  "string",
			"arn": "string",
		},
	})

	// ---- Listener ----
	registerComputeSchema("listener", &services.ResourceSchema{
		Label: "listener",
		Fields: []services.FieldDescriptor{
			{Name: "load_balancer_arn", Type: "string", Required: true, Enum: []string{}, Description: "Load balancer ARN"},
			{Name: "port", Type: "int", Required: true, Enum: []string{}, Description: "Listener port (1-65535)"},
			{Name: "protocol", Type: "string", Required: true, Enum: []string{"HTTP", "HTTPS", "TCP", "TLS"}, Description: "Listener protocol"},
			{Name: "default_action", Type: "object", Required: true, Enum: []string{}, Description: "Default action (forward, redirect, or fixed-response)"},
			{Name: "certificate_arn", Type: "string", Required: false, Enum: []string{}, Description: "ACM certificate ARN (required for HTTPS/TLS)"},
			{Name: "ssl_policy", Type: "string", Required: false, Enum: []string{}, Description: "SSL policy for HTTPS/TLS listeners"},
		},
		Outputs: map[string]string{
			"id":  "string",
			"arn": "string",
		},
	})

	// ---- Auto Scaling Group ----
	registerComputeSchema("auto_scaling_group", &services.ResourceSchema{
		Label: "auto_scaling_group",
		Fields: []services.FieldDescriptor{
			{Name: "auto_scaling_group_name", Type: "string", Required: false, Enum: []string{}, Description: "Exact name for the ASG"},
			{Name: "auto_scaling_group_name_prefix", Type: "string", Required: false, Enum: []string{}, Description: "Name prefix (alternative to exact name)"},
			{Name: "min_size", Type: "int", Required: true, Enum: []string{}, Description: "Minimum number of instances"},
			{Name: "max_size", Type: "int", Required: true, Enum: []string{}, Description: "Maximum number of instances (max 10000)"},
			{Name: "desired_capacity", Type: "int", Required: false, Enum: []string{}, Description: "Desired number of instances"},
			{Name: "vpc_zone_identifier", Type: "[]string", Required: true, Enum: []string{}, Description: "Subnet IDs for the ASG"},
			{Name: "launch_template", Type: "object", Required: true, Enum: []string{}, Description: "Launch template specification (id + version)"},
			{Name: "health_check_type", Type: "string", Required: false, Enum: []string{"EC2", "ELB"}, Default: "EC2", Description: "Health check type"},
			{Name: "health_check_grace_period", Type: "int", Required: false, Enum: []string{}, Default: 300, Description: "Health check grace period in seconds"},
			{Name: "target_group_arns", Type: "[]string", Required: false, Enum: []string{}, Description: "Target group ARNs (required for ELB health checks)"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Tags with propagation settings"},
		},
		Outputs: map[string]string{
			"id":               "string",
			"arn":              "string",
			"name":             "string",
			"desired_capacity": "int",
			"min_size":         "int",
			"max_size":         "int",
		},
	})

	// ---- Lambda Function ----
	registerComputeSchema("lambda_function", &services.ResourceSchema{
		Label: "lambda_function",
		Fields: []services.FieldDescriptor{
			{Name: "function_name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the Lambda function (1-64 chars)"},
			{Name: "role_arn", Type: "string", Required: true, Enum: []string{}, Description: "IAM role ARN for the function"},
			{Name: "runtime", Type: "string", Required: false, Enum: []string{"nodejs18.x", "nodejs20.x", "python3.9", "python3.10", "python3.11", "python3.12", "java17", "java21", "go1.x", "dotnet6", "dotnet8", "ruby3.2", "ruby3.3"}, Description: "Runtime (required for S3/zip deployment)"},
			{Name: "handler", Type: "string", Required: false, Enum: []string{}, Description: "Handler function (required for S3/zip deployment)"},
			{Name: "memory_size", Type: "int", Required: false, Enum: []string{}, Default: 128, Description: "Memory in MB (128-10240, must be multiple of 64)"},
			{Name: "timeout", Type: "int", Required: false, Enum: []string{}, Default: 3, Description: "Timeout in seconds (1-900)"},
			{Name: "s3_bucket", Type: "string", Required: false, Enum: []string{}, Description: "S3 bucket with deployment package"},
			{Name: "s3_key", Type: "string", Required: false, Enum: []string{}, Description: "S3 key for deployment package"},
			{Name: "s3_object_version", Type: "string", Required: false, Enum: []string{}, Description: "S3 object version"},
			{Name: "package_type", Type: "string", Required: false, Enum: []string{"Zip", "Image"}, Default: "Zip", Description: "Deployment package type"},
			{Name: "image_uri", Type: "string", Required: false, Enum: []string{}, Description: "ECR image URI (for container deployment)"},
			{Name: "environment", Type: "object", Required: false, Enum: []string{}, Description: "Environment variables"},
			{Name: "layers", Type: "[]string", Required: false, Enum: []string{}, Description: "Lambda layer ARNs"},
			{Name: "vpc_config", Type: "object", Required: false, Enum: []string{}, Description: "VPC configuration (subnet_ids + security_group_ids)"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"arn":           "string",
			"function_name": "string",
			"version":       "string",
			"qualified_arn": "string",
			"invoke_arn":    "string",
			"last_modified": "string",
		},
	})
}
