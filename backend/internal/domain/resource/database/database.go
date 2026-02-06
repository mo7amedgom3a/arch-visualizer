package database

import (
	"context"
)

// DatabaseService defines the interface for database operations in the domain
// This abstraction allows the domain to interact with database resources without knowing the cloud provider specifics
type DatabaseService interface {
	// Resource operations
	CreateRDSInstance(ctx context.Context, instance *RDSInstance) (*RDSInstance, error)
	GetRDSInstance(ctx context.Context, id string) (*RDSInstance, error)
}
