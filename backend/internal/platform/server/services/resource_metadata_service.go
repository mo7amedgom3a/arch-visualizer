package services

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/iam"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/storage"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// ResourceMetadataServiceImpl implements ResourceMetadataService.
type ResourceMetadataServiceImpl struct {
	networkingMeta networking.NetworkingMetadataService
	computeMeta    compute.ComputeMetadataService
	storageMeta    storage.StorageMetadataService
	databaseMeta   database.DatabaseMetadataService
	iamMeta        iam.IAMMetadataService
}

// NewResourceMetadataService creates a new resource metadata service.
func NewResourceMetadataService(
	networkingMeta networking.NetworkingMetadataService,
	computeMeta compute.ComputeMetadataService,
	storageMeta storage.StorageMetadataService,
	databaseMeta database.DatabaseMetadataService,
	iamMeta iam.IAMMetadataService,
) serverinterfaces.ResourceMetadataService {
	return &ResourceMetadataServiceImpl{
		networkingMeta: networkingMeta,
		computeMeta:    computeMeta,
		storageMeta:    storageMeta,
		databaseMeta:   databaseMeta,
		iamMeta:        iamMeta,
	}
}

func (s *ResourceMetadataServiceImpl) GetResourceSchema(ctx context.Context, provider, service, resource string) (*serverinterfaces.ResourceSchemaDTO, error) {
	if provider != "aws" {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	var schema *services.ResourceSchema
	var err error

	switch service {
	case "networking":
		schema, err = s.networkingMeta.GetResourceSchema(ctx, resource)
	case "compute":
		schema, err = s.computeMeta.GetResourceSchema(ctx, resource)
	case "storage":
		schema, err = s.storageMeta.GetResourceSchema(ctx, resource)
	case "database":
		schema, err = s.databaseMeta.GetResourceSchema(ctx, resource)
	case "iam":
		schema, err = s.iamMeta.GetResourceSchema(ctx, resource)
	default:
		return nil, fmt.Errorf("unsupported service: %s", service)
	}

	if err != nil {
		return nil, err
	}
	return mapSchemaToDTO(schema), nil
}

func (s *ResourceMetadataServiceImpl) ListResourceSchemas(ctx context.Context, provider, service string) ([]*serverinterfaces.ResourceSchemaDTO, error) {
	if provider != "aws" {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	var schemas []*services.ResourceSchema
	var err error

	switch service {
	case "networking":
		schemas, err = s.networkingMeta.ListResourceSchemas(ctx)
	case "compute":
		schemas, err = s.computeMeta.ListResourceSchemas(ctx)
	case "storage":
		schemas, err = s.storageMeta.ListResourceSchemas(ctx)
	case "database":
		schemas, err = s.databaseMeta.ListResourceSchemas(ctx)
	case "iam":
		schemas, err = s.iamMeta.ListResourceSchemas(ctx)
	default:
		return nil, fmt.Errorf("unsupported service: %s", service)
	}

	if err != nil {
		return nil, err
	}

	dtos := make([]*serverinterfaces.ResourceSchemaDTO, 0, len(schemas))
	for _, schema := range schemas {
		dtos = append(dtos, mapSchemaToDTO(schema))
	}
	return dtos, nil
}

// mapSchemaToDTO converts a cloud-layer ResourceSchema to a service-layer DTO.
func mapSchemaToDTO(schema *services.ResourceSchema) *serverinterfaces.ResourceSchemaDTO {
	fields := make([]serverinterfaces.FieldDescriptorDTO, 0, len(schema.Fields))
	for _, f := range schema.Fields {
		fields = append(fields, serverinterfaces.FieldDescriptorDTO{
			Name:        f.Name,
			Type:        f.Type,
			Required:    f.Required,
			Enum:        f.Enum,
			Default:     f.Default,
			Description: f.Description,
		})
	}
	return &serverinterfaces.ResourceSchemaDTO{
		Label:   schema.Label,
		Fields:  fields,
		Outputs: schema.Outputs,
	}
}
