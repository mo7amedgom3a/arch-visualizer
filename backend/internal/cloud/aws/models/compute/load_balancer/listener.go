package load_balancer

import (
	"errors"
	"fmt"
	"strings"
)

// ListenerActionType represents the type of listener action
type ListenerActionType string

const (
	ListenerActionTypeForward       ListenerActionType = "forward"
	ListenerActionTypeRedirect      ListenerActionType = "redirect"
	ListenerActionTypeFixedResponse ListenerActionType = "fixed-response"
)

// RedirectConfig represents redirect configuration
type RedirectConfig struct {
	Protocol   *string `json:"protocol,omitempty"`   // HTTP or HTTPS
	Port       *string `json:"port,omitempty"`        // Port number or "443"
	StatusCode string  `json:"status_code"`          // HTTP_301 or HTTP_302
	Host       *string `json:"host,omitempty"`
	Path       *string `json:"path,omitempty"`
	Query      *string `json:"query,omitempty"`
}

// FixedResponseConfig represents fixed response configuration
type FixedResponseConfig struct {
	ContentType string  `json:"content_type"`         // text/plain, text/css, text/html, application/json
	MessageBody *string `json:"message_body,omitempty"`
	StatusCode  string  `json:"status_code"`          // 200-599
}

// ListenerAction represents an action for a listener
type ListenerAction struct {
	Type               ListenerActionType   `json:"type"` // forward, redirect, fixed-response
	TargetGroupARN     *string              `json:"target_group_arn,omitempty"` // Required if type is forward
	RedirectConfig     *RedirectConfig       `json:"redirect_config,omitempty"` // Required if type is redirect
	FixedResponseConfig *FixedResponseConfig `json:"fixed_response_config,omitempty"` // Required if type is fixed-response
}

// Listener represents an AWS Load Balancer Listener configuration
type Listener struct {
	LoadBalancerARN string         `json:"load_balancer_arn"` // Required
	Port            int            `json:"port"`              // Required
	Protocol        string         `json:"protocol"`          // Required: HTTP, HTTPS, TCP, TLS
	DefaultAction   ListenerAction `json:"default_action"`    // Required
	CertificateARN  *string        `json:"certificate_arn,omitempty"` // Required for HTTPS/TLS
	SSLPolicy       *string        `json:"ssl_policy,omitempty"`      // Optional for HTTPS/TLS
}

// Validate performs AWS-specific validation
func (l *Listener) Validate() error {
	if l.LoadBalancerARN == "" {
		return errors.New("listener load balancer ARN is required")
	}
	if !strings.HasPrefix(l.LoadBalancerARN, "arn:aws:elasticloadbalancing:") {
		return errors.New("invalid load balancer ARN format")
	}

	// Validate port
	if l.Port < 1 || l.Port > 65535 {
		return errors.New("listener port must be between 1 and 65535")
	}

	// Validate protocol
	if l.Protocol == "" {
		return errors.New("listener protocol is required")
	}
	protocol := strings.ToUpper(l.Protocol)
	if protocol != "HTTP" && protocol != "HTTPS" && protocol != "TCP" && protocol != "TLS" {
		return errors.New("listener protocol must be HTTP, HTTPS, TCP, or TLS")
	}

	// Validate certificate ARN for HTTPS/TLS
	if (protocol == "HTTPS" || protocol == "TLS") && (l.CertificateARN == nil || *l.CertificateARN == "") {
		return errors.New("certificate ARN is required for HTTPS/TLS listeners")
	}

	// Validate default action
	if err := l.DefaultAction.Validate(); err != nil {
		return fmt.Errorf("default action validation failed: %w", err)
	}

	return nil
}

// Validate performs validation on listener action
func (la *ListenerAction) Validate() error {
	if la.Type == "" {
		return errors.New("listener action type is required")
	}

	switch la.Type {
	case ListenerActionTypeForward:
		if la.TargetGroupARN == nil || *la.TargetGroupARN == "" {
			return errors.New("target group ARN is required for forward action")
		}
		if !strings.HasPrefix(*la.TargetGroupARN, "arn:aws:elasticloadbalancing:") {
			return errors.New("invalid target group ARN format")
		}
	case ListenerActionTypeRedirect:
		if la.RedirectConfig == nil {
			return errors.New("redirect config is required for redirect action")
		}
		if err := la.RedirectConfig.Validate(); err != nil {
			return fmt.Errorf("redirect config validation failed: %w", err)
		}
	case ListenerActionTypeFixedResponse:
		if la.FixedResponseConfig == nil {
			return errors.New("fixed response config is required for fixed-response action")
		}
		if err := la.FixedResponseConfig.Validate(); err != nil {
			return fmt.Errorf("fixed response config validation failed: %w", err)
		}
	default:
		return fmt.Errorf("invalid listener action type: %s", la.Type)
	}

	return nil
}

// Validate performs validation on redirect config
func (rc *RedirectConfig) Validate() error {
	if rc.StatusCode == "" {
		return errors.New("redirect status code is required")
	}
	if rc.StatusCode != "HTTP_301" && rc.StatusCode != "HTTP_302" {
		return errors.New("redirect status code must be HTTP_301 or HTTP_302")
	}
	return nil
}

// Validate performs validation on fixed response config
func (frc *FixedResponseConfig) Validate() error {
	if frc.ContentType == "" {
		return errors.New("fixed response content type is required")
	}
	validContentTypes := []string{"text/plain", "text/css", "text/html", "application/json"}
	valid := false
	for _, ct := range validContentTypes {
		if frc.ContentType == ct {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("fixed response content type must be one of: %v", validContentTypes)
	}

	if frc.StatusCode == "" {
		return errors.New("fixed response status code is required")
	}
	// Status code should be 200-599
	if len(frc.StatusCode) != 3 {
		return errors.New("fixed response status code must be 3 digits")
	}
	return nil
}
