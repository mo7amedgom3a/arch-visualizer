package services

import (
	"context"

	platformerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/errors"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
	resourcerepo "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository/resource"
	serverinterfaces "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/server/interfaces"
)

// ConstraintService handles logic for resource constraints
type ConstraintService struct {
	repo *resourcerepo.ResourceConstraintRepository
}

// NewConstraintService creates a new constraint service
func NewConstraintService(repo *resourcerepo.ResourceConstraintRepository) *ConstraintService {
	return &ConstraintService{
		repo: repo,
	}
}

// GetConstraintsForResource retrieves all constraints for a given resource type
func (s *ConstraintService) GetConstraintsForResource(ctx context.Context, resourceTypeID uint) ([]*models.ResourceConstraint, error) {
	constraints, err := s.repo.FindByResourceType(ctx, resourceTypeID)
	if err != nil {
		return nil, platformerrors.NewDatabaseQueryFailed("get_constraints_for_resource", err)
	}
	return constraints, nil
}

// GetAllConstraints retrieves all constraints
func (s *ConstraintService) GetAllConstraints(ctx context.Context) ([]serverinterfaces.ConstraintRecord, error) {
	// We need a method in repo to get all. For now let's assume we can list all or hack it.
	// Since repo doesn't have FindAll, let's look at BaseRepository or add it.
	// Actually better to iterate over known resource types or add GetAll to repo.
	// For this task, I should probably add FindAll to ResourceConstraintRepository.

	// Let's rely on repo.GetAll() if it exists?
	// The BaseRepository usually has generic Find.

	// Use the repository method to find all constraints
	constraints, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, platformerrors.NewDatabaseQueryFailed("get_all_constraints", err)
	}

	var records []serverinterfaces.ConstraintRecord
	for _, c := range constraints {
		// ResourceType is already preloaded by the repository
		if c.ResourceType.Name == "" {
			// Skip valid check for now, but in theory preloading ensures it's there
		}

		records = append(records, serverinterfaces.ConstraintRecord{
			ResourceType:    c.ResourceType.Name,
			ConstraintType:  c.ConstraintType,
			ConstraintValue: c.ConstraintValue,
		})
	}

	return records, nil
}

// SaveConstraint saves a new constraint
func (s *ConstraintService) SaveConstraint(ctx context.Context, constraint *models.ResourceConstraint) error {
	// Here you might add validation logic before saving
	return s.repo.Create(ctx, constraint)
}
