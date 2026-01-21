package compute

import (
	"testing"

	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	awslambda "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda"
	awslambdaoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda/outputs"
)

func TestFromDomainLambdaFunction(t *testing.T) {
	tests := []struct {
		name     string
		domain   *domaincompute.LambdaFunction
		checkFunc func(*testing.T, *awslambda.Function)
	}{
		{
			name: "s3-code-deployment",
			domain: &domaincompute.LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				S3Bucket:     stringPtr("my-bucket"),
				S3Key:        stringPtr("code.zip"),
				Runtime:      stringPtr("python3.9"),
				Handler:      stringPtr("index.handler"),
				MemorySize:   intPtr(256),
				Timeout:      intPtr(30),
			},
			checkFunc: func(t *testing.T, aws *awslambda.Function) {
				if aws == nil {
					t.Fatal("expected AWS function, got nil")
				}
				if aws.FunctionName != "test-function" {
					t.Errorf("expected function name 'test-function', got '%s'", aws.FunctionName)
				}
				if aws.S3Bucket == nil || *aws.S3Bucket != "my-bucket" {
					t.Errorf("expected S3 bucket 'my-bucket', got '%v'", aws.S3Bucket)
				}
				if aws.Runtime == nil || *aws.Runtime != "python3.9" {
					t.Errorf("expected runtime 'python3.9', got '%v'", aws.Runtime)
				}
			},
		},
		{
			name: "container-image-deployment",
			domain: &domaincompute.LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				PackageType:  stringPtr("Image"),
				ImageURI:     stringPtr("123456789012.dkr.ecr.us-east-1.amazonaws.com/my-app:latest"),
				MemorySize:   intPtr(512),
			},
			checkFunc: func(t *testing.T, aws *awslambda.Function) {
				if aws == nil {
					t.Fatal("expected AWS function, got nil")
				}
				if aws.PackageType == nil || *aws.PackageType != "Image" {
					t.Errorf("expected package type 'Image', got '%v'", aws.PackageType)
				}
				if aws.ImageURI == nil || *aws.ImageURI != "123456789012.dkr.ecr.us-east-1.amazonaws.com/my-app:latest" {
					t.Errorf("expected image URI, got '%v'", aws.ImageURI)
				}
			},
		},
		{
			name:   "nil-input",
			domain: nil,
			checkFunc: func(t *testing.T, aws *awslambda.Function) {
				if aws != nil {
					t.Errorf("expected nil, got %v", aws)
				}
			},
		},
		{
			name: "with-environment-variables",
			domain: &domaincompute.LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				S3Bucket:     stringPtr("my-bucket"),
				S3Key:        stringPtr("code.zip"),
				Runtime:      stringPtr("python3.9"),
				Handler:      stringPtr("index.handler"),
				Environment: map[string]string{
					"ENV_VAR_1": "value1",
					"ENV_VAR_2": "value2",
				},
			},
			checkFunc: func(t *testing.T, aws *awslambda.Function) {
				if aws.Environment == nil {
					t.Fatal("expected environment variables, got nil")
				}
				if aws.Environment["ENV_VAR_1"] != "value1" {
					t.Errorf("expected ENV_VAR_1='value1', got '%s'", aws.Environment["ENV_VAR_1"])
				}
			},
		},
		{
			name: "with-vpc-config",
			domain: &domaincompute.LambdaFunction{
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				S3Bucket:     stringPtr("my-bucket"),
				S3Key:        stringPtr("code.zip"),
				Runtime:      stringPtr("python3.9"),
				Handler:      stringPtr("index.handler"),
				VPCConfig: &domaincompute.LambdaVPCConfig{
					SubnetIDs:        []string{"subnet-123", "subnet-456"},
					SecurityGroupIDs: []string{"sg-123"},
				},
			},
			checkFunc: func(t *testing.T, aws *awslambda.Function) {
				if aws.VPCConfig == nil {
					t.Fatal("expected VPC config, got nil")
				}
				if len(aws.VPCConfig.SubnetIDs) != 2 {
					t.Errorf("expected 2 subnets, got %d", len(aws.VPCConfig.SubnetIDs))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aws := FromDomainLambdaFunction(tt.domain)
			if tt.checkFunc != nil {
				tt.checkFunc(t, aws)
			}
		})
	}
}

func TestToDomainLambdaFunctionFromOutput(t *testing.T) {
	tests := []struct {
		name     string
		output   *awslambdaoutputs.FunctionOutput
		checkFunc func(*testing.T, *domaincompute.LambdaFunction)
	}{
		{
			name: "basic-output",
			output: &awslambdaoutputs.FunctionOutput{
				ARN:          "arn:aws:lambda:us-east-1:123456789012:function:test-function",
				InvokeARN:    "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789012:function:test-function/invocations",
				FunctionName: "test-function",
				RoleARN:      "arn:aws:iam::123456789012:role/test-role",
				Region:       "us-east-1",
				Version:      "1",
				Runtime:      stringPtr("python3.9"),
				Handler:      stringPtr("index.handler"),
			},
			checkFunc: func(t *testing.T, domain *domaincompute.LambdaFunction) {
				if domain == nil {
					t.Fatal("expected domain function, got nil")
				}
				if domain.ARN == nil || *domain.ARN != "arn:aws:lambda:us-east-1:123456789012:function:test-function" {
					t.Errorf("expected ARN, got '%v'", domain.ARN)
				}
				if domain.InvokeARN == nil || *domain.InvokeARN == "" {
					t.Error("expected InvokeARN, got nil")
				}
			},
		},
		{
			name:   "nil-output",
			output: nil,
			checkFunc: func(t *testing.T, domain *domaincompute.LambdaFunction) {
				if domain != nil {
					t.Errorf("expected nil, got %v", domain)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := ToDomainLambdaFunctionFromOutput(tt.output)
			if tt.checkFunc != nil {
				tt.checkFunc(t, domain)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
