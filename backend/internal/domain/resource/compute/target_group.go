package compute

import (
	"errors"
)

// TargetGroupProtocol represents the protocol for a target group
type TargetGroupProtocol string

const (
	TargetGroupProtocolHTTP  TargetGroupProtocol = "HTTP"
	TargetGroupProtocolHTTPS TargetGroupProtocol = "HTTPS"
	TargetGroupProtocolTCP   TargetGroupProtocol = "TCP"
	TargetGroupProtocolTLS   TargetGroupProtocol = "TLS"
)

// TargetType represents the type of target
type TargetType string

const (
	TargetTypeInstance TargetType = "instance"
	TargetTypeIP       TargetType = "ip"
	TargetTypeLambda   TargetType = "lambda"
)

// TargetGroupState represents the state of a target group
type TargetGroupState string

const (
	TargetGroupStateActive         TargetGroupState = "active"
	TargetGroupStateDraining       TargetGroupState = "draining"
	TargetGroupStateDeleting       TargetGroupState = "deleting"
	TargetGroupStateDeleted        TargetGroupState = "deleted"
)

// HealthCheckConfig represents health check configuration
type HealthCheckConfig struct {
	Path                *string // Default: "/"
	Matcher             *string // Default: "200"
	Interval            *int    // Default: 30 seconds
	Timeout             *int    // Default: 5 seconds
	HealthyThreshold    *int    // Default: 2
	UnhealthyThreshold  *int    // Default: 2
	Protocol            *string // Default: HTTP
	Port                *string // Default: "traffic-port"
}

// TargetGroup represents a cloud-agnostic target group
type TargetGroup struct {
	ID         string
	ARN        *string // Cloud-specific ARN
	Name       string
	VPCID      string
	Port       int
	Protocol   TargetGroupProtocol
	TargetType TargetType
	HealthCheck HealthCheckConfig
	State      TargetGroupState
}

// Validate performs domain-level validation
func (tg *TargetGroup) Validate() error {
	if tg.Name == "" {
		return errors.New("target group name is required")
	}
	if tg.VPCID == "" {
		return errors.New("target group VPC ID is required")
	}
	if tg.Port < 1 || tg.Port > 65535 {
		return errors.New("target group port must be between 1 and 65535")
	}
	if tg.Protocol == "" {
		return errors.New("target group protocol is required")
	}
	if tg.Protocol != TargetGroupProtocolHTTP &&
		tg.Protocol != TargetGroupProtocolHTTPS &&
		tg.Protocol != TargetGroupProtocolTCP &&
		tg.Protocol != TargetGroupProtocolTLS {
		return errors.New("target group protocol must be HTTP, HTTPS, TCP, or TLS")
	}
	if tg.TargetType == "" {
		tg.TargetType = TargetTypeInstance // Default
	}
	return nil
}
