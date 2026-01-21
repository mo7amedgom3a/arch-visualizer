package storage

import (
	"context"
	"fmt"

	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/storage"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/storage"
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
)

// AWSStorageAdapter adapts AWS-specific storage service to domain storage service
// This implements the Adapter pattern, allowing the domain layer to work with cloud-specific implementations
type AWSStorageAdapter struct {
	awsService awsservice.AWSStorageService
}

// NewAWSStorageAdapter creates a new AWS storage adapter
func NewAWSStorageAdapter(awsService awsservice.AWSStorageService) domainstorage.StorageService {
	return &AWSStorageAdapter{
		awsService: awsService,
	}
}

// Ensure AWSStorageAdapter implements StorageService
var _ domainstorage.StorageService = (*AWSStorageAdapter)(nil)

// EBS Volume Operations

func (a *AWSStorageAdapter) CreateEBSVolume(ctx context.Context, volume *domainstorage.EBSVolume) (*domainstorage.EBSVolume, error) {
	if err := volume.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsVolume := awsmapper.FromDomainEBSVolume(volume)
	if err := awsVolume.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsVolumeOutput, err := a.awsService.CreateEBSVolume(ctx, awsVolume)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainVolume := awsmapper.ToDomainEBSVolumeFromOutput(awsVolumeOutput)
	// Preserve region from input
	domainVolume.Region = volume.Region
	return domainVolume, nil
}

func (a *AWSStorageAdapter) GetEBSVolume(ctx context.Context, id string) (*domainstorage.EBSVolume, error) {
	awsVolumeOutput, err := a.awsService.GetEBSVolume(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainEBSVolumeFromOutput(awsVolumeOutput), nil
}

func (a *AWSStorageAdapter) UpdateEBSVolume(ctx context.Context, id string, volume *domainstorage.EBSVolume) (*domainstorage.EBSVolume, error) {
	if err := volume.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsVolume := awsmapper.FromDomainEBSVolume(volume)
	if err := awsVolume.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsVolumeOutput, err := a.awsService.UpdateEBSVolume(ctx, id, awsVolume)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainVolume := awsmapper.ToDomainEBSVolumeFromOutput(awsVolumeOutput)
	// Preserve region from input
	domainVolume.Region = volume.Region
	return domainVolume, nil
}

func (a *AWSStorageAdapter) DeleteEBSVolume(ctx context.Context, id string) error {
	if err := a.awsService.DeleteEBSVolume(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSStorageAdapter) ListEBSVolumes(ctx context.Context, filters map[string]string) ([]*domainstorage.EBSVolume, error) {
	// Convert domain filters to AWS filters format
	awsFilters := make(map[string][]string)
	for key, value := range filters {
		awsFilters[key] = []string{value}
	}

	awsVolumeOutputs, err := a.awsService.ListEBSVolumes(ctx, awsFilters)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainVolumes := make([]*domainstorage.EBSVolume, len(awsVolumeOutputs))
	for i, awsVolumeOutput := range awsVolumeOutputs {
		domainVolumes[i] = awsmapper.ToDomainEBSVolumeFromOutput(awsVolumeOutput)
	}

	return domainVolumes, nil
}

// Volume Attachment Operations

func (a *AWSStorageAdapter) AttachVolume(ctx context.Context, volumeID, instanceID, deviceName string) error {
	if err := a.awsService.AttachVolume(ctx, volumeID, instanceID, deviceName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSStorageAdapter) DetachVolume(ctx context.Context, volumeID, instanceID string) error {
	if err := a.awsService.DetachVolume(ctx, volumeID, instanceID); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

// S3 Bucket Operations

func (a *AWSStorageAdapter) CreateS3Bucket(ctx context.Context, bucket *domainstorage.S3Bucket) (*domainstorage.S3Bucket, error) {
	if err := bucket.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsBucket := awsmapper.FromDomainS3Bucket(bucket)
	if err := awsBucket.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsBucketOutput, err := a.awsService.CreateS3Bucket(ctx, awsBucket)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainBucket := awsmapper.ToDomainS3BucketFromOutput(awsBucketOutput)
	// Preserve region from input
	domainBucket.Region = bucket.Region
	return domainBucket, nil
}

func (a *AWSStorageAdapter) GetS3Bucket(ctx context.Context, id string) (*domainstorage.S3Bucket, error) {
	awsBucketOutput, err := a.awsService.GetS3Bucket(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainS3BucketFromOutput(awsBucketOutput), nil
}

func (a *AWSStorageAdapter) UpdateS3Bucket(ctx context.Context, id string, bucket *domainstorage.S3Bucket) (*domainstorage.S3Bucket, error) {
	if err := bucket.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsBucket := awsmapper.FromDomainS3Bucket(bucket)
	if err := awsBucket.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsBucketOutput, err := a.awsService.UpdateS3Bucket(ctx, id, awsBucket)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainBucket := awsmapper.ToDomainS3BucketFromOutput(awsBucketOutput)
	// Preserve region from input
	domainBucket.Region = bucket.Region
	return domainBucket, nil
}

func (a *AWSStorageAdapter) DeleteS3Bucket(ctx context.Context, id string) error {
	if err := a.awsService.DeleteS3Bucket(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSStorageAdapter) ListS3Buckets(ctx context.Context, filters map[string]string) ([]*domainstorage.S3Bucket, error) {
	// Convert domain filters to AWS filters format
	awsFilters := make(map[string][]string)
	for key, value := range filters {
		awsFilters[key] = []string{value}
	}

	awsBucketOutputs, err := a.awsService.ListS3Buckets(ctx, awsFilters)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainBuckets := make([]*domainstorage.S3Bucket, len(awsBucketOutputs))
	for i, awsBucketOutput := range awsBucketOutputs {
		domainBuckets[i] = awsmapper.ToDomainS3BucketFromOutput(awsBucketOutput)
	}

	return domainBuckets, nil
}

// S3 Bucket ACL Operations

func (a *AWSStorageAdapter) UpdateS3BucketACL(ctx context.Context, bucket string, acl *domainstorage.S3BucketACL) (*domainstorage.S3BucketACL, error) {
	if acl == nil {
		return nil, fmt.Errorf("acl is required")
	}
	if err := acl.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsACL := awsmapper.FromDomainS3BucketACL(acl)
	if err := awsACL.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsOutput, err := a.awsService.UpdateS3BucketACL(ctx, bucket, awsACL)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainACL := awsmapper.ToDomainS3BucketACLFromOutput(awsOutput)
	return domainACL, nil
}

func (a *AWSStorageAdapter) GetS3BucketACL(ctx context.Context, bucket string) (*domainstorage.S3BucketACL, error) {
	awsOutput, err := a.awsService.GetS3BucketACL(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}
	return awsmapper.ToDomainS3BucketACLFromOutput(awsOutput), nil
}

// S3 Bucket Versioning Operations

func (a *AWSStorageAdapter) UpdateS3BucketVersioning(ctx context.Context, bucket string, versioning *domainstorage.S3BucketVersioning) (*domainstorage.S3BucketVersioning, error) {
	if versioning == nil {
		return nil, fmt.Errorf("versioning is required")
	}
	if err := versioning.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsVersioning := awsmapper.FromDomainS3BucketVersioning(versioning)
	if err := awsVersioning.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsOutput, err := a.awsService.UpdateS3BucketVersioning(ctx, bucket, awsVersioning)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainS3BucketVersioningFromOutput(awsOutput), nil
}

func (a *AWSStorageAdapter) GetS3BucketVersioning(ctx context.Context, bucket string) (*domainstorage.S3BucketVersioning, error) {
	awsOutput, err := a.awsService.GetS3BucketVersioning(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}
	return awsmapper.ToDomainS3BucketVersioningFromOutput(awsOutput), nil
}

// S3 Bucket Encryption Operations

func (a *AWSStorageAdapter) UpdateS3BucketEncryption(ctx context.Context, bucket string, encryption *domainstorage.S3BucketEncryption) (*domainstorage.S3BucketEncryption, error) {
	if encryption == nil {
		return nil, fmt.Errorf("encryption is required")
	}
	if err := encryption.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsEncryption := awsmapper.FromDomainS3BucketEncryption(encryption)
	if err := awsEncryption.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsOutput, err := a.awsService.UpdateS3BucketEncryption(ctx, bucket, awsEncryption)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainS3BucketEncryptionFromOutput(awsOutput), nil
}

func (a *AWSStorageAdapter) GetS3BucketEncryption(ctx context.Context, bucket string) (*domainstorage.S3BucketEncryption, error) {
	awsOutput, err := a.awsService.GetS3BucketEncryption(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}
	return awsmapper.ToDomainS3BucketEncryptionFromOutput(awsOutput), nil
}
