package autoscaling

import (
	"context"
	"fmt"
	"time"

	awsautoscaling "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/autoscaling"
	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
	awscomputeservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
)

// ASGRunner demonstrates Auto Scaling Group setup with Launch Template and scaling policies
func ASGRunner() {
	ctx := context.Background()

	fmt.Println("============================================")
	fmt.Println("AUTO SCALING GROUP SETUP")
	fmt.Println("============================================")

	// Initialize AWS client
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		fmt.Printf("Error creating AWS client: %v\n", err)
		return
	}

	region := "us-east-1"
	fmt.Printf("\nRegion: %s\n", region)

	// Initialize compute service for ASG operations
	computeService := awscomputeservice.NewComputeService()

	// Use mock mode - skip real SDK calls
	useMockData := true

	// Step 1: Create Launch Template (mock)
	fmt.Println("\n--- Step 1: Creating Launch Template ---")
	fmt.Println("  [MOCK MODE] Using mock launch template data")

	// Create mock launch template
	name := "asg-launch-template"
	arn := fmt.Sprintf("arn:aws:ec2:%s:123456789012:launch-template/lt-mock-1234567890abcdef0", region)
	launchTemplate := &domaincompute.LaunchTemplate{
		ID:               "lt-mock-1234567890abcdef0",
		ARN:              &arn,
		Name:             name,
		Region:           region,
		ImageID:          "ami-0c55b159cbfafe1f0",
		InstanceType:     "t3.micro",
		SecurityGroupIDs: []string{"sg-123"},
	}

	fmt.Printf("Launch Template created successfully:\n")
	fmt.Printf("  ID: %s\n", launchTemplate.ID)
	if launchTemplate.ARN != nil {
		fmt.Printf("  ARN: %s\n", *launchTemplate.ARN)
	}
	fmt.Printf("  Name: %s\n", launchTemplate.Name)

	// Step 2: Create Auto Scaling Group
	fmt.Println("\n--- Step 2: Creating Auto Scaling Group ---")
	autoScalingGroup, err := createAutoScalingGroup(ctx, computeService, client, region, launchTemplate, useMockData)
	if err != nil {
		fmt.Printf("Error creating auto scaling group: %v\n", err)
		return
	}

	if autoScalingGroup == nil {
		fmt.Println("Auto Scaling Group creation failed. Cannot proceed.")
		return
	}

	fmt.Printf("Auto Scaling Group created successfully:\n")
	fmt.Printf("  Name: %s\n", autoScalingGroup.Name)
	if autoScalingGroup.ARN != nil {
		fmt.Printf("  ARN: %s\n", *autoScalingGroup.ARN)
	}
	fmt.Printf("  Min Size: %d\n", autoScalingGroup.MinSize)
	fmt.Printf("  Max Size: %d\n", autoScalingGroup.MaxSize)
	if autoScalingGroup.DesiredCapacity != nil {
		fmt.Printf("  Desired Capacity: %d\n", *autoScalingGroup.DesiredCapacity)
	}
	fmt.Printf("  Health Check Type: %s\n", autoScalingGroup.HealthCheckType)
	fmt.Printf("  State: %s\n", autoScalingGroup.State)

	// Step 3: Create Target Group (optional, for ELB health checks)
	fmt.Println("\n--- Step 3: Creating Target Group (for ELB health checks) ---")
	targetGroup, err := createTargetGroup(ctx, computeService, region, useMockData)
	if err != nil {
		fmt.Printf("Warning: Error creating target group: %v\n", err)
		fmt.Println("Continuing without target group (using EC2 health checks)")
	} else if targetGroup != nil {
		fmt.Printf("Target Group created successfully:\n")
		fmt.Printf("  Name: %s\n", targetGroup.Name)
		if targetGroup.ARN != nil {
			fmt.Printf("  ARN: %s\n", *targetGroup.ARN)
		}

		// Update ASG to use ELB health checks
		fmt.Println("\n--- Step 3.1: Updating ASG to use ELB health checks ---")
		if autoScalingGroup.ARN != nil {
			autoScalingGroup.HealthCheckType = domaincompute.AutoScalingGroupHealthCheckTypeELB
			autoScalingGroup.TargetGroupARNs = []string{*targetGroup.ARN}
			updatedASG, err := updateAutoScalingGroup(ctx, computeService, client, autoScalingGroup, useMockData)
			if err != nil {
				fmt.Printf("Warning: Error updating ASG: %v\n", err)
			} else {
				fmt.Printf("ASG updated to use ELB health checks\n")
				autoScalingGroup = updatedASG
			}
		}
	}

	// Step 4: Create Scaling Policy (Target Tracking)
	fmt.Println("\n--- Step 4: Creating Scaling Policy ---")
	err = createScalingPolicy(ctx, computeService, client, autoScalingGroup, useMockData)
	if err != nil {
		fmt.Printf("Warning: Error creating scaling policy: %v\n", err)
		fmt.Println("Continuing without scaling policy")
	} else {
		fmt.Println("  Scaling Policy created successfully (mock mode)")
	}

	// Step 5: Demonstrate scaling operations
	fmt.Println("\n--- Step 5: Demonstrating Scaling Operations ---")
	demonstrateScalingOperations(ctx, computeService, client, autoScalingGroup.Name, useMockData)

	// Step 6: Verify setup
	fmt.Println("\n--- Step 6: Verifying Setup ---")
	verifyASGSetup(ctx, computeService, client, autoScalingGroup.Name, useMockData)

	fmt.Println("\n============================================")
	fmt.Println("AUTO SCALING GROUP SETUP COMPLETE")
	fmt.Println("============================================")
}

// createAutoScalingGroup creates an Auto Scaling Group
func createAutoScalingGroup(ctx context.Context, computeService *awscomputeservice.ComputeService, client *awssdk.AWSClient, region string, launchTemplate *domaincompute.LaunchTemplate, useMockData bool) (*domaincompute.AutoScalingGroup, error) {
	if useMockData {
		fmt.Println("  [MOCK MODE] Using mock auto scaling group data")
		// Return mock ASG
		version := "$Latest"
		desiredCapacity := 2
		gracePeriod := 300
		arn := "arn:aws:autoscaling:us-east-1:123456789012:autoScalingGroup:uuid:autoScalingGroupName/test-asg"
		createdTime := time.Now().Format("2006-01-02T15:04:05Z07:00")
		return &domaincompute.AutoScalingGroup{
			ID:                "test-asg",
			ARN:               &arn,
			Name:              "test-asg",
			Region:            region,
			MinSize:           1,
			MaxSize:           5,
			DesiredCapacity:   &desiredCapacity,
			VPCZoneIdentifier: []string{"subnet-123", "subnet-456"},
			LaunchTemplate: &domaincompute.LaunchTemplateSpecification{
				ID:      launchTemplate.ID,
				Version: &version,
			},
			HealthCheckType:        domaincompute.AutoScalingGroupHealthCheckTypeEC2,
			HealthCheckGracePeriod: &gracePeriod,
			TargetGroupARNs:        []string{},
			Tags:                   []domaincompute.Tag{},
			State:                  domaincompute.AutoScalingGroupStateActive,
			CreatedTime:            &createdTime,
		}, nil
	}

	// Real SDK call
	fmt.Println("  Creating Auto Scaling Group via SDK...")
	version := "$Latest"
	desiredCapacity := 2
	gracePeriod := 300
	awsASG := &awsautoscaling.AutoScalingGroup{
		AutoScalingGroupName: stringPtr("test-asg"),
		MinSize:              1,
		MaxSize:              5,
		DesiredCapacity:      &desiredCapacity,
		VPCZoneIdentifier:    []string{"subnet-123", "subnet-456"},
		LaunchTemplate: &awsautoscaling.LaunchTemplateSpecification{
			LaunchTemplateId: launchTemplate.ID,
			Version:          &version,
		},
		HealthCheckType:        stringPtr("EC2"),
		HealthCheckGracePeriod: &gracePeriod,
		Tags: []awsautoscaling.Tag{
			{
				Key:               "Name",
				Value:             "test-asg",
				PropagateAtLaunch: true,
			},
		},
	}

	awsASGOutput, err := awssdk.CreateAutoScalingGroup(ctx, client, awsASG)
	if err != nil {
		return nil, fmt.Errorf("failed to create ASG: %w", err)
	}

	// Convert to domain model (simplified - would use mapper in real implementation)
	arn := awsASGOutput.AutoScalingGroupARN
	createdTime := awsASGOutput.CreatedTime.Format("2006-01-02T15:04:05Z07:00")
	return &domaincompute.AutoScalingGroup{
		ID:                awsASGOutput.AutoScalingGroupName,
		ARN:               &arn,
		Name:              awsASGOutput.AutoScalingGroupName,
		Region:            region,
		MinSize:           awsASGOutput.MinSize,
		MaxSize:           awsASGOutput.MaxSize,
		DesiredCapacity:   &awsASGOutput.DesiredCapacity,
		VPCZoneIdentifier: awsASGOutput.VPCZoneIdentifier,
		LaunchTemplate: &domaincompute.LaunchTemplateSpecification{
			ID:      awsASGOutput.LaunchTemplate.LaunchTemplateId,
			Version: awsASGOutput.LaunchTemplate.Version,
		},
		HealthCheckType:        domaincompute.AutoScalingGroupHealthCheckType(awsASGOutput.HealthCheckType),
		HealthCheckGracePeriod: awsASGOutput.HealthCheckGracePeriod,
		TargetGroupARNs:        awsASGOutput.TargetGroupARNs,
		State:                  domaincompute.AutoScalingGroupStateActive,
		CreatedTime:            &createdTime,
	}, nil
}

func stringPtr(s string) *string {
	return &s
}

// updateAutoScalingGroup updates an Auto Scaling Group
func updateAutoScalingGroup(ctx context.Context, computeService *awscomputeservice.ComputeService, client *awssdk.AWSClient, asg *domaincompute.AutoScalingGroup, useMockData bool) (*domaincompute.AutoScalingGroup, error) {
	if useMockData {
		fmt.Println("  [MOCK MODE] Using mock updated auto scaling group data")
		return asg, nil
	}

	// Real SDK call would go here
	version := "$Latest"
	awsASG := &awsautoscaling.AutoScalingGroup{
		AutoScalingGroupName: stringPtr(asg.Name),
		MinSize:              asg.MinSize,
		MaxSize:              asg.MaxSize,
		DesiredCapacity:      asg.DesiredCapacity,
		VPCZoneIdentifier:    asg.VPCZoneIdentifier,
		LaunchTemplate: &awsautoscaling.LaunchTemplateSpecification{
			LaunchTemplateId: asg.LaunchTemplate.ID,
			Version:          &version,
		},
		HealthCheckType:        stringPtr(string(asg.HealthCheckType)),
		HealthCheckGracePeriod: asg.HealthCheckGracePeriod,
		TargetGroupARNs:        asg.TargetGroupARNs,
	}

	awsASGOutput, err := awssdk.UpdateAutoScalingGroup(ctx, client, asg.Name, awsASG)
	if err != nil {
		return nil, fmt.Errorf("failed to update ASG: %w", err)
	}

	// Convert to domain model
	arn := awsASGOutput.AutoScalingGroupARN
	createdTime := awsASGOutput.CreatedTime.Format("2006-01-02T15:04:05Z07:00")
	region := "us-east-1"
	return &domaincompute.AutoScalingGroup{
		ID:                awsASGOutput.AutoScalingGroupName,
		ARN:               &arn,
		Name:              awsASGOutput.AutoScalingGroupName,
		Region:            region,
		MinSize:           awsASGOutput.MinSize,
		MaxSize:           awsASGOutput.MaxSize,
		DesiredCapacity:   &awsASGOutput.DesiredCapacity,
		VPCZoneIdentifier: awsASGOutput.VPCZoneIdentifier,
		LaunchTemplate: &domaincompute.LaunchTemplateSpecification{
			ID:      awsASGOutput.LaunchTemplate.LaunchTemplateId,
			Version: awsASGOutput.LaunchTemplate.Version,
		},
		HealthCheckType:        domaincompute.AutoScalingGroupHealthCheckType(awsASGOutput.HealthCheckType),
		HealthCheckGracePeriod: awsASGOutput.HealthCheckGracePeriod,
		TargetGroupARNs:        awsASGOutput.TargetGroupARNs,
		State:                  domaincompute.AutoScalingGroupStateActive,
		CreatedTime:            &createdTime,
	}, nil
}

// createTargetGroup creates a target group for ELB health checks
func createTargetGroup(ctx context.Context, computeService *awscomputeservice.ComputeService, region string, useMockData bool) (*domaincompute.TargetGroup, error) {
	if useMockData {
		fmt.Println("  [MOCK MODE] Using mock target group data")
		// Return mock target group
		arn := "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef"
		return &domaincompute.TargetGroup{
			ID:         arn,
			ARN:        &arn,
			Name:       "test-tg",
			VPCID:      "vpc-123",
			Port:       80,
			Protocol:   domaincompute.TargetGroupProtocolHTTP,
			TargetType: domaincompute.TargetTypeInstance,
			State:      domaincompute.TargetGroupStateActive,
		}, nil
	}

	// Real SDK call would go here
	// For now, return mock data
	fmt.Println("  [MOCK MODE] Using mock target group data")
	arn := "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/test-tg/1234567890abcdef"
	return &domaincompute.TargetGroup{
		ID:         arn,
		ARN:        &arn,
		Name:       "test-tg",
		VPCID:      "vpc-123",
		Port:       80,
		Protocol:   domaincompute.TargetGroupProtocolHTTP,
		TargetType: domaincompute.TargetTypeInstance,
		State:      domaincompute.TargetGroupStateActive,
	}, nil
}

// createScalingPolicy creates a target tracking scaling policy
func createScalingPolicy(ctx context.Context, computeService *awscomputeservice.ComputeService, client *awssdk.AWSClient, asg *domaincompute.AutoScalingGroup, useMockData bool) error {
	if useMockData {
		fmt.Println("  [MOCK MODE] Using mock scaling policy data")
		fmt.Println("  Scaling policy would be created via AWS service layer")
		return nil
	}

	// Real SDK call
	fmt.Println("  Creating scaling policy via SDK...")
	policy := &awsautoscaling.ScalingPolicy{
		PolicyName:           "target-tracking-policy",
		AutoScalingGroupName: asg.Name,
		PolicyType:           awsautoscaling.ScalingPolicyTypeTargetTrackingScaling,
		TargetTrackingConfiguration: &awsautoscaling.TargetTrackingConfiguration{
			TargetValue: 70.0,
			PredefinedMetricSpecification: &awsautoscaling.PredefinedMetricSpecification{
				PredefinedMetricType: "ASGAverageCPUUtilization",
			},
		},
	}

	_, err := computeService.PutScalingPolicy(ctx, policy)
	if err != nil {
		return fmt.Errorf("failed to create scaling policy: %w", err)
	}

	return nil
}

// demonstrateScalingOperations demonstrates scaling operations
func demonstrateScalingOperations(ctx context.Context, computeService *awscomputeservice.ComputeService, client *awssdk.AWSClient, asgName string, useMockData bool) {
	if useMockData {
		fmt.Println("  [MOCK MODE] Simulating scaling operations")

		fmt.Println("  Setting desired capacity to 3...")
		fmt.Println("  ✓ Desired capacity set to 3 (simulated)")

		fmt.Println("  Attaching instances...")
		instanceIDs := []string{"i-1234567890abcdef0", "i-0987654321fedcba0"}
		fmt.Printf("  ✓ Attached %d instances (simulated)\n", len(instanceIDs))
		for i, instanceID := range instanceIDs {
			fmt.Printf("    Instance %d: %s (simulated)\n", i+1, instanceID)
		}

		fmt.Println("  Detaching instances...")
		fmt.Println("  ✓ Detached 1 instance (simulated)")
		fmt.Printf("    Detached: i-1234567890abcdef0 (simulated)\n")
		return
	}

	// Real SDK calls
	fmt.Println("  Setting desired capacity to 3...")
	err := computeService.SetDesiredCapacity(ctx, asgName, 3)
	if err != nil {
		fmt.Printf("  Error setting desired capacity: %v\n", err)
	} else {
		fmt.Println("  ✓ Desired capacity set to 3")
	}

	fmt.Println("  Attaching instances...")
	instanceIDs := []string{"i-1234567890abcdef0", "i-0987654321fedcba0"}
	err = computeService.AttachInstances(ctx, asgName, instanceIDs)
	if err != nil {
		fmt.Printf("  Error attaching instances: %v\n", err)
	} else {
		fmt.Printf("  ✓ Attached %d instances\n", len(instanceIDs))
	}

	fmt.Println("  Detaching instances...")
	err = computeService.DetachInstances(ctx, asgName, []string{"i-1234567890abcdef0"})
	if err != nil {
		fmt.Printf("  Error detaching instances: %v\n", err)
	} else {
		fmt.Println("  ✓ Detached 1 instance")
	}
}

// verifyASGSetup verifies the ASG setup
func verifyASGSetup(ctx context.Context, computeService *awscomputeservice.ComputeService, client *awssdk.AWSClient, asgName string, useMockData bool) {
	if useMockData {
		fmt.Println("  [MOCK MODE] Using mock ASG data for verification")

		// Mock verification - just print what we created
		fmt.Printf("  ✓ ASG verified: %s (State: active) [MOCK]\n", asgName)
		fmt.Printf("  ASG Name: %s [MOCK]\n", asgName)
		fmt.Printf("  Min Size: 1 [MOCK]\n")
		fmt.Printf("  Max Size: 5 [MOCK]\n")
		fmt.Printf("  Desired Capacity: 2 [MOCK]\n")
		fmt.Printf("  Health Check Type: EC2 [MOCK]\n")
		fmt.Printf("  Status: active [MOCK]\n")
		fmt.Printf("  Instances: 2 [MOCK]\n")
		fmt.Printf("    Instance 1: i-1234567890abcdef0 (Health: healthy) [MOCK]\n")
		fmt.Printf("    Instance 2: i-0987654321fedcba0 (Health: healthy) [MOCK]\n")
		fmt.Printf("  Total ASGs in account: 1 [MOCK]\n")
		return
	}

	// Real SDK calls
	asgOutput, err := computeService.GetAutoScalingGroup(ctx, asgName)
	if err != nil {
		fmt.Printf("  Error getting ASG: %v\n", err)
		return
	}

	fmt.Printf("  ASG Name: %s\n", asgOutput.AutoScalingGroupName)
	fmt.Printf("  Min Size: %d\n", asgOutput.MinSize)
	fmt.Printf("  Max Size: %d\n", asgOutput.MaxSize)
	fmt.Printf("  Desired Capacity: %d\n", asgOutput.DesiredCapacity)
	fmt.Printf("  Health Check Type: %s\n", asgOutput.HealthCheckType)
	fmt.Printf("  Status: %s\n", asgOutput.Status)
	fmt.Printf("  Instances: %d\n", len(asgOutput.Instances))

	// List all ASGs
	asgs, err := computeService.ListAutoScalingGroups(ctx, map[string][]string{})
	if err != nil {
		fmt.Printf("  Error listing ASGs: %v\n", err)
	} else {
		fmt.Printf("  Total ASGs in account: %d\n", len(asgs))
	}
}
