package rules

import (
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
)

func TestAWSRuleFactory_CreateRule_ForbiddenDependencies(t *testing.T) {
	factory := NewAWSRuleFactory()

	tests := []struct {
		name           string
		resourceType   string
		constraintType string
		constraintValue string
		expectedType    rules.RuleType
	}{
		{
			name:           "forbidden dependencies rule",
			resourceType:   "Subnet",
			constraintType:  "forbidden_dependencies",
			constraintValue: "VPC,Subnet",
			expectedType:   rules.RuleTypeForbiddenDependencies,
		},
		{
			name:           "allowed dependencies rule",
			resourceType:   "Subnet",
			constraintType:  "allowed_dependencies",
			constraintValue: "RouteTable,NATGateway",
			expectedType:   rules.RuleTypeAllowedDependencies,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := factory.CreateRule(tt.resourceType, tt.constraintType, tt.constraintValue)
			if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}

			if rule.GetType() != tt.expectedType {
				t.Errorf("Expected rule type %s but got %s", tt.expectedType, rule.GetType())
			}

			if rule.GetResourceType() != tt.resourceType {
				t.Errorf("Expected resource type %s but got %s", tt.resourceType, rule.GetResourceType())
			}
		})
	}
}

func TestAWSRuleFactory_MapResourceTypeToAWS(t *testing.T) {
	factory := NewAWSRuleFactory()

	tests := []struct {
		domainType string
		expected   string
	}{
		{"VPC", "aws_vpc"},
		{"Subnet", "aws_subnet"},
		{"InternetGateway", "aws_internet_gateway"},
		{"RouteTable", "aws_route_table"},
		{"SecurityGroup", "aws_security_group"},
		{"NATGateway", "aws_nat_gateway"},
		{"EC2", "aws_instance"},
		{"EC2Instance", "aws_instance"},
		{"VirtualMachine", "aws_instance"},
		{"UnknownType", "UnknownType"}, // Should return as-is
	}

	for _, tt := range tests {
		t.Run(tt.domainType, func(t *testing.T) {
			result := factory.mapResourceTypeToAWS(tt.domainType)
			if result != tt.expected {
				t.Errorf("Expected %s but got %s", tt.expected, result)
			}
		})
	}
}
