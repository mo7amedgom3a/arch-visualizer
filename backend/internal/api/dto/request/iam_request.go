package request

// CreateIAMUserRequest represents a request to create an IAM user
type CreateIAMUserRequest struct {
	Name      string `json:"name" binding:"required"`
	Path      string `json:"path"`
	IsVirtual bool   `json:"is_virtual"`
}

// CreateIAMRoleRequest represents a request to create an IAM role
type CreateIAMRoleRequest struct {
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description"`
	Path             string `json:"path"`
	AssumeRolePolicy string `json:"assume_role_policy" binding:"required"`
	IsVirtual        bool   `json:"is_virtual"`
}
