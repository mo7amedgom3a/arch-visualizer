package constraints

import (
	"context"
	"testing"
	
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
)

func TestAllowedDependenciesRule(t *testing.T) {
	tests := []struct {
		name          string
		resource      *resource.Resource
		dependencies  []*resource.Resource
		rule          *AllowedDependenciesRule
		expectedError bool
	}{
		{
			name: "allowed dependency",
			resource: &resource.Resource{
				ID:   "ec2-1",
				Name: "test-ec2",
				Type: resource.ResourceType{Name: "EC2", Kind: "EC2"},
			},
			dependencies: []*resource.Resource{
				{
					ID:   "sg-1",
					Name: "test-sg",
					Type: resource.ResourceType{Name: "SecurityGroup", Kind: "SecurityGroup"},
				},
			},
			rule:          NewAllowedDependenciesRule("EC2", []string{"SecurityGroup"}),
			expectedError: false,
		},
		{
			name: "forbidden dependency",
			resource: &resource.Resource{
				ID:   "ec2-1",
				Name: "test-ec2",
				Type: resource.ResourceType{Name: "EC2", Kind: "EC2"},
			},
			dependencies: []*resource.Resource{
				{
					ID:   "vpc-1",
					Name: "test-vpc",
					Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
				},
			},
			rule:          NewAllowedDependenciesRule("EC2", []string{"SecurityGroup"}),
			expectedError: true,
		},
		{
			name: "no dependencies specified",
			resource: &resource.Resource{
				ID:   "resource-1",
				Name: "test-resource",
				Type: resource.ResourceType{Name: "Resource", Kind: "Resource"},
			},
			dependencies: []*resource.Resource{},
			rule:          NewAllowedDependenciesRule("Resource", []string{"SecurityGroup"}),
			expectedError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evalCtx := &rules.EvaluationContext{
				Resource:     tt.resource,
				Dependencies: tt.dependencies,
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

func TestForbiddenDependenciesRule(t *testing.T) {
	tests := []struct {
		name          string
		resource      *resource.Resource
		dependencies  []*resource.Resource
		rule          *AllowedDependenciesRule
		expectedError bool
	}{
		{
			name: "forbidden dependency present",
			resource: &resource.Resource{
				ID:   "resource-1",
				Name: "test-resource",
				Type: resource.ResourceType{Name: "Resource", Kind: "Resource"},
			},
			dependencies: []*resource.Resource{
				{
					ID:   "forbidden-1",
					Name: "forbidden",
					Type: resource.ResourceType{Name: "Forbidden", Kind: "Forbidden"},
				},
			},
			rule:          NewForbiddenDependenciesRule("Resource", []string{"Forbidden"}),
			expectedError: true,
		},
		{
			name: "no forbidden dependencies",
			resource: &resource.Resource{
				ID:   "resource-1",
				Name: "test-resource",
				Type: resource.ResourceType{Name: "Resource", Kind: "Resource"},
			},
			dependencies: []*resource.Resource{
				{
					ID:   "allowed-1",
					Name: "allowed",
					Type: resource.ResourceType{Name: "Allowed", Kind: "Allowed"},
				},
			},
			rule:          NewForbiddenDependenciesRule("Resource", []string{"Forbidden"}),
			expectedError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evalCtx := &rules.EvaluationContext{
				Resource:     tt.resource,
				Dependencies: tt.dependencies,
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
