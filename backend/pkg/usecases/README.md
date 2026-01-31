# Solution Architect Use Cases

This directory contains comprehensive solution architect scenarios that demonstrate building complete AWS architectures using networking, compute, and IAM services. All use cases use mock data to test domain models without making real AWS SDK calls.

## Overview

These use cases showcase:
- **Complete architectures** combining multiple AWS services
- **Domain model validation** and business logic
- **Resource orchestration** and dependencies
- **Cost estimation** for entire architectures
- **Best practices** for AWS solution architecture

## Available Scenarios

### Scenario 1: Basic Web Application
**File**: `scenario1_basic_web_app/basic_web_app.go`

A simple 3-tier web application architecture demonstrating:
- VPC creation with CIDR block
- Public and private subnets across multiple AZs
- Internet Gateway for public internet access
- Route tables (public and private)
- Security groups for different tiers
- EC2 instances in public subnets (web tier)
- EC2 instances in private subnets (app tier)
- Cost estimation

**Architecture**:
```
Internet
   ↓
Internet Gateway
   ↓
VPC (10.0.0.0/16)
   ├── Public Subnets (10.0.1.0/24, 10.0.2.0/24)
   │   └── Web Tier EC2 Instances (t3.micro)
   └── Private Subnets (10.0.10.0/24, 10.0.11.0/24)
       └── App Tier EC2 Instances (t3.small)
```

**Run**:
```bash
go run backend/pkg/usecases/scenario1_basic_web_app/basic_web_app.go
```

### Scenario 2: High Availability Architecture
**File**: `scenario2_high_availability/high_availability.go`

A multi-AZ high availability architecture with load balancing:
- VPC with subnets across 3 availability zones
- Internet Gateway and NAT Gateway
- Application Load Balancer across public subnets
- Target Group for backend instances
- EC2 instances in private subnets behind ALB
- Launch Template for consistent instance configuration
- Auto Scaling Group with min=2, max=6
- Cost estimation including ALB and NAT Gateway

**Architecture**:
```
Internet
   ↓
Internet Gateway
   ↓
Application Load Balancer (Multi-AZ)
   ↓
Target Group
   ↓
VPC (10.0.0.0/16)
   ├── Public Subnets (3 AZs)
   │   └── NAT Gateway
   └── Private Subnets (3 AZs)
       └── Auto Scaling Group (2-6 instances)
           └── EC2 Instances (t3.small)
```

**Run**:
```bash
go run backend/pkg/usecases/scenario2_high_availability/high_availability.go
```

### Scenario 3: Scalable API Architecture
**File**: `scenario3_scalable_api/scalable_api.go`

A scalable API backend with auto-scaling and IAM integration:

### Scenario 5: Terraform Code Generation
**File**: `scenario5_terraform_codegen/terraform_codegen.go`

End-to-end pipeline from diagram IR JSON to generated Terraform files:
- Parse IR JSON into diagram graph
- Validate diagram (structure, schemas, relationships)
- Map to domain Architecture aggregate
- Validate domain rules/constraints (AWS networking defaults)
- Build domain graph + topologically sort resources
- Run Terraform engine to produce IaC files
- Write Terraform files to `./terraform_output/`

**Run**:
```bash
go run ./cmd/api/main.go -scenario=5
```

### Scenario 6: Terraform with Database Persistence
**File**: `scenario6_terraform_with_persistence/terraform_with_persistence.go`

Extends Scenario 5 by adding database persistence:
- All steps from Scenario 5
- Persist project, resources, containments, and dependencies to database
- Uses transactions for atomicity
- Creates demo user and IaC target if not exists

**Database Entities Persisted**:
- Project (with cloud provider, region, IaC tool)
- Resources (with configurations and positions)
- Resource Containments (parent-child relationships)
- Resource Dependencies (dependency relationships)

**Run**:
```bash
go run ./cmd/api/main.go -scenario=6
```

---

### Scenario 3: Scalable API Architecture (continued)
- VPC and networking infrastructure (similar to Scenario 2)
- IAM Role for EC2 instances
- IAM Instance Profile
- Launch Template with IAM role
- Application Load Balancer
- Target Group with health checks
- Auto Scaling Group with scaling policies (simulated)
- ELB-based health checks
- Scaling event simulation
- Cost estimation

**Architecture**:
```
Internet
   ↓
Internet Gateway
   ↓
Application Load Balancer (Multi-AZ)
   ↓
Target Group (Health Checks)
   ↓
VPC (10.0.0.0/16)
   ├── Public Subnets (3 AZs)
   │   └── NAT Gateway
   └── Private Subnets (3 AZs)
       └── Auto Scaling Group (2-10 instances)
           └── EC2 Instances (t3.medium)
               └── IAM Instance Profile
                   └── IAM Role
```

**Run**:
```bash
go run backend/pkg/usecases/scenario3_scalable_api/scalable_api.go
```

## Common Utilities

### Mock Helpers (`common/mock_helpers.go`)

Provides helper functions to generate mock data for all resource types:
- `MockIDGenerator`: Generates consistent mock IDs following AWS naming conventions
- `CreateMockVPC()`: Creates mock VPC domain model
- `CreateMockSubnet()`: Creates mock subnet domain model
- `CreateMockInternetGateway()`: Creates mock Internet Gateway
- `CreateMockSecurityGroup()`: Creates mock security group
- `CreateMockEC2Instance()`: Creates mock EC2 instance
- `CreateMockLoadBalancer()`: Creates mock load balancer
- `CreateMockTargetGroup()`: Creates mock target group
- `CreateMockLaunchTemplate()`: Creates mock launch template
- `CreateMockAutoScalingGroup()`: Creates mock Auto Scaling Group
- `CreateMockIAMRole()`: Creates mock IAM role
- `CreateMockIAMInstanceProfile()`: Creates mock instance profile
- `GetDefaultAvailabilityZones()`: Returns default AZs for a region

### Region Helper (`common/region_helper.go`)

Provides region selection and validation:
- `SupportedRegions`: List of commonly used AWS regions
- `ValidateRegion()`: Validates if a region is supported
- `SelectRegion()`: Selects a region (defaults to us-east-1)
- `FormatRegionName()`: Formats region code to human-readable name
- `DisplayRegions()`: Prints all available regions

## Output Models and Output Service Interfaces

The codebase now includes dedicated output models and output service interfaces that provide a cleaner separation between input configuration and output metadata.

### Output Models

Each resource type has dedicated output DTOs that focus on cloud-generated fields:

- **InstanceOutput**: Cloud-generated identifiers, runtime state, IP addresses, DNS names
- **LoadBalancerOutput**: DNS names, zone IDs, state information
- **VPCOutput**: State, creation time, owner information
- **SubnetOutput**: Available IP count, state, creation time
- And more for other resource types

### Output Service Interfaces

In addition to standard service interfaces, output service interfaces are available:

```go
// Standard service returns full domain model
instance, err := computeService.CreateInstance(ctx, instance)
// instance is *compute.Instance with all fields

// Output service returns focused output DTO
output, err := computeOutputService.CreateInstanceOutput(ctx, instance)
// output is *compute.InstanceOutput with cloud-generated fields
```

### Benefits

1. **Clear Separation**: Input configuration vs. output metadata
2. **Focused Models**: Output DTOs contain only relevant fields
3. **Type Safety**: Dedicated types prevent confusion
4. **Backward Compatible**: Original service methods still available

### Usage in Use Cases

Use cases can leverage output service interfaces when they only need cloud-generated fields:

```go
// Get instance output (focused on runtime state)
output, err := computeOutputService.GetInstanceOutput(ctx, instanceID)
fmt.Printf("Instance State: %s\n", output.State)
fmt.Printf("Created: %s\n", output.CreatedAt)

// Get load balancer output (focused on DNS and state)
lbOutput, err := computeOutputService.GetLoadBalancerOutput(ctx, lbARN)
fmt.Printf("DNS Name: %s\n", *lbOutput.DNSName)
fmt.Printf("Zone ID: %s\n", *lbOutput.ZoneID)
```

## Mock Data Pattern

All use cases follow the same pattern:
1. **No AWS SDK calls**: All resources are created using mock helpers
2. **Domain models**: Use domain models directly (not cloud-specific)
3. **Realistic data**: Mock IDs and ARNs follow AWS naming conventions
4. **Validation**: Domain models are validated before use
5. **Cost estimation**: Uses real pricing service with mock resource metadata
6. **Output models**: Can use output service interfaces for focused output DTOs

## Resource Dependencies

The use cases demonstrate proper resource dependencies:

```
VPC
  ├── Subnets (requires VPC)
  │   └── EC2 Instances (requires Subnet)
  ├── Internet Gateway (attached to VPC)
  ├── Route Tables (requires VPC)
  │   └── Routes (requires Route Table)
  ├── Security Groups (requires VPC)
  └── NAT Gateway (requires Subnet)
      └── Elastic IP (for NAT Gateway)

Load Balancer
  ├── Subnets (requires VPC)
  └── Security Groups (requires VPC)
      └── Target Group (requires VPC)
          └── EC2 Instances

Auto Scaling Group
  ├── Launch Template
  │   └── Security Groups
  ├── Subnets (requires VPC)
  └── Target Groups (optional, for ELB health checks)

IAM
  ├── IAM Role
  └── Instance Profile (requires IAM Role)
      └── EC2 Instance (uses Instance Profile)
```

## Cost Estimation

Each scenario includes cost estimation for:
- **EC2 Instances**: Based on instance type and runtime
- **Load Balancers**: Hourly base rate
- **NAT Gateways**: Hourly rate + data processing
- **Auto Scaling Groups**: Based on average capacity (min+max)/2

**Free Resources**:
- VPC
- Subnets
- Route Tables
- Security Groups
- Internet Gateway (when attached)
- IAM Roles and Policies

## Running Use Cases

### Prerequisites
- Go 1.19 or later
- No AWS credentials required (uses mock data)

### Run Individual Scenario

```bash
# Scenario 1: Basic Web Application
go run backend/pkg/usecases/scenario1_basic_web_app/basic_web_app.go

# Scenario 2: High Availability
go run backend/pkg/usecases/scenario2_high_availability/high_availability.go

# Scenario 3: Scalable API
go run backend/pkg/usecases/scenario3_scalable_api/scalable_api.go
```

### Build All Scenarios

```bash
go build ./pkg/usecases/...
```

## Architecture Diagrams

### Scenario 1: Basic Web Application
```
┌─────────────────────────────────────────┐
│           Internet Gateway              │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         VPC (10.0.0.0/16)                │
│                                          │
│  ┌──────────────────────────────────┐  │
│  │  Public Subnet 1 (10.0.1.0/24)   │  │
│  │  └── Web Server 1 (t3.micro)     │  │
│  └──────────────────────────────────┘  │
│                                          │
│  ┌──────────────────────────────────┐  │
│  │  Public Subnet 2 (10.0.2.0/24)   │  │
│  │  └── Web Server 2 (t3.micro)     │  │
│  └──────────────────────────────────┘  │
│                                          │
│  ┌──────────────────────────────────┐  │
│  │  Private Subnet 1 (10.0.10.0/24) │  │
│  │  └── App Server 1 (t3.small)     │  │
│  └──────────────────────────────────┘  │
│                                          │
│  ┌──────────────────────────────────┐  │
│  │  Private Subnet 2 (10.0.11.0/24) │  │
│  │  └── App Server 2 (t3.small)     │  │
│  └──────────────────────────────────┘  │
└──────────────────────────────────────────┘
```

### Scenario 2: High Availability
```
┌─────────────────────────────────────────┐
│           Internet Gateway              │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│    Application Load Balancer (ALB)     │
│         (Multi-AZ)                      │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            Target Group                 │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         VPC (10.0.0.0/16)               │
│                                          │
│  ┌──────────┐  ┌──────────┐  ┌────────┐│
│  │Public AZ1│  │Public AZ2│  │Public  ││
│  │  └─NAT   │  │          │  │  AZ3   ││
│  └──────────┘  └──────────┘  └────────┘│
│       │            │            │      │
│  ┌────▼────┐  ┌────▼────┐  ┌────▼────┐│
│  │Private  │  │Private  │  │Private  ││
│  │  AZ1    │  │  AZ2    │  │  AZ3    ││
│  │  └─ASG  │  │  └─ASG  │  │  └─ASG  ││
│  │(2-6 EC2)│  │(2-6 EC2)│  │(2-6 EC2)││
│  └─────────┘  └─────────┘  └─────────┘│
└─────────────────────────────────────────┘
```

## Future Extensions

Potential additions to these use cases:
- **Serverless architectures** (Lambda, API Gateway)
- **Containerized applications** (ECS, EKS)
- **Database integration** (RDS, DynamoDB)
- **CDN integration** (CloudFront)
- **Monitoring and logging** (CloudWatch, CloudTrail)
- **Disaster recovery** scenarios
- **Multi-region** architectures
- **Cost optimization** suggestions
- **Architecture validation** rules
- **Visual diagram** generation
- **Terraform export** functionality

## Notes

- All resources are **simulated** - no actual AWS resources are created
- Mock data follows **AWS naming conventions** for realism
- Domain models are **validated** before use
- Cost estimates are **approximate** and based on static pricing
- Use cases demonstrate **best practices** but may need customization for production

## Related Documentation

- [Domain Models](../../internal/domain/resource/) - Core domain models
- [AWS Adapters](../../internal/cloud/aws/adapters/) - AWS adapter implementations
- [Pricing Service](../../internal/cloud/aws/pricing/) - Cost estimation service
- [Example Runners](../../pkg/aws/) - Individual service examples
