package storage

import (
	"context"
)

// StorageOutputService defines the interface for storage resource operations that return output DTOs
// This is a parallel interface to StorageService, providing output-specific models
type StorageOutputService interface {
	// EBS Volume operations
	CreateEBSVolumeOutput(ctx context.Context, volume *EBSVolume) (*EBSVolumeOutput, error)
	GetEBSVolumeOutput(ctx context.Context, id string) (*EBSVolumeOutput, error)
	UpdateEBSVolumeOutput(ctx context.Context, id string, volume *EBSVolume) (*EBSVolumeOutput, error)
	ListEBSVolumesOutput(ctx context.Context, filters map[string]string) ([]*EBSVolumeOutput, error)

	// S3 Bucket operations
	CreateS3BucketOutput(ctx context.Context, bucket *S3Bucket) (*S3BucketOutput, error)
	GetS3BucketOutput(ctx context.Context, id string) (*S3BucketOutput, error)
	UpdateS3BucketOutput(ctx context.Context, id string, bucket *S3Bucket) (*S3BucketOutput, error)
	ListS3BucketsOutput(ctx context.Context, filters map[string]string) ([]*S3BucketOutput, error)

	// S3 Bucket ACL operations
	UpdateS3BucketACLOutput(ctx context.Context, bucket string, acl *S3BucketACL) (*S3BucketACLOutput, error)
	GetS3BucketACLOutput(ctx context.Context, bucket string) (*S3BucketACLOutput, error)

	// S3 Bucket Versioning operations
	UpdateS3BucketVersioningOutput(ctx context.Context, bucket string, versioning *S3BucketVersioning) (*S3BucketVersioningOutput, error)
	GetS3BucketVersioningOutput(ctx context.Context, bucket string) (*S3BucketVersioningOutput, error)

	// S3 Bucket Encryption operations
	UpdateS3BucketEncryptionOutput(ctx context.Context, bucket string, encryption *S3BucketEncryption) (*S3BucketEncryptionOutput, error)
	GetS3BucketEncryptionOutput(ctx context.Context, bucket string) (*S3BucketEncryptionOutput, error)
}
