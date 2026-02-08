package containers

import (
	"encoding/json"
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/containers"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

// TaskDefinitionFromResource converts a generic domain resource to an ECS Task Definition model
func TaskDefinitionFromResource(res *resource.Resource) (*containers.ECSTaskDefinition, error) {
	if res.Type.Name != "ECSTaskDefinition" {
		return nil, fmt.Errorf("invalid resource type for ECS Task Definition mapper: %s", res.Type.Name)
	}

	taskDef := &containers.ECSTaskDefinition{
		Family: res.Name,
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

	// Check both snake_case and camelCase for requires_compatibilities
	taskDef.RequiresCompatibilities = getStringSlice("requires_compatibilities")
	if len(taskDef.RequiresCompatibilities) == 0 {
		taskDef.RequiresCompatibilities = getStringSlice("requiresCompatibilities")
	}

	// Check both snake_case and camelCase for network_mode
	taskDef.NetworkMode = getString("network_mode")
	if taskDef.NetworkMode == "" {
		taskDef.NetworkMode = getString("networkMode")
	}

	taskDef.CPU = getString("cpu")
	taskDef.Memory = getString("memory")

	// Check both naming conventions for role ARNs
	taskDef.ExecutionRoleARN = getString("execution_role_arn")
	if taskDef.ExecutionRoleARN == "" {
		taskDef.ExecutionRoleARN = getString("executionRoleArn")
	}

	taskDef.TaskRoleARN = getString("task_role_arn")
	if taskDef.TaskRoleARN == "" {
		taskDef.TaskRoleARN = getString("taskRoleArn")
	}

	// Parse container definitions from metadata - check both naming conventions
	var containerDefsRaw interface{}
	if raw, ok := res.Metadata["container_definitions"]; ok {
		containerDefsRaw = raw
	} else if raw, ok := res.Metadata["containerDefinitions"]; ok {
		containerDefsRaw = raw
	}

	if containerDefsRaw != nil {
		if containerDefs, err := parseContainerDefinitions(containerDefsRaw); err == nil && len(containerDefs) > 0 {
			taskDef.ContainerDefinitions = containerDefs
		}
	}

	// If no container definitions provided, create a default one with mock image
	if len(taskDef.ContainerDefinitions) == 0 {
		defaultContainer := containers.ContainerDefinition{
			Name:      res.Name,
			Image:     "nginx:latest", // Mock image since ECR is not implemented yet
			Essential: true,
			CPU:       256,
			Memory:    512,
			PortMappings: []containers.PortMapping{
				{ContainerPort: 80, HostPort: 80, Protocol: "tcp"},
			},
		}
		taskDef.ContainerDefinitions = []containers.ContainerDefinition{defaultContainer}
	}

	return taskDef, nil
}

// parseContainerDefinitions parses container definitions from metadata
func parseContainerDefinitions(raw interface{}) ([]containers.ContainerDefinition, error) {
	var defs []containers.ContainerDefinition

	switch v := raw.(type) {
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				def := containers.ContainerDefinition{
					Name:      getMapString(m, "name"),
					Image:     getMapString(m, "image"),
					Essential: getMapBool(m, "essential"),
					CPU:       getMapInt(m, "cpu"),
					Memory:    getMapInt(m, "memory"),
				}

				// Parse port mappings
				if pmRaw, ok := m["port_mappings"]; ok {
					if pms, ok := pmRaw.([]interface{}); ok {
						for _, pm := range pms {
							if pmMap, ok := pm.(map[string]interface{}); ok {
								def.PortMappings = append(def.PortMappings, containers.PortMapping{
									ContainerPort: getMapInt(pmMap, "container_port"),
									HostPort:      getMapInt(pmMap, "host_port"),
									Protocol:      getMapString(pmMap, "protocol"),
								})
							}
						}
					}
				}

				// Parse environment variables
				if envRaw, ok := m["environment"]; ok {
					if envs, ok := envRaw.([]interface{}); ok {
						for _, env := range envs {
							if envMap, ok := env.(map[string]interface{}); ok {
								def.Environment = append(def.Environment, containers.KeyValuePair{
									Name:  getMapString(envMap, "name"),
									Value: getMapString(envMap, "value"),
								})
							}
						}
					}
				}

				// Parse log configuration
				if logRaw, ok := m["log_configuration"]; ok {
					if logMap, ok := logRaw.(map[string]interface{}); ok {
						logConfig := &containers.LogConfiguration{
							LogDriver: getMapString(logMap, "log_driver"),
						}
						if opts, ok := logMap["options"].(map[string]interface{}); ok {
							logConfig.Options = make(map[string]string)
							for k, v := range opts {
								if str, ok := v.(string); ok {
									logConfig.Options[k] = str
								}
							}
						}
						def.LogConfig = logConfig
					}
				}

				defs = append(defs, def)
			}
		}
	case string:
		// If it's already a JSON string, parse it
		if err := json.Unmarshal([]byte(v), &defs); err != nil {
			return nil, err
		}
	}

	return defs, nil
}

func getMapString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

func getMapInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		if i, ok := v.(int); ok {
			return i
		}
		if f, ok := v.(float64); ok {
			return int(f)
		}
	}
	return 0
}

func getMapBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// MapECSTaskDefinition maps an ECS Task Definition to a TerraformBlock
func MapECSTaskDefinition(taskDef *containers.ECSTaskDefinition) (*mapper.TerraformBlock, error) {
	if taskDef == nil {
		return nil, fmt.Errorf("ecs task definition is nil")
	}

	attributes := make(map[string]mapper.TerraformValue)
	attributes["family"] = strVal(taskDef.Family)

	if len(taskDef.RequiresCompatibilities) > 0 {
		attributes["requires_compatibilities"] = listStrVal(taskDef.RequiresCompatibilities)
	}

	if taskDef.NetworkMode != "" {
		attributes["network_mode"] = strVal(taskDef.NetworkMode)
	}

	if taskDef.CPU != "" {
		attributes["cpu"] = strVal(taskDef.CPU)
	}

	if taskDef.Memory != "" {
		attributes["memory"] = strVal(taskDef.Memory)
	}

	if taskDef.ExecutionRoleARN != "" {
		attributes["execution_role_arn"] = strVal(taskDef.ExecutionRoleARN)
	}

	if taskDef.TaskRoleARN != "" {
		attributes["task_role_arn"] = strVal(taskDef.TaskRoleARN)
	}

	// Build container definitions as JSON-encoded string using jsonencode function
	if len(taskDef.ContainerDefinitions) > 0 {
		containerDefsForTF := buildContainerDefinitionsForTerraform(taskDef.ContainerDefinitions)
		jsonBytes, err := json.MarshalIndent(containerDefsForTF, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal container definitions: %w", err)
		}
		// Use Expr for raw function call representation
		funcCall := fmt.Sprintf("jsonencode(%s)", string(jsonBytes))
		attributes["container_definitions"] = exprVal(funcCall)
	}

	return &mapper.TerraformBlock{
		Kind:       "resource",
		Labels:     []string{"aws_ecs_task_definition", taskDef.Family},
		Attributes: attributes,
	}, nil
}

// buildContainerDefinitionsForTerraform converts container definitions to a format suitable for Terraform
func buildContainerDefinitionsForTerraform(defs []containers.ContainerDefinition) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(defs))

	for _, def := range defs {
		containerDef := map[string]interface{}{
			"name":      def.Name,
			"image":     def.Image,
			"essential": def.Essential,
		}

		if def.CPU > 0 {
			containerDef["cpu"] = def.CPU
		}
		if def.Memory > 0 {
			containerDef["memory"] = def.Memory
		}

		// Port mappings
		if len(def.PortMappings) > 0 {
			pms := make([]map[string]interface{}, 0, len(def.PortMappings))
			for _, pm := range def.PortMappings {
				pmMap := map[string]interface{}{
					"containerPort": pm.ContainerPort,
				}
				if pm.HostPort > 0 {
					pmMap["hostPort"] = pm.HostPort
				}
				if pm.Protocol != "" {
					pmMap["protocol"] = pm.Protocol
				}
				pms = append(pms, pmMap)
			}
			containerDef["portMappings"] = pms
		}

		// Environment variables
		if len(def.Environment) > 0 {
			envs := make([]map[string]interface{}, 0, len(def.Environment))
			for _, env := range def.Environment {
				envs = append(envs, map[string]interface{}{
					"name":  env.Name,
					"value": env.Value,
				})
			}
			containerDef["environment"] = envs
		}

		// Log configuration
		if def.LogConfig != nil {
			logConfig := map[string]interface{}{
				"logDriver": def.LogConfig.LogDriver,
			}
			if len(def.LogConfig.Options) > 0 {
				logConfig["options"] = def.LogConfig.Options
			}
			containerDef["logConfiguration"] = logConfig
		}

		result = append(result, containerDef)
	}

	return result
}
