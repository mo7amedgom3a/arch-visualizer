package storage

import (
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
	awss3 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3"
	awss3outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// FromDomainS3Bucket converts domain S3Bucket to AWS Bucket input model
func FromDomainS3Bucket(domain *domainstorage.S3Bucket) *awss3.Bucket {
	if domain == nil {
		return nil
	}

	awsBucket := &awss3.Bucket{
		ForceDestroy: domain.ForceDestroy,
	}

	// Set bucket name or prefix
	if domain.Name != "" {
		awsBucket.Bucket = &domain.Name
	} else if domain.NamePrefix != nil && *domain.NamePrefix != "" {
		awsBucket.BucketPrefix = domain.NamePrefix
	}

	// Convert tags from map[string]string to []configs.Tag
	if domain.Tags != nil && len(domain.Tags) > 0 {
		awsBucket.Tags = make([]configs.Tag, 0, len(domain.Tags))
		for key, value := range domain.Tags {
			awsBucket.Tags = append(awsBucket.Tags, configs.Tag{
				Key:   key,
				Value: value,
			})
		}
	}

	// Add Name tag if bucket name is set
	if domain.Name != "" {
		// Check if Name tag already exists
		hasNameTag := false
		for _, tag := range awsBucket.Tags {
			if tag.Key == "Name" {
				hasNameTag = true
				break
			}
		}
		if !hasNameTag {
			awsBucket.Tags = append(awsBucket.Tags, configs.Tag{
				Key:   "Name",
				Value: domain.Name,
			})
		}
	}

	return awsBucket
}

// ToDomainS3Bucket converts AWS Bucket input model to domain S3Bucket
// This is useful for backward compatibility or when reading existing buckets
func ToDomainS3Bucket(aws *awss3.Bucket) *domainstorage.S3Bucket {
	if aws == nil {
		return nil
	}

	domain := &domainstorage.S3Bucket{
		ForceDestroy: aws.ForceDestroy,
	}

	// Set bucket name or prefix
	if aws.Bucket != nil && *aws.Bucket != "" {
		domain.Name = *aws.Bucket
	} else if aws.BucketPrefix != nil && *aws.BucketPrefix != "" {
		domain.NamePrefix = aws.BucketPrefix
	}

	// Convert tags from []configs.Tag to map[string]string
	if aws.Tags != nil && len(aws.Tags) > 0 {
		domain.Tags = make(map[string]string)
		for _, tag := range aws.Tags {
			domain.Tags[tag.Key] = tag.Value
		}
	}

	return domain
}

// ToDomainS3BucketFromOutput converts AWS BucketOutput to domain S3Bucket
// This populates the domain model with AWS-generated identifiers (ID, ARN, domain names)
func ToDomainS3BucketFromOutput(output *awss3outputs.BucketOutput) *domainstorage.S3Bucket {
	if output == nil {
		return nil
	}

	domain := &domainstorage.S3Bucket{
		ID:                     output.ID,
		ARN:                    &output.ARN,
		Name:                   output.Name,
		NamePrefix:             output.NamePrefix,
		ForceDestroy:           output.ForceDestroy,
		Region:                 output.Region,
		BucketDomainName:       &output.BucketDomainName,
		BucketRegionalDomainName: &output.BucketRegionalDomainName,
	}

	// Convert tags from []struct to map[string]string
	if output.Tags != nil && len(output.Tags) > 0 {
		domain.Tags = make(map[string]string)
		for _, tag := range output.Tags {
			domain.Tags[tag.Key] = tag.Value
		}
	}

	return domain
}
