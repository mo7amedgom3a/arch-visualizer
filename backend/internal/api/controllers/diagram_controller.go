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
	diagramService       serverinterfaces.DiagramService
	architectureService  serverinterfaces.ArchitectureService
}

// NewDiagramController creates a new DiagramController
func NewDiagramController(
	pipelineOrchestrator serverinterfaces.PipelineOrchestrator,
	diagramService serverinterfaces.DiagramService,
	architectureService serverinterfaces.ArchitectureService,
) *DiagramController {
	return &DiagramController{
		pipelineOrchestrator: pipelineOrchestrator,
		diagramService:       diagramService,
		architectureService:  architectureService,
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

// ValidateDiagram validates the structure of a diagram
// @Summary      Validate diagram structure
// @Description  Validate correct structure and schema of a diagram JSON
// @Tags         diagrams
// @Accept       json
// @Produce      json
// @Param        diagram  body      object  true   "Diagram JSON"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /diagrams/validate [post]
func (ctrl *DiagramController) ValidateDiagram(c *gin.Context) {
	// Read body
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

	// Parse
	graph, err := ctrl.diagramService.Parse(c.Request.Context(), jsonData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse diagram: " + err.Error()})
		return
	}

	// Validate
	result, err := ctrl.diagramService.Validate(c.Request.Context(), graph, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate diagram: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ValidateDomainRules validates domain rules for a diagram
// @Summary      Validate domain rules
// @Description  Validate architectural rules and constraints (e.g., security, valid configs)
// @Tags         diagrams
// @Accept       json
// @Produce      json
// @Param        diagram  body      object  true   "Diagram JSON"
// @Param        provider query     string  false  "Cloud Provider (default: aws)"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /diagrams/validate-rules [post]
func (ctrl *DiagramController) ValidateDomainRules(c *gin.Context) {
	// Read body
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

	providerStr := c.Query("provider")
	if providerStr == "" {
		providerStr = "aws"
	}

	// Parse
	graph, err := ctrl.diagramService.Parse(c.Request.Context(), jsonData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse diagram: " + err.Error()})
		return
	}

	// Map to Architecture
	arch, err := ctrl.architectureService.MapFromDiagram(c.Request.Context(), graph, "aws") // Hardcoding aws for now as primary support
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to map diagram to architecture: " + err.Error()})
		return
	}

	// Validate Rules
	result, err := ctrl.architectureService.ValidateRules(c.Request.Context(), arch, "aws")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate rules: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
