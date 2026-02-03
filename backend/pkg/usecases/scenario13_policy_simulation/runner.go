package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/architecture" // Register AWS architecture generator and mappers
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/iam"   // Register IAM mappers
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/terraform"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/parser"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
	tfgen "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/generator"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	// 1) Read IR JSON from file
	jsonPath, err := resolveDiagramJSONPath("json-request-edges-s3-lambda.json")
	if err != nil {
		return err
	}
	fmt.Printf("Reading JSON from: %s\n", jsonPath)

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("read IR json %s: %w", jsonPath, err)
	}

	// 2) Parse & normalize to diagram graph
	irDiagram, err := parser.ParseIRDiagram(data)
	if err != nil {
		return fmt.Errorf("parse IR diagram: %w", err)
	}
	fmt.Printf("Parsed %d policies from IR\n", len(irDiagram.Policies))

	diagramGraph, err := parser.NormalizeToGraph(irDiagram)
	if err != nil {
		return fmt.Errorf("normalize diagram: %w", err)
	}
	fmt.Printf("Graph has %d policies\n", len(diagramGraph.Policies))

	// 3) Validate diagram
	validationResult := validator.Validate(diagramGraph, nil)
	if !validationResult.Valid {
		return fmt.Errorf("diagram validation failed:\n%s", formatValidationErrors(validationResult))
	}

	// 4) Map to domain architecture
	arch, err := architecture.MapDiagramToArchitecture(diagramGraph, resource.AWS)
	if err != nil {
		return fmt.Errorf("map diagram to architecture: %w", err)
	}
	fmt.Printf("Mapped architecture with %d resources\n", len(arch.Resources))
	for _, res := range arch.Resources {
		fmt.Printf(" - %s (%s)\n", res.Name, res.Type.Name)
	}

	// 5) Build graph and sort
	graph := architecture.NewGraph(arch)
	sorted, err := graph.GetSortedResources()
	if err != nil {
		return fmt.Errorf("topological sort failed: %w", err)
	}

	// 6) Generate Terraform
	output, err := generateTerraform(ctx, arch, sorted)
	if err != nil {
		return err
	}

	// 7) Write to output
	outDir := "terraform_output_policy_sim"
	if err := writeTerraformOutput(outDir, output); err != nil {
		return err
	}

	fmt.Printf("\nSUCCESS! Terraform code generated in ./%s\n", outDir)
	return nil
}

func generateTerraform(ctx context.Context, arch *architecture.Architecture, sorted []*resource.Resource) (*iac.Output, error) {
	// Wire Terraform mapper registry with AWS Terraform mapper.
	mapperRegistry := tfmapper.NewRegistry()
	if err := mapperRegistry.Register(terraform.New()); err != nil {
		return nil, fmt.Errorf("register aws terraform mapper: %w", err)
	}

	engine := tfgen.NewEngine(mapperRegistry)
	output, err := engine.Generate(ctx, arch, sorted)
	if err != nil {
		return nil, fmt.Errorf("terraform engine generate: %w", err)
	}

	return output, nil
}

func writeTerraformOutput(dir string, out *iac.Output) error {
	if out == nil {
		return fmt.Errorf("nil terraform output")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create output dir %s: %w", dir, err)
	}

	for _, f := range out.Files {
		target := filepath.Join(dir, f.Path)
		if err := os.WriteFile(target, []byte(f.Content), 0o644); err != nil {
			return fmt.Errorf("write file %s: %w", target, err)
		}
		fmt.Printf("Generated %s\n", target)
	}
	return nil
}

func formatValidationErrors(result *validator.ValidationResult) string {
	if result == nil || len(result.Errors) == 0 {
		return ""
	}
	var b strings.Builder
	for _, e := range result.Errors {
		b.WriteString(fmt.Sprintf("- %s: %s\n", e.Code, e.Message))
	}
	return b.String()
}

func resolveDiagramJSONPath(filename string) (string, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine caller")
	}
	dir := filepath.Dir(thisFile)
	// assuming this file is in backend/pkg/usecases/scenario13_policy_simulation/runner.go
	// and json is in backend/json-request-edges-s3-lambda.json
	// we need to go up to backend/
	root := filepath.Clean(filepath.Join(dir, "..", "..", ".."))
	jsonPath := filepath.Join(root, filename)
	return jsonPath, nil
}
