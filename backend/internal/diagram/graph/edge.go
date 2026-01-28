package graph

// Edge represents a normalized edge/relationship in the diagram graph
type Edge struct {
	ID     string
	Source string // Source node ID
	Target string // Target node ID
	Type   string // "containment", "dependency", "reference"
}

// IsContainment returns true if the edge represents containment
func (e *Edge) IsContainment() bool {
	return e.Type == "containment"
}

// IsDependency returns true if the edge represents a dependency
func (e *Edge) IsDependency() bool {
	return e.Type == "dependency"
}

// IsReference returns true if the edge represents a reference
func (e *Edge) IsReference() bool {
	return e.Type == "reference"
}
