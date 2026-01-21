package s3

import "fmt"

// BucketEncryption represents AWS S3 bucket encryption configuration
type BucketEncryption struct {
	Bucket string             `json:"bucket"`
	Rule   BucketEncryptionRule `json:"rule"`
}

// BucketEncryptionRule represents encryption rule
type BucketEncryptionRule struct {
	BucketKeyEnabled   bool                      `json:"bucket_key_enabled"`
	DefaultEncryption  BucketDefaultEncryption   `json:"default_encryption"`
}

// BucketDefaultEncryption represents default encryption settings
type BucketDefaultEncryption struct {
	SSEAlgorithm   string  `json:"sse_algorithm"` // AES256 | aws:kms
	KMSMasterKeyID *string `json:"kms_master_key_id,omitempty"`
}

// Validate validates encryption configuration
func (e *BucketEncryption) Validate() error {
	if e.Bucket == "" {
		return fmt.Errorf("bucket is required")
	}
	return e.Rule.Validate()
}

// Validate validates rule
func (r *BucketEncryptionRule) Validate() error {
	return r.DefaultEncryption.Validate()
}

// Validate validates default encryption
func (d *BucketDefaultEncryption) Validate() error {
	validAlgorithms := map[string]bool{
		"AES256":  true,
		"aws:kms": true,
	}

	if !validAlgorithms[d.SSEAlgorithm] {
		return fmt.Errorf("invalid sse_algorithm: %s", d.SSEAlgorithm)
	}

	if d.SSEAlgorithm == "aws:kms" {
		if d.KMSMasterKeyID == nil || *d.KMSMasterKeyID == "" {
			return fmt.Errorf("kms_master_key_id is required when sse_algorithm is aws:kms")
		}
	}

	return nil
}
