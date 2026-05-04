# Cloud Canvas — Agentic IaC Generator (Go MVP)

A multi-agent AI system that takes your visual cloud architecture (as a JSON IR)
and runs a pipeline of specialist AI agents to produce Terraform code, security
findings, cost estimates, and architectural recommendations.

---

## Architecture

```
Frontend Canvas (React)
        │
        │  emits canvas JSON
        ▼
┌─────────────────────────────────────────────────────┐
│                 IR PARSER                           │
│  Normalizes frontend JSON → Architecture IR         │
└───────────────────┬─────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────┐
│              ORCHESTRATOR                           │
│  Sequential pipeline — passes EnrichedArchitecture  │
│  through each agent in order                        │
└──┬──────────┬──────────┬─────────┬─────────┬───────┘
   │          │          │         │         │
   ▼          ▼          ▼         ▼         ▼
[1. Validator] [2. Completion] [3. Security] [4. Terraform] [5. Cost] [6. Recommendations]
   │          │          │         │         │         │
   └──────────┴──────────┴─────────┴─────────┴─────────┘
                    │
                    ▼
            PipelineResult
        {code, findings, cost, recs}
```

## Agent Pipeline

| # | Agent | Input | Output |
|---|-------|-------|--------|
| 1 | **ValidatorAgent** | Architecture IR | `[]ValidationIssue` |
| 2 | **CompletionAgent** | IR + Issues | Completed IR + `[]InferredResource` |
| 3 | **SecurityAgent** | Completed IR | `[]SecurityFinding` |
| 4 | **TerraformAgent** | Completed IR + Security | HCL code string |
| 5 | **CostAgent** | Completed IR + Inferred | `CostReport` |
| 6 | **RecommendationsAgent** | All above | `[]Recommendation` |

---

## IR JSON Schema

The **Intermediate Representation (IR)** is the contract between your frontend and agents.

### Architecture (root)
```json
{
  "architecture_id": "arch-001",
  "name": "My Production App",
  "cloud_provider": "aws",
  "iac_target": "terraform",
  "constraints": {
    "compliance": ["HIPAA"],
    "budget_limit_monthly": 1000,
    "regions": ["us-east-1"]
  },
  "components": [ ...Component[] ],
  "edges": [ ...Edge[] ],
  "variables": [],
  "outputs": []
}
```

### Component
```json
{
  "id": "vpc-2",
  "type": "vpc",
  "parent_id": "region-1",
  "children": ["subnet-4", "subnet-5"],
  "is_visual_only": false,
  "properties": {
    "name": "project-vpc",
    "cidr": "10.0.0.0/16"
  },
  "status": "valid"
}
```

### Edge
```json
{
  "id": "edge-1",
  "source": "alb-6",
  "target": "autoscaling-group-9",
  "protocol": "HTTP",
  "port": 80
}
```

---

## Setup & Usage

### Prerequisites
- Go 1.22+
- `GEMINI_API_KEY` environment variable

### Build
```bash
cd tmp-agents-test
go build
```

### CLI Mode (run on a JSON file)
```bash
export GEMINI_API_KEY=your-gemini-key
./tmp-agents-test run sample_ir.json
# Terraform code written to output/main.tf
```

### Server Mode (HTTP API)
```bash
export GEMINI_API_KEY=your-gemini-key
./tmp-agents-test server
# POST http://localhost:8080/analyze  with IR JSON body
```

### Via Makefile
```bash
make build
make run-cli          # CLI on sample IR
make run-server       # HTTP server on :8080
make curl-test        # test the running server
```

---

## Output Structure

```json
{
  "architecture_id": "arch-demo-001",
  "run_duration_ms": 12400,
  "validation_issues": [
    {
      "component_id": "subnet-4",
      "severity": "warning",
      "field": "name",
      "message": "Subnet name 'lijhiol' looks auto-generated",
      "suggestion": "Use a descriptive name like 'project-subnet-public-us-east-1a'"
    }
  ],
  "security_findings": [
    {
      "component_id": "alb-6",
      "severity": "high",
      "rule": "AWS-SEC-1",
      "description": "ALB listener uses HTTP (port 80) — traffic is unencrypted",
      "remediation": "Add HTTPS listener on port 443 with ACM certificate",
      "auto_fixed": false
    }
  ],
  "inferred_resources": [
    { "id": "igw-inferred", "type": "internet-gateway", ... },
    { "id": "rt-public-inferred", "type": "route-table", ... },
    { "id": "sg-alb-inferred", "type": "security-group", ... },
    { "id": "tg-inferred", "type": "target-group", ... },
    { "id": "lt-inferred", "type": "launch-template", ... }
  ],
  "generated_code": "terraform {\n  required_version = ...\n  ...",
  "cost_estimate": {
    "total_monthly_usd": 87.50,
    "total_yearly_usd": 1050.00,
    "breakdown": [
      { "component_id": "ec2-10", "name": "web-server (x2 via ASG)", "monthly_usd": 15.18, "basis": "2 × t3.micro × 730hrs × $0.0104/hr" },
      { "component_id": "alb-6",  "name": "my-alb", "monthly_usd": 22.00, "basis": "$0.008/hr + LCU charges" },
      { "component_id": "igw-inferred", "name": "NAT Gateway", "monthly_usd": 35.92, "basis": "$0.045/hr + 100GB × $0.045/GB" }
    ],
    "warnings": ["NAT Gateway cost depends heavily on data transfer volume"]
  },
  "recommendations": [
    {
      "type": "cost",
      "priority": "high",
      "title": "Use Reserved Instances for EC2",
      "description": "With 1-year Reserved Instances, t3.micro drops from $0.0104 to ~$0.006/hr",
      "saving": "~$25/month"
    },
    {
      "type": "reliability",
      "priority": "high",
      "title": "Add second AZ for high availability",
      "description": "Add subnets in us-east-1b and configure ALB across both AZs"
    }
  ]
}
```

---

## Extending the Pipeline

### Add a new agent
1. Create `internal/agents/myagent.go` implementing the `Agent` interface
2. Add it to the pipeline in `internal/orchestrator/orchestrator.go`
3. Add its output type to `internal/ir/types.go`

### Add a new IaC target (Pulumi, CDK)
1. Create `internal/agents/pulumi.go` mirroring `terraform.go`
2. The orchestrator selects which generator to run based on `arch.IaCTarget`

### Support Azure / GCP
1. Extend `ir.Architecture.CloudProvider` routing in the orchestrator
2. Each agent's system prompt handles multi-cloud via provider context injection
