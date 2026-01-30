# Diagram Parser

The parser transforms frontend IR (Intermediate Representation) JSON into normalized graph structures. It handles various JSON formats, extracts node metadata, and filters visual-only nodes.

## Overview

The parser performs two main operations:

1. **Parse**: Converts JSON bytes into structured `IRDiagram`
2. **Normalize**: Transforms `IRDiagram` into `DiagramGraph`

## IR JSON Format

The frontend sends diagrams in IR format:

```json
{
  "nodes": [
    {
      "id": "region-1",
      "type": "containerNode",
      "position": { "x": 100, "y": 200 },
      "data": {
        "label": "US East 1",
        "resourceType": "region",
        "config": { "name": "us-east-1" },
        "isVisualOnly": false
      }
    },
    {
      "id": "vpc-1",
      "type": "containerNode",
      "parentId": "region-1",
      "data": {
        "label": "Project VPC",
        "resourceType": "vpc",
        "config": { "name": "project-vpc" }
      }
    }
  ],
  "edges": [
    {
      "source": "ec2-1",
      "target": "sg-1",
      "type": "dependency"
    }
  ],
  "variables": [],
  "timestamp": 1234567890
}
```

## Data Structures

### IRDiagram

Top-level structure representing the parsed JSON:

```go
type IRDiagram struct {
    Nodes     []IRNode      // All nodes in the diagram
    Edges     []IREdge      // Relationships between nodes
    Variables []interface{} // Frontend variables
    Timestamp int64         // Creation timestamp
}
```

### IRNode

Represents a single node:

```go
type IRNode struct {
    ID       string                 // Unique identifier
    Type     string                 // "containerNode" or "resourceNode"
    Position *IRPosition            // X, Y coordinates
    ParentID *string                // Parent node ID (for containment)
    Data     IRNodeData             // Node payload
    Style    map[string]interface{} // Visual styling (ignored)
}
```

### IRNodeData

Node data payload:

```go
type IRNodeData struct {
    Label        string                 // Display label
    ResourceType string                 // Resource type (e.g., "vpc", "ec2")
    Config       map[string]interface{} // Configuration object
    Status       string                 // Node status
    IsVisualOnly *bool                  // Visual-only flag
}
```

### IREdge

Represents a relationship:

```go
type IREdge struct {
    ID     string  // Edge identifier
    Source string  // Source node ID
    Target string  // Target node ID
    Type   *string // "containment", "dependency", "reference"
}
```

## Functions

### ParseIRDiagram

Parses JSON bytes into `IRDiagram`, handling multiple formats:

```go
irDiagram, err := parser.ParseIRDiagram(jsonData)
```

**Supported Formats**:
1. Direct object: `{"nodes": [...], "edges": [...]}`
2. JSON string: `"{\"nodes\": [...]}"`
3. Array of nodes: `[{...}, {...}]`

**Error Handling**:
- Returns error if JSON is malformed
- Handles missing fields gracefully
- Validates structure before returning

### NormalizeToGraph

Transforms `IRDiagram` into normalized `DiagramGraph`:

```go
diagramGraph, err := parser.NormalizeToGraph(irDiagram)
```

**Normalization Steps**:

1. **Filter Visual-Only Nodes**: 
   - Nodes with `isVisualOnly: true` are kept in graph for reference
   - They are filtered out later when creating domain resources

2. **Extract Positions**:
   - Extracts `x`, `y` from `position` object
   - Defaults to `0, 0` if missing

3. **Create Containment Edges**:
   - Implicitly creates edges from `parentId` relationships
   - Adds to graph's edge list

4. **Normalize Resource Types**:
   - Preserves resource type as-is
   - Converts to lowercase for consistency

5. **Build Node Map**:
   - Creates `map[string]*Node` for O(1) lookups

## Node Types

### Container Nodes

Nodes that scope other resources:

- `type: "containerNode"`
- Examples: Region, VPC, Subnet, Security Group
- Can have children via `parentId`

### Resource Nodes

Actual infrastructure resources:

- `type: "resourceNode"`
- Examples: EC2, Lambda, S3, Route Table
- Become real Terraform/Pulumi resources

## Visual-Only Flag

The `isVisualOnly` flag indicates whether a node is:
- **`false`**: Real infrastructure (processed and persisted)
- **`true`**: Visual icon only (tracked but not persisted as infrastructure)

**Behavior**:
- Visual-only nodes are kept in the graph
- They are filtered out when creating domain resources
- The flag is stored in database for reference

## Position Extraction

Positions are extracted from the `position` object:

```json
{
  "position": {
    "x": 150.5,
    "y": 200.3
  }
}
```

Converted to integers and stored in `Node.PositionX` and `Node.PositionY`.

## Implicit Containment Edges

The parser automatically creates containment edges from `parentId`:

```json
{
  "id": "subnet-1",
  "parentId": "vpc-1"  // Creates containment edge
}
```

This edge is added to `graph.Edges` with type `"containment"`.

## Error Handling

### Parse Errors

```go
irDiagram, err := parser.ParseIRDiagram(jsonData)
if err != nil {
    // Handle: malformed JSON, missing fields, etc.
}
```

### Normalization Errors

```go
graph, err := parser.NormalizeToGraph(irDiagram)
if err != nil {
    // Handle: invalid structure, missing required fields
}
```

## Example Usage

```go
import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/parser"

// Read JSON from file or HTTP request
jsonData := []byte(`{"nodes": [...], "edges": [...]}`)

// Parse
irDiagram, err := parser.ParseIRDiagram(jsonData)
if err != nil {
    return fmt.Errorf("parse failed: %w", err)
}

// Normalize
diagramGraph, err := parser.NormalizeToGraph(irDiagram)
if err != nil {
    return fmt.Errorf("normalize failed: %w", err)
}

// Use graph for validation and processing
```

## Testing

Run parser tests:

```bash
go test ./internal/diagram/parser/... -v
```

Test coverage includes:
- ✅ Direct JSON object parsing
- ✅ JSON string parsing
- ✅ Array format parsing
- ✅ Position extraction
- ✅ Visual-only filtering
- ✅ Containment edge creation

## Edge Cases

### Missing Position

If `position` is missing, defaults to `(0, 0)`.

### Missing ParentID

If `parentId` is `null` or missing, node is treated as root.

### Empty ResourceType

Empty `resourceType` is preserved (validation will catch it).

### Visual-Only Nodes

Visual-only nodes are kept in graph but filtered during domain mapping.

## Integration with Architecture Generation

The parser outputs a `DiagramGraph` which is then processed by cloud provider-specific architecture generators:

```
IR JSON
    ↓ [ParseIRDiagram]
IRDiagram
    ↓ [NormalizeToGraph]
DiagramGraph
    ↓ [Cloud Provider Architecture Generator]
Domain Architecture
```

**Key Points**:
- Parser is cloud-agnostic - outputs normalized graph
- Architecture generators (provider-specific) convert graph to domain architecture
- Each provider implements its own generator in `internal/cloud/{provider}/architecture/`

## Related

- [Diagram Module README](../README.md)
- [Graph Module README](../graph/README.md)
- [Validator README](../validator/README.md)
- [Domain Architecture](../../domain/architecture/README.md) - Architecture generation documentation
- [AWS Architecture](../../cloud/aws/architecture/README.md) - AWS architecture generator