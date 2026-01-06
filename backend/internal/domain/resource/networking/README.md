# Domain Networking Layer

This package contains cloud-agnostic networking resource models following Domain-Driven Design (DDD) principles.

## Architecture Principles

- **Cloud-Agnostic**: No AWS, GCP, or Azure-specific code
- **Domain-First**: Represents business concepts, not implementation details
- **Validation**: Domain-level validation rules
- **Interfaces**: Polymorphic handling via `NetworkResource` interface

## Resources

### VPC (Virtual Private Cloud)
- Core networking container
- Contains subnets, gateways, route tables
- Region-scoped

### Subnet
- Network segment within a VPC
- Can be public or private
- Availability zone scoped

### Internet Gateway
- Provides internet access to VPC
- Attached to VPC

### Route Table
- Defines routing rules
- Associated with subnets
- Contains routes to gateways, NAT, peering connections

### Security Group
- Stateful firewall rules
- Ingress and egress rules
- VPC-scoped

### NAT Gateway
- Provides outbound internet access for private subnets
- Subnet-scoped

## Usage Example

```go
import domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"

// Create a domain VPC
vpc := &domainnetworking.VPC{
    Name:   "my-vpc",
    Region: "us-east-1",
    CIDR:   "10.0.0.0/16",
    EnableDNS: true,
    EnableDNSHostnames: true,
}

// Validate
if err := vpc.Validate(); err != nil {
    // handle error
}

// Use polymorphic interface
var resource domainnetworking.NetworkResource = vpc
fmt.Println(resource.GetName()) // "my-vpc"
```

## Mapping to Cloud Providers

Domain resources are mapped to cloud-specific implementations via mappers:
- `internal/cloud/aws/mapper/networking/` - AWS mappers
- `internal/cloud/gcp/mapper/networking/` - GCP mappers (future)
- `internal/cloud/azure/mapper/networking/` - Azure mappers (future)
