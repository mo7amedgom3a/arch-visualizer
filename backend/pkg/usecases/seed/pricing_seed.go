package seed

import (
	"context"
	"fmt"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
	"gorm.io/datatypes"
)

// SeedPricingData seeds pricing rates and hidden dependencies
func SeedPricingData(ctx context.Context) error {
	fmt.Println("\n💰 Seeding pricing data...")

	// Seed pricing rates
	if err := seedPricingRates(ctx); err != nil {
		return fmt.Errorf("failed to seed pricing rates: %w", err)
	}

	// Seed hidden dependencies
	if err := seedHiddenDependencies(ctx); err != nil {
		return fmt.Errorf("failed to seed hidden dependencies: %w", err)
	}

	fmt.Println("  ✓ Pricing data seeded successfully")
	return nil
}

// seedPricingRates seeds pricing rates for AWS resources
func seedPricingRates(ctx context.Context) error {
	pricingRateRepo, err := repository.NewPricingRateRepository()
	if err != nil {
		return fmt.Errorf("failed to create pricing rate repository: %w", err)
	}

	now := time.Now()
	rates := []*models.PricingRate{
		// NAT Gateway
		{
			Provider:      "aws",
			ResourceType:  "nat_gateway",
			ComponentName: "NAT Gateway Hourly",
			PricingModel:  "per_hour",
			Unit:          "hour",
			Rate:          0.045,
			Currency:      "USD",
			Region:        nil, // Default for all regions
			EffectiveFrom: now,
		},
		{
			Provider:      "aws",
			ResourceType:  "nat_gateway",
			ComponentName: "NAT Gateway Data Processing",
			PricingModel:  "per_gb",
			Unit:          "GB",
			Rate:          0.045,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
		},
		// Elastic IP
		{
			Provider:      "aws",
			ResourceType:  "elastic_ip",
			ComponentName: "Elastic IP Hourly (Unattached)",
			PricingModel:  "per_hour",
			Unit:          "hour",
			Rate:          0.005,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
		},
		// EC2 Instance - t3.micro
		{
			Provider:      "aws",
			ResourceType:  "ec2_instance",
			ComponentName: "EC2 Instance Hourly",
			PricingModel:  "per_hour",
			Unit:          "hour",
			Rate:          0.0104,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
			Metadata:      datatypes.JSON(`{"instance_type":"t3.micro"}`),
		},
		// EC2 Instance - t3.small
		{
			Provider:      "aws",
			ResourceType:  "ec2_instance",
			ComponentName: "EC2 Instance Hourly",
			PricingModel:  "per_hour",
			Unit:          "hour",
			Rate:          0.0208,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
			Metadata:      datatypes.JSON(`{"instance_type":"t3.small"}`),
		},
		// EC2 Instance - m5.large
		{
			Provider:      "aws",
			ResourceType:  "ec2_instance",
			ComponentName: "EC2 Instance Hourly",
			PricingModel:  "per_hour",
			Unit:          "hour",
			Rate:          0.096,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
			Metadata:      datatypes.JSON(`{"instance_type":"m5.large"}`),
		},
		// EBS Volume - gp3
		{
			Provider:      "aws",
			ResourceType:  "ebs_volume",
			ComponentName: "EBS Volume Storage",
			PricingModel:  "per_gb",
			Unit:          "GB",
			Rate:          0.08, // Per GB per month
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
			Metadata:      datatypes.JSON(`{"volume_type":"gp3"}`),
		},
		// Load Balancer - Application
		{
			Provider:      "aws",
			ResourceType:  "load_balancer",
			ComponentName: "Load Balancer Hourly",
			PricingModel:  "per_hour",
			Unit:          "hour",
			Rate:          0.0225,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
			Metadata:      datatypes.JSON(`{"load_balancer_type":"application"}`),
		},
		// Load Balancer - Network
		{
			Provider:      "aws",
			ResourceType:  "load_balancer",
			ComponentName: "Load Balancer Hourly",
			PricingModel:  "per_hour",
			Unit:          "hour",
			Rate:          0.0225,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
			Metadata:      datatypes.JSON(`{"load_balancer_type":"network"}`),
		},
		// S3 Bucket - Standard Storage
		{
			Provider:      "aws",
			ResourceType:  "s3_bucket",
			ComponentName: "S3 Storage",
			PricingModel:  "per_gb",
			Unit:          "GB",
			Rate:          0.023, // Per GB per month
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
			Metadata:      datatypes.JSON(`{"storage_class":"standard"}`),
		},
		// S3 Bucket - PUT Requests
		{
			Provider:      "aws",
			ResourceType:  "s3_bucket",
			ComponentName: "S3 PUT Requests",
			PricingModel:  "per_request",
			Unit:          "request",
			Rate:          0.005, // Per 1000 requests
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
		},
		// S3 Bucket - GET Requests
		{
			Provider:      "aws",
			ResourceType:  "s3_bucket",
			ComponentName: "S3 GET Requests",
			PricingModel:  "per_request",
			Unit:          "request",
			Rate:          0.0004, // Per 1000 requests
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
		},
		// S3 Bucket - Data Transfer Out
		{
			Provider:      "aws",
			ResourceType:  "s3_bucket",
			ComponentName: "S3 Data Transfer Out",
			PricingModel:  "per_gb",
			Unit:          "GB",
			Rate:          0.09,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
		},
		// Lambda - Compute (GB-seconds)
		{
			Provider:      "aws",
			ResourceType:  "lambda_function",
			ComponentName: "Lambda Compute",
			PricingModel:  "per_gb",
			Unit:          "GB-second",
			Rate:          0.0000166667, // Per GB-second
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
		},
		// Lambda - Requests
		{
			Provider:      "aws",
			ResourceType:  "lambda_function",
			ComponentName: "Lambda Requests",
			PricingModel:  "per_request",
			Unit:          "request",
			Rate:          0.20, // Per million requests
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
		},
		// Lambda - Data Transfer
		{
			Provider:      "aws",
			ResourceType:  "lambda_function",
			ComponentName: "Lambda Data Transfer Out",
			PricingModel:  "per_gb",
			Unit:          "GB",
			Rate:          0.09,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
		},
		// RDS Instance - db.t3.micro
		{
			Provider:      "aws",
			ResourceType:  "rds_instance",
			ComponentName: "RDS Instance Hourly",
			PricingModel:  "per_hour",
			Unit:          "hour",
			Rate:          0.017,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
			Metadata:      datatypes.JSON(`{"instance_class":"db.t3.micro"}`),
		},
		// RDS Instance - db.t3.small
		{
			Provider:      "aws",
			ResourceType:  "rds_instance",
			ComponentName: "RDS Instance Hourly",
			PricingModel:  "per_hour",
			Unit:          "hour",
			Rate:          0.034,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
			Metadata:      datatypes.JSON(`{"instance_class":"db.t3.small"}`),
		},
		// RDS Instance - db.t3.medium
		{
			Provider:      "aws",
			ResourceType:  "rds_instance",
			ComponentName: "RDS Instance Hourly",
			PricingModel:  "per_hour",
			Unit:          "hour",
			Rate:          0.068,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
			Metadata:      datatypes.JSON(`{"instance_class":"db.t3.medium"}`),
		},
		// RDS Instance - db.m5.large
		{
			Provider:      "aws",
			ResourceType:  "rds_instance",
			ComponentName: "RDS Instance Hourly",
			PricingModel:  "per_hour",
			Unit:          "hour",
			Rate:          0.176,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
			Metadata:      datatypes.JSON(`{"instance_class":"db.m5.large"}`),
		},
		// RDS Instance - db.r5.large
		{
			Provider:      "aws",
			ResourceType:  "rds_instance",
			ComponentName: "RDS Instance Hourly",
			PricingModel:  "per_hour",
			Unit:          "hour",
			Rate:          0.24,
			Currency:      "USD",
			Region:        nil,
			EffectiveFrom: now,
			Metadata:      datatypes.JSON(`{"instance_class":"db.r5.large"}`),
		},
	}

	for _, rate := range rates {
		// Check if rate already exists
		existing, _ := pricingRateRepo.FindActiveRates(ctx, rate.Provider, rate.ResourceType, rate.Region)
		exists := false
		for _, e := range existing {
			if e.ComponentName == rate.ComponentName {
				exists = true
				break
			}
		}
		if !exists {
			if err := pricingRateRepo.Create(ctx, rate); err != nil {
				return fmt.Errorf("failed to create pricing rate for %s/%s: %w", rate.ResourceType, rate.ComponentName, err)
			}
		}
	}

	fmt.Printf("  ✓ Seeded pricing rates\n")
	return nil
}

// seedHiddenDependencies seeds hidden dependency definitions
func seedHiddenDependencies(ctx context.Context) error {
	hiddenDepRepo, err := repository.NewHiddenDependencyRepository()
	if err != nil {
		return fmt.Errorf("failed to create hidden dependency repository: %w", err)
	}

	deps := []*models.HiddenDependency{
		// NAT Gateway -> Elastic IP
		{
			Provider:            "aws",
			ParentResourceType:  "nat_gateway",
			ChildResourceType:   "elastic_ip",
			QuantityExpression:  "1",
			ConditionExpression: "metadata.allocationId == null",
			IsAttached:          true,
			Description:         "NAT Gateway requires an Elastic IP. If allocationId is not provided, one is automatically created and attached (free when attached).",
		},
		// EC2 -> EBS Root Volume
		{
			Provider:            "aws",
			ParentResourceType:  "ec2_instance",
			ChildResourceType:   "ebs_volume",
			QuantityExpression:  "metadata.size_gb",
			ConditionExpression: "",
			IsAttached:          true,
			Description:         "EC2 instance requires a root EBS volume. Default size is 8GB if not specified in metadata.size_gb.",
		},
		// EC2 -> Network Interface
		{
			Provider:            "aws",
			ParentResourceType:  "ec2_instance",
			ChildResourceType:   "network_interface",
			QuantityExpression:  "1",
			ConditionExpression: "",
			IsAttached:          true,
			Description:         "EC2 instance requires a network interface (free when attached).",
		},
		// RDS -> EBS Storage
		{
			Provider:            "aws",
			ParentResourceType:  "rds_instance",
			ChildResourceType:   "ebs_volume",
			QuantityExpression:  "metadata.allocated_storage",
			ConditionExpression: "",
			IsAttached:          true,
			Description:         "RDS instance requires storage volume based on allocated_storage. Default is 20GB if not specified.",
		},
		// RDS -> S3 Backup
		{
			Provider:            "aws",
			ParentResourceType:  "rds_instance",
			ChildResourceType:   "s3_bucket",
			QuantityExpression:  "metadata.allocated_storage",
			ConditionExpression: "metadata.backup_retention_period > 0",
			IsAttached:          false,
			Description:         "RDS automated backups stored in S3 (assumed equal to DB size for estimation).",
		},
	}

	for _, dep := range deps {
		// Check if dependency already exists
		existing, _ := hiddenDepRepo.FindByParentResourceType(ctx, dep.Provider, dep.ParentResourceType)
		exists := false
		for _, e := range existing {
			if e.ChildResourceType == dep.ChildResourceType {
				exists = true
				break
			}
		}
		if !exists {
			if err := hiddenDepRepo.Create(ctx, dep); err != nil {
				return fmt.Errorf("failed to create hidden dependency %s -> %s: %w", dep.ParentResourceType, dep.ChildResourceType, err)
			}
		}
	}

	fmt.Printf("  ✓ Seeded hidden dependencies\n")
	return nil
}
