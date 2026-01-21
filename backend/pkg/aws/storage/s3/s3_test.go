package s3

import (
	"context"
	"testing"

	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
)

func TestNewS3DemoAdapter_CreateBucket(t *testing.T) {
	adapter := newS3DemoAdapter()
	ctx := context.Background()

	bucket := &domainstorage.S3Bucket{
		Name:         "test-bucket-12345",
		Region:       "us-east-1",
		ForceDestroy: true,
	}

	created, err := adapter.CreateS3Bucket(ctx, bucket)
	if err != nil {
		t.Fatalf("CreateS3Bucket returned error: %v", err)
	}
	if created == nil {
		t.Fatalf("CreateS3Bucket returned nil bucket")
	}
	if created.ID == "" {
		t.Errorf("expected bucket ID to be set")
	}
	if created.ARN == nil || *created.ARN == "" {
		t.Errorf("expected bucket ARN to be set")
	}
}

func TestS3Demo_ACL_Versioning_Encryption_Flows(t *testing.T) {
	adapter := newS3DemoAdapter()
	ctx := context.Background()

	// Create bucket
	bucket := &domainstorage.S3Bucket{
		Name:         "test-bucket-xyz",
		Region:       "us-east-1",
		ForceDestroy: true,
	}
	created, err := adapter.CreateS3Bucket(ctx, bucket)
	if err != nil {
		t.Fatalf("CreateS3Bucket error: %v", err)
	}

	// ACL
	canned := "private"
	acl := &domainstorage.S3BucketACL{
		Bucket: created.ID,
		ACL:    &canned,
	}
	if err := acl.Validate(); err != nil {
		t.Fatalf("ACL validation error: %v", err)
	}
	updatedACL, err := adapter.UpdateS3BucketACL(ctx, created.ID, acl)
	if err != nil {
		t.Fatalf("UpdateS3BucketACL error: %v", err)
	}
	if updatedACL == nil || updatedACL.ACL == nil || *updatedACL.ACL != "private" {
		t.Fatalf("expected ACL private, got %#v", updatedACL)
	}

	// Versioning
	versioning := &domainstorage.S3BucketVersioning{
		Bucket: created.ID,
		Status: "Enabled",
	}
	if err := versioning.Validate(); err != nil {
		t.Fatalf("Versioning validation error: %v", err)
	}
	updatedVer, err := adapter.UpdateS3BucketVersioning(ctx, created.ID, versioning)
	if err != nil {
		t.Fatalf("UpdateS3BucketVersioning error: %v", err)
	}
	if updatedVer.Status != "Enabled" {
		t.Fatalf("expected versioning Enabled, got %s", updatedVer.Status)
	}

	// Encryption
	encryption := &domainstorage.S3BucketEncryption{
		Bucket: created.ID,
		Rule: domainstorage.S3BucketEncryptionRule{
			BucketKeyEnabled: false,
			DefaultEncryption: domainstorage.S3BucketDefaultEncryption{
				SSEAlgorithm: "AES256",
			},
		},
	}
	if err := encryption.Validate(); err != nil {
		t.Fatalf("Encryption validation error: %v", err)
	}
	updatedEnc, err := adapter.UpdateS3BucketEncryption(ctx, created.ID, encryption)
	if err != nil {
		t.Fatalf("UpdateS3BucketEncryption error: %v", err)
	}
	if updatedEnc.Rule.DefaultEncryption.SSEAlgorithm != "AES256" {
		t.Fatalf("expected AES256 algorithm, got %s", updatedEnc.Rule.DefaultEncryption.SSEAlgorithm)
	}
}

