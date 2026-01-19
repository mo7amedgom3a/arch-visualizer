package iam

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
)

// IAMPoliciesRunner demonstrates IAM policy operations using AWS SDK
func IAMPoliciesRunner() {
	ctx := context.Background()

	// Initialize AWS client
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		fmt.Printf("Error creating AWS client: %v\n", err)
		return
	}

	fmt.Println("============================================")
	fmt.Println("IAM POLICIES OPERATIONS")
	fmt.Println("============================================")

	region := client.GetRegion()
	fmt.Printf("\nRegion: %s\n", region)
	fmt.Println("Note: IAM is a global service, but using region for consistency")
	fmt.Println()

	// 1. List AWS Managed Policies
	fmt.Println("--- 1. Listing AWS Managed Policies ---")
	listAWSManagedPolicies(ctx, client)

	// 2. Get a specific AWS Managed Policy
	fmt.Println("\n--- 2. Getting AWS Managed Policy Details ---")
	getAWSManagedPolicyExample(ctx, client, "arn:aws:iam::aws:policy/ReadOnlyAccess")

	// 3. List Customer Managed Policies
	fmt.Println("\n--- 3. Listing Customer Managed Policies ---")
	listCustomerPolicies(ctx, client)

	// 4. Create a Customer Managed Policy (example - will be cleaned up)
	fmt.Println("\n--- 4. Creating Customer Managed Policy (Example) ---")
	policyARN := createPolicyExample(ctx, client)

	// 5. Get the created policy
	if policyARN != "" {
		fmt.Println("\n--- 5. Getting Created Policy ---")
		getPolicyExample(ctx, client, policyARN)

		// 6. Update the policy
		fmt.Println("\n--- 6. Updating Policy (Creating New Version) ---")
		updatePolicyExample(ctx, client, policyARN)

		// 7. Cleanup: Delete the policy
		fmt.Println("\n--- 7. Cleaning Up: Deleting Policy ---")
		deletePolicyExample(ctx, client, policyARN)
	}

	// 8. Display Policy Information
	fmt.Println("\n--- 8. Policy Type Information ---")
	displayPolicyTypes()
}

// listAWSManagedPolicies lists AWS managed policies
func listAWSManagedPolicies(ctx context.Context, client *awssdk.AWSClient) {
	policies, err := awssdk.ListAWSManagedPolicies(ctx, client, nil)
	if err != nil {
		fmt.Printf("Error listing AWS managed policies: %v\n", err)
		return
	}

	fmt.Printf("Found %d AWS managed policies\n", len(policies))
	if len(policies) > 0 {
		fmt.Println("\nFirst 10 AWS Managed Policies:")
		maxDisplay := 10
		if len(policies) < maxDisplay {
			maxDisplay = len(policies)
		}
		for i := 0; i < maxDisplay; i++ {
			policy := policies[i]
			fmt.Printf("  %d. %s\n", i+1, policy.Name)
			fmt.Printf("     ARN: %s\n", policy.ARN)
			if policy.Description != nil {
				fmt.Printf("     Description: %s\n", *policy.Description)
			}
			fmt.Printf("     Path: %s\n", policy.Path)
			fmt.Printf("     Attachments: %d\n", policy.AttachmentCount)
			fmt.Println()
		}
	}
}

// getAWSManagedPolicyExample retrieves details of an AWS managed policy
func getAWSManagedPolicyExample(ctx context.Context, client *awssdk.AWSClient, arn string) {
	policy, err := awssdk.GetAWSManagedPolicy(ctx, client, arn)
	if err != nil {
		fmt.Printf("Error getting AWS managed policy: %v\n", err)
		return
	}

	fmt.Printf("Policy Name: %s\n", policy.Name)
	fmt.Printf("ARN: %s\n", policy.ARN)
	if policy.Description != nil {
		fmt.Printf("Description: %s\n", *policy.Description)
	}
	fmt.Printf("Path: %s\n", policy.Path)
	fmt.Printf("Is AWS Managed: %v\n", policy.IsAWSManaged)
	fmt.Printf("Is Attachable: %v\n", policy.IsAttachable)
	fmt.Printf("Attachment Count: %d\n", policy.AttachmentCount)
	fmt.Printf("Created: %s\n", policy.CreateDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", policy.UpdateDate.Format("2006-01-02 15:04:05"))

	// Display policy document (truncated if too long)
	if policy.PolicyDocument != "" {
		doc := policy.PolicyDocument
		if len(doc) > 500 {
			doc = doc[:500] + "... (truncated)"
		}
		fmt.Printf("Policy Document:\n%s\n", doc)
	}
}

// listCustomerPolicies lists customer managed policies
func listCustomerPolicies(ctx context.Context, client *awssdk.AWSClient) {
	// Use Local scope for customer managed policies
	// Don't fetch policy documents for listing (faster performance)
	policies, err := awssdk.ListPolicies(ctx, client, nil, types.PolicyScopeTypeLocal, false)
	if err != nil {
		fmt.Printf("Error listing customer policies: %v\n", err)
		return
	}

	fmt.Printf("Found %d customer managed policies\n", len(policies))
	if len(policies) > 0 {
		fmt.Println("\nCustomer Managed Policies:")
		for i, policy := range policies {
			fmt.Printf("  %d. %s\n", i+1, policy.Name)
			fmt.Printf("     ARN: %s\n", policy.ARN)
			if policy.Description != nil {
				fmt.Printf("     Description: %s\n", *policy.Description)
			}
			fmt.Printf("     Path: %s\n", policy.Path)
			fmt.Printf("     Attachments: %d\n", policy.AttachmentCount)
			fmt.Println()
		}
	} else {
		fmt.Println("No customer managed policies found")
	}
}

// createPolicyExample creates a sample customer managed policy
func createPolicyExample(ctx context.Context, client *awssdk.AWSClient) string {
	policy := &awsiam.Policy{
		Name:        "test-policy-example",
		Description: stringPtr("Example policy for testing IAM SDK integration"),
		Path:        stringPtr("/test/"),
		PolicyDocument: `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": [
						"s3:GetObject",
						"s3:ListBucket"
					],
					"Resource": [
						"arn:aws:s3:::example-bucket/*",
						"arn:aws:s3:::example-bucket"
					]
				}
			]
		}`,
		Tags: []configs.Tag{
			{Key: "Environment", Value: "Test"},
			{Key: "Purpose", Value: "SDK-Example"},
		},
	}

	output, err := awssdk.CreatePolicy(ctx, client, policy)
	if err != nil {
		fmt.Printf("Error creating policy: %v\n", err)
		fmt.Println("Note: This might fail if the policy already exists or due to permissions")
		return ""
	}

	fmt.Printf("✓ Policy created successfully!\n")
	fmt.Printf("  Name: %s\n", output.Name)
	fmt.Printf("  ARN: %s\n", output.ARN)
	fmt.Printf("  Path: %s\n", output.Path)
	if output.Description != nil {
		fmt.Printf("  Description: %s\n", *output.Description)
	}
	fmt.Printf("  Created: %s\n", output.CreateDate.Format("2006-01-02 15:04:05"))

	return output.ARN
}

// getPolicyExample retrieves a policy by ARN
func getPolicyExample(ctx context.Context, client *awssdk.AWSClient, arn string) {
	policy, err := awssdk.GetPolicy(ctx, client, arn)
	if err != nil {
		fmt.Printf("Error getting policy: %v\n", err)
		return
	}

	fmt.Printf("✓ Policy retrieved successfully!\n")
	fmt.Printf("  Name: %s\n", policy.Name)
	fmt.Printf("  ARN: %s\n", policy.ARN)
	fmt.Printf("  Path: %s\n", policy.Path)
	fmt.Printf("  Is AWS Managed: %v\n", policy.IsAWSManaged)
	fmt.Printf("  Attachment Count: %d\n", policy.AttachmentCount)
	if policy.DefaultVersionID != nil {
		fmt.Printf("  Default Version ID: %s\n", *policy.DefaultVersionID)
	}

	// Display policy document
	if policy.PolicyDocument != "" {
		fmt.Printf("  Policy Document:\n%s\n", policy.PolicyDocument)
	}
}

// updatePolicyExample updates a policy by creating a new version
func updatePolicyExample(ctx context.Context, client *awssdk.AWSClient, arn string) {
	updatedPolicy := &awsiam.Policy{
		Name:        "test-policy-example",
		Description: stringPtr("Updated example policy for testing IAM SDK integration"),
		Path:        stringPtr("/test/"),
		PolicyDocument: `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": [
						"s3:GetObject",
						"s3:ListBucket",
						"s3:PutObject"
					],
					"Resource": [
						"arn:aws:s3:::example-bucket/*",
						"arn:aws:s3:::example-bucket"
					]
				}
			]
		}`,
	}

	output, err := awssdk.UpdatePolicy(ctx, client, arn, updatedPolicy)
	if err != nil {
		fmt.Printf("Error updating policy: %v\n", err)
		return
	}

	fmt.Printf("✓ Policy updated successfully!\n")
	fmt.Printf("  ARN: %s\n", output.ARN)
	fmt.Printf("  Updated: %s\n", output.UpdateDate.Format("2006-01-02 15:04:05"))
	if output.DefaultVersionID != nil {
		fmt.Printf("  New Default Version ID: %s\n", *output.DefaultVersionID)
	}
}

// deletePolicyExample deletes a policy
func deletePolicyExample(ctx context.Context, client *awssdk.AWSClient, arn string) {
	err := awssdk.DeletePolicy(ctx, client, arn)
	if err != nil {
		fmt.Printf("Error deleting policy: %v\n", err)
		fmt.Println("Note: Policy might have attachments or might not exist")
		return
	}

	fmt.Printf("✓ Policy deleted successfully!\n")
	fmt.Printf("  ARN: %s\n", arn)
}

// displayPolicyTypes displays information about policy types
func displayPolicyTypes() {
	fmt.Println("IAM Policy Types:")
	fmt.Println("  1. AWS Managed Policies")
	fmt.Println("     - Pre-defined by AWS")
	fmt.Println("     - Cannot be modified or deleted")
	fmt.Println("     - ARN format: arn:aws:iam::aws:policy/<name>")
	fmt.Println("     - Examples: ReadOnlyAccess, PowerUserAccess, AdministratorAccess")
	fmt.Println()
	fmt.Println("  2. Customer Managed Policies")
	fmt.Println("     - Created by users/organizations")
	fmt.Println("     - Can be modified (creates new versions)")
	fmt.Println("     - Can be deleted (if no attachments)")
	fmt.Println("     - ARN format: arn:aws:iam::<account-id>:policy/<path><name>")
	fmt.Println()
	fmt.Println("  3. Inline Policies")
	fmt.Println("     - Embedded directly in users, roles, or groups")
	fmt.Println("     - Cannot be reused")
	fmt.Println("     - Deleted when the identity is deleted")
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
