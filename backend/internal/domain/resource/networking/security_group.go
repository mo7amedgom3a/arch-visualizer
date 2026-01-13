package networking

import (
	"errors"
	"fmt"
	"strconv"
)

// Protocol represents network protocol
type Protocol string

const (
	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
	ProtocolICMP Protocol = "icmp"
	ProtocolAll  Protocol = "-1" // All protocols
)

// SecurityGroupRule represents a single security group rule
type SecurityGroupRule struct {
	Type            string   // "ingress" or "egress"
	Protocol       Protocol
	FromPort       *int     // nil means all ports
	ToPort         *int     // nil means all ports
	CIDRBlocks     []string // Source/destination CIDR blocks
	SourceGroupIDs []string // Source security group IDs (for ingress)
	Description    string
}

// SecurityGroup represents a cloud-agnostic security group
type SecurityGroup struct {
	ID          string
	ARN         *string // Cloud-specific ARN (AWS, Azure, etc.) - optional
	Name        string
	Description string
	VPCID       string
	Rules       []SecurityGroupRule
}

// Validate performs domain-level validation
func (sg *SecurityGroup) Validate() error {
	if sg.Name == "" {
		return errors.New("security group name is required")
	}
	if sg.Description == "" {
		return errors.New("security group description is required")
	}
	if sg.VPCID == "" {
		return errors.New("security group vpc_id is required")
	}
	
	// Validate rules
	for i, rule := range sg.Rules {
		if rule.Type != "ingress" && rule.Type != "egress" {
			return fmt.Errorf("rule %d: type must be 'ingress' or 'egress'", i)
		}
		
		if rule.Protocol == "" {
			return fmt.Errorf("rule %d: protocol is required", i)
		}
		
		// Validate port range
		if rule.FromPort != nil && rule.ToPort != nil {
			if *rule.FromPort < 0 || *rule.FromPort > 65535 {
				return fmt.Errorf("rule %d: from_port must be between 0 and 65535", i)
			}
			if *rule.ToPort < 0 || *rule.ToPort > 65535 {
				return fmt.Errorf("rule %d: to_port must be between 0 and 65535", i)
			}
			if *rule.FromPort > *rule.ToPort {
				return fmt.Errorf("rule %d: from_port must be <= to_port", i)
			}
		}
		
		// For ingress, require either CIDR blocks or source group IDs
		if rule.Type == "ingress" {
			if len(rule.CIDRBlocks) == 0 && len(rule.SourceGroupIDs) == 0 {
				return fmt.Errorf("rule %d: ingress rule requires cidr_blocks or source_group_ids", i)
			}
		}
		
		// For egress, require CIDR blocks
		if rule.Type == "egress" {
			if len(rule.CIDRBlocks) == 0 {
				return fmt.Errorf("rule %d: egress rule requires cidr_blocks", i)
			}
		}
	}
	
	return nil
}

// AddIngressRule adds an ingress rule to the security group
func (sg *SecurityGroup) AddIngressRule(protocol Protocol, fromPort, toPort *int, cidrBlocks []string, description string) {
	sg.Rules = append(sg.Rules, SecurityGroupRule{
		Type:         "ingress",
		Protocol:     protocol,
		FromPort:     fromPort,
		ToPort:       toPort,
		CIDRBlocks:   cidrBlocks,
		Description:  description,
	})
}

// AddEgressRule adds an egress rule to the security group
func (sg *SecurityGroup) AddEgressRule(protocol Protocol, fromPort, toPort *int, cidrBlocks []string, description string) {
	sg.Rules = append(sg.Rules, SecurityGroupRule{
		Type:        "egress",
		Protocol:    protocol,
		FromPort:    fromPort,
		ToPort:      toPort,
		CIDRBlocks:  cidrBlocks,
		Description: description,
	})
}

// ParsePort parses a port string to int
func ParsePort(portStr string) (*int, error) {
	if portStr == "" || portStr == "all" {
		return nil, nil
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}
	return &port, nil
}
