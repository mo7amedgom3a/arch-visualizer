# Scenario 6: Terraform Code Generation with Database Persistence

This use case demonstrates the complete end-to-end pipeline from a diagram IR JSON to generated Terraform files, while also persisting the project and architecture to the database.

## Overview

This scenario extends Scenario 5 (Terraform Code Generation) by adding database persistence. It saves the following entities to the database:

- **Project**: The main project entity with cloud provider, region, and IaC tool information
- **Resources**: All infrastructure resources (VPC, Subnets, EC2, etc.) with their configurations
- **Resource Containments**: Parent-child relationships (e.g., VPC → Subnet → EC2)
- **Resource Dependencies**: Dependency relationships between resources

## Pipeline Steps

1. **Parse IR JSON** - Read and parse the diagram JSON file
2. **Normalize to Graph** - Convert to normalized diagram graph structure
3. **Validate Diagram** - Validate structure, schemas, and relationships
4. **Map to Architecture** - Convert to cloud-agnostic domain architecture
5. **Validate AWS Rules** - Apply AWS-specific networking rules
6. **Topological Sort** - Sort resources by dependencies
7. **Generate Terraform** - Produce HCL files using the Terraform engine
8. **Persist to Database** - Save project, resources, and relationships
9. **Write Files** - Output Terraform files to disk

## Database Entities

### Project
```go
type Project struct {
    ID            uuid.UUID
    UserID        uuid.UUID
    InfraToolID   uint        // References IACTarget (Terraform, Pulumi, etc.)
    Name          string
    CloudProvider string      // "aws", "azure", "gcp"
    Region        string
}
```

### Resource
```go
type Resource struct {
    ID             uuid.UUID
    ProjectID      uuid.UUID
    ResourceTypeID uint
    Name           string
    PosX, PosY     int         // Canvas position
    IsVisualOnly   bool
    Config         JSON        // Full resource configuration
}
```

### ResourceContainment
```go
type ResourceContainment struct {
    ParentResourceID uuid.UUID
    ChildResourceID  uuid.UUID
}
```

### ResourceDependency
```go
type ResourceDependency struct {
    FromResourceID   uuid.UUID
    ToResourceID     uuid.UUID
    DependencyTypeID uint
}
```

## Usage

### Run from Command Line

```bash
# From the backend directory
go run ./cmd/api/main.go -scenario=6
```

### Run Programmatically

```go
import (
    "context"
    "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario6_terraform_with_persistence"
)

func main() {
    ctx := context.Background()
    err := scenario6_terraform_with_persistence.TerraformWithPersistenceRunner(ctx)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Prerequisites

1. **Database Connection**: Ensure the database is configured and accessible
2. **Database Schema**: The required tables must exist:
   - `users`
   - `iac_targets`
   - `projects`
   - `resource_types`
   - `resources`
   - `resource_containment`
   - `resource_dependencies`
   - `dependency_types`

## Output

### Console Output
```
====================================================================================================
SCENARIO 6: Terraform Code Generation with Database Persistence
====================================================================================================
✓ Read diagram JSON from: /path/to/json-request-fiagram-complete.json
✓ Parsed IR diagram: 12 nodes, 5 variables, 6 outputs
✓ Normalized to graph: 12 nodes, 11 edges
✓ Diagram validation passed
✓ Mapped to architecture: 10 resources, region=us-east-1, provider=aws
✓ AWS rules validation passed
✓ Topologically sorted: 10 resources
✓ Generated Terraform: 3 files
  → Created demo user
  → Created IaC target: Terraform
  → Created project: Generated Architecture Project
  → Created 10 resources
  → Created 7 containment relationships
  → Created 0 dependency relationships
✓ Persisted to database: project_id=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
====================================================================================================
SUCCESS: Terraform code generated and project persisted!
  Project ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  Output directory: ./terraform_output/
    - main.tf
    - variables.tf
    - outputs.tf
====================================================================================================
```

### Generated Files
- `terraform_output/main.tf` - Main Terraform configuration
- `terraform_output/variables.tf` - Input variables
- `terraform_output/outputs.tf` - Output values

## Error Handling

The use case uses database transactions to ensure atomicity:
- If any step fails, the entire transaction is rolled back
- No partial data is left in the database on failure

## Comparison with Scenario 5

| Feature | Scenario 5 | Scenario 6 |
|---------|------------|------------|
| Parse IR JSON | ✓ | ✓ |
| Validate Diagram | ✓ | ✓ |
| Map to Architecture | ✓ | ✓ |
| AWS Rules Validation | ✓ | ✓ |
| Generate Terraform | ✓ | ✓ |
| Write Files | ✓ | ✓ |
| Database Persistence | ✗ | ✓ |
| Transaction Support | ✗ | ✓ |

## Related Files

- `backend/json-request-fiagram-complete.json` - Sample diagram JSON
- `backend/internal/platform/models/` - Database models
- `backend/internal/platform/repository/` - Repository implementations
- `backend/pkg/usecases/scenario5_terraform_codegen/` - Base scenario without persistence
