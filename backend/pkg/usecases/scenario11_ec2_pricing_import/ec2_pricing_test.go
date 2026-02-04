package scenario11_ec2_pricing_import

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
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/services/pricing_importer"
)

// EC2PricingImportTestRunner tests the new EC2 pricing import and calculation feature:
// 1. Runs database migration to add instance_type and operating_system fields
// 2. Imports EC2 pricing data from scraper JSON (if file provided)
// 3. Tests pricing calculation with different EC2 instance types
// 4. Compares DB rates vs hardcoded rates
// 5. Creates a test architecture and calculates pricing
func EC2PricingImportTestRunner(ctx context.Context, scraperJSONPath string) error {
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("SCENARIO 11: EC2 Pricing Import and Calculation Test")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("\nThis scenario demonstrates:")
	fmt.Println("  • Database migration for instance_type and operating_system fields")
	fmt.Println("  • Importing EC2 pricing data from scraper JSON")
	fmt.Println("  • Testing pricing calculation with different EC2 instance types")
	fmt.Println("  • Comparing DB rates vs hardcoded rates")
	fmt.Println(strings.Repeat("=", 100))

	// Step 1: Run database migrations
	fmt.Println("\n[Step 1] Running database migrations...")
	migrationsDir := filepath.Join(getBackendRoot(), "migrations")
	if err := database.RunMigrations(migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	fmt.Println("✓ Database migrations completed successfully")

	// Step 2: Import EC2 pricing data (if file provided)
	if scraperJSONPath != "" {
		fmt.Printf("\n[Step 2] Importing EC2 pricing data from: %s\n", scraperJSONPath)

		// Connect to database
		if _, err := database.Connect(); err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		// Create importer
		importer, err := pricing_importer.NewImporter()
		if err != nil {
			return fmt.Errorf("failed to create importer: %w", err)
		}

		// Import pricing data
		stats, err := importer.ImportEC2Pricing(ctx, scraperJSONPath)
		if err != nil {
			return fmt.Errorf("failed to import pricing data: %w", err)
		}

		fmt.Println("✓ EC2 pricing data imported successfully")
		fmt.Printf("  Total Instances Processed: %d\n", stats.TotalInstances)
		fmt.Printf("  Total Rates Imported: %d\n", stats.TotalRates)
		if len(stats.RegionsProcessed) > 0 {
			fmt.Printf("  Regions: %d\n", len(stats.RegionsProcessed))
		}
		if len(stats.OSProcessed) > 0 {
			fmt.Printf("  Operating Systems: %d\n", len(stats.OSProcessed))
		}
	} else {
		fmt.Println("\n[Step 2] Skipping import (no scraper JSON path provided)")
		fmt.Println("  To import pricing data, provide the path to scraper output:")
		fmt.Println("  Example: www/instances.json")
	}

	// Step 3: Test pricing calculation with different instance types
	fmt.Println("\n[Step 3] Testing pricing calculation with different EC2 instance types...")

	// Connect to database if not already connected
	if _, err := database.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create repositories
	pricingRateRepo, err := repository.NewPricingRateRepository()
	if err != nil {
		return fmt.Errorf("failed to create pricing rate repository: %w", err)
	}

	hiddenDepRepo, err := repository.NewHiddenDependencyRepository()
	if err != nil {
		return fmt.Errorf("failed to create hidden dependency repository: %w", err)
	}

	// Create pricing service with repositories
	pricingService := pricing.NewAWSPricingServiceWithRepos(pricingRateRepo, hiddenDepRepo)
	calculator := pricingService.GetCalculator()

	// Test instance types
	testInstances := []struct {
		instanceType string
		region       string
		os           string
	}{
		{"t3.micro", "us-east-1", "linux"},
		{"t3.small", "us-east-1", "linux"},
		{"m5.large", "us-east-1", "linux"},
		{"c5.xlarge", "us-east-1", "linux"},
		{"t3.micro", "us-east-1", "mswin"},
		{"t3.micro", "us-west-2", "linux"},
	}

	duration := 720 * time.Hour // Monthly

	fmt.Println("\n📊 Pricing Calculation Results:")
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("%-20s %-15s %-10s %-15s %-15s %-10s\n", "Instance Type", "Region", "OS", "Hourly Rate", "Monthly Cost", "Source")
	fmt.Println(strings.Repeat("-", 100))

	for _, test := range testInstances {
		// Create resource
		res := &resource.Resource{
			ID:       uuid.New().String(),
			Name:     fmt.Sprintf("test-%s", test.instanceType),
			Provider: "aws",
			Region:   test.region,
			Type:     resource.ResourceType{Name: "EC2"},
			Metadata: map[string]interface{}{
				"instance_type":    test.instanceType,
				"operating_system": test.os,
			},
		}

		// Calculate cost
		estimate, err := calculator.CalculateResourceCost(ctx, res, duration)
		if err != nil {
			fmt.Printf("%-20s %-15s %-10s %-15s %-15s %-10s\n",
				test.instanceType, test.region, test.os,
				"ERROR", "ERROR", "ERROR")
			continue
		}

		// Extract hourly rate
		hourlyRate := 0.0
		if len(estimate.Breakdown) > 0 {
			hourlyRate = estimate.Breakdown[0].UnitRate
		}

		// Determine source
		source := "DB"
		if hourlyRate == 0 {
			source = "Fallback"
		}

		fmt.Printf("%-20s %-15s %-10s $%-14.6f $%-14.2f %-10s\n",
			test.instanceType,
			test.region,
			test.os,
			hourlyRate,
			estimate.TotalCost,
			source,
		)
	}

	// Step 4: Create test architecture and calculate pricing
	fmt.Println("\n[Step 4] Creating test architecture with multiple EC2 instance types...")

	// Initialize server
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	srv, err := server.NewServer(logger)
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}

	// Create architecture with different instance types
	diagramJSON := createTestArchitectureWithMultipleInstanceTypes()

	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	processReq := &serverinterfaces.ProcessDiagramRequest{
		JSONData:        diagramJSON,
		UserID:          userID,
		ProjectName:     "EC2 Pricing Test - Multiple Instance Types",
		IACToolID:       1, // Terraform
		CloudProvider:   "aws",
		Region:          "us-east-1",
		PricingDuration: duration,
	}

	result, err := srv.PipelineOrchestrator.ProcessDiagram(ctx, processReq)
	if err != nil {
		return fmt.Errorf("failed to process diagram: %w", err)
	}

	fmt.Printf("✓ Architecture processed successfully\n")
	fmt.Printf("  Project ID: %s\n", result.ProjectID.String())

	// Step 5: Display pricing breakdown
	if result.PricingEstimate != nil {
		fmt.Println("\n[Step 5] Pricing Breakdown:")
		fmt.Println(strings.Repeat("=", 100))
		fmt.Printf("💰 TOTAL MONTHLY COST: $%.2f %s\n",
			result.PricingEstimate.TotalCost,
			result.PricingEstimate.Currency,
		)
		fmt.Printf("   Provider: %s\n", result.PricingEstimate.Provider)
		fmt.Printf("   Region: %s\n", result.PricingEstimate.Region)
		fmt.Printf("   Period: %s\n", result.PricingEstimate.Period)

		fmt.Println("\n📋 Resource-by-Resource Breakdown:")
		fmt.Println(strings.Repeat("-", 100))
		for i, resEstimate := range result.PricingEstimate.ResourceEstimates {
			fmt.Printf("\n[Resource #%s] %s\n", i, strings.Repeat("-", 90))
			fmt.Printf("  Name:        %s\n", resEstimate.ResourceName)
			fmt.Printf("  Type:        %s\n", resEstimate.ResourceType)
			fmt.Printf("  Total Cost:  $%.4f %s\n", resEstimate.TotalCost, resEstimate.Currency)

			if len(resEstimate.Breakdown) > 0 {
				fmt.Printf("  Components:\n")
				for _, comp := range resEstimate.Breakdown {
					fmt.Printf("    • %-40s %10.4f × $%10.6f = $%10.4f\n",
						comp.ComponentName,
						comp.Quantity,
						comp.UnitRate,
						comp.Subtotal,
					)
				}
			}
		}
	}

	// Summary
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("✅ SUCCESS: EC2 Pricing Import and Calculation Test completed!")
	fmt.Println(strings.Repeat("=", 100))

	return nil
}

// createTestArchitectureWithMultipleInstanceTypes creates a test architecture with different EC2 instance types
func createTestArchitectureWithMultipleInstanceTypes() []byte {
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
						"name": "test-vpc",
						"cidr": "10.0.0.0/16",
					},
					"status":       "valid",
					"isVisualOnly": false,
				},
				"parentId": "region-1",
			},
			// Subnet
			{
				"id":       "subnet-1",
				"type":     "containerNode",
				"position": map[string]float64{"x": 50, "y": 200},
				"data": map[string]interface{}{
					"label":        "Subnet",
					"resourceType": "subnet",
					"config": map[string]interface{}{
						"name":               "test-subnet",
						"cidr":               "10.0.1.0/24",
						"availabilityZoneId": "us-east-1a",
					},
					"status":       "valid",
					"isVisualOnly": false,
				},
				"parentId": "vpc-1",
			},
			// EC2 Instance 1 - t3.micro
			{
				"id":       "ec2-1",
				"type":     "resourceNode",
				"position": map[string]float64{"x": 100, "y": 250},
				"data": map[string]interface{}{
					"label":        "Web Server 1",
					"resourceType": "ec2",
					"config": map[string]interface{}{
						"name":         "web-server-1",
						"instanceType": "t3.micro",
						"ami":          "ami-0123456789abcdef0",
					},
					"status":       "valid",
					"isVisualOnly": false,
				},
				"parentId": "subnet-1",
			},
			// EC2 Instance 2 - t3.small
			{
				"id":       "ec2-2",
				"type":     "resourceNode",
				"position": map[string]float64{"x": 200, "y": 250},
				"data": map[string]interface{}{
					"label":        "Web Server 2",
					"resourceType": "ec2",
					"config": map[string]interface{}{
						"name":         "web-server-2",
						"instanceType": "t3.small",
						"ami":          "ami-0123456789abcdef0",
					},
					"status":       "valid",
					"isVisualOnly": false,
				},
				"parentId": "subnet-1",
			},
			// EC2 Instance 3 - m5.large
			{
				"id":       "ec2-3",
				"type":     "resourceNode",
				"position": map[string]float64{"x": 300, "y": 250},
				"data": map[string]interface{}{
					"label":        "App Server",
					"resourceType": "ec2",
					"config": map[string]interface{}{
						"name":         "app-server",
						"instanceType": "m5.large",
						"ami":          "ami-0123456789abcdef0",
					},
					"status":       "valid",
					"isVisualOnly": false,
				},
				"parentId": "subnet-1",
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

// getBackendRoot returns the backend root directory
func getBackendRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	// Go up from pkg/usecases/scenario11_ec2_pricing_import/ec2_pricing_test.go
	// to backend root
	dir := filepath.Dir(filename)
	dir = filepath.Dir(filepath.Dir(filepath.Dir(dir))) // pkg -> backend
	return dir
}
