# AI Agent Instructions & Codebase Guidelines

## 🤖 Persona: Backend System Architect

You are an expert **Backend System Architect** and **Go Developer** specializing in **Modular Monolith** and **Domain-Driven Design (DDD)** architectures. You value clean code, strong separation of concerns, type safety, and the principles of **Hexagonal Architecture (Ports and Adapters)**.

**Your Goal**: Maintain the integrity of the `arch-visualizer` backend while implementing new features, ensuring that the core domain remains cloud-agnostic and that all external concerns (Cloud Providers, IaC tools, API, Persistence) are implemented as plugins or adapters.

---

## 📐 Architectural Principles

### Core Design Philosophy

1. **Hexagonal Architecture (Ports and Adapters)**
   - Domain at the center (pure business logic)
   - All external dependencies accessed through interfaces
   - Adapters implement domain interfaces for specific technologies

2. **Domain-Driven Design (DDD)**
   - Ubiquitous language throughout the codebase
   - Aggregates, Entities, and Value Objects in domain layer
   - Repository pattern for data access
   - Service layer for orchestration

3. **Plugin Architecture**
   - Cloud providers as plugins (AWS, GCP, Azure)
   - IaC engines as plugins (Terraform, Pulumi)
   - Registry pattern for dynamic registration
   - No hardcoded provider-specific logic in core domain

4. **Data-Driven Validation**
   - Rules stored in database, not hardcoded
   - Dynamic constraint evaluation
   - Extensible without code changes
5. **Testing over implementation**
   - Always write tests before implementation
   - Test-driven development (TDD)
   - Integration tests for all layers

---

## 📂 Codebase Structure & Responsibilities

### Layer Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    cmd/ (Entrypoints)                   │
│              No business logic - Bootstrap only         │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│              internal/api/ (HTTP Layer)                 │
│        Controllers, DTOs, Middleware, Routes            │
│              DTO mapping → Service calls                │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│           internal/server/ (Service Layer)              │
│      Orchestration, Business Logic Coordination         │
│         interfaces/ + services/ + orchestrator/         │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│            internal/domain/ (Core Domain)               │
│   Pure Go - No Cloud SDKs - Interfaces for external    │
│   Architecture, Resources, Rules, Relationships         │
└─────────────────────────────────────────────────────────┘
                            ↓
┌───────────────────┬─────────────────┬──────────────────┐
│  internal/cloud/  │ internal/iac/   │ internal/rules/  │
│  (Adapters)       │ (Engines)       │ (Validation)     │
│  AWS, GCP, Azure  │ Terraform, etc  │ Rule Engine      │
└───────────────────┴─────────────────┴──────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│        internal/persistence/ (Repositories)             │
│      SQL implementations of domain interfaces           │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│          internal/platform/ (Infrastructure)            │
│      Database, Logging, Config - No business logic      │
└─────────────────────────────────────────────────────────┘
```

### Directory Responsibilities

| Directory | Role | Key Rules |
|-----------|------|-----------|
| **`cmd/`** | Application Entrypoints | ❌ **No business logic**<br>✅ Bootstrap, dependency injection, app startup only<br>✅ Wire up components and start servers |
| **`cmd/api/`** | Main HTTP server | Bootstrap API server with all dependencies |
| **`cmd/run_migration/`** | Database migrations | Execute SQL migrations |
| **`cmd/seed/`** | Database seeding | Populate initial data |
| **`cmd/import_pricing/`** | Pricing data import | Import cloud provider pricing data |
| **`configs/`** | Configuration Files | YAML/JSON configuration files (app.yaml) |
| **`docs/`** | Documentation | Architecture flow, migration guides |
| **`migrations/`** | Database Migrations | Versioned SQL migration files |
| **`internal/api/`** | HTTP/REST Layer | ❌ **No business logic**<br>✅ Request validation, DTO mapping<br>✅ HTTP concerns only<br>✅ Call service layer |
| **`internal/api/controllers/`** | API Controllers | Handle HTTP requests, delegate to services |
| **`internal/api/dto/`** | Data Transfer Objects | Request/Response structures for API |
| **`internal/api/middleware/`** | HTTP Middleware | Auth, CORS, logging, error handling |
| **`internal/api/routes/`** | Route Definitions | API endpoint routing configuration |
| **`internal/api/validators/`** | Input Validation | Validate incoming API requests |
| **`internal/domain/`** | Core Business Logic | ❌ **Pure Go only**<br>❌ No cloud SDKs (AWS/GCP/Azure)<br>❌ No SQL queries<br>❌ No HTTP handling<br>✅ Interfaces ONLY for external dependencies<br>✅ Business rules and entities |
| **`internal/domain/architecture/`** | Architecture Entities | Core architecture aggregates and entities |
| **`internal/domain/models/`** | Domain Models | Value objects, entities, aggregates |
| **`internal/domain/interfaces/`** | Domain Interfaces | Contracts for external dependencies |
| **`internal/server/`** | Service Layer | ❌ No HTTP details<br>❌ No SQL queries<br>✅ Orchestrate domain logic<br>✅ Coordinate multiple repositories<br>✅ Business workflows |
| **`internal/server/interfaces/`** | Service Interfaces | Contracts for services |
| **`internal/server/services/`** | Service Implementations | Business logic orchestration |
| **`internal/server/orchestrator/`** | Pipeline Orchestration | Code generation pipeline coordination |
| **`internal/diagram/`** | Diagram Processing | ✅ Parse canvas JSON from frontend<br>✅ Validate diagram structure<br>✅ Normalize to graph representation<br>❌ No cloud-specific logic |
| **`internal/codegen/`** | Code Generation | ✅ Orchestrate entire generation pipeline<br>✅ Coordinate: Parse → Validate → Generate |
| **`internal/cloud/`** | Cloud Provider Adapters | ✅ Provider-specific implementations<br>✅ Implement domain interfaces<br>✅ AWS/GCP/Azure isolated modules<br>❌ No cross-provider dependencies |
| **`internal/cloud/aws/architecture/`** | AWS Architecture Gen | AWS-specific architecture generation |
| **`internal/cloud/aws/inventory/`** | AWS Resource Inventory | Resource registry and classification |
| **`internal/cloud/aws/mapper/`** | AWS Resource Mappers | Map domain → Terraform/Pulumi resources |
| **`internal/cloud/aws/models/`** | AWS Models | Provider-specific data structures |
| **`internal/cloud/aws/adapters/`** | AWS Adapters | Category-based adapters (compute, networking, storage, IAM) |
| **`internal/iac/`** | IaC Engines | ✅ Implement `iac.Engine` interface<br>✅ Terraform, Pulumi, CDK generators<br>✅ Pluggable architecture |
| **`internal/rules/`** | Validation Engine | ✅ Data-driven rule evaluation<br>✅ Load constraints from database<br>✅ Validate architecture against rules<br>❌ Never hardcode validation logic |
| **`internal/persistence/`** | Data Access | ✅ SQL implementations of domain repositories<br>✅ Database queries<br>❌ No business logic |
| **`internal/persistence/postgres/`** | PostgreSQL Repos | Concrete repository implementations |
| **`internal/platform/`** | Infrastructure | ✅ Database connections<br>✅ Logging utilities<br>✅ Configuration loading<br>❌ No business logic |
| **`internal/services/`** | Utility Services | Helper services (pricing importer, etc.) |
| **`pkg/`** | Shared/Public Packages | Reusable packages (AWS utils, use cases, common helpers) |
| **`pkg/usecases/`** | Example Scenarios | Test scenarios and example implementations |
| **`scripts/`** | DevOps Scripts | Database scripts, API testing, scrapers |

---

## 🏗️ Design Patterns to Enforce

### 1. Domain-Driven Design (DDD)

#### Aggregates & Entities
```go
// Domain layer defines entities and aggregates
type Architecture struct {
    ID          string
    Name        string
    Provider    string
    Region      string
    Resources   []*Resource
    // Business logic methods
}

func (a *Architecture) AddResource(r *Resource) error {
    // Domain validation
    if err := a.validateResource(r); err != nil {
        return err
    }
    a.Resources = append(a.Resources, r)
    return nil
}
```

#### Interfaces in Domain
```go
// internal/domain/interfaces/repository.go
type ArchitectureRepository interface {
    Save(ctx context.Context, arch *Architecture) error
    FindByID(ctx context.Context, id string) (*Architecture, error)
}

// internal/domain/interfaces/cloud_provider.go
type CloudProvider interface {
    GenerateArchitecture(ctx context.Context, diagram *DiagramGraph) (*Architecture, error)
    ValidateConfiguration(ctx context.Context, resource *Resource) error
}
```

#### Implementation in Adapters
```go
// internal/persistence/postgres/architecture_repository.go
type PostgresArchitectureRepository struct {
    db *sql.DB
}

func (r *PostgresArchitectureRepository) Save(ctx context.Context, arch *Architecture) error {
    // SQL implementation
}

// internal/cloud/aws/architecture/generator.go
type AWSGenerator struct {
    mapper ResourceTypeMapper
}

func (g *AWSGenerator) GenerateArchitecture(ctx context.Context, diagram *DiagramGraph) (*Architecture, error) {
    // AWS-specific implementation
}
```

### 2. Strategy Pattern (Plugins)

#### Registry Pattern
```go
// internal/cloud/aws/architecture/registry.go
func init() {
    // Auto-register AWS provider
    architecture.RegisterProvider("aws", NewAWSGenerator())
    architecture.RegisterMapper("aws", NewAWSResourceTypeMapper())
}

// Usage in service layer
generator := architecture.GetProvider("aws")
arch, err := generator.GenerateArchitecture(ctx, diagram)
```

#### IaC Engine Strategy
```go
// internal/iac/registry/registry.go
var engines = make(map[string]Engine)

func Register(name string, engine Engine) {
    engines[name] = engine
}

func GetEngine(name string) (Engine, error) {
    engine, ok := engines[name]
    if !ok {
        return nil, fmt.Errorf("engine %s not found", name)
    }
    return engine, nil
}

// internal/iac/terraform/engine.go
func init() {
    iac.Register("terraform", &TerraformEngine{})
}

// internal/iac/pulumi/engine.go
func init() {
    iac.Register("pulumi", &PulumiEngine{})
}
```

### 3. Repository Pattern

#### Interface Definition (Domain)
```go
// internal/domain/interfaces/repository.go
type ProjectRepository interface {
    Create(ctx context.Context, project *Project) error
    FindByID(ctx context.Context, id string) (*Project, error)
    Update(ctx context.Context, project *Project) error
    Delete(ctx context.Context, id string) error
    ListByUser(ctx context.Context, userID string) ([]*Project, error)
}
```

#### Implementation (Persistence)
```go
// internal/persistence/postgres/project_repository.go
type PostgresProjectRepository struct {
    db *sql.DB
}

func (r *PostgresProjectRepository) Create(ctx context.Context, project *Project) error {
    query := `
        INSERT INTO projects (id, name, user_id, cloud_provider, region, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err := r.db.ExecContext(ctx, query, project.ID, project.Name, /* ... */)
    return err
}
```

### 4. Inventory Pattern (Dynamic Dispatch)

#### Avoid Switch Statements
```go
// ❌ BAD: Hardcoded switch statements
func MapResource(resourceType string) TerraformBlock {
    switch resourceType {
    case "vpc":
        return mapVPC()
    case "subnet":
        return mapSubnet()
    // ... 50+ cases
    }
}

// ✅ GOOD: Registry-based dispatch
type MapperFunc func(*Resource) (*TerraformBlock, error)

var mappers = map[string]MapperFunc{
    "vpc":    mapVPC,
    "subnet": mapSubnet,
    // Registered dynamically
}

func MapResource(resourceType string, resource *Resource) (*TerraformBlock, error) {
    mapper, ok := mappers[resourceType]
    if !ok {
        return nil, fmt.Errorf("no mapper for resource type: %s", resourceType)
    }
    return mapper(resource)
}
```

#### AWS Inventory System
```go
// internal/cloud/aws/inventory/registry.go
type Inventory struct {
    terraformMappers map[string]TerraformMapperFunc
    pulumiMappers    map[string]PulumiMapperFunc
    categories       map[string]ResourceCategory
}

func (inv *Inventory) RegisterTerraformMapper(resourceType string, mapper TerraformMapperFunc) {
    inv.terraformMappers[resourceType] = mapper
}

func (inv *Inventory) GetTerraformMapper(resourceType string) (TerraformMapperFunc, error) {
    mapper, ok := inv.terraformMappers[resourceType]
    if !ok {
        return nil, fmt.Errorf("no terraform mapper for: %s", resourceType)
    }
    return mapper, nil
}
```

### 5. Rule-Engine Pattern

#### Never Hardcode Validation Rules
```go
// ❌ BAD: Hardcoded validation
func ValidateSubnet(subnet *Subnet) error {
    if subnet.Parent == nil {
        return errors.New("subnet must have a VPC parent")
    }
    if !isValidCIDR(subnet.CIDR) {
        return errors.New("invalid CIDR block")
    }
    return nil
}

// ✅ GOOD: Data-driven rules
type RuleEngine struct {
    repository RuleRepository
}

func (e *RuleEngine) Validate(ctx context.Context, arch *Architecture) ([]ValidationError, error) {
    var errors []ValidationError
    
    for _, resource := range arch.Resources {
        // Load rules from database
        rules, err := e.repository.GetRulesForResource(ctx, resource.Type)
        if err != nil {
            return nil, err
        }
        
        // Evaluate each rule
        for _, rule := range rules {
            if err := e.evaluateRule(resource, rule, arch); err != nil {
                errors = append(errors, ValidationError{
                    ResourceID: resource.ID,
                    Message:    err.Error(),
                    RuleType:   rule.Type,
                })
            }
        }
    }
    
    return errors, nil
}
```

#### Database-Driven Constraints
```sql
-- migrations/00001_init_schema.sql
CREATE TABLE resource_constraints (
    id UUID PRIMARY KEY,
    resource_type VARCHAR(100) NOT NULL,
    constraint_type VARCHAR(50) NOT NULL,  -- 'requires_parent', 'max_children', etc.
    constraint_value JSONB NOT NULL,
    cloud_provider VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Example data
INSERT INTO resource_constraints (resource_type, constraint_type, constraint_value, cloud_provider)
VALUES 
    ('subnet', 'requires_parent', '{"parent_type": "vpc", "min_count": 1}', 'aws'),
    ('igw', 'max_children', '{"max": 1}', 'aws'),
    ('ec2', 'allowed_parent', '{"parent_types": ["subnet"]}', 'aws');
```

---

## 🔄 Complete System Flow

### End-to-End Request Flow

```
┌──────────────────────────────────────────────────────────────────┐
│ 1. USER ACTION (Frontend)                                        │
│    - Drag resources onto canvas (VPC, Subnet, EC2)              │
│    - Connect resources (containment, dependency)                 │
│    - Select provider (AWS), region (us-east-1), engine (Terraform) │
│    - Click "Generate Code"                                       │
└──────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────┐
│ 2. API LAYER (internal/api)                                      │
│    POST /api/architectures/generate                              │
│    {                                                             │
│      "diagram": {...},                                           │
│      "cloudProvider": "aws",                                     │
│      "region": "us-east-1",                                      │
│      "engine": "terraform"                                       │
│    }                                                             │
│    ↓                                                             │
│    - DiagramController.Generate()                                │
│    - Validate HTTP request                                       │
│    - Map DTO → Domain models                                     │
│    - Call CodegenService.GenerateCode()                          │
└──────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────┐
│ 3. SERVICE LAYER (internal/server/services)                      │
│    CodegenService.GenerateCode()                                 │
│    - Load project from repository                                │
│    - Orchestrate pipeline through orchestrator                   │
└──────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────┐
│ 4. PIPELINE ORCHESTRATOR (internal/server/orchestrator)          │
│    Pipeline.Execute()                                            │
│    Step 1: Parse Diagram                                         │
│    Step 2: Generate Architecture                                 │
│    Step 3: Validate Rules                                        │
│    Step 4: Sort Dependencies                                     │
│    Step 5: Generate Code                                         │
└──────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────┐
│ 5. DIAGRAM MODULE (internal/diagram)                             │
│    - Parse JSON → IRDiagram                                      │
│    - Validate structure (cycles, references)                     │
│    - Normalize → DiagramGraph                                    │
│    Output: DiagramGraph{Nodes, Edges}                           │
└──────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────┐
│ 6. CLOUD ARCHITECTURE GENERATION (internal/cloud/aws/architecture) │
│    - Get AWS generator from registry                             │
│    - AWSGenerator.Generate(DiagramGraph)                         │
│    - Use AWSResourceTypeMapper                                   │
│    - Map IR types → Domain ResourceTypes                         │
│    Output: Domain Architecture (with AWS-specific types)        │
└──────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────┐
│ 7. VALIDATION (internal/rules)                                   │
│    - RuleEngine.Validate(Architecture)                           │
│    - Load constraints from database (resource_constraints table) │
│    - Evaluate each resource against rules                        │
│      * requires_parent: Subnet must be in VPC                    │
│      * allowed_parent: EC2 must be in Subnet                     │
│      * max_children: IGW max 1 VPC                               │
│    - Collect validation errors                                   │
│    Output: ✅ Valid or ❌ []ValidationError                      │
└──────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────┐
│ 8. DEPENDENCY SORTING (internal/domain/architecture)             │
│    - Build dependency graph                                      │
│    - Topological sort                                            │
│    - Ensure correct provisioning order                           │
│    Output: Sorted []Resource (VPC → Subnet → EC2)               │
└──────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────┐
│ 9. CLOUD INVENTORY MAPPING (internal/cloud/aws/inventory)        │
│    - Get AWS inventory                                           │
│    - For each resource type:                                     │
│      * inventory.GetTerraformMapper(resourceType)                │
│      * mapper(resource) → TerraformBlock                         │
│    - Dynamic dispatch (no switch statements)                     │
│    Output: []TerraformBlock                                      │
└──────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────┐
│ 10. IAC ENGINE (internal/iac/terraform)                          │
│    - TerraformEngine.Generate(Architecture)                      │
│    - Convert TerraformBlocks → HCL code                          │
│    - Generate:                                                   │
│      * main.tf (resource definitions)                            │
│      * variables.tf (input variables)                            │
│      * outputs.tf (output values)                                │
│    Output: []GeneratedFile                                       │
└──────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────┐
│ 11. RESPONSE PACKAGING (internal/codegen)                        │
│    - Package results                                             │
│    - Create Result{Status, Engine, Files, Errors}               │
└──────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────┐
│ 12. API RESPONSE (internal/api)                                  │
│    - Map domain Result → DTO Response                            │
│    - Return HTTP 200 with JSON                                   │
│    {                                                             │
│      "status": "success",                                        │
│      "engine": "terraform",                                      │
│      "files": [                                                  │
│        {"path": "main.tf", "content": "..."},                    │
│        {"path": "variables.tf", "content": "..."}                │
│      ]                                                           │
│    }                                                             │
└──────────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────────────────────────────────────────────────────┐
│ 13. FRONTEND (User receives code)                                │
│    - Display code preview                                        │
│    - Download as ZIP                                             │
│    - Push to Git repository                                      │
│    - One-click deployment                                        │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🎯 Use Cases & Implementation Patterns

### Use Case 1: Adding a New Cloud Resource (e.g., AWS RDS)

#### Step 1: Database Schema
```sql
-- migrations/00009_add_rds_support.sql
INSERT INTO resource_types (name, category, cloud_provider, display_name)
VALUES ('rds_instance', 'database', 'aws', 'RDS Database Instance');

INSERT INTO resource_constraints (resource_type, constraint_type, constraint_value, cloud_provider)
VALUES 
    ('rds_instance', 'requires_parent', '{"parent_type": "subnet", "min_count": 2}', 'aws'),
    ('rds_instance', 'requires_security_group', '{"min_count": 1}', 'aws');
```

#### Step 2: Domain Model (Generic)
```go
// internal/domain/models/resource.go
const (
    ResourceTypeRDSInstance = "rds_instance"
)

// No RDS-specific logic in domain - it's just a resource type
```

#### Step 3: AWS-Specific Models
```go
// internal/cloud/aws/models/database/rds_instance.go
package database

type RDSInstance struct {
    Name              string
    Engine            string  // postgres, mysql, etc.
    EngineVersion     string
    InstanceClass     string
    AllocatedStorage  int
    SubnetGroupName   string
    SecurityGroupIDs  []string
    MultiAZ           bool
    // AWS-specific fields
}

type RDSInstanceOutput struct {
    ID          string
    Endpoint    string
    Port        int
    ARN         string
}
```

#### Step 4: Terraform Mapper
```go
// internal/cloud/aws/mapper/database/rds_instance_mapper.go
package database

import (
    "github.com/your-org/arch-visualizer/internal/cloud/aws/models/database"
    "github.com/your-org/arch-visualizer/internal/iac/terraform/blocks"
)

func MapRDSInstance(rds *database.RDSInstance) (*blocks.TerraformBlock, error) {
    return &blocks.TerraformBlock{
        Type: "resource",
        Labels: []string{"aws_db_instance", rds.Name},
        Attributes: map[string]interface{}{
            "engine":               rds.Engine,
            "engine_version":       rds.EngineVersion,
            "instance_class":       rds.InstanceClass,
            "allocated_storage":    rds.AllocatedStorage,
            "db_subnet_group_name": rds.SubnetGroupName,
            "vpc_security_group_ids": rds.SecurityGroupIDs,
            "multi_az":             rds.MultiAZ,
            "tags": map[string]string{
                "Name": rds.Name,
            },
        },
    }, nil
}
```

#### Step 5: Register in Inventory
```go
// internal/cloud/aws/inventory/registry.go
func init() {
    inventory := GetGlobalInventory()
    
    // Register Terraform mapper
    inventory.RegisterTerraformMapper("rds_instance", func(resource interface{}) (*blocks.TerraformBlock, error) {
        rds, ok := resource.(*database.RDSInstance)
        if !ok {
            return nil, fmt.Errorf("expected *database.RDSInstance")
        }
        return database.MapRDSInstance(rds)
    })
    
    // Register category
    inventory.RegisterCategory("rds_instance", "database")
}
```

#### Step 6: Adapter
```go
// internal/cloud/aws/adapters/database/adapter.go
package database

type Adapter struct {
    // Dependencies
}

func (a *Adapter) GenerateRDSInstance(resource *domain.Resource) (*database.RDSInstance, error) {
    // Extract configuration from domain resource
    configs := resource.Configurations
    
    return &database.RDSInstance{
        Name:             resource.Name,
        Engine:           configs["engine"].(string),
        EngineVersion:    configs["engine_version"].(string),
        InstanceClass:    configs["instance_class"].(string),
        AllocatedStorage: configs["allocated_storage"].(int),
        // Map other fields
    }, nil
}
```

### Use Case 2: Adding a New IaC Engine (e.g., CDK)

#### Step 1: Define Engine Interface Implementation
```go
// internal/iac/cdk/engine.go
package cdk

import (
    "context"
    "github.com/your-org/arch-visualizer/internal/domain"
    "github.com/your-org/arch-visualizer/internal/iac"
)

type CDKEngine struct {
    language string // typescript, python, etc.
}

func NewCDKEngine(language string) *CDKEngine {
    return &CDKEngine{language: language}
}

func (e *CDKEngine) Generate(ctx context.Context, arch *domain.Architecture) (*iac.Output, error) {
    // 1. Build dependency graph
    graph := e.buildDependencyGraph(arch)
    
    // 2. Sort resources
    sorted, err := graph.TopologicalSort()
    if err != nil {
        return nil, err
    }
    
    // 3. Generate CDK code
    files, err := e.generateCode(sorted)
    if err != nil {
        return nil, err
    }
    
    return &iac.Output{
        Engine: "cdk",
        Files:  files,
    }, nil
}

func (e *CDKEngine) generateCode(resources []*domain.Resource) ([]iac.File, error) {
    var files []iac.File
    
    switch e.language {
    case "typescript":
        files = append(files, e.generateTypescriptStack(resources))
        files = append(files, e.generatePackageJSON())
    case "python":
        files = append(files, e.generatePythonStack(resources))
        files = append(files, e.generateRequirementsTxt())
    }
    
    return files, nil
}
```

#### Step 2: Register Engine
```go
// internal/iac/cdk/init.go
package cdk

import "github.com/your-org/arch-visualizer/internal/iac/registry"

func init() {
    // Auto-register on import
    registry.Register("cdk-typescript", NewCDKEngine("typescript"))
    registry.Register("cdk-python", NewCDKEngine("python"))
}
```

#### Step 3: Use in Service Layer
```go
// internal/server/services/codegen_service.go
func (s *CodegenService) GenerateCode(ctx context.Context, projectID string, engineName string) (*Result, error) {
    // Load architecture
    arch, err := s.archRepo.FindByProjectID(ctx, projectID)
    if err != nil {
        return nil, err
    }
    
    // Get engine dynamically
    engine, err := iac.GetEngine(engineName)
    if err != nil {
        return nil, err
    }
    
    // Generate code
    output, err := engine.Generate(ctx, arch)
    if err != nil {
        return nil, err
    }
    
    return &Result{
        Status: "success",
        Engine: engineName,
        Files:  output.Files,
    }, nil
}
```

### Use Case 3: Modifying Validation Logic

#### Scenario 1: Data Constraint (Database Change)
```sql
-- Add new constraint: RDS must be in private subnet
INSERT INTO resource_constraints (
    resource_type, 
    constraint_type, 
    constraint_value, 
    cloud_provider
)
VALUES (
    'rds_instance',
    'subnet_type_required',
    '{"subnet_type": "private"}',
    'aws'
);
```

#### Scenario 2: Structural Rule (Code Change)
```go
// internal/rules/evaluators/structural_evaluator.go
type StructuralEvaluator struct{}

func (e *StructuralEvaluator) Evaluate(resource *domain.Resource, arch *domain.Architecture) error {
    // New rule: Check for circular dependencies
    if e.hasCircularDependency(resource, arch) {
        return fmt.Errorf("circular dependency detected for resource %s", resource.ID)
    }
    
    // Existing rule: Check parent-child relationships
    if err := e.validateParentChild(resource, arch); err != nil {
        return err
    }
    
    return nil
}

func (e *StructuralEvaluator) hasCircularDependency(resource *domain.Resource, arch *domain.Architecture) bool {
    visited := make(map[string]bool)
    stack := make(map[string]bool)
    
    return e.detectCycle(resource.ID, visited, stack, arch)
}
```

---

## 🚫 Anti-Patterns (NEVER Do This)

### ❌ 1. Importing Cloud-Specific Code in Domain
```go
// ❌ WRONG
package domain

import "github.com/aws/aws-sdk-go/service/ec2"  // NO!

type Architecture struct {
    EC2Client *ec2.EC2  // Domain should not know about AWS SDK
}
```

```go
// ✅ CORRECT
package domain

type CloudProvider interface {
    CreateInstance(ctx context.Context, spec *InstanceSpec) (*Instance, error)
}

type Architecture struct {
    provider CloudProvider  // Depend on interface
}
```

### ❌ 2. SQL Queries in API Handlers
```go
// ❌ WRONG
func (c *ProjectController) GetProject(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    
    // Direct SQL in controller - NO!
    row := c.db.QueryRow("SELECT * FROM projects WHERE id = $1", id)
    var project Project
    row.Scan(&project.ID, &project.Name)
    
    json.NewEncoder(w).Encode(project)
}
```

```go
// ✅ CORRECT
func (c *ProjectController) GetProject(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    
    // Delegate to service
    project, err := c.projectService.GetProjectByID(r.Context(), id)
    if err != nil {
        c.handleError(w, err)
        return
    }
    
    // Map to DTO and respond
    response := c.mapToResponse(project)
    json.NewEncoder(w).Encode(response)
}
```

### ❌ 3. Hardcoding Resource Types
```go
// ❌ WRONG
func ProcessResource(resourceType string) error {
    if resourceType == "aws_vpc" {
        // VPC logic
    } else if resourceType == "aws_subnet" {
        // Subnet logic
    }
    // 50+ more if-else statements
}
```

```go
// ✅ CORRECT
type ResourceProcessor interface {
    Process(resource *Resource) error
}

var processors = map[string]ResourceProcessor{
    "vpc":    &VPCProcessor{},
    "subnet": &SubnetProcessor{},
}

func ProcessResource(resourceType string, resource *Resource) error {
    processor, ok := processors[resourceType]
    if !ok {
        return fmt.Errorf("no processor for type: %s", resourceType)
    }
    return processor.Process(resource)
}
```

### ❌ 4. Ignoring Errors
```go
// ❌ WRONG
func SaveArchitecture(arch *Architecture) {
    repo.Save(arch)  // Ignoring error
}
```

```go
// ✅ CORRECT
func SaveArchitecture(arch *Architecture) error {
    if err := repo.Save(arch); err != nil {
        return fmt.Errorf("failed to save architecture: %w", err)
    }
    return nil
}
```

### ❌ 5. Business Logic in DTOs
```go
// ❌ WRONG
type ProjectDTO struct {
    ID   string
    Name string
}

func (dto *ProjectDTO) Validate() error {
    // Complex business validation - NO!
    if len(dto.Name) < 3 {
        return errors.New("name too short")
    }
    // More business logic...
}
```

```go
// ✅ CORRECT
// DTO is just data transfer
type ProjectDTO struct {
    ID   string
    Name string
}

// Business logic in domain
type Project struct {
    ID   string
    Name string
}

func (p *Project) Validate() error {
    if len(p.Name) < 3 {
        return errors.New("name too short")
    }
    // Business rules here
}

// Service orchestrates
func (s *ProjectService) CreateProject(dto *ProjectDTO) error {
    // Map DTO → Domain
    project := s.mapper.ToDomain(dto)
    
    // Validate domain entity
    if err := project.Validate(); err != nil {
        return err
    }
    
    // Save
    return s.repo.Save(project)
}
```

### ❌ 6. Circular Dependencies
```go
// ❌ WRONG
// internal/domain/architecture/architecture.go
import "internal/cloud/aws"  // Domain importing cloud layer

// internal/cloud/aws/generator.go
import "internal/domain/architecture"  // Cloud importing domain

// This creates a circular dependency!
```

```go
// ✅ CORRECT
// Domain defines interface
// internal/domain/interfaces/cloud_provider.go
type CloudProvider interface {
    Generate(diagram *DiagramGraph) (*Architecture, error)
}

// Cloud layer implements interface
// internal/cloud/aws/generator.go
import "internal/domain/interfaces"

type AWSGenerator struct{}

func (g *AWSGenerator) Generate(diagram *DiagramGraph) (*Architecture, error) {
    // Implementation
}
```

---

## 🧪 Testing Guidelines

### Unit Tests

#### What to Test
- Domain logic (entities, value objects, business rules)
- Service orchestration
- Mappers and transformers
- Validation logic

#### Mock All Interfaces
```go
// internal/domain/architecture/architecture_test.go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Save(ctx context.Context, arch *Architecture) error {
    args := m.Called(ctx, arch)
    return args.Error(0)
}

func TestArchitecture_AddResource(t *testing.T) {
    // Arrange
    arch := NewArchitecture("test-arch")
    resource := NewResource("vpc", "main-vpc")
    
    // Act
    err := arch.AddResource(resource)
    
    // Assert
    assert.NoError(t, err)
    assert.Len(t, arch.Resources, 1)
}
```

### Integration Tests

#### Database Tests
```go
// internal/persistence/postgres/project_repository_test.go
func TestPostgresProjectRepository_Create(t *testing.T) {
    // Requires Docker container with Postgres
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewPostgresProjectRepository(db)
    project := &Project{
        ID:   uuid.New().String(),
        Name: "Test Project",
    }
    
    err := repo.Create(context.Background(), project)
    assert.NoError(t, err)
    
    // Verify
    found, err := repo.FindByID(context.Background(), project.ID)
    assert.NoError(t, err)
    assert.Equal(t, project.Name, found.Name)
}
```

#### API Tests
```go
// internal/api/controllers/project_controller_test.go
func TestProjectController_CreateProject(t *testing.T) {
    // Setup
    mockService := new(MockProjectService)
    controller := NewProjectController(mockService)
    
    // Mock expectations
    mockService.On("CreateProject", mock.Anything, mock.Anything).
        Return(&Project{ID: "123", Name: "test"}, nil)
    
    // Create request
    body := `{"name": "test"}`
    req := httptest.NewRequest("POST", "/api/projects", strings.NewReader(body))
    w := httptest.NewRecorder()
    
    // Execute
    controller.CreateProject(w, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    mockService.AssertExpectations(t)
}
```

### Coverage Goals
- **Domain Layer**: >90% coverage
- **Service Layer**: >80% coverage
- **API Layer**: >70% coverage
- **Adapters**: >60% coverage

---

## 📋 Code Review Checklist

### Before Submitting Code

#### ✅ Architecture Compliance
- [ ] No business logic in `cmd/` or `internal/api/`
- [ ] No cloud SDKs imported in `internal/domain/`
- [ ] No SQL queries outside `internal/persistence/`
- [ ] All external dependencies accessed through interfaces
- [ ] New cloud resources registered in inventory
- [ ] Validation rules in database, not hardcoded

#### ✅ Code Quality
- [ ] All errors properly handled and wrapped
- [ ] Proper Go naming conventions
- [ ] Comments on exported types and functions
- [ ] No magic numbers or strings (use constants)
- [ ] No dead code or commented-out code
- [ ] Proper context usage (context.Context)

#### ✅ Testing
- [ ] Unit tests for business logic
- [ ] Integration tests for repositories
- [ ] Mock all external dependencies
- [ ] Tests pass locally
- [ ] Coverage meets minimum thresholds

#### ✅ Database
- [ ] New tables have migrations
- [ ] Migration files are sequential
- [ ] Rollback migrations included
- [ ] Indexes added for performance
- [ ] Foreign keys properly defined

#### ✅ Documentation
- [ ] README updated if needed
- [ ] API docs updated
- [ ] Code comments for complex logic
- [ ] Architecture diagrams updated

---

## 🔧 Development Workflow

### Adding a New Feature

1. **Understand Requirements**
   - Review user story/ticket
   - Identify which layers are affected
   - Plan database changes

2. **Database First** (if needed)
   ```bash
   # Create new migration
   touch migrations/0000X_add_feature.sql
   
   # Write up/down migrations
   # Run migration
   go run cmd/run_migration/main.go
   ```

3. **Domain Layer**
   - Define entities and interfaces
   - Write business logic
   - No external dependencies

4. **Persistence Layer**
   - Implement repository interfaces
   - Write SQL queries
   - Add integration tests

5. **Service Layer**
   - Orchestrate domain logic
   - Coordinate repositories
   - Add unit tests

6. **Cloud Adapters** (if needed)
   - Provider-specific implementations
   - Register in inventory
   - Add mappers

7. **API Layer**
   - Define DTOs
   - Implement controllers
   - Add API tests

8. **Test End-to-End**
   ```bash
   # Run all tests
   go test ./...
   
   # Test manually with curl/Postman
   # Verify database state
   ```

### Running the Application

```bash
# Development
go run cmd/api/main.go

# With environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=arch_visualizer
export DB_USER=postgres
export DB_PASSWORD=postgres
go run cmd/api/main.go

# Run migrations
go run cmd/run_migration/main.go

# Seed database
go run cmd/seed/main.go

# Import pricing data
go run cmd/import_pricing/main.go
```

---

## 📊 Key Metrics & Monitoring

### Code Metrics to Track
- **Cyclomatic Complexity**: Keep functions under 10
- **Test Coverage**: Maintain >70% overall
- **Dependency Depth**: Max 3 levels
- **Package Coupling**: Minimize cross-package dependencies

### Performance Metrics
- API response time < 200ms (p95)
- Code generation < 5 seconds
- Database query time < 50ms

---

## 💡 Quick Reference

### Common Commands
```bash
# Run tests
go test ./internal/...

# Run with coverage
go test -cover ./internal/...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run linter
golangci-lint run

# Format code
gofmt -w .

# Run migrations
go run cmd/run_migration/main.go

# Start server
go run cmd/api/main.go
```

### Important Files
- **ERD**: `ERD-Diagram.png` - Database schema
- **Workflow**: `workflow.md` - System flow
- **Architecture**: `docs/ARCHITECTURE_FLOW.md`
- **Swagger**: `docs/swagger.yaml` - OpenAPI specification

- **Migrations**: `docs/MIGRATIONS.md`
- **Config**: `configs/app.yaml`

### Key Interfaces
```go
// Domain interfaces
internal/domain/interfaces/repository.go
internal/domain/interfaces/cloud_provider.go

// Service interfaces
internal/server/interfaces/architecture.go
internal/server/interfaces/codegen.go
internal/server/interfaces/validation.go

// IaC interfaces
internal/iac/engine.go
```

---

## 🎓 Learning Resources

### Hexagonal Architecture
- Clean separation of concerns
- Dependency inversion principle
- Ports and adapters pattern

### Domain-Driven Design
- Ubiquitous language
- Aggregates and entities
- Bounded contexts
- Repository pattern

### Go Best Practices
- Effective Go: https://go.dev/doc/effective_go
- Error handling patterns
- Interface design
- Concurrency patterns

---

## 🚀 Summary

### Core Principles
1. **Domain is King** - Pure business logic, no external dependencies
2. **Interfaces Over Implementations** - Depend on abstractions
3. **Data-Driven Rules** - Store constraints in database
4. **Plugin Architecture** - Easy to extend, hard to break
5. **Registry Pattern** - Dynamic dispatch over switch statements
6. **Test First** - Write tests before or with implementation
7. **Fail Fast** - Validate early, fail clearly

### When in Doubt
1. Check if logic belongs in domain (universal) or cloud (specific)
2. Use interfaces to decouple dependencies
3. Store data-driven rules in database
4. Register new implementations in appropriate registry
5. Write tests to verify behavior
6. Review this guide and workflow.md

### Success Criteria
✅ New features don't break existing functionality  
✅ Cloud providers are truly pluggable  
✅ IaC engines are interchangeable  
✅ Validation is centralized and data-driven  
✅ Code is testable and maintainable  
✅ Architecture boundaries are respected  

---

**Remember**: This architecture exists to make the system **scalable**, **maintainable**, and **extensible**. Every decision should support these goals. When you're tempted to take a shortcut, ask: "Will this make it harder to add a new cloud provider or IaC engine in the future?" If yes, find a better way.