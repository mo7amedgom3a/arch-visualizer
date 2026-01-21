package lambda

import (
	"context"
	"fmt"
	"time"

	awscomputeadapter "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/adapters/compute"
	awsautoscaling "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling"
	awsautoscalingoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling/outputs"
	awsec2 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2"
	awslttemplate "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template"
	awslttemplateoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/launch_template/outputs"
	awsec2outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2/outputs"
	awsmodel "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/instance_types"
	awslambda "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda"
	awslambdaoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/lambda/outputs"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awsloadbalanceroutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer/outputs"
	awspricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	domainresource "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
)

// inMemoryAWSComputeService is an in-memory implementation of AWSComputeService.
// It is used by the LambdaRunner to simulate Lambda function operations without calling real AWS.
type inMemoryAWSComputeService struct {
	// Lambda functions state
	functions map[string]*awslambdaoutputs.FunctionOutput
}

// Ensure inMemoryAWSComputeService implements AWSComputeService
var _ awsservice.AWSComputeService = (*inMemoryAWSComputeService)(nil)

// newInMemoryAWSComputeService creates a new in-memory service instance.
func newInMemoryAWSComputeService() *inMemoryAWSComputeService {
	return &inMemoryAWSComputeService{
		functions: make(map[string]*awslambdaoutputs.FunctionOutput),
	}
}

// newLambdaDemoAdapter constructs a domain ComputeService backed by the in-memory AWS compute service.
func newLambdaDemoAdapter() domaincompute.ComputeService {
	service := newInMemoryAWSComputeService()
	return awscomputeadapter.NewAWSComputeAdapter(service)
}

// =========================
// Lambda Function methods
// =========================

func (s *inMemoryAWSComputeService) CreateLambdaFunction(ctx context.Context, function *awslambda.Function) (*awslambdaoutputs.FunctionOutput, error) {
	if function == nil {
		return nil, fmt.Errorf("function is nil")
	}

	arn := fmt.Sprintf("arn:aws:lambda:us-east-1:123456789012:function:%s", function.FunctionName)
	invokeARN := fmt.Sprintf("arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/%s/invocations", arn)
	qualifiedARN := fmt.Sprintf("%s:$LATEST", arn)

	out := &awslambdaoutputs.FunctionOutput{
		ARN:             arn,
		InvokeARN:       invokeARN,
		QualifiedARN:    &qualifiedARN,
		Version:         "$LATEST",
		FunctionName:    function.FunctionName,
		Region:          "us-east-1",
		RoleARN:         function.RoleARN,
		S3Bucket:        function.S3Bucket,
		S3Key:           function.S3Key,
		S3ObjectVersion: function.S3ObjectVersion,
		PackageType:     function.PackageType,
		ImageURI:        function.ImageURI,
		Runtime:         function.Runtime,
		Handler:         function.Handler,
		MemorySize:      function.MemorySize,
		Timeout:         function.Timeout,
		Environment:     function.Environment,
		Layers:          function.Layers,
		LastModified:    stringPtr(time.Now().Format(time.RFC3339)),
		CodeSize:        int64Ptr(1024 * 1024), // 1MB
		CodeSHA256:      stringPtr("abc123def456..."),
	}

	if function.VPCConfig != nil {
		out.VPCConfig = &awslambdaoutputs.FunctionVPCConfigOutput{
			SubnetIDs:        function.VPCConfig.SubnetIDs,
			SecurityGroupIDs: function.VPCConfig.SecurityGroupIDs,
			VPCID:            stringPtr("vpc-12345678"),
		}
	}

	s.functions[function.FunctionName] = out
	return out, nil
}

func (s *inMemoryAWSComputeService) GetLambdaFunction(ctx context.Context, name string) (*awslambdaoutputs.FunctionOutput, error) {
	if f, ok := s.functions[name]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("lambda function %s not found", name)
}

func (s *inMemoryAWSComputeService) UpdateLambdaFunction(ctx context.Context, name string, function *awslambda.Function) (*awslambdaoutputs.FunctionOutput, error) {
	// For simplicity, reuse CreateLambdaFunction logic and overwrite
	out, err := s.CreateLambdaFunction(ctx, function)
	if err != nil {
		return nil, err
	}
	s.functions[name] = out
	return out, nil
}

func (s *inMemoryAWSComputeService) DeleteLambdaFunction(ctx context.Context, name string) error {
	delete(s.functions, name)
	return nil
}

func (s *inMemoryAWSComputeService) ListLambdaFunctions(ctx context.Context, filters map[string][]string) ([]*awslambdaoutputs.FunctionOutput, error) {
	results := make([]*awslambdaoutputs.FunctionOutput, 0, len(s.functions))
	for _, f := range s.functions {
		results = append(results, f)
	}
	return results, nil
}

// Stub implementations for other AWSComputeService methods (required by interface)
func (s *inMemoryAWSComputeService) CreateInstance(ctx context.Context, instance *awsec2.Instance) (*awsec2outputs.InstanceOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) GetInstance(ctx context.Context, id string) (*awsec2outputs.InstanceOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) UpdateInstance(ctx context.Context, instance *awsec2.Instance) (*awsec2outputs.InstanceOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) DeleteInstance(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) ListInstances(ctx context.Context, filters map[string][]string) ([]*awsec2outputs.InstanceOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) StartInstance(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) StopInstance(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) RebootInstance(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) GetInstanceTypeInfo(ctx context.Context, instanceType string, region string) (*awsmodel.InstanceTypeInfo, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) ListInstanceTypesByCategory(ctx context.Context, category awsmodel.InstanceCategory, region string) ([]*awsmodel.InstanceTypeInfo, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) CreateLaunchTemplate(ctx context.Context, template *awslttemplate.LaunchTemplate) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) GetLaunchTemplate(ctx context.Context, id string) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) UpdateLaunchTemplate(ctx context.Context, id string, template *awslttemplate.LaunchTemplate) (*awslttemplateoutputs.LaunchTemplateOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) DeleteLaunchTemplate(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) ListLaunchTemplates(ctx context.Context, filters map[string][]string) ([]*awslttemplateoutputs.LaunchTemplateOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) GetLaunchTemplateVersion(ctx context.Context, id string, version int) (*awslttemplate.LaunchTemplateVersion, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) ListLaunchTemplateVersions(ctx context.Context, id string) ([]*awslttemplate.LaunchTemplateVersion, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) CreateLoadBalancer(ctx context.Context, lb *awsloadbalancer.LoadBalancer) (*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) GetLoadBalancer(ctx context.Context, arn string) (*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) UpdateLoadBalancer(ctx context.Context, arn string, lb *awsloadbalancer.LoadBalancer) (*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) DeleteLoadBalancer(ctx context.Context, arn string) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) ListLoadBalancers(ctx context.Context, filters map[string][]string) ([]*awsloadbalanceroutputs.LoadBalancerOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) CreateTargetGroup(ctx context.Context, tg *awsloadbalancer.TargetGroup) (*awsloadbalanceroutputs.TargetGroupOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) GetTargetGroup(ctx context.Context, arn string) (*awsloadbalanceroutputs.TargetGroupOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) UpdateTargetGroup(ctx context.Context, arn string, tg *awsloadbalancer.TargetGroup) (*awsloadbalanceroutputs.TargetGroupOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) DeleteTargetGroup(ctx context.Context, arn string) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) ListTargetGroups(ctx context.Context, filters map[string][]string) ([]*awsloadbalanceroutputs.TargetGroupOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) CreateListener(ctx context.Context, listener *awsloadbalancer.Listener) (*awsloadbalanceroutputs.ListenerOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) GetListener(ctx context.Context, arn string) (*awsloadbalanceroutputs.ListenerOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) UpdateListener(ctx context.Context, arn string, listener *awsloadbalancer.Listener) (*awsloadbalanceroutputs.ListenerOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) DeleteListener(ctx context.Context, arn string) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) ListListeners(ctx context.Context, loadBalancerARN string) ([]*awsloadbalanceroutputs.ListenerOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) AttachTargetToGroup(ctx context.Context, attachment *awsloadbalancer.TargetGroupAttachment) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) DetachTargetFromGroup(ctx context.Context, targetGroupARN, targetID string) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) ListTargetGroupTargets(ctx context.Context, targetGroupARN string) ([]*awsloadbalanceroutputs.TargetGroupAttachmentOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) CreateAutoScalingGroup(ctx context.Context, asg *awsautoscaling.AutoScalingGroup) (*awsautoscalingoutputs.AutoScalingGroupOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) GetAutoScalingGroup(ctx context.Context, name string) (*awsautoscalingoutputs.AutoScalingGroupOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) UpdateAutoScalingGroup(ctx context.Context, name string, asg *awsautoscaling.AutoScalingGroup) (*awsautoscalingoutputs.AutoScalingGroupOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) DeleteAutoScalingGroup(ctx context.Context, name string) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) ListAutoScalingGroups(ctx context.Context, filters map[string][]string) ([]*awsautoscalingoutputs.AutoScalingGroupOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) SetDesiredCapacity(ctx context.Context, asgName string, capacity int) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) AttachInstances(ctx context.Context, asgName string, instanceIDs []string) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) DetachInstances(ctx context.Context, asgName string, instanceIDs []string) error {
	return fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) PutScalingPolicy(ctx context.Context, policy *awsautoscaling.ScalingPolicy) (*awsautoscalingoutputs.ScalingPolicyOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) DescribeScalingPolicies(ctx context.Context, asgName string) ([]*awsautoscalingoutputs.ScalingPolicyOutput, error) {
	return nil, fmt.Errorf("not implemented")
}
func (s *inMemoryAWSComputeService) DeleteScalingPolicy(ctx context.Context, policyName, asgName string) error {
	return fmt.Errorf("not implemented")
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

// LambdaRunner demonstrates Lambda function operations using the domain compute service and in-memory AWS implementation.
// This is intended to be called from a main() or manual test harness.
func LambdaRunner() {
	ctx := context.Background()
	computeService := newLambdaDemoAdapter()

	fmt.Println("============================================")
	fmt.Println("LAMBDA FUNCTION DEMO (IN-MEMORY)")
	fmt.Println("============================================")

	// 1. Create Lambda function with S3 code deployment
	fmt.Println("\n--- Creating Lambda Function (S3 Code) ---")
	s3Bucket := "my-lambda-code-bucket"
	s3Key := "functions/my-function.zip"
	runtime := "python3.9"
	handler := "index.handler"
	memorySize := 256
	timeout := 30

	function := &domaincompute.LambdaFunction{
		FunctionName: "my-test-function",
		RoleARN:      "arn:aws:iam::123456789012:role/lambda-execution-role",
		Region:       "us-east-1",
		S3Bucket:     &s3Bucket,
		S3Key:        &s3Key,
		Runtime:      &runtime,
		Handler:      &handler,
		MemorySize:   &memorySize,
		Timeout:      &timeout,
		Environment: map[string]string{
			"ENV":       "dev",
			"LOG_LEVEL": "INFO",
		},
		Tags: map[string]string{
			"Environment": "dev",
			"Service":     "lambda-demo",
		},
	}

	if err := function.Validate(); err != nil {
		fmt.Printf("Function validation error: %v\n", err)
		return
	}

	createdFunction, err := computeService.CreateLambdaFunction(ctx, function)
	if err != nil {
		fmt.Printf("CreateLambdaFunction error: %v\n", err)
		return
	}

	fmt.Printf("Function Name: %s\n", createdFunction.FunctionName)
	if createdFunction.ARN != nil {
		fmt.Printf("ARN: %s\n", *createdFunction.ARN)
	}
	if createdFunction.InvokeARN != nil {
		fmt.Printf("Invoke ARN: %s\n", *createdFunction.InvokeARN)
	}
	if createdFunction.Runtime != nil {
		fmt.Printf("Runtime: %s\n", *createdFunction.Runtime)
	}
	if createdFunction.Handler != nil {
		fmt.Printf("Handler: %s\n", *createdFunction.Handler)
	}
	if createdFunction.MemorySize != nil {
		fmt.Printf("Memory: %d MB\n", *createdFunction.MemorySize)
	}
	if createdFunction.Timeout != nil {
		fmt.Printf("Timeout: %d seconds\n", *createdFunction.Timeout)
	}

	// 2. Create Lambda function with container image
	fmt.Println("\n--- Creating Lambda Function (Container Image) ---")
	packageType := "Image"
	imageURI := "123456789012.dkr.ecr.us-east-1.amazonaws.com/my-lambda:latest"
	memorySize2 := 512
	timeout2 := 60

	function2 := &domaincompute.LambdaFunction{
		FunctionName: "my-container-function",
		RoleARN:      "arn:aws:iam::123456789012:role/lambda-execution-role",
		Region:       "us-east-1",
		PackageType:  &packageType,
		ImageURI:     &imageURI,
		MemorySize:   &memorySize2,
		Timeout:      &timeout2,
		Tags: map[string]string{
			"Environment": "prod",
			"Service":     "lambda-container",
		},
	}

	if err := function2.Validate(); err != nil {
		fmt.Printf("Function validation error: %v\n", err)
		return
	}

	createdFunction2, err := computeService.CreateLambdaFunction(ctx, function2)
	if err != nil {
		fmt.Printf("CreateLambdaFunction error: %v\n", err)
		return
	}

	fmt.Printf("Function Name: %s\n", createdFunction2.FunctionName)
	if createdFunction2.ARN != nil {
		fmt.Printf("ARN: %s\n", *createdFunction2.ARN)
	}
	if createdFunction2.PackageType != nil {
		fmt.Printf("Package Type: %s\n", *createdFunction2.PackageType)
	}
	if createdFunction2.ImageURI != nil {
		fmt.Printf("Image URI: %s\n", *createdFunction2.ImageURI)
	}

	// 3. Get Lambda function
	fmt.Println("\n--- Getting Lambda Function ---")
	retrievedFunction, err := computeService.GetLambdaFunction(ctx, "my-test-function")
	if err != nil {
		fmt.Printf("GetLambdaFunction error: %v\n", err)
		return
	}
	fmt.Printf("Retrieved Function: %s\n", retrievedFunction.FunctionName)
	if retrievedFunction.ARN != nil {
		fmt.Printf("ARN: %s\n", *retrievedFunction.ARN)
	}

	// 4. Update Lambda function
	fmt.Println("\n--- Updating Lambda Function ---")
	newMemorySize := 512
	newTimeout := 60
	updatedFunction := &domaincompute.LambdaFunction{
		FunctionName: "my-test-function",
		RoleARN:      "arn:aws:iam::123456789012:role/lambda-execution-role",
		Region:       "us-east-1",
		S3Bucket:     &s3Bucket,
		S3Key:        &s3Key,
		Runtime:      &runtime,
		Handler:      &handler,
		MemorySize:   &newMemorySize,
		Timeout:      &newTimeout,
		Environment: map[string]string{
			"ENV":       "prod",
			"LOG_LEVEL": "DEBUG",
		},
	}

	updated, err := computeService.UpdateLambdaFunction(ctx, updatedFunction)
	if err != nil {
		fmt.Printf("UpdateLambdaFunction error: %v\n", err)
		return
	}
	fmt.Printf("Updated Function: %s\n", updated.FunctionName)
	if updated.MemorySize != nil {
		fmt.Printf("New Memory: %d MB\n", *updated.MemorySize)
	}
	if updated.Timeout != nil {
		fmt.Printf("New Timeout: %d seconds\n", *updated.Timeout)
	}

	// 5. List Lambda functions
	fmt.Println("\n--- Listing Lambda Functions ---")
	functions, err := computeService.ListLambdaFunctions(ctx, map[string]string{})
	if err != nil {
		fmt.Printf("ListLambdaFunctions error: %v\n", err)
		return
	}
	fmt.Printf("Found %d functions:\n", len(functions))
	for _, f := range functions {
		fmt.Printf("- %s", f.FunctionName)
		if f.ARN != nil {
			fmt.Printf(" (ARN: %s)", *f.ARN)
		}
		fmt.Println()
	}
}

// LambdaPricingRunner demonstrates Lambda pricing calculations
func LambdaPricingRunner() {
	ctx := context.Background()

	fmt.Println("============================================")
	fmt.Println("LAMBDA PRICING DEMO")
	fmt.Println("============================================")

	pricingService := awspricing.NewAWSPricingService()

	// Test 1: Get Lambda pricing information
	fmt.Println("\n--- Getting Lambda Pricing Information ---")
	pricing, err := pricingService.GetPricing(ctx, "lambda_function", "aws", "us-east-1")
	if err != nil {
		fmt.Printf("Error getting pricing: %v\n", err)
		return
	}

	fmt.Printf("Resource Type: %s\n", pricing.ResourceType)
	fmt.Printf("Provider: %s\n", pricing.Provider)
	fmt.Println("\nPricing Components:")
	for i, component := range pricing.Components {
		fmt.Printf("  %d. %s\n", i+1, component.Name)
		fmt.Printf("     Model: %s\n", component.Model)
		fmt.Printf("     Unit: %s\n", component.Unit)
		fmt.Printf("     Rate: $%.10f\n", component.Rate)
		fmt.Printf("     Description: %s\n", component.Description)
		fmt.Println()
	}

	// Test 2: Calculate cost for a Lambda function
	fmt.Println("--- Calculating Lambda Function Cost ---")

	// Example: 128 MB memory, 100ms average duration, 1M requests/month, 10 GB data transfer
	memorySizeMB := 128.0
	averageDurationMs := 100.0
	requestCount := 1000000.0 // 1M requests
	dataTransferGB := 10.0

	lambdaResource := &domainresource.Resource{
		Type: domainresource.ResourceType{
			Name: "lambda_function",
		},
		Provider: "aws",
		Region:   "us-east-1",
		Metadata: map[string]interface{}{
			"memory_size_mb":      memorySizeMB,
			"average_duration_ms": averageDurationMs,
			"request_count":       requestCount,
			"data_transfer_gb":    dataTransferGB,
		},
	}

	duration := 720 * time.Hour // 1 month
	estimate, err := pricingService.EstimateCost(ctx, lambdaResource, duration)
	if err != nil {
		fmt.Printf("Error estimating cost: %v\n", err)
		return
	}

	fmt.Printf("\nCost Estimate:\n")
	fmt.Printf("  Total Cost: $%.6f\n", estimate.TotalCost)
	fmt.Printf("  Currency: %s\n", estimate.Currency)
	fmt.Printf("  Period: %s\n", estimate.Period)
	fmt.Printf("  Duration: %s\n", estimate.Duration)
	fmt.Println("\nCost Breakdown:")
	for i, component := range estimate.Breakdown {
		fmt.Printf("  %d. %s\n", i+1, component.ComponentName)
		fmt.Printf("     Model: %s\n", component.Model)
		fmt.Printf("     Quantity: %.2f\n", component.Quantity)
		fmt.Printf("     Unit Rate: $%.10f\n", component.UnitRate)
		fmt.Printf("     Subtotal: $%.6f\n", component.Subtotal)
		fmt.Println()
	}

	// Test 3: Calculate cost with requests over free tier
	fmt.Println("--- Calculating Lambda Cost (Over Free Tier) ---")

	requestCount2 := 5000000.0 // 5M requests, first 1M free
	lambdaResource2 := &domainresource.Resource{
		Type: domainresource.ResourceType{
			Name: "lambda_function",
		},
		Provider: "aws",
		Region:   "us-east-1",
		Metadata: map[string]interface{}{
			"memory_size_mb":      512.0,
			"average_duration_ms": 300.0,
			"request_count":       requestCount2,
			"data_transfer_gb":    20.0,
		},
	}

	estimate2, err := pricingService.EstimateCost(ctx, lambdaResource2, duration)
	if err != nil {
		fmt.Printf("Error estimating cost: %v\n", err)
		return
	}

	fmt.Printf("\nCost Estimate (High Usage):\n")
	fmt.Printf("  Total Cost: $%.6f\n", estimate2.TotalCost)
	fmt.Println("\nCost Breakdown:")
	for i, component := range estimate2.Breakdown {
		fmt.Printf("  %d. %s\n", i+1, component.ComponentName)
		fmt.Printf("     Quantity: %.2f\n", component.Quantity)
		fmt.Printf("     Unit Rate: $%.10f\n", component.UnitRate)
		fmt.Printf("     Subtotal: $%.6f\n", component.Subtotal)
		fmt.Println()
	}
}
