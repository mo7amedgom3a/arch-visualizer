package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram"
)

// DiagramHandler handles diagram processing requests
type DiagramHandler struct {
	diagramService *diagram.Service
}

// NewDiagramHandler creates a new diagram handler
func NewDiagramHandler() (*DiagramHandler, error) {
	service, err := diagram.NewService(slog.Default())
	if err != nil {
		return nil, fmt.Errorf("failed to create diagram service: %w", err)
	}
	return &DiagramHandler{diagramService: service}, nil
}

// ProcessDiagram handles POST /api/diagrams/process
// Request body should contain:
// - The diagram JSON (from frontend)
// - project_name: string
// - iac_tool_id: uint (optional, defaults to 1 for Terraform)
// - user_id: string (UUID)
func (h *DiagramHandler) ProcessDiagram(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Failed to read request body", err)
		return
	}
	defer r.Body.Close()

	// Parse the request - the body contains the diagram JSON
	// We need to extract metadata (project_name, iac_tool_id, user_id) from query params or headers
	// For now, we'll use query params for metadata and body for diagram JSON

	// Get metadata from query parameters
	projectName := r.URL.Query().Get("project_name")
	if projectName == "" {
		projectName = "Untitled Project" // Default name
	}

	iacToolID := uint(1) // Default to Terraform
	if iacToolIDStr := r.URL.Query().Get("iac_tool_id"); iacToolIDStr != "" {
		var parsedID uint
		if _, err := fmt.Sscanf(iacToolIDStr, "%d", &parsedID); err == nil {
			iacToolID = parsedID
		}
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		// For testing, create a default user ID
		// In production, this should come from authentication
		userIDStr = "00000000-0000-0000-0000-000000000001"
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid user_id format", err)
		return
	}

	// Process the diagram
	projectID, err := h.diagramService.ProcessDiagramRequest(r.Context(), body, userID, projectName, iacToolID)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to process diagram", err)
		return
	}

	// Send success response
	response := dto.ProcessDiagramResponse{
		Success:   true,
		ProjectID: projectID.String(),
		Message:   fmt.Sprintf("Diagram processed successfully. Project created with ID: %s", projectID.String()),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// sendErrorResponse sends an error response
func (h *DiagramHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	response := dto.ProcessDiagramResponse{
		Success: false,
		Error:   fmt.Sprintf("%s: %v", message, err),
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
