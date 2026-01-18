package launch_template

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2"
)

// LaunchTemplate represents an AWS EC2 Launch Template configuration
// Note: Uses image_id (not ami) and name_prefix (recommended over name)
type LaunchTemplate struct {
	// Naming
	NamePrefix *string `json:"name_prefix,omitempty"` // Recommended: creates unique name
	Name       *string `json:"name,omitempty"`         // Alternative to name_prefix

	// Compute Configuration
	ImageID      string `json:"image_id"`       // AMI ID (note: image_id not ami)
	InstanceType string `json:"instance_type"` // e.g., "t3.micro", "m5.large"

	// Networking & Security
	VpcSecurityGroupIds []string `json:"vpc_security_group_ids"` // Required: at least one

	// Access & Permissions
	KeyName            *string            `json:"key_name,omitempty"`
	IAMInstanceProfile *IAMInstanceProfile `json:"iam_instance_profile,omitempty"` // Structured block

	// Storage
	RootVolumeID         *string  `json:"root_volume_id,omitempty"`         // Reference to storage volume for root device
	AdditionalVolumeIDs []string `json:"additional_volume_ids,omitempty"` // References to additional storage volumes

	// Configuration
	UserData        *string          `json:"user_data,omitempty"`         // Base64 encoded (max 16KB)
	MetadataOptions *MetadataOptions `json:"metadata_options,omitempty"`  // IMDSv2 settings

	// Version Management
	UpdateDefaultVersion *bool `json:"update_default_version,omitempty"` // Default: true

	// Tags
	Tags []configs.Tag `json:"tags,omitempty"` // Tags for the template itself
}

// Validate performs AWS-specific validation
func (lt *LaunchTemplate) Validate() error {
	// Name or NamePrefix must be provided
	if (lt.Name == nil || *lt.Name == "") && (lt.NamePrefix == nil || *lt.NamePrefix == "") {
		return errors.New("launch template name or name_prefix is required")
	}

	// ImageID validation
	if lt.ImageID == "" {
		return errors.New("image_id is required")
	}
	if !strings.HasPrefix(lt.ImageID, "ami-") {
		return errors.New("image_id must start with 'ami-'")
	}
	if len(lt.ImageID) < 12 || len(lt.ImageID) > 21 {
		return errors.New("image_id format is invalid")
	}

	// Instance type validation - use helper from ec2 package
	if lt.InstanceType == "" {
		return errors.New("instance type is required")
	}
	if !ec2.IsValidInstanceType(lt.InstanceType) {
		return fmt.Errorf("invalid instance type format: %s", lt.InstanceType)
	}

	// Security group validation
	if len(lt.VpcSecurityGroupIds) == 0 {
		return errors.New("at least one security group is required")
	}
	for _, sgID := range lt.VpcSecurityGroupIds {
		if !strings.HasPrefix(sgID, "sg-") {
			return fmt.Errorf("security group id must start with 'sg-': %s", sgID)
		}
	}

	// UserData validation (16KB base64 encoded limit)
	if lt.UserData != nil {
		if len(*lt.UserData) > 16384 {
			return errors.New("user data cannot exceed 16KB when base64 encoded")
		}
	}

	// Root volume ID validation
	if lt.RootVolumeID != nil && *lt.RootVolumeID != "" {
		if !strings.HasPrefix(*lt.RootVolumeID, "vol-") {
			return errors.New("root volume id must start with 'vol-'")
		}
	}

	// Additional volume IDs validation
	for i, volID := range lt.AdditionalVolumeIDs {
		if volID == "" {
			return fmt.Errorf("additional volume id at index %d cannot be empty", i)
		}
		if !strings.HasPrefix(volID, "vol-") {
			return fmt.Errorf("additional volume id at index %d must start with 'vol-'", i)
		}
	}

	// IAM instance profile validation
	if lt.IAMInstanceProfile != nil {
		if err := lt.IAMInstanceProfile.Validate(); err != nil {
			return fmt.Errorf("IAM instance profile validation failed: %w", err)
		}
	}

	// Metadata options validation
	if lt.MetadataOptions != nil {
		if err := lt.MetadataOptions.Validate(); err != nil {
			return fmt.Errorf("metadata options validation failed: %w", err)
		}
	}

	return nil
}
