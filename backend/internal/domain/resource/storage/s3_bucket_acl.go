package storage

import (
	"errors"
	"fmt"
)

// S3BucketACL models legacy bucket ACL configuration
// Prefer IAM policies where possible; this is provided for compatibility.
type S3BucketACL struct {
	Bucket               string
	ACL                  *string
	AccessControlPolicy  *S3AccessControlPolicy
}

// S3AccessControlPolicy represents a detailed ACL policy
type S3AccessControlPolicy struct {
	Owner  S3Owner
	Grants []S3Grant
}

// S3Owner represents the ACL owner
type S3Owner struct {
	ID          string
	DisplayName *string
}

// S3Grant represents a single ACL grant
type S3Grant struct {
	Grantee    S3Grantee
	Permission string
}

// S3Grantee represents the grantee of a grant
type S3Grantee struct {
	Type        string // CanonicalUser | AmazonCustomerByEmail | Group
	ID          *string
	URI         *string
	EmailAddress *string
	DisplayName *string
}

// Validate ensures the ACL configuration is valid
func (a *S3BucketACL) Validate() error {
	if a.Bucket == "" {
		return errors.New("bucket is required")
	}

	if (a.ACL == nil || *a.ACL == "") && a.AccessControlPolicy == nil {
		return errors.New("either acl or access_control_policy is required")
	}

	if a.ACL != nil && *a.ACL != "" && a.AccessControlPolicy != nil {
		return errors.New("cannot specify both acl and access_control_policy")
	}

	if a.ACL != nil && *a.ACL != "" {
		if err := validateCannedACL(*a.ACL); err != nil {
			return err
		}
	}

	if a.AccessControlPolicy != nil {
		if err := a.AccessControlPolicy.Validate(); err != nil {
			return fmt.Errorf("access_control_policy invalid: %w", err)
		}
	}

	return nil
}

// Validate validates the access control policy
func (p *S3AccessControlPolicy) Validate() error {
	if err := p.Owner.Validate(); err != nil {
		return fmt.Errorf("owner: %w", err)
	}

	if len(p.Grants) == 0 {
		return errors.New("at least one grant is required")
	}

	for i, g := range p.Grants {
		if err := g.Validate(); err != nil {
			return fmt.Errorf("grant %d: %w", i, err)
		}
	}

	return nil
}

// Validate validates the owner
func (o *S3Owner) Validate() error {
	if o.ID == "" {
		return errors.New("owner id is required")
	}
	return nil
}

// Validate validates a grant
func (g *S3Grant) Validate() error {
	if err := g.Grantee.Validate(); err != nil {
		return fmt.Errorf("grantee: %w", err)
	}

	validPermissions := map[string]bool{
		"FULL_CONTROL": true,
		"READ":         true,
		"WRITE":        true,
		"READ_ACP":     true,
		"WRITE_ACP":    true,
	}

	if !validPermissions[g.Permission] {
		return fmt.Errorf("invalid permission: %s", g.Permission)
	}

	return nil
}

// Validate validates the grantee
func (g *S3Grantee) Validate() error {
	validTypes := map[string]bool{
		"CanonicalUser":        true,
		"AmazonCustomerByEmail": true,
		"Group":                true,
	}

	if !validTypes[g.Type] {
		return fmt.Errorf("invalid grantee type: %s", g.Type)
	}

	switch g.Type {
	case "CanonicalUser":
		if g.ID == nil || *g.ID == "" {
			return errors.New("canonical user requires id")
		}
	case "AmazonCustomerByEmail":
		if g.EmailAddress == nil || *g.EmailAddress == "" {
			return errors.New("amazon customer by email requires email_address")
		}
	case "Group":
		if g.URI == nil || *g.URI == "" {
			return errors.New("group grantee requires uri")
		}
	}

	return nil
}

// validateCannedACL ensures canned ACL value is acceptable
func validateCannedACL(acl string) error {
	validACLs := map[string]bool{
		"private":                  true,
		"public-read":              true,
		"public-read-write":        true,
		"authenticated-read":       true,
		"log-delivery-write":       true,
		"bucket-owner-read":        true,
		"bucket-owner-full-control": true,
		"aws-exec-read":            true,
	}

	if !validACLs[acl] {
		return fmt.Errorf("invalid acl: %s", acl)
	}
	return nil
}
