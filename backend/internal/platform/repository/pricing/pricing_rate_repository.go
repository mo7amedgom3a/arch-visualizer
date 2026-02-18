package pricingrepo

import (
	"context"
	"time"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"

	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	"gorm.io/gorm"
)

// PricingRateRepository provides operations for pricing rates
type PricingRateRepository struct {
	*repository.BaseRepository
}

// NewPricingRateRepository creates a new pricing rate repository
func NewPricingRateRepository() (*PricingRateRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, platformerrors.NewDatabaseConnectionFailed(err)
	}
	return &PricingRateRepository{BaseRepository: base}, nil
}

// Create creates a new pricing rate
func (r *PricingRateRepository) Create(ctx context.Context, rate *models.PricingRate) error {
	return r.GetDB(ctx).Create(rate).Error
}

// FindByID finds a pricing rate by ID
func (r *PricingRateRepository) FindByID(ctx context.Context, id uint) (*models.PricingRate, error) {
	var rate models.PricingRate
	err := r.GetDB(ctx).First(&rate, "id = ?", id).Error
	if err != nil {
		return nil, platformerrors.HandleGormError(err, "pricing_rate", "PricingRateRepository.FindByID")
	}
	return &rate, nil
}

// FindActiveRates finds active pricing rates for a resource type
func (r *PricingRateRepository) FindActiveRates(ctx context.Context, provider, resourceType string, region *string) ([]*models.PricingRate, error) {
	var rates []*models.PricingRate
	now := time.Now()

	query := r.GetDB(ctx).
		Where("provider = ?", provider).
		Where("resource_type = ?", resourceType).
		Where("effective_from <= ?", now).
		Where("(effective_to IS NULL OR effective_to > ?)", now)

	if region != nil {
		query = query.Where("(region IS NULL OR region = ?)", *region)
	} else {
		query = query.Where("region IS NULL")
	}

	err := query.Order("component_name ASC").Find(&rates).Error
	return rates, err
}

// FindByProviderAndResourceType finds all rates for a provider and resource type
func (r *PricingRateRepository) FindByProviderAndResourceType(ctx context.Context, provider, resourceType string) ([]*models.PricingRate, error) {
	var rates []*models.PricingRate
	err := r.GetDB(ctx).
		Where("provider = ?", provider).
		Where("resource_type = ?", resourceType).
		Order("effective_from DESC").
		Find(&rates).Error
	return rates, err
}

// Update updates a pricing rate
func (r *PricingRateRepository) Update(ctx context.Context, rate *models.PricingRate) error {
	return r.GetDB(ctx).Save(rate).Error
}

// Delete deletes a pricing rate
func (r *PricingRateRepository) Delete(ctx context.Context, id uint) error {
	return r.GetDB(ctx).Delete(&models.PricingRate{}, id).Error
}

// FindByInstanceType finds active pricing rates for a specific EC2 instance type
func (r *PricingRateRepository) FindByInstanceType(ctx context.Context, provider, instanceType, region, operatingSystem string) ([]*models.PricingRate, error) {
	var rates []*models.PricingRate
	now := time.Now()

	query := r.GetDB(ctx).
		Where("provider = ?", provider).
		Where("resource_type = ?", "ec2_instance").
		Where("instance_type = ?", instanceType).
		Where("operating_system = ?", operatingSystem).
		Where("effective_from <= ?", now).
		Where("(effective_to IS NULL OR effective_to > ?)", now)

	if region != "" {
		query = query.Where("(region IS NULL OR region = ?)", region)
	} else {
		query = query.Where("region IS NULL")
	}

	err := query.Order("component_name ASC").Find(&rates).Error
	return rates, err
}

// UpsertRate upserts a pricing rate (inserts if not exists, updates if exists)
func (r *PricingRateRepository) UpsertRate(ctx context.Context, rate *models.PricingRate) error {
	// Build unique key conditions
	conditions := map[string]interface{}{
		"provider":       rate.Provider,
		"resource_type":  rate.ResourceType,
		"component_name": rate.ComponentName,
	}

	if rate.Region != nil {
		conditions["region"] = *rate.Region
	} else {
		conditions["region"] = nil
	}

	if rate.InstanceType != nil {
		conditions["instance_type"] = *rate.InstanceType
	} else {
		conditions["instance_type"] = nil
	}

	if rate.OperatingSystem != nil {
		conditions["operating_system"] = *rate.OperatingSystem
	} else {
		conditions["operating_system"] = nil
	}

	// Check if rate exists
	var existing models.PricingRate
	query := r.GetDB(ctx)
	for k, v := range conditions {
		if v == nil {
			query = query.Where(k + " IS NULL")
		} else {
			query = query.Where(k+" = ?", v)
		}
	}

	err := query.First(&existing).Error
	if err != nil {
		// Not found, create new
		return r.GetDB(ctx).Create(rate).Error
	}

	// Update existing
	rate.ID = existing.ID
	return r.GetDB(ctx).Save(rate).Error
}

// BulkUpsert performs bulk upsert of pricing rates
func (r *PricingRateRepository) BulkUpsert(ctx context.Context, rates []*models.PricingRate) error {
	if len(rates) == 0 {
		return nil
	}

	// Use transaction for bulk operations
	return r.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		for _, rate := range rates {
			// Build unique key conditions
			conditions := map[string]interface{}{
				"provider":       rate.Provider,
				"resource_type":  rate.ResourceType,
				"component_name": rate.ComponentName,
			}

			if rate.Region != nil {
				conditions["region"] = *rate.Region
			} else {
				conditions["region"] = nil
			}

			if rate.InstanceType != nil {
				conditions["instance_type"] = *rate.InstanceType
			} else {
				conditions["instance_type"] = nil
			}

			if rate.OperatingSystem != nil {
				conditions["operating_system"] = *rate.OperatingSystem
			} else {
				conditions["operating_system"] = nil
			}

			// Check if rate exists
			var existing models.PricingRate
			query := tx
			for k, v := range conditions {
				if v == nil {
					query = query.Where(k + " IS NULL")
				} else {
					query = query.Where(k+" = ?", v)
				}
			}

			err := query.First(&existing).Error
			if err != nil {
				// Not found, create new
				if err := tx.Create(rate).Error; err != nil {
					return err
				}
			} else {
				// Update existing
				rate.ID = existing.ID
				if err := tx.Save(rate).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
