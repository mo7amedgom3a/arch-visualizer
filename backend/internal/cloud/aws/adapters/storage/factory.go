package storage

import (
	"fmt"
	domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/storage"
)

// StorageAdapterFactory creates storage adapters for different cloud providers
// This factory pattern allows for easy extension to other cloud providers
type StorageAdapterFactory interface {
	CreateStorageAdapter() domainstorage.StorageService
}

// AWSStorageAdapterFactory creates AWS-specific storage adapters
type AWSStorageAdapterFactory struct {
	awsService awsservice.AWSStorageService
}

// NewAWSStorageAdapterFactory creates a new AWS adapter factory
func NewAWSStorageAdapterFactory(awsService awsservice.AWSStorageService) StorageAdapterFactory {
	return &AWSStorageAdapterFactory{
		awsService: awsService,
	}
}

// CreateStorageAdapter creates an AWS storage adapter
func (f *AWSStorageAdapterFactory) CreateStorageAdapter() domainstorage.StorageService {
	return NewAWSStorageAdapter(f.awsService)
}

// CreateStorageAdapterForProvider is a convenience function that creates the appropriate adapter
// based on the cloud provider. This can be extended for GCP, Azure, etc.
func CreateStorageAdapterForProvider(provider string, awsService awsservice.AWSStorageService) (domainstorage.StorageService, error) {
	switch provider {
	case "aws":
		factory := NewAWSStorageAdapterFactory(awsService)
		return factory.CreateStorageAdapter(), nil
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
