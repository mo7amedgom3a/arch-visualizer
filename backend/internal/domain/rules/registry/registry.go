package registry

import (
	"fmt"
	"strconv"
	"strings"
	
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/constraints"
)

// RuleRegistry stores and retrieves rules by resource type
type RuleRegistry interface {
	// RegisterRule registers a rule for a resource type
	RegisterRule(resourceType string, rule rules.Rule) error
	
	// GetRules returns all rules for a resource type
	GetRules(resourceType string) []rules.Rule
	
	// GetRulesByType returns all rules of a specific type
	GetRulesByType(ruleType rules.RuleType) []rules.Rule
	
	// GetAllRules returns all registered rules
	GetAllRules() map[string][]rules.Rule
	
	// Clear removes all registered rules
	Clear()
}

// InMemoryRuleRegistry is an in-memory implementation of RuleRegistry
type InMemoryRuleRegistry struct {
	rules map[string][]rules.Rule // resourceType -> []Rule
}

// NewRuleRegistry creates a new in-memory rule registry
func NewRuleRegistry() RuleRegistry {
	return &InMemoryRuleRegistry{
		rules: make(map[string][]rules.Rule),
	}
}

// RegisterRule registers a rule for a resource type
func (r *InMemoryRuleRegistry) RegisterRule(resourceType string, rule rules.Rule) error {
	if rule == nil {
		return fmt.Errorf("cannot register nil rule")
	}
	
	if resourceType == "" {
		return fmt.Errorf("resource type cannot be empty")
	}
	
	r.rules[resourceType] = append(r.rules[resourceType], rule)
	return nil
}

// GetRules returns all rules for a resource type
func (r *InMemoryRuleRegistry) GetRules(resourceType string) []rules.Rule {
	return r.rules[resourceType]
}

// GetRulesByType returns all rules of a specific type
func (r *InMemoryRuleRegistry) GetRulesByType(ruleType rules.RuleType) []rules.Rule {
	var result []rules.Rule
	for _, ruleList := range r.rules {
		for _, rule := range ruleList {
			if rule.GetType() == ruleType {
				result = append(result, rule)
			}
		}
	}
	return result
}

// GetAllRules returns all registered rules
func (r *InMemoryRuleRegistry) GetAllRules() map[string][]rules.Rule {
	// Return a copy to prevent external modification
	result := make(map[string][]rules.Rule)
	for k, v := range r.rules {
		result[k] = make([]rules.Rule, len(v))
		copy(result[k], v)
	}
	return result
}

// Clear removes all registered rules
func (r *InMemoryRuleRegistry) Clear() {
	r.rules = make(map[string][]rules.Rule)
}

// RuleFactory creates rules from database constraint records
// This allows cloud providers to implement their own rule creation logic
type RuleFactory interface {
	// CreateRule creates a rule from constraint data
	CreateRule(resourceType string, constraintType string, constraintValue string) (rules.Rule, error)
}

// DefaultRuleFactory is a default implementation that creates standard rules
type DefaultRuleFactory struct{}

// NewRuleFactory creates a new rule factory
func NewRuleFactory() RuleFactory {
	return &DefaultRuleFactory{}
}

// CreateRule creates a rule from constraint data
func (f *DefaultRuleFactory) CreateRule(resourceType string, constraintType string, constraintValue string) (rules.Rule, error) {
	ruleType := rules.RuleType(constraintType)
	
	switch ruleType {
	case rules.RuleTypeRequiresParent:
		return constraints.NewRequiresParentRule(resourceType, constraintValue), nil
	case rules.RuleTypeAllowedParent:
		// Parse comma-separated list
		allowedTypes := parseCommaSeparated(constraintValue)
		return constraints.NewAllowedParentRule(resourceType, allowedTypes), nil
	case rules.RuleTypeRequiresRegion:
		required := constraintValue == "true"
		return constraints.NewRequiresRegionRule(resourceType, required), nil
	case rules.RuleTypeMaxChildren:
		maxCount := parseInt(constraintValue)
		return constraints.NewMaxChildrenRule(resourceType, maxCount), nil
	case rules.RuleTypeMinChildren:
		minCount := parseInt(constraintValue)
		return constraints.NewMinChildrenRule(resourceType, minCount), nil
	case rules.RuleTypeAllowedDependencies:
		allowedTypes := parseCommaSeparated(constraintValue)
		return constraints.NewAllowedDependenciesRule(resourceType, allowedTypes), nil
	case rules.RuleTypeForbiddenDependencies:
		forbiddenTypes := parseCommaSeparated(constraintValue)
		return constraints.NewForbiddenDependenciesRule(resourceType, forbiddenTypes), nil
	default:
		return nil, fmt.Errorf("unknown constraint type: %s", constraintType)
	}
}

// Helper functions
func parseCommaSeparated(value string) []string {
	if value == "" {
		return []string{}
	}
	// Simple parsing - can be enhanced
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func parseInt(value string) int {
	// Simple parsing - can be enhanced with proper error handling
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return result
}
