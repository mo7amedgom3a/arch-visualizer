package graph

// Node represents a normalized node in the diagram graph
type Node struct {
	ID           string
	Type         string // "containerNode" or "resourceNode"
	ResourceType string // e.g., "vpc", "ec2", "region", "route-table"
	Label        string
	Config       map[string]interface{}
	PositionX    int
	PositionY    int
	ParentID     *string
	Status       string
	IsVisualOnly bool // Track if this is a visual-only node (even if filtered, we track the flag)
	UI           *UIState
}

// UIState represents the full UI state of a node
type UIState struct {
	Position     Position               `json:"position"`
	Width        *float64               `json:"width,omitempty"`
	Height       *float64               `json:"height,omitempty"`
	Style        map[string]interface{} `json:"style,omitempty"`
	Measured     map[string]interface{} `json:"measured,omitempty"`
	Selected     bool                   `json:"selected"`
	Dragging     bool                   `json:"dragging"`
	Resizing     bool                   `json:"resizing"`
	DragHandle   string                 `json:"dragHandle,omitempty"`
	Focusable    bool                   `json:"focusable"`
	Selectable   bool                   `json:"selectable"`
	IsVisualOnly bool                   `json:"isVisualOnly"`
	ZIndex       int                    `json:"zIndex"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// IsContainer returns true if the node is a container node
func (n *Node) IsContainer() bool {
	return n.Type == "containerNode"
}

// IsResource returns true if the node is a resource node
func (n *Node) IsResource() bool {
	return n.Type == "resourceNode"
}

// IsRegion returns true if the node represents a region
func (n *Node) IsRegion() bool {
	return n.ResourceType == "region"
}
