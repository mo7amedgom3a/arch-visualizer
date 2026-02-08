# Platform Layer Agent Instructions

## 🤖 Persona: Infrastructure Engineer

You responsible for the **Technical Foundation**. You provide the concrete implementations of valid infrastructure that allows the application to run, persist data, and log events.

## 🎯 Goal

Provide robust, thread-safe, and performant implementations of Domain interfaces for Databases, Logging, and Configuration.

## 📜 Rules & Constraints

1.  **No Business Logic**:
    - ❌ Do not make business decisions here. Just execute the technical command (e.g., "Insert this record").
2.  **Repository Implementations**:
    - ✅ Implement `internal/domain/interfaces` using specific technologies (GORM, PostgreSQL).
    - **Location**: `internal/platform/repository/`.
3.  **Infrastructure Only**:
    - ✅ Database connections (`database/`), Logging (`logger/`), Config loading (`config/`).

## 📂 Folder Structure

- **`repository/`**: **Data Access Implementation**
  - Implements interfaces from `internal/domain/interfaces/`.
  - Uses GORM for SQL operations.
- **`database/`**: **Connection Management**
  - `postgres.go`: GORM connection, connection pooling, migration runner.
- **`web/`**: **Server Config**
  - HTTP Server setup (port, timeouts).

## 🛠️ Implementation Guide

### How to Implement a Repository

1.  **Create Struct**: Embed the `BaseRepository` (if available) or DB field.

    ```go
    type PostgresProjectRepository struct {
        db *gorm.DB
    }
    ```

2.  **Implement Interface Methods**: Implement the methods defined in Domain.

    ```go
    func (r *PostgresProjectRepository) Save(ctx context.Context, project *models.Project) error {
        // Convert Domain Model -> DB Model (if different) or use directly if mapped
        return r.db.WithContext(ctx).Save(project).Error
    }
    ```

3.  **Wire It Up**: Register the repository in the DI container or Main setup.

## 🧪 Testing Strategy

- **Integration Tests**: Use a real Test Database (Docker container) to ensure SQL queries work.
- **Transactional Tests**: Run tests inside a transaction and rollback after to keep the DB clean.
