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

// ListenerRunner demonstrates Listener operations
func ListenerRunner() {
	ctx := context.Background()

	fmt.Println("============================================")
	fmt.Println("LISTENER OPERATIONS")
	fmt.Println("============================================")

	// Initialize AWS client
	client, err := awssdk.NewAWSClient(ctx)
	if err != nil {
		fmt.Printf("Error creating AWS client: %v\n", err)
		return
	}

	// Initialize compute service
	computeService := awscomputeservice.NewComputeService(client)

	// Mock load balancer ARN (in real scenario, this would be from a created LB)
	loadBalancerARN := "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/demo-alb/1234567890abcdef"
	targetGroupARN := "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/demo-target-group/1234567890abcdef"

	// Create a listener
	fmt.Println("\n--- Creating Listener ---")
	listener := &domaincompute.Listener{
		LoadBalancerARN: loadBalancerARN,
		Port:            80,
		Protocol:        domaincompute.ListenerProtocolHTTP,
		DefaultAction: domaincompute.ListenerAction{
			Type:           domaincompute.ListenerActionTypeForward,
			TargetGroupARN: &targetGroupARN,
		},
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

	awsListenerOutput, err := computeService.CreateListener(ctx, awsListener)
	if err != nil {
		fmt.Printf("Error creating listener: %v\n", err)
		return
	}

	createdListener := awsmapper.ToDomainListenerFromOutput(awsListenerOutput)
	fmt.Printf("Listener created:\n")
	if createdListener.ARN != nil {
		fmt.Printf("  ARN: %s\n", *createdListener.ARN)
	}
	fmt.Printf("  Port: %d\n", createdListener.Port)
	fmt.Printf("  Protocol: %s\n", createdListener.Protocol)

	// List listeners
	fmt.Println("\n--- Listing Listeners ---")
	awsListeners, err := computeService.ListListeners(ctx, loadBalancerARN)
	if err != nil {
		fmt.Printf("Error listing listeners: %v\n", err)
		return
	}

	fmt.Printf("Found %d listener(s):\n", len(awsListeners))
	for i, awsListener := range awsListeners {
		l := awsmapper.ToDomainListenerFromOutput(awsListener)
		fmt.Printf("  %d. Port %d (%s)\n", i+1, l.Port, l.Protocol)
	}
}
