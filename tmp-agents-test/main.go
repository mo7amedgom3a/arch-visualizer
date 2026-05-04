package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/parser"
)

func main() {
	mode := "server"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "server":
		runServer()
	case "run":
		runCLI()
	default:
		fmt.Fprintf(os.Stderr, "Usage: cloud-canvas-agents [server|run]\n")
		os.Exit(1)
	}
}

// ============================================================
// HTTP SERVER MODE — accepts POST /analyze with canvas JSON
// ============================================================

func runServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/analyze", handleAnalyze)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("[Server] Cloud Canvas Agent Server listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Determine if the body is raw frontend JSON or already normalized IR
	arch, err := parseInput(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid input: %v", err), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	orch := NewOrchestrator()
	result, err := orch.Run(ctx, arch)
	if err != nil {
		http.Error(w, fmt.Sprintf("Pipeline error: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println(result.Summary())

	w.Header().Set("Content-Type", "application/json")
	resultJSON, _ := result.ToJSON()
	w.Write(resultJSON)
}

// ============================================================
// CLI MODE — reads from stdin or file, writes to stdout
// ============================================================

func runCLI() {
	var input []byte
	var err error

	if len(os.Args) > 2 {
		// Read from file
		input, err = os.ReadFile(os.Args[2])
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", os.Args[2], err)
		}
	} else {
		// Read from stdin
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Failed to read stdin: %v", err)
		}
	}

	arch, err := parseInput(input)
	if err != nil {
		log.Fatalf("Invalid input: %v", err)
	}

	ctx := context.Background()
	orch := NewOrchestrator()
	result, err := orch.Run(ctx, arch)
	if err != nil {
		log.Fatalf("Pipeline failed: %v", err)
	}

	fmt.Println(result.Summary())

	// Write full result JSON to stdout
	resultJSON, _ := result.ToJSON()
	fmt.Println("\n--- FULL RESULT JSON ---")
	fmt.Println(string(resultJSON))

	// Write generated Terraform to a file
	if result.GeneratedCode != "" {
		if err := os.WriteFile("output/main.tf", []byte(result.GeneratedCode), 0644); err != nil {
			log.Printf("Warning: could not write main.tf: %v", err)
		} else {
			fmt.Println("\nTerraform code written to: output/main.tf")
		}
	}
}

// parseInput detects whether the input is raw frontend JSON or normalized IR
func parseInput(body []byte) (*Architecture, error) {
	// Try normalized IR first
	var arch Architecture
	if err := json.Unmarshal(body, &arch); err == nil && arch.Components != nil {
		return &arch, nil
	}

	// Fall back to frontend canvas JSON parser
	return ParseFromFrontend(body)
}

// ParseFromFrontend parses the raw frontend JSON using the core diagram parser
func ParseFromFrontend(body []byte) (*Architecture, error) {
	graph, err := parser.ParseAndNormalize(body)
	if err != nil {
		return nil, fmt.Errorf("core parser failed: %w", err)
	}

	arch := &Architecture{
		ArchitectureID:  "frontend-arch",
		Name:            "Parsed from Frontend",
		CloudProvider:   "aws", // Default
		IaCTarget:       "terraform",
		Components:      make([]Component, 0, len(graph.Nodes)),
		Edges:           make([]Edge, 0, len(graph.Edges)),
		Variables:       make([]Variable, 0, len(graph.Variables)),
		Outputs:         make([]Output, 0, len(graph.Outputs)),
	}

	// Map nodes
	for id, node := range graph.Nodes {
		compType := node.ResourceType
		if compType == "" {
			compType = node.Type
		}
		
		comp := Component{
			ID:           id,
			Type:         compType,
			Properties:   node.Config,
			IsVisualOnly: node.IsVisualOnly,
			Status:       node.Status,
		}
		if node.ParentID != nil {
			comp.ParentID = *node.ParentID
		}
		
		children := graph.GetChildren(id)
		if len(children) > 0 {
			comp.Children = make([]string, len(children))
			for i, child := range children {
				comp.Children[i] = child.ID
			}
		}
		
		if comp.Status == "" {
			comp.Status = "valid"
		}
		
		arch.Components = append(arch.Components, comp)
	}

	// Map edges
	for _, edge := range graph.Edges {
		newEdge := Edge{
			ID:     edge.ID,
			Source: edge.Source,
			Target: edge.Target,
		}
		if edge.Config != nil {
			if label, ok := edge.Config["label"].(string); ok {
				newEdge.Label = label
			}
		}
		arch.Edges = append(arch.Edges, newEdge)
	}

	// Map variables
	for _, v := range graph.Variables {
		arch.Variables = append(arch.Variables, Variable{
			Name:        v.Name,
			Type:        v.Type,
			Description: v.Description,
			Default:     v.Default,
			Sensitive:   v.Sensitive,
		})
	}

	// Map outputs
	for _, o := range graph.Outputs {
		arch.Outputs = append(arch.Outputs, Output{
			Name:        o.Name,
			Value:       o.Value,
			Description: o.Description,
			Sensitive:   o.Sensitive,
		})
	}

	return arch, nil
}
