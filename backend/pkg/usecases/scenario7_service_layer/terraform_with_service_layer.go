package scenario7_service_layer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// TerraformWithServiceLayerRunner demonstrates the end-to-end pipeline using the service layer.
// This use case shows how to use the service layer orchestrator to process diagrams,
// generate Terraform code, and persist to the database.
//
// Steps:
//  1. Initialize the service layer server
//  2. Read diagram JSON from file
//  3. Use PipelineOrchestrator.ProcessDiagram() to process diagram (parse, validate, map, validate rules, persist)
//  4. Use PipelineOrchestrator.GenerateCode() to load architecture from database and generate Terraform code
//  5. Write Terraform files to ./terraform_output/
func TerraformWithServiceLayerRunner(ctx context.Context) error {
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("SCENARIO 7: Terraform Code Generation with Service Layer")
	fmt.Println(strings.Repeat("=", 100))

	// Step 1: Initialize service layer server
	fmt.Println("\n[Step 1] Initializing service layer server...")
	srv, err := server.NewServer(slog.Default())
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

	// Step 3: Process diagram using PipelineOrchestrator
	fmt.Println("\n[Step 3] Processing diagram through service layer orchestrator...")

	// Create demo user ID
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// Process diagram request
	processReq := &serverinterfaces.ProcessDiagramRequest{
		JSONData:      diagramData,
		UserID:        userID,
		ProjectName:   "Service Layer Architecture Project",
		IACToolID:     1, // Terraform
		CloudProvider: "aws",
		Region:        "us-east-1",
	}

	result, err := srv.PipelineOrchestrator.ProcessDiagram(ctx, processReq)
	if err != nil {
		return fmt.Errorf("failed to process diagram: %w", err)
	}

	fmt.Printf("✓ Diagram processed successfully\n")
	fmt.Printf("  Project ID: %s\n", result.ProjectID.String())
	fmt.Printf("  Message: %s\n", result.Message)

	// Step 4: Generate Terraform code using PipelineOrchestrator
	fmt.Println("\n[Step 4] Generating Terraform code from persisted project...")

	// Get project to verify it was created and get provider info
	project, err := srv.ProjectService.GetByID(ctx, result.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	fmt.Printf("✓ Project retrieved: %s (Provider: %s, Region: %s)\n",
		project.Name, project.CloudProvider, project.Region)

	// Use PipelineOrchestrator.GenerateCode to load architecture and generate code
	generateReq := &serverinterfaces.GenerateCodeRequest{
		ProjectID:     result.ProjectID,
		Engine:        "terraform",
		CloudProvider: project.CloudProvider,
	}

	output, err := srv.PipelineOrchestrator.GenerateCode(ctx, generateReq)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}
	fmt.Printf("✓ Generated Terraform code: %d files\n", len(output.Files))

	// Load architecture to get resource count for summary
	arch, err := srv.ProjectService.LoadArchitecture(ctx, result.ProjectID)
	if err != nil {
		// Log warning but continue - architecture loading is optional for summary
		fmt.Printf("⚠ Warning: failed to load architecture for summary: %v\n", err)
		arch = nil
	}

	// Step 5: Write Terraform files to disk
	fmt.Println("\n[Step 5] Writing Terraform files to disk...")
	outDir := "terraform_output"
	if err := writeTerraformOutput(outDir, output); err != nil {
		return fmt.Errorf("failed to write Terraform output: %w", err)
	}

	fmt.Printf("✓ Terraform files written to: ./%s/\n", outDir)
	for _, f := range output.Files {
		fmt.Printf("    - %s (%d bytes)\n", f.Path, len(f.Content))
	}

	// Summary
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("SUCCESS: Complete pipeline executed using service layer!")
	fmt.Printf("  Project ID: %s\n", result.ProjectID.String())
	fmt.Printf("  Project Name: %s\n", project.Name)
	fmt.Printf("  Cloud Provider: %s\n", project.CloudProvider)
	fmt.Printf("  Region: %s\n", project.Region)
	if arch != nil {
		fmt.Printf("  Resources Persisted: %d\n", len(arch.Resources))
	}
	fmt.Printf("  Terraform Files Generated: %d\n", len(output.Files))
	fmt.Printf("  Output Directory: ./%s/\n", outDir)
	fmt.Println(strings.Repeat("=", 100))

	return nil
}

// resolveDiagramJSONPath resolves the JSON file path
// relative to the backend module root, regardless of the current working dir.
func resolveDiagramJSONPath(filename string) (string, error) {
	// Use runtime.Caller to get this file's directory, then walk up to the backend root.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine caller for resolving diagram JSON path")
	}

	// thisFile = .../backend/pkg/usecases/scenario7_service_layer/terraform_with_service_layer.go
	// backend root = thisFile/../../../..
	dir := filepath.Dir(thisFile)
	root := filepath.Clean(filepath.Join(dir, "..", "..", ".."))
	jsonPath := filepath.Join(root, filename)

	return jsonPath, nil
}

// extractDiagramFromProjectJSON extracts the diagram structure from project-wrapped JSON.
// Handles both formats:
//   - Direct format: {"nodes": [...], "edges": [...]}
//   - Project-wrapped: {"cloud-canvas-project-...": {"nodes": [...], "edges": [...]}}
func extractDiagramFromProjectJSON(data []byte) ([]byte, error) {
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Check if this is a direct diagram format (has "nodes" at root)
	if _, hasNodes := rawData["nodes"]; hasNodes {
		// Already in the correct format, return as-is
		return data, nil
	}

	// Otherwise, look for project-wrapped structure
	// Find the first key that contains a nested object with "nodes"
	for _, value := range rawData {
		if projectData, ok := value.(map[string]interface{}); ok {
			if _, hasNodes := projectData["nodes"]; hasNodes {
				// Found the diagram structure, extract it
				diagramBytes, err := json.Marshal(projectData)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal extracted diagram: %w", err)
				}
				return diagramBytes, nil
			}
		}
	}

	// If we get here, couldn't find the diagram structure
	return nil, fmt.Errorf("could not find diagram structure in JSON (expected 'nodes' field at root or nested under project key)")
}

// writeTerraformOutput writes the generated Terraform files to disk
func writeTerraformOutput(dir string, output *iac.Output) error {
	if output == nil {
		return fmt.Errorf("terraform output is nil")
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, f := range output.Files {
		target := filepath.Join(dir, f.Path)
		if err := os.WriteFile(target, []byte(f.Content), 0o644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", target, err)
		}
	}

	return nil
}
