# Domain Compute Layer

This package contains cloud-agnostic compute resource models following Domain-Driven Design (DDD) principles.

## Architecture Principles

- **Cloud-Agnostic**: No AWS, GCP, or Azure-specific code
- **Domain-First**: Represents business concepts, not implementation details
- **Validation**: Domain-level validation rules
- **Interfaces**: Polymorphic handling via `ComputeResource` interface

## Resources

### Instance (EC2, VM, etc.)

Represents a compute instance that can run applications and workloads.

**Key Attributes:**
- **Compute Configuration**: Instance type, AMI/image
- **Networking**: Subnet, security groups, IP addresses
- **Access & Permissions**: SSH keys, IAM roles
- **Storage**: Root volume configuration
- **State**: Instance lifecycle state (pending, running, stopped, etc.)

## Usage Example

```go
import domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"

// Create a domain instance
instance := &domaincompute.Instance{
    Name:         "web-server",
    Region:       "us-east-1",
    InstanceType: "t3.micro",
    AMI:          "ami-0c55b159cbfafe1f0",
    SubnetID:     "subnet-123",
    SecurityGroupIDs: []string{"sg-123"},
    RootVolume: &domaincompute.RootVolume{
        VolumeType: "gp3",
        VolumeSize:  20,
        DeleteOnTermination: true,
        Encrypted: false,
    },
}

// Validate
if err := instance.Validate(); err != nil {
    // handle error
}

// Use polymorphic interface
var resource domaincompute.ComputeResource = instance
fmt.Println(resource.GetName()) // "web-server"
```

## Interfaces

### ComputeResource

All compute resources implement the `ComputeResource` interface:

```go
type ComputeResource interface {
    GetID() string
    GetName() string
    GetSubnetID() string
    Validate() error
}
```

This allows polymorphic handling of different compute resource types.

### ComputeService

The `ComputeService` interface provides cloud-agnostic operations for compute resources:

```go
type ComputeService interface {
    CreateInstance(ctx context.Context, instance *Instance) (*Instance, error)
    GetInstance(ctx context.Context, id string) (*Instance, error)
    // ... other operations
}
```

This interface returns domain models with output fields (ID, ARN, State) populated after operations.

### ComputeOutputService

The `ComputeOutputService` interface provides operations that return dedicated output DTOs:

```go
type ComputeOutputService interface {
    CreateInstanceOutput(ctx context.Context, instance *Instance) (*InstanceOutput, error)
    GetInstanceOutput(ctx context.Context, id string) (*InstanceOutput, error)
    // ... other operations
}
```

This interface returns specialized output models that focus on cloud-generated fields and calculated outputs, providing a cleaner separation between input configuration and output metadata.

## Output Models

The package includes dedicated output DTOs for each resource type that focus on cloud-generated fields and runtime state:

### InstanceOutput

Represents the output data for a compute instance after creation/update:

```go
type InstanceOutput struct {
    ID               string
    ARN              *string
    Name             string
    Region           string
    AvailabilityZone *string
    InstanceType     string
    AMI              string
    SubnetID         string
    SecurityGroupIDs []string
    PrivateIP        *string
    PublicIP         *string
    PrivateDNS       *string
    PublicDNS        *string
    VPCID            *string
    KeyName          *string
    IAMInstanceProfile *string
    State            InstanceState
    CreatedAt        *time.Time
}
```

### Other Output Models

- `LoadBalancerOutput` - Load balancer output with DNS names and zone IDs
- `TargetGroupOutput` - Target group output with health check state
- `ListenerOutput` - Listener output with action configuration
- `LaunchTemplateOutput` - Launch template output with version information
- `AutoScalingGroupOutput` - Auto scaling group output with capacity and state
- `LambdaFunctionOutput` - Lambda function output with ARNs and metadata

### Converting Domain Models to Output

Helper functions are provided to convert domain models to output DTOs:

```go
instance := &domaincompute.Instance{
    ID: "i-123",
    ARN: stringPtr("arn:aws:..."),
    // ... other fields
}

output := domaincompute.ToInstanceOutput(instance)
// output is now an InstanceOutput with all relevant fields
```

## Validation Rules

### Instance Validation

- **Name**: Required, cannot be empty
- **Region**: Required, cannot be empty
- **InstanceType**: Required, cannot be empty
- **AMI**: Required, cannot be empty
- **SubnetID**: Required, cannot be empty
- **SecurityGroupIDs**: Required, must have at least one
- **RootVolume**: Optional, but if provided must be valid
  - Volume size must be > 0 and <= 16384 GB
  - Volume type must be valid (gp2, gp3, io1, io2, sc1, st1, standard)
  - IOPS validation for io1/io2/gp3
  - Throughput validation for gp3

## Mapping to Cloud Providers

Domain resources are mapped to cloud-specific implementations via mappers:
- `internal/cloud/aws/mapper/compute/` - AWS mappers
- `internal/cloud/gcp/mapper/compute/` - GCP mappers (future)
- `internal/cloud/azure/mapper/compute/` - Azure mappers (future)

## Dependencies

### Implemented Dependencies

- **Subnet** (`domain/resource/networking/subnet.go`) - Required for instance placement
- **Security Group** (`domain/resource/networking/security_group.go`) - Required for network access control
- **VPC** (`domain/resource/networking/vpc.go`) - Implicitly required via subnet

### Not Yet Implemented

The following resources are referenced by instances but not yet implemented:

- **EBS Volumes** (separate volumes) - only root volume is supported
- **Launch Templates** - Cannot use launch templates for instance configuration
- **IAM Roles/Instance Profiles** - Field exists but cannot validate IAM role existence
- **Key Pairs** - Field exists but cannot validate key pair existence
- **Placement Groups** - Cannot specify placement groups

See AWS EC2 model documentation for details on these dependencies.
