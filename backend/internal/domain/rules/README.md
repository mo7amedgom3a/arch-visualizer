# Domain Rules System

This package provides a **cloud-agnostic rules engine** that allows different cloud providers to implement validation rules in their own way while maintaining a consistent interface.

## Purpose

The domain rules system enables:
- **Multi-Cloud Support**: Same rule interfaces work across AWS, GCP, Azure
- **Provider-Specific Implementation**: Each provider implements rules according to their constraints
- **Data-Driven Validation**: Rules can be loaded from database without code changes
- **Extensibility**: Easy to add new rule types and providers

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│              Domain Rules (This Package)                     │
│  - Rule interfaces                                          │
│  - Constraint types                                         │
│  - Rule engine & evaluator                                  │
│  - Rule registry                                            │
└──────────────────────┬──────────────────────────────────────┘
                        │
                        │ implemented by
                        ▼
┌─────────────────────────────────────────────────────────────┐
│          Cloud Provider Rule Implementations                │
│  - AWS: internal/cloud/aws/rules/                           │
│  - GCP: internal/cloud/gcp/rules/ (future)                 │
│  - Azure: internal/cloud/azure/rules/ (future)              │
└─────────────────────────────────────────────────────────────┘
```

## Structure

```
rules/
├── types.go                    # Core rule interfaces and types
├── constraints/                # Constraint implementations
│   ├── parent.go              # Parent/containment rules
│   ├── region.go              # Region requirements
│   ├── limits.go              # Children limits
│   └── dependency.go          # Dependency rules
├── engine/                     # Rule evaluation engine
│   ├── context.go             # Evaluation context builder
│   └── evaluator.go           # Rule evaluator
├── registry/                   # Rule registry
│   └── registry.go            # Rule storage and factory
└── README.md                   # This file
```

## Core Concepts

### Rule Interface

All rules implement the `Rule` interface:

```go
type Rule interface {
    GetType() RuleType
    GetResourceType() string
    GetValue() string
    Evaluate(ctx context.Context, evalCtx *EvaluationContext) error
}
```

### Rule Types

- `requires_parent` - Resource must have a parent
- `allowed_parent` - Only specific parent types allowed
- `requires_region` - Resource must (or must not) have a region
- `max_children` - Maximum number of children
- `min_children` - Minimum number of children
- `allowed_dependencies` - Allowed dependency types
- `forbidden_dependencies` - Forbidden dependency types

### Evaluation Context

Provides context for rule evaluation:

```go
type EvaluationContext struct {
    Resource      *resource.Resource
    Parents       []*resource.Resource
    Children      []*resource.Resource
    Dependencies  []*resource.Resource
    Architecture  *Architecture
    Provider      string
    Metadata      map[string]interface{}
}
```

## Usage Example

### Creating Rules

```go
import (
    "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/constraints"
)

// Create a requires_parent rule
rule := constraints.NewRequiresParentRule("Subnet", "VPC")

// Create an allowed_parent rule
rule := constraints.NewAllowedParentRule("EC2", []string{"Subnet"})

// Create a max_children rule
rule := constraints.NewMaxChildrenRule("VPC", 10)
```

### Evaluating Rules

```go
import (
    "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/engine"
    "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/registry"
)

// Create registry and register rules
registry := registry.NewRuleRegistry()
registry.RegisterRule("Subnet", rule)

// Create evaluator
evaluator := engine.NewRuleEvaluator()

// Build evaluation context
evalCtx := engine.BuildEvaluationContext(resource, architecture, "aws")

// Evaluate rules
rules := registry.GetRules("Subnet")
result := engine.EvaluateAllRules(ctx, evaluator, rules, evalCtx)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Println(err.Message)
    }
}
```

## Cloud Provider Implementation

Each cloud provider implements rules using their own factory:

### AWS Implementation

```go
// internal/cloud/aws/rules/factory.go
factory := rules.NewAWSRuleFactory()
rule, err := factory.CreateRule("Subnet", "requires_parent", "VPC")
```

AWS can:
- Map resource types to AWS-specific names
- Apply AWS-specific limits
- Add AWS-specific validation logic

### GCP Implementation (Future)

```go
// internal/cloud/gcp/rules/factory.go
factory := rules.NewGCPRuleFactory()
rule, err := factory.CreateRule("Subnet", "requires_parent", "VPC")
```

GCP can:
- Map to GCP resource types (e.g., "VPC" → "google_compute_network")
- Apply GCP-specific limits
- Add GCP-specific validation

## Rule Registry

The registry stores and retrieves rules:

```go
registry := registry.NewRuleRegistry()

// Register a rule
registry.RegisterRule("Subnet", rule)

// Get all rules for a resource type
rules := registry.GetRules("Subnet")

// Get all rules of a specific type
rules := registry.GetRulesByType(rules.RuleTypeRequiresParent)
```

## Rule Factory

Factories create rules from database constraint records:

```go
factory := registry.NewRuleFactory()
rule, err := factory.CreateRule(
    "Subnet",              // resource type
    "requires_parent",     // constraint type
    "VPC",                 // constraint value
)
```

Cloud providers can implement their own factories to customize rule creation.

## Database Integration

Rules are typically loaded from the `resource_constraints` table:

```sql
SELECT resource_type, constraint_type, constraint_value
FROM resource_constraints
WHERE cloud_provider = 'aws';
```

Then converted to rules:

```go
for _, constraint := range constraints {
    rule, err := factory.CreateRule(
        constraint.ResourceType,
        constraint.ConstraintType,
        constraint.ConstraintValue,
    )
    registry.RegisterRule(constraint.ResourceType, rule)
}
```

## Adding New Rule Types

1. **Add RuleType constant** in `types.go`:
   ```go
   RuleTypeNewRule RuleType = "new_rule"
   ```

2. **Create constraint implementation** in `constraints/`:
   ```go
   type NewRule struct {
       ResourceType string
       // ... fields
   }
   
   func (r *NewRule) Evaluate(...) error {
       // ... validation logic
   }
   ```

3. **Add to factory** in `registry/registry.go`:
   ```go
   case rules.RuleTypeNewRule:
       return constraints.NewNewRule(...), nil
   ```

4. **Update cloud provider factories** if needed

## Benefits

1. **Cloud-Agnostic**: Same rule interfaces across all providers
2. **Provider Flexibility**: Each provider implements rules their way
3. **Data-Driven**: Rules loaded from database, no code changes needed
4. **Extensible**: Easy to add new rule types
5. **Testable**: Rules can be tested independently
6. **Type-Safe**: Compile-time checks ensure correctness

## Related Documentation

- [AWS Rules Implementation](../../cloud/aws/rules/README.md) - AWS-specific rule implementation
- [Resource Domain Models](../resource/README.md) - Domain resource models
- [Architecture Validation](../architecture/validation.go) - Architecture-level validation
