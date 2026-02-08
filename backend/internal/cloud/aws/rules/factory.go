package rules

import (
	"strconv"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/constraints"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/registry"
)

// AWSRuleFactory creates AWS-specific rules
// This allows AWS to implement rules in its own way while following domain interfaces
type AWSRuleFactory struct {
	registry.RuleFactory
}

// NewAWSRuleFactory creates a new AWS rule factory
func NewAWSRuleFactory() *AWSRuleFactory {
	return &AWSRuleFactory{
		RuleFactory: registry.NewRuleFactory(),
	}
}

// CreateRule creates an AWS-specific rule from constraint data
// AWS can override default behavior or add AWS-specific rules
func (f *AWSRuleFactory) CreateRule(resourceType string, constraintType string, constraintValue string) (rules.Rule, error) {
	ruleType := rules.RuleType(constraintType)

	// AWS-specific rule handling
	switch ruleType {
	case rules.RuleTypeRequiresParent:
		// AWS might have different parent requirements
		return f.createRequiresParentRule(resourceType, constraintValue)
	case rules.RuleTypeAllowedParent:
		return f.createAllowedParentRule(resourceType, constraintValue)
	case rules.RuleTypeRequiresRegion:
		return f.createRequiresRegionRule(resourceType, constraintValue)
	case rules.RuleTypeMaxChildren:
		return f.createMaxChildrenRule(resourceType, constraintValue)
	case rules.RuleTypeMinChildren:
		return f.createMinChildrenRule(resourceType, constraintValue)
	case rules.RuleTypeAllowedDependencies:
		return f.createAllowedDependenciesRule(resourceType, constraintValue)
	case rules.RuleTypeForbiddenDependencies:
		return f.createForbiddenDependenciesRule(resourceType, constraintValue)
	case rules.RuleTypeRequiresDependency:
		return f.createRequiresDependencyRule(resourceType, constraintValue)
	default:
		// Fall back to default factory for unknown types
		return f.RuleFactory.CreateRule(resourceType, constraintType, constraintValue)
	}
}

// AWS-specific rule creation methods
func (f *AWSRuleFactory) createRequiresParentRule(resourceType, parentType string) (rules.Rule, error) {
	// Use domain resource type names for parent relationships so rules work directly
	// with the domain architecture (e.g., "VPC", "Subnet").
	return constraints.NewRequiresParentRule(resourceType, parentType), nil
}

func (f *AWSRuleFactory) createAllowedParentRule(resourceType, constraintValue string) (rules.Rule, error) {
	allowedTypes := f.parseCommaSeparated(constraintValue)
	// Keep allowed parent types in domain form (e.g., "VPC") so they match
	// resource.Type.Name directly.
	return constraints.NewAllowedParentRule(resourceType, allowedTypes), nil
}

func (f *AWSRuleFactory) createRequiresRegionRule(resourceType, constraintValue string) (rules.Rule, error) {
	required := constraintValue == "true"
	// AWS-specific: Some resources are always regional, some are global
	return constraints.NewRequiresRegionRule(resourceType, required), nil
}

func (f *AWSRuleFactory) createMaxChildrenRule(resourceType, constraintValue string) (rules.Rule, error) {
	maxCount := f.parseInt(constraintValue)
	// AWS might have different limits
	return constraints.NewMaxChildrenRule(resourceType, maxCount), nil
}

func (f *AWSRuleFactory) createMinChildrenRule(resourceType, constraintValue string) (rules.Rule, error) {
	minCount := f.parseInt(constraintValue)
	return constraints.NewMinChildrenRule(resourceType, minCount), nil
}

func (f *AWSRuleFactory) createAllowedDependenciesRule(resourceType, constraintValue string) (rules.Rule, error) {
	allowedTypes := f.parseCommaSeparated(constraintValue)
	// Use domain type names (e.g., "RouteTable", "NATGateway") so dependency
	// checks operate on the same names as the domain architecture.
	return constraints.NewAllowedDependenciesRule(resourceType, allowedTypes), nil
}

func (f *AWSRuleFactory) createForbiddenDependenciesRule(resourceType, constraintValue string) (rules.Rule, error) {
	forbiddenTypes := f.parseCommaSeparated(constraintValue)
	// Use domain type names (e.g., "VPC", "Subnet") for forbidden dependencies.
	return constraints.NewForbiddenDependenciesRule(resourceType, forbiddenTypes), nil
}

func (f *AWSRuleFactory) createRequiresDependencyRule(resourceType, constraintValue string) (rules.Rule, error) {
	requiredType := constraintValue
	// Use domain type names
	return constraints.NewRequiresDependencyRule(resourceType, requiredType), nil
}

// mapResourceTypeToAWS maps domain resource types to AWS-specific types
// This allows AWS to implement rules using its own naming conventions
func (f *AWSRuleFactory) mapResourceTypeToAWS(domainType string) string {
	// Mapping from domain types to AWS types
	mapping := map[string]string{
		"VPC":             "aws_vpc",
		"Subnet":          "aws_subnet",
		"InternetGateway": "aws_internet_gateway",
		"RouteTable":      "aws_route_table",
		"SecurityGroup":   "aws_security_group",
		"NATGateway":      "aws_nat_gateway",
		"EC2":             "aws_instance",
		"EC2Instance":     "aws_instance",
		"VirtualMachine":  "aws_instance",
	}

	if awsType, ok := mapping[domainType]; ok {
		return awsType
	}

	// Default: assume it's already an AWS type or return as-is
	return domainType
}

// Helper functions
func (f *AWSRuleFactory) parseCommaSeparated(value string) []string {
	if value == "" {
		return []string{}
	}
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

func (f *AWSRuleFactory) parseInt(value string) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return result
}
