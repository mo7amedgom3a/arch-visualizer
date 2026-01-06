package constraints

import (
	"context"
	"testing"
	
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
)

func TestRequiresParentRule(t *testing.T) {
	tests := []struct {
		name          string
		resource      *resource.Resource
		parents       []*resource.Resource
		rule          *RequiresParentRule
		expectedError bool
	}{
		{
			name: "valid parent exists",
			resource: &resource.Resource{
				ID:   "subnet-1",
				Name: "test-subnet",
				Type: resource.ResourceType{Name: "Subnet", Kind: "Subnet"},
			},
			parents: []*resource.Resource{
				{
					ID:   "vpc-1",
					Name: "test-vpc",
					Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
				},
			},
			rule:          NewRequiresParentRule("Subnet", "VPC"),
			expectedError: false,
		},
		{
			name: "no parent",
			resource: &resource.Resource{
				ID:   "subnet-1",
				Name: "test-subnet",
				Type: resource.ResourceType{Name: "Subnet", Kind: "Subnet"},
			},
			parents:       []*resource.Resource{},
			rule:         NewRequiresParentRule("Subnet", "VPC"),
			expectedError: true,
		},
		{
			name: "wrong parent type",
			resource: &resource.Resource{
				ID:   "subnet-1",
				Name: "test-subnet",
				Type: resource.ResourceType{Name: "Subnet", Kind: "Subnet"},
			},
			parents: []*resource.Resource{
				{
					ID:   "subnet-2",
					Name: "other-subnet",
					Type: resource.ResourceType{Name: "Subnet", Kind: "Subnet"},
				},
			},
			rule:          NewRequiresParentRule("Subnet", "VPC"),
			expectedError: true,
		},
		{
			name: "multiple parents with min count",
			resource: &resource.Resource{
				ID:   "resource-1",
				Name: "test-resource",
				Type: resource.ResourceType{Name: "Resource", Kind: "Resource"},
			},
			parents: []*resource.Resource{
				{
					ID:   "vpc-1",
					Name: "test-vpc",
					Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
				},
				{
					ID:   "vpc-2",
					Name: "test-vpc-2",
					Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
				},
			},
			rule: &RequiresParentRule{
				ResourceType: "Resource",
				ParentType:   "VPC",
				MinCount:     2,
			},
			expectedError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evalCtx := &rules.EvaluationContext{
				Resource: tt.resource,
				Parents:  tt.parents,
			}
			
			err := tt.rule.Evaluate(context.Background(), evalCtx)
			
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else {
					if ruleErr, ok := err.(*rules.RuleError); ok {
						if ruleErr.RuleType != rules.RuleTypeRequiresParent {
							t.Errorf("Expected RuleTypeRequiresParent, got %v", ruleErr.RuleType)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestAllowedParentRule(t *testing.T) {
	tests := []struct {
		name          string
		resource      *resource.Resource
		parents       []*resource.Resource
		rule          *AllowedParentRule
		expectedError bool
	}{
		{
			name: "valid allowed parent",
			resource: &resource.Resource{
				ID:   "ec2-1",
				Name: "test-ec2",
				Type: resource.ResourceType{Name: "EC2", Kind: "EC2"},
			},
			parents: []*resource.Resource{
				{
					ID:   "subnet-1",
					Name: "test-subnet",
					Type: resource.ResourceType{Name: "Subnet", Kind: "Subnet"},
				},
			},
			rule:          NewAllowedParentRule("EC2", []string{"Subnet"}),
			expectedError: false,
		},
		{
			name: "forbidden parent type",
			resource: &resource.Resource{
				ID:   "ec2-1",
				Name: "test-ec2",
				Type: resource.ResourceType{Name: "EC2", Kind: "EC2"},
			},
			parents: []*resource.Resource{
				{
					ID:   "vpc-1",
					Name: "test-vpc",
					Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
				},
			},
			rule:          NewAllowedParentRule("EC2", []string{"Subnet"}),
			expectedError: true,
		},
		{
			name: "multiple parents not allowed",
			resource: &resource.Resource{
				ID:   "resource-1",
				Name: "test-resource",
				Type: resource.ResourceType{Name: "Resource", Kind: "Resource"},
			},
			parents: []*resource.Resource{
				{
					ID:   "parent-1",
					Name: "parent-1",
					Type: resource.ResourceType{Name: "Parent", Kind: "Parent"},
				},
				{
					ID:   "parent-2",
					Name: "parent-2",
					Type: resource.ResourceType{Name: "Parent", Kind: "Parent"},
				},
			},
			rule: &AllowedParentRule{
				ResourceType:  "Resource",
				AllowedTypes:  []string{"Parent"},
				AllowMultiple: false,
			},
			expectedError: true,
		},
		{
			name: "multiple parents allowed",
			resource: &resource.Resource{
				ID:   "resource-1",
				Name: "test-resource",
				Type: resource.ResourceType{Name: "Resource", Kind: "Resource"},
			},
			parents: []*resource.Resource{
				{
					ID:   "parent-1",
					Name: "parent-1",
					Type: resource.ResourceType{Name: "Parent", Kind: "Parent"},
				},
				{
					ID:   "parent-2",
					Name: "parent-2",
					Type: resource.ResourceType{Name: "Parent", Kind: "Parent"},
				},
			},
			rule: &AllowedParentRule{
				ResourceType:  "Resource",
				AllowedTypes:  []string{"Parent"},
				AllowMultiple: true,
			},
			expectedError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evalCtx := &rules.EvaluationContext{
				Resource: tt.resource,
				Parents:  tt.parents,
			}
			
			err := tt.rule.Evaluate(context.Background(), evalCtx)
			
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}
