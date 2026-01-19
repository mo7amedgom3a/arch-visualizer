package alb

import (
	"context"
	"fmt"
	"time"

	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/compute"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awssdk "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/sdk"
	awscomputeservice "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/compute"
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
)

// ALBRunner demonstrates Application Load Balancer setup with EC2 instances
func ALBRunner() {
	ctx := context.Background()

	fmt.Println("============================================")
	fmt.Println("APPLICATION LOAD BALANCER SETUP")
	fmt.Println("============================================")

	// Initialize AWS client
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		fmt.Printf("Error creating AWS client: %v\n", err)
		return
	}

	region := client.GetRegion()
	fmt.Printf("\nRegion: %s\n", region)

	// Initialize compute service for load balancer operations
	// Note: Since AWS account doesn't support creating load balancers,
	// we'll use mock data for all operations
	computeService := awscomputeservice.NewComputeService(client)

	// Use mock mode - skip real SDK calls
	useMockData := true

	// Step 1: Create EC2 instances
	fmt.Println("\n--- Step 1: Creating EC2 Instances ---")
	instances, err := createEC2Instances(ctx, region)
	if err != nil {
		fmt.Printf("Error creating EC2 instances: %v\n", err)
		return
	}

	if len(instances) == 0 {
		fmt.Println("No instances created. Cannot proceed with load balancer setup.")
		return
	}

	fmt.Printf("Successfully created %d instance(s)\n", len(instances))
	for i, instance := range instances {
		fmt.Printf("  Instance %d: ID=%s, Name=%s\n", i+1, instance.ID, instance.Name)
	}

	// Step 2: Create Application Load Balancer
	fmt.Println("\n--- Step 2: Creating Application Load Balancer ---")
	loadBalancer, err := createLoadBalancer(ctx, computeService, client, region, instances, useMockData)
	if err != nil {
		fmt.Printf("Error creating load balancer: %v\n", err)
		return
	}

	if loadBalancer == nil {
		fmt.Println("Load balancer creation failed. Cannot proceed.")
		return
	}

	fmt.Printf("Load Balancer created successfully:\n")
	fmt.Printf("  Name: %s\n", loadBalancer.Name)
	if loadBalancer.ARN != nil {
		fmt.Printf("  ARN: %s\n", *loadBalancer.ARN)
	}
	if loadBalancer.DNSName != nil {
		fmt.Printf("  DNS Name: %s\n", *loadBalancer.DNSName)
	}
	fmt.Printf("  State: %s\n", loadBalancer.State)

	// Step 3: Create Target Group
	fmt.Println("\n--- Step 3: Creating Target Group ---")
	targetGroup, err := createTargetGroup(ctx, computeService, client, region, useMockData)
	if err != nil {
		fmt.Printf("Error creating target group: %v\n", err)
		return
	}

	if targetGroup == nil {
		fmt.Println("Target group creation failed. Cannot proceed.")
		return
	}

	fmt.Printf("Target Group created successfully:\n")
	fmt.Printf("  Name: %s\n", targetGroup.Name)
	if targetGroup.ARN != nil {
		fmt.Printf("  ARN: %s\n", *targetGroup.ARN)
	}
	fmt.Printf("  Port: %d\n", targetGroup.Port)
	fmt.Printf("  Protocol: %s\n", targetGroup.Protocol)
	fmt.Printf("  State: %s\n", targetGroup.State)

	// Step 4: Create Listener
	fmt.Println("\n--- Step 4: Creating Listener ---")
	listener, err := createListener(ctx, computeService, client, loadBalancer, targetGroup, useMockData)
	if err != nil {
		fmt.Printf("Error creating listener: %v\n", err)
		return
	}

	if listener == nil {
		fmt.Println("Listener creation failed.")
		return
	}

	fmt.Printf("Listener created successfully:\n")
	if listener.ARN != nil {
		fmt.Printf("  ARN: %s\n", *listener.ARN)
	}
	fmt.Printf("  Port: %d\n", listener.Port)
	fmt.Printf("  Protocol: %s\n", listener.Protocol)

	// Step 5: Attach instances to Target Group
	fmt.Println("\n--- Step 5: Attaching Instances to Target Group ---")
	err = attachInstancesToTargetGroup(ctx, computeService, client, targetGroup, instances, useMockData)
	if err != nil {
		fmt.Printf("Error attaching instances to target group: %v\n", err)
		return
	}

	fmt.Printf("Successfully attached %d instance(s) to target group\n", len(instances))

	// Step 6: Verify setup
	fmt.Println("\n--- Step 6: Verifying Setup ---")
	err = verifySetup(ctx, computeService, client, loadBalancer, targetGroup, instances, useMockData)
	if err != nil {
		fmt.Printf("Error verifying setup: %v\n", err)
		return
	}

	fmt.Println("\n============================================")
	fmt.Println("LOAD BALANCER SETUP COMPLETE!")
	fmt.Println("============================================")
	fmt.Printf("\nLoad Balancer DNS: %s\n", *loadBalancer.DNSName)
	fmt.Printf("Target Group: %s\n", targetGroup.Name)
	fmt.Printf("Instances attached: %d\n", len(instances))
}

// createEC2Instances creates mock EC2 instances for the load balancer
func createEC2Instances(ctx context.Context, region string) ([]*domaincompute.Instance, error) {
	// Mock subnet and security group IDs (in real scenario, these would be created/fetched)
	subnetID := "subnet-12345678"
	securityGroupID := "sg-12345678"

	instances := []*domaincompute.Instance{
		{
			Name:             "alb-backend-instance-1",
			Region:           region,
			InstanceType:     "t3.micro",
			AMI:              "ami-0c55b159cbfafe1f0", // Amazon Linux 2 AMI
			SubnetID:         subnetID,
			SecurityGroupIDs: []string{securityGroupID},
		},
		{
			Name:             "alb-backend-instance-2",
			Region:           region,
			InstanceType:     "t3.micro",
			AMI:              "ami-0c55b159cbfafe1f0",
			SubnetID:         subnetID,
			SecurityGroupIDs: []string{securityGroupID},
		},
	}

	var createdInstances []*domaincompute.Instance

	for _, instance := range instances {
		fmt.Printf("Creating instance: %s...\n", instance.Name)

		// Note: In a real scenario, this would create actual instances
		// For demonstration, we'll simulate the creation
		// Note: CreateInstance may not be available if ComputeService doesn't implement it
		// For demonstration, we'll create mock instances
		fmt.Printf("  Creating mock instance for demonstration (real instance creation requires full AWSComputeService implementation)\n")

		// Create mock instance with generated ID
		mockInstance := &domaincompute.Instance{
			ID:               fmt.Sprintf("i-%s", generateMockID()),
			ARN:              stringPtr(fmt.Sprintf("arn:aws:ec2:%s:123456789012:instance/i-%s", region, generateMockID())),
			Name:             instance.Name,
			Region:           instance.Region,
			InstanceType:     instance.InstanceType,
			AMI:              instance.AMI,
			SubnetID:         instance.SubnetID,
			SecurityGroupIDs: instance.SecurityGroupIDs,
			State:            domaincompute.InstanceStateRunning,
			PrivateIP:        stringPtr(fmt.Sprintf("10.0.1.%d", len(createdInstances)+100)),
		}
		createdInstances = append(createdInstances, mockInstance)
		fmt.Printf("  Mock instance created: %s\n", mockInstance.ID)

		// Small delay to simulate real-world creation time
		time.Sleep(500 * time.Millisecond)
	}

	return createdInstances, nil
}

// createLoadBalancer creates an Application Load Balancer
func createLoadBalancer(ctx context.Context, service *awscomputeservice.ComputeService, client *awssdk.AWSClient, region string, instances []*domaincompute.Instance, useMockData bool) (*domaincompute.LoadBalancer, error) {
	// Extract subnet IDs from instances (in real scenario, use multiple subnets across AZs)
	subnetIDs := []string{}
	if len(instances) > 0 {
		subnetIDs = append(subnetIDs, instances[0].SubnetID)
		// Add a second subnet for high availability (mock)
		subnetIDs = append(subnetIDs, "subnet-87654321")
	}

	// Extract security group IDs from instances
	securityGroupIDs := []string{}
	if len(instances) > 0 && len(instances[0].SecurityGroupIDs) > 0 {
		securityGroupIDs = instances[0].SecurityGroupIDs
	}

	loadBalancer := &domaincompute.LoadBalancer{
		Name:             "demo-alb",
		Region:           region,
		Type:             domaincompute.LoadBalancerTypeApplication,
		Internal:         false, // Internet-facing
		SecurityGroupIDs: securityGroupIDs,
		SubnetIDs:        subnetIDs,
	}

	fmt.Printf("Creating load balancer: %s...\n", loadBalancer.Name)

	if useMockData {
		fmt.Printf("  Using mock load balancer data (AWS account doesn't support load balancer creation)\n")
		mockLB := &domaincompute.LoadBalancer{
			ID:               fmt.Sprintf("arn:aws:elasticloadbalancing:%s:123456789012:loadbalancer/app/demo-alb/%s", region, generateMockID()),
			ARN:              stringPtr(fmt.Sprintf("arn:aws:elasticloadbalancing:%s:123456789012:loadbalancer/app/demo-alb/%s", region, generateMockID())),
			Name:             loadBalancer.Name,
			Region:           region,
			Type:             loadBalancer.Type,
			Internal:         loadBalancer.Internal,
			SecurityGroupIDs: loadBalancer.SecurityGroupIDs,
			SubnetIDs:        loadBalancer.SubnetIDs,
			DNSName:          stringPtr(fmt.Sprintf("demo-alb-%s.%s.elb.amazonaws.com", generateMockID()[:8], region)),
			ZoneID:           stringPtr("Z35SXDOTRQ7X7K"),
			State:            domaincompute.LoadBalancerStateActive,
		}
		return mockLB, nil
	}

	// Convert domain model to AWS model
	awsLB := &awsloadbalancer.LoadBalancer{
		Name:             loadBalancer.Name,
		LoadBalancerType: "application",
		Internal:         &loadBalancer.Internal,
		SecurityGroupIDs: loadBalancer.SecurityGroupIDs,
		SubnetIDs:        loadBalancer.SubnetIDs,
	}

	awsLBOutput, err := service.CreateLoadBalancer(ctx, awsLB)
	if err != nil {
		return nil, fmt.Errorf("failed to create load balancer: %w", err)
	}

	// Convert AWS output to domain model
	return awsmapper.ToDomainLoadBalancerFromOutput(awsLBOutput), nil
}

// createTargetGroup creates a Target Group for the load balancer
func createTargetGroup(ctx context.Context, service *awscomputeservice.ComputeService, client *awssdk.AWSClient, region string, useMockData bool) (*domaincompute.TargetGroup, error) {
	// Mock VPC ID (in real scenario, this would be fetched)
	vpcID := "vpc-12345678"

	healthCheckPath := "/health"
	healthCheckMatcher := "200"

	targetGroup := &domaincompute.TargetGroup{
		Name:       "demo-target-group",
		VPCID:      vpcID,
		Port:       80,
		Protocol:   domaincompute.TargetGroupProtocolHTTP,
		TargetType: domaincompute.TargetTypeInstance,
		HealthCheck: domaincompute.HealthCheckConfig{
			Path:               &healthCheckPath,
			Matcher:            &healthCheckMatcher,
			Interval:           intPtr(30),
			Timeout:            intPtr(5),
			HealthyThreshold:   intPtr(2),
			UnhealthyThreshold: intPtr(3),
		},
	}

	fmt.Printf("Creating target group: %s...\n", targetGroup.Name)

	if useMockData {
		fmt.Printf("  Using mock target group data (AWS account doesn't support load balancer creation)\n")
		mockTG := &domaincompute.TargetGroup{
			ID:          fmt.Sprintf("arn:aws:elasticloadbalancing:%s:123456789012:targetgroup/demo-target-group/%s", region, generateMockID()),
			ARN:         stringPtr(fmt.Sprintf("arn:aws:elasticloadbalancing:%s:123456789012:targetgroup/demo-target-group/%s", region, generateMockID())),
			Name:        targetGroup.Name,
			VPCID:       targetGroup.VPCID,
			Port:        targetGroup.Port,
			Protocol:    targetGroup.Protocol,
			TargetType:  targetGroup.TargetType,
			HealthCheck: targetGroup.HealthCheck,
			State:       domaincompute.TargetGroupStateActive,
		}
		return mockTG, nil
	}

	// Convert domain model to AWS model
	awsTG := &awsloadbalancer.TargetGroup{
		Name:       targetGroup.Name,
		VPCID:      targetGroup.VPCID,
		Port:       targetGroup.Port,
		Protocol:   string(targetGroup.Protocol),
		TargetType: stringPtr(string(targetGroup.TargetType)),
		HealthCheck: awsloadbalancer.HealthCheckConfig{
			Path:               targetGroup.HealthCheck.Path,
			Matcher:            targetGroup.HealthCheck.Matcher,
			Interval:           targetGroup.HealthCheck.Interval,
			Timeout:            targetGroup.HealthCheck.Timeout,
			HealthyThreshold:   targetGroup.HealthCheck.HealthyThreshold,
			UnhealthyThreshold: targetGroup.HealthCheck.UnhealthyThreshold,
		},
	}

	awsTGOutput, err := service.CreateTargetGroup(ctx, awsTG)
	if err != nil {
		return nil, fmt.Errorf("failed to create target group: %w", err)
	}

	// Convert AWS output to domain model
	return awsmapper.ToDomainTargetGroupFromOutput(awsTGOutput), nil
}

// createListener creates a Listener for the load balancer
func createListener(ctx context.Context, service *awscomputeservice.ComputeService, client *awssdk.AWSClient, loadBalancer *domaincompute.LoadBalancer, targetGroup *domaincompute.TargetGroup, useMockData bool) (*domaincompute.Listener, error) {
	if loadBalancer.ARN == nil {
		return nil, fmt.Errorf("load balancer ARN is required")
	}
	if targetGroup.ARN == nil {
		return nil, fmt.Errorf("target group ARN is required")
	}

	listener := &domaincompute.Listener{
		LoadBalancerARN: *loadBalancer.ARN,
		Port:            80,
		Protocol:        domaincompute.ListenerProtocolHTTP,
		DefaultAction: domaincompute.ListenerAction{
			Type:           domaincompute.ListenerActionTypeForward,
			TargetGroupARN: targetGroup.ARN,
		},
	}

	fmt.Printf("Creating listener on port %d...\n", listener.Port)

	if useMockData {
		fmt.Printf("  Using mock listener data (AWS account doesn't support load balancer creation)\n")
		mockListener := &domaincompute.Listener{
			ID:              fmt.Sprintf("%s/%s", *loadBalancer.ARN, generateMockID()),
			ARN:             stringPtr(fmt.Sprintf("arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/demo-alb/%s/%s", generateMockID()[:16], generateMockID())),
			LoadBalancerARN: listener.LoadBalancerARN,
			Port:            listener.Port,
			Protocol:        listener.Protocol,
			DefaultAction:   listener.DefaultAction,
		}
		return mockListener, nil
	}

	// Convert domain model to AWS model
	actionType := awsloadbalancer.ListenerActionTypeForward
	if listener.DefaultAction.Type == domaincompute.ListenerActionTypeRedirect {
		actionType = awsloadbalancer.ListenerActionTypeRedirect
	} else if listener.DefaultAction.Type == domaincompute.ListenerActionTypeFixedResponse {
		actionType = awsloadbalancer.ListenerActionTypeFixedResponse
	}

	awsListener := &awsloadbalancer.Listener{
		LoadBalancerARN: listener.LoadBalancerARN,
		Port:            listener.Port,
		Protocol:        string(listener.Protocol),
		DefaultAction: awsloadbalancer.ListenerAction{
			Type:           actionType,
			TargetGroupARN: listener.DefaultAction.TargetGroupARN,
		},
	}

	awsListenerOutput, err := service.CreateListener(ctx, awsListener)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	// Convert AWS output to domain model
	return awsmapper.ToDomainListenerFromOutput(awsListenerOutput), nil
}

// attachInstancesToTargetGroup attaches EC2 instances to the target group
func attachInstancesToTargetGroup(ctx context.Context, service *awscomputeservice.ComputeService, client *awssdk.AWSClient, targetGroup *domaincompute.TargetGroup, instances []*domaincompute.Instance, useMockData bool) error {
	if targetGroup.ARN == nil {
		return fmt.Errorf("target group ARN is required")
	}

	for i, instance := range instances {
		fmt.Printf("Attaching instance %d: %s...\n", i+1, instance.ID)

		if useMockData {
			fmt.Printf("  Using mock attachment (AWS account doesn't support load balancer creation)\n")
			fmt.Printf("  Instance %s attached successfully (simulated)\n", instance.ID)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		attachment := &awsloadbalancer.TargetGroupAttachment{
			TargetGroupARN: *targetGroup.ARN,
			TargetID:       instance.ID,
			Port:           intPtr(80), // Default HTTP port
		}

		err := service.AttachTargetToGroup(ctx, attachment)
		if err != nil {
			return fmt.Errorf("failed to attach instance %s: %w", instance.ID, err)
		}

		fmt.Printf("  Instance %s attached successfully\n", instance.ID)
		time.Sleep(200 * time.Millisecond)
	}

	return nil
}

// verifySetup verifies the complete load balancer setup
func verifySetup(ctx context.Context, service *awscomputeservice.ComputeService, client *awssdk.AWSClient, loadBalancer *domaincompute.LoadBalancer, targetGroup *domaincompute.TargetGroup, instances []*domaincompute.Instance, useMockData bool) error {
	fmt.Println("\nVerifying load balancer configuration...")

	if useMockData {
		// Mock verification - just print what we created
		fmt.Printf("  ✓ Load balancer verified: %s (State: %s) [MOCK]\n", loadBalancer.Name, loadBalancer.State)
		fmt.Printf("  ✓ Target group verified: %s (State: %s) [MOCK]\n", targetGroup.Name, targetGroup.State)
		fmt.Printf("  ✓ Found %d target(s) attached to target group [MOCK]\n", len(instances))
		for i, instance := range instances {
			fmt.Printf("    Target %d: %s (Health: healthy) [MOCK]\n", i+1, instance.ID)
		}
		if loadBalancer.ARN != nil {
			fmt.Printf("  ✓ Found 1 listener(s) on load balancer [MOCK]\n")
			fmt.Printf("    Listener 1: Port 80 (HTTP) [MOCK]\n")
		}
		return nil
	}

	// Verify load balancer
	if loadBalancer.ARN != nil {
		awsLBOutput, err := service.GetLoadBalancer(ctx, *loadBalancer.ARN)
		if err != nil {
			fmt.Printf("  Warning: Could not verify load balancer: %v\n", err)
		} else {
			lb := awsmapper.ToDomainLoadBalancerFromOutput(awsLBOutput)
			fmt.Printf("  ✓ Load balancer verified: %s (State: %s)\n", lb.Name, lb.State)
		}
	}

	// Verify target group
	if targetGroup.ARN != nil {
		awsTGOutput, err := service.GetTargetGroup(ctx, *targetGroup.ARN)
		if err != nil {
			fmt.Printf("  Warning: Could not verify target group: %v\n", err)
		} else {
			tg := awsmapper.ToDomainTargetGroupFromOutput(awsTGOutput)
			fmt.Printf("  ✓ Target group verified: %s (State: %s)\n", tg.Name, tg.State)
		}

		// List attached targets
		awsTargets, err := service.ListTargetGroupTargets(ctx, *targetGroup.ARN)
		if err != nil {
			fmt.Printf("  Warning: Could not list targets: %v\n", err)
		} else {
			fmt.Printf("  ✓ Found %d target(s) attached to target group\n", len(awsTargets))
			for i, awsTarget := range awsTargets {
				target := awsmapper.ToDomainTargetGroupAttachmentFromOutput(awsTarget)
				fmt.Printf("    Target %d: %s (Health: %s)\n", i+1, target.TargetID, target.HealthStatus)
			}
		}
	}

	// Verify listener
	if loadBalancer.ARN != nil {
		awsListeners, err := service.ListListeners(ctx, *loadBalancer.ARN)
		if err != nil {
			fmt.Printf("  Warning: Could not list listeners: %v\n", err)
		} else {
			fmt.Printf("  ✓ Found %d listener(s) on load balancer\n", len(awsListeners))
			for i, awsListener := range awsListeners {
				listener := awsmapper.ToDomainListenerFromOutput(awsListener)
				fmt.Printf("    Listener %d: Port %d (%s)\n", i+1, listener.Port, listener.Protocol)
			}
		}
	}

	return nil
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func generateMockID() string {
	// Generate a random-looking ID for mock resources
	return fmt.Sprintf("%016x", time.Now().UnixNano()%10000000000000000)
}
