package containers

import (
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/inventory"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	tfmapper "github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

func init() {
	inv := inventory.GetDefaultInventory()

	// Register Terraform Mapper for ECS Cluster
	inv.SetTerraformMapper("ECSCluster", func(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
		cluster, err := ClusterFromResource(res)
		if err != nil {
			return nil, err
		}
		block, err := MapECSCluster(cluster)
		if err != nil {
			return nil, err
		}
		return []tfmapper.TerraformBlock{*block}, nil
	})

	// Register Terraform Mapper for ECS Task Definition
	inv.SetTerraformMapper("ECSTaskDefinition", func(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
		taskDef, err := TaskDefinitionFromResource(res)
		if err != nil {
			return nil, err
		}
		block, err := MapECSTaskDefinition(taskDef)
		if err != nil {
			return nil, err
		}
		return []tfmapper.TerraformBlock{*block}, nil
	})

	// Register Terraform Mapper for ECS Service
	inv.SetTerraformMapper("ECSService", func(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
		service, err := ServiceFromResource(res)
		if err != nil {
			return nil, err
		}
		block, err := MapECSService(service)
		if err != nil {
			return nil, err
		}
		return []tfmapper.TerraformBlock{*block}, nil
	})

	// Register Terraform Mapper for ECS Capacity Provider
	inv.SetTerraformMapper("ECSCapacityProvider", func(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
		provider, err := CapacityProviderFromResource(res)
		if err != nil {
			return nil, err
		}
		block, err := MapECSCapacityProvider(provider)
		if err != nil {
			return nil, err
		}
		return []tfmapper.TerraformBlock{*block}, nil
	})

	// Register Terraform Mapper for ECS Cluster Capacity Providers
	inv.SetTerraformMapper("ECSClusterCapacityProviders", func(res *resource.Resource) ([]tfmapper.TerraformBlock, error) {
		clusterProviders, err := ClusterCapacityProvidersFromResource(res)
		if err != nil {
			return nil, err
		}
		block, err := MapECSClusterCapacityProviders(clusterProviders)
		if err != nil {
			return nil, err
		}
		return []tfmapper.TerraformBlock{*block}, nil
	})
}
