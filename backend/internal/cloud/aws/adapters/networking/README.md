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
- Accept domain models as input (no ID/ARN initially)
- Validate domain models
- Convert to AWS input models using mappers
- Validate AWS input models
- Call AWS service (returns output models with ID/ARN)
- Convert AWS output models back to domain models (with ID/ARN populated)
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

// Create VPC using domain model (input - no ID/ARN)
vpc := &domainnetworking.VPC{
    Name:   "my-vpc",
    Region: "us-east-1",
    CIDR:   "10.0.0.0/16",
    EnableDNS: true,
    // ID: "" (empty before creation)
    // ARN: nil (nil before creation)
}

// Create VPC - adapter handles input/output conversion
createdVPC, err := networkingService.CreateVPC(ctx, vpc)
if err != nil {
    // handle error
}

// Created VPC now has ID and ARN populated!
fmt.Printf("VPC ID: %s\n", createdVPC.ID)        // "vpc-0a1b2c3d4e5f6g7h8"
fmt.Printf("VPC ARN: %s\n", *createdVPC.ARN)    // "arn:aws:ec2:us-east-1:..."

// Get VPC by ID (also returns domain model with ID/ARN)
retrievedVPC, err := networkingService.GetVPC(ctx, createdVPC.ID)
// retrievedVPC has ID and ARN populated from AWS output model
```

## Flow Diagram

### Complete Flow with Input and Output Models

```
1. Domain Layer calls: CreateVPC(domainVPC)
   Input: Domain VPC (no ID/ARN)
   ↓
2. Adapter validates: domainVPC.Validate()
   ↓
3. Adapter converts: FromDomainVPC(domainVPC) → awsVPCInput
   Output: AWS VPC Input Model (no ID/ARN)
   ↓
4. Adapter validates: awsVPCInput.Validate()
   ↓
5. Adapter calls: awsService.CreateVPC(awsVPCInput)
   ↓
6. AWS Service returns: awsVPCOutput
   Output: AWS VPC Output Model (with ID, ARN, State, CreationTime)
   ↓
7. Adapter converts: ToDomainVPCFromOutput(awsVPCOutput) → domainVPC
   Output: Domain VPC (with ID and ARN populated)
   ↓
8. Domain Layer receives: domainVPC
   Result: Domain VPC with AWS-generated ID and ARN
```

### Input vs Output Models

**Input Models** (for creation/update):
- Located in: `internal/cloud/aws/models/networking/`
- Contains: Configuration fields only (Name, CIDR, Region, etc.)
- No AWS identifiers: No `ID` or `ARN` fields
- Used when: Creating or updating resources

**Output Models** (from AWS responses):
- Located in: `internal/cloud/aws/models/networking/outputs/`
- Contains: Configuration + AWS-generated metadata
- AWS identifiers: `ID`, `ARN`, `State`, `CreationTime`, etc.
- Used when: Receiving responses from AWS services

### Example: VPC Creation Flow

```go
// Step 1: Domain input (no ID/ARN)
domainVPC := &domainnetworking.VPC{
    Name:   "my-vpc",
    Region: "us-east-1",
    CIDR:   "10.0.0.0/16",
    // ID: "" (empty)
    // ARN: nil
}

// Step 2: Adapter converts to AWS input
awsVPCInput := awsmapper.FromDomainVPC(domainVPC)
// awsVPCInput has no ID/ARN fields

// Step 3: AWS service creates and returns output
awsVPCOutput := &awsoutputs.VPCOutput{
    ID:     "vpc-0a1b2c3d4e5f6g7h8",  // AWS-generated
    ARN:    "arn:aws:ec2:...",         // AWS-generated
    State:  "available",               // AWS metadata
    // ... configuration fields
}

// Step 4: Adapter converts output to domain (with ID/ARN)
createdVPC := awsmapper.ToDomainVPCFromOutput(awsVPCOutput)
// createdVPC.ID = "vpc-0a1b2c3d4e5f6g7h8"
// createdVPC.ARN = "arn:aws:ec2:..."
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
