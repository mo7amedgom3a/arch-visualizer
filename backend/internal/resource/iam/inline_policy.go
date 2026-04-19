package iam

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

// InlinePolicy represents an inline policy attached directly to a user, role, or group
// Inline policies are embedded directly in the identity and cannot be reused
type InlinePolicy struct {
	Name   string // Policy name (unique within the identity)
	Policy string // JSON string containing the policy document
}

// Validate performs validation on an inline policy
func (ip *InlinePolicy) Validate() error {
	if ip.Name == "" {
		return errors.New("inline policy name is required")
	}

	// Validate name format: 1-128 chars, alphanumeric + +=,.@-_
	if len(ip.Name) < 1 || len(ip.Name) > 128 {
		return errors.New("inline policy name must be between 1 and 128 characters")
	}

	namePattern := regexp.MustCompile(`^[a-zA-Z0-9+=,.@_-]+$`)
	if !namePattern.MatchString(ip.Name) {
		return errors.New("inline policy name contains invalid characters. Allowed: alphanumeric, +=,.@-_")
	}

	// Validate policy document is valid JSON
	if ip.Policy == "" {
		return errors.New("inline policy document is required")
	}

	var jsonDoc interface{}
	if err := json.Unmarshal([]byte(ip.Policy), &jsonDoc); err != nil {
		return fmt.Errorf("inline policy document must be valid JSON: %w", err)
	}

	return nil
}
