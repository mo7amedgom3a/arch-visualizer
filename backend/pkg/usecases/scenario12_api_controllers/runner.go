package scenario12_api_controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto/request"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/routes"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server"
)

// Run simulates requests to the API controllers
func Run(ctx context.Context) error {
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("SCENARIO 12: API Controllers Simulation")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("This scenario simulates HTTP requests to the API controllers to verify their behavior.")

	// Step 1: Initialize server and router
	fmt.Println("\n[Step 1] Initializing server and router...")
	srv, err := server.NewServer(slog.Default())
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	router := routes.SetupRouter(srv)
	fmt.Println("✓ Server and Router initialized")

	// Step 2: Create User
	fmt.Println("\n[Step 2] Creating a new user...")
	createUserReq := request.CreateUserRequest{
		Name:    "API Simulator User",
		Email:   fmt.Sprintf("simulator-%d@example.com", time.Now().Unix()),
		Auth0ID: fmt.Sprintf("auth0|simulator-%d", time.Now().Unix()),
	}

	var userResp map[string]interface{}
	if err := performRequest(router, "POST", "/api/v1/users", createUserReq, &userResp, http.StatusCreated); err != nil {
		return err
	}

	userID, ok := userResp["id"].(string)
	if !ok {
		return fmt.Errorf("failed to extract user ID from response")
	}
	fmt.Printf("✓ User created with ID: %s\n", userID)

	// Step 3: List Providers (Static Data)
	fmt.Println("\n[Step 3] Listing supported cloud providers...")
	var providersResp []interface{}
	if err := performRequest(router, "GET", "/api/v1/static/providers", nil, &providersResp, http.StatusOK); err != nil {
		return err
	}
	fmt.Printf("✓ Retrieved %d providers\n", len(providersResp))

	// Step 4: Create Project
	fmt.Println("\n[Step 4] Creating a new project...")
	createProjectReq := request.CreateProjectRequest{
		Name:          "Simulation Project",
		CloudProvider: "aws",
		Region:        "us-east-1",
		IACToolID:     1,
		UserID:        userID,
	}

	var projectResp map[string]interface{}
	if err := performRequest(router, "POST", "/api/v1/projects", createProjectReq, &projectResp, http.StatusCreated); err != nil {
		return err
	}

	projectID, ok := projectResp["id"].(string)
	if !ok {
		return fmt.Errorf("failed to extract project ID from response")
	}
	fmt.Printf("✓ Project created with ID: %s\n", projectID)

	// Step 5: Process Diagram
	fmt.Println("\n[Step 5] Processing a diagram...")
	// Minimal diagram JSON
	diagramJSON := map[string]interface{}{
		"nodes": []interface{}{},
		"edges": []interface{}{},
	}

	// Query params: project_name, user_id
	path := fmt.Sprintf("/api/v1/diagrams/process?project_name=Diagram%%20Project&user_id=%s&iac_tool_id=1", userID)

	var diagramResp map[string]interface{}
	if err := performRequest(router, "POST", path, diagramJSON, &diagramResp, http.StatusOK); err != nil {
		return err
	}

	diagramProjectID, ok := diagramResp["project_id"].(string)
	if !ok {
		return fmt.Errorf("failed to extract project ID from diagram response")
	}
	fmt.Printf("✓ Diagram processed. Created Project ID: %s\n", diagramProjectID)

	// Step 6: Get Project Details
	fmt.Println("\n[Step 6] Retrieving project details...")
	path = fmt.Sprintf("/api/v1/projects/%s", projectID)
	var getProjectResp map[string]interface{}
	if err := performRequest(router, "GET", path, nil, &getProjectResp, http.StatusOK); err != nil {
		return err
	}
	fmt.Printf("✓ Retrieved project: %s\n", getProjectResp["name"])

	// Step 7: List User Projects
	fmt.Println("\n[Step 7] Listing user projects...")
	path = fmt.Sprintf("/api/v1/users/%s/projects", userID)
	var listProjectsResp map[string]interface{}
	if err := performRequest(router, "GET", path, nil, &listProjectsResp, http.StatusOK); err != nil {
		return err
	}
	projectsList, ok := listProjectsResp["projects"].([]interface{})
	if !ok {
		return fmt.Errorf("failed to extract projects list from response")
	}
	fmt.Printf("✓ Retrieved %d projects for user\n", len(projectsList))

	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("✅ SUCCESS: All API controller simulations passed!")
	fmt.Println(strings.Repeat("=", 100))

	return nil
}

func performRequest(r *gin.Engine, method, path string, body interface{}, target interface{}, expectedStatus int) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != expectedStatus {
		return fmt.Errorf("request %s %s failed: expected status %d, got %d. Body: %s", method, path, expectedStatus, w.Code, w.Body.String())
	}

	if target != nil {
		if err := json.Unmarshal(w.Body.Bytes(), target); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}
