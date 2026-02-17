package compute

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeMetadata_GetResourceSchema_EC2Instance(t *testing.T) {
	svc := NewComputeMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "ec2_instance")

	require.NoError(t, err)
	assert.Equal(t, "ec2_instance", schema.Label)

	requiredNames := []string{}
	for _, f := range schema.Fields {
		if f.Required {
			requiredNames = append(requiredNames, f.Name)
		}
	}
	assert.Contains(t, requiredNames, "name")
	assert.Contains(t, requiredNames, "ami")
	assert.Contains(t, requiredNames, "instance_type")
	assert.Contains(t, requiredNames, "subnet_id")
	assert.Contains(t, requiredNames, "vpc_security_group_ids")

	assert.Contains(t, schema.Outputs, "id")
	assert.Contains(t, schema.Outputs, "arn")
	assert.Contains(t, schema.Outputs, "public_ip")
}

func TestComputeMetadata_GetResourceSchema_LoadBalancer(t *testing.T) {
	svc := NewComputeMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "load_balancer")

	require.NoError(t, err)
	assert.Equal(t, "load_balancer", schema.Label)

	// Validate enum on load_balancer_type
	for _, f := range schema.Fields {
		if f.Name == "load_balancer_type" {
			assert.ElementsMatch(t, []string{"application", "network"}, f.Enum)
			return
		}
	}
	t.Fatal("load_balancer_type field not found")
}

func TestComputeMetadata_GetResourceSchema_LambdaFunction(t *testing.T) {
	svc := NewComputeMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "lambda_function")

	require.NoError(t, err)
	assert.Equal(t, "lambda_function", schema.Label)

	requiredNames := []string{}
	for _, f := range schema.Fields {
		if f.Required {
			requiredNames = append(requiredNames, f.Name)
		}
	}
	assert.Contains(t, requiredNames, "function_name")
	assert.Contains(t, requiredNames, "role_arn")

	assert.Contains(t, schema.Outputs, "arn")
	assert.Contains(t, schema.Outputs, "invoke_arn")
}

func TestComputeMetadata_ListResourceSchemas_AllRegistered(t *testing.T) {
	svc := NewComputeMetadataService()
	schemas, err := svc.ListResourceSchemas(context.Background())

	require.NoError(t, err)
	assert.Len(t, schemas, 7, "expect all 7 compute resource schemas")

	labels := map[string]bool{}
	for _, s := range schemas {
		labels[s.Label] = true
	}

	expectedLabels := []string{
		"ec2_instance", "launch_template", "load_balancer", "target_group",
		"listener", "auto_scaling_group", "lambda_function",
	}
	for _, l := range expectedLabels {
		assert.True(t, labels[l], "missing schema for %s", l)
	}
}

func TestComputeMetadata_GetResourceSchema_NotFound(t *testing.T) {
	svc := NewComputeMetadataService()
	_, err := svc.GetResourceSchema(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown compute resource")
}

func TestComputeMetadata_AutoScalingGroup_HealthCheckEnum(t *testing.T) {
	svc := NewComputeMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "auto_scaling_group")
	require.NoError(t, err)

	for _, f := range schema.Fields {
		if f.Name == "health_check_type" {
			assert.ElementsMatch(t, []string{"EC2", "ELB"}, f.Enum)
			return
		}
	}
	t.Fatal("health_check_type field not found")
}
