package validator

import (
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
)

func TestValidate(t *testing.T) {
	// Create a valid graph
	g := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Add nodes
	g.Nodes["region-1"] = &graph.Node{
		ID:           "region-1",
		Type:         "containerNode",
		ResourceType: "region",
		Label:        "Region",
		Config:       map[string]interface{}{"name": "us-east-1"},
	}

	g.Nodes["vpc-2"] = &graph.Node{
		ID:           "vpc-2",
		Type:         "containerNode",
		ResourceType: "vpc",
		Label:        "VPC",
		Config:       map[string]interface{}{"name": "project-vpc", "cidr": "10.0.0.0/16"},
		ParentID:     stringPtr("region-1"),
	}

	// Validate with options
	opts := &ValidationOptions{
		ValidResourceTypes: map[string]bool{
			"region": true,
			"vpc":    true,
		},
		Provider: "aws",
	}
	result := Validate(g, opts)
	if !result.Valid {
		t.Errorf("Expected valid graph, but got errors: %v", result.Errors)
	}
}

func TestValidateMissingParent(t *testing.T) {
	g := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Add node with missing parent
	g.Nodes["vpc-2"] = &graph.Node{
		ID:           "vpc-2",
		Type:         "containerNode",
		ResourceType: "vpc",
		Label:        "VPC",
		Config:       map[string]interface{}{"name": "project-vpc", "cidr": "10.0.0.0/16"},
		ParentID:     stringPtr("non-existent-parent"),
	}

	// Validate with options (nil opts for backward compatibility)
	result := Validate(g, nil)
	if result.Valid {
		t.Error("Expected invalid graph due to missing parent")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected at least one validation error")
	}

	foundMissingParent := false
	for _, err := range result.Errors {
		if err.Code == "MISSING_PARENT" {
			foundMissingParent = true
			break
		}
	}

	if !foundMissingParent {
		t.Error("Expected MISSING_PARENT error")
	}
}

func TestValidateContainmentCycle(t *testing.T) {
	g := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Create a cycle: node1 -> node2 -> node1
	g.Nodes["node-1"] = &graph.Node{
		ID:           "node-1",
		Type:         "containerNode",
		ResourceType: "vpc",
		Config:       map[string]interface{}{"cidr": "10.0.0.0/16"},
		ParentID:     stringPtr("node-2"),
	}

	g.Nodes["node-2"] = &graph.Node{
		ID:           "node-2",
		Type:         "containerNode",
		ResourceType: "subnet",
		Config: map[string]interface{}{
			"cidr":               "10.0.1.0/24",
			"availabilityZoneId": "us-east-1a",
		},
		ParentID:     stringPtr("node-1"),
	}

	// Validate with options
	opts := &ValidationOptions{
		ValidResourceTypes: map[string]bool{
			"vpc":    true,
			"subnet": true,
		},
		Provider: "aws",
	}
	result := Validate(g, opts)
	if result.Valid {
		t.Error("Expected invalid graph due to containment cycle")
	}

	foundCycle := false
	for _, err := range result.Errors {
		if err.Code == "CONTAINMENT_CYCLE" {
			foundCycle = true
			break
		}
	}

	if !foundCycle {
		t.Error("Expected CONTAINMENT_CYCLE error")
	}
}

func TestValidateResourceTypes(t *testing.T) {
	g := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]*graph.Edge, 0),
	}

	// Add node with valid resource type
	g.Nodes["node-1"] = &graph.Node{
		ID:           "node-1",
		Type:         "resourceNode",
		ResourceType: "ec2",
		Label:        "EC2",
		Config:       map[string]interface{}{"instanceType": "t3.micro", "ami": "ami-0123456789"},
	}

	// Add node with invalid resource type
	g.Nodes["node-2"] = &graph.Node{
		ID:           "node-2",
		Type:         "resourceNode",
		ResourceType: "unknown-type",
		Label:        "Unknown",
		Config:       map[string]interface{}{"instanceType": "t3.micro", "ami": "ami-0123456789"},
	}

	// Validate with options
	opts := &ValidationOptions{
		ValidResourceTypes: map[string]bool{
			"ec2": true,
		},
		Provider: "aws",
	}
	result := Validate(g, opts)
	// Should have warnings for unknown resource type
	if len(result.Warnings) == 0 {
		t.Error("Expected warnings for unknown resource type")
	}
}

func stringPtr(s string) *string {
	return &s
}
