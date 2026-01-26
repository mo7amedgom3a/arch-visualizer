package storage

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services"
	awsebs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs"
	awsebsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs/outputs"
	awss3 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3"
	awss3outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3/outputs"
)

// StorageService implements AWSStorageService with deterministic virtual operations
type StorageService struct{}

// NewStorageService creates a new storage service implementation
func NewStorageService() *StorageService {
	return &StorageService{}
}

// EBS Volume operations

func (s *StorageService) CreateEBSVolume(ctx context.Context, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error) {
	if volume == nil {
		return nil, fmt.Errorf("volume is nil")
	}

	volumeID := fmt.Sprintf("vol-%s", services.GenerateDeterministicID(volume.Name)[:15])
	region := "us-east-1"
	arn := services.GenerateARN("ec2", "volume", volumeID, region)

	return &awsebsoutputs.VolumeOutput{
		ID:               volumeID,
		ARN:              arn,
		Name:             volume.Name,
		AvailabilityZone: volume.AvailabilityZone,
		Size:             volume.Size,
		VolumeType:       volume.VolumeType,
		IOPS:             volume.IOPS,
		Throughput:       volume.Throughput,
		Encrypted:        volume.Encrypted,
		KMSKeyID:         volume.KMSKeyID,
		SnapshotID:       volume.SnapshotID,
		State:            "available",
		AttachedTo:       nil,
		CreateTime:       services.GetFixedTimestamp(),
		Tags: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			{Key: "Name", Value: volume.Name},
		},
	}, nil
}

func (s *StorageService) GetEBSVolume(ctx context.Context, id string) (*awsebsoutputs.VolumeOutput, error) {
	region := "us-east-1"
	arn := services.GenerateARN("ec2", "volume", id, region)

	return &awsebsoutputs.VolumeOutput{
		ID:               id,
		ARN:              arn,
		Name:             "test-volume",
		AvailabilityZone: "us-east-1a",
		Size:             40,
		VolumeType:       "gp3",
		IOPS:             services.IntPtr(3000),
		Throughput:       services.IntPtr(125),
		Encrypted:        false,
		KMSKeyID:         nil,
		SnapshotID:       nil,
		State:            "available",
		AttachedTo:       nil,
		CreateTime:       services.GetFixedTimestamp(),
		Tags: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			{Key: "Name", Value: "test-volume"},
		},
	}, nil
}

func (s *StorageService) UpdateEBSVolume(ctx context.Context, id string, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error) {
	return s.CreateEBSVolume(ctx, volume)
}

func (s *StorageService) DeleteEBSVolume(ctx context.Context, id string) error {
	return nil
}

func (s *StorageService) ListEBSVolumes(ctx context.Context, filters map[string][]string) ([]*awsebsoutputs.VolumeOutput, error) {
	return []*awsebsoutputs.VolumeOutput{
		{
			ID:               "vol-0a1b2c3d4e5f6g7h8",
			ARN:              "arn:aws:ec2:us-east-1:123456789012:volume/vol-0a1b2c3d4e5f6g7h8",
			Name:             "test-volume",
			AvailabilityZone: "us-east-1a",
			Size:             40,
			VolumeType:       "gp3",
			IOPS:             services.IntPtr(3000),
			Throughput:       services.IntPtr(125),
			Encrypted:        false,
			KMSKeyID:         nil,
			SnapshotID:       nil,
			State:            "available",
			AttachedTo:       nil,
			CreateTime:       services.GetFixedTimestamp(),
			Tags: []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}{
				{Key: "Name", Value: "test-volume"},
			},
		},
	}, nil
}

// Volume attachment operations

func (s *StorageService) AttachVolume(ctx context.Context, volumeID, instanceID, deviceName string) error {
	return nil
}

func (s *StorageService) DetachVolume(ctx context.Context, volumeID, instanceID string) error {
	return nil
}

// S3 Bucket operations

func (s *StorageService) CreateS3Bucket(ctx context.Context, bucket *awss3.Bucket) (*awss3outputs.BucketOutput, error) {
	if bucket == nil {
		return nil, fmt.Errorf("bucket is nil")
	}

	bucketName := ""
	if bucket.Bucket != nil && *bucket.Bucket != "" {
		bucketName = *bucket.Bucket
	} else if bucket.BucketPrefix != nil && *bucket.BucketPrefix != "" {
		bucketName = *bucket.BucketPrefix + "mock1234567890abcdef"
	}
	if bucketName == "" {
		return nil, fmt.Errorf("bucket name or prefix must be provided")
	}

	bucketARN := fmt.Sprintf("arn:aws:s3:::%s", bucketName)
	bucketDomainName := fmt.Sprintf("%s.s3.amazonaws.com", bucketName)
	bucketRegionalDomainName := fmt.Sprintf("%s.s3.us-east-1.amazonaws.com", bucketName)

	output := &awss3outputs.BucketOutput{
		ID:                       bucketName,
		ARN:                      bucketARN,
		Name:                     bucketName,
		NamePrefix:               bucket.BucketPrefix,
		ForceDestroy:             bucket.ForceDestroy,
		BucketDomainName:         bucketDomainName,
		BucketRegionalDomainName: bucketRegionalDomainName,
		Region:                   "us-east-1",
		Tags: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{},
	}

	// Convert tags
	if bucket.Tags != nil {
		for _, tag := range bucket.Tags {
			output.Tags = append(output.Tags, struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}{
				Key:   tag.Key,
				Value: tag.Value,
			})
		}
	}

	return output, nil
}

func (s *StorageService) GetS3Bucket(ctx context.Context, id string) (*awss3outputs.BucketOutput, error) {
	bucketARN := fmt.Sprintf("arn:aws:s3:::%s", id)
	bucketDomainName := fmt.Sprintf("%s.s3.amazonaws.com", id)
	bucketRegionalDomainName := fmt.Sprintf("%s.s3.us-east-1.amazonaws.com", id)

	return &awss3outputs.BucketOutput{
		ID:                       id,
		ARN:                      bucketARN,
		Name:                     id,
		ForceDestroy:             false,
		BucketDomainName:         bucketDomainName,
		BucketRegionalDomainName: bucketRegionalDomainName,
		Region:                   "us-east-1",
		Tags: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			{Key: "Name", Value: id},
		},
	}, nil
}

func (s *StorageService) UpdateS3Bucket(ctx context.Context, id string, bucket *awss3.Bucket) (*awss3outputs.BucketOutput, error) {
	return s.CreateS3Bucket(ctx, bucket)
}

func (s *StorageService) DeleteS3Bucket(ctx context.Context, id string) error {
	return nil
}

func (s *StorageService) ListS3Buckets(ctx context.Context, filters map[string][]string) ([]*awss3outputs.BucketOutput, error) {
	return []*awss3outputs.BucketOutput{
		{
			ID:                       "test-bucket",
			ARN:                      "arn:aws:s3:::test-bucket",
			Name:                     "test-bucket",
			ForceDestroy:             false,
			BucketDomainName:         "test-bucket.s3.amazonaws.com",
			BucketRegionalDomainName: "test-bucket.s3.us-east-1.amazonaws.com",
			Region:                   "us-east-1",
			Tags: []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}{
				{Key: "Name", Value: "test-bucket"},
			},
		},
	}, nil
}

// S3 Bucket ACL operations

func (s *StorageService) UpdateS3BucketACL(ctx context.Context, bucket string, acl *awss3.BucketACL) (*awss3outputs.BucketACLOutput, error) {
	if acl == nil {
		return nil, fmt.Errorf("acl is nil")
	}

	return &awss3outputs.BucketACLOutput{
		ID:                  bucket,
		ACL:                 acl.ACL,
		AccessControlPolicy: acl.AccessControlPolicy,
	}, nil
}

func (s *StorageService) GetS3BucketACL(ctx context.Context, bucket string) (*awss3outputs.BucketACLOutput, error) {
	canned := "private"
	return &awss3outputs.BucketACLOutput{
		ID:  bucket,
		ACL: &canned,
	}, nil
}

// S3 Bucket Versioning operations

func (s *StorageService) UpdateS3BucketVersioning(ctx context.Context, bucket string, versioning *awss3.BucketVersioning) (*awss3outputs.BucketVersioningOutput, error) {
	if versioning == nil {
		return nil, fmt.Errorf("versioning is nil")
	}

	return &awss3outputs.BucketVersioningOutput{
		ID:        bucket,
		Status:    versioning.Status,
		MFADelete: versioning.MFADelete,
	}, nil
}

func (s *StorageService) GetS3BucketVersioning(ctx context.Context, bucket string) (*awss3outputs.BucketVersioningOutput, error) {
	status := "Enabled"
	return &awss3outputs.BucketVersioningOutput{
		ID:     bucket,
		Status: status,
	}, nil
}

// S3 Bucket Encryption operations

func (s *StorageService) UpdateS3BucketEncryption(ctx context.Context, bucket string, encryption *awss3.BucketEncryption) (*awss3outputs.BucketEncryptionOutput, error) {
	if encryption == nil {
		return nil, fmt.Errorf("encryption is nil")
	}

	return &awss3outputs.BucketEncryptionOutput{
		ID:               bucket,
		BucketKeyEnabled: encryption.Rule.BucketKeyEnabled,
		SSEAlgorithm:     encryption.Rule.DefaultEncryption.SSEAlgorithm,
		KMSMasterKeyID:   encryption.Rule.DefaultEncryption.KMSMasterKeyID,
	}, nil
}

func (s *StorageService) GetS3BucketEncryption(ctx context.Context, bucket string) (*awss3outputs.BucketEncryptionOutput, error) {
	return &awss3outputs.BucketEncryptionOutput{
		ID:               bucket,
		BucketKeyEnabled: false,
		SSEAlgorithm:     "AES256",
	}, nil
}
