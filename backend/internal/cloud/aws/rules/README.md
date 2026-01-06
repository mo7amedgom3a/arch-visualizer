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

## Usage Example

```go
import (
    awsrules "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/rules"
)

// Create AWS rule service
service := awsrules.NewAWSRuleService()

// Load AWS-specific rules from database
constraints := loadConstraintsFromDB("aws")
service.LoadRulesFromConstraints(ctx, constraints)

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

## AWS-Specific Considerations

### Resource Type Mapping

AWS uses Terraform-style resource names:
- Domain: `VPC` → AWS: `aws_vpc`
- Domain: `Subnet` → AWS: `aws_subnet`
- Domain: `EC2Instance` → AWS: `aws_instance`

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
