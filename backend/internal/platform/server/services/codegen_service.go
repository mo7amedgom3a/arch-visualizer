package services

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
	tfgen "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/generator"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/terraform"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// CodegenServiceImpl implements CodegenService interface
type CodegenServiceImpl struct {
	engines map[string]iac.Engine
}

// NewCodegenService creates a new codegen service with default engines
func NewCodegenService() serverinterfaces.CodegenService {
	engines := make(map[string]iac.Engine)

	// Register Terraform engine
	terraformMapperRegistry := tfmapper.NewRegistry()
	// Register AWS Terraform mapper
	if err := terraformMapperRegistry.Register(terraform.New()); err != nil {
		// Log error but continue - engine will fail when used if mapper registration fails
		fmt.Printf("Warning: failed to register AWS Terraform mapper: %v\n", err)
	}
	terraformEngine := tfgen.NewEngine(terraformMapperRegistry)
	engines["terraform"] = terraformEngine

	// TODO: Register Pulumi engine when available
	// pulumiEngine := pulumi.NewEngine()
	// engines["pulumi"] = pulumiEngine

	return &CodegenServiceImpl{
		engines: engines,
	}
}

// NewCodegenServiceWithEngines creates a new codegen service with custom engines
func NewCodegenServiceWithEngines(engines map[string]iac.Engine) serverinterfaces.CodegenService {
	return &CodegenServiceImpl{
		engines: engines,
	}
}

// Generate generates IaC code for an architecture using the specified engine
func (s *CodegenServiceImpl) Generate(ctx context.Context, arch *architecture.Architecture, engine string) (*iac.Output, error) {
	if arch == nil {
		return nil, fmt.Errorf("architecture is nil")
	}

	if engine == "" {
		engine = "terraform" // Default to terraform
	}

	iacEngine, ok := s.engines[engine]
	if !ok {
		return nil, fmt.Errorf("unsupported engine: %s. Supported engines: %v", engine, s.SupportedEngines())
	}

	// Get sorted resources for proper dependency ordering
	graph := architecture.NewGraph(arch)
	sorted, err := graph.GetSortedResources()
	if err != nil {
		return nil, fmt.Errorf("failed to sort resources: %w", err)
	}

	// Generate code using the engine
	output, err := iacEngine.Generate(ctx, arch, sorted)
	if err != nil {
		return nil, fmt.Errorf("failed to generate %s code: %w", engine, err)
	}

	return output, nil
}

// SupportedEngines returns a list of supported IaC engines
func (s *CodegenServiceImpl) SupportedEngines() []string {
	engines := make([]string, 0, len(s.engines))
	for name := range s.engines {
		engines = append(engines, name)
	}
	return engines
}
