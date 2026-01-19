package compute

import (
	domaincompute "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/compute"
	awsloadbalancer "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/compute/load_balancer/outputs"
)

// FromDomainTargetGroupAttachment converts domain TargetGroupAttachment to AWS TargetGroupAttachment
func FromDomainTargetGroupAttachment(domainAttachment *domaincompute.TargetGroupAttachment) *awsloadbalancer.TargetGroupAttachment {
	if domainAttachment == nil {
		return nil
	}

	awsAttachment := &awsloadbalancer.TargetGroupAttachment{
		TargetGroupARN:   domainAttachment.TargetGroupARN,
		TargetID:         domainAttachment.TargetID,
		Port:             domainAttachment.Port,
		AvailabilityZone: domainAttachment.AvailabilityZone,
	}

	return awsAttachment
}

// ToDomainTargetGroupAttachmentFromOutput converts AWS TargetGroupAttachment output to domain TargetGroupAttachment
func ToDomainTargetGroupAttachmentFromOutput(output *awsoutputs.TargetGroupAttachmentOutput) *domaincompute.TargetGroupAttachment {
	if output == nil {
		return nil
	}

	// Create composite ID
	id := output.TargetGroupARN + ":" + output.TargetID

	healthStatus := domaincompute.TargetHealthStatusHealthy
	switch output.HealthStatus {
	case "unhealthy":
		healthStatus = domaincompute.TargetHealthStatusUnhealthy
	case "initial":
		healthStatus = domaincompute.TargetHealthStatusInitial
	case "draining":
		healthStatus = domaincompute.TargetHealthStatusDraining
	}

	domainAttachment := &domaincompute.TargetGroupAttachment{
		ID:               id,
		TargetGroupARN:   output.TargetGroupARN,
		TargetID:         output.TargetID,
		Port:             output.Port,
		AvailabilityZone: output.AvailabilityZone,
		HealthStatus:     healthStatus,
	}

	return domainAttachment
}
