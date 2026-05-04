// Package ir defines the Intermediate Representation (IR) JSON schema
// that the frontend canvas emits and the agent pipeline consumes.
package main

// ============================================================
// ROOT — the full document emitted by the canvas
// ============================================================

type Architecture struct {
	ArchitectureID string      `json:"architecture_id"`
	Name           string      `json:"name"`
	CloudProvider  string      `json:"cloud_provider"` // "aws" | "azure" | "gcp"
	IaCTarget      string      `json:"iac_target"`     // "terraform" | "pulumi" | "cdk"
	Constraints    Constraints `json:"constraints"`
	Components     []Component `json:"components"`
	Edges          []Edge      `json:"edges"`
	Variables      []Variable  `json:"variables"`
	Outputs        []Output    `json:"outputs"`
}

// ============================================================
// CONSTRAINTS — architect-level requirements
// ============================================================

type Constraints struct {
	Compliance         []string `json:"compliance"`           // ["HIPAA","SOC2","PCI-DSS"]
	BudgetLimitMonthly float64  `json:"budget_limit_monthly"` // USD
	Regions            []string `json:"regions"`
}

// ============================================================
// COMPONENT — every node on the canvas
// ============================================================

type Component struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"` // "vpc","subnet","ec2","alb","rds", etc.
	ParentID     string                 `json:"parent_id,omitempty"`
	Children     []string               `json:"children,omitempty"`
	Properties   map[string]interface{} `json:"properties"`
	IsVisualOnly bool                   `json:"is_visual_only"` // true = representation only (e.g. EC2 inside ASG)
	Status       string                 `json:"status"`         // "valid" | "warning" | "error"
}

// ============================================================
// EDGE — connections between components
// ============================================================

type Edge struct {
	ID       string `json:"id"`
	Source   string `json:"source"`
	Target   string `json:"target"`
	Label    string `json:"label,omitempty"`
	Protocol string `json:"protocol,omitempty"` // "HTTP","HTTPS","TCP","UDP"
	Port     int    `json:"port,omitempty"`
}

// ============================================================
// VARIABLE / OUTPUT — Terraform-level declarations
// ============================================================

type Variable struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`
	Sensitive   bool        `json:"sensitive,omitempty"`
}

type Output struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Sensitive   bool   `json:"sensitive,omitempty"`
}

// ============================================================
// AGENT PIPELINE — enriched IR passed between agents
// ============================================================

// EnrichedArchitecture wraps the original IR with agent findings.
type EnrichedArchitecture struct {
	Original          *Architecture     `json:"original"`
	Completed         *Architecture     `json:"completed"` // after completion agent
	ValidationIssues  []ValidationIssue `json:"validation_issues"`
	InferredResources []Component       `json:"inferred_resources"` // implicit resources added
	GeneratedCode     string            `json:"generated_code"`
	CostEstimate      *CostReport       `json:"cost_estimate"`
	Recommendations   []Recommendation  `json:"recommendations"`
	SecurityFindings  []SecurityFinding `json:"security_findings"`
}

// ============================================================
// AGENT OUTPUT TYPES
// ============================================================

type ValidationIssue struct {
	ComponentID string `json:"component_id"`
	Severity    string `json:"severity"` // "error" | "warning" | "info"
	Field       string `json:"field"`
	Message     string `json:"message"`
	Suggestion  string `json:"suggestion"`
}

type SecurityFinding struct {
	ComponentID string `json:"component_id"`
	Severity    string `json:"severity"` // "critical"|"high"|"medium"|"low"
	Rule        string `json:"rule"`
	Description string `json:"description"`
	Remediation string `json:"remediation"`
	AutoFixed   bool   `json:"auto_fixed"`
}

type CostReport struct {
	TotalMonthlyUSD float64        `json:"total_monthly_usd"`
	TotalYearlyUSD  float64        `json:"total_yearly_usd"`
	Breakdown       []CostLineItem `json:"breakdown"`
	Warnings        []string       `json:"warnings"`
}

type CostLineItem struct {
	ComponentID  string  `json:"component_id"`
	ResourceType string  `json:"resource_type"`
	Name         string  `json:"name"`
	MonthlyUSD   float64 `json:"monthly_usd"`
	Basis        string  `json:"basis"` // human-readable: "730 hrs × $0.0416/hr"
}

type Recommendation struct {
	Type        string `json:"type"`     // "cost"|"security"|"reliability"|"performance"
	Priority    string `json:"priority"` // "high"|"medium"|"low"
	Title       string `json:"title"`
	Description string `json:"description"`
	Saving      string `json:"saving,omitempty"` // e.g. "~$45/month"
}
