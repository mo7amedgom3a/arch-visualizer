package compute

import (
	"errors"
)

// TargetHealthStatus represents the health status of a target
type TargetHealthStatus string

const (
	TargetHealthStatusHealthy   TargetHealthStatus = "healthy"
	TargetHealthStatusUnhealthy TargetHealthStatus = "unhealthy"
	TargetHealthStatusInitial   TargetHealthStatus = "initial"
	TargetHealthStatusDraining  TargetHealthStatus = "draining"
)

// TargetGroupAttachment represents a cloud-agnostic target group attachment
type TargetGroupAttachment struct {
	ID               string // Composite: target_group_arn + target_id
	TargetGroupARN   string
	TargetID         string // EC2 instance ID or IP address
	Port             *int   // Optional port override
	AvailabilityZone *string
	HealthStatus     TargetHealthStatus
}

// Validate performs domain-level validation
func (tga *TargetGroupAttachment) Validate() error {
	if tga.TargetGroupARN == "" {
		return errors.New("target group ARN is required")
	}
	if tga.TargetID == "" {
		return errors.New("target ID is required")
	}
	if tga.Port != nil {
		if *tga.Port < 1 || *tga.Port > 65535 {
			return errors.New("port must be between 1 and 65535")
		}
	}
	return nil
}
