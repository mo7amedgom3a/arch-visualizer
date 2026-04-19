package rules

import (
	"context"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/resource"
)

// mockRule is a simple mock rule for testing
type mockRule struct {
	ruleType     RuleType
	resourceType string
	value        string
	shouldFail   bool
	errorMessage string
}

func (m *mockRule) GetType() RuleType {
	return m.ruleType
}

func (m *mockRule) GetResourceType() string {
	return m.resourceType
}

func (m *mockRule) GetValue() string {
	return m.value
}

func (m *mockRule) Evaluate(ctx context.Context, evalCtx *EvaluationContext) error {
	if m.shouldFail {
		return &RuleError{
			RuleType:     m.ruleType,
			ResourceID:   evalCtx.Resource.ID,
			ResourceName: evalCtx.Resource.Name,
			ResourceType: m.resourceType,
			Message:      m.errorMessage,
			Value:        m.value,
		}
	}
	return nil
}

func TestDefaultRuleEvaluator_EvaluateRule(t *testing.T) {
	evaluator := NewRuleEvaluator()

	tests := []struct {
		name          string
		rule          Rule
		resource      *resource.Resource
		expectedPass  bool
		expectedError bool
	}{
		{
			name: "rule passes",
			rule: &mockRule{
				ruleType:     RuleTypeRequiresParent,
				resourceType: "Subnet",
				value:        "VPC",
				shouldFail:   false,
			},
			resource: &resource.Resource{
				ID:   "subnet-1",
				Name: "test-subnet",
				Type: resource.ResourceType{Name: "Subnet"},
			},
			expectedPass:  true,
			expectedError: false,
		},
		{
			name: "rule fails",
			rule: &mockRule{
				ruleType:     RuleTypeRequiresParent,
				resourceType: "Subnet",
				value:        "VPC",
				shouldFail:   true,
				errorMessage: "Subnet requires VPC parent",
			},
			resource: &resource.Resource{
				ID:   "subnet-1",
				Name: "test-subnet",
				Type: resource.ResourceType{Name: "Subnet"},
			},
			expectedPass:  false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evalCtx := &EvaluationContext{
				Resource: tt.resource,
			}

			result := evaluator.EvaluateRule(context.Background(), tt.rule, evalCtx)

			if result.Passed != tt.expectedPass {
				t.Errorf("Expected Passed=%v, got %v", tt.expectedPass, result.Passed)
			}

			if tt.expectedError {
				if result.Error == nil {
					t.Errorf("Expected error but got nil")
				} else {
					if result.Error.RuleType != tt.rule.GetType() {
						t.Errorf("Expected RuleType %v, got %v", tt.rule.GetType(), result.Error.RuleType)
					}
				}
			} else {
				if result.Error != nil {
					t.Errorf("Expected no error but got: %v", result.Error)
				}
			}
		})
	}
}

func TestDefaultRuleEvaluator_EvaluateRules(t *testing.T) {
	evaluator := NewRuleEvaluator()

	resource := &resource.Resource{
		ID:   "subnet-1",
		Name: "test-subnet",
		Type: resource.ResourceType{Name: "Subnet"},
	}

	ruleList := []Rule{
		&mockRule{ruleType: RuleTypeRequiresParent, resourceType: "Subnet", shouldFail: false},
		&mockRule{ruleType: RuleTypeRequiresRegion, resourceType: "Subnet", shouldFail: false},
		&mockRule{ruleType: RuleTypeMaxChildren, resourceType: "Subnet", shouldFail: true, errorMessage: "Too many children"},
	}

	evalCtx := &EvaluationContext{
		Resource: resource,
	}

	results := evaluator.EvaluateRules(context.Background(), ruleList, evalCtx)

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// First two should pass
	if !results[0].Passed {
		t.Errorf("Expected first rule to pass")
	}
	if !results[1].Passed {
		t.Errorf("Expected second rule to pass")
	}

	// Third should fail
	if results[2].Passed {
		t.Errorf("Expected third rule to fail")
	}
	if results[2].Error == nil {
		t.Errorf("Expected error in third result")
	}
}

func TestEvaluateAllRules(t *testing.T) {
	evaluator := NewRuleEvaluator()

	resource := &resource.Resource{
		ID:   "subnet-1",
		Name: "test-subnet",
		Type: resource.ResourceType{Name: "Subnet"},
	}

	ruleList := []Rule{
		&mockRule{ruleType: RuleTypeRequiresParent, resourceType: "Subnet", shouldFail: false},
		&mockRule{ruleType: RuleTypeRequiresRegion, resourceType: "Subnet", shouldFail: true, errorMessage: "Region required"},
		&mockRule{ruleType: RuleTypeMaxChildren, resourceType: "Subnet", shouldFail: false},
	}

	evalCtx := &EvaluationContext{
		Resource: resource,
	}

	result := EvaluateAllRules(context.Background(), evaluator, ruleList, evalCtx)

	if result.Valid {
		t.Errorf("Expected validation to fail")
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}

	if len(result.Results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result.Results))
	}
}
