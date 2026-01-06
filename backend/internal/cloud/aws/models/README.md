# AWS Models Layer

This package contains **AWS-specific resource models** that represent cloud provider implementations. These models capture AWS-specific details, constraints, and features that are not present in the cloud-agnostic domain layer.

## Purpose

The models layer serves as the **AWS-specific representation** of cloud resources:

- **Provider-Specific Details**: Captures AWS-specific fields, constraints, and behaviors
- **Validation Rules**: Implements AWS-specific validation (limits, formats, requirements)
- **API Compatibility**: Models align with AWS API structures and Terraform resource schemas
- **Configuration Support**: JSON-serializable for configuration files and API responses
- **Tag Management**: Standardized AWS tagging support across all resources

## Architecture Position

```
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                              │
│  (Cloud-Agnostic Business Logic)                             │
│  - VPC, Subnet, SecurityGroup, etc.                         │
└──────────────────────┬──────────────────────────────────────┘
                        │
                        │ mapped via
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                  Mapper Layer                                │
│  (Domain ↔ AWS Conversion)                                   │
└──────────────────────┬──────────────────────────────────────┘
                        │
                        │ converts to/from
                        ▼
┌─────────────────────────────────────────────────────────────┐
│              AWS Models Layer (This Package)                 │
│  - AWS VPC, Subnet, SecurityGroup, etc.                     │
│  - AWS-specific fields & validation                         │
└──────────────────────┬──────────────────────────────────────┘
                        │
                        │ used by
                        ▼
┌─────────────────────────────────────────────────────────────┐
│              AWS Services & APIs                             │
│  (Terraform, AWS SDK, etc.)                                  │
└─────────────────────────────────────────────────────────────┘
```

## Structure

```
models/
├── README.md                    # This file
├── networking/                  # Networking resource models
│   ├── vpc.go
│   ├── subnet.go
│   ├── internet_gateway.go
│   ├── route_table.go
│   ├── security_group.go
│   ├── nat_gateway.go
│   ├── tests/
│   │   └── vpc_test.go
│   └── README.md
├── compute/                     # Compute resource models (future)
│   └── .gitkeep
└── storage/                      # Storage resource models (future)
    └── .gitkeep
```

Each resource category has its own subdirectory with models for that category's resources.

## Design Principles

### 1. AWS-Specific Implementation

Models contain AWS-specific details that don't exist in the domain layer:

```go
type VPC struct {
    Name               string `json:"name"`
    Region             string `json:"region"`
    CIDR               string `json:"cidr"`
    EnableDNSHostnames bool   `json:"enable_dns_hostnames"`
    EnableDNSSupport   bool   `json:"enable_dns_support"`
    InstanceTenancy    string `json:"instance_tenancy"`  // AWS-specific
    Tags               []configs.Tag `json:"tags"`       // AWS tags
}
```

### 2. Validation

Every model implements a `Validate()` method with AWS-specific validation rules:

```go
func (vpc *VPC) Validate() error {
    if vpc.Name == "" {
        return errors.New("name is required")
    }
    if vpc.Region == "" {
        return errors.New("region is required")
    }
    if vpc.CIDR == "" {
        return errors.New("cidr is required")
    }
    
    // AWS-specific validation
    _, ipNet, err := net.ParseCIDR(vpc.CIDR)
    if err != nil {
        return fmt.Errorf("invalid cidr format: %w", err)
    }
    
    // AWS VPCs only support IPv4
    if ipNet.IP.To4() == nil {
        return errors.New("cidr must be IPv4")
    }
    
    return nil
}
```

**Validation Rules Include:**
- Required field checks
- Format validation (CIDR, port ranges, etc.)
- AWS-specific constraints (IPv4-only, description length limits, etc.)
- Business rule validation (subnet CIDR within VPC CIDR, etc.)

### 3. JSON Serialization

All models are JSON-serializable for:
- API requests/responses
- Configuration files
- Terraform/Pulumi code generation
- State persistence

```go
type VPC struct {
    Name   string `json:"name"`
    Region string `json:"region"`
    CIDR   string `json:"cidr"`
    // ...
}
```

### 4. Tag Support

All resources support AWS tags via `configs.Tag`:

```go
import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"

type VPC struct {
    // ...
    Tags []configs.Tag `json:"tags"`
}
```

Tags are automatically added during domain-to-AWS conversion via mappers.

### 5. Optional Fields

Fields marked with `// +optional` are not required but provide AWS-specific functionality:

```go
type VPC struct {
    Name   string `json:"name"`           // Required
    Region string `json:"region"`         // Required
    CIDR   string `json:"cidr"`           // Required
    // +optional
    EnableDNSHostnames bool `json:"enable_dns_hostnames"`
    // +optional
    InstanceTenancy string `json:"instance_tenancy"`
}
```

## Usage Example

### Creating and Validating a Model

```go
import (
    awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
    "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// Create AWS VPC
vpc := &awsnetworking.VPC{
    Name:               "my-vpc",
    Region:             "us-east-1",
    CIDR:               "10.0.0.0/16",
    EnableDNSSupport:   true,
    EnableDNSHostnames: true,
    InstanceTenancy:    "default",
    Tags: []configs.Tag{
        {Key: "Name", Value: "my-vpc"},
        {Key: "Environment", Value: "production"},
    },
}

// Validate AWS-specific rules
if err := vpc.Validate(); err != nil {
    // Handle validation error
    log.Fatal(err)
}

// Use with AWS service
awsService.CreateVPC(ctx, vpc)
```

### JSON Serialization

```go
// Serialize to JSON
jsonData, err := json.Marshal(vpc)
if err != nil {
    log.Fatal(err)
}

// Deserialize from JSON
var vpc awsnetworking.VPC
err = json.Unmarshal(jsonData, &vpc)
if err != nil {
    log.Fatal(err)
}
```

## Model Categories

### Networking Models

Located in `networking/` directory:

- **VPC**: Virtual Private Cloud
- **Subnet**: Network subnet within VPC
- **InternetGateway**: Internet gateway for VPC
- **RouteTable**: Route table with routing rules
- **SecurityGroup**: Security group with ingress/egress rules
- **NATGateway**: NAT gateway for private subnet internet access

See [Networking Models README](./networking/README.md) for detailed documentation.

### Compute Models (Future)

Located in `compute/` directory:

- EC2 instances
- Lambda functions
- ECS services
- EKS clusters

### Storage Models (Future)

Located in `storage/` directory:

- S3 buckets
- EBS volumes
- EFS file systems
- RDS databases

## Validation Patterns

### Required Fields

```go
func (vpc *VPC) Validate() error {
    if vpc.Name == "" {
        return errors.New("name is required")
    }
    // ...
}
```

### Format Validation

```go
func (vpc *VPC) Validate() error {
    // Validate CIDR format
    _, ipNet, err := net.ParseCIDR(vpc.CIDR)
    if err != nil {
        return fmt.Errorf("invalid cidr format: %w", err)
    }
    // ...
}
```

### AWS-Specific Constraints

```go
func (sg *SecurityGroup) Validate() error {
    // AWS description length limit
    if len(sg.Description) > 255 {
        return errors.New("security group description must be 255 characters or less")
    }
    // ...
}
```

### Business Rules

```go
func (rt *RouteTable) Validate() error {
    // Validate routes
    for _, route := range rt.Routes {
        // At least one target must be specified
        targetCount := 0
        if route.GatewayID != nil && *route.GatewayID != "" {
            targetCount++
        }
        // ...
        if targetCount == 0 {
            return errors.New("route must have at least one target")
        }
    }
    // ...
}
```

## Relationship to Domain Models

### Domain Models (Cloud-Agnostic)

```go
// Domain VPC - no AWS-specific fields
type VPC struct {
    ID               string
    Name             string
    Region           string
    CIDR             string
    EnableDNS        bool
    EnableDNSHostnames bool
}
```

### AWS Models (Provider-Specific)

```go
// AWS VPC - includes AWS-specific fields
type VPC struct {
    Name               string
    Region             string
    CIDR               string
    EnableDNSHostnames bool
    EnableDNSSupport   bool
    InstanceTenancy    string  // AWS-specific
    Tags               []configs.Tag  // AWS tags
}
```

**Key Differences:**
- AWS models include provider-specific fields (`InstanceTenancy`, `Tags`)
- AWS models use AWS naming conventions (`EnableDNSSupport` vs `EnableDNS`)
- AWS models have AWS-specific validation rules
- AWS models are JSON-serializable for Terraform/API usage

## Conversion Flow

### Domain → AWS

```
Domain VPC
    ↓ [Mapper: FromDomainVPC]
AWS VPC (with defaults)
    ↓ [Validate]
AWS Service/API
```

### AWS → Domain

```
AWS Service/API Response
    ↓ [Mapper: ToDomainVPC]
Domain VPC
```

## Testing

Models should be tested for:

- **Validation**: All validation rules are enforced
- **Edge Cases**: Empty strings, nil values, invalid formats
- **AWS Constraints**: Provider-specific limits and rules
- **JSON Serialization**: Round-trip serialization works correctly

See `networking/tests/vpc_test.go` for examples.

## Best Practices

### 1. Always Validate

```go
vpc := &awsnetworking.VPC{...}
if err := vpc.Validate(); err != nil {
    return err
}
```

### 2. Use Mappers for Conversion

Don't manually convert between domain and AWS models. Use mappers:

```go
// ✅ Good
awsVPC := awsmapper.FromDomainVPC(domainVPC)

// ❌ Bad
awsVPC := &awsnetworking.VPC{
    Name: domainVPC.Name,
    // ... manual conversion
}
```

### 3. Include Tags

Always include tags for AWS resources:

```go
vpc.Tags = []configs.Tag{
    {Key: "Name", Value: vpc.Name},
    {Key: "Environment", Value: "production"},
}
```

### 4. Handle Optional Fields

Optional fields should have sensible defaults:

```go
if vpc.InstanceTenancy == "" {
    vpc.InstanceTenancy = "default"
}
```

## Adding New Models

To add a new AWS model:

1. **Create model file**: `models/{category}/{resource}.go`

2. **Define struct with JSON tags**:
   ```go
   type Resource struct {
       Name string `json:"name"`
       // ...
   }
   ```

3. **Implement Validate() method**:
   ```go
   func (r *Resource) Validate() error {
       // AWS-specific validation
   }
   ```

4. **Add tags support**:
   ```go
   Tags []configs.Tag `json:"tags"`
   ```

5. **Create mapper**: Add conversion functions in `mapper/{category}/`

6. **Write tests**: Add validation tests

7. **Update documentation**: Add to category README

## Related Documentation

- [Networking Models](./networking/README.md) - Detailed networking model documentation
- [Mapper Layer](../mapper/README.md) - How mappers convert domain ↔ AWS models
- [Domain Layer](../../../domain/resource/networking/README.md) - Cloud-agnostic domain models
- [Adapter Layer](../adapters/networking/README.md) - How adapters use models
- [AWS Configs](../configs/utils.go) - Tag and configuration utilities

## AWS-Specific Features

### Instance Tenancy

VPCs support different instance tenancy options:
- `"default"`: Shared tenancy (default)
- `"dedicated"`: Dedicated instance tenancy

### Tag Management

All resources support AWS tags for:
- Resource organization
- Cost allocation
- Automation and filtering

### Region and Availability Zones

- Resources are region-scoped
- Subnets require availability zones
- Some resources are global (e.g., IAM)

### Limits and Constraints

AWS models enforce AWS-specific limits:
- Security group description: 255 characters max
- CIDR blocks: IPv4 only for VPCs
- Route targets: One target per route
- Port ranges: 0-65535

## Future Enhancements

- **Model Registry**: Centralized registry for all models
- **Validation Metadata**: Reflection-based validation with metadata
- **Default Value Injection**: Automatic default value assignment
- **Model Versioning**: Support for AWS API versioning
- **Batch Validation**: Efficient bulk validation
