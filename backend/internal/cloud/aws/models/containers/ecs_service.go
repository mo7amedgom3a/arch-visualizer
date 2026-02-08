package containers

import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"

// ECSService represents an AWS ECS Service
type ECSService struct {
	Name                     string                         `json:"name"`
	ClusterID                string                         `json:"cluster_id"`                    // Reference to ECS Cluster
	TaskDefinitionARN        string                         `json:"task_definition_arn,omitempty"` // Reference to Task Definition
	DesiredCount             int                            `json:"desired_count,omitempty"`
	LaunchType               string                         `json:"launch_type,omitempty"`         // FARGATE or EC2
	PlatformVersion          string                         `json:"platform_version,omitempty"`    // LATEST, 1.4.0, etc.
	SchedulingStrategy       string                         `json:"scheduling_strategy,omitempty"` // REPLICA or DAEMON
	NetworkConfiguration     *NetworkConfiguration          `json:"network_configuration,omitempty"`
	LoadBalancers            []LoadBalancerConfig           `json:"load_balancers,omitempty"`
	CapacityProviderStrategy []CapacityProviderStrategyItem `json:"capacity_provider_strategy,omitempty"`
	DeploymentConfiguration  *DeploymentConfiguration       `json:"deployment_configuration,omitempty"`
	ServiceRegistries        []ServiceRegistry              `json:"service_registries,omitempty"`
	EnableExecuteCommand     bool                           `json:"enable_execute_command,omitempty"`
	ForceNewDeployment       bool                           `json:"force_new_deployment,omitempty"`
	Tags                     []configs.Tag                  `json:"tags,omitempty"`
}

// NetworkConfiguration for awsvpc network mode
type NetworkConfiguration struct {
	Subnets        []string `json:"subnets"`
	SecurityGroups []string `json:"security_groups,omitempty"`
	AssignPublicIP bool     `json:"assign_public_ip,omitempty"`
}

// LoadBalancerConfig connects service to ALB/NLB
type LoadBalancerConfig struct {
	TargetGroupARN string `json:"target_group_arn"`
	ContainerName  string `json:"container_name"`
	ContainerPort  int    `json:"container_port"`
}

// CapacityProviderStrategyItem defines capacity provider allocation
type CapacityProviderStrategyItem struct {
	CapacityProvider string `json:"capacity_provider"`
	Weight           int    `json:"weight,omitempty"`
	Base             int    `json:"base,omitempty"`
}

// DeploymentConfiguration defines deployment settings
type DeploymentConfiguration struct {
	MaximumPercent        int             `json:"maximum_percent,omitempty"`
	MinimumHealthyPercent int             `json:"minimum_healthy_percent,omitempty"`
	CircuitBreaker        *CircuitBreaker `json:"circuit_breaker,omitempty"`
}

// CircuitBreaker for deployment rollback
type CircuitBreaker struct {
	Enable   bool `json:"enable"`
	Rollback bool `json:"rollback"`
}

// ServiceRegistry for service discovery
type ServiceRegistry struct {
	RegistryARN   string `json:"registry_arn"`
	Port          int    `json:"port,omitempty"`
	ContainerName string `json:"container_name,omitempty"`
	ContainerPort int    `json:"container_port,omitempty"`
}
