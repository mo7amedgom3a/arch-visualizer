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
