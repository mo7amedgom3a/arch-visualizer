package compute

import (
	"errors"
	"fmt"
	"strings"
)

// LaunchTemplate represents a cloud-agnostic launch template (blueprint for instances)
// This is the domain model - no cloud-specific details
type LaunchTemplate struct {
	ID       string
	ARN      *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	Name     string
	Region   string
	NamePrefix *string // Optional prefix for unique naming

	// Compute Configuration
	ImageID      string // AMI ID or image identifier
	InstanceType string // e.g., "t3.micro", "m5.large"

	// Networking
	SecurityGroupIDs []string

	// Access & Permissions
	KeyName            *string
	IAMInstanceProfile *string

	// Storage
	// RootVolumeID references a storage volume resource (EBS volume) for the root device
	// If nil, a default root volume will be created by the cloud provider
	RootVolumeID *string

	// AdditionalVolumeIDs references additional storage volumes to attach
	// These volumes must be created separately and attached to instances launched from this template
	AdditionalVolumeIDs []string

	// Configuration
	UserData        *string // Base64 encoded user data script
	MetadataOptions *MetadataOptions

	// Version Management
	Version       *int // Current/default version number
	LatestVersion *int // Latest version number
}

// MetadataOptions represents instance metadata service options (IMDSv2)
type MetadataOptions struct {
	HTTPEndpoint            *string // "enabled" or "disabled"
	HTTPTokens              *string // "required" (IMDSv2) or "optional" (IMDSv1)
	HTTPPutResponseHopLimit *int    // Hop limit for PUT requests (1-64)
}

// Validate performs domain-level validation
func (lt *LaunchTemplate) Validate() error {
	// Name or NamePrefix must be provided
	if lt.Name == "" && (lt.NamePrefix == nil || *lt.NamePrefix == "") {
		return errors.New("launch template name or name_prefix is required")
	}

	if lt.Region == "" {
		return errors.New("launch template region is required")
	}

	if lt.ImageID == "" {
		return errors.New("image_id is required")
	}

	if lt.InstanceType == "" {
		return errors.New("instance type is required")
	}

	if len(lt.SecurityGroupIDs) == 0 {
		return errors.New("at least one security group is required")
	}

	// Validate root volume ID format if provided
	if lt.RootVolumeID != nil && *lt.RootVolumeID != "" {
		if !strings.HasPrefix(*lt.RootVolumeID, "vol-") {
			return errors.New("root volume id must start with 'vol-'")
		}
	}

	// Validate additional volume IDs
	for i, volID := range lt.AdditionalVolumeIDs {
		if volID == "" {
			return fmt.Errorf("additional volume id at index %d cannot be empty", i)
		}
		if !strings.HasPrefix(volID, "vol-") {
			return fmt.Errorf("additional volume id at index %d must start with 'vol-'", i)
		}
	}

	// Validate metadata options if provided
	if lt.MetadataOptions != nil {
		if err := lt.validateMetadataOptions(); err != nil {
			return err
		}
	}

	// Validate user data size (16KB base64 encoded limit)
	if lt.UserData != nil {
		if len(*lt.UserData) > 16384 {
			return errors.New("user data cannot exceed 16KB when base64 encoded")
		}
	}

	return nil
}

// validateMetadataOptions validates metadata options
func (lt *LaunchTemplate) validateMetadataOptions() error {
	mo := lt.MetadataOptions

	if mo.HTTPEndpoint != nil {
		validEndpoints := map[string]bool{
			"enabled":  true,
			"disabled": true,
		}
		if !validEndpoints[*mo.HTTPEndpoint] {
			return fmt.Errorf("invalid http_endpoint: %s (must be 'enabled' or 'disabled')", *mo.HTTPEndpoint)
		}
	}

	if mo.HTTPTokens != nil {
		validTokens := map[string]bool{
			"required": true,
			"optional": true,
		}
		if !validTokens[*mo.HTTPTokens] {
			return fmt.Errorf("invalid http_tokens: %s (must be 'required' or 'optional')", *mo.HTTPTokens)
		}
	}

	if mo.HTTPPutResponseHopLimit != nil {
		if *mo.HTTPPutResponseHopLimit < 1 || *mo.HTTPPutResponseHopLimit > 64 {
			return errors.New("http_put_response_hop_limit must be between 1 and 64")
		}
	}

	return nil
}
