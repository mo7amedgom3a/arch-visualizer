package storage

import (
	"testing"

	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
	awss3 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3"
	awss3outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3/outputs"
)

func TestFromDomainS3Bucket(t *testing.T) {
	tests := []struct {
		name      string
		domain    *domainstorage.S3Bucket
		wantNil   bool
		checkFunc func(*testing.T, *awss3.Bucket)
	}{
		{
			name:    "nil input",
			domain:  nil,
			wantNil: true,
		},
		{
			name: "with bucket name",
			domain: &domainstorage.S3Bucket{
				Name:         "test-bucket",
				Region:       "us-east-1",
				ForceDestroy: false,
				Tags: map[string]string{
					"Environment": "test",
				},
			},
			wantNil: false,
			checkFunc: func(t *testing.T, aws *awss3.Bucket) {
				if aws.Bucket == nil || *aws.Bucket != "test-bucket" {
					t.Errorf("Expected bucket name 'test-bucket', got %v", aws.Bucket)
				}
				if aws.BucketPrefix != nil {
					t.Error("Expected bucket prefix to be nil")
				}
				if aws.ForceDestroy != false {
					t.Error("Expected ForceDestroy to be false")
				}
				if len(aws.Tags) != 2 { // Environment tag + Name tag
					t.Errorf("Expected 2 tags, got %d", len(aws.Tags))
				}
			},
		},
		{
			name: "with bucket prefix",
			domain: &domainstorage.S3Bucket{
				NamePrefix:   stringPtr("my-app-logs-"),
				Region:       "us-east-1",
				ForceDestroy: true,
			},
			wantNil: false,
			checkFunc: func(t *testing.T, aws *awss3.Bucket) {
				if aws.BucketPrefix == nil || *aws.BucketPrefix != "my-app-logs-" {
					t.Errorf("Expected bucket prefix 'my-app-logs-', got %v", aws.BucketPrefix)
				}
				if aws.Bucket != nil {
					t.Error("Expected bucket name to be nil")
				}
				if aws.ForceDestroy != true {
					t.Error("Expected ForceDestroy to be true")
				}
			},
		},
		{
			name: "with tags",
			domain: &domainstorage.S3Bucket{
				Name:   "test-bucket",
				Region: "us-east-1",
				Tags: map[string]string{
					"Environment": "prod",
					"Project":     "test",
				},
			},
			wantNil: false,
			checkFunc: func(t *testing.T, aws *awss3.Bucket) {
				if len(aws.Tags) != 3 { // Environment + Project + Name tags
					t.Errorf("Expected 3 tags, got %d", len(aws.Tags))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromDomainS3Bucket(tt.domain)
			if tt.wantNil {
				if got != nil {
					t.Errorf("FromDomainS3Bucket() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("FromDomainS3Bucket() = nil, want non-nil")
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, got)
			}
		})
	}
}

func TestToDomainS3Bucket(t *testing.T) {
	tests := []struct {
		name      string
		aws       *awss3.Bucket
		wantNil   bool
		checkFunc func(*testing.T, *domainstorage.S3Bucket)
	}{
		{
			name:    "nil input",
			aws:     nil,
			wantNil: true,
		},
		{
			name: "with bucket name",
			aws: &awss3.Bucket{
				Bucket:       stringPtr("test-bucket"),
				ForceDestroy: false,
			},
			wantNil: false,
			checkFunc: func(t *testing.T, domain *domainstorage.S3Bucket) {
				if domain.Name != "test-bucket" {
					t.Errorf("Expected bucket name 'test-bucket', got %s", domain.Name)
				}
				if domain.NamePrefix != nil {
					t.Error("Expected name prefix to be nil")
				}
			},
		},
		{
			name: "with bucket prefix",
			aws: &awss3.Bucket{
				BucketPrefix: stringPtr("my-app-logs-"),
				ForceDestroy: true,
			},
			wantNil: false,
			checkFunc: func(t *testing.T, domain *domainstorage.S3Bucket) {
				if domain.NamePrefix == nil || *domain.NamePrefix != "my-app-logs-" {
					t.Errorf("Expected bucket prefix 'my-app-logs-', got %v", domain.NamePrefix)
				}
				if domain.Name != "" {
					t.Error("Expected bucket name to be empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToDomainS3Bucket(tt.aws)
			if tt.wantNil {
				if got != nil {
					t.Errorf("ToDomainS3Bucket() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ToDomainS3Bucket() = nil, want non-nil")
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, got)
			}
		})
	}
}

func TestToDomainS3BucketFromOutput(t *testing.T) {
	tests := []struct {
		name      string
		output    *awss3outputs.BucketOutput
		wantNil   bool
		checkFunc func(*testing.T, *domainstorage.S3Bucket)
	}{
		{
			name:    "nil input",
			output:  nil,
			wantNil: true,
		},
		{
			name: "with all fields",
			output: &awss3outputs.BucketOutput{
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
					{Key: "Environment", Value: "test"},
				},
			},
			wantNil: false,
			checkFunc: func(t *testing.T, domain *domainstorage.S3Bucket) {
				if domain.ID != "test-bucket" {
					t.Errorf("Expected ID 'test-bucket', got %s", domain.ID)
				}
				if domain.ARN == nil || *domain.ARN != "arn:aws:s3:::test-bucket" {
					t.Errorf("Expected ARN 'arn:aws:s3:::test-bucket', got %v", domain.ARN)
				}
				if domain.BucketDomainName == nil || *domain.BucketDomainName != "test-bucket.s3.amazonaws.com" {
					t.Errorf("Expected domain name 'test-bucket.s3.amazonaws.com', got %v", domain.BucketDomainName)
				}
				if domain.BucketRegionalDomainName == nil || *domain.BucketRegionalDomainName != "test-bucket.s3.us-east-1.amazonaws.com" {
					t.Errorf("Expected regional domain name 'test-bucket.s3.us-east-1.amazonaws.com', got %v", domain.BucketRegionalDomainName)
				}
				if domain.Region != "us-east-1" {
					t.Errorf("Expected region 'us-east-1', got %s", domain.Region)
				}
				if domain.Tags["Environment"] != "test" {
					t.Errorf("Expected tag Environment=test, got %v", domain.Tags)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToDomainS3BucketFromOutput(tt.output)
			if tt.wantNil {
				if got != nil {
					t.Errorf("ToDomainS3BucketFromOutput() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ToDomainS3BucketFromOutput() = nil, want non-nil")
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, got)
			}
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
