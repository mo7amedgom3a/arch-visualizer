package registry

import (
	"testing"
	
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/constraints"
)

func TestInMemoryRuleRegistry(t *testing.T) {
	registry := NewRuleRegistry()
	
	// Test RegisterRule
	rule1 := constraints.NewRequiresParentRule("Subnet", "VPC")
	err := registry.RegisterRule("Subnet", rule1)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	rule2 := constraints.NewRequiresRegionRule("VPC", true)
	err = registry.RegisterRule("VPC", rule2)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Test GetRules
	subnetRules := registry.GetRules("Subnet")
	if len(subnetRules) != 1 {
		t.Errorf("Expected 1 rule for Subnet, got %d", len(subnetRules))
	}
	
	if subnetRules[0].GetType() != rules.RuleTypeRequiresParent {
		t.Errorf("Expected RuleTypeRequiresParent, got %v", subnetRules[0].GetType())
	}
	
	// Test GetRulesByType
	parentRuleType := rules.RuleTypeRequiresParent
	parentRules := registry.GetRulesByType(parentRuleType)
	if len(parentRules) != 1 {
		t.Errorf("Expected 1 requires_parent rule, got %d", len(parentRules))
	}
	
	// Test GetAllRules
	allRules := registry.GetAllRules()
	if len(allRules) != 2 {
		t.Errorf("Expected 2 resource types, got %d", len(allRules))
	}
	
	// Test Clear
	registry.Clear()
	subnetRules = registry.GetRules("Subnet")
	if len(subnetRules) != 0 {
		t.Errorf("Expected 0 rules after clear, got %d", len(subnetRules))
	}
}

func TestInMemoryRuleRegistry_ErrorCases(t *testing.T) {
	registry := NewRuleRegistry()
	
	// Test nil rule
	err := registry.RegisterRule("Subnet", nil)
	if err == nil {
		t.Errorf("Expected error for nil rule")
	}
	
	// Test empty resource type
	rule := constraints.NewRequiresParentRule("Subnet", "VPC")
	err = registry.RegisterRule("", rule)
	if err == nil {
		t.Errorf("Expected error for empty resource type")
	}
}

func TestDefaultRuleFactory(t *testing.T) {
	factory := NewRuleFactory()
	
	tests := []struct {
		name           string
		resourceType   string
		constraintType string
		constraintValue string
		expectedType   rules.RuleType
		expectError    bool
	}{
		{
			name:           "requires_parent",
			resourceType:   "Subnet",
			constraintType: "requires_parent",
			constraintValue: "VPC",
			expectedType:   rules.RuleTypeRequiresParent,
			expectError:    false,
		},
		{
			name:           "allowed_parent",
			resourceType:   "EC2",
			constraintType: "allowed_parent",
			constraintValue: "Subnet",
			expectedType:   rules.RuleTypeAllowedParent,
			expectError:    false,
		},
		{
			name:           "requires_region",
			resourceType:   "VPC",
			constraintType: "requires_region",
			constraintValue: "true",
			expectedType:   rules.RuleTypeRequiresRegion,
			expectError:    false,
		},
		{
			name:           "max_children",
			resourceType:   "VPC",
			constraintType: "max_children",
			constraintValue: "10",
			expectedType:   rules.RuleTypeMaxChildren,
			expectError:    false,
		},
		{
			name:           "min_children",
			resourceType:   "VPC",
			constraintType: "min_children",
			constraintValue: "1",
			expectedType:   rules.RuleTypeMinChildren,
			expectError:    false,
		},
		{
			name:           "allowed_dependencies",
			resourceType:   "EC2",
			constraintType: "allowed_dependencies",
			constraintValue: "SecurityGroup",
			expectedType:   rules.RuleTypeAllowedDependencies,
			expectError:    false,
		},
		{
			name:           "unknown type",
			resourceType:   "Resource",
			constraintType: "unknown_type",
			constraintValue: "value",
			expectError:    true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := factory.CreateRule(tt.resourceType, tt.constraintType, tt.constraintValue)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				
				if rule == nil {
					t.Errorf("Expected rule but got nil")
				}
				
				if rule.GetType() != tt.expectedType {
					t.Errorf("Expected rule type %v, got %v", tt.expectedType, rule.GetType())
				}
				
				if rule.GetResourceType() != tt.resourceType {
					t.Errorf("Expected resource type %s, got %s", tt.resourceType, rule.GetResourceType())
				}
			}
		})
	}
}
