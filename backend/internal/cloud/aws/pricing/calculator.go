package pricing

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/inventory"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing/compute"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing/networking"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/pricing/storage"
	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// AWSPricingCalculator implements the PricingCalculator interface for AWS
type AWSPricingCalculator struct {
	service *AWSPricingService
}

// NewAWSPricingCalculator creates a new AWS pricing calculator
func NewAWSPricingCalculator(service *AWSPricingService) *AWSPricingCalculator {
	return &AWSPricingCalculator{
		service: service,
	}
}

// CalculateResourceCost calculates the cost for a single resource over a given duration
func (c *AWSPricingCalculator) CalculateResourceCost(ctx context.Context, res *resource.Resource, duration time.Duration) (*domainpricing.CostEstimate, error) {
	if res.Provider != "aws" {
		return nil, fmt.Errorf("unsupported provider: %s", res.Provider)
	}

	// Try to use inventory first
	inv := inventory.GetDefaultInventory()
	if functions, ok := inv.GetFunctions(res.Type.Name); ok && functions.PricingCalculator != nil {
		return functions.PricingCalculator(res, duration)
	}

	// Fallback to switch-based calculation
	return c.calculateResourceCostFallback(ctx, res, duration)
}

// calculateResourceCostFallback provides backward compatibility with switch-based calculation
func (c *AWSPricingCalculator) calculateResourceCostFallback(ctx context.Context, res *resource.Resource, duration time.Duration) (*domainpricing.CostEstimate, error) {
	// Get pricing information for the resource type
	pricingInfo, err := c.GetResourcePricing(ctx, res.Type.Name, "aws", res.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to get pricing for resource type %s: %w", res.Type.Name, err)
	}

	// Calculate cost based on resource type
	var totalCost float64
	var breakdown []domainpricing.CostComponent

	switch res.Type.Name {
	case "nat_gateway":
		// For NAT Gateway, we need to estimate data processing
		// In a real scenario, this would come from usage metrics
		// For now, we default to 0 - this can be extended to accept usage metrics
		estimatedDataGB := 0.0
		cost := networking.CalculateNATGatewayCost(duration, estimatedDataGB, res.Region)
		totalCost = cost
		breakdown = []domainpricing.CostComponent{
			{
				ComponentName: "NAT Gateway Hourly",
				Model:         domainpricing.PerHour,
				Quantity:      duration.Hours(),
				UnitRate:      pricingInfo.Components[0].Rate,
				Subtotal:      pricingInfo.Components[0].Rate * duration.Hours(),
				Currency:      domainpricing.USD,
			},
		}
		if estimatedDataGB > 0 {
			breakdown = append(breakdown, domainpricing.CostComponent{
				ComponentName: "NAT Gateway Data Processing",
				Model:         domainpricing.PerGB,
				Quantity:      estimatedDataGB,
				UnitRate:      pricingInfo.Components[1].Rate,
				Subtotal:      pricingInfo.Components[1].Rate * estimatedDataGB,
				Currency:      domainpricing.USD,
			})
		}

	case "elastic_ip":
		// Check if EIP is attached
		// For now, we default to false (unattached) which incurs charges
		// In a real scenario, this would come from resource state
		isAttached := false
		cost := networking.CalculateElasticIPCost(duration, isAttached, res.Region)
		totalCost = cost
		if !isAttached {
			breakdown = []domainpricing.CostComponent{
				{
					ComponentName: "Elastic IP Hourly (Unattached)",
					Model:         domainpricing.PerHour,
					Quantity:      duration.Hours(),
					UnitRate:      pricingInfo.Components[0].Rate,
					Subtotal:      cost,
					Currency:      domainpricing.USD,
				},
			}
		}

	case "network_interface":
		// Check if ENI is attached
		// For now, we default to false (unattached) which incurs charges
		// In a real scenario, this would come from resource state
		isAttached := false
		cost := networking.CalculateNetworkInterfaceCost(duration, isAttached, res.Region)
		totalCost = cost
		if !isAttached {
			breakdown = []domainpricing.CostComponent{
				{
					ComponentName: "Network Interface Hourly (Unattached)",
					Model:         domainpricing.PerHour,
					Quantity:      duration.Hours(),
					UnitRate:      pricingInfo.Components[0].Rate,
					Subtotal:      cost,
					Currency:      domainpricing.USD,
				},
			}
		}

	case "ec2_instance":
		// Extract instance type from metadata
		instanceType := "t3.micro" // Default
		if res.Metadata != nil {
			if it, ok := res.Metadata["instance_type"].(string); ok && it != "" {
				instanceType = it
			}
		}

		// Get pricing for the specific instance type
		instancePricing := compute.GetEC2InstancePricing(instanceType, res.Region)
		hourlyRate := instancePricing.Components[0].Rate

		cost := compute.CalculateEC2InstanceCost(duration, instanceType, res.Region)
		totalCost = cost
		breakdown = []domainpricing.CostComponent{
			{
				ComponentName: "EC2 Instance Hourly",
				Model:         domainpricing.PerHour,
				Quantity:      duration.Hours(),
				UnitRate:      hourlyRate,
				Subtotal:      cost,
				Currency:      domainpricing.USD,
			},
		}

	case "ebs_volume":
		// Extract size and volume type from metadata
		sizeGB := 0.0
		volumeType := "gp3" // Default
		if res.Metadata != nil {
			if s, ok := res.Metadata["size_gb"].(float64); ok {
				sizeGB = s
			} else if s, ok := res.Metadata["size_gb"].(int); ok {
				sizeGB = float64(s)
			}
			if vt, ok := res.Metadata["volume_type"].(string); ok && vt != "" {
				volumeType = vt
			}
		}

		if sizeGB <= 0 {
			return nil, fmt.Errorf("ebs_volume requires size_gb in metadata")
		}

		// Get pricing for the specific volume type
		volumePricing := storage.GetEBSVolumePricing(volumeType, res.Region)
		ratePerGBMonth := volumePricing.Components[0].Rate

		cost := storage.CalculateEBSVolumeCost(duration, sizeGB, volumeType, res.Region)
		totalCost = cost

		// Calculate months for breakdown
		hoursPerMonth := 720.0
		months := duration.Hours() / hoursPerMonth

		breakdown = []domainpricing.CostComponent{
			{
				ComponentName: "EBS Volume Storage",
				Model:         domainpricing.PerGB,
				Quantity:      sizeGB * months,
				UnitRate:      ratePerGBMonth,
				Subtotal:      cost,
				Currency:      domainpricing.USD,
			},
		}

	case "load_balancer":
		// Extract load balancer type from metadata
		lbType := "application" // Default to ALB
		if res.Metadata != nil {
			if lt, ok := res.Metadata["load_balancer_type"].(string); ok && lt != "" {
				lbType = lt
			}
		}

		// Get pricing for the specific LB type
		lbPricing := compute.GetLoadBalancerPricing(lbType, res.Region)
		hourlyRate := lbPricing.Components[0].Rate

		cost := compute.CalculateLoadBalancerCost(duration, lbType, res.Region)
		totalCost = cost
		breakdown = []domainpricing.CostComponent{
			{
				ComponentName: "Load Balancer Hourly",
				Model:         domainpricing.PerHour,
				Quantity:      duration.Hours(),
				UnitRate:      hourlyRate,
				Subtotal:      cost,
				Currency:      domainpricing.USD,
			},
		}

	case "auto_scaling_group":
		// Extract instance type and capacity from metadata
		instanceType := "t3.micro" // Default
		minSize := 1
		maxSize := 3
		if res.Metadata != nil {
			if it, ok := res.Metadata["instance_type"].(string); ok && it != "" {
				instanceType = it
			}
			if ms, ok := res.Metadata["min_size"].(int); ok {
				minSize = ms
			} else if ms, ok := res.Metadata["min_size"].(float64); ok {
				minSize = int(ms)
			}
			if ms, ok := res.Metadata["max_size"].(int); ok {
				maxSize = ms
			} else if ms, ok := res.Metadata["max_size"].(float64); ok {
				maxSize = int(ms)
			}
		}

		// Get pricing for the ASG
		asgPricing := compute.GetAutoScalingGroupPricing(instanceType, minSize, maxSize, res.Region)
		hourlyRate := asgPricing.Components[0].Rate

		cost := compute.CalculateAutoScalingGroupCost(duration, instanceType, minSize, maxSize, res.Region)
		totalCost = cost

		breakdown = []domainpricing.CostComponent{
			{
				ComponentName: "Auto Scaling Group Hourly",
				Model:         domainpricing.PerHour,
				Quantity:      duration.Hours(),
				UnitRate:      hourlyRate,
				Subtotal:      cost,
				Currency:      domainpricing.USD,
			},
		}

	case "s3_bucket":
		// Extract S3 bucket parameters from metadata
		sizeGB := 0.0
		storageClass := "standard" // Default
		putRequests := 0.0
		getRequests := 0.0
		dataTransferGB := 0.0

		if res.Metadata != nil {
			if s, ok := res.Metadata["size_gb"].(float64); ok {
				sizeGB = s
			} else if s, ok := res.Metadata["size_gb"].(int); ok {
				sizeGB = float64(s)
			}
			if sc, ok := res.Metadata["storage_class"].(string); ok && sc != "" {
				storageClass = sc
			}
			if pr, ok := res.Metadata["put_requests"].(float64); ok {
				putRequests = pr
			} else if pr, ok := res.Metadata["put_requests"].(int); ok {
				putRequests = float64(pr)
			}
			if gr, ok := res.Metadata["get_requests"].(float64); ok {
				getRequests = gr
			} else if gr, ok := res.Metadata["get_requests"].(int); ok {
				getRequests = float64(gr)
			}
			if dt, ok := res.Metadata["data_transfer_gb"].(float64); ok {
				dataTransferGB = dt
			} else if dt, ok := res.Metadata["data_transfer_gb"].(int); ok {
				dataTransferGB = float64(dt)
			}
		}

		// Get pricing for S3 bucket
		s3Pricing := storage.GetS3BucketPricing(storageClass, res.Region)

		// Calculate total cost
		cost := storage.CalculateS3BucketCost(
			duration,
			sizeGB,
			putRequests,
			getRequests,
			dataTransferGB,
			storageClass,
			res.Region,
		)
		totalCost = cost

		// Build breakdown with multiple components
		hoursPerMonth := 720.0
		months := duration.Hours() / hoursPerMonth

		breakdown = []domainpricing.CostComponent{}

		// Storage component
		if sizeGB > 0 {
			storageRate := s3Pricing.Components[0].Rate
			storageCost := storageRate * sizeGB * months
			breakdown = append(breakdown, domainpricing.CostComponent{
				ComponentName: "S3 Storage",
				Model:         domainpricing.PerGB,
				Quantity:      sizeGB * months,
				UnitRate:      storageRate,
				Subtotal:      storageCost,
				Currency:      domainpricing.USD,
			})
		}

		// PUT requests component
		if putRequests > 0 {
			putRate := s3Pricing.Components[1].Rate
			putCost := (putRate / 1000.0) * putRequests
			breakdown = append(breakdown, domainpricing.CostComponent{
				ComponentName: "S3 PUT Requests",
				Model:         domainpricing.PerRequest,
				Quantity:      putRequests,
				UnitRate:      putRate / 1000.0,
				Subtotal:      putCost,
				Currency:      domainpricing.USD,
			})
		}

		// GET requests component
		if getRequests > 0 {
			getRate := s3Pricing.Components[2].Rate
			getCost := (getRate / 1000.0) * getRequests
			breakdown = append(breakdown, domainpricing.CostComponent{
				ComponentName: "S3 GET Requests",
				Model:         domainpricing.PerRequest,
				Quantity:      getRequests,
				UnitRate:      getRate / 1000.0,
				Subtotal:      getCost,
				Currency:      domainpricing.USD,
			})
		}

		// Data transfer component
		if dataTransferGB > 0 {
			dataTransferRate := s3Pricing.Components[3].Rate
			// First 1GB per month is free
			freeTierPerMonth := 1.0
			chargeableGB := 0.0
			if dataTransferGB > (freeTierPerMonth * months) {
				chargeableGB = dataTransferGB - (freeTierPerMonth * months)
			}
			if chargeableGB > 0 {
				dataTransferCost := dataTransferRate * chargeableGB
				breakdown = append(breakdown, domainpricing.CostComponent{
					ComponentName: "S3 Data Transfer Out",
					Model:         domainpricing.PerGB,
					Quantity:      chargeableGB,
					UnitRate:      dataTransferRate,
					Subtotal:      dataTransferCost,
					Currency:      domainpricing.USD,
				})
			}
		}

	case "lambda_function":
		// Extract Lambda function parameters from metadata
		memorySizeMB := 128.0      // Default
		averageDurationMs := 100.0 // Default 100ms
		requestCount := 0.0
		dataTransferGB := 0.0

		if res.Metadata != nil {
			if m, ok := res.Metadata["memory_size_mb"].(float64); ok && m > 0 {
				memorySizeMB = m
			} else if m, ok := res.Metadata["memory_size_mb"].(int); ok && m > 0 {
				memorySizeMB = float64(m)
			}
			if d, ok := res.Metadata["average_duration_ms"].(float64); ok && d > 0 {
				averageDurationMs = d
			} else if d, ok := res.Metadata["average_duration_ms"].(int); ok && d > 0 {
				averageDurationMs = float64(d)
			}
			if rc, ok := res.Metadata["request_count"].(float64); ok {
				requestCount = rc
			} else if rc, ok := res.Metadata["request_count"].(int); ok {
				requestCount = float64(rc)
			}
			if dt, ok := res.Metadata["data_transfer_gb"].(float64); ok {
				dataTransferGB = dt
			} else if dt, ok := res.Metadata["data_transfer_gb"].(int); ok {
				dataTransferGB = float64(dt)
			}
		}

		// Get pricing for Lambda function
		lambdaPricing := compute.GetLambdaFunctionPricing(memorySizeMB, res.Region)

		// Calculate total cost
		cost := compute.CalculateLambdaFunctionCost(
			duration,
			memorySizeMB,
			averageDurationMs,
			requestCount,
			dataTransferGB,
			res.Region,
		)
		totalCost = cost

		// Build breakdown with multiple components
		hoursPerMonth := 720.0
		months := duration.Hours() / hoursPerMonth

		breakdown = []domainpricing.CostComponent{}

		// Compute component (GB-seconds)
		if requestCount > 0 {
			memorySizeGB := memorySizeMB / 1024.0
			durationSeconds := averageDurationMs / 1000.0
			totalGBSeconds := memorySizeGB * durationSeconds * requestCount
			computeRate := lambdaPricing.Components[0].Rate
			computeCost := computeRate * totalGBSeconds
			breakdown = append(breakdown, domainpricing.CostComponent{
				ComponentName: "Lambda Compute",
				Model:         domainpricing.PerGB,
				Quantity:      totalGBSeconds,
				UnitRate:      computeRate,
				Subtotal:      computeCost,
				Currency:      domainpricing.USD,
			})
		}

		// Request component
		if requestCount > 0 {
			requestRate := lambdaPricing.Components[1].Rate
			// Free tier: 1M requests per month
			freeTierPerMonth := 1000000.0
			freeTierTotal := freeTierPerMonth * months
			chargeableRequests := math.Max(0, requestCount-freeTierTotal)
			if chargeableRequests > 0 {
				requestCost := (requestRate / 1000000.0) * chargeableRequests
				breakdown = append(breakdown, domainpricing.CostComponent{
					ComponentName: "Lambda Requests",
					Model:         domainpricing.PerRequest,
					Quantity:      chargeableRequests,
					UnitRate:      requestRate / 1000000.0,
					Subtotal:      requestCost,
					Currency:      domainpricing.USD,
				})
			}
		}

		// Data transfer component
		if dataTransferGB > 0 {
			dataTransferRate := lambdaPricing.Components[2].Rate
			// First 1GB per month is free
			freeTierPerMonth := 1.0
			chargeableGB := 0.0
			if dataTransferGB > (freeTierPerMonth * months) {
				chargeableGB = dataTransferGB - (freeTierPerMonth * months)
			}
			if chargeableGB > 0 {
				dataTransferCost := dataTransferRate * chargeableGB
				breakdown = append(breakdown, domainpricing.CostComponent{
					ComponentName: "Lambda Data Transfer Out",
					Model:         domainpricing.PerGB,
					Quantity:      chargeableGB,
					UnitRate:      dataTransferRate,
					Subtotal:      dataTransferCost,
					Currency:      domainpricing.USD,
				})
			}
		}

	default:
		// For other resource types, use generic calculation
		// This can be extended for other resource types
		return nil, fmt.Errorf("pricing calculation not yet implemented for resource type: %s", res.Type.Name)
	}

	// Determine period based on duration
	var period domainpricing.Period
	if duration.Hours() <= 24 {
		period = domainpricing.Hourly
	} else if duration.Hours() <= 720 {
		period = domainpricing.Monthly
	} else {
		period = domainpricing.Yearly
	}

	return &domainpricing.CostEstimate{
		TotalCost:    totalCost,
		Currency:     domainpricing.USD,
		Breakdown:    breakdown,
		Period:       period,
		Duration:     duration,
		CalculatedAt: time.Now(),
		ResourceType: &res.Type.Name,
		Provider:     domainpricing.AWS,
		Region:       &res.Region,
	}, nil
}

// CalculateArchitectureCost calculates the total cost for multiple resources over a given duration
func (c *AWSPricingCalculator) CalculateArchitectureCost(ctx context.Context, resources []*resource.Resource, duration time.Duration) (*domainpricing.CostEstimate, error) {
	var totalCost float64
	var allBreakdown []domainpricing.CostComponent

	for _, res := range resources {
		estimate, err := c.CalculateResourceCost(ctx, res, duration)
		if err != nil {
			// Log error but continue with other resources
			continue
		}
		totalCost += estimate.TotalCost
		allBreakdown = append(allBreakdown, estimate.Breakdown...)
	}

	// Determine period based on duration
	var period domainpricing.Period
	if duration.Hours() <= 24 {
		period = domainpricing.Hourly
	} else if duration.Hours() <= 720 {
		period = domainpricing.Monthly
	} else {
		period = domainpricing.Yearly
	}

	return &domainpricing.CostEstimate{
		TotalCost:    totalCost,
		Currency:     domainpricing.USD,
		Breakdown:    allBreakdown,
		Period:       period,
		Duration:     duration,
		CalculatedAt: time.Now(),
		Provider:     domainpricing.AWS,
	}, nil
}

// GetResourcePricing retrieves the pricing information for a specific resource type
func (c *AWSPricingCalculator) GetResourcePricing(ctx context.Context, resourceType string, provider string, region string) (*domainpricing.ResourcePricing, error) {
	return c.service.GetPricing(ctx, resourceType, provider, region)
}
