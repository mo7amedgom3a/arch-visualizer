package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto"
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

// ListProjects lists all projects with pagination and filtering
// @Summary      List projects
// @Description  Get a list of projects with pagination, sorting, and searching
// @Tags         projects
// @Produce      json
// @Param        page     query     int     false  "Page number"
// @Param        limit    query     int     false  "Items per page"
// @Param        sort     query     string  false  "Sort field"
// @Param        order    query     string  false  "Sort order (asc/desc)"
// @Param        search   query     string  false  "Search term"
// @Param        user_id  query     string  false  "User ID filter"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /projects [get]
func (ctrl *ProjectController) ListProjects(c *gin.Context) {
	// Parse query params
	var query struct {
		Page   int    `form:"page,default=1"`
		Limit  int    `form:"limit,default=20"`
		Sort   string `form:"sort"`
		Order  string `form:"order"`
		Search string `form:"search"`
		UserID string `form:"user_id"`
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	var userID uuid.UUID
	var err error
	if query.UserID != "" {
		userID, err = uuid.Parse(query.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id filter"})
			return
		}
	}

	// TODO: Get authenticated user ID from context if not admin, to enforce RLS-like logic
	// For now we trust the query or return all if nil (depending on service logic)
	// Service List handles nil userID by returning all (or user-scoped if enforced there)

	projects, total, err := ctrl.projectService.List(c.Request.Context(), userID, query.Page, query.Limit, query.Sort, query.Order, query.Search)
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
			// TODO: Add new response fields if needed in DTO
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projectResps,
		"total":    total,
		"page":     query.Page,
		"limit":    query.Limit,
	})
}

// ListUserProjects (Deprecated/Alias to ListProjects with user_id)
func (ctrl *ProjectController) ListUserProjects(c *gin.Context) {
	userIDStr := c.Param("id")
	// logical redirect to ListProjects with query param?
	// For now, implement finding by calling List
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	// Defaults
	page := 1
	limit := 100 // Legacy behavior might expect all?
	projects, total, err := ctrl.projectService.List(c.Request.Context(), userID, page, limit, "", "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list projects: " + err.Error()})
		return
	}

	// Legacy response format
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
		"count":    total,
	})
}

// DuplicateProject duplicates a project
// @Summary      Duplicate a project
// @Description  Create a copy of an existing project and its architecture
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Project ID"
// @Param        body  body      map[string]string  true  "Duplicate Request (requires 'name')"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /projects/{id}/duplicate [post]
func (ctrl *ProjectController) DuplicateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := ctrl.projectService.Duplicate(c.Request.Context(), id, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to duplicate project: " + err.Error()})
		return
	}

	// Return simplified response or full response
	c.JSON(http.StatusCreated, gin.H{
		"id":      project.ID.String(),
		"name":    project.Name,
		"message": "Project duplicated successfully",
	})
}

// GetProjectVersions retrieves version history
// @Summary      Get project versions
// @Description  Retrieve the history of versions for a project
// @Tags         projects
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /projects/{id}/versions [get]
func (ctrl *ProjectController) GetProjectVersions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	versions, err := ctrl.projectService.GetVersions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch versions: " + err.Error()})
		return
	}

	// Map to response DTOs if needed, or return directly
	c.JSON(http.StatusOK, gin.H{
		"versions": versions,
	})
}

// RestoreProjectVersion restores a version
// @Summary      Restore project version
// @Description  Restore a project to a specific version
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Project ID"
// @Param        body  body      map[string]string  true  "Restore Request (requires 'versionId')"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /projects/{id}/restore [post]
func (ctrl *ProjectController) RestoreProjectVersion(c *gin.Context) {
	// It's a POST /projects/:id/restore with body? Or /projects/:id/restore/:versionId?
	// Spec says: POST /projects/:id/restore Body: { "versionId": "..." }

	var req struct {
		VersionID string `json:"versionId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	versionID, err := uuid.Parse(req.VersionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}

	project, err := ctrl.projectService.RestoreVersion(c.Request.Context(), versionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore version: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Project restored successfully",
		"project_id": project.ID.String(),
	})
}

// GetArchitecture retrieves architecture
// @Summary      Get architecture
// @Description  Get full architecture for a project
// @Tags         projects
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  dto.ArchitectureResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /projects/{id}/architecture [get]
func (ctrl *ProjectController) GetArchitecture(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	arch, err := ctrl.projectService.GetArchitecture(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get architecture: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, arch)
}

// UpdateArchitecture updates full architecture
// @Summary      Update architecture
// @Description  Save/update full architecture
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        id    path      string                         true  "Project ID"
// @Param        body  body      dto.UpdateArchitectureRequest  true  "Architecture Data"
// @Success      200  {object}  dto.ArchitectureResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /projects/{id}/architecture [put]
func (ctrl *ProjectController) UpdateArchitecture(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req dto.UpdateArchitectureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arch, err := ctrl.projectService.SaveArchitecture(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save architecture: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"architecture": arch,
		"savedAt":      time.Now().UTC(),
	})
}

// UpdateArchitectureNode updates a single node
// @Summary      Update architecture node
// @Description  Update a specific node in architecture
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        id      path      string                 true  "Project ID"
// @Param        nodeId  path      string                 true  "Node ID"
// @Param        body    body      dto.UpdateNodeRequest  true  "Node Data"
// @Success      200     {object}  dto.ArchitectureNode
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /projects/{id}/architecture/nodes/{nodeId} [patch]
func (ctrl *ProjectController) UpdateArchitectureNode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	nodeID := c.Param("nodeId")
	if nodeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node ID required"})
		return
	}

	var req dto.UpdateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node, err := ctrl.projectService.UpdateNode(c.Request.Context(), id, nodeID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"node": node})
}

// DeleteArchitectureNode deletes a node
// @Summary      Delete architecture node
// @Description  Delete a specific node
// @Tags         projects
// @Produce      json
// @Param        id      path      string  true  "Project ID"
// @Param        nodeId  path      string  true  "Node ID"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /projects/{id}/architecture/nodes/{nodeId} [delete]
func (ctrl *ProjectController) DeleteArchitectureNode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	nodeID := c.Param("nodeId")
	if nodeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node ID required"})
		return
	}

	if err := ctrl.projectService.DeleteNode(c.Request.Context(), id, nodeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete node: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Node deleted"})
}

// ValidateArchitecture validates architecture
// @Summary      Validate architecture
// @Description  Validate project architecture
// @Tags         projects
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  dto.ValidationResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /projects/{id}/architecture/validate [post]
func (ctrl *ProjectController) ValidateArchitecture(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	res, err := ctrl.projectService.ValidateArchitecture(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
