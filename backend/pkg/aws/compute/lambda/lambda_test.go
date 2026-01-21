package lambda

import (
	"context"
	"testing"

	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/compute"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
)

func TestLambdaRunner(t *testing.T) {
	// This test verifies that the Lambda runner can be called without errors
	// In a real scenario, you might want to capture output or verify state
	// For now, we just ensure it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LambdaRunner panicked: %v", r)
		}
	}()

	// Note: This will print to stdout, but that's okay for a runner test
	LambdaRunner()
}

func TestLambdaPricingRunner(t *testing.T) {
	// This test verifies that the Lambda pricing runner can be called without errors
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LambdaPricingRunner panicked: %v", r)
		}
	}()

	// Note: This will print to stdout, but that's okay for a runner test
	LambdaPricingRunner()
}

func TestInMemoryAWSComputeService_CreateLambdaFunction(t *testing.T) {
	service := newInMemoryAWSComputeService()
	ctx := context.Background()

	s3Bucket := "my-bucket"
	s3Key := "code.zip"
	runtime := "python3.9"
	handler := "index.handler"
	memorySize := 256
	timeout := 30

	function := &domaincompute.LambdaFunction{
		FunctionName: "test-function",
		RoleARN:      "arn:aws:iam::123456789012:role/test-role",
		Region:       "us-east-1",
		S3Bucket:     &s3Bucket,
		S3Key:        &s3Key,
		Runtime:      &runtime,
		Handler:      &handler,
		MemorySize:   &memorySize,
		Timeout:      &timeout,
	}

	// Convert to AWS model
	awsFunction := awsmapper.FromDomainLambdaFunction(function)

	output, err := service.CreateLambdaFunction(ctx, awsFunction)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output == nil {
		t.Fatal("Expected output, got nil")
	}

	if output.FunctionName != function.FunctionName {
		t.Errorf("Expected function name %s, got %s", function.FunctionName, output.FunctionName)
	}

	if output.ARN == "" {
		t.Error("Expected ARN to be set")
	}

	if output.InvokeARN == "" {
		t.Error("Expected InvokeARN to be set")
	}
}
