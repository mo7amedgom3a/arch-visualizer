package networking

import (
	"errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ACLRuleAction represents the action for an ACL rule
type ACLRuleAction string

const (
	ACLRuleActionAllow ACLRuleAction = "allow"
	ACLRuleActionDeny  ACLRuleAction = "deny"
)

// ACLRuleType represents the type of ACL rule
type ACLRuleType string

const (
	ACLRuleTypeIngress ACLRuleType = "ingress" // Inbound
	ACLRuleTypeEgress  ACLRuleType = "egress"  // Outbound
)

// ACLRule represents a single Network ACL rule
type ACLRule struct {
	RuleNumber int           `json:"rule_number"` // Evaluated by rule number (lowest first, 1-32766)
	Type       ACLRuleType   `json:"type"`        // ingress or egress
	Protocol   string        `json:"protocol"`    // tcp, udp, icmp, -1 (all)
	PortRange  *PortRange    `json:"port_range,omitempty"` // Optional port range
	CIDR       string        `json:"cidr"`        // Source (for ingress) or Destination (for egress)
	Action     ACLRuleAction `json:"action"`      // allow or deny
}

// PortRange represents a port range
type PortRange struct {
	From *int `json:"from,omitempty"` // Starting port (nil means all ports)
	To   *int `json:"to,omitempty"`   // Ending port (nil means all ports)
}

// NetworkACL represents an AWS-specific Network ACL
// Network ACLs are stateless and operate at the subnet level
type NetworkACL struct {
	Name          string     `json:"name"`
	VPCID         string     `json:"vpc_id"`
	InboundRules  []ACLRule  `json:"inbound_rules,omitempty"`  // Ingress rules
	OutboundRules []ACLRule  `json:"outbound_rules,omitempty"`  // Egress rules
	Tags          []configs.Tag `json:"tags,omitempty"`
}

// Validate performs AWS-specific validation
func (acl *NetworkACL) Validate() error {
	if acl.Name == "" {
		return errors.New("network acl name is required")
	}
	if acl.VPCID == "" {
		return errors.New("network acl vpc_id is required")
	}

	// Validate inbound rules
	for i, rule := range acl.InboundRules {
		if err := validateACLRule(rule, ACLRuleTypeIngress, i); err != nil {
			return err
		}
	}

	// Validate outbound rules
	for i, rule := range acl.OutboundRules {
		if err := validateACLRule(rule, ACLRuleTypeEgress, i); err != nil {
			return err
		}
	}

	return nil
}

// validateACLRule validates a single ACL rule
func validateACLRule(rule ACLRule, expectedType ACLRuleType, index int) error {
	if rule.Type != expectedType {
		return errors.New("rule type mismatch")
	}

	// AWS rule number range: 1-32766
	if rule.RuleNumber < 1 || rule.RuleNumber > 32766 {
		return errors.New("rule number must be between 1 and 32766")
	}

	if rule.Protocol == "" {
		return errors.New("protocol is required")
	}

	if rule.CIDR == "" {
		return errors.New("cidr is required")
	}

	if rule.Action != ACLRuleActionAllow && rule.Action != ACLRuleActionDeny {
		return errors.New("action must be 'allow' or 'deny'")
	}

	// Validate port range if provided
	if rule.PortRange != nil {
		if rule.PortRange.From != nil && rule.PortRange.To != nil {
			if *rule.PortRange.From < 0 || *rule.PortRange.From > 65535 {
				return errors.New("port range 'from' must be between 0 and 65535")
			}
			if *rule.PortRange.To < 0 || *rule.PortRange.To > 65535 {
				return errors.New("port range 'to' must be between 0 and 65535")
			}
			if *rule.PortRange.From > *rule.PortRange.To {
				return errors.New("port range 'from' must be <= 'to'")
			}
		}
	}

	return nil
}
