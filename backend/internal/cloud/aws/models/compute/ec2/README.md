# AWS EC2 Models

This package contains AWS-specific EC2 instance models that represent cloud provider implementations. These models capture AWS-specific details, constraints, and features.

## Purpose

The EC2 models layer serves as the **AWS-specific representation** of compute instances:

- **Provider-Specific Details**: Captures AWS-specific fields, constraints, and behaviors
- **Validation Rules**: Implements AWS-specific validation (AMI format, instance types, etc.)
- **API Compatibility**: Models align with AWS API structures and Terraform resource schemas
- **Configuration Support**: JSON-serializable for configuration files and API responses
- **Tag Management**: Standardized AWS tagging support

## Models

### Instance (`instance.go`)

AWS EC2 instance input model with all configuration options:

**Required Fields:**
- `Name`: Instance name
- `AMI`: AMI ID (must start with "ami-")
- `InstanceType`: Instance type (e.g., "t3.micro", "m5.large")
- `SubnetID`: Subnet ID (must start with "subnet-")
- `VpcSecurityGroupIds`: List of security group IDs (must start with "sg-")

**Optional Fields:**
- `AssociatePublicIPAddress`: Whether to associate a public IP
- `KeyName`: SSH key pair name
- `IAMInstanceProfile`: IAM instance profile name
- `UserData`: User data script (max 12KB raw, 16KB base64 encoded)
- `RootBlockDevice`: Root volume configuration
- `Tags`: AWS tags

### RootBlockDevice (`block_device.go`)

EBS root volume configuration:

- `VolumeType`: Volume type (gp2, gp3, io1, io2, sc1, st1, standard)
- `VolumeSize`: Size in GB (1-16384)
- `DeleteOnTermination`: Whether to delete on instance termination
- `Encrypted`: Whether volume is encrypted
- `IOPS`: Optional IOPS (required for io1/io2, optional for gp3)
- `Throughput`: Optional throughput in MB/s (gp3 only)

### InstanceOutput (`outputs/instance_output.go`)

Post-creation EC2 instance data including AWS-generated identifiers:

- `ID`: Instance ID (e.g., "i-1234567890abcdef0")
- `ARN`: Amazon Resource Name
- `State`: Instance state (pending, running, stopped, etc.)
- `PublicIP`, `PrivateIP`: IP addresses
- `PublicDNS`, `PrivateDNS`: DNS hostnames
- `AvailabilityZone`: AZ where instance was launched
- `CreationTime`: When instance was created

## Validation Rules

### Instance Validation

1. **Required Fields**: Name, AMI, InstanceType, SubnetID, at least one SecurityGroupID
2. **AMI Format**: Must start with "ami-" and be 12-21 characters
3. **Instance Type Format**: Must follow pattern `{family}.{size}` (e.g., "t3.micro")
4. **Subnet ID Format**: Must start with "subnet-"
5. **Security Group ID Format**: Must start with "sg-"
6. **UserData Size**: Cannot exceed 12KB raw (16KB base64 encoded)
7. **RootBlockDevice**: Validated separately (see RootBlockDevice validation)

### RootBlockDevice Validation

1. **Volume Size**: Must be > 0 and <= 16384 GB
2. **Volume Type**: Must be one of: gp2, gp3, io1, io2, sc1, st1, standard
3. **IOPS**:
   - gp3: 3000-16000
   - io1/io2: 100-64000
   - Other types: Cannot specify IOPS
4. **Throughput**: Only for gp3, must be 125-1000 MB/s

## Input vs Output Models

### Input Model (`Instance`)

Used when **creating** or **updating** an instance:
- Contains user-provided configuration
- No AWS-generated fields (ID, ARN, etc.)
- Used in `CreateInstance()` and `UpdateInstance()` operations

### Output Model (`InstanceOutput`)

Returned **after** instance creation/retrieval:
- Contains AWS-generated identifiers (ID, ARN)
- Includes runtime state (State, IP addresses, DNS names)
- Includes timestamps (CreationTime)
- Used in service responses

## Dependencies and Missing Resources

### Implemented Dependencies

- **Subnet** - Required for `SubnetID`
- **Security Group** - Required for `VpcSecurityGroupIds`
- **VPC** - Implicitly required via subnet

### Not Yet Implemented (Dependencies)

The following resources are referenced by EC2 but not yet implemented:

#### 1. EBS Volumes (Separate Volumes)
**Status**: Not implemented  
**Use Case**: Attaching additional EBS volumes beyond the root volume  
**Impact**: EC2 can only configure root volume. Additional volumes cannot be attached.  
**Future Implementation**: 
- Create `domain/resource/storage/ebs_volume.go`
- Create `cloud/aws/models/storage/ebs_volume.go`
- Add `AttachVolume()` and `DetachVolume()` methods to compute service

#### 2. Launch Templates
**Status**: Not implemented  
**Use Case**: Reusable EC2 configurations for Auto Scaling Groups, Spot Instances  
**Impact**: EC2 instances must be configured individually. Cannot use launch templates.  
**Future Implementation**:
- Create `domain/resource/compute/launch_template.go`
- Create `cloud/aws/models/compute/ec2/launch_template.go`
- Add `LaunchTemplateID` field to Instance model

#### 3. IAM Roles/Instance Profiles
**Status**: Not implemented  
**Use Case**: IAM instance profiles for EC2 instances to access AWS services  
**Impact**: `IAMInstanceProfile` field accepts string but cannot validate IAM role existence.  
**Future Implementation**:
- Create `domain/resource/iam/instance_profile.go`
- Create `cloud/aws/models/iam/instance_profile.go`
- Add validation to check IAM role exists

#### 4. Key Pairs
**Status**: Not implemented  
**Use Case**: SSH key pairs for EC2 instance access  
**Impact**: `KeyName` field accepts string but cannot validate key pair existence.  
**Future Implementation**:
- Create `domain/resource/compute/key_pair.go`
- Create `cloud/aws/models/compute/ec2/key_pair.go`
- Add validation to check key pair exists

#### 5. Placement Groups
**Status**: Not implemented  
**Use Case**: Control EC2 instance placement for performance/clustering  
**Impact**: Cannot specify placement groups for EC2 instances.  
**Future Implementation**:
- Add `PlacementGroup` field to Instance model
- Create placement group models if needed

## Example Usage

### Creating an Instance

```go
import awsec2 "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/ec2"

instance := &awsec2.Instance{
    Name:                "web-server",
    AMI:                 "ami-0c55b159cbfafe1f0",
    InstanceType:        "t3.micro",
    SubnetID:            "subnet-12345678",
    VpcSecurityGroupIds: []string{"sg-12345678"},
    AssociatePublicIPAddress: boolPtr(true),
    RootBlockDevice: &awsec2.RootBlockDevice{
        VolumeType:          "gp3",
        VolumeSize:          20,
        DeleteOnTermination: true,
        Encrypted:           false,
    },
    Tags: []configs.Tag{{Key: "Name", Value: "web-server"}},
}

if err := instance.Validate(); err != nil {
    // handle error
}
```

## Related Documentation

- [Domain Compute Models](../../../../domain/resource/compute/README.md) - Domain layer compute models
- [AWS Models README](../../README.md) - Overview of AWS models
- [EC2 Mapper](../../../mapper/compute/README.md) - Domain ↔ AWS conversion
