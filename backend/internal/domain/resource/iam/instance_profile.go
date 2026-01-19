package iam

import (
	"errors"
	"regexp"
	"strings"
)

// InstanceProfile represents a cloud-agnostic IAM instance profile
// Instance profiles are containers for IAM roles that can be attached to EC2 instances
type InstanceProfile struct {
	ID       string  // AWS-generated name (same as Name after creation)
	ARN      *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	Name     string
	Path     *string // Default is "/" if not specified
	RoleName *string // IAM Role name attached to this profile
	Tags     []PolicyTag
}

// Validate performs domain-level validation
func (ip *InstanceProfile) Validate() error {
	if ip.Name == "" {
		return errors.New("instance profile name is required")
	}

	// Validate name format: 1-128 chars, alphanumeric + +=,.@-_
	if len(ip.Name) < 1 || len(ip.Name) > 128 {
		return errors.New("instance profile name must be between 1 and 128 characters")
	}

	namePattern := regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`)
	if !namePattern.MatchString(ip.Name) {
		return errors.New("instance profile name contains invalid characters. Allowed: alphanumeric, +=,.@-_")
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

	return nil
}
