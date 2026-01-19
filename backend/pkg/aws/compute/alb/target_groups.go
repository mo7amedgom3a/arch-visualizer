package alb

import (
	"context"
	"fmt"

	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/compute"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
	awscomputeservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
)

// TargetGroupRunner demonstrates Target Group operations
func TargetGroupRunner() {
	ctx := context.Background()

	fmt.Println("============================================")
	fmt.Println("TARGET GROUP OPERATIONS")
	fmt.Println("============================================")

	// Initialize AWS client
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		fmt.Printf("Error creating AWS client: %v\n", err)
		return
	}

	// Initialize compute service
	computeService := awscomputeservice.NewComputeService(client)

	region := client.GetRegion()
	fmt.Printf("\nRegion: %s\n", region)

	// Create a target group
	fmt.Println("\n--- Creating Target Group ---")
	targetGroup := &domaincompute.TargetGroup{
		Name:       "demo-target-group",
		VPCID:      "vpc-12345678",
		Port:       80,
		Protocol:   domaincompute.TargetGroupProtocolHTTP,
		TargetType: domaincompute.TargetTypeInstance,
		HealthCheck: domaincompute.HealthCheckConfig{
			Path:    stringPtr("/health"),
			Matcher: stringPtr("200"),
		},
	}

	// Convert domain model to AWS model
	awsTG := &awsloadbalancer.TargetGroup{
		Name:       targetGroup.Name,
		VPCID:      targetGroup.VPCID,
		Port:       targetGroup.Port,
		Protocol:   string(targetGroup.Protocol),
		TargetType: stringPtr(string(targetGroup.TargetType)),
		HealthCheck: awsloadbalancer.HealthCheckConfig{
			Path:    targetGroup.HealthCheck.Path,
			Matcher: targetGroup.HealthCheck.Matcher,
		},
	}

	awsTGOutput, err := computeService.CreateTargetGroup(ctx, awsTG)
	if err != nil {
		fmt.Printf("Error creating target group: %v\n", err)
		return
	}

	createdTG := awsmapper.ToDomainTargetGroupFromOutput(awsTGOutput)
	fmt.Printf("Target Group created:\n")
	fmt.Printf("  Name: %s\n", createdTG.Name)
	if createdTG.ARN != nil {
		fmt.Printf("  ARN: %s\n", *createdTG.ARN)
	}

	// List target groups
	fmt.Println("\n--- Listing Target Groups ---")
	awsTargetGroups, err := computeService.ListTargetGroups(ctx, map[string][]string{})
	if err != nil {
		fmt.Printf("Error listing target groups: %v\n", err)
		return
	}

	fmt.Printf("Found %d target group(s):\n", len(awsTargetGroups))
	for i, awsTG := range awsTargetGroups {
		tg := awsmapper.ToDomainTargetGroupFromOutput(awsTG)
		fmt.Printf("  %d. %s (Port: %d, Protocol: %s)\n", i+1, tg.Name, tg.Port, tg.Protocol)
	}
}
