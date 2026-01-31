# Scenario 7: Terraform Code Generation with Service Layer

This scenario demonstrates the complete end-to-end pipeline using the **service layer** (`internal/platform/server/`). It showcases how the service layer orchestrates all layers (diagram, domain, cloud, validation, IaC, and repository) to process diagrams, generate Terraform code, and persist to the database.

## Overview

This use case replaces the manual orchestration in Scenario 6 with the service layer's `PipelineOrchestrator`, demonstrating:

- **Service Layer Usage**: Using the centralized service layer for all operations
- **Dependency Injection**: How services are wired together
- **Orchestration**: Complete workflow orchestration through `PipelineOrchestrator`
- **Individual Services**: How to use individual services when needed
- **Clean Architecture**: Separation of concerns through service interfaces

## What This Scenario Demonstrates

1. **Service Layer Initialization**: Creating and wiring all services
2. **Pipeline Orchestration**: Using `PipelineOrchestrator.ProcessDiagram()` for complete workflow
3. **Individual Service Usage**: Using `DiagramService`, `ArchitectureService`, and `CodegenService` directly
4. **Terraform Generation**: Generating IaC code through the service layer
5. **Database Persistence**: Persisting architecture through `ProjectService`

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Use Case (This)                          │
│                                                             │
│  1. Initialize Service Layer Server                         │
│  2. Read Diagram JSON                                       │
│  3. Use PipelineOrchestrator.ProcessDiagram()               │
│  4. Use CodegenService.Generate()                           │
│  5. Write Terraform Files                                   │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│              Service Layer (Orchestration)                  │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │  PipelineOrchestrator                              │   │
│  │  - ProcessDiagram()                                │   │
│  │  - GenerateCode()                                  │   │
│  └────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │   Diagram    │  │ Architecture │  │   Codegen    │   │
│  │   Service   │  │   Service    │  │   Service    │   │
│  └──────────────┘  └──────────────┘  └──────────────┘   │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │  ProjectService                                    │   │
│  │  - Create()                                       │   │
│  │  - PersistArchitecture()                          │   │
│  └────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Workflow

### Step 1: Initialize Service Layer

```go
srv, err := server.NewServer()
```

This initializes:
- All repositories (Project, Resource, ResourceType, etc.)
- Repository adapters (for dependency injection)
- All services (Diagram, Architecture, Codegen, Project)
- Pipeline orchestrator (composed of all services)

### Step 2: Process Diagram

```go
result, err := srv.PipelineOrchestrator.ProcessDiagram(ctx, &ProcessDiagramRequest{
    JSONData:      diagramData,
    UserID:        userID,
    ProjectName:   "Service Layer Architecture Project",
    IACToolID:     1,
    CloudProvider: "aws",
    Region:        "us-east-1",
})
```

This single call orchestrates:
1. **Parse**: Diagram JSON → DiagramGraph
2. **Validate**: Diagram structure and schemas
3. **Map**: DiagramGraph → Architecture (using cloud provider generator)
4. **Validate Rules**: Architecture against domain rules and constraints
5. **Create Project**: Persist project to database
6. **Persist Architecture**: Save resources, containments, and dependencies

### Step 3: Generate Terraform Code

```go
// Use individual services for code generation
diagramGraph, err := srv.DiagramService.Parse(ctx, diagramData)
arch, err := srv.ArchitectureService.MapFromDiagram(ctx, diagramGraph, "aws")
output, err := srv.CodegenService.Generate(ctx, arch, "terraform")
```

### Step 4: Write Files

```go
writeTerraformOutput("terraform_output", output)
```

## Key Differences from Scenario 6

| Aspect | Scenario 6 | Scenario 7 |
|--------|-----------|------------|
| **Orchestration** | Manual (direct layer calls) | Service Layer (`PipelineOrchestrator`) |
| **Dependency Management** | Direct repository creation | Dependency injection via service layer |
| **Error Handling** | Manual at each step | Centralized in services |
| **Testing** | Hard to mock | Easy to mock (interface-based) |
| **Code Organization** | Procedural | Service-oriented |
| **Reusability** | Limited | High (services can be reused) |

## Benefits of Using Service Layer

1. **Separation of Concerns**: Business logic separated from infrastructure
2. **Testability**: Easy to mock services for unit testing
3. **Maintainability**: Changes isolated to specific services
4. **Reusability**: Services can be used across different use cases
5. **Consistency**: Centralized orchestration ensures consistent workflows
6. **Flexibility**: Can use orchestrator or individual services as needed

## Running the Scenario

```bash
# Run scenario 7
go run ./cmd/api/main.go -scenario=7

# Or build and run
go build ./cmd/api/main.go
./api -scenario=7
```

## Expected Output

```
====================================================================================================
SCENARIO 7: Terraform Code Generation with Service Layer
====================================================================================================

[Step 1] Initializing service layer server...
✓ Service layer server initialized successfully

[Step 2] Reading diagram JSON file...
✓ Read diagram JSON from: .../json-request-fiagram-complete.json (XXXX bytes)

[Step 3] Processing diagram through service layer orchestrator...
✓ Diagram processed successfully
  Project ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  Message: Diagram processed successfully. Project created with ID: ...

[Step 4] Generating Terraform code...
✓ Project retrieved: Service Layer Architecture Project (Provider: aws, Region: us-east-1)

[Step 5] Demonstrating individual service usage...
✓ Diagram parsed using DiagramService
✓ Architecture mapped: X resources

[Step 6] Writing Terraform files to disk...
✓ Terraform files written to: ./terraform_output/
    - main.tf (XXXX bytes)
    - variables.tf (XXXX bytes)
    - outputs.tf (XXXX bytes)

====================================================================================================
SUCCESS: Complete pipeline executed using service layer!
  Project ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  Project Name: Service Layer Architecture Project
  Cloud Provider: aws
  Region: us-east-1
  Resources Persisted: X
  Terraform Files Generated: X
  Output Directory: ./terraform_output/
====================================================================================================
```

## Files Generated

The scenario generates Terraform files in `./terraform_output/`:

- `main.tf` - Main Terraform resources
- `variables.tf` - Input variables (if any)
- `outputs.tf` - Output values (if any)

## Database Persistence

The scenario persists:

- **Project**: Project metadata (name, cloud provider, region, IaC tool)
- **Resources**: All resources with configurations and positions
- **Containments**: Parent-child relationships
- **Dependencies**: Resource dependency relationships

## Code Structure

```
scenario7_service_layer/
├── terraform_with_service_layer.go  # Main use case implementation
└── README.md                        # This file
```

## Integration with Service Layer

This use case demonstrates the proper way to use the service layer:

1. **Initialize once**: Create server instance at the start
2. **Use orchestrator**: For complete workflows, use `PipelineOrchestrator`
3. **Use individual services**: For specific operations, use individual services
4. **Error handling**: Services return errors that can be handled appropriately

## Implemented Features

✅ **LoadArchitecture**: `ProjectService.LoadArchitecture()` reconstructs architecture from database
✅ **GenerateCode**: `PipelineOrchestrator.GenerateCode()` loads architecture and generates IaC code
✅ **Complete Workflow**: End-to-end pipeline from diagram to persisted project to generated code

## Future Enhancements

- [ ] Add support for updating existing projects
- [ ] Add support for variables and outputs persistence
- [ ] Add support for multiple cloud providers in rule validation
- [ ] Add support for Pulumi code generation
- [ ] Add support for batch operations
- [ ] Add caching layer for improved performance

## Related Documentation

- [Service Layer README](../../../internal/platform/server/README.md) - Complete service layer documentation
- [Scenario 6 README](../scenario6_terraform_with_persistence/README.md) - Manual orchestration approach
- [Use Cases README](../README.md) - All use cases overview
