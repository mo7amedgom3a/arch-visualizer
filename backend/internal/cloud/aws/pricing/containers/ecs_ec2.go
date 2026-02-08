package containers

import (
	domainpricing "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/pricing"
)

// ECS on EC2 does not have additional ECS charges.
// The cost is based on the underlying EC2 instances.
// This file provides helper functions for ECS on EC2 pricing context.

// GetECSEC2Pricing returns pricing information for ECS on EC2
// ECS on EC2 has no additional charges - cost is based on EC2 instances
func GetECSEC2Pricing(region string) *domainpricing.ResourcePricing {
	return &domainpricing.ResourcePricing{
		ResourceType: "ecs_ec2",
		Provider:     domainpricing.AWS,
		Components:   []domainpricing.PriceComponent{},
		Metadata: map[string]interface{}{
			"launch_type": "EC2",
			"note":        "ECS on EC2 has no additional ECS charges. Cost is based on underlying EC2 instances.",
		},
	}
}

// GetECSClusterPricing returns pricing information for an ECS Cluster
// ECS Clusters themselves don't have a direct charge
func GetECSClusterPricing(clusterName, region string) *domainpricing.ResourcePricing {
	return &domainpricing.ResourcePricing{
		ResourceType: "ecs_cluster",
		Provider:     domainpricing.AWS,
		Components:   []domainpricing.PriceComponent{},
		Metadata: map[string]interface{}{
			"cluster_name": clusterName,
			"note":         "ECS Clusters have no direct charge. Cost is based on tasks/services running in the cluster.",
		},
	}
}

// GetECSServicePricing returns pricing for an ECS Service
// Service pricing depends on launch type (Fargate vs EC2)
func GetECSServicePricing(serviceName, launchType, taskCPU, taskMemory, region string, desiredCount int) *domainpricing.ResourcePricing {
	if launchType == "FARGATE" {
		vcpu := parseCPUToVCPU(taskCPU)
		memoryGB := parseMemoryToGB(taskMemory)

		// Get per-task pricing and scale by desired count
		taskPricing := GetFargatePricing(vcpu*float64(desiredCount), memoryGB*float64(desiredCount), region, false)

		taskPricing.ResourceType = "ecs_service"
		taskPricing.Metadata["service_name"] = serviceName
		taskPricing.Metadata["desired_count"] = desiredCount
		taskPricing.Metadata["launch_type"] = launchType

		return taskPricing
	}

	// EC2 launch type
	return &domainpricing.ResourcePricing{
		ResourceType: "ecs_service",
		Provider:     domainpricing.AWS,
		Components:   []domainpricing.PriceComponent{},
		Metadata: map[string]interface{}{
			"service_name":  serviceName,
			"launch_type":   launchType,
			"desired_count": desiredCount,
			"note":          "EC2 launch type: no additional ECS charges, uses EC2 instance pricing",
		},
	}
}

// GetECSCapacityProviderPricing returns pricing info for capacity providers
func GetECSCapacityProviderPricing(providerName string) *domainpricing.ResourcePricing {
	return &domainpricing.ResourcePricing{
		ResourceType: "ecs_capacity_provider",
		Provider:     domainpricing.AWS,
		Components:   []domainpricing.PriceComponent{},
		Metadata: map[string]interface{}{
			"provider_name": providerName,
			"note":          "Capacity providers have no direct charge. Cost is based on underlying ASG/EC2 instances.",
		},
	}
}
