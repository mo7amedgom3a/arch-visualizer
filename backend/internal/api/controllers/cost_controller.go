package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

type CostController struct {
	pricingService      serverinterfaces.PricingService
	projectService      serverinterfaces.ProjectService
	optimizationService serverinterfaces.OptimizationService
}

func NewCostController(
	pricingService serverinterfaces.PricingService,
	projectService serverinterfaces.ProjectService,
	optimizationService serverinterfaces.OptimizationService,
) *CostController {
	return &CostController{
		pricingService:      pricingService,
		projectService:      projectService,
		optimizationService: optimizationService,
	}
}

// GetProjectEstimate godoc
// @Summary      Get project cost estimate
// @Description  Calculate the estimated cost for the entire project architecture
// @Tags         cost
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  serverinterfaces.ArchitectureCostEstimate
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /projects/{id}/cost/estimate [get]
func (cc *CostController) GetProjectEstimate(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Load architecture from database (returns domain model *architecture.Architecture)
	arch, err := cc.projectService.LoadArchitecture(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project or architecture not found"})
		return
	}

	// 720 hours = 30 days * 24 hours (standard monthly estimation)
	duration := 720 * time.Hour

	estimate, err := cc.pricingService.CalculateArchitectureCost(c.Request.Context(), arch, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to estimate cost: %v", err)})
		return
	}

	c.JSON(http.StatusOK, estimate)
}

// GetResourceEstimate godoc
// @Summary      Get resource cost estimate
// @Description  Calculate the estimated cost for a specific resource in the project
// @Tags         cost
// @Produce      json
// @Param        id            path      string  true  "Project ID"
// @Param        resourceName  path      string  true  "Resource Name"
// @Success      200  {object}  serverinterfaces.ResourceCostEstimate
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /projects/{id}/cost/estimate/resources/{resourceName} [get]
func (cc *CostController) GetResourceEstimate(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	resourceName := c.Param("resourceName")
	if resourceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resource name is required"})
		return
	}

	// Load architecture from database
	arch, err := cc.projectService.LoadArchitecture(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project or architecture not found"})
		return
	}

	duration := 720 * time.Hour

	// Iterate over resources to find the one matching resourceName
	for _, res := range arch.Resources {
		if res.Name == resourceName {
			estimate, err := cc.pricingService.CalculateResourceCost(c.Request.Context(), res, duration)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to estimate resource cost: %v", err)})
				return
			}
			c.JSON(http.StatusOK, estimate)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Resource '%s' not found in project", resourceName)})
}

// GetProjectOptimization godoc
// @Summary      Get project cost optimization suggestions
// @Description  Get suggestions for cost optimization based on architectural patterns
// @Tags         cost
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  serverinterfaces.OptimizationWithSavings
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /projects/{id}/cost/optimize [get]
func (cc *CostController) GetProjectOptimization(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Load architecture from database
	arch, err := cc.projectService.LoadArchitecture(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project or architecture not found"})
		return
	}

	suggestions, err := cc.optimizationService.OptimizeArchitecture(c.Request.Context(), arch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get optimization suggestions: %v", err)})
		return
	}

	c.JSON(http.StatusOK, suggestions)
}
