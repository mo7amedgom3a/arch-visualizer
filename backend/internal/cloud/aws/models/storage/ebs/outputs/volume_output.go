package outputs

import (
	"time"
)

// VolumeOutput represents AWS EBS volume output/response data after creation
type VolumeOutput struct {
	// AWS-generated identifiers
	ID  string `json:"id"`  // e.g., "vol-0123456789abcdef0"
	ARN string `json:"arn"` // e.g., "arn:aws:ec2:us-east-1:123456789012:volume/vol-0123456789abcdef0"

	// Basic information
	Name             string `json:"name"`
	AvailabilityZone string `json:"availability_zone"`

	// Volume configuration
	Size       int    `json:"size"`        // Size in GiB
	VolumeType string `json:"volume_type"` // gp3, gp2, io1, io2, etc.

	// Performance
	IOPS      *int `json:"iops,omitempty"`       // IOPS for gp3/io1/io2
	Throughput *int `json:"throughput,omitempty"` // Throughput for gp3 (MB/s)

	// Security & Backups
	Encrypted  bool    `json:"encrypted"`
	KMSKeyID   *string `json:"kms_key_id,omitempty"`
	SnapshotID *string `json:"snapshot_id,omitempty"`

	// State
	State      string  `json:"state"`        // creating, available, in-use, deleting, deleted, error
	AttachedTo *string `json:"attached_to,omitempty"` // Instance ID if attached

	// Metadata
	CreateTime time.Time `json:"create_time"`

	// Tags (using struct to match AWS SDK Tag format)
	Tags []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"tags"`
}
