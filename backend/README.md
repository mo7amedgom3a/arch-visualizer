# Backend Architecture – Folder Structure

This backend is built using a **modular monolith architecture** in Go, designed to be scalable, extensible, and cloud-agnostic.

The system enables solution architects to visually design cloud architectures and automatically generate **Infrastructure as Code (IaC)** using tools like **Terraform** and **Pulumi**, across multiple cloud providers such as **AWS**, **GCP**, and **Azure**.

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

- Contains application entrypoints
- Responsible for bootstrapping the app
- Wires dependencies together
- **No business logic**

Future entrypoints may include:
- CLI
- Workers
- Background jobs

## 🔹 internal/ – Core Application Code

All internal business logic lives here.  
Packages under `internal/` cannot be imported externally, enforcing boundaries.

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

Cross-cutting technical concerns:
- Configuration loading
- Logging
- Database connections
- HTTP server setup
- Authentication helpers
- Common error handling

**No business logic here.**

### 🔸 internal/domain/ – Cloud-Agnostic Core

```
internal/domain/
├── architecture/
├── resource/
├── constraint/
├── project/
└── errors/
```

**This is the heart of the system.**

Contains:
- Architecture models
- Resource definitions
- Relationships and dependencies
- Constraint abstractions
- Domain-level errors

- 🚫 **No AWS, GCP, Terraform, or Pulumi code**
- ✅ **Pure business logic**

### 🔸 internal/cloud/ – Cloud Provider Implementations

```
internal/cloud/
├── aws/
├── gcp/
└── azure/
```

Each provider follows the same structure:

```
aws/
├── models/
├── services/
├── repositories/
└── mapper/
```

**Responsibilities:**
- Provider-specific resource models
- Mapping domain resources to cloud resources
- Handling provider differences (e.g. AWS VPC vs GCP VPC)

### 🔸 internal/iac/ – Infrastructure as Code Engines

```
internal/iac/
├── engine.go
├── terraform/
├── pulumi/
└── registry/
```

**Purpose:**
- Generate IaC from validated architecture
- Support multiple engines via a common interface

Each engine is fully pluggable:

```
terraform/
├── generator/
├── templates/
├── mapper/
└── writer/
```

This allows:
- Terraform today
- Pulumi tomorrow
- CDK in the future

### 🔸 internal/rules/ – Rules & Constraints Engine

```
internal/rules/
├── engine/
├── constraints/
└── registry/
```

**Responsibilities:**
- Load constraints from the database
- Validate architectures before code generation
- Enforce rules like:
  - `requires_parent`
  - `allowed_parent`
  - `requires_region`
  - `max_children`

This replaces hardcoded cloud rules with **data-driven validation**.

### 🔸 internal/diagram/ – Visual Architecture Logic

```
internal/diagram/
├── graph/
├── parser/
└── validator/
```

**Handles:**
- Parsing diagram/canvas JSON from frontend
- Building an internal graph
- Validating structural correctness
- Preparing domain-ready architectures

### 🔸 internal/codegen/ – Orchestration Layer

```
internal/codegen/
├── service.go
├── pipeline.go
└── result.go
```

**Coordinates the full pipeline:**
1. Diagram parsing
2. Domain conversion
3. Rules validation
4. Cloud provider mapping
5. IaC generation

This is the single entry point for **"Generate Code"**.

### 🔸 internal/api/ – Transport Layer

```
internal/api/
├── http/
│   ├── handlers/
│   ├── middleware/
│   └── routes.go
└── dto/
```

**Responsibilities:**
- HTTP routing
- Request/response DTOs
- Middleware
- Authentication & authorization

- 🚫 **No business logic**
- ✅ **Delegates to internal services**

### 🔸 internal/persistence/ – Data Access Layer

```
internal/persistence/
└── postgres/
    ├── resource_repository.go
    ├── constraint_repository.go
    └── project_repository.go
```

**Responsibilities:**
- Database access
- Repository implementations
- PostgreSQL-specific logic

## 🔹 pkg/ – Reusable Public Packages

```
pkg/
```

Optional shared utilities that may be reused by:
- Other services
- CLI tools
- External projects

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

The system uses **PostgreSQL** to store projects, resources, constraints, and relationships. The schema is designed to be flexible, cloud-agnostic, and support complex architectural graphs.

### Core Tables

#### **users**
Stores user authentication and profile information.

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `email` | TEXT | Unique email address |
| `name` | TEXT | User's display name |
| `created_at` | TIMESTAMP | Account creation timestamp |

---

#### **projects**
Each project represents a cloud architecture design.

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `user_id` | UUID | Owner of the project (FK → `users.id`) |
| `infra_tool` | INT | IaC tool to use (FK → `iac_targets.id`) |
| `name` | TEXT | Project name |
| `cloud_provider` | TEXT | Target cloud (`aws`, `azure`, `gcp`) |
| `region` | TEXT | Default region for resources |
| `created_at` | TIMESTAMP | Project creation timestamp |

---

#### **iac_targets**
Defines available Infrastructure as Code tools.

| Column | Type | Description |
|--------|------|-------------|
| `id` | SERIAL | Primary key |
| `name` | TEXT | IaC tool name (e.g., `terraform`, `pulumi`) |

Examples: `Terraform`, `Pulumi`, `CDK`

---

### Resource Type System

#### **resource_categories**
High-level categories of cloud resources.

| Column | Type | Description |
|--------|------|-------------|
| `id` | SERIAL | Primary key |
| `name` | TEXT | Category name |

Examples: `Compute`, `Networking`, `Storage`, `Database`, `Security`

---

#### **resource_kinds**
Defines the kind of resource.

| Column | Type | Description |
|--------|------|-------------|
| `id` | SERIAL | Primary key |
| `name` | TEXT | Kind name |

Examples: `VirtualMachine`, `Container`, `Function`, `Network`, `LoadBalancer`

---

#### **resource_types**
Cloud-specific resource implementations.

| Column | Type | Description |
|--------|------|-------------|
| `id` | SERIAL | Primary key |
| `name` | TEXT | Resource type name |
| `cloud_provider` | TEXT | Cloud provider (`aws`, `azure`, `gcp`) |
| `category_id` | INT | FK → `resource_categories.id` |
| `kind_id` | INT | FK → `resource_kinds.id` |
| `is_regional` | BOOLEAN | Whether resource is region-specific |
| `is_global` | BOOLEAN | Whether resource is global |

**Examples:**
- `aws` / `EC2Instance` → Category: `Compute`, Kind: `VirtualMachine`
- `gcp` / `GCE Instance` → Category: `Compute`, Kind: `VirtualMachine`
- `azure` / `Virtual Machine` → Category: `Compute`, Kind: `VirtualMachine`

This enables **cloud-agnostic mapping** at the domain level.

---

### Architecture Graph

#### **resources**
Actual resource instances within a project.

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key |
| `project_id` | UUID | Parent project (FK → `projects.id`) |
| `resource_type_id` | INT | Type of resource (FK → `resource_types.id`) |
| `name` | TEXT | User-defined resource name |
| `pos_x` | INT | X position on canvas |
| `pos_y` | INT | Y position on canvas |
| `config` | JSONB | Resource-specific configuration (CIDR, tags, instance type, etc.) |
| `created_at` | TIMESTAMP | Creation timestamp |

**Visual positioning** (`pos_x`, `pos_y`) enables diagram persistence and recreation.

---

#### **resource_containment**
Parent-child relationships (e.g., VPC → Subnet → EC2).

| Column | Type | Description |
|--------|------|-------------|
| `parent_resource_id` | UUID | Container resource (FK → `resources.id`) |
| `child_resource_id` | UUID | Contained resource (FK → `resources.id`) |

**Examples:**
- VPC contains Subnets
- Subnet contains EC2 instances
- Kubernetes Cluster contains Pods

---

#### **dependency_types**
Types of dependencies between resources.

| Column | Type | Description |
|--------|------|-------------|
| `id` | SERIAL | Primary key |
| `name` | TEXT | Dependency type name |

Examples: `uses`, `depends_on`, `connects_to`, `references`

---

#### **resource_dependencies**
Directed graph of resource dependencies.

| Column | Type | Description |
|--------|------|-------------|
| `from_resource_id` | UUID | Source resource (FK → `resources.id`) |
| `to_resource_id` | UUID | Target resource (FK → `resources.id`) |
| `dependency_type_id` | INT | Type of dependency (FK → `dependency_types.id`) |

**Examples:**
- `EC2 Instance` → `depends_on` → `Security Group`
- `Lambda Function` → `connects_to` → `DynamoDB Table`
- `Load Balancer` → `uses` → `Target Group`

---

### Rules & Constraints

#### **resource_constraints**
Database-driven validation rules.

| Column | Type | Description |
|--------|------|-------------|
| `id` | SERIAL | Primary key |
| `resource_type_id` | INT | Resource type this applies to (FK → `resource_types.id`) |
| `constraint_type` | TEXT | Type of constraint |
| `constraint_value` | TEXT | Constraint value/rule |

**Constraint Types:**
- `requires_parent` → Must be inside another resource
- `allowed_parent` → Only specific parents allowed
- `requires_region` → Must have a region set
- `max_children` → Maximum number of child resources
- `min_children` → Minimum number of child resources
- `allowed_dependencies` → Valid dependency types

**Examples:**
```sql
-- EC2 must be inside a subnet
(resource_type_id: EC2, constraint_type: 'requires_parent', constraint_value: 'Subnet')

-- Subnet must be inside a VPC
(resource_type_id: Subnet, constraint_type: 'allowed_parent', constraint_value: 'VPC')

-- S3 Bucket is global
(resource_type_id: S3, constraint_type: 'requires_region', constraint_value: 'false')
```

This enables **rule-driven validation** without hardcoding logic into the application.

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
```

---

### 📊 Data Flow Example

1. User creates a **project** (`AWS`, `us-east-1`, `Terraform`)
2. User drags resources onto the canvas:
   - VPC (`10.0.0.0/16`)
   - Subnet (`10.0.1.0/24`) inside VPC
   - EC2 instance inside Subnet
3. System validates via **constraints**:
   - ✅ Subnet has VPC as parent
   - ✅ EC2 has Subnet as parent
4. System maps to **cloud-specific types**:
   - `VPC` → `aws/VPC`
   - `Subnet` → `aws/Subnet`
   - `VirtualMachine` → `aws/EC2Instance`
5. IaC engine generates **Terraform code**

---

## 🧩 Mental Model

- **Domain** → What the architecture means
- **Rules** → What is allowed
- **Cloud** → How providers differ
- **IaC** → How code is generated
- **Diagram** → How users design
- **API** → How the outside world talks to the system

## 🚀 Scalability & Future-Proofing

This structure supports:
- ✅ Adding new cloud providers
- ✅ Adding new IaC tools
- ✅ Migrating to microservices
- ✅ Rule-driven extensibility
- ✅ Clean separation of concerns
