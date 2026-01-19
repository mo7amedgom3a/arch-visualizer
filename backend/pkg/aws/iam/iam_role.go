package iam

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
)

// IAMRolesRunner demonstrates IAM role operations using AWS SDK
// Note: Role SDK functions are not yet implemented, so this shows the structure
func IAMRolesRunner() {
	ctx := context.Background()

	// Initialize AWS client
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		fmt.Printf("Error creating AWS client: %v\n", err)
		return
	}

	fmt.Println("============================================")
	fmt.Println("IAM ROLES OPERATIONS")
	fmt.Println("============================================")

	region := client.GetRegion()
	fmt.Printf("\nRegion: %s\n", region)
	fmt.Println("Note: IAM is a global service, but using region for consistency")
	fmt.Println()

	// Display role information
	fmt.Println("--- IAM Role Information ---")
	displayRoleInformation()

	// Example role configuration (for future SDK implementation)
	fmt.Println("\n--- Example Role Configuration ---")
	displayExampleRoleConfiguration()

	// Note about implementation status
	fmt.Println("\n--- Implementation Status ---")
	fmt.Println("IAM Role SDK functions are not yet implemented.")
	fmt.Println("When implemented, this runner will demonstrate:")
	fmt.Println("  - Creating IAM roles")
	fmt.Println("  - Getting role details")
	fmt.Println("  - Listing roles")
	fmt.Println("  - Updating roles")
	fmt.Println("  - Attaching/detaching policies to roles")
	fmt.Println("  - Managing inline policies on roles")
	fmt.Println("  - Deleting roles")
}

// displayRoleInformation displays information about IAM roles
func displayRoleInformation() {
	fmt.Println("IAM Roles:")
	fmt.Println("  - Roles are identities that can be assumed by trusted entities")
	fmt.Println("  - Used for cross-account access, service-to-service communication")
	fmt.Println("  - Require a trust policy (assume role policy)")
	fmt.Println("  - Can have managed policies and inline policies attached")
	fmt.Println()
	fmt.Println("Common Use Cases:")
	fmt.Println("  1. EC2 Instance Roles - Allow EC2 instances to access AWS services")
	fmt.Println("  2. Lambda Execution Roles - Grant permissions to Lambda functions")
	fmt.Println("  3. Cross-Account Access - Allow resources in one account to access another")
	fmt.Println("  4. Service Roles - Allow AWS services to perform actions on your behalf")
}

// displayExampleRoleConfiguration shows an example role configuration
func displayExampleRoleConfiguration() {
	exampleRole := &awsiam.Role{
		Name:        "example-ec2-role",
		Description: stringPtr("Example role for EC2 instances to access S3"),
		Path:        stringPtr("/"),
		AssumeRolePolicy: `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"Service": "ec2.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}
			]
		}`,
		ManagedPolicyARNs: []string{
			"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
		},
		Tags: []configs.Tag{
			{Key: "Environment", Value: "Example"},
			{Key: "Purpose", Value: "EC2-S3-Access"},
		},
	}

	fmt.Println("Example Role Configuration:")
	fmt.Printf("  Name: %s\n", exampleRole.Name)
	if exampleRole.Description != nil {
		fmt.Printf("  Description: %s\n", *exampleRole.Description)
	}
	if exampleRole.Path != nil {
		fmt.Printf("  Path: %s\n", *exampleRole.Path)
	}
	fmt.Printf("  Assume Role Policy:\n%s\n", exampleRole.AssumeRolePolicy)
	if len(exampleRole.ManagedPolicyARNs) > 0 {
		fmt.Println("  Managed Policies:")
		for _, arn := range exampleRole.ManagedPolicyARNs {
			fmt.Printf("    - %s\n", arn)
		}
	}
	if len(exampleRole.Tags) > 0 {
		fmt.Println("  Tags:")
		for _, tag := range exampleRole.Tags {
			fmt.Printf("    - %s: %s\n", tag.Key, tag.Value)
		}
	}
}

// IAMRunner is the main runner that calls both policy and role runners
func IAMRunner() {
	fmt.Println()
	fmt.Println("******************************************************************")
	IAMPoliciesRunner()
	fmt.Println()
	fmt.Println("******************************************************************")
	IAMRolesRunner()
	fmt.Println()
	fmt.Println("******************************************************************")
	IAMInstanceProfileRunner()
}
