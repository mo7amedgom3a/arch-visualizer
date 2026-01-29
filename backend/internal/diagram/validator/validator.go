package validator

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator/schema"
)

// ValidationError represents a validation error
type ValidationError struct {
	Code    string
	Message string
	NodeID  string
}

func (e *ValidationError) Error() string {
	if e.NodeID != "" {
		return fmt.Sprintf("[%s] %s (node: %s)", e.Code, e.Message, e.NodeID)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// ValidationResult contains the result of validation
type ValidationResult struct {
	Valid    bool
	Errors   []*ValidationError
	Warnings []*ValidationError
}

// ValidationOptions contains options for validation
type ValidationOptions struct {
	// ValidResourceTypes maps resource type names (in IR format, lowercase) to validity
	// e.g., map["vpc"] = true, map["route-table"] = true
	ValidResourceTypes map[string]bool
	// Provider is the cloud provider (aws, azure, gcp)
	Provider string
	// SchemaRegistry is the schema registry to use for dynamic config validation
	// If nil, uses the default registry
	SchemaRegistry *schema.InMemorySchemaRegistry
}

// Validate performs comprehensive validation on the diagram graph
// If opts is nil, uses default validation (backward compatible)
func Validate(g *graph.DiagramGraph, opts *ValidationOptions) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   make([]*ValidationError, 0),
		Warnings: make([]*ValidationError, 0),
	}

	// Run all validation checks
	validateMissingParents(g, result)
	validateContainmentCycles(g, result)
	validateEdgeReferences(g, result)
	validateDependencyEdges(g, result)
	validateResourceTypes(g, result, opts)
	validateRegionNode(g, result)
	validateConfigSchemaWithRegistry(g, result, opts)
	validateCIDRs(g, result)
	validateContainmentTypes(g, result, opts)

	// Set valid to false if there are any errors
	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return result
}

// validateMissingParents checks that all parentId references point to existing nodes
func validateMissingParents(g *graph.DiagramGraph, result *ValidationResult) {
	for _, node := range g.Nodes {
		if node.ParentID != nil {
			if _, exists := g.Nodes[*node.ParentID]; !exists {
				result.Errors = append(result.Errors, &ValidationError{
					Code:    "MISSING_PARENT",
					Message: fmt.Sprintf("Node '%s' references non-existent parent '%s'", node.ID, *node.ParentID),
					NodeID:  node.ID,
				})
			}
		}
	}
}

// validateContainmentCycles checks for cycles in the containment tree
func validateContainmentCycles(g *graph.DiagramGraph, result *ValidationResult) {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var hasCycle func(nodeID string) bool
	hasCycle = func(nodeID string) bool {
		visited[nodeID] = true
		recStack[nodeID] = true

		// Check children (containment goes parent -> child)
		children := g.GetChildren(nodeID)
		for _, child := range children {
			if !visited[child.ID] {
				if hasCycle(child.ID) {
					return true
				}
			} else if recStack[child.ID] {
				// Found a cycle: we're trying to visit a node that's already in the recursion stack
				return true
			}
		}

		recStack[nodeID] = false
		return false
	}

	// Check all nodes for cycles (not just roots, since cycles might not include roots)
	// But we need to check from each unvisited node to catch all cycles
	for nodeID := range g.Nodes {
		if !visited[nodeID] {
			if hasCycle(nodeID) {
				result.Errors = append(result.Errors, &ValidationError{
					Code:    "CONTAINMENT_CYCLE",
					Message: fmt.Sprintf("Cycle detected in containment tree involving node '%s'", nodeID),
					NodeID:  nodeID,
				})
				// Continue scanning to collect other validation errors as well.
			}
		}
	}
}

// validateEdgeReferences checks that all edges reference existing nodes
func validateEdgeReferences(g *graph.DiagramGraph, result *ValidationResult) {
	for _, edge := range g.Edges {
		if _, exists := g.Nodes[edge.Source]; !exists {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "INVALID_EDGE_SOURCE",
				Message: fmt.Sprintf("Edge references non-existent source node '%s'", edge.Source),
			})
		}
		if _, exists := g.Nodes[edge.Target]; !exists {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "INVALID_EDGE_TARGET",
				Message: fmt.Sprintf("Edge references non-existent target node '%s'", edge.Target),
			})
		}
	}
}

// validateDependencyEdges validates dependency edges (self-loops, invalid endpoints, cycles)
func validateDependencyEdges(g *graph.DiagramGraph, result *ValidationResult) {
	// Self-loop and endpoint validation
	depEdges := make([]*graph.Edge, 0)
	for _, edge := range g.Edges {
		if !edge.IsDependency() {
			continue
		}
		depEdges = append(depEdges, edge)

		if edge.Source == "" || edge.Target == "" {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "DEPENDENCY_INVALID_ENDPOINT",
				Message: "Dependency edge has empty source or target",
			})
			continue
		}

		if edge.Source == edge.Target {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "DEPENDENCY_SELF_LOOP",
				Message: fmt.Sprintf("Dependency edge creates a self-loop on '%s'", edge.Source),
				NodeID:  edge.Source,
			})
			continue
		}

		srcNode, srcExists := g.Nodes[edge.Source]
		tgtNode, tgtExists := g.Nodes[edge.Target]
		if !srcExists || !tgtExists {
			// Covered by validateEdgeReferences; keep going.
			continue
		}

		// Dependencies should generally be between actual resources (not purely containers).
		// We warn (not error) to keep the system flexible.
		if !srcNode.IsResource() || !tgtNode.IsResource() {
			result.Warnings = append(result.Warnings, &ValidationError{
				Code:    "DEPENDENCY_NON_RESOURCE",
				Message: fmt.Sprintf("Dependency edge connects non-resource nodes: '%s' -> '%s'", edge.Source, edge.Target),
			})
		}
	}

	// Dependency cycle detection (DAG expected)
	adj := make(map[string][]string)
	for _, edge := range depEdges {
		// Only consider edges with valid endpoints
		if _, ok := g.Nodes[edge.Source]; !ok {
			continue
		}
		if _, ok := g.Nodes[edge.Target]; !ok {
			continue
		}
		adj[edge.Source] = append(adj[edge.Source], edge.Target)
	}

	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var dfs func(string) bool
	dfs = func(n string) bool {
		visited[n] = true
		recStack[n] = true
		for _, next := range adj[n] {
			if !visited[next] {
				if dfs(next) {
					return true
				}
			} else if recStack[next] {
				return true
			}
		}
		recStack[n] = false
		return false
	}

	for nodeID := range adj {
		if !visited[nodeID] {
			if dfs(nodeID) {
				result.Errors = append(result.Errors, &ValidationError{
					Code:    "DEPENDENCY_CYCLE",
					Message: "Cycle detected in dependency graph",
				})
				break
			}
		}
	}
}

// validateResourceTypes checks that resource types are valid
// If opts is nil or ValidResourceTypes is nil, skips validation (backward compatible)
func validateResourceTypes(g *graph.DiagramGraph, result *ValidationResult, opts *ValidationOptions) {
	// If no options provided, skip resource type validation (backward compatible)
	if opts == nil || opts.ValidResourceTypes == nil || len(opts.ValidResourceTypes) == 0 {
		return
	}

	for _, node := range g.Nodes {
		if node.ResourceType == "" {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "MISSING_RESOURCE_TYPE",
				Message: fmt.Sprintf("Node '%s' has no resource type", node.ID),
				NodeID:  node.ID,
			})
		} else {
			resourceTypeLower := strings.ToLower(node.ResourceType)
			if !opts.ValidResourceTypes[resourceTypeLower] {
				providerMsg := ""
				if opts.Provider != "" {
					providerMsg = fmt.Sprintf(" for provider '%s'", opts.Provider)
				}
				result.Warnings = append(result.Warnings, &ValidationError{
					Code:    "UNKNOWN_RESOURCE_TYPE",
					Message: fmt.Sprintf("Node '%s' has unknown resource type '%s'%s", node.ID, node.ResourceType, providerMsg),
					NodeID:  node.ID,
				})
			}
		}
	}
}

// validateRegionNode checks that there is exactly one region node
func validateRegionNode(g *graph.DiagramGraph, result *ValidationResult) {
	regionNodes := make([]*graph.Node, 0)
	for _, node := range g.Nodes {
		if node.IsRegion() {
			regionNodes = append(regionNodes, node)
		}
	}

	if len(regionNodes) == 0 {
		result.Warnings = append(result.Warnings, &ValidationError{
			Code:    "NO_REGION_NODE",
			Message: "No region node found in diagram",
		})
	} else if len(regionNodes) > 1 {
		result.Warnings = append(result.Warnings, &ValidationError{
			Code:    "MULTIPLE_REGION_NODES",
			Message: fmt.Sprintf("Multiple region nodes found (%d), only one will be used", len(regionNodes)),
		})
	}
}

// validateConfigSchemaWithRegistry performs dynamic config validation using the schema registry.
// This replaces the hardcoded switch statement with schema-driven validation.
func validateConfigSchemaWithRegistry(g *graph.DiagramGraph, result *ValidationResult, opts *ValidationOptions) {
	// Get schema registry
	registry := schema.DefaultRegistry
	if opts != nil && opts.SchemaRegistry != nil {
		registry = opts.SchemaRegistry
	}

	provider := "aws" // default
	if opts != nil && opts.Provider != "" {
		provider = opts.Provider
	}

	for _, node := range g.Nodes {
		rt := strings.ToLower(node.ResourceType)
		if node.Config == nil {
			node.Config = map[string]interface{}{}
		}

		// Get schema for this resource type
		resourceSchema, exists := registry.Get(rt, provider)
		if !exists {
			// No schema found - skip validation (or could add warning)
			continue
		}

		// Validate each field according to schema
		for _, fieldSpec := range resourceSchema.Fields {
			validateField(node, result, fieldSpec)
		}
	}
}

// validateField validates a single field against its specification
func validateField(node *graph.Node, result *ValidationResult, fieldSpec schema.FieldSpec) {
	value, exists := node.Config[fieldSpec.Name]

	// Check required fields
	if fieldSpec.Required && !exists {
		result.Errors = append(result.Errors, &ValidationError{
			Code:    "CONFIG_MISSING_FIELD",
			Message: fmt.Sprintf("Missing required config field '%s' on node '%s' (%s)", fieldSpec.Name, node.ID, node.ResourceType),
			NodeID:  node.ID,
		})
		return
	}

	// If field doesn't exist and isn't required, skip further validation
	if !exists {
		return
	}

	// Type validation
	if !validateFieldType(value, fieldSpec.Type) {
		result.Errors = append(result.Errors, &ValidationError{
			Code:    "CONFIG_INVALID_TYPE",
			Message: fmt.Sprintf("Config field '%s' must be %s on node '%s' (%s)", fieldSpec.Name, fieldSpec.Type, node.ID, node.ResourceType),
			NodeID:  node.ID,
		})
		return
	}

	// Constraint validation
	if fieldSpec.Constraints != nil {
		validateConstraints(node, result, fieldSpec, value)
	}
}

// validateFieldType checks if a value matches the expected type
func validateFieldType(value interface{}, expectedType schema.FieldType) bool {
	if value == nil {
		return false
	}

	switch expectedType {
	case schema.FieldTypeString, schema.FieldTypeCIDR:
		_, ok := value.(string)
		return ok
	case schema.FieldTypeInt:
		// JSON numbers are float64, accept both
		switch value.(type) {
		case int, int32, int64, float64:
			return true
		}
		return false
	case schema.FieldTypeFloat:
		switch value.(type) {
		case float32, float64:
			return true
		}
		return false
	case schema.FieldTypeBool:
		_, ok := value.(bool)
		return ok
	case schema.FieldTypeArray:
		_, ok := value.([]interface{})
		return ok
	case schema.FieldTypeObject:
		_, ok := value.(map[string]interface{})
		return ok
	case schema.FieldTypeAny:
		return true
	}
	return false
}

// validateConstraints validates field constraints
func validateConstraints(node *graph.Node, result *ValidationResult, fieldSpec schema.FieldSpec, value interface{}) {
	c := fieldSpec.Constraints
	if c == nil {
		return
	}

	// String constraints
	if strVal, ok := value.(string); ok {
		// Min length
		if c.MinLength != nil && len(strVal) < *c.MinLength {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "CONFIG_CONSTRAINT_VIOLATION",
				Message: fmt.Sprintf("Field '%s' must be at least %d characters on node '%s'", fieldSpec.Name, *c.MinLength, node.ID),
				NodeID:  node.ID,
			})
		}

		// Max length
		if c.MaxLength != nil && len(strVal) > *c.MaxLength {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "CONFIG_CONSTRAINT_VIOLATION",
				Message: fmt.Sprintf("Field '%s' must be at most %d characters on node '%s'", fieldSpec.Name, *c.MaxLength, node.ID),
				NodeID:  node.ID,
			})
		}

		// Pattern (regex)
		if c.Pattern != nil {
			matched, err := regexp.MatchString(*c.Pattern, strVal)
			if err != nil || !matched {
				result.Errors = append(result.Errors, &ValidationError{
					Code:    "CONFIG_CONSTRAINT_VIOLATION",
					Message: fmt.Sprintf("Field '%s' does not match required pattern on node '%s'", fieldSpec.Name, node.ID),
					NodeID:  node.ID,
				})
			}
		}

		// Prefix
		if c.Prefix != nil && !strings.HasPrefix(strVal, *c.Prefix) {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "CONFIG_CONSTRAINT_VIOLATION",
				Message: fmt.Sprintf("Field '%s' must start with '%s' on node '%s'", fieldSpec.Name, *c.Prefix, node.ID),
				NodeID:  node.ID,
			})
		}

		// Enum
		if len(c.Enum) > 0 {
			found := false
			for _, allowed := range c.Enum {
				if strVal == allowed {
					found = true
					break
				}
			}
			if !found {
				result.Errors = append(result.Errors, &ValidationError{
					Code:    "CONFIG_CONSTRAINT_VIOLATION",
					Message: fmt.Sprintf("Field '%s' must be one of %v on node '%s', got '%s'", fieldSpec.Name, c.Enum, node.ID, strVal),
					NodeID:  node.ID,
				})
			}
		}

		// CIDR validation
		if fieldSpec.Type == schema.FieldTypeCIDR {
			_, _, err := net.ParseCIDR(strVal)
			if err != nil {
				result.Errors = append(result.Errors, &ValidationError{
					Code:    "CONFIG_INVALID_CIDR",
					Message: fmt.Sprintf("Field '%s' is not a valid CIDR on node '%s': %s", fieldSpec.Name, node.ID, strVal),
					NodeID:  node.ID,
				})
			}
		}
	}

	// Number constraints
	var numVal float64
	switch v := value.(type) {
	case int:
		numVal = float64(v)
	case int32:
		numVal = float64(v)
	case int64:
		numVal = float64(v)
	case float64:
		numVal = v
	default:
		return // Not a number, skip number constraints
	}

	if c.MinValue != nil && numVal < *c.MinValue {
		result.Errors = append(result.Errors, &ValidationError{
			Code:    "CONFIG_CONSTRAINT_VIOLATION",
			Message: fmt.Sprintf("Field '%s' must be at least %.0f on node '%s'", fieldSpec.Name, *c.MinValue, node.ID),
			NodeID:  node.ID,
		})
	}

	if c.MaxValue != nil && numVal > *c.MaxValue {
		result.Errors = append(result.Errors, &ValidationError{
			Code:    "CONFIG_CONSTRAINT_VIOLATION",
			Message: fmt.Sprintf("Field '%s' must be at most %.0f on node '%s'", fieldSpec.Name, *c.MaxValue, node.ID),
			NodeID:  node.ID,
		})
	}
}

// validateContainmentTypes validates that parent-child relationships follow schema rules
func validateContainmentTypes(g *graph.DiagramGraph, result *ValidationResult, opts *ValidationOptions) {
	registry := schema.DefaultRegistry
	if opts != nil && opts.SchemaRegistry != nil {
		registry = opts.SchemaRegistry
	}

	provider := "aws"
	if opts != nil && opts.Provider != "" {
		provider = opts.Provider
	}

	for _, node := range g.Nodes {
		if node.ParentID == nil {
			continue // Root node, no parent to validate
		}

		parent, exists := g.Nodes[*node.ParentID]
		if !exists {
			continue // Missing parent handled elsewhere
		}

		childType := strings.ToLower(node.ResourceType)
		parentType := strings.ToLower(parent.ResourceType)

		// Get schema for child
		childSchema, childExists := registry.Get(childType, provider)
		if !childExists {
			continue // No schema, skip validation
		}

		// Check if parent type is valid for this child
		if len(childSchema.ValidParentTypes) > 0 {
			validParent := false
			for _, validType := range childSchema.ValidParentTypes {
				if validType == parentType {
					validParent = true
					break
				}
			}
			if !validParent {
				result.Warnings = append(result.Warnings, &ValidationError{
					Code:    "INVALID_CONTAINMENT",
					Message: fmt.Sprintf("Node '%s' (%s) should not be contained in '%s' (%s). Valid parents: %v", node.ID, childType, parent.ID, parentType, childSchema.ValidParentTypes),
					NodeID:  node.ID,
				})
			}
		}
	}
}

// validateCIDRs checks CIDR validity and overlaps (VPC/Subnet focused).
func validateCIDRs(g *graph.DiagramGraph, result *ValidationResult) {
	// Track VPC CIDRs by vpc node id
	vpcCIDR := make(map[string]*net.IPNet)

	// Parse VPC CIDRs
	for _, node := range g.Nodes {
		if strings.ToLower(node.ResourceType) != "vpc" {
			continue
		}
		cidrStr, ok := node.Config["cidr"].(string)
		if !ok || cidrStr == "" {
			continue // handled by validateConfigSchema
		}
		_, ipnet, err := net.ParseCIDR(cidrStr)
		if err != nil {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "CIDR_INVALID",
				Message: fmt.Sprintf("Invalid VPC CIDR '%s' on node '%s'", cidrStr, node.ID),
				NodeID:  node.ID,
			})
			continue
		}
		vpcCIDR[node.ID] = ipnet
	}

	// Collect subnet CIDRs per parent VPC
	subnetsByVPC := make(map[string][]struct {
		id   string
		cidr *net.IPNet
		raw  string
	})

	for _, node := range g.Nodes {
		if strings.ToLower(node.ResourceType) != "subnet" {
			continue
		}
		cidrStr, ok := node.Config["cidr"].(string)
		if !ok || cidrStr == "" {
			continue // handled by validateConfigSchema
		}
		_, ipnet, err := net.ParseCIDR(cidrStr)
		if err != nil {
			result.Errors = append(result.Errors, &ValidationError{
				Code:    "CIDR_INVALID",
				Message: fmt.Sprintf("Invalid subnet CIDR '%s' on node '%s'", cidrStr, node.ID),
				NodeID:  node.ID,
			})
			continue
		}

		// Associate subnet with its direct parent if it is a VPC (common pattern)
		parentVPC := ""
		if node.ParentID != nil {
			if p, ok := g.Nodes[*node.ParentID]; ok && strings.ToLower(p.ResourceType) == "vpc" {
				parentVPC = p.ID
			}
		}

		// If parent isn't a VPC, we still store under empty key, but we can't do VPC containment checks.
		subnetsByVPC[parentVPC] = append(subnetsByVPC[parentVPC], struct {
			id   string
			cidr *net.IPNet
			raw  string
		}{id: node.ID, cidr: ipnet, raw: cidrStr})

		// If we have a known parent VPC CIDR, validate subnet is within it.
		if parentVPC != "" {
			if vpcNet, ok := vpcCIDR[parentVPC]; ok {
				if !cidrWithin(vpcNet, ipnet) {
					result.Errors = append(result.Errors, &ValidationError{
						Code:    "CIDR_OUTSIDE_VPC",
						Message: fmt.Sprintf("Subnet CIDR '%s' on node '%s' is outside VPC '%s' CIDR", cidrStr, node.ID, parentVPC),
						NodeID:  node.ID,
					})
				}
			}
		}
	}

	// Overlap detection per VPC group
	for vpcID, subnets := range subnetsByVPC {
		for i := 0; i < len(subnets); i++ {
			for j := i + 1; j < len(subnets); j++ {
				if cidrOverlaps(subnets[i].cidr, subnets[j].cidr) {
					result.Errors = append(result.Errors, &ValidationError{
						Code:    "CIDR_OVERLAP",
						Message: fmt.Sprintf("Overlapping subnet CIDRs '%s' (%s) and '%s' (%s) under VPC '%s'", subnets[i].raw, subnets[i].id, subnets[j].raw, subnets[j].id, vpcID),
					})
				}
			}
		}
	}
}

func cidrOverlaps(a, b *net.IPNet) bool {
	if a == nil || b == nil {
		return false
	}
	return a.Contains(b.IP) || b.Contains(a.IP)
}

func cidrWithin(parent, child *net.IPNet) bool {
	if parent == nil || child == nil {
		return false
	}
	// Approx containment check: parent's net must contain child's network IP and child's last IP.
	if !parent.Contains(child.IP) {
		return false
	}
	last := lastIP(child)
	if last == nil {
		return false
	}
	return parent.Contains(last)
}

func lastIP(n *net.IPNet) net.IP {
	ip := n.IP.To4()
	if ip == nil {
		// IPv6 not expected in our current diagram inputs
		return nil
	}
	mask := net.IP(n.Mask).To4()
	if mask == nil {
		return nil
	}
	out := make(net.IP, len(ip))
	for i := 0; i < 4; i++ {
		out[i] = ip[i] | ^mask[i]
	}
	return out
}
