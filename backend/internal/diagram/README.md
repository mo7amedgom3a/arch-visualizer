# Diagram Processing Module

The diagram processing module is the "first brain" of the architecture visualizer. It transforms frontend visual diagrams (IR JSON) into validated, normalized graph structures that can be mapped to cloud-agnostic domain models and persisted to the database.

## Overview

This module implements a multi-stage pipeline:

```
Frontend IR JSON → Parse → Normalize → Validate → Map to Domain → Persist to DB
```

## Architecture

The module is organized into four main components:

### 1. **Parser** (`parser/`)
- Parses IR (Intermediate Representation) JSON from the frontend
- Handles various JSON formats (direct objects, strings, arrays)
- Extracts nodes, edges, positions, and metadata
- Filters visual-only nodes

### 2. **Graph** (`graph/`)
- Represents the normalized diagram structure
- Provides graph traversal and query operations
- Manages containment trees and dependency graphs

### 3. **Validator** (`validator/`)
- Validates structural correctness of the diagram
- Checks for cycles, missing references, invalid types
- Provides detailed error and warning messages

### 4. **Service** (`service.go`)
- Orchestrates the entire processing pipeline
- Loads valid resource types from database
- Maps validated graphs to domain models
- Persists projects and resources to database

## Processing Pipeline

### Step 1: Parse IR Diagram

```go
irDiagram, err := parser.ParseIRDiagram(jsonData)
```

Converts frontend JSON into structured `IRDiagram` with:
- **Nodes**: Container and resource nodes with metadata
- **Edges**: Relationships between nodes
- **Variables**: Frontend-specific variables
- **Timestamp**: Diagram creation timestamp

### Step 2: Normalize to Graph

```go
diagramGraph, err := parser.NormalizeToGraph(irDiagram)
```

Transforms IR into normalized `DiagramGraph`:
- Filters out visual-only nodes (if needed)
- Extracts positions (x, y coordinates)
- Implicitly creates containment edges from `parentId`
- Normalizes resource types and labels

### Step 3: Validate Graph

```go
validResourceTypes, _ := service.buildValidResourceTypesMap(ctx, provider)
validationOpts := &validator.ValidationOptions{
    ValidResourceTypes: validResourceTypes,
    Provider:           string(provider),
}
validationResult := validator.Validate(diagramGraph, validationOpts)
```

Validates:
- ✅ No missing parent references
- ✅ No containment cycles
- ✅ Valid edge references
- ✅ Valid resource types (from database)
- ✅ Region node presence

### Step 4: Map to Domain Architecture

```go
domainArch, err := architecture.MapDiagramToArchitecture(diagramGraph, provider)
```

Converts graph to cloud-agnostic domain model:
- **Resources**: Domain resource objects
- **Containments**: Parent-child relationships
- **Dependencies**: Logical dependencies

### Step 5: Persist to Database

```go
projectID, err := service.persistArchitecture(ctx, domainArch, userID, projectName, iacToolID, provider, region)
```

Creates database records:
- Project
- Resources (with positions and `isVisualOnly` flag)
- Resource containments
- Resource dependencies

## Usage Example

```go
// Create service
service, err := diagram.NewService()
if err != nil {
    log.Fatal(err)
}

// Process diagram
jsonData := []byte(`{"nodes": [...], "edges": [...]}`)
projectID, err := service.ProcessDiagramRequest(
    ctx,
    jsonData,
    userID,
    "My Project",
    iacToolID,
)
```

## Data Flow

```
┌─────────────────┐
│  Frontend IR    │
│     JSON        │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Parser        │ → IRDiagram
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Normalizer    │ → DiagramGraph
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Validator     │ → ValidationResult
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Domain Mapper   │ → Architecture
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Persistence   │ → Database
└─────────────────┘
```

## Key Concepts

### Node Types

- **Container Nodes**: Scope other resources (Region, VPC, Subnet, Security Group)
- **Resource Nodes**: Actual infrastructure resources (EC2, Lambda, S3)

### Visual-Only Flag

Nodes with `isVisualOnly: true` are:
- Tracked in the graph for reference
- Filtered out when creating domain resources
- Stored in database with the flag set

### Containment vs Dependencies

- **Containment**: Hierarchical relationships (VPC contains Subnet)
- **Dependencies**: Logical relationships (EC2 depends on Security Group)

### Position Tracking

Node positions (`x`, `y` coordinates) are:
- Extracted from IR JSON `position` object
- Stored in graph nodes
- Persisted to database `pos_x` and `pos_y` columns

## Error Handling

The module provides detailed error messages:
- **Validation Errors**: Block processing (missing parents, cycles)
- **Warnings**: Informational (unknown resource types, multiple regions)

## Testing

Run tests for each component:

```bash
go test ./internal/diagram/parser/...
go test ./internal/diagram/validator/...
go test ./internal/diagram/...
```

## Related Documentation

- [Parser README](parser/README.md) - Detailed parser documentation
- [Graph README](graph/README.md) - Graph structure and operations
- [Validator README](validator/README.md) - Validation rules and error codes
