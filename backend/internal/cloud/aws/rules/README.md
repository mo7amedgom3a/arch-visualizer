# AWS Rules Implementation

This package implements AWS-specific rule validation using the domain rules interfaces.

## Purpose

AWS rules implementation:
- **Implements Domain Interfaces**: Uses domain rule interfaces
- **AWS-Specific Mapping**: Maps domain types to AWS resource types
- **AWS Constraints**: Applies AWS-specific limits and constraints
- **Provider-Specific Logic**: Handles AWS-specific validation requirements

## Structure

```
rules/
├── factory.go          # AWS rule factory
├── service.go          # AWS rule service
├── defaults.go         # Default AWS networking rules
└── README.md           # This file
```

## AWS Rule Factory

The `AWSRuleFactory` creates AWS-specific rules:

```go
factory := rules.NewAWSRuleFactory()
rule, err := factory.CreateRule("Subnet", "requires_parent", "VPC")
```

### Resource Type Mapping

AWS maps domain resource types to AWS-specific types:

```go
"VPC"              → "aws_vpc"
"Subnet"           → "aws_subnet"
"InternetGateway"  → "aws_internet_gateway"
"EC2"              → "aws_instance"
```

### AWS-Specific Rules

AWS can override default rule behavior:

```go
// AWS might have different parent requirements
func (f *AWSRuleFactory) createRequiresParentRule(...) {
    awsParentType := f.mapResourceTypeToAWS(parentType)
    return constraints.NewRequiresParentRule(resourceType, awsParentType)
}
```

## AWS Rule Service

The `AWSRuleService` provides rule evaluation for AWS:

```go
service := rules.NewAWSRuleService()

// Load rules from database
constraints := []rules.ConstraintRecord{
    {ResourceType: "Subnet", ConstraintType: "requires_parent", ConstraintValue: "VPC"},
    {ResourceType: "EC2", ConstraintType: "allowed_parent", ConstraintValue: "Subnet"},
}
service.LoadRulesFromConstraints(ctx, constraints)

// Validate a resource
result, err := service.ValidateResource(ctx, resource, architecture)
```

### Default Rules and Merging

AWS provides default networking rules that can be merged with database constraints:

```go
service := rules.NewAWSRuleService()

// Load DB constraints and merge with defaults
// DB constraints override defaults when they have the same resource type + constraint type
dbConstraints := loadConstraintsFromDB("aws")
err := service.LoadRulesWithDefaults(ctx, dbConstraints)
```

**Default Rules Include:**
- **VPC**: Requires region, max 200 children, allowed dependencies on InternetGateway/NATGateway
- **Subnet**: Requires VPC parent, requires region, allowed dependencies on RouteTable/NATGateway, forbidden dependencies on Subnet/VPC
- **InternetGateway**: Requires VPC parent, requires region, allowed dependencies on RouteTable, forbidden dependencies on Subnet/NATGateway
- **RouteTable**: Requires VPC parent, requires region, allowed dependencies on InternetGateway/NATGateway/Subnet, forbidden dependencies on SecurityGroup
- **SecurityGroup**: Requires VPC parent, requires region, allowed dependencies on SecurityGroup, forbidden dependencies on networking infrastructure
- **NATGateway**: Requires Subnet parent, requires region, allowed dependencies on RouteTable, forbidden dependencies on InternetGateway/VPC

**Merge Strategy:**
- Default rules are loaded first
- Database constraints override defaults when they have the same `resource_type` + `constraint_type` combination
- This allows customization while maintaining sensible defaults

## Usage Example

```go
import (
    awsrules "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/rules"
)

// Create AWS rule service
service := awsrules.NewAWSRuleService()

// Option 1: Load only DB constraints (no defaults)
constraints := loadConstraintsFromDB("aws")
service.LoadRulesFromConstraints(ctx, constraints)

// Option 2: Load DB constraints merged with defaults (recommended)
constraints := loadConstraintsFromDB("aws")
service.LoadRulesWithDefaults(ctx, constraints)

// Validate architecture
results, err := service.ValidateArchitecture(ctx, architecture)
for resourceID, result := range results {
    if !result.Valid {
        fmt.Printf("Resource %s has validation errors:\n", resourceID)
        for _, err := range result.Errors {
            fmt.Printf("  - %s\n", err.Message)
        }
    }
}
```

### Dependency Rules Example

Dependency rules validate relationships between resources:

```go
// Subnet can depend on RouteTable and NATGateway
// But cannot depend on VPC or itself (prevents circular dependencies)
service := awsrules.NewAWSRuleService()
service.LoadRulesWithDefaults(ctx, []awsrules.ConstraintRecord{})

// This will pass: Subnet depends on RouteTable
subnet.DependsOn = []string{routeTable.ID}

// This will fail: Subnet depends on VPC (forbidden)
subnet.DependsOn = []string{vpc.ID}
```

## AWS-Specific Considerations

### Resource Type Mapping

AWS uses Terraform-style resource names:
- Domain: `VPC` → AWS: `aws_vpc`
- Domain: `Subnet` → AWS: `aws_subnet`
- Domain: `InternetGateway` → AWS: `aws_internet_gateway`
- Domain: `RouteTable` → AWS: `aws_route_table`
- Domain: `SecurityGroup` → AWS: `aws_security_group`
- Domain: `NATGateway` → AWS: `aws_nat_gateway`
- Domain: `EC2Instance` → AWS: `aws_instance`

All domain types are automatically mapped to AWS types in the factory.

### Limits

AWS has specific limits that may differ from other providers:
- VPC: Max 5 VPCs per region (soft limit)
- Subnet: Max 200 subnets per VPC
- Security Group: Max 5 security groups per instance

### Regional vs Global

AWS-specific regional requirements:
- VPC: Always regional
- S3: Global (but region for data storage)
- IAM: Global

## Extending AWS Rules

To add AWS-specific rules:

1. **Add to factory mapping**:
   ```go
   mapping := map[string]string{
       "NewResource": "aws_new_resource",
   }
   ```

2. **Add custom rule creation**:
   ```go
   func (f *AWSRuleFactory) createCustomRule(...) (rules.Rule, error) {
       // AWS-specific logic
   }
   ```

3. **Update CreateRule switch**:
   ```go
   case rules.RuleTypeCustom:
       return f.createCustomRule(...)
   ```

## Related Documentation

- [Domain Rules](../../../domain/rules/README.md) - Domain rules system
- [AWS Models](../models/README.md) - AWS resource models
- [AWS Adapters](../adapters/README.md) - AWS adapters
