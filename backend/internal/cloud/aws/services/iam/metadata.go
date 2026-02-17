package iam

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services"
)

// ---------------------------------------------------------------------------
// Internal registry
// ---------------------------------------------------------------------------

var iamSchemaRegistry = map[string]*services.ResourceSchema{}

func registerIAMSchema(label string, schema *services.ResourceSchema) {
	iamSchemaRegistry[label] = schema
}

// ---------------------------------------------------------------------------
// Service interface & implementation
// ---------------------------------------------------------------------------

// IAMMetadataService exposes structured schemas for AWS IAM resources.
type IAMMetadataService interface {
	GetResourceSchema(ctx context.Context, resource string) (*services.ResourceSchema, error)
	ListResourceSchemas(ctx context.Context) ([]*services.ResourceSchema, error)
}

type iamMetadataServiceImpl struct{}

// NewIAMMetadataService returns a ready-to-use metadata service.
func NewIAMMetadataService() IAMMetadataService {
	return &iamMetadataServiceImpl{}
}

func (s *iamMetadataServiceImpl) GetResourceSchema(_ context.Context, resource string) (*services.ResourceSchema, error) {
	schema, ok := iamSchemaRegistry[resource]
	if !ok {
		return nil, fmt.Errorf("unknown IAM resource: %s", resource)
	}
	return schema, nil
}

func (s *iamMetadataServiceImpl) ListResourceSchemas(_ context.Context) ([]*services.ResourceSchema, error) {
	out := make([]*services.ResourceSchema, 0, len(iamSchemaRegistry))
	for _, schema := range iamSchemaRegistry {
		out = append(out, schema)
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Schema registrations (run once at import time)
// ---------------------------------------------------------------------------

func init() {
	// ---- IAM Role ----
	registerIAMSchema("iam_role", &services.ResourceSchema{
		Label: "iam_role",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Role name (1-64 chars, alphanumeric + +=,.@-_)"},
			{Name: "description", Type: "string", Required: false, Enum: []string{}, Description: "Description of the role"},
			{Name: "path", Type: "string", Required: false, Enum: []string{}, Default: "/", Description: "IAM path (must start with /, max 512 chars)"},
			{Name: "assume_role_policy", Type: "string", Required: true, Enum: []string{}, Description: "Trust policy document (JSON string)"},
			{Name: "managed_policy_arns", Type: "[]string", Required: false, Enum: []string{}, Description: "ARNs of managed policies to attach"},
			{Name: "inline_policies", Type: "[]object", Required: false, Enum: []string{}, Description: "Inline policies (each with name and policy JSON)"},
			{Name: "permissions_boundary", Type: "string", Required: false, Enum: []string{}, Description: "ARN of permissions boundary policy"},
			{Name: "force_detach_policies", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Force detach policies on deletion"},
			{Name: "is_virtual", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Virtual resource for simulation only"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"arn":                  "string",
			"id":                   "string",
			"name":                 "string",
			"unique_id":            "string",
			"create_date":          "string",
			"max_session_duration": "int",
		},
	})

	// ---- IAM User ----
	registerIAMSchema("iam_user", &services.ResourceSchema{
		Label: "iam_user",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "User name (1-64 chars, alphanumeric + +=,.@-_)"},
			{Name: "path", Type: "string", Required: false, Enum: []string{}, Default: "/", Description: "IAM path (must start with /, max 512 chars)"},
			{Name: "permissions_boundary", Type: "string", Required: false, Enum: []string{}, Description: "ARN of permissions boundary policy"},
			{Name: "force_destroy", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Delete user even if it has non-Terraform-managed access keys, login profile, or MFA devices"},
			{Name: "is_virtual", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Virtual resource for simulation only"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"arn":                "string",
			"id":                 "string",
			"name":               "string",
			"unique_id":          "string",
			"create_date":        "string",
			"password_last_used": "string",
		},
	})

	// ---- IAM Policy ----
	registerIAMSchema("iam_policy", &services.ResourceSchema{
		Label: "iam_policy",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Policy name (1-128 chars, alphanumeric + +=,.@-_)"},
			{Name: "description", Type: "string", Required: false, Enum: []string{}, Description: "Description of the policy"},
			{Name: "path", Type: "string", Required: false, Enum: []string{}, Default: "/", Description: "IAM path (must start with /, max 512 chars)"},
			{Name: "policy_document", Type: "string", Required: true, Enum: []string{}, Description: "Policy document (JSON string)"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"arn":                "string",
			"id":                 "string",
			"name":               "string",
			"create_date":        "string",
			"update_date":        "string",
			"default_version_id": "string",
			"attachment_count":   "int",
			"is_attachable":      "bool",
		},
	})

	// ---- IAM Group ----
	registerIAMSchema("iam_group", &services.ResourceSchema{
		Label: "iam_group",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Group name (1-128 chars, alphanumeric + +=,.@-_)"},
			{Name: "path", Type: "string", Required: false, Enum: []string{}, Default: "/", Description: "IAM path (must start with /, max 512 chars)"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"arn":         "string",
			"id":          "string",
			"name":        "string",
			"unique_id":   "string",
			"create_date": "string",
		},
	})

	// ---- IAM Instance Profile ----
	registerIAMSchema("iam_instance_profile", &services.ResourceSchema{
		Label: "iam_instance_profile",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: false, Enum: []string{}, Description: "Instance profile name (1-128 chars, conflicts with name_prefix)"},
			{Name: "name_prefix", Type: "string", Required: false, Enum: []string{}, Description: "Name prefix for unique naming (1-38 chars, conflicts with name)"},
			{Name: "path", Type: "string", Required: false, Enum: []string{}, Default: "/", Description: "IAM path (must start with /, max 512 chars)"},
			{Name: "role", Type: "string", Required: false, Enum: []string{}, Description: "IAM role name to attach (1-64 chars)"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"arn":         "string",
			"id":          "string",
			"name":        "string",
			"create_date": "string",
			"roles":       "[]object",
		},
	})
}
