package pricingrepo

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// PricingRepository defines operations for pricing management
type PricingRepository struct {
	*repository.BaseRepository
}

// NewPricingRepository creates a new pricing repository
func NewPricingRepository() (*PricingRepository, error) {
	base, err := repository.NewBaseRepository()
	if err != nil {
		return nil, err
	}
	return &PricingRepository{BaseRepository: base}, nil
}

// CreateProjectPricing creates project-level pricing
func (r *PricingRepository) CreateProjectPricing(ctx context.Context, pricing *models.ProjectPricing) error {
	return r.GetDB(ctx).Create(pricing).Error
}

// FindProjectPricingByProjectID finds pricing for a project
func (r *PricingRepository) FindProjectPricingByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectPricing, error) {
	var pricings []*models.ProjectPricing
	err := r.GetDB(ctx).
		Where("project_id = ?", projectID).
		Order("calculated_at DESC").
		Find(&pricings).Error
	return pricings, err
}

// CreateServicePricing creates service-level pricing
func (r *PricingRepository) CreateServicePricing(ctx context.Context, pricing *models.ServicePricing) error {
	return r.GetDB(ctx).Create(pricing).Error
}

// FindServicePricingByProjectID finds service pricing for a project
func (r *PricingRepository) FindServicePricingByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ServicePricing, error) {
	var pricings []*models.ServicePricing
	err := r.GetDB(ctx).
		Where("project_id = ?", projectID).
		Preload("Category").
		Order("calculated_at DESC").
		Find(&pricings).Error
	return pricings, err
}

// CreateServiceTypePricing creates service type-level pricing
func (r *PricingRepository) CreateServiceTypePricing(ctx context.Context, pricing *models.ServiceTypePricing) error {
	return r.GetDB(ctx).Create(pricing).Error
}

// FindServiceTypePricingByProjectID finds service type pricing for a project
func (r *PricingRepository) FindServiceTypePricingByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ServiceTypePricing, error) {
	var pricings []*models.ServiceTypePricing
	err := r.GetDB(ctx).
		Where("project_id = ?", projectID).
		Preload("ResourceType").
		Order("calculated_at DESC").
		Find(&pricings).Error
	return pricings, err
}

// CreateResourcePricing creates resource-level pricing
func (r *PricingRepository) CreateResourcePricing(ctx context.Context, pricing *models.ResourcePricing) error {
	return r.GetDB(ctx).Create(pricing).Error
}

// FindResourcePricingByResourceID finds pricing for a resource
func (r *PricingRepository) FindResourcePricingByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*models.ResourcePricing, error) {
	var pricings []*models.ResourcePricing
	err := r.GetDB(ctx).
		Where("resource_id = ?", resourceID).
		Preload("Components").
		Order("calculated_at DESC").
		Find(&pricings).Error
	return pricings, err
}

// FindResourcePricingByProjectID finds all resource pricing for a project
func (r *PricingRepository) FindResourcePricingByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ResourcePricing, error) {
	var pricings []*models.ResourcePricing
	err := r.GetDB(ctx).
		Where("project_id = ?", projectID).
		Preload("Resource").
		Preload("Components").
		Order("calculated_at DESC").
		Find(&pricings).Error
	return pricings, err
}

// CreatePricingComponent creates a pricing component
func (r *PricingRepository) CreatePricingComponent(ctx context.Context, component *models.PricingComponent) error {
	return r.GetDB(ctx).Create(component).Error
}
