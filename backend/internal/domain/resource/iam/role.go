package iam

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Role represents a cloud-agnostic IAM role
// This is the domain model - no cloud-specific details
type Role struct {
	ID                string
	ARN               *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	Name              string
	Description       *string
	Path              *string // Default is "/" if not specified
	AssumeRolePolicy  string  // JSON string containing the trust policy
	PermissionsBoundary *string // ARN of permissions boundary policy
	Tags              []PolicyTag
}

// Validate performs domain-level validation
func (r *Role) Validate() error {
	if r.Name == "" {
		return errors.New("role name is required")
	}

	// Validate name format: 1-64 chars, alphanumeric + +=,.@-_
	if len(r.Name) < 1 || len(r.Name) > 64 {
		return errors.New("role name must be between 1 and 64 characters")
	}

	namePattern := regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`)
	if !namePattern.MatchString(r.Name) {
		return errors.New("role name contains invalid characters. Allowed: alphanumeric, +=,.@-_")
	}

	// Validate assume role policy is valid JSON
	if r.AssumeRolePolicy == "" {
		return errors.New("assume role policy is required")
	}

	var jsonDoc interface{}
	if err := json.Unmarshal([]byte(r.AssumeRolePolicy), &jsonDoc); err != nil {
		return fmt.Errorf("assume role policy must be valid JSON: %w", err)
	}

	// Validate path if provided
	if r.Path != nil && *r.Path != "" {
		path := *r.Path
		if !strings.HasPrefix(path, "/") {
			return errors.New("role path must start with '/'")
		}
		if len(path) > 512 {
			return errors.New("role path cannot exceed 512 characters")
		}
	}

	// Validate permissions boundary ARN format if provided
	if r.PermissionsBoundary != nil && *r.PermissionsBoundary != "" {
		arn := *r.PermissionsBoundary
		if !strings.HasPrefix(arn, "arn:aws:iam::") {
			return errors.New("permissions boundary must be a valid IAM policy ARN")
		}
	}

	return nil
}
