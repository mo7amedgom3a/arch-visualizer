package database

import (
	"context"

	awsdatabase "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/database"
	awsdatabaseoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/database/outputs"
)

// AWSDatabaseService defines AWS-specific database operations
type AWSDatabaseService interface {
	// RDS Instance operations
	CreateRDSInstance(ctx context.Context, instance *awsdatabase.RDSInstance) (*awsdatabaseoutputs.RDSInstanceOutput, error)
	GetRDSInstance(ctx context.Context, id string) (*awsdatabaseoutputs.RDSInstanceOutput, error)
	UpdateRDSInstance(ctx context.Context, id string, instance *awsdatabase.RDSInstance) (*awsdatabaseoutputs.RDSInstanceOutput, error)
	DeleteRDSInstance(ctx context.Context, id string) error
	ListRDSInstances(ctx context.Context, filters map[string][]string) ([]*awsdatabaseoutputs.RDSInstanceOutput, error)
}
