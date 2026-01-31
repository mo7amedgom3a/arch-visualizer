package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram/graph"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/engine"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// ArchitectureServiceImpl implements ArchitectureService interface
type ArchitectureServiceImpl struct {
	ruleService serverinterfaces.RuleService
}

// NewArchitectureService creates a new architecture service
func NewArchitectureService(ruleService serverinterfaces.RuleService) serverinterfaces.ArchitectureService {
	return &ArchitectureServiceImpl{
		ruleService: ruleService,
	}
}

// MapFromDiagram converts a diagram graph to a domain architecture
func (s *ArchitectureServiceImpl) MapFromDiagram(ctx context.Context, graph *graph.DiagramGraph, provider resource.CloudProvider) (*architecture.Architecture, error) {
	if graph == nil {
		return nil, fmt.Errorf("diagram graph is nil")
	}

	arch, err := architecture.MapDiagramToArchitecture(graph, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to map diagram to architecture: %w", err)
	}

	// Basic domain validation
	if err := arch.Validate(); err != nil {
		return nil, fmt.Errorf("architecture validation failed: %w", err)
	}

	return arch, nil
}

// ValidateRules validates an architecture against domain rules and constraints
func (s *ArchitectureServiceImpl) ValidateRules(ctx context.Context, arch *architecture.Architecture, provider resource.CloudProvider) (*serverinterfaces.RuleValidationResult, error) {
	if arch == nil {
		return nil, fmt.Errorf("architecture is nil")
	}

	if s.ruleService == nil {
		// No rule service provided, return valid result
		return &serverinterfaces.RuleValidationResult{
			Valid:   true,
			Results: make(map[string]*serverinterfaces.ResourceValidationResult),
		}, nil
	}

	// Adapt domain architecture to rules engine architecture view
	engineArch := &engine.Architecture{
		Resources: arch.Resources,
	}

	// Validate architecture using rule service
	resultsMap, err := s.ruleService.ValidateArchitecture(ctx, engineArch)
	if err != nil {
		return nil, fmt.Errorf("failed to validate architecture rules: %w", err)
	}

	// Convert results to our format
	validationResult := &serverinterfaces.RuleValidationResult{
		Valid:   true,
		Results: make(map[string]*serverinterfaces.ResourceValidationResult),
		Errors:  make([]string, 0),
	}

	for resID, resResult := range resultsMap {
		// Type assert to get the actual result type
		// The rule service returns map[string]interface{}, we need to convert it
		var evalResult *engine.EvaluationResult
		if resultMap, ok := resResult.(map[string]interface{}); ok {
			// Try to extract the actual EvaluationResult
			// This is a bit of a workaround - ideally the interface would be more specific
			evalResult = convertToEvaluationResult(resultMap)
		} else if er, ok := resResult.(*engine.EvaluationResult); ok {
			evalResult = er
		} else {
			// Fallback: create a valid result
			evalResult = &engine.EvaluationResult{
				Valid:   true,
				Results: []*rules.RuleResult{},
			}
		}

		resourceResult := &serverinterfaces.ResourceValidationResult{
			ResourceID: resID,
			Valid:      evalResult.Valid,
			Errors:     make([]serverinterfaces.ValidationError, 0),
		}

		// Find the resource to get its type
		for _, res := range arch.Resources {
			if res.ID == resID {
				resourceResult.ResourceType = res.Type.Name
				break
			}
		}

		if !evalResult.Valid {
			validationResult.Valid = false
			for _, ruleResult := range evalResult.Results {
				if !ruleResult.Passed && ruleResult.Error != nil {
					validationResult.Errors = append(validationResult.Errors, fmt.Sprintf("resource %s (%s): %s", resID, ruleResult.Error.ResourceType, ruleResult.Error.Message))
					resourceResult.Errors = append(resourceResult.Errors, serverinterfaces.ValidationError{
						ResourceID:   resID,
						ResourceType: ruleResult.Error.ResourceType,
						Message:      ruleResult.Error.Message,
						Code:         string(ruleResult.Error.RuleType),
					})
				}
			}
		}

		validationResult.Results[resID] = resourceResult
	}

	return validationResult, nil
}

// convertToEvaluationResult converts a map to EvaluationResult
// This is a workaround for the interface returning map[string]interface{}
func convertToEvaluationResult(m map[string]interface{}) *engine.EvaluationResult {
	result := &engine.EvaluationResult{
		Valid:   true,
		Results: []*rules.RuleResult{},
	}

	if valid, ok := m["valid"].(bool); ok {
		result.Valid = valid
	}

	if results, ok := m["results"].([]interface{}); ok {
		for _, r := range results {
			if ruleResultMap, ok := r.(map[string]interface{}); ok {
				ruleResult := &rules.RuleResult{
					Passed: true,
				}
				if passed, ok := ruleResultMap["passed"].(bool); ok {
					ruleResult.Passed = passed
				}
				if errorMap, ok := ruleResultMap["error"].(map[string]interface{}); ok {
					ruleError := &rules.RuleError{}
					if msg, ok := errorMap["message"].(string); ok {
						ruleError.Message = msg
					}
					if rt, ok := errorMap["resourceType"].(string); ok {
						ruleError.ResourceType = rt
					}
					ruleResult.Error = ruleError
				}
				result.Results = append(result.Results, ruleResult)
			}
		}
	}

	return result
}

// GetSortedResources returns resources sorted by dependencies (topological sort)
func (s *ArchitectureServiceImpl) GetSortedResources(ctx context.Context, arch *architecture.Architecture) ([]*resource.Resource, error) {
	if arch == nil {
		return nil, fmt.Errorf("architecture is nil")
	}

	graph := architecture.NewGraph(arch)
	sorted, err := graph.GetSortedResources()
	if err != nil {
		return nil, fmt.Errorf("topological sort failed: %w", err)
	}

	return sorted, nil
}

// Helper function to format validation errors
func formatValidationErrors(result *serverinterfaces.RuleValidationResult) string {
	if result == nil || len(result.Errors) == 0 {
		return ""
	}
	return strings.Join(result.Errors, "\n")
}
