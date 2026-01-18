package compute

import (
	"fmt"

	awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
)

// ComputeAdapterFactory creates compute adapters for different cloud providers
// This factory pattern allows for easy extension to other cloud providers
type ComputeAdapterFactory interface {
	CreateComputeAdapter() domaincompute.ComputeService
}

// AWSComputeAdapterFactory creates AWS-specific compute adapters
type AWSComputeAdapterFactory struct {
	awsService awsservice.AWSComputeService
}

// NewAWSComputeAdapterFactory creates a new AWS adapter factory
func NewAWSComputeAdapterFactory(awsService awsservice.AWSComputeService) ComputeAdapterFactory {
	return &AWSComputeAdapterFactory{
		awsService: awsService,
	}
}

// CreateComputeAdapter creates an AWS compute adapter
func (f *AWSComputeAdapterFactory) CreateComputeAdapter() domaincompute.ComputeService {
	return NewAWSComputeAdapter(f.awsService)
}

// CreateComputeAdapterForProvider is a convenience function that creates the appropriate adapter
// based on the cloud provider. This can be extended for GCP, Azure, etc.
func CreateComputeAdapterForProvider(provider string, awsService awsservice.AWSComputeService) (domaincompute.ComputeService, error) {
	switch provider {
	case "aws":
		factory := NewAWSComputeAdapterFactory(awsService)
		return factory.CreateComputeAdapter(), nil
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
