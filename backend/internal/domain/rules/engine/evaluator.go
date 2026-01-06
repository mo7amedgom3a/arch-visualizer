package engine

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
)

// RuleEvaluator evaluates rules against resources
type RuleEvaluator interface {
	// EvaluateRule evaluates a single rule against a resource
	EvaluateRule(ctx context.Context, rule rules.Rule, evalCtx *rules.EvaluationContext) *rules.RuleResult

	// EvaluateRules evaluates multiple rules against a resource
	EvaluateRules(ctx context.Context, rules []rules.Rule, evalCtx *rules.EvaluationContext) []*rules.RuleResult
}

// DefaultRuleEvaluator is the default implementation of RuleEvaluator
type DefaultRuleEvaluator struct{}

// NewRuleEvaluator creates a new rule evaluator
func NewRuleEvaluator() RuleEvaluator {
	return &DefaultRuleEvaluator{}
}

// EvaluateRule evaluates a single rule
func (e *DefaultRuleEvaluator) EvaluateRule(ctx context.Context, rule rules.Rule, evalCtx *rules.EvaluationContext) *rules.RuleResult {
	result := &rules.RuleResult{
		Rule:     rule,
		Passed:   false,
		Severity: rules.SeverityError,
	}

	err := rule.Evaluate(ctx, evalCtx)
	if err != nil {
		if ruleErr, ok := err.(*rules.RuleError); ok {
			result.Error = ruleErr
			result.Severity = rules.SeverityError
		} else {
			result.Error = &rules.RuleError{
				RuleType:     rule.GetType(),
				ResourceID:   evalCtx.Resource.ID,
				ResourceName: evalCtx.Resource.Name,
				ResourceType: rule.GetResourceType(),
				Message:      err.Error(),
				Value:        rule.GetValue(),
			}
			result.Severity = rules.SeverityError
		}
	} else {
		result.Passed = true
	}

	return result
}

// EvaluateRules evaluates multiple rules
func (e *DefaultRuleEvaluator) EvaluateRules(ctx context.Context, ruleList []rules.Rule, evalCtx *rules.EvaluationContext) []*rules.RuleResult {
	results := make([]*rules.RuleResult, len(ruleList))
	for i, rule := range ruleList {
		results[i] = e.EvaluateRule(ctx, rule, evalCtx)
	}
	return results
}

// EvaluationResult represents the overall result of evaluating all rules
type EvaluationResult struct {
	Valid    bool
	Results  []*rules.RuleResult
	Errors   []*rules.RuleError
	Warnings []*rules.RuleError
	Info     []*rules.RuleError
}

// EvaluateAllRules evaluates all rules for a resource and returns a summary
func EvaluateAllRules(
	ctx context.Context,
	evaluator RuleEvaluator,
	ruleList []rules.Rule,
	evalCtx *rules.EvaluationContext,
) *EvaluationResult {
	results := evaluator.EvaluateRules(ctx, ruleList, evalCtx)

	evalResult := &EvaluationResult{
		Valid:    true,
		Results:  results,
		Errors:   []*rules.RuleError{},
		Warnings: []*rules.RuleError{},
		Info:     []*rules.RuleError{},
	}

	for _, result := range results {
		if !result.Passed && result.Error != nil {
			evalResult.Valid = false

			switch result.Severity {
			case rules.SeverityError:
				evalResult.Errors = append(evalResult.Errors, result.Error)
			case rules.SeverityWarning:
				evalResult.Warnings = append(evalResult.Warnings, result.Error)
			case rules.SeverityInfo:
				evalResult.Info = append(evalResult.Info, result.Error)
			}
		}
	}

	return evalResult
}
