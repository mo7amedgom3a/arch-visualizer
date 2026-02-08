package containers

import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/cloud/aws/configs"

// ECSTaskDefinition represents an AWS ECS Task Definition
type ECSTaskDefinition struct {
	Family                  string                `json:"family"`
	RequiresCompatibilities []string              `json:"requires_compatibilities,omitempty"` // FARGATE, EC2
	NetworkMode             string                `json:"network_mode,omitempty"`             // awsvpc, bridge, host, none
	CPU                     string                `json:"cpu,omitempty"`                      // 256, 512, 1024, 2048, 4096
	Memory                  string                `json:"memory,omitempty"`                   // 512, 1024, 2048, etc.
	ExecutionRoleARN        string                `json:"execution_role_arn,omitempty"`       // Role for ECS agent
	TaskRoleARN             string                `json:"task_role_arn,omitempty"`            // Role for container app
	ContainerDefinitions    []ContainerDefinition `json:"container_definitions"`
	Volumes                 []Volume              `json:"volumes,omitempty"`
	Tags                    []configs.Tag         `json:"tags,omitempty"`
}

// ContainerDefinition defines a container within a task definition
type ContainerDefinition struct {
	Name              string            `json:"name"`
	Image             string            `json:"image"`
	Essential         bool              `json:"essential"`
	CPU               int               `json:"cpu,omitempty"`
	Memory            int               `json:"memory,omitempty"`
	MemoryReservation int               `json:"memory_reservation,omitempty"`
	PortMappings      []PortMapping     `json:"port_mappings,omitempty"`
	Environment       []KeyValuePair    `json:"environment,omitempty"`
	Secrets           []SecretPair      `json:"secrets,omitempty"`
	LogConfig         *LogConfiguration `json:"log_configuration,omitempty"`
	Command           []string          `json:"command,omitempty"`
	EntryPoint        []string          `json:"entry_point,omitempty"`
	WorkingDir        string            `json:"working_directory,omitempty"`
	HealthCheck       *HealthCheck      `json:"health_check,omitempty"`
}

// PortMapping maps a container port to a host port
type PortMapping struct {
	ContainerPort int    `json:"container_port"`
	HostPort      int    `json:"host_port,omitempty"`
	Protocol      string `json:"protocol,omitempty"` // tcp, udp
}

// KeyValuePair represents an environment variable
type KeyValuePair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// SecretPair represents a secret reference
type SecretPair struct {
	Name      string `json:"name"`
	ValueFrom string `json:"value_from"` // ARN of secret or parameter
}

// LogConfiguration defines container logging
type LogConfiguration struct {
	LogDriver string            `json:"log_driver"` // awslogs, splunk, fluentd, etc.
	Options   map[string]string `json:"options,omitempty"`
}

// HealthCheck defines container health check
type HealthCheck struct {
	Command     []string `json:"command"`
	Interval    int      `json:"interval,omitempty"` // Seconds
	Timeout     int      `json:"timeout,omitempty"`  // Seconds
	Retries     int      `json:"retries,omitempty"`
	StartPeriod int      `json:"start_period,omitempty"` // Seconds
}

// Volume defines a task volume
type Volume struct {
	Name                   string                  `json:"name"`
	HostPath               string                  `json:"host_path,omitempty"`
	EFSVolumeConfiguration *EFSVolumeConfiguration `json:"efs_volume_configuration,omitempty"`
}

// EFSVolumeConfiguration defines EFS mount configuration
type EFSVolumeConfiguration struct {
	FileSystemID          string `json:"file_system_id"`
	RootDirectory         string `json:"root_directory,omitempty"`
	TransitEncryption     string `json:"transit_encryption,omitempty"` // ENABLED, DISABLED
	TransitEncryptionPort int    `json:"transit_encryption_port,omitempty"`
	AccessPointID         string `json:"access_point_id,omitempty"`
}
