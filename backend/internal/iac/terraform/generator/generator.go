package generator

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/writer"
)

type Engine struct {
	mappers *tfmapper.MapperRegistry
}

func NewEngine(mappers *tfmapper.MapperRegistry) *Engine {
	return &Engine{mappers: mappers}
}

func (e *Engine) Name() string {
	return "terraform"
}

func (e *Engine) Generate(ctx context.Context, arch *architecture.Architecture, sortedResources []*resource.Resource) (*iac.Output, error) {
	_ = ctx

	if arch == nil {
		return nil, fmt.Errorf("architecture is nil")
	}
	if e.mappers == nil {
		return nil, fmt.Errorf("terraform mapper registry is nil")
	}

	provider := string(arch.Provider)
	mapper, ok := e.mappers.Get(provider)
	if !ok {
		return nil, fmt.Errorf("no terraform mapper registered for provider %q", provider)
	}

	blocks := make([]tfmapper.TerraformBlock, 0, len(sortedResources)+1)
	// if not provider block, add it explicitly ex) provider "aws" {
	//   region = "us-east-1"
	// }
	if pb, ok := providerBlock(provider, arch.Region); ok {
		blocks = append(blocks, pb)
	}

	for _, res := range sortedResources {
		if res == nil {
			continue
		}
		if !mapper.SupportsResource(res.Type.Name) {
			return nil, fmt.Errorf("terraform mapper for %q does not support resource type %q (resource id %q)", provider, res.Type.Name, res.ID)
		}
		bs, err := mapper.MapResource(res)
		if err != nil {
			return nil, fmt.Errorf("map resource %q (%s): %w", res.ID, res.Type.Name, err)
		}
		blocks = append(blocks, bs...)
	}

	mainTF, err := writer.RenderMainTF(blocks)
	if err != nil {
		return nil, err
	}

	out := &iac.Output{
		Files: []iac.GeneratedFile{
			{Path: "main.tf", Content: mainTF, Type: "hcl"},
		},
	}
	return out, nil
}

func providerBlock(provider, region string) (tfmapper.TerraformBlock, bool) {
	if provider == "" {
		return tfmapper.TerraformBlock{}, false
	}
	attrs := map[string]tfmapper.TerraformValue{}
	if region != "" {
		r := region
		attrs["region"] = tfmapper.TerraformValue{String: &r}
	}

	return tfmapper.TerraformBlock{
		Kind:       "provider",
		Labels:     []string{provider},
		Attributes: attrs,
	}, true
}

var _ iac.Engine = (*Engine)(nil)

