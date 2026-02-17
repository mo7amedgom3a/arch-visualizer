package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageMetadata_GetResourceSchema_EBSVolume(t *testing.T) {
	svc := NewStorageMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "ebs_volume")

	require.NoError(t, err)
	assert.Equal(t, "ebs_volume", schema.Label)

	requiredNames := []string{}
	for _, f := range schema.Fields {
		if f.Required {
			requiredNames = append(requiredNames, f.Name)
		}
	}
	assert.Contains(t, requiredNames, "name")
	assert.Contains(t, requiredNames, "availability_zone")
	assert.Contains(t, requiredNames, "size")
	assert.Contains(t, requiredNames, "volume_type")

	assert.Contains(t, schema.Outputs, "id")
	assert.Contains(t, schema.Outputs, "arn")
}

func TestStorageMetadata_GetResourceSchema_S3Bucket(t *testing.T) {
	svc := NewStorageMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "s3_bucket")

	require.NoError(t, err)
	assert.Equal(t, "s3_bucket", schema.Label)

	assert.Contains(t, schema.Outputs, "id")
	assert.Contains(t, schema.Outputs, "arn")
	assert.Contains(t, schema.Outputs, "bucket_domain_name")
}

func TestStorageMetadata_ListResourceSchemas_AllRegistered(t *testing.T) {
	svc := NewStorageMetadataService()
	schemas, err := svc.ListResourceSchemas(context.Background())

	require.NoError(t, err)
	assert.Len(t, schemas, 2, "expect all 2 storage resource schemas")

	labels := map[string]bool{}
	for _, s := range schemas {
		labels[s.Label] = true
	}

	expectedLabels := []string{"ebs_volume", "s3_bucket"}
	for _, l := range expectedLabels {
		assert.True(t, labels[l], "missing schema for %s", l)
	}
}

func TestStorageMetadata_GetResourceSchema_NotFound(t *testing.T) {
	svc := NewStorageMetadataService()
	_, err := svc.GetResourceSchema(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown storage resource")
}

func TestStorageMetadata_EBSVolume_VolumeTypeEnum(t *testing.T) {
	svc := NewStorageMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "ebs_volume")
	require.NoError(t, err)

	for _, f := range schema.Fields {
		if f.Name == "volume_type" {
			assert.ElementsMatch(t, []string{"gp2", "gp3", "io1", "io2", "sc1", "st1", "standard"}, f.Enum)
			return
		}
	}
	t.Fatal("volume_type field not found")
}
