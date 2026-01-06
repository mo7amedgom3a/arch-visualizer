package networking

import (
	"errors"
	"fmt"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// SecurityGroupRule represents an AWS security group rule
type SecurityGroupRule struct {
	Type              string   `json:"type"` // "ingress" or "egress"
	Protocol          string   `json:"protocol"` // "tcp", "udp", "icmp", "-1" (all)
	FromPort          *int     `json:"from_port,omitempty"`
	ToPort            *int     `json:"to_port,omitempty"`
	CIDRBlocks        []string `json:"cidr_blocks,omitempty"`
	SourceSecurityGroupID *string `json:"source_security_group_id,omitempty"`
	Description       string   `json:"description"`
}

// SecurityGroup represents an AWS-specific Security Group
type SecurityGroup struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	VPCID       string `json:"vpc_id"`
	Rules       []SecurityGroupRule `json:"rules"`
	Tags        []configs.Tag `json:"tags"`
}

// Validate performs AWS-specific validation
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
	
	// AWS description length limit
	if len(sg.Description) > 255 {
		return errors.New("security group description must be 255 characters or less")
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
		
		// For ingress, require either CIDR blocks or source security group
		if rule.Type == "ingress" {
			if len(rule.CIDRBlocks) == 0 && (rule.SourceSecurityGroupID == nil || *rule.SourceSecurityGroupID == "") {
				return fmt.Errorf("rule %d: ingress rule requires cidr_blocks or source_security_group_id", i)
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
