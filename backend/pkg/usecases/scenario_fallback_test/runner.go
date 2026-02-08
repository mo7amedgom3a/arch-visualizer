package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/parser"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"

	// Register AWS generator and mapper
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/architecture"
	// Register AWS inventory (for IR type mapping)
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/inventory"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("Fallback Mechanism Test Runner")
	fmt.Println(strings.Repeat("=", 100))

	// 1. Load JSON file
	jsonPath, _ := resolvePath("../../../json-request-fallback-test.json")
	fmt.Printf("Reading JSON from: %s\n", jsonPath)

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON: %w", err)
	}

	// 2. Parse and Normalize
	fmt.Println("Parsing Diagram...")
	ir, err := parser.ParseIRDiagram(data)
	if err != nil {
		return fmt.Errorf("failed to parse IR: %w", err)
	}

	g, err := parser.NormalizeToGraph(ir)
	if err != nil {
		return fmt.Errorf("failed to normalize graph: %w", err)
	}
	fmt.Printf("✓ JSON loaded. Nodes: %d\n", len(g.Nodes))

	// 3. Map to Architecture (triggers fallback logic)
	fmt.Println("Mapping to Architecture...")
	arch, err := architecture.MapDiagramToArchitecture(g, resource.AWS)
	if err != nil {
		return fmt.Errorf("failed to map architecture: %w", err)
	}

	// 4. Inspect Results
	fmt.Printf("✓ Mapped to Architecture. Resources: %d\n", len(arch.Resources))
	fmt.Printf("  Warnings: %d\n", len(arch.Warnings))

	fmt.Println("\n--- Warnings ---")
	for _, w := range arch.Warnings {
		fmt.Printf("⚠️  [%s]: %s\n", w.ResourceID, w.Message)
	}

	fmt.Println("\n--- Generated Defaults ---")
	for _, r := range arch.Resources {
		if strings.HasPrefix(r.Name, "default-") {
			fmt.Printf("✓ Created Default Resource: %s (%s)\n", r.Name, r.Type.Name)
		}
	}

	// Verify specific fallbacks
	expectedWarnings := []string{
		"missing a security group",
		"missing a Launch Template",
		"missing an Elastic IP",
	}

	missingWarnings := []string{}
	for _, expected := range expectedWarnings {
		found := false
		for _, w := range arch.Warnings {
			if strings.Contains(w.Message, expected) {
				found = true
				break
			}
		}
		if !found {
			missingWarnings = append(missingWarnings, expected)
		}
	}

	if len(missingWarnings) > 0 {
		return fmt.Errorf("FAILED: Missing expected warnings: %v", missingWarnings)
	}

	fmt.Println("\nSUCCESS: All expected fallbacks triggered and warnings generated.")
	return nil
}

func resolvePath(relPath string) (string, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine caller")
	}
	dir := filepath.Dir(thisFile)
	return filepath.Clean(filepath.Join(dir, relPath)), nil
}
