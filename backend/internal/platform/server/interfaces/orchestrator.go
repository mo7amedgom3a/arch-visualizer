package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
)

// PipelineOrchestrator orchestrates the complete diagram-to-code pipeline
type PipelineOrchestrator interface {
	// ProcessDiagram processes a diagram JSON and persists it as a project
	ProcessDiagram(ctx context.Context, req *ProcessDiagramRequest) (*ProcessDiagramResult, error)

	// GenerateCode generates IaC code for an existing project
	GenerateCode(ctx context.Context, req *GenerateCodeRequest) (*iac.Output, error)
}

// ProcessDiagramRequest contains data needed to process a diagram
type ProcessDiagramRequest struct {
	JSONData      []byte
	UserID        uuid.UUID
	ProjectName   string
	IACToolID     uint
	CloudProvider string
	Region        string
	// PricingDuration is the duration for pricing calculation (e.g., 720h for monthly)
	// If zero, pricing calculation is skipped
	PricingDuration time.Duration
}

// ProcessDiagramResult contains the result of diagram processing
type ProcessDiagramResult struct {
	ProjectID uuid.UUID `json:"project_id"`
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	// PricingEstimate contains the architecture cost estimate (if pricing was calculated)
	PricingEstimate *ArchitectureCostEstimate `json:"pricing_estimate,omitempty"`
}

// GenerateCodeRequest contains data needed to generate code
type GenerateCodeRequest struct {
	ProjectID     uuid.UUID
	Engine        string
	CloudProvider string
}
