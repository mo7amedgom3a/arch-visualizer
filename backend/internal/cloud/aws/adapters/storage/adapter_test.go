package storage

import (
	"context"
	"errors"
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

// mockAWSStorageService is a mock implementation of AWSStorageService for testing
type mockAWSStorageService struct {
	volume         *awsebs.Volume
	bucket         *awss3.Bucket
	bucketACL      *awss3.BucketACL
	bucketVersioning *awss3.BucketVersioning
	bucketEncryption *awss3.BucketEncryption
	createError    error
	getError       error
	createS3Error  error
	getS3Error     error
	aclError       error
	versioningError error
	encryptionError error
}

// Ensure mockAWSStorageService implements AWSStorageService
var _ awsservice.AWSStorageService = (*mockAWSStorageService)(nil)

// Helper function to convert Volume input to output
func volumeToOutput(volume *awsebs.Volume) *awsebsoutputs.VolumeOutput {
	if volume == nil {
		return nil
	}
	return &awsebsoutputs.VolumeOutput{
		ID:               "vol-mock-1234567890abcdef0",
		ARN:              "arn:aws:ec2:us-east-1:123456789012:volume/vol-mock-1234567890abcdef0",
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
	}
}

// EBS Volume Operations

func (m *mockAWSStorageService) CreateEBSVolume(ctx context.Context, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	m.volume = volume
	return volumeToOutput(volume), nil
}

func (m *mockAWSStorageService) GetEBSVolume(ctx context.Context, id string) (*awsebsoutputs.VolumeOutput, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return volumeToOutput(m.volume), nil
}

func (m *mockAWSStorageService) UpdateEBSVolume(ctx context.Context, id string, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error) {
	m.volume = volume
	return volumeToOutput(volume), nil
}

func (m *mockAWSStorageService) DeleteEBSVolume(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSStorageService) ListEBSVolumes(ctx context.Context, filters map[string][]string) ([]*awsebsoutputs.VolumeOutput, error) {
	if m.volume != nil {
		return []*awsebsoutputs.VolumeOutput{volumeToOutput(m.volume)}, nil
	}
	return []*awsebsoutputs.VolumeOutput{}, nil
}

// Volume Attachment Operations

func (m *mockAWSStorageService) AttachVolume(ctx context.Context, volumeID, instanceID, deviceName string) error {
	return nil
}

func (m *mockAWSStorageService) DetachVolume(ctx context.Context, volumeID, instanceID string) error {
	return nil
}

// S3 Bucket Operations

// Helper function to convert Bucket input to output
func bucketToOutput(bucket *awss3.Bucket, region string) *awss3outputs.BucketOutput {
	if bucket == nil {
		return nil
	}

	bucketName := ""
	if bucket.Bucket != nil && *bucket.Bucket != "" {
		bucketName = *bucket.Bucket
	} else if bucket.BucketPrefix != nil && *bucket.BucketPrefix != "" {
		// Generate a mock bucket name with prefix
		bucketName = *bucket.BucketPrefix + "mock1234567890abcdef"
	}

	bucketARN := "arn:aws:s3:::" + bucketName
	bucketDomainName := bucketName + ".s3.amazonaws.com"
	bucketRegionalDomainName := bucketName + ".s3." + region + ".amazonaws.com"

	output := &awss3outputs.BucketOutput{
		ID:                       bucketName,
		ARN:                      bucketARN,
		Name:                     bucketName,
		NamePrefix:               bucket.BucketPrefix,
		ForceDestroy:             bucket.ForceDestroy,
		BucketDomainName:         bucketDomainName,
		BucketRegionalDomainName: bucketRegionalDomainName,
		Region:                   region,
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

	return output
}

func (m *mockAWSStorageService) CreateS3Bucket(ctx context.Context, bucket *awss3.Bucket) (*awss3outputs.BucketOutput, error) {
	if m.createS3Error != nil {
		return nil, m.createS3Error
	}
	m.bucket = bucket
	// Use us-east-1 as default region for mock
	return bucketToOutput(bucket, "us-east-1"), nil
}

func (m *mockAWSStorageService) GetS3Bucket(ctx context.Context, id string) (*awss3outputs.BucketOutput, error) {
	if m.getS3Error != nil {
		return nil, m.getS3Error
	}
	if m.bucket == nil {
		return nil, errors.New("bucket not found")
	}
	return bucketToOutput(m.bucket, "us-east-1"), nil
}

func (m *mockAWSStorageService) UpdateS3Bucket(ctx context.Context, id string, bucket *awss3.Bucket) (*awss3outputs.BucketOutput, error) {
	m.bucket = bucket
	return bucketToOutput(bucket, "us-east-1"), nil
}

func (m *mockAWSStorageService) DeleteS3Bucket(ctx context.Context, id string) error {
	return nil
}

func (m *mockAWSStorageService) ListS3Buckets(ctx context.Context, filters map[string][]string) ([]*awss3outputs.BucketOutput, error) {
	if m.bucket != nil {
		return []*awss3outputs.BucketOutput{bucketToOutput(m.bucket, "us-east-1")}, nil
	}
	return []*awss3outputs.BucketOutput{}, nil
}

// S3 Bucket ACL Operations

func (m *mockAWSStorageService) UpdateS3BucketACL(ctx context.Context, bucket string, acl *awss3.BucketACL) (*awss3outputs.BucketACLOutput, error) {
	if m.aclError != nil {
		return nil, m.aclError
	}
	m.bucketACL = acl
	return &awss3outputs.BucketACLOutput{
		ID:                  bucket,
		ACL:                 acl.ACL,
		AccessControlPolicy: acl.AccessControlPolicy,
	}, nil
}

func (m *mockAWSStorageService) GetS3BucketACL(ctx context.Context, bucket string) (*awss3outputs.BucketACLOutput, error) {
	if m.bucketACL == nil {
		return nil, errors.New("bucket acl not found")
	}
	return &awss3outputs.BucketACLOutput{
		ID:                  bucket,
		ACL:                 m.bucketACL.ACL,
		AccessControlPolicy: m.bucketACL.AccessControlPolicy,
	}, nil
}

// S3 Bucket Versioning Operations

func (m *mockAWSStorageService) UpdateS3BucketVersioning(ctx context.Context, bucket string, versioning *awss3.BucketVersioning) (*awss3outputs.BucketVersioningOutput, error) {
	if m.versioningError != nil {
		return nil, m.versioningError
	}
	m.bucketVersioning = versioning
	return &awss3outputs.BucketVersioningOutput{
		ID:        bucket,
		Status:    versioning.Status,
		MFADelete: versioning.MFADelete,
	}, nil
}

func (m *mockAWSStorageService) GetS3BucketVersioning(ctx context.Context, bucket string) (*awss3outputs.BucketVersioningOutput, error) {
	if m.bucketVersioning == nil {
		return nil, errors.New("bucket versioning not found")
	}
	return &awss3outputs.BucketVersioningOutput{
		ID:        bucket,
		Status:    m.bucketVersioning.Status,
		MFADelete: m.bucketVersioning.MFADelete,
	}, nil
}

// S3 Bucket Encryption Operations

func (m *mockAWSStorageService) UpdateS3BucketEncryption(ctx context.Context, bucket string, encryption *awss3.BucketEncryption) (*awss3outputs.BucketEncryptionOutput, error) {
	if m.encryptionError != nil {
		return nil, m.encryptionError
	}
	m.bucketEncryption = encryption
	return &awss3outputs.BucketEncryptionOutput{
		ID:               bucket,
		BucketKeyEnabled: encryption.Rule.BucketKeyEnabled,
		SSEAlgorithm:     encryption.Rule.DefaultEncryption.SSEAlgorithm,
		KMSMasterKeyID:   encryption.Rule.DefaultEncryption.KMSMasterKeyID,
	}, nil
}

func (m *mockAWSStorageService) GetS3BucketEncryption(ctx context.Context, bucket string) (*awss3outputs.BucketEncryptionOutput, error) {
	if m.bucketEncryption == nil {
		return nil, errors.New("bucket encryption not found")
	}
	return &awss3outputs.BucketEncryptionOutput{
		ID:               bucket,
		BucketKeyEnabled: m.bucketEncryption.Rule.BucketKeyEnabled,
		SSEAlgorithm:     m.bucketEncryption.Rule.DefaultEncryption.SSEAlgorithm,
		KMSMasterKeyID:   m.bucketEncryption.Rule.DefaultEncryption.KMSMasterKeyID,
	}, nil
}

// Tests

func TestAWSStorageAdapter_CreateEBSVolume(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	domainVolume := &domainstorage.EBSVolume{
		Name:             "test-volume",
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

	if createdVolume.Name != domainVolume.Name {
		t.Errorf("Expected name %s, got %s", domainVolume.Name, createdVolume.Name)
	}

	if createdVolume.ID == "" {
		t.Error("Expected volume ID to be populated")
	}
}

func TestAWSStorageAdapter_CreateEBSVolume_ValidationError(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	invalidVolume := &domainstorage.EBSVolume{
		Name:   "", // Invalid: empty name
		Region: "us-east-1",
		Size:   40,
		Type:   "gp3",
	}

	ctx := context.Background()
	_, err := adapter.CreateEBSVolume(ctx, invalidVolume)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
}

func TestAWSStorageAdapter_GetEBSVolume(t *testing.T) {
	mockService := &mockAWSStorageService{
		volume: &awsebs.Volume{
			Name:             "test-volume",
			AvailabilityZone: "us-east-1a",
			Size:             40,
			VolumeType:       "gp3",
			Encrypted:        false,
		},
	}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	volume, err := adapter.GetEBSVolume(ctx, "vol-123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if volume == nil {
		t.Fatal("Expected volume, got nil")
	}

	if volume.Name != "test-volume" {
		t.Errorf("Expected name test-volume, got %s", volume.Name)
	}
}

func TestAWSStorageAdapter_ListEBSVolumes(t *testing.T) {
	mockService := &mockAWSStorageService{
		volume: &awsebs.Volume{
			Name:             "test-volume",
			AvailabilityZone: "us-east-1a",
			Size:             40,
			VolumeType:       "gp3",
			Encrypted:        false,
		},
	}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	volumes, err := adapter.ListEBSVolumes(ctx, map[string]string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(volumes) != 1 {
		t.Errorf("Expected 1 volume, got %d", len(volumes))
	}

	if volumes[0].Name != "test-volume" {
		t.Errorf("Expected name test-volume, got %s", volumes[0].Name)
	}
}

func TestAWSStorageAdapter_AttachVolume(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	err := adapter.AttachVolume(ctx, "vol-123", "i-123", "/dev/xvdf")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAWSStorageAdapter_DetachVolume(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	err := adapter.DetachVolume(ctx, "vol-123", "i-123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAWSStorageAdapter_ErrorHandling(t *testing.T) {
	mockService := &mockAWSStorageService{
		getError: errors.New("aws service error"),
	}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	_, err := adapter.GetEBSVolume(ctx, "vol-123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Verify error is wrapped
	if err.Error() == "" {
		t.Error("Expected error message, got empty string")
	}
}

// S3 Bucket Tests

func TestAWSStorageAdapter_CreateS3Bucket(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	domainBucket := &domainstorage.S3Bucket{
		Name:         "test-bucket",
		Region:       "us-east-1",
		ForceDestroy: false,
		Tags: map[string]string{
			"Environment": "test",
		},
	}

	ctx := context.Background()
	createdBucket, err := adapter.CreateS3Bucket(ctx, domainBucket)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdBucket == nil {
		t.Fatal("Expected created bucket, got nil")
	}

	if createdBucket.Name != domainBucket.Name {
		t.Errorf("Expected name %s, got %s", domainBucket.Name, createdBucket.Name)
	}

	if createdBucket.ID == "" {
		t.Error("Expected bucket ID to be populated")
	}

	if createdBucket.ARN == nil || *createdBucket.ARN == "" {
		t.Error("Expected bucket ARN to be populated")
	}

	if createdBucket.BucketDomainName == nil || *createdBucket.BucketDomainName == "" {
		t.Error("Expected bucket domain name to be populated")
	}

	if createdBucket.BucketRegionalDomainName == nil || *createdBucket.BucketRegionalDomainName == "" {
		t.Error("Expected bucket regional domain name to be populated")
	}
}

func TestAWSStorageAdapter_CreateS3Bucket_WithPrefix(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	prefix := "my-app-logs-"
	domainBucket := &domainstorage.S3Bucket{
		NamePrefix:   &prefix,
		Region:       "us-east-1",
		ForceDestroy: false,
	}

	ctx := context.Background()
	createdBucket, err := adapter.CreateS3Bucket(ctx, domainBucket)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if createdBucket == nil {
		t.Fatal("Expected created bucket, got nil")
	}

	if createdBucket.NamePrefix == nil || *createdBucket.NamePrefix != prefix {
		t.Errorf("Expected name prefix %s, got %v", prefix, createdBucket.NamePrefix)
	}
}

func TestAWSStorageAdapter_CreateS3Bucket_ValidationError(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	invalidBucket := &domainstorage.S3Bucket{
		Name:   "", // Invalid: empty name
		Region: "", // Invalid: empty region
	}

	ctx := context.Background()
	_, err := adapter.CreateS3Bucket(ctx, invalidBucket)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
}

func TestAWSStorageAdapter_GetS3Bucket(t *testing.T) {
	bucketName := "test-bucket"
	mockService := &mockAWSStorageService{
		bucket: &awss3.Bucket{
			Bucket: &bucketName,
		},
	}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	bucket, err := adapter.GetS3Bucket(ctx, "test-bucket")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if bucket == nil {
		t.Fatal("Expected bucket, got nil")
	}

	if bucket.Name != bucketName {
		t.Errorf("Expected name %s, got %s", bucketName, bucket.Name)
	}
}

func TestAWSStorageAdapter_ListS3Buckets(t *testing.T) {
	bucketName := "test-bucket"
	mockService := &mockAWSStorageService{
		bucket: &awss3.Bucket{
			Bucket: &bucketName,
		},
	}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	buckets, err := adapter.ListS3Buckets(ctx, map[string]string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(buckets) != 1 {
		t.Errorf("Expected 1 bucket, got %d", len(buckets))
	}

	if buckets[0].Name != bucketName {
		t.Errorf("Expected name %s, got %s", bucketName, buckets[0].Name)
	}
}

func TestAWSStorageAdapter_DeleteS3Bucket(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	err := adapter.DeleteS3Bucket(ctx, "test-bucket")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAWSStorageAdapter_S3Bucket_ErrorHandling(t *testing.T) {
	mockService := &mockAWSStorageService{
		getS3Error: errors.New("aws service error"),
	}
	adapter := NewAWSStorageAdapter(mockService)

	ctx := context.Background()
	_, err := adapter.GetS3Bucket(ctx, "test-bucket")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Verify error is wrapped
	if err.Error() == "" {
		t.Error("Expected error message, got empty string")
	}
}

func TestAWSStorageAdapter_UpdateS3BucketACL(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	canned := "private"
	domainACL := &domainstorage.S3BucketACL{
		Bucket: "test-bucket",
		ACL:    &canned,
	}

	ctx := context.Background()
	result, err := adapter.UpdateS3BucketACL(ctx, "test-bucket", domainACL)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil || result.ACL == nil || *result.ACL != "private" {
		t.Fatalf("Expected ACL private, got %v", result)
	}
}

func TestAWSStorageAdapter_UpdateS3BucketVersioning(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	domainVersioning := &domainstorage.S3BucketVersioning{
		Bucket: "test-bucket",
		Status: "Enabled",
	}

	ctx := context.Background()
	result, err := adapter.UpdateS3BucketVersioning(ctx, "test-bucket", domainVersioning)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil || result.Status != "Enabled" {
		t.Fatalf("Expected status Enabled, got %v", result)
	}
}

func TestAWSStorageAdapter_UpdateS3BucketEncryption(t *testing.T) {
	mockService := &mockAWSStorageService{}
	adapter := NewAWSStorageAdapter(mockService)

	kmsKey := "arn:aws:kms:us-east-1:123456789012:key/abc"
	domainEncryption := &domainstorage.S3BucketEncryption{
		Bucket: "test-bucket",
		Rule: domainstorage.S3BucketEncryptionRule{
			BucketKeyEnabled: true,
			DefaultEncryption: domainstorage.S3BucketDefaultEncryption{
				SSEAlgorithm:   "aws:kms",
				KMSMasterKeyID: &kmsKey,
			},
		},
	}

	ctx := context.Background()
	result, err := adapter.UpdateS3BucketEncryption(ctx, "test-bucket", domainEncryption)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil || result.Rule.DefaultEncryption.SSEAlgorithm != "aws:kms" {
		t.Fatalf("Expected aws:kms algorithm, got %v", result)
	}
	if result.Rule.DefaultEncryption.KMSMasterKeyID == nil || *result.Rule.DefaultEncryption.KMSMasterKeyID != kmsKey {
		t.Fatalf("Expected kms key %s, got %v", kmsKey, result.Rule.DefaultEncryption.KMSMasterKeyID)
	}
}
