package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
)

// IRDiagram represents the intermediate representation JSON from frontend
type IRDiagram struct {
	Nodes     []IRNode          `json:"nodes"`
	Edges     []IREdge          `json:"edges"`
	Variables []IRVariable      `json:"variables"`
	Outputs   []IROutput        `json:"outputs"`
	Policies  []IRPolicy        `json:"policies"`
	Timestamp int64             `json:"timestamp"`
	UIState   *IRProjectUIState `json:"uiState,omitempty"`
}

// IRProjectUIState represents the project-level UI state in IR JSON
type IRProjectUIState struct {
	Zoom            float64  `json:"zoom"`
	ViewportX       float64  `json:"viewportX"`
	ViewportY       float64  `json:"viewportY"`
	SelectedNodeIDs []string `json:"selectedNodeIds"`
	SelectedEdgeIDs []string `json:"selectedEdgeIds"`
}

// IRPolicy represents a policy definition in the IR JSON
type IRPolicy struct {
	EdgeID           string                 `json:"edgeId"`
	SourceID         string                 `json:"sourceId"`
	TargetID         string                 `json:"targetId"`
	SourceType       string                 `json:"sourceType"`
	TargetType       string                 `json:"targetType"`
	Role             string                 `json:"role"`
	PolicyTemplateID string                 `json:"policyTemplateId"`
	PolicyDocument   map[string]interface{} `json:"policyDocument"`
}

// IRVariable represents a Terraform input variable in the IR JSON
type IRVariable struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // e.g. "string", "number", "bool", "list(string)"
	Description string      `json:"description"`
	Default     interface{} `json:"default,omitempty"`
	Sensitive   bool        `json:"sensitive,omitempty"`
}

// IROutput represents a Terraform output value in the IR JSON
type IROutput struct {
	Name        string `json:"name"`
	Value       string `json:"value"` // Expression like "aws_vpc.main.id"
	Description string `json:"description"`
	Sensitive   bool   `json:"sensitive,omitempty"`
}

// IRNode represents a node in the IR JSON
type IRNode struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"` // "containerNode" or "resourceNode"
	Position   *IRPosition            `json:"position,omitempty"`
	ParentID   *string                `json:"parentId,omitempty"`
	Data       IRNodeData             `json:"data"`
	Style      map[string]interface{} `json:"style,omitempty"`
	Measured   map[string]interface{} `json:"measured,omitempty"`
	Selected   bool                   `json:"selected,omitempty"`
	Dragging   bool                   `json:"dragging,omitempty"`
	Resizing   bool                   `json:"resizing,omitempty"`
	Width      *float64               `json:"width,omitempty"`
	Height     *float64               `json:"height,omitempty"`
	Focusable  *bool                  `json:"focusable,omitempty"`
	Selectable *bool                  `json:"selectable,omitempty"`
	ZIndex     int                    `json:"zIndex,omitempty"`
	UI         *IRNodeUI              `json:"ui,omitempty"`
}

// IRNodeUI represents the UI state nested object in IR JSON
type IRNodeUI struct {
	Position     IRPosition             `json:"position"`
	Width        *float64               `json:"width,omitempty"`
	Height       *float64               `json:"height,omitempty"`
	Measured     map[string]interface{} `json:"measured,omitempty"`
	Style        map[string]interface{} `json:"style,omitempty"`
	Resizing     bool                   `json:"resizing,omitempty"`
	DragHandle   string                 `json:"dragHandle,omitempty"`
	Focusable    bool                   `json:"focusable,omitempty"`
	Selectable   bool                   `json:"selectable,omitempty"`
	IsVisualOnly bool                   `json:"isVisualOnly,omitempty"`
	ZIndex       int                    `json:"zIndex,omitempty"`
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
// IREdge represents an edge/connection in the IR JSON
type IREdge struct {
	ID        string                 `json:"id,omitempty"`
	Source    string                 `json:"source"`
	Target    string                 `json:"target"`
	Type      *string                `json:"type,omitempty"`      // "containment", "dependency", "reference"
	Label     string                 `json:"label,omitempty"`     // e.g. "uses"
	Data      map[string]interface{} `json:"data,omitempty"`      // Custom data
	Style     map[string]interface{} `json:"style,omitempty"`     // Visual style
	MarkerEnd map[string]interface{} `json:"markerEnd,omitempty"` // Arrow marker
}

// ParseIRDiagram parses the IR JSON string into an IRDiagram struct
// Handles both direct JSON objects and JSON strings
// Also resolves var.variable_name references using variable defaults
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
			// Resolve variable references before returning
			ResolveVariables(&diagram)
			return &diagram, nil
		}
	}

	// Try parsing as an array of nodes (the file might have nodes at root level)
	var nodesArray []IRNode
	if err := json.Unmarshal(jsonData, &nodesArray); err == nil && len(nodesArray) > 0 {
		// Successfully parsed as array of nodes
		// Now try to find edges, variables, outputs, timestamp in the remaining data
		diagram = IRDiagram{
			Nodes:     nodesArray,
			Edges:     make([]IREdge, 0),
			Variables: make([]IRVariable, 0),
			Outputs:   make([]IROutput, 0),
			Policies:  make([]IRPolicy, 0),
		}

		// Try to parse the full structure again to get edges/variables/outputs/timestamp
		var rawData map[string]interface{}
		if err := json.Unmarshal(jsonData, &rawData); err == nil {
			// Extract edges
			if edgesArray, ok := rawData["edges"].([]interface{}); ok {
				edgesBytes, _ := json.Marshal(edgesArray)
				json.Unmarshal(edgesBytes, &diagram.Edges)
			}

			// Extract variables
			if varsArray, ok := rawData["variables"].([]interface{}); ok {
				varsBytes, _ := json.Marshal(varsArray)
				json.Unmarshal(varsBytes, &diagram.Variables)
			}

			// Extract outputs
			if outputsArray, ok := rawData["outputs"].([]interface{}); ok {
				outputsBytes, _ := json.Marshal(outputsArray)
				json.Unmarshal(outputsBytes, &diagram.Outputs)
			}

			// Extract policies
			if policiesArray, ok := rawData["policies"].([]interface{}); ok {
				policiesBytes, _ := json.Marshal(policiesArray)
				json.Unmarshal(policiesBytes, &diagram.Policies)
			}

			// Extract timestamp
			if ts, ok := rawData["timestamp"].(float64); ok {
				diagram.Timestamp = int64(ts)
			} else if ts, ok := rawData["timestamp"].(int64); ok {
				diagram.Timestamp = ts
			}
		}

		// Resolve variable references before returning
		ResolveVariables(&diagram)
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
		Variables: make([]IRVariable, 0),
		Outputs:   make([]IROutput, 0),
		Policies:  make([]IRPolicy, 0),
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
		varsBytes, _ := json.Marshal(varsArray)
		json.Unmarshal(varsBytes, &diagram.Variables)
	}

	// Extract outputs
	if outputsArray, ok := rawData["outputs"].([]interface{}); ok {
		outputsBytes, _ := json.Marshal(outputsArray)
		json.Unmarshal(outputsBytes, &diagram.Outputs)
	}

	// Extract policies
	if policiesArray, ok := rawData["policies"].([]interface{}); ok {
		policiesBytes, _ := json.Marshal(policiesArray)
		json.Unmarshal(policiesBytes, &diagram.Policies)
	}

	// Extract timestamp
	if ts, ok := rawData["timestamp"].(float64); ok {
		diagram.Timestamp = int64(ts)
	} else if ts, ok := rawData["timestamp"].(int64); ok {
		diagram.Timestamp = ts
	}

	// Resolve variable references before returning
	ResolveVariables(&diagram)
	return &diagram, nil
}

// NormalizeToGraph converts the IR diagram into a normalized graph structure
// Filters out visual-only nodes and builds the graph representation
func NormalizeToGraph(ir *IRDiagram) (*graph.DiagramGraph, error) {
	g := &graph.DiagramGraph{
		Nodes:     make(map[string]*graph.Node),
		Edges:     make([]*graph.Edge, 0),
		Variables: make([]graph.Variable, 0),
		Outputs:   make([]graph.Output, 0),
		Policies:  make([]graph.Policy, 0),
	}

	// Extract Project UI State
	if ir.UIState != nil {
		g.UI = &graph.ProjectUIState{
			Zoom:            ir.UIState.Zoom,
			ViewportX:       ir.UIState.ViewportX,
			ViewportY:       ir.UIState.ViewportY,
			SelectedNodeIDs: ir.UIState.SelectedNodeIDs,
			SelectedEdgeIDs: ir.UIState.SelectedEdgeIDs,
		}
	}

	// Convert IR variables to graph variables
	for _, v := range ir.Variables {
		g.Variables = append(g.Variables, graph.Variable{
			Name:        v.Name,
			Type:        v.Type,
			Description: v.Description,
			Default:     v.Default,
			Sensitive:   v.Sensitive,
		})
	}

	// Convert IR outputs to graph outputs
	for _, o := range ir.Outputs {
		g.Outputs = append(g.Outputs, graph.Output{
			Name:        o.Name,
			Value:       o.Value,
			Description: o.Description,
			Sensitive:   o.Sensitive,
		})
	}

	// Convert IR policies to graph policies
	for _, p := range ir.Policies {
		g.Policies = append(g.Policies, graph.Policy{
			EdgeID:           p.EdgeID,
			SourceID:         p.SourceID,
			TargetID:         p.TargetID,
			SourceType:       p.SourceType,
			TargetType:       p.TargetType,
			Role:             p.Role,
			PolicyTemplateID: p.PolicyTemplateID,
			PolicyDocument:   p.PolicyDocument,
		})
	}

	// First pass: create nodes, tracking visual-only flag for all nodes
	for _, irNode := range ir.Nodes {
		// Extract isVisualOnly flag
		isVisualOnly := false
		if irNode.Data.IsVisualOnly != nil {
			isVisualOnly = *irNode.Data.IsVisualOnly
		}
		if irNode.UI != nil {
			isVisualOnly = irNode.UI.IsVisualOnly
		}

		// Extract UI State
		var ui *graph.UIState
		if irNode.UI != nil {
			// New nested structure
			ui = &graph.UIState{
				Position: graph.Position{
					X: irNode.UI.Position.X,
					Y: irNode.UI.Position.Y,
				},
				Width:        irNode.UI.Width,
				Height:       irNode.UI.Height,
				Measured:     irNode.UI.Measured,
				Style:        irNode.UI.Style,
				Resizing:     irNode.UI.Resizing,
				DragHandle:   irNode.UI.DragHandle,
				Focusable:    irNode.UI.Focusable,
				Selectable:   irNode.UI.Selectable,
				IsVisualOnly: irNode.UI.IsVisualOnly,
				ZIndex:       irNode.UI.ZIndex,
			}
		} else {
			// Backward compatibility for flat structure
			ui = &graph.UIState{
				// Position handled below
				Measured:   irNode.Measured,
				Style:      irNode.Style,
				Selected:   irNode.Selected,
				Dragging:   irNode.Dragging,
				Resizing:   irNode.Resizing,
				ZIndex:     irNode.ZIndex,
				Focusable:  true,
				Selectable: true,
			}

			if irNode.Position != nil {
				ui.Position = graph.Position{
					X: irNode.Position.X,
					Y: irNode.Position.Y,
				}
			}

			if irNode.Width != nil {
				ui.Width = irNode.Width
			}
			if irNode.Height != nil {
				ui.Height = irNode.Height
			}
			if irNode.Focusable != nil {
				ui.Focusable = *irNode.Focusable
			}
			if irNode.Selectable != nil {
				ui.Selectable = *irNode.Selectable
			}
			ui.IsVisualOnly = isVisualOnly
		}

		node := &graph.Node{
			ID:           irNode.ID,
			Type:         irNode.Type,
			ResourceType: irNode.Data.ResourceType,
			Label:        irNode.Data.Label,
			Config:       irNode.Data.Config,
			UI:           ui,
			ParentID:     irNode.ParentID,
			Status:       irNode.Data.Status,
			IsVisualOnly: isVisualOnly,
		}

		g.Nodes[irNode.ID] = node
	}

	// Second pass: create edges
	for _, irEdge := range ir.Edges {
		edgeType := "dependency" // default
		if irEdge.Type != nil {
			edgeType = *irEdge.Type
		}
		// React Flow uses "default" or "smoothstep" commonly. Normalize to "dependency" if it serves that role.
		if edgeType == "default" || edgeType == "smoothstep" {
			edgeType = "dependency"
		}

		// Aggregate configuration/metadata
		config := make(map[string]interface{})
		if irEdge.Label != "" {
			config["label"] = irEdge.Label
		}
		if irEdge.Data != nil {
			for k, v := range irEdge.Data {
				config[k] = v
			}
		}
		if irEdge.Style != nil {
			config["style"] = irEdge.Style
		}
		if irEdge.MarkerEnd != nil {
			config["markerEnd"] = irEdge.MarkerEnd
		}

		edge := &graph.Edge{
			ID:     irEdge.ID,
			Source: irEdge.Source,
			Target: irEdge.Target,
			Type:   edgeType,
			Config: config,
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

// ResolveVariables resolves all var.variable_name references in the IR diagram
// by replacing them with the variable's default value for internal processing.
// It also stores the original variable references in a special "_varRefs" metadata field
// so that Terraform generation can use the variable syntax (var.name) instead of resolved values.
func ResolveVariables(ir *IRDiagram) {
	// Build a map of variable names to their default values
	varDefaults := make(map[string]interface{})
	for _, v := range ir.Variables {
		if v.Default != nil {
			varDefaults[v.Name] = v.Default
		}
	}

	if len(varDefaults) == 0 {
		return // No variables with defaults to resolve
	}

	// Resolve variables in nodes
	for i := range ir.Nodes {
		// Initialize _varRefs metadata to track original variable references
		if ir.Nodes[i].Data.Config == nil {
			ir.Nodes[i].Data.Config = make(map[string]interface{})
		}
		varRefs := make(map[string]string) // field path -> original var reference

		// Track and resolve node ID
		if hasVarRef(ir.Nodes[i].ID) {
			varRefs["_id"] = ir.Nodes[i].ID
			ir.Nodes[i].ID = resolveVarInString(ir.Nodes[i].ID, varDefaults)
		}

		// Track and resolve parent ID
		if ir.Nodes[i].ParentID != nil && hasVarRef(*ir.Nodes[i].ParentID) {
			varRefs["_parentId"] = *ir.Nodes[i].ParentID
			resolved := resolveVarInString(*ir.Nodes[i].ParentID, varDefaults)
			ir.Nodes[i].ParentID = &resolved
		}

		// Track and resolve config values, collecting variable references
		ir.Nodes[i].Data.Config = resolveVarInMapWithTracking(ir.Nodes[i].Data.Config, varDefaults, "", varRefs)

		// Store the variable references in metadata for Terraform generation
		if len(varRefs) > 0 {
			ir.Nodes[i].Data.Config["_varRefs"] = varRefs
		}
	}

	// Resolve variables in edges
	for i := range ir.Edges {
		ir.Edges[i].Source = resolveVarInString(ir.Edges[i].Source, varDefaults)
		ir.Edges[i].Target = resolveVarInString(ir.Edges[i].Target, varDefaults)
	}
}

// hasVarRef checks if a string contains a variable reference
func hasVarRef(s string) bool {
	return strings.Contains(s, "var.")
}

// resolveVarInMapWithTracking recursively resolves var.variable_name references in a map
// and tracks the original variable references in varRefs map
func resolveVarInMapWithTracking(m map[string]interface{}, varDefaults map[string]interface{}, prefix string, varRefs map[string]string) map[string]interface{} {
	if m == nil {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range m {
		fieldPath := k
		if prefix != "" {
			fieldPath = prefix + "." + k
		}
		result[k] = resolveVarInValueWithTracking(v, varDefaults, fieldPath, varRefs)
	}
	return result
}

// resolveVarInValueWithTracking resolves var.variable_name references in any value type
// and tracks original variable references
func resolveVarInValueWithTracking(v interface{}, varDefaults map[string]interface{}, fieldPath string, varRefs map[string]string) interface{} {
	switch val := v.(type) {
	case string:
		if hasVarRef(val) {
			// Store the original variable reference
			varRefs[fieldPath] = val
		}
		resolved := resolveVarInString(val, varDefaults)
		// If the entire string was a variable reference, try to return the actual type
		if strings.HasPrefix(val, "var.") && !strings.Contains(resolved, "var.") {
			varName := strings.TrimPrefix(val, "var.")
			if defaultVal, ok := varDefaults[varName]; ok {
				return defaultVal
			}
		}
		return resolved
	case map[string]interface{}:
		return resolveVarInMapWithTracking(val, varDefaults, fieldPath, varRefs)
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			itemPath := fmt.Sprintf("%s[%d]", fieldPath, i)
			result[i] = resolveVarInValueWithTracking(item, varDefaults, itemPath, varRefs)
		}
		return result
	default:
		return v
	}
}

// varPattern matches var.variable_name patterns
var varPattern = regexp.MustCompile(`var\.([a-zA-Z_][a-zA-Z0-9_]*)`)

// resolveVarInString resolves var.variable_name references in a string
func resolveVarInString(s string, varDefaults map[string]interface{}) string {
	if !strings.Contains(s, "var.") {
		return s
	}

	return varPattern.ReplaceAllStringFunc(s, func(match string) string {
		// Extract variable name (after "var.")
		varName := strings.TrimPrefix(match, "var.")
		if defaultVal, ok := varDefaults[varName]; ok {
			// Convert the default value to string
			return fmt.Sprintf("%v", defaultVal)
		}
		// If variable not found, keep the original reference
		return match
	})
}

// resolveVarInMap recursively resolves var.variable_name references in a map
func resolveVarInMap(m map[string]interface{}, varDefaults map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = resolveVarInValue(v, varDefaults)
	}
	return result
}

// resolveVarInValue resolves var.variable_name references in any value type
func resolveVarInValue(v interface{}, varDefaults map[string]interface{}) interface{} {
	switch val := v.(type) {
	case string:
		resolved := resolveVarInString(val, varDefaults)
		// If the entire string was a variable reference, try to return the actual type
		if strings.HasPrefix(val, "var.") && !strings.Contains(resolved, "var.") {
			varName := strings.TrimPrefix(val, "var.")
			if defaultVal, ok := varDefaults[varName]; ok {
				return defaultVal
			}
		}
		return resolved
	case map[string]interface{}:
		return resolveVarInMap(val, varDefaults)
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			result[i] = resolveVarInValue(item, varDefaults)
		}
		return result
	default:
		return v
	}
}
