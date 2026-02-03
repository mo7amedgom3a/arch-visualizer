package dto

import "encoding/json"

// ArchitectureResponse represents the full architecture response
type ArchitectureResponse struct {
	Nodes     []ArchitectureNode     `json:"nodes"`
	Edges     []ArchitectureEdge     `json:"edges"`
	Variables []ArchitectureVariable `json:"variables"`
}

// ArchitectureNode represents a node in the diagram/architecture
type ArchitectureNode struct {
	ID       string               `json:"id"`
	Type     string               `json:"type"`
	Position NodePosition         `json:"position"`
	Data     ArchitectureNodeData `json:"data"`
	ParentID *string              `json:"parentId,omitempty"`
}

// NodePosition represents the x, y coordinates
type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// ArchitectureNodeData contains the node's payload
type ArchitectureNodeData struct {
	Label        string                 `json:"label"`
	ResourceType string                 `json:"resourceType"`
	Config       map[string]interface{} `json:"config"`
	IsVisualOnly bool                   `json:"isVisualOnly,omitempty"`
}

// ArchitectureEdge represents a connection between nodes
type ArchitectureEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"` // "contains" or "depends_on" (or default for generic)
	Label  string `json:"label,omitempty"`
}

// ArchitectureVariable represents an input variable
type ArchitectureVariable struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Value       interface{} `json:"value"` // Current value or default
	Description string      `json:"description,omitempty"`
	Sensitive   bool        `json:"sensitive,omitempty"`
}

// UpdateArchitectureRequest represents the payload for saving architecture
type UpdateArchitectureRequest struct {
	Nodes     []ArchitectureNode     `json:"nodes"`
	Edges     []ArchitectureEdge     `json:"edges"`
	Variables []ArchitectureVariable `json:"variables"`
}

// UpdateNodeRequest represents the payload for patching a node
type UpdateNodeRequest struct {
	Position *NodePosition         `json:"position,omitempty"`
	Data     *ArchitectureNodeData `json:"data,omitempty"`
}

// ValidationResponse represents the result of validating the architecture
type ValidationResponse struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationIssue `json:"errors"`
	Warnings []ValidationIssue `json:"warnings"`
}

// ValidationIssue represents a single validation finding
type ValidationIssue struct {
	Type     string `json:"type"` // e.g., "cost", "security", "structural"
	Message  string `json:"message"`
	NodeID   string `json:"nodeId,omitempty"`
	Severity string `json:"severity"` // "error", "warning"
}

// ToJSON converts the response to JSON bytes
func (r *ArchitectureResponse) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
