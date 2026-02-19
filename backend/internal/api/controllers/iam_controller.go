package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto/request"
	awsiam "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/iam/outputs"
	iamservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/iam"
)

// IAMController handles IAM-related requests
type IAMController struct {
	iamService iamservice.AWSIAMService
}

// NewIAMController creates a new IAMController
func NewIAMController(iamService iamservice.AWSIAMService) *IAMController {
	return &IAMController{
		iamService: iamService,
	}
}

// ListPolicies retrieves AWS managed policies
// @Summary      List AWS Managed Policies
// @Description  Get a list of AWS managed policies, optionally filtered by service
// @Tags         iam
// @Produce      json
// @Param        service  query     string  false  "Filter by service (e.g. s3, lambda)"
// @Success      200      {array}   outputs.PolicyOutput
// @Failure      500      {object}  map[string]interface{}
// @Router       /iam/policies [get]
func (ctrl *IAMController) ListPolicies(c *gin.Context) {
	service := c.Query("service")
	var scope *string
	if service != "" {
		scope = &service
	}

	var policies []*outputs.PolicyOutput
	policies, err := ctrl.iamService.ListAWSManagedPolicies(c.Request.Context(), scope, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list policies: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, policies)
	c.JSON(http.StatusOK, policies)
}

// ListPoliciesBetweenServices retrieves policies relevant to source and destination services
// @Summary      List Policies Between Services
// @Description  Get a list of policies that allow source service to access destination service
// @Tags         iam
// @Produce      json
// @Param        source       query     string  true  "Source service (e.g. lambda)"
// @Param        destination  query     string  true  "Destination service (e.g. s3)"
// @Success      200          {array}   outputs.PolicyOutput
// @Failure      400          {object}  map[string]interface{}
// @Failure      500          {object}  map[string]interface{}
// @Router       /iam/policies/between [get]
func (ctrl *IAMController) ListPoliciesBetweenServices(c *gin.Context) {
	source := c.Query("source")
	destination := c.Query("destination")

	if source == "" || destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both source and destination query parameters are required"})
		return
	}

	policies, err := ctrl.iamService.ListPoliciesBetweenServices(c.Request.Context(), source, destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list policies: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, policies)
}

// CreateUser creates a new IAM user (supports virtual)
// @Summary      Create IAM User
// @Description  Create a new IAM user, potentially virtual
// @Tags         iam
// @Accept       json
// @Produce      json
// @Param        user  body      request.CreateIAMUserRequest  true  "User creation request"
// @Success      201   {object}  outputs.UserOutput
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /iam/users [post]
func (ctrl *IAMController) CreateUser(c *gin.Context) {
	var req request.CreateIAMUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userModel := &awsiam.User{
		Name:      req.Name,
		IsVirtual: req.IsVirtual,
	}
	if req.Path != "" {
		userModel.Path = &req.Path
	}

	user, err := ctrl.iamService.CreateUser(c.Request.Context(), userModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// CreateRole creates a new IAM role (supports virtual)
// @Summary      Create IAM Role
// @Description  Create a new IAM role, potentially virtual
// @Tags         iam
// @Accept       json
// @Produce      json
// @Param        role  body      request.CreateIAMRoleRequest  true  "Role creation request"
// @Success      201   {object}  outputs.RoleOutput
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /iam/roles [post]
func (ctrl *IAMController) CreateRole(c *gin.Context) {
	var req request.CreateIAMRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleModel := &awsiam.Role{
		Name:             req.Name,
		IsVirtual:        req.IsVirtual,
		AssumeRolePolicy: req.AssumeRolePolicy,
	}
	if req.Description != "" {
		roleModel.Description = &req.Description
	}
	if req.Path != "" {
		roleModel.Path = &req.Path
	}

	role, err := ctrl.iamService.CreateRole(c.Request.Context(), roleModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}
