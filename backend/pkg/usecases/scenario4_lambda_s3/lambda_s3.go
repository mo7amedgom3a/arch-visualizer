package scenario4_lambda_s3

import (
	"context"
	"fmt"
	"time"

	awspricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	domainiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/iam"
	usecasescommon "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/common"
)

// LambdaS3Runner demonstrates a serverless architecture with Lambda function accessing S3 bucket
func LambdaS3Runner() {
	ctx := context.Background()
	region := usecasescommon.SelectRegion("us-east-1")

	fmt.Println("============================================")
	fmt.Println("SCENARIO 4: LAMBDA + S3 INTEGRATION")
	fmt.Println("============================================")
	fmt.Printf("Region: %s\n", usecasescommon.FormatRegionName(region))
	fmt.Println("\n[MOCK MODE] All resources are simulated - no AWS SDK calls")

	// Initialize mock ID generator
	gen := usecasescommon.NewMockIDGenerator()

	// Step 1: Create S3 bucket for Lambda code and data
	fmt.Println("\n--- Step 1: Creating S3 Bucket ---")
	s3Bucket := usecasescommon.CreateMockS3Bucket("lambda-data-bucket-12345", region, gen)
	if err := s3Bucket.Validate(); err != nil {
		fmt.Printf("✗ S3 bucket validation failed: %v\n", err)
		return
	}
	fmt.Printf("✓ S3 Bucket created: %s\n", s3Bucket.Name)
	if s3Bucket.ARN != nil {
		fmt.Printf("  ARN: %s\n", *s3Bucket.ARN)
	}
	if s3Bucket.BucketDomainName != nil {
		fmt.Printf("  Domain: %s\n", *s3Bucket.BucketDomainName)
	}
	if s3Bucket.BucketRegionalDomainName != nil {
		fmt.Printf("  Regional Domain: %s\n", *s3Bucket.BucketRegionalDomainName)
	}

	// Step 2: Create IAM Role for Lambda function
	fmt.Println("\n--- Step 2: Creating IAM Role for Lambda ---")
	lambdaRole := usecasescommon.CreateMockIAMRoleForLambda(
		"lambda-s3-access-role",
		"IAM role for Lambda function to access S3 bucket",
		gen,
	)
	fmt.Printf("✓ IAM Role created: %s\n", lambdaRole.Name)
	if lambdaRole.ARN != nil {
		fmt.Printf("  ARN: %s\n", *lambdaRole.ARN)
	}
	fmt.Printf("  Trust Policy: Allows lambda.amazonaws.com to assume role\n")

	// Step 3: Create IAM Policy for S3 access
	fmt.Println("\n--- Step 3: Creating IAM Policy for S3 Access ---")
	s3PolicyDocument := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"s3:GetObject",
					"s3:PutObject",
					"s3:DeleteObject",
					"s3:ListBucket"
				],
				"Resource": [
					"arn:aws:s3:::%s",
					"arn:aws:s3:::%s/*"
				]
			}
		]
	}`, s3Bucket.Name, s3Bucket.Name)

	s3Policy := &domainiam.Policy{
		ID:             "lambda-s3-access-policy",
		Name:           "LambdaS3AccessPolicy",
		Description:    stringPtr("Policy allowing Lambda to access S3 bucket"),
		Path:           stringPtr("/"),
		PolicyDocument: s3PolicyDocument,
	}

	fmt.Printf("✓ IAM Policy created: %s\n", s3Policy.Name)
	fmt.Printf("  Permissions:\n")
	fmt.Printf("    - s3:GetObject (read objects from bucket)\n")
	fmt.Printf("    - s3:PutObject (write objects to bucket)\n")
	fmt.Printf("    - s3:DeleteObject (delete objects from bucket)\n")
	fmt.Printf("    - s3:ListBucket (list bucket contents)\n")
	fmt.Printf("  Resource: %s and %s/*\n", s3Bucket.Name, s3Bucket.Name)

	// Step 4: Attach policy to role (simulated)
	fmt.Println("\n--- Step 4: Attaching IAM Policy to Role ---")
	fmt.Printf("✓ Policy '%s' attached to role '%s'\n", s3Policy.Name, lambdaRole.Name)
	fmt.Printf("  Lambda function will now have permissions to:\n")
	fmt.Printf("    - Read from S3 bucket: %s\n", s3Bucket.Name)
	fmt.Printf("    - Write to S3 bucket: %s\n", s3Bucket.Name)
	fmt.Printf("    - Delete from S3 bucket: %s\n", s3Bucket.Name)
	fmt.Printf("    - List objects in S3 bucket: %s\n", s3Bucket.Name)

	// Step 5: Create Lambda function with S3 code deployment
	fmt.Println("\n--- Step 5: Creating Lambda Function ---")
	s3CodeBucket := "lambda-code-bucket-12345"
	s3CodeKey := "functions/data-processor.zip"
	runtime := "python3.9"
	handler := "index.handler"
	memorySize := 256
	timeout := 30

	lambdaFunction := usecasescommon.CreateMockLambdaFunction(
		"data-processor-function",
		*lambdaRole.ARN,
		region,
		&s3CodeBucket,
		&s3CodeKey,
		&runtime,
		&handler,
		&memorySize,
		&timeout,
		gen,
	)

	// Add environment variables
	lambdaFunction.Environment = map[string]string{
		"S3_BUCKET_NAME": s3Bucket.Name,
		"REGION":         region,
		"LOG_LEVEL":      "INFO",
	}

	// Add tags
	lambdaFunction.Tags = map[string]string{
		"Environment": "dev",
		"Service":     "data-processing",
		"ManagedBy":   "arch-visualizer",
	}

	if err := lambdaFunction.Validate(); err != nil {
		fmt.Printf("✗ Lambda function validation failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Lambda Function created: %s\n", lambdaFunction.FunctionName)
	if lambdaFunction.ARN != nil {
		fmt.Printf("  ARN: %s\n", *lambdaFunction.ARN)
	}
	if lambdaFunction.InvokeARN != nil {
		fmt.Printf("  Invoke ARN: %s\n", *lambdaFunction.InvokeARN)
	}
	fmt.Printf("  Runtime: %s\n", *lambdaFunction.Runtime)
	fmt.Printf("  Handler: %s\n", *lambdaFunction.Handler)
	fmt.Printf("  Memory: %d MB\n", *lambdaFunction.MemorySize)
	fmt.Printf("  Timeout: %d seconds\n", *lambdaFunction.Timeout)
	fmt.Printf("  Code Location: s3://%s/%s\n", s3CodeBucket, s3CodeKey)
	fmt.Printf("  IAM Role: %s\n", lambdaFunction.RoleARN)
	fmt.Printf("  Environment Variables:\n")
	for k, v := range lambdaFunction.Environment {
		fmt.Printf("    - %s: %s\n", k, v)
	}

	// Step 6: Display architecture summary
	fmt.Println("\n============================================")
	fmt.Println("ARCHITECTURE SUMMARY")
	fmt.Println("============================================")
	fmt.Printf("S3 Bucket: %s\n", s3Bucket.Name)
	fmt.Printf("  - Purpose: Data storage for Lambda function\n")
	fmt.Printf("  - Region: %s\n", s3Bucket.Region)
	fmt.Printf("\nIAM Role: %s\n", lambdaRole.Name)
	fmt.Printf("  - Purpose: Grants Lambda function permissions\n")
	fmt.Printf("  - Trust Policy: lambda.amazonaws.com\n")
	fmt.Printf("  - Attached Policies:\n")
	fmt.Printf("    - %s (S3 access)\n", s3Policy.Name)
	fmt.Printf("\nLambda Function: %s\n", lambdaFunction.FunctionName)
	fmt.Printf("  - Purpose: Process data from S3 bucket\n")
	fmt.Printf("  - Runtime: %s\n", *lambdaFunction.Runtime)
	fmt.Printf("  - Memory: %d MB\n", *lambdaFunction.MemorySize)
	fmt.Printf("  - IAM Role: %s\n", lambdaFunction.RoleARN)
	fmt.Printf("  - Environment: S3_BUCKET_NAME=%s\n", s3Bucket.Name)

	// Step 7: Display IAM policy details
	fmt.Println("\n============================================")
	fmt.Println("IAM POLICY DETAILS")
	fmt.Println("============================================")
	fmt.Printf("Policy Name: %s\n", s3Policy.Name)
	fmt.Printf("Policy Document:\n")
	fmt.Println(s3Policy.PolicyDocument)
	fmt.Printf("\nThis policy allows the Lambda function to:\n")
	fmt.Printf("  1. Read objects from S3 bucket (%s)\n", s3Bucket.Name)
	fmt.Printf("  2. Write objects to S3 bucket (%s)\n", s3Bucket.Name)
	fmt.Printf("  3. Delete objects from S3 bucket (%s)\n", s3Bucket.Name)
	fmt.Printf("  4. List objects in S3 bucket (%s)\n", s3Bucket.Name)

	// Step 8: Calculate estimated costs
	fmt.Println("\n============================================")
	fmt.Println("COST ESTIMATION (30 days)")
	fmt.Println("============================================")
	pricingService := awspricing.NewAWSPricingService()
	duration := 30 * 24 * time.Hour // 30 days

	totalCost := 0.0

	// S3 Bucket costs
	fmt.Println("\nS3 Bucket Costs:")
	fmt.Println("  Note: S3 storage costs depend on usage (size, requests, data transfer)")
	fmt.Println("  For this example, assuming:")
	fmt.Println("    - Storage: 10 GB")
	fmt.Println("    - PUT requests: 10,000/month")
	fmt.Println("    - GET requests: 50,000/month")
	fmt.Println("    - Data transfer: 5 GB/month")

	s3Resource := &resource.Resource{
		Type:     resource.ResourceType{Name: "s3_bucket"},
		Provider: "aws",
		Region:   region,
		Metadata: map[string]interface{}{
			"size_gb":          10.0,
			"storage_class":    "standard",
			"put_requests":     10000.0,
			"get_requests":     50000.0,
			"data_transfer_gb": 5.0,
		},
	}
	s3Estimate, err := pricingService.EstimateCost(ctx, s3Resource, duration)
	if err == nil {
		fmt.Printf("  S3 Bucket (%s): $%.2f\n", s3Bucket.Name, s3Estimate.TotalCost)
		totalCost += s3Estimate.TotalCost
	}

	// Lambda Function costs
	fmt.Println("\nLambda Function Costs:")
	fmt.Println("  Note: Lambda costs depend on:")
	fmt.Println("    - Memory: 256 MB")
	fmt.Println("    - Average duration: 200ms")
	fmt.Println("    - Requests: 1,000,000/month (within free tier)")
	fmt.Println("    - Data transfer: 2 GB/month")

	lambdaResource := &resource.Resource{
		Type:     resource.ResourceType{Name: "lambda_function"},
		Provider: "aws",
		Region:   region,
		Metadata: map[string]interface{}{
			"memory_size_mb":      256.0,
			"average_duration_ms": 200.0,
			"request_count":       1000000.0, // 1M requests (within free tier)
			"data_transfer_gb":    2.0,
		},
	}
	lambdaEstimate, err := pricingService.EstimateCost(ctx, lambdaResource, duration)
	if err == nil {
		fmt.Printf("  Lambda Function (%s): $%.2f\n", lambdaFunction.FunctionName, lambdaEstimate.TotalCost)
		totalCost += lambdaEstimate.TotalCost
		fmt.Println("\n  Cost Breakdown:")
		for i, component := range lambdaEstimate.Breakdown {
			fmt.Printf("    %d. %s: $%.6f\n", i+1, component.ComponentName, component.Subtotal)
		}
	}

	// IAM costs (free)
	fmt.Println("\nIAM Costs:")
	fmt.Printf("  IAM Role: $0.00 (free)\n")
	fmt.Printf("  IAM Policy: $0.00 (free)\n")

	// Total cost
	fmt.Println("\n============================================")
	fmt.Printf("TOTAL ESTIMATED COST (30 days): $%.2f\n", totalCost)
	fmt.Println("============================================")
	fmt.Println("\nNote:")
	fmt.Println("  - Costs are estimates based on provided usage metrics")
	fmt.Println("  - Actual costs may vary based on real usage patterns")
	fmt.Println("  - Lambda free tier includes 1M requests and 400,000 GB-seconds per month")
	fmt.Println("  - S3 free tier includes 5 GB storage and 20,000 GET requests per month")
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
