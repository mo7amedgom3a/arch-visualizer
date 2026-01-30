# AWS Inventory Layer

The inventory layer provides a **cloud provider-specific resource classification and function registry system** that replaces hardcoded switch statements with dynamic dispatch. This enables extensible, maintainable resource handling across all layers of the AWS provider implementation.

## Purpose

The inventory system serves as a **single source of truth** for:
- **Resource Classifications**: Organizing resources by category (Networking, Compute, Storage, etc.)
- **Function Registry**: Dynamic dispatch of mapper, pricing, and transformation functions
- **Type Mapping**: Mapping IR types (from diagrams) to domain resource names
- **Alias Support**: Supporting multiple naming conventions (kebab-case, snake_case, abbreviations)

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│              Domain Layer (Cloud-Agnostic)              │
│  - Category constants (Networking, Compute, etc.)        │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│         AWS Inventory (This Package)                     │
│  - ResourceClassification: metadata about resources     │
│  - FunctionRegistry: function references                 │
│  - Inventory: manages classifications & functions         │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│         Consumer Layers                                   │
│  - Terraform Mapper: uses inventory for mapping         │
│  - Pricing Service: uses inventory for pricing lookups  │
│  - Resource Mapper: uses inventory for transformations  │
└─────────────────────────────────────────────────────────┘
```

## Components

### ResourceClassification

Defines how a resource is classified and identified:

```go
type ResourceClassification struct {
    Category     string   // "Networking", "Compute", etc.
    ResourceName string   // "VPC", "Subnet", "EC2", etc.
    Aliases      []string // ["vpc", "internet-gateway"] for IR mapping
    IRType       string   // IR resource type (kebab-case) for diagram parsing
}
```

**Example**:
```go
{
    Category:     resource.CategoryNetworking,
    ResourceName: "VPC",
    IRType:       "vpc",
    Aliases:      []string{"vpc"},
}
```

### FunctionRegistry

Holds function references for dynamic dispatch:

```go
type FunctionRegistry struct {
    TerraformMapper    func(*resource.Resource) ([]tfmapper.TerraformBlock, error)
    PricingCalculator  func(*resource.Resource, time.Duration) (*pricing.CostEstimate, error)
    GetPricingInfo     func(string) (*pricing.ResourcePricing, error)
    DomainMapper       func(interface{}) (interface{}, error)  // optional
    AWSModelMapper     func(interface{}) (interface{}, error) // optional
}
```

### Inventory

Manages all resource classifications and function mappings:

```go
type Inventory struct {
    Classifications map[string]ResourceClassification // resourceName -> classification
    Functions       map[string]FunctionRegistry       // resourceName -> functions
    ByCategory      map[string][]string               // category -> []resourceNames
    ByIRType        map[string]string                 // IR type -> resourceName
}
```

## File Structure

```
inventory/
├── README.md           # This file
├── inventory.go        # Core inventory structures and methods
├── resources.go         # AWS resource classifications
└── registry.go         # Inventory initialization and function setters
```

## Usage

### 1. Getting the Default Inventory

```go
import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/inventory"

inv := inventory.GetDefaultInventory()
```

### 2. Checking Resource Support

```go
if inv.SupportsResource("VPC") {
    // Resource is supported
}
```

### 3. Getting Resource Classification

```go
classification, ok := inv.GetResourceClassification("VPC")
if ok {
    fmt.Printf("Category: %s\n", classification.Category)
    fmt.Printf("Resource Name: %s\n", classification.ResourceName)
}
```

### 4. Using Function Registry

```go
functions, ok := inv.GetFunctions("VPC")
if ok && functions.TerraformMapper != nil {
    blocks, err := functions.TerraformMapper(resource)
    // Use Terraform blocks
}
```

### 5. Getting Resources by Category

```go
networkingResources := inv.GetResourcesByCategory(resource.CategoryNetworking)
// Returns: ["VPC", "Subnet", "RouteTable", "SecurityGroup", ...]
```

### 6. Mapping IR Type to Resource Name

```go
resourceName, ok := inv.GetResourceNameByIRType("vpc")
if ok {
    // resourceName = "VPC"
}
```

## Integration Points

### Terraform Mapper

The Terraform mapper uses inventory for dynamic dispatch:

```go
// Before (switch statement):
switch res.Type.Name {
case "VPC":
    return m.mapVPC(res)
case "Subnet":
    return m.mapSubnet(res)
// ...
}

// After (inventory-based):
inv := inventory.GetDefaultInventory()
functions, ok := inv.GetFunctions(res.Type.Name)
if ok && functions.TerraformMapper != nil {
    return functions.TerraformMapper(res)
}
```

**File**: `internal/cloud/aws/mapper/terraform/mapper.go`

### Pricing Service

The pricing service uses inventory for resource type lookups:

```go
// Before (switch statement):
switch resourceType {
case "nat_gateway":
    return networking.GetNATGatewayPricing(region), nil
// ...
}

// After (inventory-based):
inv := inventory.GetDefaultInventory()
resourceName := mapPricingTypeToResourceName(resourceType)
if functions, ok := inv.GetFunctions(resourceName); ok && functions.GetPricingInfo != nil {
    return functions.GetPricingInfo(region)
}
```

**File**: `internal/cloud/aws/pricing/service.go`

### Pricing Calculator

The pricing calculator uses inventory for cost calculation:

```go
// Before (switch statement):
switch res.Type.Name {
case "nat_gateway":
    // Calculate NAT Gateway cost
// ...
}

// After (inventory-based):
inv := inventory.GetDefaultInventory()
if functions, ok := inv.GetFunctions(res.Type.Name); ok && functions.PricingCalculator != nil {
    return functions.PricingCalculator(res, duration)
}
```

**File**: `internal/cloud/aws/pricing/calculator.go`

## Adding New Resources

To add a new AWS resource to the inventory:

### 1. Add Resource Classification

Edit `resources.go`:

```go
{
    Category:     resource.CategoryNetworking,
    ResourceName: "NetworkACL",
    IRType:       "network-acl",
    Aliases:      []string{"network-acl", "network_acl", "nacl"},
},
```

### 2. Register Mapper Function

In the Terraform mapper's `New()` function (`mapper/terraform/mapper.go`):

```go
inv.SetTerraformMapper("NetworkACL", mapper.mapNetworkACL)
```

### 3. Implement Mapper Function

```go
func (m *AWSMapper) mapNetworkACL(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
    // Implementation
}
```

The inventory system will automatically:
- Index the resource by category
- Map IR types and aliases to the resource name
- Enable dynamic dispatch for all registered functions

## Benefits

### 1. Extensibility
- Add new resources by updating inventory, not by modifying switch statements
- No need to touch multiple files when adding support for a new resource

### 2. Maintainability
- Single source of truth for resource classifications
- Clear separation between resource definitions and implementation logic

### 3. Testability
- Easy to mock inventory for testing
- Can create test inventories with specific resource subsets

### 4. Provider Independence
- Each cloud provider (AWS, GCP, Azure) defines its own inventory
- No hardcoded provider-specific logic in shared code

### 5. Type Safety
- Compile-time checks for function signatures
- Clear contracts for function registries

## Resource Categories

Resources are organized into the following categories (defined in `internal/domain/resource/categories.go`):

- **Networking**: VPC, Subnet, RouteTable, SecurityGroup, InternetGateway, NATGateway, ElasticIP
- **Compute**: EC2, Lambda, LoadBalancer, AutoScalingGroup
- **Storage**: S3, EBS
- **Database**: RDS, DynamoDB
- **IAM**: (future)
- **Monitoring**: (future)
- **Security**: (future)
- **Analytics**: (future)
- **Application**: (future)

## Initialization Flow

1. **Package Init** (`registry.go`):
   - Creates `DefaultAWSInventory`
   - Registers all resource classifications from `resources.go`
   - Functions are registered lazily by their respective services

2. **Terraform Mapper Init** (`mapper/terraform/mapper.go`):
   - Calls `New()` which registers Terraform mapper functions
   - Each resource type gets its specific mapper function

3. **Pricing Service Init** (`pricing/service.go`):
   - Calls `NewAWSPricingService()` which registers pricing functions
   - Each resource type gets pricing calculator and info getter functions

## Migration from Switch Statements

The inventory system maintains **backward compatibility** with fallback switch statements:

```go
// Primary path: use inventory
if functions, ok := inv.GetFunctions(res.Type.Name); ok && functions.TerraformMapper != nil {
    return functions.TerraformMapper(res)
}

// Fallback: switch statement (for resources not yet migrated)
return m.mapResourceFallback(res)
```

This allows gradual migration and ensures existing code continues to work during the transition.

## Future Enhancements

- **Validation Rules**: Store validation rules in inventory
- **Schema Definitions**: Include resource schemas in classifications
- **Dependency Rules**: Define resource dependencies and constraints
- **Multi-Provider Support**: Extend to GCP, Azure with their own inventories
- **Dynamic Loading**: Load resource definitions from configuration files

## See Also

- `internal/domain/resource/categories.go` - Category constants
- `internal/cloud/aws/mapper/terraform/mapper.go` - Terraform mapper usage
- `internal/cloud/aws/pricing/service.go` - Pricing service usage
- `internal/cloud/aws/pricing/calculator.go` - Pricing calculator usage
