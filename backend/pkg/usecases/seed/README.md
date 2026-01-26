# Database Seeding

This package provides database seeding functionality to populate the database with realistic use case data based on the solution architect scenarios.

## What Gets Seeded

### Reference Data
- **Resource Categories**: Compute, Networking, Storage, Database, Security, Analytics, Application Integration
- **Resource Kinds**: VirtualMachine, Container, Function, Network, LoadBalancer, Database, Storage, Gateway
- **Resource Types**: AWS resource types (EC2, Lambda, VPC, Subnet, S3, RDS, etc.)
- **Dependency Types**: uses, depends_on, connects_to, references, contains
- **IAC Targets**: Terraform, Pulumi, CDK, CloudFormation

### Users
- 3 sample users with different email addresses

### Projects & Resources
Based on the use case scenarios:

1. **Basic Web Application** (User: alice@example.com)
   - VPC with CIDR 10.0.0.0/16
   - Public and private subnets
   - Internet Gateway
   - Security Groups
   - EC2 instances
   - Resource containment relationships

2. **High Availability Architecture** (User: bob@example.com)
   - VPC
   - Application Load Balancer
   - Auto Scaling Group
   - NAT Gateway

3. **Serverless Lambda + S3** (User: charlie@example.com)
   - Lambda function
   - S3 bucket
   - Dependency relationship (Lambda uses S3)

## Usage

### Run Seeding

```bash
# From backend directory
go run cmd/seed/main.go
```

Or build and run:

```bash
go build -o bin/seed cmd/seed/main.go
./bin/seed
```

### Prerequisites

1. Database must be running and accessible
2. Migrations must be run first:
   ```bash
   go run cmd/platform/main.go
   ```

### Idempotency

The seed function is idempotent:
- Reference data uses `FirstOrCreate` - won't duplicate if already exists
- Users check for existence before creating
- Projects are created fresh each time (you may want to clear them first if re-seeding)

### Clearing Data

To start fresh, you can:
1. Drop and recreate the database
2. Run migrations again
3. Run the seed command

## Seeding Approaches

This package provides two seeding approaches:

### 1. Basic Seeding (`seed.go`)
Uses `interface{}` types and JSON marshaling to create resources directly in the database. This approach is simpler but less type-safe.

**Function**: `SeedDatabase()`

### 2. Adapter-Based Seeding (`seed_with_adapters.go`)
Uses strong domain types and adapters to create resources through the service layer, then maps them to database resources. This approach:
- Uses strongly-typed domain models (e.g., `domaincompute.Instance`, `domainnetworking.VPC`)
- Leverages adapters to interact with virtual AWS services
- Provides better type safety and validation
- More closely mirrors how the application actually creates resources

**Function**: `SeedDatabaseWithAdapters()`

**Key Differences:**
- **Type Safety**: Adapter-based uses domain models instead of `map[string]interface{}`
- **Service Layer**: Resources are created through adapters using virtual services
- **Validation**: Domain and AWS validation is performed automatically
- **Realistic**: More closely matches production code paths

## Structure

```
pkg/usecases/seed/
├── seed.go              # Basic seeding with interface{} types
├── seed_with_adapters.go # Adapter-based seeding with strong types
└── README.md            # This file

cmd/seed/
└── main.go              # Entry point for seeding command
```

## Extending

To add more seed data:

### For Basic Seeding (`seed.go`):
1. Add new reference data in `seedReferenceData()`
2. Add new users in `seedUsers()`
3. Add new projects/scenarios in `seedScenarios()`

Make sure to:
- Use repositories for data access
- Handle errors properly
- Use JSON marshaling for resource configs
- Create proper relationships (containment, dependencies)

### For Adapter-Based Seeding (`seed_with_adapters.go`):
1. Create domain models (e.g., `domaincompute.Instance`, `domainnetworking.VPC`)
2. Use adapters to create resources through virtual services
3. Convert domain resources to database resources using `domainResourceToJSON()`
4. Save to database using repositories

**Example:**
```go
// Create domain model
vpcDomain := &domainnetworking.VPC{
    Name:   "my-vpc",
    Region: "us-east-1",
    CIDR:   "10.0.0.0/16",
}

// Create via adapter (uses virtual service)
createdVPC, err := networkingAdapter.CreateVPC(ctx, vpcDomain)

// Convert to database resource
vpcConfigJSON, _ := domainResourceToJSON(createdVPC)
vpcResource := &models.Resource{
    ProjectID:      project.ID,
    ResourceTypeID: vpcType.ID,
    Name:           createdVPC.Name,
    Config:         vpcConfigJSON,
}

// Save to database
resourceRepo.Create(ctx, vpcResource)
```
