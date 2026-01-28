package parser

import (
	"encoding/json"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
)

// IRDiagram represents the intermediate representation JSON from frontend
type IRDiagram struct {
	Nodes     []IRNode      `json:"nodes"`
	Edges     []IREdge      `json:"edges"`
	Variables []interface{} `json:"variables"`
	Timestamp int64         `json:"timestamp"`
}

// IRNode represents a node in the IR JSON
type IRNode struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"` // "containerNode" or "resourceNode"
	Position *IRPosition            `json:"position,omitempty"`
	ParentID *string                `json:"parentId,omitempty"`
	Data     IRNodeData             `json:"data"`
	Style    map[string]interface{} `json:"style,omitempty"`
}

// IRPosition represents node position coordinates
type IRPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// IRNodeData represents the data payload of a node
type IRNodeData struct {
	Label        string                 `json:"label"`
	ResourceType string                 `json:"resourceType"`
	Config       map[string]interface{} `json:"config"`
	Status       string                 `json:"status,omitempty"`
	IsVisualOnly *bool                  `json:"isVisualOnly,omitempty"`
}

// IREdge represents an edge/connection in the IR JSON
type IREdge struct {
	ID     string  `json:"id,omitempty"`
	Source string  `json:"source"`
	Target string  `json:"target"`
	Type   *string `json:"type,omitempty"` // "containment", "dependency", "reference"
}

// ParseIRDiagram parses the IR JSON string into an IRDiagram struct
// Handles both direct JSON objects and JSON strings
func ParseIRDiagram(jsonData []byte) (*IRDiagram, error) {
	// First, try to unmarshal as a JSON string (if the data is double-encoded)
	var jsonStr string
	if err := json.Unmarshal(jsonData, &jsonStr); err == nil {
		// It's a JSON string, unmarshal again
		jsonData = []byte(jsonStr)
	}

	// Try to parse as direct object first (standard format: {"nodes": [...], "edges": [...], ...})
	var diagram IRDiagram
	if err := json.Unmarshal(jsonData, &diagram); err == nil {
		// Check if we got nodes - if so, this is the correct format
		if len(diagram.Nodes) > 0 || (diagram.Nodes != nil && len(diagram.Nodes) == 0) {
			return &diagram, nil
		}
	}

	// Try parsing as an array of nodes (the file might have nodes at root level)
	var nodesArray []IRNode
	if err := json.Unmarshal(jsonData, &nodesArray); err == nil && len(nodesArray) > 0 {
		// Successfully parsed as array of nodes
		// Now try to find edges, variables, timestamp in the remaining data
		diagram = IRDiagram{
			Nodes:     nodesArray,
			Edges:     make([]IREdge, 0),
			Variables: make([]interface{}, 0),
		}

		// Try to parse the full structure again to get edges/variables/timestamp
		var rawData map[string]interface{}
		if err := json.Unmarshal(jsonData, &rawData); err == nil {
			// Extract edges
			if edgesArray, ok := rawData["edges"].([]interface{}); ok {
				edgesBytes, _ := json.Marshal(edgesArray)
				json.Unmarshal(edgesBytes, &diagram.Edges)
			}

			// Extract variables
			if varsArray, ok := rawData["variables"].([]interface{}); ok {
				diagram.Variables = varsArray
			}

			// Extract timestamp
			if ts, ok := rawData["timestamp"].(float64); ok {
				diagram.Timestamp = int64(ts)
			} else if ts, ok := rawData["timestamp"].(int64); ok {
				diagram.Timestamp = ts
			}
		}

		return &diagram, nil
	}

	// Try parsing as raw map to extract fields manually
	var rawData map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse IR diagram JSON: %w", err)
	}

	// Reconstruct the diagram structure from raw map
	diagram = IRDiagram{
		Edges:     make([]IREdge, 0),
		Variables: make([]interface{}, 0),
	}

	// Extract nodes (might be in an array at root or in "nodes" field)
	if nodesArray, ok := rawData["nodes"].([]interface{}); ok {
		nodesBytes, _ := json.Marshal(nodesArray)
		if err := json.Unmarshal(nodesBytes, &diagram.Nodes); err != nil {
			return nil, fmt.Errorf("failed to parse nodes: %w", err)
		}
	}

	// Extract edges
	if edgesArray, ok := rawData["edges"].([]interface{}); ok {
		edgesBytes, _ := json.Marshal(edgesArray)
		if err := json.Unmarshal(edgesBytes, &diagram.Edges); err != nil {
			return nil, fmt.Errorf("failed to parse edges: %w", err)
		}
	}

	// Extract variables
	if varsArray, ok := rawData["variables"].([]interface{}); ok {
		diagram.Variables = varsArray
	}

	// Extract timestamp
	if ts, ok := rawData["timestamp"].(float64); ok {
		diagram.Timestamp = int64(ts)
	} else if ts, ok := rawData["timestamp"].(int64); ok {
		diagram.Timestamp = ts
	}

	return &diagram, nil
}

// NormalizeToGraph converts the IR diagram into a normalized graph structure
// Filters out visual-only nodes and builds the graph representation
func NormalizeToGraph(ir *IRDiagram) (*graph.DiagramGraph, error) {
	g := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// First pass: create nodes, tracking visual-only flag for all nodes
	for _, irNode := range ir.Nodes {
		// Extract isVisualOnly flag (default to false if not specified)
		isVisualOnly := false
		if irNode.Data.IsVisualOnly != nil {
			isVisualOnly = *irNode.Data.IsVisualOnly
		}

		// Extract position
		posX, posY := 0, 0
		if irNode.Position != nil {
			posX = int(irNode.Position.X)
			posY = int(irNode.Position.Y)
		}

		node := &graph.Node{
			ID:           irNode.ID,
			Type:         irNode.Type,
			ResourceType: irNode.Data.ResourceType,
			Label:        irNode.Data.Label,
			Config:       irNode.Data.Config,
			PositionX:    posX,
			PositionY:    posY,
			ParentID:     irNode.ParentID,
			Status:       irNode.Data.Status,
			IsVisualOnly: isVisualOnly,
		}

		g.Nodes[irNode.ID] = node
	}

	// Second pass: create edges
	for _, irEdge := range ir.Edges {
		// Validate that both source and target nodes exist
		if _, exists := g.Nodes[irEdge.Source]; !exists {
			continue // Skip edges to filtered nodes
		}
		if _, exists := g.Nodes[irEdge.Target]; !exists {
			continue // Skip edges to filtered nodes
		}

		edgeType := "dependency" // default
		if irEdge.Type != nil {
			edgeType = *irEdge.Type
		}

		edge := &graph.Edge{
			ID:     irEdge.ID,
			Source: irEdge.Source,
			Target: irEdge.Target,
			Type:   edgeType,
		}

		g.Edges = append(g.Edges, edge)
	}

	// Third pass: build containment relationships from parentId
	// This creates implicit containment edges
	for _, node := range g.Nodes {
		if node.ParentID != nil {
			// Check if parent exists in graph
			if _, exists := g.Nodes[*node.ParentID]; exists {
				// Create containment edge if not already present
				edgeExists := false
				for _, edge := range g.Edges {
					if edge.Source == *node.ParentID && edge.Target == node.ID && edge.Type == "containment" {
						edgeExists = true
						break
					}
				}
				if !edgeExists {
					containmentEdge := &graph.Edge{
						Source: *node.ParentID,
						Target: node.ID,
						Type:   "containment",
					}
					g.Edges = append(g.Edges, containmentEdge)
				}
			}
		}
	}

	return g, nil
}

// ParseAndNormalize is a convenience function that parses and normalizes in one step
func ParseAndNormalize(jsonData []byte) (*graph.DiagramGraph, error) {
	ir, err := ParseIRDiagram(jsonData)
	if err != nil {
		return nil, err
	}
	return NormalizeToGraph(ir)
}
