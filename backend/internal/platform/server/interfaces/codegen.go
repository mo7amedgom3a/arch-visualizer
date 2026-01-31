package interfaces

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
)

// CodegenService handles Infrastructure as Code generation
type CodegenService interface {
	// Generate generates IaC code for an architecture using the specified engine
	Generate(ctx context.Context, arch *architecture.Architecture, engine string) (*iac.Output, error)

	// SupportedEngines returns a list of supported IaC engines (e.g., "terraform", "pulumi")
	SupportedEngines() []string
}
