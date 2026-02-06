package database

import (
	"fmt"

	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/database"
	domaindatabase "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/database"
)

// DatabaseAdapterFactory creates database adapters for different cloud providers
// This factory pattern allows for easy extension to other cloud providers
type DatabaseAdapterFactory interface {
	CreateDatabaseAdapter() domaindatabase.DatabaseService
}

// AWSDatabaseAdapterFactory creates AWS-specific database adapters
type AWSDatabaseAdapterFactory struct {
	awsService awsservice.AWSDatabaseService
}

// NewAWSDatabaseAdapterFactory creates a new AWS adapter factory
func NewAWSDatabaseAdapterFactory(awsService awsservice.AWSDatabaseService) DatabaseAdapterFactory {
	return &AWSDatabaseAdapterFactory{
		awsService: awsService,
	}
}

// CreateDatabaseAdapter creates an AWS database adapter
func (f *AWSDatabaseAdapterFactory) CreateDatabaseAdapter() domaindatabase.DatabaseService {
	// We need to return the Adapter struct which implements domaindatabase.DatabaseService
	// But header of adapter.go says `type Adapter struct`
	// Let's assume Adapter struct satisfies the interface.
	// Wait, Adapter struct in adapter.go needs to be exported or used via NewAdapter
	return NewAdapter(f.awsService)
}

// CreateDatabaseAdapterForProvider is a convenience function that creates the appropriate adapter
// based on the cloud provider. This can be extended for GCP, Azure, etc.
func CreateDatabaseAdapterForProvider(provider string, awsService awsservice.AWSDatabaseService) (domaindatabase.DatabaseService, error) {
	switch provider {
	case "aws":
		factory := NewAWSDatabaseAdapterFactory(awsService)
		return factory.CreateDatabaseAdapter(), nil
	case "gcp":
		// TODO: Implement GCP adapter
		return nil, fmt.Errorf("gcp adapter not yet implemented")
	case "azure":
		// TODO: Implement Azure adapter
		return nil, fmt.Errorf("azure adapter not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
