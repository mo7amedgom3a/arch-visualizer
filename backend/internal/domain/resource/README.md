# Domain Resource Layer

This package contains cloud-agnostic resource models following Domain-Driven Design (DDD) principles.

## Architecture

The domain resource layer provides:

- **Cloud-Agnostic Models**: No provider-specific code
- **Domain-First Design**: Business concepts, not implementation details
- **Validation**: Domain-level validation rules
- **Service Interfaces**: Polymorphic handling via service interfaces
- **Output Models**: Dedicated DTOs for cloud-generated fields and runtime state

## Structure

```
resource/
├── compute/          # Compute resources (instances, load balancers, Lambda, etc.)
├── networking/       # Networking resources (VPCs, subnets, security groups, etc.)
├── storage/          # Storage resources (EBS volumes, S3 buckets, etc.)
├── iam/              # IAM resources (roles, policies, users, groups, etc.)
├── resource.go       # Base resource types and ResourceOutput
└── types.go          # Common types
```

## Service Interfaces

Each resource category provides two service interfaces:

### Standard Service Interface

Returns full domain models with output fields populated:

```go
type ComputeService interface {
    CreateInstance(ctx context.Context, instance *Instance) (*Instance, error)
    GetInstance(ctx context.Context, id string) (*Instance, error)
    // ... other operations
}
```

### Output Service Interface

Returns dedicated output DTOs focused on cloud-generated fields:

```go
type ComputeOutputService interface {
    CreateInstanceOutput(ctx context.Context, instance *Instance) (*InstanceOutput, error)
    GetInstanceOutput(ctx context.Context, id string) (*InstanceOutput, error)
    // ... other operations
}
```

## Output Models

Each resource package includes dedicated output DTOs that focus on:

- **Cloud-Generated Identifiers**: ID, ARN, unique identifiers
- **Runtime State**: Current state, health status
- **Calculated Outputs**: DNS names, zone IDs, connection strings
- **Timestamps**: Created at, updated at, last modified

### Example: InstanceOutput

```go
type InstanceOutput struct {
    ID               string
    ARN              *string
    Name             string
    Region           string
    AvailabilityZone *string
    InstanceType     string
    AMI              string
    SubnetID         string
    SecurityGroupIDs []string
    PrivateIP        *string
    PublicIP         *string
    PrivateDNS       *string
    PublicDNS        *string
    VPCID            *string
    State            InstanceState
    CreatedAt        *time.Time
}
```

### Converting Domain Models to Output

Helper functions are provided in each package:

```go
// Convert instance to output
instance := &compute.Instance{...}
output := compute.ToInstanceOutput(instance)

// Convert VPC to output
vpc := &networking.VPC{...}
output := networking.ToVPCOutput(vpc)
```

## Generic Output Envelope

The `internal/domain/output` package provides a generic result envelope:

```go
type Result[T any] struct {
    Data     T
    Metadata *Metadata
}
```

This can be used to wrap output DTOs with additional metadata when needed.

## Resource Output

The base `ResourceOutput` type provides common output fields:

```go
type ResourceOutput struct {
    ID        string
    ARN       *string
    Name      string
    Type      ResourceType
    Provider  CloudProvider
    Region    string
    State     *string
    CreatedAt *time.Time
    // ... other fields
}
```

## Usage Patterns

### Standard Service Pattern

```go
// Create resource
instance, err := computeService.CreateInstance(ctx, instance)
// instance is *compute.Instance with ID/ARN populated

// Get resource
instance, err := computeService.GetInstance(ctx, id)
// instance has all fields including output fields
```

### Output Service Pattern

```go
// Create resource and get output
output, err := computeOutputService.CreateInstanceOutput(ctx, instance)
// output is *compute.InstanceOutput with cloud-generated fields

// Get resource output
output, err := computeOutputService.GetInstanceOutput(ctx, id)
// output focuses on runtime state and metadata
```

## Benefits

1. **Separation of Concerns**: Input configuration vs. output metadata
2. **Type Safety**: Dedicated types prevent confusion
3. **Backward Compatible**: Original service methods still available
4. **Focused Models**: Output DTOs contain only relevant fields
5. **Cloud-Agnostic**: Works with any cloud provider implementation

## Related Documentation

- [Compute Resources](compute/README.md) - Compute resource documentation
- [Networking Resources](networking/README.md) - Networking resource documentation
- [Storage Resources](storage/) - Storage resource documentation
- [IAM Resources](iam/) - IAM resource documentation
- [Output Package](../../output/models.go) - Generic output envelope
