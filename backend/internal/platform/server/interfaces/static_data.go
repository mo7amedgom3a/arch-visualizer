package interfaces

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

// StaticDataService handles static reference data
type StaticDataService interface {
	// ListResourceTypes retrieves all resource types
	ListResourceTypes(ctx context.Context) ([]*models.ResourceType, error)
	// ListResourceTypesByProvider retrieves resource types for a provider
	ListResourceTypesByProvider(ctx context.Context, provider string) ([]*models.ResourceType, error)
	// ListProviders retrieves supported cloud providers
	ListProviders(ctx context.Context) ([]string, error)
}
