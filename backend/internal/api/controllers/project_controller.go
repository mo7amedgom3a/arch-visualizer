package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto/request"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto/response"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// Ensure dto is imported for swag annotations
var _ dto.ArchitectureResponse

// ProjectController handles project-related requests
type ProjectController struct {
	projectService serverinterfaces.ProjectService
}

// NewProjectController creates a new ProjectController
func NewProjectController(projectService serverinterfaces.ProjectService) *ProjectController {
	return &ProjectController{projectService: projectService}
}

// ── Project CRUD (non-versioned) ──────────────────────────────────────────────

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

	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	project, err := ctrl.projectService.Create(c.Request.Context(), &serverinterfaces.CreateProjectRequest{
		Name:          req.Name,
		UserID:        userID,
		IACTargetID:   req.IACToolID,
		CloudProvider: req.CloudProvider,
		Region:        req.Region,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, projectToResponse(project))
}

// GetProject retrieves a project by ID
// @Summary      Get a project
// @Description  Get a project snapshot by its ID
// @Tags         projects
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  response.ProjectResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /projects/{id} [get]
func (ctrl *ProjectController) GetProject(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	project, err := ctrl.projectService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project: " + err.Error()})
		return
	}
	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(http.StatusOK, projectToResponse(project))
}

// UpdateProject performs an in-place metadata update (no new snapshot).
// @Summary      Update project metadata
// @Description  In-place update of project metadata (name, description, cloud_provider, region, iac_tool_id). Does NOT create a new version.
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
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	var req request.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := ctrl.projectService.GetByID(c.Request.Context(), id)
	if err != nil || project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

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

	updated, err := ctrl.projectService.UpdateMetadata(c.Request.Context(), project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, projectToResponse(updated))
}

// DeleteProject deletes a project
// @Summary      Delete a project
// @Description  Delete a project by its ID
// @Tags         projects
// @Param        id   path      string  true  "Project ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /projects/{id} [delete]
func (ctrl *ProjectController) DeleteProject(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
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
	if query.UserID != "" {
		parsed, err := uuid.Parse(query.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id filter"})
			return
		}
		userID = parsed
	}

	projects, total, err := ctrl.projectService.List(c.Request.Context(), userID, query.Page, query.Limit, query.Sort, query.Order, query.Search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list projects: " + err.Error()})
		return
	}

	resps := make([]response.ProjectResponse, len(projects))
	for i, p := range projects {
		resps[i] = projectToResponse(p)
	}
	c.JSON(http.StatusOK, gin.H{"projects": resps, "total": total, "page": query.Page, "limit": query.Limit})
}

// ListUserProjects – alias for /users/:id/projects
func (ctrl *ProjectController) ListUserProjects(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}
	projects, total, err := ctrl.projectService.List(c.Request.Context(), userID, 1, 100, "", "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list projects: " + err.Error()})
		return
	}
	resps := make([]response.ProjectResponse, len(projects))
	for i, p := range projects {
		resps[i] = projectToResponse(p)
	}
	c.JSON(http.StatusOK, gin.H{"projects": resps, "count": total})
}

// DuplicateProject duplicates a project
// @Summary      Duplicate a project
// @Description  Create an independent copy of an existing project and its architecture (new root, version 1)
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "Project ID"
// @Param        body  body      map[string]string   true  "Requires 'name'"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /projects/{id}/duplicate [post]
func (ctrl *ProjectController) DuplicateProject(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project, version, err := ctrl.projectService.Duplicate(c.Request.Context(), id, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to duplicate project: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":         project.ID.String(),
		"version_id": version.ID.String(),
		"name":       project.Name,
		"message":    "Project duplicated successfully",
	})
}

// ── Architecture (read-only) ──────────────────────────────────────────────────

// GetArchitecture retrieves a project's latest architecture (read-only).
// @Summary      Get architecture
// @Description  Get the full architecture for a project snapshot. For a specific historical state, use GET /projects/{id}/versions/{version_id}.
// @Tags         projects
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  dto.ArchitectureResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /projects/{id}/architecture [get]
func (ctrl *ProjectController) GetArchitecture(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	arch, err := ctrl.projectService.GetArchitecture(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get architecture: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, arch)
}

// ── Version CRUD ──────────────────────────────────────────────────────────────

// CreateVersion creates a new immutable architecture snapshot (version).
// @Summary      Create new version
// @Description  Save a full architecture state as a new immutable snapshot. Returns the version metadata including the new project_id that encodes this snapshot.
// @Tags         versioning
// @Accept       json
// @Produce      json
// @Param        id    path      string                              true  "Project ID (any version in the lineage)"
// @Param        body  body      serverinterfaces.CreateVersionRequest  true  "Architecture state and optional message"
// @Success      201   {object}  serverinterfaces.ProjectVersionDetail
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /projects/{id}/versions [post]
func (ctrl *ProjectController) CreateVersion(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req serverinterfaces.CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	detail, err := ctrl.projectService.CreateVersion(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create version: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, detail)
}

// ListVersions lists all versions in a project's lineage.
// @Summary      List versions
// @Description  Returns the full ordered version chain. Any project_id in the lineage may be used.
// @Tags         versioning
// @Produce      json
// @Param        id   path      string  true  "Project ID (any version in the lineage)"
// @Success      200  {array}   serverinterfaces.ProjectVersionSummary
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /projects/{id}/versions [get]
func (ctrl *ProjectController) ListVersions(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	versions, err := ctrl.projectService.GetVersions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch versions: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, versions)
}

// GetLatestVersion returns the most recent version with full state.
// @Summary      Get latest version
// @Description  Returns the latest version in the chain including full architecture state.
// @Tags         versioning
// @Produce      json
// @Param        id   path      string  true  "Project ID (any version in the lineage)"
// @Success      200  {object}  serverinterfaces.ProjectVersionDetail
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /projects/{id}/versions/latest [get]
func (ctrl *ProjectController) GetLatestVersion(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	detail, err := ctrl.projectService.GetLatestVersion(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest version: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, detail)
}

// GetVersionDetail returns a specific version with its full architecture state.
// @Summary      Get specific version
// @Description  Returns a specific version entry including its full architecture snapshot.
// @Tags         versioning
// @Produce      json
// @Param        id          path      string  true  "Project ID"
// @Param        version_id  path      string  true  "Version ID"
// @Success      200         {object}  serverinterfaces.ProjectVersionDetail
// @Failure      400         {object}  map[string]interface{}
// @Failure      404         {object}  map[string]interface{}
// @Failure      500         {object}  map[string]interface{}
// @Router       /projects/{id}/versions/{version_id} [get]
func (ctrl *ProjectController) GetVersionDetail(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	versionID, ok := parseID(c, "version_id")
	if !ok {
		return
	}
	detail, err := ctrl.projectService.GetVersionByID(c.Request.Context(), id, versionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get version: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, detail)
}

// GetVersionArchitecture returns the architecture snapshot of a specific version.
// @Summary      Get architecture for version
// @Description  Returns the full architecture state captured in a specific version.
// @Tags         versioning
// @Produce      json
// @Param        id          path      string  true  "Project ID"
// @Param        version_id  path      string  true  "Version ID"
// @Success      200         {object}  dto.ArchitectureResponse
// @Failure      400         {object}  map[string]interface{}
// @Failure      404         {object}  map[string]interface{}
// @Failure      500         {object}  map[string]interface{}
// @Router       /projects/{id}/versions/{version_id}/architecture [get]
func (ctrl *ProjectController) GetVersionArchitecture(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	versionID, ok := parseID(c, "version_id")
	if !ok {
		return
	}

	// GetVersionDetail already returns the version along with the State field (dto.ArchitectureResponse)
	detail, err := ctrl.projectService.GetVersionByID(c.Request.Context(), id, versionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get version: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, detail.State)
}

// DeleteVersion removes a version entry (does not delete the snapshot itself).
// @Summary      Delete version
// @Description  Removes a version entry from the chain. The underlying project snapshot is preserved.
// @Tags         versioning
// @Param        id          path      string  true  "Project ID"
// @Param        version_id  path      string  true  "Version ID"
// @Success      204         {object}  nil
// @Failure      400         {object}  map[string]interface{}
// @Failure      404         {object}  map[string]interface{}
// @Failure      500         {object}  map[string]interface{}
// @Router       /projects/{id}/versions/{version_id} [delete]
func (ctrl *ProjectController) DeleteVersion(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	versionID, ok := parseID(c, "version_id")
	if !ok {
		return
	}
	if err := ctrl.projectService.DeleteVersion(c.Request.Context(), id, versionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete version: " + err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ── Version-scoped utility actions ────────────────────────────────────────────

// ValidateVersion validates the architecture of a specific version.
// @Summary      Validate version architecture
// @Description  Validates the architecture state captured in the specified version.
// @Tags         versioning
// @Produce      json
// @Param        id          path      string  true  "Project ID"
// @Param        version_id  path      string  true  "Version ID"
// @Success      200         {object}  dto.ValidationResponse
// @Failure      400         {object}  map[string]interface{}
// @Failure      404         {object}  map[string]interface{}
// @Failure      500         {object}  map[string]interface{}
// @Router       /projects/{id}/versions/{version_id}/validate [post]
func (ctrl *ProjectController) ValidateVersion(c *gin.Context) {
	_, ok := parseID(c, "id")
	if !ok {
		return
	}
	versionID, ok := parseID(c, "version_id")
	if !ok {
		return
	}
	res, err := ctrl.projectService.ValidateVersionArchitecture(c.Request.Context(), versionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// ── helper ────────────────────────────────────────────────────────────────────

func parseID(c *gin.Context, param string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(param))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid " + param})
		return uuid.Nil, false
	}
	return id, true
}

func projectToResponse(p *models.Project) response.ProjectResponse {
	return response.ProjectResponse{
		ID:            p.ID.String(),
		Name:          p.Name,
		CloudProvider: p.CloudProvider,
		Region:        p.Region,
		IACToolID:     p.InfraToolID,
		UserID:        p.UserID.String(),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}
