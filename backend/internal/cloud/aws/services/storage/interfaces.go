package storage

import (
	"context"

	awsebs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs"
	awsebsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs/outputs"
	awss3 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3"
	awss3outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3/outputs"
)

// AWSStorageService defines AWS-specific storage operations
// This implements cloud provider-specific logic while maintaining domain compatibility
type AWSStorageService interface {
	// EBS Volume operations
	CreateEBSVolume(ctx context.Context, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error)
	GetEBSVolume(ctx context.Context, id string) (*awsebsoutputs.VolumeOutput, error)
	UpdateEBSVolume(ctx context.Context, id string, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error)
	DeleteEBSVolume(ctx context.Context, id string) error
	ListEBSVolumes(ctx context.Context, filters map[string][]string) ([]*awsebsoutputs.VolumeOutput, error)

	// Volume attachment operations
	AttachVolume(ctx context.Context, volumeID, instanceID, deviceName string) error
	DetachVolume(ctx context.Context, volumeID, instanceID string) error

	// S3 Bucket operations
	CreateS3Bucket(ctx context.Context, bucket *awss3.Bucket) (*awss3outputs.BucketOutput, error)
	GetS3Bucket(ctx context.Context, id string) (*awss3outputs.BucketOutput, error)
	UpdateS3Bucket(ctx context.Context, id string, bucket *awss3.Bucket) (*awss3outputs.BucketOutput, error)
	DeleteS3Bucket(ctx context.Context, id string) error
	ListS3Buckets(ctx context.Context, filters map[string][]string) ([]*awss3outputs.BucketOutput, error)

	// S3 Bucket ACL operations
	UpdateS3BucketACL(ctx context.Context, bucket string, acl *awss3.BucketACL) (*awss3outputs.BucketACLOutput, error)
	GetS3BucketACL(ctx context.Context, bucket string) (*awss3outputs.BucketACLOutput, error)

	// S3 Bucket Versioning operations
	UpdateS3BucketVersioning(ctx context.Context, bucket string, versioning *awss3.BucketVersioning) (*awss3outputs.BucketVersioningOutput, error)
	GetS3BucketVersioning(ctx context.Context, bucket string) (*awss3outputs.BucketVersioningOutput, error)

	// S3 Bucket Encryption operations
	UpdateS3BucketEncryption(ctx context.Context, bucket string, encryption *awss3.BucketEncryption) (*awss3outputs.BucketEncryptionOutput, error)
	GetS3BucketEncryption(ctx context.Context, bucket string) (*awss3outputs.BucketEncryptionOutput, error)
}
