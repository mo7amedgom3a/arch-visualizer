package containers

import (
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/containers"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

// ClusterCapacityProvidersFromResource converts a generic domain resource to an ECS Cluster Capacity Providers model
func ClusterCapacityProvidersFromResource(res *resource.Resource) (*containers.ECSClusterCapacityProviders, error) {
	if res.Type.Name != "ECSClusterCapacityProviders" {
		return nil, fmt.Errorf("invalid resource type for ECS Cluster Capacity Providers mapper: %s", res.Type.Name)
	}

	clusterProviders := &containers.ECSClusterCapacityProviders{
		ClusterName: res.Name,
	}

	getString := func(key string) string {
		if val, ok := res.Metadata[key]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
		return ""
	}

	getStringSlice := func(key string) []string {
		if val, ok := res.Metadata[key]; ok {
			if sl, ok := val.([]interface{}); ok {
				result := make([]string, 0, len(sl))
				for _, s := range sl {
					if str, ok := s.(string); ok {
						result = append(result, str)
					}
				}
				return result
			}
			if sl, ok := val.([]string); ok {
				return sl
			}
		}
		return nil
	}

	clusterProviders.ClusterName = getString("cluster_name")
	if clusterProviders.ClusterName == "" {
		clusterProviders.ClusterName = res.Name
	}

	clusterProviders.CapacityProviders = getStringSlice("capacity_providers")

	// Parse default capacity provider strategies
	if strategiesRaw, ok := res.Metadata["default_capacity_provider_strategy"]; ok {
		if strategies, ok := strategiesRaw.([]interface{}); ok {
			for _, stratRaw := range strategies {
				if stratMap, ok := stratRaw.(map[string]interface{}); ok {
					strat := containers.CapacityProviderStrategy{
						CapacityProvider: getMapString(stratMap, "capacity_provider"),
						Weight:           getMapInt(stratMap, "weight"),
						Base:             getMapInt(stratMap, "base"),
					}
					clusterProviders.DefaultCapacityProviderStrategies = append(
						clusterProviders.DefaultCapacityProviderStrategies,
						strat,
					)
				}
			}
		}
	}

	return clusterProviders, nil
}

// MapECSClusterCapacityProviders maps ECS Cluster Capacity Providers to a TerraformBlock
func MapECSClusterCapacityProviders(clusterProviders *containers.ECSClusterCapacityProviders) (*mapper.TerraformBlock, error) {
	if clusterProviders == nil {
		return nil, fmt.Errorf("ecs cluster capacity providers is nil")
	}

	attributes := make(map[string]mapper.TerraformValue)
	nestedBlocks := make(map[string][]mapper.NestedBlock)

	attributes["cluster_name"] = strVal(clusterProviders.ClusterName)

	if len(clusterProviders.CapacityProviders) > 0 {
		attributes["capacity_providers"] = listStrVal(clusterProviders.CapacityProviders)
	}

	// Default capacity provider strategy blocks
	if len(clusterProviders.DefaultCapacityProviderStrategies) > 0 {
		stratBlocks := make([]mapper.NestedBlock, 0, len(clusterProviders.DefaultCapacityProviderStrategies))
		for _, strat := range clusterProviders.DefaultCapacityProviderStrategies {
			stratAttrs := make(map[string]mapper.TerraformValue)
			stratAttrs["capacity_provider"] = strVal(strat.CapacityProvider)

			if strat.Weight > 0 {
				stratAttrs["weight"] = intVal(strat.Weight)
			}
			if strat.Base > 0 {
				stratAttrs["base"] = intVal(strat.Base)
			}

			stratBlocks = append(stratBlocks, mapper.NestedBlock{Attributes: stratAttrs})
		}
		nestedBlocks["default_capacity_provider_strategy"] = stratBlocks
	}

	return &mapper.TerraformBlock{
		Kind:         "resource",
		Labels:       []string{"aws_ecs_cluster_capacity_providers", clusterProviders.ClusterName},
		Attributes:   attributes,
		NestedBlocks: nestedBlocks,
	}, nil
}
