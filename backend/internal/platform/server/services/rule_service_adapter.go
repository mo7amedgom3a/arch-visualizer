package services

import (
	"context"
	"fmt"

	awsrules "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/rules"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/engine"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// awsRuleServiceAdapter adapts awsrules.AWSRuleService to serverinterfaces.RuleService
type awsRuleServiceAdapter struct {
	service *awsrules.AWSRuleService
}

// NewAWSRuleServiceAdapter creates a new AWS rule service adapter
func NewAWSRuleServiceAdapter() serverinterfaces.RuleService {
	return &awsRuleServiceAdapter{
		service: awsrules.NewAWSRuleService(),
	}
}

// LoadRulesWithDefaults loads rules from database constraints and merges with defaults
func (a *awsRuleServiceAdapter) LoadRulesWithDefaults(ctx context.Context, dbConstraints []serverinterfaces.ConstraintRecord) error {
	// Convert serverinterfaces.ConstraintRecord to awsrules.ConstraintRecord
	constraints := make([]awsrules.ConstraintRecord, len(dbConstraints))
	for i, c := range dbConstraints {
		constraints[i] = awsrules.ConstraintRecord{
			ResourceType:    c.ResourceType,
			ConstraintType:  c.ConstraintType,
			ConstraintValue: c.ConstraintValue,
		}
	}

	return a.service.LoadRulesWithDefaults(ctx, constraints)
}

// ValidateArchitecture validates all resources in an architecture
func (a *awsRuleServiceAdapter) ValidateArchitecture(ctx context.Context, architecture interface{}) (map[string]interface{}, error) {
	// Type assert to get the actual architecture type
	engineArch, ok := architecture.(*engine.Architecture)
	if !ok {
		return nil, fmt.Errorf("invalid architecture type, expected *engine.Architecture")
	}

	// Validate using AWS rule service
	results, err := a.service.ValidateArchitecture(ctx, engineArch)
	if err != nil {
		return nil, err
	}

	// Convert results to map[string]interface{} format
	// The AWS rule service returns map[string]*engine.EvaluationResult
	resultMap := make(map[string]interface{})
	for resID, result := range results {
		resultMap[resID] = result
	}

	return resultMap, nil
}
