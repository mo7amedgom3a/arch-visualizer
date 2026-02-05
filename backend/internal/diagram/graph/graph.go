package graph

// DiagramGraph represents the normalized diagram structure
// It contains nodes and edges after filtering and normalization
type DiagramGraph struct {
	Nodes     map[string]*Node
	Edges     []*Edge
	Variables []Variable
	Outputs   []Output
	Policies  []Policy
	UI        *ProjectUIState
}

// ProjectUIState represents the global UI state for the diagram
type ProjectUIState struct {
	Zoom            float64  `json:"zoom"`
	ViewportX       float64  `json:"viewportX"`
	ViewportY       float64  `json:"viewportY"`
	SelectedNodeIDs []string `json:"selectedNodeIds"`
	SelectedEdgeIDs []string `json:"selectedEdgeIds"`
}

// Variable represents a Terraform input variable from the diagram
type Variable struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // e.g. "string", "number", "bool", "list(string)"
	Description string      `json:"description"`
	Default     interface{} `json:"default,omitempty"`
	Sensitive   bool        `json:"sensitive,omitempty"`
}

// Output represents a Terraform output value from the diagram
type Output struct {
	Name        string `json:"name"`
	Value       string `json:"value"` // Expression like "aws_vpc.main.id"
	Description string `json:"description"`
	Sensitive   bool   `json:"sensitive,omitempty"`
}

// Policy represents a policy definition from the diagram
type Policy struct {
	EdgeID           string                 `json:"edgeId"`
	SourceID         string                 `json:"sourceId"`
	TargetID         string                 `json:"targetId"`
	SourceType       string                 `json:"sourceType"`
	TargetType       string                 `json:"targetType"`
	Role             string                 `json:"role"`
	PolicyTemplateID string                 `json:"policyTemplateId"`
	PolicyDocument   map[string]interface{} `json:"policyDocument"`
}

// GetNode returns a node by ID
func (g *DiagramGraph) GetNode(id string) (*Node, bool) {
	node, exists := g.Nodes[id]
	return node, exists
}

// GetChildren returns all child nodes of a given parent node ID
func (g *DiagramGraph) GetChildren(parentID string) []*Node {
	children := make([]*Node, 0)
	for _, node := range g.Nodes {
		if node.ParentID != nil && *node.ParentID == parentID {
			children = append(children, node)
		}
	}
	return children
}

// GetParent returns the parent node of a given node ID
func (g *DiagramGraph) GetParent(childID string) (*Node, bool) {
	node, exists := g.Nodes[childID]
	if !exists || node.ParentID == nil {
		return nil, false
	}
	return g.GetNode(*node.ParentID)
}

// GetContainmentEdges returns all containment edges
func (g *DiagramGraph) GetContainmentEdges() []*Edge {
	edges := make([]*Edge, 0)
	for _, edge := range g.Edges {
		if edge.IsContainment() {
			edges = append(edges, edge)
		}
	}
	return edges
}

// GetDependencyEdges returns all dependency edges
func (g *DiagramGraph) GetDependencyEdges() []*Edge {
	edges := make([]*Edge, 0)
	for _, edge := range g.Edges {
		if edge.IsDependency() {
			edges = append(edges, edge)
		}
	}
	return edges
}

// BuildContainmentTree builds a tree structure from containment relationships
// Returns a map of parent ID -> list of child nodes
func (g *DiagramGraph) BuildContainmentTree() map[string][]*Node {
	tree := make(map[string][]*Node)
	for _, node := range g.Nodes {
		if node.ParentID != nil {
			parentID := *node.ParentID
			if _, exists := tree[parentID]; !exists {
				tree[parentID] = make([]*Node, 0)
			}
			tree[parentID] = append(tree[parentID], node)
		}
	}
	return tree
}

// GetRootNodes returns all nodes that have no parent (top-level nodes)
func (g *DiagramGraph) GetRootNodes() []*Node {
	roots := make([]*Node, 0)
	for _, node := range g.Nodes {
		if node.ParentID == nil {
			roots = append(roots, node)
		}
	}
	return roots
}

// FindRegionNode finds the region node in the graph
func (g *DiagramGraph) FindRegionNode() (*Node, bool) {
	for _, node := range g.Nodes {
		if node.IsRegion() {
			return node, true
		}
	}
	return nil, false
}
