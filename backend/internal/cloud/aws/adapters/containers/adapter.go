package containers

import (
	"context"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/containers"
	containerservices "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services/containers"
)

// AWSContainerAdapter adapts AWS-specific container service to domain container service
// This implements the Adapter pattern for ECS resources using the Container Service layer
type AWSContainerAdapter struct {
	service containerservices.AWSContainerService
}

// NewAWSContainerAdapter creates a new AWS container adapter with the default service
func NewAWSContainerAdapter() *AWSContainerAdapter {
	return &AWSContainerAdapter{
		service: containerservices.NewContainerService(),
	}
}

// NewAWSContainerAdapterWithService creates a new AWS container adapter with a custom service
// Useful for testing with mock services
func NewAWSContainerAdapterWithService(service containerservices.AWSContainerService) *AWSContainerAdapter {
	return &AWSContainerAdapter{
		service: service,
	}
}

// ============================================================================
// Cluster Operations
// ============================================================================

// CreateCluster creates a new ECS cluster
func (a *AWSContainerAdapter) CreateCluster(ctx context.Context, cluster *containers.ECSCluster) (*containers.ECSCluster, error) {
	_, err := a.service.CreateCluster(ctx, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

// GetCluster retrieves an ECS cluster by name
func (a *AWSContainerAdapter) GetCluster(ctx context.Context, clusterName string) (*containers.ECSCluster, error) {
	output, err := a.service.GetCluster(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	// Convert tags map to []configs.Tag
	tags := make([]configs.Tag, 0, len(output.Tags))
	for k, v := range output.Tags {
		tags = append(tags, configs.Tag{Key: k, Value: v})
	}

	return &containers.ECSCluster{
		Name:                     output.Name,
		ContainerInsightsEnabled: output.ContainerInsightsEnabled,
		Tags:                     tags,
	}, nil
}

// DeleteCluster deletes an ECS cluster
func (a *AWSContainerAdapter) DeleteCluster(ctx context.Context, clusterName string) error {
	return a.service.DeleteCluster(ctx, clusterName)
}

// ListClusters lists all ECS clusters
func (a *AWSContainerAdapter) ListClusters(ctx context.Context) ([]*containers.ECSCluster, error) {
	outputs, err := a.service.ListClusters(ctx, nil)
	if err != nil {
		return nil, err
	}
	clusters := make([]*containers.ECSCluster, len(outputs))
	for i, output := range outputs {
		// Convert tags map to []configs.Tag
		tags := make([]configs.Tag, 0, len(output.Tags))
		for k, v := range output.Tags {
			tags = append(tags, configs.Tag{Key: k, Value: v})
		}
		clusters[i] = &containers.ECSCluster{
			Name:                     output.Name,
			ContainerInsightsEnabled: output.ContainerInsightsEnabled,
			Tags:                     tags,
		}
	}
	return clusters, nil
}

// ============================================================================
// Task Definition Operations
// ============================================================================

// RegisterTaskDefinition registers a new task definition
func (a *AWSContainerAdapter) RegisterTaskDefinition(ctx context.Context, taskDef *containers.ECSTaskDefinition) (*containers.ECSTaskDefinition, error) {
	_, err := a.service.RegisterTaskDefinition(ctx, taskDef)
	if err != nil {
		return nil, err
	}
	return taskDef, nil
}

// GetTaskDefinition retrieves a task definition by family
func (a *AWSContainerAdapter) GetTaskDefinition(ctx context.Context, family string) (*containers.ECSTaskDefinition, error) {
	output, err := a.service.GetTaskDefinition(ctx, family)
	if err != nil {
		return nil, err
	}
	return &containers.ECSTaskDefinition{
		Family:                  output.Family,
		ContainerDefinitions:    output.ContainerDefinitions,
		RequiresCompatibilities: output.RequiresCompatibilities,
		NetworkMode:             output.NetworkMode,
		CPU:                     output.CPU,
		Memory:                  output.Memory,
		ExecutionRoleARN:        output.ExecutionRoleARN,
		TaskRoleARN:             output.TaskRoleARN,
	}, nil
}

// DeregisterTaskDefinition deregisters a task definition
func (a *AWSContainerAdapter) DeregisterTaskDefinition(ctx context.Context, taskDefARN string) error {
	return a.service.DeregisterTaskDefinition(ctx, taskDefARN)
}

// ListTaskDefinitions lists task definitions
func (a *AWSContainerAdapter) ListTaskDefinitions(ctx context.Context, family string) ([]*containers.ECSTaskDefinition, error) {
	outputs, err := a.service.ListTaskDefinitions(ctx, family)
	if err != nil {
		return nil, err
	}
	taskDefs := make([]*containers.ECSTaskDefinition, len(outputs))
	for i, output := range outputs {
		taskDefs[i] = &containers.ECSTaskDefinition{
			Family:                  output.Family,
			RequiresCompatibilities: output.RequiresCompatibilities,
			NetworkMode:             output.NetworkMode,
			CPU:                     output.CPU,
			Memory:                  output.Memory,
		}
	}
	return taskDefs, nil
}

// ============================================================================
// Service Operations
// ============================================================================

// CreateService creates a new ECS service
func (a *AWSContainerAdapter) CreateService(ctx context.Context, service *containers.ECSService) (*containers.ECSService, error) {
	_, err := a.service.CreateService(ctx, service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// GetService retrieves an ECS service
func (a *AWSContainerAdapter) GetService(ctx context.Context, clusterName, serviceName string) (*containers.ECSService, error) {
	output, err := a.service.GetService(ctx, clusterName, serviceName)
	if err != nil {
		return nil, err
	}

	// Convert tags map to []configs.Tag
	tags := make([]configs.Tag, 0, len(output.Tags))
	for k, v := range output.Tags {
		tags = append(tags, configs.Tag{Key: k, Value: v})
	}

	ecsService := &containers.ECSService{
		Name:              output.Name,
		ClusterID:         clusterName,
		TaskDefinitionARN: output.TaskDefinitionARN,
		DesiredCount:      output.DesiredCount,
		LaunchType:        output.LaunchType,
		PlatformVersion:   output.PlatformVersion,
		Tags:              tags,
	}

	if output.NetworkConfiguration != nil {
		ecsService.NetworkConfiguration = &containers.NetworkConfiguration{
			Subnets:        output.NetworkConfiguration.SubnetIDs,
			SecurityGroups: output.NetworkConfiguration.SecurityGroupIDs,
			AssignPublicIP: output.NetworkConfiguration.AssignPublicIP,
		}
	}

	return ecsService, nil
}

// UpdateService updates an ECS service
func (a *AWSContainerAdapter) UpdateService(ctx context.Context, service *containers.ECSService) (*containers.ECSService, error) {
	_, err := a.service.UpdateService(ctx, service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// DeleteService deletes an ECS service
func (a *AWSContainerAdapter) DeleteService(ctx context.Context, clusterName, serviceName string) error {
	return a.service.DeleteService(ctx, clusterName, serviceName)
}

// ListServices lists services in a cluster
func (a *AWSContainerAdapter) ListServices(ctx context.Context, clusterName string) ([]*containers.ECSService, error) {
	outputs, err := a.service.ListServices(ctx, clusterName)
	if err != nil {
		return nil, err
	}
	services := make([]*containers.ECSService, len(outputs))
	for i, output := range outputs {
		// Convert tags map to []configs.Tag
		tags := make([]configs.Tag, 0, len(output.Tags))
		for k, v := range output.Tags {
			tags = append(tags, configs.Tag{Key: k, Value: v})
		}
		services[i] = &containers.ECSService{
			Name:              output.Name,
			ClusterID:         clusterName,
			TaskDefinitionARN: output.TaskDefinitionARN,
			DesiredCount:      output.DesiredCount,
			LaunchType:        output.LaunchType,
			Tags:              tags,
		}
	}
	return services, nil
}

// ============================================================================
// Capacity Provider Operations
// ============================================================================

// CreateCapacityProvider creates a new capacity provider
func (a *AWSContainerAdapter) CreateCapacityProvider(ctx context.Context, provider *containers.ECSCapacityProvider) (*containers.ECSCapacityProvider, error) {
	_, err := a.service.CreateCapacityProvider(ctx, provider)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

// GetCapacityProvider retrieves a capacity provider
func (a *AWSContainerAdapter) GetCapacityProvider(ctx context.Context, providerName string) (*containers.ECSCapacityProvider, error) {
	output, err := a.service.GetCapacityProvider(ctx, providerName)
	if err != nil {
		return nil, err
	}

	provider := &containers.ECSCapacityProvider{
		Name: output.Name,
	}

	if output.AutoScalingGroupARN != "" {
		provider.AutoScalingGroupProvider = &containers.AutoScalingGroupProvider{
			AutoScalingGroupARN:          output.AutoScalingGroupARN,
			ManagedTerminationProtection: output.ManagedTerminationProtection,
		}

		if output.ManagedScaling != nil {
			provider.AutoScalingGroupProvider.ManagedScaling = &containers.ManagedScaling{
				Status:                 output.ManagedScaling.Status,
				TargetCapacity:         output.ManagedScaling.TargetCapacity,
				MinimumScalingStepSize: output.ManagedScaling.MinimumScalingStepSize,
				MaximumScalingStepSize: output.ManagedScaling.MaximumScalingStepSize,
				InstanceWarmupPeriod:   output.ManagedScaling.InstanceWarmupPeriod,
			}
		}
	}

	return provider, nil
}

// DeleteCapacityProvider deletes a capacity provider
func (a *AWSContainerAdapter) DeleteCapacityProvider(ctx context.Context, providerName string) error {
	return a.service.DeleteCapacityProvider(ctx, providerName)
}

// ============================================================================
// Cluster Capacity Providers Operations
// ============================================================================

// PutClusterCapacityProviders associates capacity providers with a cluster
func (a *AWSContainerAdapter) PutClusterCapacityProviders(ctx context.Context, config *containers.ECSClusterCapacityProviders) (*containers.ECSClusterCapacityProviders, error) {
	_, err := a.service.PutClusterCapacityProviders(ctx, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// GetClusterCapacityProviders gets capacity provider configuration for a cluster
func (a *AWSContainerAdapter) GetClusterCapacityProviders(ctx context.Context, clusterName string) (*containers.ECSClusterCapacityProviders, error) {
	output, err := a.service.GetClusterCapacityProviders(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	config := &containers.ECSClusterCapacityProviders{
		ClusterName:       clusterName,
		CapacityProviders: output.CapacityProviders,
	}

	if len(output.DefaultCapacityProviderStrategy) > 0 {
		config.DefaultCapacityProviderStrategies = make([]containers.CapacityProviderStrategy, len(output.DefaultCapacityProviderStrategy))
		for i, strat := range output.DefaultCapacityProviderStrategy {
			config.DefaultCapacityProviderStrategies[i] = containers.CapacityProviderStrategy{
				CapacityProvider: strat.CapacityProvider,
				Weight:           strat.Weight,
				Base:             strat.Base,
			}
		}
	}

	return config, nil
}
