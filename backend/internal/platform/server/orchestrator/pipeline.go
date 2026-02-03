package orchestrator

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// PipelineOrchestratorImpl implements PipelineOrchestrator interface
type PipelineOrchestratorImpl struct {
	diagramService      serverinterfaces.DiagramService
	architectureService serverinterfaces.ArchitectureService
	codegenService      serverinterfaces.CodegenService
	projectService      serverinterfaces.ProjectService
}

// NewPipelineOrchestrator creates a new pipeline orchestrator
func NewPipelineOrchestrator(
	diagramService serverinterfaces.DiagramService,
	architectureService serverinterfaces.ArchitectureService,
	codegenService serverinterfaces.CodegenService,
	projectService serverinterfaces.ProjectService,
) serverinterfaces.PipelineOrchestrator {
	return &PipelineOrchestratorImpl{
		diagramService:      diagramService,
		architectureService: architectureService,
		codegenService:      codegenService,
		projectService:      projectService,
	}
}

// ProcessDiagram processes a diagram JSON and persists it as a project
func (o *PipelineOrchestratorImpl) ProcessDiagram(ctx context.Context, req *serverinterfaces.ProcessDiagramRequest) (*serverinterfaces.ProcessDiagramResult, error) {
	if req == nil {
		return nil, fmt.Errorf("process diagram request is nil")
	}

	// Step 1: Parse diagram JSON
	diagramGraph, err := o.diagramService.Parse(ctx, req.JSONData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse diagram: %w", err)
	}

	// Step 2: Validate diagram
	// Note: In production, you'd want to build ValidResourceTypes from the database
	// For now, we'll use nil which means the validator will use default validation
	validationResult, err := o.diagramService.Validate(ctx, diagramGraph, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to validate diagram: %w", err)
	}

	if !validationResult.Valid {
		return nil, fmt.Errorf("diagram validation failed: %v", validationResult.Errors)
	}

	// Step 3: Map to domain architecture
	provider := resource.CloudProvider(req.CloudProvider)
	if provider == "" {
		provider = resource.AWS // Default to AWS
	}

	arch, err := o.architectureService.MapFromDiagram(ctx, diagramGraph, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to map diagram to architecture: %w", err)
	}

	// Step 4: Validate architecture rules
	ruleValidationResult, err := o.architectureService.ValidateRules(ctx, arch, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to validate architecture rules: %w", err)
	}

	if !ruleValidationResult.Valid {
		return nil, fmt.Errorf("architecture rule validation failed: %v", ruleValidationResult.Errors)
	}

	// Step 5: Create project
	createProjectReq := &serverinterfaces.CreateProjectRequest{
		UserID:        req.UserID,
		Name:          req.ProjectName,
		IACTargetID:   req.IACToolID,
		CloudProvider: string(provider),
		Region:        req.Region,
	}

	if createProjectReq.Region == "" {
		createProjectReq.Region = "us-east-1" // Default
	}

	project, err := o.projectService.Create(ctx, createProjectReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Step 6: Persist architecture with pricing (if duration specified)
	var pricingEstimate *serverinterfaces.ArchitectureCostEstimate
	if req.PricingDuration > 0 {
		// Use PersistArchitectureWithPricing for pricing calculation
		result, err := o.projectService.PersistArchitectureWithPricing(ctx, project.ID, arch, diagramGraph, req.PricingDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to persist architecture with pricing: %w", err)
		}
		pricingEstimate = result.PricingEstimate
	} else {
		// Use regular PersistArchitecture without pricing
		if err := o.projectService.PersistArchitecture(ctx, project.ID, arch, diagramGraph); err != nil {
			return nil, fmt.Errorf("failed to persist architecture: %w", err)
		}
	}

	message := fmt.Sprintf("Diagram processed successfully. Project created with ID: %s", project.ID.String())
	if pricingEstimate != nil {
		message = fmt.Sprintf("%s. Estimated monthly cost: $%.2f %s", message, pricingEstimate.TotalCost, pricingEstimate.Currency)
	}

	return &serverinterfaces.ProcessDiagramResult{
		ProjectID:       project.ID,
		Success:         true,
		Message:         message,
		PricingEstimate: pricingEstimate,
	}, nil
}

// GenerateCode generates IaC code for an existing project
func (o *PipelineOrchestratorImpl) GenerateCode(ctx context.Context, req *serverinterfaces.GenerateCodeRequest) (*iac.Output, error) {
	if req == nil {
		return nil, fmt.Errorf("generate code request is nil")
	}

	// Step 1: Get project (validate it exists and get provider info)
	project, err := o.projectService.GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Step 2: Load architecture from project
	arch, err := o.projectService.LoadArchitecture(ctx, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load architecture: %w", err)
	}

	// Step 3: Validate architecture rules (optional but recommended)
	provider := resource.CloudProvider(req.CloudProvider)
	if provider == "" {
		provider = resource.CloudProvider(project.CloudProvider)
		if provider == "" {
			provider = resource.AWS // Default
		}
	}

	ruleValidationResult, err := o.architectureService.ValidateRules(ctx, arch, provider)
	if err != nil {
		// Log warning but continue - rules validation failure shouldn't block code generation
		// In production, you might want to make this configurable
	} else if !ruleValidationResult.Valid {
		// Log warnings but continue
		// In production, you might want to return errors or warnings
	}

	// Step 4: Generate code
	engine := req.Engine
	if engine == "" {
		engine = "terraform" // Default
	}

	output, err := o.codegenService.Generate(ctx, arch, engine)
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}

	return output, nil
}
