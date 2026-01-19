package sdk

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer/outputs"
)

// CreateLoadBalancer creates a new Load Balancer using AWS SDK
func CreateLoadBalancer(ctx context.Context, client *AWSClient, lb *awsloadbalancer.LoadBalancer) (*awsoutputs.LoadBalancerOutput, error) {
	if err := lb.Validate(); err != nil {
		return nil, fmt.Errorf("load balancer validation failed: %w", err)
	}

	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	// Build CreateLoadBalancerInput
	input := &elasticloadbalancingv2.CreateLoadBalancerInput{
		Name:           aws.String(lb.Name),
		Type:           elbv2types.LoadBalancerTypeEnum(lb.LoadBalancerType),
		Subnets:        lb.SubnetIDs,
		SecurityGroups: lb.SecurityGroupIDs,
	}

	// Set internal flag
	if lb.Internal != nil && *lb.Internal {
		input.Scheme = elbv2types.LoadBalancerSchemeEnumInternal
	} else {
		input.Scheme = elbv2types.LoadBalancerSchemeEnumInternetFacing
	}

	// Set IP address type if provided
	if lb.IPAddressType != nil {
		input.IpAddressType = elbv2types.IpAddressType(*lb.IPAddressType)
	}

	// Add tags if provided
	if len(lb.Tags) > 0 {
		var tagList []elbv2types.Tag
		for _, tag := range lb.Tags {
			tagList = append(tagList, elbv2types.Tag{
				Key:   aws.String(tag.Key),
				Value: aws.String(tag.Value),
			})
		}
		input.Tags = tagList
	}

	// Create the load balancer
	result, err := client.ELBv2.CreateLoadBalancer(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create load balancer: %w", err)
	}

	if len(result.LoadBalancers) == 0 {
		return nil, fmt.Errorf("load balancer creation returned no load balancers")
	}

	return convertLoadBalancerToOutput(&result.LoadBalancers[0]), nil
}

// GetLoadBalancer retrieves a Load Balancer by ARN
func GetLoadBalancer(ctx context.Context, client *AWSClient, arn string) (*awsoutputs.LoadBalancerOutput, error) {
	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{
		LoadBalancerArns: []string{arn},
	}

	result, err := client.ELBv2.DescribeLoadBalancers(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer %s: %w", arn, err)
	}

	if len(result.LoadBalancers) == 0 {
		return nil, fmt.Errorf("load balancer %s not found", arn)
	}

	return convertLoadBalancerToOutput(&result.LoadBalancers[0]), nil
}

// UpdateLoadBalancer updates a Load Balancer
// Note: AWS ELBv2 has limited update capabilities (mainly security groups and subnets)
func UpdateLoadBalancer(ctx context.Context, client *AWSClient, arn string, lb *awsloadbalancer.LoadBalancer) (*awsoutputs.LoadBalancerOutput, error) {
	if err := lb.Validate(); err != nil {
		return nil, fmt.Errorf("load balancer validation failed: %w", err)
	}

	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	// Update security groups if provided
	if len(lb.SecurityGroupIDs) > 0 {
		setSecurityGroupsInput := &elasticloadbalancingv2.SetSecurityGroupsInput{
			LoadBalancerArn: aws.String(arn),
			SecurityGroups:  lb.SecurityGroupIDs,
		}

		_, err := client.ELBv2.SetSecurityGroups(ctx, setSecurityGroupsInput)
		if err != nil {
			return nil, fmt.Errorf("failed to update security groups: %w", err)
		}
	}

	// Update subnets if provided (add/remove subnets)
	if len(lb.SubnetIDs) > 0 {
		setSubnetsInput := &elasticloadbalancingv2.SetSubnetsInput{
			LoadBalancerArn: aws.String(arn),
			Subnets:         lb.SubnetIDs,
		}

		_, err := client.ELBv2.SetSubnets(ctx, setSubnetsInput)
		if err != nil {
			return nil, fmt.Errorf("failed to update subnets: %w", err)
		}
	}

	// Update IP address type if provided
	if lb.IPAddressType != nil {
		setIPAddressTypeInput := &elasticloadbalancingv2.SetIpAddressTypeInput{
			LoadBalancerArn: aws.String(arn),
			IpAddressType:   elbv2types.IpAddressType(*lb.IPAddressType),
		}

		_, err := client.ELBv2.SetIpAddressType(ctx, setIPAddressTypeInput)
		if err != nil {
			return nil, fmt.Errorf("failed to update IP address type: %w", err)
		}
	}

	// Return updated load balancer
	return GetLoadBalancer(ctx, client, arn)
}

// DeleteLoadBalancer deletes a Load Balancer
func DeleteLoadBalancer(ctx context.Context, client *AWSClient, arn string) error {
	if client == nil || client.ELBv2 == nil {
		return fmt.Errorf("AWS client not available")
	}

	input := &elasticloadbalancingv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(arn),
	}

	_, err := client.ELBv2.DeleteLoadBalancer(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete load balancer %s: %w", arn, err)
	}

	return nil
}

// ListLoadBalancers lists Load Balancers with optional filters
func ListLoadBalancers(ctx context.Context, client *AWSClient, filters map[string][]string) ([]*awsoutputs.LoadBalancerOutput, error) {
	if client == nil || client.ELBv2 == nil {
		return nil, fmt.Errorf("AWS client not available")
	}

	var allLoadBalancers []*awsoutputs.LoadBalancerOutput
	var nextToken *string

	for {
		input := &elasticloadbalancingv2.DescribeLoadBalancersInput{}

		// Apply filters if provided
		if names, ok := filters["name"]; ok && len(names) > 0 {
			input.Names = names
		}

		if nextToken != nil {
			input.Marker = nextToken
		}

		result, err := client.ELBv2.DescribeLoadBalancers(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list load balancers: %w", err)
		}

		// Convert each load balancer to output model
		for _, lb := range result.LoadBalancers {
			allLoadBalancers = append(allLoadBalancers, convertLoadBalancerToOutput(&lb))
		}

		// Check if there are more pages
		if result.NextMarker == nil {
			break
		}
		nextToken = result.NextMarker
	}

	return allLoadBalancers, nil
}

// convertLoadBalancerToOutput converts AWS SDK LoadBalancer to output model
func convertLoadBalancerToOutput(lb *elbv2types.LoadBalancer) *awsoutputs.LoadBalancerOutput {
	output := &awsoutputs.LoadBalancerOutput{
		ARN:             aws.ToString(lb.LoadBalancerArn),
		ID:              aws.ToString(lb.LoadBalancerArn),
		Name:            aws.ToString(lb.LoadBalancerName),
		DNSName:         aws.ToString(lb.DNSName),
		ZoneID:          aws.ToString(lb.CanonicalHostedZoneId),
		Type:            string(lb.Type),
		Internal:        lb.Scheme == elbv2types.LoadBalancerSchemeEnumInternal,
		State:           string(lb.State.Code),
		CreatedTime:     aws.ToTime(lb.CreatedTime),
	}

	// Convert security groups
	if len(lb.SecurityGroups) > 0 {
		output.SecurityGroupIDs = make([]string, len(lb.SecurityGroups))
		for i, sg := range lb.SecurityGroups {
			output.SecurityGroupIDs[i] = sg
		}
	}

	// Convert availability zones and subnets
	if len(lb.AvailabilityZones) > 0 {
		output.SubnetIDs = make([]string, 0)
		for _, az := range lb.AvailabilityZones {
			if az.SubnetId != nil {
				output.SubnetIDs = append(output.SubnetIDs, *az.SubnetId)
			}
		}
	}

	return output
}
