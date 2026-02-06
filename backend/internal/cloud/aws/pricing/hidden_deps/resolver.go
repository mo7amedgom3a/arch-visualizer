package hiddendeps

import (
	"context"
	"fmt"
	"strconv"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/repository"
)

// AWSHiddenDependencyResolver implements HiddenDependencyResolver for AWS
type AWSHiddenDependencyResolver struct {
	hiddenDepRepo *repository.HiddenDependencyRepository
}

// NewAWSHiddenDependencyResolver creates a new AWS hidden dependency resolver
func NewAWSHiddenDependencyResolver(hiddenDepRepo *repository.HiddenDependencyRepository) *AWSHiddenDependencyResolver {
	return &AWSHiddenDependencyResolver{
		hiddenDepRepo: hiddenDepRepo,
	}
}

// ResolveHiddenDependencies resolves all hidden dependencies for a resource
func (r *AWSHiddenDependencyResolver) ResolveHiddenDependencies(ctx context.Context, res *resource.Resource, architecture interface{}) ([]*domainpricing.HiddenDependencyResource, error) {
	// Get hidden dependencies from database
	dbDeps, err := r.hiddenDepRepo.FindByParentResourceType(ctx, "aws", res.Type.Name)
	if err != nil {
		// If not found in DB, fall back to hardcoded definitions
		return r.resolveFromDefinitions(ctx, res, architecture)
	}

	// Convert database models to domain models and resolve
	var resolved []*domainpricing.HiddenDependencyResource
	for _, dbDep := range dbDeps {
		dep := &domainpricing.HiddenDependency{
			ParentResourceType:  dbDep.ParentResourceType,
			ChildResourceType:   dbDep.ChildResourceType,
			QuantityExpression:  dbDep.QuantityExpression,
			ConditionExpression: dbDep.ConditionExpression,
			IsAttached:          dbDep.IsAttached,
			Description:         dbDep.Description,
		}

		// Check condition if specified
		if dep.ConditionExpression != "" {
			if !r.evaluateCondition(dep.ConditionExpression, res) {
				continue
			}
		}

		// Calculate quantity
		quantity, err := r.calculateQuantity(dep.QuantityExpression, res)
		if err != nil {
			continue // Skip if quantity calculation fails
		}

		// Create virtual resource for the hidden dependency
		hiddenRes := r.createHiddenResource(dep, res, quantity)

		resolved = append(resolved, &domainpricing.HiddenDependencyResource{
			Dependency: dep,
			Resource:   hiddenRes,
			Quantity:   quantity,
		})
	}

	// Also check hardcoded definitions for any additional dependencies
	hardcoded, err := r.resolveFromDefinitions(ctx, res, architecture)
	if err == nil && len(hardcoded) > 0 {
		// Merge hardcoded dependencies (avoid duplicates)
		for _, hc := range hardcoded {
			exists := false
			for _, r := range resolved {
				if r.Dependency.ChildResourceType == hc.Dependency.ChildResourceType {
					exists = true
					break
				}
			}
			if !exists {
				resolved = append(resolved, hc)
			}
		}
	}

	return resolved, nil
}

// GetHiddenDependenciesForResourceType returns all hidden dependencies for a resource type
func (r *AWSHiddenDependencyResolver) GetHiddenDependenciesForResourceType(ctx context.Context, provider domainpricing.CloudProvider, resourceType string) ([]*domainpricing.HiddenDependency, error) {
	if provider != domainpricing.AWS {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	dbDeps, err := r.hiddenDepRepo.FindByParentResourceType(ctx, "aws", resourceType)
	if err != nil {
		// Fall back to hardcoded definitions
		return GetHiddenDependenciesForResourceType(resourceType), nil
	}

	var deps []*domainpricing.HiddenDependency
	for _, dbDep := range dbDeps {
		deps = append(deps, &domainpricing.HiddenDependency{
			ParentResourceType:  dbDep.ParentResourceType,
			ChildResourceType:   dbDep.ChildResourceType,
			QuantityExpression:  dbDep.QuantityExpression,
			ConditionExpression: dbDep.ConditionExpression,
			IsAttached:          dbDep.IsAttached,
			Description:         dbDep.Description,
		})
	}

	// Merge with hardcoded definitions
	hardcoded := GetHiddenDependenciesForResourceType(resourceType)
	if len(hardcoded) > 0 {
		for _, hc := range hardcoded {
			exists := false
			for _, d := range deps {
				if d.ChildResourceType == hc.ChildResourceType {
					exists = true
					break
				}
			}
			if !exists {
				deps = append(deps, hc)
			}
		}
	}

	return deps, nil
}

// resolveFromDefinitions resolves dependencies from hardcoded definitions
func (r *AWSHiddenDependencyResolver) resolveFromDefinitions(ctx context.Context, res *resource.Resource, architecture interface{}) ([]*domainpricing.HiddenDependencyResource, error) {
	return ResolveForResource(res, architecture)
}

// evaluateCondition evaluates a condition expression
func (r *AWSHiddenDependencyResolver) evaluateCondition(condition string, res *resource.Resource) bool {
	// Simple condition evaluation - can be enhanced with expression parser
	// For now, support simple checks like "metadata.allocationId == null"
	if condition == "" {
		return true
	}

	// Check if allocationId is missing (NAT Gateway case)
	if condition == "metadata.allocationId == null" || condition == "!metadata.allocationId" {
		if res.Metadata != nil {
			if _, ok := res.Metadata["allocationId"]; !ok {
				return true
			}
			if allocId, ok := res.Metadata["allocationId"].(string); ok && allocId == "" {
				return true
			}
		} else {
			return true
		}
	}

	// Check if backup_retention_period > 0 (RDS case)
	if condition == "metadata.backup_retention_period > 0" {
		if res.Metadata != nil {
			if br, ok := res.Metadata["backup_retention_period"].(float64); ok && br > 0 {
				return true
			}
			if br, ok := res.Metadata["backup_retention_period"].(int); ok && br > 0 {
				return true
			}
		}
	}

	return false
}

// calculateQuantity calculates quantity from an expression
func (r *AWSHiddenDependencyResolver) calculateQuantity(expression string, res *resource.Resource) (float64, error) {
	if expression == "" || expression == "1" {
		return 1.0, nil
	}

	// Simple expression evaluation - can be enhanced
	// Support metadata.field access
	if expression == "metadata.size_gb" {
		if res.Metadata != nil {
			if size, ok := res.Metadata["size_gb"].(float64); ok {
				return size, nil
			}
			if size, ok := res.Metadata["size_gb"].(int); ok {
				return float64(size), nil
			}
		}
		return 8.0, nil // Default EBS root volume size
	}

	// Try to parse as number
	if val, err := strconv.ParseFloat(expression, 64); err == nil {
		return val, nil
	}

	return 1.0, nil // Default
}

// createHiddenResource creates a virtual resource for a hidden dependency
func (r *AWSHiddenDependencyResolver) createHiddenResource(dep *domainpricing.HiddenDependency, parent *resource.Resource, quantity float64) *resource.Resource {
	// Map child resource type to domain resource type
	childType := r.mapToDomainResourceType(dep.ChildResourceType)

	// Create metadata based on dependency type
	metadata := make(map[string]interface{})
	if dep.ChildResourceType == "ebs_volume" {
		metadata["size_gb"] = quantity
		metadata["volume_type"] = "gp3"
		metadata["is_root"] = true
	} else if dep.ChildResourceType == "elastic_ip" {
		metadata["domain"] = "vpc"
		metadata["is_attached"] = dep.IsAttached
	} else if dep.ChildResourceType == "s3_bucket" {
		metadata["size_gb"] = quantity
		metadata["storage_class"] = "standard"
	}

	return &resource.Resource{
		ID:       fmt.Sprintf("%s-hidden-%s", parent.ID, dep.ChildResourceType),
		Name:     fmt.Sprintf("%s-%s", parent.Name, dep.ChildResourceType),
		Type:     childType,
		Provider: parent.Provider,
		Region:   parent.Region,
		ParentID: &parent.ID,
		Metadata: metadata,
	}
}

// mapToDomainResourceType maps pricing resource type to domain resource type
func (r *AWSHiddenDependencyResolver) mapToDomainResourceType(pricingType string) resource.ResourceType {
	mapping := map[string]resource.ResourceType{
		"elastic_ip": {
			ID:   "ElasticIP",
			Name: "ElasticIP",
		},
		"ebs_volume": {
			ID:   "EBS",
			Name: "EBS",
		},
		"network_interface": {
			ID:   "NetworkInterface",
			Name: "NetworkInterface",
		},
	}

	if rt, ok := mapping[pricingType]; ok {
		return rt
	}

	// Default
	return resource.ResourceType{
		ID:   pricingType,
		Name: pricingType,
	}
}
