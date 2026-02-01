package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	awspricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// PricingServiceImpl implements PricingService interface
type PricingServiceImpl struct {
	pricingRepo     serverinterfaces.PricingRepository
	pricingRateRepo *repository.PricingRateRepository
	hiddenDepRepo   *repository.HiddenDependencyRepository
	awsCalculator   *awspricing.AWSPricingCalculator
	useDBRates      bool
}

// NewPricingService creates a new pricing service (backward compatible)
func NewPricingService(pricingRepo serverinterfaces.PricingRepository) serverinterfaces.PricingService {
	// Create AWS pricing service and calculator (without DB rates)
	awsPricingService := awspricing.NewAWSPricingService()
	awsCalculator := awspricing.NewAWSPricingCalculator(awsPricingService)

	return &PricingServiceImpl{
		pricingRepo:   pricingRepo,
		awsCalculator: awsCalculator,
		useDBRates:    false,
	}
}

// NewPricingServiceWithRepos creates a new pricing service with database-driven rates
func NewPricingServiceWithRepos(
	pricingRepo serverinterfaces.PricingRepository,
	pricingRateRepo *repository.PricingRateRepository,
	hiddenDepRepo *repository.HiddenDependencyRepository,
) serverinterfaces.PricingService {
	// Create AWS pricing service with repositories
	awsPricingService := awspricing.NewAWSPricingServiceWithRepos(pricingRateRepo, hiddenDepRepo)
	awsCalculator := awsPricingService.GetCalculator()

	return &PricingServiceImpl{
		pricingRepo:     pricingRepo,
		pricingRateRepo: pricingRateRepo,
		hiddenDepRepo:   hiddenDepRepo,
		awsCalculator:   awsCalculator,
		useDBRates:      true,
	}
}

// CalculateResourceCost calculates the cost for a single resource over a given duration
func (s *PricingServiceImpl) CalculateResourceCost(ctx context.Context, res *resource.Resource, duration time.Duration) (*domainpricing.CostEstimate, error) {
	if res == nil {
		return nil, fmt.Errorf("resource is nil")
	}

	// Map resource type name to pricing calculator expected name
	mappedRes := mapResourceTypeForPricing(res)

	// Route to the appropriate calculator based on provider
	switch mappedRes.Provider {
	case resource.AWS:
		return s.awsCalculator.CalculateResourceCost(ctx, mappedRes, duration)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", mappedRes.Provider)
	}
}

// mapResourceTypeForPricing maps resource type names to pricing calculator expected names
func mapResourceTypeForPricing(res *resource.Resource) *resource.Resource {
	// Create a copy with mapped type name
	mappedRes := *res
	mappedType := res.Type

	// Map resource type names to pricing calculator expected names
	typeMapping := map[string]string{
		"EC2":              "ec2_instance",
		"NATGateway":       "nat_gateway",
		"ElasticIP":        "elastic_ip",
		"LoadBalancer":     "load_balancer",
		"AutoScalingGroup": "auto_scaling_group",
		"Lambda":           "lambda_function",
		"S3":               "s3_bucket",
		"EBS":              "ebs_volume",
		"RDS":              "rds_instance",
		"DynamoDB":         "dynamodb_table",
		"NetworkInterface": "network_interface",
	}

	if mapped, ok := typeMapping[res.Type.Name]; ok {
		mappedType.Name = mapped
		mappedType.ID = mapped
	}

	mappedRes.Type = mappedType
	return &mappedRes
}

// CalculateArchitectureCost calculates the total cost for an architecture over a given duration
func (s *PricingServiceImpl) CalculateArchitectureCost(ctx context.Context, arch *architecture.Architecture, duration time.Duration) (*serverinterfaces.ArchitectureCostEstimate, error) {
	if arch == nil {
		return nil, fmt.Errorf("architecture is nil")
	}

	resourceEstimates := make(map[string]*serverinterfaces.ResourceCostEstimate)
	var totalCost float64

	fmt.Println("\n💵 Calculating pricing for resources...")
	fmt.Println(strings.Repeat("-", 100))

	for _, res := range arch.Resources {
		// Skip visual-only resources (they don't have pricing)
		if isVisualOnly, ok := res.Metadata["isVisualOnly"].(bool); ok && isVisualOnly {
			continue
		}

		estimate, err := s.CalculateResourceCost(ctx, res, duration)
		if err != nil {
			// Log error but continue with other resources
			// Some resources may not have pricing (e.g., VPC, subnets)
			fmt.Printf("  ⚠️  %s (%s): Pricing calculation skipped - %v\n", res.Name, res.Type.Name, err)
			continue
		}

		// Log resource cost calculation
		baseCost := 0.0
		hiddenCost := 0.0
		for _, comp := range estimate.Breakdown {
			baseCost += comp.Subtotal
		}
		for _, hiddenDep := range estimate.HiddenDependencyCosts {
			for _, comp := range hiddenDep.Breakdown {
				hiddenCost += comp.Subtotal
			}
		}

		fmt.Printf("  ✓ %s (%s): $%.2f", res.Name, res.Type.Name, estimate.TotalCost)
		if hiddenCost > 0 {
			fmt.Printf(" (Base: $%.2f + Hidden: $%.2f)", baseCost, hiddenCost)
		}
		fmt.Println()

		// Convert domain estimate to resource cost estimate
		breakdown := make([]serverinterfaces.CostBreakdownComponent, len(estimate.Breakdown))
		for i, comp := range estimate.Breakdown {
			breakdown[i] = serverinterfaces.CostBreakdownComponent{
				ComponentName: comp.ComponentName,
				Model:         string(comp.Model),
				Quantity:      comp.Quantity,
				UnitRate:      comp.UnitRate,
				Subtotal:      comp.Subtotal,
				Currency:      string(comp.Currency),
			}
		}

		// Include hidden dependency costs in breakdown
		for _, hiddenDep := range estimate.HiddenDependencyCosts {
			for _, comp := range hiddenDep.Breakdown {
				breakdown = append(breakdown, serverinterfaces.CostBreakdownComponent{
					ComponentName: fmt.Sprintf("%s (%s)", comp.ComponentName, hiddenDep.DependencyResourceType),
					Model:         string(comp.Model),
					Quantity:      comp.Quantity,
					UnitRate:      comp.UnitRate,
					Subtotal:      comp.Subtotal,
					Currency:      string(comp.Currency),
				})
			}
		}

		resourceEstimates[res.ID] = &serverinterfaces.ResourceCostEstimate{
			ResourceID:   res.ID,
			ResourceName: res.Name,
			ResourceType: res.Type.Name,
			TotalCost:    estimate.TotalCost,
			Currency:     string(estimate.Currency),
			Breakdown:    breakdown,
		}

		totalCost += estimate.TotalCost
	}

	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("  💰 Total Architecture Cost: $%.2f\n", totalCost)
	fmt.Println()

	// Determine period based on duration
	var period string
	if duration.Hours() <= 24 {
		period = "hourly"
	} else if duration.Hours() <= 720 {
		period = "monthly"
	} else {
		period = "yearly"
	}

	return &serverinterfaces.ArchitectureCostEstimate{
		TotalCost:         totalCost,
		Currency:          "USD",
		Period:            period,
		Duration:          duration,
		ResourceEstimates: resourceEstimates,
		Provider:          string(arch.Provider),
		Region:            arch.Region,
	}, nil
}

// PersistResourcePricing saves resource pricing to the database
func (s *PricingServiceImpl) PersistResourcePricing(ctx context.Context, projectID, resourceID uuid.UUID, estimate *domainpricing.CostEstimate, provider, region string) error {
	if estimate == nil {
		return fmt.Errorf("estimate is nil")
	}

	// Create resource pricing record
	resourcePricing := &models.ResourcePricing{
		ProjectID:       projectID,
		ResourceID:      resourceID,
		TotalCost:       estimate.TotalCost,
		Currency:        string(estimate.Currency),
		Period:          string(estimate.Period),
		DurationSeconds: int64(estimate.Duration.Seconds()),
		Provider:        provider,
		Region:          &region,
		CalculatedAt:    estimate.CalculatedAt,
	}

	if err := s.pricingRepo.CreateResourcePricing(ctx, resourcePricing); err != nil {
		return fmt.Errorf("failed to create resource pricing: %w", err)
	}

	// Create pricing components for base resource
	for _, comp := range estimate.Breakdown {
		component := &models.PricingComponent{
			ResourcePricingID: resourcePricing.ID,
			ComponentName:     comp.ComponentName,
			Model:             string(comp.Model),
			Unit:              getUnitFromModel(comp.Model),
			Quantity:          comp.Quantity,
			UnitRate:          comp.UnitRate,
			Subtotal:          comp.Subtotal,
			Currency:          string(comp.Currency),
		}

		if err := s.pricingRepo.CreatePricingComponent(ctx, component); err != nil {
			return fmt.Errorf("failed to create pricing component: %w", err)
		}
	}

	// Create pricing components for hidden dependencies
	for _, hiddenDep := range estimate.HiddenDependencyCosts {
		for _, comp := range hiddenDep.Breakdown {
			component := &models.PricingComponent{
				ResourcePricingID: resourcePricing.ID,
				ComponentName:     fmt.Sprintf("%s (%s)", comp.ComponentName, hiddenDep.DependencyResourceType),
				Model:             string(comp.Model),
				Unit:              getUnitFromModel(comp.Model),
				Quantity:          comp.Quantity,
				UnitRate:          comp.UnitRate,
				Subtotal:          comp.Subtotal,
				Currency:          string(comp.Currency),
			}

			if err := s.pricingRepo.CreatePricingComponent(ctx, component); err != nil {
				return fmt.Errorf("failed to create hidden dependency pricing component: %w", err)
			}
		}
	}

	return nil
}

// PersistProjectPricing saves project-level pricing to the database
func (s *PricingServiceImpl) PersistProjectPricing(ctx context.Context, projectID uuid.UUID, estimate *domainpricing.CostEstimate, provider, region string) error {
	if estimate == nil {
		return fmt.Errorf("estimate is nil")
	}

	projectPricing := &models.ProjectPricing{
		ProjectID:       projectID,
		TotalCost:       estimate.TotalCost,
		Currency:        string(estimate.Currency),
		Period:          string(estimate.Period),
		DurationSeconds: int64(estimate.Duration.Seconds()),
		Provider:        provider,
		Region:          &region,
		CalculatedAt:    estimate.CalculatedAt,
	}

	if err := s.pricingRepo.CreateProjectPricing(ctx, projectPricing); err != nil {
		return fmt.Errorf("failed to create project pricing: %w", err)
	}

	return nil
}

// GetProjectPricing retrieves pricing for a project
func (s *PricingServiceImpl) GetProjectPricing(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectPricing, error) {
	return s.pricingRepo.FindProjectPricingByProjectID(ctx, projectID)
}

// GetResourcePricing retrieves pricing for a resource
func (s *PricingServiceImpl) GetResourcePricing(ctx context.Context, resourceID uuid.UUID) ([]*models.ResourcePricing, error) {
	return s.pricingRepo.FindResourcePricingByResourceID(ctx, resourceID)
}

// getUnitFromModel returns the unit string based on the pricing model
func getUnitFromModel(model domainpricing.PricingModel) string {
	switch model {
	case domainpricing.PerHour:
		return "hour"
	case domainpricing.PerGB:
		return "GB"
	case domainpricing.PerRequest:
		return "request"
	case domainpricing.OneTime:
		return "unit"
	case domainpricing.Tiered:
		return "unit"
	case domainpricing.Percentage:
		return "percent"
	default:
		return "unit"
	}
}
