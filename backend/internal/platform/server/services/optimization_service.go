package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

type OptimizationServiceImpl struct {
}

func NewOptimizationService() serverinterfaces.OptimizationService {
	return &OptimizationServiceImpl{}
}

func (s *OptimizationServiceImpl) OptimizeArchitecture(ctx context.Context, arch *architecture.Architecture) (*serverinterfaces.OptimizationWithSavings, error) {
	suggestions := []serverinterfaces.OptimizationSuggestion{}
	totalSavings := 0.0

	// 1. Check for older generation instances (e.g., t2 vs t3/t4g)
	for _, res := range arch.Resources {
		if res.Type.Name == "EC2" {
			if instanceType, ok := res.Metadata["instance_type"].(string); ok {
				if strings.HasPrefix(instanceType, "t2.") {
					// t3 is generally cheaper and better performance
					suggestion := serverinterfaces.OptimizationSuggestion{
						ID:               fmt.Sprintf("opt-upgrade-%s", res.ID),
						Severity:         "medium",
						Title:            "Upgrade to next generation instance",
						Description:      fmt.Sprintf("Resource '%s' uses %s. Consider upgrading to t3.%s for better price/performance.", res.Name, instanceType, strings.TrimPrefix(instanceType, "t2.")),
						EstimatedSavings: 0.0, // Needs pricing calculation difference, defaulting to 0 or mock value for now
						ResourceType:     "EC2",
						ResourceID:       res.ID,
					}
					// Estimate 20% savings roughly
					suggestion.EstimatedSavings = 5.0 // Mock savings
					suggestions = append(suggestions, suggestion)
					totalSavings += suggestion.EstimatedSavings
				}
			}
		}
	}

	// 2. Check for potential Unattached Elastic IPs (if we can infer from graph)
	// In the architecture graph, if an EIP node exists but has no dependencies or is not depended upon by an instance, it MIGHT be unattached.
	// But EIPs are often "depended on" by NAT Gateways or "attached" to instances via ID references in config.
	for _, res := range arch.Resources {
		if res.Type.Name == "ElasticIP" {
			// Check if it's used by anything
			isUsed := false

			// Check incoming dependencies (something implies they need this IP)
			// arch.Dependencies maps resourceID -> []dependencyIDs (outgoing)
			// We need checked who depends on us.
			for _, deps := range arch.Dependencies {
				for _, depID := range deps {
					if depID == res.ID {
						isUsed = true
						break
					}
				}
				if isUsed {
					break
				}
			}

			if !isUsed {
				// Also check if it depends on an instance (attachment)
				if len(arch.Dependencies[res.ID]) > 0 {
					isUsed = true
				}
			}

			if !isUsed {
				suggestion := serverinterfaces.OptimizationSuggestion{
					ID:               fmt.Sprintf("opt-eip-unattached-%s", res.ID),
					Severity:         "high",
					Title:            "Potential Unattached Elastic IP",
					Description:      fmt.Sprintf("Elastic IP '%s' appears to be unattached. Unattached EIPs incur hourly charges.", res.Name),
					EstimatedSavings: 3.60, // $0.005 * 720 hours
					ResourceType:     "ElasticIP",
					ResourceID:       res.ID,
				}
				suggestions = append(suggestions, suggestion)
				totalSavings += suggestion.EstimatedSavings
			}
		}
	}

	// 3. Generic Savings Plan / Reserved Instances suggestion
	// If there are EC2 instances running 24/7 (which we assume for architecture diagrams)
	hasEC2 := false
	for _, res := range arch.Resources {
		if res.Type.Name == "EC2" {
			hasEC2 = true
			break
		}
	}

	if hasEC2 {
		suggestions = append(suggestions, serverinterfaces.OptimizationSuggestion{
			ID:               "opt-ri-savings",
			Severity:         "high",
			Title:            "Use Reserved Instances or Savings Plans",
			Description:      "For steady-state workloads, purchasing Reserved Instances or Savings Plans can reduce costs by up to 72% compared to On-Demand prices.",
			EstimatedSavings: 50.0, // Mock value
			ResourceType:     "EC2",
		})
		totalSavings += 50.0
	}

	return &serverinterfaces.OptimizationWithSavings{
		Suggestions:           suggestions,
		TotalPotentialSavings: totalSavings,
		Currency:              "USD",
	}, nil
}
