# Plan: Adding a New Cloud Resource

## Overview
Adding a new resource requires modifications across **6 layers**:

| # | Layer | Purpose |
|---|-------|---------|
| 1 | **Domain** | Cloud-agnostic resource definitions |
| 2 | **Cloud Provider** | AWS-specific models & adapters, Virtual Services |
| 3 | **Terraform Mapping** | Generate IaC code |
| 4 | **Rules/Constraints** | Validation logic |
| 5 | **Pricing** | Cost estimation |
| 6 | **Diagram Validation** | UI input validation |

---

## Layer-by-Layer Implementation Plan

### Layer 1: Domain (`internal/domain/`)

| File | Action |
|------|--------|
| `resource/categories.go` | Add category if needed |
| `resource/<category>/newresource.go` | Create domain model |

**Example:** VPC endpoint defined at `internal/domain/resource/networking/vpc_endpoint.go`

---

### Layer 2: AWS Models (`internal/cloud/aws/models/`)

| File | Action |
|------|--------|
| `models/<category>/newresource.go` | Create AWS model |
| `models/<category>/outputs/...` | Create output model |

---

### Layer 3: AWS Provider Integration

| File | Action |
|------|--------|
| `cloud/aws/inventory/resources.go` | Add resource classification (IRType, aliases) |
| `cloud/aws/architecture/resource_type_mapper.go` | Map IRType → ResourceType |
| `cloud/aws/adapters/<category>/adapter.go` | Create domain→AWS adapter |
| `cloud/aws/adapters/<category>/factory.go` | Create factory function |
| `cloud/aws/services/<category>/service.go` | Create AWS SDK service |

---

### Layer 4: Terraform Mapping (`internal/cloud/aws/mapper/`)

| File | Action |
|------|--------|
| `mapper/terraform/mapper.go` | Register mapper function |
| `mapper/<category>/newresource_mapper.go` | Create mapper |

**Pattern:** Convert domain resource → Terraform block with `aws_<resource>` syntax

---

### Layer 5: Rules & Constraints (`internal/cloud/aws/rules/`)

| File | Action |
|------|--------|
| `rules/defaults.go` | Add constraint records |

**Constraint types available:**
- `requires_parent` - Must have parent resource
- `allowed_parent` - Only specific parents allowed
- `requires_region` - Region required (true/false)
- `max_children` / `min_children` - Child count limits
- `allowed_dependencies` / `forbidden_dependencies`

---

### Layer 6: Pricing (`internal/cloud/aws/pricing/`)

| File | Action |
|------|--------|
| `pricing/<category>/newresource.go` | Create pricing calculator |
| `pricing/service.go` | Add fallback handler |

---

### Layer 7: Diagram Validation (`internal/diagram/validator/schema/`)

| File | Action |
|------|--------|
| `schema/aws_schemas.go` | Register validation schema |

**Schema defines:** required fields, field types, valid parent/child types

---

## Example Files to Reference

For a complete reference, look at the recent **VPC Endpoint** addition:
- Domain: `internal/domain/resource/networking/vpc_endpoint.go`
- AWS Model: `internal/cloud/aws/models/networking/vpc_endpoint.go`
- Terraform Mapper: `internal/cloud/aws/mapper/terraform/vpc_endpoint_mapper.go`
- Pricing: `internal/cloud/aws/pricing/networking/vpc_endpoint.go`

---

## Summary Table

| Layer | # of Files | Complexity |
|-------|-----------|------------|
| Domain | 1-2 | Low |
| AWS Models | 1-2 | Low |
| AWS Integration | 3-4 | Medium |
| Terraform Mapping | 1-2 | Medium |
| Rules | 1 | Low |
| Pricing | 1-2 | Medium |
| Schema Validation | 1 | Low |
| **Total** | **9-14** | - |

---

## Step-by-Step Implementation Guide

### Step 1: Define Domain Resource
Create the cloud-agnostic domain model in `internal/domain/resource/<category>/newresource.go`

### Step 2: Create AWS Models
Create AWS-specific models in `internal/cloud/aws/models/<category>/newresource.go`

### Step 3: Register in Inventory
Add resource classification in `internal/cloud/aws/inventory/resources.go`

### Step 4: Add Architecture Mapping
Map the resource type in `internal/cloud/aws/architecture/resource_type_mapper.go`

### Step 5: Create Adapter
Build the adapter in `internal/cloud/aws/adapters/<category>/adapter.go`

### Step 6: Implement Terraform Mapper
Create the mapper in `internal/cloud/aws/mapper/<category>/newresource_mapper.go` and register in `mapper/terraform/mapper.go`

### Step 7: Add Constraints
Define rules in `internal/cloud/aws/rules/defaults.go`

### Step 8: Implement Pricing
Add pricing logic in `internal/cloud/aws/pricing/<category>/newresource.go`

### Step 9: Add Validation Schema
Register the schema in `internal/diagram/validator/schema/aws_schemas.go`

---

## Common Patterns

### Resource Classification Pattern
```go
{
    Category:     resource.CategoryNetworking,
    ResourceName: "VPCEndpoint",
    IRType:       "vpc-endpoint",
    Aliases:      []string{"vpc-endpoint", "vpc_endpoint", "vpce"},
}
```

### Terraform Mapper Pattern
```go
func (m *AWSMapper) mapNewResource(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
    attrs := map[string]tfmapper.TerraformValue{
        "vpc_id": tfReference(res.ParentID, "aws_vpc", "id"),
    }
    
    block := tfmapper.TerraformBlock{
        Kind:   "resource",
        Labels: []string{"aws_new_resource", sanitizeName(res.Name)},
        Attributes: attrs,
    }
    
    return []tfmapper.TerraformBlock{block}, nil
}
```

### Constraint Pattern
```go
{ResourceType: "NewResource", ConstraintType: "requires_parent", ConstraintValue: "VPC"},
{ResourceType: "NewResource", ConstraintType: "allowed_parent", ConstraintValue: "VPC"},
{ResourceType: "NewResource", ConstraintType: "requires_region", ConstraintValue: "true"},
```

### Pricing Pattern
```go
func CalculateNewResourceCost(duration time.Duration, region string) float64 {
    multiplier := GetRegionMultiplier(region)
    return baseRate * multiplier * duration.Hours()
}
```
