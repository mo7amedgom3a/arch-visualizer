package interfaces

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
)

// OptimizationService provides cost optimization suggestions
type OptimizationService interface {
	// OptimizeArchitecture analyzes the architecture and returns cost optimization suggestions
	OptimizeArchitecture(ctx context.Context, arch *architecture.Architecture) (*OptimizationWithSavings, error)
}

// OptimizationWithSavings contains suggestions and total potential savings
type OptimizationWithSavings struct {
	Suggestions           []OptimizationSuggestion `json:"suggestions"`
	TotalPotentialSavings float64                  `json:"total_potential_savings"`
	Currency              string                   `json:"currency"`
}

// OptimizationSuggestion represents a single cost optimization recommendation
type OptimizationSuggestion struct {
	ID               string  `json:"id"`
	Severity         string  `json:"severity"` // "high", "medium", "low"
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	EstimatedSavings float64 `json:"estimated_savings"`
	ResourceType     string  `json:"resource_type"`
	ResourceID       string  `json:"resource_id,omitempty"`
}
