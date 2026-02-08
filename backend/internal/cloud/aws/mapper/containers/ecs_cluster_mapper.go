package containers

import (
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/containers"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/resource"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/iac/terraform/mapper"
)

// ClusterFromResource converts a generic domain resource to an ECS Cluster model
func ClusterFromResource(res *resource.Resource) (*containers.ECSCluster, error) {
	if res.Type.Name != "ECSCluster" {
		return nil, fmt.Errorf("invalid resource type for ECS Cluster mapper: %s", res.Type.Name)
	}

	cluster := &containers.ECSCluster{
		Name: res.Name,
	}

	// Helper to safely extract values from metadata
	getString := func(key string) string {
		if val, ok := res.Metadata[key]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
		return ""
	}

	getBool := func(key string) bool {
		if val, ok := res.Metadata[key]; ok {
			if b, ok := val.(bool); ok {
				return b
			}
		}
		return false
	}

	cluster.ContainerInsightsEnabled = getBool("container_insights_enabled")
	cluster.ExecuteCommandEnabled = getBool("execute_command_enabled")
	cluster.KMSKeyID = getString("kms_key_id")
	cluster.LogGroup = getString("log_group")

	return cluster, nil
}

// MapECSCluster maps an ECS Cluster to a TerraformBlock
func MapECSCluster(cluster *containers.ECSCluster) (*mapper.TerraformBlock, error) {
	if cluster == nil {
		return nil, fmt.Errorf("ecs cluster is nil")
	}

	attributes := make(map[string]mapper.TerraformValue)
	nestedBlocks := make(map[string][]mapper.NestedBlock)

	attributes["name"] = strVal(cluster.Name)

	// Container Insights setting
	if cluster.ContainerInsightsEnabled {
		settingBlock := mapper.NestedBlock{
			Attributes: map[string]mapper.TerraformValue{
				"name":  strVal("containerInsights"),
				"value": strVal("enabled"),
			},
		}
		nestedBlocks["setting"] = []mapper.NestedBlock{settingBlock}
	}

	// Execute command configuration
	if cluster.ExecuteCommandEnabled {
		execCmdAttrs := make(map[string]mapper.TerraformValue)
		if cluster.KMSKeyID != "" {
			execCmdAttrs["kms_key_id"] = strVal(cluster.KMSKeyID)
		}
		if cluster.LogGroup != "" {
			execCmdAttrs["logging"] = strVal("OVERRIDE")
		}
		// Note: log_configuration would need nested blocks within conf block
		// For simplicity, we'll handle it as a flat structure

		configBlock := mapper.NestedBlock{
			Attributes: execCmdAttrs,
		}
		nestedBlocks["configuration"] = []mapper.NestedBlock{configBlock}
	}

	return &mapper.TerraformBlock{
		Kind:         "resource",
		Labels:       []string{"aws_ecs_cluster", cluster.Name},
		Attributes:   attributes,
		NestedBlocks: nestedBlocks,
	}, nil
}

// Helper functions for creating TerraformValue

func strVal(s string) mapper.TerraformValue {
	return mapper.TerraformValue{String: &s}
}

func intVal(i int) mapper.TerraformValue {
	f := float64(i)
	return mapper.TerraformValue{Number: &f}
}

func boolVal(b bool) mapper.TerraformValue {
	return mapper.TerraformValue{Bool: &b}
}

func listStrVal(list []string) mapper.TerraformValue {
	vals := make([]mapper.TerraformValue, len(list))
	for i, s := range list {
		str := s
		vals[i] = mapper.TerraformValue{String: &str}
	}
	return mapper.TerraformValue{List: vals}
}

func exprVal(expr string) mapper.TerraformValue {
	e := mapper.TerraformExpr(expr)
	return mapper.TerraformValue{Expr: &e}
}
