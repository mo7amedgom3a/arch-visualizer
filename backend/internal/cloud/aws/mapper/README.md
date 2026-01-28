# AWS Mapper Layer

The mapper layer provides **bidirectional transformation** between cloud-agnostic domain models and AWS-specific resource models. This layer is essential for maintaining clean separation between the domain layer and cloud provider implementations.

## Purpose

The mapper layer serves as a **translation layer** that:

- **Converts Domain → AWS**: Transforms cloud-agnostic domain models into AWS-specific models
- **Converts AWS → Domain**: Transforms AWS models back into domain models
- **Handles Provider Differences**: Manages field name differences, default values, and provider-specific features
- **Maintains Separation**: Keeps domain layer completely free of AWS-specific code
- **Enables Multi-Cloud**: Allows the same domain models to work with different cloud providers

## Architecture Position

```
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                              │
│  (Cloud-Agnostic Business Logic)                             │
│  - VPC, Subnet, SecurityGroup, etc.                         │
└──────────────────────┬──────────────────────────────────────┘
                        │
                        │ uses
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                  Adapter Layer                               │
│  (Implements Domain Interfaces)                              │
└──────────────────────┬──────────────────────────────────────┘
                        │
                        │ uses
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                  Mapper Layer (This Package)                 │
│  - Domain ↔ AWS Conversion                                  │
│  - Field Mapping                                            │
│  - Default Value Injection                                   │
└──────────────────────┬──────────────────────────────────────┘
                        │
                        │ uses
                        ▼
┌─────────────────────────────────────────────────────────────┐
│              AWS Models & Services                          │
│  (AWS-Specific Implementation)                              │
└─────────────────────────────────────────────────────────────┘
```

## Structure

```
mapper/
├── README.md                    # This file
└── networking/                  # Networking resource mappers
    ├── vpc_mapper.go
    ├── subnet_mapper.go
    ├── internet_gateway_mapper.go
    ├── route_table_mapper.go
    ├── security_group_mapper.go
    ├── nat_gateway_mapper.go
    ├── mapper_test.go
    └── README.md                # Detailed networking mapper docs
```

Each resource category has its own subdirectory with mappers for that category's resources.

## Design Principles

### 1. Bidirectional Mapping

Every mapper provides multiple conversion functions:
- `ToDomain*()` - Converts AWS input model → Domain model (backward compatibility)
- `FromDomain*()` - Converts Domain model → AWS input model
- `ToDomain*FromOutput()` - Converts AWS output model → Domain model (with ID/ARN)
- `ToDomain*OutputFromOutput()` - Converts AWS output model → Domain output DTO (new)

### 2. Null Safety

All mappers handle `nil` inputs gracefully:

```go
func ToDomainVPC(awsVPC *awsnetworking.VPC) *domainnetworking.VPC {
    if awsVPC == nil {
        return nil
    }
    // ... conversion logic
}
```

### 3. Default Value Injection

When converting from domain to AWS, mappers inject AWS-specific defaults:

```go
func FromDomainVPC(domainVPC *domainnetworking.VPC) *awsnetworking.VPC {
    // ...
    awsVPC := &awsnetworking.VPC{
        // ... mapped fields
        InstanceTenancy: "default", // AWS default
        Tags: []configs.Tag{{Key: "Name", Value: domainVPC.Name}}, // Default tags
    }
    return awsVPC
}
```

### 4. Field Transformation

Mappers handle field name and structure differences:

- **Name Mapping**: `EnableDNS` (domain) → `EnableDNSSupport` (AWS)
- **Structure Flattening**: Nested domain structures → AWS flat structures
- **Type Conversion**: Domain types → AWS-specific types

### 5. No Validation

Mappers **do not validate** - they only transform. Validation happens in:
- Domain models: `domainVPC.Validate()`
- AWS models: `awsVPC.Validate()`

## Usage Example

### Basic Conversion

```go
import (
    domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
    awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
    awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/networking"
)

// Domain VPC (cloud-agnostic)
domainVPC := &domainnetworking.VPC{
    Name:   "my-vpc",
    Region: "us-east-1",
    CIDR:   "10.0.0.0/16",
    EnableDNS: true,
    EnableDNSHostnames: true,
}

// Convert to AWS VPC
awsVPC := awsmapper.FromDomainVPC(domainVPC)
// awsVPC now has AWS-specific fields like InstanceTenancy, Tags, etc.

// Validate AWS model
if err := awsVPC.Validate(); err != nil {
    // handle validation error
}

// Use AWS service
awsService.CreateVPC(ctx, awsVPC)

// Convert AWS response back to domain
createdAWSVPC, _ := awsService.GetVPC(ctx, "vpc-123")
domainVPC = awsmapper.ToDomainVPC(createdAWSVPC)
```

### In Adapter Pattern

Mappers are typically used within adapters:

```go
// In adapter.go
func (a *AWSNetworkingAdapter) CreateVPC(ctx context.Context, vpc *domainnetworking.VPC) (*domainnetworking.VPC, error) {
    // Validate domain model
    if err := vpc.Validate(); err != nil {
        return nil, err
    }
    
    // Convert to AWS model
    awsVPC := awsmapper.FromDomainVPC(vpc)
    
    // Validate AWS model
    if err := awsVPC.Validate(); err != nil {
        return nil, err
    }
    
    // Call AWS service
    createdAWSVPC, err := a.awsService.CreateVPC(ctx, awsVPC)
    if err != nil {
        return nil, err
    }
    
    // Convert back to domain
    return awsmapper.ToDomainVPC(createdAWSVPC), nil
}
```

## Mapping Patterns

### 1:1 Field Mapping

Simple fields that map directly:

```go
domainVPC.Name → awsVPC.Name
domainVPC.Region → awsVPC.Region
domainVPC.CIDR → awsVPC.CIDR
```

### Field Name Transformation

Fields with different names but same meaning:

```go
domainVPC.EnableDNS → awsVPC.EnableDNSSupport
domainSubnet.IsPublic → awsSubnet.MapPublicIPOnLaunch
```

### Structure Transformation

Complex structures that need transformation:

```go
// Domain: Single target type
domainRoute.TargetType = "internet_gateway"
domainRoute.TargetID = "igw-123"

// AWS: Multiple optional fields
awsRoute.GatewayID = &"igw-123"  // if TargetType == "internet_gateway"
awsRoute.NatGatewayID = &"nat-123"  // if TargetType == "nat_gateway"
```

### Default Value Injection

AWS-specific defaults added during conversion:

```go
// Domain model doesn't have these
awsVPC.InstanceTenancy = "default"
awsVPC.Tags = []configs.Tag{{Key: "Name", Value: domainVPC.Name}}
```

### List to Single Value

Handling provider limitations:

```go
// Domain: Multiple source groups
domainRule.SourceGroupIDs = []string{"sg-1", "sg-2"}

// AWS: Only one source group (limitation)
if len(domainRule.SourceGroupIDs) > 0 {
    awsRule.SourceSecurityGroupID = &domainRule.SourceGroupIDs[0]
}
```

## Common Mapping Scenarios

### Scenario 1: Creating a Resource (Standard Service)

```
Domain Model (User Input)
    ↓ [FromDomain*]
AWS Model (with defaults)
    ↓ [Validate]
AWS Service (create)
    ↓ [ToDomain*FromOutput]
Domain Model (with ID/ARN)
```

### Scenario 2: Creating a Resource (Output Service)

```
Domain Model (User Input)
    ↓ [FromDomain*]
AWS Model (with defaults)
    ↓ [Validate]
AWS Service (create)
    ↓ [ToDomain*OutputFromOutput]
Domain Output DTO (focused on cloud-generated fields)
```

### Scenario 3: Reading a Resource

```
AWS Service (get)
    ↓ [ToDomain*FromOutput] or [ToDomain*OutputFromOutput]
Domain Model or Output DTO
```

### Scenario 4: Updating a Resource

```
Domain Model (updates)
    ↓ [FromDomain*]
AWS Model
    ↓ [Validate]
AWS Service (update)
    ↓ [ToDomain*FromOutput] or [ToDomain*OutputFromOutput]
Domain Model or Output DTO (updated)
```

## Output Model Mapping

### Purpose

Output mapper functions (`ToDomain*OutputFromOutput`) convert AWS output models directly to domain output DTOs. These DTOs focus on:

- Cloud-generated identifiers (ID, ARN)
- Runtime state information
- Calculated outputs (DNS names, zone IDs, etc.)
- Timestamps (created at, updated at)

### When to Use Output Mappers

Use output mappers when:
- Implementing output service interfaces
- You only need cloud-generated fields
- You want a cleaner separation between input and output
- You're building APIs that return focused response models

### Example

```go
// AWS service returns output model
awsOutput := &awsec2outputs.InstanceOutput{
    ID: "i-123",
    ARN: "arn:aws:...",
    State: "running",
    CreationTime: time.Now(),
    // ... other fields
}

// Convert to output DTO
instanceOutput := awsmapper.ToDomainInstanceOutputFromOutput(awsOutput)
// instanceOutput is *domaincompute.InstanceOutput
// Contains only cloud-generated fields and state
```

## Adding New Mappers

To add mappers for a new resource type:

1. **Create mapper file**: `mapper/{category}/{resource}_mapper.go`

2. **Implement bidirectional functions**:
   ```go
   func ToDomain{Resource}(aws{Resource} *awsmodels.{Resource}) *domain.{Resource}
   func FromDomain{Resource}(domain{Resource} *domain.{Resource}) *awsmodels.{Resource}
   ```

3. **Handle nil inputs**:
   ```go
   if awsResource == nil {
       return nil
   }
   ```

4. **Map fields**:
   - Direct mappings
   - Field name transformations
   - Default value injection
   - Structure transformations

5. **Add tests**: Verify bidirectional conversion

6. **Update documentation**: Add to category README

## Testing

Mappers should be tested for:

- **Bidirectional Conversion**: `ToDomain` then `FromDomain` should preserve data
- **Nil Handling**: Both functions handle nil gracefully
- **Default Values**: `FromDomain` injects correct defaults
- **Field Mapping**: All fields are correctly mapped
- **Edge Cases**: Empty strings, zero values, etc.

See `mapper/networking/mapper_test.go` for examples.

## Benefits

1. **Separation of Concerns**: Domain layer never imports AWS code
2. **Testability**: Easy to test conversion logic independently
3. **Maintainability**: Changes to AWS models don't affect domain layer
4. **Extensibility**: Easy to add new cloud providers with their own mappers
5. **Type Safety**: Compile-time checks ensure correct conversions

## Related Documentation

- [Networking Mappers](../mapper/networking/README.md) - Detailed networking mapper documentation
- [Domain Layer](../../../domain/resource/networking/README.md) - Domain model documentation
- [AWS Models](../models/networking/README.md) - AWS model documentation
- [Adapter Layer](../adapters/networking/README.md) - How adapters use mappers

## Future Enhancements

- **Mapper Registry**: Centralized registry for all mappers
- **Validation Integration**: Optional validation during mapping
- **Batch Mapping**: Efficient bulk conversions
- **Mapping Metadata**: Reflection-based mapping with metadata
- **Error Context**: Enhanced error messages with mapping context
