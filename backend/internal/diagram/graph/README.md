# Diagram Graph

The graph module provides the normalized internal representation of a diagram. It offers efficient traversal, query operations, and relationship management for nodes and edges.

## Overview

`DiagramGraph` is the core data structure after parsing and normalization. It provides:
- **Node Management**: O(1) node lookups by ID
- **Relationship Queries**: Parent/child, containment, dependencies
- **Tree Operations**: Root nodes, containment trees
- **Graph Traversal**: Efficient navigation of the structure

## Data Structures

### DiagramGraph

Main graph structure:

```go
type DiagramGraph struct {
    Nodes map[string]*Node  // Node ID -> Node
    Edges []*Edge           // All edges in the graph
}
```

### Node

Represents a normalized node:

```go
type Node struct {
    ID           string                 // Unique identifier
    Type         string                 // "containerNode" or "resourceNode"
    ResourceType string                 // e.g., "vpc", "ec2", "region"
    Label        string                 // Display label
    Config       map[string]interface{} // Configuration
    PositionX    int                    // X coordinate
    PositionY    int                    // Y coordinate
    ParentID     *string                // Parent node ID (containment)
    Status       string                 // Node status
    IsVisualOnly bool                   // Visual-only flag
}
```

### Edge

Represents a relationship:

```go
type Edge struct {
    Source string  // Source node ID
    Target string  // Target node ID
    Type   string  // "containment", "dependency", "reference"
}
```

## Graph Operations

### Node Lookup

```go
node, exists := graph.GetNode("node-id")
```

Returns node by ID with existence check.

### Parent/Child Queries

```go
// Get all children of a node
children := graph.GetChildren("parent-id")

// Get parent of a node
parent, exists := graph.GetParent("child-id")
```

### Root Nodes

```go
roots := graph.GetRootNodes()
```

Returns all nodes without a parent (top-level nodes).

### Containment Tree

```go
tree := graph.BuildContainmentTree()
// Returns: map[string][]*Node (parent ID -> children)
```

Builds a complete containment hierarchy.

### Edge Filtering

```go
// Get only containment edges
containmentEdges := graph.GetContainmentEdges()

// Get only dependency edges
dependencyEdges := graph.GetDependencyEdges()
```

### Region Node

```go
regionNode, found := graph.FindRegionNode()
```

Finds the region node (if present).

## Node Helper Methods

### Type Checks

```go
node.IsContainer()  // true if type == "containerNode"
node.IsResource()   // true if type == "resourceNode"
node.IsRegion()     // true if resourceType == "region"
```

## Edge Helper Methods

### Type Checks

```go
edge.IsContainment()  // true if type == "containment"
edge.IsDependency()   // true if type == "dependency"
```

## Usage Examples

### Traverse Containment Tree

```go
// Get root nodes
roots := graph.GetRootNodes()

// Recursively traverse children
var traverse func(nodeID string, depth int)
traverse = func(nodeID string, depth int) {
    node, _ := graph.GetNode(nodeID)
    indent := strings.Repeat("  ", depth)
    fmt.Printf("%s%s (%s)\n", indent, node.Label, node.ResourceType)
    
    children := graph.GetChildren(nodeID)
    for _, child := range children {
        traverse(child.ID, depth+1)
    }
}

for _, root := range roots {
    traverse(root.ID, 0)
}
```

### Find All Dependencies

```go
dependencies := graph.GetDependencyEdges()
for _, edge := range dependencies {
    source, _ := graph.GetNode(edge.Source)
    target, _ := graph.GetNode(edge.Target)
    fmt.Printf("%s depends on %s\n", source.Label, target.Label)
}
```

### Check Node Hierarchy

```go
// Check if node is in a VPC
node, _ := graph.GetNode("ec2-1")
current := node
for current.ParentID != nil {
    parent, _ := graph.GetParent(current.ID)
    if parent.ResourceType == "vpc" {
        fmt.Printf("Node is in VPC: %s\n", parent.Label)
        break
    }
    current = parent
}
```

## Graph Structure

### Containment Hierarchy

```
Region (root)
└── VPC
    ├── Subnet
    │   └── EC2
    └── RouteTable
```

### Dependency Graph

```
EC2 ──depends_on──> SecurityGroup
EC2 ──depends_on──> Subnet
Subnet ──depends_on──> VPC
```

## Performance Characteristics

- **Node Lookup**: O(1) via map
- **GetChildren**: O(n) where n = total nodes
- **GetParent**: O(1) via direct lookup
- **BuildContainmentTree**: O(n)
- **GetRootNodes**: O(n)

## Common Patterns

### Find All Resources of Type

```go
func findResourcesByType(graph *DiagramGraph, resourceType string) []*Node {
    var results []*Node
    for _, node := range graph.Nodes {
        if node.ResourceType == resourceType {
            results = append(results, node)
        }
    }
    return results
}
```

### Get Ancestors

```go
func getAncestors(graph *DiagramGraph, nodeID string) []*Node {
    var ancestors []*Node
    node, exists := graph.GetNode(nodeID)
    if !exists {
        return ancestors
    }
    
    current := node
    for current.ParentID != nil {
        parent, exists := graph.GetParent(current.ID)
        if !exists {
            break
        }
        ancestors = append(ancestors, parent)
        current = parent
    }
    return ancestors
}
```

### Get Descendants

```go
func getDescendants(graph *DiagramGraph, nodeID string) []*Node {
    var descendants []*Node
    children := graph.GetChildren(nodeID)
    
    for _, child := range children {
        descendants = append(descendants, child)
        descendants = append(descendants, getDescendants(graph, child.ID)...)
    }
    
    return descendants
}
```

## Integration

The graph is created by the parser:

```go
// In parser
diagramGraph, err := parser.NormalizeToGraph(irDiagram)
```

Used by the validator:

```go
// In validator
result := validator.Validate(diagramGraph, opts)
```

Mapped to domain models:

```go
// In domain/architecture
arch := architecture.MapDiagramToArchitecture(diagramGraph, provider)
```

## Testing

Run graph tests:

```bash
go test ./internal/diagram/graph/... -v
```

## Best Practices

1. **Use helper methods**: Prefer `GetChildren()` over manual iteration
2. **Check existence**: Always check `exists` from `GetNode()`
3. **Handle nil ParentID**: Root nodes have `nil` parent
4. **Use appropriate queries**: Use `GetContainmentEdges()` vs `GetDependencyEdges()`

## Related

- [Diagram Module README](../README.md)
- [Parser README](../parser/README.md)
- [Validator README](../validator/README.md)
