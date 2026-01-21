package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	awsebs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs"
	awsebsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs/outputs"
	awss3 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3"
	awss3outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3/outputs"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/storage"
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
)

// realisticAWSStorageService is a realistic implementation that returns proper output models
type realisticAWSStorageService struct{}

var _ awsservice.AWSStorageService = (*realisticAWSStorageService)(nil)

// EBS Volume Operations

func (s *realisticAWSStorageService) CreateEBSVolume(ctx context.Context, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error) {
	// Simulate realistic AWS EBS volume creation
	return &awsebsoutputs.VolumeOutput{
		ID:               "vol-0a1b2c3d4e5f6g7h8",
		ARN:              "arn:aws:ec2:us-east-1:123456789012:volume/vol-0a1b2c3d4e5f6g7h8",
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
		CreateTime:       time.Now(),
		Tags: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			{Key: "Name", Value: volume.Name},
		},
	}, nil
}

func (s *realisticAWSStorageService) GetEBSVolume(ctx context.Context, id string) (*awsebsoutputs.VolumeOutput, error) {
	return &awsebsoutputs.VolumeOutput{
		ID:               id,
		ARN:              fmt.Sprintf("arn:aws:ec2:us-east-1:123456789012:volume/%s", id),
		Name:             "test-volume",
		AvailabilityZone: "us-east-1a",
		Size:             40,
		VolumeType:       "gp3",
		IOPS:             intPtr(3000),
		Throughput:       intPtr(125),
		Encrypted:        false,
		KMSKeyID:         nil,
		SnapshotID:       nil,
		State:            "available",
		AttachedTo:       nil,
		CreateTime:       time.Now(),
		Tags: []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{
			{Key: "Name", Value: "test-volume"},
		},
	}, nil
}

func (s *realisticAWSStorageService) UpdateEBSVolume(ctx context.Context, id string, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error) {
	return s.CreateEBSVolume(context.Background(), volume)
}

func (s *realisticAWSStorageService) DeleteEBSVolume(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSStorageService) ListEBSVolumes(ctx context.Context, filters map[string][]string) ([]*awsebsoutputs.VolumeOutput, error) {
	return []*awsebsoutputs.VolumeOutput{
		{
			ID:               "vol-0a1b2c3d4e5f6g7h8",
			ARN:              "arn:aws:ec2:us-east-1:123456789012:volume/vol-0a1b2c3d4e5f6g7h8",
			Name:             "test-volume",
			AvailabilityZone: "us-east-1a",
			Size:             40,
			VolumeType:       "gp3",
			IOPS:             intPtr(3000),
			Throughput:       intPtr(125),
			Encrypted:        false,
			KMSKeyID:         nil,
			SnapshotID:       nil,
			State:            "available",
			AttachedTo:       nil,
			CreateTime:       time.Now(),
			Tags: []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}{
				{Key: "Name", Value: "test-volume"},
			},
		},
	}, nil
}

// Volume Attachment Operations

func (s *realisticAWSStorageService) AttachVolume(ctx context.Context, volumeID, instanceID, deviceName string) error {
	return nil
}

func (s *realisticAWSStorageService) DetachVolume(ctx context.Context, volumeID, instanceID string) error {
	return nil
}

// S3 Bucket Operations

func (s *realisticAWSStorageService) CreateS3Bucket(ctx context.Context, bucket *awss3.Bucket) (*awss3outputs.BucketOutput, error) {
	bucketName := ""
	if bucket.Bucket != nil && *bucket.Bucket != "" {
		bucketName = *bucket.Bucket
	} else if bucket.BucketPrefix != nil && *bucket.BucketPrefix != "" {
		bucketName = *bucket.BucketPrefix + "mock1234567890abcdef"
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

func (s *realisticAWSStorageService) GetS3Bucket(ctx context.Context, id string) (*awss3outputs.BucketOutput, error) {
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

func (s *realisticAWSStorageService) UpdateS3Bucket(ctx context.Context, id string, bucket *awss3.Bucket) (*awss3outputs.BucketOutput, error) {
	return s.CreateS3Bucket(context.Background(), bucket)
}

func (s *realisticAWSStorageService) DeleteS3Bucket(ctx context.Context, id string) error {
	return nil
}

func (s *realisticAWSStorageService) ListS3Buckets(ctx context.Context, filters map[string][]string) ([]*awss3outputs.BucketOutput, error) {
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

// S3 Bucket ACL Operations
func (s *realisticAWSStorageService) UpdateS3BucketACL(ctx context.Context, bucket string, acl *awss3.BucketACL) (*awss3outputs.BucketACLOutput, error) {
	return &awss3outputs.BucketACLOutput{
		ID:                  bucket,
		ACL:                 acl.ACL,
		AccessControlPolicy: acl.AccessControlPolicy,
	}, nil
}

func (s *realisticAWSStorageService) GetS3BucketACL(ctx context.Context, bucket string) (*awss3outputs.BucketACLOutput, error) {
	canned := "private"
	return &awss3outputs.BucketACLOutput{
		ID:  bucket,
		ACL: &canned,
	}, nil
}

// S3 Bucket Versioning Operations
func (s *realisticAWSStorageService) UpdateS3BucketVersioning(ctx context.Context, bucket string, versioning *awss3.BucketVersioning) (*awss3outputs.BucketVersioningOutput, error) {
	return &awss3outputs.BucketVersioningOutput{
		ID:        bucket,
		Status:    versioning.Status,
		MFADelete: versioning.MFADelete,
	}, nil
}

func (s *realisticAWSStorageService) GetS3BucketVersioning(ctx context.Context, bucket string) (*awss3outputs.BucketVersioningOutput, error) {
	status := "Enabled"
	return &awss3outputs.BucketVersioningOutput{
		ID:     bucket,
		Status: status,
	}, nil
}

// S3 Bucket Encryption Operations
func (s *realisticAWSStorageService) UpdateS3BucketEncryption(ctx context.Context, bucket string, encryption *awss3.BucketEncryption) (*awss3outputs.BucketEncryptionOutput, error) {
	return &awss3outputs.BucketEncryptionOutput{
		ID:               bucket,
		BucketKeyEnabled: encryption.Rule.BucketKeyEnabled,
		SSEAlgorithm:     encryption.Rule.DefaultEncryption.SSEAlgorithm,
		KMSMasterKeyID:   encryption.Rule.DefaultEncryption.KMSMasterKeyID,
	}, nil
}

func (s *realisticAWSStorageService) GetS3BucketEncryption(ctx context.Context, bucket string) (*awss3outputs.BucketEncryptionOutput, error) {
	return &awss3outputs.BucketEncryptionOutput{
		ID:               bucket,
		BucketKeyEnabled: false,
		SSEAlgorithm:     "AES256",
	}, nil
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

// Integration Tests

func TestAWSStorageAdapter_OutputIntegration_CreateEBSVolume(t *testing.T) {
	realisticService := &realisticAWSStorageService{}
	adapter := NewAWSStorageAdapter(realisticService)

	domainVolume := &domainstorage.EBSVolume{
		Name:             "integration-test-volume",
		Region:           "us-east-1",
		AvailabilityZone: "us-east-1a",
		Size:             40,
		Type:             "gp3",
		Encrypted:        false,
	}

	ctx := context.Background()
	createdVolume, err := adapter.CreateEBSVolume(ctx, domainVolume)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdVolume == nil {
		t.Fatal("Expected created volume, got nil")
	}

	// Verify AWS-generated identifiers are populated
	if createdVolume.ID == "" {
		t.Error("Expected volume ID to be populated")
	}

	if createdVolume.ID != "vol-0a1b2c3d4e5f6g7h8" {
		t.Errorf("Expected volume ID vol-0a1b2c3d4e5f6g7h8, got %s", createdVolume.ID)
	}

	if createdVolume.ARN == nil {
		t.Error("Expected volume ARN to be populated")
	}

	if createdVolume.ARN != nil && *createdVolume.ARN == "" {
		t.Error("Expected volume ARN to be non-empty")
	}

	// Verify state is populated
	if createdVolume.State != domainstorage.EBSVolumeStateAvailable {
		t.Errorf("Expected state available, got %s", createdVolume.State)
	}
}

func TestAWSStorageAdapter_OutputIntegration_GetEBSVolume(t *testing.T) {
	realisticService := &realisticAWSStorageService{}
	adapter := NewAWSStorageAdapter(realisticService)

	ctx := context.Background()
	volume, err := adapter.GetEBSVolume(ctx, "vol-0a1b2c3d4e5f6g7h8")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if volume == nil {
		t.Fatal("Expected volume, got nil")
	}

	// Verify all output fields are populated
	if volume.ID != "vol-0a1b2c3d4e5f6g7h8" {
		t.Errorf("Expected ID vol-0a1b2c3d4e5f6g7h8, got %s", volume.ID)
	}

	if volume.ARN == nil {
		t.Error("Expected ARN to be populated")
	}

	if volume.State == "" {
		t.Error("Expected state to be populated")
	}

	if volume.Size == 0 {
		t.Error("Expected size to be populated")
	}
}

func TestAWSStorageAdapter_OutputIntegration_ListEBSVolumes(t *testing.T) {
	realisticService := &realisticAWSStorageService{}
	adapter := NewAWSStorageAdapter(realisticService)

	ctx := context.Background()
	volumes, err := adapter.ListEBSVolumes(ctx, map[string]string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(volumes) == 0 {
		t.Fatal("Expected at least one volume, got none")
	}

	// Verify first volume has all output fields
	volume := volumes[0]
	if volume.ID == "" {
		t.Error("Expected volume ID to be populated")
	}

	if volume.ARN == nil {
		t.Error("Expected volume ARN to be populated")
	}

	if volume.State == "" {
		t.Error("Expected volume state to be populated")
	}
}

func TestAWSStorageAdapter_OutputIntegration_UpdateEBSVolume(t *testing.T) {
	realisticService := &realisticAWSStorageService{}
	adapter := NewAWSStorageAdapter(realisticService)

	domainVolume := &domainstorage.EBSVolume{
		Name:             "updated-volume",
		Region:           "us-east-1",
		AvailabilityZone: "us-east-1a",
		Size:             80, // Increased size
		Type:             "gp3",
		Encrypted:        true,
	}

	ctx := context.Background()
	updatedVolume, err := adapter.UpdateEBSVolume(ctx, "vol-0a1b2c3d4e5f6g7h8", domainVolume)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if updatedVolume == nil {
		t.Fatal("Expected updated volume, got nil")
	}

	// Verify updated fields
	if updatedVolume.Size != 80 {
		t.Errorf("Expected size 80, got %d", updatedVolume.Size)
	}

	if !updatedVolume.Encrypted {
		t.Error("Expected encrypted to be true")
	}
}

func TestAWSStorageAdapter_OutputIntegration_AttachVolume(t *testing.T) {
	realisticService := &realisticAWSStorageService{}
	adapter := NewAWSStorageAdapter(realisticService)

	ctx := context.Background()
	err := adapter.AttachVolume(ctx, "vol-0a1b2c3d4e5f6g7h8", "i-0a1b2c3d4e5f6g7h8", "/dev/xvdf")

	if err != nil {
		t.Fatalf("Expected no error attaching volume, got: %v", err)
	}
}

func TestAWSStorageAdapter_OutputIntegration_DetachVolume(t *testing.T) {
	realisticService := &realisticAWSStorageService{}
	adapter := NewAWSStorageAdapter(realisticService)

	ctx := context.Background()
	err := adapter.DetachVolume(ctx, "vol-0a1b2c3d4e5f6g7h8", "i-0a1b2c3d4e5f6g7h8")

	if err != nil {
		t.Fatalf("Expected no error detaching volume, got: %v", err)
	}
}

func TestAWSStorageAdapter_OutputIntegration_VolumeWithIOPS(t *testing.T) {
	realisticService := &realisticAWSStorageService{}
	adapter := NewAWSStorageAdapter(realisticService)

	iops := 6000
	domainVolume := &domainstorage.EBSVolume{
		Name:             "high-iops-volume",
		Region:           "us-east-1",
		AvailabilityZone: "us-east-1a",
		Size:             100,
		Type:             "gp3",
		IOPS:              &iops,
		Encrypted:        false,
	}

	ctx := context.Background()
	createdVolume, err := adapter.CreateEBSVolume(ctx, domainVolume)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdVolume == nil {
		t.Fatal("Expected created volume, got nil")
	}

	// Verify IOPS is preserved
	if createdVolume.IOPS == nil {
		t.Error("Expected IOPS to be populated")
	}

	if createdVolume.IOPS != nil && *createdVolume.IOPS != iops {
		t.Errorf("Expected IOPS %d, got %d", iops, *createdVolume.IOPS)
	}
}

func TestAWSStorageAdapter_OutputIntegration_VolumeWithThroughput(t *testing.T) {
	realisticService := &realisticAWSStorageService{}
	adapter := NewAWSStorageAdapter(realisticService)

	throughput := 250
	domainVolume := &domainstorage.EBSVolume{
		Name:             "high-throughput-volume",
		Region:           "us-east-1",
		AvailabilityZone: "us-east-1a",
		Size:             100,
		Type:             "gp3",
		Throughput:       &throughput,
		Encrypted:        false,
	}

	ctx := context.Background()
	createdVolume, err := adapter.CreateEBSVolume(ctx, domainVolume)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdVolume == nil {
		t.Fatal("Expected created volume, got nil")
	}

	// Verify throughput is preserved
	if createdVolume.Throughput == nil {
		t.Error("Expected throughput to be populated")
	}

	if createdVolume.Throughput != nil && *createdVolume.Throughput != throughput {
		t.Errorf("Expected throughput %d, got %d", throughput, *createdVolume.Throughput)
	}
}

func TestAWSStorageAdapter_OutputIntegration_EncryptedVolume(t *testing.T) {
	realisticService := &realisticAWSStorageService{}
	adapter := NewAWSStorageAdapter(realisticService)

	kmsKeyID := "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
	domainVolume := &domainstorage.EBSVolume{
		Name:             "encrypted-volume",
		Region:           "us-east-1",
		AvailabilityZone: "us-east-1a",
		Size:             40,
		Type:             "gp3",
		Encrypted:        true,
		KMSKeyID:         &kmsKeyID,
	}

	ctx := context.Background()
	createdVolume, err := adapter.CreateEBSVolume(ctx, domainVolume)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdVolume == nil {
		t.Fatal("Expected created volume, got nil")
	}

	// Verify encryption settings
	if !createdVolume.Encrypted {
		t.Error("Expected encrypted to be true")
	}

	if createdVolume.KMSKeyID == nil {
		t.Error("Expected KMS key ID to be populated")
	}

	if createdVolume.KMSKeyID != nil && *createdVolume.KMSKeyID != kmsKeyID {
		t.Errorf("Expected KMS key ID %s, got %s", kmsKeyID, *createdVolume.KMSKeyID)
	}
}

func TestAWSStorageAdapter_OutputIntegration_VolumeFromSnapshot(t *testing.T) {
	realisticService := &realisticAWSStorageService{}
	adapter := NewAWSStorageAdapter(realisticService)

	snapshotID := "snap-0a1b2c3d4e5f6g7h8"
	domainVolume := &domainstorage.EBSVolume{
		Name:             "volume-from-snapshot",
		Region:           "us-east-1",
		AvailabilityZone: "us-east-1a",
		Size:             40,
		Type:             "gp3",
		SnapshotID:       &snapshotID,
		Encrypted:        false,
	}

	ctx := context.Background()
	createdVolume, err := adapter.CreateEBSVolume(ctx, domainVolume)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdVolume == nil {
		t.Fatal("Expected created volume, got nil")
	}

	// Verify snapshot ID is preserved
	if createdVolume.SnapshotID == nil {
		t.Error("Expected snapshot ID to be populated")
	}

	if createdVolume.SnapshotID != nil && *createdVolume.SnapshotID != snapshotID {
		t.Errorf("Expected snapshot ID %s, got %s", snapshotID, *createdVolume.SnapshotID)
	}
}
