package storage

import (
	"time"
)

// S3BucketOutput represents the output data for an S3 bucket after creation/update
type S3BucketOutput struct {
	// Core identifiers
	ID           string
	ARN          *string
	Name         string
	NamePrefix   *string
	Region       string

	// Configuration
	ForceDestroy bool
	Tags         map[string]string

	// Output fields (cloud-generated)
	BucketDomainName         *string // Standard DNS name
	BucketRegionalDomainName *string // Region-specific DNS name

	// CreatedAt timestamp
	CreatedAt *time.Time
}

// EBSVolumeOutput represents the output data for an EBS volume after creation/update
type EBSVolumeOutput struct {
	// Core identifiers
	ID               string
	ARN              *string
	Name             string
	Region           string
	AvailabilityZone string

	// Configuration
	Size      int
	Type      string
	IOPS      *int
	Throughput *int

	// Security & Backups
	Encrypted  bool
	KMSKeyID   *string
	SnapshotID *string

	// State
	State      EBSVolumeState
	AttachedTo *string

	// CreatedAt timestamp
	CreatedAt *time.Time
}

// S3BucketACLOutput represents the output data for S3 bucket ACL operations
type S3BucketACLOutput struct {
	ID                  string
	Bucket              string
	ACL                 *string
	GrantRead           *string
	GrantReadACP        *string
	GrantWrite          *string
	GrantWriteACP       *string
	GrantFullControl    *string
	GrantOwnerID        *string
	GrantOwnerDisplayName *string
	CreatedAt           *time.Time
}

// S3BucketVersioningOutput represents the output data for S3 bucket versioning operations
type S3BucketVersioningOutput struct {
	ID        string
	Bucket    string
	Status    *string // "Enabled", "Suspended", or nil
	MfaDelete *string // "Enabled" or "Disabled"
	CreatedAt *time.Time
}

// S3BucketEncryptionOutput represents the output data for S3 bucket encryption operations
type S3BucketEncryptionOutput struct {
	ID                 string
	Bucket             string
	SSEAlgorithm       *string // "AES256" or "aws:kms"
	KMSMasterKeyID     *string
	BucketKeyEnabled   *bool
	CreatedAt          *time.Time
}

// ToS3BucketOutput converts an S3Bucket domain model to S3BucketOutput
func ToS3BucketOutput(bucket *S3Bucket) *S3BucketOutput {
	if bucket == nil {
		return nil
	}
	return &S3BucketOutput{
		ID:                       bucket.ID,
		ARN:                      bucket.ARN,
		Name:                     bucket.Name,
		NamePrefix:               bucket.NamePrefix,
		Region:                   bucket.Region,
		ForceDestroy:             bucket.ForceDestroy,
		Tags:                     bucket.Tags,
		BucketDomainName:         bucket.BucketDomainName,
		BucketRegionalDomainName: bucket.BucketRegionalDomainName,
	}
}

// ToEBSVolumeOutput converts an EBSVolume domain model to EBSVolumeOutput
func ToEBSVolumeOutput(volume *EBSVolume) *EBSVolumeOutput {
	if volume == nil {
		return nil
	}
	return &EBSVolumeOutput{
		ID:               volume.ID,
		ARN:              volume.ARN,
		Name:             volume.Name,
		Region:           volume.Region,
		AvailabilityZone: volume.AvailabilityZone,
		Size:             volume.Size,
		Type:             volume.Type,
		IOPS:             volume.IOPS,
		Throughput:       volume.Throughput,
		Encrypted:        volume.Encrypted,
		KMSKeyID:         volume.KMSKeyID,
		SnapshotID:       volume.SnapshotID,
		State:            volume.State,
		AttachedTo:       volume.AttachedTo,
	}
}
