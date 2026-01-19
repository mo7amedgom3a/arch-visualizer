package iam

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// Policy represents an AWS IAM policy configuration
type Policy struct {
	Name          string        `json:"name"`
	Description   *string       `json:"description,omitempty"`
	Path          *string       `json:"path,omitempty"` // Default is "/"
	PolicyDocument string       `json:"policy_document"` // JSON string
	Tags          []configs.Tag `json:"tags,omitempty"`
}

// Validate performs AWS-specific validation
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
