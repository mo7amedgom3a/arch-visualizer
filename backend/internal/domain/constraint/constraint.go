package constraint

// This package provides constraint-related types and utilities
// Constraints are used to define validation rules for resources
// See internal/domain/rules/ for the actual rule implementations

// ConstraintType represents the type of constraint
type ConstraintType string

const (
	ConstraintTypeRequiresParent       ConstraintType = "requires_parent"
	ConstraintTypeAllowedParent        ConstraintType = "allowed_parent"
	ConstraintTypeRequiresRegion       ConstraintType = "requires_region"
	ConstraintTypeMaxChildren          ConstraintType = "max_children"
	ConstraintTypeMinChildren          ConstraintType = "min_children"
	ConstraintTypeAllowedDependencies  ConstraintType = "allowed_dependencies"
	ConstraintTypeForbiddenDependencies ConstraintType = "forbidden_dependencies"
)

// Constraint represents a constraint definition
// This is a data structure that can be stored in the database
type Constraint struct {
	ResourceType    string
	ConstraintType  ConstraintType
	ConstraintValue string
}
