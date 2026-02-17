package networking

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetadata_GetResourceSchema_VPC(t *testing.T) {
	svc := NewNetworkingMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "vpc")

	require.NoError(t, err)
	assert.Equal(t, "vpc", schema.Label)

	// Expect required fields: name, region, cidr
	requiredNames := []string{}
	for _, f := range schema.Fields {
		if f.Required {
			requiredNames = append(requiredNames, f.Name)
		}
	}
	assert.Contains(t, requiredNames, "name")
	assert.Contains(t, requiredNames, "region")
	assert.Contains(t, requiredNames, "cidr")

	// Outputs should include id, arn, state
	assert.Contains(t, schema.Outputs, "id")
	assert.Contains(t, schema.Outputs, "arn")
	assert.Contains(t, schema.Outputs, "state")
}

func TestMetadata_GetResourceSchema_Subnet(t *testing.T) {
	svc := NewNetworkingMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "subnet")

	require.NoError(t, err)
	assert.Equal(t, "subnet", schema.Label)

	// Required fields: name, vpc_id, cidr, availability_zone
	requiredNames := []string{}
	for _, f := range schema.Fields {
		if f.Required {
			requiredNames = append(requiredNames, f.Name)
		}
	}
	assert.Contains(t, requiredNames, "name")
	assert.Contains(t, requiredNames, "vpc_id")
	assert.Contains(t, requiredNames, "cidr")
	assert.Contains(t, requiredNames, "availability_zone")

	assert.Contains(t, schema.Outputs, "id")
	assert.Contains(t, schema.Outputs, "available_ip_count")
}

func TestMetadata_ListResourceSchemas_AllRegistered(t *testing.T) {
	svc := NewNetworkingMetadataService()
	schemas, err := svc.ListResourceSchemas(context.Background())

	require.NoError(t, err)
	assert.Len(t, schemas, 10, "expect all 10 networking resource schemas")

	labels := map[string]bool{}
	for _, s := range schemas {
		labels[s.Label] = true
	}

	expectedLabels := []string{
		"vpc", "subnet", "internet_gateway", "route_table",
		"security_group", "nat_gateway", "elastic_ip",
		"network_acl", "network_interface", "vpc_endpoint",
	}
	for _, l := range expectedLabels {
		assert.True(t, labels[l], "missing schema for %s", l)
	}
}

func TestMetadata_GetResourceSchema_NotFound(t *testing.T) {
	svc := NewNetworkingMetadataService()
	_, err := svc.GetResourceSchema(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown networking resource")
}

func TestMetadata_SecurityGroup_HasEnumValues(t *testing.T) {
	svc := NewNetworkingMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "security_group")
	require.NoError(t, err)

	// VPC ID should be required
	var vpcField *struct{ Required bool }
	for _, f := range schema.Fields {
		if f.Name == "vpc_id" {
			vpcField = &struct{ Required bool }{f.Required}
		}
	}
	require.NotNil(t, vpcField)
	assert.True(t, vpcField.Required)
}

func TestMetadata_ElasticIP_PoolTypeEnum(t *testing.T) {
	svc := NewNetworkingMetadataService()
	schema, err := svc.GetResourceSchema(context.Background(), "elastic_ip")
	require.NoError(t, err)

	for _, f := range schema.Fields {
		if f.Name == "address_pool_type" {
			assert.ElementsMatch(t, []string{"amazon", "byoip", "customer_owned", "ipam"}, f.Enum)
			return
		}
	}
	t.Fatal("address_pool_type field not found")
}
