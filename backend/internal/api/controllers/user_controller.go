package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto/request"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/dto/response"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// UserController handles user-related requests
type UserController struct {
	userService serverinterfaces.UserService
}

// NewUserController creates a new UserController
func NewUserController(userService serverinterfaces.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser creates a new user
// @Summary      Create a new user
// @Description  Create a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      request.CreateUserRequest  true  "User creation request"
// @Success      201   {object}  response.UserResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svcReq := &serverinterfaces.CreateUserRequest{
		Name:      req.Name,
		Email:     req.Email,
		Auth0ID:   req.Auth0ID,
		AvatarURL: req.AvatarURL,
	}

	user, err := ctrl.userService.Create(c.Request.Context(), svcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	resp := response.UserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	if user.Avatar != nil {
		resp.AvatarURL = *user.Avatar
	}

	c.JSON(http.StatusCreated, resp)
}

// GetUser retrieves a user by ID
// @Summary      Get a user
// @Description  Get a user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.UserResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/{id} [get]
func (ctrl *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := ctrl.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user: " + err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	resp := response.UserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	if user.Avatar != nil {
		resp.AvatarURL = *user.Avatar
	}

	c.JSON(http.StatusOK, resp)
}
