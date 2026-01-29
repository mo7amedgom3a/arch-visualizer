# Resource Schema Registry

The schema package provides a dynamic, registry-based approach to resource validation. Instead of hardcoding validation rules per resource type, schemas are defined declaratively and loaded into a registry.

## Overview

This approach enables:
- **Dynamic validation**: Load schemas at runtime from code, files, or database
- **Provider-specific schemas**: Different rules for AWS, Azure, GCP
- **Extensibility**: Add new resource types without modifying validator code
- **Consistency**: Single source of truth for resource specifications
- **Self-documenting**: Schemas describe fields, types, and constraints

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                  Schema Registry                        │
│  ┌─────────────────────────────────────────────────┐   │
│  │  aws                                             │   │
│  │  ├── vpc: {fields: [name, cidr, ...]}           │   │
│  │  ├── subnet: {fields: [name, cidr, az, ...]}    │   │
│  │  ├── ec2: {fields: [name, ami, type, ...]}      │   │
│  │  └── ...                                         │   │
│  ├─────────────────────────────────────────────────┤   │
│  │  azure                                           │   │
│  │  ├── vnet: {fields: [...]}                      │   │
│  │  └── ...                                         │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                    Validator                            │
│  1. Get schema for resource type + provider             │
│  2. Validate required fields                            │
│  3. Validate field types                                │
│  4. Validate constraints (min/max, patterns, enums)     │
│  5. Validate containment relationships                  │
└─────────────────────────────────────────────────────────┘
```

## Data Structures

### ResourceSchema

Defines a resource type's specification:

```go
type ResourceSchema struct {
    ResourceType     string      // e.g., "vpc", "ec2"
    Provider         string      // e.g., "aws", "azure"
    Category         string      // e.g., "networking", "compute"
    Description      string      // Human-readable description
    Fields           []FieldSpec // Field specifications
    ValidParentTypes []string    // What can contain this resource
    ValidChildTypes  []string    // What this resource can contain
}
```

### FieldSpec

Defines a single field:

```go
type FieldSpec struct {
    Name        string           // Field name (e.g., "cidr")
    Type        FieldType        // string, int, bool, cidr, array, object
    Required    bool             // Is this field required?
    Description string           // Human-readable description
    Default     interface{}      // Default value
    Constraints *FieldConstraint // Validation constraints
}
```

### FieldConstraint

Defines validation constraints:

```go
type FieldConstraint struct {
    MinLength   *int     // Minimum string length
    MaxLength   *int     // Maximum string length
    Pattern     *string  // Regex pattern
    MinValue    *float64 // Minimum numeric value
    MaxValue    *float64 // Maximum numeric value
    Enum        []string // Allowed values
    Prefix      *string  // Required prefix (e.g., "ami-")
    CIDRVersion *string  // "ipv4" or "ipv6"
}
```

## Field Types

| Type | Description | Example |
|------|-------------|---------|
| `string` | Text value | `"my-vpc"` |
| `int` | Integer | `8080` |
| `float` | Decimal | `3.14` |
| `bool` | Boolean | `true` |
| `cidr` | CIDR block (validated) | `"10.0.0.0/16"` |
| `array` | List of values | `["sg-1", "sg-2"]` |
| `object` | Nested object | `{"key": "value"}` |
| `any` | Any type | - |

## Usage

### Registering Schemas

```go
import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator/schema"

// Use default registry (pre-loaded with AWS schemas)
schema.DefaultRegistry.Get("vpc", "aws")

// Or create custom registry
registry := schema.NewSchemaRegistry()
registry.Register(&schema.ResourceSchema{
    ResourceType: "custom-resource",
    Provider:     "aws",
    Fields: []schema.FieldSpec{
        {Name: "name", Type: schema.FieldTypeString, Required: true},
    },
})
```

### Using with Validator

```go
import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator"

opts := &validator.ValidationOptions{
    Provider:       "aws",
    SchemaRegistry: schema.DefaultRegistry, // or custom registry
}

result := validator.Validate(graph, opts)
```

### Adding New Provider Schemas

Create a new file (e.g., `azure_schemas.go`):

```go
package schema

func RegisterAzureSchemas(registry *InMemorySchemaRegistry) {
    registry.Register(&ResourceSchema{
        ResourceType: "vnet",
        Provider:     "azure",
        Category:     "networking",
        Fields: []FieldSpec{
            {Name: "name", Type: FieldTypeString, Required: true},
            {Name: "addressSpace", Type: FieldTypeArray, Required: true},
            {Name: "location", Type: FieldTypeString, Required: true},
        },
        ValidChildTypes: []string{"subnet"},
    })
    // ... more Azure resources
}
```

## Pre-registered AWS Schemas

The following AWS resource schemas are pre-registered:

### Networking
- `region` - AWS Region (container)
- `vpc` - Virtual Private Cloud
- `subnet` - VPC Subnet
- `security-group` - Security Group
- `route-table` - Route Table
- `internet-gateway` - Internet Gateway
- `nat-gateway` - NAT Gateway
- `elastic-ip` - Elastic IP Address

### Compute
- `ec2` - EC2 Instance
- `lambda` - Lambda Function
- `load-balancer` - ALB/NLB
- `auto-scaling-group` - Auto Scaling Group

### Storage
- `s3` - S3 Bucket
- `ebs` - EBS Volume

### Database
- `rds` - RDS Instance
- `dynamodb` - DynamoDB Table

## Validation Errors

The schema-driven validator produces these error codes:

| Code | Description |
|------|-------------|
| `CONFIG_MISSING_FIELD` | Required field is missing |
| `CONFIG_INVALID_TYPE` | Field has wrong type |
| `CONFIG_INVALID_CIDR` | CIDR field is not valid |
| `CONFIG_CONSTRAINT_VIOLATION` | Field violates constraint |
| `INVALID_CONTAINMENT` | Invalid parent-child relationship |

## Example Schema

```go
&ResourceSchema{
    ResourceType: "ec2",
    Provider:     "aws",
    Category:     "compute",
    Description:  "EC2 Instance",
    Fields: []FieldSpec{
        {
            Name:        "name",
            Type:        FieldTypeString,
            Required:    true,
            Description: "Instance name",
        },
        {
            Name:        "ami",
            Type:        FieldTypeString,
            Required:    true,
            Description: "AMI ID",
            Constraints: &FieldConstraint{
                Prefix:    strPtr("ami-"),
                MinLength: intPtr(12),
                MaxLength: intPtr(21),
            },
        },
        {
            Name:        "instanceType",
            Type:        FieldTypeString,
            Required:    true,
            Description: "Instance type (e.g., t3.micro)",
        },
    },
    ValidParentTypes: []string{"subnet"},
}
```

## Extending the Schema System

### Loading from JSON/YAML

You could extend the system to load schemas from files:

```go
func LoadSchemasFromFile(registry *InMemorySchemaRegistry, path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }
    
    var schemas []ResourceSchema
    if err := json.Unmarshal(data, &schemas); err != nil {
        return err
    }
    
    for _, s := range schemas {
        if err := registry.Register(&s); err != nil {
            return err
        }
    }
    return nil
}
```

### Loading from Database

```go
func LoadSchemasFromDB(registry *InMemorySchemaRegistry, db *gorm.DB) error {
    var dbSchemas []models.ResourceSchema
    if err := db.Find(&dbSchemas).Error; err != nil {
        return err
    }
    
    for _, s := range dbSchemas {
        registry.Register(convertToSchema(s))
    }
    return nil
}
```

## Best Practices

1. **Register early**: Register schemas during application initialization
2. **Use descriptive names**: Field names should match config keys
3. **Add descriptions**: Help users understand what fields mean
4. **Set appropriate constraints**: Catch errors early
5. **Define containment rules**: Validate parent-child relationships
6. **Test schemas**: Write tests for custom schemas

## Related

- [Validator README](../README.md)
- [AWS Cloud Models](../../../cloud/aws/models/README.md)
