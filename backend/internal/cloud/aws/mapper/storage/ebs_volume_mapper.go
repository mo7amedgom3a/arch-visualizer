package storage

import (
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
	awsebs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs"
	awsebsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// FromDomainEBSVolume converts domain EBSVolume to AWS Volume input model
func FromDomainEBSVolume(domain *domainstorage.EBSVolume) *awsebs.Volume {
	if domain == nil {
		return nil
	}

	awsVolume := &awsebs.Volume{
		Name:             domain.Name,
		AvailabilityZone: domain.AvailabilityZone,
		Size:             domain.Size,
		VolumeType:       domain.Type,
		Encrypted:        domain.Encrypted,
		IOPS:             domain.IOPS,
		Throughput:       domain.Throughput,
		KMSKeyID:         domain.KMSKeyID,
		SnapshotID:       domain.SnapshotID,
	}

	// Add Name tag
	awsVolume.Tags = []configs.Tag{
		{Key: "Name", Value: domain.Name},
	}

	return awsVolume
}

// ToDomainEBSVolume converts AWS Volume input model to domain EBSVolume
// This is useful for backward compatibility or when reading existing volumes
func ToDomainEBSVolume(aws *awsebs.Volume) *domainstorage.EBSVolume {
	if aws == nil {
		return nil
	}

	domain := &domainstorage.EBSVolume{
		Name:             aws.Name,
		AvailabilityZone: aws.AvailabilityZone,
		Size:             aws.Size,
		Type:             aws.VolumeType,
		Encrypted:        aws.Encrypted,
		IOPS:             aws.IOPS,
		Throughput:       aws.Throughput,
		KMSKeyID:         aws.KMSKeyID,
		SnapshotID:       aws.SnapshotID,
	}

	return domain
}

// ToDomainEBSVolumeFromOutput converts AWS VolumeOutput to domain EBSVolume
// This populates the domain model with AWS-generated identifiers (ID, ARN, state)
func ToDomainEBSVolumeFromOutput(output *awsebsoutputs.VolumeOutput) *domainstorage.EBSVolume {
	if output == nil {
		return nil
	}

	domain := &domainstorage.EBSVolume{
		ID:               output.ID,
		ARN:              &output.ARN,
		Name:             output.Name,
		Region:           "", // Region should be set from context
		AvailabilityZone: output.AvailabilityZone,
		Size:             output.Size,
		Type:             output.VolumeType,
		IOPS:             output.IOPS,
		Throughput:       output.Throughput,
		Encrypted:        output.Encrypted,
		KMSKeyID:         output.KMSKeyID,
		SnapshotID:       output.SnapshotID,
		State:            domainstorage.EBSVolumeState(output.State),
		AttachedTo:       output.AttachedTo,
	}

	return domain
}
