package iam

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIAMMetadata_GetResourceSchema_Role(t *testing.T) {
	svc := NewIAMMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "iam_role")

	require.NoError(t, err)
	assert.Equal(t, "iam_role", schema.Label)

	requiredNames := []string{}
	for _, f := range schema.Fields {
		if f.Required {
			requiredNames = append(requiredNames, f.Name)
		}
	}
	assert.Contains(t, requiredNames, "name")
	assert.Contains(t, requiredNames, "assume_role_policy")

	assert.Contains(t, schema.Outputs, "arn")
	assert.Contains(t, schema.Outputs, "unique_id")
}

func TestIAMMetadata_GetResourceSchema_User(t *testing.T) {
	svc := NewIAMMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "iam_user")

	require.NoError(t, err)
	assert.Equal(t, "iam_user", schema.Label)

	requiredNames := []string{}
	for _, f := range schema.Fields {
		if f.Required {
			requiredNames = append(requiredNames, f.Name)
		}
	}
	assert.Contains(t, requiredNames, "name")
	assert.Contains(t, schema.Outputs, "arn")
	assert.Contains(t, schema.Outputs, "unique_id")
}

func TestIAMMetadata_GetResourceSchema_Policy(t *testing.T) {
	svc := NewIAMMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "iam_policy")

	require.NoError(t, err)
	assert.Equal(t, "iam_policy", schema.Label)

	requiredNames := []string{}
	for _, f := range schema.Fields {
		if f.Required {
			requiredNames = append(requiredNames, f.Name)
		}
	}
	assert.Contains(t, requiredNames, "name")
	assert.Contains(t, requiredNames, "policy_document")

	assert.Contains(t, schema.Outputs, "arn")
	assert.Contains(t, schema.Outputs, "attachment_count")
}

func TestIAMMetadata_GetResourceSchema_Group(t *testing.T) {
	svc := NewIAMMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "iam_group")

	require.NoError(t, err)
	assert.Equal(t, "iam_group", schema.Label)

	requiredNames := []string{}
	for _, f := range schema.Fields {
		if f.Required {
			requiredNames = append(requiredNames, f.Name)
		}
	}
	assert.Contains(t, requiredNames, "name")
	assert.Contains(t, schema.Outputs, "arn")
	assert.Contains(t, schema.Outputs, "unique_id")
}

func TestIAMMetadata_GetResourceSchema_InstanceProfile(t *testing.T) {
	svc := NewIAMMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "iam_instance_profile")

	require.NoError(t, err)
	assert.Equal(t, "iam_instance_profile", schema.Label)

	assert.Contains(t, schema.Outputs, "arn")
	assert.Contains(t, schema.Outputs, "roles")
}

func TestIAMMetadata_ListResourceSchemas_AllRegistered(t *testing.T) {
	svc := NewIAMMetadataService()
	schemas, err := svc.ListResourceSchemas(context.Background())

	require.NoError(t, err)
	assert.Len(t, schemas, 5, "expect all 5 IAM resource schemas")

	labels := map[string]bool{}
	for _, s := range schemas {
		labels[s.Label] = true
	}

	expectedLabels := []string{
		"iam_role", "iam_user", "iam_policy", "iam_group", "iam_instance_profile",
	}
	for _, l := range expectedLabels {
		assert.True(t, labels[l], "missing schema for %s", l)
	}
}

func TestIAMMetadata_GetResourceSchema_NotFound(t *testing.T) {
	svc := NewIAMMetadataService()
	_, err := svc.GetResourceSchema(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown IAM resource")
}
