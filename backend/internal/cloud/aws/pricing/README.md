# AWS Pricing Feature

This package implements a flexible, scalable pricing system for AWS networking resources. The pricing feature supports multiple pricing models (per hour, per GB, per request, etc.) and can be easily extended to other resource categories (compute, serverless, storage, etc.) in the future.

## Overview

The pricing system provides:
- **Cost Estimation**: Estimate costs before resource creation
- **Resource Pricing Information**: Get detailed pricing breakdowns for resources
- **Architecture Cost Calculation**: Calculate total costs for multiple resources
- **Multiple Pricing Models**: Support for various pricing structures (hourly, per GB, per request, etc.)
- **Extensible Design**: Easy to add new resource types and pricing models

## Architecture

### Directory Structure

```
backend/internal/
├── domain/
│   └── pricing/
│       ├── models.go          # Core pricing models (PriceComponent, ResourcePricing, CostEstimate)
│       ├── calculator.go       # PricingCalculator interface
│       ├── service.go         # PricingService interface
│       └── types.go           # Pricing types and enums (PricingModel, Currency, Period)
│
└── cloud/
    └── aws/
        └── pricing/
            ├── models.go              # AWS-specific pricing models
            ├── calculator.go          # AWS pricing calculator implementation
            ├── service.go             # AWS pricing service implementation
            ├── rates.go               # Static pricing rates
            └── networking/
                ├── nat_gateway.go     # NAT Gateway pricing calculations
                ├── elastic_ip.go      # Elastic IP pricing calculations
                ├── network_interface.go # Network Interface pricing calculations
                ├── data_transfer.go   # Data transfer pricing calculations
                └── *_test.go          # Test files for each resource type
```

### Component Overview

#### Domain Layer (`internal/domain/pricing/`)

The domain layer defines cloud-agnostic pricing interfaces and models:

- **`types.go`**: Defines pricing model types (PerHour, PerGB, PerRequest, OneTime, Tiered, Percentage), currencies, and periods
- **`models.go`**: Core pricing data structures:
  - `PriceComponent`: Individual pricing component (e.g., "NAT Gateway Hourly", "Data Transfer Out")
  - `ResourcePricing`: Complete pricing information for a resource type
  - `CostComponent`: Calculated cost component in an estimate
  - `CostEstimate`: Final cost estimate with breakdown
- **`calculator.go`**: `PricingCalculator` interface for cost calculations
- **`service.go`**: `PricingService` interface for pricing operations

#### AWS Implementation Layer (`internal/cloud/aws/pricing/`)

The AWS layer implements the domain interfaces with AWS-specific pricing:

- **`models.go`**: AWS-specific pricing models (`AWSPricingRate`, `AWSResourcePricing`)
- **`rates.go`**: Static pricing rates for networking resources
- **`calculator.go`**: `AWSPricingCalculator` - implements cost calculation logic
- **`service.go`**: `AWSPricingService` - implements pricing service operations
- **`networking/`**: Resource-specific pricing calculations

## Pricing Models

The system supports six pricing models:

### 1. Per Hour (PerHour)
**Use Case**: Resources charged by the hour (e.g., NAT Gateway, Elastic IP when unattached)

**Example**: NAT Gateway = $0.045/hour

**Calculation**: `rate * hours`

**Resources Using This Model**:
- NAT Gateway: $0.045/hour
- Elastic IP (unattached): $0.005/hour
- Network Interface (unattached): $0.01/hour

### 2. Per GB (PerGB)
**Use Case**: Data transfer and data processing charges

**Example**: Data Transfer Out = $0.09/GB (after free tier)

**Calculation**: `rate * gigabytes`

**Resources Using This Model**:
- NAT Gateway Data Processing: $0.045/GB
- Data Transfer Outbound: $0.09/GB (first 1GB free/month)
- Data Transfer Inter-AZ: $0.01/GB

### 3. Per Request (PerRequest)
**Use Case**: API calls, function invocations (future use)

**Example**: API Gateway = $0.000001/request

**Calculation**: `rate * request_count`

**Note**: Currently not used for networking resources, but available for future compute/serverless resources.

### 4. One-Time (OneTime)
**Use Case**: Setup fees, provisioning fees

**Example**: Some services have one-time setup costs

**Calculation**: `rate` (single charge)

**Note**: Currently not used for networking resources.

### 5. Tiered (Tiered)
**Use Case**: Pricing with free tiers or volume discounts

**Example**: First 1GB free, then $0.09/GB

**Calculation**: Complex logic based on tiers

**Resources Using This Model**:
- Data Transfer Outbound: First 1GB/month free, then $0.09/GB

### 6. Percentage (Percentage)
**Use Case**: Percentage-based fees

**Example**: 2% of base resource cost

**Calculation**: `base_cost * (percentage / 100)`

**Note**: Currently not used for networking resources.

## Supported Resources

### Current Networking Resources

1. **NAT Gateway**
   - Hourly rate: $0.045/hour
   - Data processing: $0.045/GB
   - **Example**: 30 days (720 hours) = $32.40 + data processing costs

2. **Elastic IP**
   - Hourly rate (unattached): $0.005/hour
   - **Free when attached** to a running instance
   - **Example**: 30 days unattached = $3.60

3. **Network Interface (ENI)**
   - Hourly rate (unattached): $0.01/hour
   - **Free when attached** to an instance
   - **Example**: 30 days unattached = $7.20

4. **Data Transfer**
   - Inbound: **Free**
   - Outbound: $0.09/GB (first 1GB/month free)
   - Inter-AZ: $0.01/GB
   - **Example**: 100GB outbound = (100 - 1) * $0.09 = $8.91

## Usage Examples

### Get Resource Pricing Information

```go
import (
    "context"
    awspricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing"
)

// Create pricing service
pricingService := awspricing.NewAWSPricingService()
ctx := context.Background()

// Get NAT Gateway pricing
pricing, err := pricingService.GetPricing(ctx, "nat_gateway", "aws", "us-east-1")
if err != nil {
    // handle error
}

// Access pricing components
for _, component := range pricing.Components {
    fmt.Printf("Component: %s, Rate: $%.3f/%s\n", 
        component.Name, component.Rate, component.Unit)
}
```

### Estimate Resource Cost

```go
import (
    "time"
    "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// Create a resource
res := &resource.Resource{
    Type: resource.ResourceType{
        Name: "nat_gateway",
    },
    Provider: "aws",
    Region:   "us-east-1",
}

// Estimate cost for 30 days
duration := 720 * time.Hour // 30 days
estimate, err := pricingService.EstimateCost(ctx, res, duration)
if err != nil {
    // handle error
}

fmt.Printf("Total Cost: $%.2f\n", estimate.TotalCost)
fmt.Printf("Currency: %s\n", estimate.Currency)
fmt.Printf("Period: %s\n", estimate.Period)

// View breakdown
for _, component := range estimate.Breakdown {
    fmt.Printf("%s: %.2f %s @ $%.3f/%s = $%.2f\n",
        component.ComponentName,
        component.Quantity,
        component.Model,
        component.UnitRate,
        component.Model,
        component.Subtotal)
}
```

### Estimate Architecture Cost

```go
// Create multiple resources
resources := []*resource.Resource{
    {
        Type: resource.ResourceType{Name: "nat_gateway"},
        Provider: "aws",
        Region:   "us-east-1",
    },
    {
        Type: resource.ResourceType{Name: "elastic_ip"},
        Provider: "aws",
        Region:   "us-east-1",
    },
    {
        Type: resource.ResourceType{Name: "network_interface"},
        Provider: "aws",
        Region:   "us-east-1",
    },
}

// Calculate total architecture cost
duration := 720 * time.Hour // 30 days
estimate, err := pricingService.EstimateArchitectureCost(ctx, resources, duration)
if err != nil {
    // handle error
}

fmt.Printf("Total Architecture Cost: $%.2f\n", estimate.TotalCost)
```

### Using Adapter for Pricing

```go
import (
    domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
)

// Get adapter (already configured with AWS service)
adapter := domainnetworking.NewAWSNetworkingAdapter(awsService)

// Estimate cost via adapter
config := map[string]interface{}{
    "region": "us-east-1",
}
estimate, err := adapter.EstimateResourceCost(ctx, "nat_gateway", config, 720*time.Hour)

// Get pricing information
pricing, err := adapter.GetResourcePricing(ctx, "nat_gateway", "us-east-1")
```

## Testing

### Running Tests

#### Run All Pricing Tests
```bash
go test ./internal/cloud/aws/pricing/... -v
```

#### Run Specific Test Package
```bash
# Test networking pricing functions
go test ./internal/cloud/aws/pricing/networking/... -v

# Test calculator
go test ./internal/cloud/aws/pricing/... -run TestAWSPricingCalculator -v

# Test service
go test ./internal/cloud/aws/pricing/... -run TestAWSPricingService -v
```

#### Run with Coverage
```bash
go test ./internal/cloud/aws/pricing/... -cover
```

#### Run Specific Test
```bash
go test ./internal/cloud/aws/pricing/networking/... -run TestCalculateNATGatewayCost -v
```

### Test Structure

Tests follow a **table-driven design** pattern for consistency and maintainability:

```go
func TestCalculateNATGatewayCost(t *testing.T) {
    tests := []struct {
        name            string
        duration        time.Duration
        dataProcessedGB float64
        region          string
        expectedCost    float64
    }{
        {
            name:            "nat-gateway-1-hour-no-data",
            duration:        1 * time.Hour,
            dataProcessedGB: 0.0,
            region:          "us-east-1",
            expectedCost:    0.045,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cost := CalculateNATGatewayCost(tt.duration, tt.dataProcessedGB, tt.region)
            if cost != tt.expectedCost {
                t.Errorf("Expected cost %.2f, got %.2f", tt.expectedCost, cost)
            }
        })
    }
}
```

### Test Coverage

Current test coverage:
- **Main pricing package**: 79.7% coverage
- **Networking pricing**: 98.0% coverage

### Test Files

#### Networking Pricing Tests (`networking/*_test.go`)

1. **`data_transfer_test.go`**
   - `TestCalculateDataTransferCost`: 8 test cases
     - Inbound (free)
     - Outbound within free tier
     - Outbound exceeding free tier
     - Inter-AZ transfer
     - Large amounts
   - `TestGetDataTransferPricing`: Validates pricing information structure

2. **`nat_gateway_test.go`**
   - `TestCalculateNATGatewayCost`: 5 test cases
     - 1 hour, 24 hours, 720 hours (30 days)
     - With and without data processing
   - `TestGetNATGatewayPricing`: Validates pricing components

3. **`elastic_ip_test.go`**
   - `TestCalculateElasticIPCost`: 5 test cases
     - Attached (free) vs unattached
     - Various durations
   - `TestGetElasticIPPricing`: Validates pricing information

4. **`network_interface_test.go`**
   - `TestCalculateNetworkInterfaceCost`: 6 test cases
     - Attached (free) vs unattached
     - Various durations (1 hour, 24 hours, 168 hours, 720 hours)
   - `TestGetNetworkInterfacePricing`: Validates pricing information

#### Calculator Tests (`calculator_test.go`)

- `TestAWSPricingCalculator_CalculateResourceCost`: 6 test cases
  - NAT Gateway, Elastic IP, Network Interface
  - Error cases (unsupported provider, unsupported resource type)
- `TestAWSPricingCalculator_CalculateArchitectureCost`: Multiple resources aggregation
- `TestAWSPricingCalculator_GetResourcePricing`: 5 test cases for all resource types

#### Service Tests (`service_test.go`)

- `TestAWSPricingService_GetPricing`: 6 test cases
- `TestAWSPricingService_EstimateCost`: 2 test cases
- `TestAWSPricingService_EstimateArchitectureCost`: Multiple resources
- `TestAWSPricingService_ListSupportedResources`: 2 test cases

#### Adapter Integration Tests (`adapters/networking/adapter_pricing_test.go`)

- `TestAWSNetworkingAdapter_EstimateResourceCost`: 3 test cases
- `TestAWSNetworkingAdapter_GetResourcePricing`: 3 test cases

### Test Scenarios Covered

✅ **Hourly Pricing Calculations**
- 1 hour, 24 hours, 720 hours (30 days)
- Multiple resource types

✅ **Free Tier Allowances**
- Data transfer first 1GB free
- Attached resources (Elastic IP, Network Interface)

✅ **Data Processing Costs**
- NAT Gateway data processing

✅ **Multiple Resource Aggregation**
- Architecture cost calculation
- Cost breakdown by component

✅ **Error Handling**
- Unsupported providers
- Unsupported resource types
- Invalid configurations

✅ **Regional Variations**
- Different AWS regions
- Regional pricing multipliers

✅ **Resource State-Based Pricing**
- Attached vs unattached resources
- Free when attached scenarios

## Adding New Resources

To add pricing for a new resource:

1. **Add pricing rates** in `rates.go`:
```go
var NetworkingPricingRates = map[string]AWSPricingRate{
    "new_resource": {
        BaseHourlyRate: 0.05,
        // ... other rates
    },
}
```

2. **Create pricing calculation function** in `networking/new_resource.go`:
```go
func CalculateNewResourceCost(duration time.Duration, region string) float64 {
    // Calculation logic
}

func GetNewResourcePricing(region string) *domainpricing.ResourcePricing {
    // Return pricing information
}
```

3. **Add to service** in `service.go`:
```go
case "new_resource":
    return networking.GetNewResourcePricing(region), nil
```

4. **Add to calculator** in `calculator.go`:
```go
case "new_resource":
    cost := networking.CalculateNewResourceCost(duration, region)
    // ... calculation logic
```

5. **Write tests** in `networking/new_resource_test.go`:
```go
func TestCalculateNewResourceCost(t *testing.T) {
    // Test cases
}
```

## Future Extensions

The pricing system is designed to be easily extended:

- **Compute Resources**: EC2, Lambda, ECS pricing
- **Serverless Resources**: API Gateway, Lambda invocations
- **Storage Resources**: S3, EBS pricing
- **Database Resources**: RDS, DynamoDB pricing
- **Dynamic Pricing**: Integration with AWS Pricing API for real-time rates
- **Cost Tracking**: Track actual costs for created resources
- **Cost Alerts**: Set up alerts for cost thresholds
- **Cost Optimization**: Suggest cheaper alternatives
- **Multi-Provider**: Support GCP, Azure pricing

## Key Design Decisions

1. **Separation of Concerns**: Pricing is separate from resource models but can be attached
2. **Provider-Specific**: AWS pricing in `cloud/aws/pricing/`, can add GCP/Azure later
3. **Extensible**: Easy to add new resource types and pricing models
4. **Flexible**: Supports static rates now, can add dynamic API integration later
5. **Scalable**: Registry pattern allows easy extension to compute, serverless, etc.

## Related Documentation

- [Domain Pricing Models](../../../domain/pricing/) - Core pricing interfaces
- [AWS Networking Models](../models/networking/) - Networking resource models
- [Networking Adapters](../adapters/networking/) - Adapter implementation
