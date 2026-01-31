package scenario8_architecture_roundtrip

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// ArchitectureRoundtripRunner tests the complete roundtrip:
// 1. Read diagram JSON from file
// 2. Process and save to database
// 3. Load from database
// 4. Convert back to diagram JSON format
// 5. Save response JSON
func ArchitectureRoundtripRunner(ctx context.Context) error {
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("SCENARIO 8: Architecture Roundtrip (Save → Load → Convert to JSON)")
	fmt.Println(strings.Repeat("=", 100))

	// Step 1: Initialize service layer server
	fmt.Println("\n[Step 1] Initializing service layer server...")
	srv, err := server.NewServer()
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}
	fmt.Println("✓ Service layer server initialized successfully")

	// Step 2: Read IR JSON from file
	fmt.Println("\n[Step 2] Reading diagram JSON file...")
	jsonPath, err := resolveDiagramJSONPath("json-request-fiagram-complete.json")
	if err != nil {
		return fmt.Errorf("failed to resolve diagram JSON path: %w", err)
	}

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read diagram JSON file: %w", err)
	}
	fmt.Printf("✓ Read diagram JSON from: %s (%d bytes)\n", jsonPath, len(jsonData))

	// Extract diagram from project-wrapped JSON structure if needed
	diagramData, err := extractDiagramFromProjectJSON(jsonData)
	if err != nil {
		return fmt.Errorf("failed to extract diagram from project JSON: %w", err)
	}

	// Step 3: Process diagram and save to database
	fmt.Println("\n[Step 3] Processing diagram and saving to database...")

	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	processReq := &serverinterfaces.ProcessDiagramRequest{
		JSONData:      diagramData,
		UserID:        userID,
		ProjectName:   "Architecture Roundtrip Test Project",
		IACToolID:     1, // Terraform
		CloudProvider: "aws",
		Region:        "us-east-1",
	}

	result, err := srv.PipelineOrchestrator.ProcessDiagram(ctx, processReq)
	if err != nil {
		return fmt.Errorf("failed to process diagram: %w", err)
	}

	fmt.Printf("✓ Diagram processed and saved to database\n")
	fmt.Printf("  Project ID: %s\n", result.ProjectID.String())

	// Step 4: Load architecture from database
	fmt.Println("\n[Step 4] Loading architecture from database...")

	arch, err := srv.ProjectService.LoadArchitecture(ctx, result.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to load architecture: %w", err)
	}

	fmt.Printf("✓ Architecture loaded from database\n")
	fmt.Printf("  Resources: %d\n", len(arch.Resources))
	fmt.Printf("  Containments: %d\n", len(arch.Containments))
	fmt.Printf("  Dependencies: %d\n", len(arch.Dependencies))
	fmt.Printf("  Variables: %d\n", len(arch.Variables))
	fmt.Printf("  Outputs: %d\n", len(arch.Outputs))

	// Step 5: Convert architecture back to diagram JSON format
	fmt.Println("\n[Step 5] Converting architecture to diagram JSON format...")

	diagramJSON, err := convertArchitectureToDiagramJSON(arch, result.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to convert architecture to diagram JSON: %w", err)
	}

	fmt.Println("✓ Architecture converted to diagram JSON format")

	// Step 6: Save response JSON to file
	fmt.Println("\n[Step 6] Saving response JSON to file...")

	outputPath := filepath.Join(filepath.Dir(jsonPath), "json-response-architecture-loaded.json")
	outputData, err := json.MarshalIndent(diagramJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal response JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, outputData, 0o644); err != nil {
		return fmt.Errorf("failed to write response JSON file: %w", err)
	}

	fmt.Printf("✓ Response JSON saved to: %s (%d bytes)\n", outputPath, len(outputData))

	// Summary
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("SUCCESS: Architecture roundtrip completed!")
	fmt.Printf("  Project ID: %s\n", result.ProjectID.String())
	fmt.Printf("  Resources Loaded: %d\n", len(arch.Resources))
	fmt.Printf("  Response JSON: %s\n", outputPath)
	fmt.Println(strings.Repeat("=", 100))

	return nil
}

// convertArchitectureToDiagramJSON converts a domain Architecture back to the diagram JSON format
func convertArchitectureToDiagramJSON(arch *architecture.Architecture, projectID uuid.UUID) (map[string]interface{}, error) {
	// Build project key
	projectKey := fmt.Sprintf("cloud-canvas-project-%s", projectID.String())

	// Convert resources to nodes
	nodes := make([]map[string]interface{}, 0)

	// Create region node first
	regionNode := map[string]interface{}{
		"id":       "var.aws_region",
		"type":     "containerNode",
		"position": map[string]int{"x": 440, "y": 40},
		"data": map[string]interface{}{
			"label":        "Region",
			"resourceType": "region",
			"config": map[string]interface{}{
				"0":    map[string]string{"id": "us-east-1", "name": "US East (N. Virginia)"},
				"1":    map[string]string{"id": "us-east-2", "name": "US East (Ohio)"},
				"2":    map[string]string{"id": "us-west-1", "name": "US West (N. California)"},
				"3":    map[string]string{"id": "us-west-2", "name": "US West (Oregon)"},
				"4":    map[string]string{"id": "eu-west-1", "name": "EU (Ireland)"},
				"5":    map[string]string{"id": "eu-west-2", "name": "EU (London)"},
				"6":    map[string]string{"id": "eu-central-1", "name": "EU (Frankfurt)"},
				"7":    map[string]string{"id": "ap-southeast-1", "name": "Asia Pacific (Singapore)"},
				"8":    map[string]string{"id": "ap-southeast-2", "name": "Asia Pacific (Sydney)"},
				"9":    map[string]string{"id": "ap-northeast-1", "name": "Asia Pacific (Tokyo)"},
				"name": arch.Region,
			},
			"status":       "valid",
			"isVisualOnly": false,
		},
		"selectable": true,
		"focusable":  true,
		"style":      map[string]int{"width": 600, "height": 400},
		"measured":   map[string]int{"width": 600, "height": 400},
		"selected":   false,
		"dragging":   false,
	}
	nodes = append(nodes, regionNode)

	// Build resource ID to node map for parent lookups
	resourceIDToNode := make(map[string]map[string]interface{})

	// Convert resources to nodes
	for _, res := range arch.Resources {
		// Get position from metadata
		posX, posY := 0, 0
		if pos, ok := res.Metadata["position"].(map[string]interface{}); ok {
			if x, ok := pos["x"].(float64); ok {
				posX = int(x)
			} else if x, ok := pos["x"].(int); ok {
				posX = x
			}
			if y, ok := pos["y"].(float64); ok {
				posY = int(y)
			} else if y, ok := pos["y"].(int); ok {
				posY = y
			}
		}

		// Get isVisualOnly from metadata
		isVisualOnly := false
		if v, ok := res.Metadata["isVisualOnly"].(bool); ok {
			isVisualOnly = v
		}

		// Determine node type based on resource type
		nodeType := "resourceNode"
		if res.Type.Category == "Networking" && (res.Type.Kind == "Network" || res.Type.Name == "VPC" || res.Type.Name == "Subnet") {
			nodeType = "containerNode"
		}

		// Build config from metadata (excluding position and isVisualOnly)
		config := make(map[string]interface{})
		for k, v := range res.Metadata {
			if k != "position" && k != "isVisualOnly" {
				config[k] = v
			}
		}
		// Ensure name is in config
		if _, ok := config["name"]; !ok {
			config["name"] = res.Name
		}

		// Determine parent ID
		var parentID interface{} = nil
		if res.ParentID != nil {
			// Use the parent ID directly (it's already a domain resource ID)
			parentID = *res.ParentID
		}

		// Build node
		node := map[string]interface{}{
			"id":       res.ID,
			"type":     nodeType,
			"position": map[string]int{"x": posX, "y": posY},
			"data": map[string]interface{}{
				"label":        res.Name,
				"resourceType": res.Type.ID, // Use ID as resourceType (e.g., "vpc", "ec2")
				"config":       config,
				"status":       "valid",
				"isVisualOnly": isVisualOnly,
			},
			"selectable": true,
			"focusable":  true,
		}

		// Add parentId if exists
		if parentID != nil {
			node["parentId"] = parentID
			node["extent"] = "parent"
		}

		// Add style for container nodes
		if nodeType == "containerNode" {
			node["style"] = map[string]int{"width": 300, "height": 200}
			node["measured"] = map[string]int{"width": 300, "height": 200}
		} else {
			node["style"] = map[string]int{"width": 80, "height": 80}
			node["measured"] = map[string]int{"width": 80, "height": 80}
		}

		node["selected"] = false
		node["dragging"] = false

		nodes = append(nodes, node)
		resourceIDToNode[res.ID] = node
	}

	// Build edges from dependencies
	edges := make([]map[string]interface{}, 0)
	for fromID, toIDs := range arch.Dependencies {
		for _, toID := range toIDs {
			edge := map[string]interface{}{
				"id":     fmt.Sprintf("%s-%s", fromID, toID),
				"source": fromID,
				"target": toID,
				"type":   "default",
			}
			edges = append(edges, edge)
		}
	}

	// Convert variables
	variables := make([]map[string]interface{}, 0, len(arch.Variables))
	for _, v := range arch.Variables {
		variables = append(variables, map[string]interface{}{
			"name":        v.Name,
			"type":        v.Type,
			"description": v.Description,
			"default":     v.Default,
			"sensitive":   v.Sensitive,
		})
	}

	// Convert outputs
	outputs := make([]map[string]interface{}, 0, len(arch.Outputs))
	for _, o := range arch.Outputs {
		outputs = append(outputs, map[string]interface{}{
			"name":        o.Name,
			"value":       o.Value,
			"description": o.Description,
			"sensitive":   o.Sensitive,
		})
	}

	// Build the complete diagram structure
	diagramData := map[string]interface{}{
		"nodes":     nodes,
		"edges":     edges,
		"variables": variables,
		"outputs":   outputs,
		"timestamp": time.Now().UnixMilli(),
	}

	// Wrap in project key
	result := map[string]interface{}{
		projectKey: diagramData,
	}

	return result, nil
}

// resolveDiagramJSONPath resolves the JSON file path
func resolveDiagramJSONPath(filename string) (string, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine caller for resolving diagram JSON path")
	}

	dir := filepath.Dir(thisFile)
	root := filepath.Clean(filepath.Join(dir, "..", "..", ".."))
	jsonPath := filepath.Join(root, filename)

	return jsonPath, nil
}

// extractDiagramFromProjectJSON extracts the diagram structure from project-wrapped JSON
func extractDiagramFromProjectJSON(data []byte) ([]byte, error) {
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Check if this is a direct diagram format (has "nodes" at root)
	if _, hasNodes := rawData["nodes"]; hasNodes {
		return data, nil
	}

	// Otherwise, look for project-wrapped structure
	for _, value := range rawData {
		if projectData, ok := value.(map[string]interface{}); ok {
			if _, hasNodes := projectData["nodes"]; hasNodes {
				diagramBytes, err := json.Marshal(projectData)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal extracted diagram: %w", err)
				}
				return diagramBytes, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find diagram structure in JSON")
}
