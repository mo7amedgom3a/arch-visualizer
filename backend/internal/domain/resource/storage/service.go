package storage

import (
	"context"
)

// StorageService defines the interface for storage resource operations
// This is cloud-agnostic and can be implemented by any cloud provider
type StorageService interface {
	// EBS Volume operations
	CreateEBSVolume(ctx context.Context, volume *EBSVolume) (*EBSVolume, error)
	GetEBSVolume(ctx context.Context, id string) (*EBSVolume, error)
	UpdateEBSVolume(ctx context.Context, id string, volume *EBSVolume) (*EBSVolume, error)
	DeleteEBSVolume(ctx context.Context, id string) error
	ListEBSVolumes(ctx context.Context, filters map[string]string) ([]*EBSVolume, error)

	// Volume attachment operations
	AttachVolume(ctx context.Context, volumeID, instanceID, deviceName string) error
	DetachVolume(ctx context.Context, volumeID, instanceID string) error
}
