package networking

import (
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awsoutputs "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking/outputs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
)

// ToDomainSecurityGroup converts AWS Security Group to domain Security Group (for backward compatibility)
func ToDomainSecurityGroup(awsSG *awsnetworking.SecurityGroup) *domainnetworking.SecurityGroup {
	if awsSG == nil {
		return nil
	}
	
	domainRules := make([]domainnetworking.SecurityGroupRule, len(awsSG.Rules))
	for i, awsRule := range awsSG.Rules {
		domainRule := domainnetworking.SecurityGroupRule{
			Type:         awsRule.Type,
			Protocol:     domainnetworking.Protocol(awsRule.Protocol),
			FromPort:     awsRule.FromPort,
			ToPort:       awsRule.ToPort,
			CIDRBlocks:   awsRule.CIDRBlocks,
			Description:  awsRule.Description,
		}
		
		// Convert source security group ID to list
		if awsRule.SourceSecurityGroupID != nil && *awsRule.SourceSecurityGroupID != "" {
			domainRule.SourceGroupIDs = []string{*awsRule.SourceSecurityGroupID}
		}
		
		domainRules[i] = domainRule
	}
	
	return &domainnetworking.SecurityGroup{
		Name:        awsSG.Name,
		Description: awsSG.Description,
		VPCID:       awsSG.VPCID,
		Rules:       domainRules,
	}
}

// ToDomainSecurityGroupFromOutput converts AWS Security Group output to domain Security Group with ID and ARN
func ToDomainSecurityGroupFromOutput(output *awsoutputs.SecurityGroupOutput) *domainnetworking.SecurityGroup {
	if output == nil {
		return nil
	}
	
	arn := &output.ARN
	if output.ARN == "" {
		arn = nil
	}
	
	domainRules := make([]domainnetworking.SecurityGroupRule, len(output.Rules))
	for i, awsRule := range output.Rules {
		domainRule := domainnetworking.SecurityGroupRule{
			Type:         awsRule.Type,
			Protocol:     domainnetworking.Protocol(awsRule.Protocol),
			FromPort:     awsRule.FromPort,
			ToPort:       awsRule.ToPort,
			CIDRBlocks:   awsRule.CIDRBlocks,
			Description:  awsRule.Description,
		}
		
		// Convert source security group ID to list
		if awsRule.SourceSecurityGroupID != nil && *awsRule.SourceSecurityGroupID != "" {
			domainRule.SourceGroupIDs = []string{*awsRule.SourceSecurityGroupID}
		}
		
		domainRules[i] = domainRule
	}
	
	return &domainnetworking.SecurityGroup{
		ID:          output.ID,
		ARN:         arn,
		Name:        output.Name,
		Description: output.Description,
		VPCID:       output.VPCID,
		Rules:       domainRules,
	}
}

// FromDomainSecurityGroup converts domain Security Group to AWS Security Group
func FromDomainSecurityGroup(domainSG *domainnetworking.SecurityGroup) *awsnetworking.SecurityGroup {
	if domainSG == nil {
		return nil
	}
	
	awsRules := make([]awsnetworking.SecurityGroupRule, len(domainSG.Rules))
	for i, domainRule := range domainSG.Rules {
		awsRule := awsnetworking.SecurityGroupRule{
			Type:        domainRule.Type,
			Protocol:    string(domainRule.Protocol),
			FromPort:    domainRule.FromPort,
			ToPort:      domainRule.ToPort,
			CIDRBlocks:  domainRule.CIDRBlocks,
			Description: domainRule.Description,
		}
		
		// Convert source group IDs to single source security group ID (AWS limitation)
		if len(domainRule.SourceGroupIDs) > 0 {
			awsRule.SourceSecurityGroupID = &domainRule.SourceGroupIDs[0]
		}
		
		awsRules[i] = awsRule
	}
	
	return &awsnetworking.SecurityGroup{
		Name:        domainSG.Name,
		Description: domainSG.Description,
		VPCID:       domainSG.VPCID,
		Rules:       awsRules,
		Tags:        []configs.Tag{{Key: "Name", Value: domainSG.Name}},
	}
}
