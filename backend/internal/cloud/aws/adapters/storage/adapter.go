package storage

import (
	"context"
	"fmt"

	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/storage"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/storage"
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
)

// AWSStorageAdapter adapts AWS-specific storage service to domain storage service
// This implements the Adapter pattern, allowing the domain layer to work with cloud-specific implementations
type AWSStorageAdapter struct {
	awsService awsservice.AWSStorageService
}

// NewAWSStorageAdapter creates a new AWS storage adapter
func NewAWSStorageAdapter(awsService awsservice.AWSStorageService) domainstorage.StorageService {
	return &AWSStorageAdapter{
		awsService: awsService,
	}
}

// Ensure AWSStorageAdapter implements StorageService
var _ domainstorage.StorageService = (*AWSStorageAdapter)(nil)

// EBS Volume Operations

func (a *AWSStorageAdapter) CreateEBSVolume(ctx context.Context, volume *domainstorage.EBSVolume) (*domainstorage.EBSVolume, error) {
	if err := volume.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsVolume := awsmapper.FromDomainEBSVolume(volume)
	if err := awsVolume.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsVolumeOutput, err := a.awsService.CreateEBSVolume(ctx, awsVolume)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainVolume := awsmapper.ToDomainEBSVolumeFromOutput(awsVolumeOutput)
	// Preserve region from input
	domainVolume.Region = volume.Region
	return domainVolume, nil
}

func (a *AWSStorageAdapter) GetEBSVolume(ctx context.Context, id string) (*domainstorage.EBSVolume, error) {
	awsVolumeOutput, err := a.awsService.GetEBSVolume(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainEBSVolumeFromOutput(awsVolumeOutput), nil
}

func (a *AWSStorageAdapter) UpdateEBSVolume(ctx context.Context, id string, volume *domainstorage.EBSVolume) (*domainstorage.EBSVolume, error) {
	if err := volume.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	awsVolume := awsmapper.FromDomainEBSVolume(volume)
	if err := awsVolume.Validate(); err != nil {
		return nil, fmt.Errorf("aws validation failed: %w", err)
	}

	awsVolumeOutput, err := a.awsService.UpdateEBSVolume(ctx, id, awsVolume)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainVolume := awsmapper.ToDomainEBSVolumeFromOutput(awsVolumeOutput)
	// Preserve region from input
	domainVolume.Region = volume.Region
	return domainVolume, nil
}

func (a *AWSStorageAdapter) DeleteEBSVolume(ctx context.Context, id string) error {
	if err := a.awsService.DeleteEBSVolume(ctx, id); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSStorageAdapter) ListEBSVolumes(ctx context.Context, filters map[string]string) ([]*domainstorage.EBSVolume, error) {
	// Convert domain filters to AWS filters format
	awsFilters := make(map[string][]string)
	for key, value := range filters {
		awsFilters[key] = []string{value}
	}

	awsVolumeOutputs, err := a.awsService.ListEBSVolumes(ctx, awsFilters)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	domainVolumes := make([]*domainstorage.EBSVolume, len(awsVolumeOutputs))
	for i, awsVolumeOutput := range awsVolumeOutputs {
		domainVolumes[i] = awsmapper.ToDomainEBSVolumeFromOutput(awsVolumeOutput)
	}

	return domainVolumes, nil
}

// Volume Attachment Operations

func (a *AWSStorageAdapter) AttachVolume(ctx context.Context, volumeID, instanceID, deviceName string) error {
	if err := a.awsService.AttachVolume(ctx, volumeID, instanceID, deviceName); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}

func (a *AWSStorageAdapter) DetachVolume(ctx context.Context, volumeID, instanceID string) error {
	if err := a.awsService.DetachVolume(ctx, volumeID, instanceID); err != nil {
		return fmt.Errorf("aws service error: %w", err)
	}
	return nil
}
