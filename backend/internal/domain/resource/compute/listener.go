package compute

import (
	"errors"
)

// ListenerProtocol represents the protocol for a listener
type ListenerProtocol string

const (
	ListenerProtocolHTTP  ListenerProtocol = "HTTP"
	ListenerProtocolHTTPS ListenerProtocol = "HTTPS"
	ListenerProtocolTCP   ListenerProtocol = "TCP"
	ListenerProtocolTLS   ListenerProtocol = "TLS"
)

// ListenerActionType represents the type of listener action
type ListenerActionType string

const (
	ListenerActionTypeForward       ListenerActionType = "forward"
	ListenerActionTypeRedirect      ListenerActionType = "redirect"
	ListenerActionTypeFixedResponse ListenerActionType = "fixed-response"
)

// ListenerAction represents an action for a listener
type ListenerAction struct {
	Type              ListenerActionType
	TargetGroupARN    *string
	RedirectConfig    *RedirectConfig
	FixedResponseConfig *FixedResponseConfig
}

// RedirectConfig represents redirect configuration
type RedirectConfig struct {
	Protocol   *string // HTTP or HTTPS
	Port       *string // Port number or "443"
	StatusCode string  // HTTP_301 or HTTP_302
	Host       *string
	Path       *string
	Query      *string
}

// FixedResponseConfig represents fixed response configuration
type FixedResponseConfig struct {
	ContentType string // text/plain, text/css, text/html, application/json
	MessageBody *string
	StatusCode  string // 200-599
}

// ListenerRule represents a listener rule (for future use)
type ListenerRule struct {
	Priority  int
	Conditions []string // Path patterns, host headers, etc.
	Actions   []ListenerAction
}

// Listener represents a cloud-agnostic load balancer listener
type Listener struct {
	ID             string
	ARN            *string // Cloud-specific ARN
	LoadBalancerARN string
	Port           int
	Protocol       ListenerProtocol
	DefaultAction  ListenerAction
	Rules          []ListenerRule // Optional: for future use
}

// Validate performs domain-level validation
func (l *Listener) Validate() error {
	if l.LoadBalancerARN == "" {
		return errors.New("listener load balancer ARN is required")
	}
	if l.Port < 1 || l.Port > 65535 {
		return errors.New("listener port must be between 1 and 65535")
	}
	if l.Protocol == "" {
		return errors.New("listener protocol is required")
	}
	if l.Protocol != ListenerProtocolHTTP &&
		l.Protocol != ListenerProtocolHTTPS &&
		l.Protocol != ListenerProtocolTCP &&
		l.Protocol != ListenerProtocolTLS {
		return errors.New("listener protocol must be HTTP, HTTPS, TCP, or TLS")
	}
	if l.DefaultAction.Type == "" {
		return errors.New("listener default action type is required")
	}
	if l.DefaultAction.Type == ListenerActionTypeForward {
		if l.DefaultAction.TargetGroupARN == nil || *l.DefaultAction.TargetGroupARN == "" {
			return errors.New("target group ARN is required for forward action")
		}
	}
	return nil
}
