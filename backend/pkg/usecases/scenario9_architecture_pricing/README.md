# Scenario 9: Architecture Pricing

This scenario demonstrates the complete architecture pricing flow, including:

1. **Diagram Processing**: Read and parse a diagram JSON file
2. **Pricing Calculation**: Calculate cost estimates for each resource in the architecture
3. **Pricing Persistence**: Save pricing data to the database
4. **Pricing Report**: Generate a detailed pricing breakdown report

## Features Tested

- `ProcessDiagram` with `PricingDuration` parameter
- `CalculateArchitectureCost` for total architecture cost estimation
- `CalculateResourceCost` for individual resource cost estimation
- `PersistResourcePricing` for saving resource-level pricing to database
- `PersistProjectPricing` for saving project-level pricing to database
- `GetProjectPricing` for retrieving pricing from database

## Supported Resource Types

The pricing calculator supports the following AWS resource types:

| Resource Type | Pricing Model | Description |
|--------------|---------------|-------------|
| `ec2_instance` | Per Hour | EC2 instance hourly rate based on instance type |
| `ebs_volume` | Per GB/Month | EBS volume storage cost based on volume type |
| `s3_bucket` | Per GB/Month + Requests | S3 storage, PUT/GET requests, data transfer |
| `load_balancer` | Per Hour | ALB/NLB/CLB hourly rate |
| `auto_scaling_group` | Per Hour | Based on average capacity × instance rate |
| `lambda_function` | Per Request + GB-seconds | Compute time and request count |
| `nat_gateway` | Per Hour + Per GB | Hourly rate + data processing |
| `elastic_ip` | Per Hour (unattached) | Free when attached to running instance |

## Usage

### Running the Scenario

```go
package main

import (
    "context"
    "log"
    
    scenario9 "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario9_architecture_pricing"
)

func main() {
    ctx := context.Background()
    if err := scenario9.ArchitecturePricingRunner(ctx); err != nil {
        log.Fatalf("Scenario 9 failed: %v", err)
    }
}
```

### Using Pricing in Your Code

```go
// Process diagram with monthly pricing calculation
processReq := &serverinterfaces.ProcessDiagramRequest{
    JSONData:        diagramJSON,
    UserID:          userID,
    ProjectName:     "My Project",
    IACToolID:       1,
    CloudProvider:   "aws",
    Region:          "us-east-1",
    PricingDuration: 720 * time.Hour, // 30 days for monthly estimate
}

result, err := srv.PipelineOrchestrator.ProcessDiagram(ctx, processReq)
if err != nil {
    return err
}

// Access pricing estimate
if result.PricingEstimate != nil {
    fmt.Printf("Total Monthly Cost: $%.2f\n", result.PricingEstimate.TotalCost)
    
    // Access individual resource costs
    for _, resEstimate := range result.PricingEstimate.ResourceEstimates {
        fmt.Printf("  %s: $%.2f\n", resEstimate.ResourceName, resEstimate.TotalCost)
    }
}
```

## Output

The scenario generates:

1. **Console Output**: Detailed pricing breakdown with resource costs
2. **JSON Report**: `json-response-architecture-pricing.json` with structured pricing data

### Sample Console Output

```
====================================================================================================
SCENARIO 9: Architecture Pricing (Process Diagram → Calculate Pricing → Persist)
====================================================================================================

[Step 1] Initializing service layer server...
✓ Service layer server initialized successfully

[Step 2] Reading diagram JSON file...
✓ Read diagram JSON from: /path/to/json-request-fiagram-complete.json (12345 bytes)

[Step 3] Processing diagram with pricing calculation...
✓ Diagram processed and saved to database
  Project ID: 12345678-1234-1234-1234-123456789012
  Success: true
  Message: Diagram processed successfully. Project created with ID: 12345678-1234-1234-1234-123456789012. Estimated monthly cost: $45.67 USD

[Step 4] Pricing Breakdown:
--------------------------------------------------------------------------------

📊 Architecture Cost Estimate (Monthly - 720h0m0s)
--------------------------------------------------------------------------------

Resource                                 Type            Monthly Cost
--------------------------------------------------------------------------------
web-server                               ec2_instance    $7.4880 USD
  └─ EC2 Instance Hourly                 per_hour        720.0000 × $0.010400 = $7.4880
database-server                          ec2_instance    $69.1200 USD
  └─ EC2 Instance Hourly                 per_hour        720.0000 × $0.096000 = $69.1200
--------------------------------------------------------------------------------

💰 TOTAL MONTHLY COST: $76.61 USD
   Provider: aws
   Region: us-east-1
   Period: monthly

📅 ESTIMATED YEARLY COST: $919.32 USD
```

## Database Tables

Pricing data is stored in the following tables:

- `project_pricing`: Project-level cost estimates
- `resource_pricing`: Individual resource cost estimates
- `pricing_components`: Detailed cost breakdown components

## Notes

- Pricing is calculated based on static rates (not real-time AWS pricing API)
- Visual-only resources (VPC, Subnet, etc.) are skipped in pricing calculation
- Free tier is considered for some resources (Lambda, S3 data transfer)
- Regional pricing multipliers are applied where applicable
