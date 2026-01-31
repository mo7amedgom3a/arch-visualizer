# Scenario 8: Architecture Roundtrip

This scenario demonstrates a complete roundtrip: **Save → Load → Convert to JSON**. It reads a diagram JSON file, processes and saves it to the database, loads it back, and converts it back to the original diagram JSON format.

## Overview

This use case tests the complete persistence and retrieval cycle:

1. **Read** diagram JSON from `json-request-fiagram-complete.json`
2. **Process** diagram through service layer (parse, validate, map, persist)
3. **Load** architecture from database
4. **Convert** architecture back to diagram JSON format
5. **Save** response JSON to `json-response-architecture-loaded.json`

## What This Scenario Demonstrates

- **Complete Persistence Cycle**: Save architecture to database and load it back
- **Data Integrity**: Verify that all data (resources, containments, dependencies, variables, outputs) is preserved
- **JSON Conversion**: Convert domain Architecture back to diagram JSON format
- **Roundtrip Testing**: Ensure data can be saved and retrieved without loss

## Workflow

```
┌─────────────────────────────────────┐
│  json-request-fiagram-complete.json │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  ProcessDiagram (Service Layer)     │
│  - Parse                            │
│  - Validate                         │
│  - Map to Architecture              │
│  - Persist to Database              │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Database (PostgreSQL)              │
│  - Project                          │
│  - Resources                        │
│  - Containments                     │
│  - Dependencies                     │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  LoadArchitecture (Service Layer)   │
│  - Load Project                     │
│  - Load Resources                   │
│  - Load Containments                │
│  - Load Dependencies                │
│  - Reconstruct Architecture         │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Convert Architecture to JSON       │
│  - Resources → Nodes                │
│  - Dependencies → Edges             │
│  - Variables & Outputs              │
│  - Format as Diagram JSON           │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  json-response-architecture-loaded  │
│  .json                              │
└─────────────────────────────────────┘
```

## Architecture to JSON Conversion

The `convertArchitectureToDiagramJSON` function converts a domain `Architecture` back to the diagram JSON format:

### Nodes

- **Region Node**: Created from architecture region and provider
- **Resource Nodes**: Created from architecture resources with:
  - Position from metadata
  - Resource type from domain resource type
  - Config from metadata
  - Parent relationships from containments
  - Visual-only flag from metadata

### Edges

- Created from architecture dependencies
- Each dependency becomes an edge from source to target

### Variables & Outputs

- Directly copied from architecture variables and outputs

## Running the Scenario

```bash
# Run scenario 8
go run ./cmd/api/main.go -scenario=8

# Or build and run
go build ./cmd/api/main.go
./api -scenario=8
```

## Expected Output

```
====================================================================================================
SCENARIO 8: Architecture Roundtrip (Save → Load → Convert to JSON)
====================================================================================================

[Step 1] Initializing service layer server...
✓ Service layer server initialized successfully

[Step 2] Reading diagram JSON file...
✓ Read diagram JSON from: .../json-request-fiagram-complete.json (XXXX bytes)

[Step 3] Processing diagram and saving to database...
✓ Diagram processed and saved to database
  Project ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

[Step 4] Loading architecture from database...
✓ Architecture loaded from database
  Resources: X
  Containments: X
  Dependencies: X
  Variables: X
  Outputs: X

[Step 5] Converting architecture to diagram JSON format...
✓ Architecture converted to diagram JSON format

[Step 6] Saving response JSON to file...
✓ Response JSON saved to: .../json-response-architecture-loaded.json (XXXX bytes)

====================================================================================================
SUCCESS: Architecture roundtrip completed!
  Project ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  Resources Loaded: X
  Response JSON: .../json-response-architecture-loaded.json
====================================================================================================
```

## Output Files

### Input File
- `json-request-fiagram-complete.json` - Original diagram JSON

### Output File
- `json-response-architecture-loaded.json` - Architecture loaded from database and converted back to JSON format

## JSON Format

The response JSON follows the same structure as the input:

```json
{
  "cloud-canvas-project-{project-id}": {
    "nodes": [
      {
        "id": "var.aws_region",
        "type": "containerNode",
        "position": { "x": 440, "y": 40 },
        "data": {
          "label": "Region",
          "resourceType": "region",
          "config": { ... },
          "status": "valid",
          "isVisualOnly": false
        },
        ...
      },
      {
        "id": "{resource-id}",
        "type": "containerNode" | "resourceNode",
        "parentId": "{parent-id}",
        "position": { "x": 20, "y": 40 },
        "data": {
          "label": "{resource-name}",
          "resourceType": "{resource-type}",
          "config": { ... },
          "status": "valid",
          "isVisualOnly": false
        },
        ...
      }
    ],
    "edges": [
      {
        "id": "{edge-id}",
        "source": "{from-resource-id}",
        "target": "{to-resource-id}",
        "type": "default"
      }
    ],
    "variables": [
      {
        "name": "{variable-name}",
        "type": "{variable-type}",
        "description": "{description}",
        "default": {default-value},
        "sensitive": false
      }
    ],
    "outputs": [
      {
        "name": "{output-name}",
        "value": "{output-value}",
        "description": "{description}",
        "sensitive": false
      }
    ],
    "timestamp": {timestamp}
  }
}
```

## Data Preservation

The roundtrip preserves:

- ✅ **Resources**: All resources with their configurations
- ✅ **Positions**: X/Y coordinates from metadata
- ✅ **Containments**: Parent-child relationships
- ✅ **Dependencies**: Resource dependencies
- ✅ **Variables**: Terraform input variables
- ✅ **Outputs**: Terraform output values
- ✅ **Visual-only Flag**: For visual-only nodes

## Limitations

1. **Resource IDs**: When loading from database, resource IDs are UUIDs. The original node IDs (like "vpc-2", "subnet-4") are not preserved unless stored in metadata.

2. **Region Node Config**: The region node configuration (list of available regions) is hardcoded in the converter. In a full implementation, this would be loaded from configuration.

3. **Node Types**: Node type (containerNode vs resourceNode) is inferred from resource type category/kind. This may not always match the original exactly.

## Future Enhancements

- [ ] Store original node IDs in metadata to preserve exact IDs
- [ ] Load region configuration from database or config
- [ ] Preserve exact node type from original diagram
- [ ] Add validation to compare input and output JSON
- [ ] Support for partial updates (update existing project)

## Related Documentation

- [Service Layer README](../../../internal/platform/server/README.md) - Service layer documentation
- [Scenario 7 README](../scenario7_service_layer/README.md) - Service layer usage
- [Use Cases README](../README.md) - All use cases overview
