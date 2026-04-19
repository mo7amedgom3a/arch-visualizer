package rules

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/resource"
)

// AWSRuleService provides AWS-specific rule evaluation services
type AWSRuleService struct {
	registry  RuleRegistry
	factory   *AWSRuleFactory
	evaluator RuleEvaluator
}

// NewAWSRuleService creates a new AWS rule service
// Use LoadRulesWithDefaults() to load default rules and merge with DB constraints
func NewAWSRuleService() *AWSRuleService {
	return &AWSRuleService{
		registry:  NewRuleRegistry(),
		factory:   NewAWSRuleFactory(),
		evaluator: NewRuleEvaluator(),
	}
}

// LoadRulesFromConstraints loads rules from database constraints
// This merges DB constraints with code defaults (DB constraints override defaults)
func (s *AWSRuleService) LoadRulesFromConstraints(ctx context.Context, constraints []ConstraintRecord) error {
	for _, constraint := range constraints {
		rule, err := s.factory.CreateRule(
			constraint.ResourceType,
			constraint.ConstraintType,
			constraint.ConstraintValue,
		)
		if err != nil {
			return err
		}

		if err := s.registry.RegisterRule(constraint.ResourceType, rule); err != nil {
			return err
		}
	}
	return nil
}

// LoadRulesWithDefaults loads rules from database constraints and merges with defaults
// DB constraints override defaults when they have the same resource type + constraint type
func (s *AWSRuleService) LoadRulesWithDefaults(ctx context.Context, dbConstraints []ConstraintRecord) error {
	// Start with defaults
	defaultRules := DefaultNetworkingRules()
	defaultRules = append(defaultRules, DefaultComputeRules()...)
	defaultRules = append(defaultRules, DefaultStorageRules()...)
	defaultRules = append(defaultRules, DefaultDatabaseRules()...)
	defaultRules = append(defaultRules, DefaultIAMRules()...)

	// Create a map to track which default rules should be overridden
	overrideMap := make(map[string]bool)
	for _, dbConstraint := range dbConstraints {
		key := constraintKey(dbConstraint.ResourceType, dbConstraint.ConstraintType)
		overrideMap[key] = true
	}

	// Load defaults that aren't overridden
	for _, defaultRule := range defaultRules {
		key := constraintKey(defaultRule.ResourceType, defaultRule.ConstraintType)
		if !overrideMap[key] {
			rule, err := s.factory.CreateRule(
				defaultRule.ResourceType,
				defaultRule.ConstraintType,
				defaultRule.ConstraintValue,
			)
			if err != nil {
				return err
			}
			if err := s.registry.RegisterRule(defaultRule.ResourceType, rule); err != nil {
				return err
			}
		}
	}

	// Load DB constraints (these override defaults)
	return s.LoadRulesFromConstraints(ctx, dbConstraints)
}

// constraintKey creates a unique key for a constraint
func constraintKey(resourceType, constraintType string) string {
	return resourceType + ":" + constraintType
}

// ValidateResource validates a resource against all applicable rules
func (s *AWSRuleService) ValidateResource(
	ctx context.Context,
	res *resource.Resource,
	architecture *Architecture,
) (*EvaluationResult, error) {
	// Check if resource is visual-only
	if isVisualOnly, ok := res.Metadata["isVisualOnly"].(bool); ok && isVisualOnly {
		return &EvaluationResult{
			Valid:   true,
			Results: []*RuleResult{},
		}, nil
	}

	// Get rules for this resource type
	resourceRules := s.registry.GetRules(res.Type.Name)
	if len(resourceRules) == 0 {
		// No rules to evaluate
		return &EvaluationResult{
			Valid:   true,
			Results: []*RuleResult{},
		}, nil
	}

	// Build evaluation context
	evalCtx := BuildEvaluationContext(res, architecture, "aws")

	// Evaluate all rules
	result := EvaluateAllRules(ctx, s.evaluator, resourceRules, evalCtx)

	return result, nil
}

// ValidateArchitecture validates all resources in an architecture
func (s *AWSRuleService) ValidateArchitecture(
	ctx context.Context,
	architecture *Architecture,
) (map[string]*EvaluationResult, error) {
	results := make(map[string]*EvaluationResult)

	for _, res := range architecture.Resources {
		result, err := s.ValidateResource(ctx, res, architecture)
		if err != nil {
			return nil, err
		}
		results[res.ID] = result
	}

	return results, nil
}

// ConstraintRecord represents a constraint from the database
type ConstraintRecord struct {
	ResourceType    string
	ConstraintType  string
	ConstraintValue string
}
