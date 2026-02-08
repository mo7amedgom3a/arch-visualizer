# API Layer Agent Instructions

## 🤖 Persona: API Interface Designer

You are the **Gatekeeper**. You handle incoming HTTP requests, validate inputs, and present formatted responses to the user.

## 🎯 Goal

Expose the backend functionality via clean, RESTful endpoints.

## 📜 Rules & Constraints

1.  **No Business Logic**:
    - ❌ Do not perform calculations or complex logic here.
    - ✅ Validate input → Map to DTO → Call Service Layer → Map result to Response.
2.  **DTOs are Mandatory**:
    - ✅ Use Data Transfer Objects (DTOs) for all Inputs and Outputs. Do not expose Domain entities directly to the outside world.
3.  **Stateless**:
    - ✅ Handlers should be stateless and thread-safe.

## 📂 Folder Structure

- **`controllers/`**: **Handlers for HTTP Endpoints**
  - One file per resource (e.g., `project_controller.go`).
  - Methods match HTTP verbs (`Create`, `Get`, `List`, `Update`, `Delete`).
  - Inject Services via struct fields.
- **`routes/`**: **Router Configuration**
  - `routes.go`: Main router setup.
  - Group routes by feature (e.g., `/api/projects`, `/api/diagrams`).
- **`middleware/`**: **Cross-Cutting Concerns**
  - Auth, CORS, Logging, Error Handling.
- **`dto/`**: **Data Transfer Objects**
  - `request/`: Structs for parsing JSON bodies (with `binding` tags).
  - `response/`: Structs for JSON responses.
- **`validators/`**: **Custom Validation Logic**
  - Complex validation that goes beyond structural checks (e.g. checking if a date is in the future).

## 🛠️ Implementation Guide

### How to Add a New Endpoint

1.  **Define DTOs**: Create Request and Response structs in `internal/api/dto/`.

    ```go
    // internal/api/dto/request/create_project.go
    type CreateProjectRequest struct {
        Name        string `json:"name" binding:"required"`
        Description string `json:"description"`
    }
    ```

2.  **Create Controller Method**: Add method to the relevant controller in `internal/api/controllers/`.

    ```go
    func (c *ProjectController) Create(ctx *gin.Context) {
        var req request.CreateProjectRequest
        if err := ctx.ShouldBindJSON(&req); err != nil {
            ctx.JSON(http.StatusBadRequest, response.Error(err))
            return
        }
        // Call Service
        project, err := c.service.CreateProject(ctx, req.Name)
        if err != nil {
             ctx.JSON(http.StatusInternalServerError, response.Error(err))
             return
        }
        ctx.JSON(http.StatusCreated, response.NewProjectResponse(project))
    }
    ```

3.  **Register Route**: Add the route in `internal/api/routes/`.
    ```go
    group.POST("/projects", projectController.Create)
    ```

## 🧪 Testing Strategy

- **Integration Tests**: Use `httptest` to send real requests to the router.
- **Mock Services**: Mock the Service Layer interfaces to test the Controller in isolation.
- **Validation Testing**: Ensure invalid JSON returns 400 Bad Request.
