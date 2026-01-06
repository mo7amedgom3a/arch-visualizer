package rules

import (
	"context"
	
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// RuleType represents the type of rule being evaluated
type RuleType string

const (
	RuleTypeRequiresParent      RuleType = "requires_parent"
	RuleTypeAllowedParent        RuleType = "allowed_parent"
	RuleTypeRequiresRegion      RuleType = "requires_region"
	RuleTypeMaxChildren         RuleType = "max_children"
	RuleTypeMinChildren         RuleType = "min_children"
	RuleTypeAllowedDependencies RuleType = "allowed_dependencies"
	RuleTypeForbiddenDependencies RuleType = "forbidden_dependencies"
	RuleTypeRequiresTag          RuleType = "requires_tag"
	RuleTypeCIDRConstraint       RuleType = "cidr_constraint"
	RuleTypePortRange            RuleType = "port_range"
)

// Rule represents a validation rule that can be evaluated
// This is the cloud-agnostic interface that all providers must implement
type Rule interface {
	// GetType returns the type of this rule
	GetType() RuleType
	
	// GetResourceType returns the resource type this rule applies to
	GetResourceType() string
	
	// GetValue returns the rule value/configuration
	GetValue() string
	
	// Evaluate evaluates the rule against a resource in the given context
	// Returns nil if rule passes, or an error describing the violation
	Evaluate(ctx context.Context, evalCtx *EvaluationContext) error
}

// EvaluationContext provides context for rule evaluation
type EvaluationContext struct {
	// Resource being evaluated
	Resource *resource.Resource
	
	// Parent resources (for containment rules)
	Parents []*resource.Resource
	
	// Child resources (for limits rules)
	Children []*resource.Resource
	
	// Dependencies (for dependency rules)
	Dependencies []*resource.Resource
	
	// Architecture graph for complex evaluations
	//Architecture *resource.Architecture
	
	// Cloud provider context (optional, for provider-specific rules)
	Provider string
	
	// Additional metadata
	Metadata map[string]interface{}
}

// RuleError represents a rule validation error
type RuleError struct {
	RuleType     RuleType
	ResourceID   string
	ResourceName string
	ResourceType string
	Message      string
	Value        string
}

func (e *RuleError) Error() string {
	return e.Message
}

// RuleResult represents the result of evaluating a rule
type RuleResult struct {
	Rule     Rule
	Passed   bool
	Error    *RuleError
	Severity Severity
}

// Severity indicates the severity of a rule violation
type Severity string

const (
	SeverityError   Severity = "error"   // Must be fixed
	SeverityWarning Severity = "warning"  // Should be fixed
	SeverityInfo    Severity = "info"     // Informational
)
