package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/diagram"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
)

func main() {
	// Parse command line flags
	diagramFile := flag.String("diagram", "json-request-diagram-valid.json", "Path to diagram JSON file")
	projectName := flag.String("project", "Test Project", "Project name")
	iacToolID := flag.Uint("iac-tool", 1, "IaC tool ID (1=Terraform, 2=Pulumi, etc.)")
	userIDStr := flag.String("user-id", "00000000-0000-0000-0000-000000000001", "User ID (UUID)")
	flag.Parse()

	// Connect to database
	if _, err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✓ Database connected successfully")

	// Read diagram JSON file
	log.Printf("Reading diagram file: %s", *diagramFile)
	jsonData, err := os.ReadFile(*diagramFile)
	if err != nil {
		log.Fatalf("Failed to read diagram file: %v", err)
	}
	log.Printf("✓ Read %d bytes from diagram file", len(jsonData))

	// Parse user ID
	userID, err := uuid.Parse(*userIDStr)
	if err != nil {
		log.Fatalf("Invalid user ID format: %v", err)
	}

	// Ensure user exists (create if doesn't exist)
	ctx := context.Background()
	userRepo, err := repository.NewUserRepository()
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		// User doesn't exist, create it
		log.Printf("User %s not found, creating new user...", userID.String())
		newUser := &models.User{
			ID:         userID,
			Name:       "Test User",
			IsVerified: false,
		}
		if err := userRepo.Create(ctx, newUser); err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		log.Printf("✓ Created user: %s", userID.String())
		user = newUser
	} else {
		log.Printf("✓ User found: %s (%s)", user.Name, userID.String())
	}

	// Create diagram service
	log.Println("Initializing diagram service...")
	diagramService, err := diagram.NewService()
	if err != nil {
		log.Fatalf("Failed to create diagram service: %v", err)
	}
	log.Println("✓ Diagram service created")

	// Process diagram
	log.Println("Processing diagram...")
	log.Printf("  Project Name: %s", *projectName)
	log.Printf("  IaC Tool ID: %d", *iacToolID)
	log.Printf("  User ID: %s", userID.String())
	log.Println()

	projectID, err := diagramService.ProcessDiagramRequest(
		ctx,
		jsonData,
		userID,
		*projectName,
		uint(*iacToolID),
	)
	if err != nil {
		log.Fatalf("Failed to process diagram: %v", err)
	}

	// Success
	fmt.Println()
	log.Println(strings.Repeat("=", 52))
	log.Println("✓ Diagram processed successfully!")
	log.Printf("  Project ID: %s", projectID.String())
	log.Printf("  Project Name: %s", *projectName)
	log.Println(strings.Repeat("=", 52))
}
