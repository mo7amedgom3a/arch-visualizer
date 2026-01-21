# AWS Storage Adapter

This package implements the **Adapter Pattern** to bridge the domain layer and AWS-specific storage implementations.

## Purpose

The adapter layer:
- **Implements domain interfaces** using AWS-specific services
- **Converts between domain and AWS models** using mappers
- **Handles validation** at both domain and AWS levels
- **Provides error translation** from AWS-specific to domain-level errors
- **Enables provider swapping** without changing domain code

## Architecture

```
Domain Layer (Cloud-Agnostic)
    ↓
StorageService Interface
    ↓
AWSStorageAdapter (This Package)
    ↓
AWSStorageService (AWS-Specific)
    ↓
AWS Models & API
```

## Components

### AWSStorageAdapter

The main adapter that implements `domainstorage.StorageService` interface.

**Responsibilities:**
- Accept domain models as input (no ID/ARN initially)
- Validate domain models
- Convert to AWS input models using mappers
- Validate AWS input models
- Call AWS service (returns output models with ID/ARN)
- Convert AWS output models back to domain models (with ID/ARN populated)
- Wrap errors with context

### Factory Pattern

`AWSStorageAdapterFactory` provides a factory pattern for creating adapters:

```go
factory := NewAWSStorageAdapterFactory(awsService)
adapter := factory.CreateStorageAdapter()
```

## Usage Example

```go
import (
    domainstorage "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/storage"
    awsservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/storage"
    awsadapter "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/adapters/storage"
)

// Create AWS service (implementation)
awsService := &awsStorageServiceImpl{...}

// Create adapter
adapter := awsadapter.NewAWSStorageAdapter(awsService)

// Use domain interface (cloud-agnostic)
var storageService domainstorage.StorageService = adapter

// Create EBS Volume using domain model (input - no ID/ARN)
volume := &domainstorage.EBSVolume{
    Name:             "my-volume",
    Region:           "us-east-1",
    AvailabilityZone: "us-east-1a",
    Size:             40,
    Type:             "gp3",
    Encrypted:        false,
    // ID: "" (empty before creation)
    // ARN: nil (nil before creation)
}

// Create volume - adapter handles input/output conversion
createdVolume, err := storageService.CreateEBSVolume(ctx, volume)
if err != nil {
    // handle error
}

// Created volume now has ID and ARN populated!
fmt.Printf("Volume ID: %s\n", createdVolume.ID)        // "vol-0a1b2c3d4e5f6g7h8"
fmt.Printf("Volume ARN: %s\n", *createdVolume.ARN)    // "arn:aws:ec2:us-east-1:..."

// Get volume by ID (also returns domain model with ID/ARN)
retrievedVolume, err := storageService.GetEBSVolume(ctx, createdVolume.ID)
// retrievedVolume has ID and ARN populated from AWS output model

// Create S3 Bucket using domain model (input - no ID/ARN)
bucket := &domainstorage.S3Bucket{
    Name:         "my-bucket",
    Region:       "us-east-1",
    ForceDestroy: false,
    Tags: map[string]string{
        "Environment": "prod",
    },
    // ID: "" (empty before creation)
    // ARN: nil (nil before creation)
}

// Create bucket - adapter handles input/output conversion
createdBucket, err := storageService.CreateS3Bucket(ctx, bucket)
if err != nil {
    // handle error
}

// Created bucket now has ID, ARN, and domain names populated!
fmt.Printf("Bucket ID: %s\n", createdBucket.ID)                    // "my-bucket"
fmt.Printf("Bucket ARN: %s\n", *createdBucket.ARN)                 // "arn:aws:s3:::my-bucket"
fmt.Printf("Domain Name: %s\n", *createdBucket.BucketDomainName)    // "my-bucket.s3.amazonaws.com"
fmt.Printf("Regional Domain: %s\n", *createdBucket.BucketRegionalDomainName) // "my-bucket.s3.us-east-1.amazonaws.com"
```

## Flow Diagram

### Complete Flow with Input and Output Models

```
1. Domain Layer calls: CreateEBSVolume(domainVolume)
   Input: Domain EBSVolume (no ID/ARN)
   ↓
2. Adapter validates: domainVolume.Validate()
   ↓
3. Adapter converts: FromDomainEBSVolume(domainVolume) → awsVolumeInput
   Output: AWS Volume Input Model (no ID/ARN)
   ↓
4. Adapter validates: awsVolumeInput.Validate()
   ↓
5. Adapter calls: awsService.CreateEBSVolume(awsVolumeInput)
   ↓
6. AWS Service returns: awsVolumeOutput
   Output: AWS Volume Output Model (with ID, ARN, State, CreateTime)
   ↓
7. Adapter converts: ToDomainEBSVolumeFromOutput(awsVolumeOutput) → domainVolume
   Output: Domain EBSVolume (with ID and ARN populated)
   ↓
8. Domain Layer receives: domainVolume
   Result: Domain EBSVolume with AWS-generated ID and ARN
```

### Input vs Output Models

**Input Models** (for creation/update):
- Located in: `internal/cloud/aws/models/storage/ebs/`
- Contains: Configuration fields only (Name, Size, Type, AvailabilityZone, etc.)
- No AWS identifiers: No `ID` or `ARN` fields
- Used when: Creating or updating resources

**Output Models** (from AWS responses):
- Located in: `internal/cloud/aws/models/storage/ebs/outputs/`
- Contains: Configuration + AWS-generated metadata
- AWS identifiers: `ID`, `ARN`, `State`, `CreateTime`, `AttachedTo`, etc.
- Used when: Receiving responses from AWS services

### Example: EBS Volume Creation Flow

```go
// Step 1: Domain input (no ID/ARN)
domainVolume := &domainstorage.EBSVolume{
    Name:             "my-volume",
    Region:           "us-east-1",
    AvailabilityZone: "us-east-1a",
    Size:             40,
    Type:             "gp3",
    // ID: "" (empty)
    // ARN: nil
}

// Step 2: Adapter converts to AWS input
awsVolumeInput := awsmapper.FromDomainEBSVolume(domainVolume)
// awsVolumeInput has no ID/ARN fields

// Step 3: AWS service creates and returns output
awsVolumeOutput := &awsebsoutputs.VolumeOutput{
    ID:     "vol-0a1b2c3d4e5f6g7h8",  // AWS-generated
    ARN:    "arn:aws:ec2:...",         // AWS-generated
    State:  "available",               // AWS metadata
    // ... configuration fields
}

// Step 4: Adapter converts output to domain (with ID/ARN)
createdVolume := awsmapper.ToDomainEBSVolumeFromOutput(awsVolumeOutput)
// createdVolume.ID = "vol-0a1b2c3d4e5f6g7h8"
// createdVolume.ARN = "arn:aws:ec2:..."
```

### Example: S3 Bucket Creation Flow

```go
// Step 1: Domain input (no ID/ARN)
domainBucket := &domainstorage.S3Bucket{
    Name:         "my-bucket",
    Region:       "us-east-1",
    ForceDestroy: false,
    Tags: map[string]string{
        "Environment": "prod",
    },
    // ID: "" (empty)
    // ARN: nil
}

// Step 2: Adapter converts to AWS input
awsBucketInput := awsmapper.FromDomainS3Bucket(domainBucket)
// awsBucketInput has no ID/ARN fields

// Step 3: AWS service creates and returns output
awsBucketOutput := &awss3outputs.BucketOutput{
    ID:                       "my-bucket",                          // AWS-generated (bucket name)
    ARN:                      "arn:aws:s3:::my-bucket",            // AWS-generated
    BucketDomainName:         "my-bucket.s3.amazonaws.com",       // AWS-generated
    BucketRegionalDomainName: "my-bucket.s3.us-east-1.amazonaws.com", // AWS-generated
    Region:                   "us-east-1",
    // ... configuration fields
}

// Step 4: Adapter converts output to domain (with ID/ARN/domain names)
createdBucket := awsmapper.ToDomainS3BucketFromOutput(awsBucketOutput)
// createdBucket.ID = "my-bucket"
// createdBucket.ARN = "arn:aws:s3:::my-bucket"
// createdBucket.BucketDomainName = "my-bucket.s3.amazonaws.com"
// createdBucket.BucketRegionalDomainName = "my-bucket.s3.us-east-1.amazonaws.com"
```

### S3 Bucket Input/Output Models

**Input Models** (for creation/update):
- Located in: `internal/cloud/aws/models/storage/s3/`
- Contains: Configuration fields only (Bucket/BucketPrefix, ForceDestroy, Tags)
- No AWS identifiers: No `ID` or `ARN` fields
- Used when: Creating or updating buckets

**Output Models** (from AWS responses):
- Located in: `internal/cloud/aws/models/storage/s3/outputs/`
- Contains: Configuration + AWS-generated metadata
- AWS identifiers: `ID` (bucket name), `ARN`, `BucketDomainName`, `BucketRegionalDomainName`, `Region`
- Used when: Receiving responses from AWS services

## Available Operations

### EBS Volume Operations

- `CreateEBSVolume()` - Create a new EBS volume
- `GetEBSVolume()` - Retrieve volume by ID
- `UpdateEBSVolume()` - Update volume configuration (size, type, IOPS, throughput)
- `DeleteEBSVolume()` - Delete a volume
- `ListEBSVolumes()` - List volumes with filters

### Volume Attachment Operations

- `AttachVolume()` - Attach volume to EC2 instance
- `DetachVolume()` - Detach volume from EC2 instance

### S3 Bucket Operations

- `CreateS3Bucket()` - Create a new S3 bucket
- `GetS3Bucket()` - Retrieve bucket by ID (bucket name)
- `UpdateS3Bucket()` - Update bucket configuration (tags, force_destroy)
- `DeleteS3Bucket()` - Delete a bucket
- `ListS3Buckets()` - List buckets with filters
- `UpdateS3BucketACL()` / `GetS3BucketACL()` - Manage legacy bucket ACLs
- `UpdateS3BucketVersioning()` / `GetS3BucketVersioning()` - Manage bucket versioning state
- `UpdateS3BucketEncryption()` / `GetS3BucketEncryption()` - Manage bucket encryption defaults

## Error Handling

The adapter wraps errors at each layer:

```go
// Domain validation error
if err := volume.Validate(); err != nil {
    return nil, fmt.Errorf("domain validation failed: %w", err)
}

// AWS validation error
if err := awsVolume.Validate(); err != nil {
    return nil, fmt.Errorf("aws validation failed: %w", err)
}

// AWS service error
if err != nil {
    return nil, fmt.Errorf("aws service error: %w", err)
}
```

This provides clear error context while maintaining error chain for debugging.

## Benefits

1. **Separation of Concerns**: Domain layer never knows about AWS
2. **Testability**: Easy to mock AWS service for testing
3. **Extensibility**: Add new providers by creating new adapters
4. **Type Safety**: Compile-time checks ensure interface compliance
5. **Error Context**: Clear error messages with full context

## Testing

The adapter is tested with a mock AWS service to verify:
- Domain-to-AWS conversion
- AWS-to-domain conversion
- Validation at both layers
- Error handling and wrapping
- All CRUD operations

See `adapter_test.go` and `output_integration_test.go` for examples.

## Future Extensions

To add a new cloud provider (e.g., GCP):

1. Create `internal/cloud/gcp/adapters/storage/adapter.go`
2. Implement `domainstorage.StorageService`
3. Use GCP-specific service and mappers
4. Add to factory pattern

The domain layer remains unchanged!
