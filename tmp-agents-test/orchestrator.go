// Package orchestrator implements the master agent pipeline.
// It runs specialist agents in the correct order, passing enriched
// state between them — modelled after LangChain's SequentialChain.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// PipelineResult is the final output returned to the caller.
type PipelineResult struct {
	ArchitectureID    string            `json:"architecture_id"`
	RunDurationMs     int64             `json:"run_duration_ms"`
	ValidationIssues  []ValidationIssue `json:"validation_issues"`
	SecurityFindings  []SecurityFinding `json:"security_findings"`
	InferredResources []Component       `json:"inferred_resources"`
	GeneratedCode     string            `json:"generated_code"`
	CostEstimate      *CostReport       `json:"cost_estimate"`
	Recommendations   []Recommendation  `json:"recommendations"`
	Errors            []string          `json:"errors,omitempty"`
}

// Orchestrator manages the agent pipeline execution.
type Orchestrator struct {
	agents []Agent
}

// NewOrchestrator creates the default agent pipeline in execution order:
// Validate → Complete → Security → Generate IaC → Cost → Recommend
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		agents: []Agent{
			NewValidatorAgent(),       // 1. Find all issues
			NewCompletionAgent(),      // 2. Fill gaps + add implicit resources
			NewSecurityAgent(),        // 3. Audit security (uses completed arch)
			NewTerraformAgent(),       // 4. Generate IaC (bakes in security fixes)
			NewCostAgent(),            // 5. Estimate costs (uses final resource set)
			NewRecommendationsAgent(), // 6. Synthesize recommendations
		},
	}
}

// Run executes the full agent pipeline on the given architecture IR.
func (o *Orchestrator) Run(ctx context.Context, arch *Architecture) (*PipelineResult, error) {
	start := time.Now()
	log.Printf("[Orchestrator] Starting pipeline for architecture: %s\n", arch.ArchitectureID)

	// Initialize the enriched state that flows through all agents
	enriched := &EnrichedArchitecture{
		Original: arch,
	}

	var pipelineErrors []string

	// Run each agent sequentially, passing enriched state forward
	// (For parallelism: Security + Cost can run in goroutines after Completion)
	for i, agent := range o.agents {
		log.Printf("[Orchestrator] Running agent: %s\n", agent.Name())

		var err error
		enriched, err = agent.Run(ctx, enriched)
		if err != nil {
			// Log and continue — don't let one agent failure kill the whole pipeline
			errMsg := fmt.Sprintf("[%s] %v", agent.Name(), err)
			log.Printf("[Orchestrator] WARNING: %s\n", errMsg)
			pipelineErrors = append(pipelineErrors, errMsg)
		}

		// Add a delay between agents to avoid rate limits
		if i < len(o.agents)-1 {
			log.Printf("[Orchestrator] Waiting 15 seconds before next agent to prevent rate limiting...\n")
			time.Sleep(15 * time.Second)
		}
	}

	duration := time.Since(start).Milliseconds()
	log.Printf("[Orchestrator] Pipeline complete in %dms\n", duration)

	return &PipelineResult{
		ArchitectureID:    arch.ArchitectureID,
		RunDurationMs:     duration,
		ValidationIssues:  enriched.ValidationIssues,
		SecurityFindings:  enriched.SecurityFindings,
		InferredResources: enriched.InferredResources,
		GeneratedCode:     enriched.GeneratedCode,
		CostEstimate:      enriched.CostEstimate,
		Recommendations:   enriched.Recommendations,
		Errors:            pipelineErrors,
	}, nil
}

// Summary prints a human-readable pipeline summary to stdout.
func (r *PipelineResult) Summary() string {
	errors := len(r.ValidationIssues)
	criticalSec := 0
	for _, f := range r.SecurityFindings {
		if f.Severity == "critical" {
			criticalSec++
		}
	}

	monthly := 0.0
	if r.CostEstimate != nil {
		monthly = r.CostEstimate.TotalMonthlyUSD
	}

	return fmt.Sprintf(`
╔══════════════════════════════════════════════════════════╗
║           CLOUD CANVAS — AGENT PIPELINE RESULT           ║
╠══════════════════════════════════════════════════════════╣
║  Architecture ID  : %-36s ║
║  Run Duration     : %dms%-*s ║
║  Validation Issues: %-36d ║
║  Security Findings: %-36d ║
║  Critical Security: %-36d ║
║  Inferred Resources: %-35d ║
║  Est. Monthly Cost: $%-35.2f ║
║  Recommendations  : %-36d ║
║  IaC Code Length  : %-36d ║
╚══════════════════════════════════════════════════════════╝`,
		r.ArchitectureID,
		r.RunDurationMs, 36-len(fmt.Sprintf("%dms", r.RunDurationMs)), "",
		errors,
		len(r.SecurityFindings),
		criticalSec,
		len(r.InferredResources),
		monthly,
		len(r.Recommendations),
		len(r.GeneratedCode),
	)
}

// ToJSON serializes the full result to JSON.
func (r *PipelineResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
