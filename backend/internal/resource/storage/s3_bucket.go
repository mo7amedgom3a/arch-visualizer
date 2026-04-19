package storage

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// S3Bucket represents a cloud-agnostic S3 bucket (object storage)
// This is the domain model - no cloud-specific details
type S3Bucket struct {
	ID           string
	ARN          *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	Name         string  // Bucket name (required if NamePrefix is not set)
	NamePrefix   *string // Bucket name prefix (required if Name is not set)
	Region       string  // Required: AWS region
	ForceDestroy bool    // If true, allows deletion of non-empty bucket
	Tags         map[string]string

	// Output fields (populated after creation)
	BucketDomainName         *string // Standard DNS name (e.g., bucket-name.s3.amazonaws.com)
	BucketRegionalDomainName *string // Region-specific DNS (e.g., bucket-name.s3.us-east-1.amazonaws.com)
}

// Validate performs domain-level validation
func (b *S3Bucket) Validate() error {
	// Name or NamePrefix must be provided
	if b.Name == "" && (b.NamePrefix == nil || *b.NamePrefix == "") {
		return errors.New("bucket name or name_prefix is required")
	}

	// Cannot have both Name and NamePrefix
	if b.Name != "" && b.NamePrefix != nil && *b.NamePrefix != "" {
		return errors.New("cannot specify both bucket name and name_prefix")
	}

	// Validate bucket name if provided
	if b.Name != "" {
		if err := validateBucketName(b.Name); err != nil {
			return fmt.Errorf("invalid bucket name: %w", err)
		}
	}

	// Validate bucket name prefix if provided
	if b.NamePrefix != nil && *b.NamePrefix != "" {
		if err := validateBucketNamePrefix(*b.NamePrefix); err != nil {
			return fmt.Errorf("invalid bucket name prefix: %w", err)
		}
	}

	// Region is required
	if b.Region == "" {
		return errors.New("bucket region is required")
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

// Ensure S3Bucket implements StorageResource
var _ StorageResource = (*S3Bucket)(nil)

// GetID returns the bucket ID (name)
func (b *S3Bucket) GetID() string {
	return b.ID
}

// GetName returns the bucket name
func (b *S3Bucket) GetName() string {
	return b.Name
}

// GetRegion returns the bucket region
func (b *S3Bucket) GetRegion() string {
	return b.Region
}
