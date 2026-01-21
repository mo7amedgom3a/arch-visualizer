package storage

import (
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
	awss3 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3"
	awss3outputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/s3/outputs"
)

// FromDomainS3BucketACL converts domain ACL to AWS ACL model
func FromDomainS3BucketACL(domain *domainstorage.S3BucketACL) *awss3.BucketACL {
	if domain == nil {
		return nil
	}

	aws := &awss3.BucketACL{
		Bucket: domain.Bucket,
		ACL:    domain.ACL,
	}

	if domain.AccessControlPolicy != nil {
		aws.AccessControlPolicy = &awss3.AccessControlPolicy{
			Owner: awss3.Owner{
				ID:          domain.AccessControlPolicy.Owner.ID,
				DisplayName: domain.AccessControlPolicy.Owner.DisplayName,
			},
		}

		for _, g := range domain.AccessControlPolicy.Grants {
			awsGrant := awss3.Grant{
				Grantee: awss3.Grantee{
					Type:         g.Grantee.Type,
					ID:           g.Grantee.ID,
					URI:          g.Grantee.URI,
					EmailAddress: g.Grantee.EmailAddress,
					DisplayName:  g.Grantee.DisplayName,
				},
				Permission: g.Permission,
			}
			aws.AccessControlPolicy.Grants = append(aws.AccessControlPolicy.Grants, awsGrant)
		}
	}

	return aws
}

// ToDomainS3BucketACL converts AWS ACL model to domain
func ToDomainS3BucketACL(awsACL *awss3.BucketACL) *domainstorage.S3BucketACL {
	if awsACL == nil {
		return nil
	}

	domain := &domainstorage.S3BucketACL{
		Bucket: awsACL.Bucket,
		ACL:    awsACL.ACL,
	}

	if awsACL.AccessControlPolicy != nil {
		domain.AccessControlPolicy = &domainstorage.S3AccessControlPolicy{
			Owner: domainstorage.S3Owner{
				ID:          awsACL.AccessControlPolicy.Owner.ID,
				DisplayName: awsACL.AccessControlPolicy.Owner.DisplayName,
			},
		}
		for _, g := range awsACL.AccessControlPolicy.Grants {
			domainGrant := domainstorage.S3Grant{
				Grantee: domainstorage.S3Grantee{
					Type:         g.Grantee.Type,
					ID:           g.Grantee.ID,
					URI:          g.Grantee.URI,
					EmailAddress: g.Grantee.EmailAddress,
					DisplayName:  g.Grantee.DisplayName,
				},
				Permission: g.Permission,
			}
			domain.AccessControlPolicy.Grants = append(domain.AccessControlPolicy.Grants, domainGrant)
		}
	}

	return domain
}

// ToDomainS3BucketACLFromOutput converts AWS ACL output to domain
func ToDomainS3BucketACLFromOutput(output *awss3outputs.BucketACLOutput) *domainstorage.S3BucketACL {
	if output == nil {
		return nil
	}

	domain := &domainstorage.S3BucketACL{
		Bucket: output.ID,
		ACL:    output.ACL,
	}

	if output.AccessControlPolicy != nil {
		domain.AccessControlPolicy = &domainstorage.S3AccessControlPolicy{
			Owner: domainstorage.S3Owner{
				ID:          output.AccessControlPolicy.Owner.ID,
				DisplayName: output.AccessControlPolicy.Owner.DisplayName,
			},
		}

		for _, g := range output.AccessControlPolicy.Grants {
			domainGrant := domainstorage.S3Grant{
				Grantee: domainstorage.S3Grantee{
					Type:         g.Grantee.Type,
					ID:           g.Grantee.ID,
					URI:          g.Grantee.URI,
					EmailAddress: g.Grantee.EmailAddress,
					DisplayName:  g.Grantee.DisplayName,
				},
				Permission: g.Permission,
			}
			domain.AccessControlPolicy.Grants = append(domain.AccessControlPolicy.Grants, domainGrant)
		}
	}

	return domain
}
