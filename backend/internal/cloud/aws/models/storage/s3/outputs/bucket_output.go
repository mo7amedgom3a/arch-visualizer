package outputs

// BucketOutput represents AWS S3 bucket output/response data after creation
type BucketOutput struct {
	// AWS-generated identifiers
	ID  string `json:"id"`  // Bucket name (e.g., "my-bucket-name")
	ARN string `json:"arn"` // e.g., "arn:aws:s3:::my-bucket-name"

	// Basic information
	Name       string  `json:"name"`        // Bucket name
	NamePrefix *string `json:"name_prefix,omitempty"` // Bucket name prefix (if used)

	// Configuration
	ForceDestroy bool `json:"force_destroy"` // Force destroy flag

	// Output fields
	BucketDomainName       string `json:"bucket_domain_name"`        // Standard DNS name (e.g., bucket-name.s3.amazonaws.com)
	BucketRegionalDomainName string `json:"bucket_regional_domain_name"` // Region-specific DNS (e.g., bucket-name.s3.us-east-1.amazonaws.com)
	Region                 string `json:"region"`                   // AWS region

	// Tags
	Tags []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"tags"`
}
