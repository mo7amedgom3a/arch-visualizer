# AWS Networking Mapper

This package provides bidirectional mapping between domain networking resources and AWS-specific networking models.

## Purpose

The mapper layer enables:
- **Domain → AWS**: Convert cloud-agnostic domain models to AWS-specific models
- **AWS → Domain**: Convert AWS models back to domain models
- **Separation of Concerns**: Keeps domain and cloud layers decoupled

## Mappers

### VPC Mapper
- `ToDomainVPC()` - AWS VPC → Domain VPC
- `FromDomainVPC()` - Domain VPC → AWS VPC

### Subnet Mapper
- `ToDomainSubnet()` - AWS Subnet → Domain Subnet
- `FromDomainSubnet()` - Domain Subnet → AWS Subnet

### Internet Gateway Mapper
- `ToDomainInternetGateway()` - AWS IGW → Domain IGW
- `FromDomainInternetGateway()` - Domain IGW → AWS IGW

### Route Table Mapper
- `ToDomainRouteTable()` - AWS Route Table → Domain Route Table
- `FromDomainRouteTable()` - Domain Route Table → AWS Route Table

### Security Group Mapper
- `ToDomainSecurityGroup()` - AWS SG → Domain SG
- `FromDomainSecurityGroup()` - Domain SG → AWS SG

### NAT Gateway Mapper
- `ToDomainNATGateway()` - AWS NAT → Domain NAT
- `FromDomainNATGateway()` - Domain NAT → AWS NAT

## Usage Example

```go
import (
    domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
    awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
    awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/networking"
)

// Domain VPC (cloud-agnostic)
domainVPC := &domainnetworking.VPC{
    Name:   "my-vpc",
    Region: "us-east-1",
    CIDR:   "10.0.0.0/16",
    EnableDNS: true,
}

// Convert to AWS VPC
awsVPC := awsmapper.FromDomainVPC(domainVPC)

// Convert back to domain
convertedDomainVPC := awsmapper.ToDomainVPC(awsVPC)
```

## Mapping Rules

### VPC
- Domain `EnableDNS` → AWS `EnableDNSSupport`
- Domain `EnableDNSHostnames` → AWS `EnableDNSHostnames`
- AWS adds default `InstanceTenancy` and `Tags`

### Subnet
- Domain `IsPublic` → AWS `MapPublicIPOnLaunch`
- Domain `AvailabilityZone` → AWS `AvailabilityZone` (required)
- AWS adds default `Tags`

### Route Table
- Domain route `TargetType` → AWS route target fields (gateway_id, nat_gateway_id, etc.)
- AWS supports multiple target types, domain uses single target per route

### Security Group
- Domain `SourceGroupIDs` (list) → AWS `SourceSecurityGroupID` (single)
- AWS limitation: only one source security group per rule

## Design Notes

- **Null Safety**: All mappers handle nil inputs gracefully
- **Default Values**: AWS mappers add provider-specific defaults (tags, tenancy)
- **Field Mapping**: Some fields map 1:1, others require transformation
- **Validation**: Mappers don't validate - validation happens in model `Validate()` methods
