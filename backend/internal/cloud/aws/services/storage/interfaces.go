package storage

import (
	"context"

	awsebs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs"
	awsebsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/storage/ebs/outputs"
)

// AWSStorageService defines AWS-specific storage operations
// This implements cloud provider-specific logic while maintaining domain compatibility
type AWSStorageService interface {
	// EBS Volume operations
	CreateEBSVolume(ctx context.Context, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error)
	GetEBSVolume(ctx context.Context, id string) (*awsebsoutputs.VolumeOutput, error)
	UpdateEBSVolume(ctx context.Context, id string, volume *awsebs.Volume) (*awsebsoutputs.VolumeOutput, error)
	DeleteEBSVolume(ctx context.Context, id string) error
	ListEBSVolumes(ctx context.Context, filters map[string][]string) ([]*awsebsoutputs.VolumeOutput, error)

	// Volume attachment operations
	AttachVolume(ctx context.Context, volumeID, instanceID, deviceName string) error
	DetachVolume(ctx context.Context, volumeID, instanceID string) error
}
