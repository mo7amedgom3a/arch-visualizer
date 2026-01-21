package s3

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// Bucket represents an AWS S3 bucket configuration
type Bucket struct {
	Bucket       *string       `json:"bucket,omitempty"`        // Optional: bucket name
	BucketPrefix *string       `json:"bucket_prefix,omitempty"` // Optional: bucket name prefix
	ForceDestroy bool          `json:"force_destroy"`           // Default: false
	Tags         []configs.Tag `json:"tags,omitempty"`          // Optional tags
}

// Validate performs AWS-specific validation
func (b *Bucket) Validate() error {
	// Bucket or BucketPrefix must be provided
	if (b.Bucket == nil || *b.Bucket == "") && (b.BucketPrefix == nil || *b.BucketPrefix == "") {
		return errors.New("bucket name or bucket_prefix is required")
	}

	// Cannot have both Bucket and BucketPrefix
	if b.Bucket != nil && *b.Bucket != "" && b.BucketPrefix != nil && *b.BucketPrefix != "" {
		return errors.New("cannot specify both bucket name and bucket_prefix")
	}

	// Validate bucket name if provided
	if b.Bucket != nil && *b.Bucket != "" {
		if err := validateBucketName(*b.Bucket); err != nil {
			return fmt.Errorf("invalid bucket name: %w", err)
		}
	}

	// Validate bucket name prefix if provided
	if b.BucketPrefix != nil && *b.BucketPrefix != "" {
		if err := validateBucketNamePrefix(*b.BucketPrefix); err != nil {
			return fmt.Errorf("invalid bucket name prefix: %w", err)
		}
	}

	return nil
}

// validateBucketName validates S3 bucket naming rules
// Based on AWS S3 bucket naming rules:
// - 3-63 characters long
// - Can contain lowercase letters, numbers, dots (.), and hyphens (-)
// - Must begin and end with a letter or number
// - Must not be formatted as an IP address
// - Must not start with "xn--" or "sthree-"
func validateBucketName(name string) error {
	if len(name) < 3 {
		return errors.New("bucket name must be at least 3 characters long")
	}
	if len(name) > 63 {
		return errors.New("bucket name must be at most 63 characters long")
	}

	// Must begin and end with a letter or number
	firstChar := name[0]
	lastChar := name[len(name)-1]
	if !isAlphanumeric(firstChar) {
		return errors.New("bucket name must begin with a letter or number")
	}
	if !isAlphanumeric(lastChar) {
		return errors.New("bucket name must end with a letter or number")
	}

	// Can only contain lowercase letters, numbers, dots, and hyphens
	validPattern := regexp.MustCompile(`^[a-z0-9.-]+$`)
	if !validPattern.MatchString(name) {
		return errors.New("bucket name can only contain lowercase letters, numbers, dots (.), and hyphens (-)")
	}

	// Must not start with "xn--" or "sthree-"
	if strings.HasPrefix(name, "xn--") {
		return errors.New("bucket name must not start with 'xn--'")
	}
	if strings.HasPrefix(name, "sthree-") {
		return errors.New("bucket name must not start with 'sthree-'")
	}

	// Must not be formatted as an IP address (basic check)
	ipPattern := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if ipPattern.MatchString(name) {
		return errors.New("bucket name must not be formatted as an IP address")
	}

	// Must not contain consecutive dots
	if strings.Contains(name, "..") {
		return errors.New("bucket name must not contain consecutive dots")
	}

	return nil
}

// validateBucketNamePrefix validates S3 bucket name prefix
// Prefixes can end with a hyphen since they will be appended to
func validateBucketNamePrefix(prefix string) error {
	if len(prefix) < 1 {
		return errors.New("bucket name prefix must be at least 1 character long")
	}
	if len(prefix) > 63 {
		return errors.New("bucket name prefix must be at most 63 characters long")
	}

	// Must begin with a letter or number
	firstChar := prefix[0]
	if !isAlphanumeric(firstChar) {
		return errors.New("bucket name prefix must begin with a letter or number")
	}

	// Can only contain lowercase letters, numbers, dots, and hyphens
	validPattern := regexp.MustCompile(`^[a-z0-9.-]+$`)
	if !validPattern.MatchString(prefix) {
		return errors.New("bucket name prefix can only contain lowercase letters, numbers, dots (.), and hyphens (-)")
	}

	// Must not start with "xn--" or "sthree-"
	if strings.HasPrefix(prefix, "xn--") {
		return errors.New("bucket name prefix must not start with 'xn--'")
	}
	if strings.HasPrefix(prefix, "sthree-") {
		return errors.New("bucket name prefix must not start with 'sthree-'")
	}

	// Must not contain consecutive dots
	if strings.Contains(prefix, "..") {
		return errors.New("bucket name prefix must not contain consecutive dots")
	}

	return nil
}

// isAlphanumeric checks if a character is alphanumeric
func isAlphanumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}
