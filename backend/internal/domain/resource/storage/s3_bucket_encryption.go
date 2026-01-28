package storage

import (
	"fmt"

	domainerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/errors"
)

// S3BucketEncryption models server-side encryption configuration for a bucket
type S3BucketEncryption struct {
	Bucket string
	Rule   S3BucketEncryptionRule
}

// S3BucketEncryptionRule represents the encryption rule
type S3BucketEncryptionRule struct {
	BucketKeyEnabled bool
	DefaultEncryption S3BucketDefaultEncryption
}

// S3BucketDefaultEncryption represents default encryption settings
type S3BucketDefaultEncryption struct {
	SSEAlgorithm   string // AES256 | aws:kms
	KMSMasterKeyID *string
}

// Validate ensures encryption configuration is valid
func (e *S3BucketEncryption) Validate() error {
	if e.Bucket == "" {
		return domainerrors.New(domainerrors.CodeS3BucketRequired, domainerrors.KindValidation, "bucket is required")
	}
	return e.Rule.Validate()
}

// Validate ensures the encryption rule is valid
func (r *S3BucketEncryptionRule) Validate() error {
	return r.DefaultEncryption.Validate()
}

// Validate ensures default encryption settings are valid
func (d *S3BucketDefaultEncryption) Validate() error {
	validAlgorithms := map[string]bool{
		"AES256":  true,
		"aws:kms": true,
	}

	if !validAlgorithms[d.SSEAlgorithm] {
		return domainerrors.New(domainerrors.CodeS3InvalidSSEAlgorithm, domainerrors.KindValidation, fmt.Sprintf("invalid sse_algorithm: %s", d.SSEAlgorithm)).
			WithMeta("sse_algorithm", d.SSEAlgorithm)
	}

	if d.SSEAlgorithm == "aws:kms" {
		if d.KMSMasterKeyID == nil || *d.KMSMasterKeyID == "" {
			return domainerrors.New(domainerrors.CodeRequiredFieldMissing, domainerrors.KindValidation, "kms_master_key_id is required when sse_algorithm is aws:kms").
				WithMeta("field", "kms_master_key_id")
		}
	}

	return nil
}
