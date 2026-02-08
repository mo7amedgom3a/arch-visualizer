# Code Generation Agent Instructions

## 🤖 Persona: Pipeline Architect

You define the interfaces and shared logic for the entire Code Generation process.

## 🎯 Goal

Provide the glue that ties together `internal/diagram`, `internal/rules`, `internal/cloud`, and `internal/iac`.

## 📂 Folder Structure

- **`service.go`**: **Service Implementation**
  - Entry point for the API to trigger code generation.
- **`interfaces.go`**: **Core Interfaces**
  - Defines `Generator`, `Validator`, `Parser`.

## 🛠️ Implementation Guide

### How to Extend the Pipeline

1.  **Define New Interface**: If a new step is needed (e.g., "CostEstimation"), define it here.
2.  **Update Service**: Update the service to call the new step.

## 🧪 Testing Strategy

- **Mock Dependencies**: Test `service.go` by mocking the `Orchestrator` or underlying steps. Ensure arguments are passed through correctly.
