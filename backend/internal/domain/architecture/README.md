## Architecture Domain (Unified Dependency Graph)

This package contains the core domain model that represents a **logical cloud architecture** built from a diagram, plus helpers to **query**, **validate**, and **order** resources by their dependencies.

It is the bridge between the visual diagram (`internal/diagram`) and downstream engines (rules, pricing, IaC codegen).

---

## Files Overview

- `aggregate.go`  
  Defines the main `Architecture` aggregate and mapping from the normalized diagram graph.

- `graph.go`  
  Thin wrapper (`Graph`) around `Architecture` that provides convenient graph-style queries.

- `toposort.go`  
  Implements topological sorting over the unified dependency graph to produce a safe provisioning order.

- `toposort_test.go`  
  Unit tests that cover typical AWS scenarios, containment, mixed dependencies, edge cases, and cycle detection.

- `validation.go`  
  Placeholder for domain-level validation on an `Architecture` (e.g. parent existence).

---

## Architecture Aggregate (`aggregate.go`)

`Architecture` is the in-memory representation of a project’s intent:

- `Resources []*resource.Resource`  
  All logical resources (VPC, Subnet, EC2, etc.) derived from the diagram.

- `Region string` / `Provider resource.CloudProvider`  
  High-level context extracted from the diagram (usually from a region node).

- `Containments map[string][]string`  
  Parent → children relationships (e.g. `vpc-1 -> [subnet-1, subnet-2]`).

- `Dependencies map[string][]string`  
  Resource → dependencies it **depends on** (e.g. `ec2-1 -> [subnet-1, sg-1]`).

The helper `MapDiagramToArchitecture(...)`:
- Takes a normalized `DiagramGraph`.
- Skips visual-only and region-only nodes.
- Maps IR resource types (e.g. `"vpc"`, `"subnet"`) to domain `ResourceType`.
- Builds:
  - `Resources`
  - `Containments` from node parenthood
  - `Dependencies` from dependency edges

This is the **single source of truth** for logical relationships between resources.

---

## Graph Helper (`graph.go`)

`Graph` is a light wrapper around `*Architecture` that exposes query methods:

- `GetResource(id string) (*resource.Resource, bool)`  
- `GetChildren(parentID string) []*resource.Resource`  
- `GetParent(childID string) (*resource.Resource, bool)`  
- `GetDependencies(resourceID string) []*resource.Resource`  
- `GetRootResources() []*resource.Resource`  
- `BuildContainmentTree() map[string][]*resource.Resource`

These operate purely in-memory and are used by:
- Rules engine (for relationship-aware validations).
- Pricing and analysis features.
- Topological sorting (see below).

---

## Topological Sorting (`toposort.go`)

The goal is to derive a **logical provisioning order** for resources, such that:
- **All dependencies are created before their dependents**, and
- **Parents are created before children** (containment).

This file exposes:

- `type TopologicalSortResult struct`  
  - `Resources []*resource.Resource` – sorted list in provisioning order.  
  - `Levels [][]string` – IDs grouped by dependency “layer”.  
  - `HasCycle bool` – whether a cycle was detected.  
  - `CycleInfo []string` – IDs that participate in a cycle (if any).

- `func (g *Graph) TopologicalSort() (*TopologicalSortResult, error)`  
  - Uses **Kahn's algorithm** (BFS-based topological sort).
  - Combines:
    - `Dependencies` map (explicit “depends_on” edges).
    - `Containments` map (parent → child).
    - `ParentID` on each `resource.Resource` as a fallback.
  - Builds an adjacency list where **edges point from dependency to dependent**:
    - If `subnet-1` depends on `vpc-1`, edge is `vpc-1 -> subnet-1`.
    - If `ec2-1` is contained in `subnet-1`, edge is `subnet-1 -> ec2-1`.
  - Computes in-degrees and repeatedly removes nodes with in-degree `0`,
    producing:
    - A globally ordered `Resources` slice.
    - `Levels` that capture each “wave” of zero-in-degree nodes.
  - Detects cycles when not all nodes can be processed.

- `func (g *Graph) GetSortedResources() ([]*resource.Resource, error)`  
  - Convenience wrapper over `TopologicalSort`.  
  - Returns only the sorted resources.  
  - Fails with an error if `HasCycle` is `true`.

**Example (AWS networking / compute):**

Expected high-level order:

1. VPC  
2. Subnets  
3. Internet Gateway / NAT Gateway  
4. Route Tables  
5. Security Groups  
6. EC2 instances

The exact order within each level can vary, but **no resource is scheduled before its prerequisites**.

---

## Validation (`validation.go`)

Currently `Validate()` performs minimal domain checks:

- Ensures parent references (when present) point to some existing resource.  
- Treats missing parents as soft (to allow region-level parents not modeled as resources).

This file is the extensibility point for **cross-resource invariants** that are
not specific to any single cloud provider (provider-specific rules live under
`internal/cloud/<provider>/rules`).

---

## Typical Usage

High-level flow from diagram to ordered resources:

1. Parse & normalize diagram → `DiagramGraph` (`internal/diagram`).
2. `MapDiagramToArchitecture(diagramGraph, provider)` → `*Architecture`.
3. Wrap aggregate: `g := NewGraph(arch)`.
4. Get provisioning order:

   ```go
   sorted, err := g.GetSortedResources()
   if err != nil {
       // handle circular dependencies
   }
   // use sorted for IaC codegen / planning / previews
   ```

This keeps **diagram parsing**, **domain modeling**, and **execution planning**
cleanly separated while still giving downstream components (Terraform, Pulumi,
rules, pricing) a unified, dependency-aware view of the architecture.

