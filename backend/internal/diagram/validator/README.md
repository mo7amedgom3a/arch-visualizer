# Diagram Validator

The validator ensures structural correctness of diagram graphs before they are processed and persisted. It performs comprehensive checks on containment relationships, edge references, resource types, and region nodes.

## Overview

The validator operates on a `DiagramGraph` and returns a `ValidationResult` containing:
- **Valid**: Boolean indicating if the graph passed all error checks
- **Errors**: Critical issues that prevent processing
- **Warnings**: Non-critical issues that should be reviewed

## Validation Checks

### 1. Missing Parent References

**Error Code**: `MISSING_PARENT`

Checks that all `parentId` references point to existing nodes.

```go
// ❌ Invalid: Node references non-existent parent
{
  "id": "subnet-1",
  "parentId": "vpc-999"  // vpc-999 doesn't exist
}
```

### 2. Containment Cycles

**Error Code**: `CONTAINMENT_CYCLE`

Detects cycles in the containment tree using DFS with recursion stack tracking.

```go
// ❌ Invalid: Circular containment
node1.parentId = "node2"
node2.parentId = "node1"  // Cycle!
```

### 3. Invalid Edge References

**Error Codes**: `INVALID_EDGE_SOURCE`, `INVALID_EDGE_TARGET`

Ensures all edges reference existing nodes.

```go
// ❌ Invalid: Edge references non-existent node
{
  "source": "ec2-1",
  "target": "sg-999"  // sg-999 doesn't exist
}
```

### 4. Resource Type Validation

**Error Code**: `MISSING_RESOURCE_TYPE`  
**Warning Code**: `UNKNOWN_RESOURCE_TYPE`

Validates resource types against database entries for the specified cloud provider.

```go
// ❌ Error: Missing resource type
{
  "id": "node-1",
  "resourceType": ""  // Empty!
}

// ⚠️ Warning: Unknown resource type
{
  "id": "node-2",
  "resourceType": "unknown-type"  // Not in database
}
```

### 5. Region Node Validation

**Warning Codes**: `NO_REGION_NODE`, `MULTIPLE_REGION_NODES`

Checks for the presence and uniqueness of region nodes.

```go
// ⚠️ Warning: No region found
// ⚠️ Warning: Multiple regions found (only one will be used)
```

## Usage

### Basic Validation

```go
import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator"

// Validate without resource type checking (backward compatible)
result := validator.Validate(diagramGraph, nil)

if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Error: %s\n", err.Message)
    }
}
```

### Validation with Resource Types

```go
// Load valid resource types from database
validTypes := map[string]bool{
    "vpc":         true,
    "subnet":      true,
    "ec2":         true,
    "route-table": true,
    // ... more types
}

opts := &validator.ValidationOptions{
    ValidResourceTypes: validTypes,
    Provider:           "aws",
}

result := validator.Validate(diagramGraph, opts)
```

## ValidationResult Structure

```go
type ValidationResult struct {
    Valid    bool              // true if no errors
    Errors   []*ValidationError // Critical issues
    Warnings []*ValidationError // Non-critical issues
}
```

## ValidationError Structure

```go
type ValidationError struct {
    Code    string // Error code (e.g., "MISSING_PARENT")
    Message string // Human-readable message
    NodeID  string // Node ID (if applicable)
}
```

## Error Codes Reference

| Code | Severity | Description |
|------|----------|-------------|
| `MISSING_PARENT` | Error | Node references non-existent parent |
| `CONTAINMENT_CYCLE` | Error | Cycle detected in containment tree |
| `INVALID_EDGE_SOURCE` | Error | Edge source node doesn't exist |
| `INVALID_EDGE_TARGET` | Error | Edge target node doesn't exist |
| `MISSING_RESOURCE_TYPE` | Error | Node has no resource type |
| `UNKNOWN_RESOURCE_TYPE` | Warning | Resource type not found in database |
| `NO_REGION_NODE` | Warning | No region node found |
| `MULTIPLE_REGION_NODES` | Warning | Multiple region nodes found |

## Example: Handling Validation Results

```go
result := validator.Validate(diagramGraph, opts)

if !result.Valid {
    fmt.Println("Validation failed with errors:")
    for _, err := range result.Errors {
        fmt.Printf("  [%s] %s", err.Code, err.Message)
        if err.NodeID != "" {
            fmt.Printf(" (node: %s)", err.NodeID)
        }
        fmt.Println()
    }
    return fmt.Errorf("validation failed")
}

if len(result.Warnings) > 0 {
    fmt.Println("Validation passed with warnings:")
    for _, warn := range result.Warnings {
        fmt.Printf("  [%s] %s\n", warn.Code, warn.Message)
    }
}
```

## Algorithm: Cycle Detection

The containment cycle detection uses DFS with recursion stack:

```go
1. For each unvisited node:
   a. Mark as visited
   b. Add to recursion stack
   c. Recursively check all children
   d. If child is in recursion stack → cycle found
   e. Remove from recursion stack
```

This ensures we detect cycles even if they don't include root nodes.

## Integration with Service

The validator is typically called from the diagram service:

```go
// In diagram/service.go
validResourceTypes, err := s.buildValidResourceTypesMap(ctx, provider)
if err != nil {
    return uuid.Nil, err
}

opts := &validator.ValidationOptions{
    ValidResourceTypes: validResourceTypes,
    Provider:           string(provider),
}

result := validator.Validate(diagramGraph, opts)
if !result.Valid {
    return uuid.Nil, fmt.Errorf("validation failed: %v", result.Errors)
}
```

## Testing

Run validator tests:

```bash
go test ./internal/diagram/validator/... -v
```

Test coverage includes:
- ✅ Valid graphs
- ✅ Missing parent references
- ✅ Containment cycles
- ✅ Unknown resource types
- ✅ Edge reference validation

## Best Practices

1. **Always validate before processing**: Never skip validation
2. **Handle warnings appropriately**: Log warnings but don't block processing
3. **Use provider-specific validation**: Load resource types from database
4. **Provide context in errors**: Include node IDs for easier debugging

## Related

- [Diagram Module README](../README.md)
- [Graph Module README](../graph/README.md)
