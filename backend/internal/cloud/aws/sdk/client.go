package sdk

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// AWSClient wraps AWS SDK clients for different services
type AWSClient struct {
	Config aws.Config
	EC2    *ec2.Client
	IAM    *iam.Client
	ELBv2  *elasticloadbalancingv2.Client
}

// NewAWSClient creates a new AWS client with configuration from environment variables
// Reads credentials from .env file or environment variables:
// - AWS_ACCESS_KEY_ID
// - AWS_SECRET_ACCESS_KEY
func NewAWSClient(ctx context.Context) (*AWSClient, error) {
	// Get credentials provider (may be nil to use default chain)
	credsProvider := getCredentialsProvider()

	// Build config options
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(getRegion()),
	}

	// Add credentials provider if provided (from env vars)
	if credsProvider != nil {
		opts = append(opts, config.WithCredentialsProvider(credsProvider))
	}

	// Load default config (reads from environment variables, ~/.aws/credentials, etc.)
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &AWSClient{
		Config: cfg,
		EC2:    ec2.NewFromConfig(cfg),
		IAM:    iam.NewFromConfig(cfg),
		ELBv2:  elasticloadbalancingv2.NewFromConfig(cfg),
	}, nil
}

// NewAWSClientWithConfig creates a new AWS client with explicit configuration
func NewAWSClientWithConfig(ctx context.Context, region string, accessKeyID, secretAccessKey, sessionToken string) (*AWSClient, error) {
	var credsProvider aws.CredentialsProvider

	if accessKeyID != "" && secretAccessKey != "" {
		credsProvider = credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			sessionToken,
		)
	}

	if region == "" {
		region = getRegion()
	}

	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	if credsProvider != nil {
		opts = append(opts, config.WithCredentialsProvider(credsProvider))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &AWSClient{
		Config: cfg,
		EC2:    ec2.NewFromConfig(cfg),
		IAM:    iam.NewFromConfig(cfg),
		ELBv2:  elasticloadbalancingv2.NewFromConfig(cfg),
	}, nil
}

// getRegion returns the AWS region from environment variable or default
func getRegion() string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1" // Default region
	}
	return region
}

// getCredentialsProvider returns credentials provider from environment variables
func getCredentialsProvider() aws.CredentialsProvider {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if accessKeyID != "" && secretAccessKey != "" {
		return credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"", // session token
		)
	}

	// Return nil to use default credential chain (IAM roles, ~/.aws/credentials, etc.)
	return nil
}

// GetRegion returns the configured AWS region
func (c *AWSClient) GetRegion() string {
	return c.Config.Region
}
