package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/routes"
	_ "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/architecture" // Register AWS architecture generator
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario10_pricing_with_hidden_costs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario12_api_controllers"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario5_terraform_codegen"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario6_terraform_with_persistence"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario7_service_layer"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario8_architecture_roundtrip"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario9_architecture_pricing"
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
	scenario := flag.Int("scenario", 0, "Scenario to run (0=API Server, 5=Terraform codegen, 6=Terraform with DB persistence, 7=Service Layer, 8=Architecture Roundtrip, 9=Architecture Pricing, 10=Pricing with Hidden Costs, 12=API Controllers Simulation)")
	flag.Parse()

	var err error
	switch *scenario {
	case 0:
		err = runServer()
	case 5:
		err = scenario5_terraform_codegen.TerraformCodegenRunner(context.Background())
	case 6:
		err = scenario6_terraform_with_persistence.TerraformWithPersistenceRunner(context.Background())
	case 7:
		err = scenario7_service_layer.TerraformWithServiceLayerRunner(context.Background())
	case 8:
		err = scenario8_architecture_roundtrip.ArchitectureRoundtripRunner(context.Background())
	case 9:
		err = scenario9_architecture_pricing.ArchitecturePricingRunner(context.Background())
	case 10:
		err = scenario10_pricing_with_hidden_costs.PricingWithHiddenCostsRunner(context.Background())
	case 12:
		err = scenario12_api_controllers.Run(context.Background())
	default:
		fmt.Printf("Unknown scenario: %d\n", *scenario)
		os.Exit(1)
	}

	if err != nil {
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

	// Initialize Server
	srv, err := server.NewServer()
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
