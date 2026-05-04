package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// ============================================================
// RECOMMENDATIONS AGENT
// Synthesizes all previous agent findings and produces
// actionable recommendations across 4 pillars:
//   - Cost optimization
//   - Security improvements
//   - Reliability / HA
//   - Performance
// ============================================================

type RecommendationsAgent struct {
	llm *LLMClient
}

func NewRecommendationsAgent() *RecommendationsAgent {
	return &RecommendationsAgent{llm: NewLLMClient()}
}

func (a *RecommendationsAgent) Name() string { return "RecommendationsAgent" }

const recommendationsSystemPrompt = `You are a senior AWS solutions architect reviewing a cloud architecture.
Given the full context (architecture, security findings, cost estimate), produce actionable recommendations.

Focus on HIGH-VALUE recommendations only — things the architect should seriously consider.
Do NOT repeat issues already flagged by the security agent.

Categories:
- cost: ways to reduce the monthly bill
- security: additional security improvements beyond critical findings
- reliability: HA, multi-AZ, backup, DR improvements  
- performance: caching, CDN, database optimization

You MUST respond with ONLY a valid JSON array:
[
  {
    "type": "cost|security|reliability|performance",
    "priority": "high|medium|low",
    "title": "short title",
    "description": "detailed explanation",
    "saving": "optional: estimated saving e.g. ~$30/month"
  }
]`

func (a *RecommendationsAgent) Run(ctx context.Context, enriched *EnrichedArchitecture) (*EnrichedArchitecture, error) {
	log.Println("[RecommendationsAgent] Generating recommendations...")

	archJSON := ArchToJSON(enriched.Completed)
	costJSON, _ := json.MarshalIndent(enriched.CostEstimate, "", "  ")
	secJSON, _ := json.MarshalIndent(enriched.SecurityFindings, "", "  ")

	userPrompt := fmt.Sprintf(`Review this architecture and provide recommendations:

ARCHITECTURE:
%s

COST ESTIMATE:
%s

SECURITY FINDINGS ALREADY REPORTED:
%s

Provide recommendations that ADD VALUE beyond what's already flagged.
Consider: Reserved Instances savings, CloudFront for the ALB, RDS instead of EC2-hosted DB,
multi-AZ for the ASG, S3 for static assets, WAF for the ALB, backup policies, etc.`,
		archJSON, string(costJSON), string(secJSON))

	raw, err := a.llm.Invoke(ctx, recommendationsSystemPrompt, userPrompt)
	if err != nil {
		return enriched, fmt.Errorf("RecommendationsAgent LLM call failed: %w", err)
	}

	jsonBytes, err := ParseJSONBlock(raw)
	if err != nil {
		return enriched, fmt.Errorf("RecommendationsAgent parse failed: %w", err)
	}

	var recs []Recommendation
	if err := json.Unmarshal(jsonBytes, &recs); err != nil {
		return enriched, fmt.Errorf("RecommendationsAgent unmarshal failed: %w", err)
	}

	enriched.Recommendations = recs
	log.Printf("[RecommendationsAgent] Generated %d recommendations\n", len(recs))

	return enriched, nil
}
