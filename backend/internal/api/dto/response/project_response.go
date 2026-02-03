package response

import "time"

// ProjectResponse represents the project data returned to the client.
type ProjectResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	CloudProvider string    `json:"cloud_provider"`
	Region        string    `json:"region"`
	IACToolID     uint      `json:"iac_tool_id"`
	UserID        string    `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ProjectListResponse represents a list of projects with pagination metadata.
type ProjectListResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}
