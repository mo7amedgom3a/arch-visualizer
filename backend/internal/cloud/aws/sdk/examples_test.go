package sdk_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/joho/godotenv"
	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
)

func init() {
	// Load .env file from project root
	// The test file is in backend/internal/cloud/aws/sdk/
	// So we need to go up 4 levels to reach the project root (backend/)
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	// Try to find .env file by going up directories
	// Try multiple paths in case the working directory varies
	envPaths := []string{
		filepath.Join(wd, "..", "..", "..", "..", ".env"), // From sdk/ to backend/
		filepath.Join(wd, ".env"),                         // Current directory
		filepath.Join(wd, "..", ".env"),                   // Parent
		filepath.Join(wd, "..", "..", ".env"),             // 2 levels up
		filepath.Join(wd, "..", "..", "..", ".env"),       // 3 levels up
		filepath.Join(wd, "..", "..", "..", "..", ".env"), // 4 levels up
	}

	// Try loading from each path
	for _, envPath := range envPaths {
		if err := godotenv.Load(envPath); err == nil {
			return // Successfully loaded
		}
	}

	// If all paths failed, try loading from current directory (godotenv default behavior)
	_ = godotenv.Load()
}

// ExampleNewAWSClient demonstrates how to create an AWS client
// This example requires AWS credentials in environment variables or .env file
func ExampleNewAWSClient() {
	ctx := context.Background()

	// Create AWS client (reads from environment variables)
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		fmt.Printf("Error creating AWS client: %v\n", err)
		return
	}

	fmt.Printf("AWS Client initialized for region: %s\n", client.GetRegion())
}

// TestAWSClientInitialization tests AWS client initialization
// This test requires AWS credentials to be set in environment variables
func TestAWSClientInitialization(t *testing.T) {
	// Skip if credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping test: AWS_ACCESS_KEY_ID not set")
	}

	ctx := context.Background()

	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	if client == nil {
		t.Fatal("AWS client is nil")
	}

	if client.EC2 == nil {
		t.Fatal("EC2 client is nil")
	}

	region := client.GetRegion()
	if region == "" {
		t.Fatal("Region is empty")
	}

	t.Logf("AWS Client initialized successfully for region: %s", region)
}

// TestEC2DescribeInstances tests EC2 DescribeInstances operation
// This test requires AWS credentials and will make actual API calls
func TestEC2DescribeInstances(t *testing.T) {
	// Skip if credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping test: AWS_ACCESS_KEY_ID not set")
	}

	ctx := context.Background()

	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	// Describe instances
	input := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int32(5), // Limit to 5 results for testing
	}

	output, err := client.EC2.DescribeInstances(ctx, input)
	if err != nil {
		t.Fatalf("Failed to describe instances: %v", err)
	}

	t.Logf("Found %d reservations", len(output.Reservations))

	instanceCount := 0
	for _, reservation := range output.Reservations {
		instanceCount += len(reservation.Instances)
		for _, instance := range reservation.Instances {
			t.Logf("Instance ID: %s, State: %s, Type: %s",
				aws.ToString(instance.InstanceId),
				instance.State.Name,
				instance.InstanceType)
		}
	}

	t.Logf("Total instances found: %d", instanceCount)
}

// TestEC2DescribeRegions tests EC2 DescribeRegions operation
// This is a safe operation that doesn't require any resources
func TestEC2DescribeRegions(t *testing.T) {
	// Skip if credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping test: AWS_ACCESS_KEY_ID not set")
	}

	ctx := context.Background()

	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	// Describe regions
	input := &ec2.DescribeRegionsInput{}

	output, err := client.EC2.DescribeRegions(ctx, input)
	if err != nil {
		t.Fatalf("Failed to describe regions: %v", err)
	}

	if len(output.Regions) == 0 {
		t.Fatal("No regions found")
	}

	t.Logf("Found %d regions", len(output.Regions))
	for _, region := range output.Regions {
		t.Logf("Region: %s, Endpoint: %s",
			aws.ToString(region.RegionName),
			aws.ToString(region.Endpoint))
	}
}

// TestEC2DescribeInstanceTypes tests EC2 DescribeInstanceTypes operation
// This is a safe operation that lists available instance types
func TestEC2DescribeInstanceTypes(t *testing.T) {
	// Skip if credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping test: AWS_ACCESS_KEY_ID not set")
	}

	ctx := context.Background()

	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	// Describe instance types (limit to t3 family for testing)
	input := &ec2.DescribeInstanceTypesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-type"),
				Values: []string{"t3.*"},
			},
		},
		MaxResults: aws.Int32(10),
	}

	output, err := client.EC2.DescribeInstanceTypes(ctx, input)
	if err != nil {
		t.Fatalf("Failed to describe instance types: %v", err)
	}

	t.Logf("Found %d instance types", len(output.InstanceTypes))
	for _, instanceType := range output.InstanceTypes {
		t.Logf("Instance Type: %s",
			string(instanceType.InstanceType))
	}
}

// TestEC2DescribeVPCs tests EC2 DescribeVPCs operation
func TestEC2DescribeVPCs(t *testing.T) {
	// Skip if credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping test: AWS_ACCESS_KEY_ID not set")
	}

	ctx := context.Background()

	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	// Describe VPCs
	input := &ec2.DescribeVpcsInput{
		MaxResults: aws.Int32(10),
	}

	output, err := client.EC2.DescribeVpcs(ctx, input)
	if err != nil {
		t.Fatalf("Failed to describe VPCs: %v", err)
	}

	t.Logf("Found %d VPCs", len(output.Vpcs))
	for _, vpc := range output.Vpcs {
		t.Logf("VPC ID: %s, CIDR: %s, State: %s",
			aws.ToString(vpc.VpcId),
			aws.ToString(vpc.CidrBlock),
			vpc.State)
	}
}

// TestEC2DescribeSubnets tests EC2 DescribeSubnets operation
func TestEC2DescribeSubnets(t *testing.T) {
	// Skip if credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping test: AWS_ACCESS_KEY_ID not set")
	}

	ctx := context.Background()

	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	// Describe subnets
	input := &ec2.DescribeSubnetsInput{
		MaxResults: aws.Int32(10),
	}

	output, err := client.EC2.DescribeSubnets(ctx, input)
	if err != nil {
		t.Fatalf("Failed to describe subnets: %v", err)
	}

	t.Logf("Found %d subnets", len(output.Subnets))
	for _, subnet := range output.Subnets {
		t.Logf("Subnet ID: %s, VPC: %s, CIDR: %s, AZ: %s",
			aws.ToString(subnet.SubnetId),
			aws.ToString(subnet.VpcId),
			aws.ToString(subnet.CidrBlock),
			aws.ToString(subnet.AvailabilityZone))
	}
}

// TestEC2DescribeSecurityGroups tests EC2 DescribeSecurityGroups operation
func TestEC2DescribeSecurityGroups(t *testing.T) {
	// Skip if credentials are not available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("Skipping test: AWS_ACCESS_KEY_ID not set")
	}

	ctx := context.Background()

	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	// Describe security groups
	input := &ec2.DescribeSecurityGroupsInput{
		MaxResults: aws.Int32(10),
	}

	output, err := client.EC2.DescribeSecurityGroups(ctx, input)
	if err != nil {
		t.Fatalf("Failed to describe security groups: %v", err)
	}

	t.Logf("Found %d security groups", len(output.SecurityGroups))
	for _, sg := range output.SecurityGroups {
		t.Logf("Security Group ID: %s, VPC: %s, Name: %s",
			aws.ToString(sg.GroupId),
			aws.ToString(sg.VpcId),
			aws.ToString(sg.GroupName))
	}
}
