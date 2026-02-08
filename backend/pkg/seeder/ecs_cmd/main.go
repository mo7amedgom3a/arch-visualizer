package main

import (
	"context"
	"log"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/seeder"
)

func main() {
	log.Println("ECS Data Seeder")
	log.Println("================")

	// Initialize database
	if _, err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run ECS seeder
	ctx := context.Background()
	if err := seeder.SeedECSData(ctx); err != nil {
		log.Fatalf("ECS seeder failed: %v", err)
	}

	log.Println("✓ ECS data seeding complete!")
}
