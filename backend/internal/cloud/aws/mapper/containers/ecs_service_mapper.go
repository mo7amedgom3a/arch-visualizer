package containers

import (
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/containers"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

// ServiceFromResource converts a generic domain resource to an ECS Service model
func ServiceFromResource(res *resource.Resource) (*containers.ECSService, error) {
	if res.Type.Name != "ECSService" {
		return nil, fmt.Errorf("invalid resource type for ECS Service mapper: %s", res.Type.Name)
	}

	service := &containers.ECSService{
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

	service.ClusterID = getString("cluster_id")
	service.TaskDefinitionARN = getString("task_definition_arn")
	service.DesiredCount = getInt("desired_count")
	service.LaunchType = getString("launch_type")
	service.PlatformVersion = getString("platform_version")
	service.SchedulingStrategy = getString("scheduling_strategy")
	service.EnableExecuteCommand = getBool("enable_execute_command")
	service.ForceNewDeployment = getBool("force_new_deployment")

	// Network configuration
	subnets := getStringSlice("subnets")
	securityGroups := getStringSlice("security_groups")
	assignPublicIP := getBool("assign_public_ip")

	if len(subnets) > 0 || len(securityGroups) > 0 {
		service.NetworkConfiguration = &containers.NetworkConfiguration{
			Subnets:        subnets,
			SecurityGroups: securityGroups,
			AssignPublicIP: assignPublicIP,
		}
	}

	// Load balancer configuration
	targetGroupARN := getString("target_group_arn")
	containerName := getString("container_name")
	containerPort := getInt("container_port")

	if targetGroupARN != "" && containerName != "" {
		service.LoadBalancers = []containers.LoadBalancerConfig{
			{
				TargetGroupARN: targetGroupARN,
				ContainerName:  containerName,
				ContainerPort:  containerPort,
			},
		}
	}

	// Deployment circuit breaker
	if enableCircuitBreaker := getBool("enable_circuit_breaker"); enableCircuitBreaker {
		service.DeploymentConfiguration = &containers.DeploymentConfiguration{
			CircuitBreaker: &containers.CircuitBreaker{
				Enable:   true,
				Rollback: getBool("circuit_breaker_rollback"),
			},
		}
	}

	return service, nil
}

// MapECSService maps an ECS Service to a TerraformBlock
func MapECSService(service *containers.ECSService) (*mapper.TerraformBlock, error) {
	if service == nil {
		return nil, fmt.Errorf("ecs service is nil")
	}

	attributes := make(map[string]mapper.TerraformValue)
	nestedBlocks := make(map[string][]mapper.NestedBlock)

	attributes["name"] = strVal(service.Name)

	if service.ClusterID != "" {
		attributes["cluster"] = strVal(service.ClusterID)
	}

	if service.TaskDefinitionARN != "" {
		attributes["task_definition"] = strVal(service.TaskDefinitionARN)
	}

	if service.DesiredCount > 0 {
		attributes["desired_count"] = intVal(service.DesiredCount)
	}

	if service.LaunchType != "" {
		attributes["launch_type"] = strVal(service.LaunchType)
	}

	if service.PlatformVersion != "" {
		attributes["platform_version"] = strVal(service.PlatformVersion)
	}

	if service.SchedulingStrategy != "" {
		attributes["scheduling_strategy"] = strVal(service.SchedulingStrategy)
	}

	if service.EnableExecuteCommand {
		attributes["enable_execute_command"] = boolVal(true)
	}

	if service.ForceNewDeployment {
		attributes["force_new_deployment"] = boolVal(true)
	}

	// Network configuration block
	if service.NetworkConfiguration != nil {
		netConfigAttrs := make(map[string]mapper.TerraformValue)
		if len(service.NetworkConfiguration.Subnets) > 0 {
			netConfigAttrs["subnets"] = listStrVal(service.NetworkConfiguration.Subnets)
		}
		if len(service.NetworkConfiguration.SecurityGroups) > 0 {
			netConfigAttrs["security_groups"] = listStrVal(service.NetworkConfiguration.SecurityGroups)
		}
		if service.NetworkConfiguration.AssignPublicIP {
			netConfigAttrs["assign_public_ip"] = boolVal(true)
		}

		nestedBlocks["network_configuration"] = []mapper.NestedBlock{
			{Attributes: netConfigAttrs},
		}
	}

	// Load balancer blocks
	if len(service.LoadBalancers) > 0 {
		lbBlocks := make([]mapper.NestedBlock, 0, len(service.LoadBalancers))
		for _, lb := range service.LoadBalancers {
			lbAttrs := map[string]mapper.TerraformValue{
				"target_group_arn": strVal(lb.TargetGroupARN),
				"container_name":   strVal(lb.ContainerName),
				"container_port":   intVal(lb.ContainerPort),
			}
			lbBlocks = append(lbBlocks, mapper.NestedBlock{Attributes: lbAttrs})
		}
		nestedBlocks["load_balancer"] = lbBlocks
	}

	// Deployment circuit breaker
	if service.DeploymentConfiguration != nil && service.DeploymentConfiguration.CircuitBreaker != nil {
		cbAttrs := map[string]mapper.TerraformValue{
			"enable":   boolVal(service.DeploymentConfiguration.CircuitBreaker.Enable),
			"rollback": boolVal(service.DeploymentConfiguration.CircuitBreaker.Rollback),
		}
		nestedBlocks["deployment_circuit_breaker"] = []mapper.NestedBlock{
			{Attributes: cbAttrs},
		}
	}

	return &mapper.TerraformBlock{
		Kind:         "resource",
		Labels:       []string{"aws_ecs_service", service.Name},
		Attributes:   attributes,
		NestedBlocks: nestedBlocks,
	}, nil
}
