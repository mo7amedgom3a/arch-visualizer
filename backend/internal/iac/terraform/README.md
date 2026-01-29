# `internal/iac/terraform/` — Terraform Engine (HCL Generation)

This package implements the **Terraform IaC engine** using a small Terraform “IR” and an `hclwrite`-based writer.

## Pipeline

Terraform codegen follows:

```
Domain Architecture (validated)
        ↓
Domain Graph / DAG builder (`architecture.NewGraph`)
        ↓
Topological Sort (`graph.GetSortedResources`)
        ↓
Terraform Engine (`internal/iac/terraform/generator`)
        ↓
HCL Writer (`internal/iac/terraform/writer`)
        ↓
main.tf (today), variables.tf/outputs.tf (optional in future)
```

## Packages

### `mapper/`

Defines the Terraform intermediate representation and the mapper registry:

- `TerraformBlock`: a generic HCL block (`Kind`, `Labels`, `Attributes`)
- `TerraformValue`: supports string/number/bool/list/map and **expressions** via `TerraformExpr`
- `ResourceMapper`: provider-specific mapping interface:
  - `Provider() string`
  - `SupportsResource(resourceType string) bool`
  - `MapResource(*domainResource) ([]TerraformBlock, error)`
- `MapperRegistry`: provider → mapper registration/lookup

This boundary is what keeps the Terraform engine **cloud-agnostic**.

### `generator/`

Implements the Terraform `iac.Engine`:

- Selects a mapper by `arch.Provider`
- Emits a `provider "<provider>" { region = "<arch.Region>" }` block (if region is set)
- Maps resources **in topological order** into Terraform blocks
- Uses `writer.RenderMainTF` to produce `main.tf`

See: [`generator/generator.go`](generator/generator.go)

### `writer/`

Renders Terraform blocks into formatted HCL using:

- `github.com/hashicorp/hcl/v2/hclwrite`
- `github.com/hashicorp/hcl/v2/hclsyntax`

Expressions (like `aws_vpc.vpc_1.id`) are rendered as traversals, not quoted strings.

See: [`writer/writer.go`](writer/writer.go)

## Current provider implementation

AWS provider mapping is implemented in:

- `internal/cloud/aws/mapper/terraform/mapper.go`

It maps domain resource types (e.g. `VPC`, `Subnet`, `EC2`) into Terraform resources (e.g. `aws_vpc`, `aws_subnet`, `aws_instance`) and uses the **domain resource ID** (sanitized) as Terraform local name to make references stable.

## Extending

- **Add a new provider**: implement `mapper.ResourceMapper`, then register it in the mapper registry used by the Terraform engine.
- **Add a new resource type**: extend the provider mapper’s `SupportsResource` and `MapResource`.

