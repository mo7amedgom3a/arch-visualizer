# `internal/codegen/terraform/` — Terraform Compiler

This package is the **Terraform-specific orchestration layer**.

## What it does

`Compiler.Compile(ctx, arch)` performs:

1. **DAG builder**: wraps the domain architecture using `architecture.NewGraph(arch)`
2. **Topological sort**: `graph.GetSortedResources()` (dependencies first)
3. **Terraform compilation**: calls the injected Terraform `iac.Engine`

See: [`compiler.go`](compiler.go)

## What it does NOT do

- It does **not** perform validation/rules evaluation.
  - The caller must ensure the architecture has already passed all validations and rules.
- It does **not** contain cloud-provider logic.
  - Provider mapping is done inside the Terraform engine via provider mappers (e.g. AWS mapper).

## Composition example (wiring)

At application bootstrap time you typically:

- Create a Terraform mapper registry and register provider mappers (e.g. AWS)
- Create the Terraform engine with that mapper registry
- Create this compiler with the Terraform engine

Then runtime code calls `Compile()` with a validated architecture.

