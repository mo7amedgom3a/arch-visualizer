package networking
import (
	"fmt"
	awsmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/mapper/networking"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/networking"
	awspricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing"
	networkingpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	domainnetworking "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource/networking"
	domainrules "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules"
	domainconstraints "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/constraints"
	registry "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/rules/registry"
	"context"
	"time"
)

func NetworkingRunner() {
	// new vpc
	vpc := &domainnetworking.VPC{
		Name:   "test-vpc",
		Region: "us-east-1",
		CIDR:   "10.0.0.0/16",
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
	ruleRegistry := registry.NewRuleRegistry()                        // InMemoryRuleRegistry
	rule1 := domainconstraints.NewRequiresParentRule("VPC", "Subnet") // RequiresParentRule
	err := ruleRegistry.RegisterRule("VPC", rule1)                    // RegisterRule
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
	// ============================================
	// PRICING FEATURE TESTS
	// ============================================
	fmt.Println("\n============================================")
	fmt.Println("PRICING FEATURE TESTS")
	fmt.Println("============================================")

	ctx := context.Background()
	pricingService := awspricing.NewAWSPricingService()

	// Test 1: Get Pricing Information for Resources
	fmt.Println("\n--- Test 1: Get Pricing Information ---")

	// Get NAT Gateway pricing
	natPricing, err := pricingService.GetPricing(ctx, "nat_gateway", "aws", "us-east-1")
	if err != nil {
		fmt.Printf("Error getting NAT Gateway pricing: %v\n", err)
	} else {
		fmt.Printf("NAT Gateway Pricing:\n")
		fmt.Printf("  Resource Type: %s\n", natPricing.ResourceType)
		fmt.Printf("  Provider: %s\n", natPricing.Provider)
		fmt.Printf("  Components:\n")
		for _, comp := range natPricing.Components {
			fmt.Printf("    - %s: $%.3f/%s\n", comp.Name, comp.Rate, comp.Unit)
		}
	}

	// Get Elastic IP pricing
	eipPricing, err := pricingService.GetPricing(ctx, "elastic_ip", "aws", "us-east-1")
	if err != nil {
		fmt.Printf("Error getting Elastic IP pricing: %v\n", err)
	} else {
		fmt.Printf("\nElastic IP Pricing:\n")
		fmt.Printf("  Resource Type: %s\n", eipPricing.ResourceType)
		for _, comp := range eipPricing.Components {
			fmt.Printf("    - %s: $%.3f/%s\n", comp.Name, comp.Rate, comp.Unit)
		}
	}

	// Get Network Interface pricing
	eniPricing, err := pricingService.GetPricing(ctx, "network_interface", "aws", "us-east-1")
	if err != nil {
		fmt.Printf("Error getting Network Interface pricing: %v\n", err)
	} else {
		fmt.Printf("\nNetwork Interface Pricing:\n")
		fmt.Printf("  Resource Type: %s\n", eniPricing.ResourceType)
		for _, comp := range eniPricing.Components {
			fmt.Printf("    - %s: $%.3f/%s\n", comp.Name, comp.Rate, comp.Unit)
		}
	}

	// Get Data Transfer pricing
	dtPricing, err := pricingService.GetPricing(ctx, "data_transfer", "aws", "us-east-1")
	if err != nil {
		fmt.Printf("Error getting Data Transfer pricing: %v\n", err)
	} else {
		fmt.Printf("\nData Transfer Pricing:\n")
		fmt.Printf("  Resource Type: %s\n", dtPricing.ResourceType)
		for _, comp := range dtPricing.Components {
			fmt.Printf("    - %s: $%.3f/%s\n", comp.Name, comp.Rate, comp.Unit)
		}
	}

	// Test 2: Calculate Costs Using Networking Functions
	fmt.Println("\n--- Test 2: Calculate Costs (Direct Functions) ---")

	// NAT Gateway cost for 1 hour with no data
	natCost1h := networkingpricing.CalculateNATGatewayCost(1*time.Hour, 0.0, "us-east-1")
	fmt.Printf("NAT Gateway (1 hour, no data): $%.3f\n", natCost1h)

	// NAT Gateway cost for 1 hour with 100GB data
	natCost1hData := networkingpricing.CalculateNATGatewayCost(1*time.Hour, 100.0, "us-east-1")
	fmt.Printf("NAT Gateway (1 hour, 100GB data): $%.3f\n", natCost1hData)

	// NAT Gateway cost for 30 days (720 hours)
	natCost30d := networkingpricing.CalculateNATGatewayCost(720*time.Hour, 0.0, "us-east-1")
	fmt.Printf("NAT Gateway (30 days, no data): $%.2f\n", natCost30d)

	// Elastic IP cost (unattached) for 1 hour
	eipCost1h := networkingpricing.CalculateElasticIPCost(1*time.Hour, false, "us-east-1")
	fmt.Printf("Elastic IP (unattached, 1 hour): $%.3f\n", eipCost1h)

	// Elastic IP cost (attached) for 1 hour (should be free)
	eipCostAttached := networkingpricing.CalculateElasticIPCost(1*time.Hour, true, "us-east-1")
	fmt.Printf("Elastic IP (attached, 1 hour): $%.3f (free)\n", eipCostAttached)

	// Elastic IP cost for 30 days (unattached)
	eipCost30d := networkingpricing.CalculateElasticIPCost(720*time.Hour, false, "us-east-1")
	fmt.Printf("Elastic IP (unattached, 30 days): $%.2f\n", eipCost30d)

	// Network Interface cost (unattached) for 1 hour
	eniCost1h := networkingpricing.CalculateNetworkInterfaceCost(1*time.Hour, false, "us-east-1")
	fmt.Printf("Network Interface (unattached, 1 hour): $%.2f\n", eniCost1h)

	// Network Interface cost (attached) for 1 hour (should be free)
	eniCostAttached := networkingpricing.CalculateNetworkInterfaceCost(1*time.Hour, true, "us-east-1")
	fmt.Printf("Network Interface (attached, 1 hour): $%.2f (free)\n", eniCostAttached)

	// Network Interface cost for 30 days (unattached)
	eniCost30d := networkingpricing.CalculateNetworkInterfaceCost(720*time.Hour, false, "us-east-1")
	fmt.Printf("Network Interface (unattached, 30 days): $%.2f\n", eniCost30d)

	// Data Transfer costs
	dtInbound := networkingpricing.CalculateDataTransferCost(100.0, networkingpricing.Inbound, "us-east-1")
	fmt.Printf("Data Transfer (100GB inbound): $%.2f (free)\n", dtInbound)

	dtOutboundSmall := networkingpricing.CalculateDataTransferCost(0.5, networkingpricing.Outbound, "us-east-1")
	fmt.Printf("Data Transfer (0.5GB outbound): $%.2f (within free tier)\n", dtOutboundSmall)

	dtOutboundLarge := networkingpricing.CalculateDataTransferCost(100.0, networkingpricing.Outbound, "us-east-1")
	fmt.Printf("Data Transfer (100GB outbound): $%.2f\n", dtOutboundLarge)

	dtInterAZ := networkingpricing.CalculateDataTransferCost(50.0, networkingpricing.InterAZ, "us-east-1")
	fmt.Printf("Data Transfer (50GB inter-AZ): $%.2f\n", dtInterAZ)

	// Test 3: Estimate Costs Using Pricing Service
	fmt.Println("\n--- Test 3: Estimate Costs (Pricing Service) ---")

	// Create resources for cost estimation
	natResource := &resource.Resource{
		Type: resource.ResourceType{
			Name: "nat_gateway",
		},
		Provider: "aws",
		Region:   "us-east-1",
	}

	eipResource := &resource.Resource{
		Type: resource.ResourceType{
			Name: "elastic_ip",
		},
		Provider: "aws",
		Region:   "us-east-1",
	}

	eniResource := &resource.Resource{
		Type: resource.ResourceType{
			Name: "network_interface",
		},
		Provider: "aws",
		Region:   "us-east-1",
	}

	// Estimate NAT Gateway cost for 30 days
	natEstimate, err := pricingService.EstimateCost(ctx, natResource, 720*time.Hour)
	if err != nil {
		fmt.Printf("Error estimating NAT Gateway cost: %v\n", err)
	} else {
		fmt.Printf("NAT Gateway Estimate (30 days):\n")
		fmt.Printf("  Total Cost: $%.2f\n", natEstimate.TotalCost)
		fmt.Printf("  Currency: %s\n", natEstimate.Currency)
		fmt.Printf("  Period: %s\n", natEstimate.Period)
		fmt.Printf("  Breakdown:\n")
		for _, comp := range natEstimate.Breakdown {
			fmt.Printf("    - %s: %.2f %s @ $%.3f = $%.2f\n",
				comp.ComponentName, comp.Quantity, comp.Model, comp.UnitRate, comp.Subtotal)
		}
	}

	// Estimate Elastic IP cost for 30 days
	eipEstimate, err := pricingService.EstimateCost(ctx, eipResource, 720*time.Hour)
	if err != nil {
		fmt.Printf("Error estimating Elastic IP cost: %v\n", err)
	} else {
		fmt.Printf("\nElastic IP Estimate (30 days, unattached):\n")
		fmt.Printf("  Total Cost: $%.2f\n", eipEstimate.TotalCost)
		if len(eipEstimate.Breakdown) > 0 {
			fmt.Printf("  Breakdown:\n")
			for _, comp := range eipEstimate.Breakdown {
				fmt.Printf("    - %s: %.2f %s @ $%.3f = $%.2f\n",
					comp.ComponentName, comp.Quantity, comp.Model, comp.UnitRate, comp.Subtotal)
			}
		}
	}

	// Test 4: Architecture Cost Calculation
	fmt.Println("\n--- Test 4: Architecture Cost Calculation ---")

	architectureResources := []*resource.Resource{
		natResource,
		eipResource,
		eniResource,
	}

	archEstimate, err := pricingService.EstimateArchitectureCost(ctx, architectureResources, 720*time.Hour)
	if err != nil {
		fmt.Printf("Error estimating architecture cost: %v\n", err)
	} else {
		fmt.Printf("Architecture Estimate (30 days):\n")
		fmt.Printf("  Total Cost: $%.2f\n", archEstimate.TotalCost)
		fmt.Printf("  Currency: %s\n", archEstimate.Currency)
		fmt.Printf("  Period: %s\n", archEstimate.Period)
		fmt.Printf("  Number of Components: %d\n", len(archEstimate.Breakdown))
		fmt.Printf("  Breakdown:\n")
		for _, comp := range archEstimate.Breakdown {
			fmt.Printf("    - %s: %.2f %s @ $%.3f = $%.2f\n",
				comp.ComponentName, comp.Quantity, comp.Model, comp.UnitRate, comp.Subtotal)
		}
	}

	// Test 5: List Supported Resources
	fmt.Println("\n--- Test 5: List Supported Resources ---")

	supportedResources, err := pricingService.ListSupportedResources(ctx, "aws")
	if err != nil {
		fmt.Printf("Error listing supported resources: %v\n", err)
	} else {
		fmt.Printf("Supported AWS Resources:\n")
		for i, res := range supportedResources {
			fmt.Printf("  %d. %s\n", i+1, res)
		}
	}

	fmt.Println("\n============================================")
	fmt.Println("PRICING TESTS COMPLETED")
	fmt.Println("============================================")
}