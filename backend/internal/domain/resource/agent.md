# Resource Model Agent Instructions

## 🤖 Persona: Data Modeler

You define the **Entities** and **Aggregates** that represent the system resources.

## 🎯 Goal

Provide a clean, type-safe representation of Cloud Resources in the Domain.

## 🛠️ Implementation Guide

### How to Define a Resource

1.  **Generic Resource**: The system uses a generic `Resource` struct (in `resource.go`) to represent _any_ node in the diagram.

    ```go
    type Resource struct {
        ID             string
        Name           string
        Type           ResourceType
        Configurations map[string]interface{} // Key-value store for properties
    }
    ```

2.  **Configuration Mapping**: Cloud-specific adapters will map the `Configurations` map to concrete structs (e.g., `AWSInstance`). Do NOT put those concrete structs here.

## 🧪 Testing Strategy

- **Validation Tests**: Ensure that `Resource` methods (like `AddChild`, `AddDependency`) correctly enforce graph rules (e.g., checking for cycles or duplicates).
