# API Endpoints & Usage

This document lists the available API endpoints for the Arch Visualizer Backend.

**Base URL**: `http://localhost:9000/api/v1`
**Swagger UI**: `http://localhost:9000/swagger/index.html`

## Static Data

### List Cloud Providers

Get a list of supported cloud providers.

- **Method**: `GET`
- **Path**: `/static/providers`

**Curl commands**:

```bash
curl -X GET "http://localhost:9000/api/v1/static/providers"
```

### List Resource Types

Get a list of resource types, optionally filtered by provider.

- **Method**: `GET`
- **Path**: `/static/resource-types`
- **Query Params**:
  - `provider` (optional): e.g., `aws`, `gcp`

**Curl commands**:

````bash
# List all
curl -X GET "http://localhost:9000/api/v1/static/resource-types"

```bash
# Filter by AWS
curl -X GET "http://localhost:9000/api/v1/static/resource-types?provider=aws"
````

### List Resource Output Models

Get the JSON structure of resource output models grouped by service category with default mock data.

- **Method**: `GET`
- **Path**: `/static/resource-models`
- **Query Params**:
  - `provider` (optional): e.g., `aws`

**Curl commands**:

```bash
curl -X GET "http://localhost:9000/api/v1/static/resource-models?provider=aws"
```

---

## Users

### Create User

Create a new user.

- **Method**: `POST`
- **Path**: `/users`

**Body (`application/json`)**:

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "auth0_id": "auth0|123456",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

**Curl commands**:

```bash
curl -X POST "http://localhost:9000/api/v1/users" \
     -H "Content-Type: application/json" \
     -d '{
           "name": "Jane Doe",
           "email": "jane@example.com",
           "auth0_id": "auth0|123456"
         }'
```

### Get User

Get a user by ID.

- **Method**: `GET`
- **Path**: `/users/:id`

**Curl commands**:

```bash
curl -X GET "http://localhost:9000/api/v1/users/YOUR_USER_ID_HERE"
```

### List User Projects

Get all projects belonging to a user.

- **Method**: `GET`
- **Path**: `/users/:id/projects`

**Curl commands**:

```bash
curl -X GET "http://localhost:9000/api/v1/users/YOUR_USER_ID_HERE/projects"
```

---

## Projects

### Create Project

Create a new project.

- **Method**: `POST`
- **Path**: `/projects`

**Body (`application/json`)**:

```json
{
  "name": "My Architecture Project",
  "user_id": "YOUR_USER_ID_UUID",
  "cloud_provider": "aws",
  "region": "us-east-1",
  "iac_tool_id": 1
}
```

**Curl commands**:

```bash
curl -X POST "http://localhost:9000/api/v1/projects" \
     -H "Content-Type: application/json" \
     -d '{
           "name": "My Demo Project",
           "user_id": "00000000-0000-0000-0000-000000000001",
           "cloud_provider": "aws",
           "region": "us-east-1",
           "iac_tool_id": 1
         }'
```

### Get Project

Get project details.

- **Method**: `GET`
- **Path**: `/projects/:id`

**Curl commands**:

```bash
curl -X GET "http://localhost:9000/api/v1/projects/YOUR_PROJECT_ID_HERE"
```

### Update Project

Update an existing project.

- **Method**: `PUT`
- **Path**: `/projects/:id`

**Body (`application/json`)**:

```json
{
  "name": "Updated Project Name",
  "region": "us-west-2"
}
```

**Curl commands**:

```bash
curl -X PUT "http://localhost:9000/api/v1/projects/YOUR_PROJECT_ID_HERE" \
     -H "Content-Type: application/json" \
     -d '{
           "name": "Updated Project Name",
           "region": "us-west-2"
         }'
```

### Delete Project

Delete a project.

- **Method**: `DELETE`
- **Path**: `/projects/:id`

**Curl commands**:

```bash
curl -X DELETE "http://localhost:9000/api/v1/projects/YOUR_PROJECT_ID_HERE"
```

---

## Diagrams

### Process Diagram

Process a diagram JSON payload to generate architecture/IaC.

- **Method**: `POST`
- **Path**: `/diagrams/process`
- **Query Params**:
  - `project_name` (optional)
  - `user_id` (optional)
  - `iac_tool_id` (optional)

**Body (`application/json`)**:
The raw diagram JSON structure.

**Curl commands**:

```bash
curl -X POST "http://localhost:9000/api/v1/diagrams/process?project_name=MyDiagram&user_id=...&iac_tool_id=1" \
     -H "Content-Type: application/json" \
     -d @path_to_diagram.json
```
