package request

// CreateProjectRequest represents the request payload for creating a new project.
type CreateProjectRequest struct {
	Name          string `json:"name" binding:"required,min=3,max=100"`
	CloudProvider string `json:"cloud_provider" binding:"required,oneof=aws azure gcp"`
	Region        string `json:"region" binding:"required"`
	IACToolID     uint   `json:"iac_tool_id" binding:"required,min=1"`
	UserID        string `json:"user_id"` // Temporary for testing without auth
}

// UpdateProjectRequest represents the request payload for updating an existing project.
type UpdateProjectRequest struct {
	Name          string `json:"name,omitempty" binding:"omitempty,min=3,max=100"`
	CloudProvider string `json:"cloud_provider,omitempty" binding:"omitempty,oneof=aws azure gcp"`
	Region        string `json:"region,omitempty" binding:"omitempty"`
	IACToolID     uint   `json:"iac_tool_id,omitempty" binding:"omitempty,min=1"`
}
