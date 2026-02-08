package containers

// Factory provides methods to create container adapters

// NewContainerAdapter creates a new container adapter with default configuration
func NewContainerAdapter() *AWSContainerAdapter {
	return NewAWSContainerAdapter()
}

// ContainerAdapterConfig holds configuration for the container adapter
type ContainerAdapterConfig struct {
	Region string
	// In a full implementation, add AWS credentials config
}

// NewContainerAdapterWithConfig creates a new container adapter with custom configuration
func NewContainerAdapterWithConfig(config ContainerAdapterConfig) *AWSContainerAdapter {
	adapter := NewAWSContainerAdapter()
	// In a full implementation, configure AWS SDK client with region and credentials
	return adapter
}
