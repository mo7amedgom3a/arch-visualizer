package containers

import (
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/containers"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

// CapacityProviderFromResource converts a generic domain resource to an ECS Capacity Provider model
func CapacityProviderFromResource(res *resource.Resource) (*containers.ECSCapacityProvider, error) {
	if res.Type.Name != "ECSCapacityProvider" {
		return nil, fmt.Errorf("invalid resource type for ECS Capacity Provider mapper: %s", res.Type.Name)
	}

	provider := &containers.ECSCapacityProvider{
		Name: res.Name,
	}

	getString := func(key string) string {
		if val, ok := res.Metadata[key]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
		return ""
	}

	getInt := func(key string) int {
		if val, ok := res.Metadata[key]; ok {
			if i, ok := val.(int); ok {
				return i
			}
			if f, ok := val.(float64); ok {
				return int(f)
			}
		}
		return 0
	}

	getBool := func(key string) bool {
		if val, ok := res.Metadata[key]; ok {
			if b, ok := val.(bool); ok {
				return b
			}
		}
		return false
	}

	// Auto scaling group provider configuration
	asgARN := getString("auto_scaling_group_arn")
	if asgARN != "" {
		provider.AutoScalingGroupProvider = &containers.AutoScalingGroupProvider{
			AutoScalingGroupARN:          asgARN,
			ManagedTerminationProtection: getString("managed_termination_protection"),
			ManagedDraining:              getString("managed_draining"),
		}

		// Managed scaling configuration
		if getBool("managed_scaling_enabled") {
			provider.AutoScalingGroupProvider.ManagedScaling = &containers.ManagedScaling{
				Status:                 "ENABLED",
				TargetCapacity:         getInt("target_capacity"),
				MinimumScalingStepSize: getInt("minimum_scaling_step_size"),
				MaximumScalingStepSize: getInt("maximum_scaling_step_size"),
				InstanceWarmupPeriod:   getInt("instance_warmup_period"),
			}
		}
	}

	return provider, nil
}

// MapECSCapacityProvider maps an ECS Capacity Provider to a TerraformBlock
func MapECSCapacityProvider(provider *containers.ECSCapacityProvider) (*mapper.TerraformBlock, error) {
	if provider == nil {
		return nil, fmt.Errorf("ecs capacity provider is nil")
	}

	attributes := make(map[string]mapper.TerraformValue)
	nestedBlocks := make(map[string][]mapper.NestedBlock)

	attributes["name"] = strVal(provider.Name)

	// Auto scaling group provider block
	if provider.AutoScalingGroupProvider != nil {
		asgProviderAttrs := make(map[string]mapper.TerraformValue)
		asgProviderAttrs["auto_scaling_group_arn"] = strVal(provider.AutoScalingGroupProvider.AutoScalingGroupARN)

		if provider.AutoScalingGroupProvider.ManagedTerminationProtection != "" {
			asgProviderAttrs["managed_termination_protection"] = strVal(provider.AutoScalingGroupProvider.ManagedTerminationProtection)
		}

		if provider.AutoScalingGroupProvider.ManagedDraining != "" {
			asgProviderAttrs["managed_draining"] = strVal(provider.AutoScalingGroupProvider.ManagedDraining)
		}

		nestedBlocks["auto_scaling_group_provider"] = []mapper.NestedBlock{
			{Attributes: asgProviderAttrs},
		}

		// Managed scaling block - nested within auto_scaling_group_provider
		// Note: In Terraform, managed_scaling is nested inside auto_scaling_group_provider
		// We'll add it as a separate nested block for simplicity; the HCL writer should handle nesting
		if provider.AutoScalingGroupProvider.ManagedScaling != nil {
			ms := provider.AutoScalingGroupProvider.ManagedScaling
			msAttrs := make(map[string]mapper.TerraformValue)

			if ms.Status != "" {
				msAttrs["status"] = strVal(ms.Status)
			}
			if ms.TargetCapacity > 0 {
				msAttrs["target_capacity"] = intVal(ms.TargetCapacity)
			}
			if ms.MinimumScalingStepSize > 0 {
				msAttrs["minimum_scaling_step_size"] = intVal(ms.MinimumScalingStepSize)
			}
			if ms.MaximumScalingStepSize > 0 {
				msAttrs["maximum_scaling_step_size"] = intVal(ms.MaximumScalingStepSize)
			}
			if ms.InstanceWarmupPeriod > 0 {
				msAttrs["instance_warmup_period"] = intVal(ms.InstanceWarmupPeriod)
			}

			nestedBlocks["managed_scaling"] = []mapper.NestedBlock{
				{Attributes: msAttrs},
			}
		}
	}

	return &mapper.TerraformBlock{
		Kind:         "resource",
		Labels:       []string{"aws_ecs_capacity_provider", provider.Name},
		Attributes:   attributes,
		NestedBlocks: nestedBlocks,
	}, nil
}
