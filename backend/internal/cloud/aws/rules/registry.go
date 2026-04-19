package rules

import (
	"fmt"
	"strconv"
	"strings"
)

// RuleRegistry stores and retrieves rules by resource type
type RuleRegistry interface {
	// RegisterRule registers a rule for a resource type
	RegisterRule(resourceType string, rule Rule) error

	// GetRules returns all rules for a resource type
	GetRules(resourceType string) []Rule

	// GetRulesByType returns all rules of a specific type
	GetRulesByType(ruleType RuleType) []Rule

	// GetAllRules returns all registered rules
	GetAllRules() map[string][]Rule

	// Clear removes all registered rules
	Clear()
}

// InMemoryRuleRegistry is an in-memory implementation of RuleRegistry
type InMemoryRuleRegistry struct {
	rules map[string][]Rule // resourceType -> []Rule
}

// NewRuleRegistry creates a new in-memory rule registry
func NewRuleRegistry() RuleRegistry {
	return &InMemoryRuleRegistry{
		rules: make(map[string][]Rule),
	}
}

// RegisterRule registers a rule for a resource type
func (r *InMemoryRuleRegistry) RegisterRule(resourceType string, rule Rule) error {
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
func (r *InMemoryRuleRegistry) GetRules(resourceType string) []Rule {
	return r.rules[resourceType]
}

// GetRulesByType returns all rules of a specific type
func (r *InMemoryRuleRegistry) GetRulesByType(ruleType RuleType) []Rule {
	var result []Rule
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
func (r *InMemoryRuleRegistry) GetAllRules() map[string][]Rule {
	// Return a copy to prevent external modification
	result := make(map[string][]Rule)
	for k, v := range r.rules {
		result[k] = make([]Rule, len(v))
		copy(result[k], v)
	}
	return result
}

// Clear removes all registered rules
func (r *InMemoryRuleRegistry) Clear() {
	r.rules = make(map[string][]Rule)
}

// RuleFactory creates rules from database constraint records
// This allows cloud providers to implement their own rule creation logic
type RuleFactory interface {
	// CreateRule creates a rule from constraint data
	CreateRule(resourceType string, constraintType string, constraintValue string) (Rule, error)
}

// DefaultRuleFactory is a default implementation that creates standard rules
type DefaultRuleFactory struct{}

// NewRuleFactory creates a new rule factory
func NewRuleFactory() RuleFactory {
	return &DefaultRuleFactory{}
}

// CreateRule creates a rule from constraint data
func (f *DefaultRuleFactory) CreateRule(resourceType string, constraintType string, constraintValue string) (Rule, error) {
	ruleType := RuleType(constraintType)

	switch ruleType {
	case RuleTypeRequiresParent:
		return NewRequiresParentRule(resourceType, constraintValue), nil
	case RuleTypeAllowedParent:
		// Parse comma-separated list
		allowedTypes := parseCommaSeparated(constraintValue)
		return NewAllowedParentRule(resourceType, allowedTypes), nil
	case RuleTypeRequiresRegion:
		required := constraintValue == "true"
		return NewRequiresRegionRule(resourceType, required), nil
	case RuleTypeMaxChildren:
		maxCount := parseInt(constraintValue)
		return NewMaxChildrenRule(resourceType, maxCount), nil
	case RuleTypeMinChildren:
		minCount := parseInt(constraintValue)
		return NewMinChildrenRule(resourceType, minCount), nil
	case RuleTypeAllowedDependencies:
		allowedTypes := parseCommaSeparated(constraintValue)
		return NewAllowedDependenciesRule(resourceType, allowedTypes), nil
	case RuleTypeForbiddenDependencies:
		forbiddenTypes := parseCommaSeparated(constraintValue)
		return NewForbiddenDependenciesRule(resourceType, forbiddenTypes), nil
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
