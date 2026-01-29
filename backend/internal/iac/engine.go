package iac

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// GeneratedFile represents a single generated IaC file.
type GeneratedFile struct {
	// Path is the file path relative to the output root (e.g. "main.tf").
	Path string
	// Content is the file contents.
	Content string
	// Type is the file type/language hint (e.g. "hcl", "python").
	Type string
}

// Output is the result of an IaC engine generation.
type Output struct {
	Files []GeneratedFile
}

// Engine is a pluggable Infrastructure-as-Code generator (Terraform, Pulumi, ...).
//
// IMPORTANT: Engines must not hardcode cloud-provider specifics; they delegate
// provider-specific mapping to provider adapters/mappers.
type Engine interface {
	Name() string

	// Generate compiles a validated architecture into IaC files.
	// The caller is responsible for ensuring validations/rules have passed.
	Generate(ctx context.Context, arch *architecture.Architecture, sortedResources []*resource.Resource) (*Output, error)
}
