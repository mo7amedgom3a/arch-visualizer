package load_balancer

import (
	"errors"
	"strings"
)

// TargetGroupAttachment represents an AWS Target Group Attachment configuration
type TargetGroupAttachment struct {
	TargetGroupARN   string  `json:"target_group_arn"`   // Required
	TargetID         string  `json:"target_id"`          // Required (instance ID or IP)
	Port             *int    `json:"port,omitempty"`     // Optional
	AvailabilityZone *string `json:"availability_zone,omitempty"` // Optional
}

// Validate performs AWS-specific validation
func (tga *TargetGroupAttachment) Validate() error {
	if tga.TargetGroupARN == "" {
		return errors.New("target group ARN is required")
	}
	if !strings.HasPrefix(tga.TargetGroupARN, "arn:aws:elasticloadbalancing:") {
		return errors.New("invalid target group ARN format")
	}

	if tga.TargetID == "" {
		return errors.New("target ID is required")
	}

	// Validate port if provided
	if tga.Port != nil {
		if *tga.Port < 1 || *tga.Port > 65535 {
			return errors.New("port must be between 1 and 65535")
		}
	}

	return nil
}
