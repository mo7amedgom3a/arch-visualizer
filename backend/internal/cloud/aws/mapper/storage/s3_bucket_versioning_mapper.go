package storage

import (
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
	awss3 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3"
	awss3outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3/outputs"
)

// FromDomainS3BucketVersioning converts domain versioning to AWS model
func FromDomainS3BucketVersioning(domain *domainstorage.S3BucketVersioning) *awss3.BucketVersioning {
	if domain == nil {
		return nil
	}

	return &awss3.BucketVersioning{
		Bucket:    domain.Bucket,
		Status:    domain.Status,
		MFADelete: domain.MFADelete,
	}
}

// ToDomainS3BucketVersioning converts AWS model to domain
func ToDomainS3BucketVersioning(awsVersioning *awss3.BucketVersioning) *domainstorage.S3BucketVersioning {
	if awsVersioning == nil {
		return nil
	}

	return &domainstorage.S3BucketVersioning{
		Bucket:    awsVersioning.Bucket,
		Status:    awsVersioning.Status,
		MFADelete: awsVersioning.MFADelete,
	}
}

// ToDomainS3BucketVersioningFromOutput converts AWS output to domain
func ToDomainS3BucketVersioningFromOutput(output *awss3outputs.BucketVersioningOutput) *domainstorage.S3BucketVersioning {
	if output == nil {
		return nil
	}

	return &domainstorage.S3BucketVersioning{
		Bucket:    output.ID,
		Status:    output.Status,
		MFADelete: output.MFADelete,
	}
}
