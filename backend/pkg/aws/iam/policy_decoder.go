package iam

import (
	"fmt"

	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
)

// DecodeAllPolicyFiles decodes all URL-encoded policy documents in service-specific JSON files
// This function should be run after fetching policies from SDK to convert them to readable JSON
func DecodeAllPolicyFiles() {
	fmt.Println("============================================")
	fmt.Println("DECODING POLICY DOCUMENTS")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("This will:")
	fmt.Println("  1. Read all service-specific policy JSON files")
	fmt.Println("  2. Decode URL-encoded policy_document fields")
	fmt.Println("  3. Save the decoded JSON back to the files")
	fmt.Println()

	if err := awssdk.DecodeAllPolicyFiles(); err != nil {
		fmt.Printf("Error decoding policy files: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("============================================")
	fmt.Println("SUCCESS: All policy documents decoded")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("All policy_document fields are now in readable JSON format")
	fmt.Println("instead of URL-encoded strings.")
	fmt.Println()
}
