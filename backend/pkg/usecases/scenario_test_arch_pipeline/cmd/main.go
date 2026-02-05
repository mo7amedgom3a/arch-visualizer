package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario_test_arch_pipeline"
)

func main() {
	// Initialize Database Connection
	if _, err := database.Connect(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	ctx := context.Background()
	if err := scenario_test_arch_pipeline.Run(ctx); err != nil {
		log.Fatalf("Scenario failed: %v", err)
	}
	fmt.Println("Scenario completed successfully!")
}
