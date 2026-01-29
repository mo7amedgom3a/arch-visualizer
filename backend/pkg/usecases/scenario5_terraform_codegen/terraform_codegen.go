package scenario5_terraform_codegen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/terraform"
	awsrules "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/rules"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/parser"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/validator"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	rulesengine "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/engine"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
	tfgen "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/generator"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

// This use case demonstrates the end-to-end pipeline from a diagram IR JSON
// to generated Terraform files using the internal codegen and IaC layers.
//
// Steps:
//  1. Parse IR JSON into diagram graph
//  2. Validate diagram (structure, schemas, relationships)
//  3. Map to domain Architecture aggregate
//  4. Validate domain rules/constraints (AWS networking defaults)
//  5. Build domain graph + topologically sort resources
//  6. Run Terraform engine to produce IaC files
//  7. Write Terraform files to ./terraform_output/
func TerraformCodegenRunner(ctx context.Context) error {
	// 1) Read IR JSON from file
	jsonPath, err := resolveDiagramJSONPath("json-request-diagram-valid.json")
	if err != nil {
		return err
	}
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Println(strings.Repeat("=", 100))
		fmt.Println("Error", err)
		return fmt.Errorf("read IR json %s: %w", jsonPath, err)
	}

	// 2) Parse & normalize to diagram graph
	irDiagram, err := parser.ParseIRDiagram(data)
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("IR Diagram", irDiagram)
	if err != nil {
		return fmt.Errorf("parse IR diagram: %w", err)
	}

	diagramGraph, err := parser.NormalizeToGraph(irDiagram)
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("Diagram Graph", diagramGraph)
	if err != nil {
		return fmt.Errorf("normalize diagram: %w", err)
	}

	// 3) Validate diagram (structure + schemas + relationships)
	validationResult := validator.Validate(diagramGraph, nil)
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("Validation Result", validationResult)
	if !validationResult.Valid {
		return fmt.Errorf("diagram validation failed:\n%s", formatValidationErrors(validationResult))
	}

	// 4) Map to domain architecture (cloud-agnostic core)
	arch, err := architecture.MapDiagramToArchitecture(diagramGraph, resource.AWS)
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("Architecture", arch)
	if err != nil {

		return fmt.Errorf("map diagram to architecture: %w", err)
	}

	// Basic domain validation hook (currently soft checks).
	if err := arch.Validate(); err != nil {
		return fmt.Errorf("architecture validation failed: %w", err)
	}

	// 5) Validate domain rules/constraints using AWS default networking rules
	if err := validateAWSRules(ctx, arch); err != nil {
		fmt.Println(err)
		return err
	}

	// 6) Build domain graph + topologically sort resources
	graph := architecture.NewGraph(arch)
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("Graph", graph)
	sorted, err := graph.GetSortedResources()
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("Sorted", sorted)
	if err != nil {
		return fmt.Errorf("topological sort failed: %w", err)
	}

	// 7) Generate Terraform using Terraform engine (HCL generation)
	output, err := generateTerraform(ctx, arch, sorted)
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("Output", output)
	if err != nil {
		return err
	}

	// 8) Write generated files to disk
	outDir := "terraform_output"
	if err := writeTerraformOutput(outDir, output); err != nil {
		return err
	}

	fmt.Printf("Terraform code generated in ./%s:\n", outDir)
	for _, f := range output.Files {
		fmt.Printf("- %s\n", filepath.Join(outDir, f.Path))
	}

	return nil
}

func validateAWSRules(ctx context.Context, arch *architecture.Architecture) error {
	ruleService := awsrules.NewAWSRuleService()

	// Start with code-defined defaults; no DB overrides for this use case.
	if err := ruleService.LoadRulesWithDefaults(ctx, nil); err != nil {
		return fmt.Errorf("load AWS default rules: %w", err)
	}

	// Adapt domain architecture to rules engine architecture view.
	engineArch := &rulesengine.Architecture{
		Resources: arch.Resources,
	}

	results, err := ruleService.ValidateArchitecture(ctx, engineArch)
	if err != nil {
		return fmt.Errorf("validate architecture rules: %w", err)
	}

	var messages []string
	for resID, resResult := range results {
		if !resResult.Valid {
			for _, re := range resResult.Errors {
				messages = append(messages, fmt.Sprintf("resource %s (%s): %s", resID, re.ResourceType, re.Message))
			}
		}
	}

	if len(messages) > 0 {
		return fmt.Errorf("rule/constraint validation failed:\n%s", strings.Join(messages, "\n"))
	}

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

	hasVars := false
	hasOutputs := false

	for _, f := range out.Files {
		target := filepath.Join(dir, f.Path)
		if err := os.WriteFile(target, []byte(f.Content), 0o644); err != nil {
			return fmt.Errorf("write file %s: %w", target, err)
		}
		switch f.Path {
		case "variables.tf":
			hasVars = true
		case "outputs.tf":
			hasOutputs = true
		}
	}

	// Optionally create placeholder files for variables.tf and outputs.tf if they
	// are not yet generated by the engine. This keeps the structure familiar.
	if !hasVars {
		varsPath := filepath.Join(dir, "variables.tf")
		if err := os.WriteFile(varsPath, []byte("# variables.tf - no input variables defined for this scenario yet\n"), 0o644); err != nil {
			return fmt.Errorf("write variables.tf: %w", err)
		}
	}
	if !hasOutputs {
		outPath := filepath.Join(dir, "outputs.tf")
		if err := os.WriteFile(outPath, []byte("# outputs.tf - no outputs defined for this scenario yet\n"), 0o644); err != nil {
			return fmt.Errorf("write outputs.tf: %w", err)
		}
	}

	return nil
}

func formatValidationErrors(result *validator.ValidationResult) string {
	if result == nil || len(result.Errors) == 0 {
		return ""
	}
	var b strings.Builder
	for _, e := range result.Errors {
		if e == nil {
			continue
		}
		b.WriteString("- ")
		b.WriteString(e.Code)
		if e.NodeID != "" {
			b.WriteString(" (node ")
			b.WriteString(e.NodeID)
			b.WriteString(")")
		}
		if e.Message != "" {
			b.WriteString(": ")
			b.WriteString(e.Message)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// resolveDiagramJSONPath resolves the json-request-diagram-valid.json path
// relative to the backend module root, regardless of the current working dir.
func resolveDiagramJSONPath(filename string) (string, error) {
	// Use runtime.Caller to get this file's directory, then walk up to the backend root.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine caller for resolving diagram JSON path")
	}

	// thisFile = .../backend/pkg/usecases/scenario5_terraform_codegen/terraform_codegen.go
	// backend root = thisFile/../../../..
	dir := filepath.Dir(thisFile)
	root := filepath.Clean(filepath.Join(dir, "..", "..", ".."))
	jsonPath := filepath.Join(root, filename)

	return jsonPath, nil
}
