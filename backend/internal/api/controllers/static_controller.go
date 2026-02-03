package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// StaticController handles static data requests
type StaticController struct {
	staticService serverinterfaces.StaticDataService
}

// NewStaticController creates a new StaticController
func NewStaticController(staticService serverinterfaces.StaticDataService) *StaticController {
	return &StaticController{
		staticService: staticService,
	}
}

// ListProviders retrieves supported cloud providers
// @Summary      List supported cloud providers
// @Description  Get a list of all supported cloud providers (AWS, GCP, etc.)
// @Tags         static
// @Produce      json
// @Success      200  {object}  []string  "List of providers"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /static/providers [get]
func (ctrl *StaticController) ListProviders(c *gin.Context) {
	providers, err := ctrl.staticService.ListProviders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list providers"})
		return
	}
	c.JSON(http.StatusOK, providers)
}

// ListResourceTypes retrieves resource types
// @Summary      List resource types
// @Description  Get a list of all supported resource types, optionally filtered by provider
// @Tags         static
// @Produce      json
// @Param        provider  query     string  false  "Provider name (e.g. aws)"
// @Success      200       {object}  interface{}
// @Failure      500       {object}  map[string]string "Internal Server Error"
// @Router       /static/resource-types [get]
func (ctrl *StaticController) ListResourceTypes(c *gin.Context) {
	provider := c.Query("provider")

	var err error
	var types interface{} // using interface to avoid strict typing issues for now

	if provider != "" {
		types, err = ctrl.staticService.ListResourceTypesByProvider(c.Request.Context(), provider)
	} else {
		types, err = ctrl.staticService.ListResourceTypes(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list resource types: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, types)
}
