package iam

import (
	"context"
	"fmt"

	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
)

// FetchAndSavePolicies fetches AWS managed policies from SDK and saves them to service-specific JSON files
// This is a utility function that can be called to refresh the static policy data
func FetchAndSavePolicies() {
	ctx := context.Background()

	// Initialize AWS client
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		fmt.Printf("Error creating AWS client: %v\n", err)
		fmt.Println("Note: This function requires AWS credentials to fetch policies from SDK")
		return
	}

	fmt.Println("============================================")
	fmt.Println("FETCHING AND SAVING AWS MANAGED POLICIES")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("This will:")
	fmt.Println("  1. Fetch all AWS managed policies from AWS SDK")
	fmt.Println("  2. Categorize them by service (EC2, VPC, EBS, IAM, General)")
	fmt.Println("  3. Save them to service-specific JSON files")
	fmt.Println()
	fmt.Println("Note: This may take a while as it fetches policy documents for all policies...")
	fmt.Println()

	// Fetch and save policies
	if err := awssdk.FetchAndSaveServicePolicies(ctx, client); err != nil {
		fmt.Printf("Error fetching and saving policies: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("============================================")
	fmt.Println("SUCCESS: Policies saved to service-specific files")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("Files created:")
	fmt.Println("  - backend/internal/cloud/aws/models/iam/data/ec2/policies.json")
	fmt.Println("  - backend/internal/cloud/aws/models/iam/data/vpc/policies.json")
	fmt.Println("  - backend/internal/cloud/aws/models/iam/data/ebs/policies.json")
	fmt.Println("  - backend/internal/cloud/aws/models/iam/data/s3/polices.json")
	fmt.Println("  - backend/internal/cloud/aws/models/iam/data/iam/policies.json")
	fmt.Println("  - backend/internal/cloud/aws/models/iam/data/general/policies.json")
	fmt.Println()
}

// FetchAndSaveS3Policies fetches S3-related AWS managed policies from SDK and saves them to s3/polices.json
// This is a utility function that can be called to refresh the S3 policy data
func FetchAndSaveS3Policies() {
	ctx := context.Background()

	// Initialize AWS client
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		fmt.Printf("Error creating AWS client: %v\n", err)
		fmt.Println("Note: This function requires AWS credentials to fetch policies from SDK")
		fmt.Println("Make sure you have set:")
		fmt.Println("  - AWS_ACCESS_KEY_ID")
		fmt.Println("  - AWS_SECRET_ACCESS_KEY")
		fmt.Println("  - AWS_REGION (optional, defaults to us-east-1)")
		return
	}

	fmt.Println("============================================")
	fmt.Println("FETCHING AND SAVING S3 POLICIES")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("This will:")
	fmt.Println("  1. Fetch all AWS managed policies from AWS SDK")
	fmt.Println("  2. Filter S3-related policies (AmazonS3, S3, SimpleStorageService)")
	fmt.Println("  3. Fetch policy documents for each S3 policy")
	fmt.Println("  4. Save them to s3/polices.json")
	fmt.Println()
	fmt.Println("Note: This may take a while as it fetches policy documents...")
	fmt.Println()

	// Fetch and save S3 policies
	if err := awssdk.FetchAndSaveS3Policies(ctx, client); err != nil {
		fmt.Printf("Error fetching and saving S3 policies: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("============================================")
	fmt.Println("SUCCESS: S3 policies saved")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("File created:")
	fmt.Println("  - backend/internal/cloud/aws/models/iam/data/s3/polices.json")
	fmt.Println()
}
