package storage

import (
	"errors"
	"fmt"
	"strings"
)

// EBSVolumeState represents the state of an EBS volume
type EBSVolumeState string

const (
	EBSVolumeStateCreating EBSVolumeState = "creating"
	EBSVolumeStateAvailable EBSVolumeState = "available"
	EBSVolumeStateInUse    EBSVolumeState = "in-use"
	EBSVolumeStateDeleting EBSVolumeState = "deleting"
	EBSVolumeStateDeleted  EBSVolumeState = "deleted"
	EBSVolumeStateError    EBSVolumeState = "error"
)

// EBSVolume represents a cloud-agnostic EBS volume (block storage)
// This is the domain model - no cloud-specific details
type EBSVolume struct {
	ID               string
	ARN              *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	Name             string
	Region           string
	AvailabilityZone string // Required: must match instance AZ

	// Volume Configuration
	Size      int    // Size in GiB (required)
	Type      string // Volume type: gp3, gp2, io1, io2, sc1, st1, standard

	// Performance (Optional)
	IOPS      *int // Optional IOPS for gp3/io1/io2
	Throughput *int // Optional throughput for gp3 (MB/s)

	// Security & Backups
	Encrypted  bool
	KMSKeyID   *string // Optional KMS key ARN for encryption
	SnapshotID *string // Optional snapshot ID to create from

	// State
	State     EBSVolumeState
	AttachedTo *string // Instance ID if attached

	// Metadata
	CreateTime *string // Creation timestamp (optional)
}

// Validate performs domain-level validation
func (v *EBSVolume) Validate() error {
	if v.Name == "" {
		return errors.New("volume name is required")
	}

	if v.Region == "" {
		return errors.New("volume region is required")
	}

	if v.AvailabilityZone == "" {
		return errors.New("availability zone is required")
	}

	// Validate AZ format (e.g., us-east-1a)
	if !isValidAvailabilityZone(v.AvailabilityZone) {
		return errors.New("invalid availability zone format")
	}

	// Size validation
	if v.Size <= 0 {
		return errors.New("volume size must be greater than 0")
	}
	if v.Size > 16384 {
		return errors.New("volume size cannot exceed 16384 GiB")
	}

	// Volume type validation
	validVolumeTypes := map[string]bool{
		"gp2":      true,
		"gp3":      true,
		"io1":      true,
		"io2":      true,
		"sc1":      true,
		"st1":      true,
		"standard": true,
	}

	if !validVolumeTypes[v.Type] {
		return fmt.Errorf("invalid volume type: %s", v.Type)
	}

	// IOPS validation
	if v.IOPS != nil {
		if v.Type == "gp3" {
			if *v.IOPS < 3000 || *v.IOPS > 16000 {
				return errors.New("gp3 IOPS must be between 3000 and 16000")
			}
		} else if v.Type == "io1" || v.Type == "io2" {
			if *v.IOPS < 100 || *v.IOPS > 64000 {
				return errors.New("io1/io2 IOPS must be between 100 and 64000")
			}
		} else {
			return fmt.Errorf("IOPS can only be specified for gp3, io1, or io2 volume types")
		}
	}

	// Throughput validation (gp3 only)
	if v.Throughput != nil {
		if v.Type != "gp3" {
			return errors.New("throughput can only be specified for gp3 volume type")
		}
		if *v.Throughput < 125 || *v.Throughput > 1000 {
			return errors.New("gp3 throughput must be between 125 and 1000 MB/s")
		}
	}

	// Snapshot ID format validation if provided
	if v.SnapshotID != nil && *v.SnapshotID != "" {
		if !strings.HasPrefix(*v.SnapshotID, "snap-") {
			return errors.New("snapshot id must start with 'snap-'")
		}
	}

	return nil
}

// isValidAvailabilityZone validates AZ format (e.g., us-east-1a)
func isValidAvailabilityZone(az string) bool {
	// Basic format check: region-letter (e.g., us-east-1a, eu-west-2b)
	parts := strings.Split(az, "-")
	if len(parts) < 3 {
		return false
	}

	// Last part should be region + letter (e.g., "1a", "2b")
	lastPart := parts[len(parts)-1]
	if len(lastPart) < 2 {
		return false
	}

	// Should end with a letter
	lastChar := lastPart[len(lastPart)-1]
	if lastChar < 'a' || lastChar > 'z' {
		return false
	}

	return true
}

// Ensure EBSVolume implements StorageResource
var _ StorageResource = (*EBSVolume)(nil)

// GetID returns the volume ID
func (v *EBSVolume) GetID() string {
	return v.ID
}

// GetName returns the volume name
func (v *EBSVolume) GetName() string {
	return v.Name
}

// GetRegion returns the volume region
func (v *EBSVolume) GetRegion() string {
	return v.Region
}
