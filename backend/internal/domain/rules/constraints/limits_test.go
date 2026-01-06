package constraints

import (
	"context"
	"testing"
	
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
)

func TestMaxChildrenRule(t *testing.T) {
	tests := []struct {
		name          string
		resource      *resource.Resource
		children      []*resource.Resource
		rule          *MaxChildrenRule
		expectedError bool
	}{
		{
			name: "within max limit",
			resource: &resource.Resource{
				ID:   "vpc-1",
				Name: "test-vpc",
				Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
			},
			children: []*resource.Resource{
				{ID: "subnet-1", Type: resource.ResourceType{Name: "Subnet"}},
				{ID: "subnet-2", Type: resource.ResourceType{Name: "Subnet"}},
			},
			rule:          NewMaxChildrenRule("VPC", 5),
			expectedError: false,
		},
		{
			name: "exceeds max limit",
			resource: &resource.Resource{
				ID:   "vpc-1",
				Name: "test-vpc",
				Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
			},
			children: []*resource.Resource{
				{ID: "subnet-1", Type: resource.ResourceType{Name: "Subnet"}},
				{ID: "subnet-2", Type: resource.ResourceType{Name: "Subnet"}},
				{ID: "subnet-3", Type: resource.ResourceType{Name: "Subnet"}},
				{ID: "subnet-4", Type: resource.ResourceType{Name: "Subnet"}},
				{ID: "subnet-5", Type: resource.ResourceType{Name: "Subnet"}},
				{ID: "subnet-6", Type: resource.ResourceType{Name: "Subnet"}},
			},
			rule:          NewMaxChildrenRule("VPC", 5),
			expectedError: true,
		},
		{
			name: "exactly at max limit",
			resource: &resource.Resource{
				ID:   "vpc-1",
				Name: "test-vpc",
				Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
			},
			children: []*resource.Resource{
				{ID: "subnet-1", Type: resource.ResourceType{Name: "Subnet"}},
				{ID: "subnet-2", Type: resource.ResourceType{Name: "Subnet"}},
				{ID: "subnet-3", Type: resource.ResourceType{Name: "Subnet"}},
			},
			rule:          NewMaxChildrenRule("VPC", 3),
			expectedError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evalCtx := &rules.EvaluationContext{
				Resource: tt.resource,
				Children: tt.children,
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

func TestMinChildrenRule(t *testing.T) {
	tests := []struct {
		name          string
		resource      *resource.Resource
		children      []*resource.Resource
		rule          *MinChildrenRule
		expectedError bool
	}{
		{
			name: "meets min requirement",
			resource: &resource.Resource{
				ID:   "vpc-1",
				Name: "test-vpc",
				Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
			},
			children: []*resource.Resource{
				{ID: "subnet-1", Type: resource.ResourceType{Name: "Subnet"}},
				{ID: "subnet-2", Type: resource.ResourceType{Name: "Subnet"}},
			},
			rule:          NewMinChildrenRule("VPC", 2),
			expectedError: false,
		},
		{
			name: "below min requirement",
			resource: &resource.Resource{
				ID:   "vpc-1",
				Name: "test-vpc",
				Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
			},
			children: []*resource.Resource{
				{ID: "subnet-1", Type: resource.ResourceType{Name: "Subnet"}},
			},
			rule:          NewMinChildrenRule("VPC", 2),
			expectedError: true,
		},
		{
			name: "exactly at min requirement",
			resource: &resource.Resource{
				ID:   "vpc-1",
				Name: "test-vpc",
				Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
			},
			children: []*resource.Resource{
				{ID: "subnet-1", Type: resource.ResourceType{Name: "Subnet"}},
			},
			rule:          NewMinChildrenRule("VPC", 1),
			expectedError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evalCtx := &rules.EvaluationContext{
				Resource: tt.resource,
				Children: tt.children,
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
