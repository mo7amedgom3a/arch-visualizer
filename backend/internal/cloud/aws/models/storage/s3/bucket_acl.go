package s3

import (
	"errors"
	"fmt"
)

// BucketACL represents AWS S3 bucket ACL configuration
type BucketACL struct {
	Bucket              string               `json:"bucket"`
	ACL                 *string              `json:"acl,omitempty"`
	AccessControlPolicy *AccessControlPolicy `json:"access_control_policy,omitempty"`
}

// AccessControlPolicy represents an ACL policy
type AccessControlPolicy struct {
	Owner  Owner   `json:"owner"`
	Grants []Grant `json:"grants"`
}

// Owner represents the bucket owner
type Owner struct {
	ID          string  `json:"id"`
	DisplayName *string `json:"display_name,omitempty"`
}

// Grant represents a single grant
type Grant struct {
	Grantee    Grantee `json:"grantee"`
	Permission string  `json:"permission"`
}

// Grantee represents the grantee of a grant
type Grantee struct {
	Type         string  `json:"type"` // CanonicalUser | AmazonCustomerByEmail | Group
	ID           *string `json:"id,omitempty"`
	URI          *string `json:"uri,omitempty"`
	EmailAddress *string `json:"email_address,omitempty"`
	DisplayName  *string `json:"display_name,omitempty"`
}

// Validate validates ACL configuration
func (b *BucketACL) Validate() error {
	if b.Bucket == "" {
		return errors.New("bucket is required")
	}

	if (b.ACL == nil || *b.ACL == "") && b.AccessControlPolicy == nil {
		return errors.New("either acl or access_control_policy is required")
	}

	if b.ACL != nil && *b.ACL != "" && b.AccessControlPolicy != nil {
		return errors.New("cannot specify both acl and access_control_policy")
	}

	if b.ACL != nil && *b.ACL != "" {
		if err := validateCannedACL(*b.ACL); err != nil {
			return err
		}
	}

	if b.AccessControlPolicy != nil {
		if err := b.AccessControlPolicy.Validate(); err != nil {
			return fmt.Errorf("access_control_policy invalid: %w", err)
		}
	}

	return nil
}

// Validate validates access control policy
func (p *AccessControlPolicy) Validate() error {
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

// Validate validates owner
func (o *Owner) Validate() error {
	if o.ID == "" {
		return errors.New("owner id is required")
	}
	return nil
}

// Validate validates grant
func (g *Grant) Validate() error {
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

// Validate validates grantee
func (g *Grantee) Validate() error {
	validTypes := map[string]bool{
		"CanonicalUser":         true,
		"AmazonCustomerByEmail": true,
		"Group":                 true,
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

// validateCannedACL validates canned ACL value
func validateCannedACL(acl string) error {
	validACLs := map[string]bool{
		"private":                   true,
		"public-read":               true,
		"public-read-write":         true,
		"authenticated-read":        true,
		"log-delivery-write":        true,
		"bucket-owner-read":         true,
		"bucket-owner-full-control": true,
		"aws-exec-read":             true,
	}

	if !validACLs[acl] {
		return fmt.Errorf("invalid acl: %s", acl)
	}
	return nil
}
