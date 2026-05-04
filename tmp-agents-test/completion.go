package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// ============================================================
// COMPLETION AGENT
// Infers ALL implicit AWS resources that are required but not drawn:
//   - Internet Gateway (when public subnet exists)
//   - Route Tables + Associations
//   - NAT Gateway + EIP (when private subnet needs outbound)
//   - Security Groups (default deny-all with required rules)
//   - IAM Instance Profile for EC2
//   - Target Groups for ALB
//   - Launch Template for ASG
// ============================================================

type CompletionAgent struct {
	llm *LLMClient
}

func NewCompletionAgent() *CompletionAgent {
	return &CompletionAgent{llm: NewLLMClient()}
}

func (a *CompletionAgent) Name() string { return "CompletionAgent" }

const completionSystemPrompt = `You are an expert AWS infrastructure engineer.
Given an architecture IR JSON (which represents what an architect drew on a canvas),
your job is to identify ALL implicit AWS resources that are REQUIRED but not explicitly present.

Think like a Terraform engineer: what resources must exist for this architecture to actually work?

Examples of implicit resources:
- A public subnet always needs: aws_internet_gateway + aws_route_table (with 0.0.0.0/0 → IGW) + aws_route_table_association
- An ALB always needs: aws_security_group (allowing HTTP/HTTPS) + aws_target_group + aws_lb_listener
- An EC2 in a private subnet needing internet: NAT Gateway + EIP in public subnet + private route table
- An ASG always needs: aws_launch_template
- RDS always needs: aws_db_subnet_group + aws_security_group

You MUST respond with ONLY a valid JSON object with two fields:
{
  "inferred_resources": [ array of Component objects to ADD ],
  "completed_architecture": { the full Architecture object with missing values filled in }
}

For completed_architecture, fix:
- Generic names like "lijhiol" → use the parent context to generate a real name
- Missing AMI IDs → use "ami-0c02fb55956c7d316" (Amazon Linux 2023 us-east-1)
- Missing security group references → reference the security groups you're creating
- Ensure all component IDs are consistent`

func (a *CompletionAgent) Run(ctx context.Context, enriched *EnrichedArchitecture) (*EnrichedArchitecture, error) {
	log.Println("[CompletionAgent] Inferring implicit resources...")

	archJSON := ArchToJSON(enriched.Original)
	issuesJSON, _ := json.MarshalIndent(enriched.ValidationIssues, "", "  ")

	userPrompt := fmt.Sprintf(`Here is the architecture IR:
%s

These validation issues were already found:
%s

Please:
1. Infer all implicit AWS resources needed
2. Return the completed architecture with all missing values filled in
3. Fix all the validation issues found above`, archJSON, string(issuesJSON))

	raw, err := a.llm.Invoke(ctx, completionSystemPrompt, userPrompt)
	if err != nil {
		return enriched, fmt.Errorf("CompletionAgent LLM call failed: %w", err)
	}

	jsonBytes, err := ParseJSONBlock(raw)
	if err != nil {
		return enriched, fmt.Errorf("CompletionAgent parse failed: %w", err)
	}

	var result struct {
		InferredResources     []Component   `json:"inferred_resources"`
		CompletedArchitecture *Architecture `json:"completed_architecture"`
	}

	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return enriched, fmt.Errorf("CompletionAgent unmarshal failed: %w", err)
	}

	enriched.InferredResources = result.InferredResources
	if result.CompletedArchitecture != nil {
		enriched.Completed = result.CompletedArchitecture
	} else {
		enriched.Completed = enriched.Original
	}

	log.Printf("[CompletionAgent] Added %d implicit resources\n", len(result.InferredResources))
	return enriched, nil
}
