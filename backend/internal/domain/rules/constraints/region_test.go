package constraints

import (
	"context"
	"testing"
	
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
)

func TestRequiresRegionRule(t *testing.T) {
	tests := []struct {
		name          string
		resource      *resource.Resource
		rule          *RequiresRegionRule
		expectedError bool
	}{
		{
			name: "region required and provided",
			resource: &resource.Resource{
				ID:     "vpc-1",
				Name:   "test-vpc",
				Region: "us-east-1",
				Type:   resource.ResourceType{Name: "VPC", Kind: "VPC"},
			},
			rule:          NewRequiresRegionRule("VPC", true),
			expectedError: false,
		},
		{
			name: "region required but not provided",
			resource: &resource.Resource{
				ID:   "vpc-1",
				Name: "test-vpc",
				Type: resource.ResourceType{Name: "VPC", Kind: "VPC"},
			},
			rule:          NewRequiresRegionRule("VPC", true),
			expectedError: true,
		},
		{
			name: "region forbidden and not provided",
			resource: &resource.Resource{
				ID:   "s3-1",
				Name: "test-bucket",
				Type: resource.ResourceType{Name: "S3", Kind: "S3"},
			},
			rule:          NewRequiresRegionRule("S3", false),
			expectedError: false,
		},
		{
			name: "region forbidden but provided",
			resource: &resource.Resource{
				ID:     "s3-1",
				Name:   "test-bucket",
				Region: "us-east-1",
				Type:   resource.ResourceType{Name: "S3", Kind: "S3"},
			},
			rule:          NewRequiresRegionRule("S3", false),
			expectedError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evalCtx := &rules.EvaluationContext{
				Resource: tt.resource,
			}
			
			err := tt.rule.Evaluate(context.Background(), evalCtx)
			
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else {
					if ruleErr, ok := err.(*rules.RuleError); ok {
						if ruleErr.RuleType != rules.RuleTypeRequiresRegion {
							t.Errorf("Expected RuleTypeRequiresRegion, got %v", ruleErr.RuleType)
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
