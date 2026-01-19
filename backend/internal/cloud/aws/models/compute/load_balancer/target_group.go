package load_balancer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// HealthCheckConfig represents health check configuration
type HealthCheckConfig struct {
	Path               *string `json:"path,omitempty"`                // Default: "/"
	Matcher            *string `json:"matcher,omitempty"`             // Default: "200"
	Interval           *int    `json:"interval,omitempty"`           // Default: 30 seconds
	Timeout            *int    `json:"timeout,omitempty"`            // Default: 5 seconds
	HealthyThreshold   *int    `json:"healthy_threshold,omitempty"`  // Default: 2
	UnhealthyThreshold *int    `json:"unhealthy_threshold,omitempty"` // Default: 2
	Protocol           *string `json:"protocol,omitempty"`           // Default: HTTP
	Port               *string `json:"port,omitempty"`               // Default: "traffic-port"
}

// TargetGroup represents an AWS Target Group configuration
type TargetGroup struct {
	Name       string        `json:"name"`        // Required
	Port       int           `json:"port"`        // Required
	Protocol   string        `json:"protocol"`   // Required: HTTP, HTTPS, TCP, TLS
	VPCID      string        `json:"vpc_id"`     // Required
	TargetType *string       `json:"target_type,omitempty"` // Optional: instance, ip, lambda (default: instance)
	HealthCheck HealthCheckConfig `json:"health_check"`
	Tags       []configs.Tag `json:"tags,omitempty"`
}

// Validate performs AWS-specific validation
func (tg *TargetGroup) Validate() error {
	if tg.Name == "" {
		return errors.New("target group name is required")
	}

	// Validate name format: 1-32 chars, alphanumeric and hyphens
	if len(tg.Name) < 1 || len(tg.Name) > 32 {
		return errors.New("target group name must be between 1 and 32 characters")
	}

	namePattern := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	if !namePattern.MatchString(tg.Name) {
		return errors.New("target group name can only contain alphanumeric characters and hyphens")
	}

	// Validate port
	if tg.Port < 1 || tg.Port > 65535 {
		return errors.New("target group port must be between 1 and 65535")
	}

	// Validate protocol
	if tg.Protocol == "" {
		return errors.New("target group protocol is required")
	}
	protocol := strings.ToUpper(tg.Protocol)
	if protocol != "HTTP" && protocol != "HTTPS" && protocol != "TCP" && protocol != "TLS" {
		return errors.New("target group protocol must be HTTP, HTTPS, TCP, or TLS")
	}

	// Validate VPC ID
	if tg.VPCID == "" {
		return errors.New("target group VPC ID is required")
	}
	if !strings.HasPrefix(tg.VPCID, "vpc-") {
		return errors.New("VPC ID must start with 'vpc-'")
	}

	// Validate target type if provided
	if tg.TargetType != nil {
		targetType := strings.ToLower(*tg.TargetType)
		if targetType != "instance" && targetType != "ip" && targetType != "lambda" {
			return errors.New("target type must be 'instance', 'ip', or 'lambda'")
		}
	}

	// Validate health check
	if err := tg.HealthCheck.Validate(); err != nil {
		return fmt.Errorf("health check validation failed: %w", err)
	}

	return nil
}

// Validate performs validation on health check configuration
func (hc *HealthCheckConfig) Validate() error {
	// Validate interval if provided
	if hc.Interval != nil {
		if *hc.Interval < 5 || *hc.Interval > 300 {
			return errors.New("health check interval must be between 5 and 300 seconds")
		}
	}

	// Validate timeout if provided
	if hc.Timeout != nil {
		if *hc.Timeout < 2 || *hc.Timeout > 120 {
			return errors.New("health check timeout must be between 2 and 120 seconds")
		}
		// Timeout must be less than interval
		if hc.Interval != nil && *hc.Timeout >= *hc.Interval {
			return errors.New("health check timeout must be less than interval")
		}
	}

	// Validate healthy threshold if provided
	if hc.HealthyThreshold != nil {
		if *hc.HealthyThreshold < 2 || *hc.HealthyThreshold > 10 {
			return errors.New("health check healthy threshold must be between 2 and 10")
		}
	}

	// Validate unhealthy threshold if provided
	if hc.UnhealthyThreshold != nil {
		if *hc.UnhealthyThreshold < 2 || *hc.UnhealthyThreshold > 10 {
			return errors.New("health check unhealthy threshold must be between 2 and 10")
		}
	}

	return nil
}
