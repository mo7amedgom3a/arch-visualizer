package containers

import (
	"context"
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/containers"
	"github.com/stretchr/testify/assert"
)

func TestNewAWSContainerAdapter(t *testing.T) {
	adapter := NewAWSContainerAdapter()
	assert.NotNil(t, adapter)
}

func TestNewContainerAdapter(t *testing.T) {
	adapter := NewContainerAdapter()
	assert.NotNil(t, adapter)
}

func TestNewContainerAdapterWithConfig(t *testing.T) {
	config := ContainerAdapterConfig{
		Region: "us-east-1",
	}
	adapter := NewContainerAdapterWithConfig(config)
	assert.NotNil(t, adapter)
}

func TestClusterOperations(t *testing.T) {
	adapter := NewAWSContainerAdapter()
	ctx := context.Background()

	t.Run("CreateCluster", func(t *testing.T) {
		cluster := &containers.ECSCluster{
			Name:                     "test-cluster",
			ContainerInsightsEnabled: true,
		}
		result, err := adapter.CreateCluster(ctx, cluster)
		assert.NoError(t, err)
		assert.Equal(t, "test-cluster", result.Name)
	})

	t.Run("GetCluster", func(t *testing.T) {
		result, err := adapter.GetCluster(ctx, "test-cluster")
		assert.NoError(t, err)
		assert.Equal(t, "test-cluster", result.Name)
	})

	t.Run("ListClusters", func(t *testing.T) {
		result, err := adapter.ListClusters(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("DeleteCluster", func(t *testing.T) {
		err := adapter.DeleteCluster(ctx, "test-cluster")
		assert.NoError(t, err)
	})
}

func TestTaskDefinitionOperations(t *testing.T) {
	adapter := NewAWSContainerAdapter()
	ctx := context.Background()

	t.Run("RegisterTaskDefinition", func(t *testing.T) {
		taskDef := &containers.ECSTaskDefinition{
			Family:      "test-task",
			NetworkMode: "awsvpc",
		}
		result, err := adapter.RegisterTaskDefinition(ctx, taskDef)
		assert.NoError(t, err)
		assert.Equal(t, "test-task", result.Family)
	})

	t.Run("GetTaskDefinition", func(t *testing.T) {
		result, err := adapter.GetTaskDefinition(ctx, "test-task")
		assert.NoError(t, err)
		assert.Equal(t, "test-task", result.Family)
	})

	t.Run("ListTaskDefinitions", func(t *testing.T) {
		result, err := adapter.ListTaskDefinitions(ctx, "test-task")
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("DeregisterTaskDefinition", func(t *testing.T) {
		err := adapter.DeregisterTaskDefinition(ctx, "arn:aws:ecs:...")
		assert.NoError(t, err)
	})
}

func TestServiceOperations(t *testing.T) {
	adapter := NewAWSContainerAdapter()
	ctx := context.Background()

	t.Run("CreateService", func(t *testing.T) {
		service := &containers.ECSService{
			Name:         "test-service",
			ClusterID:    "test-cluster",
			DesiredCount: 3,
		}
		result, err := adapter.CreateService(ctx, service)
		assert.NoError(t, err)
		assert.Equal(t, "test-service", result.Name)
	})

	t.Run("GetService", func(t *testing.T) {
		result, err := adapter.GetService(ctx, "test-cluster", "test-service")
		assert.NoError(t, err)
		assert.Equal(t, "test-service", result.Name)
		assert.Equal(t, "test-cluster", result.ClusterID)
	})

	t.Run("UpdateService", func(t *testing.T) {
		service := &containers.ECSService{
			Name:         "test-service",
			ClusterID:    "test-cluster",
			DesiredCount: 5,
		}
		result, err := adapter.UpdateService(ctx, service)
		assert.NoError(t, err)
		assert.Equal(t, 5, result.DesiredCount)
	})

	t.Run("ListServices", func(t *testing.T) {
		result, err := adapter.ListServices(ctx, "test-cluster")
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("DeleteService", func(t *testing.T) {
		err := adapter.DeleteService(ctx, "test-cluster", "test-service")
		assert.NoError(t, err)
	})
}

func TestCapacityProviderOperations(t *testing.T) {
	adapter := NewAWSContainerAdapter()
	ctx := context.Background()

	t.Run("CreateCapacityProvider", func(t *testing.T) {
		provider := &containers.ECSCapacityProvider{
			Name: "test-provider",
		}
		result, err := adapter.CreateCapacityProvider(ctx, provider)
		assert.NoError(t, err)
		assert.Equal(t, "test-provider", result.Name)
	})

	t.Run("GetCapacityProvider", func(t *testing.T) {
		result, err := adapter.GetCapacityProvider(ctx, "test-provider")
		assert.NoError(t, err)
		assert.Equal(t, "test-provider", result.Name)
	})

	t.Run("DeleteCapacityProvider", func(t *testing.T) {
		err := adapter.DeleteCapacityProvider(ctx, "test-provider")
		assert.NoError(t, err)
	})
}

func TestClusterCapacityProvidersOperations(t *testing.T) {
	adapter := NewAWSContainerAdapter()
	ctx := context.Background()

	t.Run("PutClusterCapacityProviders", func(t *testing.T) {
		config := &containers.ECSClusterCapacityProviders{
			ClusterName:       "test-cluster",
			CapacityProviders: []string{"FARGATE", "FARGATE_SPOT"},
		}
		result, err := adapter.PutClusterCapacityProviders(ctx, config)
		assert.NoError(t, err)
		assert.Equal(t, "test-cluster", result.ClusterName)
	})

	t.Run("GetClusterCapacityProviders", func(t *testing.T) {
		result, err := adapter.GetClusterCapacityProviders(ctx, "test-cluster")
		assert.NoError(t, err)
		assert.Equal(t, "test-cluster", result.ClusterName)
	})
}
