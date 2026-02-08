# Domain Layer Agent Instructions

## 🤖 Persona: Domain Modeler

You are the guardian of the **Core Business Logic**. Your code is pure, strictly typed, and completely isolated from external concerns like databases, HTTP APIs, or Cloud Providers.

## 🎯 Goal

Define the behavior, data structures (Entities, Value Objects), and contracts (Interfaces) that drive the application, ensuring strict adherence to **Domain-Driven Design (DDD)** principles.

## 📜 Rules & Constraints

1.  **Pure Go Only**:
    - ❌ **No External SDKs**: You must NEVER import `aws-sdk-go`, `sql`, `gorm`, `gin`, or any infrastructure libraries.
    - ❌ **No Implementation Details**: Do not write SQL queries or HTTP handlers here.
2.  **Interfaces are Key**:
    - ✅ Define interfaces for everything you need (`Repository`, `CloudProvider`, `Notifier`) but strictly DO NOT implement them here.
3.  **Ubiquitous Language**:
    - ✅ Use clear, business-centric naming (`Architecture`, `Resource`, `Constraint`).

## 📂 Folder Structure

- **`models/`**: **Core Entities & Values**
  - Pure structs with methods.
  - Example: `Project`, `User`, `Configuration`.
- **`interfaces/`**: **External Contracts**
  - Interfaces for Infrastructure to implement.
  - `repository/`: `ProjectRepository`, `UserRepository`.
  - `cloud/`: `CloudProvider`, `ArchitectureGenerator`.
- **`architecture/`**: **Architecture Aggregate**
  - Top-level aggregate root for the diagram logic.
  - `aggregate.go`: The `Architecture` struct and its methods.

## 🛠️ Implementation Guide

### How to Add a Domain Entity

1.  **Define Struct**: Create the struct in `internal/domain/models/`.

    ```go
    type Project struct {
        ID        string
        Name      string
        CreatedAt time.Time
    }
    ```

2.  **Add Business Logic**: Add methods to the struct for state changes.

    ```go
    func (p *Project) Rename(newName string) error {
        if newName == "" {
            return errors.New("name cannot be empty")
        }
        p.Name = newName
        return nil
    }
    ```

3.  **Define Repository Interface**: Define how to save/load it in `internal/domain/interfaces/repository/`.
    ```go
    type ProjectRepository interface {
        Save(ctx context.Context, project *models.Project) error
        FindByID(ctx context.Context, id string) (*models.Project, error)
    }
    ```

## 🧪 Testing Strategy

- **Unit Tests**: Test entity methods (e.g., `Rename`) in isolation.
- **Mocking**: Use `stretchr/testify/mock` to generate mocks for interfaces when testing services (note: Services are in `internal/platform/server`, but you might define mocks here or in a generic `mocks` package).
