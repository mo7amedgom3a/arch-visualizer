package main

import (
	"context"
	"fmt"
	"log"
)

// ============================================================
// TERRAFORM GENERATOR AGENT
// Produces production-grade Terraform HCL from the completed IR.
// Security findings are baked directly into the generated code.
// Output is a modular, DRY, well-commented Terraform project.
// ============================================================

type TerraformAgent struct {
	llm *LLMClient
}

func NewTerraformAgent() *TerraformAgent {
	return &TerraformAgent{llm: NewLLMClient()}
}

func (a *TerraformAgent) Name() string { return "TerraformAgent" }

const terraformSystemPrompt = `You are a senior Terraform engineer with 10 years of AWS experience.
Generate production-grade Terraform HCL from the given architecture IR.

RULES:
1. Generate a COMPLETE, working Terraform project — not snippets
2. Use terraform best practices: variables, locals, outputs, for_each where appropriate
3. Bake in ALL security fixes from the security findings — do not leave vulnerabilities in the code
4. Add standard tags to every resource: Name, Environment, ManagedBy=Terraform, Project
5. Include implicit resources identified by the completion agent
6. Structure output as a single main.tf with: terraform{}, provider{}, locals{}, then resources
7. Add inline comments explaining non-obvious decisions
8. Use data sources where appropriate (e.g., data "aws_availability_zones")

OUTPUT FORMAT:
Return ONLY the raw Terraform HCL code. No markdown fences, no explanation.
Start directly with: terraform {`

func (a *TerraformAgent) Run(ctx context.Context, enriched *EnrichedArchitecture) (*EnrichedArchitecture, error) {
	log.Println("[TerraformAgent] Generating Terraform code...")

	arch := enriched.Completed
	if arch == nil {
		arch = enriched.Original
	}

	archJSON := ArchToJSON(arch)

	// Build security context string
	securityContext := ""
	for _, f := range enriched.SecurityFindings {
		if f.Severity == "critical" || f.Severity == "high" {
			securityContext += fmt.Sprintf("- [%s] %s: %s\n", f.Severity, f.Rule, f.Remediation)
		}
	}

	// Build inferred resources context
	inferredContext := ""
	for _, r := range enriched.InferredResources {
		inferredContext += fmt.Sprintf("- %s (%s)\n", r.ID, r.Type)
	}

	userPrompt := fmt.Sprintf(`Generate Terraform HCL for this architecture:

ARCHITECTURE IR:
%s

SECURITY FIXES TO BAKE IN (do not skip these):
%s

IMPLICIT RESOURCES TO INCLUDE:
%s

Generate a complete main.tf that provisions ALL of the above.
Remember to include: VPC, subnets, IGW, route tables, security groups,
ALB with target group + listener, ASG with launch template, and all
the implicit resources listed above.`, archJSON, securityContext, inferredContext)

	code, err := a.llm.Invoke(ctx, terraformSystemPrompt, userPrompt)
	if err != nil {
		return enriched, fmt.Errorf("TerraformAgent LLM call failed: %w", err)
	}

	// Strip any accidental markdown fences
	code = stripMarkdownFences(code)

	enriched.GeneratedCode = code
	log.Printf("[TerraformAgent] Generated %d characters of Terraform code\n", len(code))

	return enriched, nil
}

func stripMarkdownFences(s string) string {
	// Remove ```hcl, ```terraform, ``` fences
	result := []byte(s)
	fences := []string{"```hcl\n", "```terraform\n", "```\n", "```hcl", "```terraform", "```"}
	for _, fence := range fences {
		for i := 0; i < len(result)-len(fence); i++ {
			if string(result[i:i+len(fence)]) == fence {
				result = append(result[:i], result[i+len(fence):]...)
			}
		}
	}
	return string(result)
}
