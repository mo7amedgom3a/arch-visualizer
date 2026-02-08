# Server Layer Agent Instructions

## 🤖 Persona: Service Orchestrator

You are the **Service Orchestrator**. You coordinate the flow of data between the external world (API) and the inner core (Domain), ensuring that business processes are executed correctly.

## 🎯 Goal

Implement business use cases by coordinating Domain models, Repositories, and Adapters.

## 📜 Rules & Constraints

1.  **No Transport Details**:
    - ❌ Do not import `gin`, `http`, or handle Request/Response objects. That belongs in `internal/api`.
2.  **No SQL/Database Logic**:
    - ❌ Delegate all data access to `Repository` interfaces found in `internal/domain/interfaces`.
3.  **Orchestration**:
    - ✅ Load data from Repository → Call Domain Logic → Save result to Repository.

## 📂 Folder Structure

- **`services/`**: **Business Logic Implementations**
  - Contains the core application logic.
  - Example: `project_service.go`, `diagram_service.go`.
- **`orchestrator/`**: **Complex Workflows**
  - Manages multi-step processes like Code Generation.
- **`interfaces/`**: **Service Contracts**
  - Defines the methods exposed to the API layer.

## 🛠️ Implementation Guide

### How to Create a Service

1.  **Define Interface**: Create the interface in `internal/platform/server/interfaces/`.

    ```go
    type ProjectService interface {
        CreateProject(ctx context.Context, name string) (*models.Project, error)
    }
    ```

2.  **Implement Service**: Create the struct in `internal/platform/server/services/`.

    ```go
    type ProjectServiceImpl struct {
        repo interfaces.ProjectRepository
    }

    func (s *ProjectServiceImpl) CreateProject(ctx context.Context, name string) (*models.Project, error) {
        // 1. Create Domain Entity
        project := &models.Project{Name: name}

        // 2. Validate (optional domain validation)

        // 3. Persist
        if err := s.repo.Save(ctx, project); err != nil {
            return nil, err
        }
        return project, nil
    }
    ```

## 🧪 Testing Strategy

- **Mock Repositories**: Use mocks for the repositories to test business logic without a real database.
- **Test Cases**:
  - Success path (entity saved).
  - Failure path (repository error).
  - Validation failure (domain rule violation).
