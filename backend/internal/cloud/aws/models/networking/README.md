# AWS Networking Models

This package contains AWS-specific networking resource models that implement cloud provider-specific details.

## Structure

```
networking/
├── vpc.go                 # AWS VPC model
├── subnet.go              # AWS Subnet model
├── internet_gateway.go    # AWS Internet Gateway model
├── route_table.go         # AWS Route Table model
├── security_group.go      # AWS Security Group model
├── nat_gateway.go         # AWS NAT Gateway model
└── tests/
    └── vpc_test.go        # VPC validation tests
```

## Design Principles

1. **AWS-Specific**: Contains AWS-specific fields and validation rules
2. **Validation**: Each model has `Validate()` method for AWS-specific rules
3. **JSON Serialization**: Models are JSON-serializable for API/configuration
4. **Tags**: All resources support AWS tags via `configs.Tag`

## Mapping to Domain

AWS models are mapped to domain models via:
- `internal/cloud/aws/mapper/networking/` - Conversion functions

### Example Flow

```
Domain VPC (cloud-agnostic)
    ↓ [FromDomainVPC]
AWS VPC (AWS-specific)
    ↓ [Validate]
AWS API/Configuration
```

## Usage Example

```go
import awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"

// Create AWS VPC
vpc := &awsnetworking.VPC{
    Name:               "my-vpc",
    Region:             "us-east-1",
    CIDR:               "10.0.0.0/16",
    EnableDNSSupport:   true,
    EnableDNSHostnames: true,
    InstanceTenancy:    "default",
    Tags: []configs.Tag{
        {Key: "Name", Value: "my-vpc"},
        {Key: "Environment", Value: "production"},
    },
}

// Validate AWS-specific rules
if err := vpc.Validate(); err != nil {
    // handle error
}
```

## AWS-Specific Features

### VPC
- Instance tenancy (default, dedicated)
- DNS support flags
- Region-scoped

### Subnet
- Availability zone required
- Map public IP on launch
- VPC-scoped

### Internet Gateway
- VPC attachment required
- Single attachment per VPC

### Route Table
- Multiple route targets supported
- Subnet associations
- VPC-scoped

### Security Group
- Ingress/egress rules
- Source security group references
- Description required (max 255 chars)

### NAT Gateway
- Elastic IP allocation required
- Subnet-scoped (must be public subnet)
