package database

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services"
)

// ---------------------------------------------------------------------------
// Internal registry
// ---------------------------------------------------------------------------

var databaseSchemaRegistry = map[string]*services.ResourceSchema{}

func registerDatabaseSchema(label string, schema *services.ResourceSchema) {
	databaseSchemaRegistry[label] = schema
}

// ---------------------------------------------------------------------------
// Service interface & implementation
// ---------------------------------------------------------------------------

// DatabaseMetadataService exposes structured schemas for AWS database resources.
type DatabaseMetadataService interface {
	GetResourceSchema(ctx context.Context, resource string) (*services.ResourceSchema, error)
	ListResourceSchemas(ctx context.Context) ([]*services.ResourceSchema, error)
}

type databaseMetadataServiceImpl struct{}

// NewDatabaseMetadataService returns a ready-to-use metadata service.
func NewDatabaseMetadataService() DatabaseMetadataService {
	return &databaseMetadataServiceImpl{}
}

func (s *databaseMetadataServiceImpl) GetResourceSchema(_ context.Context, resource string) (*services.ResourceSchema, error) {
	schema, ok := databaseSchemaRegistry[resource]
	if !ok {
		return nil, fmt.Errorf("unknown database resource: %s", resource)
	}
	return schema, nil
}

func (s *databaseMetadataServiceImpl) ListResourceSchemas(_ context.Context) ([]*services.ResourceSchema, error) {
	out := make([]*services.ResourceSchema, 0, len(databaseSchemaRegistry))
	for _, schema := range databaseSchemaRegistry {
		out = append(out, schema)
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Schema registrations (run once at import time)
// ---------------------------------------------------------------------------

func init() {
	// ---- RDS Instance ----
	registerDatabaseSchema("rds_instance", &services.ResourceSchema{
		Label: "rds_instance",
		Fields: []services.FieldDescriptor{
			{Name: "name", Type: "string", Required: true, Enum: []string{}, Description: "DB instance identifier"},
			{Name: "engine", Type: "string", Required: true, Enum: []string{"mysql", "postgres", "mariadb", "oracle-ee", "oracle-se2", "sqlserver-ex", "sqlserver-web", "sqlserver-se", "sqlserver-ee"}, Description: "Database engine"},
			{Name: "engine_version", Type: "string", Required: true, Enum: []string{
				"5.7.38", "8.0.28", "8.0.30",
				"13.7", "14.1", "14.2", "15.0",
				"10.6.8", "10.5.16",
				"19.0.0.0.ru-2022-01.rur-2022-01.r1", "21.0.0.0.ru-2022-01.rur-2022-01.r1",
				"15.00.4043.16.v1", "14.00.3421.10.v1",
			}, Description: "Database engine version"},
			{Name: "instance_class", Type: "string", Required: true, Enum: []string{
				"db.t3.micro", "db.t3.small", "db.m5.large", "db.m5.xlarge", "db.r5.large",
			}, Description: "DB instance class (e.g. db.t3.micro, db.m5.large)"},
			{Name: "allocated_storage", Type: "int", Required: true, Enum: []string{}, Description: "Storage size in GiB"},
			{Name: "storage_type", Type: "string", Required: false, Enum: []string{"gp2", "gp3", "io1", "standard"}, Default: "gp2", Description: "Storage type"},
			{Name: "username", Type: "string", Required: false, Enum: []string{}, Description: "Master username"},
			{Name: "password", Type: "string", Required: false, Enum: []string{}, Description: "Master password"},
			{Name: "db_name", Type: "string", Required: false, Enum: []string{}, Description: "Initial database name"},
			{Name: "subnet_group_name", Type: "string", Required: false, Enum: []string{}, Description: "DB subnet group name"},
			{Name: "vpc_security_group_ids", Type: "[]string", Required: false, Enum: []string{}, Description: "VPC security group IDs"},
			{Name: "publicly_accessible", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Whether the instance is publicly accessible"},
			{Name: "multi_az", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Enable Multi-AZ deployment"},
			{Name: "backup_retention_period", Type: "int", Required: false, Enum: []string{}, Default: 0, Description: "Backup retention period in days (0-35)"},
			{Name: "skip_final_snapshot", Type: "bool", Required: false, Enum: []string{}, Default: false, Description: "Skip final snapshot on deletion"},
			{Name: "replicate_source_db", Type: "string", Required: false, Enum: []string{}, Description: "Source DB identifier for read replica"},
			{Name: "tags", Type: "object", Required: false, Enum: []string{}, Description: "Resource tags"},
		},
		Outputs: map[string]string{
			"id":       "string",
			"arn":      "string",
			"endpoint": "string",
			"port":     "int",
			"status":   "string",
			"address":  "string",
		},
	})
}
