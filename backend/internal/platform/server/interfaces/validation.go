package interfaces

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// ValidationService handles validation operations (rules and constraints)
type ValidationService interface {
	// ValidateArchitecture validates an architecture against rules and constraints
	ValidateArchitecture(ctx context.Context, arch *architecture.Architecture, provider resource.CloudProvider) (*RuleValidationResult, error)
}

// RuleService handles rule evaluation (provider-specific)
type RuleService interface {
	// LoadRulesWithDefaults loads rules from database constraints and merges with defaults
	LoadRulesWithDefaults(ctx context.Context, dbConstraints []ConstraintRecord) error

	// ValidateArchitecture validates all resources in an architecture
	ValidateArchitecture(ctx context.Context, architecture interface{}) (map[string]interface{}, error)
}

// ConstraintRecord represents a constraint from the database
type ConstraintRecord struct {
	ResourceType    string
	ConstraintType  string
	ConstraintValue string
}
