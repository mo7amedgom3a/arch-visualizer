package networking

import (
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomainNetworkACL converts AWS Network ACL to domain Network ACL (for backward compatibility)
func ToDomainNetworkACL(awsACL *awsnetworking.NetworkACL) *domainnetworking.NetworkACL {
	if awsACL == nil {
		return nil
	}
	
	domainInboundRules := make([]domainnetworking.ACLRule, len(awsACL.InboundRules))
	for i, awsRule := range awsACL.InboundRules {
		domainInboundRules[i] = convertACLRule(awsRule)
	}
	
	domainOutboundRules := make([]domainnetworking.ACLRule, len(awsACL.OutboundRules))
	for i, awsRule := range awsACL.OutboundRules {
		domainOutboundRules[i] = convertACLRule(awsRule)
	}
	
	return &domainnetworking.NetworkACL{
		Name:          awsACL.Name,
		VPCID:         awsACL.VPCID,
		InboundRules:  domainInboundRules,
		OutboundRules: domainOutboundRules,
	}
}

// ToDomainNetworkACLFromOutput converts AWS Network ACL output to domain Network ACL with ID and ARN
func ToDomainNetworkACLFromOutput(output *awsoutputs.NetworkACLOutput) *domainnetworking.NetworkACL {
	if output == nil {
		return nil
	}
	
	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}
	
	domainInboundRules := make([]domainnetworking.ACLRule, len(output.InboundRules))
	for i, awsRule := range output.InboundRules {
		domainInboundRules[i] = convertACLRule(awsRule)
	}
	
	domainOutboundRules := make([]domainnetworking.ACLRule, len(output.OutboundRules))
	for i, awsRule := range output.OutboundRules {
		domainOutboundRules[i] = convertACLRule(awsRule)
	}
	
	// Extract subnet IDs from associations
	subnetIDs := make([]string, 0, len(output.Associations))
	for _, assoc := range output.Associations {
		subnetIDs = append(subnetIDs, assoc.SubnetID)
	}
	
	return &domainnetworking.NetworkACL{
		ID:            output.ID,
		ARN:           arn,
		Name:          output.Name,
		VPCID:         output.VPCID,
		IsDefault:     output.IsDefault,
		InboundRules:  domainInboundRules,
		OutboundRules: domainOutboundRules,
		Subnets:       subnetIDs,
	}
}

// FromDomainNetworkACL converts domain Network ACL to AWS Network ACL
func FromDomainNetworkACL(domainACL *domainnetworking.NetworkACL) *awsnetworking.NetworkACL {
	if domainACL == nil {
		return nil
	}
	
	awsInboundRules := make([]awsnetworking.ACLRule, len(domainACL.InboundRules))
	for i, domainRule := range domainACL.InboundRules {
		awsInboundRules[i] = convertDomainACLRule(domainRule)
	}
	
	awsOutboundRules := make([]awsnetworking.ACLRule, len(domainACL.OutboundRules))
	for i, domainRule := range domainACL.OutboundRules {
		awsOutboundRules[i] = convertDomainACLRule(domainRule)
	}
	
	return &awsnetworking.NetworkACL{
		Name:          domainACL.Name,
		VPCID:         domainACL.VPCID,
		InboundRules:  awsInboundRules,
		OutboundRules: awsOutboundRules,
		Tags:          []configs.Tag{{Key: "Name", Value: domainACL.Name}},
	}
}

// convertACLRule converts AWS ACL rule to domain ACL rule
func convertACLRule(awsRule awsnetworking.ACLRule) domainnetworking.ACLRule {
	domainRule := domainnetworking.ACLRule{
		RuleNumber: awsRule.RuleNumber,
		Type:       domainnetworking.ACLRuleType(awsRule.Type),
		Protocol:   awsRule.Protocol,
		CIDR:       awsRule.CIDR,
		Action:     domainnetworking.ACLRuleAction(awsRule.Action),
	}
	
	if awsRule.PortRange != nil {
		domainRule.PortRange = &domainnetworking.PortRange{
			From: awsRule.PortRange.From,
			To:   awsRule.PortRange.To,
		}
	}
	
	return domainRule
}

// convertDomainACLRule converts domain ACL rule to AWS ACL rule
func convertDomainACLRule(domainRule domainnetworking.ACLRule) awsnetworking.ACLRule {
	awsRule := awsnetworking.ACLRule{
		RuleNumber: domainRule.RuleNumber,
		Type:       awsnetworking.ACLRuleType(domainRule.Type),
		Protocol:   domainRule.Protocol,
		CIDR:       domainRule.CIDR,
		Action:     awsnetworking.ACLRuleAction(domainRule.Action),
	}
	
	if domainRule.PortRange != nil {
		awsRule.PortRange = &awsnetworking.PortRange{
			From: domainRule.PortRange.From,
			To:   domainRule.PortRange.To,
		}
	}
	
	return awsRule
}

// ConvertDomainACLRuleToAWS converts a domain ACL rule to AWS ACL rule
// This is a public function for use in adapters when adding rules
func ConvertDomainACLRuleToAWS(domainRule domainnetworking.ACLRule) awsnetworking.ACLRule {
	return convertDomainACLRule(domainRule)
}

// ConvertDomainACLRuleTypeToAWS converts a domain ACL rule type to AWS ACL rule type
func ConvertDomainACLRuleTypeToAWS(domainType domainnetworking.ACLRuleType) awsnetworking.ACLRuleType {
	return awsnetworking.ACLRuleType(domainType)
}
