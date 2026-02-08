package architecture

import (
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/stretchr/testify/assert"
)

func TestEnrichArchitecture_EC2Fallback(t *testing.T) {
	// Ensure mapper is registered
	mapper := NewAWSResourceTypeMapper()
	architecture.RegisterResourceTypeMapper(resource.AWS, mapper)

	gen := NewAWSArchitectureGenerator()

	// Create graph
	g := &graph.DiagramGraph{
		Nodes: map[string]*graph.Node{
			"region": {
				ID: "region", ResourceType: "region", Config: map[string]interface{}{"name": "us-east-1"},
			},
			"vpc-1": {
				ID: "vpc-1", ResourceType: "VPC", ParentID: strPtr("region"),
				Config: map[string]interface{}{"cidr": "10.0.0.0/16"}, Label: "vpc",
			},
			"subnet-1": {
				ID: "subnet-1", ResourceType: "Subnet", ParentID: strPtr("vpc-1"),
				Config: map[string]interface{}{"cidr": "10.0.1.0/24"}, Label: "subnet",
			},
			"ec2-1": {
				ID: "ec2-1", ResourceType: "EC2", ParentID: strPtr("subnet-1"),
				Config: map[string]interface{}{"instanceType": "t2.micro"}, Label: "web-server",
			},
		},
		Edges:     []*graph.Edge{},
		Variables: []graph.Variable{},
		Outputs:   []graph.Output{},
	}

	arch, err := gen.Generate(g)
	assert.NoError(t, err)
	assert.NotNil(t, arch)

	foundEC2 := false
	foundSG := false
	var depID string

	for _, r := range arch.Resources {
		if r.Type.Name == "EC2" {
			foundEC2 = true
			if len(r.DependsOn) > 0 {
				depID = r.DependsOn[0]
			}
			// Check metadata
			_, hasSG := r.Metadata["securityGroupIds"]
			assert.True(t, hasSG, "EC2 should have securityGroupIds metadata set")
		}
	}

	// Iterate again to find the SG
	for _, r := range arch.Resources {
		// depID is the ID of the security group the EC2 depends on
		if r.ID == depID {
			if r.Type.Name == "SecurityGroup" {
				foundSG = true
			}
		}
	}

	assert.True(t, foundEC2, "Should find EC2 resource")
	assert.True(t, foundSG, "Should find dependent Security Group resource")
	assert.Greater(t, len(arch.Warnings), 0, "Should generate warnings")
	if len(arch.Warnings) > 0 {
		assert.Contains(t, arch.Warnings[0].Message, "missing a security group", "Warning message incorrect")
	}
}

func TestEnrichArchitecture_ASGFallback(t *testing.T) {
	// Ensure mapper is registered
	mapper := NewAWSResourceTypeMapper()
	architecture.RegisterResourceTypeMapper(resource.AWS, mapper)

	gen := NewAWSArchitectureGenerator()

	g := &graph.DiagramGraph{
		Nodes: map[string]*graph.Node{
			"region": {
				ID: "region", ResourceType: "region", Config: map[string]interface{}{"name": "us-east-1"},
			},
			"vpc-1": {
				ID: "vpc-1", ResourceType: "VPC", ParentID: strPtr("region"),
				Config: map[string]interface{}{"cidr": "10.0.0.0/16"},
			},
			"subnet-1": {
				ID: "subnet-1", ResourceType: "Subnet", ParentID: strPtr("vpc-1"),
				Config: map[string]interface{}{"cidr": "10.0.1.0/24"},
			},
			"asg-1": {
				ID: "asg-1", ResourceType: "AutoScalingGroup", ParentID: strPtr("subnet-1"),
				Config: map[string]interface{}{"minSize": 1, "maxSize": 3}, Label: "app-asg",
			},
		},
		Edges: []*graph.Edge{},
	}

	arch, err := gen.Generate(g)
	assert.NoError(t, err)
	assert.NotNil(t, arch)

	foundASG := false
	foundLT := false
	var ltID string

	for _, r := range arch.Resources {
		if r.Type.Name == "AutoScalingGroup" {
			foundASG = true
			if len(r.DependsOn) > 0 {
				ltID = r.DependsOn[0]
			}
			_, hasLT := r.Metadata["launchTemplate"]
			assert.True(t, hasLT, "ASG should have launchTemplate metadata set")
		}
	}

	for _, r := range arch.Resources {
		if r.ID == ltID {
			if r.Type.Name == "LaunchTemplate" {
				foundLT = true
			}
		}
	}

	assert.True(t, foundASG, "Should find ASG resource")
	assert.True(t, foundLT, "Should find generated Launch Template resource")
	assert.Greater(t, len(arch.Warnings), 0, "Should generate warnings")
}

func strPtr(s string) *string { return &s }
