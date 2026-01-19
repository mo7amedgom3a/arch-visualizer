package iam

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Policy represents a cloud-agnostic IAM policy
// This is the domain model - no cloud-specific details
type Policy struct {
	ID            string
	ARN           *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	Name          string
	Description   *string
	Path          *string // Default is "/" if not specified
	PolicyDocument string // JSON string containing the policy document
	Tags          []PolicyTag
	Type          PolicyType // Type of policy (AWS managed, customer managed, or inline)
	IsAttachable  *bool      // Whether the policy can be attached (for AWS managed policies)
}

// PolicyTag represents a key-value tag for IAM resources
type PolicyTag struct {
	Key   string
	Value string
}

// Validate performs domain-level validation
func (p *Policy) Validate() error {
	if p.Name == "" {
		return errors.New("policy name is required")
	}

	// Validate name format: 1-128 chars, alphanumeric + +=,.@-_
	if len(p.Name) < 1 || len(p.Name) > 128 {
		return errors.New("policy name must be between 1 and 128 characters")
	}

	namePattern := regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`)
	if !namePattern.MatchString(p.Name) {
		return errors.New("policy name contains invalid characters. Allowed: alphanumeric, +=,.@-_")
	}

	// Validate policy document is valid JSON
	if p.PolicyDocument == "" {
		return errors.New("policy document is required")
	}

	var jsonDoc interface{}
	if err := json.Unmarshal([]byte(p.PolicyDocument), &jsonDoc); err != nil {
		return fmt.Errorf("policy document must be valid JSON: %w", err)
	}

	// Validate path if provided
	if p.Path != nil && *p.Path != "" {
		path := *p.Path
		if !strings.HasPrefix(path, "/") {
			return errors.New("policy path must start with '/'")
		}
		if len(path) > 512 {
			return errors.New("policy path cannot exceed 512 characters")
		}
	}

	return nil
}
