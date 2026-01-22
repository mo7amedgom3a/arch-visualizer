# Backend Architecture – Folder Structure

This backend is built using a **modular monolith architecture** in Go, designed for scalability, extensibility, and cloud-agnostic operations.

The system allows solution architects to visually compose cloud architectures and automatically generate **Infrastructure as Code (IaC)** using tools like **Terraform** and **Pulumi**. It supports multiple cloud providers such as **AWS**, **GCP**, and **Azure**.

## 🧠 Architectural Principles

- **Domain-first design** (cloud-agnostic core)
- **Strong module boundaries**
- **Provider-specific implementations**
- **Rule-driven validation** (no hardcoded logic)
- **Pluggable IaC engines**
- **Monolith today, microservices tomorrow**

## 📁 High-Level Structure

```
backend/
├── cmd/
├── internal/
├── pkg/
├── configs/
├── migrations/
├── scripts/
├── go.mod
├── go.sum
└── README.md
```

## 🔹 cmd/ – Application Entrypoints

```
cmd/
└── api/
    └── main.go
```

- Application entrypoints (API, CLI, workers, etc.)
- Bootstraps dependencies and services
- **No business logic**

## 🔹 internal/ – Core Application Code

This is where all server-side business logic lives.  
Packages under `internal/` are for internal use only and cannot be imported externally.

### 🔸 internal/platform/ – Technical Foundation

```
internal/platform/
├── config/
├── logger/
├── database/
├── server/
├── auth/
└── errors/
```

Responsibilities:
- Configuration, logging, DB connection pooling
- HTTP server setup, authentication helpers
- Common error handling
- **No business logic**

### 🔸 internal/domain/ – Cloud-Agnostic Core

```
internal/domain/
├── architecture/
├── resource/
├── constraint/
├── project/
├── errors/
└── rules/
```

**The heart of the system.**

Contains:
- Domain models and resource abstractions
- Relationships and dependencies
- Constraint interfaces
- Domain-level errors
- Cloud-agnostic rules logic

**About `rules/`:**
- The `rules` folder in `internal/domain/` defines interfaces and constructs for all rule types (e.g., constraint validators, rule engines).
- These interfaces describe *what* kind of validations or roles can exist, but are not tied to any specific cloud provider.
- This empowers any provider to use, extend, or implement these roles via their own code.
- The domain rules engine enforces constraints like: required parents, allowed parents, required region, max/min children, and more — for any resource, regardless of cloud provider.
- By keeping the interfaces in domain, we enable **all providers** (AWS, GCP, Azure, etc.) to leverage a shared, consistent rules contract.

**Key Points:**
- 🚫 No provider-specific or IaC tool code in domain
- ✅ All cloud providers implement and use these domain rules

### 🔸 internal/cloud/ – Cloud Provider Implementations

```
internal/cloud/
├── aws/
├── gcp/
└── azure/
```

Each provider follows a similar structure:

```
aws/
├── models/
│   ├── networking/
│   └── compute/
├── mapper/
├── adapters/
│   ├── networking/
│   └── compute/
├── services/
├── repositories/
└── README.md
```

Responsibilities:
- Provider-specific resource models and logic
- Mapping domain concepts to cloud resources
- Handling provider differences (e.g., AWS VPC vs GCP VPC)
- Implementing interfaces and rules defined in the domain layer

### 🔸 internal/iac/ – Infrastructure as Code Engines

```
internal/iac/
├── engine.go
├── terraform/
├── pulumi/
└── registry/
```

Purpose:
- Generate IaC from validated architectures
- Support multiple engines with a plug-in interface

Each engine is fully pluggable:

```
terraform/
├── generator/
├── templates/
├── mapper/
└── writer/
```

(Terraform today, Pulumi tomorrow, CDK in the future)

### 🔸 internal/rules/ – Legacy/Cloud-Linked Rules Logic

```
internal/rules/
├── engine/
├── constraints/
└── registry/
```

**Note:**  
With the new structure, the *domain* now owns the APIs and interfaces for constraints/rules — used by all providers and services.  
The `internal/rules/` folder supports legacy or backward-compatible logic that may interact with persistence or external data sources, but core rule definitions now live under `internal/domain/rules/`.  
All cloud providers or engines consume these domain interfaces, ensuring universal, cloud-agnostic usage.

### 🔸 internal/diagram/ – Visual Architecture Logic

```
internal/diagram/
├── graph/
├── parser/
└── validator/
```

Handles:
- Diagram/canvas parsing and validation from frontend
- Converts UI designs into validation-ready domain models

### 🔸 internal/codegen/ – Orchestration Layer

```
internal/codegen/
├── service.go
├── pipeline.go
└── result.go
```

Orchestrates the full pipeline:
1. Diagram parsing
2. Domain conversion
3. Rule/constraint validation (via domain rules)
4. Cloud provider mapping
5. IaC generation

Single entry point for **"Generate Code"**

### 🔸 internal/api/ – API/Transport Layer

```
internal/api/
├── http/
│   ├── handlers/
│   ├── middleware/
│   └── routes.go
└── dto/
```

Responsibilities:
- HTTP endpoints and routing
- Request/response DTOs
- Middleware
- Authentication & authorization
- Delegation to internal services
- **No business logic**

### 🔸 internal/persistence/ – Data Access Layer

```
internal/persistence/
└── postgres/
    ├── resource_repository.go
    ├── constraint_repository.go
    └── project_repository.go
```

Responsibilities:
- Database reads/writes
- Repository implementations
- PostgreSQL-specific logic

## 🔹 pkg/ – Reusable Public Packages

```
pkg/
```

Optionally shared helper utilities for CLI, other services, or open-source projects.

## 🔹 configs/

```
configs/
└── app.yaml
```

Application configuration files.

## 🔹 migrations/

```
migrations/
```

Database schema migrations.

## 🔹 scripts/

```
scripts/
```

Utility scripts (setup, tooling, CI helpers).

## 🗄️ Database Schema

The system uses **PostgreSQL** to store projects, resources, constraints, and their relationships. The schema is flexible and cloud-agnostic, supporting complex architectural graphs.

### Core Tables

#### **users**
Stores user authentication and profiles.

| Column     | Type      | Description                    |
|------------|-----------|--------------------------------|
| `id`       | UUID      | Primary key                    |
| `email`    | TEXT      | Unique email address           |
| `name`     | TEXT      | User's display name            |
| `created_at` | TIMESTAMP | Account creation timestamp    |

---

#### **projects**
Each project is a cloud architecture design.

| Column         | Type  | Description                                         |
|----------------|-------|-----------------------------------------------------|
| `id`           | UUID  | Primary key                                         |
| `user_id`      | UUID  | Owner (FK → `users.id`)                             |
| `infra_tool`   | INT   | IaC tool (FK → `iac_targets.id`)                    |
| `name`         | TEXT  | Project name                                        |
| `cloud_provider` | TEXT | Target cloud (`aws`, `azure`, `gcp`)               |
| `region`       | TEXT  | Default region for resources                        |
| `created_at`   | TIMESTAMP | Project creation timestamp                       |

---

#### **iac_targets**
IaC tools supported.

| Column | Type   | Description                              |
|--------|--------|------------------------------------------|
| `id`   | SERIAL | Primary key                              |
| `name` | TEXT   | IaC tool name (e.g., `terraform`, `pulumi`) |

---

### Pricing & Cost Estimates

#### **project_pricing**
Stores total pricing estimate per project.

| Column            | Type    | Description                                    |
|-------------------|---------|------------------------------------------------|
| `id`              | SERIAL  | Primary key                                    |
| `project_id`      | UUID    | Project (FK → `projects.id`)                   |
| `total_cost`      | NUMERIC | Total estimated cost                           |
| `currency`        | TEXT    | Currency (`USD`, `EUR`, `GBP`)                 |
| `period`          | TEXT    | Period (`hourly`, `monthly`, `yearly`)         |
| `duration_seconds` | BIGINT | Duration used for the estimate                 |
| `provider`        | TEXT    | Cloud provider (`aws`, `azure`, `gcp`)         |
| `region`          | TEXT    | Region (nullable)                              |
| `calculated_at`   | TIMESTAMP | Calculation timestamp                        |

---

#### **service_pricing**
Pricing per service (resource category).

| Column            | Type    | Description                                    |
|-------------------|---------|------------------------------------------------|
| `id`              | SERIAL  | Primary key                                    |
| `project_id`      | UUID    | Project (FK → `projects.id`)                   |
| `category_id`     | INT     | Category (FK → `resource_categories.id`)       |
| `total_cost`      | NUMERIC | Total estimated cost                           |
| `currency`        | TEXT    | Currency (`USD`, `EUR`, `GBP`)                 |
| `period`          | TEXT    | Period (`hourly`, `monthly`, `yearly`)         |
| `duration_seconds` | BIGINT | Duration used for the estimate                 |
| `provider`        | TEXT    | Cloud provider (`aws`, `azure`, `gcp`)         |
| `region`          | TEXT    | Region (nullable)                              |
| `calculated_at`   | TIMESTAMP | Calculation timestamp                        |

---

#### **service_type_pricing**
Pricing per service type (resource type).

| Column            | Type    | Description                                    |
|-------------------|---------|------------------------------------------------|
| `id`              | SERIAL  | Primary key                                    |
| `project_id`      | UUID    | Project (FK → `projects.id`)                   |
| `resource_type_id`| INT     | Resource type (FK → `resource_types.id`)       |
| `total_cost`      | NUMERIC | Total estimated cost                           |
| `currency`        | TEXT    | Currency (`USD`, `EUR`, `GBP`)                 |
| `period`          | TEXT    | Period (`hourly`, `monthly`, `yearly`)         |
| `duration_seconds` | BIGINT | Duration used for the estimate                 |
| `provider`        | TEXT    | Cloud provider (`aws`, `azure`, `gcp`)         |
| `region`          | TEXT    | Region (nullable)                              |
| `calculated_at`   | TIMESTAMP | Calculation timestamp                        |

---

#### **resource_pricing**
Pricing per resource instance.

| Column            | Type    | Description                                    |
|-------------------|---------|------------------------------------------------|
| `id`              | SERIAL  | Primary key                                    |
| `project_id`      | UUID    | Project (FK → `projects.id`)                   |
| `resource_id`     | UUID    | Resource (FK → `resources.id`)                 |
| `total_cost`      | NUMERIC | Total estimated cost                           |
| `currency`        | TEXT    | Currency (`USD`, `EUR`, `GBP`)                 |
| `period`          | TEXT    | Period (`hourly`, `monthly`, `yearly`)         |
| `duration_seconds` | BIGINT | Duration used for the estimate                 |
| `provider`        | TEXT    | Cloud provider (`aws`, `azure`, `gcp`)         |
| `region`          | TEXT    | Region (nullable)                              |
| `calculated_at`   | TIMESTAMP | Calculation timestamp                        |

---

#### **pricing_components**
Pricing breakdown per component (per-hour, per-GB, per-request, etc.).

| Column               | Type    | Description                                   |
|----------------------|---------|-----------------------------------------------|
| `id`                 | SERIAL  | Primary key                                   |
| `resource_pricing_id`| INT     | Resource pricing (FK → `resource_pricing.id`) |
| `component_name`     | TEXT    | Component name                                |
| `model`              | TEXT    | Pricing model (`per_hour`, `per_gb`, etc.)     |
| `unit`               | TEXT    | Unit of measure                               |
| `quantity`           | NUMERIC | Quantity used                                 |
| `unit_rate`          | NUMERIC | Rate per unit                                 |
| `subtotal`           | NUMERIC | Subtotal cost for this component              |
| `currency`           | TEXT    | Currency (`USD`, `EUR`, `GBP`)                |

---

### Resource Type System

#### **resource_categories**

| Column | Type   | Description      |
|--------|--------|-----------------|
| `id`   | SERIAL | Primary key     |
| `name` | TEXT   | Category name   |

Examples: `Compute`, `Networking`, `Storage`, `Database`, `Security`

---

#### **resource_kinds**

| Column | Type   | Description      |
|--------|--------|-----------------|
| `id`   | SERIAL | Primary key     |
| `name` | TEXT   | Kind name       |

Examples: `VirtualMachine`, `Container`, `Function`, `Network`, `LoadBalancer`

---

#### **resource_types**
Cloud-specific resource implementations.

| Column          | Type    | Description                                  |
|-----------------|---------|----------------------------------------------|
| `id`            | SERIAL  | Primary key                                  |
| `name`          | TEXT    | Resource type name                           |
| `cloud_provider`| TEXT    | Cloud provider (`aws`, `azure`, `gcp`)       |
| `category_id`   | INT     | FK → `resource_categories.id`                |
| `kind_id`       | INT     | FK → `resource_kinds.id`                     |
| `is_regional`   | BOOLEAN | Whether resource is region-specific          |
| `is_global`     | BOOLEAN | Whether resource is global                   |

This mapping enables cloud-agnostic handling at the domain level.

---

### Architecture Graph

#### **resources**
Resource instances in a project.

| Column            | Type    | Description                                 |
|-------------------|---------|---------------------------------------------|
| `id`              | UUID    | Primary key                                 |
| `project_id`      | UUID    | Parent project (FK → `projects.id`)         |
| `resource_type_id`| INT     | Resource type (FK → `resource_types.id`)    |
| `name`            | TEXT    | User-defined name                           |
| `pos_x`           | INT     | X coordinate on canvas                      |
| `pos_y`           | INT     | Y coordinate on canvas                      |
| `config`          | JSONB   | Resource-specific config (CIDR, tags, etc.) |
| `created_at`      | TIMESTAMP | Creation timestamp                        |

---

#### **resource_containment**
Parent-child relationships (VPC → Subnet → EC2, etc.).

| Column            | Type    | Description                                 |
|-------------------|---------|---------------------------------------------|
| `parent_resource_id` | UUID | Parent  (FK → `resources.id`)               |
| `child_resource_id`  | UUID | Contained resource (FK → `resources.id`)    |

---

#### **dependency_types**
Types of resource dependencies.

| Column | Type   | Description      |
|--------|--------|-----------------|
| `id`   | SERIAL | Primary key     |
| `name` | TEXT   | Dependency type |

E.g. `uses`, `depends_on`, `connects_to`, `references`

---

#### **resource_dependencies**
Directed dependencies between resources.

| Column           | Type   | Description                                   |
|------------------|--------|-----------------------------------------------|
| `from_resource_id` | UUID | Source      (FK → `resources.id`)             |
| `to_resource_id`   | UUID | Target      (FK → `resources.id`)             |
| `dependency_type_id` | INT | Dependency type (FK → `dependency_types.id`) |

---

### Rules & Constraints

#### **resource_constraints**
Database-driven validation rules and constraints.

| Column            | Type   | Description                                               |
|-------------------|--------|----------------------------------------------------------|
| `id`              | SERIAL | Primary key                                              |
| `resource_type_id`| INT    | Applies to resource type (FK → `resource_types.id`)      |
| `constraint_type` | TEXT   | Type of constraint (interface in domain/rules/)          |
| `constraint_value`| TEXT   | The constraint value itself                              |

**Constraint Types:**  
(defined as interfaces in `internal/domain/rules/`, implemented by any provider)

- `requires_parent`   → Must be inside another resource
- `allowed_parent`    → Only specific parents allowed
- `requires_region`   → Must specify a region
- `max_children`      → Max number of children
- `min_children`      → Min number of children
- `allowed_dependencies` → Valid dependency types

Any cloud provider can use these roles via implementing the interfaces from the domain rules contract.

**Examples:**
```sql
-- EC2 must be inside a subnet
(resource_type_id: EC2, constraint_type: 'requires_parent', constraint_value: 'Subnet')

-- Subnet must be inside a VPC
(resource_type_id: Subnet, constraint_type: 'allowed_parent', constraint_value: 'VPC')

-- S3 Bucket is global
(resource_type_id: S3, constraint_type: 'requires_region', constraint_value: 'false')
```
This enables **rule-driven, provider-independent validation** with no logic hardcoded to a specific cloud.

---

### 🔗 Relationships Overview

```
users (1) ──→ (N) projects
projects (1) ──→ (N) resources
resources (N) ──→ (1) resource_types
resource_types (N) ──→ (1) resource_categories
resource_types (N) ──→ (1) resource_kinds
resource_types (1) ──→ (N) resource_constraints

resources (parent) ──→ (children) resource_containment
resources (from) ──→ (to) resource_dependencies
resource_dependencies (N) ──→ (1) dependency_types
projects (N) ──→ (1) iac_targets
projects (1) ──→ (N) project_pricing
projects (1) ──→ (N) service_pricing
projects (1) ──→ (N) service_type_pricing
projects (1) ──→ (N) resource_pricing
resource_categories (1) ──→ (N) service_pricing
resource_types (1) ──→ (N) service_type_pricing
resources (1) ──→ (N) resource_pricing
resource_pricing (1) ──→ (N) pricing_components
```

---

### 📊 Data Flow Example

1. User creates a **project** (`AWS`, `us-east-1`, `Terraform`)
2. User adds resources via UI:
   - VPC (`10.0.0.0/16`)
   - Subnet (`10.0.1.0/24`, inside VPC)
   - EC2 (inside Subnet)
3. System validates all constraints using domain rules:
   - ✅ Subnet requires VPC as parent
   - ✅ EC2 requires Subnet as parent
4. Domain logic maps to provider types:
   - `VPC` → `aws/VPC`
   - `Subnet` → `aws/Subnet`
   - `VirtualMachine` → `aws/EC2Instance`
5. IaC engine generates **Terraform code**

---

## 🧩 Mental Model

- **Domain** → What the architecture means (interfaces/roles for all)
- **Rules** → What is allowed (enforced by domain/rules interfaces)
- **Cloud** → How providers differ (cloud-specific implementations)
- **IaC** → How code is generated
- **Diagram** → How users design
- **API** → How the outside world interacts

## 🚀 Scalability & Future-Proofing

This structure supports:
- ✅ Adding new cloud providers
- ✅ Adding/plugging in new IaC tools
- ✅ Migrating to microservices later
- ✅ Centralized, role/interface-driven validation
- ✅ Clean separation of concerns, cloud-agnostic at core

