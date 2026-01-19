package iam

import (
	"errors"
	"regexp"
	"strings"
)

// Group represents a cloud-agnostic IAM group
// This is the domain model - no cloud-specific details
type Group struct {
	ID   string
	ARN  *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	Name string
	Path *string // Default is "/" if not specified
	Tags []PolicyTag
}

// Validate performs domain-level validation
func (g *Group) Validate() error {
	if g.Name == "" {
		return errors.New("group name is required")
	}

	// Validate name format: 1-128 chars, alphanumeric + +=,.@-_
	if len(g.Name) < 1 || len(g.Name) > 128 {
		return errors.New("group name must be between 1 and 128 characters")
	}

	namePattern := regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`)
	if !namePattern.MatchString(g.Name) {
		return errors.New("group name contains invalid characters. Allowed: alphanumeric, +=,.@-_")
	}

	// Validate path if provided
	if g.Path != nil && *g.Path != "" {
		path := *g.Path
		if !strings.HasPrefix(path, "/") {
			return errors.New("group path must start with '/'")
		}
		if len(path) > 512 {
			return errors.New("group path cannot exceed 512 characters")
		}
	}

	return nil
}
