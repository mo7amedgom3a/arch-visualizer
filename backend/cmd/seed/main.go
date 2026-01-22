package main

import (
	"log"

	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/seed"
)

func main() {
	if err := seed.SeedDatabase(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
}
