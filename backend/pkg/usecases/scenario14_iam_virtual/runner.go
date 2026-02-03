package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/iam" // Register IAM mappers
	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/terraform"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/generator"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("Starting IAM Virtual User Simulation...")

	// 1. Setup Architecture manually with an IAM User
	arch := architecture.NewArchitecture()
	arch.Provider = resource.AWS

	userRes := &resource.Resource{
		ID:   "virtual-user-1",
		Name: "virtual-user-1",
		Type: resource.ResourceType{
			ID:       "IAMUser",
			Name:     "IAMUser",
			Category: "IAM",
		},
		Provider: "aws",
		Metadata: map[string]interface{}{
			"name":        "my-virtual-user",
			"path":        "/system/",
			"is_virtual":  true,
			"description": "A virtual user for testing",
		},
	}
	arch.Resources = append(arch.Resources, userRes)

	fmt.Printf("Added resource: %s (%s)\n", userRes.ID, userRes.Type.Name)

	// 2. Configure Generator
	mapper := awsmapper.New()
	registry := tfmapper.NewRegistry()
	if err := registry.Register(mapper); err != nil {
		return fmt.Errorf("register mapper failed: %w", err)
	}

	tfGenerator := generator.NewEngine(registry)

	// 3. Generate Terraform
	outputDir := "./terraform_output_iam_sim"
	_ = os.RemoveAll(outputDir) // Clean up

	// Note: Generator expects sorted list, we just pass the slice
	output, err := tfGenerator.Generate(context.Background(), arch, arch.Resources)
	if err != nil {
		return fmt.Errorf("generate failed: %w", err)
	}

	// Write output manually since Generate returns *iac.Output
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	for _, file := range output.Files {
		if err := os.WriteFile(outputDir+"/"+file.Path, []byte(file.Content), 0644); err != nil {
			return err
		}
	}

	fmt.Printf("Generated %d files in %s\n", len(output.Files), outputDir)

	// 4. Validate output
	mainTfPath := outputDir + "/main.tf"
	contentBytes, err := os.ReadFile(mainTfPath)
	if err != nil {
		return fmt.Errorf("read main.tf failed: %w", err)
	}
	content := string(contentBytes)

	fmt.Println("\n--- Generated main.tf Content ---")
	fmt.Println(content)
	fmt.Println("-------------------------------")

	if !strings.Contains(content, "resource \"aws_iam_user\"") {
		return fmt.Errorf("FAILED: output does not contain aws_iam_user resource")
	}
	if !strings.Contains(content, "\"my-virtual-user\"") {
		return fmt.Errorf("FAILED: output does not contain correct user name")
	}
	if !strings.Contains(content, "\"/system/\"") {
		return fmt.Errorf("FAILED: output does not contain correct path")
	}

	fmt.Println("SUCCESS! Terraform code generated correctly for IAMUser.")
	return nil
}
