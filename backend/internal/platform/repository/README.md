# Repository Layer

This package provides GORM-based repository implementations for database operations.

## Repositories

- **BaseRepository**: Common database operations and transaction support

### Core domain

- **UserRepository**: User management operations
- **ProjectRepository**: Project management operations
- **ResourceRepository**: Resource management and relationships
- **PricingRepository**: Pricing data management
- **ResourceTypeRepository**: Lookup for cloud resource types
- **DependencyTypeRepository**: Lookup for dependency types
- **ResourceConstraintRepository**: Validation rules for resource types

### Infrastructure lookups

- **IACTargetRepository**: Lookup for supported IaC tools (Terraform, Pulumi, CDK, etc.)
- **ResourceCategoryRepository**: Lookup for resource categories (Compute, Storage, Networking)
- **ResourceKindRepository**: Lookup for resource kinds (VirtualMachine, Container, Function)

### Marketplace

- **CategoryRepository**: Marketplace template categories
- **TemplateRepository**: Templates with rich relationships and search helpers
- **ReviewRepository**: Template reviews and helpful votes
- **IACFormatRepository**: Marketplace IaC formats (Terraform, CDK, etc.)
- **TechnologyRepository**: Marketplace technology tags
- **ComplianceStandardRepository**: Compliance standards (SOC2, HIPAA, PCI, etc.)

### Join tables & relationships

- **TemplateComplianceRepository**: Template ↔ compliance associations
- **TemplateIACFormatRepository**: Template ↔ IaC format associations
- **TemplateTechnologyRepository**: Template ↔ technology associations
- **ResourceContainmentRepository**: Resource containment hierarchy (parent/child)
- **ResourceDependencyRepository**: Directed resource dependencies

## Usage Example

```go
import (
    "context"
    "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
    "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
    "github.com/google/uuid"
)

// Create a repository
userRepo, err := repository.NewUserRepository()
if err != nil {
    log.Fatal(err)
}

// Create a user
user := &models.User{
    Email: "user@example.com",
    Name:  stringPtr("John Doe"),
}
err = userRepo.Create(context.Background(), user)

// Find a user
foundUser, err := userRepo.FindByEmail(context.Background(), "user@example.com")

// Create a project
projectRepo, _ := repository.NewProjectRepository()
project := &models.Project{
    UserID:        user.ID,
    InfraToolID:   1, // Terraform
    Name:          "My Project",
    CloudProvider: "aws",
    Region:        "us-east-1",
}
err = projectRepo.Create(context.Background(), project)
```

## Transaction Support

```go
baseRepo, _ := repository.NewBaseRepository()
ctx := context.Background()

tx, ctx := baseRepo.BeginTransaction(ctx)
defer func() {
    if err != nil {
        baseRepo.RollbackTransaction(tx)
    } else {
        baseRepo.CommitTransaction(tx)
    }
}()

// Use repositories with transaction context
userRepo.Create(ctx, user)
projectRepo.Create(ctx, project)
```
