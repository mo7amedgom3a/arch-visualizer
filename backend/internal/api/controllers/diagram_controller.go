package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// DiagramController handles diagram-related requests
type DiagramController struct {
	pipelineOrchestrator serverinterfaces.PipelineOrchestrator
}

// NewDiagramController creates a new DiagramController
func NewDiagramController(pipelineOrchestrator serverinterfaces.PipelineOrchestrator) *DiagramController {
	return &DiagramController{
		pipelineOrchestrator: pipelineOrchestrator,
	}
}

// ProcessDiagram processes a diagram JSON
// @Summary      Process a collected diagram
// @Description  Process a diagram JSON and create/update a project
// @Tags         diagrams
// @Accept       json
// @Produce      json
// @Param        project_name  query     string  false  "Project Name"
// @Param        iac_tool_id   query     int     false  "IaC Tool ID (1=Terraform)"
// @Param        user_id       query     string  false  "User ID"
// @Param        diagram       body      object  true   "Diagram JSON"
// @Success      200           {object}  map[string]interface{}
// @Failure      400           {object}  map[string]interface{}
// @Failure      500           {object}  map[string]interface{}
// @Router       /diagrams/process [post]
func (ctrl *DiagramController) ProcessDiagram(c *gin.Context) {
	// Read body (diagram JSON)
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body: " + err.Error()})
		return
	}
	defer c.Request.Body.Close()

	if len(jsonData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty request body"})
		return
	}

	// Get metadata from query parameters
	projectName := c.Query("project_name")
	if projectName == "" {
		projectName = "Untitled Project"
	}

	iacToolID := uint(1) // Default to Terraform
	if iacToolIDStr := c.Query("iac_tool_id"); iacToolIDStr != "" {
		var parsedID uint
		if _, err := fmt.Sscanf(iacToolIDStr, "%d", &parsedID); err == nil {
			iacToolID = parsedID
		}
	}

	userIDStr := c.Query("user_id")
	// TODO: Get from auth middleware
	if userIDStr == "" {
		userIDStr = "00000000-0000-0000-0000-000000000001"
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	// Process
	req := &serverinterfaces.ProcessDiagramRequest{
		JSONData:    jsonData,
		UserID:      userID,
		ProjectName: projectName,
		IACToolID:   iacToolID,
		// CloudProvider and Region are currently default in pipeline or extracted from diagram,
		// they can be enhanced to be passed in query params too
	}

	result, err := ctrl.pipelineOrchestrator.ProcessDiagram(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process diagram: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"project_id": result.ProjectID.String(),
		"message":    result.Message,
	})
}
