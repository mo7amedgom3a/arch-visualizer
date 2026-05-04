package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// ============================================================
// COST ESTIMATION AGENT
// Estimates monthly AWS costs for the full architecture including:
//   - Compute (EC2, ASG)
//   - Load Balancing (ALB)
//   - NAT Gateway
//   - Data Transfer
//   - Storage (EBS)
//   - Any other inferred resources
// ============================================================

type CostAgent struct {
	llm *LLMClient
}

func NewCostAgent() *CostAgent {
	return &CostAgent{llm: NewLLMClient()}
}

func (a *CostAgent) Name() string { return "CostAgent" }

const costSystemPrompt = `You are an AWS FinOps specialist with deep knowledge of AWS pricing (us-east-1 region, 2024 prices).
Estimate the monthly cost for the given architecture.

PRICING REFERENCE (us-east-1):
- EC2 t3.micro:    $0.0104/hr  → $7.59/mo
- EC2 t3.small:    $0.0208/hr  → $15.18/mo
- EC2 t3.medium:   $0.0416/hr  → $30.37/mo
- EC2 t3.large:    $0.0832/hr  → $60.74/mo
- ALB:             $0.008/hr + $0.008/LCU-hr → ~$16/mo base
- NAT Gateway:     $0.045/hr + $0.045/GB    → ~$32/mo base
- EBS gp3 (20GB):  $0.08/GB/mo             → $1.60/mo
- EBS gp3 (100GB): $0.08/GB/mo             → $8.00/mo
- Data Transfer out: $0.09/GB (first 10TB)

ASSUMPTIONS when not specified:
- EC2 instances: 730 hrs/month (always on)
- NAT Gateway data: 100GB/month
- ALB: 10 LCU average

You MUST respond with ONLY a valid JSON object:
{
  "total_monthly_usd": number,
  "total_yearly_usd": number,
  "breakdown": [
    {
      "component_id": "string",
      "resource_type": "string",
      "name": "string",
      "monthly_usd": number,
      "basis": "human readable explanation"
    }
  ],
  "warnings": ["string array of cost warnings"]
}`

func (a *CostAgent) Run(ctx context.Context, enriched *EnrichedArchitecture) (*EnrichedArchitecture, error) {
	log.Println("[CostAgent] Estimating costs...")

	arch := enriched.Completed
	if arch == nil {
		arch = enriched.Original
	}

	archJSON := ArchToJSON(arch)
	inferredJSON, _ := json.MarshalIndent(enriched.InferredResources, "", "  ")

	userPrompt := fmt.Sprintf(`Estimate monthly AWS costs for this architecture:

ARCHITECTURE:
%s

INFERRED RESOURCES (also include these in cost):
%s

Include costs for ALL resources — drawn AND inferred.
Flag any potential hidden costs (data transfer, CloudWatch logs, etc.)`, archJSON, string(inferredJSON))

	raw, err := a.llm.Invoke(ctx, costSystemPrompt, userPrompt)
	if err != nil {
		return enriched, fmt.Errorf("CostAgent LLM call failed: %w", err)
	}

	jsonBytes, err := ParseJSONBlock(raw)
	if err != nil {
		return enriched, fmt.Errorf("CostAgent parse failed: %w", err)
	}

	var report CostReport
	if err := json.Unmarshal(jsonBytes, &report); err != nil {
		return enriched, fmt.Errorf("CostAgent unmarshal failed: %w", err)
	}

	enriched.CostEstimate = &report
	log.Printf("[CostAgent] Estimated total: $%.2f/month ($%.2f/year)\n",
		report.TotalMonthlyUSD, report.TotalYearlyUSD)

	return enriched, nil
}
