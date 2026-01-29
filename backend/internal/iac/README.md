# `internal/iac/` — IaC Engine Abstractions

This package defines the **pluggable IaC engine interface** used by code generation.

## What lives here

- **Engine contract**: [`engine.go`](engine.go)
  - `iac.Engine` exposes:
    - `Name() string`
    - `Generate(ctx, arch, sortedResources) (*iac.Output, error)`
  - `iac.Output` is a list of generated files (`Path`, `Content`, `Type`).

- **Engine registry**: [`registry/registry.go`](registry/registry.go)
  - `registry.EngineRegistry` allows selecting an engine by name (e.g. `"terraform"`),
    enabling **loose coupling** between the orchestration layer and engine implementations.

## Key design principle

- **No cloud/provider hardcoding here**. Engines are selected by name, and provider-specific
  logic belongs in provider packages (e.g. AWS) behind mapper interfaces.

## Typical usage

At runtime, the application composes:

- a Terraform mapper registry (provider → mapper)
- a Terraform engine (implements `iac.Engine`)
- an engine registry (engine name → engine instance)

Then codegen chooses `"terraform"` (or other engines) without importing engine code directly.

