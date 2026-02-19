package validator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/parser"
)

func TestInvalidDiagramFixtureReturnsMultipleErrors(t *testing.T) {
	// Read invalid fixture
	data, err := os.ReadFile("../../../pkg/usecases/json-request-diagram-invalid.json")
	if err != nil {
		t.Fatalf("failed to read invalid fixture: %v", err)
	}

	// Parse + normalize
	ir, err := parser.ParseIRDiagram(data)
	if err != nil {
		t.Fatalf("failed to parse invalid fixture: %v", err)
	}
	g, err := parser.NormalizeToGraph(ir)
	fmt.Println("Graph", g)
	if err != nil {
		t.Fatalf("failed to normalize invalid fixture: %v", err)
	}

	// Validate with a provider-aware resource type list (unit test, no DB)
	opts := &ValidationOptions{
		Provider: "aws",
		ValidResourceTypes: map[string]bool{
			"region":      true,
			"vpc":         true,
			"subnet":      true,
			"ec2":         true,
			"route-table": true,
			// intentionally exclude "not-a-real-type"
		},
	}

	result := Validate(g, opts)
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("Result", result)
	if result.Valid {
		t.Fatalf("expected invalid diagram, got valid")
	}

	// Loop on the errors and print to console
	fmt.Println("Validation Errors:")
	for i, err := range result.Errors {
		if err != nil {
			fmt.Printf("Error %d: %+v\n", i+1, err)
		}
	}

	// We expect multiple errors (not just one).
	if len(result.Errors) < 5 {
		pretty, _ := json.MarshalIndent(result, "", "  ")
		t.Fatalf("expected multiple errors (>=5), got %d\nresult=%s", len(result.Errors), string(pretty))
	}

	// Assert key error codes exist
	expectCodes := []string{
		"MISSING_PARENT",
		"CONTAINMENT_CYCLE",
		"INVALID_EDGE_SOURCE",
		"INVALID_EDGE_TARGET",
		"DEPENDENCY_SELF_LOOP",
		"DEPENDENCY_CYCLE",
		"CIDR_INVALID",
		"CIDR_OVERLAP",
		"CONFIG_MISSING_FIELD",
	}
	for _, code := range expectCodes {
		if !hasCode(result, code) {
			pretty, _ := json.MarshalIndent(result, "", "  ")
			t.Fatalf("expected error code %q to be present\nresult=%s", code, string(pretty))
		}
	}
}

func hasCode(res *ValidationResult, code string) bool {
	for _, e := range res.Errors {
		if e != nil && e.Code == code {
			return true
		}
	}
	for _, w := range res.Warnings {
		if w != nil && w.Code == code {
			return true
		}
	}
	return false
}
