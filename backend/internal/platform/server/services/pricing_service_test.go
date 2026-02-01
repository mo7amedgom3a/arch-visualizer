package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/architecture"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// mockPricingRepository implements PricingRepository interface for testing
type mockPricingRepository struct {
	createProjectPricingFunc            func(ctx context.Context, pricing *models.ProjectPricing) error
	findProjectPricingByProjectIDFunc   func(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectPricing, error)
	createResourcePricingFunc           func(ctx context.Context, pricing *models.ResourcePricing) error
	findResourcePricingByResourceIDFunc func(ctx context.Context, resourceID uuid.UUID) ([]*models.ResourcePricing, error)
	findResourcePricingByProjectIDFunc  func(ctx context.Context, projectID uuid.UUID) ([]*models.ResourcePricing, error)
	createPricingComponentFunc          func(ctx context.Context, component *models.PricingComponent) error
}

func (m *mockPricingRepository) CreateProjectPricing(ctx context.Context, pricing *models.ProjectPricing) error {
	if m.createProjectPricingFunc != nil {
		return m.createProjectPricingFunc(ctx, pricing)
	}
	return nil
}

func (m *mockPricingRepository) FindProjectPricingByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectPricing, error) {
	if m.findProjectPricingByProjectIDFunc != nil {
		return m.findProjectPricingByProjectIDFunc(ctx, projectID)
	}
	return []*models.ProjectPricing{}, nil
}

func (m *mockPricingRepository) CreateResourcePricing(ctx context.Context, pricing *models.ResourcePricing) error {
	if m.createResourcePricingFunc != nil {
		return m.createResourcePricingFunc(ctx, pricing)
	}
	return nil
}

func (m *mockPricingRepository) FindResourcePricingByResourceID(ctx context.Context, resourceID uuid.UUID) ([]*models.ResourcePricing, error) {
	if m.findResourcePricingByResourceIDFunc != nil {
		return m.findResourcePricingByResourceIDFunc(ctx, resourceID)
	}
	return []*models.ResourcePricing{}, nil
}

func (m *mockPricingRepository) FindResourcePricingByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.ResourcePricing, error) {
	if m.findResourcePricingByProjectIDFunc != nil {
		return m.findResourcePricingByProjectIDFunc(ctx, projectID)
	}
	return []*models.ResourcePricing{}, nil
}

func (m *mockPricingRepository) CreatePricingComponent(ctx context.Context, component *models.PricingComponent) error {
	if m.createPricingComponentFunc != nil {
		return m.createPricingComponentFunc(ctx, component)
	}
	return nil
}

func TestPricingService_CalculateResourceCost(t *testing.T) {
	ctx := context.Background()
	repo := &mockPricingRepository{}
	service := NewPricingService(repo)

	tests := []struct {
		name      string
		resource  *resource.Resource
		duration  time.Duration
		wantError bool
	}{
		{
			name:      "nil resource",
			resource:  nil,
			duration:  720 * time.Hour,
			wantError: true,
		},
		{
			name: "unsupported provider",
			resource: &resource.Resource{
				ID:       "res-1",
				Name:     "test-resource",
				Provider: "unsupported",
				Region:   "us-east-1",
				Type: resource.ResourceType{
					Name: "ec2_instance",
				},
			},
			duration:  720 * time.Hour,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimate, err := service.CalculateResourceCost(ctx, tt.resource, tt.duration)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if estimate == nil {
				t.Errorf("expected estimate but got nil")
			}
		})
	}
}

func TestPricingService_CalculateArchitectureCost(t *testing.T) {
	ctx := context.Background()
	repo := &mockPricingRepository{}
	service := NewPricingService(repo)

	tests := []struct {
		name      string
		arch      *architecture.Architecture
		duration  time.Duration
		wantError bool
	}{
		{
			name:      "nil architecture",
			arch:      nil,
			duration:  720 * time.Hour,
			wantError: true,
		},
		{
			name: "empty architecture",
			arch: &architecture.Architecture{
				Provider:  resource.AWS,
				Region:    "us-east-1",
				Resources: []*resource.Resource{},
			},
			duration:  720 * time.Hour,
			wantError: false,
		},
		{
			name: "architecture with visual-only resource (should be skipped)",
			arch: &architecture.Architecture{
				Provider: resource.AWS,
				Region:   "us-east-1",
				Resources: []*resource.Resource{
					{
						ID:       "vpc-1",
						Name:     "test-vpc",
						Provider: resource.AWS,
						Region:   "us-east-1",
						Type: resource.ResourceType{
							Name: "vpc",
						},
						Metadata: map[string]interface{}{
							"isVisualOnly": true,
						},
					},
				},
			},
			duration:  720 * time.Hour,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimate, err := service.CalculateArchitectureCost(ctx, tt.arch, tt.duration)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if estimate == nil {
				t.Errorf("expected estimate but got nil")
			}
		})
	}
}

func TestPricingService_GetProjectPricing(t *testing.T) {
	ctx := context.Background()
	projectID := uuid.New()

	tests := []struct {
		name      string
		setupMock func(*mockPricingRepository)
		wantError bool
		wantCount int
	}{
		{
			name: "returns pricing records",
			setupMock: func(repo *mockPricingRepository) {
				repo.findProjectPricingByProjectIDFunc = func(ctx context.Context, id uuid.UUID) ([]*models.ProjectPricing, error) {
					return []*models.ProjectPricing{
						{ID: 1, ProjectID: id, TotalCost: 100.0},
						{ID: 2, ProjectID: id, TotalCost: 150.0},
					}, nil
				}
			},
			wantError: false,
			wantCount: 2,
		},
		{
			name: "returns empty list",
			setupMock: func(repo *mockPricingRepository) {
				repo.findProjectPricingByProjectIDFunc = func(ctx context.Context, id uuid.UUID) ([]*models.ProjectPricing, error) {
					return []*models.ProjectPricing{}, nil
				}
			},
			wantError: false,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockPricingRepository{}
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}
			service := NewPricingService(repo)

			pricings, err := service.GetProjectPricing(ctx, projectID)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if len(pricings) != tt.wantCount {
				t.Errorf("expected %d pricing records, got %d", tt.wantCount, len(pricings))
			}
		})
	}
}
