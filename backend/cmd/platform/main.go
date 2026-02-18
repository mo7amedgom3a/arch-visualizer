package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	userrepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/user"
)

func main() {
	// Get the migrations directory path relative to backend root
	migrationsDir := filepath.Join("migrations")

	// Run database migrations
	fmt.Println("Running database migrations...")
	if err := database.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	fmt.Println("Migrations completed successfully!")

	// Connect to database
	fmt.Println("Connecting to database...")
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Database connected successfully!")

	// Test database connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Example: You can now use repositories here
	userRepo, err := userrepo.NewUserRepository()
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	ctx := context.Background()
	users, err := userRepo.List(ctx, 10, 0)
	if err != nil {
		log.Printf("Error listing users: %v", err)
	} else {
		fmt.Printf("Found %d users\n", len(users))
	}

	fmt.Println("Application initialized successfully!")

	// Keep connection alive for testing (remove in production)
	// In production, you would start your HTTP server here
	_ = context.Background()
}
