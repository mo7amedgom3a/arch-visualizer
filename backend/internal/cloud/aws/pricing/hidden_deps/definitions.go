package hiddendeps

import (
	"fmt"

	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
)

// GetHiddenDependenciesForResourceType returns hardcoded hidden dependencies for a resource type
func GetHiddenDependenciesForResourceType(resourceType string) []*domainpricing.HiddenDependency {
	switch resourceType {
	case "nat_gateway", "NATGateway":
		return []*domainpricing.HiddenDependency{
			{
				ParentResourceType:  resourceType,
				ChildResourceType:   "elastic_ip",
				QuantityExpression:  "1",
				ConditionExpression: "metadata.allocationId == null",
				IsAttached:          true,
				Description:         "NAT Gateway requires an Elastic IP. If not provided, one is automatically created and attached (free when attached).",
			},
		}
	case "ec2_instance", "EC2":
		return []*domainpricing.HiddenDependency{
			{
				ParentResourceType:  resourceType,
				ChildResourceType:   "ebs_volume",
				QuantityExpression:  "metadata.size_gb",
				ConditionExpression: "",
				IsAttached:          true,
				Description:         "EC2 instance requires a root EBS volume. Default size is 8GB if not specified.",
			},
			{
				ParentResourceType:  resourceType,
				ChildResourceType:   "network_interface",
				QuantityExpression:  "1",
				ConditionExpression: "",
				IsAttached:          true,
				Description:         "EC2 instance requires a network interface (free when attached).",
			},
		}
	case "load_balancer", "LoadBalancer":
		// Load balancers don't have direct hidden costs, but may have target groups
		// Target groups themselves don't have direct costs
		return []*domainpricing.HiddenDependency{}
	case "rds_instance", "RDS":
		return []*domainpricing.HiddenDependency{
			{
				ParentResourceType:  resourceType,
				ChildResourceType:   "ebs_volume",
				QuantityExpression:  "metadata.allocated_storage",
				ConditionExpression: "",
				IsAttached:          true,
				Description:         "RDS instance requires storage volume based on allocated_storage.",
			},
			{
				ParentResourceType:  resourceType,
				ChildResourceType:   "s3_bucket",
				QuantityExpression:  "metadata.allocated_storage",
				ConditionExpression: "metadata.backup_retention_period > 0",
				IsAttached:          false,
				Description:         "RDS automated backups stored in S3 (assumed equal to DB size for estimation).",
			},
		}
	case "lambda_function", "Lambda":
		// Lambda has CloudWatch Logs costs, but these are usage-based and hard to estimate
		// We'll skip for now or add as optional
		return []*domainpricing.HiddenDependency{}
	case "auto_scaling_group", "AutoScalingGroup":
		// ASG costs are derived from EC2 instances, not hidden dependencies
		return []*domainpricing.HiddenDependency{}
	default:
		return []*domainpricing.HiddenDependency{}
	}
}

// ResolveForResource resolves hidden dependencies for a specific resource
func ResolveForResource(res *resource.Resource, architecture interface{}) ([]*domainpricing.HiddenDependencyResource, error) {
	deps := GetHiddenDependenciesForResourceType(res.Type.Name)
	if len(deps) == 0 {
		return []*domainpricing.HiddenDependencyResource{}, nil
	}

	var resolved []*domainpricing.HiddenDependencyResource
	for _, dep := range deps {
		// Check condition
		if dep.ConditionExpression != "" {
			if !evaluateCondition(dep.ConditionExpression, res) {
				continue
			}
		}

		// Calculate quantity
		quantity := calculateQuantity(dep.QuantityExpression, res)

		// Create virtual resource
		hiddenRes := createHiddenResource(dep, res, quantity)

		resolved = append(resolved, &domainpricing.HiddenDependencyResource{
			Dependency: dep,
			Resource:   hiddenRes,
			Quantity:   quantity,
		})
	}

	return resolved, nil
}

// evaluateCondition evaluates a condition expression
func evaluateCondition(condition string, res *resource.Resource) bool {
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

	return false
}

// calculateQuantity calculates quantity from an expression
func calculateQuantity(expression string, res *resource.Resource) float64 {
	if expression == "" || expression == "1" {
		return 1.0
	}

	// Support metadata.field access
	if expression == "metadata.size_gb" {
		if res.Metadata != nil {
			if size, ok := res.Metadata["size_gb"].(float64); ok {
				return size
			}
			if size, ok := res.Metadata["size_gb"].(int); ok {
				return float64(size)
			}
		}
		return 8.0 // Default EBS root volume size
	}

	if expression == "metadata.allocated_storage" {
		if res.Metadata != nil {
			if size, ok := res.Metadata["allocated_storage"].(float64); ok {
				return size
			}
			if size, ok := res.Metadata["allocated_storage"].(int); ok {
				return float64(size)
			}
		}
		return 20.0 // Default RDS storage
	}

	return 1.0
}

// createHiddenResource creates a virtual resource for a hidden dependency
func createHiddenResource(dep *domainpricing.HiddenDependency, parent *resource.Resource, quantity float64) *resource.Resource {
	childType := mapToDomainResourceType(dep.ChildResourceType)

	metadata := make(map[string]interface{})
	if dep.ChildResourceType == "ebs_volume" {
		metadata["size_gb"] = quantity
		metadata["volume_type"] = "gp3"
		metadata["is_root"] = true
	} else if dep.ChildResourceType == "elastic_ip" {
		metadata["domain"] = "vpc"
		metadata["is_attached"] = dep.IsAttached
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
func mapToDomainResourceType(pricingType string) resource.ResourceType {
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

	return resource.ResourceType{
		ID:   pricingType,
		Name: pricingType,
	}
}
