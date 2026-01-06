package networking

import (
	"fmt"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/networking"
)

// NetworkingAdapterFactory creates networking adapters for different cloud providers
// This factory pattern allows for easy extension to other cloud providers
type NetworkingAdapterFactory interface {
	CreateNetworkingAdapter() domainnetworking.NetworkingService
}

// AWSNetworkingAdapterFactory creates AWS-specific networking adapters
type AWSNetworkingAdapterFactory struct {
	awsService awsservice.AWSNetworkingService
}

// NewAWSNetworkingAdapterFactory creates a new AWS adapter factory
func NewAWSNetworkingAdapterFactory(awsService awsservice.AWSNetworkingService) NetworkingAdapterFactory {
	return &AWSNetworkingAdapterFactory{
		awsService: awsService,
	}
}

// CreateNetworkingAdapter creates an AWS networking adapter
func (f *AWSNetworkingAdapterFactory) CreateNetworkingAdapter() domainnetworking.NetworkingService {
	return NewAWSNetworkingAdapter(f.awsService)
}

// CreateNetworkingAdapterForProvider is a convenience function that creates the appropriate adapter
// based on the cloud provider. This can be extended for GCP, Azure, etc.
func CreateNetworkingAdapterForProvider(provider string, awsService awsservice.AWSNetworkingService) (domainnetworking.NetworkingService, error) {
	switch provider {
	case "aws":
		factory := NewAWSNetworkingAdapterFactory(awsService)
		return factory.CreateNetworkingAdapter(), nil
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
