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
	fmt.Println("  - backend/internal/cloud/aws/models/iam/data/iam/policies.json")
	fmt.Println("  - backend/internal/cloud/aws/models/iam/data/general/policies.json")
	fmt.Println()
}
