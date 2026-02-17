package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseMetadata_GetResourceSchema_RDSInstance(t *testing.T) {
	svc := NewDatabaseMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "rds_instance")

	require.NoError(t, err)
	assert.Equal(t, "rds_instance", schema.Label)

	requiredNames := []string{}
	for _, f := range schema.Fields {
		if f.Required {
			requiredNames = append(requiredNames, f.Name)
		}
	}
	assert.Contains(t, requiredNames, "name")
	assert.Contains(t, requiredNames, "engine")
	assert.Contains(t, requiredNames, "engine_version")
	assert.Contains(t, requiredNames, "instance_class")
	assert.Contains(t, requiredNames, "allocated_storage")

	assert.Contains(t, schema.Outputs, "id")
	assert.Contains(t, schema.Outputs, "arn")
	assert.Contains(t, schema.Outputs, "endpoint")
}

func TestDatabaseMetadata_ListResourceSchemas_AllRegistered(t *testing.T) {
	svc := NewDatabaseMetadataService()
	schemas, err := svc.ListResourceSchemas(context.Background())

	require.NoError(t, err)
	assert.Len(t, schemas, 1, "expect 1 database resource schema")

	labels := map[string]bool{}
	for _, s := range schemas {
		labels[s.Label] = true
	}

	assert.True(t, labels["rds_instance"], "missing schema for rds_instance")
}

func TestDatabaseMetadata_GetResourceSchema_NotFound(t *testing.T) {
	svc := NewDatabaseMetadataService()
	_, err := svc.GetResourceSchema(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown database resource")
}

func TestDatabaseMetadata_RDSInstance_EngineEnum(t *testing.T) {
	svc := NewDatabaseMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "rds_instance")
	require.NoError(t, err)

	for _, f := range schema.Fields {
		if f.Name == "engine" {
			assert.Contains(t, f.Enum, "mysql")
			assert.Contains(t, f.Enum, "postgres")
			assert.Contains(t, f.Enum, "mariadb")
			return
		}
	}
	t.Fatal("engine field not found")
}
