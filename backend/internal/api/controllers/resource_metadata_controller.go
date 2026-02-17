package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// ResourceMetadataController handles resource schema requests.
type ResourceMetadataController struct {
	metadataService serverinterfaces.ResourceMetadataService
}

// NewResourceMetadataController creates a new ResourceMetadataController.
func NewResourceMetadataController(metadataService serverinterfaces.ResourceMetadataService) *ResourceMetadataController {
	return &ResourceMetadataController{
		metadataService: metadataService,
	}
}

// ListSchemas returns all resource schemas for a provider/service pair.
// @Summary      List resource schemas
// @Description  Get structured schema definitions for all resources under a provider/service
// @Tags         metadata
// @Produce      json
// @Success      200  {array}   serverinterfaces.ResourceSchemaDTO
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /aws/networking/schemas [get]
func (ctrl *ResourceMetadataController) ListSchemas(c *gin.Context) {
	schemas, err := ctrl.metadataService.ListResourceSchemas(c.Request.Context(), "aws", "networking")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schemas)
}

// GetSchema returns a single resource schema by resource label.
// @Summary      Get resource schema
// @Description  Get the structured schema definition for a specific networking resource
// @Tags         metadata
// @Produce      json
// @Param        resource  path  string  true  "Resource label (e.g. vpc, subnet)"
// @Success      200  {object}  serverinterfaces.ResourceSchemaDTO
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /aws/networking/schemas/{resource} [get]
func (ctrl *ResourceMetadataController) GetSchema(c *gin.Context) {
	resource := c.Param("resource")

	schema, err := ctrl.metadataService.GetResourceSchema(c.Request.Context(), "aws", "networking", resource)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schema)
}

// ListSchemasByService returns all resource schemas for a given AWS service.
// @Summary      List resource schemas by service
// @Description  Get structured schema definitions for all resources under a specific AWS service (networking, compute, storage, database)
// @Tags         metadata
// @Produce      json
// @Param        service  path  string  true  "AWS service (networking, compute, storage, database)"
// @Success      200  {array}   serverinterfaces.ResourceSchemaDTO
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /aws/{service}/schemas [get]
func (ctrl *ResourceMetadataController) ListSchemasByService(c *gin.Context) {
	service := c.Param("service")

	schemas, err := ctrl.metadataService.ListResourceSchemas(c.Request.Context(), "aws", service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schemas)
}

// GetSchemaByService returns a single resource schema for a given AWS service and resource label.
// @Summary      Get resource schema by service
// @Description  Get the structured schema definition for a specific resource under an AWS service
// @Tags         metadata
// @Produce      json
// @Param        service   path  string  true  "AWS service (networking, compute, storage, database)"
// @Param        resource  path  string  true  "Resource label (e.g. ec2_instance, s3_bucket)"
// @Success      200  {object}  serverinterfaces.ResourceSchemaDTO
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /aws/{service}/schemas/{resource} [get]
func (ctrl *ResourceMetadataController) GetSchemaByService(c *gin.Context) {
	service := c.Param("service")
	resource := c.Param("resource")

	schema, err := ctrl.metadataService.GetResourceSchema(c.Request.Context(), "aws", service, resource)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schema)
}
