package containers

import (
	"testing"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/containers"
	"github.com/stretchr/testify/assert"
)

func TestMapECSCluster(t *testing.T) {
	tests := []struct {
		name     string
		input    *containers.ECSCluster
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name: "Basic ECS Cluster",
			input: &containers.ECSCluster{
				Name: "my-cluster",
			},
			expected: map[string]interface{}{
				"name": "my-cluster",
			},
			wantErr: false,
		},
		{
			name: "ECS Cluster with Container Insights",
			input: &containers.ECSCluster{
				Name:                     "my-cluster",
				ContainerInsightsEnabled: true,
			},
			expected: map[string]interface{}{
				"name": "my-cluster",
			},
			wantErr: false,
		},
		{
			name: "ECS Cluster with Execute Command",
			input: &containers.ECSCluster{
				Name:                  "my-cluster",
				ExecuteCommandEnabled: true,
				KMSKeyID:              "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012",
				LogGroup:              "/ecs/execute-command-logs",
			},
			expected: map[string]interface{}{
				"name": "my-cluster",
			},
			wantErr: false,
		},
		{
			name:    "Nil Cluster",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapECSCluster(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, "resource", got.Kind)
			assert.Equal(t, "aws_ecs_cluster", got.Labels[0])
			assert.Equal(t, tt.input.Name, got.Labels[1])

			// Check name attribute
			assert.NotNil(t, got.Attributes["name"].String)
			assert.Equal(t, tt.expected["name"], *got.Attributes["name"].String)

			// Check Container Insights setting
			if tt.input.ContainerInsightsEnabled {
				assert.Contains(t, got.NestedBlocks, "setting")
				assert.Len(t, got.NestedBlocks["setting"], 1)
			}

			// Check Execute Command configuration
			if tt.input.ExecuteCommandEnabled {
				assert.Contains(t, got.NestedBlocks, "configuration")
			}
		})
	}
}

func TestMapECSService(t *testing.T) {
	tests := []struct {
		name     string
		input    *containers.ECSService
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name: "Basic ECS Service",
			input: &containers.ECSService{
				Name:         "my-service",
				ClusterID:    "my-cluster",
				DesiredCount: 2,
				LaunchType:   "FARGATE",
			},
			expected: map[string]interface{}{
				"name":          "my-service",
				"cluster":       "my-cluster",
				"desired_count": 2,
				"launch_type":   "FARGATE",
			},
			wantErr: false,
		},
		{
			name: "ECS Service with Network Configuration",
			input: &containers.ECSService{
				Name:         "my-service",
				ClusterID:    "my-cluster",
				DesiredCount: 1,
				LaunchType:   "FARGATE",
				NetworkConfiguration: &containers.NetworkConfiguration{
					Subnets:        []string{"subnet-1", "subnet-2"},
					SecurityGroups: []string{"sg-1"},
					AssignPublicIP: true,
				},
			},
			wantErr: false,
		},
		{
			name: "ECS Service with Load Balancer",
			input: &containers.ECSService{
				Name:         "my-service",
				ClusterID:    "my-cluster",
				DesiredCount: 2,
				LaunchType:   "FARGATE",
				LoadBalancers: []containers.LoadBalancerConfig{
					{
						TargetGroupARN: "arn:aws:elasticloadbalancing:us-east-1:123456789012:targetgroup/my-tg/1234567890123456",
						ContainerName:  "app",
						ContainerPort:  8080,
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Nil Service",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapECSService(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, "resource", got.Kind)
			assert.Equal(t, "aws_ecs_service", got.Labels[0])
			assert.Equal(t, tt.input.Name, got.Labels[1])

			// Check network configuration
			if tt.input.NetworkConfiguration != nil {
				assert.Contains(t, got.NestedBlocks, "network_configuration")
			}

			// Check load balancer
			if len(tt.input.LoadBalancers) > 0 {
				assert.Contains(t, got.NestedBlocks, "load_balancer")
			}
		})
	}
}

func TestMapECSTaskDefinition(t *testing.T) {
	tests := []struct {
		name     string
		input    *containers.ECSTaskDefinition
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name: "Basic Task Definition",
			input: &containers.ECSTaskDefinition{
				Family:                  "my-task",
				RequiresCompatibilities: []string{"FARGATE"},
				NetworkMode:             "awsvpc",
				CPU:                     "256",
				Memory:                  "512",
			},
			expected: map[string]interface{}{
				"family":       "my-task",
				"network_mode": "awsvpc",
				"cpu":          "256",
				"memory":       "512",
			},
			wantErr: false,
		},
		{
			name: "Task Definition with Container",
			input: &containers.ECSTaskDefinition{
				Family:                  "my-task",
				RequiresCompatibilities: []string{"FARGATE"},
				NetworkMode:             "awsvpc",
				CPU:                     "256",
				Memory:                  "512",
				ContainerDefinitions: []containers.ContainerDefinition{
					{
						Name:      "app",
						Image:     "nginx:latest",
						Essential: true,
						CPU:       256,
						Memory:    512,
						PortMappings: []containers.PortMapping{
							{ContainerPort: 80, Protocol: "tcp"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Nil Task Definition",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapECSTaskDefinition(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, "resource", got.Kind)
			assert.Equal(t, "aws_ecs_task_definition", got.Labels[0])
			assert.Equal(t, tt.input.Family, got.Labels[1])

			// Check container definitions
			if len(tt.input.ContainerDefinitions) > 0 {
				assert.Contains(t, got.Attributes, "container_definitions")
			}
		})
	}
}

func TestMapECSCapacityProvider(t *testing.T) {
	tests := []struct {
		name    string
		input   *containers.ECSCapacityProvider
		wantErr bool
	}{
		{
			name: "Basic Capacity Provider",
			input: &containers.ECSCapacityProvider{
				Name: "my-capacity-provider",
				AutoScalingGroupProvider: &containers.AutoScalingGroupProvider{
					AutoScalingGroupARN:          "arn:aws:autoscaling:us-east-1:123456789012:autoScalingGroup:12345678-1234-1234-1234-123456789012:autoScalingGroupName/my-asg",
					ManagedTerminationProtection: "DISABLED",
				},
			},
			wantErr: false,
		},
		{
			name: "Capacity Provider with Managed Scaling",
			input: &containers.ECSCapacityProvider{
				Name: "my-capacity-provider",
				AutoScalingGroupProvider: &containers.AutoScalingGroupProvider{
					AutoScalingGroupARN: "arn:aws:autoscaling:us-east-1:123456789012:autoScalingGroup:12345678-1234-1234-1234-123456789012:autoScalingGroupName/my-asg",
					ManagedScaling: &containers.ManagedScaling{
						Status:         "ENABLED",
						TargetCapacity: 100,
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Nil Capacity Provider",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapECSCapacityProvider(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, "resource", got.Kind)
			assert.Equal(t, "aws_ecs_capacity_provider", got.Labels[0])
			assert.Equal(t, tt.input.Name, got.Labels[1])

			// Check ASG provider block
			if tt.input.AutoScalingGroupProvider != nil {
				assert.Contains(t, got.NestedBlocks, "auto_scaling_group_provider")
			}

			// Check managed scaling block
			if tt.input.AutoScalingGroupProvider != nil && tt.input.AutoScalingGroupProvider.ManagedScaling != nil {
				asgBlock := got.NestedBlocks["auto_scaling_group_provider"][0]
				assert.Contains(t, asgBlock.NestedBlocks, "managed_scaling")
			}
		})
	}
}

func TestMapECSClusterCapacityProviders(t *testing.T) {
	tests := []struct {
		name    string
		input   *containers.ECSClusterCapacityProviders
		wantErr bool
	}{
		{
			name: "Basic Cluster Capacity Providers",
			input: &containers.ECSClusterCapacityProviders{
				ClusterName:       "my-cluster",
				CapacityProviders: []string{"FARGATE", "FARGATE_SPOT"},
			},
			wantErr: false,
		},
		{
			name: "Cluster Capacity Providers with Strategy",
			input: &containers.ECSClusterCapacityProviders{
				ClusterName:       "my-cluster",
				CapacityProviders: []string{"FARGATE", "FARGATE_SPOT"},
				DefaultCapacityProviderStrategies: []containers.CapacityProviderStrategy{
					{CapacityProvider: "FARGATE", Weight: 1, Base: 1},
					{CapacityProvider: "FARGATE_SPOT", Weight: 4},
				},
			},
			wantErr: false,
		},
		{
			name:    "Nil Cluster Capacity Providers",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapECSClusterCapacityProviders(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, "resource", got.Kind)
			assert.Equal(t, "aws_ecs_cluster_capacity_providers", got.Labels[0])
			assert.Equal(t, tt.input.ClusterName, got.Labels[1])

			// Check capacity providers list
			if len(tt.input.CapacityProviders) > 0 {
				assert.Contains(t, got.Attributes, "capacity_providers")
			}

			// Check strategy blocks
			if len(tt.input.DefaultCapacityProviderStrategies) > 0 {
				assert.Contains(t, got.NestedBlocks, "default_capacity_provider_strategy")
				assert.Len(t, got.NestedBlocks["default_capacity_provider_strategy"], len(tt.input.DefaultCapacityProviderStrategies))
			}
		})
	}
}
