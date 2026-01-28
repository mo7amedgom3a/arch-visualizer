package compute

import (
	"strings"

	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// FromDomainLoadBalancer converts domain LoadBalancer to AWS LoadBalancer
func FromDomainLoadBalancer(domainLB *domaincompute.LoadBalancer) *awsloadbalancer.LoadBalancer {
	if domainLB == nil {
		return nil
	}

	internal := domainLB.Internal
	lbType := "application"
	if domainLB.Type == domaincompute.LoadBalancerTypeNetwork {
		lbType = "network"
	}

	awsLB := &awsloadbalancer.LoadBalancer{
		Name:             domainLB.Name,
		LoadBalancerType: lbType,
		Internal:         &internal,
		SecurityGroupIDs: domainLB.SecurityGroupIDs,
		SubnetIDs:        domainLB.SubnetIDs,
	}

	// Convert tags if any (domain doesn't have tags, but we keep structure consistent)
	awsLB.Tags = []configs.Tag{}

	return awsLB
}

// ToDomainLoadBalancerFromOutput converts AWS LoadBalancer output to domain LoadBalancer
func ToDomainLoadBalancerFromOutput(output *awsoutputs.LoadBalancerOutput) *domaincompute.LoadBalancer {
	if output == nil {
		return nil
	}

	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}

	dnsName := &output.DNSName
	if output.DNSName == "" {
		dnsName = nil
	}

	zoneID := &output.ZoneID
	if output.ZoneID == "" {
		zoneID = nil
	}

	lbType := domaincompute.LoadBalancerTypeApplication
	if output.Type == "network" {
		lbType = domaincompute.LoadBalancerTypeNetwork
	}

	state := domaincompute.LoadBalancerStateActive
	switch output.State {
	case "provisioning":
		state = domaincompute.LoadBalancerStateProvisioning
	case "active_impaired":
		state = domaincompute.LoadBalancerStateActiveImpaired
	case "failed":
		state = domaincompute.LoadBalancerStateFailed
	}

	// Extract region from ARN if available
	// ARN format: arn:aws:elasticloadbalancing:REGION:ACCOUNT:loadbalancer/...
	region := ""
	if output.ARN != "" {
		parts := strings.Split(output.ARN, ":")
		if len(parts) >= 4 {
			region = parts[3]
		}
	}

	domainLB := &domaincompute.LoadBalancer{
		ID:               output.ID,
		ARN:              arn,
		Name:             output.Name,
		Region:           region,
		Type:             lbType,
		Internal:         output.Internal,
		SecurityGroupIDs: output.SecurityGroupIDs,
		SubnetIDs:        output.SubnetIDs,
		DNSName:          dnsName,
		ZoneID:           zoneID,
		State:            state,
	}

	return domainLB
}

// ToDomainLoadBalancerOutputFromOutput converts AWS LoadBalancer output directly to domain LoadBalancerOutput
func ToDomainLoadBalancerOutputFromOutput(output *awsoutputs.LoadBalancerOutput) *domaincompute.LoadBalancerOutput {
	if output == nil {
		return nil
	}

	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}

	dnsName := &output.DNSName
	if output.DNSName == "" {
		dnsName = nil
	}

	zoneID := &output.ZoneID
	if output.ZoneID == "" {
		zoneID = nil
	}

	lbType := domaincompute.LoadBalancerTypeApplication
	if output.Type == "network" {
		lbType = domaincompute.LoadBalancerTypeNetwork
	}

	state := domaincompute.LoadBalancerStateActive
	switch output.State {
	case "provisioning":
		state = domaincompute.LoadBalancerStateProvisioning
	case "active_impaired":
		state = domaincompute.LoadBalancerStateActiveImpaired
	case "failed":
		state = domaincompute.LoadBalancerStateFailed
	}

	// Extract region from ARN if available
	region := ""
	if output.ARN != "" {
		parts := strings.Split(output.ARN, ":")
		if len(parts) >= 4 {
			region = parts[3]
		}
	}

	createdAt := &output.CreatedTime

	return &domaincompute.LoadBalancerOutput{
		ID:               output.ID,
		ARN:              arn,
		Name:             output.Name,
		Region:           region,
		Type:             lbType,
		Internal:         output.Internal,
		SecurityGroupIDs: output.SecurityGroupIDs,
		SubnetIDs:        output.SubnetIDs,
		DNSName:          dnsName,
		ZoneID:           zoneID,
		State:            state,
		CreatedAt:        createdAt,
	}
}
