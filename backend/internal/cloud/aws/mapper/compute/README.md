# AWS Compute Mapper

This package provides bidirectional conversion between domain compute models and AWS-specific compute models.

## Purpose

The mapper layer translates between:
- **Domain Models**: Cloud-agnostic, business-focused representations
- **AWS Models**: AWS-specific implementations with provider details

## Mapper Functions

### Instance Mapper (`instance_mapper.go`)

#### `FromDomainInstance(domainInstance *domaincompute.Instance) *awsec2.Instance`

Converts domain instance to AWS EC2 input model.

**Key Conversions:**
- Domain `SecurityGroupIDs` → AWS `VpcSecurityGroupIds`
- Domain `PublicIP` presence → AWS `AssociatePublicIPAddress` boolean
- Domain `RootVolume` → AWS `RootBlockDevice`
- Adds default AWS tags

**Example:**
```go
domainInstance := &domaincompute.Instance{
    Name: "web-server",
    AMI: "ami-123",
    InstanceType: "t3.micro",
    SubnetID: "subnet-123",
    SecurityGroupIDs: []string{"sg-123"},
}

awsInstance := awsmapper.FromDomainInstance(domainInstance)
// awsInstance is ready for AWS service call
```

#### `ToDomainInstance(awsInstance *awsec2.Instance) *domaincompute.Instance`

Converts AWS EC2 input model to domain instance (backward compatibility).

**Note**: This function does not populate ID/ARN as those are only available in output models.

#### `ToDomainInstanceFromOutput(output *awsec2outputs.InstanceOutput) *domaincompute.Instance`

Converts AWS EC2 output model to domain instance with ID and ARN populated.

**Key Conversions:**
- AWS `ID` → Domain `ID`
- AWS `ARN` → Domain `ARN` (as pointer, nil if empty)
- AWS `State` → Domain `State` (converted to InstanceState enum)
- AWS `AvailabilityZone` → Domain `AvailabilityZone` (as pointer)
- AWS `PublicIP`, `PrivateIP` → Domain IP addresses
- AWS `SecurityGroupIDs` → Domain `SecurityGroupIDs`

**Example:**
```go
awsOutput := &awsec2outputs.InstanceOutput{
    ID: "i-1234567890abcdef0",
    ARN: "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
    Name: "web-server",
    State: "running",
    // ... other fields
}

domainInstance := awsmapper.ToDomainInstanceFromOutput(awsOutput)
// domainInstance has ID and ARN populated
```

## Input vs Output Models

### When to Use Each Function

1. **Creating Instance**: `FromDomainInstance()` → AWS service → `ToDomainInstanceFromOutput()`
2. **Getting Instance**: AWS service → `ToDomainInstanceFromOutput()`
3. **Updating Instance**: `FromDomainInstance()` → AWS service → `ToDomainInstanceFromOutput()`

### Mapping Flow

```
Domain Instance (Input)
    ↓ FromDomainInstance()
AWS Instance (Input Model)
    ↓ AWS Service CreateInstance()
AWS InstanceOutput (Output Model)
    ↓ ToDomainInstanceFromOutput()
Domain Instance (with ID/ARN)
```

## Field Mapping Details

### Root Volume / Block Device

| Domain Field | AWS Field | Notes |
|-------------|-----------|-------|
| `RootVolume.VolumeType` | `RootBlockDevice.VolumeType` | Direct mapping |
| `RootVolume.VolumeSize` | `RootBlockDevice.VolumeSize` | Direct mapping |
| `RootVolume.DeleteOnTermination` | `RootBlockDevice.DeleteOnTermination` | Direct mapping |
| `RootVolume.Encrypted` | `RootBlockDevice.Encrypted` | Direct mapping |
| `RootVolume.IOPS` | `RootBlockDevice.IOPS` | Pointer mapping |
| `RootVolume.Throughput` | `RootBlockDevice.Throughput` | Pointer mapping |

### Networking

| Domain Field | AWS Field | Notes |
|-------------|-----------|-------|
| `SubnetID` | `SubnetID` | Direct mapping |
| `SecurityGroupIDs` | `VpcSecurityGroupIds` | Array mapping |
| `PublicIP` (presence) | `AssociatePublicIPAddress` | Boolean derived from PublicIP field |
| `PrivateIP` | `PrivateIP` | From output only |

### Access & Permissions

| Domain Field | AWS Field | Notes |
|-------------|-----------|-------|
| `KeyName` | `KeyName` | Pointer mapping |
| `IAMInstanceProfile` | `IAMInstanceProfile` | Pointer mapping |

## Related Documentation

- [Domain Compute Models](../../../domain/resource/compute/README.md) - Domain layer models
- [AWS EC2 Models](../../models/compute/ec2/README.md) - AWS EC2 models
- [Networking Mapper](../networking/README.md) - Similar pattern for networking resources
