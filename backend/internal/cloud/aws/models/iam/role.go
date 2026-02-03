package iam

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// InlinePolicy represents an inline policy attached directly to a role
type InlinePolicy struct {
	Name   string `json:"name"`
	Policy string `json:"policy"` // JSON string
}

// Role represents an AWS IAM role configuration
type Role struct {
	Name                string         `json:"name"`
	Description         *string        `json:"description,omitempty"`
	Path                *string        `json:"path,omitempty"`     // Default is "/"
	AssumeRolePolicy    string         `json:"assume_role_policy"` // JSON string (trust policy)
	ManagedPolicyARNs   []string       `json:"managed_policy_arns,omitempty"`
	InlinePolicies      []InlinePolicy `json:"inline_policies,omitempty"`
	PermissionsBoundary *string        `json:"permissions_boundary,omitempty"` // ARN of permissions boundary
	ForceDetachPolicies *bool          `json:"force_detach_policies,omitempty"`
	Tags                []configs.Tag  `json:"tags,omitempty"`
	IsVirtual           bool           `json:"is_virtual,omitempty"` // If true, this resource exists only for simulation/terraform generation
}

// Validate performs AWS-specific validation
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

	// Validate managed policy ARNs
	for i, arn := range r.ManagedPolicyARNs {
		if !strings.HasPrefix(arn, "arn:aws:iam::") {
			return fmt.Errorf("managed policy ARN at index %d must be a valid IAM policy ARN", i)
		}
	}

	// Validate permissions boundary ARN format if provided
	if r.PermissionsBoundary != nil && *r.PermissionsBoundary != "" {
		arn := *r.PermissionsBoundary
		if !strings.HasPrefix(arn, "arn:aws:iam::") {
			return errors.New("permissions boundary must be a valid IAM policy ARN")
		}
	}

	// Validate inline policies
	for i, inlinePolicy := range r.InlinePolicies {
		if inlinePolicy.Name == "" {
			return fmt.Errorf("inline policy at index %d must have a name", i)
		}
		var jsonDoc interface{}
		if err := json.Unmarshal([]byte(inlinePolicy.Policy), &jsonDoc); err != nil {
			return fmt.Errorf("inline policy %s at index %d must be valid JSON: %w", inlinePolicy.Name, i, err)
		}
	}

	return nil
}
