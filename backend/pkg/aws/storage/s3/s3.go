package s3

import (
	"context"
	"fmt"

	awss3adapter "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/adapters/storage"
	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/storage"
	awsebs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs"
	awsebsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs/outputs"
	awss3 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3"
	awss3outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3/outputs"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/storage"
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
)

// inMemoryAWSStorageService is an in-memory implementation of AWSStorageService.
// It is used by the S3Runner to simulate bucket operations without calling real AWS.
type inMemoryAWSStorageService struct {
	// EBS (unused, but required to satisfy interface)
	volumes map[string]*awsebsoutputs.VolumeOutput

	// S3 state
	buckets    map[string]*awss3outputs.BucketOutput
	acls       map[string]*awss3outputs.BucketACLOutput
	versioning map[string]*awss3outputs.BucketVersioningOutput
	encryption map[string]*awss3outputs.BucketEncryptionOutput
}

// Ensure inMemoryAWSStorageService implements AWSStorageService
var _ awsservice.AWSStorageService = (*inMemoryAWSStorageService)(nil)

// newInMemoryAWSStorageService creates a new in-memory service instance.
func newInMemoryAWSStorageService() *inMemoryAWSStorageService {
	return &inMemoryAWSStorageService{
		volumes:    make(map[string]*awsebsoutputs.VolumeOutput),
		buckets:    make(map[string]*awss3outputs.BucketOutput),
		acls:       make(map[string]*awss3outputs.BucketACLOutput),
		versioning: make(map[string]*awss3outputs.BucketVersioningOutput),
		encryption: make(map[string]*awss3outputs.BucketEncryptionOutput),
	}
}

// newS3DemoAdapter constructs a domain StorageService backed by the in-memory AWS storage service.
func newS3DemoAdapter() domainstorage.StorageService {
	service := newInMemoryAWSStorageService()
	return awss3adapter.NewAWSStorageAdapter(service)
}

// =========================
// EBS methods (no-op stubs)
// =========================

func (s *inMemoryAWSStorageService) CreateEBSVolume(ctx context.Context, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error) {
	if volume == nil {
		return nil, fmt.Errorf("volume is nil")
	}
	out := &awsebsoutputs.VolumeOutput{
		ID:               "vol-mem-1234567890",
		ARN:              "arn:aws:ec2:us-east-1:123456789012:volume/vol-mem-1234567890",
		Name:             volume.Name,
		AvailabilityZone: volume.AvailabilityZone,
		Size:             volume.Size,
		VolumeType:       volume.VolumeType,
		Encrypted:        volume.Encrypted,
	}
	s.volumes[out.ID] = out
	return out, nil
}

func (s *inMemoryAWSStorageService) GetEBSVolume(ctx context.Context, id string) (*awsebsoutputs.VolumeOutput, error) {
	if v, ok := s.volumes[id]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("volume %s not found", id)
}

func (s *inMemoryAWSStorageService) UpdateEBSVolume(ctx context.Context, id string, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error) {
	return s.CreateEBSVolume(ctx, volume)
}

func (s *inMemoryAWSStorageService) DeleteEBSVolume(ctx context.Context, id string) error {
	delete(s.volumes, id)
	return nil
}

func (s *inMemoryAWSStorageService) ListEBSVolumes(ctx context.Context, filters map[string][]string) ([]*awsebsoutputs.VolumeOutput, error) {
	results := make([]*awsebsoutputs.VolumeOutput, 0, len(s.volumes))
	for _, v := range s.volumes {
		results = append(results, v)
	}
	return results, nil
}

func (s *inMemoryAWSStorageService) AttachVolume(ctx context.Context, volumeID, instanceID, deviceName string) error {
	return nil
}

func (s *inMemoryAWSStorageService) DetachVolume(ctx context.Context, volumeID, instanceID string) error {
	return nil
}

// =================
// S3 Bucket methods
// =================

func (s *inMemoryAWSStorageService) CreateS3Bucket(ctx context.Context, bucket *awss3.Bucket) (*awss3outputs.BucketOutput, error) {
	if bucket == nil {
		return nil, fmt.Errorf("bucket is nil")
	}
	// Derive name from bucket or prefix
	name := ""
	if bucket.Bucket != nil && *bucket.Bucket != "" {
		name = *bucket.Bucket
	} else if bucket.BucketPrefix != nil && *bucket.BucketPrefix != "" {
		name = *bucket.BucketPrefix + "demo1234567890"
	}
	if name == "" {
		return nil, fmt.Errorf("bucket name or prefix must be provided")
	}

	arn := fmt.Sprintf("arn:aws:s3:::%s", name)
	domain := fmt.Sprintf("%s.s3.amazonaws.com", name)
	regionalDomain := fmt.Sprintf("%s.s3.us-east-1.amazonaws.com", name)

	out := &awss3outputs.BucketOutput{
		ID:                       name,
		ARN:                      arn,
		Name:                     name,
		NamePrefix:               bucket.BucketPrefix,
		ForceDestroy:             bucket.ForceDestroy,
		BucketDomainName:         domain,
		BucketRegionalDomainName: regionalDomain,
		Region:                   "us-east-1",
	}

	// Copy tags
	for _, tag := range bucket.Tags {
		out.Tags = append(out.Tags, struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{Key: tag.Key, Value: tag.Value})
	}

	s.buckets[name] = out
	return out, nil
}

func (s *inMemoryAWSStorageService) GetS3Bucket(ctx context.Context, id string) (*awss3outputs.BucketOutput, error) {
	if b, ok := s.buckets[id]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("bucket %s not found", id)
}

func (s *inMemoryAWSStorageService) UpdateS3Bucket(ctx context.Context, id string, bucket *awss3.Bucket) (*awss3outputs.BucketOutput, error) {
	// For simplicity, reuse CreateS3Bucket logic and overwrite
	out, err := s.CreateS3Bucket(ctx, bucket)
	if err != nil {
		return nil, err
	}
	s.buckets[out.ID] = out
	return out, nil
}

func (s *inMemoryAWSStorageService) DeleteS3Bucket(ctx context.Context, id string) error {
	delete(s.buckets, id)
	delete(s.acls, id)
	delete(s.versioning, id)
	delete(s.encryption, id)
	return nil
}

func (s *inMemoryAWSStorageService) ListS3Buckets(ctx context.Context, filters map[string][]string) ([]*awss3outputs.BucketOutput, error) {
	results := make([]*awss3outputs.BucketOutput, 0, len(s.buckets))
	for _, b := range s.buckets {
		results = append(results, b)
	}
	return results, nil
}

// S3 ACL

func (s *inMemoryAWSStorageService) UpdateS3BucketACL(ctx context.Context, bucket string, acl *awss3.BucketACL) (*awss3outputs.BucketACLOutput, error) {
	if acl == nil {
		return nil, fmt.Errorf("acl is nil")
	}
	out := &awss3outputs.BucketACLOutput{
		ID:                  bucket,
		ACL:                 acl.ACL,
		AccessControlPolicy: acl.AccessControlPolicy,
	}
	s.acls[bucket] = out
	return out, nil
}

func (s *inMemoryAWSStorageService) GetS3BucketACL(ctx context.Context, bucket string) (*awss3outputs.BucketACLOutput, error) {
	if a, ok := s.acls[bucket]; ok {
		return a, nil
	}
	return nil, fmt.Errorf("acl for bucket %s not found", bucket)
}

// S3 Versioning

func (s *inMemoryAWSStorageService) UpdateS3BucketVersioning(ctx context.Context, bucket string, versioning *awss3.BucketVersioning) (*awss3outputs.BucketVersioningOutput, error) {
	if versioning == nil {
		return nil, fmt.Errorf("versioning is nil")
	}
	out := &awss3outputs.BucketVersioningOutput{
		ID:        bucket,
		Status:    versioning.Status,
		MFADelete: versioning.MFADelete,
	}
	s.versioning[bucket] = out
	return out, nil
}

func (s *inMemoryAWSStorageService) GetS3BucketVersioning(ctx context.Context, bucket string) (*awss3outputs.BucketVersioningOutput, error) {
	if v, ok := s.versioning[bucket]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("versioning for bucket %s not found", bucket)
}

// S3 Encryption

func (s *inMemoryAWSStorageService) UpdateS3BucketEncryption(ctx context.Context, bucket string, encryption *awss3.BucketEncryption) (*awss3outputs.BucketEncryptionOutput, error) {
	if encryption == nil {
		return nil, fmt.Errorf("encryption is nil")
	}
	out := &awss3outputs.BucketEncryptionOutput{
		ID:               bucket,
		BucketKeyEnabled: encryption.Rule.BucketKeyEnabled,
		SSEAlgorithm:     encryption.Rule.DefaultEncryption.SSEAlgorithm,
		KMSMasterKeyID:   encryption.Rule.DefaultEncryption.KMSMasterKeyID,
	}
	s.encryption[bucket] = out
	return out, nil
}

func (s *inMemoryAWSStorageService) GetS3BucketEncryption(ctx context.Context, bucket string) (*awss3outputs.BucketEncryptionOutput, error) {
	if e, ok := s.encryption[bucket]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("encryption for bucket %s not found", bucket)
}

// =====================
// Public runner helpers
// =====================

// S3Runner demonstrates S3 bucket operations using the domain storage service and in-memory AWS implementation.
// This is intended to be called from a main() or manual test harness.
func S3Runner() {
	ctx := context.Background()
	storageService := newS3DemoAdapter()

	fmt.Println("============================================")
	fmt.Println("S3 BUCKET STORAGE DEMO (IN-MEMORY)")
	fmt.Println("============================================")

	// 1. Create bucket
	bucket := &domainstorage.S3Bucket{
		Name:         "example-bucket-12345",
		Region:       "us-east-1",
		ForceDestroy: true,
		Tags: map[string]string{
			"Environment": "dev",
			"Service":     "storage-demo",
		},
	}

	fmt.Println("\n--- Creating S3 Bucket ---")
	createdBucket, err := storageService.CreateS3Bucket(ctx, bucket)
	if err != nil {
		fmt.Printf("CreateS3Bucket error: %v\n", err)
		return
	}

	fmt.Printf("Bucket ID: %s\n", createdBucket.ID)
	if createdBucket.ARN != nil {
		fmt.Printf("Bucket ARN: %s\n", *createdBucket.ARN)
	}
	if createdBucket.BucketDomainName != nil {
		fmt.Printf("Domain: %s\n", *createdBucket.BucketDomainName)
	}
	if createdBucket.BucketRegionalDomainName != nil {
		fmt.Printf("Regional Domain: %s\n", *createdBucket.BucketRegionalDomainName)
	}

	// 2. Configure ACL (canned private)
	fmt.Println("\n--- Updating Bucket ACL (private) ---")
	canned := "private"
	acl := &domainstorage.S3BucketACL{
		Bucket: createdBucket.ID,
		ACL:    &canned,
	}
	if err := acl.Validate(); err != nil {
		fmt.Printf("ACL validation error: %v\n", err)
		return
	}

	updatedACL, err := storageService.UpdateS3BucketACL(ctx, createdBucket.ID, acl)
	if err != nil {
		fmt.Printf("UpdateS3BucketACL error: %v\n", err)
		return
	}
	fmt.Printf("ACL updated for bucket %s: %v\n", updatedACL.Bucket, awsmapper.ToDomainS3BucketACL(awsmapper.FromDomainS3BucketACL(updatedACL)))

	// 3. Enable versioning
	fmt.Println("\n--- Enabling Bucket Versioning ---")
	versioning := &domainstorage.S3BucketVersioning{
		Bucket: createdBucket.ID,
		Status: "Enabled",
	}
	if err := versioning.Validate(); err != nil {
		fmt.Printf("Versioning validation error: %v\n", err)
		return
	}

	updatedVersioning, err := storageService.UpdateS3BucketVersioning(ctx, createdBucket.ID, versioning)
	if err != nil {
		fmt.Printf("UpdateS3BucketVersioning error: %v\n", err)
		return
	}
	fmt.Printf("Versioning status for bucket %s: %s\n", updatedVersioning.Bucket, updatedVersioning.Status)

	// 4. Configure encryption (SSE-S3)
	fmt.Println("\n--- Configuring Bucket Encryption (AES256) ---")
	encryption := &domainstorage.S3BucketEncryption{
		Bucket: createdBucket.ID,
		Rule: domainstorage.S3BucketEncryptionRule{
			BucketKeyEnabled: false,
			DefaultEncryption: domainstorage.S3BucketDefaultEncryption{
				SSEAlgorithm: "AES256",
			},
		},
	}
	if err := encryption.Validate(); err != nil {
		fmt.Printf("Encryption validation error: %v\n", err)
		return
	}

	updatedEnc, err := storageService.UpdateS3BucketEncryption(ctx, createdBucket.ID, encryption)
	if err != nil {
		fmt.Printf("UpdateS3BucketEncryption error: %v\n", err)
		return
	}
	fmt.Printf("Encryption for bucket %s: algorithm=%s, bucket_key=%v\n",
		updatedEnc.Bucket,
		updatedEnc.Rule.DefaultEncryption.SSEAlgorithm,
		updatedEnc.Rule.BucketKeyEnabled,
	)

	// 5. List buckets
	fmt.Println("\n--- Listing Buckets ---")
	buckets, err := storageService.ListS3Buckets(ctx, map[string]string{})
	if err != nil {
		fmt.Printf("ListS3Buckets error: %v\n", err)
		return
	}
	for _, b := range buckets {
		fmt.Printf("- %s (region=%s)\n", b.ID, b.Region)
	}
}
