# AWS Adapters Layer

This directory contains adapter implementations that bridge the domain layer and AWS-specific services.

## Purpose

The adapter layer implements the **Adapter Pattern** to:
- Decouple domain layer from cloud-specific implementations
- Enable easy provider swapping
- Provide a clean abstraction for domain services
- Handle conversion between domain and cloud models

## Structure

```
adapters/
└── networking/
    ├── adapter.go          # Main adapter implementation
    ├── adapter_test.go     # Adapter tests
    ├── factory.go          # Factory pattern for adapter creation
    └── README.md           # Detailed documentation
```

## Architecture Flow

### Complete Flow with Input and Output Models

```
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                            │
│  (Cloud-Agnostic Business Logic)                            │
│                                                             │
│  NetworkingService Interface                               │
│  - Domain models (no ID/ARN initially)                     │
│  - Domain models (with ID/ARN after creation)              │
└──────────────────────┬────────────────────────────────────┘
                        │
                        │ implements
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                  Adapter Layer                              │
│  (This Package)                                             │
│                                                             │
│  AWSNetworkingAdapter                                       │
│  - Implements NetworkingService                             │
│  - Uses Mappers for conversion                              │
│  - Handles validation at both layers                        │
│  - Maps: Domain → AWS Input → AWS Output → Domain          │
└──────────────────────┬────────────────────────────────────┘
                        │
                        │ uses
                        ▼
┌─────────────────────────────────────────────────────────────┐
│              AWS Service Layer                              │
│  (Cloud-Specific Implementation)                            │
│                                                             │
│  AWSNetworkingService Interface                            │
│  - Accepts: AWS Input Models                               │
│  - Returns: AWS Output Models (with ID/ARN/State)          │
└──────────────────────┬────────────────────────────────────┘
                        │
                        │ uses
                        ▼
┌─────────────────────────────────────────────────────────────┐
│              AWS Models & API                               │
│  (AWS SDK, Terraform, etc.)                                 │
│                                                             │
│  Input Models:  Configuration only                          │
│  Output Models: Configuration + AWS metadata               │
└─────────────────────────────────────────────────────────────┘
```

### Detailed Flow Example: Creating a VPC

```
1. Domain Input (no ID/ARN)
   ┌─────────────────────────┐
   │ Domain VPC              │
   │ - Name: "my-vpc"        │
   │ - Region: "us-east-1"   │
   │ - CIDR: "10.0.0.0/16"   │
   │ - ID: "" (empty)        │
   │ - ARN: nil              │
   └───────────┬─────────────┘
               │
               │ [FromDomainVPC]
               ▼
2. AWS Input Model
   ┌─────────────────────────┐
   │ AWS VPC Input           │
   │ - Name: "my-vpc"        │
   │ - Region: "us-east-1"   │
   │ - CIDR: "10.0.0.0/16"   │
   │ - No ID/ARN fields      │
   └───────────┬─────────────┘
               │
               │ [CreateVPC]
               ▼
3. AWS Output Model
   ┌─────────────────────────┐
   │ AWS VPC Output          │
   │ - Name: "my-vpc"        │
   │ - Region: "us-east-1"   │
   │ - CIDR: "10.0.0.0/16"   │
   │ - ID: "vpc-0a1b2c3..."  │ ← AWS-generated
   │ - ARN: "arn:aws:..."    │ ← AWS-generated
   │ - State: "available"    │ ← AWS metadata
   │ - CreationTime: ...     │ ← AWS metadata
   └───────────┬─────────────┘
               │
               │ [ToDomainVPCFromOutput]
               ▼
4. Domain Output (with ID/ARN)
   ┌─────────────────────────┐
   │ Domain VPC              │
   │ - Name: "my-vpc"        │
   │ - Region: "us-east-1"   │
   │ - CIDR: "10.0.0.0/16"   │
   │ - ID: "vpc-0a1b2c3..."  │ ← Populated!
   │ - ARN: "arn:aws:..."    │ ← Populated!
   └─────────────────────────┘
```

## Key Benefits

1. **Separation of Concerns**: Domain layer never imports AWS code
2. **Testability**: Easy to mock AWS services
3. **Extensibility**: Add new providers without changing domain code
4. **Type Safety**: Compile-time interface compliance
5. **Error Context**: Clear error messages with full context

## Usage Example

### Basic Usage

```go
// 1. Create AWS service implementation
awsService := &awsNetworkingServiceImpl{...}

// 2. Create adapter
adapter := networking.NewAWSNetworkingAdapter(awsService)

// 3. Use as domain service (cloud-agnostic)
var service domainnetworking.NetworkingService = adapter

// 4. Domain layer uses service without knowing about AWS
vpc, err := service.CreateVPC(ctx, domainVPC)
```

### Complete Example with Output Models

```go
import (
    "context"
    domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
    awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/adapters/networking"
)

// 1. Create domain VPC (input - no ID/ARN)
domainVPC := &domainnetworking.VPC{
    Name:   "production-vpc",
    Region: "us-east-1",
    CIDR:   "10.0.0.0/16",
    EnableDNS: true,
}

// 2. Create adapter
adapter := networking.NewAWSNetworkingAdapter(awsService)

// 3. Create VPC through adapter
ctx := context.Background()
createdVPC, err := adapter.CreateVPC(ctx, domainVPC)
if err != nil {
    // handle error
}

// 4. Created VPC now has ID and ARN populated!
fmt.Printf("VPC ID: %s\n", createdVPC.ID)        // "vpc-0a1b2c3d4e5f6g7h8"
fmt.Printf("VPC ARN: %s\n", *createdVPC.ARN)    // "arn:aws:ec2:us-east-1:..."

// 5. Get VPC by ID (returns domain model with ID/ARN)
retrievedVPC, err := adapter.GetVPC(ctx, createdVPC.ID)
// retrievedVPC.ID and retrievedVPC.ARN are populated
```

## Input vs Output Models in Adapters

### What Happens Inside the Adapter

When you call `adapter.CreateVPC()`:

1. **Input Conversion**: Domain model → AWS input model
   ```go
   awsVPCInput := awsmapper.FromDomainVPC(domainVPC)
   // No ID/ARN in awsVPCInput
   ```

2. **Service Call**: AWS service creates resource
   ```go
   awsVPCOutput, err := awsService.CreateVPC(ctx, awsVPCInput)
   // awsVPCOutput has ID, ARN, State, CreationTime, etc.
   ```

3. **Output Conversion**: AWS output model → Domain model
   ```go
   domainVPCWithID := awsmapper.ToDomainVPCFromOutput(awsVPCOutput)
   // domainVPCWithID has ID and ARN populated
   ```

### Key Points

- **Input Models**: Used for creation/update operations (no ID/ARN)
- **Output Models**: Returned by AWS services (with ID/ARN/metadata)
- **Adapter Handles Both**: Automatically converts between input/output models
- **Domain Layer**: Always receives domain models, but with ID/ARN after operations

## Output Service Interfaces

In addition to the standard service interfaces, adapters can also implement output service interfaces that return dedicated output DTOs. This provides a cleaner separation between input configuration and output metadata.

### Output Service Pattern

```go
// Standard service returns full domain model
instance, err := adapter.CreateInstance(ctx, domainInstance)
// instance is *domaincompute.Instance with all fields

// Output service returns focused output DTO
output, err := adapter.CreateInstanceOutput(ctx, domainInstance)
// output is *domaincompute.InstanceOutput with cloud-generated fields
```

### Benefits

1. **Clear Separation**: Input configuration vs. output metadata
2. **Focused Models**: Output DTOs contain only cloud-generated fields
3. **Type Safety**: Dedicated types for outputs prevent confusion
4. **Backward Compatible**: Original service methods still available

### Implementation

Adapters implement both interfaces:

```go
type AWSComputeAdapter struct {
    awsService awsservice.AWSComputeService
}

// Implements ComputeService (standard interface)
var _ domaincompute.ComputeService = (*AWSComputeAdapter)(nil)

// Implements ComputeOutputService (output interface)
var _ domaincompute.ComputeOutputService = (*AWSComputeAdapter)(nil)
```

## Adding New Adapters

To add adapters for other resource types:

1. Create new directory: `adapters/{resource_type}/`
2. Implement domain service interface
3. Implement domain output service interface (optional but recommended)
4. Use mappers for conversion (including output mappers)
5. Add factory pattern
6. Write tests

Example: `adapters/compute/` for EC2, Lambda, etc.

## Future Providers

The adapter pattern makes it easy to add new cloud providers:

1. Create `internal/cloud/{provider}/adapters/networking/`
2. Implement `domainnetworking.NetworkingService`
3. Use provider-specific service and mappers
4. Add to factory pattern

The domain layer remains completely unchanged!
