package database

import (
	"context"
	"fmt"

	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/database"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/database"
	domaindatabase "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/database"
)

// Adapter handles the conversion between domain resources and AWS RDS models
type Adapter struct {
	awsService awsservice.AWSDatabaseService
}

// NewAdapter creates a new database adapter
func NewAdapter(awsService awsservice.AWSDatabaseService) *Adapter {
	return &Adapter{
		awsService: awsService,
	}
}

// CreateRDSInstance converts domain entity to AWS model, calls service, and returns updated domain entity
func (a *Adapter) CreateRDSInstance(ctx context.Context, instance *domaindatabase.RDSInstance) (*domaindatabase.RDSInstance, error) {
	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	// Convert Domain -> AWS
	awsInstance := awsmapper.FromDomainRDSInstance(instance)

	// Call Service
	output, err := a.awsService.CreateRDSInstance(ctx, awsInstance)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	// Convert Output -> Domain
	// Note: The output usually only contains ID/ARN. We should merge back into the original instance or return a new one.
	// Ideally we merge. For now let's create a new one from output and copy config back if needed, OR just return what we have.
	domainInstance := awsmapper.ToDomainRDSInstanceFromOutput(output)

	// Copy back input configuration that isn't in output (mock service doesn't verify existence of all fields in output)
	domainInstance.Name = instance.Name
	domainInstance.Engine = instance.Engine
	domainInstance.EngineVersion = instance.EngineVersion
	domainInstance.InstanceClass = instance.InstanceClass
	domainInstance.AllocatedStorage = instance.AllocatedStorage
	domainInstance.StorageType = instance.StorageType
	domainInstance.Username = instance.Username
	domainInstance.Password = instance.Password
	domainInstance.DBName = instance.DBName
	domainInstance.SubnetGroupName = instance.SubnetGroupName
	domainInstance.VpcSecurityGroupIds = instance.VpcSecurityGroupIds
	domainInstance.SkipFinalSnapshot = instance.SkipFinalSnapshot
	domainInstance.PubliclyAccessible = instance.PubliclyAccessible
	domainInstance.MultiAZ = instance.MultiAZ
	domainInstance.BackupRetentionPeriod = instance.BackupRetentionPeriod
	domainInstance.Tags = instance.Tags

	return domainInstance, nil
}

func (a *Adapter) GetRDSInstance(ctx context.Context, id string) (*domaindatabase.RDSInstance, error) {
	output, err := a.awsService.GetRDSInstance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("aws service error: %w", err)
	}

	return awsmapper.ToDomainRDSInstanceFromOutput(output), nil
}
