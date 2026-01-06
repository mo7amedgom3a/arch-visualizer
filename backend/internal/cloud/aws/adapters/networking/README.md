# AWS Networking Adapter

This package implements the **Adapter Pattern** to bridge the domain layer and AWS-specific implementations.

## Purpose

The adapter layer:
- **Implements domain interfaces** using AWS-specific services
- **Converts between domain and AWS models** using mappers
- **Handles validation** at both domain and AWS levels
- **Provides error translation** from AWS-specific to domain-level errors
- **Enables provider swapping** without changing domain code

## Architecture

```
Domain Layer (Cloud-Agnostic)
    ↓
NetworkingService Interface
    ↓
AWSNetworkingAdapter (This Package)
    ↓
AWSNetworkingService (AWS-Specific)
    ↓
AWS Models & API
```

## Components

### AWSNetworkingAdapter

The main adapter that implements `domainnetworking.NetworkingService` interface.

**Responsibilities:**
- Accept domain models as input
- Validate domain models
- Convert to AWS models using mappers
- Validate AWS models
- Call AWS service
- Convert AWS response back to domain models
- Wrap errors with context

### Factory Pattern

`AWSNetworkingAdapterFactory` provides a factory pattern for creating adapters:

```go
factory := NewAWSNetworkingAdapterFactory(awsService)
adapter := factory.CreateNetworkingAdapter()
```

## Usage Example

```go
import (
    domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
    awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/networking"
    awsadapter "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/adapters/networking"
)

// Create AWS service (implementation)
awsService := &awsNetworkingServiceImpl{...}

// Create adapter
adapter := awsadapter.NewAWSNetworkingAdapter(awsService)

// Use domain interface (cloud-agnostic)
var networkingService domainnetworking.NetworkingService = adapter

// Create VPC using domain model
vpc := &domainnetworking.VPC{
    Name:   "my-vpc",
    Region: "us-east-1",
    CIDR:   "10.0.0.0/16",
    EnableDNS: true,
}

createdVPC, err := networkingService.CreateVPC(ctx, vpc)
```

## Flow Diagram

```
1. Domain Layer calls: CreateVPC(domainVPC)
   ↓
2. Adapter validates: domainVPC.Validate()
   ↓
3. Adapter converts: FromDomainVPC(domainVPC) → awsVPC
   ↓
4. Adapter validates: awsVPC.Validate()
   ↓
5. Adapter calls: awsService.CreateVPC(awsVPC)
   ↓
6. AWS Service returns: awsVPC
   ↓
7. Adapter converts: ToDomainVPC(awsVPC) → domainVPC
   ↓
8. Domain Layer receives: domainVPC
```

## Error Handling

The adapter wraps errors at each layer:

```go
// Domain validation error
if err := vpc.Validate(); err != nil {
    return nil, fmt.Errorf("domain validation failed: %w", err)
}

// AWS validation error
if err := awsVPC.Validate(); err != nil {
    return nil, fmt.Errorf("aws validation failed: %w", err)
}

// AWS service error
if err != nil {
    return nil, fmt.Errorf("aws service error: %w", err)
}
```

This provides clear error context while maintaining error chain for debugging.

## Benefits

1. **Separation of Concerns**: Domain layer never knows about AWS
2. **Testability**: Easy to mock AWS service for testing
3. **Extensibility**: Add new providers by creating new adapters
4. **Type Safety**: Compile-time checks ensure interface compliance
5. **Error Context**: Clear error messages with full context

## Testing

The adapter is tested with a mock AWS service to verify:
- Domain-to-AWS conversion
- AWS-to-domain conversion
- Validation at both layers
- Error handling and wrapping
- All CRUD operations

See `adapter_test.go` for examples.

## Future Extensions

To add a new cloud provider (e.g., GCP):

1. Create `internal/cloud/gcp/adapters/networking/adapter.go`
2. Implement `domainnetworking.NetworkingService`
3. Use GCP-specific service and mappers
4. Add to factory pattern

The domain layer remains unchanged!
