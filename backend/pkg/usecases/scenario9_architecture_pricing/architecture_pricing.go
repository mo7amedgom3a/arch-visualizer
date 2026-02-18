package scenario9_architecture_pricing

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

// ArchitecturePricingRunner tests the complete architecture pricing flow:
// 1. Read diagram JSON from file
// 2. Process and save to database with pricing calculation
// 3. Display pricing breakdown for each resource
// 4. Display total project cost estimate
func ArchitecturePricingRunner(ctx context.Context) error {
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("SCENARIO 9: Architecture Pricing (Process Diagram → Calculate Pricing → Persist)")
	fmt.Println(strings.Repeat("=", 100))

	// Step 1: Initialize service layer server
	fmt.Println("\n[Step 1] Initializing service layer server...")
	srv, err := server.NewServer(slog.Default())
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}
	fmt.Println("✓ Service layer server initialized successfully")

	// Step 2: Read IR JSON from file
	fmt.Println("\n[Step 2] Reading diagram JSON file...")
	jsonPath, err := resolveDiagramJSONPath("json-request-fiagram-complete.json")
	if err != nil {
		return fmt.Errorf("failed to resolve diagram JSON path: %w", err)
	}

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read diagram JSON file: %w", err)
	}
	fmt.Printf("✓ Read diagram JSON from: %s (%d bytes)\n", jsonPath, len(jsonData))

	// Extract diagram from project-wrapped JSON structure if needed
	diagramData, err := extractDiagramFromProjectJSON(jsonData)
	if err != nil {
		return fmt.Errorf("failed to extract diagram from project JSON: %w", err)
	}

	// Step 3: Process diagram with pricing calculation
	fmt.Println("\n[Step 3] Processing diagram with pricing calculation...")

	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// Use 720 hours (30 days) for monthly pricing estimate
	pricingDuration := 720 * time.Hour

	processReq := &serverinterfaces.ProcessDiagramRequest{
		JSONData:        diagramData,
		UserID:          userID,
		ProjectName:     "Architecture Pricing Test Project",
		IACToolID:       1, // Terraform
		CloudProvider:   "aws",
		Region:          "us-east-1",
		PricingDuration: pricingDuration,
	}

	result, err := srv.PipelineOrchestrator.ProcessDiagram(ctx, processReq)
	if err != nil {
		return fmt.Errorf("failed to process diagram: %w", err)
	}

	fmt.Printf("✓ Diagram processed and saved to database\n")
	fmt.Printf("  Project ID: %s\n", result.ProjectID.String())
	fmt.Printf("  Success: %v\n", result.Success)
	fmt.Printf("  Message: %s\n", result.Message)

	// Step 4: Display pricing breakdown
	fmt.Println("\n[Step 4] Pricing Breakdown:")
	fmt.Println(strings.Repeat("-", 80))

	if result.PricingEstimate != nil {
		fmt.Printf("\n📊 Architecture Cost Estimate (Monthly - %v)\n", pricingDuration)
		fmt.Println(strings.Repeat("-", 80))

		// Display individual resource costs
		fmt.Printf("\n%-40s %-15s %s\n", "Resource", "Type", "Monthly Cost")
		fmt.Println(strings.Repeat("-", 80))

		for _, resEstimate := range result.PricingEstimate.ResourceEstimates {
			fmt.Printf("%-40s %-15s $%.4f %s\n",
				truncateString(resEstimate.ResourceName, 38),
				truncateString(resEstimate.ResourceType, 13),
				resEstimate.TotalCost,
				resEstimate.Currency,
			)

			// Display breakdown components
			for _, comp := range resEstimate.Breakdown {
				// Check if this is a hidden dependency cost (contains resource type in parentheses)
				isHiddenDep := strings.Contains(comp.ComponentName, "(") && strings.Contains(comp.ComponentName, ")")
				prefix := "  └─"
				if isHiddenDep {
					prefix = "  🔗" // Use special marker for hidden dependencies
				}
				fmt.Printf("%s %-36s %-10s %.4f × $%.6f = $%.4f\n",
					prefix,
					truncateString(comp.ComponentName, 34),
					comp.Model,
					comp.Quantity,
					comp.UnitRate,
					comp.Subtotal,
				)
			}
		}

		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("\n💰 TOTAL MONTHLY COST: $%.2f %s\n",
			result.PricingEstimate.TotalCost,
			result.PricingEstimate.Currency,
		)
		fmt.Printf("   Provider: %s\n", result.PricingEstimate.Provider)
		fmt.Printf("   Region: %s\n", result.PricingEstimate.Region)
		fmt.Printf("   Period: %s\n", result.PricingEstimate.Period)

		// Calculate yearly estimate
		yearlyEstimate := result.PricingEstimate.TotalCost * 12
		fmt.Printf("\n📅 ESTIMATED YEARLY COST: $%.2f %s\n", yearlyEstimate, result.PricingEstimate.Currency)
	} else {
		fmt.Println("⚠️  No pricing estimate available")
	}

	// Step 5: Load architecture and verify pricing was persisted
	fmt.Println("\n[Step 5] Verifying pricing persistence...")

	arch, err := srv.ProjectService.LoadArchitecture(ctx, result.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to load architecture: %w", err)
	}

	fmt.Printf("✓ Architecture loaded from database\n")
	fmt.Printf("  Resources: %d\n", len(arch.Resources))

	// Get project pricing from database
	projectPricings, err := srv.ProjectService.GetProjectPricing(ctx, result.ProjectID)
	if err != nil {
		fmt.Printf("⚠️  Could not retrieve project pricing: %v\n", err)
	} else {
		fmt.Printf("✓ Project pricing records in database: %d\n", len(projectPricings))
		for _, pricing := range projectPricings {
			fmt.Printf("  └─ Total: $%.2f %s (%s)\n",
				pricing.TotalCost,
				pricing.Currency,
				pricing.Period,
			)
		}
	}

	// Step 6: Save pricing report to file
	fmt.Println("\n[Step 6] Saving pricing report to file...")

	report := buildPricingReport(result, pricingDuration)
	outputPath := filepath.Join(filepath.Dir(jsonPath), "json-response-architecture-pricing.json")
	outputData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal pricing report: %w", err)
	}

	if err := os.WriteFile(outputPath, outputData, 0o644); err != nil {
		return fmt.Errorf("failed to write pricing report file: %w", err)
	}

	fmt.Printf("✓ Pricing report saved to: %s (%d bytes)\n", outputPath, len(outputData))

	// Summary
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("SUCCESS: Architecture pricing completed!")
	fmt.Printf("  Project ID: %s\n", result.ProjectID.String())
	fmt.Printf("  Resources: %d\n", len(arch.Resources))
	if result.PricingEstimate != nil {
		fmt.Printf("  Monthly Cost: $%.2f %s\n", result.PricingEstimate.TotalCost, result.PricingEstimate.Currency)
		fmt.Printf("  Yearly Cost: $%.2f %s\n", result.PricingEstimate.TotalCost*12, result.PricingEstimate.Currency)
	}
	fmt.Printf("  Pricing Report: %s\n", outputPath)
	fmt.Println(strings.Repeat("=", 100))

	return nil
}

// buildPricingReport creates a structured pricing report
func buildPricingReport(result *serverinterfaces.ProcessDiagramResult, duration time.Duration) map[string]interface{} {
	report := map[string]interface{}{
		"project_id":   result.ProjectID.String(),
		"success":      result.Success,
		"message":      result.Message,
		"generated_at": time.Now().Format(time.RFC3339),
		"duration":     duration.String(),
	}

	if result.PricingEstimate != nil {
		resourceBreakdown := make([]map[string]interface{}, 0)
		for _, resEstimate := range result.PricingEstimate.ResourceEstimates {
			components := make([]map[string]interface{}, 0)
			for _, comp := range resEstimate.Breakdown {
				components = append(components, map[string]interface{}{
					"component_name": comp.ComponentName,
					"model":          comp.Model,
					"quantity":       comp.Quantity,
					"unit_rate":      comp.UnitRate,
					"subtotal":       comp.Subtotal,
					"currency":       comp.Currency,
				})
			}

			resourceBreakdown = append(resourceBreakdown, map[string]interface{}{
				"resource_id":   resEstimate.ResourceID,
				"resource_name": resEstimate.ResourceName,
				"resource_type": resEstimate.ResourceType,
				"total_cost":    resEstimate.TotalCost,
				"currency":      resEstimate.Currency,
				"components":    components,
			})
		}

		report["pricing"] = map[string]interface{}{
			"total_cost":         result.PricingEstimate.TotalCost,
			"currency":           result.PricingEstimate.Currency,
			"period":             result.PricingEstimate.Period,
			"provider":           result.PricingEstimate.Provider,
			"region":             result.PricingEstimate.Region,
			"resource_breakdown": resourceBreakdown,
			"yearly_estimate":    result.PricingEstimate.TotalCost * 12,
		}
	}

	return report
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// resolveDiagramJSONPath resolves the JSON file path
func resolveDiagramJSONPath(filename string) (string, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine caller for resolving diagram JSON path")
	}

	dir := filepath.Dir(thisFile)
	root := filepath.Clean(filepath.Join(dir, ".."))
	jsonPath := filepath.Join(root, filename)

	return jsonPath, nil
}

// extractDiagramFromProjectJSON extracts the diagram structure from project-wrapped JSON
func extractDiagramFromProjectJSON(data []byte) ([]byte, error) {
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Check if this is a direct diagram format (has "nodes" at root)
	if _, hasNodes := rawData["nodes"]; hasNodes {
		return data, nil
	}

	// Otherwise, look for project-wrapped structure
	for _, value := range rawData {
		if projectData, ok := value.(map[string]interface{}); ok {
			if _, hasNodes := projectData["nodes"]; hasNodes {
				diagramBytes, err := json.Marshal(projectData)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal extracted diagram: %w", err)
				}
				return diagramBytes, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find diagram structure in JSON")
}
