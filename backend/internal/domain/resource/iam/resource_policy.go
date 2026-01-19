package iam

import "errors"

// ResourcePolicy represents a resource-based policy attached to an AWS resource
// (e.g., S3 bucket policy, Lambda function policy, etc.)
type ResourcePolicy struct {
	ResourceType string // Type of resource (e.g., "s3_bucket", "lambda_function")
	ResourceARN  string // ARN of the resource
	Policy       string // JSON string containing the resource policy document
}

// Validate performs basic validation on a resource policy
func (rp *ResourcePolicy) Validate() error {
	if rp.ResourceType == "" {
		return errors.New("resource type is required")
	}
	if rp.ResourceARN == "" {
		return errors.New("resource ARN is required")
	}
	if rp.Policy == "" {
		return errors.New("resource policy document is required")
	}
	// JSON validation would be done at the AWS adapter level
	return nil
}
