package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario10_pricing_with_hidden_costs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario12_api_controllers"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario13_resource_constraints"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario5_terraform_codegen"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario6_terraform_with_persistence"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario7_service_layer"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario8_architecture_roundtrip"
	"github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario9_architecture_pricing"
)

func main() {
	scenario := flag.Int("scenario", 0, "Scenario to run (5=Terraform codegen, 6=Terraform with DB persistence, 7=Service Layer, 8=Architecture Roundtrip, 9=Architecture Pricing, 10=Pricing with Hidden Costs, 12=API Controllers Simulation, 13=Resource Constraints Verification)")
	flag.Parse()

	if *scenario == 0 {
		fmt.Println("Please specify a scenario ID greater than 0")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var err error
	switch *scenario {
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
	case 13:
		err = scenario13_resource_constraints.Run(context.Background())
	default:
		fmt.Printf("Unknown scenario: %d\n", *scenario)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
