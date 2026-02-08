package containers

// ECSClusterCapacityProviders represents the attachment of capacity providers to an ECS Cluster
type ECSClusterCapacityProviders struct {
	ClusterName                       string                     `json:"cluster_name"`
	CapacityProviders                 []string                   `json:"capacity_providers,omitempty"` // FARGATE, FARGATE_SPOT, custom
	DefaultCapacityProviderStrategies []CapacityProviderStrategy `json:"default_capacity_provider_strategy,omitempty"`
}

// CapacityProviderStrategy defines the default strategy for task placement
type CapacityProviderStrategy struct {
	CapacityProvider string `json:"capacity_provider"`
	Weight           int    `json:"weight,omitempty"` // 0-1000
	Base             int    `json:"base,omitempty"`   // Min tasks before splitting by weight
}
