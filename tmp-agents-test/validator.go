package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// ============================================================
// VALIDATOR AGENT
// Checks every component for:
//   - Missing required fields
//   - Invalid CIDR ranges
//   - Broken parent/child relationships
//   - Conflicting configurations
// ============================================================

type ValidatorAgent struct {
	llm *LLMClient
}

func NewValidatorAgent() *ValidatorAgent {
	return &ValidatorAgent{llm: NewLLMClient()}
}

func (a *ValidatorAgent) Name() string { return "ValidatorAgent" }

const validatorSystemPrompt = `You are a senior AWS cloud architect specializing in infrastructure validation.
Your job is to analyze an architecture IR JSON and find ALL problems:
- Missing required fields (e.g., VPC missing CIDR, subnet missing availability zone)
- Invalid values (e.g., invalid CIDR notation, unsupported instance types)
- Broken relationships (e.g., subnet references a VPC that does not exist)
- Dependency issues (e.g., ALB needs at least 2 subnets in different AZs)
- Configuration conflicts (e.g., private subnet marked as public)

You MUST respond with ONLY a valid JSON array of ValidationIssue objects. No explanation, no markdown.
Each object has these fields:
{
  "component_id": "string",
  "severity": "error|warning|info",
  "field": "string",
  "message": "string",
  "suggestion": "string"
}`

func (a *ValidatorAgent) Run(ctx context.Context, enriched *EnrichedArchitecture) (*EnrichedArchitecture, error) {
	log.Println("[ValidatorAgent] Starting validation...")

	archJSON := ArchToJSON(enriched.Original)
	userPrompt := fmt.Sprintf(`Validate this AWS architecture IR and return all issues as a JSON array:

%s

Focus especially on:
1. Any component with empty or default names like "lijhiol" — these need real names
2. ALB listener targets — do they reference existing component IDs?
3. Subnet CIDR blocks — are they valid subnets of the parent VPC CIDR?
4. Auto Scaling Group — does it have a valid launch template or AMI?
5. Missing security groups on EC2 and ALB`, archJSON)

	raw, err := a.llm.Invoke(ctx, validatorSystemPrompt, userPrompt)
	if err != nil {
		return enriched, fmt.Errorf("ValidatorAgent LLM call failed: %w", err)
	}

	jsonBytes, err := ParseJSONBlock(raw)
	if err != nil {
		return enriched, fmt.Errorf("ValidatorAgent parse failed: %w", err)
	}

	var issues []ValidationIssue
	if err := json.Unmarshal(jsonBytes, &issues); err != nil {
		return enriched, fmt.Errorf("ValidatorAgent unmarshal failed: %w", err)
	}

	enriched.ValidationIssues = issues
	log.Printf("[ValidatorAgent] Found %d issues (%d errors)\n",
		len(issues), countBySeverity(issues, "error"))

	return enriched, nil
}

func countBySeverity(issues []ValidationIssue, severity string) int {
	count := 0
	for _, i := range issues {
		if i.Severity == severity {
			count++
		}
	}
	return count
}
