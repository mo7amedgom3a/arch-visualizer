package dto

// ProcessDiagramRequest represents the request to process a diagram
type ProcessDiagramRequest struct {
	ProjectName string `json:"project_name"`
	IACToolID   uint   `json:"iac_tool_id"` // 1 = Terraform, 2 = Pulumi, etc.
	UserID      string `json:"user_id"`      // UUID string
}

// ProcessDiagramResponse represents the response after processing a diagram
type ProcessDiagramResponse struct {
	Success   bool   `json:"success"`
	ProjectID string `json:"project_id,omitempty"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
}
