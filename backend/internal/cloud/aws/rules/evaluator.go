package rules

import (
	"context"
)

// RuleEvaluator evaluates rules against resources
type RuleEvaluator interface {
	// EvaluateRule evaluates a single rule against a resource
	EvaluateRule(ctx context.Context, rule Rule, evalCtx *EvaluationContext) *RuleResult

	// EvaluateRules evaluates multiple rules against a resource
	EvaluateRules(ctx context.Context, rules []Rule, evalCtx *EvaluationContext) []*RuleResult
}

// DefaultRuleEvaluator is the default implementation of RuleEvaluator
type DefaultRuleEvaluator struct{}

// NewRuleEvaluator creates a new rule evaluator
func NewRuleEvaluator() RuleEvaluator {
	return &DefaultRuleEvaluator{}
}

// EvaluateRule evaluates a single rule
func (e *DefaultRuleEvaluator) EvaluateRule(ctx context.Context, rule Rule, evalCtx *EvaluationContext) *RuleResult {
	result := &RuleResult{
		Rule:     rule,
		Passed:   false,
		Severity: SeverityError,
	}

	err := rule.Evaluate(ctx, evalCtx)
	if err != nil {
		if ruleErr, ok := err.(*RuleError); ok {
			result.Error = ruleErr
			result.Severity = SeverityError
		} else {
			result.Error = &RuleError{
				RuleType:     rule.GetType(),
				ResourceID:   evalCtx.Resource.ID,
				ResourceName: evalCtx.Resource.Name,
				ResourceType: rule.GetResourceType(),
				Message:      err.Error(),
				Value:        rule.GetValue(),
			}
			result.Severity = SeverityError
		}
	} else {
		result.Passed = true
	}

	return result
}

// EvaluateRules evaluates multiple rules
func (e *DefaultRuleEvaluator) EvaluateRules(ctx context.Context, ruleList []Rule, evalCtx *EvaluationContext) []*RuleResult {
	results := make([]*RuleResult, len(ruleList))
	for i, rule := range ruleList {
		results[i] = e.EvaluateRule(ctx, rule, evalCtx)
	}
	return results
}

// EvaluationResult represents the overall result of evaluating all rules
type EvaluationResult struct {
	Valid    bool
	Results  []*RuleResult
	Errors   []*RuleError
	Warnings []*RuleError
	Info     []*RuleError
}

// EvaluateAllRules evaluates all rules for a resource and returns a summary
func EvaluateAllRules(
	ctx context.Context,
	evaluator RuleEvaluator,
	ruleList []Rule,
	evalCtx *EvaluationContext,
) *EvaluationResult {
	results := evaluator.EvaluateRules(ctx, ruleList, evalCtx)

	evalResult := &EvaluationResult{
		Valid:    true,
		Results:  results,
		Errors:   []*RuleError{},
		Warnings: []*RuleError{},
		Info:     []*RuleError{},
	}

	for _, result := range results {
		if !result.Passed && result.Error != nil {
			evalResult.Valid = false

			switch result.Severity {
			case SeverityError:
				evalResult.Errors = append(evalResult.Errors, result.Error)
			case SeverityWarning:
				evalResult.Warnings = append(evalResult.Warnings, result.Error)
			case SeverityInfo:
				evalResult.Info = append(evalResult.Info, result.Error)
			}
		}
	}

	return evalResult
}
