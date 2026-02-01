# Scenario 10: Pricing with Hidden Dependency Costs

This scenario demonstrates the complete pricing calculation system including hidden/implicit dependency costs.

## Overview

This use case creates a new architecture and calculates pricing that includes:
- **Base resource costs**: Direct costs for resources (NAT Gateway, EC2 instances)
- **Hidden dependency costs**: Implicit costs for resources automatically created or required

## What It Demonstrates

### Architecture Created
- **VPC** (10.0.0.0/16)
- **Public Subnet** (10.0.1.0/24)
- **Private Subnet** (10.0.2.0/24)
- **NAT Gateway** (creates hidden Elastic IP)
- **EC2 Instance t3.micro** (creates hidden EBS root volume and Network Interface)
- **EC2 Instance t3.small** (creates hidden EBS root volume and Network Interface)

### Hidden Dependencies Demonstrated

1. **NAT Gateway → Elastic IP**
   - When a NAT Gateway is created without an `allocationId`, AWS automatically creates an Elastic IP
   - This Elastic IP is attached to the NAT Gateway (free when attached)
   - The hidden dependency resolver detects this and includes it in pricing

2. **EC2 Instance → EBS Root Volume**
   - Every EC2 instance requires a root EBS volume
   - Default size is 8GB if not specified in metadata
   - The hidden dependency resolver creates a virtual EBS volume resource and calculates its cost

3. **EC2 Instance → Network Interface**
   - Every EC2 instance requires a network interface
   - Free when attached to an instance
   - The hidden dependency resolver detects this but marks it as free

## Running the Scenario

```bash
cd backend
go run ./cmd/api/main.go -scenario=10
```

## Output

The scenario produces:
1. **Console output** with detailed pricing breakdown showing:
   - Base costs for each resource
   - Hidden dependency costs (marked with 🔗)
   - Total base costs
   - Total hidden dependency costs
   - Percentage of total cost from hidden dependencies

2. **JSON report** saved to:
   - `backend/pkg/usecases/json-response-pricing-with-hidden-costs.json`

## Example Output

```
📊 Architecture Cost Estimate (Monthly - 720h0m0s)
----------------------------------------------------------------------------------------------------
Resource                                 Type                       Base Cost      Total Cost
----------------------------------------------------------------------------------------------------
web-server-1                             EC2                  $       7.4880 $       8.1280
    └─ EC2 Instance Hourly                 per_hour   720.0000 × $  0.010400 = $    7.4880
    🔗 EBS Volume Storage (ebs_volume)     per_gb       8.0000 × $  0.080000 = $    0.6400

demo-nat-gateway                         NATGateway           $      32.4450 $      36.0450
    └─ NAT Gateway Hourly                  per_hour   720.0000 × $  0.045000 = $   32.4000
    🔗 Elastic IP Hourly (Unattached)...   per_hour   720.0000 × $  0.005000 = $    3.6000

----------------------------------------------------------------------------------------------------
TOTAL BASE COSTS                                              $      47.4210 $      47.4210
TOTAL HIDDEN DEPENDENCY COSTS                                                 $       5.8400
====================================================================================================

💰 TOTAL MONTHLY COST: $53.26 USD
📅 ESTIMATED YEARLY COST: $639.13 USD

🔗 HIDDEN DEPENDENCY COSTS BREAKDOWN:
  Total Hidden Costs: $5.84 USD
  Percentage of Total: 11.0%
```

## Key Features

1. **Automatic Detection**: Hidden dependencies are automatically detected based on resource types
2. **Database-Driven**: Uses pricing rates from the `pricing_rates` table
3. **Conditional Logic**: Hidden dependencies can have conditions (e.g., only create Elastic IP if `allocationId` is not provided)
4. **Quantity Expressions**: Supports dynamic quantity calculation (e.g., EBS volume size from metadata)
5. **Visual Markers**: Hidden dependency costs are marked with 🔗 in the output

## Database Requirements

This scenario requires:
- Migration `00005_add_pricing_tables.sql` to be applied
- Pricing seed data to be populated (run `go run ./cmd/seed_pricing/main.go`)

## Technical Details

### Hidden Dependency Resolution Flow

1. **Resource Processing**: When a resource is processed for pricing
2. **Dependency Lookup**: System queries `hidden_dependencies` table for the resource type
3. **Condition Evaluation**: Evaluates any condition expressions (e.g., `metadata.allocationId == null`)
4. **Quantity Calculation**: Calculates quantity using expressions (e.g., `metadata.size_gb`)
5. **Virtual Resource Creation**: Creates a virtual resource for the hidden dependency
6. **Cost Calculation**: Calculates cost for the virtual resource using pricing rates
7. **Cost Aggregation**: Adds hidden dependency costs to the base resource cost

### Pricing Rate Lookup

The system:
1. Queries `pricing_rates` table for active rates matching:
   - Provider (aws)
   - Resource type
   - Region (or NULL for default)
   - Effective date range
2. Falls back to hardcoded rates if database rates not found
3. Supports multiple pricing models: `per_hour`, `per_gb`, `per_request`, etc.

## Extending Hidden Dependencies

To add new hidden dependencies:

1. **Add to Database**: Insert into `hidden_dependencies` table
2. **Or Add to Code**: Update `backend/internal/cloud/aws/pricing/hidden_deps/definitions.go`

Example:
```sql
INSERT INTO hidden_dependencies (
    provider, parent_resource_type, child_resource_type,
    quantity_expression, condition_expression, is_attached, description
) VALUES (
    'aws', 'rds_instance', 'ebs_volume',
    'metadata.allocated_storage', '', true,
    'RDS instance requires storage volume based on allocated_storage'
);
```
