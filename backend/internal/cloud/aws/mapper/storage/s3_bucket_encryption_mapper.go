package storage

import (
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
	awss3 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3"
	awss3outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3/outputs"
)

// FromDomainS3BucketEncryption converts domain encryption to AWS model
func FromDomainS3BucketEncryption(domain *domainstorage.S3BucketEncryption) *awss3.BucketEncryption {
	if domain == nil {
		return nil
	}

	aws := &awss3.BucketEncryption{
		Bucket: domain.Bucket,
		Rule: awss3.BucketEncryptionRule{
			BucketKeyEnabled: domain.Rule.BucketKeyEnabled,
			DefaultEncryption: awss3.BucketDefaultEncryption{
				SSEAlgorithm:   domain.Rule.DefaultEncryption.SSEAlgorithm,
				KMSMasterKeyID: domain.Rule.DefaultEncryption.KMSMasterKeyID,
			},
		},
	}

	return aws
}

// ToDomainS3BucketEncryption converts AWS model to domain
func ToDomainS3BucketEncryption(awsEncryption *awss3.BucketEncryption) *domainstorage.S3BucketEncryption {
	if awsEncryption == nil {
		return nil
	}

	return &domainstorage.S3BucketEncryption{
		Bucket: awsEncryption.Bucket,
		Rule: domainstorage.S3BucketEncryptionRule{
			BucketKeyEnabled: awsEncryption.Rule.BucketKeyEnabled,
			DefaultEncryption: domainstorage.S3BucketDefaultEncryption{
				SSEAlgorithm:   awsEncryption.Rule.DefaultEncryption.SSEAlgorithm,
				KMSMasterKeyID: awsEncryption.Rule.DefaultEncryption.KMSMasterKeyID,
			},
		},
	}
}

// ToDomainS3BucketEncryptionFromOutput converts AWS output to domain
func ToDomainS3BucketEncryptionFromOutput(output *awss3outputs.BucketEncryptionOutput) *domainstorage.S3BucketEncryption {
	if output == nil {
		return nil
	}

	return &domainstorage.S3BucketEncryption{
		Bucket: output.ID,
		Rule: domainstorage.S3BucketEncryptionRule{
			BucketKeyEnabled: output.BucketKeyEnabled,
			DefaultEncryption: domainstorage.S3BucketDefaultEncryption{
				SSEAlgorithm:   output.SSEAlgorithm,
				KMSMasterKeyID: output.KMSMasterKeyID,
			},
		},
	}
}
