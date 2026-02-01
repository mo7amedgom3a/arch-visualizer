package pricing_importer

import (
	"context"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
)

// Importer handles importing pricing data from scraper output
type Importer struct {
	pricingRateRepo *repository.PricingRateRepository
}

// NewImporter creates a new pricing importer
func NewImporter() (*Importer, error) {
	repo, err := repository.NewPricingRateRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create pricing rate repository: %w", err)
	}

	return &Importer{
		pricingRateRepo: repo,
	}, nil
}

// ImportEC2Pricing imports EC2 pricing data from a scraper JSON file
func (i *Importer) ImportEC2Pricing(ctx context.Context, filePath string) (*ImportStats, error) {
	stats := &ImportStats{
		RegionsProcessed: make(map[string]int),
		OSProcessed:      make(map[string]int),
		Errors:           []string{},
	}

	// Parse instances from JSON file
	instances, err := ParseEC2Instances(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EC2 instances: %w", err)
	}

	stats.TotalInstances = len(instances)

	// Convert to pricing rates
	rates, err := ConvertToPricingRates(instances)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to pricing rates: %w", err)
	}

	stats.TotalRates = len(rates)

	// Track regions and OS
	for _, rate := range rates {
		if rate.Region != nil {
			stats.RegionsProcessed[*rate.Region]++
		}
		if rate.OperatingSystem != nil {
			stats.OSProcessed[*rate.OperatingSystem]++
		}
	}

	// Bulk upsert to database
	if err := i.pricingRateRepo.BulkUpsert(ctx, rates); err != nil {
		stats.Errors = append(stats.Errors, fmt.Sprintf("bulk upsert failed: %v", err))
		return stats, fmt.Errorf("failed to bulk upsert rates: %w", err)
	}

	return stats, nil
}

// ImportEC2PricingFromReader imports EC2 pricing data from a reader (for testing)
func (i *Importer) ImportEC2PricingFromReader(ctx context.Context, reader func() ([]EC2Instance, error)) (*ImportStats, error) {
	stats := &ImportStats{
		RegionsProcessed: make(map[string]int),
		OSProcessed:      make(map[string]int),
		Errors:           []string{},
	}

	// Parse instances
	instances, err := reader()
	if err != nil {
		return nil, fmt.Errorf("failed to read instances: %w", err)
	}

	stats.TotalInstances = len(instances)

	// Convert to pricing rates
	rates, err := ConvertToPricingRates(instances)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to pricing rates: %w", err)
	}

	stats.TotalRates = len(rates)

	// Track regions and OS
	for _, rate := range rates {
		if rate.Region != nil {
			stats.RegionsProcessed[*rate.Region]++
		}
		if rate.OperatingSystem != nil {
			stats.OSProcessed[*rate.OperatingSystem]++
		}
	}

	// Bulk upsert to database
	if err := i.pricingRateRepo.BulkUpsert(ctx, rates); err != nil {
		stats.Errors = append(stats.Errors, fmt.Sprintf("bulk upsert failed: %v", err))
		return stats, fmt.Errorf("failed to bulk upsert rates: %w", err)
	}

	return stats, nil
}
