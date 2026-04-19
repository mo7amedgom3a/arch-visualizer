package storage

import (
	"context"
)

// StorageService defines the interface for storage resource operations
// This is cloud-agnostic and can be implemented by any cloud provider
type StorageService interface {
	// EBS Volume operations
	CreateEBSVolume(ctx context.Context, volume *EBSVolume) (*EBSVolume, error)
	GetEBSVolume(ctx context.Context, id string) (*EBSVolume, error)
	UpdateEBSVolume(ctx context.Context, id string, volume *EBSVolume) (*EBSVolume, error)
	DeleteEBSVolume(ctx context.Context, id string) error
	ListEBSVolumes(ctx context.Context, filters map[string]string) ([]*EBSVolume, error)

	// Volume attachment operations
	AttachVolume(ctx context.Context, volumeID, instanceID, deviceName string) error
	DetachVolume(ctx context.Context, volumeID, instanceID string) error

	// S3 Bucket operations
	CreateS3Bucket(ctx context.Context, bucket *S3Bucket) (*S3Bucket, error)
	GetS3Bucket(ctx context.Context, id string) (*S3Bucket, error)
	UpdateS3Bucket(ctx context.Context, id string, bucket *S3Bucket) (*S3Bucket, error)
	DeleteS3Bucket(ctx context.Context, id string) error
	ListS3Buckets(ctx context.Context, filters map[string]string) ([]*S3Bucket, error)

	// S3 Bucket ACL operations
	UpdateS3BucketACL(ctx context.Context, bucket string, acl *S3BucketACL) (*S3BucketACL, error)
	GetS3BucketACL(ctx context.Context, bucket string) (*S3BucketACL, error)

	// S3 Bucket Versioning operations
	UpdateS3BucketVersioning(ctx context.Context, bucket string, versioning *S3BucketVersioning) (*S3BucketVersioning, error)
	GetS3BucketVersioning(ctx context.Context, bucket string) (*S3BucketVersioning, error)

	// S3 Bucket Encryption operations
	UpdateS3BucketEncryption(ctx context.Context, bucket string, encryption *S3BucketEncryption) (*S3BucketEncryption, error)
	GetS3BucketEncryption(ctx context.Context, bucket string) (*S3BucketEncryption, error)
}
