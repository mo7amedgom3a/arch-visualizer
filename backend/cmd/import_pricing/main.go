package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/services/pricing_importer"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "file", "", "Path to the scraper EC2 instances JSON file (e.g., www/instances.json)")
	flag.Parse()

	if filePath == "" {
		fmt.Fprintf(os.Stderr, "Error: -file flag is required\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -file <path-to-instances.json>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s -file ../../scripts/scraper/www/instances.json\n", os.Args[0])
		os.Exit(1)
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: File not found: %s\n", filePath)
		os.Exit(1)
	}

	// Connect to database
	if _, err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create importer
	importer, err := pricing_importer.NewImporter()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create importer: %v\n", err)
		os.Exit(1)
	}

	// Import pricing data
	ctx := context.Background()
	fmt.Printf("Importing EC2 pricing data from: %s\n", filePath)
	fmt.Println("This may take a few minutes...")

	stats, err := importer.ImportEC2Pricing(ctx, filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Import failed: %v\n", err)
		if len(stats.Errors) > 0 {
			fmt.Fprintf(os.Stderr, "Errors encountered:\n")
			for _, e := range stats.Errors {
				fmt.Fprintf(os.Stderr, "  - %s\n", e)
			}
		}
		os.Exit(1)
	}

	// Print statistics
	fmt.Println("\n✅ Import completed successfully!")
	fmt.Printf("\n📊 Import Statistics:\n")
	fmt.Printf("  Total Instances Processed: %d\n", stats.TotalInstances)
	fmt.Printf("  Total Rates Imported: %d\n", stats.TotalRates)
	
	if len(stats.RegionsProcessed) > 0 {
		fmt.Printf("\n  Regions Processed:\n")
		for region, count := range stats.RegionsProcessed {
			fmt.Printf("    %s: %d rates\n", region, count)
		}
	}
	
	if len(stats.OSProcessed) > 0 {
		fmt.Printf("\n  Operating Systems Processed:\n")
		for os, count := range stats.OSProcessed {
			fmt.Printf("    %s: %d rates\n", os, count)
		}
	}

	if len(stats.Errors) > 0 {
		fmt.Printf("\n⚠️  Warnings:\n")
		for _, e := range stats.Errors {
			fmt.Printf("  - %s\n", e)
		}
	}

	fmt.Println("\n✨ Pricing data is now available in the database!")
}
