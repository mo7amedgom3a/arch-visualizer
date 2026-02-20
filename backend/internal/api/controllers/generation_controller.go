package controllers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

type GenerationController struct {
	orchestrator serverinterfaces.PipelineOrchestrator
	logger       *slog.Logger
}

func NewGenerationController(orchestrator serverinterfaces.PipelineOrchestrator, logger *slog.Logger) *GenerationController {
	return &GenerationController{
		orchestrator: orchestrator,
		logger:       logger,
	}
}

// GenerateCode triggers code generation for a project
// @Summary Generate Infrastructure-as-Code
// @Description Generates IaC files for the specified project and tool
// @Tags Code Generation
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body dto.GenerateCodeRequest true "Generation options"
// @Success 200 {object} dto.GenerationResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id}/generate [post]
func (ctrl *GenerationController) GenerateCode(c *gin.Context) {
	ctrl.logger.Info("Generating code request")
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req dto.GenerateCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call orchestrator to generate code
	out, err := ctrl.orchestrator.GenerateCode(c.Request.Context(), &serverinterfaces.GenerateCodeRequest{
		ProjectID:     projectID,
		Engine:        req.Tool,
		CloudProvider: "aws", // TODO: Get from project or request? Assuming stored in project or inferred.
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate code: " + err.Error()})
		return
	}

	// Map output to response
	var files []dto.GeneratedFileResponse
	for _, f := range out.Files {
		files = append(files, dto.GeneratedFileResponse{
			Name:     f.Path,
			Language: f.Type,
			Content:  f.Content,
			Size:     len(f.Content),
		})
	}

	resp := &dto.GenerationResponse{
		GenerationID: uuid.New(), // ephemeral ID
		ProjectID:    projectID,
		Status:       "completed",
		Tool:         req.Tool,
		Files:        files,
		CreatedAt:    time.Now(),
	}
	resp.ID = resp.GenerationID

	c.JSON(http.StatusOK, resp)
}

// DownloadCode downloads the generated code as a ZIP file
// @Summary Download generated code
// @Description Generates and downloads the ZIP file of the code
// @Tags Code Generation
// @Produce application/zip
// @Param id path string true "Project ID"
// @Param tool query string true "IaC Tool (terraform, pulumi, etc)"
// @Success 200 {file} file
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id}/download [get]
func (ctrl *GenerationController) DownloadCode(c *gin.Context) {
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	tool := c.Query("tool")
	if tool == "" {
		tool = "terraform" // default
	}

	// Generate code
	out, err := ctrl.orchestrator.GenerateCode(c.Request.Context(), &serverinterfaces.GenerateCodeRequest{
		ProjectID:     projectID,
		Engine:        tool,
		CloudProvider: "aws",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate code: " + err.Error()})
		return
	}

	// Create ZIP
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, f := range out.Files {
		zipFile, err := zipWriter.Create(f.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create zip entry: " + err.Error()})
			return
		}
		_, err = zipFile.Write([]byte(f.Content))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write content to zip: " + err.Error()})
			return
		}
	}

	if err := zipWriter.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close zip: " + err.Error()})
		return
	}

	// Serve
	fileName := fmt.Sprintf("project-%s-%s.zip", projectID.String(), tool)
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}

// GenerateCodeForVersion generates IaC for a specific version snapshot.
// @Summary Generate IaC for version
// @Description Generates IaC files for the architecture captured in a specific version
// @Tags versioning
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param version_id path string true "Version ID"
// @Param request body dto.GenerateCodeRequest false "Generation options"
// @Success 200 {object} dto.GenerationResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id}/versions/{version_id}/export/terraform [post]
func (ctrl *GenerationController) GenerateCodeForVersion(c *gin.Context) {
	ctrl.logger.Info("Generating code for version")
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	if _, err := uuid.Parse(c.Param("version_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}
	var req dto.GenerateCodeRequest
	_ = c.ShouldBindJSON(&req)
	if req.Tool == "" {
		req.Tool = "terraform"
	}
	out, err := ctrl.orchestrator.GenerateCode(c.Request.Context(), &serverinterfaces.GenerateCodeRequest{
		ProjectID:     projectID,
		Engine:        req.Tool,
		CloudProvider: "aws",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate code: " + err.Error()})
		return
	}
	var files []dto.GeneratedFileResponse
	for _, f := range out.Files {
		files = append(files, dto.GeneratedFileResponse{Name: f.Path, Language: f.Type, Content: f.Content, Size: len(f.Content)})
	}
	resp := &dto.GenerationResponse{
		GenerationID: uuid.New(),
		ProjectID:    projectID,
		Status:       "completed",
		Tool:         req.Tool,
		Files:        files,
		CreatedAt:    time.Now(),
	}
	resp.ID = resp.GenerationID
	c.JSON(http.StatusOK, resp)
}
