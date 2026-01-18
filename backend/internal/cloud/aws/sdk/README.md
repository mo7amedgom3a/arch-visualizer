# AWS SDK Integration

This package provides AWS SDK v2 client initialization and configuration for programmatic interaction with AWS services.

## Purpose

The SDK package:
- **Initializes AWS SDK clients** with proper configuration
- **Loads credentials** from environment variables (.env file)
- **Provides reusable client instances** for different AWS services
- **Handles credential chain** (env vars, IAM roles, credentials file)

## Setup

### 1. Environment Variables

Create a `.env` file in the project root with AWS credentials:

```env
AWS_ACCESS_KEY_ID=your-access-key-id
AWS_SECRET_ACCESS_KEY=your-secret-access-key
AWS_REGION=us-east-1
AWS_SESSION_TOKEN=your-session-token  # Optional, for temporary credentials
```

### 2. Loading Environment Variables

The SDK automatically reads from environment variables. The test files (`examples_test.go`) automatically load the `.env` file from the project root using `godotenv`.

**For tests:** The `.env` file is automatically loaded when running tests. Place your `.env` file in the project root (`backend/.env`).

**For application code:** If you want to load `.env` in your application, you can use:

```go
import "github.com/joho/godotenv"

func init() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }
}
```

Or use the autoload feature:

```go
import _ "github.com/joho/godotenv/autoload" // Auto-load .env on import
```

## Usage

### Basic Client Initialization

```go
import (
    "context"
    awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
)

ctx := context.Background()

// Create AWS client (reads from environment variables)
client, err := awssdk.NewAWSClient(ctx)
if err != nil {
    // handle error
}

// Use EC2 client
ec2Client := client.EC2
```

### Explicit Credentials

```go
client, err := awssdk.NewAWSClientWithConfig(
    ctx,
    "us-east-1",                    // region
    "AKIAIOSFODNN7EXAMPLE",        // access key ID
    "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", // secret access key
    "",                             // session token (optional)
)
```

### Using AWS Services

```go
// EC2 Operations
import (
    "github.com/aws/aws-sdk-go-v2/service/ec2"
    "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Describe instances
input := &ec2.DescribeInstancesInput{}
output, err := client.EC2.DescribeInstances(ctx, input)
if err != nil {
    // handle error
}

for _, reservation := range output.Reservations {
    for _, instance := range reservation.Instances {
        fmt.Printf("Instance ID: %s\n", *instance.InstanceId)
        fmt.Printf("State: %s\n", instance.State.Name)
    }
}
```

## Credential Chain

The SDK uses the following credential chain (in order):

1. **Environment Variables** (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. **Shared Credentials File** (`~/.aws/credentials`)
3. **Shared Config File** (`~/.aws/config`)
4. **IAM Roles** (when running on EC2/ECS/Lambda)
5. **Container Credentials** (when running in ECS)
6. **Instance Profile Credentials** (when running on EC2)

## Available Services

### EC2 Client

```go
client.EC2 // *ec2.Client
```

**Operations:**
- `RunInstances()` - Launch EC2 instances
- `DescribeInstances()` - Describe instances
- `TerminateInstances()` - Terminate instances
- `StartInstances()` - Start stopped instances
- `StopInstances()` - Stop running instances
- `RebootInstances()` - Reboot instances

### Adding More Services

To add more AWS services, extend the `AWSClient` struct:

```go
import (
    "github.com/aws/aws-sdk-go-v2/service/vpc"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSClient struct {
    Config aws.Config
    EC2    *ec2.Client
    VPC    *vpc.Client  // Add VPC client
    S3     *s3.Client   // Add S3 client
}

func NewAWSClient(ctx context.Context) (*AWSClient, error) {
    // ... existing code ...
    
    return &AWSClient{
        Config: cfg,
        EC2:    ec2.NewFromConfig(cfg),
        VPC:    vpc.NewFromConfig(cfg),  // Initialize VPC client
        S3:     s3.NewFromConfig(cfg),   // Initialize S3 client
    }, nil
}
```

## Testing with AWS SDK

### Integration Tests

```go
func TestEC2InstanceCreation(t *testing.T) {
    ctx := context.Background()
    
    // Initialize AWS client
    client, err := awssdk.NewAWSClient(ctx)
    if err != nil {
        t.Fatalf("Failed to create AWS client: %v", err)
    }
    
    // Create instance
    input := &ec2.RunInstancesInput{
        ImageId:      aws.String("ami-0c55b159cbfafe1f0"),
        InstanceType: types.InstanceTypeT3Micro,
        MinCount:     aws.Int32(1),
        MaxCount:     aws.Int32(1),
    }
    
    output, err := client.EC2.RunInstances(ctx, input)
    if err != nil {
        t.Fatalf("Failed to create instance: %v", err)
    }
    
    instanceID := output.Instances[0].InstanceId
    t.Logf("Created instance: %s", *instanceID)
    
    // Cleanup
    defer func() {
        _, _ = client.EC2.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
            InstanceIds: []string{*instanceID},
        })
    }()
}
```

### Mock Testing

For unit tests, use interfaces and mocks instead of real AWS SDK calls:

```go
type EC2Client interface {
    RunInstances(ctx context.Context, params *ec2.RunInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)
    DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}
```

## Environment Variables Reference

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `AWS_ACCESS_KEY_ID` | AWS access key ID | Yes* | - |
| `AWS_SECRET_ACCESS_KEY` | AWS secret access key | Yes* | - |
| `AWS_REGION` | AWS region | No | `us-east-1` |
| `AWS_SESSION_TOKEN` | Session token for temporary credentials | No | - |
| `AWS_PROFILE` | AWS profile name (for shared credentials) | No | `default` |

\* Required if not using IAM roles or shared credentials file

## Security Best Practices

1. **Never commit `.env` files** - Add to `.gitignore`
2. **Use IAM roles** when running on AWS infrastructure (EC2, ECS, Lambda)
3. **Rotate credentials** regularly
4. **Use least privilege** IAM policies
5. **Use temporary credentials** when possible (AWS STS)

## Example .env File

```env
# AWS Credentials
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_REGION=us-east-1

# Optional: Session token for temporary credentials
# AWS_SESSION_TOKEN=your-session-token

# Optional: AWS Profile (for shared credentials)
# AWS_PROFILE=my-profile
```

## Related Documentation

- [AWS SDK for Go v2 Documentation](https://aws.github.io/aws-sdk-go-v2/docs/)
- [AWS Credentials Guide](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials)
- [EC2 Service API Reference](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/)
