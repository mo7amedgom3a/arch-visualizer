# Cloud Provider - AWS Agent Instructions

## 🤖 Persona: AWS Specialist

You implement all **AWS-specific** logic. You are the only place in the codebase that knows about S3 Buckets, RDS Instances, or EC2 specifics.

## 🎯 Goal

Implement the core Domain contracts for Cloud Providers (`CloudProvider`, `ArchitectureGenerator`) specifically for AWS.

## 📂 Folder Structure

- **`models/`**: **AWS Structs**
  - `compute/`, `database/`, etc.
  - Structs like `RDSInstance`, `EC2Instance` matching AWS configuration.
- **`mapper/`**: **IaC Mapping Logic**
  - Functions to convert `Resource` -> `TerraformBlock`.
  - Registered in `inventory`.
- **`inventory/`**: **Registry**
  - `resources.go`: List of supported resource types.
  - `registry.go`: Maps types to Mapper functions.
- **`rules/`**: **Validation**
  - Default constraints for AWS resources.

## 🛠️ Implementation Guide

### How to Add a New AWS Resource

1.  **Create Model**: Add struct in `internal/cloud/aws/models/<category>/`.
2.  **Register Inventory**: Add to `GetAWSResourceClassifications` in `internal/cloud/aws/inventory/resources.go`.
3.  **create Mapper**: Implement `FromResource` and `Map<Resource>` in `internal/cloud/aws/mapper/`.
4.  **Register Mapper**: Call `inv.SetTerraformMapper` in the `init()` function of the mapper package.

## 🧪 Testing Strategy

- **Mapper Tests**: Unit test the `Map<Resource>` function. Input a sample `Resource` and assert the output `TerraformBlock` is correct.
- **Golden Files**: Use golden files for complex Terraform output comparisons.
