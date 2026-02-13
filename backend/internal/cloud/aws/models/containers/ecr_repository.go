package containers

import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"

// ECRRepository represents an AWS ECR Repository
type ECRRepository struct {
	Name               string        `json:"name"`
	ImageTagMutability string        `json:"image_tag_mutability,omitempty"` // MUTABLE or IMMUTABLE
	ScanOnPush         bool          `json:"scan_on_push,omitempty"`
	EncryptionType     string        `json:"encryption_type,omitempty"` // AES256 or KMS
	KMSKey             string        `json:"kms_key,omitempty"`
	ForceDelete        bool          `json:"force_delete,omitempty"`
	Tags               []configs.Tag `json:"tags,omitempty"`
}
