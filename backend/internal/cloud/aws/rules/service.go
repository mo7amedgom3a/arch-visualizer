package rules

import (
	"context"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/engine"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/registry"
)

// AWSRuleService provides AWS-specific rule evaluation services
type AWSRuleService struct {
	registry    registry.RuleRegistry
	factory     *AWSRuleFactory
	evaluator   engine.RuleEvaluator
}

// NewAWSRuleService creates a new AWS rule service
func NewAWSRuleService() *AWSRuleService {
	return &AWSRuleService{
		registry:  registry.NewRuleRegistry(),
		factory:   NewAWSRuleFactory(),
		evaluator: engine.NewRuleEvaluator(),
	}
}

// LoadRulesFromConstraints loads rules from database constraints
// This is where AWS can load its specific rules
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

// ValidateResource validates a resource against all applicable rules
func (s *AWSRuleService) ValidateResource(
	ctx context.Context,
	res *resource.Resource,
	architecture *engine.Architecture,
) (*engine.EvaluationResult, error) {
	// Get rules for this resource type
	resourceRules := s.registry.GetRules(res.Type.Name)
	if len(resourceRules) == 0 {
		// No rules to evaluate
		return &engine.EvaluationResult{
			Valid:   true,
			Results: []*rules.RuleResult{},
		}, nil
	}
	
	// Build evaluation context
	evalCtx := engine.BuildEvaluationContext(res, architecture, "aws")
	
	// Evaluate all rules
	result := engine.EvaluateAllRules(ctx, s.evaluator, resourceRules, evalCtx)
	
	return result, nil
}

// ValidateArchitecture validates all resources in an architecture
func (s *AWSRuleService) ValidateArchitecture(
	ctx context.Context,
	architecture *engine.Architecture,
) (map[string]*engine.EvaluationResult, error) {
	results := make(map[string]*engine.EvaluationResult)
	
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
