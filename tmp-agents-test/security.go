package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// ============================================================
// SECURITY AGENT
// Audits the completed architecture for:
//   - Overly permissive security groups (0.0.0.0/0 on sensitive ports)
//   - Unencrypted storage (S3, RDS, EBS)
//   - Missing VPC Flow Logs
//   - Public RDS instances
//   - IAM over-provisioning
//   - Missing MFA / access controls
// ============================================================

type SecurityAgent struct {
	llm *LLMClient
}

func NewSecurityAgent() *SecurityAgent {
	return &SecurityAgent{llm: NewLLMClient()}
}

func (a *SecurityAgent) Name() string { return "SecurityAgent" }

const securitySystemPrompt = `You are an AWS security engineer and CIS Benchmark expert.
Analyze the given architecture IR for security vulnerabilities and misconfigurations.

For EACH finding, assess:
- Severity: critical (immediate risk) | high (serious) | medium (should fix) | low (best practice)
- Whether it can be auto-fixed in code (auto_fixed: true)
- A clear remediation instruction

You MUST respond with ONLY a valid JSON array of SecurityFinding objects:
[
  {
    "component_id": "string",
    "severity": "critical|high|medium|low",
    "rule": "string (e.g. CIS-AWS-4.1)",
    "description": "string",
    "remediation": "string",
    "auto_fixed": boolean
  }
]

Key rules to check:
- CIS-AWS-4.1: No security group allows unrestricted SSH (port 22) from 0.0.0.0/0
- CIS-AWS-4.2: No security group allows unrestricted RDP (port 3389) from 0.0.0.0/0
- AWS-SEC-1: ALB should have HTTPS listener, not just HTTP (port 80)
- AWS-SEC-2: EC2 instances should not have public IPs unless in public subnet with explicit need
- AWS-SEC-3: ASG launch template should not use overly permissive IAM roles
- AWS-SEC-4: VPC Flow Logs should be enabled
- AWS-SEC-5: EC2 should have IMDSv2 enforced (http_tokens = required)`

func (a *SecurityAgent) Run(ctx context.Context, enriched *EnrichedArchitecture) (*EnrichedArchitecture, error) {
	log.Println("[SecurityAgent] Running security audit...")

	// Use completed architecture if available (has more context)
	arch := enriched.Completed
	if arch == nil {
		arch = enriched.Original
	}

	archJSON := ArchToJSON(arch)
	userPrompt := fmt.Sprintf(`Audit this AWS architecture for security issues:

%s

Pay special attention to:
1. The ALB is configured with HTTP (port 80) — should it be HTTPS?
2. EC2 instances inside the ASG — do they have security groups?
3. Are there any resources directly exposed to the internet that shouldn't be?
4. Are VPC flow logs enabled?`, archJSON)

	raw, err := a.llm.Invoke(ctx, securitySystemPrompt, userPrompt)
	if err != nil {
		return enriched, fmt.Errorf("SecurityAgent LLM call failed: %w", err)
	}

	jsonBytes, err := ParseJSONBlock(raw)
	if err != nil {
		return enriched, fmt.Errorf("SecurityAgent parse failed: %w", err)
	}

	var findings []SecurityFinding
	if err := json.Unmarshal(jsonBytes, &findings); err != nil {
		return enriched, fmt.Errorf("SecurityAgent unmarshal failed: %w", err)
	}

	enriched.SecurityFindings = findings
	log.Printf("[SecurityAgent] Found %d security findings (%d critical)\n",
		len(findings), countFindingsBySeverity(findings, "critical"))

	return enriched, nil
}

func countFindingsBySeverity(findings []SecurityFinding, severity string) int {
	count := 0
	for _, f := range findings {
		if f.Severity == severity {
			count++
		}
	}
	return count
}
