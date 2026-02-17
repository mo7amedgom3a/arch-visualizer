# Codebase Layers Understanding

This document provides a comprehensive overview of the `arch-visualizer` backend architecture, explaining how each layer works and how they interact.

## 🏗️ Architecture Overview

The codebase follows **Hexagonal Architecture (Ports and Adapters)** combined with **Domain-Driven Design (DDD)** principles. The architecture is organized in clear layers with strict separation of concerns.

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
│      internal/platform/server/ (Service Layer)          │
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
│  internal/cloud/ │ internal/iac/   │ internal/rules/  │
│  (Adapters)       │ (Engines)       │ (Validation)     │
│  AWS, GCP, Azure  │ Terraform, etc  │ Rule Engine      │
└───────────────────┴─────────────────┴──────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│    internal/platform/repository/ (Repositories)         │
│      SQL implementations of domain interfaces           │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│          internal/platform/ (Infrastructure)            │
│      Database, Logging, Config - No business logic      │
└─────────────────────────────────────────────────────────┘
```

---

## 📂 Layer-by-Layer Breakdown

### 1. Entry Points (`cmd/`)

**Purpose**: Application bootstrap and dependency injection only. **NO business logic**.

**Structure**:
```
cmd/
├── api/              # Main HTTP server entry point
├── run_migration/    # Database migration runner
├── seed/             # Database seeding
├── import_pricing/   # Pricing data import
└── runner/           # Other utilities
```

**Key File**: `cmd/api/main.go`
- Connects to database
- Initializes logger
- Creates server instance
- Sets up routes
- Starts HTTP server

**Example**:
```go
func main() {
    // Connect to database
    database.Connect()
    
    // Initialize logger
    logger.Init(logger.Config{LogDir: "log"})
    
    // Create server (wires all dependencies)
    srv, _ := server.NewServer(log)
    
    // Setup routes
    r := routes.SetupRouter(srv)
    
    // Start server
    r.Run(":9000")
}
```

**Key Principle**: This layer only wires components together. All business logic is delegated to lower layers.

---

### 2. API Layer (`internal/api/`)

**Purpose**: HTTP request handling, DTO mapping, and routing. **NO business logic**.

**Structure**:
```
internal/api/
├── controllers/      # HTTP request handlers
├── dto/              # Data Transfer Objects (request/response)
├── middleware/       # HTTP middleware (auth, CORS, logging)
├── routes/           # Route definitions
└── validators/       # Input validation
```

**Key Components**:

#### Controllers (`controllers/`)
- Handle HTTP requests
- Validate input
- Map DTOs to domain models
- Call service layer
- Map domain results to DTOs
- Return HTTP responses

**Example**: `generation_controller.go`
```go
func (ctrl *GenerationController) GenerateCode(c *gin.Context) {
    // 1. Parse request
    var req dto.GenerateCodeRequest
    c.ShouldBindJSON(&req)
    
    // 2. Call orchestrator (service layer)
    out, err := ctrl.orchestrator.GenerateCode(ctx, &GenerateCodeRequest{
        ProjectID: projectID,
        Engine: req.Tool,
    })
    
    // 3. Map to response DTO
    resp := &dto.GenerationResponse{
        Status: "completed",
        Files: files,
    }
    
    // 4. Return HTTP response
    c.JSON(http.StatusOK, resp)
}
```

#### DTOs (`dto/`)
- Request structures (what the API receives)
- Response structures (what the API returns)
- No business logic, just data containers

**Key Principle**: Controllers are thin - they only handle HTTP concerns and delegate to services.

---

### 3. Service Layer (`internal/platform/server/`)

**Purpose**: Orchestrates business workflows, coordinates multiple repositories, and coordinates domain logic.

**Structure**:
```
internal/platform/server/
├── interfaces/       # Service interfaces (contracts)
├── services/         # Service implementations
└── orchestrator/    # Pipeline orchestration
```

**Key Components**:

#### Service Interfaces (`interfaces/`)
Define contracts for services:
- `CodegenService` - Code generation
- `ArchitectureService` - Architecture operations
- `DiagramService` - Diagram parsing/validation
- `ProjectService` - Project management
- `PricingService` - Cost estimation
- `ValidationService` - Rule validation

**Example**: `interfaces/codegen.go`
```go
type CodegenService interface {
    Generate(ctx context.Context, arch *architecture.Architecture, engine string) (*iac.Output, error)
    SupportedEngines() []string
}
```

#### Service Implementations (`services/`)
Implement business workflows:
- Coordinate multiple repositories
- Call domain logic
- Handle transactions
- Transform data between layers

**Example**: `services/codegen_service.go`
```go
func (s *CodegenService) Generate(ctx context.Context, arch *architecture.Architecture, engine string) (*iac.Output, error) {
    // 1. Get IaC engine from registry
    iacEngine, err := iac.GetEngine(engine)
    
    // 2. Generate code
    output, err := iacEngine.Generate(ctx, arch)
    
    return output, nil
}
```

#### Pipeline Orchestrator (`orchestrator/`)
Coordinates the entire code generation pipeline:

**Flow**:
1. Parse Diagram (JSON → DiagramGraph)
2. Validate Diagram (structural validation)
3. Generate Architecture (DiagramGraph → Architecture)
4. Validate Rules (constraint validation)
5. Generate Code (Architecture → IaC files)

**Example**: `orchestrator/pipeline.go`
```go
func (o *PipelineOrchestratorImpl) GenerateCode(ctx context.Context, req *GenerateCodeRequest) (*iac.Output, error) {
    // Step 1: Get project
    project, err := o.projectService.GetByID(ctx, req.ProjectID)
    
    // Step 2: Load architecture
    arch, err := o.projectService.LoadArchitecture(ctx, req.ProjectID)
    
    // Step 3: Validate rules
    ruleValidationResult, err := o.architectureService.ValidateRules(ctx, arch, provider)
    
    // Step 4: Generate code
    output, err := o.codegenService.Generate(ctx, arch, engine)
    
    return output, nil
}
```

**Key Principle**: Services orchestrate workflows but don't contain domain logic - they delegate to domain entities.

---

### 4. Domain Layer (`internal/domain/`)

**Purpose**: Pure business logic, cloud-agnostic, no external dependencies.

**Structure**:
```
internal/domain/
├── architecture/     # Architecture aggregate root
├── resource/         # Resource entities (compute, networking, storage, IAM)
├── rules/            # Rule engine interfaces
├── constraint/       # Constraint evaluation
├── pricing/          # Pricing domain models
└── errors/            # Domain errors
```

**Key Components**:

#### Architecture Aggregate (`architecture/`)
- Top-level aggregate root
- Contains resources and relationships
- Business logic methods (AddResource, Validate, etc.)

**Example**: `architecture/aggregate.go`
```go
type Architecture struct {
    ID        string
    Name      string
    Provider  CloudProvider
    Region    string
    Resources []*Resource
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

#### Resource Entities (`resource/`)
- Cloud-agnostic resource models
- Organized by category:
  - `compute/` - EC2, Lambda, LoadBalancer
  - `networking/` - VPC, Subnet, SecurityGroup
  - `storage/` - S3, EBS
  - `iam/` - Role, Policy, User
  - `database/` - RDS

**Example**: `resource/compute/instance.go`
```go
type Instance struct {
    ID          string
    Name        string
    InstanceType string
    AMI         string
    // ... cloud-agnostic fields
}
```

#### Domain Interfaces (`resource/*/service.go`)
Define contracts for external operations:
- `ComputeService` - Compute operations
- `NetworkingService` - Networking operations
- `StorageService` - Storage operations

**Example**: `resource/compute/service.go`
```go
type ComputeService interface {
    CreateInstance(ctx context.Context, instance *Instance) (*Instance, error)
    GetInstance(ctx context.Context, id string) (*Instance, error)
}
```

**Key Principle**: Domain layer is pure Go - no AWS SDK, no SQL, no HTTP. Only interfaces for external dependencies.

---

### 5. Cloud Adapters (`internal/cloud/`)

**Purpose**: Provider-specific implementations that adapt domain interfaces to cloud providers.

**Structure**:
```
internal/cloud/
├── aws/
│   ├── architecture/     # AWS architecture generator
│   ├── inventory/        # AWS resource registry
│   ├── mapper/           # Domain → IaC mappers
│   ├── models/           # AWS-specific models
│   ├── adapters/         # AWS SDK adapters
│   ├── pricing/          # AWS pricing calculators
│   └── rules/            # AWS rule implementations
├── azure/                # Azure implementations (stub)
└── gcp/                  # GCP implementations (stub)
```

**Key Components**:

#### Architecture Generator (`aws/architecture/`)
Converts diagram graphs to domain architectures:
- Implements `ArchitectureGenerator` interface
- Maps IR types to AWS resource types
- Builds relationships (containment, dependencies)

**Example**: `aws/architecture/generator.go`
```go
type AWSArchitectureGenerator struct {
    mapper ResourceTypeMapper
}

func (g *AWSArchitectureGenerator) Generate(ctx context.Context, diagram *DiagramGraph) (*Architecture, error) {
    // Map diagram nodes to AWS resources
    // Build relationships
    // Return domain architecture
}
```

**Registration**: Auto-registers via `init()`:
```go
func init() {
    generator := NewAWSArchitectureGenerator()
    architecture.RegisterGenerator(generator)
    
    mapper := NewAWSResourceTypeMapper()
    architecture.RegisterResourceTypeMapper(resource.AWS, mapper)
}
```

#### Inventory System (`aws/inventory/`)
Registry-based resource mapping (no switch statements):
- Maps IR types → Resource Names
- Registers Terraform mappers
- Registers Pulumi mappers
- Dynamic dispatch pattern

**Example**: `aws/inventory/inventory.go`
```go
type Inventory struct {
    terraformMappers map[string]TerraformMapperFunc
    pulumiMappers    map[string]PulumiMapperFunc
}

func (inv *Inventory) GetTerraformMapper(resourceType string) (TerraformMapperFunc, error) {
    mapper, ok := inv.terraformMappers[resourceType]
    if !ok {
        return nil, fmt.Errorf("no mapper for: %s", resourceType)
    }
    return mapper, nil
}
```

#### Resource Mappers (`aws/mapper/`)
Convert domain resources to IaC blocks:
- Category-based organization (compute, networking, storage)
- Each resource has a mapper function
- Registered in inventory

**Example**: `aws/mapper/compute/ec2_mapper.go`
```go
func MapEC2Instance(instance *compute.Instance) (*tfmapper.TerraformBlock, error) {
    return &tfmapper.TerraformBlock{
        Kind:   "resource",
        Labels: []string{"aws_instance", instance.Name},
        Attributes: map[string]tfmapper.TerraformValue{
            "instance_type": strVal(instance.InstanceType),
            "ami":           strVal(instance.AMI),
        },
    }, nil
}
```

**Key Principle**: Cloud adapters implement domain interfaces. Domain doesn't know about AWS - it only knows about interfaces.

---

### 6. IaC Engines (`internal/iac/`)

**Purpose**: Generate Infrastructure-as-Code files (Terraform, Pulumi, etc.)

**Structure**:
```
internal/iac/
├── engine.go          # Engine interface
├── registry/          # Engine registry
├── terraform/          # Terraform engine
└── pulumi/            # Pulumi engine (stub)
```

**Key Components**:

#### Engine Interface (`engine.go`)
```go
type Engine interface {
    Name() string
    Generate(ctx context.Context, arch *architecture.Architecture) (*Output, error)
}
```

#### Engine Registry (`registry/`)
Dynamic engine registration:
```go
type EngineRegistry struct {
    engines map[string]iac.Engine
}

func (r *EngineRegistry) Register(engine iac.Engine) error {
    r.engines[engine.Name()] = engine
    return nil
}

func (r *EngineRegistry) Get(name string) (iac.Engine, bool) {
    e, ok := r.engines[name]
    return e, ok
}
```

#### Terraform Engine (`terraform/`)
- Converts Architecture → Terraform blocks
- Generates HCL code
- Creates multiple files (main.tf, variables.tf, outputs.tf)

**Key Principle**: Engines are pluggable - add a new engine by implementing the `Engine` interface and registering it.

---

### 7. Rules Engine (`internal/domain/rules/`)

**Purpose**: Data-driven validation - rules stored in database, not hardcoded.

**Structure**:
```
internal/domain/rules/
├── types.go           # Rule types and interfaces
├── engine/            # Rule evaluation engine
├── constraints/       # Constraint evaluators
└── registry/          # Rule registry
```

**Key Components**:

#### Rule Interface (`types.go`)
```go
type Rule interface {
    GetType() RuleType
    GetResourceType() string
    GetValue() string
    Evaluate(ctx context.Context, evalCtx *EvaluationContext) error
}
```

#### Rule Types
- `requires_parent` - Resource must have parent
- `allowed_parent` - Allowed parent types
- `max_children` - Maximum child count
- `allowed_dependencies` - Allowed dependency types
- `forbidden_dependencies` - Forbidden dependency types

#### Database-Driven Rules
Rules are stored in `resource_constraints` table:
```sql
CREATE TABLE resource_constraints (
    id UUID PRIMARY KEY,
    resource_type VARCHAR(100) NOT NULL,
    constraint_type VARCHAR(50) NOT NULL,
    constraint_value JSONB NOT NULL,
    cloud_provider VARCHAR(20) NOT NULL
);

-- Example
INSERT INTO resource_constraints VALUES (
    'subnet', 'requires_parent', '{"parent_type": "vpc"}', 'aws'
);
```

**Key Principle**: Never hardcode validation rules - always load from database.

---

### 8. Repository Layer (`internal/platform/repository/`)

**Purpose**: SQL implementations of domain repository interfaces.

**Structure**:
```
internal/platform/repository/
├── project_repository.go
├── resource_repository.go
├── resource_type_repository.go
├── constraint_repository.go
└── ... (many more)
```

**Key Components**:

#### Base Repository (`base_repository.go`)
Provides common database operations:
```go
type BaseRepository struct {
    db *gorm.DB
}

func (r *BaseRepository) GetDB(ctx context.Context) *gorm.DB {
    return r.db.WithContext(ctx)
}
```

#### Repository Implementations
Implement domain repository interfaces:
```go
type ProjectRepository struct {
    *BaseRepository
}

func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
    return r.GetDB(ctx).Create(project).Error
}

func (r *ProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
    var project models.Project
    err := r.GetDB(ctx).First(&project, "id = ?", id).Error
    return &project, err
}
```

**Key Principle**: Repositories contain SQL queries but no business logic - they're pure data access.

---

### 9. Infrastructure Layer (`internal/platform/`)

**Purpose**: Cross-cutting concerns - database, logging, configuration.

**Structure**:
```
internal/platform/
├── database/          # Database connection
├── logger/            # Logging utilities
├── config/            # Configuration loading
├── models/            # Database models (GORM)
└── errors/            # Error handling
```

**Key Components**:

#### Database (`database/`)
- Connection management
- Migration execution
- Connection pooling

#### Logger (`logger/`)
- Structured logging
- File and console output
- Log rotation

#### Models (`models/`)
- GORM database models
- Table definitions
- Relationships

**Key Principle**: Infrastructure provides utilities but no business logic.

---

## 🔄 Complete Request Flow

### Example: Generate Code Request

```
1. HTTP Request
   POST /api/v1/projects/{id}/generate
   Body: {"tool": "terraform"}

2. API Layer (generation_controller.go)
   - Parse request
   - Validate input
   - Call orchestrator.GenerateCode()

3. Service Layer (orchestrator/pipeline.go)
   - Get project from repository
   - Load architecture from project
   - Validate rules
   - Call codegenService.Generate()

4. Codegen Service (services/codegen_service.go)
   - Get IaC engine from registry
   - Call engine.Generate(architecture)

5. IaC Engine (iac/terraform/engine.go)
   - Convert Architecture → Terraform blocks
   - Generate HCL code
   - Return Output{Files: [...]}

6. Response Flow
   - Output → DTO mapping
   - HTTP 200 with JSON response
```

---

## 🎯 Key Design Patterns

### 1. Registry Pattern
Dynamic registration of providers and engines:
- Cloud providers register themselves via `init()`
- IaC engines register themselves via `init()`
- No hardcoded switch statements

### 2. Repository Pattern
Data access abstraction:
- Domain defines interfaces
- Platform implements with SQL
- Services depend on interfaces

### 3. Strategy Pattern
Pluggable implementations:
- Cloud providers as strategies
- IaC engines as strategies
- Rule evaluators as strategies

### 4. Adapter Pattern
Cloud SDK adaptation:
- AWS SDK wrapped in adapters
- Domain interfaces implemented by adapters
- Domain doesn't know about AWS SDK

### 5. Inventory Pattern
Dynamic dispatch:
- Resource mappers registered in inventory
- No switch statements
- Easy to add new resources

---

## 🚫 Architecture Boundaries

### Domain Layer MUST NOT:
- ❌ Import cloud SDKs (AWS, GCP, Azure)
- ❌ Import SQL/database libraries
- ❌ Import HTTP frameworks
- ❌ Contain implementation details

### Domain Layer MUST:
- ✅ Define interfaces for external dependencies
- ✅ Contain pure business logic
- ✅ Use ubiquitous language
- ✅ Be cloud-agnostic

### API Layer MUST NOT:
- ❌ Contain business logic
- ❌ Access database directly
- ❌ Know about domain entities (only DTOs)

### Service Layer MUST NOT:
- ❌ Contain domain logic (delegate to domain)
- ❌ Know about HTTP details
- ❌ Know about SQL queries

---

## 📊 Layer Dependencies

```
cmd/ → api/ → platform/server/ → domain/
                              ↓
                         cloud/ (implements domain interfaces)
                         iac/ (implements domain interfaces)
                         platform/repository/ (implements domain interfaces)
```

**Dependency Rule**: Inner layers (domain) don't depend on outer layers. Outer layers depend on inner layers through interfaces.

---

## 🔍 How to Navigate the Codebase

### Adding a New Cloud Resource (e.g., RDS)

1. **Domain** (`internal/domain/resource/database/rds_instance.go`)
   - Define RDSInstance struct
   - Define DatabaseService interface

2. **Cloud Adapter** (`internal/cloud/aws/models/database/rds_instance.go`)
   - Define AWS RDS model

3. **Inventory** (`internal/cloud/aws/inventory/resources.go`)
   - Register RDS classification

4. **Mapper** (`internal/cloud/aws/mapper/database/rds_mapper.go`)
   - Implement Terraform mapper
   - Register in inventory

5. **Rules** (`internal/cloud/aws/rules/defaults.go`)
   - Add default constraints

6. **Database** (`migrations/`)
   - Seed resource type
   - Seed constraints

### Adding a New IaC Engine (e.g., CDK)

1. **Domain** (`internal/iac/engine.go`)
   - Engine interface already defined

2. **CDK Engine** (`internal/iac/cdk/engine.go`)
   - Implement Engine interface
   - Generate CDK code

3. **Registry** (`internal/iac/cdk/init.go`)
   - Register engine in `init()`

4. **Service** (`internal/platform/server/services/codegen_service.go`)
   - Already uses registry - no changes needed!

---

## ✅ Summary

The codebase follows a **clean, layered architecture** with:

1. **Clear separation of concerns** - Each layer has a single responsibility
2. **Dependency inversion** - Depend on interfaces, not implementations
3. **Plugin architecture** - Easy to add new providers/engines
4. **Data-driven rules** - Validation rules in database
5. **Registry pattern** - Dynamic dispatch, no switch statements
6. **Domain-first design** - Business logic isolated from infrastructure

This architecture makes the system:
- ✅ **Scalable** - Easy to add new features
- ✅ **Maintainable** - Clear boundaries and responsibilities
- ✅ **Testable** - Each layer can be tested independently
- ✅ **Extensible** - New providers/engines added without changing core
