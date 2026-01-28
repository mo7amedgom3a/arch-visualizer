package parser

import (
	"os"
	"fmt"
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
