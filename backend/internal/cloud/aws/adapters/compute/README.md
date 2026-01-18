# AWS Compute Adapter

This package implements the Adapter pattern to bridge the domain compute service interface with AWS-specific compute service implementations.

## Purpose

The adapter layer:
- **Translates** domain models to AWS models and vice versa
- **Validates** resources at both domain and AWS levels
- **Handles Errors** with clear, contextual error messages
- **Maintains Separation** between domain and cloud-specific code

## Architecture

```
Domain Layer (Cloud-Agnostic)
    ↓
Adapter Layer (Translation & Validation)
    ↓
Mapper Layer (Model Conversion)
    ↓
AWS Service Layer (Provider-Specific)
```

## Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                              │
│  domaincompute.Instance (Cloud-Agnostic)                    │
└──────────────────────┬──────────────────────────────────────┘
                        │
                        │ CreateInstance()
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                  Adapter Layer                               │
│  AWSComputeAdapter                                          │
│  1. Validate domain model                                    │
│  2. Convert to AWS model (via mapper)                        │
│  3. Validate AWS model                                      │
│  4. Call AWS service                                        │
│  5. Convert output to domain (via mapper)                   │
└──────────────────────┬──────────────────────────────────────┘
                        │
                        │ AWS Service Call
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                  AWS Service Layer                           │
│  AWSComputeService.CreateInstance()                          │
│  Returns: awsec2outputs.InstanceOutput                       │
└──────────────────────┬──────────────────────────────────────┘
                        │
                        │ Output Model
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                  Adapter Layer                               │
│  Converts output to domain model                             │
└──────────────────────┬──────────────────────────────────────┘
                        │
                        │ Domain Instance (with ID/ARN)
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                              │
│  domaincompute.Instance (with populated ID/ARN)            │
└─────────────────────────────────────────────────────────────┘
```

## Step-by-Step Example: Creating an Instance

### 1. Domain Model (Input)

```go
domainInstance := &domaincompute.Instance{
    Name:         "web-server",
    Region:       "us-east-1",
    InstanceType: "t3.micro",
    AMI:          "ami-0c55b159cbfafe1f0",
    SubnetID:     "subnet-123",
    SecurityGroupIDs: []string{"sg-123"},
    RootVolume: &domaincompute.RootVolume{
        VolumeType: "gp3",
        VolumeSize: 20,
    },
}
```

### 2. Adapter Validation

```go
adapter := NewAWSComputeAdapter(awsService)
instance, err := adapter.CreateInstance(ctx, domainInstance)
```

**Adapter performs:**
- ✅ Domain validation (`domainInstance.Validate()`)
- ✅ Converts to AWS model (`FromDomainInstance()`)
- ✅ AWS validation (`awsInstance.Validate()`)
- ✅ Calls AWS service (`awsService.CreateInstance()`)
- ✅ Converts output to domain (`ToDomainInstanceFromOutput()`)

### 3. AWS Service Call

The AWS service receives:
```go
awsInstance := &awsec2.Instance{
    Name:                "web-server",
    AMI:                 "ami-0c55b159cbfafe1f0",
    InstanceType:        "t3.micro",
    SubnetID:            "subnet-123",
    VpcSecurityGroupIds: []string{"sg-123"},
    RootBlockDevice: &awsec2.RootBlockDevice{
        VolumeType: "gp3",
        VolumeSize: 20,
    },
}
```

### 4. Output Model

AWS service returns:
```go
awsOutput := &awsec2outputs.InstanceOutput{
    ID:     "i-1234567890abcdef0",
    ARN:    "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
    Name:   "web-server",
    State:  "running",
    PublicIP: stringPtr("54.123.45.67"),
    PrivateIP: "10.0.1.100",
    // ... other fields
}
```

### 5. Domain Model (Output)

Adapter converts and returns:
```go
domainInstance := &domaincompute.Instance{
    ID:     "i-1234567890abcdef0",
    ARN:    stringPtr("arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0"),
    Name:   "web-server",
    State:  domaincompute.InstanceStateRunning,
    PublicIP: stringPtr("54.123.45.67"),
    PrivateIP: stringPtr("10.0.1.100"),
    // ... other fields
}
```

## Available Operations

### Instance Operations

- `CreateInstance()` - Create a new EC2 instance
- `GetInstance()` - Retrieve instance by ID
- `UpdateInstance()` - Update instance configuration
- `DeleteInstance()` - Delete an instance
- `ListInstances()` - List instances with filters

### Instance Lifecycle Operations

- `StartInstance()` - Start a stopped instance
- `StopInstance()` - Stop a running instance
- `RebootInstance()` - Reboot an instance

## Error Handling

The adapter provides contextual error messages:

```go
// Domain validation error
err := adapter.CreateInstance(ctx, invalidInstance)
// Error: "domain validation failed: instance name is required"

// AWS validation error
err := adapter.CreateInstance(ctx, instanceWithInvalidAMI)
// Error: "aws validation failed: ami must start with 'ami-'"

// AWS service error
err := adapter.CreateInstance(ctx, validInstance)
// Error: "aws service error: subnet not found"
```

## Usage Example

```go
import (
    "context"
    domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
    awsadapter "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/adapters/compute"
    awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
)

// Create AWS service (mock or real implementation)
awsService := awsservice.NewAWSComputeService()

// Create adapter
adapter := awsadapter.NewAWSComputeAdapter(awsService)

// Create domain instance
instance := &domaincompute.Instance{
    Name:         "web-server",
    Region:       "us-east-1",
    InstanceType: "t3.micro",
    AMI:          "ami-0c55b159cbfafe1f0",
    SubnetID:     "subnet-123",
    SecurityGroupIDs: []string{"sg-123"},
}

// Create instance via adapter
createdInstance, err := adapter.CreateInstance(ctx, instance)
if err != nil {
    // handle error
}

// Instance now has ID and ARN populated
fmt.Printf("Created instance: %s\n", createdInstance.ID)
fmt.Printf("ARN: %s\n", *createdInstance.ARN)
```

## Related Documentation

- [Domain Compute Service](../../../domain/resource/compute/service.go) - Domain service interface
- [AWS Compute Service](../../services/compute/interfaces.go) - AWS service interface
- [Compute Mapper](../../mapper/compute/README.md) - Model conversion
- [Networking Adapter](../networking/README.md) - Similar pattern for networking
