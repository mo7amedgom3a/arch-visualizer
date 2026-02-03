package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto/request"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto/response"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// ProjectController handles project-related requests
type ProjectController struct {
	projectService serverinterfaces.ProjectService
}

// NewProjectController creates a new ProjectController
func NewProjectController(projectService serverinterfaces.ProjectService) *ProjectController {
	return &ProjectController{
		projectService: projectService,
	}
}

// CreateProject creates a new project
// @Summary      Create a new project
// @Description  Create a new project for the authenticated user
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        project  body      request.CreateProjectRequest  true  "Project creation request"
// @Success      201      {object}  response.ProjectResponse
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /projects [post]
func (ctrl *ProjectController) CreateProject(c *gin.Context) {
	var req request.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Get user ID from context (from auth middleware)
	// For now, we'll use user_id from request or generate dummy
	var userID uuid.UUID
	var err error
	if req.UserID != "" {
		userID, err = uuid.Parse(req.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}
	} else {
		// Fallback to dummy which will likely fail FK constraint unless user exists
		// But in unit tests we might allow it. In integration, we need valid user.
		// Let's create a specific placeholder logic or error
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	svcReq := &serverinterfaces.CreateProjectRequest{
		Name:          req.Name,
		UserID:        userID,
		IACTargetID:   req.IACToolID,
		CloudProvider: req.CloudProvider,
		Region:        req.Region,
	}

	project, err := ctrl.projectService.Create(c.Request.Context(), svcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project: " + err.Error()})
		return
	}

	resp := response.ProjectResponse{
		ID:            project.ID.String(),
		Name:          project.Name,
		CloudProvider: project.CloudProvider,
		Region:        project.Region,
		IACToolID:     project.InfraToolID,
		UserID:        project.UserID.String(),
		CreatedAt:     project.CreatedAt,
		UpdatedAt:     project.UpdatedAt,
	}

	c.JSON(http.StatusCreated, resp)
}

// GetProject retrieves a project by ID
// @Summary      Get a project
// @Description  Get a project by its ID
// @Tags         projects
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  response.ProjectResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /projects/{id} [get]
func (ctrl *ProjectController) GetProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := ctrl.projectService.GetByID(c.Request.Context(), id)
	if err != nil {
		// Differentiate between not found and internal error if possible
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project: " + err.Error()})
		return
	}
	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	resp := response.ProjectResponse{
		ID:            project.ID.String(),
		Name:          project.Name,
		CloudProvider: project.CloudProvider,
		Region:        project.Region,
		IACToolID:     project.InfraToolID,
		UserID:        project.UserID.String(),
		CreatedAt:     project.CreatedAt,
		UpdatedAt:     project.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateProject updates an existing project
// @Summary      Update a project
// @Description  Update an existing project's details
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        id       path      string                        true  "Project ID"
// @Param        project  body      request.UpdateProjectRequest  true  "Project update request"
// @Success      200      {object}  response.ProjectResponse
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /projects/{id} [put]
func (ctrl *ProjectController) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req request.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch existing project first
	project, err := ctrl.projectService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		return
	}
	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Update fields
	if req.Name != "" {
		project.Name = req.Name
	}
	if req.CloudProvider != "" {
		project.CloudProvider = req.CloudProvider
	}
	if req.Region != "" {
		project.Region = req.Region
	}
	if req.IACToolID != 0 {
		project.InfraToolID = req.IACToolID
	}

	if err := ctrl.projectService.Update(c.Request.Context(), project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project: " + err.Error()})
		return
	}

	resp := response.ProjectResponse{
		ID:            project.ID.String(),
		Name:          project.Name,
		CloudProvider: project.CloudProvider,
		Region:        project.Region,
		IACToolID:     project.InfraToolID,
		UserID:        project.UserID.String(),
		CreatedAt:     project.CreatedAt,
		UpdatedAt:     project.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteProject deletes a project
// @Summary      Delete a project
// @Description  Delete a project by its ID
// @Tags         projects
// @Produces     json
// @Param        id   path      string  true  "Project ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /projects/{id} [delete]
func (ctrl *ProjectController) DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := ctrl.projectService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListUserProjects lists all projects for a user
// @Summary      List user projects
// @Description  Get a list of all projects for a specific user
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/{id}/projects [get]
func (ctrl *ProjectController) ListUserProjects(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	projects, err := ctrl.projectService.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list projects: " + err.Error()})
		return
	}

	projectResps := make([]response.ProjectResponse, len(projects))
	for i, project := range projects {
		projectResps[i] = response.ProjectResponse{
			ID:            project.ID.String(),
			Name:          project.Name,
			CloudProvider: project.CloudProvider,
			Region:        project.Region,
			IACToolID:     project.InfraToolID,
			UserID:        project.UserID.String(),
			CreatedAt:     project.CreatedAt,
			UpdatedAt:     project.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projectResps,
		"count":    len(projects),
	})
}
