package main

import (
	"fmt"
	"os"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/routes"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/architecture" // Register AWS architecture generator
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/logger"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server"
)

// @title           Arch Visualizer Backend API
// @version         1.0
// @description     Backend API for Arch Visualizer.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:9000
// @BasePath        /api/v1

func main() {
	if err := runServer(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runServer() error {
	// Connect to database
	if _, err := database.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	fmt.Println("✓ Database connected successfully")

	// Initialize Logger
	if err := logger.Init(logger.Config{LogDir: "log"}); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	log := logger.Get()
	log.Info("Logger initialized")

	// Initialize Server
	srv, err := server.NewServer(log)
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}

	// Setup Router
	r := routes.SetupRouter(srv)

	// Run Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	return r.Run(":" + port)
}

/*
	// Old code preserved for reference:
	// // Parse command line flags
	// diagramFile := flag.String("diagram", "json-request-diagram-valid.json", "Path to diagram JSON file")
	// projectName := flag.String("project", "Test Project", "Project name")
	// iacToolID := flag.Uint("iac-tool", 1, "IaC tool ID (1=Terraform, 2=Pulumi, etc.)")
	// userIDStr := flag.String("user-id", "00000000-0000-0000-0000-000000000001", "User ID (UUID)")
	// flag.Parse()

	// // Connect to database
	// if _, err := database.Connect(); err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }
	// log.Println("✓ Database connected successfully")

	// // Read diagram JSON file
	// log.Printf("Reading diagram file: %s", *diagramFile)
	// jsonData, err := os.ReadFile(*diagramFile)
	// if err != nil {
	// 	log.Fatalf("Failed to read diagram file: %v", err)
	// }
	// log.Printf("✓ Read %d bytes from diagram file", len(jsonData))

	// // Parse user ID
	// userID, err := uuid.Parse(*userIDStr)
	// if err != nil {
	// 	log.Fatalf("Invalid user ID format: %v", err)
	// }

	// // Ensure user exists (create if doesn't exist)
	// ctx := context.Background()
	// userRepo, err := repository.NewUserRepository()
	// if err != nil {
	// 	log.Fatalf("Failed to create user repository: %v", err)
	// }

	// user, err := userRepo.FindByID(ctx, userID)
	// if err != nil {
	// 	// User doesn't exist, create it
	// 	log.Printf("User %s not found, creating new user...", userID.String())
	// 	newUser := &models.User{
	// 		ID:         userID,
	// 		Name:       "Test User",
	// 		IsVerified: false,
	// 	}
	// 	if err := userRepo.Create(ctx, newUser); err != nil {
	// 		log.Fatalf("Failed to create user: %v", err)
	// 	}
	// 	log.Printf("✓ Created user: %s", userID.String())
	// 	user = newUser
	// } else {
	// 	log.Printf("✓ User found: %s (%s)", user.Name, userID.String())
	// }

	// // Create diagram service
	// log.Println("Initializing diagram service...")
	// diagramService, err := diagram.NewService()
	// if err != nil {
	// 	log.Fatalf("Failed to create diagram service: %v", err)
	// }
	// log.Println("✓ Diagram service created")

	// // Process diagram
	// log.Println("Processing diagram...")
	// log.Printf("  Project Name: %s", *projectName)
	// log.Printf("  IaC Tool ID: %d", *iacToolID)
	// log.Printf("  User ID: %s", userID.String())
	// log.Println()

	// projectID, err := diagramService.ProcessDiagramRequest(
	// 	ctx,
	// 	jsonData,
	// 	userID,
	// 	*projectName,
	// 	uint(*iacToolID),
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to process diagram: %v", err)
	// }

	// // Success
	// fmt.Println()
	// log.Println(strings.Repeat("=", 52))
	// log.Println("✓ Diagram processed successfully!")
	// log.Printf("  Project ID: %s", projectID.String())
	// log.Printf("  Project Name: %s", *projectName)
	// log.Println(strings.Repeat("=", 52))
*/
