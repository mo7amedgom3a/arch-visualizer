package iam

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// InstanceProfile represents an AWS IAM instance profile configuration
type InstanceProfile struct {
	Name      string        `json:"name,omitempty"`       // Optional if NamePrefix is provided
	NamePrefix *string      `json:"name_prefix,omitempty"` // Optional, conflicts with Name
	Path      *string       `json:"path,omitempty"`     // Default is "/"
	Role      *string       `json:"role,omitempty"`     // IAM Role name to attach
	Tags      []configs.Tag `json:"tags,omitempty"`
}

// Validate performs AWS-specific validation
func (ip *InstanceProfile) Validate() error {
	// Name or NamePrefix must be provided
	if ip.Name == "" && (ip.NamePrefix == nil || *ip.NamePrefix == "") {
		return errors.New("either name or name_prefix must be provided")
	}

	// Name and NamePrefix cannot both be provided
	if ip.Name != "" && ip.NamePrefix != nil && *ip.NamePrefix != "" {
		return errors.New("name and name_prefix cannot both be provided")
	}

	// Validate name if provided
	if ip.Name != "" {
		if len(ip.Name) < 1 || len(ip.Name) > 128 {
			return errors.New("instance profile name must be between 1 and 128 characters")
		}

		namePattern := regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`)
		if !namePattern.MatchString(ip.Name) {
			return errors.New("instance profile name contains invalid characters. Allowed: alphanumeric, +=,.@-_")
		}
	}

	// Validate namePrefix if provided
	if ip.NamePrefix != nil && *ip.NamePrefix != "" {
		if len(*ip.NamePrefix) < 1 || len(*ip.NamePrefix) > 38 {
			return errors.New("instance profile name_prefix must be between 1 and 38 characters")
		}

		namePattern := regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`)
		if !namePattern.MatchString(*ip.NamePrefix) {
			return errors.New("instance profile name_prefix contains invalid characters. Allowed: alphanumeric, +=,.@-_")
		}
	}

	// Validate path if provided
	if ip.Path != nil && *ip.Path != "" {
		path := *ip.Path
		if !strings.HasPrefix(path, "/") {
			return errors.New("instance profile path must start with '/'")
		}
		if len(path) > 512 {
			return errors.New("instance profile path cannot exceed 512 characters")
		}
	}

	// Validate role name if provided
	if ip.Role != nil && *ip.Role != "" {
		if len(*ip.Role) < 1 || len(*ip.Role) > 64 {
			return errors.New("role name must be between 1 and 64 characters")
		}

		rolePattern := regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`)
		if !rolePattern.MatchString(*ip.Role) {
			return errors.New("role name contains invalid characters. Allowed: alphanumeric, +=,.@-_")
		}
	}

	return nil
}

// GetName returns the name or generates one from prefix
func (ip *InstanceProfile) GetName() string {
	if ip.Name != "" {
		return ip.Name
	}
	if ip.NamePrefix != nil && *ip.NamePrefix != "" {
		// AWS will generate a unique name with the prefix
		return *ip.NamePrefix
	}
	return ""
}
