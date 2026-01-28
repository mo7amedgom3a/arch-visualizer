# Arch Visualizer API Server

This is the main API server for the Arch Visualizer application. It provides endpoints for processing cloud architecture diagrams and persisting them as projects.

## Running the Server

### Basic Usage

```bash
# Run the server on default port (8080)
go run cmd/api/main.go

# Run on a custom port
go run cmd/api/main.go -port 3000

# Run migrations before starting
go run cmd/api/main.go -migrate
```

### Build and Run

```bash
# Build the binary
go build -o bin/api cmd/api/main.go

# Run the binary
./bin/api -port 8080
```

## API Endpoints

### POST /api/diagrams/process

Processes a diagram JSON and creates a project with resources in the database.

**Query Parameters:**
- `project_name` (optional): Name for the project (default: "Untitled Project")
- `iac_tool_id` (optional): IaC tool ID (default: 1 for Terraform)
- `user_id` (optional): User UUID (default: test user ID)

**Request Body:**
The request body should contain the diagram JSON from the frontend (see `json-request-diagram-valid.json` for format).

**Example Request:**
```bash
curl -X POST "http://localhost:8080/api/diagrams/process?project_name=My%20Project&iac_tool_id=1&user_id=00000000-0000-0000-0000-000000000001" \
  -H "Content-Type: application/json" \
  -d @json-request-diagram-valid.json
```

**Response:**
```json
{
  "success": true,
  "project_id": "uuid-here",
  "message": "Diagram processed successfully. Project created with ID: uuid-here"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "Error message here"
}
```

### GET /api/health

Health check endpoint.

**Response:**
```json
{
  "status": "ok",
  "service": "arch-visualizer-api"
}
```

### GET /

Root endpoint with API information.

**Response:**
```json
{
  "message": "Arch Visualizer API",
  "version": "1.0.0"
}
```

## Testing

Use the provided test script:

```bash
./scripts/test_api.sh
```

Or manually test with curl:

```bash
curl -X POST "http://localhost:8080/api/diagrams/process?project_name=Test%20Project" \
  -H "Content-Type: application/json" \
  -d @json-request-diagram-valid.json
```

## Architecture

The API follows this flow:

1. **Request** → Handler receives diagram JSON
2. **Parse** → `diagram/parser` parses IR JSON into graph structure
3. **Validate** → `diagram/validator` validates graph structure
4. **Build Graph** → `diagram/graph` builds containment and dependency graphs
5. **Map to Domain** → `domain/architecture` maps to domain models
6. **Persist** → `diagram/service` saves project and resources to database

## Environment Variables

The server uses environment variables for database configuration (via `.env` file):

- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name (default: arch_visualizer)
- `DB_SSLMODE`: SSL mode (default: disable)
