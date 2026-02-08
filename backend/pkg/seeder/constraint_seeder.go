package seeder

import (
	"context"
	"log"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/rules"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
)

// SeedResourceConstraints seeds the database with default resource constraints
func SeedResourceConstraints(
	ctx context.Context,
	constraintRepo *repository.ResourceConstraintRepository,
	resourceTypeRepo *repository.ResourceTypeRepository,
) error {
	log.Println("Seeding resource constraints...")

	// Get default rules
	defaultRules := append(rules.DefaultNetworkingRules(), rules.DefaultComputeRules()...)
	defaultRules = append(defaultRules, rules.DefaultStorageRules()...)
	defaultRules = append(defaultRules, rules.DefaultDatabaseRules()...)
	defaultRules = append(defaultRules, rules.DefaultIAMRules()...)

	count := 0
	skipped := 0
	failed := 0

	for _, rule := range defaultRules {
		// Find resource type
		// Assuming AWS rules for now
		resourceType, err := resourceTypeRepo.FindByNameAndProvider(ctx, rule.ResourceType, "aws")
		if err != nil {
			// Resource type might not exist yet if seeding order is wrong or resource types aren't seeded.
			// For now, we'll log and skip.
			log.Printf("Warning: Resource type '%s' not found for constraint seeding. Skipping.", rule.ResourceType)
			failed++
			continue
		}

		// Check if constraint already exists
		// We can do a check by querying existing constraints for this resource type and filtering in memory,
		// or just try to insert and ignore clashes if we had a unique constraint (which we don't naturally have on value/type yet maybe).
		// Better approach: Check existence.
		existingConstraints, err := constraintRepo.FindByResourceType(ctx, resourceType.ID)
		if err != nil {
			log.Printf("Error checking existing constraints for %s: %v", rule.ResourceType, err)
			failed++
			continue
		}

		exists := false
		for _, ec := range existingConstraints {
			if ec.ConstraintType == rule.ConstraintType && ec.ConstraintValue == rule.ConstraintValue {
				exists = true
				break
			}
		}

		if exists {
			skipped++
			continue
		}

		// Create new constraint
		newConstraint := &models.ResourceConstraint{
			ResourceTypeID:  resourceType.ID,
			ConstraintType:  rule.ConstraintType,
			ConstraintValue: rule.ConstraintValue,
		}

		if err := constraintRepo.Create(ctx, newConstraint); err != nil {
			log.Printf("Error creating constraint for %s: %v", rule.ResourceType, err)
			failed++
			continue
		}
		count++
	}

	log.Printf("Constraint seeding finished: %d created, %d skipped, %d failed", count, skipped, failed)
	return nil
}
