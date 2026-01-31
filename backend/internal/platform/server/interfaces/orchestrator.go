package interfaces

import (
	"context"

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
}

// ProcessDiagramResult contains the result of diagram processing
type ProcessDiagramResult struct {
	ProjectID uuid.UUID
	Success   bool
	Message   string
}

// GenerateCodeRequest contains data needed to generate code
type GenerateCodeRequest struct {
	ProjectID     uuid.UUID
	Engine         string
	CloudProvider string
}
