package scenario10_pricing_with_hidden_costs

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// PricingWithHiddenCostsRunner demonstrates pricing calculation with hidden dependencies:
// 1. Creates a new architecture with resources that have hidden costs (NAT Gateway, EC2, etc.)
// 2. Processes the architecture through the pipeline
// 3. Calculates pricing including hidden dependency costs
// 4. Displays detailed breakdown showing base costs and hidden costs
func PricingWithHiddenCostsRunner(ctx context.Context) error {
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("SCENARIO 10: Pricing with Hidden Dependency Costs")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("\nThis scenario demonstrates:")
	fmt.Println("  • Creating a new architecture with NAT Gateway and EC2 instances")
	fmt.Println("  • Automatic calculation of hidden dependency costs")
	fmt.Println("  • NAT Gateway → Elastic IP (hidden cost)")
	fmt.Println("  • EC2 → EBS Root Volume (hidden cost)")
	fmt.Println("  • EC2 → Network Interface (hidden cost, free when attached)")
	fmt.Println(strings.Repeat("=", 100))

	// Step 1: Initialize service layer server
	fmt.Println("\n[Step 1] Initializing service layer server...")
	srv, err := server.NewServer(slog.Default())
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}
	fmt.Println("✓ Service layer server initialized successfully")

	// Step 2: Create architecture diagram JSON programmatically
	fmt.Println("\n[Step 2] Creating architecture diagram with hidden dependencies...")
	diagramJSON := createArchitectureWithHiddenDependencies()
	fmt.Println("✓ Architecture diagram created")
	fmt.Println("  Resources included:")
	fmt.Println("    • VPC (10.0.0.0/16)")
	fmt.Println("    • Public Subnet (10.0.1.0/24)")
	fmt.Println("    • Private Subnet (10.0.2.0/24)")
	fmt.Println("    • NAT Gateway (will create hidden Elastic IP)")
	fmt.Println("    • EC2 Instance t3.micro (will create hidden EBS volume and Network Interface)")
	fmt.Println("    • EC2 Instance t3.small (will create hidden EBS volume and Network Interface)")

	// Step 3: Process diagram with pricing calculation
	fmt.Println("\n[Step 3] Processing architecture and calculating pricing...")
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	pricingDuration := 720 * time.Hour // Monthly pricing

	processReq := &serverinterfaces.ProcessDiagramRequest{
		JSONData:        diagramJSON,
		UserID:          userID,
		ProjectName:     "Pricing with Hidden Costs Demo",
		IACToolID:       1, // Terraform
		CloudProvider:   "aws",
		Region:          "us-east-1",
		PricingDuration: pricingDuration,
	}

	result, err := srv.PipelineOrchestrator.ProcessDiagram(ctx, processReq)
	if err != nil {
		return fmt.Errorf("failed to process diagram: %w", err)
	}

	fmt.Printf("✓ Architecture processed and saved to database\n")
	fmt.Printf("  Project ID: %s\n", result.ProjectID.String())
	fmt.Printf("  Success: %v\n", result.Success)
	fmt.Printf("  Message: %s\n", result.Message)

	// Step 4: Display detailed pricing breakdown with hidden costs
	fmt.Println("\n[Step 4] Detailed Pricing Breakdown (Including Hidden Costs):")
	fmt.Println(strings.Repeat("=", 100))

	if result.PricingEstimate != nil {
		// Display resource-by-resource breakdown
		displayResourceCosts(result.PricingEstimate, pricingDuration)

		// Calculate totals
		var totalBaseCost float64
		var totalHiddenCost float64
		var totalCost float64

		for _, resEstimate := range result.PricingEstimate.ResourceEstimates {
			baseCost, hiddenCost := calculateBaseAndHiddenCosts(resEstimate.Breakdown)
			totalBaseCost += baseCost
			totalHiddenCost += hiddenCost
			totalCost += resEstimate.TotalCost
		}

		// Summary table
		fmt.Printf("\n📊 COST SUMMARY TABLE\n")
		fmt.Println(strings.Repeat("=", 100))
		fmt.Printf("%-40s %-20s %15s %15s\n", "Resource", "Type", "Base Cost", "Total Cost")
		fmt.Println(strings.Repeat("-", 100))

		for _, resEstimate := range result.PricingEstimate.ResourceEstimates {
			baseCost, _ := calculateBaseAndHiddenCosts(resEstimate.Breakdown)
			fmt.Printf("%-40s %-20s $%13.4f $%13.4f\n",
				truncateString(resEstimate.ResourceName, 38),
				truncateString(resEstimate.ResourceType, 18),
				baseCost,
				resEstimate.TotalCost,
			)
		}

		fmt.Println(strings.Repeat("-", 100))
		fmt.Printf("%-40s %-20s $%13.4f $%13.4f\n",
			"TOTAL BASE COSTS",
			"",
			totalBaseCost,
			totalBaseCost,
		)
		fmt.Printf("%-40s %-20s %15s $%13.4f\n",
			"TOTAL HIDDEN DEPENDENCY COSTS",
			"",
			"",
			totalHiddenCost,
		)
		fmt.Println(strings.Repeat("=", 100))
		fmt.Printf("\n💰 TOTAL MONTHLY COST: $%.2f %s\n",
			totalCost,
			result.PricingEstimate.Currency,
		)
		fmt.Printf("   Provider: %s\n", result.PricingEstimate.Provider)
		fmt.Printf("   Region: %s\n", result.PricingEstimate.Region)
		fmt.Printf("   Period: %s\n", result.PricingEstimate.Period)

		// Calculate yearly estimate
		yearlyEstimate := totalCost * 12
		fmt.Printf("\n📅 ESTIMATED YEARLY COST: $%.2f %s\n", yearlyEstimate, result.PricingEstimate.Currency)

		// Show hidden dependency summary
		if totalHiddenCost > 0 {
			fmt.Printf("\n🔗 HIDDEN DEPENDENCY COSTS SUMMARY:\n")
			fmt.Println(strings.Repeat("-", 100))
			fmt.Printf("  Total Hidden Costs: $%.2f %s\n", totalHiddenCost, result.PricingEstimate.Currency)
			fmt.Printf("  Percentage of Total: %.1f%%\n", (totalHiddenCost/totalCost)*100)
			fmt.Println("\n  Hidden Dependencies Detected:")
			fmt.Println("    • NAT Gateway → Elastic IP (free when attached)")
			fmt.Println("    • EC2 Instances → EBS Root Volumes (8GB gp3 default)")
			fmt.Println("    • EC2 Instances → Network Interfaces (free when attached)")
		}
	} else {
		fmt.Println("⚠️  No pricing estimate available")
	}

	// Step 5: Save pricing report to file
	fmt.Println("\n[Step 5] Saving pricing report to file...")
	outputPath := filepath.Join(filepath.Dir(getCurrentDir()), "json-response-pricing-with-hidden-costs.json")
	report := map[string]interface{}{
		"project_id":   result.ProjectID.String(),
		"success":      result.Success,
		"message":      result.Message,
		"generated_at": time.Now(),
		"duration":     pricingDuration.String(),
		"pricing":      result.PricingEstimate,
	}
	outputData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal pricing report JSON: %w", err)
	}
	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write pricing report JSON file: %w", err)
	}
	fmt.Printf("✓ Pricing report saved to: %s (%d bytes)\n", outputPath, len(outputData))

	// Summary
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("✅ SUCCESS: Pricing with hidden costs calculation completed!")
	fmt.Printf("  Project ID: %s\n", result.ProjectID.String())
	if result.PricingEstimate != nil {
		fmt.Printf("  Monthly Cost: $%.2f %s\n", result.PricingEstimate.TotalCost, result.PricingEstimate.Currency)
		fmt.Printf("  Yearly Cost: $%.2f %s\n", result.PricingEstimate.TotalCost*12, result.PricingEstimate.Currency)
	}
	fmt.Printf("  Pricing Report: %s\n", outputPath)
	fmt.Println(strings.Repeat("=", 100))

	return nil
}

// createArchitectureWithHiddenDependencies creates a diagram JSON with resources that have hidden dependencies
func createArchitectureWithHiddenDependencies() []byte {
	diagram := map[string]interface{}{
		"nodes": []map[string]interface{}{
			// Region
			{
				"id":       "region-1",
				"type":     "containerNode",
				"position": map[string]float64{"x": 400, "y": 200},
				"data": map[string]interface{}{
					"label":        "Region",
					"resourceType": "region",
					"config": map[string]interface{}{
						"name": "us-east-1",
					},
					"status":       "valid",
					"isVisualOnly": false,
				},
				"parentId": nil,
			},
			// VPC
			{
				"id":       "vpc-1",
				"type":     "containerNode",
				"position": map[string]float64{"x": 100, "y": 100},
				"data": map[string]interface{}{
					"label":        "VPC",
					"resourceType": "vpc",
					"config": map[string]interface{}{
						"name": "demo-vpc",
						"cidr": "10.0.0.0/16",
					},
					"status":       "valid",
					"isVisualOnly": false,
				},
				"parentId": "region-1",
			},
			// Public Subnet
			{
				"id":       "subnet-public-1",
				"type":     "containerNode",
				"position": map[string]float64{"x": 50, "y": 200},
				"data": map[string]interface{}{
					"label":        "Public Subnet",
					"resourceType": "subnet",
					"config": map[string]interface{}{
						"name":               "public-subnet-1a",
						"cidr":               "10.0.1.0/24",
						"availabilityZoneId": "us-east-1a",
					},
					"status":       "valid",
					"isVisualOnly": false,
				},
				"parentId": "vpc-1",
			},
			// Private Subnet
			{
				"id":       "subnet-private-1",
				"type":     "containerNode",
				"position": map[string]float64{"x": 200, "y": 200},
				"data": map[string]interface{}{
					"label":        "Private Subnet",
					"resourceType": "subnet",
					"config": map[string]interface{}{
						"name":               "private-subnet-1a",
						"cidr":               "10.0.2.0/24",
						"availabilityZoneId": "us-east-1a",
					},
					"status":       "valid",
					"isVisualOnly": false,
				},
				"parentId": "vpc-1",
			},
			// NAT Gateway (will create hidden Elastic IP)
			{
				"id":       "nat-gateway-1",
				"type":     "resourceNode",
				"position": map[string]float64{"x": 80, "y": 220},
				"data": map[string]interface{}{
					"label":        "NAT Gateway",
					"resourceType": "nat-gateway",
					"config": map[string]interface{}{
						"name":     "demo-nat-gateway",
						"subnetId": "subnet-public-1",
						// Note: allocationId is NOT provided, so a hidden Elastic IP will be created
					},
					"status":       "valid",
					"isVisualOnly": false,
				},
				"parentId": "subnet-public-1",
			},
			// EC2 Instance 1 - t3.micro (will create hidden EBS volume and Network Interface)
			{
				"id":       "ec2-1",
				"type":     "resourceNode",
				"position": map[string]float64{"x": 220, "y": 220},
				"data": map[string]interface{}{
					"label":        "Web Server 1",
					"resourceType": "ec2",
					"config": map[string]interface{}{
						"name":         "web-server-1",
						"instanceType": "t3.micro",
						"ami":          "ami-0123456789abcdef0",
						// Note: size_gb not specified, will use default 8GB for EBS root volume
					},
					"status":       "valid",
					"isVisualOnly": false,
				},
				"parentId": "subnet-private-1",
			},
			// EC2 Instance 2 - t3.small (will create hidden EBS volume and Network Interface)
			{
				"id":       "ec2-2",
				"type":     "resourceNode",
				"position": map[string]float64{"x": 280, "y": 220},
				"data": map[string]interface{}{
					"label":        "Web Server 2",
					"resourceType": "ec2",
					"config": map[string]interface{}{
						"name":         "web-server-2",
						"instanceType": "t3.small",
						"ami":          "ami-0123456789abcdef0",
						"size_gb":      20, // Explicit EBS root volume size
					},
					"status":       "valid",
					"isVisualOnly": false,
				},
				"parentId": "subnet-private-1",
			},
		},
		"edges":     []interface{}{},
		"variables": []interface{}{},
	}

	jsonData, err := json.Marshal(diagram)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal diagram JSON: %v", err))
	}

	return jsonData
}

// displayResourceCosts displays each resource with its detailed cost breakdown
func displayResourceCosts(estimate *serverinterfaces.ArchitectureCostEstimate, duration time.Duration) {
	fmt.Printf("\n📊 Architecture Cost Estimate (Monthly - %v)\n", duration)
	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("\n💰 RESOURCE-BY-RESOURCE COST BREAKDOWN\n")
	fmt.Println(strings.Repeat("=", 100))

	resourceCount := 0
	for _, resEstimate := range estimate.ResourceEstimates {
		resourceCount++
		baseCost, hiddenCost := calculateBaseAndHiddenCosts(resEstimate.Breakdown)

		// Resource header
		fmt.Printf("\n[Resource #%d] %s\n", resourceCount, strings.Repeat("-", 90))
		fmt.Printf("  Name:        %s\n", resEstimate.ResourceName)
		fmt.Printf("  Type:        %s\n", resEstimate.ResourceType)
		fmt.Printf("  Resource ID: %s\n", truncateString(resEstimate.ResourceID, 36))
		fmt.Printf("  Currency:    %s\n", resEstimate.Currency)
		fmt.Println()

		// Base costs
		if baseCost > 0 {
			fmt.Printf("  📋 BASE COSTS:\n")
			for _, comp := range resEstimate.Breakdown {
				if !isHiddenDependency(comp.ComponentName) {
					fmt.Printf("     • %-45s %-12s %10.4f × $%10.6f = $%10.4f\n",
						truncateString(comp.ComponentName, 43),
						comp.Model,
						comp.Quantity,
						comp.UnitRate,
						comp.Subtotal,
					)
				}
			}
			fmt.Printf("     └─ Subtotal (Base): $%.4f\n", baseCost)
			fmt.Println()
		}

		// Hidden dependency costs
		if hiddenCost > 0 {
			fmt.Printf("  🔗 HIDDEN DEPENDENCY COSTS:\n")
			for _, comp := range resEstimate.Breakdown {
				if isHiddenDependency(comp.ComponentName) {
					depType := extractDependencyType(comp.ComponentName)
					fmt.Printf("     • %-45s %-12s %10.4f × $%10.6f = $%10.4f\n",
						truncateString(comp.ComponentName, 43),
						comp.Model,
						comp.Quantity,
						comp.UnitRate,
						comp.Subtotal,
					)
					if depType != "" {
						fmt.Printf("       └─ Dependency Type: %s\n", depType)
					}
				}
			}
			fmt.Printf("     └─ Subtotal (Hidden): $%.4f\n", hiddenCost)
			fmt.Println()
		}

		// Resource total
		fmt.Printf("  💵 RESOURCE TOTAL: $%.4f %s\n", resEstimate.TotalCost, resEstimate.Currency)
		if hiddenCost > 0 {
			hiddenPercentage := (hiddenCost / resEstimate.TotalCost) * 100
			fmt.Printf("     └─ Hidden costs represent %.1f%% of this resource's total\n", hiddenPercentage)
		}
	}

	fmt.Println(strings.Repeat("=", 100))
}

// calculateBaseAndHiddenCosts separates base costs from hidden dependency costs
func calculateBaseAndHiddenCosts(breakdown []serverinterfaces.CostBreakdownComponent) (baseCost, hiddenCost float64) {
	for _, comp := range breakdown {
		if isHiddenDependency(comp.ComponentName) {
			hiddenCost += comp.Subtotal
		} else {
			baseCost += comp.Subtotal
		}
	}
	return baseCost, hiddenCost
}

// isHiddenDependency checks if a component name indicates a hidden dependency
func isHiddenDependency(componentName string) bool {
	return strings.Contains(componentName, "(") && strings.Contains(componentName, ")")
}

// extractDependencyType extracts the dependency type from component name like "Component (elastic_ip)"
func extractDependencyType(componentName string) string {
	start := strings.Index(componentName, "(")
	end := strings.Index(componentName, ")")
	if start != -1 && end != -1 && end > start {
		return componentName[start+1 : end]
	}
	return ""
}

// Helper functions
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func getCurrentDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}
