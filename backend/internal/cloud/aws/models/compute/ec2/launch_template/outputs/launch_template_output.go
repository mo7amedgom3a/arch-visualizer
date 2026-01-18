package outputs

import "time"

// LaunchTemplateOutput represents AWS EC2 Launch Template output/response data after creation
type LaunchTemplateOutput struct {
	// AWS-generated identifiers
	ID  string `json:"id"`  // e.g., "lt-0123456789abcdef0"
	ARN string `json:"arn"` // e.g., "arn:aws:ec2:us-east-1:123456789012:launch-template/lt-0123456789abcdef0"

	// Basic information
	Name string `json:"name"`

	// Version information
	DefaultVersion int `json:"default_version"` // Default version number
	LatestVersion  int `json:"latest_version"`  // Latest version number

	// Metadata
	CreateTime time.Time `json:"create_time"`
	CreatedBy  string    `json:"created_by"` // IAM user/role ARN

	// Tags (using struct to match AWS SDK Tag format)
	Tags []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"tags"`
}
