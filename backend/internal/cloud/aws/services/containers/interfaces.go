package containers

import (
	"context"

	awscontainers "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/models/containers"
)

// AWSContainerService defines AWS-specific ECS container operations
// This implements cloud provider-specific logic while maintaining domain compatibility
type AWSContainerService interface {
	// Cluster operations
	CreateCluster(ctx context.Context, cluster *awscontainers.ECSCluster) (*ClusterOutput, error)
	GetCluster(ctx context.Context, clusterName string) (*ClusterOutput, error)
	UpdateCluster(ctx context.Context, cluster *awscontainers.ECSCluster) (*ClusterOutput, error)
	DeleteCluster(ctx context.Context, clusterName string) error
	ListClusters(ctx context.Context, filters map[string][]string) ([]*ClusterOutput, error)

	// Task Definition operations
	RegisterTaskDefinition(ctx context.Context, taskDef *awscontainers.ECSTaskDefinition) (*TaskDefinitionOutput, error)
	GetTaskDefinition(ctx context.Context, family string) (*TaskDefinitionOutput, error)
	DeregisterTaskDefinition(ctx context.Context, taskDefARN string) error
	ListTaskDefinitions(ctx context.Context, family string) ([]*TaskDefinitionOutput, error)

	// Service operations
	CreateService(ctx context.Context, service *awscontainers.ECSService) (*ServiceOutput, error)
	GetService(ctx context.Context, clusterName, serviceName string) (*ServiceOutput, error)
	UpdateService(ctx context.Context, service *awscontainers.ECSService) (*ServiceOutput, error)
	DeleteService(ctx context.Context, clusterName, serviceName string) error
	ListServices(ctx context.Context, clusterName string) ([]*ServiceOutput, error)

	// Capacity Provider operations
	CreateCapacityProvider(ctx context.Context, provider *awscontainers.ECSCapacityProvider) (*CapacityProviderOutput, error)
	GetCapacityProvider(ctx context.Context, providerName string) (*CapacityProviderOutput, error)
	DeleteCapacityProvider(ctx context.Context, providerName string) error
	ListCapacityProviders(ctx context.Context) ([]*CapacityProviderOutput, error)

	// Cluster Capacity Providers operations
	PutClusterCapacityProviders(ctx context.Context, config *awscontainers.ECSClusterCapacityProviders) (*ClusterCapacityProvidersOutput, error)
	GetClusterCapacityProviders(ctx context.Context, clusterName string) (*ClusterCapacityProvidersOutput, error)
}

// ClusterOutput represents the output from ECS Cluster operations
type ClusterOutput struct {
	ARN                               string            `json:"arn"`
	Name                              string            `json:"name"`
	Status                            string            `json:"status"`
	Region                            string            `json:"region"`
	RegisteredContainerInstancesCount int               `json:"registered_container_instances_count"`
	RunningTasksCount                 int               `json:"running_tasks_count"`
	PendingTasksCount                 int               `json:"pending_tasks_count"`
	ActiveServicesCount               int               `json:"active_services_count"`
	CapacityProviders                 []string          `json:"capacity_providers,omitempty"`
	ContainerInsightsEnabled          bool              `json:"container_insights_enabled"`
	CreatedAt                         string            `json:"created_at"`
	Tags                              map[string]string `json:"tags,omitempty"`
}

// TaskDefinitionOutput represents the output from ECS Task Definition operations
type TaskDefinitionOutput struct {
	ARN                     string                              `json:"arn"`
	Family                  string                              `json:"family"`
	Revision                int                                 `json:"revision"`
	Status                  string                              `json:"status"`
	ContainerDefinitions    []awscontainers.ContainerDefinition `json:"container_definitions"`
	RequiresCompatibilities []string                            `json:"requires_compatibilities"`
	NetworkMode             string                              `json:"network_mode"`
	CPU                     string                              `json:"cpu"`
	Memory                  string                              `json:"memory"`
	ExecutionRoleARN        string                              `json:"execution_role_arn,omitempty"`
	TaskRoleARN             string                              `json:"task_role_arn,omitempty"`
	RegisteredAt            string                              `json:"registered_at"`
}

// ServiceOutput represents the output from ECS Service operations
type ServiceOutput struct {
	ARN                    string                      `json:"arn"`
	Name                   string                      `json:"name"`
	ClusterARN             string                      `json:"cluster_arn"`
	TaskDefinitionARN      string                      `json:"task_definition_arn"`
	Status                 string                      `json:"status"`
	DesiredCount           int                         `json:"desired_count"`
	RunningCount           int                         `json:"running_count"`
	PendingCount           int                         `json:"pending_count"`
	LaunchType             string                      `json:"launch_type"`
	PlatformVersion        string                      `json:"platform_version,omitempty"`
	HealthCheckGracePeriod int                         `json:"health_check_grace_period,omitempty"`
	NetworkConfiguration   *NetworkConfigurationOutput `json:"network_configuration,omitempty"`
	LoadBalancers          []LoadBalancerOutput        `json:"load_balancers,omitempty"`
	CreatedAt              string                      `json:"created_at"`
	Tags                   map[string]string           `json:"tags,omitempty"`
}

// NetworkConfigurationOutput represents network configuration output
type NetworkConfigurationOutput struct {
	SubnetIDs        []string `json:"subnet_ids"`
	SecurityGroupIDs []string `json:"security_group_ids"`
	AssignPublicIP   bool     `json:"assign_public_ip"`
}

// LoadBalancerOutput represents load balancer configuration output
type LoadBalancerOutput struct {
	TargetGroupARN string `json:"target_group_arn"`
	ContainerName  string `json:"container_name"`
	ContainerPort  int    `json:"container_port"`
}

// CapacityProviderOutput represents the output from Capacity Provider operations
type CapacityProviderOutput struct {
	ARN                          string                `json:"arn"`
	Name                         string                `json:"name"`
	Status                       string                `json:"status"`
	AutoScalingGroupARN          string                `json:"auto_scaling_group_arn"`
	ManagedScaling               *ManagedScalingOutput `json:"managed_scaling,omitempty"`
	ManagedTerminationProtection string                `json:"managed_termination_protection"`
}

// ManagedScalingOutput represents managed scaling configuration output
type ManagedScalingOutput struct {
	Status                 string `json:"status"`
	TargetCapacity         int    `json:"target_capacity"`
	MinimumScalingStepSize int    `json:"minimum_scaling_step_size"`
	MaximumScalingStepSize int    `json:"maximum_scaling_step_size"`
	InstanceWarmupPeriod   int    `json:"instance_warmup_period"`
}

// ClusterCapacityProvidersOutput represents cluster capacity providers output
type ClusterCapacityProvidersOutput struct {
	ClusterARN                      string                                  `json:"cluster_arn"`
	CapacityProviders               []string                                `json:"capacity_providers"`
	DefaultCapacityProviderStrategy []DefaultCapacityProviderStrategyOutput `json:"default_capacity_provider_strategy,omitempty"`
}

// DefaultCapacityProviderStrategyOutput represents default capacity provider strategy output
type DefaultCapacityProviderStrategyOutput struct {
	CapacityProvider string `json:"capacity_provider"`
	Weight           int    `json:"weight"`
	Base             int    `json:"base"`
}
