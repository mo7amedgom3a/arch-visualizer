package ebs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// Volume represents an AWS EBS volume configuration
type Volume struct {
	Name             string        `json:"name"`               // Required
	AvailabilityZone string        `json:"availability_zone"`  // Required
	Size             int           `json:"size"`               // Size in GiB (required)
	VolumeType       string        `json:"volume_type"`        // gp3, gp2, io1, io2, sc1, st1, standard
	IOPS             *int          `json:"iops,omitempty"`      // Optional for gp3/io1/io2
	Throughput       *int          `json:"throughput,omitempty"` // Optional for gp3 (MB/s)
	Encrypted        bool          `json:"encrypted"`           // Default: false
	KMSKeyID         *string       `json:"kms_key_id,omitempty"` // Optional KMS key ARN
	SnapshotID       *string       `json:"snapshot_id,omitempty"` // Optional snapshot ID
	Tags             []configs.Tag `json:"tags,omitempty"`      // Optional tags
}

// Validate performs AWS-specific validation
func (v *Volume) Validate() error {
	if v.Name == "" {
		return errors.New("volume name is required")
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

	if !validVolumeTypes[v.VolumeType] {
		return fmt.Errorf("invalid volume type: %s", v.VolumeType)
	}

	// IOPS validation
	if v.IOPS != nil {
		if v.VolumeType == "gp3" {
			if *v.IOPS < 3000 || *v.IOPS > 16000 {
				return errors.New("gp3 IOPS must be between 3000 and 16000")
			}
		} else if v.VolumeType == "io1" || v.VolumeType == "io2" {
			if *v.IOPS < 100 || *v.IOPS > 64000 {
				return errors.New("io1/io2 IOPS must be between 100 and 64000")
			}
		} else {
			return fmt.Errorf("IOPS can only be specified for gp3, io1, or io2 volume types")
		}
	}

	// Throughput validation (gp3 only)
	if v.Throughput != nil {
		if v.VolumeType != "gp3" {
			return errors.New("throughput can only be specified for gp3 volume type")
		}
		if *v.Throughput < 125 || *v.Throughput > 1000 {
			return errors.New("gp3 throughput must be between 125 and 1000 MB/s")
		}
	}

	// KMS Key ID format validation if provided
	if v.KMSKeyID != nil && *v.KMSKeyID != "" {
		if !strings.HasPrefix(*v.KMSKeyID, "arn:aws:kms:") {
			return errors.New("kms key id must be a valid KMS key ARN")
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
	parts := strings.Split(az, "-")
	if len(parts) < 3 {
		return false
	}

	lastPart := parts[len(parts)-1]
	if len(lastPart) < 2 {
		return false
	}

	lastChar := lastPart[len(lastPart)-1]
	if lastChar < 'a' || lastChar > 'z' {
		return false
	}

	return true
}
