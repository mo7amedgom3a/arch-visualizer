package rules

import (
	"context"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	domainrules "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/engine"
)

func TestAWSRuleService_LoadRulesWithDefaults(t *testing.T) {
	service := NewAWSRuleService()

	// Load with empty DB constraints (should use all defaults)
	err := service.LoadRulesWithDefaults(context.Background(), []ConstraintRecord{})
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Check that default rules are loaded
	defaultRules := DefaultNetworkingRules()
	for _, defaultRule := range defaultRules {
		loadedRules := service.registry.GetRules(defaultRule.ResourceType)
		found := false
		for _, rule := range loadedRules {
			if rule.GetType() == domainrules.RuleType(defaultRule.ConstraintType) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected default rule %s:%s to be loaded", defaultRule.ResourceType, defaultRule.ConstraintType)
		}
	}
}

func TestAWSRuleService_LoadRulesWithDefaults_Override(t *testing.T) {
	service := NewAWSRuleService()

	// Override a default rule
	dbConstraints := []ConstraintRecord{
		{ResourceType: "Subnet", ConstraintType: "max_children", ConstraintValue: "100"}, // Override default
	}

	err := service.LoadRulesWithDefaults(context.Background(), dbConstraints)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Check that the override is applied
	loadedRules := service.registry.GetRules("Subnet")
	foundOverride := false
	for _, rule := range loadedRules {
		if rule.GetType() == domainrules.RuleTypeMaxChildren {
			if rule.GetValue() == "100" {
				foundOverride = true
			}
		}
	}
	if !foundOverride {
		t.Error("Expected override rule to be applied")
	}
}

func TestAWSRuleFactory_CreateForbiddenDependenciesRule(t *testing.T) {
	factory := NewAWSRuleFactory()

	rule, err := factory.CreateRule("Subnet", "forbidden_dependencies", "VPC,Subnet")
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	if rule.GetType() != domainrules.RuleTypeForbiddenDependencies {
		t.Errorf("Expected rule type %s but got %s", domainrules.RuleTypeForbiddenDependencies, rule.GetType())
	}

	if rule.GetResourceType() != "Subnet" {
		t.Errorf("Expected resource type Subnet but got %s", rule.GetResourceType())
	}
}

func TestAWSRuleService_ValidateResource_WithDependencyRules(t *testing.T) {
	service := NewAWSRuleService()
	err := service.LoadRulesWithDefaults(context.Background(), []ConstraintRecord{})
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Create a VPC
	vpc := &resource.Resource{
		ID:     "vpc-1",
		Name:   "test-vpc",
		Type:   resource.ResourceType{Name: "VPC", Kind: "VPC"},
		Region: "us-east-1",
	}

	// Create a Subnet that depends on VPC (forbidden)
	subnet := &resource.Resource{
		ID:        "subnet-1",
		Name:      "test-subnet",
		Type:      resource.ResourceType{Name: "Subnet", Kind: "Subnet"},
		DependsOn: []string{"vpc-1"},
		ParentID:  stringPtr("vpc-1"),
	}

	architecture := &engine.Architecture{
		Resources: []*resource.Resource{vpc, subnet},
	}

	// Validate subnet - should fail because Subnet cannot depend on VPC
	result, err := service.ValidateResource(context.Background(), subnet, architecture)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	if result.Valid {
		t.Error("Expected validation to fail (Subnet cannot depend on VPC) but it passed")
	}

	// Check that we have a forbidden dependency error
	hasForbiddenError := false
	for _, ruleErr := range result.Errors {
		if ruleErr.RuleType == domainrules.RuleTypeForbiddenDependencies {
			hasForbiddenError = true
			break
		}
	}
	if !hasForbiddenError {
		t.Error("Expected forbidden dependency error but didn't find one")
	}
}

func TestAWSRuleService_ValidateResource_WithAllowedDependencies(t *testing.T) {
	service := NewAWSRuleService()
	err := service.LoadRulesWithDefaults(context.Background(), []ConstraintRecord{})
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Create a VPC
	vpc := &resource.Resource{
		ID:   "vpc-1",
		Name: "test-vpc",
		Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
	}

	// Create a RouteTable
	routeTable := &resource.Resource{
		ID:       "rt-1",
		Name:     "test-rt",
		Type:     resource.ResourceType{Name: "RouteTable", Kind: "RouteTable"},
		Region:   "us-east-1",
		ParentID: stringPtr("vpc-1"),
	}

	// Create a Subnet that depends on RouteTable (allowed)
	subnet := &resource.Resource{
		ID:        "subnet-1",
		Name:      "test-subnet",
		Type:      resource.ResourceType{Name: "Subnet", Kind: "Subnet"},
		Region:    "us-east-1",
		DependsOn: []string{"rt-1"},
		ParentID:  stringPtr("vpc-1"),
	}

	architecture := &engine.Architecture{
		Resources: []*resource.Resource{vpc, routeTable, subnet},
	}

	// Validate subnet - should pass because Subnet can depend on RouteTable
	result, err := service.ValidateResource(context.Background(), subnet, architecture)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected validation to pass but it failed: %v", result.Errors)
	}
}

func TestDefaultNetworkingRules(t *testing.T) {
	defaults := DefaultNetworkingRules()

	if len(defaults) == 0 {
		t.Error("Expected default rules to be defined")
	}

	// Check that we have rules for all networking resources
	expectedResources := []string{"VPC", "Subnet", "InternetGateway", "RouteTable", "SecurityGroup", "NATGateway"}
	resourceMap := make(map[string]bool)
	for _, rule := range defaults {
		resourceMap[rule.ResourceType] = true
	}

	for _, expectedResource := range expectedResources {
		if !resourceMap[expectedResource] {
			t.Errorf("Expected default rules for %s but didn't find any", expectedResource)
		}
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
