package engine

import (
	"testing"
	
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

func TestBuildEvaluationContext(t *testing.T) {
	vpcID := "vpc-1"
	subnetID := "subnet-1"
	ec2ID := "ec2-1"
	sgID := "sg-1"
	
	architecture := &Architecture{
		Resources: []*resource.Resource{
			{
				ID:       vpcID,
				Name:     "test-vpc",
				Type:     resource.ResourceType{Name: "VPC", Kind: "VPC"},
				Provider: resource.AWS,
			},
			{
				ID:       subnetID,
				Name:     "test-subnet",
				Type:     resource.ResourceType{Name: "Subnet", Kind: "Subnet"},
				Provider: resource.AWS,
				ParentID: &vpcID,
			},
			{
				ID:         ec2ID,
				Name:       "test-ec2",
				Type:       resource.ResourceType{Name: "EC2", Kind: "EC2"},
				Provider:   resource.AWS,
				ParentID:   &subnetID,
				DependsOn:  []string{sgID},
			},
			{
				ID:       sgID,
				Name:     "test-sg",
				Type:     resource.ResourceType{Name: "SecurityGroup", Kind: "SecurityGroup"},
				Provider: resource.AWS,
			},
		},
	}
	
	// Test subnet context (has VPC as parent)
	subnet := architecture.Resources[1]
	ctx := BuildEvaluationContext(subnet, architecture, "aws")
	
	if ctx.Resource.ID != subnetID {
		t.Errorf("Expected resource ID %s, got %s", subnetID, ctx.Resource.ID)
	}
	
	if len(ctx.Parents) != 1 {
		t.Errorf("Expected 1 parent, got %d", len(ctx.Parents))
	}
	
	if ctx.Parents[0].ID != vpcID {
		t.Errorf("Expected parent ID %s, got %s", vpcID, ctx.Parents[0].ID)
	}
	
	if len(ctx.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(ctx.Children))
	}
	
	if ctx.Children[0].ID != ec2ID {
		t.Errorf("Expected child ID %s, got %s", ec2ID, ctx.Children[0].ID)
	}
	
	// Test EC2 context (has subnet as parent, depends on SG)
	ec2 := architecture.Resources[2]
	ctx = BuildEvaluationContext(ec2, architecture, "aws")
	
	if len(ctx.Parents) != 1 {
		t.Errorf("Expected 1 parent, got %d", len(ctx.Parents))
	}
	
	if ctx.Parents[0].ID != subnetID {
		t.Errorf("Expected parent ID %s, got %s", subnetID, ctx.Parents[0].ID)
	}
	
	if len(ctx.Dependencies) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(ctx.Dependencies))
	}
	
	if ctx.Dependencies[0].ID != sgID {
		t.Errorf("Expected dependency ID %s, got %s", sgID, ctx.Dependencies[0].ID)
	}
	
	// Test VPC context (has subnet as child)
	vpc := architecture.Resources[0]
	ctx = BuildEvaluationContext(vpc, architecture, "aws")
	
	if len(ctx.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(ctx.Children))
	}
	
	if ctx.Children[0].ID != subnetID {
		t.Errorf("Expected child ID %s, got %s", subnetID, ctx.Children[0].ID)
	}
	
	// Test with nil architecture
	ctx = BuildEvaluationContext(subnet, nil, "aws")
	
	if ctx.Resource.ID != subnetID {
		t.Errorf("Expected resource ID %s, got %s", subnetID, ctx.Resource.ID)
	}
	
	if len(ctx.Parents) != 0 {
		t.Errorf("Expected 0 parents with nil architecture, got %d", len(ctx.Parents))
	}
}
