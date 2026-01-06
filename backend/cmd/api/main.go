package main

import (
	"fmt"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/networking"
	domainrules "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
	domainconstraints "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/constraints"
	registry "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/registry"
)
func main() {
	// new vpc
	vpc := &domainnetworking.VPC{
		Name: "test-vpc",
		Region: "us-east-1",
		CIDR: "10.0.0.0/16",

	}
	awsVPC := awsmapper.FromDomainVPC(vpc)

	if err := awsVPC.Validate(); err != nil {
		fmt.Println("Error validating vpc:", err)
		return
	}
	fmt.Println("CIDR", awsVPC.CIDR)
	fmt.Println("EnableDNSHostnames", awsVPC.EnableDNSHostnames)
	fmt.Println("EnableDNSSupport", awsVPC.EnableDNSSupport)
	fmt.Println("Region", awsVPC.Region)
	fmt.Println("Name", awsVPC.Name)
	fmt.Println("Tags", awsVPC.Tags[0].Key, awsVPC.Tags[0].Value)
	fmt.Println("InstanceTenancy", awsVPC.InstanceTenancy)
	fmt.Println("--------------------------------")
	fmt.Println("Registering rules")
	ruleRegistry := registry.NewRuleRegistry() // InMemoryRuleRegistry
	rule1 := domainconstraints.NewRequiresParentRule("VPC", "Subnet") // RequiresParentRule
	err := ruleRegistry.RegisterRule("VPC", rule1) // RegisterRule
	if err != nil {
		fmt.Println("Error registering rule:", err)
		return
	}
	rule2 := domainconstraints.NewRequiresRegionRule("VPC", true)
	err = ruleRegistry.RegisterRule("VPC", rule2)
	if err != nil {
		fmt.Println("Error registering rule:", err)
		return
	}
	rule3 := domainconstraints.NewAllowedParentRule("Subnet", []string{"VPC"})
	err = ruleRegistry.RegisterRule("Subnet", rule3)
	if err != nil {
		fmt.Println("Error registering rule:", err)
		return
	}
	subnetRules := ruleRegistry.GetRules("Subnet")
	for _, rule := range subnetRules {
		fmt.Println("Subnet rule", rule.GetType())
	}
	parentRuleType := domainrules.RuleTypeRequiresParent
	parentRules := ruleRegistry.GetRulesByType(parentRuleType)
	for _, rule := range parentRules {
		fmt.Println("Parent rule", rule.GetType())
	}

}
