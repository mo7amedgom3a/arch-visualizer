# Infrastructure as Code (IaC) Agent Instructions

## 🤖 Persona: DevOps Automation Engineer

You implement the engines that generate executable Infrastructure as Code (Terraform, Pulumi, CDK).

## 🎯 Goal

Translate the abstract Domain Architecture into concrete, deployable code artifacts.

## 📂 Folder Structure

- **`registry/`**: **Engine Registry**
  - Allows engines to register themselves by name ("terraform", "pulumi").
- **`terraform/`**: **Terraform Engine**
  - `engine.go`: Implements `iac.Engine`.
  - `mapper/`: Utilities for creating HCL blocks.
- **`pulumi/`**: **Pulumi Engine**
  - `engine.go`: Implements `iac.Engine`.

## 🛠️ Implementation Guide

### How to Add a New IaC Engine

1.  **Implement Interface**: Create a struct implementing `iac.Engine`.

    ```go
    type APIEngine struct {}
    func (e *APIEngine) Generate(ctx, arch) (*iac.Output, error) { ... }
    ```

2.  **Generate Logic**: Iterate through the Architecture's resources and convert them to the target language string/file.

3.  **Register**: Call `registry.Register("api", &APIEngine{})`.

## 🧪 Testing Strategy

- **Snapshot Testing**: Generate code for a known architecture and compare specific output files against expected strings.
- **Syntax Check**: (Optional) Run a linter (like `terraform validate`) on the output if possible, though usually out of scope for unit tests.
