package request

// CreateUserRequest represents the request payload for creating a new user.
type CreateUserRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=100"`
	Email     string `json:"email" binding:"required,email"`
	Auth0ID   string `json:"auth0_id" binding:"required"`
	AvatarURL string `json:"avatar_url,omitempty" binding:"omitempty,url"`
}

// UpdateUserRequest represents the request payload for updating a user profile.
type UpdateUserRequest struct {
	Name      string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	AvatarURL string `json:"avatar_url,omitempty" binding:"omitempty,url"`
}
