package scenario13_resource_constraints

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server"
)

// Run executes the resource constraints verification scenario
func Run(ctx context.Context) error {
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("SCENARIO 13: Resource Constraints Verification")
	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("This scenario verifies that AWS resource constraints (Compute, Storage) are correctly seeded in the database.")

	// Step 1: Initialize server (this triggers seeding)
	fmt.Println("\n[Step 1] Initializing server to trigger seeding...")
	if _, err := server.NewServer(slog.Default()); err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}
	fmt.Println("✓ Server initialized and seeding triggered")

	// Step 2: Verify constraints in database
	fmt.Println("\n[Step 2] Verifying constraints in database...")

	targetResources := []string{"EC2", "Lambda", "S3", "EBS"}

	for _, resourceName := range targetResources {
		fmt.Printf("\nChecking constraints for %s:\n", resourceName)

		var resourceType models.ResourceType
		if err := database.DB.Where("name = ? AND cloud_provider = ?", resourceName, "aws").First(&resourceType).Error; err != nil {
			fmt.Printf("  X Resource type not found: %v\n", err)
			continue
		}

		var constraints []models.ResourceConstraint
		if err := database.DB.Where("resource_type_id = ?", resourceType.ID).Find(&constraints).Error; err != nil {
			fmt.Printf("  X Failed to query constraints: %v\n", err)
			continue
		}

		if len(constraints) == 0 {
			fmt.Println("  X No constraints found")
		} else {
			for _, c := range constraints {
				fmt.Printf("  - %s: %s\n", c.ConstraintType, c.ConstraintValue)
			}
			fmt.Printf("  ✓ Found %d constraints\n", len(constraints))
		}
	}

	return nil
}
