package storage

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services"
)

// ---------------------------------------------------------------------------
// Internal registry
// ---------------------------------------------------------------------------

var storageSchemaRegistry = map[string]*services.ResourceSchema{}

func registerStorageSchema(label string, schema *services.ResourceSchema) {
	storageSchemaRegistry[label] = schema
}

// ---------------------------------------------------------------------------
// Service interface & implementation
// ---------------------------------------------------------------------------

// StorageMetadataService exposes structured schemas for AWS storage resources.
type StorageMetadataService interface {
	GetResourceSchema(ctx context.Context, resource string) (*services.ResourceSchema, error)
	ListResourceSchemas(ctx context.Context) ([]*services.ResourceSchema, error)
}

type storageMetadataServiceImpl struct{}

// NewStorageMetadataService returns a ready-to-use metadata service.
func NewStorageMetadataService() StorageMetadataService {
	return &storageMetadataServiceImpl{}
}

func (s *storageMetadataServiceImpl) GetResourceSchema(_ context.Context, resource string) (*services.ResourceSchema, error) {
	schema, ok := storageSchemaRegistry[resource]
	if !ok {
		return nil, fmt.Errorf("unknown storage resource: %s", resource)
	}
	return schema, nil
}

func (s *storageMetadataServiceImpl) ListResourceSchemas(_ context.Context) ([]*services.ResourceSchema, error) {
	out := make([]*services.ResourceSchema, 0, len(storageSchemaRegistry))
	for _, schema := range storageSchemaRegistry {
		out = append(out, schema)
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Schema registrations (run once at import time)
// ---------------------------------------------------------------------------

func init() {
	// ---- EBS Volume ----
	registerStorageSchema("ebs_volume", &services.ResourceSchema{
		Label: "ebs_volume",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "Name of the EBS volume"},
			{Name: "availability_zone", Type: "string", Required: true, Enum: []string{}, Description: "Availability zone (e.g. us-east-1a)"},
			{Name: "size", Type: "int", Required: true, Enum: []string{}, Description: "Volume size in GiB (1-16384)"},
			{Name: "volume_type", Type: "string", Required: true, Enum: []string{"gp2", "gp3", "io1", "io2", "sc1", "st1", "standard"}, Description: "EBS volume type"},
			{Name: "iops", Type: "int", Required: false, Enum: []string{}, Description: "Provisioned IOPS (for gp3/io1/io2)"},
			{Name: "throughput", Type: "int", Required: false, Enum: []string{}, Description: "Throughput in MB/s (gp3 only, 125-1000)"},
			{Name: "encrypted", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Enable encryption"},
			{Name: "kms_key_id", Type: "string", Required: false, Enum: []string{}, Description: "KMS key ARN for encryption"},
			{Name: "snapshot_id", Type: "string", Required: false, Enum: []string{}, Description: "Snapshot ID to restore from"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":                "string",
			"arn":               "string",
			"state":             "string",
			"availability_zone": "string",
			"size":              "int",
			"volume_type":       "string",
		},
	})

	// ---- S3 Bucket ----
	registerStorageSchema("s3_bucket", &services.ResourceSchema{
		Label: "s3_bucket",
		Fields: []services.FieldDescriptor{
			{Name: "bucket", Type: "string", Required: false, Enum: []string{}, Description: "Bucket name (3-63 chars, globally unique)"},
			{Name: "bucket_prefix", Type: "string", Required: false, Enum: []string{}, Description: "Bucket name prefix (alternative to bucket)"},
			{Name: "force_destroy", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Allow deletion of non-empty bucket"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":                          "string",
			"arn":                         "string",
			"bucket_domain_name":          "string",
			"bucket_regional_domain_name": "string",
			"region":                      "string",
		},
	})
}
