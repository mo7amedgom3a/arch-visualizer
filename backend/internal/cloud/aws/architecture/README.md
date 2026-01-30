# AWS Architecture Package

This package provides AWS-specific implementations for converting diagram graphs into domain architectures. It implements the cloud provider-specific logic required for AWS resource mapping and architecture generation.

## Purpose

The AWS architecture package is responsible for:
- **Architecture Generation**: Converting diagram graphs to domain architectures with AWS-specific logic
- **Resource Type Mapping**: Mapping IR types and resource names to AWS-specific domain resource types
- **Provider Registration**: Auto-registering AWS components with the domain layer

This package follows the **provider-specific generator pattern**, where each cloud provider (AWS, Azure, GCP) implements its own architecture generation logic, allowing for provider-specific resource names and behaviors.

## Architecture

```
┌─────────────────────────────────────────┐
│  Domain Layer (Cloud-Agnostic)          │
│  - ArchitectureGenerator interface     │
│  - ResourceTypeMapper interface        │
│  - Registry system                      │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│  AWS Architecture Package (This)        │
│  - AWSArchitectureGenerator            │
│  - AWSResourceTypeMapper               │
│  - Auto-registration                    │
└─────────────────────────────────────────┘
```

## Components

### 1. AWSArchitectureGenerator (`generator.go`)

Implements the `ArchitectureGenerator` interface for AWS. Converts a diagram graph into a domain architecture.

**Key Responsibilities**:
- Extracts region from region node
- Maps diagram nodes to domain resources
- Builds containment relationships (parent-child)
- Builds dependency relationships
- Uses AWS resource type mapper for type conversion

**Example Flow**:
```go
diagramGraph → AWSArchitectureGenerator.Generate() → Architecture
```

### 2. AWSResourceTypeMapper (`resource_type_mapper.go`)

Implements the `ResourceTypeMapper` interface for AWS. Maps IR types and resource names to AWS-specific domain resource types.

**Key Responsibilities**:
- Maps IR types (e.g., `"vpc"`, `"ec2"`) to ResourceType
- Maps resource names (e.g., `"VPC"`, `"EC2"`) to ResourceType
- Uses AWS inventory for IR type → resource name mapping
- Defines AWS-specific resource type metadata (Category, Kind, IsRegional, etc.)

**Supported AWS Resources**:
- **Networking**: VPC, Subnet, RouteTable, SecurityGroup, InternetGateway, NATGateway, ElasticIP
- **Compute**: EC2, Lambda, LoadBalancer, AutoScalingGroup
- **Storage**: S3, EBS
- **Database**: RDS, DynamoDB

### 3. Registry (`registry.go`)

Auto-registers AWS components when the package is imported.

**Registration**:
- Registers `AWSArchitectureGenerator` with the domain architecture registry
- Registers `AWSResourceTypeMapper` with the domain resource type mapper registry

**Auto-Initialization**:
```go
func init() {
    // Register AWS architecture generator
    generator := NewAWSArchitectureGenerator()
    architecture.RegisterGenerator(generator)

    // Register AWS resource type mapper
    mapper := NewAWSResourceTypeMapper()
    architecture.RegisterResourceTypeMapper(resource.AWS, mapper)
}
```

## How It Works

### 1. Diagram to Architecture Conversion

When `MapDiagramToArchitecture(diagramGraph, resource.AWS)` is called:

1. **Domain Layer** checks for registered generator for AWS
2. **AWS Generator** is retrieved from registry
3. **Generate()** is called:
   - Extracts region from region node
   - First pass: Builds node ID → resource ID mapping
   - Second pass: Creates domain resources
     - Maps IR type using `AWSResourceTypeMapper`
     - Extracts name, parent ID, dependencies
     - Creates domain resource with metadata
   - Builds containment relationships
   - Builds dependency relationships
4. Returns complete `Architecture` aggregate

### 2. Resource Type Mapping

When mapping an IR type (e.g., `"vpc"`) to ResourceType:

1. **AWSResourceTypeMapper.MapIRTypeToResourceType()** is called
2. Uses AWS inventory to map IR type → resource name (`"vpc"` → `"VPC"`)
3. Maps resource name → ResourceType using AWS-specific mapping
4. Returns ResourceType with AWS-specific metadata

**Example**:
```
IR Type: "vpc"
  ↓ (via AWS inventory)
Resource Name: "VPC"
  ↓ (via AWSResourceTypeMapper)
ResourceType: {
  ID: "vpc",
  Name: "VPC",
  Category: "Networking",
  Kind: "Network",
  IsRegional: true,
  IsGlobal: false
}
```

## File Structure

```
architecture/
├── README.md                  # This file
├── generator.go               # AWSArchitectureGenerator implementation
├── resource_type_mapper.go    # AWSResourceTypeMapper implementation
└── registry.go                # Auto-registration on package import
```

## Usage

### Direct Usage

```go
import (
    "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/architecture"
    "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
    "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// The generator is automatically registered when package is imported
// Just use the domain layer API
arch, err := architecture.MapDiagramToArchitecture(diagramGraph, resource.AWS)
```

### Using Resource Type Mapper

```go
import (
    "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
    "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// Get AWS resource type mapper
mapper, ok := architecture.GetResourceTypeMapper(resource.AWS)
if !ok {
    // Handle error
}

// Map IR type to ResourceType
resourceType, err := mapper.MapIRTypeToResourceType("vpc")
if err != nil {
    // Handle error
}

// Map resource name to ResourceType
resourceType, err := mapper.MapResourceNameToResourceType("VPC")
if err != nil {
    // Handle error
}
```

## Integration Points

### Domain Architecture Layer

The AWS architecture package integrates with:
- `internal/domain/architecture/generator.go` - ArchitectureGenerator interface
- `internal/domain/architecture/resource_type_mapper.go` - ResourceTypeMapper interface
- `internal/domain/architecture/aggregate.go` - MapDiagramToArchitecture function

### AWS Inventory

The resource type mapper uses:
- `internal/cloud/aws/inventory` - For IR type → resource name mapping
- `internal/domain/architecture/ir_type_mapper.go` - IRTypeMapper interface

## Adding New AWS Resources

To add support for a new AWS resource:

### 1. Add to AWS Inventory

Edit `internal/cloud/aws/inventory/resources.go`:

```go
{
    Category:     resource.CategoryNetworking,
    ResourceName: "NetworkACL",
    IRType:       "network-acl",
    Aliases:      []string{"network-acl", "nacl"},
},
```

### 2. Add to Resource Type Mapper

Edit `internal/cloud/aws/architecture/resource_type_mapper.go`:

```go
"NetworkACL": {
    ID:         "network-acl",
    Name:       "NetworkACL",
    Category:   string(resource.CategoryNetworking),
    Kind:       "Network",
    IsRegional: true,
    IsGlobal:   false,
},
```

### 3. Register Terraform Mapper (if needed)

Edit `internal/cloud/aws/mapper/terraform/mapper.go`:

```go
inv.SetTerraformMapper("NetworkACL", mapper.mapNetworkACL)
```

The new resource will automatically be supported in architecture generation!

## Key Design Decisions

### 1. Provider-Specific Mappings

**Why**: Different cloud providers use different resource names:
- AWS: `VPC`, `EC2`, `S3`
- Azure: `VirtualNetwork`, `VirtualMachine`, `StorageAccount`
- GCP: `VPCNetwork`, `ComputeInstance`, `StorageBucket`

**Solution**: Each provider implements its own `ResourceTypeMapper` with provider-specific mappings.

### 2. No Static Fallbacks

**Why**: Static fallbacks create hardcoded assumptions about resource names.

**Solution**: All mappings are provider-specific. If a provider doesn't support a resource type, it returns an error rather than falling back to a generic mapping.

### 3. Two-Pass Resource Creation

**Why**: Parent-child relationships need to be resolved before creating resources.

**Solution**: 
- First pass: Build complete node ID → resource ID mapping
- Second pass: Create resources with resolved parent IDs

### 4. Auto-Registration

**Why**: Simplifies usage - no manual registration required.

**Solution**: Package `init()` function automatically registers components when imported.

## Testing

To test the AWS architecture package:

```bash
# Run tests
go test ./internal/cloud/aws/architecture/... -v

# Test with full integration
go run cmd/api/main.go
```

## Dependencies

- `internal/domain/architecture` - Domain layer interfaces and types
- `internal/domain/resource` - Resource domain models
- `internal/diagram/graph` - Diagram graph types
- `internal/cloud/aws/inventory` - AWS resource inventory

## Future Enhancements

- **Azure Support**: Add `AzureArchitectureGenerator` and `AzureResourceTypeMapper`
- **GCP Support**: Add `GCPArchitectureGenerator` and `GCPResourceTypeMapper`
- **Validation**: Add AWS-specific validation rules in the generator
- **Metadata Enrichment**: Add AWS-specific metadata extraction (tags, regions, etc.)

## See Also

- `internal/domain/architecture/README.md` - Domain architecture layer documentation
- `internal/cloud/aws/inventory/README.md` - AWS inventory system documentation
- `internal/cloud/aws/mapper/terraform/README.md` - Terraform mapper documentation
