package iam

import (
	"errors"
	"regexp"
	"strings"
)

// User represents a cloud-agnostic IAM user
// This is the domain model - no cloud-specific details
type User struct {
	ID                string
	ARN               *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	Name              string
	Path              *string // Default is "/" if not specified
	PermissionsBoundary *string // ARN of permissions boundary policy
	ForceDestroy      *bool   // If true, allows deletion even with attached resources
	Tags              []PolicyTag
}

// Validate performs domain-level validation
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
