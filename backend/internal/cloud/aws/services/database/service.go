package database

import (
	"context"
	"fmt"

	awserrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/errors"
	awsdatabase "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/database"
	awsdatabaseoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/database/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services"
	domainerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/errors"
)

// DatabaseService implements AWSDatabaseService with deterministic virtual operations
type DatabaseService struct{}

// NewDatabaseService creates a new database service implementation
func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

// RDS Instance operations

func (s *DatabaseService) CreateRDSInstance(ctx context.Context, instance *awsdatabase.RDSInstance) (*awsdatabaseoutputs.RDSInstanceOutput, error) {
	if instance == nil {
		return nil, domainerrors.New(awserrors.CodeRDSInstanceCreationFailed, domainerrors.KindValidation, "instance is nil").
			WithOp("DatabaseService.CreateRDSInstance")
	}

	instanceID := services.GenerateDeterministicID(instance.Name)[:16]
	// Ensure ID starts with a letter as per AWS naming conventions if needed, but here simple ID is fine
	instanceID = "db-" + instanceID
	region := "us-east-1"
	arn := services.GenerateARN("rds", "db", instanceID, region)

	// Simulate endpoint generation
	endpoint := fmt.Sprintf("%s.c1234567890.%s.rds.amazonaws.com", instanceID, region)
	port := 3306
	if instance.Engine == "postgres" {
		port = 5432
	} else if instance.Engine == "sqlserver-ex" || instance.Engine == "sqlserver-se" || instance.Engine == "sqlserver-ee" || instance.Engine == "sqlserver-web" {
		port = 1433
	} else if instance.Engine == "oracle-ee" || instance.Engine == "oracle-se2" {
		port = 1521
	}

	return &awsdatabaseoutputs.RDSInstanceOutput{
		ID:       instanceID,
		Address:  endpoint,
		Port:     port,
		Endpoint: fmt.Sprintf("%s:%d", endpoint, port),
		ARN:      arn,
	}, nil
}

func (s *DatabaseService) GetRDSInstance(ctx context.Context, id string) (*awsdatabaseoutputs.RDSInstanceOutput, error) {
	region := "us-east-1"
	arn := services.GenerateARN("rds", "db", id, region)
	endpoint := fmt.Sprintf("%s.c1234567890.%s.rds.amazonaws.com", id, region)

	return &awsdatabaseoutputs.RDSInstanceOutput{
		ID:       id,
		Address:  endpoint,
		Port:     3306, // Default mock port
		Endpoint: fmt.Sprintf("%s:%d", endpoint, 3306),
		ARN:      arn,
	}, nil
}

func (s *DatabaseService) UpdateRDSInstance(ctx context.Context, id string, instance *awsdatabase.RDSInstance) (*awsdatabaseoutputs.RDSInstanceOutput, error) {
	return s.CreateRDSInstance(ctx, instance)
}

func (s *DatabaseService) DeleteRDSInstance(ctx context.Context, id string) error {
	return nil
}

func (s *DatabaseService) ListRDSInstances(ctx context.Context, filters map[string][]string) ([]*awsdatabaseoutputs.RDSInstanceOutput, error) {
	id := "db-test-instance"
	region := "us-east-1"
	arn := services.GenerateARN("rds", "db", id, region)
	endpoint := fmt.Sprintf("%s.c1234567890.%s.rds.amazonaws.com", id, region)

	return []*awsdatabaseoutputs.RDSInstanceOutput{
		{
			ID:       id,
			Address:  endpoint,
			Port:     3306,
			Endpoint: fmt.Sprintf("%s:%d", endpoint, 3306),
			ARN:      arn,
		},
	}, nil
}
