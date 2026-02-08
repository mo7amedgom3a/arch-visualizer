package containers

import (
	"context"
	"fmt"

	awserrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/errors"
	awscontainers "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/containers"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/services"
	domainerrors "github.com/mo7amedgom3a/arch-visualizer/backend/internal/domain/errors"
)

// ContainerService implements AWSContainerService with deterministic virtual operations
type ContainerService struct{}

// NewContainerService creates a new container service implementation
func NewContainerService() *ContainerService {
	return &ContainerService{}
}

// ============================================================================
// Cluster Operations
// ============================================================================

func (s *ContainerService) CreateCluster(ctx context.Context, cluster *awscontainers.ECSCluster) (*ClusterOutput, error) {
	if cluster == nil {
		return nil, domainerrors.New(awserrors.CodeECSClusterCreationFailed, domainerrors.KindValidation, "cluster is nil").
			WithOp("ContainerService.CreateCluster")
	}

	region := "us-east-1"
	clusterARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:cluster/%s", region, cluster.Name)

	// Convert tags to map for output
	tags := make(map[string]string)
	for _, tag := range cluster.Tags {
		tags[tag.Key] = tag.Value
	}

	return &ClusterOutput{
		ARN:                               clusterARN,
		Name:                              cluster.Name,
		Status:                            "ACTIVE",
		Region:                            region,
		RegisteredContainerInstancesCount: 0,
		RunningTasksCount:                 0,
		PendingTasksCount:                 0,
		ActiveServicesCount:               0,
		CapacityProviders:                 []string{},
		ContainerInsightsEnabled:          cluster.ContainerInsightsEnabled,
		CreatedAt:                         services.GetFixedTimestamp().Format("2006-01-02T15:04:05Z"),
		Tags:                              tags,
	}, nil
}

func (s *ContainerService) GetCluster(ctx context.Context, clusterName string) (*ClusterOutput, error) {
	region := "us-east-1"
	clusterARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:cluster/%s", region, clusterName)

	return &ClusterOutput{
		ARN:                               clusterARN,
		Name:                              clusterName,
		Status:                            "ACTIVE",
		Region:                            region,
		RegisteredContainerInstancesCount: 0,
		RunningTasksCount:                 2,
		PendingTasksCount:                 0,
		ActiveServicesCount:               1,
		CapacityProviders:                 []string{"FARGATE", "FARGATE_SPOT"},
		ContainerInsightsEnabled:          true,
		CreatedAt:                         services.GetFixedTimestamp().Format("2006-01-02T15:04:05Z"),
		Tags:                              map[string]string{"Name": clusterName},
	}, nil
}

func (s *ContainerService) UpdateCluster(ctx context.Context, cluster *awscontainers.ECSCluster) (*ClusterOutput, error) {
	return s.CreateCluster(ctx, cluster)
}

func (s *ContainerService) DeleteCluster(ctx context.Context, clusterName string) error {
	return nil
}

func (s *ContainerService) ListClusters(ctx context.Context, filters map[string][]string) ([]*ClusterOutput, error) {
	return []*ClusterOutput{
		{
			ARN:                               "arn:aws:ecs:us-east-1:123456789012:cluster/test-cluster",
			Name:                              "test-cluster",
			Status:                            "ACTIVE",
			Region:                            "us-east-1",
			RegisteredContainerInstancesCount: 0,
			RunningTasksCount:                 2,
			PendingTasksCount:                 0,
			ActiveServicesCount:               1,
			CapacityProviders:                 []string{"FARGATE"},
			ContainerInsightsEnabled:          true,
			CreatedAt:                         services.GetFixedTimestamp().Format("2006-01-02T15:04:05Z"),
			Tags:                              map[string]string{"Name": "test-cluster"},
		},
	}, nil
}

// ============================================================================
// Task Definition Operations
// ============================================================================

func (s *ContainerService) RegisterTaskDefinition(ctx context.Context, taskDef *awscontainers.ECSTaskDefinition) (*TaskDefinitionOutput, error) {
	if taskDef == nil {
		return nil, domainerrors.New(awserrors.CodeECSTaskDefinitionCreationFailed, domainerrors.KindValidation, "task definition is nil").
			WithOp("ContainerService.RegisterTaskDefinition")
	}

	region := "us-east-1"
	revision := 1
	taskDefARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:task-definition/%s:%d", region, taskDef.Family, revision)

	return &TaskDefinitionOutput{
		ARN:                     taskDefARN,
		Family:                  taskDef.Family,
		Revision:                revision,
		Status:                  "ACTIVE",
		ContainerDefinitions:    taskDef.ContainerDefinitions,
		RequiresCompatibilities: taskDef.RequiresCompatibilities,
		NetworkMode:             taskDef.NetworkMode,
		CPU:                     taskDef.CPU,
		Memory:                  taskDef.Memory,
		ExecutionRoleARN:        taskDef.ExecutionRoleARN,
		TaskRoleARN:             taskDef.TaskRoleARN,
		RegisteredAt:            services.GetFixedTimestamp().Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *ContainerService) GetTaskDefinition(ctx context.Context, family string) (*TaskDefinitionOutput, error) {
	region := "us-east-1"
	revision := 1
	taskDefARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:task-definition/%s:%d", region, family, revision)

	return &TaskDefinitionOutput{
		ARN:      taskDefARN,
		Family:   family,
		Revision: revision,
		Status:   "ACTIVE",
		ContainerDefinitions: []awscontainers.ContainerDefinition{
			{
				Name:      "web",
				Image:     "nginx:latest",
				Essential: true,
				CPU:       256,
				Memory:    512,
				PortMappings: []awscontainers.PortMapping{
					{ContainerPort: 80, HostPort: 80, Protocol: "tcp"},
				},
			},
		},
		RequiresCompatibilities: []string{"FARGATE"},
		NetworkMode:             "awsvpc",
		CPU:                     "256",
		Memory:                  "512",
		RegisteredAt:            services.GetFixedTimestamp().Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *ContainerService) DeregisterTaskDefinition(ctx context.Context, taskDefARN string) error {
	return nil
}

func (s *ContainerService) ListTaskDefinitions(ctx context.Context, family string) ([]*TaskDefinitionOutput, error) {
	region := "us-east-1"
	taskDefARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:task-definition/%s:1", region, family)

	return []*TaskDefinitionOutput{
		{
			ARN:                     taskDefARN,
			Family:                  family,
			Revision:                1,
			Status:                  "ACTIVE",
			RequiresCompatibilities: []string{"FARGATE"},
			NetworkMode:             "awsvpc",
			CPU:                     "256",
			Memory:                  "512",
			RegisteredAt:            services.GetFixedTimestamp().Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

// ============================================================================
// Service Operations
// ============================================================================

func (s *ContainerService) CreateService(ctx context.Context, service *awscontainers.ECSService) (*ServiceOutput, error) {
	if service == nil {
		return nil, domainerrors.New(awserrors.CodeECSServiceCreationFailed, domainerrors.KindValidation, "service is nil").
			WithOp("ContainerService.CreateService")
	}

	region := "us-east-1"
	clusterARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:cluster/%s", region, service.ClusterID)
	serviceARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:service/%s/%s", region, service.ClusterID, service.Name)

	// Convert tags to map for output
	tags := make(map[string]string)
	for _, tag := range service.Tags {
		tags[tag.Key] = tag.Value
	}

	output := &ServiceOutput{
		ARN:               serviceARN,
		Name:              service.Name,
		ClusterARN:        clusterARN,
		TaskDefinitionARN: service.TaskDefinitionARN,
		Status:            "ACTIVE",
		DesiredCount:      service.DesiredCount,
		RunningCount:      service.DesiredCount,
		PendingCount:      0,
		LaunchType:        service.LaunchType,
		PlatformVersion:   service.PlatformVersion,
		CreatedAt:         services.GetFixedTimestamp().Format("2006-01-02T15:04:05Z"),
		Tags:              tags,
	}

	// Add network configuration if present (using correct field names: Subnets, SecurityGroups)
	if service.NetworkConfiguration != nil {
		output.NetworkConfiguration = &NetworkConfigurationOutput{
			SubnetIDs:        service.NetworkConfiguration.Subnets,
			SecurityGroupIDs: service.NetworkConfiguration.SecurityGroups,
			AssignPublicIP:   service.NetworkConfiguration.AssignPublicIP,
		}
	}

	// Add load balancers if present
	if len(service.LoadBalancers) > 0 {
		output.LoadBalancers = make([]LoadBalancerOutput, len(service.LoadBalancers))
		for i, lb := range service.LoadBalancers {
			output.LoadBalancers[i] = LoadBalancerOutput{
				TargetGroupARN: lb.TargetGroupARN,
				ContainerName:  lb.ContainerName,
				ContainerPort:  lb.ContainerPort,
			}
		}
	}

	return output, nil
}

func (s *ContainerService) GetService(ctx context.Context, clusterName, serviceName string) (*ServiceOutput, error) {
	region := "us-east-1"
	clusterARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:cluster/%s", region, clusterName)
	serviceARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:service/%s/%s", region, clusterName, serviceName)

	return &ServiceOutput{
		ARN:               serviceARN,
		Name:              serviceName,
		ClusterARN:        clusterARN,
		TaskDefinitionARN: fmt.Sprintf("arn:aws:ecs:%s:123456789012:task-definition/web-app:1", region),
		Status:            "ACTIVE",
		DesiredCount:      2,
		RunningCount:      2,
		PendingCount:      0,
		LaunchType:        "FARGATE",
		PlatformVersion:   "LATEST",
		NetworkConfiguration: &NetworkConfigurationOutput{
			SubnetIDs:        []string{"subnet-123", "subnet-456"},
			SecurityGroupIDs: []string{"sg-123"},
			AssignPublicIP:   true,
		},
		CreatedAt: services.GetFixedTimestamp().Format("2006-01-02T15:04:05Z"),
		Tags:      map[string]string{"Name": serviceName},
	}, nil
}

func (s *ContainerService) UpdateService(ctx context.Context, service *awscontainers.ECSService) (*ServiceOutput, error) {
	return s.CreateService(ctx, service)
}

func (s *ContainerService) DeleteService(ctx context.Context, clusterName, serviceName string) error {
	return nil
}

func (s *ContainerService) ListServices(ctx context.Context, clusterName string) ([]*ServiceOutput, error) {
	region := "us-east-1"
	clusterARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:cluster/%s", region, clusterName)
	serviceARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:service/%s/test-service", region, clusterName)

	return []*ServiceOutput{
		{
			ARN:               serviceARN,
			Name:              "test-service",
			ClusterARN:        clusterARN,
			TaskDefinitionARN: fmt.Sprintf("arn:aws:ecs:%s:123456789012:task-definition/web-app:1", region),
			Status:            "ACTIVE",
			DesiredCount:      2,
			RunningCount:      2,
			PendingCount:      0,
			LaunchType:        "FARGATE",
			CreatedAt:         services.GetFixedTimestamp().Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

// ============================================================================
// Capacity Provider Operations
// ============================================================================

func (s *ContainerService) CreateCapacityProvider(ctx context.Context, provider *awscontainers.ECSCapacityProvider) (*CapacityProviderOutput, error) {
	if provider == nil {
		return nil, domainerrors.New(awserrors.CodeECSCapacityProviderCreationFailed, domainerrors.KindValidation, "capacity provider is nil").
			WithOp("ContainerService.CreateCapacityProvider")
	}

	region := "us-east-1"
	providerARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:capacity-provider/%s", region, provider.Name)

	output := &CapacityProviderOutput{
		ARN:    providerARN,
		Name:   provider.Name,
		Status: "ACTIVE",
	}

	// Extract from AutoScalingGroupProvider if present
	if provider.AutoScalingGroupProvider != nil {
		output.AutoScalingGroupARN = provider.AutoScalingGroupProvider.AutoScalingGroupARN
		output.ManagedTerminationProtection = provider.AutoScalingGroupProvider.ManagedTerminationProtection

		if provider.AutoScalingGroupProvider.ManagedScaling != nil {
			output.ManagedScaling = &ManagedScalingOutput{
				Status:                 provider.AutoScalingGroupProvider.ManagedScaling.Status,
				TargetCapacity:         provider.AutoScalingGroupProvider.ManagedScaling.TargetCapacity,
				MinimumScalingStepSize: provider.AutoScalingGroupProvider.ManagedScaling.MinimumScalingStepSize,
				MaximumScalingStepSize: provider.AutoScalingGroupProvider.ManagedScaling.MaximumScalingStepSize,
				InstanceWarmupPeriod:   provider.AutoScalingGroupProvider.ManagedScaling.InstanceWarmupPeriod,
			}
		}
	}

	return output, nil
}

func (s *ContainerService) GetCapacityProvider(ctx context.Context, providerName string) (*CapacityProviderOutput, error) {
	region := "us-east-1"
	providerARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:capacity-provider/%s", region, providerName)

	return &CapacityProviderOutput{
		ARN:                          providerARN,
		Name:                         providerName,
		Status:                       "ACTIVE",
		AutoScalingGroupARN:          "arn:aws:autoscaling:us-east-1:123456789012:autoScalingGroup:uuid:autoScalingGroupName/test-asg",
		ManagedTerminationProtection: "ENABLED",
		ManagedScaling: &ManagedScalingOutput{
			Status:                 "ENABLED",
			TargetCapacity:         100,
			MinimumScalingStepSize: 1,
			MaximumScalingStepSize: 10000,
			InstanceWarmupPeriod:   300,
		},
	}, nil
}

func (s *ContainerService) DeleteCapacityProvider(ctx context.Context, providerName string) error {
	return nil
}

func (s *ContainerService) ListCapacityProviders(ctx context.Context) ([]*CapacityProviderOutput, error) {
	return []*CapacityProviderOutput{
		{
			ARN:                          "arn:aws:ecs:us-east-1:123456789012:capacity-provider/FARGATE",
			Name:                         "FARGATE",
			Status:                       "ACTIVE",
			ManagedTerminationProtection: "DISABLED",
		},
		{
			ARN:                          "arn:aws:ecs:us-east-1:123456789012:capacity-provider/FARGATE_SPOT",
			Name:                         "FARGATE_SPOT",
			Status:                       "ACTIVE",
			ManagedTerminationProtection: "DISABLED",
		},
	}, nil
}

// ============================================================================
// Cluster Capacity Providers Operations
// ============================================================================

func (s *ContainerService) PutClusterCapacityProviders(ctx context.Context, config *awscontainers.ECSClusterCapacityProviders) (*ClusterCapacityProvidersOutput, error) {
	if config == nil {
		return nil, domainerrors.New(awserrors.CodeECSClusterCreationFailed, domainerrors.KindValidation, "cluster capacity providers config is nil").
			WithOp("ContainerService.PutClusterCapacityProviders")
	}

	region := "us-east-1"
	clusterARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:cluster/%s", region, config.ClusterName)

	output := &ClusterCapacityProvidersOutput{
		ClusterARN:        clusterARN,
		CapacityProviders: config.CapacityProviders,
	}

	// Use correct field name: DefaultCapacityProviderStrategies
	if len(config.DefaultCapacityProviderStrategies) > 0 {
		output.DefaultCapacityProviderStrategy = make([]DefaultCapacityProviderStrategyOutput, len(config.DefaultCapacityProviderStrategies))
		for i, strat := range config.DefaultCapacityProviderStrategies {
			output.DefaultCapacityProviderStrategy[i] = DefaultCapacityProviderStrategyOutput{
				CapacityProvider: strat.CapacityProvider,
				Weight:           strat.Weight,
				Base:             strat.Base,
			}
		}
	}

	return output, nil
}

func (s *ContainerService) GetClusterCapacityProviders(ctx context.Context, clusterName string) (*ClusterCapacityProvidersOutput, error) {
	region := "us-east-1"
	clusterARN := fmt.Sprintf("arn:aws:ecs:%s:123456789012:cluster/%s", region, clusterName)

	return &ClusterCapacityProvidersOutput{
		ClusterARN:        clusterARN,
		CapacityProviders: []string{"FARGATE", "FARGATE_SPOT"},
		DefaultCapacityProviderStrategy: []DefaultCapacityProviderStrategyOutput{
			{CapacityProvider: "FARGATE", Weight: 1, Base: 1},
		},
	}, nil
}
