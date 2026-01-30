package parser

import (
	"fmt"
	"os"
	"testing"
)

func TestParseIRDiagram(t *testing.T) {
	// Read the test JSON file (try valid version first, fallback to original)
	jsonData, err := os.ReadFile("../../../json-request-diagram-valid.json")
	if err != nil {
		// Fallback to original file
		jsonData, err = os.ReadFile("../../../json-request-diagram.json")
		if err != nil {
			t.Fatalf("Failed to read test JSON file: %v", err)
		}
	}

	// Parse the diagram
	irDiagram, err := ParseIRDiagram(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse IR diagram: %v", err)
	}

	// Verify basic structure
	if len(irDiagram.Nodes) == 0 {
		t.Error("Expected at least one node in the diagram")
	}

	// Check that we have expected nodes
	nodeIDs := make(map[string]bool)
	for _, node := range irDiagram.Nodes {
		fmt.Println(node)
		nodeIDs[node.ID] = true
	}

	expectedNodes := []string{"region-1", "vpc-2", "route-table-3", "subnet-4", "ec2-5", "security-group-6"}
	for _, expectedID := range expectedNodes {
		if !nodeIDs[expectedID] {
			t.Errorf("Expected node %s not found in parsed diagram", expectedID)
		}
	}
}

func TestNormalizeToGraph(t *testing.T) {
	// Read the test JSON file (try valid version first, fallback to original)
	jsonData, err := os.ReadFile("../../../json-request-diagram-valid.json")
	if err != nil {
		// Fallback to original file
		jsonData, err = os.ReadFile("../../../json-request-diagram.json")
		if err != nil {
			t.Fatalf("Failed to read test JSON file: %v", err)
		}
	}

	// Parse and normalize
	irDiagram, err := ParseIRDiagram(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse IR diagram: %v", err)
	}

	diagramGraph, err := NormalizeToGraph(irDiagram)
	if err != nil {
		t.Fatalf("Failed to normalize diagram: %v", err)
	}

	// Verify graph structure
	if len(diagramGraph.Nodes) == 0 {
		t.Error("Expected at least one node in the graph")
	}

	// Verify region node is present
	regionNode, hasRegion := diagramGraph.FindRegionNode()
	if !hasRegion {
		t.Error("Expected region node in graph")
	} else if regionNode.ID != "region-1" {
		t.Errorf("Expected region node ID to be 'region-1', got '%s'", regionNode.ID)
	}

	// Verify containment relationships
	vpcNode, exists := diagramGraph.GetNode("vpc-2")
	fmt.Println(vpcNode)
	if !exists {
		t.Error("Expected VPC node in graph")
	} else if vpcNode.ParentID == nil || *vpcNode.ParentID != "region-1" {
		t.Error("Expected VPC to have region-1 as parent")
	}

	// Verify visual-only filtering (if any visual-only nodes exist, they should be filtered)
	// In the test data, all nodes have isVisualOnly: false, so all should be present
	expectedNodeCount := 6 // region-1, vpc-2, route-table-3, subnet-4, ec2-5, security-group-6
	if len(diagramGraph.Nodes) != expectedNodeCount {
		t.Errorf("Expected %d nodes in graph, got %d", expectedNodeCount, len(diagramGraph.Nodes))
	}
}

func TestParseAndNormalize(t *testing.T) {
	// Read the test JSON file (try valid version first, fallback to original)
	jsonData, err := os.ReadFile("../../../json-request-diagram-valid.json")
	if err != nil {
		// Fallback to original file
		jsonData, err = os.ReadFile("../../../json-request-diagram.json")
		if err != nil {
			t.Fatalf("Failed to read test JSON file: %v", err)
		}
	}

	// Parse and normalize in one step
	diagramGraph, err := ParseAndNormalize(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse and normalize: %v", err)
	}

	// Verify result
	if diagramGraph == nil {
		t.Error("Expected non-nil graph")
	}

	if len(diagramGraph.Nodes) == 0 {
		t.Error("Expected at least one node in the graph")
	}
}

func TestNormalizeFiltersVisualOnly(t *testing.T) {
	// Create test data with visual-only nodes
	testJSON := `{
		"nodes": [
			{
				"id": "node-1",
				"type": "resourceNode",
				"data": {
					"label": "Real Node",
					"resourceType": "ec2",
					"config": {},
					"isVisualOnly": false
				}
			},
			{
				"id": "node-2",
				"type": "resourceNode",
				"data": {
					"label": "Visual Only",
					"resourceType": "ec2",
					"config": {},
					"isVisualOnly": true
				}
			}
		],
		"edges": [],
		"variables": [],
		"timestamp": 1234567890
	}`

	irDiagram, err := ParseIRDiagram([]byte(testJSON))
	if err != nil {
		t.Fatalf("Failed to parse test JSON: %v", err)
	}

	diagramGraph, err := NormalizeToGraph(irDiagram)
	if err != nil {
		t.Fatalf("Failed to normalize: %v", err)
	}

	// Visual-only nodes are tracked in the graph (for UI parity), but are filtered later
	// at the domain mapping stage when creating real infrastructure resources.
	if len(diagramGraph.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(diagramGraph.Nodes))
	}

	visualNode, exists := diagramGraph.Nodes["node-2"]
	if !exists {
		t.Error("Visual-only node should be present in the graph")
	} else if !visualNode.IsVisualOnly {
		t.Error("Visual-only node should have IsVisualOnly=true")
	}

	if _, exists := diagramGraph.Nodes["node-1"]; !exists {
		t.Error("Real node should be present")
	}
}

func TestResolveVariables(t *testing.T) {
	// Create test data with variable references
	testJSON := `{
		"nodes": [
			{
				"id": "var.region_name",
				"type": "containerNode",
				"data": {
					"label": "Region",
					"resourceType": "region",
					"config": {
						"name": "var.region_name"
					}
				}
			},
			{
				"id": "vpc-1",
				"type": "containerNode",
				"parentId": "var.region_name",
				"data": {
					"label": "VPC",
					"resourceType": "vpc",
					"config": {
						"name": "my-vpc",
						"cidr": "var.vpc_cidr"
					}
				}
			},
			{
				"id": "ec2-1",
				"type": "resourceNode",
				"parentId": "vpc-1",
				"data": {
					"label": "EC2",
					"resourceType": "ec2",
					"config": {
						"instanceType": "var.instance_type",
						"name": "web-server"
					}
				}
			}
		],
		"edges": [],
		"variables": [
			{
				"name": "region_name",
				"type": "string",
				"description": "AWS region",
				"default": "us-west-2"
			},
			{
				"name": "vpc_cidr",
				"type": "string",
				"description": "VPC CIDR block",
				"default": "10.0.0.0/16"
			},
			{
				"name": "instance_type",
				"type": "string",
				"description": "EC2 instance type",
				"default": "t3.medium"
			}
		],
		"outputs": [],
		"timestamp": 1234567890
	}`

	irDiagram, err := ParseIRDiagram([]byte(testJSON))
	if err != nil {
		t.Fatalf("Failed to parse test JSON: %v", err)
	}

	// Check that variable references were resolved
	// Region node ID should be resolved
	if irDiagram.Nodes[0].ID != "us-west-2" {
		t.Errorf("Expected region node ID to be 'us-west-2', got '%s'", irDiagram.Nodes[0].ID)
	}

	// Region config name should be resolved
	if name, ok := irDiagram.Nodes[0].Data.Config["name"].(string); !ok || name != "us-west-2" {
		t.Errorf("Expected region config name to be 'us-west-2', got '%v'", irDiagram.Nodes[0].Data.Config["name"])
	}

	// VPC parentId should be resolved
	if irDiagram.Nodes[1].ParentID == nil || *irDiagram.Nodes[1].ParentID != "us-west-2" {
		parentID := "<nil>"
		if irDiagram.Nodes[1].ParentID != nil {
			parentID = *irDiagram.Nodes[1].ParentID
		}
		t.Errorf("Expected VPC parentId to be 'us-west-2', got '%s'", parentID)
	}

	// VPC CIDR should be resolved
	if cidr, ok := irDiagram.Nodes[1].Data.Config["cidr"].(string); !ok || cidr != "10.0.0.0/16" {
		t.Errorf("Expected VPC cidr to be '10.0.0.0/16', got '%v'", irDiagram.Nodes[1].Data.Config["cidr"])
	}

	// EC2 instance type should be resolved
	if instanceType, ok := irDiagram.Nodes[2].Data.Config["instanceType"].(string); !ok || instanceType != "t3.medium" {
		t.Errorf("Expected EC2 instanceType to be 't3.medium', got '%v'", irDiagram.Nodes[2].Data.Config["instanceType"])
	}

	// Normalize and verify the graph
	diagramGraph, err := NormalizeToGraph(irDiagram)
	if err != nil {
		t.Fatalf("Failed to normalize: %v", err)
	}

	// Region node should exist with resolved ID
	if _, exists := diagramGraph.Nodes["us-west-2"]; !exists {
		t.Error("Region node with resolved ID 'us-west-2' should exist in graph")
	}

	// VPC should have resolved parent
	vpcNode, exists := diagramGraph.Nodes["vpc-1"]
	if !exists {
		t.Error("VPC node should exist in graph")
	} else if vpcNode.ParentID == nil || *vpcNode.ParentID != "us-west-2" {
		t.Error("VPC should have 'us-west-2' as parent after resolution")
	}
}

func TestResolveVariables_NoVariables(t *testing.T) {
	// Test with no variables - should not change anything
	testJSON := `{
		"nodes": [
			{
				"id": "var.undefined_var",
				"type": "resourceNode",
				"data": {
					"label": "Test",
					"resourceType": "ec2",
					"config": {}
				}
			}
		],
		"edges": [],
		"variables": [],
		"timestamp": 1234567890
	}`

	irDiagram, err := ParseIRDiagram([]byte(testJSON))
	if err != nil {
		t.Fatalf("Failed to parse test JSON: %v", err)
	}

	// Undefined variable should remain as-is
	if irDiagram.Nodes[0].ID != "var.undefined_var" {
		t.Errorf("Expected undefined variable to remain as 'var.undefined_var', got '%s'", irDiagram.Nodes[0].ID)
	}
}

func TestResolveVariables_BoolAndNumber(t *testing.T) {
	// Test with boolean and number variable types
	testJSON := `{
		"nodes": [
			{
				"id": "vpc-1",
				"type": "containerNode",
				"data": {
					"label": "VPC",
					"resourceType": "vpc",
					"config": {
						"enableDns": "var.enable_dns",
						"maxInstances": "var.max_instances"
					}
				}
			}
		],
		"edges": [],
		"variables": [
			{
				"name": "enable_dns",
				"type": "bool",
				"description": "Enable DNS",
				"default": true
			},
			{
				"name": "max_instances",
				"type": "number",
				"description": "Max instances",
				"default": 5
			}
		],
		"timestamp": 1234567890
	}`

	irDiagram, err := ParseIRDiagram([]byte(testJSON))
	if err != nil {
		t.Fatalf("Failed to parse test JSON: %v", err)
	}

	// Boolean value should be resolved (kept as bool type)
	if enableDns, ok := irDiagram.Nodes[0].Data.Config["enableDns"].(bool); !ok || enableDns != true {
		t.Errorf("Expected enableDns to be true (bool), got '%v' (%T)", irDiagram.Nodes[0].Data.Config["enableDns"], irDiagram.Nodes[0].Data.Config["enableDns"])
	}

	// Number value should be resolved (kept as float64 from JSON)
	if maxInstances, ok := irDiagram.Nodes[0].Data.Config["maxInstances"].(float64); !ok || maxInstances != 5 {
		t.Errorf("Expected maxInstances to be 5 (number), got '%v' (%T)", irDiagram.Nodes[0].Data.Config["maxInstances"], irDiagram.Nodes[0].Data.Config["maxInstances"])
	}
}
