package iam

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// User represents an AWS IAM user configuration
type User struct {
	Name                string        `json:"name"`
	Path                *string       `json:"path,omitempty"`                 // Default is "/"
	PermissionsBoundary *string       `json:"permissions_boundary,omitempty"` // ARN of permissions boundary
	ForceDestroy        *bool         `json:"force_destroy,omitempty"`
	Tags                []configs.Tag `json:"tags,omitempty"`
	IsVirtual           bool          `json:"is_virtual,omitempty"` // If true, this resource exists only for simulation/terraform generation
}

// Validate performs AWS-specific validation
func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("user name is required")
	}

	// Validate name format: 1-64 chars, alphanumeric + +=,.@-_
	if len(u.Name) < 1 || len(u.Name) > 64 {
		return errors.New("user name must be between 1 and 64 characters")
	}

	namePattern := regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`)
	if !namePattern.MatchString(u.Name) {
		return errors.New("user name contains invalid characters. Allowed: alphanumeric, +=,.@-_")
	}

	// Validate path if provided
	if u.Path != nil && *u.Path != "" {
		path := *u.Path
		if !strings.HasPrefix(path, "/") {
			return errors.New("user path must start with '/'")
		}
		if len(path) > 512 {
			return errors.New("user path cannot exceed 512 characters")
		}
	}

	// Validate permissions boundary ARN format if provided
	if u.PermissionsBoundary != nil && *u.PermissionsBoundary != "" {
		arn := *u.PermissionsBoundary
		if !strings.HasPrefix(arn, "arn:aws:iam::") {
			return errors.New("permissions boundary must be a valid IAM policy ARN")
		}
	}

	return nil
}
