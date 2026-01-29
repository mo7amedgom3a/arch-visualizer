package terraform

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
)

// Compiler orchestrates Terraform code generation for a validated architecture:
// DAG Builder -> Topological Sort -> Terraform Engine.
type Compiler struct {
	engine iac.Engine
}

func NewCompiler(terraformEngine iac.Engine) *Compiler {
	return &Compiler{engine: terraformEngine}
}

func (c *Compiler) Compile(ctx context.Context, arch *architecture.Architecture) (*iac.Output, error) {
	if c.engine == nil {
		return nil, fmt.Errorf("terraform engine is nil")
	}
	if arch == nil {
		return nil, fmt.Errorf("architecture is nil")
	}

	graph := architecture.NewGraph(arch)
	sorted, err := graph.GetSortedResources()
	if err != nil {
		return nil, err
	}

	return c.engine.Generate(ctx, arch, sorted)
}

