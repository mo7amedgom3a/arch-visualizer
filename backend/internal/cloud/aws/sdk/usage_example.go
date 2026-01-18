package sdk

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// ExampleUsage demonstrates how to use the AWS SDK client
// This is a simple example showing common EC2 operations
func ExampleUsage() {
	ctx := context.Background()

	// Initialize AWS client
	// This will read credentials from environment variables:
	// - AWS_ACCESS_KEY_ID
	// - AWS_SECRET_ACCESS_KEY
	// - AWS_REGION (optional, defaults to us-east-1)
	client, err := NewAWSClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create AWS client: %v", err)
	}

	fmt.Printf("Connected to AWS region: %s\n", client.GetRegion())

	// Example 1: List all regions
	fmt.Println("\n--- Listing AWS Regions ---")
	regionsOutput, err := client.EC2.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		log.Printf("Error describing regions: %v", err)
	} else {
		for _, region := range regionsOutput.Regions {
			fmt.Printf("Region: %s, Endpoint: %s\n",
				aws.ToString(region.RegionName),
				aws.ToString(region.Endpoint))
		}
	}

	// Example 2: List VPCs
	fmt.Println("\n--- Listing VPCs ---")
	vpcsOutput, err := client.EC2.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		MaxResults: aws.Int32(10),
	})
	if err != nil {
		log.Printf("Error describing VPCs: %v", err)
	} else {
		for _, vpc := range vpcsOutput.Vpcs {
			fmt.Printf("VPC ID: %s, CIDR: %s, State: %s\n",
				aws.ToString(vpc.VpcId),
				aws.ToString(vpc.CidrBlock),
				vpc.State)
		}
	}

	// Example 3: List EC2 instances
	fmt.Println("\n--- Listing EC2 Instances ---")
	instancesOutput, err := client.EC2.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		MaxResults: aws.Int32(10),
	})
	if err != nil {
		log.Printf("Error describing instances: %v", err)
	} else {
		instanceCount := 0
		for _, reservation := range instancesOutput.Reservations {
			for _, instance := range reservation.Instances {
				instanceCount++
				fmt.Printf("Instance ID: %s, Type: %s, State: %s\n",
					aws.ToString(instance.InstanceId),
					instance.InstanceType,
					instance.State.Name)
			}
		}
		if instanceCount == 0 {
			fmt.Println("No instances found")
		}
	}

	// Example 4: List instance types
	fmt.Println("\n--- Listing Instance Types (t3 family) ---")
	instanceTypesOutput, err := client.EC2.DescribeInstanceTypes(ctx, &ec2.DescribeInstanceTypesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-type"),
				Values: []string{"t3.*"},
			},
		},
		MaxResults: aws.Int32(5),
	})
	if err != nil {
		log.Printf("Error describing instance types: %v", err)
	} else {
		for _, instanceType := range instanceTypesOutput.InstanceTypes {
			fmt.Printf("Instance Type: %s\n", string(instanceType.InstanceType))
		}
	}
}
