package iam

import (
	"context"
	"fmt"

	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// IAMInstanceProfileRunner demonstrates IAM instance profile operations using AWS SDK
func IAMInstanceProfileRunner() {
	ctx := context.Background()

	// Initialize AWS client
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		fmt.Printf("Error creating AWS client: %v\n", err)
		return
	}

	fmt.Println("============================================")
	fmt.Println("IAM INSTANCE PROFILE OPERATIONS")
	fmt.Println("============================================")

	region := client.GetRegion()
	fmt.Printf("\nRegion: %s\n", region)
	fmt.Println("Note: IAM is a global service, but using region for consistency")
	fmt.Println()

	// 1. List Instance Profiles
	fmt.Println("--- 1. Listing Instance Profiles ---")
	listInstanceProfiles(ctx, client)

	// 2. Create an Instance Profile (example - will be cleaned up)
	fmt.Println("\n--- 2. Creating Instance Profile (Example) ---")
	profileName := createInstanceProfileExample(ctx, client)

	// 3. Get the created profile
	if profileName != "" {
		fmt.Println("\n--- 3. Getting Created Instance Profile ---")
		getInstanceProfileExample(ctx, client, profileName)

		// 4. Display Instance Profile Information
		fmt.Println("\n--- 4. Instance Profile Information ---")
		displayInstanceProfileInfo()

		// 5. Cleanup: Delete the instance profile
		fmt.Println("\n--- 5. Cleaning Up: Deleting Instance Profile ---")
		deleteInstanceProfileExample(ctx, client, profileName)
	}
}

// listInstanceProfiles lists all instance profiles
func listInstanceProfiles(ctx context.Context, client *awssdk.AWSClient) {
	profiles, err := awssdk.ListInstanceProfiles(ctx, client, nil)
	if err != nil {
		fmt.Printf("Error listing instance profiles: %v\n", err)
		return
	}

	fmt.Printf("Found %d instance profiles\n", len(profiles))
	if len(profiles) > 0 {
		fmt.Println("\nInstance Profiles:")
		for i, profile := range profiles {
			fmt.Printf("  %d. %s\n", i+1, profile.Name)
			fmt.Printf("     ARN: %s\n", profile.ARN)
			fmt.Printf("     Path: %s\n", profile.Path)
			fmt.Printf("     Created: %s\n", profile.CreateDate.Format("2006-01-02 15:04:05"))
			if len(profile.Roles) > 0 {
				fmt.Printf("     Roles: %d attached\n", len(profile.Roles))
				for _, role := range profile.Roles {
					if role != nil {
						fmt.Printf("       - %s\n", role.Name)
					}
				}
			} else {
				fmt.Println("     Roles: None")
			}
			fmt.Println()
		}
	} else {
		fmt.Println("No instance profiles found")
	}
}

// createInstanceProfileExample creates a sample instance profile
func createInstanceProfileExample(ctx context.Context, client *awssdk.AWSClient) string {
	profile := &awsiam.InstanceProfile{
		Name: "test-instance-profile-example",
		Path: stringPtr("/test/"),
		Tags: []configs.Tag{
			{Key: "Environment", Value: "Test"},
			{Key: "Purpose", Value: "SDK-Example"},
		},
	}

	output, err := awssdk.CreateInstanceProfile(ctx, client, profile)
	if err != nil {
		fmt.Printf("Error creating instance profile: %v\n", err)
		fmt.Println("Note: This might fail if the profile already exists or due to permissions")
		return ""
	}

	fmt.Printf("✓ Instance profile created successfully!\n")
	fmt.Printf("  Name: %s\n", output.Name)
	fmt.Printf("  ARN: %s\n", output.ARN)
	fmt.Printf("  Path: %s\n", output.Path)
	fmt.Printf("  Created: %s\n", output.CreateDate.Format("2006-01-02 15:04:05"))

	return output.Name
}

// getInstanceProfileExample retrieves an instance profile by name
func getInstanceProfileExample(ctx context.Context, client *awssdk.AWSClient, name string) {
	profile, err := awssdk.GetInstanceProfile(ctx, client, name)
	if err != nil {
		fmt.Printf("Error getting instance profile: %v\n", err)
		return
	}

	fmt.Printf("✓ Instance profile retrieved successfully!\n")
	fmt.Printf("  Name: %s\n", profile.Name)
	fmt.Printf("  ARN: %s\n", profile.ARN)
	fmt.Printf("  Path: %s\n", profile.Path)
	fmt.Printf("  Created: %s\n", profile.CreateDate.Format("2006-01-02 15:04:05"))

	if len(profile.Roles) > 0 {
		fmt.Printf("  Attached Roles:\n")
		for _, role := range profile.Roles {
			if role != nil {
				fmt.Printf("    - %s (ARN: %s)\n", role.Name, role.ARN)
			}
		}
	} else {
		fmt.Println("  Attached Roles: None")
	}
}

// deleteInstanceProfileExample deletes an instance profile
func deleteInstanceProfileExample(ctx context.Context, client *awssdk.AWSClient, name string) {
	err := awssdk.DeleteInstanceProfile(ctx, client, name)
	if err != nil {
		fmt.Printf("Error deleting instance profile: %v\n", err)
		fmt.Println("Note: Instance profile might have attached roles or might not exist")
		return
	}

	fmt.Printf("✓ Instance profile deleted successfully!\n")
	fmt.Printf("  Name: %s\n", name)
}

// displayInstanceProfileInfo displays information about instance profiles
func displayInstanceProfileInfo() {
	fmt.Println("IAM Instance Profiles:")
	fmt.Println("  Instance profiles are containers for IAM roles that can be attached to EC2 instances.")
	fmt.Println("  They allow EC2 instances to assume IAM roles and access AWS services using temporary credentials.")
	fmt.Println("")
	fmt.Println("  Relationship:")
	fmt.Println("    EC2 Instance → IAMInstanceProfile (name/arn) → Instance Profile → IAM Role")
	fmt.Println("")
	fmt.Println("  Key Features:")
	fmt.Println("    - One instance profile can contain one IAM role")
	fmt.Println("    - Instance profiles are attached to EC2 instances at launch time")
	fmt.Println("    - EC2 instances use the role's permissions to access AWS services")
	fmt.Println("    - No need to store AWS credentials on the instance")
	fmt.Println("")
	fmt.Println("  Common Use Cases:")
	fmt.Println("    - Allow EC2 instances to access S3 buckets")
	fmt.Println("    - Enable instances to write to CloudWatch Logs")
	fmt.Println("    - Grant instances permissions to access other AWS services")
}

